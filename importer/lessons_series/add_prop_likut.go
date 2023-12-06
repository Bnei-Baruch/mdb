package lessons_series

import (
	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strings"
)

func (ls *LessonsSeries) RunAddLikutProp() {
	mdb := ls.openDB()
	q := `SELECT DISTINCT c.id FROM collections c
INNER JOIN collections_content_units ccu ON ccu.collection_id = c.id
INNER JOIN content_units cu ON ccu.content_unit_id = cu.id
INNER JOIN content_unit_derivations cud ON cud.source_id = cu.id
INNER JOIN content_units l ON cud.derived_id = l.id
WHERE c.type_id = $1 
  AND c.properties ? 'tags'
  AND cu.type_id = $2
  AND l.type_id = $3
ORDER BY c.id`

	rows, err := queries.Raw(
		q,
		common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSONS_SERIES].ID,
		common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSON_PART].ID,
		common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LIKUTIM].ID,
	).Query(mdb)
	utils.Must(err)
	defer rows.Close()

	for rows.Next() {
		var cId int64
		utils.Must(rows.Scan(&cId))
		c, err := models.Collections(
			models.CollectionWhere.ID.EQ(cId),
			qm.Load("CollectionsContentUnits"),
			qm.Load("CollectionsContentUnits.ContentUnit"),
			qm.Load("CollectionsContentUnits.ContentUnit.SourceContentUnitDerivations"),
			qm.Load("CollectionsContentUnits.ContentUnit.SourceContentUnitDerivations.Derived"),
		).One(mdb)
		utils.Must(err)
		l := findLikutForC(c.R.CollectionsContentUnits)
		if len(l) == 0 {
			log.Errorf("not found likut for collection %d", cId)
		}

		var props map[string]interface{}
		if err := c.Properties.Unmarshal(&props); err != nil {
			log.Errorf("cant unmarshal props for collection %d", cId)
			continue
		}
		props[strings.ToLower(common.CT_LIKUTIM)] = l
		utils.Must(api.UpdateCollectionProperties(mdb, c, props))
	}
	utils.Must(rows.Err())

}

func findLikutForC(ccus models.CollectionsContentUnitSlice) []string {
	uniq := make(map[string]bool)
	for _, ccu := range ccus {
		cu := ccu.R.ContentUnit
		for _, cul := range cu.R.SourceContentUnitDerivations {
			l := cul.R.Derived
			if common.CONTENT_TYPE_REGISTRY.ByID[l.TypeID].Name != common.CT_LIKUTIM {
				continue
			}
			if _, ok := uniq[l.UID]; !ok {
				uniq[l.UID] = true
			}
		}
	}
	uids := make([]string, 0)
	for l := range uniq {
		uids = append(uids, l)
	}
	return uids
}
