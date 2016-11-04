package dal

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func Init() bool {
	// TODO: Move this from config.
	url := "postgres://localhost/mdb?sslmode=disable"
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
