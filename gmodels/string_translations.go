package gmodels

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/vattle/sqlboiler/strmangle"
	"gopkg.in/nullbio/null.v6"
)

// StringTranslation is an object representing the database table.
type StringTranslation struct {
	ID               int64       `boil:"id" json:"id" toml:"id" yaml:"id"`
	Language         string      `boil:"language" json:"language" toml:"language" yaml:"language"`
	Text             string      `boil:"text" json:"text" toml:"text" yaml:"text"`
	OriginalLanguage null.String `boil:"original_language" json:"original_language,omitempty" toml:"original_language" yaml:"original_language,omitempty"`
	UserID           null.Int64  `boil:"user_id" json:"user_id,omitempty" toml:"user_id" yaml:"user_id,omitempty"`
	CreatedAt        time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *stringTranslationR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L stringTranslationL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// stringTranslationR is where relationships are stored.
type stringTranslationR struct {
	NameContentUnits        ContentUnitSlice
	DescriptionContentUnits ContentUnitSlice
	NamePersons             PersonSlice
	DescriptionPersons      PersonSlice
	NameContentRoles        ContentRoleSlice
	DescriptionContentRoles ContentRoleSlice
	LabelTags               TagSlice
	NameCollections         CollectionSlice
	DescriptionCollections  CollectionSlice
}

// stringTranslationL is where Load methods for each relationship are stored.
type stringTranslationL struct{}

var (
	stringTranslationColumns               = []string{"id", "language", "text", "original_language", "user_id", "created_at"}
	stringTranslationColumnsWithoutDefault = []string{"language", "text", "original_language", "user_id"}
	stringTranslationColumnsWithDefault    = []string{"id", "created_at"}
	stringTranslationPrimaryKeyColumns     = []string{"id", "language"}
)

type (
	// StringTranslationSlice is an alias for a slice of pointers to StringTranslation.
	// This should generally be used opposed to []StringTranslation.
	StringTranslationSlice []*StringTranslation

	stringTranslationQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	stringTranslationType                 = reflect.TypeOf(&StringTranslation{})
	stringTranslationMapping              = queries.MakeStructMapping(stringTranslationType)
	stringTranslationPrimaryKeyMapping, _ = queries.BindMapping(stringTranslationType, stringTranslationMapping, stringTranslationPrimaryKeyColumns)
	stringTranslationInsertCacheMut       sync.RWMutex
	stringTranslationInsertCache          = make(map[string]insertCache)
	stringTranslationUpdateCacheMut       sync.RWMutex
	stringTranslationUpdateCache          = make(map[string]updateCache)
	stringTranslationUpsertCacheMut       sync.RWMutex
	stringTranslationUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single stringTranslation record from the query, and panics on error.
func (q stringTranslationQuery) OneP() *StringTranslation {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single stringTranslation record from the query.
func (q stringTranslationQuery) One() (*StringTranslation, error) {
	o := &StringTranslation{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "gmodels: failed to execute a one query for string_translations")
	}

	return o, nil
}

// AllP returns all StringTranslation records from the query, and panics on error.
func (q stringTranslationQuery) AllP() StringTranslationSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all StringTranslation records from the query.
func (q stringTranslationQuery) All() (StringTranslationSlice, error) {
	var o StringTranslationSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "gmodels: failed to assign all query results to StringTranslation slice")
	}

	return o, nil
}

// CountP returns the count of all StringTranslation records in the query, and panics on error.
func (q stringTranslationQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all StringTranslation records in the query.
func (q stringTranslationQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "gmodels: failed to count string_translations rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q stringTranslationQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q stringTranslationQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "gmodels: failed to check if string_translations exists")
	}

	return count > 0, nil
}

// NameContentUnitsG retrieves all the content_unit's content units via name_id column.
func (o *StringTranslation) NameContentUnitsG(mods ...qm.QueryMod) contentUnitQuery {
	return o.NameContentUnits(boil.GetDB(), mods...)
}

// NameContentUnits retrieves all the content_unit's content units with an executor via name_id column.
func (o *StringTranslation) NameContentUnits(exec boil.Executor, mods ...qm.QueryMod) contentUnitQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"name_id\"=?", o.ID),
	)

	query := ContentUnits(exec, queryMods...)
	queries.SetFrom(query.Query, "\"content_units\" as \"a\"")
	return query
}

// DescriptionContentUnitsG retrieves all the content_unit's content units via description_id column.
func (o *StringTranslation) DescriptionContentUnitsG(mods ...qm.QueryMod) contentUnitQuery {
	return o.DescriptionContentUnits(boil.GetDB(), mods...)
}

