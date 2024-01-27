package api

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/gin-gonic/gin.v1/binding"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/events"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

// Start capture of AV file, i.e. morning lesson, tv program, etc...
func CaptureStartHandler(c *gin.Context) {
	log.Info(common.OP_CAPTURE_START)
	var i CaptureStartRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleCaptureStart, nil)
	}
}

// Stop capture of AV file, ending a matching capture_start event.
// This is the first time a physical file is created in the studio.
func CaptureStopHandler(c *gin.Context) {
	log.Info(common.OP_CAPTURE_STOP)
	var i CaptureStopRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleCaptureStop, nil)
	}
}

// Demux manifest file to original and low resolution proxy
func DemuxHandler(c *gin.Context) {
	log.Info(common.OP_DEMUX)
	var i DemuxRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleDemux, nil)
	}
}

// Trim demuxed files at certain points
func TrimHandler(c *gin.Context) {
	log.Info(common.OP_TRIM)
	var i TrimRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleTrim, nil)
	}
}

// Final files sent from studio
func SendHandler(c *gin.Context) {
	log.Info(common.OP_SEND)
	var i SendRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleSend, sendResultRenderer)
	}
}

// Files converted to low resolution web formats, language splitting, etc...
func ConvertHandler(c *gin.Context) {
	log.Info(common.OP_CONVERT)
	var i ConvertRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleConvert, nil)
	}
}

// File uploaded to a public accessible URL
func UploadHandler(c *gin.Context) {
	log.Info(common.OP_UPLOAD)
	var i UploadRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleUpload, nil)
	}
}

// Sirtutim archive file generated
func SirtutimHandler(c *gin.Context) {
	log.Info(common.OP_SIRTUTIM)
	var i SirtutimRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleSirtutim, nil)
	}
}

// Insert new file to archive
func InsertHandler(c *gin.Context) {
	log.Info(common.OP_INSERT)
	var i InsertRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleInsert, insertResultRenderer)
	}
}

// A file in archive has been transcoded
func TranscodeHandler(c *gin.Context) {
	log.Info(common.OP_TRANSCODE)

	var i TranscodeRequest
	if c.BindJSON(&i) == nil {
		if i.Message != "" {
			handleOperation(c, i, handleTranscode, nil)
		} else {
			if err := binding.Validator.ValidateStruct(i.AsFile()); err != nil {
				c.AbortWithError(400, err).SetType(gin.ErrorTypeBind)
			} else {
				handleOperation(c, i, handleTranscode, nil)
			}
		}
	}
}

// Join multiple files sequentially
func JoinHandler(c *gin.Context) {
	log.Info(common.OP_JOIN)
	var i JoinRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleJoin, nil)
	}
}

// Replace HLS file (on add language, quality)
func Replace(c *gin.Context) {
	log.Info(common.OP_REPLACE)
	var i ReplaceRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleReplace, replaceResultRenderer)
	}
}

// This endpoint is used in the trim admin workflow
// We need this for the "fix" / "update" content unit flow.
// When a trim from capture is made to fix some unit on capture files
// with multiple descendant units.
func DescendantUnitsHandler(c *gin.Context) {
	mdb := c.MustGet("MDB").(*sql.DB)

	f, _, err := FindFileBySHA1(mdb, c.Param("sha1"))
	if err != nil {
		if _, ok := err.(FileNotFound); ok {
			NewNotFoundError().Abort(c)
		} else if _, ok := errors.Cause(err).(hex.InvalidByteError); ok {
			NewBadRequestError(err).Abort(c)
		} else if errors.Cause(err) == hex.ErrLength {
			NewBadRequestError(err).Abort(c)
		} else {
			NewInternalError(err).Abort(c)
		}
		return
	}

	// fetch content unit IDs of all files descendant from given file.
	// We deliberately exclude derived content units as they are not
	// explicitly created in the workflow (so no reason to change them)
	var cuIDs pq.Int64Array
	q := `WITH RECURSIVE rf AS (
  SELECT f.*
  FROM files f
  WHERE f.id = $1
  UNION
  SELECT f.*
  FROM files f INNER JOIN rf ON f.parent_id = rf.id
) SELECT array_agg(DISTINCT rf.content_unit_id)
  FILTER (WHERE rf.content_unit_id IS NOT NULL)
  FROM rf 
	INNER JOIN content_units cu ON rf.content_unit_id = cu.id AND NOT (cu.type_id = ANY($2))`
	err = queries.Raw(q, f.ID, pq.Array([]int64{
		//CONTENT_TYPE_REGISTRY.ByName[CT_KITEI_MAKOR].ID,  // created in workflow
		//CONTENT_TYPE_REGISTRY.ByName[CT_LELO_MIKUD].ID,   // created in workflow
		common.CONTENT_TYPE_REGISTRY.ByName[common.CT_PUBLICATION].ID,
	})).QueryRow(mdb).Scan(&cuIDs)
	if err != nil {
		NewInternalError(err).Abort(c)
		return
	}

	if len(cuIDs) == 0 {
		c.JSON(http.StatusOK, NewContentUnitsResponse())
		return
	}

	// fetch units by previously found IDs
	units, err := models.ContentUnits(
		qm.WhereIn("id in ?", utils.ConvertArgsInt64(cuIDs)...),
		qm.Load("ContentUnitI18ns")).
		All(mdb)
	if err != nil {
		NewInternalError(err).Abort(c)
		return
	}

	// i18n
	data := make([]*ContentUnit, len(units))
	for i, cu := range units {
		x := &ContentUnit{ContentUnit: *cu}
		data[i] = x
		x.I18n = make(map[string]*models.ContentUnitI18n, len(cu.R.ContentUnitI18ns))
		for _, i18n := range cu.R.ContentUnitI18ns {
			x.I18n[i18n.Language] = i18n
		}
	}

	c.JSON(http.StatusOK, ContentUnitsResponse{
		ListResponse: ListResponse{Total: int64(len(data))},
		ContentUnits: data,
	})
	return
}

