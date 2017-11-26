package kmedia

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime/debug"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/volatiletech/null.v6"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

type FileMappings struct {
	Sha1              string
	KMediaID          int
	KMediaContainerID int
	MdbExists         bool
	MdbID             int64
	MdbUnitID         int64
	MdbKMediaID       int
}

func MapUnits() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting Kmedia unit mappings")

	log.Info("Setting up connection to MDB")
	mdb, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()
	boil.SetDB(mdb)
	//boil.DebugMode = true

	log.Info("Setting up connection to Kmedia")
	kmdb, err = sql.Open("postgres", viper.GetString("kmedia.url"))
	utils.Must(err)
	utils.Must(kmdb.Ping())
	defer kmdb.Close()

	log.Info("Initializing static data from MDB")
	utils.Must(api.InitTypeRegistries(mdb))

	log.Info("Initializing static data from Kmedia")
	kmediaLessonCT, err = kmodels.ContentTypes(kmdb, qm.Where("name = ?", "Lesson")).One()
	utils.Must(err)
	serverUrls, err = initServers(kmdb)
	utils.Must(err)

	log.Info("Loading all containers")
	containers, err := kmodels.Containers(kmdb).All()
	utils.Must(err)
	log.Infof("Got %d containers", len(containers))

	stats = NewImportStatistics()

	log.Info("Setting up workers")
	jobs := make(chan *kmodels.Container, 100)
	results := make(chan []*FileMappings, 100)
	done := make(chan []*FileMappings)

	var workersWG sync.WaitGroup
	for w := 1; w <= 5; w++ {
		workersWG.Add(1)
		go fileMappingsWorker(jobs, results, &workersWG)
	}

	go mappingsAggregator(results, done)

	log.Info("Queueing work")
	for _, container := range containers {
		jobs <- container
	}

	log.Info("Closing jobs channel")
	close(jobs)

	log.Info("Waiting for workers to finish")
	workersWG.Wait()

	log.Info("Closing results channel")
	close(results)

	// wait for aggregation results
	mappings := <-done
	log.Infof("Got %d file mappings", len(mappings))

	log.Info("Dumping file mappings")
	err = dumpFileMappings(mappings)
	utils.Must(err)

	log.Info("Fixing unit mappings")
	err = fixUnitMappings(mappings)
	utils.Must(err)

	log.Info("Fixing file mappings")
	err = fixFileMappings(mappings)
	utils.Must(err)

	stats.dump()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func fileMappingsWorker(jobs <-chan *kmodels.Container, results chan []*FileMappings, wg *sync.WaitGroup) {
	for container := range jobs {
		stats.ContainersProcessed.Inc(1)

		err := container.L.LoadFileAssets(kmdb, true, container)
		if err != nil {
			log.Error(err)
			debug.PrintStack()
			continue
		}

		sha1s := make([][]byte, 0)
		bySha1 := make(map[string]*FileMappings, 0)
		mappings := make([]*FileMappings, 0)

		for _, fileAsset := range container.R.FileAssets {
			//log.Infof("Processing file_asset %d", fileAsset.ID)
			stats.FileAssetsProcessed.Inc(1)

			if !fileAsset.Sha1.Valid {
				stats.FileAssetsMissingSHA1.Inc(1)
				continue
			}

			sha1 := fileAsset.Sha1.String

			hexSha1, err := hex.DecodeString(sha1)
			if err != nil {
				log.Error(err)
				debug.PrintStack()
				continue
			}

			sha1s = append(sha1s, hexSha1)
			fm := &FileMappings{
				Sha1:              fileAsset.Sha1.String,
				KMediaID:          fileAsset.ID,
				KMediaContainerID: container.ID,
				MdbExists:         false,
			}
			bySha1[sha1] = fm
			mappings = append(mappings, fm)
		}

		if len(sha1s) > 0 {
			// fetch mdb files
			files, err := models.Files(mdb,
				qm.WhereIn("sha1 in ?", utils.ConvertArgsBytes(sha1s)...)).
				All()
			if err != nil {
				log.Error(err)
				debug.PrintStack()
				continue
			}

			for _, file := range files {
				sha1 := hex.EncodeToString(file.Sha1.Bytes)
				fm := bySha1[sha1]
				fm.MdbExists = true
				fm.MdbID = file.ID
				if file.ContentUnitID.Valid {
					fm.MdbUnitID = file.ContentUnitID.Int64
				}
				if file.Properties.Valid {
					var props map[string]interface{}
					err = json.Unmarshal(file.Properties.JSON, &props)
					if err != nil {
						log.Error(err)
						debug.PrintStack()
						continue
					}
					if kmid, ok := props["kmedia_id"]; ok {
						fm.MdbKMediaID = int(kmid.(float64))
					}
				}
			}
		}

		results <- mappings
	}
	wg.Done()
}

func mappingsAggregator(results <-chan []*FileMappings, done chan []*FileMappings) {
	mappings := make([]*FileMappings, 0)
	for result := range results {
		mappings = append(mappings, result...)
	}
	done <- mappings
}

