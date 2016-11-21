package dal

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func Init() bool {
	url := viper.GetString("mdb.url")
	db, err := sql.Open("postgres", url)
	checkErr(err)
	defer db.Close()
	return true
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
