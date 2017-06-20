package kmedia

import (
	"database/sql"
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
	"runtime/debug"
)

func UpdateUnits() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

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

	log.Info("Loading all units with kmedia_id")
	units, err := models.ContentUnits(mdb,
		qm.Where("properties -> 'kmedia_id' is not null")).
		All()
	utils.Must(err)
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

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func updateUnitWorker(jobs <-chan *models.ContentUnit, wg *sync.WaitGroup) {
	for u := range jobs {
		var props map[string]interface{}
		u.Properties.Unmarshal(&props)
		kID := props["kmedia_id"]

		descriptions, err := kmodels.ContainerDescriptions(kmdb,
			qm.Where("container_id = ?", kID)).
			All()
		if err != nil {
			log.Error(err)
			debug.PrintStack()
			continue
		}

		tx, err := mdb.Begin()
		utils.Must(err)

		for _, d := range descriptions {
			if (d.ContainerDesc.Valid && d.ContainerDesc.String != "") ||
				(d.Descr.Valid && d.Descr.String != "") {
				cui18n := models.ContentUnitI18n{
					ContentUnitID: u.ID,
					Language:      api.LANG_MAP[d.LangID.String],
					Name:          d.ContainerDesc,
					Description:   d.Descr,
				}
				err = cui18n.Upsert(tx,
					true,
					[]string{"content_unit_id", "language"},
					[]string{"name", "description"})
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
