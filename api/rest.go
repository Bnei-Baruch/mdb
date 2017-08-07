package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/nullbio/null.v6"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

const (
	DEFAULT_PAGE_SIZE = 50
	MAX_PAGE_SIZE     = 1000
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

		resp, err = handleCollectionsList(boil.GetDB(), r)
	case http.MethodPost:
		var collection Collection
		if c.BindJSON(&collection) != nil {
			return
		}

		if _, ok := CONTENT_TYPE_REGISTRY.ByID[collection.TypeID]; !ok {
			err := errors.Errorf("Unknown content type %d", collection.TypeID)
			NewBadRequestError(err).Abort(c)
			return
		}

		for _, x := range collection.I18n {
			if StdLang(x.Language) == LANG_UNKNOWN {
				err := errors.Errorf("Unknown language %s", x.Language)
				NewBadRequestError(err).Abort(c)
				return
			}
		}

		tx := mustBeginTx()
		resp, err = handleCreateCollection(tx, collection)
		mustConcludeTx(tx, err)
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
		resp, err = handleGetCollection(boil.GetDB(), id)
	case http.MethodPut:
		var cl Collection
		if c.Bind(&cl) != nil {
			return
		}

		cl.ID = id
		tx := mustBeginTx()
		resp, err = handleUpdateCollection(tx, &cl)
		mustConcludeTx(tx, err)
	case http.MethodDelete:
		tx := mustBeginTx()
		err = handleDeleteCollection(tx, id)
		mustConcludeTx(tx, err)
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
		if StdLang(x.Language) == LANG_UNKNOWN {
			NewBadRequestError(errors.Errorf("Unknown language %s", x.Language)).Abort(c)
			return
		}
	}

	tx := mustBeginTx()
	resp, err := handleUpdateCollectionI18n(tx, id, i18ns)
	mustConcludeTx(tx, err)
	concludeRequest(c, resp, err)
}

func CollectionContentUnitsHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleCollectionCCU(boil.GetDB(), id)
	concludeRequest(c, resp, err)
}

// Toggle the active flag of a single container
func CollectionActivateHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleCollectionActivate(boil.GetDB(), id)
	concludeRequest(c, resp, err)
}

func ContentUnitsListHandler(c *gin.Context) {
	var r ContentUnitsRequest
	if c.Bind(&r) != nil {
		return
	}

	resp, err := handleContentUnitsList(boil.GetDB(), r)
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
		resp, err = handleGetContentUnit(boil.GetDB(), id)
	} else {
		if c.Request.Method == http.MethodPut {
			var cu ContentUnit
			if c.Bind(&cu) != nil {
				return
			}

			cu.ID = id
			tx := mustBeginTx()
			resp, err = handleUpdateContentUnit(tx, &cu)
			mustConcludeTx(tx, err)
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
		if StdLang(x.Language) == LANG_UNKNOWN {
			NewBadRequestError(errors.Errorf("Unknown language %s", x.Language)).Abort(c)
			return
		}
	}

	tx := mustBeginTx()
	resp, err := handleUpdateContentUnitI18n(tx, id, i18ns)
	mustConcludeTx(tx, err)
	concludeRequest(c, resp, err)
}

func ContentUnitFilesHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleContentUnitFiles(boil.GetDB(), id)
	concludeRequest(c, resp, err)
}

func ContentUnitCollectionsHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleContentUnitCCU(boil.GetDB(), id)
	concludeRequest(c, resp, err)
}

func ContentUnitSourcesHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var err *HttpError
	var resp interface{}

	switch c.Request.Method {
	case http.MethodGet, "":
		resp, err = handleGetContentUnitSources(boil.GetDB(), id)
	case http.MethodPost:
		var body map[string]int64
		if c.BindJSON(&body) != nil {
			return
		}

		sourceID, ok := body["sourceID"]
		if !ok {
			err = NewBadRequestError(errors.Wrap(e, "No sourceID given"))
			break
		}

		tx := mustBeginTx()
		err = handleContentUnitAddSource(tx, id, sourceID)
		mustConcludeTx(tx, err)
	case http.MethodDelete:
		sourceID, e := strconv.ParseInt(c.Param("sourceID"), 10, 0)
		if e != nil {
			err = NewBadRequestError(errors.Wrap(e, "sourceID expects int64"))
			break
		}

		tx := mustBeginTx()
		err = handleContentUnitRemoveSource(tx, id, sourceID)
		mustConcludeTx(tx, err)
	}

	concludeRequest(c, resp, err)
}

func ContentUnitTagsHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	var err *HttpError
	var resp interface{}

	switch c.Request.Method {
	case http.MethodGet, "":
		resp, err = handleGetContentUnitTags(boil.GetDB(), id)
	case http.MethodPost:
		var body map[string]int64
		if c.BindJSON(&body) != nil {
			return
		}

		tagID, ok := body["tagID"]
		if !ok {
			err = NewBadRequestError(errors.Wrap(e, "No tagID given"))
			break
		}

		tx := mustBeginTx()
		err = handleContentUnitAddTag(tx, id, tagID)
		mustConcludeTx(tx, err)
	case http.MethodDelete:
		tagID, e := strconv.ParseInt(c.Param("tagID"), 10, 0)
		if e != nil {
			err = NewBadRequestError(errors.Wrap(e, "tagID expects int64"))
			break
		}

		tx := mustBeginTx()
		err = handleContentUnitRemoveTag(tx, id, tagID)
		mustConcludeTx(tx, err)
	}

	concludeRequest(c, resp, err)
}

func FilesListHandler(c *gin.Context) {
	var r FilesRequest
	if c.Bind(&r) != nil {
		return
	}

	resp, err := handleFilesList(boil.GetDB(), r)
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
		resp, err = handleGetFile(boil.GetDB(), id)
	} else {
		if c.Request.Method == http.MethodPut {
			var f MFile
			if c.Bind(&f) != nil {
				return
			}

			f.ID = id
			tx := mustBeginTx()
			resp, err = handleUpdateFile(tx, &f)
			mustConcludeTx(tx, err)
		}
	}

	concludeRequest(c, resp, err)
}

func OperationsListHandler(c *gin.Context) {
	var r OperationsRequest
	if c.Bind(&r) != nil {
		return
	}

	resp, err := handleOperationsList(boil.GetDB(), r)
	concludeRequest(c, resp, err)
}

func OperationItemHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleOperationItem(boil.GetDB(), id)
	concludeRequest(c, resp, err)
}

func OperationFilesHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleOperationFiles(boil.GetDB(), id)
	concludeRequest(c, resp, err)
}

func AuthorsHandler(c *gin.Context) {
	authors, err := models.Authors(boil.GetDB(),
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
		var r SourcesRequest
		if c.Bind(&r) != nil {
			return
		}
		resp, err = handleGetSources(boil.GetDB(), r)
	} else {
		if c.Request.Method == http.MethodPost {
			var r CreateSourceRequest
			if c.Bind(&r) != nil {
				return
			}

			if _, ok := SOURCE_TYPE_REGISTRY.ByID[r.Source.TypeID]; !ok {
				err := errors.Errorf("Unknown source type %d", r.Source.TypeID)
				NewBadRequestError(err).Abort(c)
				return
			}

			for _, x := range r.Source.I18n {
				if StdLang(x.Language) == LANG_UNKNOWN {
					err := errors.Errorf("Unknown language %s", x.Language)
					NewBadRequestError(err).Abort(c)
					return
				}
			}

			tx := mustBeginTx()
			resp, err = handleCreateSource(tx, r)
			mustConcludeTx(tx, err)
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
		resp, err = handleGetSource(boil.GetDB(), id)
	} else {
		if c.Request.Method == http.MethodPut {
			var s Source
			if c.Bind(&s) != nil {
				return
			}

			s.ID = id
			tx := mustBeginTx()
			resp, err = handleUpdateSource(tx, &s)
			mustConcludeTx(tx, err)
		}
	}

	concludeRequest(c, resp, err)
}

func SourceI18nHandler(c *gin.Context) {
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
		if StdLang(x.Language) == LANG_UNKNOWN {
			NewBadRequestError(errors.Errorf("Unknown language %s", x.Language)).Abort(c)
			return
		}
	}

	tx := mustBeginTx()
	resp, err := handleUpdateSourceI18n(tx, id, i18ns)
	mustConcludeTx(tx, err)
	concludeRequest(c, resp, err)
}

