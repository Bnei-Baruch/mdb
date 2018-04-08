package api

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gopkg.in/volatiletech/null.v6"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

type HandlersSuite struct {
	suite.Suite
	utils.TestDBManager
	tx *sql.Tx
}

func (suite *HandlersSuite) SetupSuite() {
	suite.Require().Nil(suite.InitTestDB())
	suite.Require().Nil(InitTypeRegistries(suite.DB))
}

func (suite *HandlersSuite) TearDownSuite() {
	suite.Require().Nil(suite.DestroyTestDB())
}

func (suite *HandlersSuite) SetupTest() {
	var err error
	suite.tx, err = suite.DB.Begin()
	suite.Require().Nil(err)
}

func (suite *HandlersSuite) TearDownTest() {
	err := suite.tx.Rollback()
	suite.Require().Nil(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestHandlers(t *testing.T) {
	suite.Run(t, new(HandlersSuite))
}

func (suite *HandlersSuite) TestHandleCaptureStart() {
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

	op, evnts, err := handleCaptureStart(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_CAPTURE_START].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	var props map[string]interface{}
	err = op.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.Operation.WorkflowID, props["workflow_id"], "properties: workflow_id")
	suite.Equal(input.CaptureSource, props["capture_source"], "properties: capture_source")
	suite.Equal(input.CollectionUID, props["collection_uid"], "properties: collection_uid")

	// Check user
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op))
	suite.Equal(input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 1, "Number of files")
	f := op.R.Files[0]
	suite.Equal(input.FileName, f.Name, "File: Name")
	suite.False(f.Sha1.Valid, "File: SHA1")
}

func (suite *HandlersSuite) TestHandleCaptureStop() {
	// Prepare capture_start operation
	opStart, evnts, err := handleCaptureStart(suite.tx, CaptureStartRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: "c12356789",
		},
		FileName:      "heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4",
		CaptureSource: "mltcap",
		CollectionUID: "abcdefgh",
	})
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)
	suite.Require().Nil(opStart.L.LoadFiles(suite.tx, true, opStart))
	parent := opStart.R.Files[0]

	// Do capture_stop operation
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

	op, evnts, err := handleCaptureStop(suite.tx, input)
	suite.Require().Nil(err)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_CAPTURE_STOP].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	var props map[string]interface{}
	err = op.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.Operation.WorkflowID, props["workflow_id"], "properties: workflow_id")
	suite.Equal(input.CaptureSource, props["capture_source"], "properties: capture_source")
	suite.Equal(input.CollectionUID, props["collection_uid"], "properties: collection_uid")
	suite.Equal(input.Part, props["part"], "properties: part")

	// Check user
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op))
	suite.Equal(input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 1, "Number of files")
	f := op.R.Files[0]
	suite.Equal(input.FileName, f.Name, "File: Name")
	suite.Equal(input.Sha1, hex.EncodeToString(f.Sha1.Bytes), "File: SHA1")
	suite.Equal(input.Size, f.Size, "File: Size")
	suite.Equal(input.CreatedAt.Time.Unix(), f.FileCreatedAt.Time.Unix(), "File: FileCreatedAt")
	suite.Equal(parent.ID, f.ParentID.Int64, "File Parent.ID")
	suite.False(f.Properties.Valid, "properties")
}

func (suite *HandlersSuite) TestHandleDemux() {
	// Create dummy parent file
	fi := File{
		FileName:  "dummy parent file",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef0123456789",
		Size:      math.MaxInt64,
	}
	_, err := CreateFile(suite.tx, nil, fi, nil)
	suite.Require().Nil(err)

	// Do demux operation
	input := DemuxRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		Sha1: fi.Sha1,
		Original: AVFile{
			File: File{
				FileName:  "original.mp4",
				Sha1:      "012356789abcdef012356789abcdef9876543210",
				Size:      98737,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		Proxy: AVFile{
			File: File{
				FileName:  "proxy.mp4",
				Sha1:      "987653210abcdef012356789abcdef9876543210",
				Size:      987,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 891.8800,
		},
		CaptureSource: "mltcap",
	}

	op, evnts, err := handleDemux(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_DEMUX].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	var props map[string]interface{}
	err = op.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.CaptureSource, props["capture_source"], "properties: capture_source")

	// Check user
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op))
	suite.Equal(input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 3, "Number of files")
	fm := make(map[string]*models.File)
	for _, x := range op.R.Files {
		fm[x.Name] = x
	}

	parent := fm[fi.FileName]

	// Check original
	original := fm[input.Original.FileName]
	suite.Equal(input.Original.FileName, original.Name, "Original: Name")
	suite.Equal(input.Original.Sha1, hex.EncodeToString(original.Sha1.Bytes), "Original: SHA1")
	suite.Equal(input.Original.Size, original.Size, "Original: Size")
	suite.Equal(input.Original.CreatedAt.Time.Unix(), original.FileCreatedAt.Time.Unix(), "Original: FileCreatedAt")
	suite.Equal(parent.ID, original.ParentID.Int64, "Original Parent.ID")
	err = original.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.Original.Duration, props["duration"], "Original props: duration")

	// Check proxy
	proxy := fm[input.Proxy.FileName]
	suite.Equal(input.Proxy.FileName, proxy.Name, "Proxy: Name")
	suite.Equal(input.Proxy.Sha1, hex.EncodeToString(proxy.Sha1.Bytes), "Proxy: SHA1")
	suite.Equal(input.Proxy.Size, proxy.Size, "Proxy: Size")
	suite.Equal(input.Proxy.CreatedAt.Time.Unix(), proxy.FileCreatedAt.Time.Unix(), "Proxy: FileCreatedAt")
	suite.Equal(parent.ID, proxy.ParentID.Int64, "Proxy Parent.ID")
	err = proxy.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.Proxy.Duration, props["duration"], "Proxy props: duration")
}

