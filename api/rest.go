package api

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/casbin/casbin"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/volatiletech/null.v6"
	"net/http"
	"strconv"
	"strings"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/events"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/permissions"
	"github.com/Bnei-Baruch/mdb/utils"
)

const (
	DEFAULT_PAGE_SIZE = 50
	MAX_PAGE_SIZE     = 1000

	SEARCH_IN_FILES         = 1
	SEARCH_IN_CONTENT_UNITS = 2
	SEARCH_IN_COLLECTIONS   = 3
)

func CollectionsListHandler(c *gin.Context) {
	var err *HttpError
	var resp interface{}

	switch c.Request.Method {
	case http.MethodGet, "":
		var r CollectionsRequest
		if c.Bind(&r) != nil {
			return
		}

		resp, err = handleCollectionsList(c, c.MustGet("MDB").(*sql.DB), r)
	case http.MethodPost:
		var collection Collection
		if c.BindJSON(&collection) != nil {
			return
		}

		if _, ok := common.CONTENT_TYPE_REGISTRY.ByID[collection.TypeID]; !ok {
			err := errors.Errorf("Unknown content type %d", collection.TypeID)
			NewBadRequestError(err).Abort(c)
			return
		}

		for _, x := range collection.I18n {
			if common.StdLang(x.Language) == common.LANG_UNKNOWN {
				err := errors.Errorf("Unknown language %s", x.Language)
				NewBadRequestError(err).Abort(c)
				return
			}
		}

		tx := mustBeginTx(c)
		resp, err = handleCreateCollection(c, tx, collection)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, events.CollectionCreateEvent(&resp.(*Collection).Collection))
		}
	}

	concludeRequest(c, resp, err)
}

func CollectionHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var err *HttpError
	var resp interface{}

	switch c.Request.Method {
	case http.MethodGet, "":
		resp, err = handleGetCollection(c, c.MustGet("MDB").(*sql.DB), id)
	case http.MethodPut:
		var cl PartialCollection
		if c.Bind(&cl) != nil {
			return
		}

		cl.ID = id
		tx := mustBeginTx(c)
		resp, err = handleUpdateCollection(c, tx, &cl)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, events.CollectionUpdateEvent(&resp.(*Collection).Collection))
		}
	case http.MethodDelete:
		tx := mustBeginTx(c)
		var cl *models.Collection
		cl, err = handleDeleteCollection(c, tx, id)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, events.CollectionDeleteEvent(cl))
		}
	}

	concludeRequest(c, resp, err)
}

func CollectionI18nHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var i18ns []*models.CollectionI18n
	if c.Bind(&i18ns) != nil {
		return
	}
	for _, x := range i18ns {
		if common.StdLang(x.Language) == common.LANG_UNKNOWN {
			NewBadRequestError(errors.Errorf("Unknown language %s", x.Language)).Abort(c)
			return
		}
	}

	tx := mustBeginTx(c)
	resp, err := handleUpdateCollectionI18n(c, tx, id, i18ns)
	mustConcludeTx(tx, err)

	if err == nil {
		emitEvents(c, events.CollectionUpdateEvent(&resp.Collection))
	}

	concludeRequest(c, resp, err)
}

func CollectionContentUnitsHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var err *HttpError
	var resp interface{}

	switch c.Request.Method {
	case http.MethodGet, "":
		resp, err = handleCollectionCCU(c, c.MustGet("MDB").(*sql.DB), id)
	case http.MethodPost:
		var ccus []*models.CollectionsContentUnit
		if c.BindJSON(&ccus) != nil {
			return
		}

		var evnts []events.Event
		tx := mustBeginTx(c)
		evnts, err = handleCollectionAddCCU(c, tx, id, ccus)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, evnts...)
		}
	case http.MethodPut:
		var ccu models.CollectionsContentUnit
		if c.BindJSON(&ccu) != nil {
			return
		}

		cuID, e := strconv.ParseInt(c.Param("cuID"), 10, 0)
		if e != nil {
			err = NewBadRequestError(errors.Wrap(e, "cuID expects int64"))
			break
		}
		ccu.ContentUnitID = cuID

		var event *events.Event
		tx := mustBeginTx(c)
		event, err = handleCollectionUpdateCCU(c, tx, id, ccu)
		mustConcludeTx(tx, err)

		if err == nil && event != nil {
			emitEvents(c, *event)
		}
	case http.MethodDelete:
		cuID, e := strconv.ParseInt(c.Param("cuID"), 10, 0)
		if e != nil {
			err = NewBadRequestError(errors.Wrap(e, "cuID expects int64"))
			break
		}

		var evnts []events.Event
		tx := mustBeginTx(c)
		evnts, err = handleCollectionRemoveCCU(c, tx, id, cuID)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, evnts...)
		}
	}

	concludeRequest(c, resp, err)
}

// Toggle the active flag of a single container
func CollectionActivateHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleCollectionActivate(c, c.MustGet("MDB").(*sql.DB), id)
	concludeRequest(c, resp, err)
}

func ContentUnitsListHandler(c *gin.Context) {
	var err *HttpError
	var resp interface{}

	switch c.Request.Method {
	case http.MethodGet, "":
		var r ContentUnitsRequest
		if c.Bind(&r) != nil {
			return
		}

		resp, err = handleContentUnitsList(c, c.MustGet("MDB").(*sql.DB), r)
	case http.MethodPost:
		var unit ContentUnit
		if c.BindJSON(&unit) != nil {
			return
		}

		if _, ok := common.CONTENT_TYPE_REGISTRY.ByID[unit.TypeID]; !ok {
			err := errors.Errorf("Unknown content type %d", unit.TypeID)
			NewBadRequestError(err).Abort(c)
			return
		}

		for _, x := range unit.I18n {
			if common.StdLang(x.Language) == common.LANG_UNKNOWN {
				err := errors.Errorf("Unknown language %s", x.Language)
				NewBadRequestError(err).Abort(c)
				return
			}
		}

		tx := mustBeginTx(c)
		resp, err = handleCreateContentUnit(c, tx, unit)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, events.ContentUnitCreateEvent(&resp.(*ContentUnit).ContentUnit))
		}
	}

	concludeRequest(c, resp, err)
}

func ContentUnitHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var err *HttpError
	var resp interface{}

	if c.Request.Method == http.MethodGet || c.Request.Method == "" {
		resp, err = handleGetContentUnit(c, c.MustGet("MDB").(*sql.DB), id)
	} else {
		if c.Request.Method == http.MethodPut {
			var cu PartialContentUnit
			if c.Bind(&cu) != nil {
				return
			}

			cu.ID = id
			tx := mustBeginTx(c)
			resp, err = handleUpdateContentUnit(c, tx, &cu)
			mustConcludeTx(tx, err)

			if err == nil {
				emitEvents(c, events.ContentUnitUpdateEvent(&resp.(*ContentUnit).ContentUnit))
			}
		}
	}

	concludeRequest(c, resp, err)
}

func ContentUnitI18nHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var i18ns []*models.ContentUnitI18n
	if c.Bind(&i18ns) != nil {
		return
	}
	for _, x := range i18ns {
		if common.StdLang(x.Language) == common.LANG_UNKNOWN {
			NewBadRequestError(errors.Errorf("Unknown language %s", x.Language)).Abort(c)
			return
		}
	}

	tx := mustBeginTx(c)
	resp, err := handleUpdateContentUnitI18n(c, tx, id, i18ns)
	mustConcludeTx(tx, err)

	if err == nil {
		emitEvents(c, events.ContentUnitUpdateEvent(&resp.ContentUnit))
	}

	concludeRequest(c, resp, err)
}

func ContentUnitFilesHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var err *HttpError
	var resp interface{}

	if c.Request.Method == http.MethodGet || c.Request.Method == "" {
		resp, err = handleContentUnitFiles(c, c.MustGet("MDB").(*sql.DB), id)
	} else {
		if c.Request.Method == http.MethodPost {
			var fids []int64
			if c.Bind(&fids) != nil {
				return
			}

			var evnts []events.Event
			tx := mustBeginTx(c)
			resp, evnts, err = handleContentUnitAddFiles(c, tx, id, fids)
			mustConcludeTx(tx, err)

			if err == nil {
				emitEvents(c, evnts...)
			}
		}
	}

	concludeRequest(c, resp, err)
}

func ContentUnitCollectionsHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleContentUnitCCU(c, c.MustGet("MDB").(*sql.DB), id)
	concludeRequest(c, resp, err)
}

func ContentUnitDerivativesHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var err *HttpError
	var resp interface{}

	switch c.Request.Method {
	case http.MethodGet, "":
		resp, err = handleContentUnitCUD(c, c.MustGet("MDB").(*sql.DB), id)
	case http.MethodPost:
		var cud models.ContentUnitDerivation
		if c.BindJSON(&cud) != nil {
			return
		}

		tx := mustBeginTx(c)
		resp, err = handleContentUnitAddCUD(c, tx, id, cud)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, events.ContentUnitDerivativesChangeEvent(resp.(*models.ContentUnit)))
		}
	case http.MethodPut:
		var cud models.ContentUnitDerivation
		if c.BindJSON(&cud) != nil {
			return
		}

		duID, e := strconv.ParseInt(c.Param("duID"), 10, 0)
		if e != nil {
			err = NewBadRequestError(errors.Wrap(e, "duID expects int64"))
			break
		}
		cud.DerivedID = duID

		tx := mustBeginTx(c)
		resp, err = handleContentUnitUpdateCUD(c, tx, id, cud)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, events.ContentUnitDerivativesChangeEvent(resp.(*models.ContentUnit)))
		}
	case http.MethodDelete:
		duID, e := strconv.ParseInt(c.Param("duID"), 10, 0)
		if e != nil {
			err = NewBadRequestError(errors.Wrap(e, "duID expects int64"))
			break
		}

		tx := mustBeginTx(c)
		resp, err = handleContentUnitRemoveCUD(c, tx, id, duID)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, events.ContentUnitDerivativesChangeEvent(resp.(*models.ContentUnit)))
		}
	}

	concludeRequest(c, resp, err)
}

func ContentUnitOriginsHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleContentUnitOrigins(c, c.MustGet("MDB").(*sql.DB), id)
	concludeRequest(c, resp, err)
}

func ContentUnitSourcesHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	switch c.Request.Method {
	case http.MethodGet, "":
		resp, err := handleGetContentUnitSources(c, c.MustGet("MDB").(*sql.DB), id)
		concludeRequest(c, resp, err)
	case http.MethodPost:
		var body map[string]int64
		if c.BindJSON(&body) != nil {
			return
		}

		sourceID, ok := body["sourceID"]
		if !ok {
			NewBadRequestError(errors.Wrap(e, "No sourceID given")).Abort(c)
			return
		}

		tx := mustBeginTx(c)
		resp, err := handleContentUnitAddSource(c, tx, id, sourceID)
		mustConcludeTx(tx, err)

		if err == nil && resp != nil {
			emitEvents(c, events.ContentUnitSourcesChangeEvent(resp))
		}

		concludeRequest(c, resp, err)
	case http.MethodDelete:
		sourceID, e := strconv.ParseInt(c.Param("sourceID"), 10, 0)
		if e != nil {
			NewBadRequestError(errors.Wrap(e, "sourceID expects int64")).Abort(c)
			return
		}

		tx := mustBeginTx(c)
		resp, err := handleContentUnitRemoveSource(c, tx, id, sourceID)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, events.ContentUnitSourcesChangeEvent(resp))
		}
		concludeRequest(c, resp, err)
	}
}

func ContentUnitTagsHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	switch c.Request.Method {
	case http.MethodGet, "":
		resp, err := handleGetContentUnitTags(c, c.MustGet("MDB").(*sql.DB), id)
		concludeRequest(c, resp, err)
	case http.MethodPost:
		var body map[string]int64
		if c.BindJSON(&body) != nil {
			return
		}

		tagID, ok := body["tagID"]
		if !ok {
			NewBadRequestError(errors.Wrap(e, "No tagID given")).Abort(c)
			return
		}

		tx := mustBeginTx(c)
		resp, err := handleContentUnitAddTag(c, tx, id, tagID)
		mustConcludeTx(tx, err)

		if err == nil && resp != nil {
			emitEvents(c, events.ContentUnitTagsChangeEvent(resp))
		}

		concludeRequest(c, resp, err)
	case http.MethodDelete:
		tagID, e := strconv.ParseInt(c.Param("tagID"), 10, 0)
		if e != nil {
			NewBadRequestError(errors.Wrap(e, "tagID expects int64")).Abort(c)
			return
		}

		tx := mustBeginTx(c)
		resp, err := handleContentUnitRemoveTag(c, tx, id, tagID)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, events.ContentUnitTagsChangeEvent(resp))
		}

		concludeRequest(c, resp, err)
	}
}

func ContentUnitPersonsHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	switch c.Request.Method {
	case http.MethodGet, "":
		resp, err := handleGetContentUnitPersons(c, c.MustGet("MDB").(*sql.DB), id)
		concludeRequest(c, resp, err)
	case http.MethodPost:
		var cup models.ContentUnitsPerson
		if c.BindJSON(&cup) != nil {
			return
		}

		tx := mustBeginTx(c)
		resp, err := handleContentUnitAddPerson(c, tx, id, cup)
		mustConcludeTx(tx, err)

		if err == nil && resp != nil {
			emitEvents(c, events.ContentUnitPersonsChangeEvent(resp))
		}

		concludeRequest(c, resp, err)
	case http.MethodDelete:
		personID, e := strconv.ParseInt(c.Param("personID"), 10, 0)
		if e != nil {
			NewBadRequestError(errors.Wrap(e, "personID expects int64")).Abort(c)
			return
		}

		tx := mustBeginTx(c)
		resp, err := handleContentUnitRemovePerson(c, tx, id, personID)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, events.ContentUnitPersonsChangeEvent(resp))
		}

		concludeRequest(c, resp, err)
	}
}

func ContentUnitPublishersHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var err *HttpError
	var resp interface{}

	switch c.Request.Method {
	case http.MethodGet, "":
		resp, err = handleGetContentUnitPublishers(c, c.MustGet("MDB").(*sql.DB), id)
	case http.MethodPost:
		var body map[string]int64
		if c.BindJSON(&body) != nil {
			return
		}

		publisherID, ok := body["publisherID"]
		if !ok {
			err = NewBadRequestError(errors.Wrap(e, "No publisherID given"))
			break
		}

		tx := mustBeginTx(c)
		resp, err = handleContentUnitAddPublisher(c, tx, id, publisherID)
		mustConcludeTx(tx, err)

		if respCU, ok := resp.(*models.ContentUnit); ok && err == nil {
			emitEvents(c, events.ContentUnitPublishersChangeEvent(respCU))
		}
	case http.MethodDelete:
		publisherID, e := strconv.ParseInt(c.Param("publisherID"), 10, 0)
		if e != nil {
			err = NewBadRequestError(errors.Wrap(e, "publisherID expects int64"))
			break
		}

		tx := mustBeginTx(c)
		resp, err = handleContentUnitRemovePublisher(c, tx, id, publisherID)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, events.ContentUnitPublishersChangeEvent(resp.(*models.ContentUnit)))
		}
	}

	concludeRequest(c, resp, err)
}

func ContentUnitMergeHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var cuIDs []int64
	if c.Bind(&cuIDs) != nil {
		return
	}

	// filter out host unit ID
	b := cuIDs[:0]
	for _, x := range cuIDs {
		if x != id {
			b = append(b, x)
		}
	}

	tx := mustBeginTx(c)
	resp, evnts, err := handleContentUnitMerge(c, tx, id, b)
	mustConcludeTx(tx, err)

	if err == nil {
		emitEvents(c, evnts...)
	}

	concludeRequest(c, resp, err)
}

func FilesListHandler(c *gin.Context) {
	var r FilesRequest
	if c.Bind(&r) != nil {
		return
	}

	resp, err := handleFilesList(c, c.MustGet("MDB").(*sql.DB), r)
	concludeRequest(c, resp, err)
}

func FileHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var err *HttpError
	var resp interface{}

	if c.Request.Method == http.MethodGet || c.Request.Method == "" {
		resp, err = handleGetFile(c, c.MustGet("MDB").(*sql.DB), id)
	} else {
		if c.Request.Method == http.MethodPut {
			var f PartialFile
			if c.Bind(&f) != nil {
				return
			}

			f.ID = id
			var evnts []events.Event
			tx := mustBeginTx(c)
			resp, evnts, err = handleUpdateFile(c, tx, &f)
			mustConcludeTx(tx, err)

			if err == nil {
				emitEvents(c, evnts...)
			}
		}
	}

	concludeRequest(c, resp, err)
}

func FileStoragesHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleFileStorages(c, c.MustGet("MDB").(*sql.DB), id)
	concludeRequest(c, resp, err)
}

func FilesWithOperationsTreeHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	db := c.MustGet("MDB").(*sql.DB)

	files, err := FindFileTreeWithOperations(db, id)
	if err != nil {
		NewInternalError(err).Abort(c)
		return
	}

	opIdsMap := make(map[int64]bool)
	for i := range files {
		for j := range files[i].OperationIds {
			opIdsMap[files[i].OperationIds[j]] = true
		}
	}

	opIds := make([]int64, len(opIdsMap))
	i := 0
	for k := range opIdsMap {
		opIds[i] = k
		i++
	}

	ops, err := models.Operations(db,
		qm.WhereIn("id in ?", utils.ConvertArgsInt64(opIds)...)).
		All()
	if err != nil {
		NewInternalError(err).Abort(c)
		return
	}

	opsMap := make(map[int64]*models.Operation)
	for _, op := range ops {
		opsMap[op.ID] = op
	}

	resp := &struct {
		Files      []*MFile                    `json:"files"`
		Operations map[int64]*models.Operation `json:"operations"`
	}{
		files,
		opsMap,
	}

	concludeRequest(c, resp, nil)
}

func OperationsListHandler(c *gin.Context) {
	// check permissions
	if !can(c, secureToPermission(common.SEC_PUBLIC), common.PERM_READ) {
		NewForbiddenError().Abort(c)
		return
	}

	var r OperationsRequest
	if c.Bind(&r) != nil {
		return
	}

	resp, err := handleOperationsList(c.MustGet("MDB").(*sql.DB), r)
	concludeRequest(c, resp, err)
}

func OperationItemHandler(c *gin.Context) {
	// check permissions
	if !can(c, secureToPermission(common.SEC_PUBLIC), common.PERM_READ) {
		NewForbiddenError().Abort(c)
		return
	}

	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleOperationItem(c.MustGet("MDB").(*sql.DB), id)
	concludeRequest(c, resp, err)
}

func OperationFilesHandler(c *gin.Context) {
	// check permissions
	if !can(c, secureToPermission(common.SEC_PUBLIC), common.PERM_READ) {
		NewForbiddenError().Abort(c)
		return
	}

	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleOperationFiles(c, c.MustGet("MDB").(*sql.DB), id)
	concludeRequest(c, resp, err)
}

func AuthorsHandler(c *gin.Context) {
	if !can(c, secureToPermission(common.SEC_PUBLIC), common.PERM_READ) {
		NewForbiddenError().Abort(c)
		return
	}

	authors, err := models.Authors(c.MustGet("MDB").(*sql.DB),
		qm.Load("AuthorI18ns", "Sources")).
		All()
	if err != nil {
		NewInternalError(errors.Wrap(err, "Load authors from DB")).Abort(c)
		return
	}

	data := make([]*Author, len(authors))

	for i, a := range authors {
		x := &Author{Author: *a}
		data[i] = x

		x.I18n = make(map[string]*models.AuthorI18n, len(a.R.AuthorI18ns))
		for _, i18n := range a.R.AuthorI18ns {
			x.I18n[i18n.Language] = i18n
		}

		x.Sources = make([]*Source, len(a.R.Sources))
		for j, s := range a.R.Sources {
			x.Sources[j] = &Source{Source: *s}
		}
	}

	resp := AuthorsResponse{
		ListResponse: ListResponse{Total: int64(len(data))},
		Authors:      data,
	}

	concludeRequest(c, resp, nil)
}

func SourcesHandler(c *gin.Context) {
	var err *HttpError
	var resp interface{}

	if c.Request.Method == http.MethodGet || c.Request.Method == "" {
		if !can(c, secureToPermission(common.SEC_PUBLIC), common.PERM_READ) {
			NewForbiddenError().Abort(c)
			return
		}

		var r SourcesRequest
		if c.Bind(&r) != nil {
			return
		}
		resp, err = handleGetSources(c.MustGet("MDB").(*sql.DB), r)
	} else {
		if c.Request.Method == http.MethodPost {
			if !isAdmin(c) {
				NewForbiddenError().Abort(c)
				return
			}

			var r CreateSourceRequest
			if c.Bind(&r) != nil {
				return
			}

			if _, ok := common.SOURCE_TYPE_REGISTRY.ByID[r.Source.TypeID]; !ok {
				err := errors.Errorf("Unknown source type %d", r.Source.TypeID)
				NewBadRequestError(err).Abort(c)
				return
			}

			for _, x := range r.Source.I18n {
				if common.StdLang(x.Language) == common.LANG_UNKNOWN {
					err := errors.Errorf("Unknown language %s", x.Language)
					NewBadRequestError(err).Abort(c)
					return
				}
			}

			tx := mustBeginTx(c)
			resp, err = handleCreateSource(tx, r)
			mustConcludeTx(tx, err)

			if err == nil {
				emitEvents(c, events.SourceCreateEvent(&resp.(*Source).Source))
			}
		}
	}

	concludeRequest(c, resp, err)
}

func SourceHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var err *HttpError
	var resp interface{}

	if c.Request.Method == http.MethodGet || c.Request.Method == "" {
		if !can(c, secureToPermission(common.SEC_PUBLIC), common.PERM_READ) {
			NewForbiddenError().Abort(c)
			return
		}

		resp, err = handleGetSource(c.MustGet("MDB").(*sql.DB), id)
	} else {
		if c.Request.Method == http.MethodPut {
			if !isAdmin(c) {
				NewForbiddenError().Abort(c)
				return
			}

			var s Source
			if c.Bind(&s) != nil {
				return
			}

			s.ID = id
			tx := mustBeginTx(c)
			resp, err = handleUpdateSource(tx, &s)
			mustConcludeTx(tx, err)

			if err == nil {
				emitEvents(c, events.SourceUpdateEvent(&resp.(*Source).Source))
			}
		}
	}

	concludeRequest(c, resp, err)
}

func SourceI18nHandler(c *gin.Context) {
	if !isAdmin(c) {
		NewForbiddenError().Abort(c)
		return
	}

	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var i18ns []*models.SourceI18n
	if c.Bind(&i18ns) != nil {
		return
	}
	for _, x := range i18ns {
		if common.StdLang(x.Language) == common.LANG_UNKNOWN {
			NewBadRequestError(errors.Errorf("Unknown language %s", x.Language)).Abort(c)
			return
		}
	}

	tx := mustBeginTx(c)
	resp, err := handleUpdateSourceI18n(tx, id, i18ns)
	mustConcludeTx(tx, err)

	if err == nil {
		emitEvents(c, events.SourceUpdateEvent(&resp.Source))
	}

	concludeRequest(c, resp, err)
}

func TagsHandler(c *gin.Context) {
	var err *HttpError
	var resp interface{}

	if c.Request.Method == http.MethodGet || c.Request.Method == "" {
		if !can(c, secureToPermission(common.SEC_PUBLIC), common.PERM_READ) {
			NewForbiddenError().Abort(c)
			return
		}

		var r TagsRequest
		if c.Bind(&r) != nil {
			return
		}
		resp, err = handleGetTags(c.MustGet("MDB").(*sql.DB), r)
	} else {
		if c.Request.Method == http.MethodPost {
			if !isAdmin(c) {
				NewForbiddenError().Abort(c)
				return
			}

			var t Tag
			if c.Bind(&t) != nil {
				return
			}

			for _, x := range t.I18n {
				if common.StdLang(x.Language) == common.LANG_UNKNOWN {
					NewBadRequestError(errors.Errorf("Unknown language %s", x.Language)).Abort(c)
					return
				}
			}

			tx := mustBeginTx(c)
			resp, err = handleCreateTag(tx, &t)
			mustConcludeTx(tx, err)

			if err == nil {
				emitEvents(c, events.TagCreateEvent(&resp.(*Tag).Tag))
			}
		}
	}

	concludeRequest(c, resp, err)
}

func TagHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var err *HttpError
	var resp interface{}

	if c.Request.Method == http.MethodGet || c.Request.Method == "" {
		if !can(c, secureToPermission(common.SEC_PUBLIC), common.PERM_READ) {
			NewForbiddenError().Abort(c)
			return
		}

		resp, err = handleGetTag(c.MustGet("MDB").(*sql.DB), id)
	} else {
		if c.Request.Method == http.MethodPut {
			if !isAdmin(c) {
				NewForbiddenError().Abort(c)
				return
			}

			var t Tag
			if c.Bind(&t) != nil {
				return
			}

			t.ID = id
			tx := mustBeginTx(c)
			resp, err = handleUpdateTag(tx, &t)
			mustConcludeTx(tx, err)

			if err == nil {
				emitEvents(c, events.TagUpdateEvent(&resp.(*Tag).Tag))
			}
		}
	}

	concludeRequest(c, resp, err)
}

func TagI18nHandler(c *gin.Context) {
	if !isAdmin(c) {
		NewForbiddenError().Abort(c)
		return
	}

	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var i18ns []*models.TagI18n
	if c.Bind(&i18ns) != nil {
		return
	}
	for _, x := range i18ns {
		if common.StdLang(x.Language) == common.LANG_UNKNOWN {
			NewBadRequestError(errors.Errorf("Unknown language %s", x.Language)).Abort(c)
			return
		}
	}

	tx := mustBeginTx(c)
	resp, err := handleUpdateTagI18n(tx, id, i18ns)
	mustConcludeTx(tx, err)

	if err == nil {
		emitEvents(c, events.TagUpdateEvent(&resp.Tag))
	}

	concludeRequest(c, resp, err)
}

func PersonsListHandler(c *gin.Context) {
	var err *HttpError
	var resp interface{}

	switch c.Request.Method {
	case http.MethodGet, "":
		if !can(c, secureToPermission(common.SEC_PUBLIC), common.PERM_READ) {
			NewForbiddenError().Abort(c)
			return
		}

		var r PersonsRequest
		if c.Bind(&r) != nil {
			return
		}

		resp, err = handlePersonsList(c.MustGet("MDB").(*sql.DB), r)
	case http.MethodPost:
		if !isAdmin(c) {
			NewForbiddenError().Abort(c)
			return
		}

		var person Person
		if c.BindJSON(&person) != nil {
			return
		}

		for _, x := range person.I18n {
			if common.StdLang(x.Language) == common.LANG_UNKNOWN {
				err := errors.Errorf("Unknown language %s", x.Language)
				NewBadRequestError(err).Abort(c)
				return
			}
		}

		tx := mustBeginTx(c)
		resp, err = handleCreatePerson(tx, &person)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, events.PersonCreateEvent(&resp.(*Person).Person))
		}
	}

	concludeRequest(c, resp, err)
}

func PersonHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var err *HttpError
	var resp interface{}

	switch c.Request.Method {
	case http.MethodGet, "":
		if !can(c, secureToPermission(common.SEC_PUBLIC), common.PERM_READ) {
			NewForbiddenError().Abort(c)
			return
		}

		resp, err = handleGetPerson(c.MustGet("MDB").(*sql.DB), id)
	case http.MethodPut:
		if !isAdmin(c) {
			NewForbiddenError().Abort(c)
			return
		}

		var p Person
		if c.Bind(&p) != nil {
			return
		}

		p.ID = id
		tx := mustBeginTx(c)
		resp, err = handleUpdatePerson(tx, &p)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, events.PersonUpdateEvent(&resp.(*Person).Person))
		}
	case http.MethodDelete:
		if !isAdmin(c) {
			NewForbiddenError().Abort(c)
			return
		}

		tx := mustBeginTx(c)
		pr, err := handleDeletePerson(tx, id)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, events.PersonDeleteEvent(pr))
		}
	}

	concludeRequest(c, resp, err)
}