func TagsHandler(c *gin.Context) {
	var err *HttpError
	var resp interface{}

	if c.Request.Method == http.MethodGet || c.Request.Method == "" {
		var r TagsRequest
		if c.Bind(&r) != nil {
			return
		}
		resp, err = handleGetTags(boil.GetDB(), r)
	} else {
		if c.Request.Method == http.MethodPost {
			var t Tag
			if c.Bind(&t) != nil {
				return
			}

			for _, x := range t.I18n {
				if StdLang(x.Language) == LANG_UNKNOWN {
					NewBadRequestError(errors.Errorf("Unknown language %s", x.Language)).Abort(c)
					return
				}
			}

			tx := mustBeginTx()
			resp, err = handleCreateTag(tx, &t)
			mustConcludeTx(tx, err)
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
		resp, err = handleGetTag(boil.GetDB(), id)
	} else {
		if c.Request.Method == http.MethodPut {
			var t Tag
			if c.Bind(&t) != nil {
				return
			}

			t.ID = id
			tx := mustBeginTx()
			resp, err = handleUpdateTag(tx, &t)
			mustConcludeTx(tx, err)
		}
	}

	concludeRequest(c, resp, err)
}

func TagI18nHandler(c *gin.Context) {
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
		if StdLang(x.Language) == LANG_UNKNOWN {
			NewBadRequestError(errors.Errorf("Unknown language %s", x.Language)).Abort(c)
			return
		}
	}

	tx := mustBeginTx()
	resp, err := handleUpdateTagI18n(tx, id, i18ns)
	mustConcludeTx(tx, err)
	concludeRequest(c, resp, err)
}

// Handlers Logic