func (suite *HandlersSuite) TestHandleTrim() {
	// Create dummy original and proxy parent files
	ofi := File{
		FileName:  "dummy original parent file",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef9876543210",
		Size:      math.MaxInt64,
	}
	_, err := CreateFile(suite.tx, nil, ofi, nil)
	suite.Require().Nil(err)

	pfi := File{
		FileName:  "dummy proxy parent file",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "987653210abcdef012356789abcdef9876543210",
		Size:      math.MaxInt64,
	}
	_, err = CreateFile(suite.tx, nil, pfi, nil)
	suite.Require().Nil(err)

	// Do trim operation
	input := TrimRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		OriginalSha1: ofi.Sha1,
		ProxySha1:    pfi.Sha1,
		Original: AVFile{
			File: File{
				FileName:  "original_trim.mp4",
				Sha1:      "012356789abcdef012356789abcdef1111111111",
				Size:      98737,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		Proxy: AVFile{
			File: File{
				FileName:  "proxy_trim.mp4",
				Sha1:      "987653210abcdef012356789abcdef2222222222",
				Size:      987,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 891.8800,
		},
		CaptureSource: "mltcap",
		In:            []float64{10.05, 249.43},
		Out:           []float64{240.51, 899.27},
	}

	op, evnts, err := handleTrim(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_TRIM].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	var props map[string]interface{}
	err = op.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.CaptureSource, props["capture_source"], "properties: capture_source")
	for i, v := range input.In {
		suite.Equal(v, props["in"].([]interface{})[i], "properties: in[%d]", i)
	}
	for i, v := range input.Out {
		suite.Equal(v, props["out"].([]interface{})[i], "properties: out[%d]", i)
	}

	// Check user
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op))
	suite.Equal(input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 4, "Number of files")
	fm := make(map[string]*models.File)
	for _, x := range op.R.Files {
		fm[x.Name] = x
	}

	originalParent := fm[ofi.FileName]
	proxyParent := fm[pfi.FileName]

	// Check original
	original := fm[input.Original.FileName]
	suite.Equal(input.Original.FileName, original.Name, "Original: Name")
	suite.Equal(input.Original.Sha1, hex.EncodeToString(original.Sha1.Bytes), "Original: SHA1")
	suite.Equal(input.Original.Size, original.Size, "Original: Size")
	suite.Equal(input.Original.CreatedAt.Time.Unix(), original.FileCreatedAt.Time.Unix(), "Original: FileCreatedAt")
	suite.Equal(originalParent.ID, original.ParentID.Int64, "Original Parent.ID")
	err = original.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.Original.Duration, props["duration"], "Original props: duration")

	// Check proxy
	proxy := fm[input.Proxy.FileName]
	suite.Equal(input.Proxy.FileName, proxy.Name, "Proxy: Name")
	suite.Equal(input.Proxy.Sha1, hex.EncodeToString(proxy.Sha1.Bytes), "Proxy: SHA1")
	suite.Equal(input.Proxy.Size, proxy.Size, "Proxy: Size")
	suite.Equal(input.Proxy.CreatedAt.Time.Unix(), proxy.FileCreatedAt.Time.Unix(), "Proxy: FileCreatedAt")
	suite.Equal(proxyParent.ID, proxy.ParentID.Int64, "Proxy Parent.ID")
	err = proxy.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.Proxy.Duration, props["duration"], "Proxy props: duration")
}

func (suite *HandlersSuite) TestHandleSend() {
	// Create dummy original and proxy trimmed files
	ofi := File{
		FileName:  "dummy original trimmed file",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef9876543210",
		Size:      math.MaxInt64,
	}
	_, err := CreateFile(suite.tx, nil, ofi, nil)
	suite.Require().Nil(err)

	pfi := File{
		FileName:  "dummy proxy trimmed file",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "987653210abcdef012356789abcdef9876543210",
		Size:      math.MaxInt64,
	}
	_, err = CreateFile(suite.tx, nil, pfi, nil)
	suite.Require().Nil(err)

	// Do send operation
	input := SendRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		Original: Rename{
			Sha1:     ofi.Sha1,
			FileName: "original_renamed.mp4",
		},
		Proxy: Rename{
			Sha1:     pfi.Sha1,
			FileName: "proxy_renamed.mp4",
		},
		Metadata: CITMetadata{
			ContentType: CT_LESSON_PART,
		},
	}

	op, evnts, err := handleSend(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().NotEmpty(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_SEND].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	var props map[string]interface{}
	err = op.Properties.Unmarshal(&props)
	suite.Require().Nil(err)

	// Check user
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op))
	suite.Equal(input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 2, "Number of files")
	fm := make(map[string]*models.File)
	for _, x := range op.R.Files {
		fm[x.Name] = x
	}

	// Check original
	original := fm[input.Original.FileName]
	suite.Equal(input.Original.FileName, original.Name, "Original: Name")
	suite.Equal(input.Original.Sha1, hex.EncodeToString(original.Sha1.Bytes), "Original: SHA1")

	// Check proxy
	proxy := fm[input.Proxy.FileName]
	suite.Equal(input.Proxy.FileName, proxy.Name, "Proxy: Name")
	suite.Equal(input.Proxy.Sha1, hex.EncodeToString(proxy.Sha1.Bytes), "Proxy: SHA1")
}

