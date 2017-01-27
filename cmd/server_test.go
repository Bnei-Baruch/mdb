package cmd

import (
	"github.com/Bnei-Baruch/mdb/dal"

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
            "bad": "json"
        }`)
	if err != nil || res == nil || res.StatusCode != 400 {
		t.Error("Should gracefully fail with status code 400.", err, res)
	}

	res, err = Post("/operations/capture_start", `
        {
            "capture_source": "mltcap",
            "station": "10.66.1.120", 
            "user": "operator@dev.com", 
            "file_name": "mlt_o_rav_2017-01-04_lesson_boker_part1", 
            "capture_id": "c1483491605226", 
            "size": 13860980880, 
            "sha1": "4d4ad60c11738630178a08560362a0b2c323f87a",
            "created_at": 1484155956
        }`)
	if err != nil || res == nil || res.StatusCode != 200 {
		t.Error("Should succeed starting capture.", err, res)
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