// Handler logic

func handleCaptureStart(exec boil.Executor, input interface{}) (*models.Operation, []events.Event, error) {
	r := input.(CaptureStartRequest)

	log.Info("Creating operation")
	props := map[string]interface{}{
		"capture_source": r.CaptureSource,
		"collection_uid": r.CollectionUID,
	}
	operation, err := CreateOperation(exec, common.OP_CAPTURE_START, r.Operation, props)
	if err != nil {
		return nil, nil, err
	}

	log.Info("Creating file and associating to operation")
	uid, err := GetFreeUID(exec, new(FileUIDChecker))
	if err != nil {
		return nil, nil, err
	}
	file := models.File{
		UID:  uid,
		Name: r.FileName,
	}
	return operation, nil, operation.AddFiles(exec, true, &file)
}

func handleCaptureStop(exec boil.Executor, input interface{}) (*models.Operation, []events.Event, error) {
	r := input.(CaptureStopRequest)

	log.Info("Creating operation")
	props := map[string]interface{}{
		"capture_source": r.CaptureSource,
		"collection_uid": r.CollectionUID, // $LID = backup capture id when lesson, capture_id when program (part=false)
		"part":           r.Part,
	}
	operation, err := CreateOperation(exec, common.OP_CAPTURE_STOP, r.Operation, props)
	if err != nil {
		return nil, nil, err
	}

	log.Info("Looking up parent file, workflow_id=", r.Operation.WorkflowID)
	var parent *models.File
	var parentID int64
	err = queries.Raw(
		`SELECT file_id FROM files_operations
		 INNER JOIN operations ON operation_id = id
		 WHERE type_id=$1 AND properties -> 'workflow_id' ? $2`,
		common.OPERATION_TYPE_REGISTRY.ByName[common.OP_CAPTURE_START].ID,
		r.Operation.WorkflowID).
		QueryRow(exec).
		Scan(&parentID)
	if err == nil {
		parent = &models.File{ID: parentID}
	} else {
		if err == sql.ErrNoRows {
			log.Warnf("capture_start operation not found for workflow_id [%s]. Skipping.",
				r.Operation.WorkflowID)
		} else {
			return nil, nil, err
		}
	}

	log.Info("Creating file")
	fProps := make(map[string]interface{})
	if r.LabelID.Valid {
		fProps["label_id"] = r.LabelID.Int
	}
	file, err := CreateFile(exec, parent, r.File, fProps)
	if err != nil {
		return nil, nil, err
	}

	log.Info("Associating file to operation")
	return operation, nil, operation.AddFiles(exec, false, file)
}

func handleDemux(exec boil.Executor, input interface{}) (*models.Operation, []events.Event, error) {
	r := input.(DemuxRequest)

	parent, _, err := FindFileBySHA1(exec, r.Sha1)
	if err != nil {
		return nil, nil, err
	}

	opFiles := []*models.File{parent}

	log.Info("Creating operation")
	props := map[string]interface{}{
		"capture_source": r.CaptureSource,
	}
	operation, err := CreateOperation(exec, common.OP_DEMUX, r.Operation, props)
	if err != nil {
		return nil, nil, err
	}

	log.Info("Creating original")
	props = map[string]interface{}{
		"duration": r.Original.Duration,
	}
	original, err := CreateFile(exec, parent, r.Original.File, props)
	if err != nil {
		return nil, nil, err
	}
	opFiles = append(opFiles, original)

	if r.Proxy != nil {
		log.Info("Creating proxy")
		props = map[string]interface{}{
			"duration": r.Proxy.Duration,
		}
		proxy, err := CreateFile(exec, parent, r.Proxy.File, props)
		if err != nil {
			return nil, nil, err
		}
		opFiles = append(opFiles, proxy)
	} else {
		log.Info("Proxy not provided. Skipping Creation of proxy file")
	}

	log.Info("Associating files to operation")
	return operation, nil, operation.AddFiles(exec, false, opFiles...)
}

func handleTrim(exec boil.Executor, input interface{}) (*models.Operation, []events.Event, error) {
	r := input.(TrimRequest)

	// Fetch parent files
	original, _, err := FindFileBySHA1(exec, r.OriginalSha1)
	if err != nil {
		return nil, nil, err
	}
	opFiles := []*models.File{original}

	var proxy *models.File
	if r.ProxySha1 != "" {
		proxy, _, err = FindFileBySHA1(exec, r.ProxySha1)
		if err != nil {
			return nil, nil, err
		}
		opFiles = append(opFiles, proxy)
	}

	// TODO: in case of re-trim with the exact same parameters we already have the files in DB.
	// No need to return an error, a warning in the log is enough.

	log.Info("Creating operation")
	props := map[string]interface{}{
		"capture_source": r.CaptureSource,
		"in":             r.In,
		"out":            r.Out,
	}
	operation, err := CreateOperation(exec, common.OP_TRIM, r.Operation, props)
	if err != nil {
		return nil, nil, err
	}

	log.Info("Creating trimmed original")
	props = map[string]interface{}{
		"duration": r.Original.Duration,
	}
	originalTrim, err := CreateFile(exec, original, r.Original.File, props)
	if err != nil {
		return nil, nil, err
	}
	opFiles = append(opFiles, originalTrim)

	if proxy != nil && r.Proxy != nil {
		log.Info("Creating trimmed proxy")
		props = map[string]interface{}{
			"duration": r.Proxy.Duration,
		}
		proxyTrim, err := CreateFile(exec, proxy, r.Proxy.File, props)
		if err != nil {
			return nil, nil, err
		}
		opFiles = append(opFiles, proxyTrim)
	} else {
		log.Info("Proxy not provided. Skipping trimmed proxy creation")
	}

	log.Info("Associating files to operation")
	return operation, nil, operation.AddFiles(exec, false, opFiles...)
}

