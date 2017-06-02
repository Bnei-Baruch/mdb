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
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/nullbio/null.v6"

	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/vattle/sqlboiler/boil"
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
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.router.Use(utils.ErrorHandlingMiddleware(), gin.Recovery())
	SetupRoutes(suite.router)

	test.RegisterURLVarExtractor(Vars)
	var err error
	suite.testServer, err = test.NewServer(suite.router)
	if err != nil {
		panic(err.Error())
	}

	suite.Require().Nil(suite.InitTestDB())
	suite.Require().Nil(InitTypeRegistries(boil.GetDB()))
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

	resp, err := suite.testOperation(OP_CAPTURE_START, input)
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
			Language:  LANG_MULTI,
		},
		CaptureSource: "mltcap",
		CollectionUID: "abcdefgh",
		Part:          "part",
	}

	resp, err := suite.testOperation(OP_CAPTURE_STOP, input)
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
				Language:  LANG_MULTI,
			},
			Duration: 892.1900,
		},
		Proxy: AVFile{
			File: File{
				FileName:  "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_p.mp4",
				Sha1:      "0987654321fedcba0987654321fedcba87654321",
				Size:      837,
				CreatedAt: &Timestamp{Time: time.Now()},
				Type:      "type",
				SubType:   "subtype",
				MimeType:  "mime_type",
				Language:  LANG_HEBREW,
			},
			Duration: 892.1900,
		},
	}

	resp, err := suite.testOperation(OP_DEMUX, input)
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
				Language:  LANG_MULTI,
			},
			Duration: 871,
		},
		Proxy: AVFile{
			File: File{
				FileName:  "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_p_trim.mp4",
				Sha1:      "0987654321fedcba0987654321fedcba22222222",
				Size:      694,
				CreatedAt: &Timestamp{Time: time.Now()},
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

	resp, err := suite.testOperation(OP_TRIM, input)
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
		Proxy: Rename{
			Sha1:     "0987654321fedcba0987654321fedcba22222222",
			FileName: "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson_rename_p.mp4",
		},
		Metadata: CITMetadata{
			ContentType:    CT_LESSON_PART,
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

	resp, err := suite.testOperation(OP_SEND, input)
	suite.Require().Nil(err)
	suite.assertJsonOK(resp)
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
					Language:  LANG_HEBREW,
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
					Language:  LANG_ENGLISH,
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
					Language:  LANG_RUSSIAN,
				},
				Duration: 871,
			},
		},
	}

	resp, err := suite.testOperation(OP_CONVERT, input)
	suite.Require().Nil(err)
	suite.assertJsonOK(resp)
}

func (suite *DocsSuite) Test7UploadHandler() {
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
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892,
		},
		Url: "https://example.com/heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
	}

	resp, err := suite.testOperation(OP_UPLOAD, input)
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
