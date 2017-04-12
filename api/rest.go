package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
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

	mods := make([]qm.QueryMod, 0)

	// filters
	if err := appendContentTypesFilterMods(&mods, r.ContentTypesFilter); err != nil {
		c.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}

	// count query
	total, err := models.CollectionsG(mods...).Count()
	if err != nil {
		internalServerError(c, err)
		return
	}
	if total == 0 {
		c.JSON(http.StatusOK, NewCollectionsResponse())
		return
	}

	// order, limit, offset
	if err = appendListMods(&mods, r.ListRequest); err != nil {
		c.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}

	// Eager loading
	mods = append(mods, qm.Load("CollectionI18ns"))

	// data query
	collections, err := models.CollectionsG(mods...).All()
	if err != nil {
		internalServerError(c, err)
		return
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

	c.JSON(http.StatusOK, CollectionsResponse{
		ListResponse: ListResponse{Total: total},
		Collections:  data,
	})
}

// Toggle the active flag of a single container
func CollectionActivateHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 0)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "id expects int64")).
			SetType(gin.ErrorTypePublic)
		return
	}

	collection, err := models.FindCollectionG(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatus(http.StatusNotFound)
			return
		} else {
			internalServerError(c, err)
			return
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
		internalServerError(c, err)
	}
	collection.Properties = null.JSONFrom(pbytes)
	err = collection.UpdateG("properties")
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		internalServerError(c, err)
	}
}

func FilesListHandler(c *gin.Context) {
	var r FilesRequest
	if c.Bind(&r) != nil {
		return
	}

	mods := make([]qm.QueryMod, 0)

	// filters
	if r.Query != "" {
		mods = append(mods, qm.Where("name LIKE ?", r.Query),
			qm.Or("uid LIKE ?", r.Query),
			qm.Or("id::TEXT LIKE ?", r.Query))
	}

	// count query
	total, err := models.FilesG(mods...).Count()
	if err != nil {
		internalServerError(c, err)
		return
	}
	if total == 0 {
		c.JSON(http.StatusOK, NewFilesResponse())
		return
	}

	// order, limit, offset
	if err = appendListMods(&mods, r.ListRequest); err != nil {
		c.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
		return
	}

	// data query
	files, err := models.FilesG(mods...).All()
	if err != nil {
		internalServerError(c, err)
		return
	}

	data := make([]*MFile, len(files))
	for i, f := range files {
		data[i] = NewMFile(f)
	}

	c.JSON(http.StatusOK, FilesResponse{
		ListResponse: ListResponse{Total: total},
		Files:        data,
	})
}

func FileItemHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 0)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "id expects int64")).
			SetType(gin.ErrorTypePublic)
		return
	}

	file, err := models.FindFileG(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatus(http.StatusNotFound)
			return
		} else {
			internalServerError(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, NewMFile(file))
}

func appendListMods(mods *[]qm.QueryMod, r ListRequest) error {
	if r.OrderBy == "" {
		*mods = append(*mods, qm.OrderBy("id"))
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