func (suite *HandlersSuite) TestHandleConvert() {
	// Create dummy input file
	fi := File{
		FileName:  "dummy input file",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef9876543210",
		Size:      math.MaxInt64,
	}
	_, err := CreateFile(suite.tx, nil, fi, nil)
	suite.Require().Nil(err)

	// Do convert operation
	input := ConvertRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		Sha1: fi.Sha1,
		Output: []AVFile{
			{
				File: File{
					FileName:  "heb_file.mp4",
					Sha1:      "0987654321fedcba0987654321fedcba33333333",
					Size:      694,
					CreatedAt: &Timestamp{Time: time.Now()},
					Type:      "type1",
					SubType:   "subtype1",
					MimeType:  "mime_type1",
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
					Type:      "type2",
					SubType:   "subtype2",
					MimeType:  "mime_type2",
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
					Type:      "type3",
					SubType:   "subtype3",
					MimeType:  "mime_type3",
					Language:  LANG_RUSSIAN,
				},
				Duration: 871,
			},
			{
				File: File{
					FileName:  "duplicate_rus_file.mp4",
					Sha1:      "0987654321fedcba0987654321fedcba55555555",
					Size:      694,
					CreatedAt: &Timestamp{Time: time.Now()},
					Type:      "type3",
					SubType:   "subtype3",
					MimeType:  "mime_type3",
					Language:  LANG_RUSSIAN,
				},
				Duration: 871,
			},
		},
	}

	op, evnts, err := handleConvert(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().Empty(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_CONVERT].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	suite.False(op.Properties.Valid, "Operation properties")

	// Check user
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op))
	suite.Equal(input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 4, "Number of files")
	fm := make(map[string]*models.File)
	for _, x := range op.R.Files {
		fm[x.Name] = x
	}

	// Check input
	in := fm[fi.FileName]
	suite.Equal(fi.Sha1, hex.EncodeToString(in.Sha1.Bytes), "In: SHA1")

	// Check output
	var props map[string]interface{}
	for i, x := range input.Output[:len(input.Output)-2] {
		f := fm[x.FileName]
		suite.Equal(x.FileName, f.Name, "Output[%d]: Name", i)
		suite.Equal(x.Sha1, hex.EncodeToString(f.Sha1.Bytes), "Output[%d]: SHA1", i)
		suite.Equal(x.Size, f.Size, "Output[%d]: Size", i)
		suite.Equal(x.CreatedAt.Time.Unix(), f.FileCreatedAt.Time.Unix(), "Output[%d]: FileCreatedAt", i)
		suite.Equal(x.Type, f.Type, "Output[%d]: Type", i)
		suite.Equal(x.SubType, f.SubType, "Output[%d]: SubType", i)
		suite.Equal(x.MimeType, f.MimeType.String, "Output[%d]: MimeType", i)
		suite.Equal(x.Language, f.Language.String, "Output[%d]: Language", i)
		suite.Equal(in.ID, f.ParentID.Int64, "Output[%d] Parent.ID", i)
		err = f.Properties.Unmarshal(&props)
		suite.Require().Nil(err)
		suite.Equal(x.Duration, props["duration"], "Output[%d] props: duration", i)
	}

	// test "reconvert" update existing files
	efi := File{
		FileName:  "dummy input file",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef9999999999",
		Size:      math.MaxInt64,
	}
	_, err = CreateFile(suite.tx, nil, efi, nil)
	suite.Require().Nil(err)

	input.Output = []AVFile{
		{
			File: File{
				FileName:  "heb_file.mp4",
				Sha1:      "012356789abcdef012356789abcdef9999999999",
				Size:      math.MaxInt64,
				CreatedAt: &Timestamp{Time: time.Now()},
				Type:      "type1",
				SubType:   "subtype1",
				MimeType:  "mime_type1",
				Language:  LANG_HEBREW,
			},
			Duration: 871,
		},
	}
	op, evnts, err = handleConvert(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().NotEmpty(evnts)

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 2, "Number of files")
	fm = make(map[string]*models.File)
	for _, x := range op.R.Files {
		fm[x.Name] = x
	}

	// Check output
	for i, x := range input.Output {
		f := fm[x.FileName]
		suite.Equal(x.FileName, f.Name, "Output[%d]: Name", i)
		suite.Equal(x.Sha1, hex.EncodeToString(f.Sha1.Bytes), "Output[%d]: SHA1", i)
		suite.Equal(x.Size, f.Size, "Output[%d]: Size", i)
		suite.Equal(x.CreatedAt.Time.Unix(), f.FileCreatedAt.Time.Unix(), "Output[%d]: FileCreatedAt", i)
		suite.Equal(x.Type, f.Type, "Output[%d]: Type", i)
		suite.Equal(x.SubType, f.SubType, "Output[%d]: SubType", i)
		suite.Equal(x.MimeType, f.MimeType.String, "Output[%d]: MimeType", i)
		suite.Equal(x.Language, f.Language.String, "Output[%d]: Language", i)
		suite.Equal(in.ID, f.ParentID.Int64, "Output[%d] Parent.ID", i)
		err = f.Properties.Unmarshal(&props)
		suite.Require().Nil(err)
		suite.Equal(x.Duration, props["duration"], "Output[%d] props: duration", i)
	}
}