// DescriptionContentUnits retrieves all the content_unit's content units with an executor via description_id column.
func (o *StringTranslation) DescriptionContentUnits(exec boil.Executor, mods ...qm.QueryMod) contentUnitQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"description_id\"=?", o.ID),
	)

	query := ContentUnits(exec, queryMods...)
	queries.SetFrom(query.Query, "\"content_units\" as \"a\"")
	return query
}

// NamePersonsG retrieves all the person's persons via name_id column.
func (o *StringTranslation) NamePersonsG(mods ...qm.QueryMod) personQuery {
	return o.NamePersons(boil.GetDB(), mods...)
}

// NamePersons retrieves all the person's persons with an executor via name_id column.
func (o *StringTranslation) NamePersons(exec boil.Executor, mods ...qm.QueryMod) personQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"name_id\"=?", o.ID),
	)

	query := Persons(exec, queryMods...)
	queries.SetFrom(query.Query, "\"persons\" as \"a\"")
	return query
}

// DescriptionPersonsG retrieves all the person's persons via description_id column.
func (o *StringTranslation) DescriptionPersonsG(mods ...qm.QueryMod) personQuery {
	return o.DescriptionPersons(boil.GetDB(), mods...)
}

// DescriptionPersons retrieves all the person's persons with an executor via description_id column.
func (o *StringTranslation) DescriptionPersons(exec boil.Executor, mods ...qm.QueryMod) personQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"description_id\"=?", o.ID),
	)

	query := Persons(exec, queryMods...)
	queries.SetFrom(query.Query, "\"persons\" as \"a\"")
	return query
}

// NameContentRolesG retrieves all the content_role's content roles via name_id column.
func (o *StringTranslation) NameContentRolesG(mods ...qm.QueryMod) contentRoleQuery {
	return o.NameContentRoles(boil.GetDB(), mods...)
}

// NameContentRoles retrieves all the content_role's content roles with an executor via name_id column.
func (o *StringTranslation) NameContentRoles(exec boil.Executor, mods ...qm.QueryMod) contentRoleQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"name_id\"=?", o.ID),
	)

	query := ContentRoles(exec, queryMods...)
	queries.SetFrom(query.Query, "\"content_roles\" as \"a\"")
	return query
}

// DescriptionContentRolesG retrieves all the content_role's content roles via description_id column.
func (o *StringTranslation) DescriptionContentRolesG(mods ...qm.QueryMod) contentRoleQuery {
	return o.DescriptionContentRoles(boil.GetDB(), mods...)
}

// DescriptionContentRoles retrieves all the content_role's content roles with an executor via description_id column.
func (o *StringTranslation) DescriptionContentRoles(exec boil.Executor, mods ...qm.QueryMod) contentRoleQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"description_id\"=?", o.ID),
	)

	query := ContentRoles(exec, queryMods...)
	queries.SetFrom(query.Query, "\"content_roles\" as \"a\"")
	return query
}

// LabelTagsG retrieves all the tag's tags via label_id column.
func (o *StringTranslation) LabelTagsG(mods ...qm.QueryMod) tagQuery {
	return o.LabelTags(boil.GetDB(), mods...)
}

// LabelTags retrieves all the tag's tags with an executor via label_id column.
func (o *StringTranslation) LabelTags(exec boil.Executor, mods ...qm.QueryMod) tagQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"label_id\"=?", o.ID),
	)

	query := Tags(exec, queryMods...)
	queries.SetFrom(query.Query, "\"tags\" as \"a\"")
	return query
}

// NameCollectionsG retrieves all the collection's collections via name_id column.
func (o *StringTranslation) NameCollectionsG(mods ...qm.QueryMod) collectionQuery {
	return o.NameCollections(boil.GetDB(), mods...)
}

// NameCollections retrieves all the collection's collections with an executor via name_id column.
func (o *StringTranslation) NameCollections(exec boil.Executor, mods ...qm.QueryMod) collectionQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"name_id\"=?", o.ID),
	)

	query := Collections(exec, queryMods...)
	queries.SetFrom(query.Query, "\"collections\" as \"a\"")
	return query
}

// DescriptionCollectionsG retrieves all the collection's collections via description_id column.
func (o *StringTranslation) DescriptionCollectionsG(mods ...qm.QueryMod) collectionQuery {
	return o.DescriptionCollections(boil.GetDB(), mods...)
}

// DescriptionCollections retrieves all the collection's collections with an executor via description_id column.
func (o *StringTranslation) DescriptionCollections(exec boil.Executor, mods ...qm.QueryMod) collectionQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"description_id\"=?", o.ID),
	)

	query := Collections(exec, queryMods...)
	queries.SetFrom(query.Query, "\"collections\" as \"a\"")
	return query
}

