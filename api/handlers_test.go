package api

import (
	"encoding/hex"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/vattle/sqlboiler/boil"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

type HandlersSuite struct {
	suite.Suite
	utils.TestDBManager
	tx boil.Transactor
}

func (suite *HandlersSuite) SetupSuite() {
	suite.Require().Nil(suite.InitTestDB())
	suite.Require().Nil(InitTypeRegistries(boil.GetDB()))
}

func (suite *HandlersSuite) TearDownSuite() {
	suite.Require().Nil(suite.DestroyTestDB())
}

func (suite *HandlersSuite) SetupTest() {
	var err error
	suite.tx, err = boil.Begin()
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

	op, err := handleCaptureStart(suite.tx, input)
	suite.Require().Nil(err)

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
	opStart, err := handleCaptureStart(suite.tx, CaptureStartRequest{
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

	op, err := handleCaptureStop(suite.tx, input)
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

	op, err := handleDemux(suite.tx, input)
	suite.Require().Nil(err)

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

	op, err := handleTrim(suite.tx, input)
	suite.Require().Nil(err)

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

	op, err := handleSend(suite.tx, input)
	suite.Require().Nil(err)

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
		},
	}

	op, err := handleConvert(suite.tx, input)
	suite.Require().Nil(err)

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
