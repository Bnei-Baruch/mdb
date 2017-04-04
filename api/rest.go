package api

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

const (
	DEFAULT_PAGE_SIZE = 50
	MAX_PAGE_SIZE     = 100
)

func CollectionsListHandler(c *gin.Context) {
	var r CollectionsRequest
	if c.Bind(&r) != nil {
		return
	}

	mods := make([]qm.QueryMod, 0)

	// filters
	if err := appendContentTypesFilterMods(&mods, r.ContentTypesFilter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
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
	appendListMods(&mods, r.ListRequest)

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

func appendListMods(mods *[]qm.QueryMod, r ListRequest) {
	if r.OrderBy == "" {
		*mods = append(*mods, qm.OrderBy("id"))
	} else {
		*mods = append(*mods, qm.OrderBy(r.OrderBy))
	}

	var pageSize int
	if r.PageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	} else {
		pageSize = utils.Min(r.PageSize, MAX_PAGE_SIZE)
	}

	*mods = append(*mods, qm.Limit(pageSize))
	if r.PageNumber > 1 {
		*mods = append(*mods, qm.Offset((r.PageNumber-1)*pageSize))
	}
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