func PersonI18nHandler(c *gin.Context) {
	if !isAdmin(c) {
		NewForbiddenError().Abort(c)
		return
	}

	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var i18ns []*models.PersonI18n
	if c.Bind(&i18ns) != nil {
		return
	}
	for _, x := range i18ns {
		if common.StdLang(x.Language) == common.LANG_UNKNOWN {
			NewBadRequestError(errors.Errorf("Unknown language %s", x.Language)).Abort(c)
			return
		}
	}

	tx := mustBeginTx(c)
	resp, err := handleUpdatePersonI18n(tx, id, i18ns)
	mustConcludeTx(tx, err)

	if err == nil {
		emitEvents(c, events.PersonUpdateEvent(&resp.Person))
	}

	concludeRequest(c, resp, err)
}

func PublishersHandler(c *gin.Context) {
	var err *HttpError
	var resp interface{}

	switch c.Request.Method {
	case http.MethodGet, "":
		if !can(c, secureToPermission(common.SEC_PUBLIC), common.PERM_READ) {
			NewForbiddenError().Abort(c)
			return
		}

		var r PublishersRequest
		if c.Bind(&r) != nil {
			return
		}

		resp, err = handlePublishersList(c.MustGet("MDB").(*sql.DB), r)
	case http.MethodPost:
		if !isAdmin(c) {
			NewForbiddenError().Abort(c)
			return
		}

		var publisher Publisher
		if c.BindJSON(&publisher) != nil {
			return
		}

		for _, x := range publisher.I18n {
			if common.StdLang(x.Language) == common.LANG_UNKNOWN {
				err := errors.Errorf("Unknown language %s", x.Language)
				NewBadRequestError(err).Abort(c)
				return
			}
		}

		tx := mustBeginTx(c)
		resp, err = handleCreatePublisher(tx, &publisher)
		mustConcludeTx(tx, err)

		if err == nil {
			emitEvents(c, events.PublisherCreateEvent(&resp.(*Publisher).Publisher))
		}
	}

	concludeRequest(c, resp, err)
}

func PublisherHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var err *HttpError
	var resp interface{}

	if c.Request.Method == http.MethodGet || c.Request.Method == "" {
		if !can(c, secureToPermission(common.SEC_PUBLIC), common.PERM_READ) {
			NewForbiddenError().Abort(c)
			return
		}

		resp, err = handleGetPublisher(c.MustGet("MDB").(*sql.DB), id)
	} else {
		if c.Request.Method == http.MethodPut {
			if !isAdmin(c) {
				NewForbiddenError().Abort(c)
				return
			}

			var p Publisher
			if c.Bind(&p) != nil {
				return
			}

			p.ID = id
			tx := mustBeginTx(c)
			resp, err = handleUpdatePublisher(tx, &p)
			mustConcludeTx(tx, err)

			if err == nil {
				emitEvents(c, events.PublisherUpdateEvent(&resp.(*Publisher).Publisher))
			}
		}
	}

	concludeRequest(c, resp, err)
}

func PublisherI18nHandler(c *gin.Context) {
	if !isAdmin(c) {
		NewForbiddenError().Abort(c)
		return
	}

	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var i18ns []*models.PublisherI18n
	if c.Bind(&i18ns) != nil {
		return
	}
	for _, x := range i18ns {
		if common.StdLang(x.Language) == common.LANG_UNKNOWN {
			NewBadRequestError(errors.Errorf("Unknown language %s", x.Language)).Abort(c)
			return
		}
	}

	tx := mustBeginTx(c)
	resp, err := handleUpdatePublisherI18n(tx, id, i18ns)
	mustConcludeTx(tx, err)

	if err == nil {
		emitEvents(c, events.PublisherUpdateEvent(&resp.Publisher))
	}

	concludeRequest(c, resp, err)
}

func StoragesHandler(c *gin.Context) {
	if !can(c, secureToPermission(common.SEC_PUBLIC), common.PERM_READ) {
		NewForbiddenError().Abort(c)
		return
	}

	var r StoragesRequest
	if c.Bind(&r) != nil {
		return
	}

	resp, err := handleStoragesList(c.MustGet("MDB").(*sql.DB), r)
	concludeRequest(c, resp, err)
}

// Handlers Logic