func handleSend(exec boil.Executor, input interface{}) (*models.Operation, []events.Event, error) {
	r := input.(SendRequest)

	// Original
	original, _, err := FindFileBySHA1(exec, r.Original.Sha1)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Lookup original file")
	}
	if original.Name == r.Original.FileName {
		log.Info("Original's name hasn't change")
	} else {
		log.Info("Renaming original")
		original.Name = r.Original.FileName
		_, err = original.Update(exec, boil.Whitelist("name"))
		if err != nil {
			return nil, nil, errors.Wrap(err, "Rename original file")
		}
	}

	opFiles := []*models.File{original}

	// Proxy
	var proxy *models.File
	if r.Proxy != nil {
		proxy, _, err = FindFileBySHA1(exec, r.Proxy.Sha1)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Lookup proxy file")
		}
		opFiles = append(opFiles, proxy)
		if proxy.Name == r.Proxy.FileName {
			log.Info("Proxy's name hasn't change")
		} else {
			log.Info("Renaming proxy")
			proxy.Name = r.Proxy.FileName
			_, err = proxy.Update(exec, boil.Whitelist("name"))
			if err != nil {
				return nil, nil, errors.Wrap(err, "Rename proxy file")
			}
		}
	} else {
		log.Info("Proxy not provided. Skipping possible rename")
	}

	// Source
	var source *models.File
	if r.Source != nil {
		source, _, err = FindFileBySHA1(exec, r.Source.Sha1)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Lookup source file")
		}
		opFiles = append(opFiles, source)
		if source.Name == r.Source.FileName {
			log.Info("Source's name hasn't change")
		} else {
			log.Info("Renaming source")
			source.Name = r.Source.FileName
			_, err = source.Update(exec, boil.Whitelist("name"))
			if err != nil {
				return nil, nil, errors.Wrap(err, "Rename source file")
			}
		}
	} else {
		log.Info("Source not provided. Skipping possible rename")
	}

	mode := "new"
	if r.Mode.Valid {
		mode = r.Mode.String
	}

	log.Infof("Processing CIT Metadata: %s mode", mode)
	var evnts []events.Event
	if mode == "new" {
		evnts, err = ProcessCITMetadata(exec, r.Metadata, original, proxy, source)
	} else {
		evnts, err = ProcessCITMetadataUpdate(exec, r.Metadata, original, proxy, source)
	}
	if err != nil {
		return nil, nil, errors.Wrap(err, "Process CIT Metadata")
	}

	log.Info("Creating operation")
	props := make(map[string]interface{})
	b, err := json.Marshal(r.Metadata)
	if err != nil {
		return nil, nil, errors.Wrap(err, "json Marshal")
	}
	if err = json.Unmarshal(b, &props); err != nil {
		return nil, nil, errors.Wrap(err, "json Unmarshal")
	}
	operation, err := CreateOperation(exec, common.OP_SEND, r.Operation, props)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Create operation")
	}

	log.Info("Associating files to operation")
	err = operation.AddFiles(exec, false, opFiles...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Associate files")
	}

	log.Infof("Updating unit workflow_id to: %s", r.WorkflowID)
	err = original.L.LoadContentUnit(exec, true, original, nil)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "refresh original's unit: %d", original.ContentUnitID.Int64)
	}
	err = UpdateContentUnitProperties(exec, original.R.ContentUnit, map[string]interface{}{
		"workflow_id": r.WorkflowID,
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unit workflow_id: %d", original.ContentUnitID.Int64)
	}

	return operation, evnts, nil
}

func sendResultRenderer(c *gin.Context, exec boil.Executor, input interface{}, op *models.Operation) error {
	i := input.(SendRequest)
	original, _, _ := FindFileBySHA1(exec, i.Original.Sha1)
	cu, err := models.FindContentUnit(exec, original.ContentUnitID.Int64)
	if err != nil {
		return errors.Wrapf(err, "Lookup content unit")
	}

	c.JSON(http.StatusOK, cu)
	return nil
}

