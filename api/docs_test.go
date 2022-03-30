package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"testing"
	"time"

	"github.com/adams-sarah/test2doc/test"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/null/v8"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/events"
	"github.com/Bnei-Baruch/mdb/permissions"
	"github.com/Bnei-Baruch/mdb/utils"
)

type DocsSuite struct {
	suite.Suite
	utils.TestDBManager
	router     *gin.Engine
	testServer *test.Server
}

// For now leave this empty. When needed make sure code
// extracts variables from request url.
func Vars(req *http.Request) map[string]string {
	return make(map[string]string)
}

func (suite *DocsSuite) SetupSuite() {
	suite.Require().Nil(suite.InitTestDB())
	suite.Require().Nil(common.InitTypeRegistries(suite.DB))
	//suite.Require().Nil(InitTypeRegistries(boil.GetDB()))

	enforcer := permissions.NewEnforcer()
	enforcer.EnableEnforce(false)

	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.router.Use(
		utils.EnvMiddleware(suite.DB, new(events.NoopEmitter), enforcer, nil),
		utils.ErrorHandlingMiddleware(),
		gin.Recovery())
	SetupRoutes(suite.router)

	test.RegisterURLVarExtractor(Vars)
	var err error
	suite.testServer, err = test.NewServer(suite.router)
	if err != nil {
		panic(err.Error())
	}
}

func (suite *DocsSuite) TearDownSuite() {
	suite.testServer.Finish()
	suite.Require().Nil(suite.DestroyTestDB())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestDocs(t *testing.T) {
	suite.Run(t, new(DocsSuite))
}

func (suite *DocsSuite) Test1CaptureStartHandler() {
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

	resp, err := suite.testOperation(common.OP_CAPTURE_START, input)
	suite.Require().Nil(err)
	suite.assertJsonOK(resp)
}

func (suite *DocsSuite) Test2CaptureStopHandler() {
	input := CaptureStopRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: "c12356789",
		},
		File: File{
			FileName:  "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
			Sha1:      "012356789abcdef012356789abcdef0123456789",
			Size:      98737,
			CreatedAt: &Timestamp{Time: time.Now()},
			Type:      "type",
			SubType:   "subtype",
			MimeType:  "mime_type",
			Language:  common.LANG_MULTI,
		},
		CaptureSource: "mltcap",
		CollectionUID: "abcdefgh",
		Part:          "part",
	}

	resp, err := suite.testOperation(common.OP_CAPTURE_STOP, input)
	suite.Require().Nil(err)
	suite.assertJsonOK(resp)
}

func (suite *DocsSuite) Test3DemuxHandler() {
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
				CreatedAt: &Timestamp{Time: time.Now()},
				Type:      "type",
				SubType:   "subtype",
				MimeType:  "mime_type",
				Language:  common.LANG_MULTI,
			},
			Duration: 892.1900,
		},
		Proxy: &AVFile{
			File: File{
				FileName:  "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_p.mp4",
				Sha1:      "0987654321fedcba0987654321fedcba87654321",
				Size:      837,
				CreatedAt: &Timestamp{Time: time.Now()},
				Type:      "type",
				SubType:   "subtype",
				MimeType:  "mime_type",
				Language:  common.LANG_HEBREW,
			},
			Duration: 892.1900,
		},
	}

	resp, err := suite.testOperation(common.OP_DEMUX, input)
	suite.Require().Nil(err)
	suite.assertJsonOK(resp)
}

func (suite *DocsSuite) Test4TrimHandler() {
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
				CreatedAt: &Timestamp{Time: time.Now()},
				Type:      "type",
				SubType:   "subtype",
				MimeType:  "mime_type",
				Language:  common.LANG_MULTI,
			},
			Duration: 871,
		},
		Proxy: &AVFile{
			File: File{
				FileName:  "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_p_trim.mp4",
				Sha1:      "0987654321fedcba0987654321fedcba22222222",
				Size:      694,
				CreatedAt: &Timestamp{Time: time.Now()},
				Type:      "type",
				SubType:   "subtype",
				MimeType:  "mime_type",
				Language:  common.LANG_HEBREW,
			},
			Duration: 871,
		},
		In:  []float64{0.00, 198.23},
		Out: []float64{10.50, 207.31},
	}

	resp, err := suite.testOperation(common.OP_TRIM, input)
	suite.Require().Nil(err)
	suite.assertJsonOK(resp)
}

