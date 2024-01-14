package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strings"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/events"
	"github.com/Bnei-Baruch/mdb/models"
)

var MinCuNumberForNewLessonSeries = 2
var DaysCheckForLessonsSeries = 30
var TES_ROOT_UID = "xtKmrbb9"
var TES_PARTS_UIDS = []string{"9xNFLSSp", "XlukqLH8", "AerA1hNN", "1kDKQxJb", "o5lXptLo", "eNwJXy4s", "ahipVtPu", "Pscnn3pP", "Lfu7W3CD", "n03vXCJl", "UGcGGSpP", "NpLQT0LX", "AUArdCkH", "tit6XNAo", "FaKUG7ru", "mW6eON0z"}
var ZOAR_UID = "AwGBQX2L"
var ZOAR_PART_ONE_UID = "43BXTx3C"

var qCusByS = fmt.Sprintf(`
  SELECT DISTINCT ON(s.id) s.uid, array_agg(DISTINCT cu.id)
  FROM content_units cu
  INNER JOIN content_units_sources cus ON cu.id = cus.content_unit_id
  INNER JOIN sources s ON s.id = cus.source_id
  WHERE cu.type_id = $1 
  AND coalesce((cu.properties->>'film_date')::date, cu.created_at) > (CURRENT_DATE - '%d day'::interval)
  AND cu.published = TRUE AND cu.secure = 0
  GROUP BY  s.id`, DaysCheckForLessonsSeries)

type AssociateBySources struct {
	tx            boil.Executor
	cu            *models.ContentUnit
	evnts         []events.Event
	seriesSources map[string]bool
	cusByS        map[string][]int64
}

func (a *AssociateBySources) Associate(sUIDs []string) ([]events.Event, error) {
	a.evnts = make([]events.Event, 0)
	a.seriesSources = make(map[string]bool)

	if err := a.prepareCUs(sUIDs); err != nil {
		return nil, NewInternalError(err)
	}
	for sUid, _ := range a.seriesSources {
		if len(a.cusByS[sUid]) < MinCuNumberForNewLessonSeries {
			continue
		}
		c, err := a.findPrevCollection(sUid)
		if errors.Is(err, sql.ErrNoRows) {
			c, err = a.createCollection(sUid)
		}
		if err != nil {
			return nil, NewInternalError(err)
		}

		if err := a.attachCollection(c, sUid); err != nil {
			return nil, NewInternalError(err)
		}
	}
	return a.evnts, nil
}

func (a *AssociateBySources) prepareCUs(cuSUids []string) error {
	rows, err := queries.Raw(qCusByS, common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSON_PART].ID).Query(a.tx)
	if err != nil {
		return NewInternalError(err)
	}
	defer rows.Close()

	_cusByS := make(map[string][]int64)
	sUIDs := make([]string, 0)
	var sUid string
	var cuIdsByS pq.Int64Array
	for rows.Next() {
		err = rows.Scan(&sUid, &cuIdsByS)
		if err != nil {
			return NewInternalError(err)
		}
		sUIDs = append(sUIDs, sUid)
		_cusByS[sUid] = append(_cusByS[sUid], cuIdsByS...)
	}
	sUIDs = append(sUIDs, cuSUids...)
	sByLeaf, err := MapParentByLeaf(a.tx, sUIDs)
	if err != nil {
		return NewInternalError(err)
	}
	for _, uid := range cuSUids {
		a.seriesSources[sByLeaf[uid]] = true
	}

	a.cusByS = make(map[string][]int64)
	for _, prevUid := range sUIDs {
		fixedUid := sByLeaf[prevUid]
		for uid, _ := range a.seriesSources {
			if _, ok := a.cusByS[uid]; !ok {
				a.cusByS[uid] = []int64{}
			}
			if fixedUid == uid {
				a.cusByS[uid] = appendUniqIds(a.cusByS[uid], _cusByS[prevUid])
			}
		}
	}
	return nil
}

func (a *AssociateBySources) findStartDate(uid string) (string, error) {
	return findStartDate(a.tx, a.cusByS[uid][0])
}

