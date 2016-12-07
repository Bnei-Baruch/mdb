package dal

import (
	"github.com/Bnei-Baruch/mdb/rest"

	// "database/sql"
    "fmt"
    "testing"
)

func CreateTmpDB() string {
    db, err := Init()
    if err == nil {
        panic(fmt.Sprintf("Could not connect to database. %s", err))
    }

	url := viper.GetString("test.url-template")
    database := utils.GenerateUID(10)
    rows, err := db.Query(fmt.Sprintf(url, database))


}

func DropTmpDB(name string) {
}

func TestInit(t *testing.T) {
	if success := Init(); !success {
		t.Error("Expected true, got ", success)
	}
}

func TestCaptureStart(t *testing.T) {
    cs := rest.CaptureStart{
        Type: "type",
        Station: "a station",
        User: "operator@dev.com",
        FileName: "some.file.name",
        CaptureID: "this.is.capture.id",
    }

    ok, err := CaptureStart(cs)
    if !ok {
        t.Error("CaptureStart should succeed.", err)
    }
}