func (suite *DocsSuite) Test5SendHandler() {
	input := SendRequest{
		Operation: Operation{
			Station:    "Trim station",
			User:       "operator@dev.com",
			WorkflowID: "t12356789",
		},
		Original: Rename{
			Sha1:     "0987654321fedcba0987654321fedcba11111111",
			FileName: "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_rename_o.mp4",
		},
		Proxy: &Rename{
			Sha1:     "0987654321fedcba0987654321fedcba22222222",
			FileName: "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_rename_p.mp4",
		},
		Metadata: CITMetadata{
			ContentType:    common.CT_LESSON_PART,
			FinalName:      "final_name",
			AutoName:       "auto_name",
			ManualName:     "manual_name",
			CaptureDate:    Date{Time: time.Now()},
			Language:       "heb",
			Lecturer:       "rav",
			HasTranslation: true,
			RequireTest:    false,
			Number:         null.IntFrom(1),
			Part:           null.IntFrom(2),
			Sources:        []string{"12345678", "87654321", "abcdefgh"},
			Tags:           []string{"12345678", "87654321"},
			Major:          &CITMetadataMajor{Type: "source", Idx: 1},
		},
	}

	resp, err := suite.testOperation(common.OP_SEND, input)
	suite.Require().Nil(err)

	suite.Equal(http.StatusOK, resp.StatusCode, "HTTP status")
	suite.Equal("application/json; charset=utf-8", resp.Header.Get("Content-Type"), "HTTP Content-Type")

	var body map[string]interface{}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&body)
	suite.Require().Nil(err)
	suite.NotNil(body["id"])
	suite.NotNil(body["uid"])
	suite.Nil(body["errors"], "HTTP body.errors")
}

func (suite *DocsSuite) Test6ConvertHandler() {
	input := ConvertRequest{
		Operation: Operation{
			Station: "Convert station",
			User:    "operator@dev.com",
		},
		Sha1: "0987654321fedcba0987654321fedcba11111111",
		Output: []AVFile{
			{
				File: File{
					FileName:  "heb_file.mp4",
					Sha1:      "0987654321fedcba0987654321fedcba33333333",
					Size:      694,
					CreatedAt: &Timestamp{Time: time.Now()},
					Type:      "type",
					SubType:   "subtype",
					MimeType:  "mime_type",
					Language:  common.LANG_HEBREW,
				},
				Duration: 871,
			},
			{
				File: File{
					FileName:  "eng_file.mp4",
					Sha1:      "0987654321fedcba0987654321fedcba44444444",
					Size:      694,
					CreatedAt: &Timestamp{Time: time.Now()},
					Type:      "type",
					SubType:   "subtype",
					MimeType:  "mime_type",
					Language:  common.LANG_ENGLISH,
				},
				Duration: 871,
			},
			{
				File: File{
					FileName:  "rus_file.mp4",
					Sha1:      "0987654321fedcba0987654321fedcba55555555",
					Size:      694,
					CreatedAt: &Timestamp{Time: time.Now()},
					Type:      "type",
					SubType:   "subtype",
					MimeType:  "mime_type",
					Language:  common.LANG_RUSSIAN,
				},
				Duration: 871,
			},
		},
	}

	resp, err := suite.testOperation(common.OP_CONVERT, input)
	suite.Require().Nil(err)
	suite.assertJsonOK(resp)
}

func (suite *DocsSuite) Test7UploadHandler() {
	input := UploadRequest{
		Operation: Operation{
			Station: "Upload station",
			User:    "operator1@dev.com",
		},
		AVFile: AVFile{
			File: File{
				FileName:  "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_o.mp4",
				Sha1:      "0987654321fedcba0987654321fedcba09876543",
				Size:      19837,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892,
		},
		Url: "https://example.com/heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
	}

	resp, err := suite.testOperation(common.OP_UPLOAD, input)
	suite.Require().Nil(err)
	suite.assertJsonOK(resp)
}

func (suite *DocsSuite) Test8SirtutimHandler() {
	input := SirtutimRequest{
		Operation: Operation{
			Station: "Upload station",
			User:    "operator1@dev.com",
		},
		File: File{
			FileName:  "heb_o_rav_2016-09-14_lesson_o.zip",
			Sha1:      "0987654321fedcba0987654321fedcba09876544",
			Size:      19837,
			CreatedAt: &Timestamp{Time: time.Now()},
			Language:  common.LANG_HEBREW,
		},
		OriginalSha1: "0987654321fedcba0987654321fedcba11111111",
	}

	resp, err := suite.testOperation(common.OP_SIRTUTIM, input)
	suite.Require().Nil(err)
	suite.assertJsonOK(resp)
}

func (suite *DocsSuite) Test9InsertHandler() {
	tx, err := suite.DB.Begin()
	suite.Require().Nil(err)
	cu, err := CreateContentUnit(tx, common.CT_LESSON_PART, nil)
	suite.Require().Nil(err)
	err = tx.Commit()
	suite.Require().Nil(err)

	input := InsertRequest{
		Operation: Operation{
			Station: "Insert station",
			User:    "operator1@dev.com",
		},
		InsertType:     "akladot",
		ContentUnitUID: cu.UID,
		AVFile: AVFile{
			File: File{
				FileName:  "heb_o_rav_2016-09-14_lesson_akladot.docx",
				Sha1:      "0987654321fedcba0987654321fedcba09876555",
				Size:      19837,
				CreatedAt: &Timestamp{Time: time.Now()},
				Language:  common.LANG_HEBREW,
			},
		},
		ParentSha1: "0987654321fedcba0987654321fedcba11111111",
		Mode:       "new",
	}

	resp, err := suite.testOperation(common.OP_INSERT, input)
	suite.Require().Nil(err)
	suite.assertJsonOK(resp)
}

func (suite *DocsSuite) Test91InsertHandlerNewUnit() {
	input := InsertRequest{
		Operation: Operation{
			Station: "Insert station",
			User:    "operator1@dev.com",
		},
		InsertType: "declamation",
		AVFile: AVFile{
			File: File{
				FileName:  "rus_o_norav_2016-09-14_declamation.mp3",
				Sha1:      "0987654321fedcba0987654321fedcba09875555",
				Size:      19837,
				CreatedAt: &Timestamp{Time: time.Now()},
				Language:  common.LANG_HEBREW,
			},
		},
		Mode: "new",
		Metadata: &CITMetadata{
			FilmDate:    &Date{time.Now()},
			ContentType: common.CT_BLOG_POST,
			FinalName:   "final_name",
			Language:    "rus",
			Lecturer:    "norav",
		},
	}

	resp, err := suite.testOperation(common.OP_INSERT, input)
	suite.Require().Nil(err)

	suite.Equal(http.StatusOK, resp.StatusCode, "HTTP status")
	suite.Equal("application/json; charset=utf-8", resp.Header.Get("Content-Type"), "HTTP Content-Type")

	var body map[string]interface{}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&body)
	suite.Require().Nil(err)
	suite.NotNil(body["id"])
	suite.NotNil(body["uid"])
	suite.Nil(body["errors"], "HTTP body.errors")
}