func (suite *HandlersSuite) TestHandleUpload() {
	// First seen, unknown, file
	input := UploadRequest{
		Operation: Operation{
			Station: "Upload station",
			User:    "operator@dev.com",
		},
		AVFile: AVFile{
			File: File{
				FileName:  "file.mp4",
				Sha1:      "0987654321fedcba0987654321fedcba33333333",
				Size:      694,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 871,
		},
		Url: "http://example.com/some/url/to/file.mp4",
	}

	op, evnts, err := handleUpload(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().NotEmpty(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_UPLOAD].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	suite.False(op.Properties.Valid, "Operation properties")

	// Check file
	err = op.L.LoadFiles(suite.tx, true, op)
	suite.Require().Nil(err)
	suite.Require().Len(op.R.Files, 1, "Operation Files length")
	f := op.R.Files[0]
	suite.Equal(input.FileName, f.Name, "File.Name")
	suite.Equal(input.Sha1, hex.EncodeToString(f.Sha1.Bytes), "File.SHA1")
	suite.Equal(input.Size, f.Size, "File.Size")
	suite.Equal(input.CreatedAt.Time.Unix(), f.FileCreatedAt.Time.Unix(), "File.FileCreatedAt")
	suite.False(f.ParentID.Valid, "File.ParentID")
	suite.True(f.Published, "File.Published")
	suite.Equal(SEC_PUBLIC, f.Secure, "File.Secure")
	var props map[string]interface{}
	err = f.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.Url, props["url"], "file props: url")
	suite.Equal(input.Duration, props["duration"], "file props: duration")

	// Existing file in a content unit and collection structure
	f2 := File{
		FileName:  "file.mp4",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef0123456789",
		Size:      math.MaxInt64,
	}
	file, err := CreateFile(suite.tx, nil, f2, nil)
	suite.Require().Nil(err)
	cu, err := CreateContentUnit(suite.tx, CT_LESSON_PART, nil)
	suite.Require().Nil(err)
	err = file.SetContentUnit(suite.tx, false, cu)
	suite.Require().Nil(err)
	c, err := CreateCollection(suite.tx, CT_DAILY_LESSON, nil)
	suite.Require().Nil(err)
	err = c.AddCollectionsContentUnits(suite.tx, true, &models.CollectionsContentUnit{ContentUnitID: cu.ID})
	suite.Require().Nil(err)

	input = UploadRequest{
		Operation: Operation{
			Station: "Upload station",
			User:    "operator@dev.com",
		},
		AVFile: AVFile{
			File:     f2,
			Duration: 871,
		},
		Url: "http://example.com/some/url/to/file.mp4",
	}

	op, evnts, err = handleUpload(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().NotEmpty(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_UPLOAD].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	suite.False(op.Properties.Valid, "Operation properties")

	// Check file
	err = op.L.LoadFiles(suite.tx, true, op)
	suite.Require().Nil(err)
	suite.Require().Len(op.R.Files, 1, "Operation Files length")
	f = op.R.Files[0]
	suite.Equal(input.FileName, f.Name, "File.Name")
	suite.Equal(input.Sha1, hex.EncodeToString(f.Sha1.Bytes), "File.SHA1")
	suite.Equal(input.Size, f.Size, "File.Size")
	suite.Equal(input.CreatedAt.Time.Unix(), f.FileCreatedAt.Time.Unix(), "File.FileCreatedAt")
	suite.False(f.ParentID.Valid, "File.ParentID")
	suite.True(f.Published, "File.Published")
	suite.Equal(SEC_PUBLIC, f.Secure, "File.Secure")
	err = f.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.Url, props["url"], "file props: url")
	suite.Equal(input.Duration, props["duration"], "file props: duration")

	// Check content unit and collection
	err = cu.Reload(suite.tx)
	suite.Require().Nil(err)
	suite.True(cu.Published)
	err = c.Reload(suite.tx)
	suite.Require().Nil(err)
	suite.True(c.Published)
}

func (suite *HandlersSuite) TestHandleSirtutim() {
	// Create dummy original parent files
	ofi := File{
		FileName:  "dummy original parent file",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef9876543210",
		Size:      math.MaxInt64,
	}
	original, err := CreateFile(suite.tx, nil, ofi, nil)
	suite.Require().Nil(err)

	// Create dummy content unit
	cu, err := CreateContentUnit(suite.tx, CT_LESSON_PART, nil)
	suite.Require().Nil(err)

	// associate original and content unit
	original.ContentUnitID = null.Int64From(cu.ID)
	err = original.Update(suite.tx, "content_unit_id")
	suite.Require().Nil(err)

	// Do sirtutim operation
	input := SirtutimRequest{
		Operation: Operation{
			Station: "Some station",
			User:    "operator@dev.com",
		},
		OriginalSha1: ofi.Sha1,
		File: File{
			FileName:  "sirtutim.zip",
			Sha1:      "012356789abcdef012356789abcdef1111111111",
			Size:      98737,
			CreatedAt: &Timestamp{Time: time.Now()},
		},
	}

	op, evnts, err := handleSirtutim(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_SIRTUTIM].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	suite.False(op.Properties.Valid, "properties.Valid")

	// Check user
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op))
	suite.Equal(input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 2, "Number of files")
	fm := make(map[string]*models.File)
	for _, x := range op.R.Files {
		fm[x.Name] = x
	}

	originalParent := fm[ofi.FileName]
	suite.Equal(original.ID, originalParent.ID, "original <-> operation")

	// Check sirtutim file
	sirtutim := fm[input.FileName]
	suite.Equal(input.FileName, sirtutim.Name, "sirtutim: Name")
	suite.Equal(input.Sha1, hex.EncodeToString(sirtutim.Sha1.Bytes), "sirtutim: SHA1")
	suite.Equal(input.Size, sirtutim.Size, "sirtutim: Size")
	suite.Equal(input.CreatedAt.Time.Unix(), sirtutim.FileCreatedAt.Time.Unix(), "sirtutim: FileCreatedAt")

	// check content unit association
	suite.True(sirtutim.ContentUnitID.Valid, "sirtutim: ContentUnitID.Valid")
	suite.Equal(cu.ID, sirtutim.ContentUnitID.Int64, "sirtutim: ContentUnitID.Int64")
}

