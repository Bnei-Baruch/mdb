package kmedia

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/pkg/errors"
	"io/ioutil"
)

type FileMappings struct {
	Sha1              string
	KMediaID          int
	KMediaContainerID int
	MdbExists         bool
	MdbID             int64
	MdbUnitID         int64
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

	log.Info("Loading all virtual_lessons")
	vls, err := kmodels.VirtualLessons(kmdb).All()
	utils.Must(err)
	log.Infof("Got %d lessons", len(vls))

	stats = NewImportStatistics()

	log.Info("Setting up workers")
	jobs := make(chan *kmodels.VirtualLesson, 100)
	results := make(chan []*FileMappings, 100)
	done := make(chan []*FileMappings)

	var workersWG sync.WaitGroup
	for w := 1; w <= 5; w++ {
		workersWG.Add(1)
		go fileMappingsWorker(jobs, results, &workersWG)
	}

	go mappingsAggregator(results, done)

	log.Info("Queueing work")
	for _, vl := range vls {
		jobs <- vl
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
	err = analyzeFileMappings(mappings)
	utils.Must(err)

	stats.dump()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func fileMappingsWorker(jobs <-chan *kmodels.VirtualLesson, results chan []*FileMappings, wg *sync.WaitGroup) {
	for vl := range jobs {
		//log.Infof("Processing virtual_lesson %d", vl.ID)
		stats.LessonsProcessed.Inc(1)

		// Validate virtual lesson data
		containers, err := getValidContainers(kmdb, vl)
		if err != nil {
			log.Error(err)
			debug.PrintStack()
			continue
		}
		if len(containers) == 0 {
			//log.Warnf("Invalid lesson [%d]", vl.ID)
			stats.InvalidLessons.Inc(1)
			continue
		}

		stats.ValidLessons.Inc(1)

		sha1s := make([][]byte, 0)
		bySha1 := make(map[string]*FileMappings, 0)
		mappings := make([]*FileMappings, 0)
		for _, container := range containers {
			//log.Infof("Processing container %d", container.ID)
			stats.ContainersProcessed.Inc(1)

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

func analyzeFileMappings(mappings []*FileMappings) error {
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

	unitMap := make(map[int]map[int64]bool, stats.ContainersProcessed.Get())
	for _, fm := range mappings {
		cm, ok := unitMap[fm.KMediaContainerID]
		if !ok {
			cm = make(map[int64]bool)
		}
		if fm.MdbUnitID > 25822 {
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
	expectedCT := map[int64]bool{
		api.CONTENT_TYPE_REGISTRY.ByName[api.CT_LESSON_PART].ID: true,
		api.CONTENT_TYPE_REGISTRY.ByName[api.CT_FULL_LESSON].ID: true,
	}
	for k, v := range bySize {
		fmt.Printf("size %d\n", k)
		for _, kcid := range v {
			for mcuid := range unitMap[kcid] {
				fmt.Printf("%d,%d\n", kcid, mcuid)
				u, err := models.FindContentUnit(mdb, mcuid)
				if err != nil {
					return errors.Wrap(err, "Load CU from DB")
				}
				if _, ok := expectedCT[u.TypeID]; !ok {
					fmt.Printf("!!!%s!!! [%d]\n", api.CONTENT_TYPE_REGISTRY.ByID[u.TypeID].Name, mcuid)
				} else {
					err = api.UpdateContentUnitProperties(tx, u, map[string]interface{}{"kmedia_id": kcid})
					if err != nil {
						utils.Must(tx.Rollback())
						return errors.Wrapf(err, "Update CU properties [%d]", u.ID)
					}
				}
			}
		}
	}
	utils.Must(tx.Commit())

	return nil
}
