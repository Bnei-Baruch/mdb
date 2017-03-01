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
	log "github.com/Sirupsen/logrus"
	"github.com/vattle/sqlboiler/queries"
	"encoding/hex"
)

// Starts capturing file, i.e., morning lesson or other program.
func CaptureStartHandler(c *gin.Context) {
	var i CaptureStartRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleCaptureStart)
	}
}

// Stops capturing file, i.e., morning lesson or other program.
func CaptureStopHandler(c *gin.Context) {
	var i CaptureStopRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleCaptureStop)
	}
}

// Demux original file to low resolution proxy
func DemuxHandler(c *gin.Context) {
	var i DemuxRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleDemux)
	}
}

// Moves file from capture machine to other storage.
func SendHandler(c *gin.Context) {
	var i SendRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleSend)
	}
}

// Enabled file to be accessible from URL.
func UploadHandler(c *gin.Context) {
	var i UploadRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleUpload)
	}
}

// Handler logic

func handleCaptureStart(exec boil.Executor, input interface{}) error {
	r := input.(CaptureStartRequest)

	// Create operation
	props, err := json.Marshal(map[string]string{
		"workflow_id":    r.Operation.WorkflowID,
		"collection_uid": r.CollectionUID,
		"capture_source": r.CaptureSource,
	})
	operation, err := createOperation(exec, OP_CAPTURE_START, r.Operation, &props)
	if err != nil {
		return err
	}

	// Create file
	file := models.File{
		UID:  utils.GenerateUID(8),
		Name: r.FileName,
	}
	err = operation.AddFiles(exec, true, &file)

	return err
}

func handleCaptureStop(exec boil.Executor, input interface{}) error {
	r := input.(CaptureStopRequest)

	// Create operation
	props, err := json.Marshal(map[string]string{
		"workflow_id":    r.Operation.WorkflowID,
		"collection_uid": r.CollectionUID,
		"capture_source": r.CaptureSource,
		"content_type":   r.ContentType,
		"part":           r.Part,
	})
	operation, err := createOperation(exec, OP_CAPTURE_STOP, r.Operation, &props)
	if err != nil {
		return err
	}

	// Create file
	sha1, err := hex.DecodeString(r.Sha1)
	if err != nil {
		return err
	}
	file := models.File{
		UID:  utils.GenerateUID(8),
		Name: r.FileName,
		Sha1: null.BytesFrom(sha1),
		Size: r.Size,
	}

	// Find parent file (capture_start with same workflow_id)
	var parentID int64
	err = queries.Raw(exec, "SELECT file_id FROM files_operations INNER JOIN operations ON operation_id = id WHERE properties -> 'workflow_id' ? $1",
		r.Operation.WorkflowID).QueryRow().Scan(&parentID)
	if err == nil {
		file.ParentID = null.Int64From(parentID)
	} else {
		if err == sql.ErrNoRows {
			logrus.Warnf("capture_start operation not found for workflow_id [%s]. Skipping.",
				r.Operation.WorkflowID)
		} else {
			return err
		}
	}

	err = operation.AddFiles(exec, true, &file)
	return err
}

func handleDemux(exec boil.Executor, input interface{}) error {
	return nil
}

func handleSend(exec boil.Executor, input interface{}) error {
	return nil
}

func handleUpload(exec boil.Executor, input interface{}) error {
	return nil
}

// Helpers

// Generic operation handler.
// 	* Manage DB transactions
// 	* Call operation logic handler
// 	* Render JSON response
func handleOperation(c *gin.Context, input interface{}, opHandler func(boil.Executor, interface{}) error) {
	tx, err := boil.Begin()
	if err == nil {
		err = opHandler(tx, input)
		if err == nil {
			tx.Commit()
		} else {
			log.Error("Error handling operation: ", err)
			if err = tx.Rollback(); err != nil {
				log.Error("Error rolling back DB transaction: ", err)
			}
		}
	}

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
	}
}

func createOperation(exec boil.Executor, name string, o Operation, properties *[]byte) (*models.Operation, error) {
	operation := models.Operation{
		TypeID:     OPERATION_TYPE_REGISTRY.ByName[name].ID,
		UID:        utils.GenerateUID(8),
		Station:    null.StringFrom(o.Station),
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
