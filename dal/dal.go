package dal

import (
	"github.com/Bnei-Baruch/mdb/rest"
	"github.com/Bnei-Baruch/mdb/models"

    "fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
    "github.com/jinzhu/gorm"
)

var db *gorm.DB

func Init() (*gorm.DB, error) {
	url := viper.GetString("mdb.url")
    return InitByUrl(url)
}

func InitByUrl(url string) (*gorm.DB, error) {
    var err error
	db, err = gorm.Open("postgres", url)
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

func CaptureStart(cs rest.CaptureStart) (error) {

    var user = models.User{Email: cs.User}
    fmt.Printf("%s\n", user)
    db.First(&user)
    fmt.Printf("%s\n", user)

    return DalError{err: "Not Implemented."}
}
