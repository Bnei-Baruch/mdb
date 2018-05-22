package tagging

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"regexp"
	"strings"
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
	utils.Must(api.InitTypeRegistries(mdb))

	return clock
}

func Shutdown() {
	utils.Must(mdb.Close())
}

type CollectionWName struct {
	*models.Collection
	name string
}

func (c *CollectionWName) Name() string {
	if c.name != "" {
		return c.name
	}

	ci18ns := make(map[string]string)
	for i := range c.R.CollectionI18ns {
		i18n := c.R.CollectionI18ns[i]
		if i18n.Name.Valid {
			ci18ns[i18n.Language] = i18n.Name.String
		}
	}
	if v, ok := ci18ns[api.LANG_HEBREW]; ok {
		c.name = v
	} else if v, ok := ci18ns[api.LANG_ENGLISH]; ok {
		c.name = v
	} else if v, ok := ci18ns[api.LANG_RUSSIAN]; ok {
		c.name = v
	}
	return c.name
}

type UnitWName struct {
	*models.ContentUnit
	name        string
	description string
}

func (cu *UnitWName) Name() string {
	if cu.name != "" {
		return cu.name
	}

	cui18ns := make(map[string]string)
	for i := range cu.R.ContentUnitI18ns {
		i18n := cu.R.ContentUnitI18ns[i]
		if i18n.Name.Valid {
			cui18ns[i18n.Language] = i18n.Name.String
		}
	}
	if v, ok := cui18ns[api.LANG_HEBREW]; ok {
		cu.name = v
	} else if v, ok := cui18ns[api.LANG_ENGLISH]; ok {
		cu.name = v
	} else if v, ok := cui18ns[api.LANG_RUSSIAN]; ok {
		cu.name = v
	}
	return cu.name
}

func (cu *UnitWName) Description() string {
	if cu.description != "" {
		return cu.description
	}

	cui18ns := make(map[string]string)
	for i := range cu.R.ContentUnitI18ns {
		i18n := cu.R.ContentUnitI18ns[i]
		if i18n.Description.Valid {
			cui18ns[i18n.Language] = i18n.Description.String
		}
	}
	if v, ok := cui18ns[api.LANG_HEBREW]; ok {
		cu.description = v
	} else if v, ok := cui18ns[api.LANG_ENGLISH]; ok {
		cu.description = v
	} else if v, ok := cui18ns[api.LANG_RUSSIAN]; ok {
		cu.description = v
	}
	return cu.description
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
