package batch

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/boil"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

var mdb *sql.DB

type UnitNames struct {
	UnitID int64
	Names  map[string]string
}

func RenameUnits() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting unit renaming")

	log.Info("Setting up connection to MDB")
	mdb, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()
	boil.SetDB(mdb)
	//boil.DebugMode = true

	log.Info("Initializing static data from MDB")
	utils.Must(api.InitTypeRegistries(mdb))

	log.Info("Loading all units")
	units, err := models.ContentUnits(mdb).All()
	utils.Must(err)
	log.Infof("Got %d units", len(units))

	log.Info("Setting up workers")
	jobs := make(chan *models.ContentUnit, 100)
	results := make(chan UnitNames, 100)
	done := make(chan bool)

	var workersWG sync.WaitGroup
	for w := 1; w <= 5; w++ {
		workersWG.Add(1)
		go namesUnitWorker(jobs, results, &workersWG)
	}
	go namesWriter(results, done)

	log.Info("Queueing work")
	for _, u := range units {
		jobs <- u
	}

	log.Info("Closing jobs channel")
	close(jobs)

	log.Info("Waiting for workers to finish")
	workersWG.Wait()

	log.Info("Closing results channel")
	close(results)

	log.Info("Waiting for writer to finish")
	<-done

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func namesUnitWorker(jobs <-chan *models.ContentUnit, results chan UnitNames, wg *sync.WaitGroup) {
	for cu := range jobs {
		metadata := api.CITMetadata{}

		describer, ok := api.CUDescribers[api.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name]
		if !ok {
			describer = api.GenericDescriber{}
		}

		i18ns, err := describer.DescribeContentUnit(mdb, cu, metadata)
		if err != nil {
			log.Errorf("Error naming unit [%d]", cu.ID)
			debug.PrintStack()
			continue
		}

		names := make(map[string]string, len(i18ns))
		for _, i18n := range i18ns {
			if i18n.Name.Valid {
				names[i18n.Language] = i18n.Name.String
			}
		}

		results <- UnitNames{UnitID: cu.ID, Names: names}
	}
	wg.Done()
}

func namesWriter(results <-chan UnitNames, done chan bool) {
	f, err := ioutil.TempFile("/tmp", "unit_names_")
	if err != nil {
		utils.Must(err)
	}
	defer f.Close()
	log.Infof("Report file: %s", f.Name())

	_, err = fmt.Fprintf(f, "ID\t%s\n", strings.Join(api.ALL_LANGS, "\t"))
	utils.Must(err)

	for un := range results {
		if len(un.Names) == 0 {
			log.Infof("No name [%d]", un.UnitID)
			continue
		}

		values := make([]string, len(api.ALL_LANGS)+1)
		values[0] = strconv.FormatInt(un.UnitID, 10)
		for i, language := range api.ALL_LANGS {
			if name, ok := un.Names[language]; ok {
				values[i+1] = name
			} else {
				values[i+1] = ""
			}
		}

		_, err = fmt.Fprintln(f, strings.Join(values, "\t"))
		utils.Must(err)
	}

	done <- true
}
