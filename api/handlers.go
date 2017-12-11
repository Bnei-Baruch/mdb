package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	//_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/gin-gonic/gin.v1/binding"
	"gopkg.in/volatiletech/null.v6"

	"github.com/Bnei-Baruch/mdb/events"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
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
		// we skip using the generic handleOperation to return some data to caller

		mdb := c.MustGet("MDB").(*sql.DB)
		tx, err := mdb.Begin()
		utils.Must(err)

		_, evnts, err := handleSend(tx, i)
		if err != nil {
			utils.Must(tx.Rollback())
			err = errors.Wrapf(err, "Handle operation")
			NewInternalError(err).Abort(c)
			return
		}

		utils.Must(tx.Commit())

		// fetch newly created content unit
		original, _, _ := FindFileBySHA1(mdb, i.Original.Sha1)
		cu, err := models.FindContentUnit(mdb, original.ContentUnitID.Int64)
		if err != nil {
			err = errors.Wrapf(err, "Lookup content unit")
		}

		if err == nil {
			emitEvents(c, evnts...)
			c.JSON(http.StatusOK, cu)
		} else {
			NewInternalError(err).Abort(c)
		}
	}
}

// Files converted to low resolution web formats, language splitting, etc...
func ConvertHandler(c *gin.Context) {
	log.Info(OP_CONVERT)
	var i ConvertRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleConvert)
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

// Sirtutim archive file generated
func SirtutimHandler(c *gin.Context) {
	log.Info(OP_SIRTUTIM)
	var i SirtutimRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleSirtutim)
	}
}

// Insert new file to archive
func InsertHandler(c *gin.Context) {
	log.Info(OP_INSERT)
	var i InsertRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleInsert)
	}
}

// A file in archive has been transcoded
func TranscodeHandler(c *gin.Context) {
	log.Info(OP_TRANSCODE)

	var i TranscodeRequest
	if c.BindJSON(&i) == nil {
		if i.Message != "" {
			handleOperation(c, i, handleTranscode)
		} else {
			if err := binding.Validator.ValidateStruct(i.AsFile()); err != nil {
				c.AbortWithError(400, err).SetType(gin.ErrorTypeBind)
			} else {
				handleOperation(c, i, handleTranscode)
			}
		}
	}
}

// Handler logic

func handleCaptureStart(exec boil.Executor, input interface{}) (*models.Operation, []events.Event, error) {
	r := input.(CaptureStartRequest)

	log.Info("Creating operation")
	props := map[string]interface{}{
		"capture_source": r.CaptureSource,
		"collection_uid": r.CollectionUID,
	}
	operation, err := CreateOperation(exec, OP_CAPTURE_START, r.Operation, props)
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
	operation, err := CreateOperation(exec, OP_CAPTURE_STOP, r.Operation, props)
	if err != nil {
		return nil, nil, err
	}

	log.Info("Looking up parent file, workflow_id=", r.Operation.WorkflowID)
	var parent *models.File
	var parentID int64
	err = queries.Raw(exec,
		`SELECT file_id FROM files_operations
		 INNER JOIN operations ON operation_id = id
		 WHERE type_id=$1 AND properties -> 'workflow_id' ? $2`,
		OPERATION_TYPE_REGISTRY.ByName[OP_CAPTURE_START].ID,
		r.Operation.WorkflowID).
		QueryRow().
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
	file, err := CreateFile(exec, parent, r.File, nil)
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

	log.Info("Creating operation")
	props := map[string]interface{}{
		"capture_source": r.CaptureSource,
	}
	operation, err := CreateOperation(exec, OP_DEMUX, r.Operation, props)
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

	log.Info("Creating proxy")
	props = map[string]interface{}{
		"duration": r.Proxy.Duration,
	}
	proxy, err := CreateFile(exec, parent, r.Proxy.File, props)
	if err != nil {
		return nil, nil, err
	}

	log.Info("Associating files to operation")
	return operation, nil, operation.AddFiles(exec, false, parent, original, proxy)
}

