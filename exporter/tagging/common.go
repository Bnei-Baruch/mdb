package tagging

import (
	"database/sql"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/utils"
)

var (
	mdb *sql.DB
)

func Init() time.Time {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting Tagging Exporter")

	log.Info("Setting up connection to MDB")
	mdb, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	//boil.DebugMode = true

	log.Info("Initializing static data from MDB")
	utils.Must(common.InitTypeRegistries(mdb))

	return clock
}

func Shutdown() {
	utils.Must(mdb.Close())
}

var separators = regexp.MustCompile(`[ &_=+:\s]`)
var dashes = regexp.MustCompile(`[\-]+`)
var illegalName = regexp.MustCompile(`[\[\]\\/?:*]`)

func CleanSheetName(s string) string {

	// Remove any trailing space to avoid ending on -
	s = strings.Trim(s, " ")

	// Replace certain joining characters with a dash
	s = separators.ReplaceAllString(s, "-")

	// Remove all other unrecognised characters - NB we do allow any printable characters
	s = illegalName.ReplaceAllString(s, "")

	// Remove any multiple dashes caused by replacements above
	s = dashes.ReplaceAllString(s, "-")

	return s

}