func (suite *HandlersSuite) TestHandleInsert() {
	// Create dummy content unit
	cu, err := CreateContentUnit(suite.tx, CT_LESSON_PART, nil)
	suite.Require().Nil(err)

	// Do insert operation
	input := InsertRequest{
		Operation: Operation{
			Station:    "Some station",
			User:       "operator@dev.com",
			WorkflowID: "workflow_id",
		},
		InsertType:     "akladot",
		ContentUnitUID: cu.UID,
		AVFile: AVFile{
			File: File{
				FileName:  "akladot.doc",
				Sha1:      "012356789abcdef012356789abcdef1111111111",
				Size:      98737,
				CreatedAt: &Timestamp{Time: time.Now()},
				MimeType:  "application/msword",
				Language:  LANG_HEBREW,
			},
			Duration: 123.4,
		},
		Mode: "new",
	}

	op, evnts, err := handleInsert(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().NotEmpty(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_INSERT].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	var props map[string]interface{}
	err = op.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.Operation.WorkflowID, props["workflow_id"].(string), "Operation workflow_id")
	suite.Equal(input.InsertType, props["insert_type"].(string), "Operation insert_type")
	suite.Equal(input.Mode, props["mode"].(string), "Operation mode")

	// Check user
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op))
	suite.Equal(input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 1, "Number of files")

	// check inserted file
	f := op.R.Files[0]
	suite.Equal(input.FileName, f.Name, "File Name")
	suite.Equal(input.Sha1, hex.EncodeToString(f.Sha1.Bytes), "File SHA1")
	suite.Equal(input.Size, f.Size, "File Size")
	suite.Equal(input.CreatedAt.Time.Unix(), f.FileCreatedAt.Time.Unix(), "File FileCreatedAt")
	suite.Equal("text", f.Type, "File Type")
	suite.Equal(input.MimeType, f.MimeType.String, "File MimeType")
	suite.Equal(input.Language, f.Language.String, "File Language")
	suite.False(f.ParentID.Valid, "File ParentID")
	err = f.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.InsertType, props["insert_type"].(string), "File insert_type")
	suite.Equal(input.AVFile.Duration, props["duration"], "File duration")

	// check content unit association
	suite.True(f.ContentUnitID.Valid, "File ContentUnitID.Valid")
	suite.Equal(cu.ID, f.ContentUnitID.Int64, "File ContentUnitID.Int64")

	// test with parent file
	pfi := File{
		FileName:  "dummy parent file",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef9876543210",
		Size:      math.MaxInt64,
	}
	parent, err := CreateFile(suite.tx, nil, pfi, nil)
	suite.Require().Nil(err)
	input.ParentSha1 = pfi.Sha1

	input.Sha1 = "012356789abcdef012356789abcdef1111111112"
	op, evnts, err = handleInsert(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().NotEmpty(evnts)

	// check file's ParentID has changed
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 2, "Number of files")
	f = op.R.Files[0]
	if op.R.Files[1].Language.String != "" {
		f = op.R.Files[1]
	}
	suite.True(f.ParentID.Valid, "File ParentID.Valid")
	suite.Equal(parent.ID, f.ParentID.Int64, "File ParentID.Int64")

	// test when file already exists
	cu2, err := CreateContentUnit(suite.tx, CT_LESSON_PART, nil)
	suite.Require().Nil(err)
	input.ContentUnitUID = cu2.UID

	op, evnts, err = handleInsert(suite.tx, input)
	suite.Require().NotNil(err)
}

func (suite *HandlersSuite) TestHandleInsertKiteiMakor() {
	// Create dummy content unit
	cu, err := CreateContentUnit(suite.tx, CT_LESSON_PART, nil)
	suite.Require().Nil(err)

	// Do insert operation
	input := InsertRequest{
		Operation: Operation{
			Station:    "Some station",
			User:       "operator@dev.com",
			WorkflowID: "workflow_id",
		},
		InsertType:     "kitei-makor",
		ContentUnitUID: cu.UID,
		AVFile: AVFile{
			File: File{
				FileName:  "kitei-makor.docx",
				Sha1:      "012356789abcdef012356789abcdef1111111111",
				Size:      98737,
				CreatedAt: &Timestamp{Time: time.Now()},
				MimeType:  MEDIA_TYPE_REGISTRY.ByExtension["docx"].MimeType,
				Language:  LANG_HEBREW,
			},
		},
		Mode: "new",
	}

	op, evnts, err := handleInsert(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().NotEmpty(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_INSERT].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	var props map[string]interface{}
	err = op.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.Operation.WorkflowID, props["workflow_id"].(string), "Operation workflow_id")
	suite.Equal(input.InsertType, props["insert_type"].(string), "Operation insert_type")
	suite.Equal(input.Mode, props["mode"].(string), "Operation mode")

	// Check user
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op))
	suite.Equal(input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 1, "Number of files")

	// check inserted file
	f := op.R.Files[0]
	suite.Equal(input.FileName, f.Name, "File Name")
	suite.Equal(input.Sha1, hex.EncodeToString(f.Sha1.Bytes), "File SHA1")
	suite.Equal(input.Size, f.Size, "File Size")
	suite.Equal(input.CreatedAt.Time.Unix(), f.FileCreatedAt.Time.Unix(), "File FileCreatedAt")
	suite.Equal("text", f.Type, "File Type")
	suite.Equal(input.MimeType, f.MimeType.String, "File MimeType")
	suite.Equal(input.Language, f.Language.String, "File Language")
	suite.False(f.ParentID.Valid, "File ParentID")
	err = f.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.InsertType, props["insert_type"].(string), "File insert_type")
	suite.Nil(props["duration"], "File duration")

	// check content unit association
	suite.Require().Nil(f.L.LoadContentUnit(suite.tx, true, f))
	ktCU := f.R.ContentUnit
	suite.Equal(CT_KITEI_MAKOR, CONTENT_TYPE_REGISTRY.ByID[ktCU.TypeID].Name, "KT CU type")
	suite.Require().Nil(ktCU.L.LoadDerivedContentUnitDerivations(suite.tx, true, ktCU))
	suite.Equal(cu.ID, ktCU.R.DerivedContentUnitDerivations[0].SourceID, "KT CU source CU")

	// test when KT cu exists
	input.FileName = "kitei-makor2.docx"
	input.Sha1 = "012356789abcdef012356789abcdef1111111112"

	op2, evnts, err := handleInsert(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().NotEmpty(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_INSERT].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	err = op.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.Operation.WorkflowID, props["workflow_id"].(string), "Operation workflow_id")
	suite.Equal(input.InsertType, props["insert_type"].(string), "Operation insert_type")
	suite.Equal(input.Mode, props["mode"].(string), "Operation mode")

	// Check user
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op))
	suite.Equal(input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 1, "Number of files")

	// check inserted file
	f2 := op2.R.Files[0]
	suite.Equal(input.FileName, f2.Name, "File Name")
	suite.Equal(input.Sha1, hex.EncodeToString(f2.Sha1.Bytes), "File SHA1")
	suite.Equal(input.Size, f2.Size, "File Size")
	suite.Equal(input.CreatedAt.Time.Unix(), f2.FileCreatedAt.Time.Unix(), "File FileCreatedAt")
	suite.Equal("text", f2.Type, "File Type")
	suite.Equal(input.MimeType, f2.MimeType.String, "File MimeType")
	suite.Equal(input.Language, f2.Language.String, "File Language")
	suite.False(f2.ParentID.Valid, "File ParentID")
	err = f.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.InsertType, props["insert_type"].(string), "File insert_type")
	suite.Nil(props["duration"], "File duration")

	// check content unit association
	suite.True(f2.ContentUnitID.Valid, "f2.ContentUnitID.Valid")
	suite.Equal(f.ContentUnitID.Int64, f2.ContentUnitID.Int64, "f2.ContentUnitID.Int64")
}

