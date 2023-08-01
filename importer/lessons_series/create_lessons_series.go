package lessons_series

import (
	"database/sql"
	"fmt"
	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/events"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"time"
)

var BATCH_DURATION = 2 * 7 * 24 * time.Hour

type LessonsSeries struct {
	from     time.Time
	tx       *sql.Tx
	bySource map[string]*bySourceItem
}

type bySourceItem struct {
	from       time.Time
	cus        []*models.ContentUnit
	source     string
	likut      *models.ContentUnit
	collection *models.Collection
}

func (ls *LessonsSeries) Run() {
	mdb := ls.openDB()
	defer mdb.Close()

	//boil.DebugMode = true
	ls.bySource = map[string]*bySourceItem{}
	//ls.from, _ = time.Parse("2006-01-02", "1999-01-01")
	ls.from, _ = time.Parse("2006-01-02", "1980-01-01")

	for ls.from.Before(time.Now()) {
		tx, err := mdb.Begin()
		utils.Must(err)
		ls.tx = tx
		//boil.DebugMode = true
		cus, err := ls.fetchBatch()
		//boil.DebugMode = false
		utils.Must(err)
		ls.group(cus)
		utils.Must(ls.save())
		mustConcludeTx(ls.tx, err)
		log.Printf("Batch from: %v was ended", ls.from)
		log.Printf("bySource length %d", len(ls.bySource))
		ls.from = ls.from.Add(BATCH_DURATION)
	}
}

func (ls *LessonsSeries) fetchBatch() ([]*models.ContentUnit, error) {
	return models.ContentUnits(
		qm.Load("CollectionsContentUnits"),
		qm.Load("CollectionsContentUnits.Collection"),
		qm.Load("Sources"),
		qm.Load("SourceContentUnitDerivations"),
		qm.Load("SourceContentUnitDerivations.Derived"),
		qm.Load("SourceContentUnitDerivations.Derived.Tags"),
		qm.Where("type_id = ? AND published = TRUE AND secure = 0", common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSON_PART].ID),
		qm.And("(properties->>'film_date')::date > ?::date", ls.from),
		qm.And("(properties->>'film_date')::date <= ?::date", ls.from.Add(BATCH_DURATION)),
		qm.OrderBy("id ASC"),
	).All(ls.tx)
}

func (ls *LessonsSeries) group(cus []*models.ContentUnit) {
	for _, cu := range cus {
		bySource, err := ls.groupBySource(cu)
		utils.Must(err)

		for key, item := range ls.groupByLikutim(cu) {
			bySource[key] = item
		}

		ls.attachCollection(cu, bySource)

		for key, item := range bySource {
			ls.mergeBySourceItem(item, key)
		}
	}
}

func (ls *LessonsSeries) mergeBySourceItem(item *bySourceItem, key string) {
	if _, ok := ls.bySource[key]; !ok {
		ls.bySource[key] = item
		return
	}
	if item.collection != nil {
		ls.bySource[key].collection = item.collection
	}
	ls.bySource[key].cus = append(ls.bySource[key].cus, item.cus...)
	ls.bySource[key].from = item.from
}

func (ls *LessonsSeries) groupBySource(cu *models.ContentUnit) (map[string]*bySourceItem, error) {
	bySource := map[string]*bySourceItem{}

	sUids := make([]string, len(cu.R.Sources))
	for i, s := range cu.R.Sources {
		sUids[i] = s.UID
	}
	sByUids, err := api.MapParentByLeaf(ls.tx, sUids)
	if err != nil {
		return nil, err
	}
	for _, uid := range sByUids {
		if _, ok := bySource[uid]; !ok {
			bySource[uid] = &bySourceItem{cus: make([]*models.ContentUnit, 0), source: uid}
		}

		bySource[uid].from = ls.from
		bySource[uid].cus = append(bySource[uid].cus, cu)
	}
	return bySource, nil
}

func (ls *LessonsSeries) groupByLikutim(cu *models.ContentUnit) map[string]*bySourceItem {
	bySource := map[string]*bySourceItem{}
	for _, dcu := range cu.R.SourceContentUnitDerivations {
		if dcu.R.Derived.TypeID != common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LIKUTIM].ID {
			continue
		}
		if _, ok := bySource[dcu.R.Derived.UID]; !ok {
			bySource[dcu.R.Derived.UID] = &bySourceItem{cus: make([]*models.ContentUnit, 0), likut: dcu.R.Derived}
		}

		bySource[dcu.R.Derived.UID].from = ls.from
		bySource[dcu.R.Derived.UID].cus = append(bySource[dcu.R.Derived.UID].cus, cu)
	}
	return bySource
}

func (ls *LessonsSeries) attachCollection(cu *models.ContentUnit, bySource map[string]*bySourceItem) {
	for _, ccu := range cu.R.CollectionsContentUnits {
		if ccu.R.Collection.TypeID != common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSONS_SERIES].ID {
			continue
		}
		var props map[string]interface{}
		if err := ccu.R.Collection.Properties.Unmarshal(&props); err != nil {
			continue
		}
		for uid, item := range bySource {
			if props["source"].(string) == uid {
				item.collection = ccu.R.Collection
				log.Printf("attach to old collection cu_id - %d, c_id - %d", cu.ID, ccu.CollectionID)
			}
			if lProp, ok := props["likutim"]; ok {
				for _, l := range lProp.([]interface{}) {
					if l.(string) == uid {
						item.collection = ccu.R.Collection
						log.Printf("attach to old collection cu_id - %d, c_id - %d", cu.ID, ccu.CollectionID)
					}
				}
			}
		}
	}
}

