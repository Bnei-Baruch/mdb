package sources

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/nullbio/null.v6"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

const (
	BASE_PATH        = "importer/sources/data/Sources - "
	AUTHORS_FILE     = BASE_PATH + "Authors.csv"
	COLLECTIONS_FILE = BASE_PATH + "Collections.csv"
)

var (
	LANGS = [7]string{
		api.LANG_ENGLISH,
		api.LANG_HEBREW,
		api.LANG_RUSSIAN,
		api.LANG_GERMAN,
		api.LANG_SPANISH,
		api.LANG_TURKISH,
		api.LANG_UKRAINIAN,
	}
)

type SourceInput struct {
	models.Source
	FileName string
	Handler  func(exec boil.Executor, author *models.Author) error
}

func ImportSources() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting sources import")

	log.Info("Setting up connection to MDB")
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()
	boil.SetDB(mdb)
	//boil.DebugMode = true

	utils.Must(api.SOURCE_TYPE_REGISTRY.Init())

	utils.Must(handleAuthors(mdb))
	utils.Must(handleCollections(mdb))

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func handleAuthors(db *sql.DB) error {
	records, err := utils.ReadCSV(AUTHORS_FILE)
	if err != nil {
		return errors.Wrap(err, "Read authors")
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
		if err = doAuthor(tx, h, x); err != nil {
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

func doAuthor(exec boil.Executor, header map[string]int, record []string) error {
	// Get or create Author
	author, err := models.Authors(exec, qm.Where("code = ?", record[header["code"]])).One()
	if err != nil {
		if err == sql.ErrNoRows {
			// Create
			author = &models.Author{
				Code:     record[header["code"]],
				Name:     record[header["name"]],
				FullName: null.NewString(record[header["full name"]], record[header["full name"]] != ""),
			}
			err = author.Insert(exec)
			if err != nil {
				return errors.Wrapf(err, "Insert author [%s]", record)
			}
		} else {
			return errors.Wrapf(err, "Lookup author in db [%s]", record)
		}
	} else {
		author.Name = record[header["name"]]
		author.FullName = null.NewString(record[header["full name"]], record[header["full name"]] != "")
		err = author.Update(exec, "name", "full_name")
		if err != nil {
			return errors.Wrapf(err, "Update author [%d] [%s]", author.ID, record)
		}
	}

	// i18n
	for _, l := range LANGS {
		n := record[header[l+".name"]]
		fn := record[header[l+".full_name"]]
		if n == "" && fn == "" {
			continue
		}

		ai18n := models.AuthorI18n{
			AuthorID: author.ID,
			Language: l,
			Name:     null.NewString(n, n != ""),
			FullName: null.NewString(fn, fn != ""),
		}
		err = ai18n.Upsert(exec, true,
			[]string{"author_id", "language"},
			[]string{"name", "full_name"})
		if err != nil {
			return errors.Wrapf(err, "Upsert author i18n")
		}
	}

	return nil
}

func handleCollections(db *sql.DB) error {
	records, err := utils.ReadCSV(COLLECTIONS_FILE)
	if err != nil {
		return errors.Wrap(err, "Read collections")
	}

	h, err := utils.ParseCSVHeader(records[0])
	if err != nil {
		return errors.Wrap(err, "Bad header")
	}

	var tx *sql.Tx
	for _, x := range records[1:] {
		tx, err = db.Begin()
		if err != nil {
			return errors.Wrap(err, "Start transaction")
		}

		err = doCollection(tx, h, x)

		if err == nil {
			err = tx.Commit()
			if err != nil {
				err = errors.Wrap(err, "Commit transaction")
				break
			}
		} else {
			if ex := tx.Rollback(); ex != nil {
				err = errors.Wrap(ex, "Rollback transaction")
			}
			break
		}
	}

	return err
}

func doCollection(exec boil.Executor, header map[string]int, record []string) error {
	authorCode := record[header["author"]]
	if authorCode == "" {
		return errors.New("Empty author code")
	}
	name := record[header["name"]]
	if name == "" {
		return errors.New("Empty collection name")
	}
	pattern := record[header["pattern"]]
	log.Infof("Author: %s, Name: %s", authorCode, name)

	// Fetch author
	author, err := models.Authors(exec, qm.Where("code = ?", authorCode)).One()
	if err != nil {
		return errors.Wrapf(err, "Fetch author [%s]", authorCode)
	}

	// Get or create collection source
	collection, err := models.Sources(exec,
		qm.InnerJoin("authors_sources x on x.source_id = sources.id and author_id = ?", author.ID),
		qm.Where("name = ? and parent_id is null", name)).
		One()
	if err == nil {
		// update
		if collection.Pattern.Valid || pattern != "" {
			collection.Pattern = null.NewString(pattern, pattern != "")
			err = collection.Update(exec, "pattern")
			if err != nil {
				return errors.Wrapf(err, "Update collection [%s]", collection.ID)
			}
		}
	} else {
		if err == sql.ErrNoRows {
			// create
			collection = &models.Source{
				UID:     utils.GenerateUID(8),
				Name:    name,
				Pattern: null.NewString(pattern, pattern != ""),
				TypeID:  api.SOURCE_TYPE_REGISTRY.ByName[api.SRC_COLLECTION].ID,
			}
			err = author.AddSources(exec, true, collection)
			if err != nil {
				return errors.Wrapf(err, "Create collection [%s %s]", authorCode, name)
			}
		} else {
			return errors.Wrapf(err, "Lookup collection in db [%s, %s]", authorCode, name)
		}
	}

	// i18n
	for _, l := range LANGS {
		n := record[header[l+".name"]]
		d := record[header[l+".description"]]
		if n == "" && d == "" {
			continue
		}

		si18n := models.SourceI18n{
			SourceID:    collection.ID,
			Language:    l,
			Name:        null.NewString(n, n != ""),
			Description: null.NewString(d, d != ""),
		}
		err = si18n.Upsert(exec, true,
			[]string{"source_id", "language"},
			[]string{"name", "description"})
		if err != nil {
			return errors.Wrapf(err, "Upsert collection i18n")
		}
	}

	// Content
	fn := fmt.Sprintf("%s%s-%s.csv",
		BASE_PATH,
		strings.ToLower(authorCode),
		strings.Replace(strings.ToLower(name), " ", "-", -1))

	records, err := utils.ReadCSV(fn)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warnf("Input missing: %s", fn)
			return nil
		}
		return errors.Wrap(err, "Read collection contents")
	}

	h, err := utils.ParseCSVHeader(records[0])
	if err != nil {
		return errors.Wrap(err, "Bad header")
	}

	var parents = []*models.Source{collection}
	for i, x := range records[1:] {
		if utils.IsEmpty(x) {
			continue
		}

		xLevel := x[h["level"]]
		level, err := strconv.Atoi(xLevel)
		if err != nil {
			return errors.Wrapf(err, "Bad level at row %d: %s", i+1, xLevel)
		}

		xType := x[h["type"]]
		sType, ok := api.SOURCE_TYPE_REGISTRY.ByName[xType]
		if !ok {
			return errors.Errorf("Unknown source type at row %d: %s", i+1, xType)
		}

		name := x[h["name"]]
		if name == "" {
			return errors.Errorf("Missing name at row %d", i)
		}

		position := -1
		xPosition := x[h["position"]]
		if xPosition != "" {
			position, err = strconv.Atoi(xPosition)
			if err != nil {
				return errors.Wrapf(err, "Bad position: %s", xPosition)
			}
		}

		pattern = x[h["pattern"]]
		description := x[h["description"]]

		// Get or Create source
		parent := parents[level-1]
		source, err := models.Sources(exec,
			qm.Where("type_id = ? and parent_id = ? and name = ?", sType.ID, parent.ID, name)).
			One()
		if err == nil {
			// update
			source.Description = null.NewString(description, description != "")
			source.Pattern = null.NewString(pattern, pattern != "")
			source.Position = null.NewInt(position, position != -1)
			err = source.Update(exec, "description", "pattern", "position")
			if err != nil {
				return errors.Wrapf(err, "Update source [%d %d %s]", sType.ID, parent.ID, name)
			}
		} else {
			if err == sql.ErrNoRows {
				// create
				source = &models.Source{
					UID:         utils.GenerateUID(8),
					TypeID:      sType.ID,
					Name:        name,
					Description: null.NewString(description, description != ""),
					Pattern:     null.NewString(pattern, pattern != ""),
					ParentID:    null.Int64From(parent.ID),
					Position:    null.NewInt(position, position != -1),
				}
				err = source.Insert(exec)
				if err != nil {
					return errors.Wrapf(err, "Insert source [%s]", x)
				}
			} else {
				return errors.Wrapf(err, "Fetch source [%d %d %s]", sType.ID, parent.ID, name)
			}
		}

		// Source i18n
		for _, l := range LANGS {
			n := x[h[l+".name"]]
			d := x[h[l+".description"]]
			if n == "" && d == "" {
				continue
			}

			si18n := models.SourceI18n{
				SourceID:    source.ID,
				Language:    l,
				Name:        null.NewString(n, n != ""),
				Description: null.NewString(d, d != ""),
			}
			err = si18n.Upsert(exec, true,
				[]string{"source_id", "language"},
				[]string{"name", "description"})
			if err != nil {
				return errors.Wrapf(err, "Upsert source i18n")
			}
		}

		if level == len(parents) {
			parents = append(parents, source)
		} else {
			parents[level] = source
		}
	}

	return nil
}
