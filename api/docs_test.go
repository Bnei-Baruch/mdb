package api

import (
	"os"
	"testing"
	"time"

	"bytes"
	"encoding/json"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/edoshor/test2doc/test"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"net/url"
	"path"
)

var (
	router     *gin.Engine
	testServer *test.Server
)

func TestCaptureStartHandler(t *testing.T) {
	input := CaptureStartRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: "c12356789",
		},
		FileName:      "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
		CaptureSource: "mltcap",
		CollectionUID: "abcdefgh",
	}

	resp, err := testOperation(OP_CAPTURE_START, input)
	if err != nil {
		t.Error("Unknown error: ", err)
	}
	assertJsonOK(t, resp)
}

func TestCaptureStopHandler(t *testing.T) {
	input := CaptureStopRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: "c12356789",
		},
		AVFile: AVFile{
			File: File{
				FileName:  "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
				Sha1:      "012356789abcdef012356789abcdef0123456789",
				Size:      98737,
				CreatedAt: &Timestamp{time.Now()},
				Type:      "type",
				SubType:   "subtype",
				MimeType:  "mime_type",
				Language:  LANG_MULTI,
			},
			Duration: 892.1900,
		},
		CaptureSource: "mltcap",
		CollectionUID: "abcdefgh",
		Part:          "part",
	}

	resp, err := testOperation(OP_CAPTURE_STOP, input)
	if err != nil {
		t.Error("Unknown error: ", err)
	}
	assertJsonOK(t, resp)
}

func TestDemuxHandler(t *testing.T) {
	input := DemuxRequest{
		Operation: Operation{
			Station: "Demux Station",
			User:    "operator@dev.com",
		},
		CaptureSource: "mltbackup",
		Sha1:          "012356789abcdef012356789abcdef0123456789",
		Original: AVFile{
			File: File{
				FileName:  "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_o.mp4",
				Sha1:      "0987654321fedcba0987654321fedcba09876543",
				Size:      19837,
				CreatedAt: &Timestamp{time.Now()},
				Type:      "type",
				SubType:   "subtype",
				MimeType:  "mime_type",
				Language:  LANG_MULTI,
			},
			Duration: 892.1900,
		},
		Proxy: AVFile{
			File: File{
				FileName:  "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_p.mp4",
				Sha1:      "0987654321fedcba0987654321fedcba87654321",
				Size:      837,
				CreatedAt: &Timestamp{time.Now()},
				Type:      "type",
				SubType:   "subtype",
				MimeType:  "mime_type",
				Language:  LANG_HEBREW,
			},
			Duration: 892.1900,
		},
	}

	resp, err := testOperation(OP_DEMUX, input)
	if err != nil {
		t.Error("Unknown error: ", err)
	}
	assertJsonOK(t, resp)
}

func TestTrimHandler(t *testing.T) {
	input := TrimRequest{
		Operation: Operation{
			Station:    "Trim station",
			User:       "operator@dev.com",
			WorkflowID: "t12356789",
		},
		CaptureSource: "capture_source",
		OriginalSha1:  "0987654321fedcba0987654321fedcba09876543",
		ProxySha1:     "0987654321fedcba0987654321fedcba87654321",
		Original: AVFile{
			File: File{
				FileName:  "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_o_trim.mp4",
				Sha1:      "0987654321fedcba0987654321fedcba11111111",
				Size:      19800,
				CreatedAt: &Timestamp{time.Now()},
				Type:      "type",
				SubType:   "subtype",
				MimeType:  "mime_type",
				Language:  LANG_MULTI,
			},
			Duration: 871,
		},
		Proxy: AVFile{
			File: File{
				FileName:  "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_p_trim.mp4",
				Sha1:      "0987654321fedcba0987654321fedcba22222222",
				Size:      694,
				CreatedAt: &Timestamp{time.Now()},
				Type:      "type",
				SubType:   "subtype",
				MimeType:  "mime_type",
				Language:  LANG_HEBREW,
			},
			Duration: 871,
		},
		In:  []float64{0.00, 198.23},
		Out: []float64{10.50, 207.31},
	}

	resp, err := testOperation(OP_TRIM, input)
	if err != nil {
		t.Error("Unknown error: ", err)
	}
	assertJsonOK(t, resp)
}

func TestUploadHandler(t *testing.T) {
	input := UploadRequest{
		Operation: Operation{
			Station: "Upload station",
			User:    "111operator@dev.com",
		},
		AVFile: AVFile{
			File: File{
				FileName:  "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_o.mp4",
				Sha1:      "0987654321fedcba0987654321fedcba09876543",
				Size:      19837,
				CreatedAt: &Timestamp{time.Now()},
			},
			Duration: 892,
		},
		Url: "https://example.com/heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
	}

	resp, err := testOperation(OP_UPLOAD, input)
	if err != nil {
		t.Error("Unknown error: ", err)
	}
	assertJsonOK(t, resp)
}

func testOperation(name string, input interface{}) (*http.Response, error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(input)
	u, _ := url.Parse(testServer.URL)
	u.Path = path.Join(u.Path, "operations", name)
	return http.Post(u.String(), "application/json", b)
}

func assertJsonOK(t *testing.T, resp *http.Response) {
	if resp.StatusCode != http.StatusOK {
		t.Errorf("HTTP status_code should be 200, was: %d", resp.StatusCode)
	}

	if h := resp.Header.Get("Content-Type"); h != "application/json; charset=utf-8" {
		t.Errorf("Content-Type should be application/json, was %s", h)
	}

	var body map[string]interface{}
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		t.Error("Error parsing JSON response: ", err)
	}

	if body["status"] != "ok" {
		t.Error("Unexpected response: ", body)
	}
}

// For now leave this empty. When needed make sure code
// extracts variables from request url.
func Vars(req *http.Request) map[string]string {
	return make(map[string]string)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	router = gin.New()
	router.Use(utils.ErrorHandlingMiddleware(), gin.Recovery())
	SetupRoutes(router)

	test.RegisterURLVarExtractor(Vars)
	var err error
	testServer, err = test.NewServer(router)
	if err != nil {
		panic(err.Error())
	}

	if err := utils.InitTestDB(); err != nil {
		panic(err)
	}
	if err := CONTENT_TYPE_REGISTRY.Init(); err != nil {
		panic(err)
	}
	if err := OPERATION_TYPE_REGISTRY.Init(); err != nil {
		panic(err)
	}

	s := m.Run()
	utils.DestroyTestDB()
	testServer.Finish()
	os.Exit(s)
}
