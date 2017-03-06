package api

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/nullbio/null.v6"
)

// Start capture of AV file, i.e. morning lesson, tv program, etc...
func CaptureStartHandler(c *gin.Context) {
	log.Info(OP_CAPTURE_START)
	var i CaptureStartRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleCaptureStart)
	}
}

// Stop capture of AV file, ending a matching capture_start event.
// This is the first time a physical file is created in the studio.
func CaptureStopHandler(c *gin.Context) {
	log.Info(OP_CAPTURE_STOP)
	var i CaptureStopRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleCaptureStop)
	}
}

// Demux manifest file to original and low resolution proxy
func DemuxHandler(c *gin.Context) {
	log.Info(OP_DEMUX)
	var i DemuxRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleDemux)
	}
}

// Trim demuxed files at certain points
func TrimHandler(c *gin.Context) {
	log.Info(OP_TRIM)
	var i TrimRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleTrim)
	}
}

// Final files sent from studio
func SendHandler(c *gin.Context) {
	log.Info(OP_SEND)
	var i SendRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleSend)
	}
}

// File uploaded to a public accessible URL
func UploadHandler(c *gin.Context) {
	log.Info(OP_UPLOAD)
	var i UploadRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleUpload)
	}
}

// Handler logic

func handleCaptureStart(exec boil.Executor, input interface{}) (*models.Operation, error) {
	r := input.(CaptureStartRequest)

	log.Info("Creating operation")
	props := map[string]interface{}{
		"capture_source": r.CaptureSource,
		"collection_uid": r.CollectionUID,
	}
	operation, err := createOperation(exec, OP_CAPTURE_START, r.Operation, props)
	if err != nil {
		return nil, err
	}

	log.Info("Creating file and associating to operation")
	file := models.File{
		UID:  utils.GenerateUID(8),
		Name: r.FileName,
	}
	return operation, operation.AddFiles(exec, true, &file)
}

func handleCaptureStop(exec boil.Executor, input interface{}) (*models.Operation, error) {
	r := input.(CaptureStopRequest)

	log.Info("Creating operation")
	props := map[string]interface{}{
		"capture_source": r.CaptureSource,
		"collection_uid": r.CollectionUID, // $LID = backup capture id when lesson, capture_id when program (part=false)
		"part":           r.Part,
	}
	operation, err := createOperation(exec, OP_CAPTURE_STOP, r.Operation, props)
	if err != nil {
		return nil, err
	}

	log.Info("Looking up parent file, workflow_id=", r.Operation.WorkflowID)
	var parentID int64
	err = queries.Raw(exec,
		`SELECT file_id FROM files_operations
		 INNER JOIN operations ON operation_id = id
		 WHERE type_id=$1 AND properties -> 'workflow_id' ? $2`,
		OPERATION_TYPE_REGISTRY.ByName[OP_CAPTURE_START].ID,
		r.Operation.WorkflowID).
		QueryRow().
		Scan(&parentID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warnf("capture_start operation not found for workflow_id [%s]. Skipping.",
				r.Operation.WorkflowID)
		} else {
			return nil, err
		}
	}

	log.Info("Creating file")
	props = map[string]interface{}{
		"duration": r.Duration,
	}
	file, err := createFile(exec, &models.File{ID: parentID}, r.File, props)
	if err != nil {
		return nil, err
	}

	log.Info("Associating file to operation")
	return operation, operation.AddFiles(exec, false, file)
}

func handleDemux(exec boil.Executor, input interface{}) (*models.Operation, error) {
	r := input.(DemuxRequest)

	parent, _, err := findFileBySHA1(exec, r.Sha1)
	if err != nil {
		return nil, err
	}

	log.Info("Creating operation")
	props := map[string]interface{}{
		"capture_source": r.CaptureSource,
	}
	operation, err := createOperation(exec, OP_DEMUX, r.Operation, props)
	if err != nil {
		return nil, err
	}

	log.Info("Creating original")
	props = map[string]interface{}{
		"duration": r.Original.Duration,
	}
	original, err := createFile(exec, parent, r.Original.File, props)
	if err != nil {
		return nil, err
	}

	log.Info("Creating proxy")
	props = map[string]interface{}{
		"duration": r.Proxy.Duration,
	}
	proxy, err := createFile(exec, parent, r.Proxy.File, props)
	if err != nil {
		return nil, err
	}

	log.Info("Associating files to operation")
	return operation, operation.AddFiles(exec, false, parent, original, proxy)
}

func handleTrim(exec boil.Executor, input interface{}) (*models.Operation, error) {
	r := input.(TrimRequest)

	original, _, err := findFileBySHA1(exec, r.OriginalSha1)
	if err != nil {
		return nil, err
	}

	proxy, _, err := findFileBySHA1(exec, r.ProxySha1)
	if err != nil {
		return nil, err
	}

	log.Info("Creating operation")
	props := map[string]interface{}{
		"capture_source": r.CaptureSource,
		"in":             r.In,
		"out":            r.Out,
	}
	operation, err := createOperation(exec, OP_TRIM, r.Operation, props)
	if err != nil {
		return nil, err
	}

	log.Info("Creating trimmed original")
	props = map[string]interface{}{
		"duration": r.Original.Duration,
	}
	originalTrim, err := createFile(exec, original, r.Original.File, props)
	if err != nil {
		return nil, err
	}

	log.Info("Creating trimmed proxy")
	props = map[string]interface{}{
		"duration": r.Proxy.Duration,
	}
	proxyTrim, err := createFile(exec, proxy, r.Proxy.File, props)
	if err != nil {
		return nil, err
	}

	log.Info("Associating files to operation")
	return operation, operation.AddFiles(exec, false, original, originalTrim, proxy, proxyTrim)
}

