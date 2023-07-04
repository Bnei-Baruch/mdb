package kmedia

import (
	"database/sql"
	"encoding/json"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/boil"
	qm4 "github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

var (
	kmediaLessonCT *kmodels.ContentType
)

func ImportVirtualLessons() {
	clock := Init()

	log.Info("Loading all virtual_lessons")
	vls, err := kmodels.VirtualLessons(kmdb,
		qm.Where("film_date between '2017-03-24' and '2017-05-08'")).
		All()
	utils.Must(err)
	log.Infof("Got %d lessons", len(vls))

	// Process lessons
	stats = NewImportStatistics()

	log.Info("Setting up workers")
	jobs := make(chan *kmodels.VirtualLesson, 100)
	var workersWG sync.WaitGroup
	for w := 1; w <= 5; w++ {
		workersWG.Add(1)
		go worker(jobs, &workersWG)
	}

	log.Info("Queueing work")
	for _, vl := range vls {
		jobs <- vl
	}

	log.Info("Closing jobs channel")
	close(jobs)

	log.Info("Waiting for workers to finish")
	workersWG.Wait()

	// TODO: clean mdb stale data that no longer exists in kmedia
	// This would require some good understanding of how data can go stale.
	// Files are removed from kmedia ?
	// Containers merged ? what happens to old container ? removed / flagged ?
	// Lessons change ?

	stats.dump()

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func worker(jobs <-chan *kmodels.VirtualLesson, wg *sync.WaitGroup) {
	for vl := range jobs {
		log.Infof("Processing virtual_lesson %d", vl.ID)
		stats.LessonsProcessed.Inc(1)

		// Validate virtual lesson data
		containers, err := getValidContainers(kmdb, vl)
		if err != nil {
			log.Error(err)
			debug.PrintStack()
			continue
		}
		if len(containers) == 0 {
			log.Warnf("Invalid lesson [%d]", vl.ID)
			stats.InvalidLessons.Inc(1)
			continue
		}
		stats.ValidLessons.Inc(1)

		// Begin mdb transaction
		tx, err := mdb.Begin()
		utils.Must(err)

		// Create import operation
		operation, err := api.CreateOperation(tx, common.OP_IMPORT_KMEDIA,
			api.Operation{WorkflowID: strconv.Itoa(vl.ID)}, nil)
		if err != nil {
			utils.Must(tx.Rollback())
			stats.TxRolledBack.Inc(1)
			log.Error(err)
			debug.PrintStack()
			continue
		}
		stats.OperationsCreated.Inc(1)

		// Handle MDB collection
		collection, err := importVirtualLesson(tx, vl)
		if err != nil {
			utils.Must(tx.Rollback())
			stats.TxRolledBack.Inc(1)
			log.Error(err)
			debug.PrintStack()
			continue
		}

		var unit *models.ContentUnit
		for _, container := range containers {
			log.Infof("Processing container %d", container.ID)
			stats.ContainersProcessed.Inc(1)

			// Create or update MDB content_unit
			unit, err = importContainer(tx, container, collection,
				common.CT_LESSON_PART, strconv.Itoa(container.Position.Int), container.Position.Int)
			if err != nil {
				log.Error(err)
				debug.PrintStack()
				break
			}

			var file *models.File
			for _, fileAsset := range container.R.FileAssets {
				log.Infof("Processing file_asset %d", fileAsset.ID)
				stats.FileAssetsProcessed.Inc(1)

				// Create or update MDB file
				file, err = importFileAsset(tx, fileAsset, unit, operation)
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
				break
			}
			if unit.Published {
				collection.Published = true
				_, err = unit.Update(tx, boil.Whitelist("published"))
			}
			if err != nil {
				break
			}
		}

		if collection.Published {
			_, err = collection.Update(tx, boil.Whitelist("published"))
		}
		if err == nil {
			utils.Must(tx.Commit())
			stats.TxCommitted.Inc(1)
		} else {
			utils.Must(tx.Rollback())
			stats.TxRolledBack.Inc(1)
		}
	}

	wg.Done()
}

func importVirtualLesson(exec boil.Executor, vl *kmodels.VirtualLesson) (*models.Collection, error) {
	collection, err := models.Collections(qm4.Where("(properties->>'kmedia_id')::int = ?", vl.ID)).One(exec)
	if err == nil {
		stats.CollectionsUpdated.Inc(1)
	} else {
		if err == sql.ErrNoRows {
			// Create new collection
			collection = &models.Collection{
				UID:    utils.GenerateUID(8),
				TypeID: common.CONTENT_TYPE_REGISTRY.ByName[common.CT_DAILY_LESSON].ID,
			}
			err = collection.Insert(exec, boil.Infer())
			if err != nil {
				return nil, errors.Wrapf(err, "Insert collection, virtual lesson [%d]", vl.ID)
			}
			stats.CollectionsCreated.Inc(1)
		} else {
			return nil, errors.Wrapf(err, "Lookup collection, virtual lesson [%d]", vl.ID)
		}
	}

	if vl.FilmDate.Time.Weekday() == 6 {
		collection.TypeID = common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SPECIAL_LESSON].ID
	} else {
		collection.TypeID = common.CONTENT_TYPE_REGISTRY.ByName[common.CT_DAILY_LESSON].ID
	}

	var props = make(map[string]interface{})
	if collection.Properties.Valid {
		collection.Properties.Unmarshal(&props)
	}
	props["kmedia_id"] = vl.ID
	props["film_date"] = vl.FilmDate.Time.Format("2006-01-02")
	p, _ := json.Marshal(props)
	collection.Properties = null.JSONFrom(p)

	// TODO: what to do with name and description ?
	_, err = collection.Update(exec, boil.Infer())
	if err != nil {
		return nil, errors.Wrapf(err, "Update collection, virtual lesson [%d]", vl.ID)
	}

	return collection, nil
}

func getValidContainers(exec boil.Executor, vl *kmodels.VirtualLesson) ([]*kmodels.Container, error) {
	// Fetch containers with file assets
	containers, err := vl.Containers(exec,
		qm.Where("content_type_id = ?", kmediaLessonCT.ID),
		qm.Load("FileAssets")).
		All()
	if err != nil {
		return nil, errors.Wrapf(err, "load containers from DB [%d]", vl.ID)
	}
	stats.ContainersVisited.Inc(int32(len(containers)))

	// Filter out containers without file_assets
	validContainers := containers[:0]
	for _, x := range containers {
		if len(x.R.FileAssets) > 0 {
			validContainers = append(validContainers, x)
			stats.ContainersWithFiles.Inc(1)
		} else {
			log.Warningf("Empty container [%d]", x.ID)
			stats.ContainersWithoutFiles.Inc(1)
		}
	}

	return validContainers, nil
}
