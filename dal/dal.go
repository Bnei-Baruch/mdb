package dal

import (
	"github.com/Bnei-Baruch/mdb/rest"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"

    "fmt"
    "path/filepath"
    "strings"
    "time"

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
    if e.reason != nil {
        return fmt.Sprintf("%s due to %s", e.err, e.reason.Error())
    }
    return e.err
}

type FileName struct {
    Name        string
    Base        string
    Type        string  // File extension, mp3 or mp4 or other.
    Language    string
    Rav         bool
    Part        string
    Date        time.Time
    DateStr     string
}

func ParseFileName(name string) (*FileName, error) {
    fn := FileName{
        Name: name,
        Base: filepath.Base(name),
        Type: filepath.Ext(name)[1:],
    }
    parts := strings.Split(strings.TrimSuffix(fn.Base, filepath.Ext(fn.Base)), "_")
    if len(parts) < 4 {
        return nil, DalError{err: fmt.Sprintf(
            "Bad filename, expected at least 4 parts, found %d: %s", len(parts), parts)}
    }
    fn.Language = parts[0]
    if parts[2] == "rav" {
        fn.Rav = true
    } else if parts[2] == "norav" {
        fn.Rav = false
    } else {
        return nil, DalError{err: fmt.Sprintf(
            "Bad filename, expected rav/norav got %s", parts[2])}
    }
    fn.Part = parts[3]
    var err error
    fn.Date, err = time.Parse("2006-01-02", parts[4])
    if err != nil {
        return nil, DalError{err: fmt.Sprintf(
            "Bad filename, could not parse date (%s): %s", parts[3], err.Error())}
    }
    fn.DateStr = parts[4]

    return &fn, nil
}

func CaptureStart(cs rest.CaptureStart) (error) {
    // Validate operation may be executed correctly.
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

    var dl = models.ContentType{Name: "DAILY_LESSON"}
    if db.Where(&dl).First(&dl).RecordNotFound() {
        return DalError{err: "Failed fetching \"DAILY_LESSON\" content type from db"}
    }

    var lp = models.ContentType{Name: "LESSON_PART"}
    if db.Where(&lp).First(&lp).RecordNotFound() {
        return DalError{err: "Failed fetching \"LESSON_PART\" content type from db"}
    }

    fn, err := ParseFileName(cs.FileName)
    if err != nil {
        return DalError{err: "Error parsing filename.", reason: err}
    }

    // Execute (change DB).
    var c = models.Collection{ExternalID: cs.CaptureID}
    if db.Where(&c).First(&c).RecordNotFound() {
        // Could not find collection by external id, create new.
        c = models.Collection{
            ExternalID: cs.CaptureID,
            UID: utils.GenerateUID(8),
            Type: dl,
            TranslatedContent: models.TranslatedContent{
                Name: models.StringTranslation{Text: "Collection name"},
                Description: models.StringTranslation{Text: "Collection description"},
            },
        }
        if err := db.Create(&c).Error; err != nil {
            return DalError{err: fmt.Sprintf("Failed adding collection to db: %s", err.Error())}
        }
    }

    var cu = models.ContentUnit{
        Type: lp,
        UID: utils.GenerateUID(8),
        TranslatedContent: models.TranslatedContent{
            Name: models.StringTranslation{Text: "Content unit name"},
            Description: models.StringTranslation{Text: "Content unit description"},
        },
    }
    if err := db.Create(&cu).Error; err != nil {
        return DalError{err: fmt.Sprintf("Failed adding content unit to db: %s", err.Error())}
    }

    var m2m = models.CollectionsContentUnit{
        Collection: c,
        ContentUnit: cu,
        Name: fn.Part,
    }
    if err := db.Create(&m2m).Error; err != nil {
        return DalError{err: fmt.Sprintf("Failed adding collections content unit relation to db: %s", err.Error())}
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