func (suite *HandlersSuite) TestHandleInsertPublication() {
	// Create dummy content unit
	cu, err := CreateContentUnit(suite.tx, CT_ARTICLE, nil)
	suite.Require().Nil(err)

	// create dummy publisher
	publisher := models.Publisher{
		UID: "12345678",
	}
	suite.Require().Nil(publisher.Insert(suite.tx))

	// Do insert operation
	input := InsertRequest{
		Operation: Operation{
			Station:    "Some station",
			User:       "operator@dev.com",
			WorkflowID: "workflow_id",
		},
		InsertType:     "publication",
		PublisherUID:   publisher.UID,
		ContentUnitUID: cu.UID,
		AVFile: AVFile{
			File: File{
				FileName:  "publication.png",
				Sha1:      "012356789abcdef012356789abcdef1111111111",
				Size:      98737,
				CreatedAt: &Timestamp{Time: time.Now()},
				MimeType:  MEDIA_TYPE_REGISTRY.ByExtension["png"].MimeType,
				Language:  LANG_HEBREW,
			},
		},
		Mode: "new",
	}

	op, evnts, err := handleInsert(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().NotEmpty(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_INSERT].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	var props map[string]interface{}
	err = op.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.Operation.WorkflowID, props["workflow_id"].(string), "Operation workflow_id")
	suite.Equal(input.InsertType, props["insert_type"].(string), "Operation insert_type")
	suite.Equal(input.Mode, props["mode"].(string), "Operation mode")

	// Check user
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op))
	suite.Equal(input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 1, "Number of files")

	// check inserted file
	f := op.R.Files[0]
	suite.Equal(input.FileName, f.Name, "File Name")
	suite.Equal(input.Sha1, hex.EncodeToString(f.Sha1.Bytes), "File SHA1")
	suite.Equal(input.Size, f.Size, "File Size")
	suite.Equal(input.CreatedAt.Time.Unix(), f.FileCreatedAt.Time.Unix(), "File FileCreatedAt")
	suite.Equal("image", f.Type, "File Type")
	suite.Equal(input.MimeType, f.MimeType.String, "File MimeType")
	suite.Equal(input.Language, f.Language.String, "File Language")
	suite.False(f.ParentID.Valid, "File ParentID")
	err = f.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.InsertType, props["insert_type"].(string), "File insert_type")
	suite.Nil(props["duration"], "File duration")

	// check content unit association
	suite.Require().Nil(f.L.LoadContentUnit(suite.tx, true, f))
	pCU := f.R.ContentUnit
	suite.Equal(CT_PUBLICATION, CONTENT_TYPE_REGISTRY.ByID[pCU.TypeID].Name, "Publication CU type")
	suite.Require().Nil(pCU.L.LoadDerivedContentUnitDerivations(suite.tx, true, pCU))
	suite.Equal(cu.ID, pCU.R.DerivedContentUnitDerivations[0].SourceID, "Publication CU source CU")

	// test when Publication cu exists
	input.FileName = "publication2.png"
	input.Sha1 = "012356789abcdef012356789abcdef1111111112"

	op2, evnts, err := handleInsert(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().NotEmpty(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_INSERT].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	err = op.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.Operation.WorkflowID, props["workflow_id"].(string), "Operation workflow_id")
	suite.Equal(input.InsertType, props["insert_type"].(string), "Operation insert_type")
	suite.Equal(input.Mode, props["mode"].(string), "Operation mode")

	// Check user
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op))
	suite.Equal(input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 1, "Number of files")

	// check inserted file
	f2 := op2.R.Files[0]
	suite.Equal(input.FileName, f2.Name, "File Name")
	suite.Equal(input.Sha1, hex.EncodeToString(f2.Sha1.Bytes), "File SHA1")
	suite.Equal(input.Size, f2.Size, "File Size")
	suite.Equal(input.CreatedAt.Time.Unix(), f2.FileCreatedAt.Time.Unix(), "File FileCreatedAt")
	suite.Equal("image", f2.Type, "File Type")
	suite.Equal(input.MimeType, f2.MimeType.String, "File MimeType")
	suite.Equal(input.Language, f2.Language.String, "File Language")
	suite.False(f2.ParentID.Valid, "File ParentID")
	err = f.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.InsertType, props["insert_type"].(string), "File insert_type")
	suite.Nil(props["duration"], "File duration")

	// check content unit association
	suite.True(f2.ContentUnitID.Valid, "f2.ContentUnitID.Valid")
	suite.Equal(f.ContentUnitID.Int64, f2.ContentUnitID.Int64, "f2.ContentUnitID.Int64")
}