func handleConvert(exec boil.Executor, input interface{}) (*models.Operation, []events.Event, error) {
	r := input.(ConvertRequest)

	in, _, err := FindFileBySHA1(exec, r.Sha1)
	if err != nil {
		if _, ok := err.(FileNotFound); ok {
			log.Infof("Parent file not found, noop.")
			return nil, nil, nil
		} else {
			return nil, nil, errors.Wrapf(err, "lookup parent")
		}
	}

	log.Info("Creating operation")
	operation, err := CreateOperation(exec, common.OP_CONVERT, r.Operation, nil)
	if err != nil {
		return nil, nil, err
	}

	// We dedup duplicate files here.
	// These are usually files without sound translation.
	// So they are identical, i.e. same SHA1
	log.Info("Deduping output files by SHA1")
	uniq := make(map[string]int)
	for i := range r.Output {
		uniq[r.Output[i].Sha1] = i
	}
	log.Infof("%d unique files out of %d", len(uniq), len(r.Output))

	log.Info("Creating output files")
	evnts := make([]events.Event, 0)
	files := make([]*models.File, len(uniq)+1)
	files[0] = in
	i := 0
	for _, v := range uniq {
		x := r.Output[v]
		props := createHLSProps(x)
		// lookup by sha1 as it might be a "reconvert"
		f, _, err := FindFileBySHA1(exec, x.Sha1)
		if err == nil {
			log.Infof("File already exists, updating: %s", x.Sha1)
			err = UpdateFile(exec, f, in, x.File, props)
			if err != nil {
				return nil, nil, errors.Wrap(err, "Update file")
			}

			// restore files that were removed as they are now recreated a fresh.
			// they were probably removed in update mode to send operation, aka fix.
			if f.RemovedAt.Valid {
				f.RemovedAt = null.NewTime(time.Unix(0, 0), false)
				_, err = f.Update(exec, boil.Whitelist("removed_at"))
				if err != nil {
					return nil, nil, errors.Wrap(err, "Restore file")
				}
			}

			// TODO: Here we might change an unit published status...

			evnts = append(evnts, events.FileUpdateEvent(f))
		} else {
			if _, ok := err.(FileNotFound); ok {
				// new file
				f, err = CreateFile(exec, in, x.File, props)
				if err != nil {
					return nil, nil, errors.Wrap(err, "Create file")
				}
			} else {
				return nil, nil, errors.Wrap(err, "Lookup file in DB")
			}
		}

		i++
		files[i] = f
	}

	log.Info("Associating files to operation")
	return operation, evnts, operation.AddFiles(exec, false, files...)
}

func handleUpload(exec boil.Executor, input interface{}) (*models.Operation, []events.Event, error) {
	r := input.(UploadRequest)

	log.Info("Creating operation")
	operation, err := CreateOperation(exec, common.OP_UPLOAD, r.Operation, nil)
	if err != nil {
		return nil, nil, err
	}

	file, _, err := FindFileBySHA1(exec, r.Sha1)
	if err != nil {
		if _, ok := err.(FileNotFound); ok {
			log.Info("File not found, creating new.")
			file, err = CreateFile(exec, nil, r.File, nil)
		} else {
			return nil, nil, err
		}
	}
	log.Info("Updating file's properties")
	var fileProps = make(map[string]interface{})
	if file.Properties.Valid {
		err = file.Properties.Unmarshal(&fileProps)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Unmarshal file properties [%d]", file.ID)
		}
	}
	fileProps["url"] = r.Url
	fileProps["duration"] = r.Duration
	fpa, _ := json.Marshal(fileProps)
	file.Properties = null.JSONFrom(fpa)

	log.Info("Saving changes to DB")
	_, err = file.Update(exec, boil.Whitelist("properties"))
	if err != nil {
		return nil, nil, err
	}

	impact, err := PublishFile(exec, file)
	if err != nil {
		return nil, nil, err
	}

	log.Info("Check if file must replace")
	if err := file.L.LoadOperations(exec, true, file, nil); err != nil {
		return nil, nil, err
	}
	var opRepl *models.Operation
	for _, o := range file.R.Operations {
		if common.OPERATION_TYPE_REGISTRY.ByName[common.OP_REPLACE].ID == o.TypeID {
			opRepl = o
		}
	}
	evnts := make([]events.Event, 0)
	if opRepl != nil {
		log.Infof("special for Replace operation %d", opRepl.ID)
		//find old files by replace operation
		if err := opRepl.L.LoadFiles(exec, true, opRepl, nil); err != nil {
			return nil, nil, err
		}
		var oldFile *models.File
		for _, f := range opRepl.R.Files {
			if f.ID != file.ID {
				oldFile = f
				break
			}
		}
		if oldFile == nil {
			return nil, nil, errors.New("no file that was replaced")
		}
		_evnts, err := removeDescendants(exec, oldFile)
		if err != nil {
			return nil, nil, err
		}
		evnts = append(evnts, _evnts...)
	}
	evnts = append(evnts, events.FilePublishedEvent(file))
	evnts = append(evnts, impact.Events()...)

	log.Info("Associating file to operation")
	return operation, evnts, operation.AddFiles(exec, false, file)
}

