package cusource

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"io/ioutil"
)

func Validator() {
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	defer mdb.Close()
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	boil.DebugMode = true
	utils.Must(common.InitTypeRegistries(mdb))

	sDB, err := models.Sources(mdb).All()
	utils.Must(err)
	for _, s := range sDB {
		compareDirWithDB(mdb, s.UID)
	}

	fmt.Print("validation is successful")
}

func compareDirWithDB(mdb *sql.DB, uid string) {
	path := viper.GetString("source-import.source-dir") + "/" + uid

	cu, err := models.ContentUnits(mdb,
		qm.InnerJoin("content_units_sources cus ON cus.content_unit_id = id"),
		qm.Where("cus.source_id = ? AND type_id = ?", uid, common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SOURCE].ID),
	).One()

	if err != nil || cu == nil {
		panic(errors.New(fmt.Sprintf("Not fount CU for Source: %s", uid)))
	}

	filesOS, err := ioutil.ReadDir(path)
	if err != nil {
		panic(errors.New(fmt.Sprintf("No folder for Source: %s on the path: %s", uid, path)))
	}
	filesDB, err := models.Files(mdb,
		qm.Where("content_unit_id = ?", cu.ID),
	).All()

	if len(filesDB) != len(filesOS) {
		panic(errors.New("not equal arrays length"))
	}

	for _, f := range filesDB {
		eq := false
		for _, fos := range filesOS {
			if fos.Name() == f.Name {
				eq = true
			}
		}
		if !eq {
			panic(errors.New(fmt.Sprintf("No file with name %s in dir: %s", f.Name, path)))
		}
	}
}