func (suite *HandlersSuite) TestHandleInsertUpdateMode() {
	// Create dummy content unit
	cu, err := CreateContentUnit(suite.tx, CT_LESSON_PART, nil)
	suite.Require().Nil(err)

	// Create dummy old file
	ofi := File{
		FileName:  "dummy old file",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef9876543210",
		Size:      math.MaxInt64,
	}
	oldFile, err := CreateFile(suite.tx, nil, ofi, nil)
	suite.Require().Nil(err)
	err = cu.AddFiles(suite.tx, false, oldFile)
	suite.Require().Nil(err)

	// Do insert operation
	input := InsertRequest{
		Operation: Operation{
			Station:    "Some station",
			User:       "operator@dev.com",
			WorkflowID: "workflow_id",
		},
		InsertType:     "akladot",
		ContentUnitUID: cu.UID,
		AVFile: AVFile{
			File: File{
				FileName:  "akladot.doc",
				Sha1:      "012356789abcdef012356789abcdef1111111111",
				Size:      98737,
				CreatedAt: &Timestamp{Time: time.Now()},
				MimeType:  "application/msword",
				Language:  LANG_HEBREW,
			},
			Duration: 123.4,
		},
		Mode:    "update",
		OldSha1: ofi.Sha1,
	}

	op, evnts, err := handleInsert(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().NotEmpty(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_INSERT].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	var props map[string]interface{}
	err = op.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.Operation.WorkflowID, props["workflow_id"].(string), "Operation workflow_id")
	suite.Equal(input.InsertType, props["insert_type"].(string), "Operation insert_type")
	suite.Equal(input.Mode, props["mode"].(string), "Operation mode")

	// Check user
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op))
	suite.Equal(input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 2, "Number of files")
	f, of := op.R.Files[0], op.R.Files[1]
	if f.Name != "akladot.doc" {
		f, of = of, f
	}

	// check inserted file
	suite.Equal(input.FileName, f.Name, "File Name")
	suite.Equal(input.Sha1, hex.EncodeToString(f.Sha1.Bytes), "File SHA1")
	suite.Equal(input.Size, f.Size, "File Size")
	suite.Equal(input.CreatedAt.Time.Unix(), f.FileCreatedAt.Time.Unix(), "File FileCreatedAt")
	suite.Equal("text", f.Type, "File Type")
	suite.Equal(input.MimeType, f.MimeType.String, "File MimeType")
	suite.Equal(input.Language, f.Language.String, "File Language")
	suite.False(f.ParentID.Valid, "File ParentID")
	err = f.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.InsertType, props["insert_type"].(string), "File insert_type")
	suite.Equal(input.AVFile.Duration, props["duration"], "File duration")

	// check content unit association
	suite.True(f.ContentUnitID.Valid, "File ContentUnitID.Valid")
	suite.Equal(cu.ID, f.ContentUnitID.Int64, "File ContentUnitID.Int64")

	// check old file removed
	suite.True(of.RemovedAt.Valid, "Old File RemovedAt.Valid")
}

func (suite *HandlersSuite) TestHandleInsertRenameMode() {
	// Create dummy content unit
	cu, err := CreateContentUnit(suite.tx, CT_LESSON_PART, nil)
	suite.Require().Nil(err)

	cu2, err := CreateContentUnit(suite.tx, CT_LESSON_PART, nil)
	suite.Require().Nil(err)

	// Create dummy old file
	ofi := File{
		FileName:  "dummy file",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef9876543210",
		Size:      math.MaxInt64,
		Type:      "video",
		SubType:   "video subtype",
		MimeType:  "video/mp4",
		Language:  LANG_RUSSIAN,
	}
	oldFile, err := CreateFile(suite.tx, nil, ofi, nil)
	suite.Require().Nil(err)
	err = cu.AddFiles(suite.tx, false, oldFile)
	suite.Require().Nil(err)

	// Do insert operation
	input := InsertRequest{
		Operation: Operation{
			Station:    "Some station",
			User:       "operator@dev.com",
			WorkflowID: "workflow_id",
		},
		InsertType:     "akladot",
		ContentUnitUID: cu2.UID,
		AVFile: AVFile{
			File: File{
				FileName:  "akladot.doc",
				Sha1:      "012356789abcdef012356789abcdef9876543210",
				Size:      98737,
				CreatedAt: &Timestamp{Time: time.Now()},
				MimeType:  "application/msword",
				Language:  LANG_HEBREW,
			},
			Duration: 123.4,
		},
		Mode:    "rename",
	}

	op, evnts, err := handleInsert(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().NotEmpty(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_INSERT].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	var props map[string]interface{}
	err = op.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.Operation.WorkflowID, props["workflow_id"].(string), "Operation workflow_id")
	suite.Equal(input.InsertType, props["insert_type"].(string), "Operation insert_type")
	suite.Equal(input.Mode, props["mode"].(string), "Operation mode")

	// Check user
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op))
	suite.Equal(input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 1, "Number of files")
	f := op.R.Files[0]

	// check inserted file
	suite.Equal(input.FileName, f.Name, "File Name")
	suite.Equal(ofi.Sha1, hex.EncodeToString(f.Sha1.Bytes), "File SHA1")
	suite.Equal(ofi.Size, f.Size, "File Size")
	suite.Equal(input.CreatedAt.Time.Unix(), f.FileCreatedAt.Time.Unix(), "File FileCreatedAt")
	suite.Equal("text", f.Type, "File Type")
	suite.Equal(input.MimeType, f.MimeType.String, "File MimeType")
	suite.Equal(input.Language, f.Language.String, "File Language")
	suite.False(f.ParentID.Valid, "File ParentID")
	err = f.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.InsertType, props["insert_type"].(string), "File insert_type")

	// check content unit association
	suite.True(f.ContentUnitID.Valid, "File ContentUnitID.Valid")
	suite.Equal(cu2.ID, f.ContentUnitID.Int64, "File ContentUnitID.Int64")
}

