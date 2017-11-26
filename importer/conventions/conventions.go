package conventions

import (
	"database/sql"
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"gopkg.in/volatiletech/null.v6"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

const (
	CONVENTIONS_FILE = "importer/convetions/data/Conventions - 2017.csv"
)

var (
	LANGS = [4]string{
		api.LANG_ENGLISH,
		api.LANG_HEBREW,
		api.LANG_RUSSIAN,
		api.LANG_SPANISH,
	}
)

func ImportConvetions() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting convetions insert")

	log.Info("Setting up connection to MDB")
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()
	boil.SetDB(mdb)
	//boil.DebugMode = true

	utils.Must(api.InitTypeRegistries(mdb))

	utils.Must(handleConventions(mdb))

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func handleConventions(db *sql.DB) error {
	records, err := utils.ReadCSV(CONVENTIONS_FILE)
	if err != nil {
		return errors.Wrap(err, "Read conventions")
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
		if err = doConvention(tx, h, x); err != nil {
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

func doConvention(exec boil.Executor, header map[string]int, record []string) error {
	// Get or create convention
	ctID := api.CONTENT_TYPE_REGISTRY.ByName[api.CT_CONGRESS].ID
	var convention models.Collection
	err := queries.Raw(exec,
		`select * from collections where type_id=$1 and properties -> 'pattern' ? $2 limit 1`,
		ctID, record[header["pattern"]],
	).Bind(&convention)
	if err != nil {
		if err == sql.ErrNoRows {
			// Create
			convention = models.Collection{
				UID:    utils.GenerateUID(8),
				TypeID: ctID,
			}
			err = convention.Insert(exec)
			if err != nil {
				return errors.Wrapf(err, "Insert convention [%s]", record)
			}
		} else {
			return errors.Wrapf(err, "Lookup convention in db [%s]", record)
		}
	}

	// Properties
	var props = make(map[string]interface{})
	if convention.Properties.Valid {
		convention.Properties.Unmarshal(&props)
	}
	props["pattern"] = record[header["pattern"]]
	props["active"] = true
	props["country"] = record[header["country"]]
	props["city"] = record[header["city"]]
	props["full_address"] = record[header["full_address"]]
	sd, err := time.Parse("2006-01-02", record[header["start_date"]])
	if err != nil {
		return errors.Wrapf(err, "Bad start_date format, expected `2006-01-02` got %s", record[header["start_date"]])
	}
	props["start_date"] = sd
	ed, err := time.Parse("2006-01-02", record[header["end_date"]])
	if err != nil {
		return errors.Wrapf(err, "Bad end_date format, expected `2006-01-02` got %s", record[header["end_date"]])
	}
	props["end_date"] = ed

	p, err := json.Marshal(props)
	if err != nil {
		return errors.Wrap(err, "Marshal convention properties")
	}
	convention.Properties = null.JSONFrom(p)

	err = convention.Update(exec)
	if err != nil {
		return errors.Wrap(err, "Update convention properties")
	}

	// i18n
	for _, l := range LANGS {
		n := record[header[l+".name"]]
		if n == "" {
			continue
		}

		ci18n := models.CollectionI18n{
			CollectionID: convention.ID,
			Language:     l,
			Name:         null.NewString(n, n != ""),
		}
		err = ci18n.Upsert(exec, true,
			[]string{"collection_id", "language"},
			[]string{"name"})
		if err != nil {
			return errors.Wrapf(err, "Upsert convention i18n")
		}
	}

	return nil
}
