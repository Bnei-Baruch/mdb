package metus

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/denisenkom/go-mssqldb"
		"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/utils"
)

var (
	mdb     *sql.DB
	metusDB *sql.DB
)

func Init() time.Time {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting metus import")

	log.Info("Setting up connection to MDB")
	mdb, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	//boil.DebugMode = true

	log.Info("Initializing static data from MDB")
	utils.Must(api.InitTypeRegistries(mdb))

	metusDB, err = sql.Open("mssql", viper.GetString("metus.url"))
	utils.Must(err)
	utils.Must(metusDB.Ping())

	return clock
}

func Shutdown() {
	utils.Must(mdb.Close())
	utils.Must(metusDB.Close())
}

