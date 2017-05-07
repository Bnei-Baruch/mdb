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
	var r CollectionsRequest
	if c.Bind(&r) != nil {
		return
	}

	resp, err := handleCollectionsList(boil.GetDB(), r)
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

func CollectionItemHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleCollectionItem(boil.GetDB(), id)
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

func CollectionContentUnitsHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleCollectionCCU(boil.GetDB(), id)
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

// Toggle the active flag of a single container
func CollectionActivateHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	err := handleCollectionActivate(boil.GetDB(), id)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		err.Abort(c)
	}
}

func ContentUnitsListHandler(c *gin.Context) {
	var r ContentUnitsRequest
	if c.Bind(&r) != nil {
		return
	}

	resp, err := handleContentUnitsList(boil.GetDB(), r)
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

func ContentUnitItemHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleContentUnitItem(boil.GetDB(), id)
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

func ContentUnitFilesHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleContentUnitFiles(boil.GetDB(), id)
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

func ContentUnitCollectionsHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleContentUnitCCU(boil.GetDB(), id)
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

func FilesListHandler(c *gin.Context) {
	var r FilesRequest
	if c.Bind(&r) != nil {
		return
	}

	resp, err := handleFilesList(boil.GetDB(), r)
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

func FileItemHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleFileItem(boil.GetDB(), id)
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

func OperationsListHandler(c *gin.Context) {
	var r OperationsRequest
	if c.Bind(&r) != nil {
		return
	}

	resp, err := handleOperationsList(boil.GetDB(), r)
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

func OperationItemHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleOperationItem(boil.GetDB(), id)
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

func OperationFilesHandler(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 0)
	if e != nil {
		NewBadRequestError(errors.Wrap(e, "id expects int64")).Abort(c)
		return
	}

	resp, err := handleOperationFiles(boil.GetDB(), id)
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

// Handlers Logic

func handleCollectionsList(exec boil.Executor, r CollectionsRequest) (*CollectionsResponse, *HttpError) {
	mods := make([]qm.QueryMod, 0)

	// filters
	if err := appendDateRangeFilterMods(&mods, r.DateRangeFilter, "(properties->>'film_date')::date"); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendContentTypesFilterMods(&mods, r.ContentTypesFilter); err != nil {
		return nil, NewBadRequestError(err)
	}

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

func handleCollectionItem(exec boil.Executor, id int64) (*Collection, *HttpError) {
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

func handleCollectionActivate(exec boil.Executor, id int64) *HttpError {
	collection, err := models.FindCollection(exec, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewNotFoundError()
		} else {
			return NewInternalError(err)
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
		return NewInternalError(err)
	}
	collection.Properties = null.JSONFrom(pbytes)
	err = collection.Update(exec, "properties")
	if err != nil {
		return NewInternalError(err)
	}

	return nil
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
	if err := appendDateRangeFilterMods(&mods, r.DateRangeFilter, "(properties->>'film_date')::date"); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendContentTypesFilterMods(&mods, r.ContentTypesFilter); err != nil {
		return nil, NewBadRequestError(err)
	}
	if err := appendSourcesFilterMods(exec, &mods, r.SourcesFilter); err != nil {
		return nil, NewInternalError(err)
	}
	if err := appendTagsFilterMods(exec, &mods, r.TagsFilter); err != nil {
		return nil, NewInternalError(err)
	}

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

func handleContentUnitItem(exec boil.Executor, id int64) (*ContentUnit, *HttpError) {
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

func handleFilesList(exec boil.Executor, r FilesRequest) (*FilesResponse, *HttpError) {

	mods := make([]qm.QueryMod, 0)

	// filters
	if err := appendDateRangeFilterMods(&mods, r.DateRangeFilter, "file_created_at"); err != nil {
		return nil, NewBadRequestError(err)
	}
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

func handleFileItem(exec boil.Executor, id int64) (*MFile, *HttpError) {
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
	// slice of all source ids we want
	source_ids := make([]int64, 0)

	// fetch source ids by authors
	if !utils.IsEmpty(f.Authors) {
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

	if len(source_ids) == 0 {
		return nil
	}

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

	*mods = append(*mods,
		qm.InnerJoin("content_units_sources cus ON id = cus.content_unit_id"),
		qm.WhereIn("cus.source_id in ?", utils.ConvertArgsInt64(ids)...))

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

	*mods = append(*mods,
		qm.InnerJoin("content_units_tags cut ON id = cut.content_unit_id"),
		qm.WhereIn("cut.tag_id in ?", utils.ConvertArgsInt64(ids)...))

	return nil
}