func handleTrim(exec boil.Executor, input interface{}) (*models.Operation, []events.Event, error) {
	r := input.(TrimRequest)

	// Fetch parent files
	original, _, err := FindFileBySHA1(exec, r.OriginalSha1)
	if err != nil {
		return nil, nil, err
	}
	proxy, _, err := FindFileBySHA1(exec, r.ProxySha1)
	if err != nil {
		return nil, nil, err
	}

	// TODO: in case of re-trim with the exact same parameters we already have the files in DB.
	// No need to return an error, a warning in the log is enough.

	log.Info("Creating operation")
	props := map[string]interface{}{
		"capture_source": r.CaptureSource,
		"in":             r.In,
		"out":            r.Out,
	}
	operation, err := CreateOperation(exec, OP_TRIM, r.Operation, props)
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

	log.Info("Creating trimmed proxy")
	props = map[string]interface{}{
		"duration": r.Proxy.Duration,
	}
	proxyTrim, err := CreateFile(exec, proxy, r.Proxy.File, props)
	if err != nil {
		return nil, nil, err
	}

	log.Info("Associating files to operation")
	return operation, nil, operation.AddFiles(exec, false, original, originalTrim, proxy, proxyTrim)
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
		err = original.Update(exec, "name")
		if err != nil {
			return nil, nil, errors.Wrap(err, "Rename original file")
		}
	}

	// Proxy
	proxy, _, err := FindFileBySHA1(exec, r.Proxy.Sha1)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Lookup proxy file")
	}
	if proxy.Name == r.Proxy.FileName {
		log.Info("Proxy's name hasn't change")
	} else {
		log.Info("Renaming proxy")
		proxy.Name = r.Proxy.FileName
		err = proxy.Update(exec, "name")
		if err != nil {
			return nil, nil, errors.Wrap(err, "Rename proxy file")
		}
	}

	log.Info("Processing CIT Metadata")
	evnts, err := ProcessCITMetadata(exec, r.Metadata, original, proxy)
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
	operation, err := CreateOperation(exec, OP_SEND, r.Operation, props)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Create operation")
	}

	log.Info("Associating files to operation")
	err = operation.AddFiles(exec, false, original, proxy)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Associate files")
	}

	return operation, evnts, nil
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
	operation, err := CreateOperation(exec, OP_CONVERT, r.Operation, nil)
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
	props := make(map[string]interface{})
	i := 0
	for _, v := range uniq {
		x := r.Output[v]
		props["duration"] = x.Duration

		// lookup by sha1 as it might be a "reconvert"
		f, _, err := FindFileBySHA1(exec, x.Sha1)
		if err == nil {
			log.Infof("File already exists, updating: %s", x.Sha1)
			err = UpdateFile(exec, f, in, x.File, props)
			if err != nil {
				return nil, nil, errors.Wrap(err, "Update file")
			}
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
	operation, err := CreateOperation(exec, OP_UPLOAD, r.Operation, nil)
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
		file.Properties.Unmarshal(&fileProps)
	}
	fileProps["url"] = r.Url
	fileProps["duration"] = r.Duration
	fpa, _ := json.Marshal(fileProps)
	file.Properties = null.JSONFrom(fpa)

	log.Info("Saving changes to DB")
	err = file.Update(exec, "properties")
	if err != nil {
		return nil, nil, err
	}

	impact, err := PublishFile(exec, file)
	if err != nil {
		return nil, nil, err
	}

	evnts := make([]events.Event, 0)
	evnts = append(evnts, events.FilePublishedEvent(file))
	if impact.ChangedContentUnit != nil {
		evnts = append(evnts, events.ContentUnitPublishedChangeEvent(impact.ChangedContentUnit))
	}
	if impact.ChangedCollections != nil {
		for i := range impact.ChangedCollections {
			evnts = append(evnts, events.CollectionPublishedChangeEvent(impact.ChangedCollections[i]))
		}
	}

	log.Info("Associating file to operation")
	return operation, evnts, operation.AddFiles(exec, false, file)
}

