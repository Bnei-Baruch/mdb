package common

import (
	"github.com/Bnei-Baruch/mdb/models"
)

type CollectionWName struct {
	*models.Collection
	name string
}

func (c *CollectionWName) Name() string {
	if c.name != "" {
		return c.name
	}

	ci18ns := make(map[string]string)
	for i := range c.R.CollectionI18ns {
		i18n := c.R.CollectionI18ns[i]
		if i18n.Name.Valid {
			ci18ns[i18n.Language] = i18n.Name.String
		}
	}
	if v, ok := ci18ns[LANG_HEBREW]; ok {
		c.name = v
	} else if v, ok := ci18ns[LANG_ENGLISH]; ok {
		c.name = v
	} else if v, ok := ci18ns[LANG_RUSSIAN]; ok {
		c.name = v
	}
	return c.name
}

type UnitWName struct {
	*models.ContentUnit
	name        string
	description string
}

func (cu *UnitWName) Name() string {
	if cu.name != "" {
		return cu.name
	}

	cui18ns := make(map[string]string)
	for i := range cu.R.ContentUnitI18ns {
		i18n := cu.R.ContentUnitI18ns[i]
		if i18n.Name.Valid {
			cui18ns[i18n.Language] = i18n.Name.String
		}
	}
	if v, ok := cui18ns[LANG_HEBREW]; ok {
		cu.name = v
	} else if v, ok := cui18ns[LANG_ENGLISH]; ok {
		cu.name = v
	} else if v, ok := cui18ns[LANG_RUSSIAN]; ok {
		cu.name = v
	}
	return cu.name
}

func (cu *UnitWName) Description() string {
	if cu.description != "" {
		return cu.description
	}

	cui18ns := make(map[string]string)
	for i := range cu.R.ContentUnitI18ns {
		i18n := cu.R.ContentUnitI18ns[i]
		if i18n.Description.Valid {
			cui18ns[i18n.Language] = i18n.Description.String
		}
	}
	if v, ok := cui18ns[LANG_HEBREW]; ok {
		cu.description = v
	} else if v, ok := cui18ns[LANG_ENGLISH]; ok {
		cu.description = v
	} else if v, ok := cui18ns[LANG_RUSSIAN]; ok {
		cu.description = v
	}
	return cu.description
}

