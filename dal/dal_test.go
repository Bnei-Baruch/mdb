package dal

import (
	"github.com/Bnei-Baruch/mdb/migrations"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/rest"
	"github.com/Bnei-Baruch/mdb/utils"

	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

func RunMigrations(tmpDb *gorm.DB) {

	var visit = func(path string, f os.FileInfo, err error) error {
		match, _ := regexp.MatchString(".*\\.sql$", path)
		if !match {
			return nil
		}

		fmt.Printf("Applying migration %s\n", path)
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

	name := fmt.Sprintf("test_%s", strings.ToLower(utils.GenerateName(10)))
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
		Operation: rest.Operation{
			Station: "a station",
			User:    "111operator@dev.com",
		},
		FileName:      "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
		CaptureID:     "this.is.capture.id",
		CaptureSource: "mltcap",
	}

	if err := CaptureStart(cs); err == nil ||
		!strings.Contains(err.Error(), "User 111operator@dev.com not found.") {
		t.Error("Expected user not found, got:", err)
	}

	cs = rest.CaptureStart{
		Operation: rest.Operation{
			Station: "a station",
			User:    "operator@dev.com",
		},
		FileName:      "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
		CaptureID:     "this.is.capture.id",
		CaptureSource: "mltcap",
	}
	if err := CaptureStart(cs); err != nil {
		t.Error("CaptureStart should succeed.", err)
	}
}

func TestCaptureStop(t *testing.T) {
	baseDb, tmpDb, name := SwitchToTmpDb()
	defer DropTmpDB(baseDb, tmpDb, name)

	op := rest.Operation{
		Station: "a station",
		User:    "operator@dev.com",
	}
	start := rest.CaptureStart{
		Operation:     op,
		FileName:      "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
		CaptureID:     "this.is.capture.id",
		CaptureSource: "mlpcap",
	}
	if err := CaptureStart(start); err != nil {
		t.Error("CaptureStart should succeed.", err)
	}
	start = rest.CaptureStart{
		Operation:     op,
		FileName:      "heb_o_rav_bs-igeret_2016-09-14_lesson.mp4",
		CaptureID:     "this.is.capture.id",
		CaptureSource: "mlpcap",
	}
	if err := CaptureStart(start); err != nil {
		t.Error("CaptureStart should succeed.", err)
	}

	sha1 := "abcd1234abcd1234abcd1234abcd1234abcd1234"
	var size uint64 = 123
	part := "full"
	stopOp := rest.Operation{
		Station: "a station",
		User:    "operator@dev.com",
	}
	stop := rest.CaptureStop{
		CaptureStart: rest.CaptureStart{
			Operation:     stopOp,
			FileName:      "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
			CaptureID:     "this.is.capture.id",
			CaptureSource: "mltcap",
		},
		Sha1: sha1,
		Size: size,
		Part: part,
	}
	if err := CaptureStop(stop); err != nil {
		t.Error("CaptureStop should succeed.", err)
	}

	f := models.File{Name: "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4"}
	db.Where(&f).First(&f)
	if !f.Sha1.Valid || hex.EncodeToString([]byte(f.Sha1.String)) != sha1 {
		t.Error(fmt.Sprintf("Expected size %d got %d", size, f.Size))
	}
	if f.Size != size {
		t.Error(fmt.Sprintf("Expected size %d got %d", size, f.Size))
	}

	var cu models.ContentUnit
	db.Model(&f).Related(&cu, "ContentUnit")

	var ccu models.CollectionsContentUnit
	db.Model(&cu).Related(&ccu, "CollectionsContentUnit")
	if ccu.Name != part {
		t.Error(fmt.Sprintf("Expected part %s, got %s", part, ccu.Name))
	}

	stop = rest.CaptureStop{
		CaptureStart: rest.CaptureStart{
			Operation:     stopOp,
			FileName:      "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
			CaptureID:     "this.is.capture.id",
			CaptureSource: "mltcap",
		},
		Sha1: "1111111111111111111111111111111111111111",
		Size: 111,
	}
	if err := CaptureStop(stop); err == nil ||
		strings.Contains(err.Error(), "CaptureStop File already has different Sha1") {
		t.Error("Expected to fail with wrong Sha1 on CaptureStop: ", err)
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

	// Make sure code does not crash.
	ParseFileName("2017-01-04_02-40-19")
}

func AddFile(FileName string, Sha1 string, Size uint64) error {
	start := rest.CaptureStart{
		Operation: rest.Operation{
			Station: "a station",
			User:    "operator@dev.com",
		},
		FileName:  FileName,
		CaptureID: "this.is.capture.id",
	}
	if err := CaptureStart(start); err != nil {
		return err
	}

	stop := rest.CaptureStop{
		CaptureStart: rest.CaptureStart{
			Operation: rest.Operation{
				Station: "a station",
				User:    "operator@dev.com",
			},
			FileName:  FileName,
			CaptureID: "this.is.capture.id",
		},
		Sha1: Sha1,
		Size: Size,
	}
	if err := CaptureStop(stop); err != nil {
		return err
	}

	return nil
}

func TestDemux(t *testing.T) {
	baseDb, tmpDb, name := SwitchToTmpDb()
	defer DropTmpDB(baseDb, tmpDb, name)

	// Prepare file
	sha1 := "abcdef123456"
	if err := AddFile("lang_o_norav_part-a_2016-12-31_source.mp4", sha1, 111); err != nil {
		t.Error("Could not create file.", err)
	}

	origFileName := "lang_o_norav_part-a_2016-12-31_orig.mp4"
	proxyFileName := "lang_o_norav_part-a_2016-12-31_proxy.mp4"
	demux := rest.Demux{
		Operation: rest.Operation{
			Station: "a station",
			User:    "operator@dev.com",
		},
		Sha1: "abcdef123456",
		Original: rest.FileUpdate{
			FileName: origFileName,
			Sha1:     "aaaaaa111111",
			Size:     111,
		},
		Proxy: rest.FileUpdate{
			FileName: proxyFileName,
			Sha1:     "bbbbbb222222",
			Size:     222,
		},
	}
	if err := Demux(demux); err != nil {
		t.Error("Demux should succeed.", err)
	}

	orig := models.File{Name: origFileName}
	db.Where(&orig).First(&orig)
	if !orig.Sha1.Valid || hex.EncodeToString([]byte(orig.Sha1.String)) != demux.Original.Sha1 {
		t.Error(fmt.Sprintf("Expected sha1 %s got %s", demux.Original.Sha1, orig.Sha1.String))
	}
	if orig.Size != demux.Original.Size {
		t.Error(fmt.Sprintf("Expected size %d got %d", demux.Original.Size, orig.Size))
	}
	proxy := models.File{Name: proxyFileName}
	db.Where(&proxy).First(&proxy)
	if !proxy.Sha1.Valid || hex.EncodeToString([]byte(proxy.Sha1.String)) != demux.Proxy.Sha1 {
		t.Error(fmt.Sprintf("Expected sha1 %s got %s", demux.Proxy.Sha1, proxy.Sha1.String))
	}
	if proxy.Size != demux.Proxy.Size {
		t.Error(fmt.Sprintf("Expected size %d got %d", demux.Proxy.Size, proxy.Size))
	}
	sha1NullString, _ := Sha1ToNullString(sha1)
	source := models.File{Sha1: sha1NullString}
	db.Where(&source).First(&source)
	if uint64(proxy.ParentID.Int64) != source.ID {
		t.Error(fmt.Sprintf("Bad proxy parent id %d, expected %d", proxy.ParentID, source.ID))
	}
}

func TestSend(t *testing.T) {
	baseDb, tmpDb, name := SwitchToTmpDb()
	defer DropTmpDB(baseDb, tmpDb, name)

	// Prepare file
	sha1 := "abcdef123456"
	if err := AddFile("lang_o_norav_part-a_2016-12-31_source.mp4", sha1, 111); err != nil {
		t.Error("Could not create file.", err)
	}

	destFileName := "lang_o_norav_part-a_2016-12-31_dest.mp4"
	send := rest.Send{
		Operation: rest.Operation{
			Station: "a station",
			User:    "operator@dev.com",
		},
		Sha1: sha1,
		Dest: rest.FileUpdate{
			FileName: destFileName,
			Sha1:     "cccccc333333",
			Size:     333,
		},
	}
	if err := Send(send); err != nil {
		t.Error("Send should succeed.", err)
	}

	dest := models.File{Name: destFileName}
	db.Where(&dest).First(&dest)
	if !dest.Sha1.Valid || hex.EncodeToString([]byte(dest.Sha1.String)) != send.Dest.Sha1 {
		t.Error(fmt.Sprintf("Expected sha1 %s got %s", send.Dest.Sha1, dest.Sha1.String))
	}
	if dest.Size != send.Dest.Size {
		t.Error(fmt.Sprintf("Expected size %d got %d", send.Dest.Size, dest.Size))
	}
}

func TestUpload(t *testing.T) {
	baseDb, tmpDb, name := SwitchToTmpDb()
	defer DropTmpDB(baseDb, tmpDb, name)

	// Prepare file
	sha1 := "abcdef123456"
	fileName := "lang_o_norav_part-a_2016-12-31_source.mp4"
	if err := AddFile(fileName, sha1, 111); err != nil {
		t.Error("Could not create file.", err)
	}

	url := "http://this/is/some/url"
	upload := rest.Upload{
		Operation: rest.Operation{
			Station: "a station",
			User:    "operator@dev.com",
		},
		FileUpdate: rest.FileUpdate{
			FileName: fileName,
			Sha1:     sha1,
			Size:     111,
		},
		Url:          url,
		ExistingSha1: "1234",
	}
	if err := Upload(upload); err != nil {
		t.Error("Upload should succeed.", err)
	}
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UTC().UnixNano())
	os.Exit(m.Run())
}
