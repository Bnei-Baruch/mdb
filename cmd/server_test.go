package cmd

import (
	"github.com/Bnei-Baruch/mdb/dal"
	"github.com/Bnei-Baruch/mdb/rest"

	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/adams-sarah/test2doc/test"
)

var testServer *test.Server

func Post(command string, body string) (*http.Response, error) {
	u, _ := url.Parse(testServer.URL)
	u.Path = path.Join(u.Path, command)
	return http.Post(u.String(), "application/json", strings.NewReader(body))
}

func TestCaptureStart(t *testing.T) {
	baseDb, tmpDb, name := dal.SwitchToTmpDb()
	defer dal.DropTmpDB(baseDb, tmpDb, name)

	fmt.Println("Testing capture start on server.")

    res, err := Post("/operations/capture_start", `
        {
            "station": "10.66.1.120", 
            "user": "operator@dev.com", 
            "file_name": "mlt_o_rav_2017-01-04_lesson_boker_part1", 
            "size": 13860980880, 
            "created_at": 1484155956,
            "capture_id": "c1483491605226", 
            "capture_source": "mltcap"
        }`)
	if err != nil || res == nil || res.StatusCode != 200 {
		t.Error("Should succeed starting capture.", err)
	}

	res, err = Post("/operations/capture_start", `
        {
            "bad": "json"
        }`)
	if err != nil || res == nil || res.StatusCode != 400 {
		t.Error("Should gracefully fail with status code 400.", err, res)
	}
}

func TestCaptureStop(t *testing.T) {
	baseDb, tmpDb, name := dal.SwitchToTmpDb()
	defer dal.DropTmpDB(baseDb, tmpDb, name)

	fmt.Println("Testing capture stop.")

	start := rest.CaptureStart{
        Operation:     rest.Operation{Station: "a station", User: "operator@dev.com"},
		FileName:      "mlt_o_rav_2017-01-04_lesson_boker_part1",
		CaptureID:     "c1483491605226",
		CaptureSource: "mlpcap",
	}
	dal.CaptureStart(start)

	res, err := Post("/operations/capture_stop", `
        {
            "capture_source": "mltcap",
            "station": "10.66.1.120", 
            "user": "operator@dev.com", 
            "file_name": "mlt_o_rav_2017-01-04_lesson_boker_part1", 
            "created_at": 1484155956,
            "size": 13860980880, 
            "sha1": "4d4ad60c11738630178a08560362a0b2c323f87a",
            "part": "1",
            "capture_id": "c1483491605226"
        }`)
	if err != nil || res == nil || res.StatusCode != 200 {
		t.Error("Should succeed stoping capture.", err)
	}
}

func TestDemux(t *testing.T) {
	baseDb, tmpDb, name := dal.SwitchToTmpDb()
	defer dal.DropTmpDB(baseDb, tmpDb, name)

	fmt.Println("Testing demux.")

	dal.AddTestFile("mlt_o_rav_2017-01-04_lesson_boker_source",
        "4d4ad60c11738630178a08560362a0b2c323f87a",
        123123)

	res, err := Post("/operations/demux", `
        {
            "station": "10.66.1.120", 
            "user": "operator@dev.com", 
            "file_name": "mlt_o_rav_2017-01-04_lesson_boker_source", 
            "created_at": 1484155956,
            "sha1": "4d4ad60c11738630178a08560362a0b2c323f87a",
            "original": {
                "file_name": "mlt_o_rav_2017-01-04_lesson_boker_orig", 
                "created_at": 1485555956,
                "sha1": "4d4ad60c111110178a08560362a0b2c323f87a",
                "size": 111
            },
            "proxy": {
                "file_name": "mlt_o_rav_2017-01-04_lesson_boker_proxy", 
                "created_at": 1486666956,
                "sha1": "4d4ad60c222220178a08560362a0b2c323f87a",
                "size": 222
            }
        }`)
	if err != nil || res == nil || res.StatusCode != 200 {
		t.Error("Should succeed demux.", err)
	}
}

func TestSend(t *testing.T) {
	baseDb, tmpDb, name := dal.SwitchToTmpDb()
	defer dal.DropTmpDB(baseDb, tmpDb, name)

	fmt.Println("Testing send.")

	dal.AddTestFile("mlt_o_rav_2017-01-04_lesson_boker_source",
        "4d4ad60c11738630178a08560362a0b2c323f87a",
        123123)

	res, err := Post("/operations/send", `
        {
            "station": "10.66.1.120", 
            "user": "operator@dev.com", 
            "file_name": "mlt_o_rav_2017-01-04_lesson_boker_source", 
            "created_at": 1484155956,
            "sha1": "4d4ad60c11738630178a08560362a0b2c323f87a",
            "dest": {
                "file_name": "mlt_o_rav_2017-01-04_lesson_boker_dest", 
                "created_at": 1485555956,
                "sha1": "4d4ad60c111110178a08560362a0b2c323f87a",
                "size": 111222333
            }
        }`)
	if err != nil || res == nil || res.StatusCode != 200 {
		t.Error("Should succeed send.", err)
	}
}

func TestUpload(t *testing.T) {
	baseDb, tmpDb, name := dal.SwitchToTmpDb()
	defer dal.DropTmpDB(baseDb, tmpDb, name)

	fmt.Println("Testing upload.")

	dal.AddTestFile("mlt_o_rav_2017-01-04_lesson_boker_source",
        "4d4ad60c11738630178a08560362a0b2c323f87a",
        123123)

	res, err := Post("/operations/upload", `
        {
            "station": "10.66.1.120", 
            "user": "operator@dev.com", 
            "file_name": "mlt_o_rav_2017-01-04_lesson_boker_source", 
            "created_at": 1484155956,
            "sha1": "4d4ad60c11738630178a08560362a0b2c323f87a",
            "size": 123123,
            "url": "http://some/url/file/was/uploaded/to",
            "duration": 9878,
            "existing": {
                "file_name": "mlt_o_rav_2017-01-04_lesson_boker_dest", 
                "created_at": 1485555956,
                "sha1": "4d4ad60c111110178a08560362a0b2c323f87a"
            }
        }`)
	if err != nil || res == nil || res.StatusCode != 200 {
		t.Error("Should succeed upload.", err)
	}
}

// For now leave this empty. When needed make sure code
// extracts variables from request url.
func Vars(req *http.Request) map[string]string {
	m := make(map[string]string)
	return m
}

func TestMain(m *testing.M) {
	test.RegisterURLVarExtractor(Vars)
	serverFn(nil, nil)
	var err error
	testServer, err = test.NewServer(router)
	if err != nil {
		panic(err.Error())
	}
	exitCode := m.Run()
	testServer.Finish()
	os.Exit(exitCode)
}
