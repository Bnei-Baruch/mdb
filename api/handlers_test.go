package api

import (
	"testing"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/stretchr/testify/assert"
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
	assert.Nil(t, err)
	defer tx.Rollback()

	err = handleCaptureStart(tx, input)
	assert.Nil(t, err)

	operation, err := models.Operations(tx,
		qm.OrderBy("created_at desc"),
		qm.Load("User", "Files")).
		One()
	assert.Nil(t, err, "Error looking up last inserted operation in DB")

	// Check operation
	assert.Equal(t, OPERATION_TYPE_REGISTRY.ByName[OP_CAPTURE_START].ID, operation.TypeID,
		"Operation TypeID")
	assert.Equal(t, input.Operation.Station, operation.Station.String,"Operation Station")
	assert.Equal(t, input.Operation.User, operation.R.User.Email,"Operation User")

	// Check operation properties
	var props  = make(map[string]interface{})
	err = operation.Properties.Unmarshal(&props)
	assert.Nil(t, err)
	assert.Equal(t, input.Operation.WorkflowID, props["workflow_id"],"properties: workflow_id")
	assert.Equal(t, input.CaptureSource, props["capture_source"],"properties: capture_source")
	assert.Equal(t, input.CollectionUID, props["collection_uid"],"properties: collection_uid")

	// Check associated files
	assert.Len(t, operation.R.Files, 1, "Number of files")
	f := operation.R.Files[0]
	assert.Equal(t, input.FileName, f.Name,"File: Name")
	assert.False(t, f.Sha1.Valid, "File: SHA1")
}