func handleSirtutim(exec boil.Executor, input interface{}) (*models.Operation, []events.Event, error) {
	r := input.(SirtutimRequest)

	log.Info("Creating operation")
	operation, err := CreateOperation(exec, common.OP_SIRTUTIM, r.Operation, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Info("Creating file")
	r.File.Type = "image"
	file, err := CreateFile(exec, nil, r.File, nil)
	if err != nil {
		return nil, nil, err
	}

	// content_unit association
	var original *models.File
	if r.OriginalSha1 != "" {
		log.Info("Looking up content unit by original sha1")
		original, _, err = FindFileBySHA1(exec, r.OriginalSha1)
		if err != nil {
			if _, ok := err.(FileNotFound); ok {
				log.Warnf("Original file not found [%s]", r.OriginalSha1)
			} else {
				return nil, nil, err
			}
		} else {
			if original.ContentUnitID.Valid {
				log.Infof("Associating to content_unit [%d]", original.ContentUnitID.Int64)
				file.ContentUnitID = original.ContentUnitID
				_, err = file.Update(exec, boil.Whitelist("content_unit_id"))
				if err != nil {
					return nil, nil, errors.Wrap(err, "Save file.content_unit_id to DB")
				}
			} else {
				log.Warn("Original file is not associated to any content unit")
			}
		}
	}

	log.Info("Associating files to operation")
	if original == nil {
		return operation, nil, operation.AddFiles(exec, false, file)
	} else {
		return operation, nil, operation.AddFiles(exec, false, original, file)
	}
}

// modes:
// new - create new file
// update - create new file and remove old file
// rename - change metadata of previously inserted file
func handleInsert(exec boil.Executor, input interface{}) (*models.Operation, []events.Event, error) {
	r := input.(InsertRequest)

	log.Infof("Lookup file by SHA1")
	file, _, err := FindFileBySHA1(exec, r.File.Sha1)
	if err != nil {
		if _, ok := err.(FileNotFound); !ok {
			return nil, nil, err
		}
	}

	// validate user input for based on r.Mode and file existence
	if r.Mode == "new" && file != nil {
		return nil, nil, errors.Errorf("File already exist")
	}
	if r.Mode == "rename" && err != nil {
		return nil, nil, err
	}

	opFiles := make([]*models.File, 0)
	var oldFile *models.File
	if r.Mode == "update" {
		log.Infof("Lookup old file by SHA1")
		oldFile, _, err = FindFileBySHA1(exec, r.OldSha1)
		if err != nil {
			if _, ok := err.(FileNotFound); ok {
				return nil, nil, errors.Wrap(err, "Old file not found")
			} else {
				return nil, nil, err
			}
		}
		opFiles = append(opFiles, oldFile)
	}

	// Content Unit to whom which the inserted file should belong to
	// In case the unit already exists, we lookup in DB with UID.
	// In case we're to create a new unit, we do so below based on CITMetadata.
	var cu *models.ContentUnit

	if r.ContentUnitUID != "" {
		log.Infof("Lookup content unit by uid %s", r.ContentUnitUID)
		cu, err = models.ContentUnits(
			qm.Where("uid = ?", r.ContentUnitUID),
			qm.Load("SourceContentUnitDerivations"),
			qm.Load("SourceContentUnitDerivations.Derived"),
		).One(exec)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Fetch unit from mdb")
		}
		log.Infof("Found content unit %d", cu.ID)
	}

	var parent *models.File
	if r.ParentSha1 != "" {
		log.Info("Looking up parent file by sha1")
		parent, _, err = FindFileBySHA1(exec, r.ParentSha1)
		if err != nil {
			if _, ok := err.(FileNotFound); ok {
				log.Warnf("Parent file not found [%s]", r.ParentSha1)
			} else {
				return nil, nil, err
			}
		} else {
			opFiles = append(opFiles, parent)
		}
	}

	log.Info("Creating operation")
	props := map[string]interface{}{
		"insert_type": r.InsertType,
		"mode":        r.Mode,
	}
	operation, err := CreateOperation(exec, common.OP_INSERT, r.Operation, props)
	if err != nil {
		return nil, nil, err
	}

	// process metadata
	log.Info("Processing metadata")
	if r.File.Type == "" {
		switch r.InsertType {
		case "akladot", "tamlil", "kitei-makor", "article":
			r.File.Type = "text"
		case "sirtutim", "publication":
			r.File.Type = "image"
		case "aricha":
			r.File.Type = "video"
		case "subtitles":
			r.File.Type = "subtitles"
		default:
			r.File.Type = ""
		}
	}

	if r.AVFile.Duration > 0 {
		props["duration"] = r.AVFile.Duration
	}
	if r.AVFile.VideoSize != "" {
		props["video_size"] = r.AVFile.VideoSize
	}

	// create new file based on mode
	if r.Mode == "new" || r.Mode == "update" {
		log.Info("Creating new file")
		file, err = CreateFile(exec, parent, r.File, props)
		if err != nil {
			return nil, nil, err
		}
	} else if r.Mode == "rename" {
		log.Info("Renaming existing file")

		// make a temp file
		mf, err := makeFile(parent, r.File, props)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Make file")
		}

		// set new attributes
		file.Name = mf.Name
		file.Type = mf.Type
		file.SubType = mf.SubType
		file.MimeType = mf.MimeType
		file.Language = mf.Language
		file.Properties = mf.Properties
		file.ParentID = mf.ParentID

		// save
		_, err = file.Update(exec, boil.Infer())
		if err != nil {
			return nil, nil, errors.Wrap(err, "Update file")
		}
	}

	opFiles = append(opFiles, file)

	var cuID int64
	if cu != nil {
		cuID = cu.ID
	}

	// special types logic
	if r.InsertType == "kitei-makor" ||
		r.InsertType == "research-material" {
		log.Infof("%s, associating to derived unit", r.InsertType)

		var ct string
		switch r.InsertType {
		case "kitei-makor":
			ct = common.CT_KITEI_MAKOR
		case "research-material":
			ct = common.CT_RESEARCH_MATERIAL
		}

		var cudID int64
		if len(cu.R.SourceContentUnitDerivations) > 0 {
			for _, cud := range cu.R.SourceContentUnitDerivations {
				if common.CONTENT_TYPE_REGISTRY.ByID[cud.R.Derived.TypeID].Name == ct {
					cudID = cud.DerivedID
					break
				}
			}
		}

		if cudID > 0 {
			cuID = cudID
			log.Infof("%s derived unit exist: %d", ct, cuID)
		} else {
			log.Infof("%s derived unit doesn't exists. Creating...", ct)
			ktCU, err := CreateContentUnit(exec, ct, nil)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "Create %s derived unit", ct)
			}

			cud := &models.ContentUnitDerivation{
				SourceID: cu.ID,
				Name:     ct,
			}
			err = ktCU.AddDerivedContentUnitDerivations(exec, true, cud)
			if err != nil {
				return nil, nil, errors.Wrap(err, "Save CUD in DB")
			}

			cuID = ktCU.ID
		}
	} else if r.InsertType == "publication" {
		log.Infof("Publication, associating to derived unit")

		publisher, err := models.Publishers(qm.Where("uid = ?", r.PublisherUID)).One(exec)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Fetch publisher from mdb")
		}

		var cudID int64
		err = queries.Raw(`SELECT cu.id
FROM content_units cu
  INNER JOIN content_unit_derivations cud ON cu.id = cud.derived_id AND cud.source_id = $1 AND cu.type_id = $2
  INNER JOIN content_units_publishers cup ON cud.derived_id = cup.content_unit_id AND cup.publisher_id = $3`,
			cu.ID, common.CONTENT_TYPE_REGISTRY.ByName[common.CT_PUBLICATION].ID, publisher.ID).
			QueryRow(exec).
			Scan(&cudID)
		if err != nil && err != sql.ErrNoRows {
			return nil, nil, errors.Wrap(err, "Lookup existing cud from mdb")
		}

		if cudID > 0 {
			cuID = cudID
			log.Infof("PUBLICATION derived unit exist: %d", cuID)
		} else {
			log.Infof("PUBLICATION derived unit doesn't exists. Creating...")
			pCU, err := CreateContentUnit(exec, common.CT_PUBLICATION, map[string]interface{}{
				"original_language": file.Language.String,
			})
			if err != nil {
				return nil, nil, errors.Wrap(err, "Create PUBLICATION derived unit")
			}

			err = pCU.AddPublishers(exec, false, publisher)
			if err != nil {
				return nil, nil, errors.Wrap(err, "Associate publisher with cud")
			}

			cud := &models.ContentUnitDerivation{
				SourceID: cu.ID,
				Name:     common.CT_PUBLICATION,
			}
			err = pCU.AddDerivedContentUnitDerivations(exec, true, cud)
			if err != nil {
				return nil, nil, errors.Wrap(err, "Save CUD in DB")
			}

			cuID = pCU.ID
		}
	} else if r.InsertType == "declamation" {
		if r.Mode == "rename" {
			log.Info("declamation, skipping new content unit creation on rename")
		} else {
			log.Info("declamation, creating new content unit")
			if r.Metadata == nil {
				return nil, nil, NewBadRequestError(errors.New("Metadata is required"))
			}

			filmDate := r.Metadata.CaptureDate
			if r.Metadata.FilmDate != nil {
				filmDate = *r.Metadata.FilmDate
			}

			cu, err = CreateContentUnit(exec, common.CT_BLOG_POST, map[string]interface{}{
				"film_date":         filmDate,
				"original_language": common.StdLang(r.Metadata.Language),
			})
			if err != nil {
				return nil, nil, errors.Wrap(err, "Create declamation content unit")
			}
			cuID = cu.ID

			log.Infof("Describing content unit [%d]", cu.ID)
			err = DescribeContentUnit(exec, cu, CITMetadata{ContentType: common.CT_BLOG_POST})
			if err != nil {
				log.Errorf("Error describing content unit: %s", err.Error())
			}
		}
	}

	log.Infof("Associating file [%d] to content_unit [%d]", file.ID, cuID)
	file.ContentUnitID = null.Int64From(cuID)
	_, err = file.Update(exec, boil.Whitelist("content_unit_id"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "Save file.content_unit_id to DB")
	}

	// TODO: shouldn't we move this up and emit events for new content units as well ?
	evnts := make([]events.Event, 0)

	// remove oldFile in update mode and collect events
	if r.Mode == "new" {
		evnts = append(evnts, events.FileInsertEvent(file, r.InsertType))
	} else if r.Mode == "rename" {
		evnts = append(evnts, events.FileUpdateEvent(file))
	} else if r.Mode == "update" {
		evnts = append(evnts, events.FileReplaceEvent(oldFile, file, r.InsertType))

		if impact, err := RemoveFile(exec, oldFile); err != nil {
			return nil, nil, errors.Wrapf(err, "Remove old file %d", oldFile.ID)
		} else {
			evnts = append(evnts, impact.Events()...)
		}
	}

	// Set unit duration if not already set
	// This is here for units not created via send operation of some trimmed file.
	if r.AVFile.Duration > 0 {
		cu, err := models.ContentUnits(
			qm.Where("id = ?", cuID),
			qm.Load("Files")).
			One(exec)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Refresh CU [%d] from DB", cuID)
		}

		if len(cu.R.Files) == 1 {
			err = UpdateContentUnitProperties(exec, cu, map[string]interface{}{"duration": int(r.AVFile.Duration)})
			if err != nil {
				return nil, nil, errors.Wrapf(err, "Update CU properties [%d]", cuID)
			}
			evnts = append(evnts, events.ContentUnitUpdateEvent(cu))
		}
	}

	log.Info("Associating files to operation")
	return operation, evnts, operation.AddFiles(exec, false, opFiles...)
}