// LoadNameContentUnits allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (stringTranslationL) LoadNameContentUnits(e boil.Executor, singular bool, maybeStringTranslation interface{}) error {
	var slice []*StringTranslation
	var object *StringTranslation

	count := 1
	if singular {
		object = maybeStringTranslation.(*StringTranslation)
	} else {
		slice = *maybeStringTranslation.(*StringTranslationSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &stringTranslationR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &stringTranslationR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"content_units\" where \"name_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load content_units")
	}
	defer results.Close()

	var resultSlice []*ContentUnit
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice content_units")
	}

	if singular {
		object.R.NameContentUnits = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.NameID {
				local.R.NameContentUnits = append(local.R.NameContentUnits, foreign)
				break
			}
		}
	}

	return nil
}

// LoadDescriptionContentUnits allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (stringTranslationL) LoadDescriptionContentUnits(e boil.Executor, singular bool, maybeStringTranslation interface{}) error {
	var slice []*StringTranslation
	var object *StringTranslation

	count := 1
	if singular {
		object = maybeStringTranslation.(*StringTranslation)
	} else {
		slice = *maybeStringTranslation.(*StringTranslationSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &stringTranslationR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &stringTranslationR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"content_units\" where \"description_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load content_units")
	}
	defer results.Close()

	var resultSlice []*ContentUnit
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice content_units")
	}

	if singular {
		object.R.DescriptionContentUnits = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.DescriptionID.Int64 {
				local.R.DescriptionContentUnits = append(local.R.DescriptionContentUnits, foreign)
				break
			}
		}
	}

	return nil
}

// LoadNamePersons allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (stringTranslationL) LoadNamePersons(e boil.Executor, singular bool, maybeStringTranslation interface{}) error {
	var slice []*StringTranslation
	var object *StringTranslation

	count := 1
	if singular {
		object = maybeStringTranslation.(*StringTranslation)
	} else {
		slice = *maybeStringTranslation.(*StringTranslationSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &stringTranslationR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &stringTranslationR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"persons\" where \"name_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load persons")
	}
	defer results.Close()

	var resultSlice []*Person
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice persons")
	}

	if singular {
		object.R.NamePersons = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.NameID {
				local.R.NamePersons = append(local.R.NamePersons, foreign)
				break
			}
		}
	}

	return nil
}

// LoadDescriptionPersons allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (stringTranslationL) LoadDescriptionPersons(e boil.Executor, singular bool, maybeStringTranslation interface{}) error {
	var slice []*StringTranslation
	var object *StringTranslation

	count := 1
	if singular {
		object = maybeStringTranslation.(*StringTranslation)
	} else {
		slice = *maybeStringTranslation.(*StringTranslationSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &stringTranslationR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &stringTranslationR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"persons\" where \"description_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load persons")
	}
	defer results.Close()

	var resultSlice []*Person
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice persons")
	}

	if singular {
		object.R.DescriptionPersons = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.DescriptionID.Int64 {
				local.R.DescriptionPersons = append(local.R.DescriptionPersons, foreign)
				break
			}
		}
	}

	return nil
}

// LoadNameContentRoles allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (stringTranslationL) LoadNameContentRoles(e boil.Executor, singular bool, maybeStringTranslation interface{}) error {
	var slice []*StringTranslation
	var object *StringTranslation

	count := 1
	if singular {
		object = maybeStringTranslation.(*StringTranslation)
	} else {
		slice = *maybeStringTranslation.(*StringTranslationSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &stringTranslationR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &stringTranslationR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"content_roles\" where \"name_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load content_roles")
	}
	defer results.Close()

	var resultSlice []*ContentRole
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice content_roles")
	}

	if singular {
		object.R.NameContentRoles = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.NameID {
				local.R.NameContentRoles = append(local.R.NameContentRoles, foreign)
				break
			}
		}
	}

	return nil
}

// LoadDescriptionContentRoles allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (stringTranslationL) LoadDescriptionContentRoles(e boil.Executor, singular bool, maybeStringTranslation interface{}) error {
	var slice []*StringTranslation
	var object *StringTranslation

	count := 1
	if singular {
		object = maybeStringTranslation.(*StringTranslation)
	} else {
		slice = *maybeStringTranslation.(*StringTranslationSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &stringTranslationR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &stringTranslationR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"content_roles\" where \"description_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load content_roles")
	}
	defer results.Close()

	var resultSlice []*ContentRole
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice content_roles")
	}

	if singular {
		object.R.DescriptionContentRoles = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.DescriptionID.Int64 {
				local.R.DescriptionContentRoles = append(local.R.DescriptionContentRoles, foreign)
				break
			}
		}
	}

	return nil
}