func handleSend(exec boil.Executor, input interface{}) (*models.Operation, error) {
	return nil, nil
}

func handleUpload(exec boil.Executor, input interface{}) (*models.Operation, error) {
	r := input.(UploadRequest)

	log.Info("Creating operation")
	operation, err := createOperation(exec, OP_UPLOAD, r.Operation, nil)
	if err != nil {
		return nil, err
	}

	file, _, err := findFileBySHA1(exec, r.Sha1)
	if err != nil {
		if _, ok := err.(FileNotFound); ok {
			log.Info("File not found, creating new.")
			file, err = createFile(exec, nil, r.File, nil)
		} else {
			return nil, err
		}
	}

	log.Info("Updating file's properties")
	var fileProps = make(map[string]interface{})
	if file.Properties.Valid {
		file.Properties.Unmarshal(&fileProps)
	}
	fileProps["url"] = r.Url
	fileProps["duration"] = r.Duration
	fpa, _ := json.Marshal(fileProps)
	file.Properties = null.JSONFrom(fpa)

	log.Info("Saving changes to DB")
	err = file.Update(exec, "properties")
	if err != nil {
		return nil, err
	}

	log.Info("Associating file to operation")
	return operation, operation.AddFiles(exec, false, file)
}

// Helpers

// Generic operation handler.
// 	* Manage DB transactions
// 	* Call operation logic handler
// 	* Render JSON response
func handleOperation(c *gin.Context, input interface{},
	opHandler func(boil.Executor, interface{}) (*models.Operation, error)) {

	tx, err := boil.Begin()
	if err == nil {
		_, err = opHandler(tx, input)
		if err == nil {
			tx.Commit()
		} else {
			log.Error("Error handling operation: ", err)
			if txErr := tx.Rollback(); txErr != nil {
				log.Error("Error rolling back DB transaction: ", txErr)
			}
		}
	}

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		switch err.(type) {
		case FileNotFound:
			c.Error(err).SetType(gin.ErrorTypePublic)
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		default:
			c.Error(err).SetType(gin.ErrorTypePrivate)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		}
	}
}

type FileNotFound struct {
	Sha1 string
}

func (f FileNotFound) Error() string {
	return fmt.Sprintf("File not found, sha1 = %s", f.Sha1)
}

func createOperation(exec boil.Executor, name string, o Operation, properties map[string]interface{}) (*models.Operation, error) {
	operation := models.Operation{
		TypeID:  OPERATION_TYPE_REGISTRY.ByName[name].ID,
		UID:     utils.GenerateUID(8),
		Station: null.StringFrom(o.Station),
	}

	// Lookup user, skip if doesn't exist
	user, err := models.Users(exec, qm.Where("email=?", o.User)).One()
	if err == nil {
		operation.UserID = null.Int64From(user.ID)
	} else {
		if err == sql.ErrNoRows {
			log.Warnf("Unknown User [%s]. Skipping.", o.User)
		} else {
			return nil, err
		}
	}

	// Handle properties
	if o.WorkflowID != "" {
		if properties == nil {
			properties = make(map[string]interface{})
		}
		properties["workflow_id"] = o.WorkflowID
	}
	if properties != nil {
		props, err := json.Marshal(properties)
		if err != nil {
			return nil, err
		}
		operation.Properties = null.JSONFrom(props)
	}

	return &operation, operation.Insert(exec)
}

func createFile(exec boil.Executor, parent *models.File, f File, properties map[string]interface{}) (*models.File, error) {
	sha1, err := hex.DecodeString(f.Sha1)
	if err != nil {
		return nil, err
	}

	file := models.File{
		UID:           utils.GenerateUID(8),
		Name:          f.FileName,
		Sha1:          null.BytesFrom(sha1),
		Size:          f.Size,
		FileCreatedAt: null.TimeFrom(f.CreatedAt.Time),
		Type:          f.Type,
		SubType:       f.SubType,
	}

	if f.MimeType != "" {
		file.MimeType = null.StringFrom(f.MimeType)
	}
	if f.Language != "" {
		file.Language = null.StringFrom(f.Language)
	}
	if parent != nil {
		file.ParentID = null.Int64From(parent.ID)
	}

	// Handle properties
	if properties != nil {
		props, err := json.Marshal(properties)
		if err != nil {
			return nil, err
		}
		file.Properties = null.JSONFrom(props)
	}

	return &file, file.Insert(exec)
}

func findFileBySHA1(exec boil.Executor, sha1 string) (*models.File, []byte, error) {
	log.Info("Looking up file, sha1=", sha1)
	s, err := hex.DecodeString(sha1)
	if err != nil {
		return nil, nil, err
	}

	f, err := models.Files(exec, qm.Where("sha1=?", s)).One()
	if err == nil {
		return f, s, nil
	} else {
		if err == sql.ErrNoRows {
			return nil, s, FileNotFound{Sha1: sha1}
		} else {
			return nil, s, err
		}
	}
}