func (a *AssociateBySources) findPrevCollection(uid string) (*models.Collection, error) {
	cus := a.cusByS[uid]
	c, err := models.Collections(
		models.CollectionWhere.TypeID.EQ(common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSONS_SERIES].ID),
		qm.InnerJoin("collections_content_units ccu ON ccu.collection_id = id"),
		qm.WhereIn("ccu.content_unit_id IN ?", utils.ConvertArgsInt64(cus)...),
		qm.Where("properties->>'source' = ?", uid),
		qm.OrderBy("(properties->>'end_date')::date DESC, id DESC"),
	).One(a.tx)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (a *AssociateBySources) createCollection(uid string) (*models.Collection, error) {
	startDate, err := a.findStartDate(uid)
	if err != nil {
		return nil, err
	}
	props := map[string]interface{}{
		"source":     uid,
		"start_date": startDate,
	}
	c, err := CreateCollection(a.tx, common.CT_LESSONS_SERIES, props)
	if err != nil {
		return nil, err
	}
	c.Published = true
	_, err = c.Update(a.tx, boil.Whitelist("published"))
	if err != nil {
		return nil, err
	}

	addCCUs := make([]*models.CollectionsContentUnit, 0)
	for i, id := range a.cusByS[uid] {
		ccu := &models.CollectionsContentUnit{
			ContentUnitID: id,
			CollectionID:  c.ID,
			Position:      i + 1,
			Name:          fmt.Sprintf("%d", i+1+1),
		}
		addCCUs = append(addCCUs, ccu)
	}
	if err := c.AddCollectionsContentUnits(a.tx, true, addCCUs...); err != nil {
		return nil, err
	}

	if err = DescribeCollection(a.tx, c); err != nil {
		return nil, err
	}

	a.evnts = append(a.evnts, events.CollectionCreateEvent(c))
	return c, nil
}

func (a *AssociateBySources) attachCollection(c *models.Collection, uid string) error {
	return attachCollection(a.tx, c, a.cu, a.cusByS[uid])
}

var qCusByL = `
  SELECT DISTINCT ON(lcu.id) lcu.uid, array_agg(DISTINCT cu.id)
  FROM content_units cu
  INNER JOIN content_unit_derivations dcu ON cu.id = dcu.source_id
  INNER JOIN content_units lcu ON lcu.id = dcu.derived_id
  WHERE cu.type_id = $1
  AND coalesce((cu.properties->>'film_date')::date, cu.created_at) > (CURRENT_DATE - '%d day'::interval)
  AND cu.published = TRUE AND cu.secure = 0
  AND lcu.uid IN (%s)
  GROUP BY  lcu.id
`

type AssociateByLikutim struct {
	tx     boil.Executor
	cu     *models.ContentUnit
	evnts  []events.Event
	cusByL map[string][]int64
	lUIDs  []string
}

func (a *AssociateByLikutim) Associate(lUIDs []string) ([]events.Event, error) {
	a.evnts = make([]events.Event, 0)
	a.cusByL = make(map[string][]int64)
	a.lUIDs = lUIDs
	if err := a.prepareCUs(); err != nil {
		return nil, NewInternalError(err)
	}

	for _, lUid := range a.lUIDs {
		if len(a.cusByL[lUid]) < MinCuNumberForNewLessonSeries {
			continue
		}
		c, err := a.findPrevCollection(lUid)
		if errors.Is(err, sql.ErrNoRows) {
			c, err = a.createCollection(lUid)
		}
		if err != nil {
			return nil, NewInternalError(err)
		}

		if err := a.attachCollection(c, lUid); err != nil {
			return nil, NewInternalError(err)
		}
	}

	return a.evnts, nil
}

func (a *AssociateByLikutim) prepareCUs() error {
	q := fmt.Sprintf(qCusByL, DaysCheckForLessonsSeries, fmt.Sprintf("'%s'", strings.Join(a.lUIDs, "','")))
	rows, err := queries.Raw(q, common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSON_PART].ID).Query(a.tx)

	if err != nil {
		return err
	}
	defer rows.Close()

	var lUId string
	var cuIdsByL pq.Int64Array
	var cuIds []int64
	for rows.Next() {
		err = rows.Scan(&lUId, &cuIdsByL)
		if err != nil {
			return err
		}
		a.cusByL[lUId] = cuIdsByL
		cuIds = append(cuIds, cuIdsByL...)
	}

	return nil
}

func (a *AssociateByLikutim) findStartDate(uid string) (string, error) {
	return findStartDate(a.tx, a.cusByL[uid][0])
}

func (a *AssociateByLikutim) findPrevCollection(uid string) (*models.Collection, error) {
	cus := a.cusByL[uid]
	c, err := models.Collections(
		models.CollectionWhere.TypeID.EQ(common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSONS_SERIES].ID),
		qm.InnerJoin("collections_content_units ccu ON ccu.collection_id = id"),
		qm.WhereIn("ccu.content_unit_id IN ?", utils.ConvertArgsInt64(cus)...),
		qm.Where(`properties->'likutim' @> ?`, fmt.Sprintf(`["%s"]`, uid)),
	).One(a.tx)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (a *AssociateByLikutim) createCollection(uid string) (*models.Collection, error) {
	startDate, err := a.findStartDate(uid)
	if err != nil {
		return nil, err
	}
	props := map[string]interface{}{
		"likutim":    []string{uid},
		"start_date": startDate,
	}
	likut, err := models.ContentUnits(
		models.ContentUnitWhere.UID.EQ(uid),
		qm.Load("Tags"),
	).One(a.tx)
	if err != nil {
		return nil, err
	}
	var tags []string
	if likut.R.Tags != nil {
		for _, t := range likut.R.Tags {
			tags = append(tags, t.UID)
		}
		props["tags"] = tags
	}
	c, err := CreateCollection(a.tx, common.CT_LESSONS_SERIES, props)
	if err != nil {
		return nil, err
	}
	c.Published = true
	_, err = c.Update(a.tx, boil.Whitelist("published"))
	if err != nil {
		return nil, err
	}

	addCCUs := make([]*models.CollectionsContentUnit, 0)
	for i, id := range a.cusByL[uid] {
		ccu := &models.CollectionsContentUnit{
			ContentUnitID: id,
			Position:      i + 1,
		}
		addCCUs = append(addCCUs, ccu)
	}
	if err := c.AddCollectionsContentUnits(a.tx, true, addCCUs...); err != nil {
		return nil, err
	}

	if err := I18nFromCU(a.tx, uid, c); err != nil {
		return nil, err
	}
	a.evnts = append(a.evnts, events.CollectionCreateEvent(c))
	return c, nil
}