// LoadLabelTags allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (stringTranslationL) LoadLabelTags(e boil.Executor, singular bool, maybeStringTranslation interface{}) error {
	var slice []*StringTranslation
	var object *StringTranslation

	count := 1
	if singular {
		object = maybeStringTranslation.(*StringTranslation)
	} else {
		slice = *maybeStringTranslation.(*StringTranslationSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &stringTranslationR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &stringTranslationR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"tags\" where \"label_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load tags")
	}
	defer results.Close()

	var resultSlice []*Tag
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice tags")
	}

	if singular {
		object.R.LabelTags = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.LabelID {
				local.R.LabelTags = append(local.R.LabelTags, foreign)
				break
			}
		}
	}

	return nil
}

// LoadNameCollections allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (stringTranslationL) LoadNameCollections(e boil.Executor, singular bool, maybeStringTranslation interface{}) error {
	var slice []*StringTranslation
	var object *StringTranslation

	count := 1
	if singular {
		object = maybeStringTranslation.(*StringTranslation)
	} else {
		slice = *maybeStringTranslation.(*StringTranslationSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &stringTranslationR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &stringTranslationR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"collections\" where \"name_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load collections")
	}
	defer results.Close()

	var resultSlice []*Collection
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice collections")
	}

	if singular {
		object.R.NameCollections = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.NameID {
				local.R.NameCollections = append(local.R.NameCollections, foreign)
				break
			}
		}
	}

	return nil
}

// LoadDescriptionCollections allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (stringTranslationL) LoadDescriptionCollections(e boil.Executor, singular bool, maybeStringTranslation interface{}) error {
	var slice []*StringTranslation
	var object *StringTranslation

	count := 1
	if singular {
		object = maybeStringTranslation.(*StringTranslation)
	} else {
		slice = *maybeStringTranslation.(*StringTranslationSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &stringTranslationR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &stringTranslationR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"collections\" where \"description_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load collections")
	}
	defer results.Close()

	var resultSlice []*Collection
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice collections")
	}

	if singular {
		object.R.DescriptionCollections = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.DescriptionID.Int64 {
				local.R.DescriptionCollections = append(local.R.DescriptionCollections, foreign)
				break
			}
		}
	}

	return nil
}

// AddNameContentUnits adds the given related objects to the existing relationships
// of the string_translation, optionally inserting them as new records.
// Appends related to o.R.NameContentUnits.
// Sets related.R.Name appropriately.
func (o *StringTranslation) AddNameContentUnits(exec boil.Executor, insert bool, related ...*ContentUnit) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.NameID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"content_units\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"name_id"}),
				strmangle.WhereClause("\"", "\"", 2, contentUnitPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.NameID = o.ID
		}
	}

	if o.R == nil {
		o.R = &stringTranslationR{
			NameContentUnits: related,
		}
	} else {
		o.R.NameContentUnits = append(o.R.NameContentUnits, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &contentUnitR{
				Name: o,
			}
		} else {
			rel.R.Name = o
		}
	}
	return nil
}

// AddDescriptionContentUnits adds the given related objects to the existing relationships
// of the string_translation, optionally inserting them as new records.
// Appends related to o.R.DescriptionContentUnits.
// Sets related.R.Description appropriately.
func (o *StringTranslation) AddDescriptionContentUnits(exec boil.Executor, insert bool, related ...*ContentUnit) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.DescriptionID.Int64 = o.ID
			rel.DescriptionID.Valid = true
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"content_units\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"description_id"}),
				strmangle.WhereClause("\"", "\"", 2, contentUnitPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.DescriptionID.Int64 = o.ID
			rel.DescriptionID.Valid = true
		}
	}

	if o.R == nil {
		o.R = &stringTranslationR{
			DescriptionContentUnits: related,
		}
	} else {
		o.R.DescriptionContentUnits = append(o.R.DescriptionContentUnits, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &contentUnitR{
				Description: o,
			}
		} else {
			rel.R.Description = o
		}
	}
	return nil
}

// SetDescriptionContentUnits removes all previously related items of the
// string_translation replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Description's DescriptionContentUnits accordingly.
// Replaces o.R.DescriptionContentUnits with related.
// Sets related.R.Description's DescriptionContentUnits accordingly.
func (o *StringTranslation) SetDescriptionContentUnits(exec boil.Executor, insert bool, related ...*ContentUnit) error {
	query := "update \"content_units\" set \"description_id\" = null where \"description_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.DescriptionContentUnits {
			rel.DescriptionID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Description = nil
		}

		o.R.DescriptionContentUnits = nil
	}
	return o.AddDescriptionContentUnits(exec, insert, related...)
}