func insertResultRenderer(c *gin.Context, exec boil.Executor, input interface{}, op *models.Operation) error {
	i := input.(InsertRequest)

	if i.InsertType != "declamation" {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return nil
	}

	f, _, _ := FindFileBySHA1(exec, i.Sha1)
	cu, err := models.FindContentUnit(exec, f.ContentUnitID.Int64)
	if err != nil {
		return errors.Wrapf(err, "Lookup content unit")
	}

	c.JSON(http.StatusOK, cu)
	return nil
}

func handleTranscode(exec boil.Executor, input interface{}) (*models.Operation, []events.Event, error) {
	r := input.(TranscodeRequest)

	if r.Message != "" {
		log.Infof("Transcode Error: %s", r.Message)

		log.Info("Creating operation")
		opProps := map[string]interface{}{
			"message": r.Message,
		}

		operation, err := CreateOperation(exec, common.OP_TRANSCODE, r.Operation, opProps)
		if err != nil {
			return nil, nil, err
		}

		log.Info("Looking up original file")
		original, _, err := FindFileBySHA1(exec, r.OriginalSha1)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Lookup original file %s", r.OriginalSha1)
		}

		log.Info("Updating queue table")
		_, err = queries.Raw("update batch_convert set operation_id=$1 where file_id=$2", operation.ID, original.ID).
			Exec(exec)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Update queue table")
		}

		return operation, nil, operation.AddFiles(exec, false, original)
	}

	log.Info("Creating operation")
	operation, err := CreateOperation(exec, common.OP_TRANSCODE, r.Operation, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Info("Looking up original file")
	original, _, err := FindFileBySHA1(exec, r.OriginalSha1)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Lookup original file %s", r.OriginalSha1)
	}

	log.Info("Creating file")
	mt := common.MEDIA_TYPE_REGISTRY.ByExtension["mp4"]
	r.MaybeFile.Type = mt.Type
	r.MaybeFile.MimeType = mt.MimeType

	if original.Language.Valid {
		r.MaybeFile.Language = original.Language.String
	}

	var props map[string]interface{}
	if original.Properties.Valid {
		var oProps map[string]interface{}
		err := json.Unmarshal(original.Properties.JSON, &oProps)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "json.Unmarshal original properties [%d]", original.ID)
		}
		if duration, ok := oProps["duration"]; ok {
			props = make(map[string]interface{})
			props["duration"] = duration
		}
	}
	file, err := CreateFile(exec, original, r.AsFile(), props)
	if err != nil {
		return nil, nil, err
	}

	log.Info("Updating file secure published information")
	file.Published = true
	if original != nil {
		file.Secure = original.Secure
	}
	_, err = file.Update(exec, boil.Whitelist("secure", "published"))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Update secure published [%d]", file.ID)
	}

	log.Info("Updating queue table")
	_, err = queries.Raw("update batch_convert set operation_id=$1 where file_id=$2", operation.ID, original.ID).
		Exec(exec)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Update queue table")
	}

	opFiles := []*models.File{original, file}

	log.Info("Loading original's children")
	err = original.L.LoadParentFiles(exec, true, original, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Load original children")
	}

	log.Infof("Original has %d children", len(original.R.ParentFiles))
	for i := range original.R.ParentFiles {
		child := original.R.ParentFiles[i]
		if child.ID == file.ID {
			continue
		}

		log.Infof("Removing child file [%d]", child.ID)
		child.RemovedAt = null.TimeFrom(time.Now().UTC())
		_, err := child.Update(exec, boil.Whitelist("removed_at"))
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Save file to DB: %d", child.ID)
		}
		opFiles = append(opFiles, child)
	}

	log.Info("Associating files to operation")
	return operation, nil, operation.AddFiles(exec, false, opFiles...)
}

