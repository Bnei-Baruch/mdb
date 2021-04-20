package cusource

import (
	"database/sql"
	"encoding/json"
	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/volatiletech/null.v6"
	"time"
)

func BuildCUSources(mdb *sql.DB) ([]*models.Source, []*models.ContentUnit) {

	rows, err := queries.Raw(mdb,
		`SELECT cus.source_id FROM content_units cu
			INNER JOIN content_units_sources cus ON cus.content_unit_id = cu.id
			WHERE cu.type_id = $1`,
		common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SOURCE].ID,
	).Query()

	utils.Must(err)
	defer rows.Close()
	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		utils.Must(err)
		ids = append(ids, id)
	}
	mods := make([]qm.QueryMod, 0)
	if len(ids) > 0 {
		mods = append(mods, qm.WhereIn("id NOT IN ?", utils.ConvertArgsInt64(ids)...))
	}

	sources, err := models.Sources(mdb, mods...).All()
	utils.Must(err)

	for _, s := range sources {
		utils.Must(err)
		_, err := createCU(s, mdb)
		if err != nil {
			log.Debug("Duplicate create CU", err)
		}
	}
	return sources, nil
}

func createCU(s *models.Source, mdb boil.Executor) (*models.ContentUnit, error) {
	cuUid := s.UID
	hasCU, err := models.ContentUnits(mdb, qm.Where("uid = ?", cuUid)).Exists()
	if err != nil {
		return nil, api.NewInternalError(err)
	}
	if hasCU {
		cuUid, err = api.GetFreeUID(mdb, new(api.ContentUnitUIDChecker))
		if err != nil {
			return nil, api.NewInternalError(err)
		}
	}

	props := make(map[string]interface{})
	props["source_id"] = s.UID
	p, _ := json.Marshal(props)
	cu := &models.ContentUnit{
		UID:        cuUid,
		TypeID:     common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SOURCE].ID,
		Secure:     common.SEC_PUBLIC,
		Published:  true,
		Properties: null.JSONFrom(p),
		CreatedAt:  time.Now(),
	}

	err = cu.Insert(mdb)
	if err != nil {
		return nil, err
	}

	err = cu.AddSources(mdb, false, s)
	utils.Must(err)
	return cu, nil
}

func haveSourceCUTypeSource(suid string, mdb boil.Executor) bool {
	tid := common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SOURCE].ID
	has, err := models.ContentUnits(
		mdb,
		qm.Where("content_units.type_id = ? AND s.uid = ?", tid, suid),
		qm.InnerJoin("content_units_sources as cus ON cus.content_unit_id = content_units.id"),
		qm.InnerJoin("sources as s ON cus.source_id = s.id"),
	).Exists()
	utils.Must(err)
	return has
}