// RemoveDescriptionContentUnits relationships from objects passed in.
// Removes related items from R.DescriptionContentUnits (uses pointer comparison, removal does not keep order)
// Sets related.R.Description.
func (o *StringTranslation) RemoveDescriptionContentUnits(exec boil.Executor, related ...*ContentUnit) error {
	var err error
	for _, rel := range related {
		rel.DescriptionID.Valid = false
		if rel.R != nil {
			rel.R.Description = nil
		}
		if err = rel.Update(exec, "description_id"); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.DescriptionContentUnits {
			if rel != ri {
				continue
			}

			ln := len(o.R.DescriptionContentUnits)
			if ln > 1 && i < ln-1 {
				o.R.DescriptionContentUnits[i] = o.R.DescriptionContentUnits[ln-1]
			}
			o.R.DescriptionContentUnits = o.R.DescriptionContentUnits[:ln-1]
			break
		}
	}

	return nil
}

// AddNamePersons adds the given related objects to the existing relationships
// of the string_translation, optionally inserting them as new records.
// Appends related to o.R.NamePersons.
// Sets related.R.Name appropriately.
func (o *StringTranslation) AddNamePersons(exec boil.Executor, insert bool, related ...*Person) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.NameID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"persons\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"name_id"}),
				strmangle.WhereClause("\"", "\"", 2, personPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.NameID = o.ID
		}
	}

	if o.R == nil {
		o.R = &stringTranslationR{
			NamePersons: related,
		}
	} else {
		o.R.NamePersons = append(o.R.NamePersons, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &personR{
				Name: o,
			}
		} else {
			rel.R.Name = o
		}
	}
	return nil
}

// AddDescriptionPersons adds the given related objects to the existing relationships
// of the string_translation, optionally inserting them as new records.
// Appends related to o.R.DescriptionPersons.
// Sets related.R.Description appropriately.
func (o *StringTranslation) AddDescriptionPersons(exec boil.Executor, insert bool, related ...*Person) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.DescriptionID.Int64 = o.ID
			rel.DescriptionID.Valid = true
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"persons\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"description_id"}),
				strmangle.WhereClause("\"", "\"", 2, personPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.DescriptionID.Int64 = o.ID
			rel.DescriptionID.Valid = true
		}
	}

	if o.R == nil {
		o.R = &stringTranslationR{
			DescriptionPersons: related,
		}
	} else {
		o.R.DescriptionPersons = append(o.R.DescriptionPersons, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &personR{
				Description: o,
			}
		} else {
			rel.R.Description = o
		}
	}
	return nil
}

// SetDescriptionPersons removes all previously related items of the
// string_translation replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Description's DescriptionPersons accordingly.
// Replaces o.R.DescriptionPersons with related.
// Sets related.R.Description's DescriptionPersons accordingly.
func (o *StringTranslation) SetDescriptionPersons(exec boil.Executor, insert bool, related ...*Person) error {
	query := "update \"persons\" set \"description_id\" = null where \"description_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.DescriptionPersons {
			rel.DescriptionID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Description = nil
		}

		o.R.DescriptionPersons = nil
	}
	return o.AddDescriptionPersons(exec, insert, related...)
}

// RemoveDescriptionPersons relationships from objects passed in.
// Removes related items from R.DescriptionPersons (uses pointer comparison, removal does not keep order)
// Sets related.R.Description.
func (o *StringTranslation) RemoveDescriptionPersons(exec boil.Executor, related ...*Person) error {
	var err error
	for _, rel := range related {
		rel.DescriptionID.Valid = false
		if rel.R != nil {
			rel.R.Description = nil
		}
		if err = rel.Update(exec, "description_id"); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.DescriptionPersons {
			if rel != ri {
				continue
			}

			ln := len(o.R.DescriptionPersons)
			if ln > 1 && i < ln-1 {
				o.R.DescriptionPersons[i] = o.R.DescriptionPersons[ln-1]
			}
			o.R.DescriptionPersons = o.R.DescriptionPersons[:ln-1]
			break
		}
	}

	return nil
}

// AddNameContentRoles adds the given related objects to the existing relationships
// of the string_translation, optionally inserting them as new records.
// Appends related to o.R.NameContentRoles.
// Sets related.R.Name appropriately.
func (o *StringTranslation) AddNameContentRoles(exec boil.Executor, insert bool, related ...*ContentRole) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.NameID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"content_roles\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"name_id"}),
				strmangle.WhereClause("\"", "\"", 2, contentRolePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.NameID = o.ID
		}
	}

	if o.R == nil {
		o.R = &stringTranslationR{
			NameContentRoles: related,
		}
	} else {
		o.R.NameContentRoles = append(o.R.NameContentRoles, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &contentRoleR{
				Name: o,
			}
		} else {
			rel.R.Name = o
		}
	}
	return nil
}

