package roza

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/utils"
)

var (
	mdb  *sql.DB
	kmdb *sql.DB
)

func Init() time.Time {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting Roza Importer")

	log.Info("Setting up connection to MDB")
	mdb, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	//boil.DebugMode = true

	log.Info("Setting up connection to Kmedia")
	kmdb, err = sql.Open("postgres", viper.GetString("kmedia.url"))
	utils.Must(err)
	utils.Must(kmdb.Ping())

	log.Info("Initializing static data from MDB")
	utils.Must(common.InitTypeRegistries(mdb))

	return clock
}

func Shutdown() {
	utils.Must(mdb.Close())
	utils.Must(kmdb.Close())
}

func LoadIndex() {
	clock := Init()

	idx := new(RozaIndex)

	utils.Must(idx.Load(mdb))

	log.Infof("%d roots", len(idx.Roots))
	for _, v := range idx.Roots {
		v.printRec("")
	}

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}