func handleJoin(exec boil.Executor, input interface{}) (*models.Operation, []events.Event, error) {
	r := input.(JoinRequest)

	// Fetch input files
	inOriginals := make([]*models.File, 0)
	for i := range r.OriginalShas {
		f, _, err := FindFileBySHA1(exec, r.OriginalShas[i])
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Original number %d, sha1 %s", i+1, r.OriginalShas[i])
		}
		inOriginals = append(inOriginals, f)
	}

	inProxies := make([]*models.File, 0)
	for i := range r.ProxyShas {
		f, _, err := FindFileBySHA1(exec, r.ProxyShas[i])
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Proxy number %d, sha1 %s", i+1, r.ProxyShas[i])
		}
		inProxies = append(inProxies, f)
	}

	opFiles := append(inOriginals, inProxies...)

	log.Info("Creating operation")
	props := map[string]interface{}{
		"original_shas": r.OriginalShas,
		"proxy_shas":    r.ProxyShas,
	}
	operation, err := CreateOperation(exec, common.OP_JOIN, r.Operation, props)
	if err != nil {
		return nil, nil, err
	}

	log.Info("Creating joined original")
	props = map[string]interface{}{
		"duration": r.Original.Duration,
	}
	original, err := CreateFile(exec, nil, r.Original.File, props)
	if err != nil {
		return nil, nil, err
	}
	opFiles = append(opFiles, original)

	if r.Proxy != nil {
		log.Info("Creating joined proxy")
		props = map[string]interface{}{
			"duration": r.Proxy.Duration,
		}
		proxy, err := CreateFile(exec, nil, r.Proxy.File, props)
		if err != nil {
			return nil, nil, err
		}
		opFiles = append(opFiles, proxy)
	} else {
		log.Info("No proxy provided. skipping joined proxy creation")
	}

	log.Info("Associating files to operation")
	return operation, nil, operation.AddFiles(exec, false, opFiles...)
}