func (ls *LessonsSeries) save() error {
	evnts := make([]events.Event, 0)
	var err error
	for sUid, item := range ls.bySource {
		//if have hand created collection that cu was missed we do nothing
		if item.collection != nil && item.collection.CreatedAt.Before(time.Now().Add(-1*24*time.Hour)) {
			continue
		}
		_evnts, err := ls.saveCollection(item, sUid)
		if err != nil {
			return err
		}
		evnts = append(evnts, _evnts...)

		_evnts, err = ls.saveCCUs(item)
		if err != nil {
			return err
		}
		evnts = append(evnts, _evnts...)

		if item.from.Before(ls.from.Add(-1 * BATCH_DURATION)) {
			delete(ls.bySource, sUid)
		}
	}
	return err
}
func (ls *LessonsSeries) saveCollection(item *bySourceItem, sUid string) ([]events.Event, error) {
	evnts := make([]events.Event, 0)
	var err error
	var c *models.Collection

	var end_props map[string]interface{}
	if err := item.cus[len(item.cus)-1].Properties.Unmarshal(&end_props); err != nil {
		return nil, err
	}
	if item.collection != nil {
		c = item.collection
		props := map[string]interface{}{
			"source":   sUid,
			"end_date": end_props["film_date"],
		}
		if err := api.UpdateCollectionProperties(ls.tx, c, props); err != nil {
			return nil, err
		}
	} else if len(item.cus) > api.MIN_CU_NUMBER_FOR_NEW_LESSON_SERIES {
		if c, err = ls.saveNewCollection(item, sUid, end_props["film_date"].(string)); err != nil {
			return nil, err
		}
	}
	item.collection = c
	return evnts, nil
}

func (ls *LessonsSeries) saveNewCollection(item *bySourceItem, sUid string, end_date string) (*models.Collection, error) {
	var startProps map[string]interface{}
	if err := item.cus[0].Properties.Unmarshal(&startProps); err != nil {
		return nil, err
	}
	props := map[string]interface{}{
		"end_date":   end_date,
		"start_date": startProps["film_date"],
	}
	if item.likut != nil {
		props["likutim"] = []string{sUid}
		var tags []string
		/*err := item.likut.L.LoadTags(ls.tx, true, item.likut, nil)
		if err != nil {
			return nil, err
		}*/
		if item.likut.R.Tags != nil {
			for _, t := range item.likut.R.Tags {
				tags = append(tags, t.UID)
			}
			props["tags"] = tags
		}
	} else {
		props["source"] = sUid
	}
	c, err := api.CreateCollection(ls.tx, common.CT_LESSONS_SERIES, props)
	if err != nil {
		return nil, err
	}
	if item.likut != nil {
		utils.Must(ls.i18nFromCU(sUid, c))
	} else {
		utils.Must(ls.i18nFromSource(sUid, c))
	}
	c.Published = true
	_, err = c.Update(ls.tx, boil.Whitelist("published"))
	return c, err
}

func (ls *LessonsSeries) saveCCUs(item *bySourceItem) ([]events.Event, error) {
	evnts := make([]events.Event, 0)
	if item.collection == nil {
		return evnts, nil
	}
	ccus := make([]*models.CollectionsContentUnit, 0)
	for i, cu := range item.cus {
		err := cu.L.LoadCollectionsContentUnits(ls.tx, true, cu, nil)
		if err != nil {
			return nil, err
		}
		has := false
		for _, _ccu := range cu.R.CollectionsContentUnits {
			if _ccu.CollectionID == item.collection.ID {
				has = true
			}
		}
		if has {
			continue
		}
		ccu := &models.CollectionsContentUnit{ContentUnitID: cu.ID, Position: i, Name: fmt.Sprint(i)}
		ccus = append(ccus, ccu)
	}

	if err := item.collection.AddCollectionsContentUnits(ls.tx, true, ccus...); err != nil {
		return nil, err
	}
	if _, err := item.collection.Update(ls.tx, boil.Infer()); err != nil {
		return nil, err
	}
	evnts = append(evnts, events.CollectionContentUnitsChangeEvent(item.collection))

	return evnts, nil
}

func (ls *LessonsSeries) i18nFromCU(uid string, c *models.Collection) error {
	cusI18ns, err := models.ContentUnits(qm.Where("uid = ?", uid), qm.Load("ContentUnitI18ns")).One(ls.tx)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if cusI18ns != nil {
		i18ns := make([]*models.CollectionI18n, 0)
		for _, sI18n := range cusI18ns.R.ContentUnitI18ns {
			i18n := &models.CollectionI18n{
				Language: sI18n.Language,
				Name:     sI18n.Name,
			}
			i18ns = append(i18ns, i18n)
		}
		if err = c.AddCollectionI18ns(ls.tx, true, i18ns...); err != nil {
			return err
		}
	}
	return nil
}

func (ls *LessonsSeries) i18nFromSource(uid string, c *models.Collection) error {
	cusI18ns, err := models.Sources(qm.Where("uid = ?", uid), qm.Load("SourceI18ns")).One(ls.tx)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if cusI18ns != nil {
		i18ns := make([]*models.CollectionI18n, 0)
		for _, sI18n := range cusI18ns.R.SourceI18ns {
			i18n := &models.CollectionI18n{
				Language: sI18n.Language,
				Name:     sI18n.Name,
			}
			i18ns = append(i18ns, i18n)
		}
		if err = c.AddCollectionI18ns(ls.tx, true, i18ns...); err != nil {
			return err
		}
	}
	return nil

}

func (ls *LessonsSeries) openDB() *sql.DB {
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	utils.Must(common.InitTypeRegistries(mdb))
	return mdb
}

func mustConcludeTx(tx *sql.Tx, err error) {
	if err == nil {
		utils.Must(tx.Commit())
	} else {
		utils.Must(tx.Rollback())
	}
}
