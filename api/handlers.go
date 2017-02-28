package api

import (
	_ "github.com/lib/pq"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"gopkg.in/nullbio/null.v6"
	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/vattle/sqlboiler/boil"
	"database/sql"
	"github.com/Sirupsen/logrus"
	"encoding/json"
)


// Starts capturing file, i.e., morning lesson or other program.
func CaptureStartHandler(c *gin.Context) {
	var r CaptureStartRequest
	if c.BindJSON(&r) != nil {
		return
	}

	// TODO: should be taken from Input soon (needs impl in WF)
	// For now we only get lessons...
	contentType := CT_LESSON_PART

	// Start DB transaction
	tx, err := boil.Begin()

	// Create operation
	props, err := json.Marshal(map[string]string{
		"workflow_id": r.Operation.WorkflowID,
		"collection_uid": r.CollectionUID,
		"capture_source": r.CaptureSource,
	})
	operation, err := createOperation(tx, OP_CAPTURE_START, r.Operation, &props)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error":  err.Error()})
		return
	}

	// Create file
	file := models.File{
		UID:         utils.GenerateUID(8),
		Name:        r.FileName,
	}
	err = operation.AddFiles(tx, true, &file)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error":  err.Error()})
		return
	}

	// Create content unit
	cType, err := models.ContentTypesG(qm.Where("name=?", contentType)).One()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error":  err.Error()})
		return
	}

	name := models.StringTranslation{Text: DEFAULT_NAMES[contentType]}
	if err = name.Insert(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error":  err.Error()})
		return
	}

	cu := models.ContentUnit{
		TypeID: cType.ID,
		UID:    utils.GenerateUID(8),
		NameID:        name.ID,
	}
	err = file.SetContentUnit(tx, true, &cu)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error":  err.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Stops capturing file, i.e., morning lesson or other program.
//
// `capture_id`: Unique identifier per collection, i.e., for morning lesson all
// parts of the lesson will have the same `capture_id`.
func CaptureStopHandler(c *gin.Context) {
	var r CaptureStopRequest
	if c.BindJSON(&r) != nil {
		return
	}

	// DO logic here

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Demux
func DemuxHandler(c *gin.Context) {
	var r DemuxRequest
	if c.BindJSON(&r) != nil {
		return
	}

	// DO logic here

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Moves file from capture machine to other storage.
func SendHandler(c *gin.Context) {
	var r SendRequest
	if c.BindJSON(&r) != nil {
		return
	}

	// DO logic here

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Enabled file to be accessible from URL.
func UploadHandler(c *gin.Context) {
	var r UploadRequest
	if c.BindJSON(&r) != nil {
		return
	}

	// DO logic here

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}


// Helpers

func createOperation(exec boil.Executor, name string, o Operation, properties *[]byte) (*models.Operation, error) {
	opType, err := models.OperationTypesG(qm.Where("name=?", name)).One()
	if err != nil {
		return nil, err
	}

	operation := models.Operation{
		TypeID: opType.ID,
		UID: utils.GenerateUID(8),
		Station: null.StringFrom(o.Station),
		Properties: null.JSONFromPtr(properties),
	}

	// Lookup user, skip if doesn't exist
	user, err := models.UsersG(qm.Where("email=?", o.User)).One()
	if err == nil {
		operation.UserID = null.Int64From(user.ID)
	} else {
		if err == sql.ErrNoRows {
			logrus.Warnf("Unknown User [%s]. Skipping.", o.User)
		} else {
			return nil, err
		}
	}

	return &operation, operation.Insert(exec)
}