// AddDescriptionContentRoles adds the given related objects to the existing relationships
// of the string_translation, optionally inserting them as new records.
// Appends related to o.R.DescriptionContentRoles.
// Sets related.R.Description appropriately.
func (o *StringTranslation) AddDescriptionContentRoles(exec boil.Executor, insert bool, related ...*ContentRole) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.DescriptionID.Int64 = o.ID
			rel.DescriptionID.Valid = true
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"content_roles\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"description_id"}),
				strmangle.WhereClause("\"", "\"", 2, contentRolePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.DescriptionID.Int64 = o.ID
			rel.DescriptionID.Valid = true
		}
	}

	if o.R == nil {
		o.R = &stringTranslationR{
			DescriptionContentRoles: related,
		}
	} else {
		o.R.DescriptionContentRoles = append(o.R.DescriptionContentRoles, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &contentRoleR{
				Description: o,
			}
		} else {
			rel.R.Description = o
		}
	}
	return nil
}

// SetDescriptionContentRoles removes all previously related items of the
// string_translation replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Description's DescriptionContentRoles accordingly.
// Replaces o.R.DescriptionContentRoles with related.
// Sets related.R.Description's DescriptionContentRoles accordingly.
func (o *StringTranslation) SetDescriptionContentRoles(exec boil.Executor, insert bool, related ...*ContentRole) error {
	query := "update \"content_roles\" set \"description_id\" = null where \"description_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.DescriptionContentRoles {
			rel.DescriptionID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Description = nil
		}

		o.R.DescriptionContentRoles = nil
	}
	return o.AddDescriptionContentRoles(exec, insert, related...)
}

// RemoveDescriptionContentRoles relationships from objects passed in.
// Removes related items from R.DescriptionContentRoles (uses pointer comparison, removal does not keep order)
// Sets related.R.Description.
func (o *StringTranslation) RemoveDescriptionContentRoles(exec boil.Executor, related ...*ContentRole) error {
	var err error
	for _, rel := range related {
		rel.DescriptionID.Valid = false
		if rel.R != nil {
			rel.R.Description = nil
		}
		if err = rel.Update(exec, "description_id"); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.DescriptionContentRoles {
			if rel != ri {
				continue
			}

			ln := len(o.R.DescriptionContentRoles)
			if ln > 1 && i < ln-1 {
				o.R.DescriptionContentRoles[i] = o.R.DescriptionContentRoles[ln-1]
			}
			o.R.DescriptionContentRoles = o.R.DescriptionContentRoles[:ln-1]
			break
		}
	}

	return nil
}

// AddLabelTags adds the given related objects to the existing relationships
// of the string_translation, optionally inserting them as new records.
// Appends related to o.R.LabelTags.
// Sets related.R.Label appropriately.
func (o *StringTranslation) AddLabelTags(exec boil.Executor, insert bool, related ...*Tag) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.LabelID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"tags\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"label_id"}),
				strmangle.WhereClause("\"", "\"", 2, tagPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.LabelID = o.ID
		}
	}

	if o.R == nil {
		o.R = &stringTranslationR{
			LabelTags: related,
		}
	} else {
		o.R.LabelTags = append(o.R.LabelTags, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &tagR{
				Label: o,
			}
		} else {
			rel.R.Label = o
		}
	}
	return nil
}

// AddNameCollections adds the given related objects to the existing relationships
// of the string_translation, optionally inserting them as new records.
// Appends related to o.R.NameCollections.
// Sets related.R.Name appropriately.
func (o *StringTranslation) AddNameCollections(exec boil.Executor, insert bool, related ...*Collection) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.NameID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"collections\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"name_id"}),
				strmangle.WhereClause("\"", "\"", 2, collectionPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.NameID = o.ID
		}
	}

	if o.R == nil {
		o.R = &stringTranslationR{
			NameCollections: related,
		}
	} else {
		o.R.NameCollections = append(o.R.NameCollections, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &collectionR{
				Name: o,
			}
		} else {
			rel.R.Name = o
		}
	}
	return nil
}

