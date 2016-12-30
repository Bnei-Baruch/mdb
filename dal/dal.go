package dal

import (
	"github.com/Bnei-Baruch/mdb/rest"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"

    "fmt"
    "encoding/hex"
    "path/filepath"
	"database/sql"
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

// Validates operation input and creates object (without inserting to db).
func CreateOperation(o rest.Operation) (models.Operation, error) {
    ret := models.Operation{UID: utils.GenerateUID(8), Station: o.Station}
    var u = models.User{Email: o.User}
    if db.Where(&u).First(&u).RecordNotFound() {
        return ret, DalError{err: fmt.Sprintf("User %s not found.", o.User)}
    }
    ret.User = u

    var t = models.OperationType{Name: o.Type}
    if db.Where(&t).First(&t).RecordNotFound() {
        return ret, DalError{err: fmt.Sprintf("Operation type %s not found.", o.Type)}
    }
    ret.Type = t

    return ret, nil
}

func ValidateCapture(start rest.CaptureStart) (
    *models.ContentType, *models.ContentType, *FileName, error) {
    var dl = models.ContentType{Name: "DAILY_LESSON"}
    if db.Where(&dl).First(&dl).RecordNotFound() {
        return nil, nil, nil, DalError{err: "Failed fetching \"DAILY_LESSON\" content type from db"}
    }

    var lp = models.ContentType{Name: "LESSON_PART"}
    if db.Where(&lp).First(&lp).RecordNotFound() {
        return nil, nil, nil, DalError{err: "Failed fetching \"LESSON_PART\" content type from db"}
    }

    fn, err := ParseFileName(start.FileName)
    if err != nil {
        return nil, nil, nil, DalError{err: "Error parsing filename.", reason: err}
    }

    return &dl, &lp, fn, err
}

func CaptureStart(start rest.CaptureStart) error {
    collectionType, contentUnitType, fileName, err := ValidateCapture(start)
    if err != nil {
        return err
    }

    var o models.Operation
    o, err = CreateOperation(start.Operation)
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

func FindFileByExternalIDAndFileName(externalID string, fileName string) (models.File, error) {
    // Select file by ExternalID and FileName
    var id uint64
    f := models.File{}
    if err := db.CommonDB().QueryRow(
        "select files.id from files, collections, collections_content_units where " +
        "files.name = $1 and collections.external_id = $2 and " +
        "collections.id = collections_content_units.collection_id and " +
        "files.content_unit_id = collections_content_units.content_unit_id",
        fileName, externalID).Scan(&id); err != nil {
        return f, DalError{err: fmt.Sprintf("Failed fetching file id due to %s", err.Error())}
    }
    f.ID = id
    if errs := db.Where(&f).First(&f).GetErrors(); len(errs) > 0 {
        return f, DalError{err: fmt.Sprintf("Failed fetching file: %+v due to %s", f, errs)}
    }
    return f, nil
}

func CaptureStop(stop rest.CaptureStop) error {
    _, _, _, err := ValidateCapture(stop.CaptureStart)
    if err != nil {
        return err
    }

    var o models.Operation
    o, err = CreateOperation(stop.Operation)
    if err != nil {
        return err
    }

    sha1, err := Sha1ToNullString(stop.Sha1)
    if err != nil {
        return DalError{err: fmt.Sprintf("Cannot convert sha1 %s to []bytes.", stop.Sha1)}
    }

    // Execute (change DB).
    f, err := FindFileByExternalIDAndFileName(stop.CaptureID, stop.FileName)
    if err != nil {
        return err
    }
    if f.Sha1.Valid && f.Sha1 != sha1 {
        return DalError{err: fmt.Sprintf("File already has different Sha1 existing: %s vs new: %s",
            hex.EncodeToString([]byte(f.Sha1.String)), stop.Sha1)}
    }
    if f.Size > 0 && f.Size != stop.Size {
        return DalError{err: fmt.Sprintf("File already has different Size existing: %d vs new: %d",
            f.Size, stop.Size)}
    }
    f.Sha1 = sha1
    f.Size = stop.Size
    if errs := utils.FilterErrors(db.Model(&f).Update(&f).GetErrors()); len(errs) > 0 {
        return DalError{err: fmt.Sprintf("Failed updating file in db: %s", errs)}
    }

    // Create operation for the update.
    if err := db.Create(&o).Error; err != nil {
        return DalError{err: fmt.Sprintf("Failed adding operation to db: %s", err.Error())}
    }

    return nil
}

func CreateFile(update rest.FileUpdate) (models.File, error) {
    f := models.File{
        UID: utils.GenerateUID(8),
        Name: update.FileName,
        Size: update.Size,
    }
    sha1, err := Sha1ToNullString(update.Sha1)
    if err != nil {
        return f, err
    }
    f.Sha1 = sha1
    return f, nil
}

func Demux(demux rest.Demux) error {
    o, err := CreateOperation(demux.Operation)
    if err != nil {
        return err
    }

    sha1, err := Sha1ToNullString(demux.Sha1)
    if err != nil {
        return DalError{err: fmt.Sprintf("Cannot convert sha1 %s to []bytes.", demux.Sha1)}
    }

    f := models.File{Sha1: sha1}
    db.Where(&f).First(&f)
    if errs := db.Where(&f).First(&f).GetErrors(); len(errs) > 0 {
        return DalError{err: fmt.Sprintf("Could not find file %+v by Sha1, got errors: %s", f, errs)}
    }

    // Create orig file
    origFile, err := CreateFile(demux.Original)
    if err != nil {
        return err
    }
    // Create operation for the update.
    if err := db.Create(&o).Error; err != nil {
        return DalError{err: fmt.Sprintf("Failed adding operation to db: %s", err.Error())}
    }
    origFile.ParentID = sql.NullInt64{Valid: true, Int64: int64(f.ID)}
    origFile.OperationID = o.ID
    origFile.ContentUnitID = f.ContentUnitID


    fmt.Printf("Operation: %+v\n", o)

    if err := db.Debug().Create(&origFile).Error; err != nil {
        return DalError{err: fmt.Sprintf("Failed creating orig file: %s", err.Error())}
    }

    // Create proxy file
    var proxyFile models.File
    proxyFile, err = CreateFile(demux.Proxy)
    if err != nil {
        return err
    }
    proxyFile.ParentID = sql.NullInt64{Valid: true, Int64: int64(f.ID)}
    proxyFile.OperationID = o.ID
    proxyFile.ContentUnitID = f.ContentUnitID

    if err := db.Create(&proxyFile).Error; err != nil {
        return DalError{err: fmt.Sprintf("Failed creating proxy file: %s", err.Error())}
    }

    return nil
}

func Sha1ToNullString(sha1 string) (sql.NullString, error) {
    bytes, err := hex.DecodeString(sha1)
    if err != nil {
        return sql.NullString{Valid: false, String: ""},
               DalError{err: fmt.Sprintf("Cannot convert sha1 %s to []bytes.", sha1)}
    }
    return sql.NullString{Valid: len(bytes) > 0, String: string(bytes)}, nil
}

func Send(send rest.Send) error {
    o, err := CreateOperation(send.Operation)
    if err != nil {
        return err
    }

    sha1, err := Sha1ToNullString(send.Sha1)
    if err != nil {
        return DalError{err: fmt.Sprintf("Cannot convert sha1 %s to []bytes.", send.Sha1)}
    }

    f := models.File{Sha1: sha1}
    db.Where(&f).First(&f)
    if errs := db.Where(&f).First(&f).GetErrors(); len(errs) > 0 {
        return DalError{err: fmt.Sprintf("Could not find file %+v by Sha1, got errors: %s", f, errs)}
    }

    // Create dest file
    destSha1, err := Sha1ToNullString(send.Dest.Sha1)
    if err != nil {
        return err
    }
    // Create operation for the update.
    if err := db.Create(&o).Error; err != nil {
        return DalError{err: fmt.Sprintf("Failed adding operation to db: %s", err.Error())}
    }
    sendFile := models.File{
        UID: utils.GenerateUID(8),
        Name: send.Dest.FileName,
        Sha1: destSha1,
        Size: send.Dest.Size,
        OperationID: o.ID,
        ContentUnitID: f.ContentUnitID,
        ParentID: sql.NullInt64{Valid: true, Int64: int64(f.ID)},
    }
    if err := db.Create(&sendFile).Error; err != nil {
        return DalError{err: fmt.Sprintf("Failed creating send file: %s", err.Error())}
    }

    return nil
}

