package tvshows

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

const (
	TV_SHOWS_FILE = "importer/tvshows/data/TV Shows - final.csv"
)

var (
	LANGS = [7]string{
		common.LANG_ENGLISH,
		common.LANG_HEBREW,
		common.LANG_RUSSIAN,
		common.LANG_SPANISH,
		common.LANG_GERMAN,
		common.LANG_UKRAINIAN,
		common.LANG_CHINESE,
	}
)

func ImportTVShows() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting tv shows import")

	log.Info("Setting up connection to MDB")
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()
	boil.SetDB(mdb)
	//boil.DebugMode = true

	utils.Must(common.InitTypeRegistries(mdb))

	utils.Must(handleTVShows(mdb))

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func handleTVShows(db *sql.DB) error {
	records, err := utils.ReadCSV(TV_SHOWS_FILE)
	if err != nil {
		return errors.Wrap(err, "Read tv shows")
	}

	h, err := utils.ParseCSVHeader(records[0])
	if err != nil {
		return errors.Wrap(err, "Bad header")
	}

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "Start transaction")
	}

	for _, x := range records[1:] {
		if err = doTVShow(tx, h, x); err != nil {
			break
		}
	}

	if err == nil {
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err, "Commit transaction")
		}
	} else {
		if ex := tx.Rollback(); ex != nil {
			return errors.Wrap(ex, "Rollback transaction")
		}
		return err
	}

	return nil
}

func doTVShow(exec boil.Executor, header map[string]int, record []string) error {
	// Get or create TV Show
	ctID := common.CONTENT_TYPE_REGISTRY.ByName[common.CT_VIDEO_PROGRAM].ID
	show, err := models.Collections(
		qm.Where("type_id = ? AND (properties->>'kmedia_id')::int = ?", ctID, record[header["kmedia_id"]])).
		One(exec)
	if err != nil {
		if err == sql.ErrNoRows {
			// Create
			show = &models.Collection{
				UID:    common.GenerateUID(8),
				TypeID: ctID,
			}
			err = show.Insert(exec, boil.Infer())
			if err != nil {
				return errors.Wrapf(err, "Insert show [%s]", record)
			}
		} else {
			return errors.Wrapf(err, "Lookup show in db [%s]", record)
		}
	}

	// Properties
	var props = make(map[string]interface{})
	if show.Properties.Valid {
		show.Properties.Unmarshal(&props)
	}
	props["kmedia_id"] = record[header["kmedia_id"]]
	props["pattern"] = record[header["mdb_pattern"]]
	props["active"] = strings.ToLower(strings.TrimSpace(record[header["active"]])) == "v"
	dl := record[header["language"]]
	if dl != "" {
		if l, ok := common.LANG_MAP[dl]; ok {
			props["default_language"] = l
		}
	}
	p, err := json.Marshal(props)
	if err != nil {
		return errors.Wrap(err, "Marshal show properties")
	}
	show.Properties = null.JSONFrom(p)

	_, err = show.Update(exec, boil.Infer())
	if err != nil {
		return errors.Wrap(err, "Update show properties")
	}

	// i18n
	for _, l := range LANGS {
		n := record[header[l+".name"]]
		if n == "" {
			continue
		}

		ci18n := models.CollectionI18n{
			CollectionID: show.ID,
			Language:     l,
			Name:         null.NewString(n, n != ""),
		}
		err = ci18n.Upsert(exec, true,
			[]string{"collection_id", "language"},
			boil.Whitelist("name"),
			boil.Infer())
		if err != nil {
			return errors.Wrapf(err, "Upsert show i18n")
		}
	}

	return nil
}
