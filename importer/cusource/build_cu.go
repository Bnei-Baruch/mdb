package cusource

import (
	"database/sql"
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func BuildCUSources(mdb *sql.DB) ([]*models.Source, []*models.ContentUnit) {

	rows, err := queries.Raw(
		`SELECT cu.properties->>'source_id' FROM content_units cu WHERE cu.type_id = $1`,
		common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SOURCE].ID,
	).Query(mdb)

	utils.Must(err)
	defer rows.Close()
	uids := make([]string, 0)
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		utils.Must(err)
		uids = append(uids, id)
	}
	mods := make([]qm.QueryMod, 0)
	if len(uids) > 0 {
		mods = append(mods, qm.WhereIn("uid NOT IN ?", utils.ConvertArgsString(uids)...))
	}

	sources, err := models.Sources(mods...).All(mdb)
	utils.Must(err)

	for _, s := range sources {
		isParent := false
		for _, sl := range sources {
			if sl.ParentID.Int64 == s.ID {
				isParent = true
			}
		}
		if isParent {
			continue
		}
		_, err := createCU(s, mdb)
		if err != nil {
			log.Debug("Duplicate create CU", err)
		}
	}
	return sources, nil
}

func createCU(s *models.Source, mdb boil.Executor) (*models.ContentUnit, error) {
	props, _ := json.Marshal(map[string]string{"source_id": s.UID, "film_date": "1980-01-01"})
	cu := &models.ContentUnit{
		UID:        s.UID,
		TypeID:     common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SOURCE].ID,
		Published:  true,
		Properties: null.JSONFrom(props),
	}

	err := cu.Insert(mdb, boil.Infer())
	if err != nil {
		return nil, err
	}
	return cu, nil
}