func handleCollectionsList(exec boil.Executor, r CollectionsRequest) (*CollectionsResponse, *HttpError) {
	mods := make([]qm.QueryMod, 0)

	// filters
	if err := appendContentTypesFilterMods(&mods, r.ContentTypesFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendDateRangeFilterMods(&mods, r.DateRangeFilter, "(properties->>'film_date')::date"); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendSecureFilterMods(&mods, r.SecureFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	appendPublishedFilterMods(&mods, r.PublishedFilter)

	// count query
	total, err := models.Collections(exec, mods...).Count()
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

func handleCreateCollection(exec boil.Executor, c Collection) (*Collection, *HttpError) {
	// unmarshal properties
	props := make(map[string]interface{})
	if c.Properties.Valid {
		err := json.Unmarshal(c.Properties.JSON, &props)
		if err != nil {
			return nil, NewBadRequestError(errors.Wrap(err, "json.Unmarshal properties"))
		}
	}

	// create collection in DB
	ct := CONTENT_TYPE_REGISTRY.ByID[c.TypeID].Name
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

	return handleGetCollection(exec, collection.ID)
}

func handleGetCollection(exec boil.Executor, id int64) (*Collection, *HttpError) {
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

	// i18n
	x := &Collection{Collection: *collection}
	x.I18n = make(map[string]*models.CollectionI18n, len(collection.R.CollectionI18ns))
	for _, i18n := range collection.R.CollectionI18ns {
		x.I18n[i18n.Language] = i18n
	}

	return x, nil
}

func handleUpdateCollection(exec boil.Executor, c *Collection) (*Collection, *HttpError) {
	collection, err := models.FindCollection(exec, c.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	// update entity attributes
	collection.Secure = c.Secure
	err = collection.Update(exec, "secure")
	if err != nil {
		return nil, NewInternalError(err)
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

	return handleGetCollection(exec, c.ID)
}

func handleDeleteCollection(exec boil.Executor, id int64) *HttpError {
	collection, err := models.FindCollection(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewNotFoundError()
		} else {
			return NewInternalError(err)
		}
	}

	err = models.CollectionsContentUnits(exec, qm.Where("collection_id = ?", id)).DeleteAll()
	if err != nil {
		return NewInternalError(err)
	}

	err = models.CollectionI18ns(exec, qm.Where("collection_id = ?", id)).DeleteAll()
	if err != nil {
		return NewInternalError(err)
	}

	err = collection.Delete(exec)
	if err != nil {
		return NewInternalError(err)
	}

	return nil
}

func handleUpdateCollectionI18n(exec boil.Executor, id int64, i18ns []*models.CollectionI18n) (*Collection, *HttpError) {
	collection, err := handleGetCollection(exec, id)
	if err != nil {
		return nil, err
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

	return handleGetCollection(exec, id)
}

func handleCollectionActivate(exec boil.Executor, id int64) (*Collection, *HttpError) {
	collection, err := models.FindCollection(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
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

	return handleGetCollection(exec, id)
}

func handleCollectionCCU(exec boil.Executor, id int64) ([]*CollectionContentUnit, *HttpError) {
	ok, err := models.CollectionExists(exec, id)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if !ok {
		return nil, NewNotFoundError()
	}

	ccus, err := models.CollectionsContentUnits(exec, qm.Where("collection_id = ?", id)).All()
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

	data := make([]*CollectionContentUnit, len(ccus))
	for i, ccu := range ccus {
		data[i] = &CollectionContentUnit{
			Name:        ccu.Name,
			ContentUnit: cusById[ccu.ContentUnitID],
		}
	}

	return data, nil
}

func handleContentUnitsList(exec boil.Executor, r ContentUnitsRequest) (*ContentUnitsResponse, *HttpError) {
	mods := make([]qm.QueryMod, 0)

	// filters
	if err := appendContentTypesFilterMods(&mods, r.ContentTypesFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendDateRangeFilterMods(&mods, r.DateRangeFilter, "(properties->>'film_date')::date"); err != nil {
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
	if err := appendSecureFilterMods(&mods, r.SecureFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	appendPublishedFilterMods(&mods, r.PublishedFilter)

	// count query
	total, err := models.ContentUnits(exec, mods...).Count()
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
		x.I18n = make(map[string]*models.ContentUnitI18n, len(cu.R.ContentUnitI18ns))
		for _, i18n := range cu.R.ContentUnitI18ns {
			x.I18n[i18n.Language] = i18n
		}
	}

	return &ContentUnitsResponse{
		ListResponse: ListResponse{Total: total},
		ContentUnits: data,
	}, nil
}

func handleGetContentUnit(exec boil.Executor, id int64) (*ContentUnit, *HttpError) {
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

	// i18n
	x := &ContentUnit{ContentUnit: *unit}
	x.I18n = make(map[string]*models.ContentUnitI18n, len(unit.R.ContentUnitI18ns))
	for _, i18n := range unit.R.ContentUnitI18ns {
		x.I18n[i18n.Language] = i18n
	}

	return x, nil
}

func handleUpdateContentUnit(exec boil.Executor, cu *ContentUnit) (*ContentUnit, *HttpError) {
	unit, err := models.FindContentUnit(exec, cu.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	unit.Secure = cu.Secure
	err = unit.Update(exec, "secure")
	if err != nil {
		return nil, NewInternalError(err)
	}

	return handleGetContentUnit(exec, cu.ID)
}

func handleUpdateContentUnitI18n(exec boil.Executor, id int64, i18ns []*models.ContentUnitI18n) (*ContentUnit, *HttpError) {
	unit, err := handleGetContentUnit(exec, id)
	if err != nil {
		return nil, err
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

	return handleGetContentUnit(exec, id)
}

func handleContentUnitFiles(exec boil.Executor, id int64) ([]*MFile, *HttpError) {
	ok, err := models.ContentUnitExists(exec, id)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if !ok {
		return nil, NewNotFoundError()
	}

	files, err := models.Files(exec, qm.Where("content_unit_id = ?", id)).All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	data := make([]*MFile, len(files))
	for i, f := range files {
		data[i] = NewMFile(f)
	}

	return data, nil
}

func handleContentUnitCCU(exec boil.Executor, id int64) ([]*CollectionContentUnit, *HttpError) {
	ok, err := models.ContentUnitExists(exec, id)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if !ok {
		return nil, NewNotFoundError()
	}

	ccus, err := models.CollectionsContentUnits(exec, qm.Where("content_unit_id = ?", id)).All()
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
			Collection: csById[ccu.CollectionID],
		}
	}

	return data, nil
}

func handleGetContentUnitSources(exec boil.Executor, id int64) ([]*Source, *HttpError) {
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

func handleContentUnitAddSource(exec boil.Executor, id int64, sourceID int64) *HttpError {
	unit, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewNotFoundError()
		} else {
			return NewInternalError(err)
		}
	}

	source, err := models.FindSource(exec, sourceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewBadRequestError(errors.Errorf("Unknown source id %d", sourceID))
		} else {
			return NewInternalError(err)
		}
	}

	var count int64
	err = queries.Raw(exec,
		`SELECT COUNT(1) FROM content_units_sources WHERE content_unit_id=$1 AND source_id=$2`,
		unit.ID, source.ID).
		QueryRow().
		Scan(&count)
	if err != nil {
		return NewInternalError(err)
	}
	if count > 0 {
		return nil // noop
	}

	err = unit.AddSources(exec, false, source)
	if err != nil {
		return NewInternalError(err)
	}

	return nil
}

func handleContentUnitRemoveSource(exec boil.Executor, id int64, sourceID int64) *HttpError {
	unit, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewNotFoundError()
		} else {
			return NewInternalError(err)
		}
	}

	source, err := models.FindSource(exec, sourceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewBadRequestError(errors.Errorf("Unknown source id %d", sourceID))
		} else {
			return NewInternalError(err)
		}
	}

	err = unit.RemoveSources(exec, source)
	if err != nil {
		return NewInternalError(err)
	}

	return nil
}

func handleGetContentUnitTags(exec boil.Executor, id int64) ([]*Tag, *HttpError) {
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

func handleContentUnitAddTag(exec boil.Executor, id int64, tagID int64) *HttpError {
	unit, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewNotFoundError()
		} else {
			return NewInternalError(err)
		}
	}

	tag, err := models.FindTag(exec, tagID)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewBadRequestError(errors.Errorf("Unknown tag id %d", tagID))
		} else {
			return NewInternalError(err)
		}
	}

	var count int64
	err = queries.Raw(exec,
		`SELECT COUNT(1) FROM content_units_tags WHERE content_unit_id=$1 AND tag_id=$2`,
		unit.ID, tag.ID).
		QueryRow().
		Scan(&count)
	if err != nil {
		return NewInternalError(err)
	}
	if count > 0 {
		return nil // noop
	}

	err = unit.AddTags(exec, false, tag)
	if err != nil {
		return NewInternalError(err)
	}

	return nil
}

