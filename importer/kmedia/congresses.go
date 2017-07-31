package kmedia

import (
	"database/sql"
	"time"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/utils"
	"strconv"
)

const CONGRESSES_FILE = "importer/kmedia/data/Conventions - congresses.csv"

type EventPart struct {
	KMediaID int
	Name     string
	Position int
}

type Congress struct {
	KMediaID    int
	Country     string
	City        string
	FullAddress string
	Name        string
	Year        int
	Month       int
	Events      []EventPart
}

var CountryCatalogs = map[int]string{
	7900: "Ukraine",
	8091: "Kazakhstan",
	4556: "Canada",
	4537: "United States",
	4536: "Russia",
	4543: "Estonia",
	4553: "Austria",
	4680: "Colombia",
	4563: "Switzerland",
	4544: "Spain",
	4549: "Argentina",
	4787: "Sweden",
	4788: "Bulgaria",
	2323: "Mexico",
	4545: "Italy",
	4552: "Chile",
	4667: "Brazil",
	4551: "England",
	4547: "Germany",
	4550: "Turkey",
	7872: "Romania",
	7910: "Czech",
	4658: "Lithuania",
	4741: "France",
	4710: "Georgia",
	4554: "Poland",
	//4555: "Israel",
	//8030: "America_2017",
}

func DumpCongresses() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting Kmedia dump congresses")

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

	congresses, err := initCongresses()
	utils.Must(err)

	utils.Must(loadEventParts(congresses))
	//utils.Must(dumpCongresses())

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func dumpCongresses() error {
	all := make([]*Congress, 0)
	for cID, country := range CountryCatalogs {
		log.Infof("Processing %s", country)

		catalogs, err := kmodels.Catalogs(kmdb, qm.Where("parent_id = ?", cID)).All()
		if err != nil {
			return errors.Wrap(err, "Load catalogs from db")
		}
		log.Infof("Got %d catalogs", len(catalogs))

		for _, catalog := range catalogs {
			c := &Congress{
				KMediaID: catalog.ID,
				Country:  country,
				Name:     catalog.Name,
			}
			all = append(all, c)
		}
	}

	log.Infof("Found %d total congresses", len(all))
	for _, congress := range all {
		fmt.Printf("%d\t%s\t%s\n", congress.KMediaID, congress.Country, congress.Name)
	}

	return nil
}

func loadEventParts(congresses map[int]*Congress) error {


	return nil
}

func initCongresses() (map[int]*Congress, error) {
	// Read mappings file
	records, err := utils.ReadCSV(CONGRESSES_FILE)
	if err != nil {
		return nil, errors.Wrap(err, "Read congresses")
	}
	log.Infof("Congresses file has %d rows", len(records))

	// Create mappings
	mappings := make(map[int]*Congress, len(records)-1)
	for i, r := range records[1:] {
		kmid, err := strconv.Atoi(r[0])
		if err != nil {
			return nil, errors.Wrapf(err, "Bad kmedia_id, row [%d]", i)
		}

		var year int
		if r[6] != "" {
			year, err = strconv.Atoi(r[6])
			if err != nil {
				return nil, errors.Wrapf(err, "Bad year, row [%d]", i)
			}
		}

		var month int
		if r[7] != "" {
			month, err = strconv.Atoi(r[7])
			if err != nil {
				return nil, errors.Wrapf(err, "Bad month, row [%d]", i)
			}
		}

		mappings[kmid] = &Congress{
			KMediaID:    kmid,
			Country:     r[1],
			City:        r[2],
			FullAddress: r[3],
			Year:        year,
			Month:       month,
		}
	}

	return mappings, nil
}
