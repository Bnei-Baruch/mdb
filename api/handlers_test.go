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
)

func TestCaptureStartHandler(t *testing.T) {
	input := CaptureStartRequest{
		Operation: Operation{
			Station:    "a station",
			User:       "111operator@dev.com",
			WorkflowID: "c12356789",
		},
		FileKey:       FileKey{FileName: "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4"},
		CollectionUID: "abcdefgh",
		CaptureSource: "mltcap",
	}

	w := testOperationHandler(CaptureStartHandler, input)
	assertJsonOK(t, w)
}

func TestCaptureStopHandler(t *testing.T) {
	input := CaptureStopRequest{
		Operation: Operation{
			Station:    "a station",
			User:       "111operator@dev.com",
			WorkflowID: "c12356789",
		},
		FileKey:       FileKey{FileName: "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4"},
		CollectionUID: "abcdefgh",
		CaptureSource: "mltcap",
		Sha1:          "012356789abcdef012356789abcdef0123456789",
		Size:          123,
		ContentType:   "content_type",
		Part:          "part",
	}

	w := testOperationHandler(CaptureStopHandler, input)
	assertJsonOK(t, w)
}

func testOperationHandler(handler gin.HandlerFunc, input interface{}) *httptest.ResponseRecorder {
	r := gin.Default()
	r.POST("/test", handler)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(input)
	req, _ := http.NewRequest("POST", "/test", b)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func assertJsonOK(t *testing.T, w *httptest.ResponseRecorder) {
	if w.Code != http.StatusOK {
		t.Errorf("HTTP status_code should be 200, was: %d", w.Code)
	}

	if w.Body.String() != "{\"status\":\"ok\"}\n" {
		t.Errorf("Response should be {\"status\":\"ok\"}, was: %s", w.Body.String())
	}

	if w.HeaderMap.Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("Content-Type should be application/json, was %s", w.HeaderMap.Get("Content-Type"))
	}
}