func handleCollectionsList(cp utils.ContextProvider, exec boil.Executor, r CollectionsRequest) (*CollectionsResponse, *HttpError) {
	mods := make([]qm.QueryMod, 0)
	appendPermissionsMods(cp, &mods)

	// filters
	if err := appendIDsFilterMods(&mods, r.IDsFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendUIDsFilterMods(&mods, r.UIDsFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendContentTypesFilterMods(&mods, r.ContentTypesFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendDateRangeFilterMods(&mods, r.DateRangeFilter,
		"(coalesce(properties->>'film_date', properties->>'start_date', created_at::text))::date"); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendSecureFilterMods(&mods, r.SecureFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendSearchTermFilterMods(exec, &mods, r.SearchTermFilter, SEARCH_IN_COLLECTIONS); err != nil {
		return nil, NewBadRequestError(err)
	}

	appendPublishedFilterMods(&mods, r.PublishedFilter)

	// count query
	var total int64
	countMods := append([]qm.QueryMod{qm.Select("count(DISTINCT id)")}, mods...)
	err := models.Collections(exec, countMods...).QueryRow().Scan(&total)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if total == 0 {
		return NewCollectionsResponse(), nil
	}

	// order, limit, offset
	if err = appendListMods(&mods, r.ListRequest); err != nil {
		return nil, NewBadRequestError(err)
	}

	// Eager loading
	mods = append(mods, qm.Load("CollectionI18ns"))

	// data query
	collections, err := models.Collections(exec, mods...).All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	// i18n
	data := make([]*Collection, len(collections))
	for i, c := range collections {
		x := &Collection{Collection: *c}
		data[i] = x
		x.I18n = make(map[string]*models.CollectionI18n, len(c.R.CollectionI18ns))
		for _, i18n := range c.R.CollectionI18ns {
			x.I18n[i18n.Language] = i18n
		}
	}

	return &CollectionsResponse{
		ListResponse: ListResponse{Total: total},
		Collections:  data,
	}, nil
}

func handleCreateCollection(cp utils.ContextProvider, exec boil.Executor, c Collection) (*Collection, *HttpError) {
	// check object level permissions
	if !can(cp, secureToPermission(c.Secure), common.PERM_WRITE) {
		return nil, NewForbiddenError()
	}

	// unmarshal properties
	props := make(map[string]interface{})
	if c.Properties.Valid {
		err := json.Unmarshal(c.Properties.JSON, &props)
		if err != nil {
			return nil, NewBadRequestError(errors.Wrap(err, "json.Unmarshal properties"))
		}
	}

	// create collection in DB
	ct := common.CONTENT_TYPE_REGISTRY.ByID[c.TypeID].Name
	collection, err := CreateCollection(exec, ct, props)
	if err != nil {
		return nil, NewInternalError(err)
	}

	// save i18n
	for _, v := range c.I18n {
		err := collection.AddCollectionI18ns(exec, true, v)
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	return handleGetCollection(cp, exec, collection.ID)
}

func handleGetCollection(cp utils.ContextProvider, exec boil.Executor, id int64) (*Collection, *HttpError) {
	collection, err := models.Collections(exec,
		qm.Where("id = ?", id),
		qm.Load("CollectionI18ns")).
		One()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(collection.Secure), common.PERM_READ) {
		return nil, NewForbiddenError()
	}

	// i18n
	x := &Collection{Collection: *collection}
	x.I18n = make(map[string]*models.CollectionI18n, len(collection.R.CollectionI18ns))
	for _, i18n := range collection.R.CollectionI18ns {
		x.I18n[i18n.Language] = i18n
	}

	return x, nil
}

func handleUpdateCollection(cp utils.ContextProvider, exec boil.Executor, c *PartialCollection) (*Collection, *HttpError) {
	collection, err := models.FindCollection(exec, c.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(collection.Secure), common.PERM_WRITE) {
		return nil, NewForbiddenError()
	}

	// update entity attributes
	if c.Secure.Valid {
		collection.Secure = c.Secure.Int16
		err = collection.Update(exec, "secure")
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	// update properties bag
	if c.Properties.Valid {
		var props map[string]interface{}
		err = c.Properties.Unmarshal(&props)
		if err != nil {
			return nil, NewInternalError(err)
		}

		err = UpdateCollectionProperties(exec, collection, props)
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	return handleGetCollection(cp, exec, c.ID)
}

func handleDeleteCollection(cp utils.ContextProvider, exec boil.Executor, id int64) (*models.Collection, *HttpError) {
	collection, err := models.FindCollection(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(collection.Secure), common.PERM_WRITE) {
		return nil, NewForbiddenError()
	}

	err = models.CollectionsContentUnits(exec, qm.Where("collection_id = ?", id)).DeleteAll()
	if err != nil {
		return nil, NewInternalError(err)
	}

	err = models.CollectionI18ns(exec, qm.Where("collection_id = ?", id)).DeleteAll()
	if err != nil {
		return nil, NewInternalError(err)
	}

	err = collection.Delete(exec)
	if err != nil {
		return nil, NewInternalError(err)
	}

	return collection, nil
}

func handleUpdateCollectionI18n(cp utils.ContextProvider, exec boil.Executor, id int64, i18ns []*models.CollectionI18n) (*Collection, *HttpError) {
	collection, err := handleGetCollection(cp, exec, id)
	if err != nil {
		return nil, err
	}

	// check object level permissions
	if !can(cp, secureToPermission(collection.Secure), common.PERM_I18N_WRITE) {
		return nil, NewForbiddenError()
	}

	// Upsert all new i18ns
	nI18n := make(map[string]*models.CollectionI18n, len(i18ns))
	for _, i18n := range i18ns {
		i18n.CollectionID = id
		nI18n[i18n.Language] = i18n
		err := i18n.Upsert(exec, true,
			[]string{"collection_id", "language"},
			[]string{"name", "description"})
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	// Delete old i18ns not in new i18ns
	for k, v := range collection.I18n {
		if _, ok := nI18n[k]; !ok {
			err := v.Delete(exec)
			if err != nil {
				return nil, NewInternalError(err)
			}
		}
	}

	return handleGetCollection(cp, exec, id)
}

func handleCollectionActivate(cp utils.ContextProvider, exec boil.Executor, id int64) (*Collection, *HttpError) {
	collection, err := models.FindCollection(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(collection.Secure), common.PERM_WRITE) {
		return nil, NewForbiddenError()
	}

	var props = make(map[string]interface{})
	if collection.Properties.Valid {
		collection.Properties.Unmarshal(&props)
	}
	active, ok := props["active"]
	if ok {
		b, _ := active.(bool)
		props["active"] = !b
	} else {
		props["active"] = false
	}

	pbytes, err := json.Marshal(props)
	if err != nil {
		return nil, NewInternalError(err)
	}
	collection.Properties = null.JSONFrom(pbytes)
	err = collection.Update(exec, "properties")
	if err != nil {
		return nil, NewInternalError(err)
	}

	return handleGetCollection(cp, exec, id)
}

func handleCollectionCCU(cp utils.ContextProvider, exec boil.Executor, id int64) ([]*CollectionContentUnit, *HttpError) {
	collection, err := models.FindCollection(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(collection.Secure), common.PERM_READ) {
		return nil, NewForbiddenError()
	}

	ccus, err := models.CollectionsContentUnits(exec,
		qm.Where("collection_id = ?", id),
		qm.OrderBy("position")).
		All()
	if err != nil {
		return nil, NewInternalError(err)
	} else if len(ccus) == 0 {
		return make([]*CollectionContentUnit, 0), nil
	}

	ids := make([]int64, len(ccus))
	for i, ccu := range ccus {
		ids[i] = ccu.ContentUnitID
	}
	cus, err := models.ContentUnits(exec,
		qm.Where("secure <= ?", allowedRead(cp)),
		qm.WhereIn("id in ?", utils.ConvertArgsInt64(ids)...),
		qm.Load("ContentUnitI18ns")).
		All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	cusById := make(map[int64]*ContentUnit, len(cus))
	for _, cu := range cus {
		x := ContentUnit{ContentUnit: *cu}
		x.I18n = make(map[string]*models.ContentUnitI18n, len(cu.R.ContentUnitI18ns))
		for _, i18n := range cu.R.ContentUnitI18ns {
			x.I18n[i18n.Language] = i18n
		}
		cusById[x.ID] = &x
	}

	data := make([]*CollectionContentUnit, 0)
	for i := range ccus {
		ccu := ccus[i]
		if cu, ok := cusById[ccu.ContentUnitID]; ok {
			data = append(data, &CollectionContentUnit{
				Name:        ccu.Name,
				Position:    ccu.Position,
				ContentUnit: cu,
			})
		}
	}

	return data, nil
}

func handleCollectionAddCCU(cp utils.ContextProvider, exec boil.Executor, id int64, ccus []*models.CollectionsContentUnit) ([]events.Event, *HttpError) {
	c, err := models.FindCollection(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(c.Secure), common.PERM_WRITE) {
		return nil, NewForbiddenError()
	}

	evnts := make([]events.Event, 1)
	evnts[0] = events.CollectionContentUnitsChangeEvent(c)

	for _, ccu := range ccus {
		cu, err := models.FindContentUnit(exec, ccu.ContentUnitID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, NewBadRequestError(errors.Errorf("Unknown content unit id %d", ccu.ContentUnitID))
			} else {
				return nil, NewInternalError(err)
			}
		}

		exists, err := models.CollectionsContentUnits(exec,
			qm.Where("collection_id = ? AND content_unit_id = ?", id, ccu.ContentUnitID)).
			Exists()
		if err != nil {
			return nil, NewInternalError(err)
		}
		if exists {
			return nil, NewBadRequestError(errors.New("Association already exists"))
		}

		err = c.AddCollectionsContentUnits(exec, true, ccu)
		if err != nil {
			return nil, NewInternalError(err)
		}

		if cu.Published && !c.Published {
			c.Published = true
			if err := c.Update(exec, "published"); err != nil {
				return nil, NewInternalError(err)
			}
			evnts = append(evnts, events.CollectionPublishedChangeEvent(c))
		}
	}

	return evnts, nil
}

func handleCollectionUpdateCCU(cp utils.ContextProvider, exec boil.Executor, id int64, ccu models.CollectionsContentUnit) (*events.Event, *HttpError) {
	c, err := models.FindCollection(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(c.Secure), common.PERM_WRITE) {
		return nil, NewForbiddenError()
	}

	mCCU, err := models.FindCollectionsContentUnit(exec, id, ccu.ContentUnitID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	mCCU.Name = ccu.Name
	mCCU.Position = ccu.Position
	err = mCCU.Update(exec, "name", "position")
	if err != nil {
		return nil, NewInternalError(err)
	}

	e := events.CollectionContentUnitsChangeEvent(c)

	return &e, nil
}

func handleCollectionRemoveCCU(cp utils.ContextProvider, exec boil.Executor, id int64, cuID int64) ([]events.Event, *HttpError) {
	c, err := models.FindCollection(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(c.Secure), common.PERM_WRITE) {
		return nil, NewForbiddenError()
	}

	ccu, err := models.FindCollectionsContentUnit(exec, id, cuID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	err = ccu.Delete(exec)
	if err != nil {
		return nil, NewInternalError(err)
	}

	var evnts = make([]events.Event, 1)
	evnts[0] = events.CollectionContentUnitsChangeEvent(c)

	if c.Published {
		var hasPublishedCUs bool
		query := `SELECT count(*) > 0
                 FROM collections_content_units ccu INNER JOIN content_units cu
                     ON ccu.content_unit_id = cu.id AND ccu.collection_id = $1 AND cu.published IS TRUE`
		if err := queries.Raw(exec, query, id).QueryRow().Scan(&hasPublishedCUs); err != nil {
			return nil, NewInternalError(err)
		}

		if !hasPublishedCUs {
			c.Published = false
			if err := c.Update(exec, "published"); err != nil {
				return nil, NewInternalError(err)
			}
			evnts = append(evnts, events.CollectionPublishedChangeEvent(c))
		}
	}

	return evnts, nil
}

func handleContentUnitsList(cp utils.ContextProvider, exec boil.Executor, r ContentUnitsRequest) (*ContentUnitsResponse, *HttpError) {
	mods := make([]qm.QueryMod, 0)
	appendPermissionsMods(cp, &mods)

	// filters
	if err := appendIDsFilterMods(&mods, r.IDsFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendUIDsFilterMods(&mods, r.UIDsFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendContentTypesFilterMods(&mods, r.ContentTypesFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendDateRangeFilterMods(&mods, r.DateRangeFilter,
		"(coalesce(properties->>'capture_date', properties->>'film_date', created_at::text))::date"); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendSourcesFilterMods(exec, &mods, r.SourcesFilter); err != nil {
		if e, ok := err.(*HttpError); ok {
			return nil, e
		} else {
			NewInternalError(err)
		}
	}
	if err := appendTagsFilterMods(exec, &mods, r.TagsFilter); err != nil {
		return nil, NewInternalError(err)
	}
	if err := appendSearchTermFilterMods(exec, &mods, r.SearchTermFilter, SEARCH_IN_CONTENT_UNITS); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendSecureFilterMods(&mods, r.SecureFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	appendPublishedFilterMods(&mods, r.PublishedFilter)

	// count query
	var total int64
	countMods := append([]qm.QueryMod{qm.Select("count(DISTINCT id)")}, mods...)
	err := models.ContentUnits(exec, countMods...).QueryRow().Scan(&total)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if total == 0 {
		return NewContentUnitsResponse(), nil
	}

	// order, limit, offset
	if err = appendListMods(&mods, r.ListRequest); err != nil {
		return nil, NewBadRequestError(err)
	}

	// Eager loading
	mods = append(mods, qm.Load("ContentUnitI18ns"))

	// data query
	units, err := models.ContentUnits(exec, mods...).All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	// i18n
	data := make([]*ContentUnit, len(units))
	for i, cu := range units {
		x := &ContentUnit{ContentUnit: *cu}
		data[i] = x

		i18ns := cu.R.ContentUnitI18ns

		x.I18n = make(map[string]*models.ContentUnitI18n, len(i18ns))
		for _, i18n := range i18ns {
			x.I18n[i18n.Language] = i18n
		}
	}

	return &ContentUnitsResponse{
		ListResponse: ListResponse{Total: total},
		ContentUnits: data,
	}, nil
}

func handleGetContentUnit(cp utils.ContextProvider, exec boil.Executor, id int64) (*ContentUnit, *HttpError) {
	unit, err := models.ContentUnits(exec,
		qm.Where("id = ?", id),
		qm.Load("ContentUnitI18ns")).
		One()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(unit.Secure), common.PERM_READ) {
		return nil, NewForbiddenError()
	}

	// i18n
	i18ns := unit.R.ContentUnitI18ns
	x := &ContentUnit{ContentUnit: *unit}
	x.I18n = make(map[string]*models.ContentUnitI18n, len(i18ns))
	for _, i18n := range i18ns {
		x.I18n[i18n.Language] = i18n
	}

	return x, nil
}

func handleCreateContentUnit(cp utils.ContextProvider, exec boil.Executor, cu ContentUnit) (*ContentUnit, *HttpError) {
	// check object level permissions
	if !can(cp, secureToPermission(cu.Secure), common.PERM_WRITE) {
		return nil, NewForbiddenError()
	}

	// unmarshal properties
	props := make(map[string]interface{})
	if cu.Properties.Valid {
		err := json.Unmarshal(cu.Properties.JSON, &props)
		if err != nil {
			return nil, NewBadRequestError(errors.Wrap(err, "json.Unmarshal properties"))
		}
	}

	// create content_unit in DB
	ct := common.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name
	unit, err := CreateContentUnit(exec, ct, props)
	if err != nil {
		return nil, NewInternalError(err)
	}

	// save i18n
	for _, v := range cu.I18n {
		err := unit.AddContentUnitI18ns(exec, true, v)
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	return handleGetContentUnit(cp, exec, unit.ID)
}

func handleUpdateContentUnit(cp utils.ContextProvider, exec boil.Executor, cu *PartialContentUnit) (*ContentUnit, *HttpError) {
	unit, err := models.FindContentUnit(exec, cu.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(unit.Secure), common.PERM_WRITE) {
		return nil, NewForbiddenError()
	}

	if unit.TypeID == common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SOURCE].ID {
		return nil, NewBadRequestError(errors.Errorf("Unit type %s is close for change", common.CT_SOURCE))
	}

	if cu.Secure.Valid {
		unit.Secure = cu.Secure.Int16
		err = unit.Update(exec, "secure")
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	// update properties bag
	if cu.Properties.Valid {
		var props map[string]interface{}
		err = cu.Properties.Unmarshal(&props)
		if err != nil {
			return nil, NewInternalError(err)
		}

		err = UpdateContentUnitProperties(exec, unit, props)
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	return handleGetContentUnit(cp, exec, cu.ID)
}

func handleUpdateContentUnitI18n(cp utils.ContextProvider, exec boil.Executor, id int64, i18ns []*models.ContentUnitI18n) (*ContentUnit, *HttpError) {
	unit, err := handleGetContentUnit(cp, exec, id)
	if err != nil {
		return nil, err
	}

	// check object level permissions
	if !can(cp, secureToPermission(unit.Secure), common.PERM_I18N_WRITE) {
		return nil, NewForbiddenError()
	}

	if unit.TypeID == common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SOURCE].ID {
		return nil, NewForbiddenError()
	}

	// Upsert all new i18ns
	nI18n := make(map[string]*models.ContentUnitI18n, len(i18ns))
	for _, i18n := range i18ns {
		i18n.ContentUnitID = id
		nI18n[i18n.Language] = i18n
		err := i18n.Upsert(exec, true,
			[]string{"content_unit_id", "language"},
			[]string{"name", "description"})
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	// Delete old i18ns not in new i18ns
	for k, v := range unit.I18n {
		if _, ok := nI18n[k]; !ok {
			err := v.Delete(exec)
			if err != nil {
				return nil, NewInternalError(err)
			}
		}
	}

	return handleGetContentUnit(cp, exec, id)
}

func handleContentUnitFiles(cp utils.ContextProvider, exec boil.Executor, id int64) ([]*MFile, *HttpError) {
	unit, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(unit.Secure), common.PERM_READ) {
		return nil, NewForbiddenError()
	}

	files, err := models.Files(exec,
		qm.Where("secure <= ?", allowedRead(cp)),
		qm.Where("content_unit_id = ?", id)).
		All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	data := make([]*MFile, len(files))
	for i, f := range files {
		data[i] = NewMFile(f)
	}

	return data, nil
}

func handleContentUnitAddFiles(cp utils.ContextProvider, exec boil.Executor, id int64, fileIDs []int64) (*ContentUnit, []events.Event, *HttpError) {
	unit, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, NewNotFoundError()
		} else {
			return nil, nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(unit.Secure), common.PERM_WRITE) {
		return nil, nil, NewForbiddenError()
	}

	// fetch files
	// With respect to write permissions (as we're about to modify them)
	files, err := models.Files(exec,
		qm.Where("secure <= ?", allowedWrite(cp)),
		qm.WhereIn("id in ?", utils.ConvertArgsInt64(fileIDs)...)).
		All()
	if err != nil {
		return nil, nil, NewInternalError(err)
	}

	if len(files) != len(fileIDs) {
		return nil, nil, NewBadRequestError(errors.New("Couldn't find all files (permissions maybe ?)"))
	}

	// look through files, some might change indeed.
	// collect a set of CU IDs whose files are gone. Any impact on published status ?
	evnts := make([]events.Event, 0)
	somePublished := false
	changedIDs := make([]interface{}, 0)
	possiblyEffectedCUs := make(map[int64]bool)
	for i := range files {
		f := files[i]
		fCUID := f.ContentUnitID.Int64
		if fCUID == id {
			continue
		} else {
			if f.ContentUnitID.Valid {
				possiblyEffectedCUs[fCUID] = possiblyEffectedCUs[fCUID] || f.Published
			}
			f.ContentUnitID = null.Int64From(id) // so we respond with it without re-fetch from db
		}

		changedIDs = append(changedIDs, f.ID)
		evnts = append(evnts, events.FileUpdateEvent(f))

		if f.Published {
			somePublished = true
		}
	}

	// actual update of files that needs to be changed
	err = models.Files(exec,
		qm.WhereIn("id in ?", changedIDs...),
	).UpdateAll(models.M{"content_unit_id": id})
	if err != nil {
		return nil, nil, NewInternalError(err)
	}

	// published status may change for host unit and it's related collections
	impact, err := FileAddedUnitImpact(exec, somePublished, unit.ID)
	if err != nil {
		return nil, nil, NewInternalError(err)
	}
	evnts = append(evnts, impact.Events()...)

	// published status may change for units we leave as well
	for k, v := range possiblyEffectedCUs {
		impact, err := FileLeftUnitImpact(exec, v, k)
		if err != nil {
			return nil, nil, NewInternalError(err)
		}
		evnts = append(evnts, impact.Events()...)
	}

	resp, herr := handleGetContentUnit(cp, exec, id)

	return resp, evnts, herr
}

func handleContentUnitCCU(cp utils.ContextProvider, exec boil.Executor, id int64) ([]*CollectionContentUnit, *HttpError) {
	unit, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(unit.Secure), common.PERM_READ) {
		return nil, NewForbiddenError()
	}

	ccus, err := models.CollectionsContentUnits(exec,
		qm.Where("content_unit_id = ?", id)).
		All()
	if err != nil {
		return nil, NewInternalError(err)
	} else if len(ccus) == 0 {
		return make([]*CollectionContentUnit, 0), nil
	}

	ids := make([]int64, len(ccus))
	for i, ccu := range ccus {
		ids[i] = ccu.CollectionID
	}
	cs, err := models.Collections(exec,
		qm.Where("secure <= ?", allowedRead(cp)),
		qm.WhereIn("id in ?", utils.ConvertArgsInt64(ids)...),
		qm.Load("CollectionI18ns")).
		All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	csById := make(map[int64]*Collection, len(cs))
	for _, c := range cs {
		x := Collection{Collection: *c}
		x.I18n = make(map[string]*models.CollectionI18n, len(c.R.CollectionI18ns))
		for _, i18n := range c.R.CollectionI18ns {
			x.I18n[i18n.Language] = i18n
		}
		csById[x.ID] = &x
	}

	data := make([]*CollectionContentUnit, len(ccus))
	for i, ccu := range ccus {
		data[i] = &CollectionContentUnit{
			Name:       ccu.Name,
			Position:   ccu.Position,
			Collection: csById[ccu.CollectionID],
		}
	}

	return data, nil
}

func handleContentUnitCUD(cp utils.ContextProvider, exec boil.Executor, id int64) ([]*ContentUnitDerivation, *HttpError) {
	unit, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(unit.Secure), common.PERM_READ) {
		return nil, NewForbiddenError()
	}

	cuds, err := models.ContentUnitDerivations(exec,
		qm.Where("source_id = ?", id)).
		All()
	if err != nil {
		return nil, NewInternalError(err)
	} else if len(cuds) == 0 {
		return make([]*ContentUnitDerivation, 0), nil
	}

	ids := make([]int64, len(cuds))
	for i := range cuds {
		ids[i] = cuds[i].DerivedID
	}
	cus, err := models.ContentUnits(exec,
		qm.Where("secure <= ?", allowedRead(cp)),
		qm.WhereIn("id in ?", utils.ConvertArgsInt64(ids)...),
		qm.Load("ContentUnitI18ns")).
		All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	cusById := make(map[int64]*ContentUnit, len(cus))
	for _, cu := range cus {
		x := ContentUnit{ContentUnit: *cu}
		x.I18n = make(map[string]*models.ContentUnitI18n, len(cu.R.ContentUnitI18ns))
		for _, i18n := range cu.R.ContentUnitI18ns {
			x.I18n[i18n.Language] = i18n
		}
		cusById[x.ID] = &x
	}

	data := make([]*ContentUnitDerivation, len(cuds))
	for i, cud := range cuds {
		data[i] = &ContentUnitDerivation{
			Derived: cusById[cud.DerivedID],
			Name:    cud.Name,
		}
	}

	return data, nil
}

func handleContentUnitAddCUD(cp utils.ContextProvider, exec boil.Executor, id int64, cud models.ContentUnitDerivation) (*models.ContentUnit, *HttpError) {
	cu, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(cu.Secure), common.PERM_WRITE) {
		return nil, NewForbiddenError()
	}

	exists, err := models.ContentUnits(exec,
		qm.Where("id = ?", cud.DerivedID)).
		Exists()
	if err != nil {
		return nil, NewInternalError(err)
	}
	if !exists {
		return nil, NewBadRequestError(errors.Errorf("Unknown content unit id %d", cud.DerivedID))
	}

	exists, err = models.ContentUnitDerivations(exec,
		qm.Where("source_id = ? AND derived_id = ?", id, cud.DerivedID)).
		Exists()
	if err != nil {
		return nil, NewInternalError(err)
	}
	if exists {
		return nil, NewBadRequestError(errors.New("Derivation already exists"))
	}

	err = cu.AddSourceContentUnitDerivations(exec, true, &cud)
	if err != nil {
		return nil, NewInternalError(err)
	}

	return cu, nil
}

func handleContentUnitUpdateCUD(cp utils.ContextProvider, exec boil.Executor, id int64, cud models.ContentUnitDerivation) (*models.ContentUnit, *HttpError) {
	cu, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(cu.Secure), common.PERM_WRITE) {
		return nil, NewForbiddenError()
	}

	mCUD, err := models.FindContentUnitDerivation(exec, id, cud.DerivedID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	mCUD.Name = cud.Name
	err = mCUD.Update(exec, "name")
	if err != nil {
		return nil, NewInternalError(err)
	}

	return cu, nil
}

func handleContentUnitRemoveCUD(cp utils.ContextProvider, exec boil.Executor, id int64, duID int64) (*models.ContentUnit, *HttpError) {
	cu, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(cu.Secure), common.PERM_WRITE) {
		return nil, NewForbiddenError()
	}

	cud, err := models.FindContentUnitDerivation(exec, id, duID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	err = cud.Delete(exec)
	if err != nil {
		return nil, NewInternalError(err)
	}

	return cu, nil
}

func handleContentUnitOrigins(cp utils.ContextProvider, exec boil.Executor, id int64) ([]*ContentUnitDerivation, *HttpError) {
	unit, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(unit.Secure), common.PERM_READ) {
		return nil, NewForbiddenError()
	}

	cuds, err := models.ContentUnitDerivations(exec,
		qm.Where("derived_id = ?", id)).
		All()
	if err != nil {
		return nil, NewInternalError(err)
	} else if len(cuds) == 0 {
		return make([]*ContentUnitDerivation, 0), nil
	}

	ids := make([]int64, len(cuds))
	for i := range cuds {
		ids[i] = cuds[i].SourceID
	}
	cus, err := models.ContentUnits(exec,
		qm.Where("secure <= ?", allowedRead(cp)),
		qm.WhereIn("id in ?", utils.ConvertArgsInt64(ids)...),
		qm.Load("ContentUnitI18ns")).
		All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	cusById := make(map[int64]*ContentUnit, len(cus))
	for _, cu := range cus {
		x := ContentUnit{ContentUnit: *cu}
		x.I18n = make(map[string]*models.ContentUnitI18n, len(cu.R.ContentUnitI18ns))
		for _, i18n := range cu.R.ContentUnitI18ns {
			x.I18n[i18n.Language] = i18n
		}
		cusById[x.ID] = &x
	}

	data := make([]*ContentUnitDerivation, len(cuds))
	for i, cud := range cuds {
		data[i] = &ContentUnitDerivation{
			Name:   cud.Name,
			Source: cusById[cud.SourceID],
		}
	}

	return data, nil
}

func handleGetContentUnitSources(cp utils.ContextProvider, exec boil.Executor, id int64) ([]*Source, *HttpError) {
	unit, err := models.ContentUnits(exec,
		qm.Where("id = ?", id),
		qm.Load("Sources", "Sources.SourceI18ns")).
		One()

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(unit.Secure), common.PERM_READ) {
		return nil, NewForbiddenError()
	}

	data := make([]*Source, len(unit.R.Sources))
	for i, source := range unit.R.Sources {
		x := &Source{Source: *source}
		x.I18n = make(map[string]*models.SourceI18n, len(source.R.SourceI18ns))
		for _, i18n := range source.R.SourceI18ns {
			x.I18n[i18n.Language] = i18n
		}
		data[i] = x
	}

	return data, nil
}

func handleContentUnitAddSource(cp utils.ContextProvider, exec boil.Executor, id int64, sourceID int64) (*models.ContentUnit, *HttpError) {
	cu, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(cu.Secure), common.PERM_METADATA_WRITE) {
		return nil, NewForbiddenError()
	}

	source, err := models.FindSource(exec, sourceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewBadRequestError(errors.Errorf("Unknown source id %d", sourceID))
		} else {
			return nil, NewInternalError(err)
		}
	}

	var count int64
	err = queries.Raw(exec,
		`SELECT COUNT(1) FROM content_units_sources WHERE content_unit_id=$1 AND source_id=$2`,
		cu.ID, source.ID).
		QueryRow().
		Scan(&count)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if count > 0 {
		return nil, nil // noop
	}

	err = cu.AddSources(exec, false, source)
	if err != nil {
		return nil, NewInternalError(err)
	}

	return cu, nil
}

func handleContentUnitRemoveSource(cp utils.ContextProvider, exec boil.Executor, id int64, sourceID int64) (*models.ContentUnit, *HttpError) {
	cu, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(cu.Secure), common.PERM_METADATA_WRITE) {
		return nil, NewForbiddenError()
	}

	source, err := models.FindSource(exec, sourceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewBadRequestError(errors.Errorf("Unknown source id %d", sourceID))
		} else {
			return nil, NewInternalError(err)
		}
	}

	err = cu.RemoveSources(exec, source)
	if err != nil {
		return nil, NewInternalError(err)
	}

	return cu, nil
}

func handleGetContentUnitTags(cp utils.ContextProvider, exec boil.Executor, id int64) ([]*Tag, *HttpError) {
	unit, err := models.ContentUnits(exec,
		qm.Where("id = ?", id),
		qm.Load("Tags", "Tags.TagI18ns")).
		One()

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(unit.Secure), common.PERM_READ) {
		return nil, NewForbiddenError()
	}

	data := make([]*Tag, len(unit.R.Tags))
	for i, tag := range unit.R.Tags {
		x := &Tag{Tag: *tag}
		x.I18n = make(map[string]*models.TagI18n, len(tag.R.TagI18ns))
		for _, i18n := range tag.R.TagI18ns {
			x.I18n[i18n.Language] = i18n
		}
		data[i] = x
	}

	return data, nil
}

func handleContentUnitAddTag(cp utils.ContextProvider, exec boil.Executor, id int64, tagID int64) (*models.ContentUnit, *HttpError) {
	cu, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(cu.Secure), common.PERM_METADATA_WRITE) {
		return nil, NewForbiddenError()
	}

	tag, err := models.FindTag(exec, tagID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewBadRequestError(errors.Errorf("Unknown tag id %d", tagID))
		} else {
			return nil, NewInternalError(err)
		}
	}

	var count int64
	err = queries.Raw(exec,
		`SELECT COUNT(1) FROM content_units_tags WHERE content_unit_id=$1 AND tag_id=$2`,
		cu.ID, tag.ID).
		QueryRow().
		Scan(&count)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if count > 0 {
		return nil, nil // noop
	}

	err = cu.AddTags(exec, false, tag)
	if err != nil {
		return nil, NewInternalError(err)
	}

	return cu, nil
}

