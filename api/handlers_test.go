package api

import (
	"regexp"
	"testing"
	
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.Equal(t, OPERATION_TYPE_REGISTRY.ByName[OP_CAPTURE_START].ID, op.TypeID,
		"Operation TypeID")
	assert.Equal(t, input.Operation.Station, op.Station.String, "Operation Station")

	// Check op user
	require.Nil(t, op.L.LoadUser(tx, true, op))
	assert.Equal(t, input.Operation.User, op.R.User.Email, "Operation User")

	// Check op properties
	var props = make(map[string]interface{})
	err = op.Properties.Unmarshal(&props)
	require.Nil(t, err)
	assert.Equal(t, input.Operation.WorkflowID, props["workflow_id"], "properties: workflow_id")
	assert.Equal(t, input.CaptureSource, props["capture_source"], "properties: capture_source")
	assert.Equal(t, input.CollectionUID, props["collection_uid"], "properties: collection_uid")

	// Check associated files
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
	assert.Regexp(t, regexp.MustCompile("[a-zA-z0-9]{8}"), op.UID, "UID regex")
	assert.True(t, op.Station.Valid, "Station null valid")
	assert.Equal(t, o.Station, op.Station.String, "Station")
	user, err := models.Users(tx, qm.Where("email=?", o.User)).One()
	assert.Equal(t, user.ID, op.UserID.Int64, "User")

	// test with workflow_id
	o.WorkflowID = "workflow_id"
	op, err = createOperation(tx, OP_CAPTURE_START, o, nil)
	require.Nil(t, err)
	assert.True(t, op.Properties.Valid)
	var props map[string]interface{}
	err = op.Properties.Unmarshal(&props)
	require.Nil(t, err)
	assert.Equal(t, o.WorkflowID, props["workflow_id"], "props: workflow_id")


	// test with custom props
	customProps := map[string]interface{}{
		"a": 1,
		"b": "2",
		"c": true,
		"d": []float64{1.2, 2.3, 3.4},
	}
	op, err = createOperation(tx, OP_CAPTURE_START, o, customProps)
	require.Nil(t, err)
	props = make(map[string]interface{})
	err = op.Properties.Unmarshal(&props)
	require.Nil(t, err)
	assert.Equal(t, o.WorkflowID, props["workflow_id"], "props: workflow_id")
	assert.EqualValues(t, customProps["a"], props["a"], "props: a")
	assert.Equal(t, customProps["b"], props["b"], "props: b")
	assert.Equal(t, customProps["c"], props["c"], "props: c")
	assert.Len(t, props["d"].([]interface{}), len(customProps["d"].([]float64)), "props: d length")
	for i, v := range customProps["d"].([]float64) {
		assert.Equal(t, v, props["d"].([]interface{})[i], "props: d[%d]", i)
	}

}

// helpers

func eagerReloadOperation(exec boil.Executor) (*models.Operation, error) {
	return models.Operations(exec,
		qm.OrderBy("created_at desc"),
		qm.Load("User", "Files")).
		One()
}
