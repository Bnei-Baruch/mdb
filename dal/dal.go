package dal

import (
	"github.com/Bnei-Baruch/mdb/rest"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"

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

    var u = models.User{Email: cs.User}
    if db.Where(&u).First(&u).RecordNotFound() {
        return DalError{err: fmt.Sprintf("User %s not found.", cs.User)}
    }

    var t = models.OperationType{Name: cs.Type}
    if db.Where(&t).First(&t).RecordNotFound() {
        return DalError{err: fmt.Sprintf("Operation type %s not found.", cs.Type)}
    }

    var o = models.Operation{UID: utils.GenerateUID(8),  User: u, Type: t, Station: cs.Station}
    if err := db.Create(&o).Error; err != nil {
        return DalError{err: fmt.Sprintf("Failed adding operation to db: %s", err.Error())}
    }

    var ct = models.ContentType{Name: "LESSON_PART"}
    if db.Where(&ct).First(&ct).RecordNotFound() {
        return DalError{err: "Failed fetching \"LESSON_PART\" content type from db"}
    }

    var cu = models.ContentUnit{
        Type: ct,
        TranslatedContent: models.TranslatedContent{
            Name: models.StringTranslation{Text: "Some name"},
            Description: models.StringTranslation{Text: "Some description"},
        },
    }
    if err := db.Create(&cu).Error; err != nil {
        return DalError{err: fmt.Sprintf("Failed adding content unit to db: %s", err.Error())}
    }

    var f = models.File{
        UID: utils.GenerateUID(8),
        Name: cs.FileName,
        Operation: o,
        ContentUnit: cu,
    }
    if err := db.Create(&f).Error; err != nil {
        return DalError{err: fmt.Sprintf("Failed adding file to db: %s", err.Error())}
    }

    return nil
}