// AddDescriptionCollections adds the given related objects to the existing relationships
// of the string_translation, optionally inserting them as new records.
// Appends related to o.R.DescriptionCollections.
// Sets related.R.Description appropriately.
func (o *StringTranslation) AddDescriptionCollections(exec boil.Executor, insert bool, related ...*Collection) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.DescriptionID.Int64 = o.ID
			rel.DescriptionID.Valid = true
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"collections\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"description_id"}),
				strmangle.WhereClause("\"", "\"", 2, collectionPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.DescriptionID.Int64 = o.ID
			rel.DescriptionID.Valid = true
		}
	}

	if o.R == nil {
		o.R = &stringTranslationR{
			DescriptionCollections: related,
		}
	} else {
		o.R.DescriptionCollections = append(o.R.DescriptionCollections, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &collectionR{
				Description: o,
			}
		} else {
			rel.R.Description = o
		}
	}
	return nil
}

// SetDescriptionCollections removes all previously related items of the
// string_translation replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Description's DescriptionCollections accordingly.
// Replaces o.R.DescriptionCollections with related.
// Sets related.R.Description's DescriptionCollections accordingly.
func (o *StringTranslation) SetDescriptionCollections(exec boil.Executor, insert bool, related ...*Collection) error {
	query := "update \"collections\" set \"description_id\" = null where \"description_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.DescriptionCollections {
			rel.DescriptionID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Description = nil
		}

		o.R.DescriptionCollections = nil
	}
	return o.AddDescriptionCollections(exec, insert, related...)
}

// RemoveDescriptionCollections relationships from objects passed in.
// Removes related items from R.DescriptionCollections (uses pointer comparison, removal does not keep order)
// Sets related.R.Description.
func (o *StringTranslation) RemoveDescriptionCollections(exec boil.Executor, related ...*Collection) error {
	var err error
	for _, rel := range related {
		rel.DescriptionID.Valid = false
		if rel.R != nil {
			rel.R.Description = nil
		}
		if err = rel.Update(exec, "description_id"); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.DescriptionCollections {
			if rel != ri {
				continue
			}

			ln := len(o.R.DescriptionCollections)
			if ln > 1 && i < ln-1 {
				o.R.DescriptionCollections[i] = o.R.DescriptionCollections[ln-1]
			}
			o.R.DescriptionCollections = o.R.DescriptionCollections[:ln-1]
			break
		}
	}

	return nil
}

// StringTranslationsG retrieves all records.
func StringTranslationsG(mods ...qm.QueryMod) stringTranslationQuery {
	return StringTranslations(boil.GetDB(), mods...)
}

// StringTranslations retrieves all the records using an executor.
func StringTranslations(exec boil.Executor, mods ...qm.QueryMod) stringTranslationQuery {
	mods = append(mods, qm.From("\"string_translations\""))
	return stringTranslationQuery{NewQuery(exec, mods...)}
}

// FindStringTranslationG retrieves a single record by ID.
func FindStringTranslationG(id int64, language string, selectCols ...string) (*StringTranslation, error) {
	return FindStringTranslation(boil.GetDB(), id, language, selectCols...)
}