func (suite *DocsSuite) Test92TranscodeHandler() {
	input := TranscodeRequest{
		Operation: Operation{
			Station: "Insert station",
			User:    "operator1@dev.com",
		},
		OriginalSha1: "0987654321fedcba0987654321fedcba11111111",
		MaybeFile: MaybeFile{
			FileName:  "heb_o_rav_2016-09-14_lesson_akladot.mp4",
			Sha1:      "0987654321fedcba0987654321fedcba09876666",
			Size:      19837,
			CreatedAt: &Timestamp{Time: time.Now()},
		},
	}

	resp, err := suite.testOperation(common.OP_TRANSCODE, input)
	suite.Require().Nil(err)
	suite.assertJsonOK(resp)
}

func (suite *DocsSuite) Test922TranscodeHandlerError() {
	input := TranscodeRequest{
		Operation: Operation{
			Station: "Insert station",
			User:    "operator1@dev.com",
		},
		OriginalSha1: "0987654321fedcba0987654321fedcba11111111",
		Message:      "Some transcoding error message",
	}

	resp, err := suite.testOperation(common.OP_TRANSCODE, input)
	suite.Require().Nil(err)
	suite.assertJsonOK(resp)
}

func (suite *DocsSuite) Test93JoinHandler() {
	input := JoinRequest{
		Operation: Operation{
			Station:    "Join station",
			User:       "operator@dev.com",
			WorkflowID: "d12356789",
		},
		OriginalShas: []string{"0987654321fedcba0987654321fedcba09876543"},
		ProxyShas:    []string{"0987654321fedcba0987654321fedcba87654321"},
		Original: AVFile{
			File: File{
				FileName:  "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_o_trim.mp4",
				Sha1:      "0987654321fedcba0987654321fedcba11111113",
				Size:      19800,
				CreatedAt: &Timestamp{Time: time.Now()},
				Type:      "type",
				SubType:   "subtype",
				MimeType:  "mime_type",
				Language:  common.LANG_MULTI,
			},
			Duration: 871,
		},
		Proxy: &AVFile{
			File: File{
				FileName:  "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_p_trim.mp4",
				Sha1:      "0987654321fedcba0987654321fedcba22222223",
				Size:      694,
				CreatedAt: &Timestamp{Time: time.Now()},
				Type:      "type",
				SubType:   "subtype",
				MimeType:  "mime_type",
				Language:  common.LANG_HEBREW,
			},
			Duration: 871,
		},
	}

	resp, err := suite.testOperation(common.OP_JOIN, input)
	suite.Require().Nil(err)
	suite.assertJsonOK(resp)
}

func (suite *DocsSuite) testOperation(name string, input interface{}) (*http.Response, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(input)
	suite.Require().Nil(err)
	u, _ := url.Parse(suite.testServer.URL)
	u.Path = path.Join(u.Path, "operations", name)
	return http.Post(u.String(), "application/json", b)
}

func (suite *DocsSuite) assertJsonOK(resp *http.Response) {
	suite.Equal(http.StatusOK, resp.StatusCode, "HTTP status")
	suite.Equal("application/json; charset=utf-8", resp.Header.Get("Content-Type"), "HTTP Content-Type")

	var body map[string]interface{}
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(&body)
	suite.Require().Nil(err)
	suite.Equal("ok", body["status"], "HTTP body.status")
	//suite.T().Log(body)
	suite.Nil(body["errors"], "HTTP body.errors")
}