//
//func TestCaptureStart(t *testing.T) {
//	baseDb, tmpDb, name := SwitchToTmpDb()
//	defer DropTmpDB(baseDb, tmpDb, name)
//
//	// User not found.
//	cs := CaptureStart{
//		Operation: Operation{
//			Station: "a station",
//			User:    "111operator@dev.com",
//		},
//		FileName:      "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
//		CaptureID:     "this.is.capture.id",
//		CaptureSource: "mltcap",
//	}
//
//	if err := CaptureStart(cs); err == nil ||
//		!strings.Contains(err.Error(), "User 111operator@dev.com not found.") {
//		t.Error("Expected user not found, got:", err)
//	}
//
//	cs = CaptureStart{
//		Operation: Operation{
//			Station: "a station",
//			User:    "operator@dev.com",
//		},
//		FileName:      "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
//		CaptureID:     "this.is.capture.id",
//		CaptureSource: "mltcap",
//	}
//	if err := CaptureStart(cs); err != nil {
//		t.Error("CaptureStart should succeed.", err)
//	}
//}
//
//func TestCaptureStop(t *testing.T) {
//	baseDb, tmpDb, name := SwitchToTmpDb()
//	defer DropTmpDB(baseDb, tmpDb, name)
//
//	op := Operation{
//		Station: "a station",
//		User:    "operator@dev.com",
//	}
//	start := CaptureStart{
//		Operation:     op,
//		FileName:      "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
//		CaptureID:     "this.is.capture.id",
//		CaptureSource: "mlpcap",
//	}
//	if err := CaptureStart(start); err != nil {
//		t.Error("CaptureStart should succeed.", err)
//	}
//	start = CaptureStart{
//		Operation:     op,
//		FileName:      "heb_o_rav_bs-igeret_2016-09-14_lesson.mp4",
//		CaptureID:     "this.is.capture.id",
//		CaptureSource: "mlpcap",
//	}
//	if err := CaptureStart(start); err != nil {
//		t.Error("CaptureStart should succeed.", err)
//	}
//
//	sha1 := "abcd1234abcd1234abcd1234abcd1234abcd1234"
//	var size uint64 = 123
//	part := "full"
//	stopOp := Operation{
//		Station: "a station",
//		User:    "operator@dev.com",
//	}
//	stop := CaptureStop{
//		CaptureStart: CaptureStart{
//			Operation:     stopOp,
//			FileName:      "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
//			CaptureID:     "this.is.capture.id",
//			CaptureSource: "mltcap",
//		},
//		Sha1: sha1,
//		Size: size,
//		Part: part,
//	}
//	if err := CaptureStop(stop); err != nil {
//		t.Error("CaptureStop should succeed.", err)
//	}
//
//	f := models.File{Name: "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4"}
//	db.Where(&f).First(&f)
//	if !f.Sha1.Valid || hex.EncodeToString([]byte(f.Sha1.String)) != sha1 {
//		t.Error(fmt.Sprintf("Expected size %d got %d", size, f.Size))
//	}
//	if f.Size != size {
//		t.Error(fmt.Sprintf("Expected size %d got %d", size, f.Size))
//	}
//
//	var cu models.ContentUnit
//	db.Model(&f).Related(&cu, "ContentUnit")
//
//	var ccu models.CollectionsContentUnit
//	db.Model(&cu).Related(&ccu, "CollectionsContentUnit")
//	if ccu.Name != part {
//		t.Error(fmt.Sprintf("Expected part %s, got %s", part, ccu.Name))
//	}
//
//	stop = CaptureStop{
//		CaptureStart: CaptureStart{
//			Operation:     stopOp,
//			FileName:      "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
//			CaptureID:     "this.is.capture.id",
//			CaptureSource: "mltcap",
//		},
//		Sha1: "1111111111111111111111111111111111111111",
//		Size: 111,
//	}
//	if err := CaptureStop(stop); err == nil ||
//		strings.Contains(err.Error(), "CaptureStop File already has different Sha1") {
//		t.Error("Expected to fail with wrong Sha1 on CaptureStop: ", err)
//	}
//}
//
//
//func TestDemux(t *testing.T) {
//	baseDb, tmpDb, name := SwitchToTmpDb()
//	defer DropTmpDB(baseDb, tmpDb, name)
//
//	// Prepare file
//	sha1 := "abcdef123456"
//	if err := AddTestFile("lang_o_norav_part-a_2016-12-31_source.mp4", sha1, 111); err != nil {
//		t.Error("Could not create file.", err)
//	}
//
//	origFileName := "lang_o_norav_part-a_2016-12-31_orig.mp4"
//	proxyFileName := "lang_o_norav_part-a_2016-12-31_proxy.mp4"
//	demux := Demux{
//		Operation: Operation{
//			Station: "a station",
//			User:    "operator@dev.com",
//		},
//		FileKey: FileKey{
//			Sha1: "abcdef123456",
//		},
//		Original: FileUpdate{
//			FileKey: FileKey{
//				FileName: origFileName,
//				Sha1:     "aaaaaa111111",
//			},
//			Size: 111,
//		},
//		Proxy: FileUpdate{
//			FileKey: FileKey{
//				FileName: proxyFileName,
//				Sha1:     "bbbbbb222222",
//			},
//			Size: 222,
//		},
//	}
//	if err := Demux(demux); err != nil {
//		t.Error("Demux should succeed.", err)
//	}
//
//	orig := models.File{Name: origFileName}
//	db.Where(&orig).First(&orig)
//	if !orig.Sha1.Valid || hex.EncodeToString([]byte(orig.Sha1.String)) != demux.Original.Sha1 {
//		t.Error(fmt.Sprintf("Expected sha1 %s got %s", demux.Original.Sha1, orig.Sha1.String))
//	}
//	if orig.Size != demux.Original.Size {
//		t.Error(fmt.Sprintf("Expected size %d got %d", demux.Original.Size, orig.Size))
//	}
//	proxy := models.File{Name: proxyFileName}
//	db.Where(&proxy).First(&proxy)
//	if !proxy.Sha1.Valid || hex.EncodeToString([]byte(proxy.Sha1.String)) != demux.Proxy.Sha1 {
//		t.Error(fmt.Sprintf("Expected sha1 %s got %s", demux.Proxy.Sha1, proxy.Sha1.String))
//	}
//	if proxy.Size != demux.Proxy.Size {
//		t.Error(fmt.Sprintf("Expected size %d got %d", demux.Proxy.Size, proxy.Size))
//	}
//	sha1NullString, _ := Sha1ToNullString(sha1)
//	source := models.File{Sha1: sha1NullString}
//	db.Where(&source).First(&source)
//	if uint64(proxy.ParentID.Int64) != source.ID {
//		t.Error(fmt.Sprintf("Bad proxy parent id %d, expected %d", proxy.ParentID, source.ID))
//	}
//}
//
//func TestSend(t *testing.T) {
//	baseDb, tmpDb, name := SwitchToTmpDb()
//	defer DropTmpDB(baseDb, tmpDb, name)
//
//	// Prepare file
//	sha1 := "abcdef123456"
//	if err := AddTestFile("lang_o_norav_part-a_2016-12-31_source.mp4", sha1, 111); err != nil {
//		t.Error("Could not create file.", err)
//	}
//
//	destFileName := "lang_o_norav_part-a_2016-12-31_dest.mp4"
//	send := Send{
//		Operation: Operation{
//			Station: "a station",
//			User:    "operator@dev.com",
//		},
//		FileKey: FileKey{
//			Sha1: sha1,
//		},
//		Dest: FileUpdate{
//			FileKey: FileKey{
//				FileName: destFileName,
//				Sha1:     "cccccc333333",
//			},
//			Size: 333,
//		},
//	}
//	if err := Send(send); err != nil {
//		t.Error("Send should succeed.", err)
//	}
//
//	dest := models.File{Name: destFileName}
//	db.Where(&dest).First(&dest)
//	if !dest.Sha1.Valid || hex.EncodeToString([]byte(dest.Sha1.String)) != send.Dest.Sha1 {
//		t.Error(fmt.Sprintf("Expected sha1 %s got %s", send.Dest.Sha1, dest.Sha1.String))
//	}
//	if dest.Size != send.Dest.Size {
//		t.Error(fmt.Sprintf("Expected size %d got %d", send.Dest.Size, dest.Size))
//	}
//}
//
//func TestUpload(t *testing.T) {
//	baseDb, tmpDb, name := SwitchToTmpDb()
//	defer DropTmpDB(baseDb, tmpDb, name)
//
//	// Prepare file
//	sha1 := "abcdef123456"
//	fileName := "lang_o_norav_part-a_2016-12-31_source.mp4"
//	if err := AddTestFile(fileName, sha1, 111); err != nil {
//		t.Error("Could not create file.", err)
//	}
//
//	o := Operation{
//		Station: "a station",
//		User:    "operator@dev.com",
//	}
//
//	url := "http://this/is/some/url"
//	upload := Upload{
//		Operation: o,
//		FileUpdate: FileUpdate{
//			FileKey: FileKey{
//				FileName: fileName,
//				Sha1:     sha1,
//			},
//			Size: 111,
//		},
//		Url: url,
//		Existing: FileKey{
//			Sha1: "1234",
//		},
//	}
//	if err := Upload(upload); err != nil {
//		t.Error("Upload should succeed.", err)
//	}
//	f := models.File{Name: fileName}
//	db.Where(&f).First(&f)
//	if f.Properties["url"] != url {
//		t.Error(fmt.Sprintf("Expected url %s got %s", url, f.Properties["url"]))
//	}
//
//	// Upload non existing file.
//	otherUrl := "http://some/other/url"
//	upload = Upload{
//		Operation: o,
//		FileUpdate: FileUpdate{
//			FileKey: FileKey{
//				FileName: "some_name",
//				Sha1:     "abcd",
//			},
//			Size: 111,
//		},
//		Url: otherUrl,
//	}
//	if err := Upload(upload); err != nil {
//		t.Error("Upload should succeed.", err)
//	}
//	f = models.File{Name: "some_name"}
//	db.Where(&f).First(&f)
//	if f.Properties["url"] != otherUrl {
//		t.Error(fmt.Sprintf("Expected url %s got %s", otherUrl, f.Properties["url"]))
//	}
//}
//
func TestMain(m *testing.M) {
	rand.Seed(time.Now().UTC().UnixNano())
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