func handleSirtutim(exec boil.Executor, input interface{}) (*models.Operation, []events.Event, error) {
	r := input.(SirtutimRequest)

	log.Info("Creating operation")
	operation, err := CreateOperation(exec, OP_SIRTUTIM, r.Operation, nil)
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
				err = file.Update(exec, "content_unit_id")
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

	log.Infof("Lookup content unit by uid %s", r.ContentUnitUID)
	cu, err := models.ContentUnits(exec,
		qm.Where("uid = ?", r.ContentUnitUID),
		qm.Load("SourceContentUnitDerivations", "SourceContentUnitDerivations.Derived"),
	).One()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Fetch unit from mdb")
	}
	log.Infof("Found content unit %d", cu.ID)

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
	operation, err := CreateOperation(exec, OP_INSERT, r.Operation, props)
	if err != nil {
		return nil, nil, err
	}

	// Create new file
	log.Info("Creating new file")
	if r.File.Type == "" {
		switch r.InsertType {
		case "akladot", "tamlil", "kitei-makor", "article":
			r.File.Type = "text"
		case "sirtutim", "publication":
			r.File.Type = "image"
		case "dgima", "aricha":
			r.File.Type = "video"
		default:
			r.File.Type = ""
		}
	}

	if r.AVFile.Duration > 0 {
		props["duration"] = r.AVFile.Duration
	}

	file, err = CreateFile(exec, parent, r.File, props)
	if err != nil {
		return nil, nil, err
	}
	opFiles = append(opFiles, file)

	// special types logic
	cuID := cu.ID
	if r.InsertType == "kitei-makor" {
		log.Infof("Kitei makor, associating to derived unit")

		var cudID int64
		if len(cu.R.SourceContentUnitDerivations) > 0 {
			for _, cud := range cu.R.SourceContentUnitDerivations {
				if CONTENT_TYPE_REGISTRY.ByID[cud.R.Derived.TypeID].Name == CT_KITEI_MAKOR {
					cudID = cud.DerivedID
					break
				}
			}
		}

		if cudID > 0 {
			cuID = cudID
			log.Infof("KITEI_MAKOR derived unit exist: %d", cuID)
		} else {
			log.Infof("KITEI_MAKOR derived unit doesn't exists. Creating...")
			ktCU, err := CreateContentUnit(exec, CT_KITEI_MAKOR, nil)
			if err != nil {
				return nil, nil, errors.Wrap(err, "Create KITEI_MAKOR derived unit")
			}

			cud := &models.ContentUnitDerivation{
				SourceID: cu.ID,
				Name:     CT_KITEI_MAKOR,
			}
			err = ktCU.AddDerivedContentUnitDerivations(exec, true, cud)
			if err != nil {
				return nil, nil, errors.Wrap(err, "Save CUD in DB")
			}

			cuID = ktCU.ID
		}
	} else if r.InsertType == "publication" {
		log.Infof("Publication, associating to derived unit")

		publisher, err := models.Publishers(exec, qm.Where("uid = ?", r.PublisherUID)).One()
		if err != nil {
			return nil, nil, errors.Wrap(err, "Fetch publisher from mdb")
		}

		var cudID int64
		err = queries.Raw(exec, `SELECT cu.id
FROM content_units cu
  INNER JOIN content_unit_derivations cud ON cu.id = cud.derived_id AND cud.source_id = $1 AND cu.type_id = $2
  INNER JOIN content_units_publishers cup ON cud.derived_id = cup.content_unit_id AND cup.publisher_id = $3`,
			cu.ID, CONTENT_TYPE_REGISTRY.ByName[CT_PUBLICATION].ID, publisher.ID).
			QueryRow().
			Scan(&cudID)
		if err != nil && err != sql.ErrNoRows {
			return nil, nil, errors.Wrap(err, "Lookup existing cud from mdb")
		}

		if cudID > 0 {
			cuID = cudID
			log.Infof("PUBLICATION derived unit exist: %d", cuID)
		} else {
			log.Infof("PUBLICATION derived unit doesn't exists. Creating...")
			pCU, err := CreateContentUnit(exec, CT_PUBLICATION, map[string]interface{}{
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
				Name:     CT_PUBLICATION,
			}
			err = pCU.AddDerivedContentUnitDerivations(exec, true, cud)
			if err != nil {
				return nil, nil, errors.Wrap(err, "Save CUD in DB")
			}

			cuID = pCU.ID
		}
	}

	log.Infof("Associating file [%d] to content_unit [%d]", file.ID, cuID)
	file.ContentUnitID = null.Int64From(cuID)
	err = file.Update(exec, "content_unit_id")
	if err != nil {
		return nil, nil, errors.Wrap(err, "Save file.content_unit_id to DB")
	}

	evnts := make([]events.Event, 0)

	// remove oldFile in update mode and collect events
	if r.Mode == "new" {
		evnts = append(evnts, events.FileInsertEvent(file, r.InsertType))
	} else if r.Mode == "update" {
		evnts = append(evnts, events.FileReplaceEvent(oldFile, file, r.InsertType))

		if impact, err := RemoveFile(exec, oldFile); err != nil {
			return nil, nil, errors.Wrapf(err, "Remove old file %d", oldFile.ID)
		} else {
			if impact.ChangedContentUnit != nil {
				evnts = append(evnts, events.ContentUnitPublishedChangeEvent(impact.ChangedContentUnit))
			}
			if impact.ChangedCollections != nil {
				for i := range impact.ChangedCollections {
					evnts = append(evnts, events.CollectionPublishedChangeEvent(impact.ChangedCollections[i]))
				}
			}
		}

	}

	log.Info("Associating files to operation")
	return operation, evnts, operation.AddFiles(exec, false, opFiles...)
}

