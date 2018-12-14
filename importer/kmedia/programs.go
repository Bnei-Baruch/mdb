package kmedia

import (
	"database/sql"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

const CHAPTERS_FILE = "importer/kmedia/data/programs_chapters.csv"

func ImportProgramsChapters() {
	clock := Init()

	stats = NewImportStatistics()
	utils.Must(importChapters())
	stats.dump()

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func importChapters() error {
	// Read mappings file
	records, err := utils.ReadCSV(CHAPTERS_FILE)
	if err != nil {
		return errors.Wrap(err, "Read CHAPTERS_FILE")
	}
	log.Infof("CHAPTERS_FILE file has %d rows", len(records))

	h, err := utils.ParseCSVHeader(records[0])
	if err != nil {
		return errors.Wrap(err, "Bad header")
	}

	for _, x := range records[1:] {
		catalogID, err := strconv.Atoi(x[h["catalog.id"]])
		if err != nil {
			return errors.Wrapf(err, "Bad catalog.id %s", x[h["catalog.id"]])
		}

		log.Infof("Processing catalog %d", catalogID)
		stats.CatalogsProcessed.Inc(1)

		collection, err := models.Collections(mdb,
			qm.Where("(properties->>'kmedia_id')::int = ?", catalogID)).
			One()
		if err != nil {
			return errors.Wrapf(err, "Lookup collection in mdb [kmid %d]", catalogID)
		}

		chaptersArr := x[h["containers.ids"]]
		chaptersArr = chaptersArr[1 : len(chaptersArr)-1]
		chaptersIDs := strings.Split(chaptersArr, ",")

		for i := range chaptersIDs {
			containerID, err := strconv.Atoi(chaptersIDs[i])
			if err != nil {
				return errors.Wrapf(err, "Bad containerID %s", chaptersIDs[i])
			}

			container, err := kmodels.FindContainer(kmdb, containerID)
			if err != nil {
				if err == sql.ErrNoRows {
					log.Warnf("Container doesn't exists in kmedia %s", chaptersIDs[i])
					continue
				}
				return errors.Wrapf(err, "Lookup container in kmedia %s", chaptersIDs[i])
			}

			tx, err := mdb.Begin()
			if err != nil {
				return errors.Wrap(err, "Start transaction")
			}

			if err = importChapter(tx, collection, container); err != nil {
				utils.Must(tx.Rollback())
				stats.TxRolledBack.Inc(1)
				log.Error(err)
				debug.PrintStack()
				continue
			} else {
				utils.Must(tx.Commit())
				stats.TxCommitted.Inc(1)
			}
		}
	}

	return nil
}

func importChapter(exec boil.Executor, collection *models.Collection, container *kmodels.Container) error {
	stats.ContainersProcessed.Inc(1)

	exists, err := models.ContentUnits(mdb,
		qm.Where("(properties->>'kmedia_id')::int = ?", container.ID),
	).Exists()
	if err != nil {
		return errors.Wrapf(err, "Lookup content unit kmid %d", container.ID)
	}

	if exists {
		return importChapterExistingUnit(exec, collection, container)
	} else {
		return importChapterNewUnit(exec, collection, container)
	}
}

func importChapterExistingUnit(exec boil.Executor, collection *models.Collection, container *kmodels.Container) error {
	_, err := importContainer(exec, container, collection, "",
		strconv.Itoa(container.Position.Int), container.Position.Int)
	return err
}

func importChapterNewUnit(exec boil.Executor, collection *models.Collection, container *kmodels.Container) error {
	err := container.L.LoadFileAssets(kmdb, true, container)
	if err != nil {
		return errors.Wrapf(err, "Load kmedia file assets %d", container.ID)
	}

	// Create import operation
	operation, err := api.CreateOperation(exec, api.OP_IMPORT_KMEDIA,
		api.Operation{WorkflowID: strconv.Itoa(container.ID)}, nil)
	if err != nil {
		return errors.Wrapf(err, "Create operation %d", container.ID)
	}
	stats.OperationsCreated.Inc(1)

	// import container
	ccuName := strconv.Itoa(container.Position.Int)
	unit, err := importContainer(exec, container, collection, api.CT_VIDEO_PROGRAM_CHAPTER,
		ccuName, container.Position.Int)
	if err != nil {
		return errors.Wrapf(err, "Import container %d", container.ID)
	}

	// import container files
	var file *models.File
	for _, fileAsset := range container.R.FileAssets {
		log.Infof("Processing file_asset %d", fileAsset.ID)
		stats.FileAssetsProcessed.Inc(1)

		// Create or update MDB file
		file, err = importFileAsset(exec, fileAsset, unit, operation)
		if err != nil {
			log.Error(err)
			debug.PrintStack()
			break
		}
		if file != nil && file.Published {
			unit.Published = true
		}
	}
	if err != nil {
		return errors.Wrapf(err, "Import container files %d", container.ID)
	}

	if unit.Published {
		collection.Published = true
		err = unit.Update(exec, "published")
		if err != nil {
			return errors.Wrapf(err, "Update unit published column %d", container.ID)
		}
	}

	return nil
}
