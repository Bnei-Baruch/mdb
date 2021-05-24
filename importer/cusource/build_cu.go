package cusource

import (
	"database/sql"
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/volatiletech/null.v6"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func BuildCUSources(mdb *sql.DB) ([]*models.Source, []*models.ContentUnit) {

	rows, err := queries.Raw(mdb,
		`SELECT cu.properties->>'source_id' FROM content_units cu WHERE cu.type_id = $1`,
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
		mods = append(mods, qm.WhereIn("uid NOT IN ?", utils.ConvertArgsInt64(ids)...))
	}

	sources, err := models.Sources(mdb, mods...).All()
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
	props, _ := json.Marshal(map[string]string{"source_id": s.UID})
	cu := &models.ContentUnit{
		UID:        s.UID,
		TypeID:     common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SOURCE].ID,
		Published:  false,
		Properties: null.JSONFrom(props),
	}

	err := cu.Insert(mdb)
	if err != nil {
		return nil, err
	}
	return cu, nil
}
