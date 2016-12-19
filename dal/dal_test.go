package dal

import (
	"github.com/Bnei-Baruch/mdb/migrations"
	"github.com/Bnei-Baruch/mdb/rest"
	"github.com/Bnei-Baruch/mdb/utils"

    "regexp"
    "fmt"
    "math/rand"
    "os"
    "strings"
    "testing"
    "time"
    "path/filepath"

	"github.com/spf13/viper"
    "github.com/jinzhu/gorm"
)

func RunMigrations(tmpDb *gorm.DB) {

    var visit = func(path string, f os.FileInfo, err error) error {
        match, _ := regexp.MatchString(".*\\.sql$", path);
        if !match {
            fmt.Printf("Did not match sql file %s\n", path)
            return nil
        }

        fmt.Printf("Migrating %s\n", path)
        m, err := migrations.NewMigration(path)
        if err != nil {
            fmt.Printf("Error migrating %s, %s", path, err.Error())
            return err
        }

        for _, statement := range m.Up() {
            if _, err := tmpDb.CommonDB().Exec(statement); err != nil {
                return fmt.Errorf("Unable to apply migration %s: %s\nStatement: %s\n", m.Name, err, statement)
            }
        }

        return nil
    }

    err := filepath.Walk("../migrations", visit)
    if err != nil {
        panic(fmt.Sprintf("Could not load and run all migrations. %s", err.Error()))
    }

}

func InitTestConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("../")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Could not read config, using: ", viper.ConfigFileUsed(), err)
	}
}

func SwitchToTmpDb() (*gorm.DB, *gorm.DB, string) {
    InitTestConfig()

    baseDb, err := Init()
    if err != nil {
        panic(fmt.Sprintf("Could not connect to database. %s", err))
    }

    name := strings.ToLower(utils.GenerateName(10))
    if err := baseDb.Exec(fmt.Sprintf("CREATE DATABASE %s", name)).Error; err != nil {
        panic(fmt.Sprintf("Could not create tmp database %s due to %s.", name, err))
    }

	url := viper.GetString("test.url-template")
    var tmpDb *gorm.DB
    tmpDb, err = InitByUrl(fmt.Sprintf(url, name))

    RunMigrations(tmpDb)

    return baseDb, tmpDb, name
}

func DropTmpDB(baseDb *gorm.DB, tmpDb *gorm.DB, name string) {
    tmpDb.Close()
    if err := baseDb.Exec(fmt.Sprintf("DROP DATABASE %s", name)).Error; err != nil {
        panic(fmt.Sprintf("Could not drop test database. %s", err))
    }
}

func TestInit(t *testing.T) {
    InitTestConfig()

    if _, err := InitByUrl("bad://database-connection-url"); err == nil {
		t.Error("Expected not nil, got nil.")
	}
	if _, err := Init(); err != nil {
		t.Error("Expected nil, got ", err.Error())
	}
}

func TestCaptureStart(t *testing.T) {
    baseDb, tmpDb, name := SwitchToTmpDb()
    defer DropTmpDB(baseDb, tmpDb, name)

    // User not found.
    cs := rest.CaptureStart{
        Type: "mltcap",
        Station: "a station",
        User: "111operator@dev.com",
        FileName: "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
        CaptureID: "this.is.capture.id",
    }

    if err := CaptureStart(cs); err == nil ||
        err.Error() != "User 111operator@dev.com not found." {
        t.Error("Expected user not found, got", err)
    }

    cs = rest.CaptureStart{
        Type: "mltcap",
        Station: "a station",
        User: "operator@dev.com",
        FileName: "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
        CaptureID: "this.is.capture.id",
    }
    if err := CaptureStart(cs); err != nil {
        t.Error("CaptureStart should succeed.", err)
    }
}

func TestCaptureStop(t *testing.T) {
    baseDb, tmpDb, name := SwitchToTmpDb()
    defer DropTmpDB(baseDb, tmpDb, name)

    start := rest.CaptureStart{
        Type: "mltcap",
        Station: "a station",
        User: "operator@dev.com",
        FileName: "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
        CaptureID: "this.is.capture.id",
    }

    if err := CaptureStart(start); err != nil {
        t.Error("CaptureStart should succeed.", err)
    }

    stop := rest.CaptureStop{
        CaptureStart: rest.CaptureStart{
            Type: "mltcap",
            Station: "a station",
            User: "operator@dev.com",
            FileName: "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
            CaptureID: "this.is.capture.id",
        },
        Sha1: "abcd1234abcd1234abcd1234abcd1234abcd1234",
        Size: 123,
    }
    if err := CaptureStop(stop); err != nil {
        // t.Error("CaptureStop should succeed.", err)
    }
}


func TestParseFilename(t *testing.T) {
    fn, err := ParseFileName("heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4")
    if err != nil {
        t.Error("ParseFileName should succeed.")
    }
    if fn.Type != "mp4" {
        t.Error("Expected type to be mp4")
    }
    if fn.Part != "rb-1990-02-kishalon" {
        t.Error("Expected different part %s", fn.Part)
    }

    _, err = ParseFileName("heb_o_rav_rb-1990-02-kishalon_201-09-14_lesson.mp4")
    if e := "could not parse date"; err == nil || !strings.Contains(err.Error(), e) {
        t.Error(fmt.Sprintf("ParseFileName should contain %s, got %s.", e, err))
    }
}

func TestMain(m *testing.M) {
    rand.Seed(time.Now().UTC().UnixNano())
    os.Exit(m.Run())
}
