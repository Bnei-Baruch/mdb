package api

import (
	"encoding/hex"
	"math"
	"regexp"
	"testing"
	"time"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"
)

func TestHandleCaptureStart(t *testing.T) {
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

	tx, err := boil.Begin()
	require.Nil(t, err)
	defer tx.Rollback()

	op, err := handleCaptureStart(tx, input)
	require.Nil(t, err)

	// Check op
	assert.Equal(t, OPERATION_TYPE_REGISTRY.ByName[OP_CAPTURE_START].ID, op.TypeID, "Operation TypeID")
	assert.Equal(t, input.Operation.Station, op.Station.String, "Operation Station")
	var props map[string]interface{}
	err = op.Properties.Unmarshal(&props)
	require.Nil(t, err)
	assert.Equal(t, input.Operation.WorkflowID, props["workflow_id"], "properties: workflow_id")
	assert.Equal(t, input.CaptureSource, props["capture_source"], "properties: capture_source")
	assert.Equal(t, input.CollectionUID, props["collection_uid"], "properties: collection_uid")

	// Check user
	require.Nil(t, op.L.LoadUser(tx, true, op))
	assert.Equal(t, input.Operation.User, op.R.User.Email, "Operation User")

	// Check associated files
	require.Nil(t, op.L.LoadFiles(tx, true, op))
	assert.Len(t, op.R.Files, 1, "Number of files")
	f := op.R.Files[0]
	assert.Equal(t, input.FileName, f.Name, "File: Name")
	assert.False(t, f.Sha1.Valid, "File: SHA1")
}

func TestCreateOperation(t *testing.T) {
	tx, err := boil.Begin()
	require.Nil(t, err)
	defer tx.Rollback()

	// test minimal input
	o := Operation{
		Station: "station",
		User:    "operator@dev.com",
	}
	op, err := createOperation(tx, OP_CAPTURE_START, o, nil)
	require.Nil(t, err)
	require.Nil(t, op.Reload(tx))

	assert.Equal(t, OPERATION_TYPE_REGISTRY.ByName[OP_CAPTURE_START].ID, op.TypeID, "TypeID")
	assert.Regexp(t, regexp.MustCompile("[a-zA-z0-9]{8}"), op.UID, "UID")
	assert.True(t, op.Station.Valid, "Station.Valid")
	assert.Equal(t, o.Station, op.Station.String, "Station.String")
	user, err := models.Users(tx, qm.Where("email=?", o.User)).One()
	assert.Equal(t, user.ID, op.UserID.Int64, "User")

	// test with unknown user
	o.User = "unknown@example.com"
	op, err = createOperation(tx, OP_CAPTURE_START, o, nil)
	require.Nil(t, err)
	require.Nil(t, op.Reload(tx))
	assert.False(t, op.UserID.Valid)
	o.User = "operator@dev.com"

	// test with workflow_id
	o.WorkflowID = "workflow_id"
	op, err = createOperation(tx, OP_CAPTURE_START, o, nil)
	require.Nil(t, err)
	require.Nil(t, op.Reload(tx))

	assert.True(t, op.Properties.Valid)
	var props map[string]interface{}
	err = op.Properties.Unmarshal(&props)
	require.Nil(t, err)
	assert.Equal(t, o.WorkflowID, props["workflow_id"], "props: workflow_id")

	// test with custom props
	customProps := getCustomProps()
	op, err = createOperation(tx, OP_CAPTURE_START, o, customProps)
	require.Nil(t, err)
	require.Nil(t, op.Reload(tx))

	err = op.Properties.Unmarshal(&props)
	require.Nil(t, err)
	assertCustomProps(t, customProps, props)
}