func handleTranscode(exec boil.Executor, input interface{}) (*models.Operation, []events.Event, error) {
	r := input.(TranscodeRequest)

	if r.Message != "" {
		log.Infof("Transcode Error: %s", r.Message)

		log.Info("Creating operation")
		opProps := map[string]interface{}{
			"message": r.Message,
		}

		operation, err := CreateOperation(exec, OP_TRANSCODE, r.Operation, opProps)
		if err != nil {
			return nil, nil, err
		}

		log.Info("Looking up original file")
		original, _, err := FindFileBySHA1(exec, r.OriginalSha1)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Lookup original file %s", r.OriginalSha1)
		}

		log.Info("Updating queue table")
		_, err = queries.Raw(exec, "update batch_convert set operation_id=$1 where file_id=$2", operation.ID, original.ID).Exec()
		if err != nil {
			return nil, nil, errors.Wrap(err, "Update queue table")
		}

		return operation, nil, operation.AddFiles(exec, false, original)
	}

	log.Info("Creating operation")
	operation, err := CreateOperation(exec, OP_TRANSCODE, r.Operation, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Info("Looking up original file")
	original, _, err := FindFileBySHA1(exec, r.OriginalSha1)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Lookup original file %s", r.OriginalSha1)
	}

	log.Info("Creating file")
	mt := MEDIA_TYPE_REGISTRY.ByExtension["mp4"]
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
	err = file.Update(exec, "secure", "published")
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Update secure published [%d]", file.ID)
	}

	log.Info("Updating queue table")
	_, err = queries.Raw(exec, "update batch_convert set operation_id=$1 where file_id=$2", operation.ID, original.ID).Exec()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Update queue table")
	}

	log.Info("Associating files to operation")
	return operation, nil, operation.AddFiles(exec, false, original, file)
}

// Helpers

type OperationHandler func(boil.Executor, interface{}) (*models.Operation, []events.Event, error)

// Generic operation handler.
// 	* Manage DB transactions
// 	* Call operation logic handler
// 	* Render JSON response
func handleOperation(c *gin.Context, input interface{}, opHandler OperationHandler) {
	mdb := c.MustGet("MDB").(*sql.DB)
	tx, err := mdb.Begin()
	//tx, err := boil.Begin()
	utils.Must(err)

	_, evnts, err := opHandler(tx, input)
	if err == nil {
		utils.Must(tx.Commit())
	} else {
		utils.Must(tx.Rollback())
		err = errors.Wrapf(err, "Handle operation")
	}

	if err == nil {
		emitter := c.MustGet("EVENTS_EMITTER").(events.EventEmitter)
		emitter.Emit(evnts...)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		switch err.(type) {
		case FileNotFound:
			NewBadRequestError(err).Abort(c)
		default:
			NewInternalError(err).Abort(c)
		}
	}
}
