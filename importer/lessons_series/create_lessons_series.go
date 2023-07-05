package lessons_series

import (
	"database/sql"
	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/events"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var BATCH_SIZE = 100

type LessonsSeries struct {
	batch    []*models.ContentUnit
	tx       *sql.Tx
	bySource map[string]*bySourceItem
}

type bySourceItem struct {
	batchIdx   int
	cus        []*models.ContentUnit
	source     string
	likut      *models.ContentUnit
	collection *models.Collection
}

func (ls *LessonsSeries) Run() {
	mdb := ls.openDB()
	defer mdb.Close()
	ls.bySource = map[string]*bySourceItem{}

	total, err := models.ContentUnits(
		qm.Where("type_id = ? AND published = TRUE AND secure = 0", common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSON_PART].ID),
	).Count(ls.tx)
	utils.Must(err)
	numBatches := int(total) / BATCH_SIZE
	for i := 0; i < numBatches; i++ {
		cus, err := ls.fetchBatch(i)
		utils.Must(err)
		ls.group(cus, i)
		utils.Must(ls.save(i))
	}
	mustConcludeTx(ls.tx, err)
}

func (ls *LessonsSeries) fetchBatch(i int) ([]*models.ContentUnit, error) {
	return models.ContentUnits(
		qm.Load("CollectionsContentUnits"),
		qm.Load("CollectionsContentUnits.Collection"),
		qm.Load("Sources"),
		qm.Load("SourceContentUnitDerivations"),
		qm.Load("SourceContentUnitDerivations.Derived"),
		qm.Where("type_id = ? AND published = TRUE AND secure = 0", common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSON_PART].ID),
		qm.OrderBy("id ASC"),
		qm.Limit(BATCH_SIZE),
		qm.Offset(BATCH_SIZE*i),
	).All(ls.tx)
}

func (ls *LessonsSeries) group(cus []*models.ContentUnit, batchIdx int) {
	for _, cu := range cus {
		bySource, err := ls.groupBySource(cu, batchIdx)
		utils.Must(err)

		for key, item := range ls.groupByLikutim(cu, batchIdx) {
			bySource[key] = item
		}

		ls.attachCollection(cu, bySource)

		for key, item := range bySource {
			ls.bySource[key] = item
		}
	}
}

func (ls *LessonsSeries) groupBySource(cu *models.ContentUnit, batchIdx int) (map[string]*bySourceItem, error) {
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

		bySource[uid].batchIdx = batchIdx
		bySource[uid].cus = append(bySource[uid].cus, cu)
	}
	return bySource, nil
}

func (ls *LessonsSeries) groupByLikutim(cu *models.ContentUnit, batchIdx int) map[string]*bySourceItem {
	bySource := map[string]*bySourceItem{}
	for _, dcu := range cu.R.SourceContentUnitDerivations {
		if dcu.R.Derived.TypeID != common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LIKUTIM].ID {
			continue
		}
		if _, ok := bySource[dcu.R.Derived.UID]; !ok {
			bySource[dcu.R.Derived.UID] = &bySourceItem{cus: make([]*models.ContentUnit, 0), likut: dcu.R.Derived}
		}

		bySource[dcu.R.Derived.UID].batchIdx = batchIdx
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
			}
		}
	}
}

func (ls *LessonsSeries) save(idx int) error {
	evnts := make([]events.Event, 0)
	var err error
	for sUid, item := range ls.bySource {
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

		if item.batchIdx-1 > idx {
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
	if err := item.cus[0].Properties.Unmarshal(&end_props); err != nil {
		return nil, err
	}
	if item.collection != nil {

		var start_props map[string]interface{}
		if err := item.cus[len(item.cus)-1].Properties.Unmarshal(&start_props); err != nil {
			return nil, err
		}
		c = item.collection
		props := map[string]interface{}{
			"source":     sUid,
			"end_date":   end_props["film_date"],
			"start_date": start_props["film_date"],
		}

		if err := api.UpdateCollectionProperties(ls.tx, c, props); err != nil {
			return nil, err
		}
	} else if len(item.cus) > api.MIN_CU_NUMBER_FOR_NEW_LESSON_SERIES {
		props := map[string]interface{}{
			"source":   sUid,
			"end_date": end_props["film_date"],
		}
		c, err = api.CreateCollection(ls.tx, common.CT_LESSONS_SERIES, props)
		if err != nil {
			return nil, err
		}
		c.Published = true
		_, err = c.Update(ls.tx, boil.Whitelist("published"))
		if err != nil {
			return nil, err
		}
	}
	item.collection = c
	return evnts, nil
}
func (ls *LessonsSeries) saveCCUs(item *bySourceItem) ([]events.Event, error) {
	evnts := make([]events.Event, 0)
	if item.collection == nil {
		return evnts, nil
	}
	ccus := make([]*models.CollectionsContentUnit, 0)
	for i, cu := range item.cus {
		has := false
		for _, _ccu := range cu.R.CollectionsContentUnits {
			if _ccu.CollectionID == item.collection.ID {
				has = true
			}
		}
		if has {
			continue
		}
		ccu := &models.CollectionsContentUnit{ContentUnitID: cu.ID, Position: i}
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

func (ls *LessonsSeries) openDB() *sql.DB {
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	boil.DebugMode = true
	utils.Must(common.InitTypeRegistries(mdb))
	tx, err := mdb.Begin()
	utils.Must(err)
	ls.tx = tx
	return mdb
}

func mustConcludeTx(tx *sql.Tx, err error) {
	if err == nil {
		utils.Must(tx.Commit())
	} else {
		utils.Must(tx.Rollback())
	}
}