func TestCreateFile(t *testing.T) {
	tx, err := boil.Begin()
	require.Nil(t, err)
	defer tx.Rollback()

	// test minimal input
	f := File{
		FileName:  "file_name",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef0123456789",
		Size:      math.MaxInt64,
	}
	file, err := createFile(tx, nil, f, nil)
	require.Nil(t, err)
	require.Nil(t, file.Reload(tx))

	assert.Regexp(t, regexp.MustCompile("[a-zA-z0-9]{8}"), file.UID, "UID")
	assert.Equal(t, f.FileName, file.Name, "Name")
	assert.True(t, file.FileCreatedAt.Valid, "FileCreatedAt.Valid")
	assert.Equal(t, f.CreatedAt.Unix(), file.FileCreatedAt.Time.Unix(), "FileCreatedAt.Time")
	assert.True(t, file.Sha1.Valid, "Sha1.Valid")
	assert.Equal(t, f.Sha1, hex.EncodeToString(file.Sha1.Bytes), "Sha1.Bytes")
	assert.Equal(t, f.Size, file.Size, "Size")
	assert.Empty(t, file.Type, "Type")
	assert.Empty(t, file.SubType, "SubType")
	assert.False(t, file.MimeType.Valid, "MimeType.Valid")
	assert.False(t, file.Language.Valid, "Language.Valid")
	assert.False(t, file.ParentID.Valid, "ParentID.Valid")

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
	file2, err := createFile(tx, file, f2, nil)
	require.Nil(t, err)
	require.Nil(t, file2.Reload(tx))

	assert.Equal(t, f2.Sha1, hex.EncodeToString(file2.Sha1.Bytes), "file2.Sha1.Bytes")
	assert.Equal(t, f2.Size, file2.Size, "file2.Size")
	assert.Equal(t, f2.Type, file2.Type, "file2.Type")
	assert.Equal(t, f2.SubType, file2.SubType, "file2.SubType")
	assert.True(t, file2.MimeType.Valid, "file2.MimeType.Valid")
	assert.Equal(t, f2.MimeType, file2.MimeType.String, "file2.MimeType.String")
	assert.True(t, file2.Language.Valid, "file2.Language.Valid")
	assert.Equal(t, f2.Language, file2.Language.String, "file2.Language.String")
	assert.True(t, file2.ParentID.Valid, "file2.ParentID.Valid")
	assert.Equal(t, file.ID, file2.ParentID.Int64, "file2.ParentID.Int64")

	// test with custom props
	f3 := File{
		FileName:  "file_name",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "abcdef012356789abcdef0123456789abcdef012",
		Size:      1,
	}
	customProps := getCustomProps()
	file3, err := createFile(tx, nil, f3, customProps)
	require.Nil(t, err)
	require.Nil(t, file3.Reload(tx))

	var props map[string]interface{}
	err = file3.Properties.Unmarshal(&props)
	require.Nil(t, err)
	assertCustomProps(t, customProps, props)
}

func TestFindFileBySHA1(t *testing.T) {
	tx, err := boil.Begin()
	require.Nil(t, err)
	defer tx.Rollback()

	const (
		sha1         = "012356789abcdef012356789abcdef0123456789"
		unknown_sha1 = "012356789abcdef012356789abcdef9876543210"
	)

	// test with empty table
	file, sha1b, err := findFileBySHA1(tx, sha1)
	assert.Exactly(t, FileNotFound{Sha1: sha1}, err, "empty error")
	assert.Equal(t, sha1, hex.EncodeToString(sha1b), "empty.sha1 bytes")
	assert.Nil(t, file, "empty.file")

	// test match with non empty table
	fm := File{
		FileName:  "file_name",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      sha1,
		Size:      math.MaxInt64,
	}
	newFile, err := createFile(tx, nil, fm, nil)
	require.Nil(t, err)
	file, _, err = findFileBySHA1(tx, sha1)
	require.Nil(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, newFile.ID, file.ID)

	// test miss with non empty table
	file, sha1b, err = findFileBySHA1(tx, unknown_sha1)
	assert.Exactly(t, FileNotFound{Sha1: unknown_sha1}, err, "miss error")
	assert.Equal(t, unknown_sha1, hex.EncodeToString(sha1b), "miss.sha1 bytes")
	assert.Nil(t, file)
}

// Helpers

func getCustomProps() map[string]interface{} {
	return map[string]interface{}{
		"a": 1,
		"b": "2",
		"c": true,
		"d": []float64{1.2, 2.3, 3.4},
	}
}

func assertCustomProps(t *testing.T, expected map[string]interface{}, actual map[string]interface{}) {
	assert.EqualValues(t, expected["a"], actual["a"], "props: a")
	assert.Equal(t, expected["b"], actual["b"], "props: b")
	assert.Equal(t, expected["c"], actual["c"], "props: c")
	assert.Len(t, actual["d"].([]interface{}), len(expected["d"].([]float64)), "props: d length")
	for i, v := range expected["d"].([]float64) {
		assert.Equal(t, v, actual["d"].([]interface{})[i], "props: d[%d]", i)
	}
}
