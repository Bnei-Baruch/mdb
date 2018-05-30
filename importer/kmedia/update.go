package kmedia

import (
	"database/sql"
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

func UpdateI18ns() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	stats = NewImportStatistics()

	log.Info("Starting Kmedia unit updates")

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

	log.Info("Updating content units")
	utils.Must(doUnits())

	//log.Info("Updating collections")
	//utils.Must(doCollections())

	stats.dump()

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func doUnits() error {
	log.Info("Loading all units with kmedia_id")
	units, err := models.ContentUnits(mdb,
		//qm.Where("created_at > '2018-01-01'"),
		qm.Where("properties -> 'kmedia_id' is not null")).
		All()
	if err != nil {
		return errors.Wrap(err, "Load units from mdb")
	}

	log.Infof("Got %d units", len(units))
	log.Info("Setting up workers")
	jobs := make(chan *models.ContentUnit, 100)
	var workersWG sync.WaitGroup
	for w := 1; w <= 5; w++ {
		workersWG.Add(1)
		go updateUnitWorker(jobs, &workersWG)
	}

	log.Info("Queueing work")
	for _, u := range units {
		jobs <- u
	}

	log.Info("Closing jobs channel")
	close(jobs)

	log.Info("Waiting for workers to finish")
	workersWG.Wait()

	return nil
}

func doCollections() error {
	log.Info("Loading all collections with kmedia_id")
	collections, err := models.Collections(mdb,
		qm.Where("properties -> 'kmedia_id' is not null")).
		All()
	if err != nil {
		return errors.Wrap(err, "Load collections from mdb")
	}

	log.Infof("Got %d collections", len(collections))
	log.Info("Setting up workers")
	jobs := make(chan *models.Collection, 100)
	var workersWG sync.WaitGroup
	for w := 1; w <= 5; w++ {
		workersWG.Add(1)
		go updateCollectionWorker(jobs, &workersWG)
	}

	log.Info("Queueing work")
	for _, u := range collections {
		jobs <- u
	}

	log.Info("Closing jobs channel")
	close(jobs)

	log.Info("Waiting for workers to finish")
	workersWG.Wait()

	return nil
}

func updateUnitWorker(jobs <-chan *models.ContentUnit, wg *sync.WaitGroup) {
	for u := range jobs {
		var props map[string]interface{}
		u.Properties.Unmarshal(&props)
		kID := props["kmedia_id"]

		var lecturerID null.Int
		err := queries.Raw(kmdb, "select lecturer_id from containers where id = $1", kID).QueryRow().Scan(&lecturerID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Warnf("no kmedia container ID: %d", int(kID.(float64)))
			} else {
				log.Error(err)
				debug.PrintStack()
				continue
			}
		}
		stats.ContainersVisited.Inc(1)

		err = u.L.LoadContentUnitsPersons(mdb, true, u)
		if err != nil {
			log.Error(err)
			debug.PrintStack()
			continue
		}

		hasRav := false
		for i := range u.R.ContentUnitsPersons {
			if u.R.ContentUnitsPersons[i].PersonID == 1 {
				hasRav = true
				break
			}
		}

		if lecturerID.IsZero() {
			// no rav
			if hasRav {
				log.Infof("norav has rav: kmID %d mdbID %d", int(kID.(float64)), u.ID)
				stats.UnkownCatalogs.Inc("norav => rav", 1)

				// TODO: fix manually as there are only few
			} else {
				// These are good
				stats.UnkownCatalogs.Inc("norav => norav", 1)
			}
		} else {
			// rav
			if !hasRav {
				//log.Infof("Missing Rav: kmID %d mdbID %d", int(kID.(float64)), u.ID)
				stats.UnkownCatalogs.Inc("rav => norav", 1)

				tx, err := mdb.Begin()
				utils.Must(err)

				// create association to RAV
				cup := models.ContentUnitsPerson{
					PersonID:      1, // rav in mdb
					ContentUnitID: u.ID,
					RoleID:        1, // lecturer
				}
				cup.Insert(tx)

				utils.Must(tx.Commit())
			} else {
				// These are good
				stats.UnkownCatalogs.Inc("rav => rav", 1)
			}
		}

		//descriptions, err := kmodels.ContainerDescriptions(kmdb,
		//	qm.Where("container_id = ?", kID)).
		//	All()
		//if err != nil {
		//	log.Error(err)
		//	debug.PrintStack()
		//	continue
		//}
		//
		//tx, err := mdb.Begin()
		//utils.Must(err)
		//
		//for _, d := range descriptions {
		//	if (d.ContainerDesc.Valid && d.ContainerDesc.String != "") ||
		//		(d.Descr.Valid && d.Descr.String != "") {
		//		cui18n := models.ContentUnitI18n{
		//			ContentUnitID: u.ID,
		//			Language:      api.LANG_MAP[d.LangID.String],
		//			Name:          d.ContainerDesc,
		//			Description:   d.Descr,
		//		}
		//		err = cui18n.Upsert(tx,
		//			true,
		//			[]string{"content_unit_id", "language"},
		//			[]string{"name", "description"})
		//		if err != nil {
		//			log.Error(err)
		//			debug.PrintStack()
		//			utils.Must(tx.Rollback())
		//		}
		//	}
		//}
		//
		//utils.Must(tx.Commit())
	}
	wg.Done()
}

func updateCollectionWorker(jobs <-chan *models.Collection, wg *sync.WaitGroup) {
	for u := range jobs {
		var props map[string]interface{}
		u.Properties.Unmarshal(&props)
		kID := props["kmedia_id"]

		descriptions, err := kmodels.CatalogDescriptions(kmdb,
			qm.Where("catalog_id = ?", kID)).
			All()
		if err != nil {
			log.Error(err)
			debug.PrintStack()
			continue
		}

		tx, err := mdb.Begin()
		utils.Must(err)

		for _, d := range descriptions {
			if d.Name.Valid && d.Name.String != "" {
				ci18n := models.CollectionI18n{
					CollectionID: u.ID,
					Language:     api.LANG_MAP[d.LangID.String],
					Name:         d.Name,
				}
				err = ci18n.Upsert(tx,
					true,
					[]string{"collection_id", "language"},
					[]string{"name"})
				if err != nil {
					log.Error(err)
					debug.PrintStack()
					utils.Must(tx.Rollback())
				}
			}
		}

		utils.Must(tx.Commit())
	}
	wg.Done()
}