func handleContentUnitRemoveTag(exec boil.Executor, id int64, tagID int64) *HttpError {
	unit, err := models.FindContentUnit(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewNotFoundError()
		} else {
			return NewInternalError(err)
		}
	}

	tag, err := models.FindTag(exec, tagID)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewBadRequestError(errors.Errorf("Unknown tag id %d", tagID))
		} else {
			return NewInternalError(err)
		}
	}

	err = unit.RemoveTags(exec, tag)
	if err != nil {
		return NewInternalError(err)
	}

	return nil
}

func handleFilesList(exec boil.Executor, r FilesRequest) (*FilesResponse, *HttpError) {

	mods := make([]qm.QueryMod, 0)

	// filters
	if err := appendDateRangeFilterMods(&mods, r.DateRangeFilter, "file_created_at"); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendSecureFilterMods(&mods, r.SecureFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	appendPublishedFilterMods(&mods, r.PublishedFilter)
	if r.Query != "" {
		mods = append(mods, qm.Where("name ~ ?", r.Query),
			qm.Or("uid ~ ?", r.Query),
			qm.Or("id::TEXT ~ ?", r.Query))
	}

	// count query
	total, err := models.Files(exec, mods...).Count()
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

func handleGetFile(exec boil.Executor, id int64) (*MFile, *HttpError) {
	file, err := models.FindFile(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	return NewMFile(file), nil
}

func handleUpdateFile(exec boil.Executor, f *MFile) (*MFile, *HttpError) {
	file, err := models.FindFile(exec, f.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError()
		} else {
			return nil, NewInternalError(err)
		}
	}

	file.Secure = f.Secure
	err = file.Update(exec, "secure")
	if err != nil {
		return nil, NewInternalError(err)
	}

	return handleGetFile(exec, f.ID)
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
	total, err := models.Operations(exec, mods...).Count()
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

func handleOperationFiles(exec boil.Executor, id int64) ([]*MFile, *HttpError) {
	ok, err := models.OperationExists(exec, id)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if !ok {
		return nil, NewNotFoundError()
	}

	files, err := models.Files(exec,
		qm.InnerJoin("files_operations fo on fo.file_id=id and fo.operation_id = ?", id)).
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
	total, err := models.Sources(exec, mods...).Count()
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
	total, err := models.Tags(exec, mods...).Count()
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

// Query Helpers

func appendListMods(mods *[]qm.QueryMod, r ListRequest) error {
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
		ct, ok := CONTENT_TYPE_REGISTRY.ByName[strings.ToUpper(x)]
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
		if x == SEC_PUBLIC || x == SEC_SENSITIVE || x == SEC_PRIVATE {
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
		ot, ok := OPERATION_TYPE_REGISTRY.ByName[strings.ToLower(x)]
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
			if _, ok := AUTHOR_REGISTRY.ByCode[strings.ToLower(x)]; !ok {
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
func mustBeginTx() boil.Transactor {
	tx, ex := boil.Begin()
	utils.Must(ex)
	return tx
}

// mustConcludeTx commits or rollback the given transaction according to given error.
// Panics if Commit() or Rollback() fails.
func mustConcludeTx(tx boil.Transactor, err *HttpError) {
	if err == nil {
		utils.Must(tx.Commit())
	} else {
		utils.Must(tx.Rollback())
	}
}

// concludeRequest responds with JSON of given response or aborts the request with the given error.
func concludeRequest(c *gin.Context, resp interface{}, err *HttpError) {
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}
