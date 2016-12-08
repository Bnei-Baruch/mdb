package dal

import (
	"github.com/Bnei-Baruch/mdb/rest"
	"github.com/Bnei-Baruch/mdb/utils"

	"database/sql"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "testing"
    "path/filepath"

	"github.com/spf13/viper"
)

func RunMigrations(tmpDb *sql.DB) {

    var visit = func(path string, f os.FileInfo, err error) error {
        file, err := ioutil.ReadFile(path)
        if err != nil {
            return err
        }

        requests := strings.Split(string(file), ";")

        for _, request := range requests {
            if _, err = tmpDb.Exec(request); err != nil {
                return err
            }
        }
        return nil
    }

    err := filepath.Walk("migrations", visit)
    if err != nil {
        panic(fmt.Sprintf("Could not load and run all migrations. %s", err))
    }

}

func SwitchToTmpDb() (*sql.DB, string) {
    baseDb, err := Init()
    if err == nil {
        panic(fmt.Sprintf("Could not connect to database. %s", err))
    }

    name := utils.GenerateUID(10)
    _, err = baseDb.Exec(fmt.Sprintf("CREATE DATABASE %s", name))
    if err != nil {
        panic(fmt.Sprintf("Could not create tmp database %s.", err))
    }

	url := viper.GetString("test.url-template")
    var tmpDb *sql.DB
    tmpDb, err = InitByUrl(fmt.Sprintf(url, name))

    RunMigrations(tmpDb)

    return baseDb, name
}

func DropTmpDB(baseDb *sql.DB, name string) {
    _, err := baseDb.Exec(fmt.Sprintf("DROP DATABASE %s", name))
    if err != nil {
        panic(fmt.Sprintf("Could not drop test database. %s", err))
    }
}

func TestInit(t *testing.T) {
	if _, err := Init(); err != nil {
		t.Error("Expected nil, got ", err)
	}
    if _, err := InitByUrl("bad://database-connection-url"); err != nil {
		t.Error("Expected nil, got ", err)
	}
}

func TestCaptureStart(t *testing.T) {
    baseDb, name := SwitchToTmpDb()

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

    DropTmpDB(baseDb, name)
}

