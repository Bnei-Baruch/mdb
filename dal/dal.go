package dal

import (
	"github.com/Bnei-Baruch/mdb/rest"

	"database/sql"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func Init() (*sql.DB, error) {
	url := viper.GetString("mdb.url")
    return InitByUrl(url)
}

func InitByUrl(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
        return nil, DalError{err: "Error opening db.", reason: err}
	}
    return db, nil
}

type DalError struct {
    err string
    reason error
}

func (e DalError) Error() string {
    return e.err
}

func CaptureStart(rest.CaptureStart) (bool, error) {

    // Implementation goes here...

    return false, DalError{err: "Not Implemented."}
}