func handleContentUnitRemoveTag(cp utils.ContextProvider, exec boil.Executor, id int64, tagID int64) (*models.ContentUnit, *HttpError) {
	cu, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(cu.Secure), common.PERM_METADATA_WRITE) {
		return nil, NewForbiddenError()
	}

	tag, err := models.FindTag(exec, tagID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewBadRequestError(errors.Errorf("Unknown tag id %d", tagID))
		} else {
			return nil, NewInternalError(err)
		}
	}

	err = cu.RemoveTags(exec, tag)
	if err != nil {
		return nil, NewInternalError(err)
	}

	return cu, nil
}

func handleGetContentUnitPersons(cp utils.ContextProvider, exec boil.Executor, id int64) ([]*ContentUnitPerson, *HttpError) {
	unit, err := models.ContentUnits(exec,
		qm.Where("id = ?", id),
		qm.Load("ContentUnitsPersons",
			"ContentUnitsPersons.Person",
			"ContentUnitsPersons.Person.PersonI18ns")).
		One()

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(unit.Secure), common.PERM_READ) {
		return nil, NewForbiddenError()
	}

	data := make([]*ContentUnitPerson, len(unit.R.ContentUnitsPersons))
	for i, cup := range unit.R.ContentUnitsPersons {
		p := &Person{Person: *cup.R.Person}
		p.I18n = make(map[string]*models.PersonI18n, len(cup.R.Person.R.PersonI18ns))
		for _, i18n := range cup.R.Person.R.PersonI18ns {
			p.I18n[i18n.Language] = i18n
		}
		data[i] = &ContentUnitPerson{Person: p, RoleID: cup.RoleID}
	}

	return data, nil
}

func handleContentUnitAddPerson(cp utils.ContextProvider, exec boil.Executor, id int64, cup models.ContentUnitsPerson) (*models.ContentUnit, *HttpError) {
	cu, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(cu.Secure), common.PERM_METADATA_WRITE) {
		return nil, NewForbiddenError()
	}

	exists, err := models.PersonExists(exec, cup.PersonID)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if !exists {
		return nil, NewBadRequestError(errors.Errorf("Unknown person id %d", cup.PersonID))
	}

	exists, err = models.ContentRoleTypeExists(exec, cup.RoleID)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if !exists {
		return nil, NewBadRequestError(errors.Errorf("Unknown role id %d", cup.RoleID))
	}

	existingCUP, err := models.FindContentUnitsPerson(exec, id, cup.PersonID)
	if err != nil {
		if err == sql.ErrNoRows {
			// create new
			cup.ContentUnitID = id
			err = cup.Insert(exec)
			if err != nil {
				return nil, NewInternalError(err)
			}
			return cu, nil
		} else {
			return nil, NewInternalError(err)
		}
	}

	// update role
	existingCUP.RoleID = cup.RoleID
	err = existingCUP.Update(exec, "role_id")
	if err != nil {
		return nil, NewInternalError(err)
	}

	return cu, nil
}

func handleContentUnitRemovePerson(cp utils.ContextProvider, exec boil.Executor, id int64, personID int64) (*models.ContentUnit, *HttpError) {
	cu, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(cu.Secure), common.PERM_METADATA_WRITE) {
		return nil, NewForbiddenError()
	}

	cup, err := models.FindContentUnitsPerson(exec, id, personID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	err = cup.Delete(exec)
	if err != nil {
		return nil, NewInternalError(err)
	}

	return cu, nil
}

func handleGetContentUnitPublishers(cp utils.ContextProvider, exec boil.Executor, id int64) ([]*Publisher, *HttpError) {
	unit, err := models.ContentUnits(exec,
		qm.Where("id = ?", id),
		qm.Load("Publishers", "Publishers.PublisherI18ns")).
		One()

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(unit.Secure), common.PERM_READ) {
		return nil, NewForbiddenError()
	}

	data := make([]*Publisher, len(unit.R.Publishers))
	for i, publisher := range unit.R.Publishers {
		p := &Publisher{Publisher: *publisher}
		p.I18n = make(map[string]*models.PublisherI18n, len(publisher.R.PublisherI18ns))
		for _, i18n := range publisher.R.PublisherI18ns {
			p.I18n[i18n.Language] = i18n
		}
		data[i] = p
	}

	return data, nil
}

func handleContentUnitAddPublisher(cp utils.ContextProvider, exec boil.Executor, id int64, publisherID int64) (*models.ContentUnit, *HttpError) {
	cu, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(cu.Secure), common.PERM_METADATA_WRITE) {
		return nil, NewForbiddenError()
	}

	publisher, err := models.FindPublisher(exec, publisherID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewBadRequestError(errors.Errorf("Unknown publisher id %d", publisherID))
		} else {
			return nil, NewInternalError(err)
		}
	}

	var count int64
	err = queries.Raw(exec,
		`SELECT COUNT(1) FROM content_units_publishers WHERE content_unit_id=$1 AND publisher_id=$2`,
		cu.ID, publisher.ID).
		QueryRow().
		Scan(&count)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if count > 0 {
		return nil, nil // noop
	}

	err = cu.AddPublishers(exec, false, publisher)
	if err != nil {
		return nil, NewInternalError(err)
	}

	return cu, nil
}

func handleContentUnitRemovePublisher(cp utils.ContextProvider, exec boil.Executor, id int64, publisherID int64) (*models.ContentUnit, *HttpError) {
	cu, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(cu.Secure), common.PERM_METADATA_WRITE) {
		return nil, NewForbiddenError()
	}

	publisher, err := models.FindPublisher(exec, publisherID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewBadRequestError(errors.Errorf("Unknown publisher id %d", publisherID))
		} else {
			return nil, NewInternalError(err)
		}
	}

	err = cu.RemovePublishers(exec, publisher)
	if err != nil {
		return nil, NewInternalError(err)
	}

	return cu, nil
}