func (a *AssociateByLikutim) attachCollection(c *models.Collection, uid string) error {
	return attachCollection(a.tx, c, a.cu, a.cusByL[uid])
}

// Helpers

func MapParentByLeaf(exec boil.Executor, uids []string) (map[string]string, error) {
	q := fmt.Sprintf(`
WITH RECURSIVE recurcive_s(id, uid, parent_id, start_uid) AS(
	SELECT id, uid, parent_id, uid
		FROM sources where uid IN (%s)
	UNION
	SELECT s.id, s.uid, s.parent_id, rs.start_uid
		FROM recurcive_s rs, sources s WHERE rs.parent_id = s.id AND rs.uid != '%s'
)
SELECT start_uid, uid FROM recurcive_s WHERE uid IN (%s)
`,
		fmt.Sprintf("'%s'", strings.Join(uids, "','")),
		ZOAR_PART_ONE_UID,
		fmt.Sprintf("'%s'", strings.Join(append(TES_PARTS_UIDS, ZOAR_UID, ZOAR_PART_ONE_UID), "','")),
	)
	rows, err := queries.Raw(q).Query(exec)
	if err != nil {
		return nil, NewInternalError(err)
	}
	defer rows.Close()
	var uid string
	var nUid string
	newByOldUid := make(map[string]string)
	for rows.Next() {
		err = rows.Scan(&uid, &nUid)
		if err != nil {
			return nil, NewInternalError(err)
		}
		//if we get more then one parent uid its mean that its zoar and its child ZOAR_PART_ONE_UID.
		// on this mode need take ZOAR_PART_ONE_UID
		if _, ok := newByOldUid[uid]; !ok {
			newByOldUid[uid] = nUid
		}
	}

	for _, uid := range uids {
		if _, ok := newByOldUid[uid]; !ok {
			newByOldUid[uid] = uid
		}
	}

	return newByOldUid, nil
}

func appendUniqIds(ids1, ids2 []int64) []int64 {
	for _, id2 := range ids2 {
		isIn := false
		for _, id1 := range ids1 {
			if id1 == id2 {
				isIn = true
				break
			}
		}
		if !isIn {
			ids1 = append(ids1, id2)
		}
	}
	return ids1
}

func findStartDate(tx boil.Executor, id int64) (string, error) {
	var _props map[string]interface{}
	cu, err := models.FindContentUnit(tx, id)
	if err != nil {
		return "", err
	}
	if err := cu.Properties.Unmarshal(&_props); err != nil {
		return "", err
	}
	fd, ok := _props["film_date"]
	if !ok {
		return "", NewInternalError(errors.New("no film date"))
	}
	return fd.(string), nil
}

func attachCollection(tx boil.Executor, c *models.Collection, cu *models.ContentUnit, cus []int64) error {
	prevCCU, err := models.FindCollectionsContentUnit(tx, c.ID, cus[len(cus)-1])
	if err != nil {
		return err
	}
	//position start from 0 when Name from 1
	ccu := &models.CollectionsContentUnit{
		ContentUnitID: cu.ID,
		CollectionID:  c.ID,
		Position:      prevCCU.Position + 1,
		Name:          fmt.Sprintf("%d", prevCCU.Position+1+1),
	}
	if err := c.AddCollectionsContentUnits(tx, true, ccu); err != nil {
		return err
	}

	var cuProps map[string]interface{}
	if err := cu.Properties.Unmarshal(&cuProps); err != nil {
		return err
	}
	if err := UpdateCollectionProperties(tx, c, map[string]interface{}{"end_date": cuProps["film_date"]}); err != nil {
		return err
	}

	return nil
}

func I18nFromCU(tx boil.Executor, uid string, c *models.Collection) error {
	cusI18ns, err := models.ContentUnits(qm.Where("uid = ?", uid), qm.Load("ContentUnitI18ns")).One(tx)
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
		if err = c.AddCollectionI18ns(tx, true, i18ns...); err != nil {
			return err
		}
	}
	return nil
}
