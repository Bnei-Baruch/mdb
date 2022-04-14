package cusource

import (
	"database/sql"

	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/utils"
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
