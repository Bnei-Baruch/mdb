package cusource

import (
	"database/sql"
	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func InitBuildCUSources() {
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	defer mdb.Close()
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	boil.DebugMode = true
	utils.Must(common.InitTypeRegistries(mdb))

	BuildCUSources(mdb)
}

func BuildCUSources(mdb *sql.DB) ([]*models.Source, []*models.ContentUnit) {
	var cus []*models.ContentUnit
	withCU, err := models.Sources(mdb,
		//qm.Select("id", "uid", "properties", ),
		qm.InnerJoin("content_units_sources cus ON id = cus.source_id"),
		qm.InnerJoin("content_units cu ON cus.content_unit_id = cu.id"),
		qm.Where("cu.type_id = ?", common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SOURCE].ID),
	).All()
	utils.Must(err)

	sources, err := filterRoots(mdb)
	utils.Must(err)

	for _, s := range sources {
		hasCU := false
		for _, sCU := range withCU {
			if s.ID == sCU.ID {
				hasCU = true
				break
			}
		}
		if !hasCU {
			cuUid, err := api.GetFreeUID(mdb, new(api.ContentUnitUIDChecker))
			utils.Must(err)
			cu, err := common.CreateCUTypeSource(s, mdb, cuUid)
			if err != nil {
				log.Debug("Duplicate create CU", err)
			}
			cus = append(cus, cu)
		}
	}
	return sources, cus
}

func filterRoots(mdb *sql.DB) ([]*models.Source, error) {
	all, err := models.Sources(mdb).All()

	if err != nil {
		return nil, err
	}
	var r []*models.Source
	for _, s := range all {
		isLeaf := true
		for _, l := range all {
			if s.ID == l.ParentID.Int64 {
				isLeaf = false
				break
			}
		}
		if isLeaf {
			r = append(r, s)
		}
	}
	return r, nil
}
