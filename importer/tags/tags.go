package tags

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/volatiletech/null.v6"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

const (
	TAGS_FILE = "importer/tags/data/Tags - All.csv"
)

var mappings map[int]int64
var LANGS = [5]string{
	common.LANG_ENGLISH,
	common.LANG_HEBREW,
	common.LANG_RUSSIAN,
	common.LANG_SPANISH,
	common.LANG_UKRAINIAN,
}

func ImportTags() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting tags import")

	log.Info("Setting up connection to MDB")
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()
	boil.SetDB(mdb)
	//boil.DebugMode = true

	mappings = make(map[int]int64)
	utils.Must(handleTopics(mdb))
	log.Infof("Here comes %d catalogs mappings", len(mappings))
	for k, v := range mappings {
		fmt.Printf("%d\t%d\n", k, v)
	}

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func handleTopics(db *sql.DB) error {
	records, err := utils.ReadCSV(TAGS_FILE)
	if err != nil {
		return errors.Wrap(err, "Read topics")
	}

	h, err := utils.ParseCSVHeader(records[0])
	if err != nil {
		return errors.Wrap(err, "Bad header")
	}

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "Start transaction")
	}

	var parents = []*models.Tag{}
	for i, x := range records[1:] {
		if utils.IsEmpty(x) {
			continue
		}

		xLevel := x[h["level"]]
		level, err := strconv.Atoi(xLevel)
		if err != nil {
			return errors.Wrapf(err, "Bad level at row %d: %s", i+1, xLevel)
		}

		kmdbID := -1
		xKmdbID := x[h["kmedia_catalog"]]
		if xKmdbID != "" {
			kmdbID, err = strconv.Atoi(xKmdbID)
			if err != nil {
				return errors.Wrapf(err, "Bad kmedia_catalog at row %d: %s", i+1, xKmdbID)
			}
		}

		pattern := x[h["pattern"]]
		description := x[h["description"]]

		// Get or create tag
		var tag *models.Tag
		var parent *models.Tag
		if level == 1 {
			tag, err = models.Tags(tx,
				qm.Where("parent_id is null and pattern = ?", pattern)).
				One()
		} else {
			parent = parents[level-2]
			tag, err = models.Tags(tx,
				qm.Where("parent_id = ? and pattern = ?", parent.ID, pattern)).
				One()
		}

		if err == nil {
			// update
			if description != "" {
				tag.Description = null.StringFrom(description)
				err = tag.Update(tx, "description")
				if err != nil {
					return errors.Wrapf(err, "Update tag %s", pattern)
				}
			}
		} else {
			if err == sql.ErrNoRows {
				log.Infof("New pattern %s", pattern)
				// create
				tag = &models.Tag{
					UID:         utils.GenerateUID(8),
					Pattern:     null.StringFrom(pattern),
					Description: null.NewString(description, description != ""),
				}
				if parent != nil {
					tag.ParentID = null.Int64From(parent.ID)
				}
				err = tag.Insert(tx)
				if err != nil {
					return errors.Wrapf(err, "Insert tag %s", pattern)
				}
			} else {
				return errors.Wrapf(err, "Fetch tag %s", pattern)
			}
		}

		// i18n
		for _, l := range LANGS {
			label := x[h[l+".label"]]
			if label == "" {
				continue
			}

			ti18n := models.TagI18n{
				TagID:    tag.ID,
				Language: l,
				Label:    null.StringFrom(label),
			}
			err = ti18n.Upsert(tx, true,
				[]string{"tag_id", "language"},
				[]string{"label"})
			if err != nil {
				return errors.Wrapf(err, "Upsert tag [%d] i18n %s", tag.ID, l)
			}
		}

		// kmedia catalogs mappings
		if kmdbID > 0 {
			mappings[kmdbID] = tag.ID
		}

		if level == len(parents)+1 {
			parents = append(parents, tag)
		} else {
			parents[level-1] = tag
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