func handleContentUnitMerge(cp utils.ContextProvider, exec boil.Executor, id int64, cuIDs []int64) (*ContentUnit, []events.Event, *HttpError) {
	unit, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, NewNotFoundError()
		} else {
			return nil, nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(unit.Secure), common.PERM_WRITE) {
		return nil, nil, NewForbiddenError()
	}

	// fetch units to be merged
	// With respect to write permissions (as we're about to modify them)
	units, err := models.ContentUnits(exec,
		qm.Where("secure <= ?", allowedWrite(cp)),
		qm.WhereIn("id in ?", utils.ConvertArgsInt64(cuIDs)...),
		qm.Load("Files", "DerivedContentUnitDerivations", "SourceContentUnitDerivations")).
		All()
	if err != nil {
		return nil, nil, NewInternalError(err)
	}

	if len(units) != len(cuIDs) {
		return nil, nil, NewBadRequestError(errors.New("Couldn't find all units (permissions maybe ?)"))
	}

	// for each merged unit we:
	// 1. move all it's files to host unit
	// 2. move all derivations to host unit
	// 3. remove unit

	evnts := make([]events.Event, 0)
	somePublished := false
	cuDerivativesChange := false
	for i := range units {
		cu := units[i]
		log.Infof("Merging CU %d into CU %d", cu.ID, unit.ID)
		if cu.Published {
			somePublished = true
		}

		// move files
		err := cu.R.Files.UpdateAll(exec, models.M{"content_unit_id": unit.ID})
		if err != nil {
			return nil, nil, NewInternalError(err)
		}
		for i := range cu.R.Files {
			evnts = append(evnts, events.FileUpdateEvent(cu.R.Files[i]))
		}

		// move derivations
		err = cu.R.DerivedContentUnitDerivations.UpdateAll(exec, models.M{"derived_id": unit.ID})
		if err != nil {
			return nil, nil, NewInternalError(err)
		}
		for i := range cu.R.DerivedContentUnitDerivations {
			sourceCU, err := cu.R.DerivedContentUnitDerivations[i].Source(exec).One()
			if err != nil {
				return nil, nil, NewInternalError(err)
			}
			evnts = append(evnts, events.ContentUnitDerivativesChangeEvent(sourceCU))
		}

		err = cu.R.SourceContentUnitDerivations.UpdateAll(exec, models.M{"source_id": unit.ID})
		if err != nil {
			return nil, nil, NewInternalError(err)
		}
		if len(cu.R.SourceContentUnitDerivations) > 0 {
			cuDerivativesChange = true
		}

		// remove unit
		err = DeleteContentUnit(exec, cu)
		if err != nil {
			return nil, nil, NewInternalError(err)
		}
		evnts = append(evnts, events.ContentUnitDeleteEvent(cu))
	}

	if cuDerivativesChange {
		evnts = append(evnts, events.ContentUnitDerivativesChangeEvent(unit))
	}

	// published status may change for host unit and it's related collections
	impact, err := FileAddedUnitImpact(exec, somePublished, unit.ID)
	if err != nil {
		return nil, nil, NewInternalError(err)
	}
	evnts = append(evnts, impact.Events()...)

	resp, herr := handleGetContentUnit(cp, exec, id)

	return resp, evnts, herr
}

