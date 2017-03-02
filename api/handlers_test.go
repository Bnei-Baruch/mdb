package api

import (
	"testing"
	"time"
	"os"
	"math/rand"

	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"net/http/httptest"
	"bytes"
	"encoding/json"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/Gurpartap/logrus-stack"
	log "github.com/Sirupsen/logrus"
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

	w := testOperationHandler(CaptureStartHandler, input)
	assertJsonOK(t, w)
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
				Type: "type",
				SubType: "subtype",
				MimeType: "mime_type",
				Language: LANG_MULTI,
			},
			Duration: 892.1900,
		},
		CaptureSource: "mltcap",
		CollectionUID: "abcdefgh",
		Part:          "part",
	}

	w := testOperationHandler(CaptureStopHandler, input)
	assertJsonOK(t, w)
}

func TestDemuxHandler(t *testing.T) {
	input := DemuxRequest{
		Operation: Operation{
			Station: "Demux Station",
			User:    "operator@dev.com",
		},
		CaptureSource: "mltbackup",
		Sha1:      "012356789abcdef012356789abcdef0123456789",
		Original: AVFile{
			File: File{
				FileName:  "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_o.mp4",
				Sha1:      "0987654321fedcba0987654321fedcba09876543",
				Size:      19837,
				CreatedAt: &Timestamp{time.Now()},
				Type: "type",
				SubType: "subtype",
				MimeType: "mime_type",
				Language: LANG_MULTI,
			},
			Duration: 892.1900,
		},
		Proxy: AVFile{
			File: File{
				FileName:  "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_p.mp4",
				Sha1:      "0987654321fedcba0987654321fedcba87654321",
				Size:      837,
				CreatedAt: &Timestamp{time.Now()},
				Type: "type",
				SubType: "subtype",
				MimeType: "mime_type",
				Language: LANG_HEBREW,
			},
			Duration: 892.1900,
		},
	}

	w := testOperationHandler(DemuxHandler, input)
	assertJsonOK(t, w)
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

	w := testOperationHandler(UploadHandler, input)
	assertJsonOK(t, w)
}

func testOperationHandler(handler gin.HandlerFunc, input interface{}) *httptest.ResponseRecorder {
	r := gin.New()
	r.Use(utils.ErrorHandlingMiddleware(), gin.Recovery())
	r.POST("/operations/test", handler)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(input)
	req, _ := http.NewRequest("POST", "/operations/test", b)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func assertJsonOK(t *testing.T, w *httptest.ResponseRecorder) {
	if w.Code != http.StatusOK {
		t.Errorf("HTTP status_code should be 200, was: %d", w.Code)
	}

	if body := w.Body.String(); body != "{\"status\":\"ok\"}\n" {
		t.Errorf("Response should be {\"status\":\"ok\"}, was: %s", body)
	}

	if h := w.HeaderMap.Get("Content-Type"); h != "application/json; charset=utf-8" {
		t.Errorf("Content-Type should be application/json, was %s", h)
	}
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UTC().UnixNano())
	log.AddHook(logrus_stack.StandardHook())
	gin.SetMode(gin.TestMode)
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
	os.Exit(s)
}
