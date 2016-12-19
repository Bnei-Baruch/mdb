package dal

import (
	"github.com/Bnei-Baruch/mdb/rest"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"

    "fmt"
    "encoding/hex"
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
            "Bad filename, could not parse date (%s): %s", parts[4], err.Error())}
    }
    fn.DateStr = parts[4]

    return &fn, nil
}

func ValidateCapture(start rest.CaptureStart) (
    *models.User, *models.OperationType, *models.ContentType, *models.ContentType, *FileName, error) {
    var u = models.User{Email: start.User}
    if db.Where(&u).First(&u).RecordNotFound() {
        return nil, nil, nil, nil, nil, DalError{err: fmt.Sprintf("User %s not found.", start.User)}
    }

    var t = models.OperationType{Name: start.Type}
    if db.Where(&t).First(&t).RecordNotFound() {
        return nil, nil, nil, nil, nil, DalError{err: fmt.Sprintf("Operation type %s not found.", start.Type)}
    }

    var dl = models.ContentType{Name: "DAILY_LESSON"}
    if db.Where(&dl).First(&dl).RecordNotFound() {
        return nil, nil, nil, nil, nil, DalError{err: "Failed fetching \"DAILY_LESSON\" content type from db"}
    }

    var lp = models.ContentType{Name: "LESSON_PART"}
    if db.Where(&lp).First(&lp).RecordNotFound() {
        return nil, nil, nil, nil, nil, DalError{err: "Failed fetching \"LESSON_PART\" content type from db"}
    }

    fn, err := ParseFileName(start.FileName)
    if err != nil {
        return nil, nil, nil, nil, nil, DalError{err: "Error parsing filename.", reason: err}
    }

    return &u, &t, &dl, &lp, fn, err
}

func CaptureStart(start rest.CaptureStart) error {
    u, t, collectionType, contentUnitType, fileName, err := ValidateCapture(start)
    if err != nil {
        return err
    }

    // Execute (change DB).
    var c = models.Collection{ExternalID: start.CaptureID}
    if db.Where(&c).First(&c).RecordNotFound() {
        // Could not find collection by external id, create new.
        c = models.Collection{
            ExternalID: start.CaptureID,
            UID: utils.GenerateUID(8),
            Type: *collectionType,
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
        Type: *contentUnitType,
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
        Name: fileName.Part,
    }
    if err := db.Create(&m2m).Error; err != nil {
        return DalError{err: fmt.Sprintf("Failed adding collections content unit relation to db: %s", err.Error())}
    }

    var o = models.Operation{UID: utils.GenerateUID(8),  User: *u, Type: *t, Station: start.Station}
    if err := db.Create(&o).Error; err != nil {
        return DalError{err: fmt.Sprintf("Failed adding operation to db: %s", err.Error())}
    }

    var f = models.File{
        UID: utils.GenerateUID(8),
        Name: start.FileName,
        Operation: o,
        ContentUnit: cu,
    }
    if err := db.Create(&f).Error; err != nil {
        return DalError{err: fmt.Sprintf("Failed adding file to db: %s", err.Error())}
    }

    return nil
}

func CaptureStop(stop rest.CaptureStop) error {
    u, t, collectionType, contentUnitType, _, err := ValidateCapture(stop.CaptureStart)
    if err != nil {
        return err
    }

    var sha1 []byte
    sha1, err = hex.DecodeString(stop.Sha1)
    if err != nil {
        return DalError{err: fmt.Sprintf("Cannot convert sha1 %s to []bytes.", stop.Sha1)}
    }

    // Extract collection, content unit and file (should exist).
    var c = models.Collection{ExternalID: stop.CaptureID}
    if db.Where(&c).First(&c).RecordNotFound() {
        return DalError{err: fmt.Sprintf("Failed finding collection with external id %s.", stop.CaptureID)}
    }

    var cu = models.ContentUnit{Type: *contentUnitType}
    if db.Where(&cu).First(&cu).RecordNotFound() {
        return DalError{err: fmt.Sprintf("Failed adding content unit to db: %s")}
    }

    // Execute (change DB).
    var f = models.File{
        Name: stop.FileName,
        ContentUnit: models.ContentUnit{
            Type: *contentUnitType,
            Collections: []models.Collection{
                models.Collection{
                    Type: *collectionType,
                    ExternalID: stop.CaptureID,
                },
            },
        },
        Sha1: sha1,
        Size: stop.Size,
    }
    if errs := utils.FilterErrors(db.Where(&f).First(&f).Update(&f).GetErrors()); len(errs) > 0 {
        return DalError{err: fmt.Sprintf("Failed fetching file from db: %s", errs)}
    }

    // Create operation for the update.
    var o = models.Operation{UID: utils.GenerateUID(8),  User: *u, Type: *t, Station: stop.Station}
    if err := db.Create(&o).Error; err != nil {
        return DalError{err: fmt.Sprintf("Failed adding operation to db: %s", err.Error())}
    }

    return nil
}