func handleFilesList(cp utils.ContextProvider, exec boil.Executor, r FilesRequest) (*FilesResponse, *HttpError) {
	mods := make([]qm.QueryMod, 0)
	appendPermissionsMods(cp, &mods)

	// filters
	if err := appendIDsFilterMods(&mods, r.IDsFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendUIDsFilterMods(&mods, r.UIDsFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendSHA1sFilterMods(&mods, r.SHA1sFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendDateRangeFilterMods(&mods, r.DateRangeFilter, "file_created_at"); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendSecureFilterMods(&mods, r.SecureFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	appendPublishedFilterMods(&mods, r.PublishedFilter)
	if err := appendSearchTermFilterMods(exec, &mods, r.SearchTermFilter, SEARCH_IN_FILES); err != nil {
		return nil, NewBadRequestError(err)
	}
	/*if r.Query != "" {
		mods = append(mods, qm.Where("name ~ ?", r.Query),
			qm.Or("uid ~ ?", r.Query),
			qm.Or("id::TEXT ~ ?", r.Query))
	}*/

	// count query
	var total int64
	countMods := append([]qm.QueryMod{qm.Select("count(DISTINCT id)")}, mods...)
	err := models.Files(exec, countMods...).QueryRow().Scan(&total)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if total == 0 {
		return NewFilesResponse(), nil
	}

	// order, limit, offset
	if err = appendListMods(&mods, r.ListRequest); err != nil {
		return nil, NewBadRequestError(err)
	}

	// data query
	files, err := models.Files(exec, mods...).All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	data := make([]*MFile, len(files))
	for i, f := range files {
		data[i] = NewMFile(f)
	}

	return &FilesResponse{
		ListResponse: ListResponse{Total: total},
		Files:        data,
	}, nil
}

func handleGetFile(cp utils.ContextProvider, exec boil.Executor, id int64) (*MFile, *HttpError) {
	file, err := models.FindFile(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(file.Secure), common.PERM_READ) {
		return nil, NewForbiddenError()
	}

	return NewMFile(file), nil
}

func handleUpdateFile(cp utils.ContextProvider, exec boil.Executor, f *PartialFile) (*MFile, []events.Event, *HttpError) {
	file, err := models.FindFile(exec, f.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, NewNotFoundError()
		} else {
			return nil, nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(file.Secure), common.PERM_WRITE) {
		return nil, nil, NewForbiddenError()
	}

	evnts := make([]events.Event, 0)

	if f.Type.Valid {
		file.Type = f.Type.String
	}
	if f.SubType.Valid {
		file.SubType = f.SubType.String
	}
	if f.MimeType.Valid {
		file.MimeType = f.MimeType
	}
	if f.Language.Valid {
		file.Language = f.Language
	}
	if f.ParentID.Valid {
		file.ParentID = f.ParentID
	}
	if f.Secure.Valid {
		file.Secure = f.Secure.Int16
	}

	prevCUID := file.ContentUnitID.Int64
	if f.ContentUnitID.Valid {
		file.ContentUnitID = f.ContentUnitID
	}

	err = file.Update(exec)
	if err != nil {
		return nil, nil, NewInternalError(err)
	}

	evnts = append(evnts, events.FileUpdateEvent(file))

	// We might be leaving some unit and joining another.
	// What should be the impact of their published status ?

	// The unit we're joining
	if f.ContentUnitID.Valid {
		impact, err := FileAddedUnitImpact(exec, file.Published, f.ContentUnitID.Int64)
		if err != nil {
			return nil, nil, NewInternalError(err)
		}
		evnts = append(evnts, impact.Events()...)
	}

	// The unit we're leaving
	if prevCUID != 0 {
		impact, err := FileLeftUnitImpact(exec, file.Published, prevCUID)
		if err != nil {
			return nil, nil, NewInternalError(err)
		}
		evnts = append(evnts, impact.Events()...)
	}

	resp, herr := handleGetFile(cp, exec, f.ID)
	return resp, evnts, herr
}

func handleFileStorages(cp utils.ContextProvider, exec boil.Executor, id int64) ([]*Storage, *HttpError) {
	file, err := models.FindFile(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// check object level permissions
	if !can(cp, secureToPermission(file.Secure), common.PERM_READ) {
		return nil, NewForbiddenError()
	}

	storages, err := models.Storages(exec,
		qm.InnerJoin("files_storages fs on fs.storage_id=id and fs.file_id = ?", id)).
		All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	data := make([]*Storage, len(storages))
	for i := range storages {
		data[i] = &Storage{Storage: *storages[i]}
	}

	return data, nil
}

func handleOperationsList(exec boil.Executor, r OperationsRequest) (*OperationsResponse, *HttpError) {
	mods := make([]qm.QueryMod, 0)

	// filters
	if err := appendDateRangeFilterMods(&mods, r.DateRangeFilter, "created_at"); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendOperationTypesFilterMods(&mods, r.OperationTypesFilter); err != nil {
		return nil, NewBadRequestError(err)
	}

	// count query
	var total int64
	countMods := append([]qm.QueryMod{qm.Select("count(DISTINCT id)")}, mods...)
	err := models.Operations(exec, countMods...).QueryRow().Scan(&total)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if total == 0 {
		return NewOperationsResponse(), nil
	}

	// order, limit, offset
	if err = appendListMods(&mods, r.ListRequest); err != nil {
		return nil, NewBadRequestError(err)
	}

	// data query
	data, err := models.Operations(exec, mods...).All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	return &OperationsResponse{
		ListResponse: ListResponse{Total: total},
		Operations:   data,
	}, nil
}

func handleOperationItem(exec boil.Executor, id int64) (*models.Operation, *HttpError) {
	operation, err := models.FindOperation(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	return operation, nil
}

func handleOperationFiles(cp utils.ContextProvider, exec boil.Executor, id int64) ([]*MFile, *HttpError) {
	ok, err := models.OperationExists(exec, id)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if !ok {
		return nil, NewNotFoundError()
	}

	files, err := models.Files(exec,
		qm.InnerJoin("files_operations fo on fo.file_id=id and fo.operation_id = ? and secure <= ?",
			id, allowedRead(cp))).
		All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	data := make([]*MFile, len(files))
	for i, f := range files {
		data[i] = NewMFile(f)
	}

	return data, nil
}

func handleGetSources(exec boil.Executor, r SourcesRequest) (*SourcesResponse, *HttpError) {
	mods := make([]qm.QueryMod, 0)

	// count query
	var total int64
	countMods := append([]qm.QueryMod{qm.Select("count(DISTINCT id)")}, mods...)
	err := models.Sources(exec, countMods...).QueryRow().Scan(&total)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if total == 0 {
		return NewSourcesResponse(), nil
	}

	// order, limit, offset
	if err = appendListMods(&mods, r.ListRequest); err != nil {
		return nil, NewBadRequestError(err)
	}

	// Eager loading
	mods = append(mods, qm.Load("SourceI18ns"))

	// data query
	sources, err := models.Sources(exec, mods...).All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	// i18n
	data := make([]*Source, len(sources))
	for i, s := range sources {
		x := &Source{Source: *s}
		data[i] = x
		x.I18n = make(map[string]*models.SourceI18n, len(s.R.SourceI18ns))
		for _, i18n := range s.R.SourceI18ns {
			x.I18n[i18n.Language] = i18n
		}
	}

	return &SourcesResponse{
		ListResponse: ListResponse{Total: total},
		Sources:      data,
	}, nil
}

func handleCreateSource(exec boil.Executor, r CreateSourceRequest) (*Source, *HttpError) {
	s := r.Source

	// check pattern unique constraint
	if s.Pattern.Valid {
		ok, err := models.Sources(exec, qm.Where("pattern = ?", s.Pattern.String)).Exists()
		if err != nil {
			return nil, NewInternalError(err)
		}
		if ok {
			err = errors.Errorf("Pattern already in use: %s", s.Pattern.String)
			return nil, NewBadRequestError(err)
		}
	}

	// make sure parent source exists if given
	if s.ParentID.Valid {
		ok, err := models.Sources(exec, qm.Where("id = ?", s.ParentID.Int64)).Exists()
		if err != nil {
			return nil, NewInternalError(err)
		}
		if !ok {
			err = errors.Errorf("Unknown parent source: %d", s.ParentID.Int64)
			return nil, NewBadRequestError(err)
		}
	}

	// make sure author exists if given
	if r.AuthorID.Valid {
		ok, err := models.Authors(exec, qm.Where("id = ?", r.AuthorID.Int64)).Exists()
		if err != nil {
			return nil, NewInternalError(err)
		}
		if !ok {
			err = errors.Errorf("Unknown author: %d", r.AuthorID.Int64)
			return nil, NewBadRequestError(err)
		}
	}

	// save source to DB
	uid, err := GetFreeUID(exec, new(SourceUIDChecker))
	if err != nil {
		return nil, NewInternalError(err)
	}
	s.UID = uid

	err = s.Source.Insert(exec)
	if err != nil {
		return nil, NewInternalError(err)
	}

	// save i18n
	for _, v := range s.I18n {
		err := s.AddSourceI18ns(exec, true, v)
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	// save author
	if r.AuthorID.Valid {
		err := s.Source.AddAuthors(exec, false, &models.Author{ID: r.AuthorID.Int64})
		if err != nil {
			return nil, NewInternalError(err)
		}
	}
	// create CU type source

	cuUid := s.UID
	hasCU, err := models.ContentUnits(exec, qm.Where("uid = ?", cuUid)).Exists()
	if err != nil {
		return nil, NewInternalError(err)
	}
	if hasCU {
		cuUid, err = GetFreeUID(exec, new(ContentUnitUIDChecker))
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	props, _ := json.Marshal(map[string]string{"source_id": s.UID})

	cu := &models.ContentUnit{
		UID:        s.UID,
		TypeID:     common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SOURCE].ID,
		Secure:     common.SEC_PUBLIC,
		Published:  true,
		Properties: null.JSONFrom(props),
	}

	err = cu.Insert(exec)
	if err != nil {
		return nil, NewInternalError(err)
	}

	return handleGetSource(exec, s.ID)
}

func handleGetSource(exec boil.Executor, id int64) (*Source, *HttpError) {
	source, err := models.Sources(exec,
		qm.Where("id = ?", id),
		qm.Load("SourceI18ns")).
		One()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// i18n
	x := &Source{Source: *source}
	x.I18n = make(map[string]*models.SourceI18n, len(source.R.SourceI18ns))
	for _, i18n := range source.R.SourceI18ns {
		x.I18n[i18n.Language] = i18n
	}
	return x, nil
}

func handleUpdateSource(exec boil.Executor, s *Source) (*Source, *HttpError) {
	source, err := models.FindSource(exec, s.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	if s.TypeID != 0 { // to allow partial updates
		source.TypeID = s.TypeID
	}
	source.Pattern = s.Pattern
	source.Description = s.Description
	err = s.Update(exec, "pattern", "description", "type_id")
	if err != nil {
		return nil, NewInternalError(err)
	}

	return handleGetSource(exec, s.ID)
}

func handleUpdateSourceI18n(exec boil.Executor, id int64, i18ns []*models.SourceI18n) (*Source, *HttpError) {
	source, err := handleGetSource(exec, id)
	if err != nil {
		return nil, err
	}

	// Upsert all new i18ns
	nI18n := make(map[string]*models.SourceI18n, len(i18ns))
	for _, i18n := range i18ns {
		i18n.SourceID = id
		nI18n[i18n.Language] = i18n
		err := i18n.Upsert(exec, true,
			[]string{"source_id", "language"},
			[]string{"name", "description"})
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	// Delete old i18ns not in new i18ns
	for k, v := range source.I18n {
		if _, ok := nI18n[k]; !ok {
			err := v.Delete(exec)
			if err != nil {
				return nil, NewInternalError(err)
			}
		}
	}

	return handleGetSource(exec, id)
}

func handleGetTags(exec boil.Executor, r TagsRequest) (*TagsResponse, *HttpError) {
	mods := make([]qm.QueryMod, 0)

	// count query
	var total int64
	countMods := append([]qm.QueryMod{qm.Select("count(DISTINCT id)")}, mods...)
	err := models.Tags(exec, countMods...).QueryRow().Scan(&total)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if total == 0 {
		return NewTagsResponse(), nil
	}

	// order, limit, offset
	if err = appendListMods(&mods, r.ListRequest); err != nil {
		return nil, NewBadRequestError(err)
	}

	// Eager loading
	mods = append(mods, qm.Load("TagI18ns"))

	// data query
	tags, err := models.Tags(exec, mods...).All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	// i18n
	data := make([]*Tag, len(tags))
	for i, t := range tags {
		x := &Tag{Tag: *t}
		data[i] = x
		x.I18n = make(map[string]*models.TagI18n, len(t.R.TagI18ns))
		for _, i18n := range t.R.TagI18ns {
			x.I18n[i18n.Language] = i18n
		}
	}

	return &TagsResponse{
		ListResponse: ListResponse{Total: total},
		Tags:         data,
	}, nil
}

func handleCreateTag(exec boil.Executor, t *Tag) (*Tag, *HttpError) {
	// make sure parent tag exists if given
	if t.ParentID.Valid {
		ok, err := models.Tags(exec, qm.Where("id = ?", t.ParentID.Int64)).Exists()
		if err != nil {
			return nil, NewInternalError(err)
		}
		if !ok {
			return nil, NewBadRequestError(errors.Errorf("Unknown parent tag %d", t.ParentID.Int64))
		}
	}

	// save tag to DB
	uid, err := GetFreeUID(exec, new(TagUIDChecker))
	if err != nil {
		return nil, NewInternalError(err)
	}
	t.UID = uid

	err = t.Tag.Insert(exec)
	if err != nil {
		return nil, NewInternalError(err)
	}

	// save i18n
	for _, v := range t.I18n {
		err := t.AddTagI18ns(exec, true, v)
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	return handleGetTag(exec, t.ID)
}

func handleGetTag(exec boil.Executor, id int64) (*Tag, *HttpError) {
	tag, err := models.Tags(exec,
		qm.Where("id = ?", id),
		qm.Load("TagI18ns")).
		One()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// i18n
	x := &Tag{Tag: *tag}
	x.I18n = make(map[string]*models.TagI18n, len(tag.R.TagI18ns))
	for _, i18n := range tag.R.TagI18ns {
		x.I18n[i18n.Language] = i18n
	}

	return x, nil
}

func handleUpdateTag(exec boil.Executor, t *Tag) (*Tag, *HttpError) {
	tag, err := models.FindTag(exec, t.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	tag.Pattern = t.Pattern
	tag.Description = t.Description
	err = t.Update(exec, "pattern", "description")
	if err != nil {
		return nil, NewInternalError(err)
	}

	return handleGetTag(exec, t.ID)
}

func handleUpdateTagI18n(exec boil.Executor, id int64, i18ns []*models.TagI18n) (*Tag, *HttpError) {
	tag, err := handleGetTag(exec, id)
	if err != nil {
		return nil, err
	}

	// Upsert all new i18ns
	nI18n := make(map[string]*models.TagI18n, len(i18ns))
	for _, i18n := range i18ns {
		i18n.TagID = id
		nI18n[i18n.Language] = i18n
		err := i18n.Upsert(exec, true, []string{"tag_id", "language"}, []string{"label"})
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	// Delete old i18ns not in new i18ns
	for k, v := range tag.I18n {
		if _, ok := nI18n[k]; !ok {
			err := v.Delete(exec)
			if err != nil {
				return nil, NewInternalError(err)
			}
		}
	}

	return handleGetTag(exec, id)
}

func handlePersonsList(exec boil.Executor, r PersonsRequest) (*PersonsResponse, *HttpError) {
	mods := make([]qm.QueryMod, 0)

	// filters
	if err := appendIDsFilterMods(&mods, r.IDsFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendUIDsFilterMods(&mods, r.UIDsFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendPatternsFilterMods(&mods, r.PatternsFilter); err != nil {
		return nil, NewBadRequestError(err)
	}

	// count query
	var total int64
	countMods := append([]qm.QueryMod{qm.Select("count(DISTINCT id)")}, mods...)
	err := models.Persons(exec, countMods...).QueryRow().Scan(&total)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if total == 0 {
		return NewPersonsResponse(), nil
	}

	// order, limit, offset
	if err = appendListMods(&mods, r.ListRequest); err != nil {
		return nil, NewBadRequestError(err)
	}

	// Eager loading
	mods = append(mods, qm.Load("PersonI18ns"))

	// data query
	persons, err := models.Persons(exec, mods...).All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	// i18n
	data := make([]*Person, len(persons))
	for i, pr := range persons {
		x := &Person{Person: *pr}
		data[i] = x
		x.I18n = make(map[string]*models.PersonI18n, len(pr.R.PersonI18ns))
		for _, i18n := range pr.R.PersonI18ns {
			x.I18n[i18n.Language] = i18n
		}
	}

	return &PersonsResponse{
		ListResponse: ListResponse{Total: total},
		Persons:      data,
	}, nil
}

func handleCreatePerson(exec boil.Executor, p *Person) (*Person, *HttpError) {

	// save person to DB
	uid, err := GetFreeUID(exec, new(PersonUIDChecker))
	if err != nil {
		return nil, NewInternalError(err)
	}
	p.UID = uid

	err = p.Person.Insert(exec)
	if err != nil {
		return nil, NewInternalError(err)
	}

	// save i18n
	for _, v := range p.I18n {
		err := p.AddPersonI18ns(exec, true, v)
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	return handleGetPerson(exec, p.ID)
}

func handleUpdatePerson(exec boil.Executor, p *Person) (*Person, *HttpError) {
	person, err := models.FindPerson(exec, p.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	person.Pattern = p.Pattern
	err = p.Update(exec, "pattern")
	if err != nil {
		return nil, NewInternalError(err)
	}

	return handleGetPerson(exec, p.ID)
}

func handleGetPerson(exec boil.Executor, id int64) (*Person, *HttpError) {
	person, err := models.Persons(exec,
		qm.Where("id = ?", id),
		qm.Load("PersonI18ns")).
		One()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// i18n
	x := &Person{Person: *person}
	x.I18n = make(map[string]*models.PersonI18n, len(person.R.PersonI18ns))
	for _, i18n := range person.R.PersonI18ns {
		x.I18n[i18n.Language] = i18n
	}

	return x, nil
}

func handleUpdatePersonI18n(exec boil.Executor, id int64, i18ns []*models.PersonI18n) (*Person, *HttpError) {
	person, err := handleGetPerson(exec, id)
	if err != nil {
		return nil, err
	}

	// Upsert all new i18ns
	nI18n := make(map[string]*models.PersonI18n, len(i18ns))
	for _, i18n := range i18ns {
		i18n.PersonID = id
		nI18n[i18n.Language] = i18n
		err := i18n.Upsert(exec, true,
			[]string{"person_id", "language"},
			[]string{"name", "description"})
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	// Delete old i18ns not in new i18ns
	for k, v := range person.I18n {
		if _, ok := nI18n[k]; !ok {
			err := v.Delete(exec)
			if err != nil {
				return nil, NewInternalError(err)
			}
		}
	}

	return handleGetPerson(exec, id)
}

func handleDeletePerson(exec boil.Executor, id int64) (*models.Person, *HttpError) {
	person, err := models.FindPerson(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	err = models.ContentUnitsPersons(exec, qm.Where("person_id = ?", id)).DeleteAll()
	if err != nil {
		return nil, NewInternalError(err)
	}

	err = models.PersonI18ns(exec, qm.Where("person_id = ?", id)).DeleteAll()
	if err != nil {
		return nil, NewInternalError(err)
	}

	err = person.Delete(exec)
	if err != nil {
		return nil, NewInternalError(err)
	}

	return person, nil
}

func handleStoragesList(exec boil.Executor, r StoragesRequest) (*StoragesResponse, *HttpError) {
	mods := make([]qm.QueryMod, 0)

	// count query
	var total int64
	countMods := append([]qm.QueryMod{qm.Select("count(DISTINCT id)")}, mods...)
	err := models.Storages(exec, countMods...).QueryRow().Scan(&total)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if total == 0 {
		return NewStoragesResponse(), nil
	}

	// order, limit, offset
	if err = appendListMods(&mods, r.ListRequest); err != nil {
		return nil, NewBadRequestError(err)
	}

	// data query
	data, err := models.Storages(exec, mods...).All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	return &StoragesResponse{
		ListResponse: ListResponse{Total: total},
		Storages:     data,
	}, nil
}

func handlePublishersList(exec boil.Executor, r PublishersRequest) (*PublishersResponse, *HttpError) {
	mods := make([]qm.QueryMod, 0)

	// filters
	if err := appendIDsFilterMods(&mods, r.IDsFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendUIDsFilterMods(&mods, r.UIDsFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendPatternsFilterMods(&mods, r.PatternsFilter); err != nil {
		return nil, NewBadRequestError(err)
	}

	// count query
	var total int64
	countMods := append([]qm.QueryMod{qm.Select("count(DISTINCT id)")}, mods...)
	err := models.Publishers(exec, countMods...).QueryRow().Scan(&total)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if total == 0 {
		return NewPublishersResponse(), nil
	}

	// order, limit, offset
	if err = appendListMods(&mods, r.ListRequest); err != nil {
		return nil, NewBadRequestError(err)
	}

	// Eager loading
	mods = append(mods, qm.Load("PublisherI18ns"))

	// data query
	publishers, err := models.Publishers(exec, mods...).All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	// i18n
	data := make([]*Publisher, len(publishers))
	for i, pr := range publishers {
		x := &Publisher{Publisher: *pr}
		data[i] = x
		x.I18n = make(map[string]*models.PublisherI18n, len(pr.R.PublisherI18ns))
		for _, i18n := range pr.R.PublisherI18ns {
			x.I18n[i18n.Language] = i18n
		}
	}

	return &PublishersResponse{
		ListResponse: ListResponse{Total: total},
		Publishers:   data,
	}, nil
}

func handleCreatePublisher(exec boil.Executor, p *Publisher) (*Publisher, *HttpError) {

	// save publisher to DB
	uid, err := GetFreeUID(exec, new(PublisherUIDChecker))
	if err != nil {
		return nil, NewInternalError(err)
	}
	p.UID = uid

	err = p.Publisher.Insert(exec)
	if err != nil {
		return nil, NewInternalError(err)
	}

	// save i18n
	for _, v := range p.I18n {
		err := p.AddPublisherI18ns(exec, true, v)
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	return handleGetPublisher(exec, p.ID)
}

func handleGetPublisher(exec boil.Executor, id int64) (*Publisher, *HttpError) {
	publisher, err := models.Publishers(exec,
		qm.Where("id = ?", id),
		qm.Load("PublisherI18ns")).
		One()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// i18n
	x := &Publisher{Publisher: *publisher}
	x.I18n = make(map[string]*models.PublisherI18n, len(publisher.R.PublisherI18ns))
	for _, i18n := range publisher.R.PublisherI18ns {
		x.I18n[i18n.Language] = i18n
	}

	return x, nil
}

func handleUpdatePublisher(exec boil.Executor, p *Publisher) (*Publisher, *HttpError) {
	publisher, err := models.FindPublisher(exec, p.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	publisher.Pattern = p.Pattern
	err = p.Update(exec, "pattern")
	if err != nil {
		return nil, NewInternalError(err)
	}

	return handleGetPublisher(exec, p.ID)
}

func handleUpdatePublisherI18n(exec boil.Executor, id int64, i18ns []*models.PublisherI18n) (*Publisher, *HttpError) {
	publisher, err := handleGetPublisher(exec, id)
	if err != nil {
		return nil, err
	}

	// Upsert all new i18ns
	nI18n := make(map[string]*models.PublisherI18n, len(i18ns))
	for _, i18n := range i18ns {
		i18n.PublisherID = id
		nI18n[i18n.Language] = i18n
		err := i18n.Upsert(exec, true,
			[]string{"publisher_id", "language"},
			[]string{"name", "description"})
		if err != nil {
			return nil, NewInternalError(err)
		}
	}

	// Delete old i18ns not in new i18ns
	for k, v := range publisher.I18n {
		if _, ok := nI18n[k]; !ok {
			err := v.Delete(exec)
			if err != nil {
				return nil, NewInternalError(err)
			}
		}
	}

	return handleGetPublisher(exec, id)
}

// Query Helpers

func appendListMods(mods *[]qm.QueryMod, r ListRequest) error {

	// group by id to remove duplicates
	*mods = append(*mods, qm.GroupBy("id"))

	if r.OrderBy == "" {
		*mods = append(*mods, qm.OrderBy("id desc"))
	} else {
		*mods = append(*mods, qm.OrderBy(r.OrderBy))
	}

	var limit, offset int

	if r.StartIndex == 0 {
		// pagination style
		if r.PageSize == 0 {
			limit = DEFAULT_PAGE_SIZE
		} else {
			limit = utils.Min(r.PageSize, MAX_PAGE_SIZE)
		}
		if r.PageNumber > 1 {
			offset = (r.PageNumber - 1) * limit
		}
	} else {
		// start & stop index style for "infinite" lists
		offset = r.StartIndex - 1
		if r.StopIndex == 0 {
			limit = MAX_PAGE_SIZE
		} else if r.StopIndex < r.StartIndex {
			return errors.Errorf("Invalid range [%d-%d]", r.StartIndex, r.StopIndex)
		} else {
			limit = r.StopIndex - r.StartIndex + 1
		}
	}

	*mods = append(*mods, qm.Limit(limit))
	if offset != 0 {
		*mods = append(*mods, qm.Offset(offset))
	}

	return nil
}

func appendPermissionsMods(cp utils.ContextProvider, mods *[]qm.QueryMod) {
	*mods = append(*mods, qm.Where("secure <= ?", allowedRead(cp)))
}

func appendSearchTermFilterMods(exec boil.Executor, mods *[]qm.QueryMod, f SearchTermFilter, entityType int) error {
	if f.Query == "" {
		return nil
	}

	var whereParts []string

	// id field - must be unsigned int
	if id, err := strconv.ParseUint(f.Query, 10, 64); err == nil {
		whereParts = append(whereParts, fmt.Sprintf("id = %d", id))
	}

	// uid field
	if len(f.Query) == 8 {
		whereParts = append(whereParts, fmt.Sprintf("uid = '%s'", f.Query))
	}

	switch entityType {
	case SEARCH_IN_FILES:
		// file name field
		whereParts = append(whereParts, fmt.Sprintf("name ~~ '%%%s%%'", f.Query))

		// file sha1
		if len(f.Query) == 40 {
			_, err := hex.DecodeString(f.Query) // make sure it's a hex string
			if err == nil {
				whereParts = append(whereParts, fmt.Sprintf("sha1 = '\\x%s'", f.Query))
			}
		}
	case SEARCH_IN_CONTENT_UNITS:

		// get CU IDs from search in i18ns
		var ids pq.Int64Array
		q := `select array_agg(cui.content_unit_id)
			  from content_units cu 
			  left join content_unit_i18n cui on cu.id = cui.content_unit_id
			  where cui.name ~ $1 or cui .description ~ $1 limit $2`
		err := queries.Raw(exec, q, f.Query, MAX_PAGE_SIZE).QueryRow().Scan(&ids)
		if err != nil {
			return err
		}

		if ids != nil && len(ids) != 0 {
			intListStr := strings.Trim(strings.Replace(fmt.Sprint(ids), " ", ",", -1), "[]")
			whereParts = append(whereParts, fmt.Sprintf("id in (%s)", intListStr))
		}

	case SEARCH_IN_COLLECTIONS:

		// get Collection IDs from search in i18ns
		var ids pq.Int64Array
		q := `select array_agg(c.id)
				from collections c 
				left join collection_i18n ci on c.id = ci.collection_id
				where ci.name ~ $1 or ci.description ~ $1 limit $2`
		err := queries.Raw(exec, q, f.Query, MAX_PAGE_SIZE).QueryRow().Scan(&ids)
		if err != nil {
			return err
		}

		if ids != nil && len(ids) != 0 {
			intListStr := strings.Trim(strings.Replace(fmt.Sprint(ids), " ", ",", -1), "[]")
			whereParts = append(whereParts, fmt.Sprintf("id in (%s)", intListStr))
		}
	}

	if len(whereParts) > 0 {
		whereQuery := fmt.Sprintf("(%s)", strings.Join(whereParts, " or "))
		*mods = append(*mods, qm.And(whereQuery))
	} else {
		*mods = append(*mods, qm.Where("id < 0")) // so we get back empty results
	}

	return nil
}

func appendIDsFilterMods(mods *[]qm.QueryMod, f IDsFilter) error {
	if len(f.IDs) == 0 {
		return nil
	}

	*mods = append(*mods, qm.WhereIn("id IN ?", utils.ConvertArgsInt64(f.IDs)...))

	return nil
}

func appendPatternsFilterMods(mods *[]qm.QueryMod, f PatternsFilter) error {
	if utils.IsEmpty(f.Patterns) {
		return nil
	}

	*mods = append(*mods, qm.WhereIn("pattern IN ?", utils.ConvertArgsString(f.Patterns)...))

	return nil
}

func appendUIDsFilterMods(mods *[]qm.QueryMod, f UIDsFilter) error {
	if utils.IsEmpty(f.UIDs) {
		return nil
	}

	*mods = append(*mods, qm.WhereIn("uid IN ?", utils.ConvertArgsString(f.UIDs)...))

	return nil
}

func appendSHA1sFilterMods(mods *[]qm.QueryMod, f SHA1sFilter) error {
	if utils.IsEmpty(f.SHA1s) {
		return nil
	}

	hexSHA1s := make([][]byte, 0)
	for i := range f.SHA1s {
		s, err := hex.DecodeString(f.SHA1s[i])
		if err != nil {
			return errors.Wrapf(err, "hex.DecodeString [%d]: %s", i, f.SHA1s[i])
		}
		hexSHA1s = append(hexSHA1s, s)
	}

	*mods = append(*mods, qm.WhereIn("sha1 IN ?", utils.ConvertArgsBytes(hexSHA1s)...))

	return nil
}

func appendDateRangeFilterMods(mods *[]qm.QueryMod, f DateRangeFilter, field string) error {
	s, e, err := f.Range()
	if err != nil {
		return err
	}

	if f.StartDate != "" && f.EndDate != "" && e.Before(s) {
		return errors.New("Invalid date range")
	}

	if field == "" {
		field = "created_at"
	}

	if f.StartDate != "" {
		*mods = append(*mods, qm.Where(fmt.Sprintf("%s >= ?", field), s))
	}
	if f.EndDate != "" {
		*mods = append(*mods, qm.Where(fmt.Sprintf("%s <= ?", field), e))
	}

	return nil
}

func appendContentTypesFilterMods(mods *[]qm.QueryMod, f ContentTypesFilter) error {
	if utils.IsEmpty(f.ContentTypes) {
		return nil
	}

	a := make([]interface{}, len(f.ContentTypes))
	for i, x := range f.ContentTypes {
		ct, ok := common.CONTENT_TYPE_REGISTRY.ByName[strings.ToUpper(x)]
		if ok {
			a[i] = ct.ID
		} else {
			return errors.Errorf("Unknown content type: %s", x)
		}
	}

	*mods = append(*mods, qm.WhereIn("type_id in ?", a...))

	return nil
}

func appendSecureFilterMods(mods *[]qm.QueryMod, f SecureFilter) error {
	if len(f.Levels) == 0 {
		return nil
	}

	a := make([]interface{}, len(f.Levels))
	for i, x := range f.Levels {
		if x == common.SEC_PUBLIC || x == common.SEC_SENSITIVE || x == common.SEC_PRIVATE {
			a[i] = x
		} else {
			return errors.Errorf("Unknown security level: %d", x)
		}
	}

	*mods = append(*mods, qm.WhereIn("secure in ?", a...))

	return nil
}

func appendPublishedFilterMods(mods *[]qm.QueryMod, f PublishedFilter) {
	var val null.Bool
	val.UnmarshalText([]byte(f.Published))
	if val.Valid {
		*mods = append(*mods, qm.Where("published = ?", val.Bool))
	}
}

func appendOperationTypesFilterMods(mods *[]qm.QueryMod, f OperationTypesFilter) error {
	if utils.IsEmpty(f.OperationTypes) {
		return nil
	}

	a := make([]interface{}, len(f.OperationTypes))
	for i, x := range f.OperationTypes {
		ot, ok := common.OPERATION_TYPE_REGISTRY.ByName[strings.ToLower(x)]
		if ok {
			a[i] = ot.ID
		} else {
			return errors.Errorf("Unknown operation type: %s", x)
		}
	}

	*mods = append(*mods, qm.WhereIn("type_id in ?", a...))

	return nil
}

func appendSourcesFilterMods(exec boil.Executor, mods *[]qm.QueryMod, f SourcesFilter) error {
	if utils.IsEmpty(f.Authors) && len(f.Sources) == 0 {
		return nil
	}

	// slice of all source ids we want
	source_ids := make([]int64, 0)

	// fetch source ids by authors
	if !utils.IsEmpty(f.Authors) {
		for _, x := range f.Authors {
			if _, ok := common.AUTHOR_REGISTRY.ByCode[strings.ToLower(x)]; !ok {
				return NewBadRequestError(errors.Errorf("Unknown author: %s", x))
			}
		}

		var ids pq.Int64Array
		q := `SELECT array_agg(DISTINCT "as".source_id)
		      FROM authors a INNER JOIN authors_sources "as" ON a.id = "as".author_id
		      WHERE a.code = ANY($1)`
		err := queries.Raw(exec, q, pq.Array(f.Authors)).QueryRow().Scan(&ids)
		if err != nil {
			return err
		}
		source_ids = append(source_ids, ids...)
	}

	// blend in requested sources
	source_ids = append(source_ids, f.Sources...)

	// find all nested source_ids
	q := `WITH RECURSIVE rec_sources AS (
		  SELECT s.id FROM sources s WHERE s.id = ANY($1)
		  UNION
		  SELECT s.id FROM sources s INNER JOIN rec_sources rs ON s.parent_id = rs.id
	      )
	      SELECT array_agg(distinct id) FROM rec_sources`
	var ids pq.Int64Array
	err := queries.Raw(exec, q, pq.Array(source_ids)).QueryRow().Scan(&ids)
	if err != nil {
		return err
	}

	if ids == nil || len(ids) == 0 {
		*mods = append(*mods, qm.Where("id < 0")) // so results would be empty
	} else {
		*mods = append(*mods,
			qm.InnerJoin("content_units_sources cus ON id = cus.content_unit_id"),
			qm.WhereIn("cus.source_id in ?", utils.ConvertArgsInt64(ids)...))
	}

	return nil
}

func appendTagsFilterMods(exec boil.Executor, mods *[]qm.QueryMod, f TagsFilter) error {
	if len(f.Tags) == 0 {
		return nil
	}

	// find all nested tag_ids
	q := `WITH RECURSIVE rec_tags AS (
	        SELECT t.id FROM tags t WHERE t.id = ANY($1)
	        UNION
	        SELECT t.id FROM tags t INNER JOIN rec_tags rt ON t.parent_id = rt.id
	      )
	      SELECT array_agg(distinct id) FROM rec_tags`
	var ids pq.Int64Array
	err := queries.Raw(exec, q, pq.Array(f.Tags)).QueryRow().Scan(&ids)
	if err != nil {
		return err
	}

	if ids == nil || len(ids) == 0 {
		*mods = append(*mods, qm.Where("id < 0")) // so results would be empty
	} else {
		*mods = append(*mods,
			qm.InnerJoin("content_units_tags cut ON id = cut.content_unit_id"),
			qm.WhereIn("cut.tag_id in ?", utils.ConvertArgsInt64(ids)...))
	}

	return nil
}

// mustBeginTx begins a transaction, panics on error.
func mustBeginTx(c utils.ContextProvider) *sql.Tx {
	tx, err := c.MustGet("MDB").(*sql.DB).Begin()
	utils.Must(err)
	return tx
}

// mustConcludeTx commits or rollback the given transaction according to given error.
// Panics if Commit() or Rollback() fails.
func mustConcludeTx(tx *sql.Tx, err *HttpError) {
	if err == nil {
		utils.Must(tx.Commit())
	} else {
		utils.Must(tx.Rollback())
	}
}

func mdbReplicationLocation(c utils.ContextProvider) (string, error) {
	var loc string
	err := c.MustGet("MDB").(*sql.DB).
		QueryRow("SELECT pg_current_xlog_insert_location();").
		Scan(&loc)
	if err != nil {
		return "", errors.Wrap(err, "Fetch Replication Position")
	}
	return loc, nil
}

// concludeRequest responds with JSON of given response or aborts the request with the given error.
func concludeRequest(c *gin.Context, resp interface{}, err *HttpError) {
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

func emitEvents(cp utils.ContextProvider, evnts ...events.Event) {
	if len(evnts) == 0 {
		return
	}

	// We attach postgresql replication log location
	// so that clients could verify that their stand-bys are synced
	// see:
	// https://blog.2ndquadrant.com/postgresql-10-transaction-traceability/
	// https://www.postgresql.org/docs/9.6/static/functions-admin.html
	// https://www.postgresql.org/docs/9.6/static/datatype-pg-lsn.html
	if rLoc, err := mdbReplicationLocation(cp); err != nil {
		log.Errorf("emitEvents: rLoc: %+v", err)
	} else {
		for i := range evnts {
			evnts[i].ReplicationLocation = rLoc
		}
	}

	cp.MustGet("EVENTS_EMITTER").(events.EventEmitter).Emit(evnts...)
}

func can(cp utils.ContextProvider, obj string, act string) bool {
	sub := []string{""}
	if v, ok := cp.Get("ID_TOKEN_CLAIMS"); ok {
		claims := v.(permissions.IDTokenClaims)
		sub = claims.RealmAccess.Roles
		//log.Infof("Subject is %s %s with roles %v", claims.Sub, claims.Name, sub)
	} else {
		// bypass hack for workflow insert station
		if act == common.PERM_READ && cp.(*gin.Context).ClientIP() == "146.185.60.45" {
			log.Info("Workflow Insert station read")
			return true
		}

		log.Infof("No subject.")
	}

	enforcer := cp.MustGet("PERMISSIONS_ENFORCER").(*casbin.Enforcer)

	for i := range sub {
		if enforcer.Enforce(sub[i], obj, act) {
			log.Infof("ALLOW %s, %s, %s", sub[i], obj, act)
			return true
		}
	}

	//log.Warnf("DENY %v, %s, %s", sub, obj, act)
	return false
}

func isAdmin(cp utils.ContextProvider) bool {
	if v, ok := cp.Get("ID_TOKEN_CLAIMS"); ok {
		claims := v.(permissions.IDTokenClaims)
		for i := range claims.RealmAccess.Roles {
			if "archive_admin" == claims.RealmAccess.Roles[i] {
				return true
			}
		}
	}
	return false
}

func secureToPermission(secure int16) string {
	switch secure {
	case common.SEC_PRIVATE:
		return "data_private"
	case common.SEC_SENSITIVE:
		return "data_sensitive"
	default:
		return "data_public"
	}
}

func allowedRead(cp utils.ContextProvider) int16 {
	return allowedSecure(cp, common.PERM_READ)
}

func allowedWrite(cp utils.ContextProvider) int16 {
	return allowedSecure(cp, common.PERM_WRITE)
}

func allowedSecure(cp utils.ContextProvider, act string) int16 {
	if can(cp, secureToPermission(common.SEC_PRIVATE), act) {
		return common.SEC_PRIVATE
	} else if can(cp, secureToPermission(common.SEC_SENSITIVE), act) {
		return common.SEC_SENSITIVE
	} else if can(cp, secureToPermission(common.SEC_PUBLIC), act) {
		return common.SEC_PUBLIC
	}

	if ginCtx, ok := cp.(*gin.Context); ok {
		clientIP := ginCtx.ClientIP()

		// workflow insert station
		if clientIP == "146.185.60.45" {
			//log.Info("Workflow Insert station secure level")
			return common.SEC_PRIVATE
		}

		// internal network (hopefully MDB-CIT [aka rename])
		if strings.HasPrefix(clientIP, "10.") {
			//log.Infof("Internal network secure level: %s", clientIP)
			return common.SEC_PRIVATE
		}
	}

	return -1
}
