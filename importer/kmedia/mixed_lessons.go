package kmedia

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Bnei-Baruch/mdb/common"
	"sort"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/emirpasic/gods/maps/treemap"
	godsutil "github.com/emirpasic/gods/utils"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

// 4728	lessons-part
// 3630	lessons-part/lesson_preparation
// 3629	lessons-part/lesson_first-part
// 3631	lessons-part/lesson_second-part
// 3632	lessons-part/lesson_third-part
// 4020	lessons-part/lesson_fourth-part
// 4541	lessons-part/lesson_fifth-part
// 4841	lessons-part/lesson_six-part
// 4862	lessons-part/lesson_sixth-part
var LessonPartsCatalogs = map[int]int{
	4728: -1,
	3630: 0,
	3629: 1,
	3631: 2,
	3632: 3,
	4020: 4,
	4541: 5,
	4841: 6,
	4862: 6,
}

var cuMapByCNID map[int]*models.ContentUnit

func ImportMixedLessons() {
	clock := Init()

	var err error
	var cnMap map[int]*kmodels.Container

	log.Info("Loading all mixed lessons containers")
	cnMap, cuMapByCNID, err = loadContainersInCatalogsAndCUs(
		11, 6932, 6933, 4772, 4541, 3629, 4020, 4700, 3630, 3631, 4841, 4862, 4728, 4016, 4761, 3632)
	utils.Must(err)
	log.Infof("Got %d containers", len(cnMap))

	// Process lessons
	stats = NewImportStatistics()

	log.Info("Setting up workers")
	jobs := make(chan *kmodels.Container, 100)
	reconcileJobs := make(chan *kmodels.Container, 100)
	var workersWG sync.WaitGroup
	for w := 1; w <= 5; w++ {
		workersWG.Add(1)
		go mixedLessonWorker(jobs, reconcileJobs, &workersWG)
	}

	log.Infof("Setting up lesson reconciler")
	var reconcilerWG sync.WaitGroup
	reconcilerWG.Add(1)
	go lessonReconciler(reconcileJobs, &reconcilerWG)

	log.Info("Queueing work")
	for _, v := range cnMap {
		jobs <- v
	}

	log.Info("Closing jobs channel")
	close(jobs)

	log.Info("Waiting for workers to finish")
	workersWG.Wait()

	log.Info("Closing reconciliation jobs channel")
	close(reconcileJobs)

	log.Info("Waiting for reconciler to finish")
	reconcilerWG.Wait()

	stats.dump()

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func mixedLessonWorker(jobs <-chan *kmodels.Container, reconcileJobs chan *kmodels.Container, wg *sync.WaitGroup) {
	for cn := range jobs {
		stats.ContainersProcessed.Inc(1)

		if unit, ok := cuMapByCNID[cn.ID]; ok {
			if err := analyzeExistingMixedContainer(cn, unit); err != nil {
				log.Errorf("Analyze existing %d: %s", cn.ID, err.Error())
			}
		} else {
			if shouldReconcile, err := analyzeNewMixedContainer(cn); err != nil {
				log.Errorf("Analyze new %d: %s", cn.ID, err.Error())
			} else if shouldReconcile {
				reconcileJobs <- cn
			}
		}
	}

	wg.Done()
}

func analyzeExistingMixedContainer(cn *kmodels.Container, cu *models.ContentUnit) error {
	stats.ContentUnitsUpdated.Inc(1)
	return nil
}

func analyzeNewMixedContainer(cn *kmodels.Container) (bool, error) {
	stats.ContentUnitsCreated.Inc(1)
	//log.Infof("New container\t%d\t%s", cn.ID, cn.Name.String)

	if cn.VirtualLessonID.Valid && cn.VirtualLessonID.Int != 0 {
		log.Warnf("New container %d %s has virtual_lesson_id %d", cn.ID, cn.Name.String, cn.VirtualLessonID.Int)

		c, err := models.Collections(mdb,
			qm.Where("(properties->>'kmedia_id')::int = ?", cn.VirtualLessonID.Int)).
			One()
		if err != nil {
			if err == sql.ErrNoRows {
				log.Warnf("No collection with kmedia_id %d", cn.VirtualLessonID.Int)
			} else {
				return false, errors.Wrapf(err, "Load Collection kmedia_id %d", cn.VirtualLessonID.Int)
			}
		} else {
			log.Warnf("Collection %d missing unit container_id %d", c.ID, cn.ID)
		}

		return false, nil
	}

	return cn.ContentTypeID.Int == 4, nil

	//err := cn.L.LoadCatalogs(kmdb, true, cn)
	//if err != nil {
	//	return false, errors.Wrapf(err, "Load catalogs cn.ID %d", cn.ID)
	//}
	//
	//isLessonPart := false
	//for i := range cn.R.Catalogs {
	//	c := cn.R.Catalogs[i]
	//	if _, ok := LessonPartsCatalogs[c.ID]; ok {
	//		isLessonPart = true
	//		break
	//	}
	//}
	//
	//if !isLessonPart {
	//	log.Infof("New non-lesson-part container\t%d\t%s", cn.ID, cn.Name.String)
	//}
	//
	//return isLessonPart, nil
}

func lessonReconciler(jobs <-chan *kmodels.Container, wg *sync.WaitGroup) {
	containers := make([]*kmodels.Container, 0)

	skipped := 0
	for cn := range jobs {
		//log.Infof("Reconcile %d\t%s\t%s", cn.ID, cn.Filmdate.Time.Format("2006-01-02"), cn.Name.String)

		skip := !cn.Filmdate.Valid ||
			strings.HasPrefix(cn.Name.String, "maamar_zohoraim") ||
			strings.HasPrefix(cn.Name.String, "ML_Maamar_zohoraim") ||
			strings.HasSuffix(cn.Name.String, "mzohoraim_bb") ||
			strings.HasPrefix(cn.Name.String, "UlpanIvrit") ||
			!strings.Contains(cn.Name.String, cn.Filmdate.Time.Format("2006-01-02"))

		if skip {
			skipped++
		} else {
			containers = append(containers, cn)
		}
	}

	log.Infof("%d containers needs reconciliation", len(containers))
	log.Infof("%d skipped", skipped)

	sort.Slice(containers, func(i, j int) bool {
		a := containers[i]
		b := containers[j]

		if a.Filmdate.Time.Before(b.Filmdate.Time) {
			return true
		} else if b.Filmdate.Time.Before(a.Filmdate.Time) {
			return false
		}

		return a.Name.String < b.Name.String
	})

	byDate := treemap.NewWith(godsutil.TimeComparator)
	for i := range containers {
		cn := containers[i]

		k := cn.Filmdate.Time
		v, ok := byDate.Get(k)
		if !ok {
			v = make([]*kmodels.Container, 0)
		}
		byDate.Put(k, append(v.([]*kmodels.Container), cn))
	}

	byDate.Each(func(k, v interface{}) {
		filmDate := k.(time.Time).Format("2006-01-02")
		log.Infof("%s\t%d", filmDate, len(v.([]*kmodels.Container)))
		for _, cn := range v.([]*kmodels.Container) {
			log.Infof("%d\t%s", cn.ID, cn.Name.String)
		}

		cus, err := models.ContentUnits(mdb,
			qm.Where("type_id = ?", 1),
			qm.Where("properties->>'film_date' = ?", filmDate),
			qm.Load("ContentUnitI18ns"),
		).
			All()
		if err != nil {
			log.Errorf("Load CUs in %s: %s", filmDate, err.Error())
			return
		}
		log.Infof("\t\tCUs in %s [%d]", filmDate, len(cus))
		for i := range cus {
			cu := cus[i]

			name := ""
			for i := range cu.R.ContentUnitI18ns {
				i18n := cu.R.ContentUnitI18ns[i]
				if i18n.Language == common.LANG_HEBREW {
					name = utils.Reverse(i18n.Name.String)
				}
			}

			if cu.Properties.Valid {
				var props map[string]interface{}
				if err := json.Unmarshal(cu.Properties.JSON, &props); err != nil {
					log.Errorf("json.Unmarshal CU properties %d", cu.ID)
				} else {
					if kmid, ok := props["kmedia_id"]; ok {
						name = fmt.Sprintf("%s\t%d", name, int(kmid.(float64)))
					}
				}
			}

			log.Infof("\t\t\t%d\t%s", cu.ID, name)
		}
	})

	wg.Done()
}