func dumpFileMappings(mappings []*FileMappings) error {
	f, err := ioutil.TempFile("/tmp", "mdb_kmedia_mappings_report")
	if err != nil {
		return errors.Wrap(err, "Create temp file")
	}
	defer f.Close()

	log.Infof("Report file: %s", f.Name())

	for _, fm := range mappings {
		fmt.Fprintf(f, "%s,%d,%d,%d,%d,%t\n",
			fm.Sha1, fm.KMediaID, fm.KMediaContainerID, fm.MdbID, fm.MdbUnitID, fm.MdbExists)
	}
	return nil
}

func fixUnitMappings(mappings []*FileMappings) error {
	unitMap := make(map[int]map[int64]bool, stats.ContainersProcessed.Get())
	for _, fm := range mappings {
		cm, ok := unitMap[fm.KMediaContainerID]
		if !ok {
			cm = make(map[int64]bool)
		}
		if fm.MdbUnitID > 25316 {
			//if fm.MdbUnitID > 0 {
			cm[fm.MdbUnitID] = true
		}
		unitMap[fm.KMediaContainerID] = cm
	}

	bySize := make(map[int][]int)
	for k, v := range unitMap {
		s, ok := bySize[len(v)]
		if !ok {
			s = make([]int, 0)
		}
		s = append(s, k)
		bySize[len(v)] = s
	}

	tx, err := mdb.Begin()
	utils.Must(err)
	unexpectedCT := map[int64]bool{
		api.CONTENT_TYPE_REGISTRY.ByName[api.CT_KITEI_MAKOR].ID: true,
	}
	for k, v := range bySize {
		fmt.Printf("size %d\n", k)
		for _, kcid := range v {
			for mcuid := range unitMap[kcid] {
				//fmt.Printf("%d,%d\n", kcid, mcuid)
				u, err := models.FindContentUnit(mdb, mcuid)
				if err != nil {
					return errors.Wrapf(err, "Load CU from DB %d", mcuid)
				}
				if _, ok := unexpectedCT[u.TypeID]; ok {
					fmt.Printf("!!!%s!!! [%d]\n", api.CONTENT_TYPE_REGISTRY.ByID[u.TypeID].Name, mcuid)
				} else {
					// check if mapping changed before updating
					if u.Properties.Valid {
						var props map[string]interface{}
						err := u.Properties.Unmarshal(&props)
						if err != nil {
							return errors.Wrapf(err, "json.Unmarshal unit properties [%d]", u.ID)
						}
						if kmid, ok := props["kmedia_id"]; ok && int(kmid.(float64)) == kcid {
							continue
						}
					}

					log.Infof("Mapping kmedia container %d to mdb unit %d", kcid, u.ID)
					err = api.UpdateContentUnitProperties(tx, u, map[string]interface{}{"kmedia_id": kcid})
					if err != nil {
						utils.Must(tx.Rollback())
						return errors.Wrapf(err, "Update CU properties [%d]", u.ID)
					}
					stats.ContentUnitsUpdated.Inc(1)
				}
			}
		}
	}
	utils.Must(tx.Commit())

	return nil
}

func fixFileMappings(mappings []*FileMappings) error {

	// load CU to container mappings
	rows, err := queries.Raw(mdb,
		"select id, (properties->>'kmedia_id')::int from content_units where properties ? 'kmedia_id';").
		Query()
	if err != nil {
		return errors.Wrap(err, "Load CU to container map")
	}
	defer rows.Close()

	cuMap := make(map[int]int64)
	for rows.Next() {
		var cuid int64
		var kmid int
		err = rows.Scan(&cuid, &kmid)
		if err != nil {
			return errors.Wrap(err, "rows.Scan")
		}
		cuMap[kmid] = cuid
	}
	if err = rows.Err(); err != nil {
		return errors.Wrap(err, "rows.Err")
	}

	// fix file mappings:
	// 1. "kmedia_id" property
	// 2. CU association based on kmedia containers mappings
	for _, fm := range mappings {
		if !fm.MdbExists {
			continue
		}

		if fm.KMediaID == fm.MdbKMediaID && fm.MdbUnitID != 0 {
			continue
		}

		f, err := models.FindFile(mdb, fm.MdbID)
		if err != nil {
			return errors.Wrapf(err, "Load file %d", fm.MdbID)
		}

		if fm.KMediaID != fm.MdbKMediaID {
			log.Infof("Updating kmedia_id property for %d to %d", f.ID, fm.KMediaID)
			err = api.UpdateFileProperties(mdb, f, map[string]interface{}{"kmedia_id": fm.KMediaID})
			if err != nil {
				return errors.Wrapf(err, "Update file properties %d", f.ID)
			}
		}

		if !f.ContentUnitID.Valid {
			if cuid, ok := cuMap[fm.KMediaContainerID]; ok {
				log.Infof("Associating file %d to CU %d [container %d]", f.ID, cuid, fm.KMediaContainerID)
				f.ContentUnitID = null.Int64From(cuid)
				err = f.Update(mdb, "content_unit_id")
				if err != nil {
					return errors.Wrapf(err, "Associate CU [file %d]", f.ID)
				}
			}
		}
		stats.FilesUpdated.Inc(1)
	}

	return nil
}