func handleReplace(exec boil.Executor, input interface{}) (*models.Operation, []events.Event, error) {
	r := input.(ReplaceRequest)

	opFiles := make([]*models.File, 0)

	log.Infof("Lookup file by SHA1 %s", r.File.Sha1)
	fileNil, _, err := FindFileBySHA1(exec, r.File.Sha1)
	if fileNil != nil {
		return nil, nil, errors.New(fmt.Sprintf("new file allready created, SHA %s", r.File.Sha1))
	}
	if _, ok := err.(FileNotFound); !ok {
		return nil, nil, err
	}

	log.Infof("Lookup old file by SHA1 %s", r.File.Sha1)
	oldFile, _, err := FindFileBySHA1(exec, r.OldSha1)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("No file for replace, SHA %s", r.OldSha1))
	}
	opFiles = append(opFiles, oldFile)

	log.Info("Creating operation")
	operation, err := CreateOperation(exec, common.OP_REPLACE, r.Operation, map[string]interface{}{"mode": "replace"})
	if err != nil {
		return nil, nil, err
	}

	// create new file based on old
	log.Info("Creating new file")
	var parent *models.File
	if oldFile.ParentID.Valid {
		parent, err = models.FindFile(exec, oldFile.ParentID.Int64)
		if _, ok := err.(FileNotFound); err != nil && !ok {
			return nil, nil, err
		}
	}
	props := createHLSProps(r.HLSFile)
	file, err := CreateFile(exec, parent, r.File, props)
	if err != nil {
		return nil, nil, err
	}
	log.Infof("Associating file [%d] to content_unit [%d]", file.ID, oldFile.ContentUnitID.Int64)
	file.ContentUnitID = oldFile.ContentUnitID
	_, err = file.Update(exec, boil.Whitelist("content_unit_id"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "Save file.content_unit_id to DB")
	}

	evnts, err := removeDescendants(exec, oldFile)
	if err != nil {
		return nil, nil, err
	}

	opFiles = append(opFiles, file)
	log.Info("Associating files to operation")
	return operation, evnts, operation.AddFiles(exec, false, opFiles...)
}

func replaceResultRenderer(c *gin.Context, exec boil.Executor, input interface{}, op *models.Operation) error {
	r := input.(ReplaceRequest)
	file, _, err := FindFileBySHA1(exec, r.File.Sha1)
	if err != nil {
		return errors.Wrapf(err, "new file not found")
	}

	c.JSON(http.StatusOK, gin.H{"new_file_uid": fmt.Sprintf("%s", file.UID)})
	return nil
}

// Helpers

type OpHandlerFunc func(boil.Executor, interface{}) (*models.Operation, []events.Event, error)
type OpResponseRenderFunc func(*gin.Context, boil.Executor, interface{}, *models.Operation) error

func defaultResultRenderer(c *gin.Context, exec boil.Executor, input interface{}, op *models.Operation) error {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return nil
}

// Generic operation handler.
//   - Manage DB transactions
//   - Call operation logic handler
//   - Handle errors
//   - Emit events
//   - Render JSON response
func handleOperation(c *gin.Context, input interface{}, opHandler OpHandlerFunc, resFunc OpResponseRenderFunc) {
	mdb := c.MustGet("MDB").(*sql.DB)
	tx, err := mdb.Begin()
	utils.Must(err)

	// recover from panics in transaction
	defer func() {
		if p := recover(); p != nil {
			if ex := tx.Rollback(); ex != nil {
				log.Error("Couldn't roll back transaction")
			}
			panic(p) // re-throw panic after Rollback
		}
	}()

	// call handler and conclude transaction
	op, evnts, err := opHandler(tx, input)
	if err == nil {
		utils.Must(tx.Commit())
	} else {
		utils.Must(tx.Rollback())
	}

	// on success, emit events and call renderer
	if err == nil {
		emitEvents(c, evnts...)

		if resFunc == nil {
			resFunc = defaultResultRenderer
		}
		err = resFunc(c, mdb, input, op)
	}

	// handle errors
	if err != nil {
		switch err.(type) {
		case FileNotFound:
			NewBadRequestError(err).Abort(c)
		case *HttpError:
			err.(*HttpError).Abort(c)
		default:
			err = errors.WithMessagef(err, "Handle operation %s", c.HandlerName())
			NewInternalError(err).Abort(c)
		}
	}
}

func createHLSProps(d HLSFile) map[string]interface{} {
	props := make(map[string]interface{})
	props["duration"] = d.Duration
	if d.VideoSize != "" {
		props["video_size"] = d.VideoSize
	}

	if d.Languages != nil {
		languages := make([]string, len(d.Languages))
		for i, l := range d.Languages {
			languages[i] = common.StdLang(l)
		}
		props["languages"] = languages
	}
	if d.Qualities != nil {
		props["video_qualities"] = d.Qualities
	}
	return props
}

func removeDescendants(exec boil.Executor, file *models.File) ([]events.Event, error) {
	evnts := make([]events.Event, 0)
	if file.RemovedAt.Valid {
		return evnts, nil
	}
	forRemove, err := FindFileDescendants(exec, file.ID)
	if err != nil {
		return nil, err
	}
	forRemoveIds := []int64{file.ID}
	now := time.Now().UTC()
	for _, f := range forRemove {
		err = UpdateFileProperties(exec, f, map[string]interface{}{"replaced": now})
		if err != nil {
			return nil, err
		}
		forRemoveIds = append(forRemoveIds, f.ID)
	}
	_, err = models.Files(models.FileWhere.ID.IN(forRemoveIds)).UpdateAll(
		exec, models.M{"removed_at": now, "published": false},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "remove descendants of file %d ", file.ID)
	}
	for _, f := range forRemove {
		evnts = append(evnts, events.FileRemoveEvent(f))
	}
	return evnts, nil
}
