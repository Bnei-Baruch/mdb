package cleanup

import (
	"context"
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/events"
	"github.com/Bnei-Baruch/mdb/utils"
)

var (
	mdb *sql.DB
)

func Init() (time.Time, *events.BufferedEmitter) {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting units cleanup")

	log.Info("Setting up connection to MDB")
	mdb, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	//boil.DebugMode = true

	log.Info("Initializing static data from MDB")
	utils.Must(common.InitTypeRegistries(mdb))

	log.Info("Setting events handler")
	emitter, err := events.InitEmitter()
	utils.Must(err)

	return clock, emitter
}

func Shutdown() {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	events.CloseEmitter(ctx)
	utils.Must(mdb.Close())
}