func (suite *HandlersSuite) TestHandleTranscodeSuccess() {
	// Create dummy original file
	ofi := File{
		FileName:  "dummy_original_file.wmv",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef9876543210",
		Size:      math.MaxInt64,
		Type:      "video",
		MimeType:  MEDIA_TYPE_REGISTRY.ByExtension["wmv"].MimeType,
		Language:  LANG_HEBREW,
	}
	oProps := map[string]interface{}{
		"duration": 1234,
	}
	original, err := CreateFile(suite.tx, nil, ofi, oProps)
	suite.Require().Nil(err)

	// Create dummy content unit
	cu, err := CreateContentUnit(suite.tx, CT_LESSON_PART, nil)
	suite.Require().Nil(err)

	// associate original and content unit
	original.ContentUnitID = null.Int64From(cu.ID)
	err = original.Update(suite.tx, "content_unit_id")
	suite.Require().Nil(err)

	// Do transcode operation
	input := TranscodeRequest{
		Operation: Operation{
			Station: "Some station",
			User:    "operator@dev.com",
		},
		OriginalSha1: ofi.Sha1,
		MaybeFile: MaybeFile{
			FileName:  "dummy_original_file.mp4",
			Sha1:      "012356789abcdef012356789abcdef1111111111",
			Size:      98737,
			CreatedAt: &Timestamp{Time: time.Now()},
		},
	}

	op, evnts, err := handleTranscode(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_TRANSCODE].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	suite.False(op.Properties.Valid, "properties.Valid")

	// Check user
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op))
	suite.Equal(input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 2, "Number of files")
	fm := make(map[string]*models.File)
	for _, x := range op.R.Files {
		fm[x.Name] = x
	}

	originalParent := fm[ofi.FileName]
	suite.Equal(original.ID, originalParent.ID, "original <-> operation")

	// Check transcoded file
	mp4 := fm[input.FileName]
	suite.Equal(input.FileName, mp4.Name, "mp4: Name")
	suite.Equal(input.Sha1, hex.EncodeToString(mp4.Sha1.Bytes), "mp4: SHA1")
	suite.Equal(input.Size, mp4.Size, "mp4: Size")
	suite.Equal(input.CreatedAt.Time.Unix(), mp4.FileCreatedAt.Time.Unix(), "mp4: FileCreatedAt")
	suite.Equal(MEDIA_TYPE_REGISTRY.ByExtension["mp4"].Type, mp4.Type, "mp4: Type")
	suite.True(mp4.MimeType.Valid, "mp4: MimeType.Valid")
	suite.Equal(MEDIA_TYPE_REGISTRY.ByExtension["mp4"].MimeType, mp4.MimeType.String, "mp4: Type")
	suite.True(mp4.ParentID.Valid, "mp4: ParentID.Valid")
	suite.Equal(original.ID, mp4.ParentID.Int64, "mp4: ParentID")
	suite.True(mp4.Properties.Valid, "mp4: Properties.Valid")
	var props map[string]interface{}
	err = json.Unmarshal(mp4.Properties.JSON, &props)
	suite.Require().Nil(err, "json.Unmarshal mp4.Properties")
	suite.EqualValues(oProps["duration"], props["duration"], "mp4.Properties duration")

	// check content unit association
	suite.True(mp4.ContentUnitID.Valid, "mp4: ContentUnitID.Valid")
	suite.Equal(cu.ID, mp4.ContentUnitID.Int64, "mp4: ContentUnitID.Int64")

	// re-transcode
	input.MaybeFile = MaybeFile{
		FileName:  "re-transcoded.mp4",
		Sha1:      "012356789abcdef012356789abcdef2222222222",
		Size:      98738,
		CreatedAt: &Timestamp{Time: time.Now()},
	}

	op, evnts, err = handleTranscode(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 3, "Number of files")
	fm = make(map[string]*models.File)
	for _, x := range op.R.Files {
		fm[x.Name] = x
	}
	originalParent = fm[ofi.FileName]
	suite.Equal(original.ID, originalParent.ID, "original <-> operation")
	suite.Equal(mp4.ID, fm[mp4.Name].ID, "old transcoded <-> operation")

	// Check re-transcoded file
	mp4New := fm[input.FileName]
	suite.Equal(input.FileName, mp4New.Name, "mp4New: Name")
	suite.Equal(input.Sha1, hex.EncodeToString(mp4New.Sha1.Bytes), "mp4New: SHA1")
	suite.Equal(input.Size, mp4New.Size, "mp4: Size")
	suite.Equal(input.CreatedAt.Time.Unix(), mp4New.FileCreatedAt.Time.Unix(), "mp4New: FileCreatedAt")
	suite.Equal(MEDIA_TYPE_REGISTRY.ByExtension["mp4"].Type, mp4New.Type, "mp4: Type")
	suite.True(mp4New.MimeType.Valid, "mp4New: MimeType.Valid")
	suite.Equal(MEDIA_TYPE_REGISTRY.ByExtension["mp4"].MimeType, mp4New.MimeType.String, "mp4New: Type")
	suite.True(mp4New.ParentID.Valid, "mp4New: ParentID.Valid")
	suite.Equal(original.ID, mp4New.ParentID.Int64, "mp4New: ParentID")
	suite.True(mp4New.Properties.Valid, "mp4New: Properties.Valid")
	err = json.Unmarshal(mp4New.Properties.JSON, &props)
	suite.Require().Nil(err, "json.Unmarshal mp4New.Properties")
	suite.EqualValues(oProps["duration"], props["duration"], "mp4New.Properties duration")

	// Check old transcoded file
	err = mp4.Reload(suite.tx)
	suite.Require().Nil(err)
	suite.True(mp4.RemovedAt.Valid, "old mp4 removed")
}

func (suite *HandlersSuite) TestHandleTranscodeError() {
	// Create dummy original file
	ofi := File{
		FileName:  "dummy_original_file.wmv",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef9876543210",
		Size:      math.MaxInt64,
	}
	original, err := CreateFile(suite.tx, nil, ofi, nil)
	suite.Require().Nil(err)

	// Do transcode operation
	input := TranscodeRequest{
		Operation: Operation{
			Station: "Some station",
			User:    "operator@dev.com",
		},
		OriginalSha1: ofi.Sha1,
		Message:      "Some error description goes here",
	}

	op, evnts, err := handleTranscode(suite.tx, input)
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	// Check op
	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_TRANSCODE].ID, op.TypeID, "Operation TypeID")
	suite.Equal(input.Operation.Station, op.Station.String, "Operation Station")
	suite.True(op.Properties.Valid, "properties.Valid")
	var props map[string]interface{}
	err = json.Unmarshal(op.Properties.JSON, &props)
	suite.Require().Nil(err, "json.Unmarshal properties")
	suite.Equal(input.Message, props["message"], "op properties message")

	// Check user
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op))
	suite.Equal(input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	suite.Require().Nil(op.L.LoadFiles(suite.tx, true, op))
	suite.Len(op.R.Files, 1, "Number of files")
	originalParent := op.R.Files[0]
	suite.Equal(original.ID, originalParent.ID, "original <-> operation")
}
