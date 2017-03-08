package api

import (
	"encoding/hex"
	"math"
	"regexp"
	"testing"
	"time"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/stretchr/testify/suite"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"
)

type HandlersSuite struct {
	suite.Suite
	utils.TestDBManager
	tx boil.Transactor
}

func (suite *HandlersSuite) SetupSuite() {
	suite.Require().Nil(suite.InitTestDB())
	suite.Require().Nil(OPERATION_TYPE_REGISTRY.Init())
	suite.Require().Nil(CONTENT_TYPE_REGISTRY.Init())
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
		AVFile: AVFile{
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
			Duration: 892.1900,
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

	err = f.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(input.Duration, props["duration"], "File props: duration")
}

func (suite *HandlersSuite) TestCreateOperation() {
	// test minimal input
	o := Operation{
		Station: "station",
		User:    "operator@dev.com",
	}
	op, err := createOperation(suite.tx, OP_CAPTURE_START, o, nil)
	suite.Require().Nil(err)
	suite.Require().Nil(op.Reload(suite.tx))

	suite.Equal(OPERATION_TYPE_REGISTRY.ByName[OP_CAPTURE_START].ID, op.TypeID, "TypeID")
	suite.Regexp(regexp.MustCompile("[a-zA-z0-9]{8}"), op.UID, "UID")
	suite.True(op.Station.Valid, "Station.Valid")
	suite.Equal(o.Station, op.Station.String, "Station.String")
	user, err := models.Users(suite.tx, qm.Where("email=?", o.User)).One()
	suite.Equal(user.ID, op.UserID.Int64, "User")

	// test with unknown user
	o.User = "unknown@example.com"
	op, err = createOperation(suite.tx, OP_CAPTURE_START, o, nil)
	suite.Require().Nil(err)
	suite.Require().Nil(op.Reload(suite.tx))
	suite.False(op.UserID.Valid)
	o.User = "operator@dev.com"

	// test with workflow_id
	o.WorkflowID = "workflow_id"
	op, err = createOperation(suite.tx, OP_CAPTURE_START, o, nil)
	suite.Require().Nil(err)
	suite.Require().Nil(op.Reload(suite.tx))

	suite.True(op.Properties.Valid)
	var props map[string]interface{}
	err = op.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(o.WorkflowID, props["workflow_id"], "props: workflow_id")

	// test with custom props
	customProps := suite.getCustomProps()
	op, err = createOperation(suite.tx, OP_CAPTURE_START, o, customProps)
	suite.Require().Nil(err)
	suite.Require().Nil(op.Reload(suite.tx))

	err = op.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.assertCustomProps(customProps, props)
}

func (suite *HandlersSuite) TestCreateFile() {
	// test minimal input
	f := File{
		FileName:  "file_name",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef0123456789",
		Size:      math.MaxInt64,
	}
	file, err := createFile(suite.tx, nil, f, nil)
	suite.Require().Nil(err)
	suite.Require().Nil(file.Reload(suite.tx))

	suite.Regexp(regexp.MustCompile("[a-zA-z0-9]{8}"), file.UID, "UID")
	suite.Equal(f.FileName, file.Name, "Name")
	suite.True(file.FileCreatedAt.Valid, "FileCreatedAt.Valid")
	suite.Equal(f.CreatedAt.Unix(), file.FileCreatedAt.Time.Unix(), "FileCreatedAt.Time")
	suite.True(file.Sha1.Valid, "Sha1.Valid")
	suite.Equal(f.Sha1, hex.EncodeToString(file.Sha1.Bytes), "Sha1.Bytes")
	suite.Equal(f.Size, file.Size, "Size")
	suite.Empty(file.Type, "Type")
	suite.Empty(file.SubType, "SubType")
	suite.False(file.MimeType.Valid, "MimeType.Valid")
	suite.False(file.Language.Valid, "Language.Valid")
	suite.False(file.ParentID.Valid, "ParentID.Valid")

	// test with optional attributes
	f2 := File{
		FileName:  "file_name",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef9876543210",
		Size:      1,
		Type:      "type",
		SubType:   "subtype",
		MimeType:  "mimetype",
		Language:  LANG_RUSSIAN,
	}
	file2, err := createFile(suite.tx, file, f2, nil)
	suite.Require().Nil(err)
	suite.Require().Nil(file2.Reload(suite.tx))

	suite.Equal(f2.Sha1, hex.EncodeToString(file2.Sha1.Bytes), "file2.Sha1.Bytes")
	suite.Equal(f2.Size, file2.Size, "file2.Size")
	suite.Equal(f2.Type, file2.Type, "file2.Type")
	suite.Equal(f2.SubType, file2.SubType, "file2.SubType")
	suite.True(file2.MimeType.Valid, "file2.MimeType.Valid")
	suite.Equal(f2.MimeType, file2.MimeType.String, "file2.MimeType.String")
	suite.True(file2.Language.Valid, "file2.Language.Valid")
	suite.Equal(f2.Language, file2.Language.String, "file2.Language.String")
	suite.True(file2.ParentID.Valid, "file2.ParentID.Valid")
	suite.Equal(file.ID, file2.ParentID.Int64, "file2.ParentID.Int64")

	// test with custom props
	f3 := File{
		FileName:  "file_name",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "abcdef012356789abcdef0123456789abcdef012",
		Size:      1,
	}
	customProps := suite.getCustomProps()
	file3, err := createFile(suite.tx, nil, f3, customProps)
	suite.Require().Nil(err)
	suite.Require().Nil(file3.Reload(suite.tx))

	var props map[string]interface{}
	err = file3.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.assertCustomProps(customProps, props)
}

func (suite *HandlersSuite) TestFindFileBySHA1() {
	const (
		sha1         = "012356789abcdef012356789abcdef0123456789"
		unknown_sha1 = "012356789abcdef012356789abcdef9876543210"
	)

	// test with empty table
	file, sha1b, err := findFileBySHA1(suite.tx, sha1)
	suite.Exactly(FileNotFound{Sha1: sha1}, err, "empty error")
	suite.Equal(sha1, hex.EncodeToString(sha1b), "empty.sha1 bytes")
	suite.Nil(file, "empty.file")

	// test match with non empty table
	fm := File{
		FileName:  "file_name",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      sha1,
		Size:      math.MaxInt64,
	}
	newFile, err := createFile(suite.tx, nil, fm, nil)
	suite.Require().Nil(err)
	file, _, err = findFileBySHA1(suite.tx, sha1)
	suite.Require().Nil(err)
	suite.NotNil(file)
	suite.Equal(newFile.ID, file.ID)

	// test miss with non empty table
	file, sha1b, err = findFileBySHA1(suite.tx, unknown_sha1)
	suite.Exactly(FileNotFound{Sha1: unknown_sha1}, err, "miss error")
	suite.Equal(unknown_sha1, hex.EncodeToString(sha1b), "miss.sha1 bytes")
	suite.Nil(file)
}

// Helpers

func (suite *HandlersSuite) getCustomProps() map[string]interface{} {
	return map[string]interface{}{
		"a": 1,
		"b": "2",
		"c": true,
		"d": []float64{1.2, 2.3, 3.4},
	}
}

func (suite *HandlersSuite) assertCustomProps(expected map[string]interface{}, actual map[string]interface{}) {
	suite.EqualValues(expected["a"], actual["a"], "props: a")
	suite.Equal(expected["b"], actual["b"], "props: b")
	suite.Equal(expected["c"], actual["c"], "props: c")
	suite.Len(actual["d"].([]interface{}), len(expected["d"].([]float64)), "props: d length")
	for i, v := range expected["d"].([]float64) {
		suite.Equal(v, actual["d"].([]interface{})[i], "props: d[%d]", i)
	}
}