// FindStringTranslationGP retrieves a single record by ID, and panics on error.
func FindStringTranslationGP(id int64, language string, selectCols ...string) *StringTranslation {
	retobj, err := FindStringTranslation(boil.GetDB(), id, language, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindStringTranslation retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindStringTranslation(exec boil.Executor, id int64, language string, selectCols ...string) (*StringTranslation, error) {
	stringTranslationObj := &StringTranslation{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"string_translations\" where \"id\"=$1 AND \"language\"=$2", sel,
	)

	q := queries.Raw(exec, query, id, language)

	err := q.Bind(stringTranslationObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "gmodels: unable to select from string_translations")
	}

	return stringTranslationObj, nil
}

// FindStringTranslationP retrieves a single record by ID with an executor, and panics on error.
func FindStringTranslationP(exec boil.Executor, id int64, language string, selectCols ...string) *StringTranslation {
	retobj, err := FindStringTranslation(exec, id, language, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *StringTranslation) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *StringTranslation) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *StringTranslation) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *StringTranslation) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("gmodels: no string_translations provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	nzDefaults := queries.NonZeroDefaultSet(stringTranslationColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	stringTranslationInsertCacheMut.RLock()
	cache, cached := stringTranslationInsertCache[key]
	stringTranslationInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			stringTranslationColumns,
			stringTranslationColumnsWithDefault,
			stringTranslationColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(stringTranslationType, stringTranslationMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(stringTranslationType, stringTranslationMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"string_translations\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

		if len(cache.retMapping) != 0 {
			cache.query += fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "gmodels: unable to insert into string_translations")
	}

	if !cached {
		stringTranslationInsertCacheMut.Lock()
		stringTranslationInsertCache[key] = cache
		stringTranslationInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single StringTranslation record. See Update for
// whitelist behavior description.
func (o *StringTranslation) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single StringTranslation record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *StringTranslation) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the StringTranslation, and panics on error.
// See Update for whitelist behavior description.
func (o *StringTranslation) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the StringTranslation.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *StringTranslation) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	stringTranslationUpdateCacheMut.RLock()
	cache, cached := stringTranslationUpdateCache[key]
	stringTranslationUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(stringTranslationColumns, stringTranslationPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("gmodels: unable to update string_translations, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"string_translations\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, stringTranslationPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(stringTranslationType, stringTranslationMapping, append(wl, stringTranslationPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to update string_translations row")
	}

	if !cached {
		stringTranslationUpdateCacheMut.Lock()
		stringTranslationUpdateCache[key] = cache
		stringTranslationUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q stringTranslationQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q stringTranslationQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to update all for string_translations")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o StringTranslationSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o StringTranslationSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o StringTranslationSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o StringTranslationSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("gmodels: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), stringTranslationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"string_translations\" SET %s WHERE (\"id\",\"language\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(stringTranslationPrimaryKeyColumns), len(colNames)+1, len(stringTranslationPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to update all in stringTranslation slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *StringTranslation) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *StringTranslation) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *StringTranslation) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *StringTranslation) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("gmodels: no string_translations provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	nzDefaults := queries.NonZeroDefaultSet(stringTranslationColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs postgres problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range updateColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range whitelist {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	stringTranslationUpsertCacheMut.RLock()
	cache, cached := stringTranslationUpsertCache[key]
	stringTranslationUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			stringTranslationColumns,
			stringTranslationColumnsWithDefault,
			stringTranslationColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			stringTranslationColumns,
			stringTranslationPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("gmodels: unable to upsert string_translations, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(stringTranslationPrimaryKeyColumns))
			copy(conflict, stringTranslationPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"string_translations\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(stringTranslationType, stringTranslationMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(stringTranslationType, stringTranslationMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to upsert string_translations")
	}

	if !cached {
		stringTranslationUpsertCacheMut.Lock()
		stringTranslationUpsertCache[key] = cache
		stringTranslationUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single StringTranslation record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *StringTranslation) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single StringTranslation record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *StringTranslation) DeleteG() error {
	if o == nil {
		return errors.New("gmodels: no StringTranslation provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single StringTranslation record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *StringTranslation) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single StringTranslation record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *StringTranslation) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("gmodels: no StringTranslation provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), stringTranslationPrimaryKeyMapping)
	sql := "DELETE FROM \"string_translations\" WHERE \"id\"=$1 AND \"language\"=$2"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete from string_translations")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q stringTranslationQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q stringTranslationQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("gmodels: no stringTranslationQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete all from string_translations")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o StringTranslationSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o StringTranslationSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("gmodels: no StringTranslation slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o StringTranslationSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o StringTranslationSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("gmodels: no StringTranslation slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), stringTranslationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"string_translations\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, stringTranslationPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(stringTranslationPrimaryKeyColumns), 1, len(stringTranslationPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete all from stringTranslation slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *StringTranslation) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *StringTranslation) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *StringTranslation) ReloadG() error {
	if o == nil {
		return errors.New("gmodels: no StringTranslation provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *StringTranslation) Reload(exec boil.Executor) error {
	ret, err := FindStringTranslation(exec, o.ID, o.Language)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *StringTranslationSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *StringTranslationSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *StringTranslationSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("gmodels: empty StringTranslationSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *StringTranslationSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	stringTranslations := StringTranslationSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), stringTranslationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"string_translations\".* FROM \"string_translations\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, stringTranslationPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(stringTranslationPrimaryKeyColumns), 1, len(stringTranslationPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&stringTranslations)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to reload all in StringTranslationSlice")
	}

	*o = stringTranslations

	return nil
}

// StringTranslationExists checks if the StringTranslation row exists.
func StringTranslationExists(exec boil.Executor, id int64, language string) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"string_translations\" where \"id\"=$1 AND \"language\"=$2 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id, language)
	}

	row := exec.QueryRow(sql, id, language)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "gmodels: unable to check if string_translations exists")
	}

	return exists, nil
}

// StringTranslationExistsG checks if the StringTranslation row exists.
func StringTranslationExistsG(id int64, language string) (bool, error) {
	return StringTranslationExists(boil.GetDB(), id, language)
}

// StringTranslationExistsGP checks if the StringTranslation row exists. Panics on error.
func StringTranslationExistsGP(id int64, language string) bool {
	e, err := StringTranslationExists(boil.GetDB(), id, language)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// StringTranslationExistsP checks if the StringTranslation row exists. Panics on error.
func StringTranslationExistsP(exec boil.Executor, id int64, language string) bool {
	e, err := StringTranslationExists(exec, id, language)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
