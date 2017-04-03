package models

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

// Source is an object representing the database table.
type Source struct {
	ID          int64       `boil:"id" json:"id" toml:"id" yaml:"id"`
	UID         string      `boil:"uid" json:"uid" toml:"uid" yaml:"uid"`
	ParentID    null.Int64  `boil:"parent_id" json:"parent_id,omitempty" toml:"parent_id" yaml:"parent_id,omitempty"`
	Pattern     null.String `boil:"pattern" json:"pattern,omitempty" toml:"pattern" yaml:"pattern,omitempty"`
	TypeID      int64       `boil:"type_id" json:"type_id" toml:"type_id" yaml:"type_id"`
	Position    null.Int    `boil:"position" json:"position,omitempty" toml:"position" yaml:"position,omitempty"`
	Name        string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	Description null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`
	CreatedAt   time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	Properties  null.JSON   `boil:"properties" json:"properties,omitempty" toml:"properties" yaml:"properties,omitempty"`

	R *sourceR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L sourceL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// sourceR is where relationships are stored.
type sourceR struct {
	Parent        *Source
	Type          *SourceType
	ParentSources SourceSlice
	ContentUnits  ContentUnitSlice
	Authors       AuthorSlice
	SourceI18ns   SourceI18nSlice
}

// sourceL is where Load methods for each relationship are stored.
type sourceL struct{}

var (
	sourceColumns               = []string{"id", "uid", "parent_id", "pattern", "type_id", "position", "name", "description", "created_at", "properties"}
	sourceColumnsWithoutDefault = []string{"uid", "parent_id", "pattern", "type_id", "position", "name", "description", "properties"}
	sourceColumnsWithDefault    = []string{"id", "created_at"}
	sourcePrimaryKeyColumns     = []string{"id"}
)

type (
	// SourceSlice is an alias for a slice of pointers to Source.
	// This should generally be used opposed to []Source.
	SourceSlice []*Source

	sourceQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	sourceType                 = reflect.TypeOf(&Source{})
	sourceMapping              = queries.MakeStructMapping(sourceType)
	sourcePrimaryKeyMapping, _ = queries.BindMapping(sourceType, sourceMapping, sourcePrimaryKeyColumns)
	sourceInsertCacheMut       sync.RWMutex
	sourceInsertCache          = make(map[string]insertCache)
	sourceUpdateCacheMut       sync.RWMutex
	sourceUpdateCache          = make(map[string]updateCache)
	sourceUpsertCacheMut       sync.RWMutex
	sourceUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single source record from the query, and panics on error.
func (q sourceQuery) OneP() *Source {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single source record from the query.
func (q sourceQuery) One() (*Source, error) {
	o := &Source{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for sources")
	}

	return o, nil
}

// AllP returns all Source records from the query, and panics on error.
func (q sourceQuery) AllP() SourceSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Source records from the query.
func (q sourceQuery) All() (SourceSlice, error) {
	var o SourceSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Source slice")
	}

	return o, nil
}

// CountP returns the count of all Source records in the query, and panics on error.
func (q sourceQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Source records in the query.
func (q sourceQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count sources rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q sourceQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q sourceQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if sources exists")
	}

	return count > 0, nil
}

// ParentG pointed to by the foreign key.
func (o *Source) ParentG(mods ...qm.QueryMod) sourceQuery {
	return o.Parent(boil.GetDB(), mods...)
}

// Parent pointed to by the foreign key.
func (o *Source) Parent(exec boil.Executor, mods ...qm.QueryMod) sourceQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.ParentID),
	}

	queryMods = append(queryMods, mods...)

	query := Sources(exec, queryMods...)
	queries.SetFrom(query.Query, "\"sources\"")

	return query
}

// TypeG pointed to by the foreign key.
func (o *Source) TypeG(mods ...qm.QueryMod) sourceTypeQuery {
	return o.Type(boil.GetDB(), mods...)
}

// Type pointed to by the foreign key.
func (o *Source) Type(exec boil.Executor, mods ...qm.QueryMod) sourceTypeQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.TypeID),
	}

	queryMods = append(queryMods, mods...)

	query := SourceTypes(exec, queryMods...)
	queries.SetFrom(query.Query, "\"source_types\"")

	return query
}

// ParentSourcesG retrieves all the source's sources via parent_id column.
func (o *Source) ParentSourcesG(mods ...qm.QueryMod) sourceQuery {
	return o.ParentSources(boil.GetDB(), mods...)
}

// ParentSources retrieves all the source's sources with an executor via parent_id column.
func (o *Source) ParentSources(exec boil.Executor, mods ...qm.QueryMod) sourceQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"parent_id\"=?", o.ID),
	)

	query := Sources(exec, queryMods...)
	queries.SetFrom(query.Query, "\"sources\" as \"a\"")
	return query
}

// ContentUnitsG retrieves all the content_unit's content units.
func (o *Source) ContentUnitsG(mods ...qm.QueryMod) contentUnitQuery {
	return o.ContentUnits(boil.GetDB(), mods...)
}

// ContentUnits retrieves all the content_unit's content units with an executor.
func (o *Source) ContentUnits(exec boil.Executor, mods ...qm.QueryMod) contentUnitQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.InnerJoin("\"content_units_sources\" as \"b\" on \"a\".\"id\" = \"b\".\"content_unit_id\""),
		qm.Where("\"b\".\"source_id\"=?", o.ID),
	)

	query := ContentUnits(exec, queryMods...)
	queries.SetFrom(query.Query, "\"content_units\" as \"a\"")
	return query
}

// AuthorsG retrieves all the author's authors.
func (o *Source) AuthorsG(mods ...qm.QueryMod) authorQuery {
	return o.Authors(boil.GetDB(), mods...)
}

// Authors retrieves all the author's authors with an executor.
func (o *Source) Authors(exec boil.Executor, mods ...qm.QueryMod) authorQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.InnerJoin("\"authors_sources\" as \"b\" on \"a\".\"id\" = \"b\".\"author_id\""),
		qm.Where("\"b\".\"source_id\"=?", o.ID),
	)

	query := Authors(exec, queryMods...)
	queries.SetFrom(query.Query, "\"authors\" as \"a\"")
	return query
}

// SourceI18nsG retrieves all the source_i18n's source i18n.
func (o *Source) SourceI18nsG(mods ...qm.QueryMod) sourceI18nQuery {
	return o.SourceI18ns(boil.GetDB(), mods...)
}

// SourceI18ns retrieves all the source_i18n's source i18n with an executor.
func (o *Source) SourceI18ns(exec boil.Executor, mods ...qm.QueryMod) sourceI18nQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"source_id\"=?", o.ID),
	)

	query := SourceI18ns(exec, queryMods...)
	queries.SetFrom(query.Query, "\"source_i18n\" as \"a\"")
	return query
}

// LoadParent allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (sourceL) LoadParent(e boil.Executor, singular bool, maybeSource interface{}) error {
	var slice []*Source
	var object *Source

	count := 1
	if singular {
		object = maybeSource.(*Source)
	} else {
		slice = *maybeSource.(*SourceSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &sourceR{}
		}
		args[0] = object.ParentID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &sourceR{}
			}
			args[i] = obj.ParentID
		}
	}

	query := fmt.Sprintf(
		"select * from \"sources\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Source")
	}
	defer results.Close()

	var resultSlice []*Source
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Source")
	}

	if singular && len(resultSlice) != 0 {
		object.R.Parent = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ParentID.Int64 == foreign.ID {
				local.R.Parent = foreign
				break
			}
		}
	}

	return nil
}

// LoadType allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (sourceL) LoadType(e boil.Executor, singular bool, maybeSource interface{}) error {
	var slice []*Source
	var object *Source

	count := 1
	if singular {
		object = maybeSource.(*Source)
	} else {
		slice = *maybeSource.(*SourceSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &sourceR{}
		}
		args[0] = object.TypeID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &sourceR{}
			}
			args[i] = obj.TypeID
		}
	}

	query := fmt.Sprintf(
		"select * from \"source_types\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load SourceType")
	}
	defer results.Close()

	var resultSlice []*SourceType
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice SourceType")
	}

	if singular && len(resultSlice) != 0 {
		object.R.Type = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.TypeID == foreign.ID {
				local.R.Type = foreign
				break
			}
		}
	}

	return nil
}

// LoadParentSources allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (sourceL) LoadParentSources(e boil.Executor, singular bool, maybeSource interface{}) error {
	var slice []*Source
	var object *Source

	count := 1
	if singular {
		object = maybeSource.(*Source)
	} else {
		slice = *maybeSource.(*SourceSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &sourceR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &sourceR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"sources\" where \"parent_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load sources")
	}
	defer results.Close()

	var resultSlice []*Source
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice sources")
	}

	if singular {
		object.R.ParentSources = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.ParentID.Int64 {
				local.R.ParentSources = append(local.R.ParentSources, foreign)
				break
			}
		}
	}

	return nil
}

// LoadContentUnits allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (sourceL) LoadContentUnits(e boil.Executor, singular bool, maybeSource interface{}) error {
	var slice []*Source
	var object *Source

	count := 1
	if singular {
		object = maybeSource.(*Source)
	} else {
		slice = *maybeSource.(*SourceSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &sourceR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &sourceR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select \"a\".*, \"b\".\"source_id\" from \"content_units\" as \"a\" inner join \"content_units_sources\" as \"b\" on \"a\".\"id\" = \"b\".\"content_unit_id\" where \"b\".\"source_id\" in (%s)",
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

	var localJoinCols []int64
	for results.Next() {
		one := new(ContentUnit)
		var localJoinCol int64

		err = results.Scan(&one.ID, &one.UID, &one.TypeID, &one.CreatedAt, &one.Properties, &localJoinCol)
		if err = results.Err(); err != nil {
			return errors.Wrap(err, "failed to plebian-bind eager loaded slice content_units")
		}

		resultSlice = append(resultSlice, one)
		localJoinCols = append(localJoinCols, localJoinCol)
	}

	if err = results.Err(); err != nil {
		return errors.Wrap(err, "failed to plebian-bind eager loaded slice content_units")
	}

	if singular {
		object.R.ContentUnits = resultSlice
		return nil
	}

	for i, foreign := range resultSlice {
		localJoinCol := localJoinCols[i]
		for _, local := range slice {
			if local.ID == localJoinCol {
				local.R.ContentUnits = append(local.R.ContentUnits, foreign)
				break
			}
		}
	}

	return nil
}

// LoadAuthors allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (sourceL) LoadAuthors(e boil.Executor, singular bool, maybeSource interface{}) error {
	var slice []*Source
	var object *Source

	count := 1
	if singular {
		object = maybeSource.(*Source)
	} else {
		slice = *maybeSource.(*SourceSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &sourceR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &sourceR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select \"a\".*, \"b\".\"source_id\" from \"authors\" as \"a\" inner join \"authors_sources\" as \"b\" on \"a\".\"id\" = \"b\".\"author_id\" where \"b\".\"source_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load authors")
	}
	defer results.Close()

	var resultSlice []*Author

	var localJoinCols []int64
	for results.Next() {
		one := new(Author)
		var localJoinCol int64

		err = results.Scan(&one.ID, &one.Code, &one.Name, &one.FullName, &one.CreatedAt, &localJoinCol)
		if err = results.Err(); err != nil {
			return errors.Wrap(err, "failed to plebian-bind eager loaded slice authors")
		}

		resultSlice = append(resultSlice, one)
		localJoinCols = append(localJoinCols, localJoinCol)
	}

	if err = results.Err(); err != nil {
		return errors.Wrap(err, "failed to plebian-bind eager loaded slice authors")
	}

	if singular {
		object.R.Authors = resultSlice
		return nil
	}

	for i, foreign := range resultSlice {
		localJoinCol := localJoinCols[i]
		for _, local := range slice {
			if local.ID == localJoinCol {
				local.R.Authors = append(local.R.Authors, foreign)
				break
			}
		}
	}

	return nil
}

// LoadSourceI18ns allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (sourceL) LoadSourceI18ns(e boil.Executor, singular bool, maybeSource interface{}) error {
	var slice []*Source
	var object *Source

	count := 1
	if singular {
		object = maybeSource.(*Source)
	} else {
		slice = *maybeSource.(*SourceSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &sourceR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &sourceR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"source_i18n\" where \"source_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load source_i18n")
	}
	defer results.Close()

	var resultSlice []*SourceI18n
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice source_i18n")
	}

	if singular {
		object.R.SourceI18ns = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.SourceID {
				local.R.SourceI18ns = append(local.R.SourceI18ns, foreign)
				break
			}
		}
	}

	return nil
}

// SetParentG of the source to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentSources.
// Uses the global database handle.
func (o *Source) SetParentG(insert bool, related *Source) error {
	return o.SetParent(boil.GetDB(), insert, related)
}

// SetParentP of the source to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentSources.
// Panics on error.
func (o *Source) SetParentP(exec boil.Executor, insert bool, related *Source) {
	if err := o.SetParent(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParentGP of the source to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentSources.
// Uses the global database handle and panics on error.
func (o *Source) SetParentGP(insert bool, related *Source) {
	if err := o.SetParent(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParent of the source to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentSources.
func (o *Source) SetParent(exec boil.Executor, insert bool, related *Source) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"sources\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"parent_id"}),
		strmangle.WhereClause("\"", "\"", 2, sourcePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ParentID.Int64 = related.ID
	o.ParentID.Valid = true

	if o.R == nil {
		o.R = &sourceR{
			Parent: related,
		}
	} else {
		o.R.Parent = related
	}

	if related.R == nil {
		related.R = &sourceR{
			ParentSources: SourceSlice{o},
		}
	} else {
		related.R.ParentSources = append(related.R.ParentSources, o)
	}

	return nil
}

// RemoveParentG relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle.
func (o *Source) RemoveParentG(related *Source) error {
	return o.RemoveParent(boil.GetDB(), related)
}

// RemoveParentP relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Panics on error.
func (o *Source) RemoveParentP(exec boil.Executor, related *Source) {
	if err := o.RemoveParent(exec, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParentGP relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle and panics on error.
func (o *Source) RemoveParentGP(related *Source) {
	if err := o.RemoveParent(boil.GetDB(), related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParent relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *Source) RemoveParent(exec boil.Executor, related *Source) error {
	var err error

	o.ParentID.Valid = false
	if err = o.Update(exec, "parent_id"); err != nil {
		o.ParentID.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.Parent = nil
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.ParentSources {
		if o.ParentID.Int64 != ri.ParentID.Int64 {
			continue
		}

		ln := len(related.R.ParentSources)
		if ln > 1 && i < ln-1 {
			related.R.ParentSources[i] = related.R.ParentSources[ln-1]
		}
		related.R.ParentSources = related.R.ParentSources[:ln-1]
		break
	}
	return nil
}

// SetTypeG of the source to the related item.
// Sets o.R.Type to related.
// Adds o to related.R.TypeSources.
// Uses the global database handle.
func (o *Source) SetTypeG(insert bool, related *SourceType) error {
	return o.SetType(boil.GetDB(), insert, related)
}

// SetTypeP of the source to the related item.
// Sets o.R.Type to related.
// Adds o to related.R.TypeSources.
// Panics on error.
func (o *Source) SetTypeP(exec boil.Executor, insert bool, related *SourceType) {
	if err := o.SetType(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetTypeGP of the source to the related item.
// Sets o.R.Type to related.
// Adds o to related.R.TypeSources.
// Uses the global database handle and panics on error.
func (o *Source) SetTypeGP(insert bool, related *SourceType) {
	if err := o.SetType(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetType of the source to the related item.
// Sets o.R.Type to related.
// Adds o to related.R.TypeSources.
func (o *Source) SetType(exec boil.Executor, insert bool, related *SourceType) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"sources\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"type_id"}),
		strmangle.WhereClause("\"", "\"", 2, sourcePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.TypeID = related.ID

	if o.R == nil {
		o.R = &sourceR{
			Type: related,
		}
	} else {
		o.R.Type = related
	}

	if related.R == nil {
		related.R = &sourceTypeR{
			TypeSources: SourceSlice{o},
		}
	} else {
		related.R.TypeSources = append(related.R.TypeSources, o)
	}

	return nil
}

// AddParentSourcesG adds the given related objects to the existing relationships
// of the source, optionally inserting them as new records.
// Appends related to o.R.ParentSources.
// Sets related.R.Parent appropriately.
// Uses the global database handle.
func (o *Source) AddParentSourcesG(insert bool, related ...*Source) error {
	return o.AddParentSources(boil.GetDB(), insert, related...)
}

// AddParentSourcesP adds the given related objects to the existing relationships
// of the source, optionally inserting them as new records.
// Appends related to o.R.ParentSources.
// Sets related.R.Parent appropriately.
// Panics on error.
func (o *Source) AddParentSourcesP(exec boil.Executor, insert bool, related ...*Source) {
	if err := o.AddParentSources(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddParentSourcesGP adds the given related objects to the existing relationships
// of the source, optionally inserting them as new records.
// Appends related to o.R.ParentSources.
// Sets related.R.Parent appropriately.
// Uses the global database handle and panics on error.
func (o *Source) AddParentSourcesGP(insert bool, related ...*Source) {
	if err := o.AddParentSources(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddParentSources adds the given related objects to the existing relationships
// of the source, optionally inserting them as new records.
// Appends related to o.R.ParentSources.
// Sets related.R.Parent appropriately.
func (o *Source) AddParentSources(exec boil.Executor, insert bool, related ...*Source) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.ParentID.Int64 = o.ID
			rel.ParentID.Valid = true
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"sources\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"parent_id"}),
				strmangle.WhereClause("\"", "\"", 2, sourcePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.ParentID.Int64 = o.ID
			rel.ParentID.Valid = true
		}
	}

	if o.R == nil {
		o.R = &sourceR{
			ParentSources: related,
		}
	} else {
		o.R.ParentSources = append(o.R.ParentSources, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &sourceR{
				Parent: o,
			}
		} else {
			rel.R.Parent = o
		}
	}
	return nil
}

// SetParentSourcesG removes all previously related items of the
// source replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentSources accordingly.
// Replaces o.R.ParentSources with related.
// Sets related.R.Parent's ParentSources accordingly.
// Uses the global database handle.
func (o *Source) SetParentSourcesG(insert bool, related ...*Source) error {
	return o.SetParentSources(boil.GetDB(), insert, related...)
}

// SetParentSourcesP removes all previously related items of the
// source replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentSources accordingly.
// Replaces o.R.ParentSources with related.
// Sets related.R.Parent's ParentSources accordingly.
// Panics on error.
func (o *Source) SetParentSourcesP(exec boil.Executor, insert bool, related ...*Source) {
	if err := o.SetParentSources(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParentSourcesGP removes all previously related items of the
// source replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentSources accordingly.
// Replaces o.R.ParentSources with related.
// Sets related.R.Parent's ParentSources accordingly.
// Uses the global database handle and panics on error.
func (o *Source) SetParentSourcesGP(insert bool, related ...*Source) {
	if err := o.SetParentSources(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParentSources removes all previously related items of the
// source replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentSources accordingly.
// Replaces o.R.ParentSources with related.
// Sets related.R.Parent's ParentSources accordingly.
func (o *Source) SetParentSources(exec boil.Executor, insert bool, related ...*Source) error {
	query := "update \"sources\" set \"parent_id\" = null where \"parent_id\" = $1"
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
		for _, rel := range o.R.ParentSources {
			rel.ParentID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Parent = nil
		}

		o.R.ParentSources = nil
	}
	return o.AddParentSources(exec, insert, related...)
}

// RemoveParentSourcesG relationships from objects passed in.
// Removes related items from R.ParentSources (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
// Uses the global database handle.
func (o *Source) RemoveParentSourcesG(related ...*Source) error {
	return o.RemoveParentSources(boil.GetDB(), related...)
}

// RemoveParentSourcesP relationships from objects passed in.
// Removes related items from R.ParentSources (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
// Panics on error.
func (o *Source) RemoveParentSourcesP(exec boil.Executor, related ...*Source) {
	if err := o.RemoveParentSources(exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParentSourcesGP relationships from objects passed in.
// Removes related items from R.ParentSources (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
// Uses the global database handle and panics on error.
func (o *Source) RemoveParentSourcesGP(related ...*Source) {
	if err := o.RemoveParentSources(boil.GetDB(), related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParentSources relationships from objects passed in.
// Removes related items from R.ParentSources (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
func (o *Source) RemoveParentSources(exec boil.Executor, related ...*Source) error {
	var err error
	for _, rel := range related {
		rel.ParentID.Valid = false
		if rel.R != nil {
			rel.R.Parent = nil
		}
		if err = rel.Update(exec, "parent_id"); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.ParentSources {
			if rel != ri {
				continue
			}

			ln := len(o.R.ParentSources)
			if ln > 1 && i < ln-1 {
				o.R.ParentSources[i] = o.R.ParentSources[ln-1]
			}
			o.R.ParentSources = o.R.ParentSources[:ln-1]
			break
		}
	}

	return nil
}

// AddContentUnitsG adds the given related objects to the existing relationships
// of the source, optionally inserting them as new records.
// Appends related to o.R.ContentUnits.
// Sets related.R.Sources appropriately.
// Uses the global database handle.
func (o *Source) AddContentUnitsG(insert bool, related ...*ContentUnit) error {
	return o.AddContentUnits(boil.GetDB(), insert, related...)
}

// AddContentUnitsP adds the given related objects to the existing relationships
// of the source, optionally inserting them as new records.
// Appends related to o.R.ContentUnits.
// Sets related.R.Sources appropriately.
// Panics on error.
func (o *Source) AddContentUnitsP(exec boil.Executor, insert bool, related ...*ContentUnit) {
	if err := o.AddContentUnits(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddContentUnitsGP adds the given related objects to the existing relationships
// of the source, optionally inserting them as new records.
// Appends related to o.R.ContentUnits.
// Sets related.R.Sources appropriately.
// Uses the global database handle and panics on error.
func (o *Source) AddContentUnitsGP(insert bool, related ...*ContentUnit) {
	if err := o.AddContentUnits(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddContentUnits adds the given related objects to the existing relationships
// of the source, optionally inserting them as new records.
// Appends related to o.R.ContentUnits.
// Sets related.R.Sources appropriately.
func (o *Source) AddContentUnits(exec boil.Executor, insert bool, related ...*ContentUnit) error {
	var err error
	for _, rel := range related {
		if insert {
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		}
	}

	for _, rel := range related {
		query := "insert into \"content_units_sources\" (\"source_id\", \"content_unit_id\") values ($1, $2)"
		values := []interface{}{o.ID, rel.ID}

		if boil.DebugMode {
			fmt.Fprintln(boil.DebugWriter, query)
			fmt.Fprintln(boil.DebugWriter, values)
		}

		_, err = exec.Exec(query, values...)
		if err != nil {
			return errors.Wrap(err, "failed to insert into join table")
		}
	}
	if o.R == nil {
		o.R = &sourceR{
			ContentUnits: related,
		}
	} else {
		o.R.ContentUnits = append(o.R.ContentUnits, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &contentUnitR{
				Sources: SourceSlice{o},
			}
		} else {
			rel.R.Sources = append(rel.R.Sources, o)
		}
	}
	return nil
}

// SetContentUnitsG removes all previously related items of the
// source replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Sources's ContentUnits accordingly.
// Replaces o.R.ContentUnits with related.
// Sets related.R.Sources's ContentUnits accordingly.
// Uses the global database handle.
func (o *Source) SetContentUnitsG(insert bool, related ...*ContentUnit) error {
	return o.SetContentUnits(boil.GetDB(), insert, related...)
}

// SetContentUnitsP removes all previously related items of the
// source replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Sources's ContentUnits accordingly.
// Replaces o.R.ContentUnits with related.
// Sets related.R.Sources's ContentUnits accordingly.
// Panics on error.
func (o *Source) SetContentUnitsP(exec boil.Executor, insert bool, related ...*ContentUnit) {
	if err := o.SetContentUnits(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetContentUnitsGP removes all previously related items of the
// source replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Sources's ContentUnits accordingly.
// Replaces o.R.ContentUnits with related.
// Sets related.R.Sources's ContentUnits accordingly.
// Uses the global database handle and panics on error.
func (o *Source) SetContentUnitsGP(insert bool, related ...*ContentUnit) {
	if err := o.SetContentUnits(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetContentUnits removes all previously related items of the
// source replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Sources's ContentUnits accordingly.
// Replaces o.R.ContentUnits with related.
// Sets related.R.Sources's ContentUnits accordingly.
func (o *Source) SetContentUnits(exec boil.Executor, insert bool, related ...*ContentUnit) error {
	query := "delete from \"content_units_sources\" where \"source_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	removeContentUnitsFromSourcesSlice(o, related)
	o.R.ContentUnits = nil
	return o.AddContentUnits(exec, insert, related...)
}

// RemoveContentUnitsG relationships from objects passed in.
// Removes related items from R.ContentUnits (uses pointer comparison, removal does not keep order)
// Sets related.R.Sources.
// Uses the global database handle.
func (o *Source) RemoveContentUnitsG(related ...*ContentUnit) error {
	return o.RemoveContentUnits(boil.GetDB(), related...)
}

// RemoveContentUnitsP relationships from objects passed in.
// Removes related items from R.ContentUnits (uses pointer comparison, removal does not keep order)
// Sets related.R.Sources.
// Panics on error.
func (o *Source) RemoveContentUnitsP(exec boil.Executor, related ...*ContentUnit) {
	if err := o.RemoveContentUnits(exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveContentUnitsGP relationships from objects passed in.
// Removes related items from R.ContentUnits (uses pointer comparison, removal does not keep order)
// Sets related.R.Sources.
// Uses the global database handle and panics on error.
func (o *Source) RemoveContentUnitsGP(related ...*ContentUnit) {
	if err := o.RemoveContentUnits(boil.GetDB(), related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveContentUnits relationships from objects passed in.
// Removes related items from R.ContentUnits (uses pointer comparison, removal does not keep order)
// Sets related.R.Sources.
func (o *Source) RemoveContentUnits(exec boil.Executor, related ...*ContentUnit) error {
	var err error
	query := fmt.Sprintf(
		"delete from \"content_units_sources\" where \"source_id\" = $1 and \"content_unit_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, len(related), 1, 1),
	)
	values := []interface{}{o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}
	removeContentUnitsFromSourcesSlice(o, related)
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.ContentUnits {
			if rel != ri {
				continue
			}

			ln := len(o.R.ContentUnits)
			if ln > 1 && i < ln-1 {
				o.R.ContentUnits[i] = o.R.ContentUnits[ln-1]
			}
			o.R.ContentUnits = o.R.ContentUnits[:ln-1]
			break
		}
	}

	return nil
}

func removeContentUnitsFromSourcesSlice(o *Source, related []*ContentUnit) {
	for _, rel := range related {
		if rel.R == nil {
			continue
		}
		for i, ri := range rel.R.Sources {
			if o.ID != ri.ID {
				continue
			}

			ln := len(rel.R.Sources)
			if ln > 1 && i < ln-1 {
				rel.R.Sources[i] = rel.R.Sources[ln-1]
			}
			rel.R.Sources = rel.R.Sources[:ln-1]
			break
		}
	}
}

// AddAuthorsG adds the given related objects to the existing relationships
// of the source, optionally inserting them as new records.
// Appends related to o.R.Authors.
// Sets related.R.Sources appropriately.
// Uses the global database handle.
func (o *Source) AddAuthorsG(insert bool, related ...*Author) error {
	return o.AddAuthors(boil.GetDB(), insert, related...)
}

// AddAuthorsP adds the given related objects to the existing relationships
// of the source, optionally inserting them as new records.
// Appends related to o.R.Authors.
// Sets related.R.Sources appropriately.
// Panics on error.
func (o *Source) AddAuthorsP(exec boil.Executor, insert bool, related ...*Author) {
	if err := o.AddAuthors(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddAuthorsGP adds the given related objects to the existing relationships
// of the source, optionally inserting them as new records.
// Appends related to o.R.Authors.
// Sets related.R.Sources appropriately.
// Uses the global database handle and panics on error.
func (o *Source) AddAuthorsGP(insert bool, related ...*Author) {
	if err := o.AddAuthors(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddAuthors adds the given related objects to the existing relationships
// of the source, optionally inserting them as new records.
// Appends related to o.R.Authors.
// Sets related.R.Sources appropriately.
func (o *Source) AddAuthors(exec boil.Executor, insert bool, related ...*Author) error {
	var err error
	for _, rel := range related {
		if insert {
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		}
	}

	for _, rel := range related {
		query := "insert into \"authors_sources\" (\"source_id\", \"author_id\") values ($1, $2)"
		values := []interface{}{o.ID, rel.ID}

		if boil.DebugMode {
			fmt.Fprintln(boil.DebugWriter, query)
			fmt.Fprintln(boil.DebugWriter, values)
		}

		_, err = exec.Exec(query, values...)
		if err != nil {
			return errors.Wrap(err, "failed to insert into join table")
		}
	}
	if o.R == nil {
		o.R = &sourceR{
			Authors: related,
		}
	} else {
		o.R.Authors = append(o.R.Authors, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &authorR{
				Sources: SourceSlice{o},
			}
		} else {
			rel.R.Sources = append(rel.R.Sources, o)
		}
	}
	return nil
}

// SetAuthorsG removes all previously related items of the
// source replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Sources's Authors accordingly.
// Replaces o.R.Authors with related.
// Sets related.R.Sources's Authors accordingly.
// Uses the global database handle.
func (o *Source) SetAuthorsG(insert bool, related ...*Author) error {
	return o.SetAuthors(boil.GetDB(), insert, related...)
}

// SetAuthorsP removes all previously related items of the
// source replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Sources's Authors accordingly.
// Replaces o.R.Authors with related.
// Sets related.R.Sources's Authors accordingly.
// Panics on error.
func (o *Source) SetAuthorsP(exec boil.Executor, insert bool, related ...*Author) {
	if err := o.SetAuthors(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetAuthorsGP removes all previously related items of the
// source replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Sources's Authors accordingly.
// Replaces o.R.Authors with related.
// Sets related.R.Sources's Authors accordingly.
// Uses the global database handle and panics on error.
func (o *Source) SetAuthorsGP(insert bool, related ...*Author) {
	if err := o.SetAuthors(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetAuthors removes all previously related items of the
// source replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Sources's Authors accordingly.
// Replaces o.R.Authors with related.
// Sets related.R.Sources's Authors accordingly.
func (o *Source) SetAuthors(exec boil.Executor, insert bool, related ...*Author) error {
	query := "delete from \"authors_sources\" where \"source_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	removeAuthorsFromSourcesSlice(o, related)
	o.R.Authors = nil
	return o.AddAuthors(exec, insert, related...)
}

// RemoveAuthorsG relationships from objects passed in.
// Removes related items from R.Authors (uses pointer comparison, removal does not keep order)
// Sets related.R.Sources.
// Uses the global database handle.
func (o *Source) RemoveAuthorsG(related ...*Author) error {
	return o.RemoveAuthors(boil.GetDB(), related...)
}

// RemoveAuthorsP relationships from objects passed in.
// Removes related items from R.Authors (uses pointer comparison, removal does not keep order)
// Sets related.R.Sources.
// Panics on error.
func (o *Source) RemoveAuthorsP(exec boil.Executor, related ...*Author) {
	if err := o.RemoveAuthors(exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveAuthorsGP relationships from objects passed in.
// Removes related items from R.Authors (uses pointer comparison, removal does not keep order)
// Sets related.R.Sources.
// Uses the global database handle and panics on error.
func (o *Source) RemoveAuthorsGP(related ...*Author) {
	if err := o.RemoveAuthors(boil.GetDB(), related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveAuthors relationships from objects passed in.
// Removes related items from R.Authors (uses pointer comparison, removal does not keep order)
// Sets related.R.Sources.
func (o *Source) RemoveAuthors(exec boil.Executor, related ...*Author) error {
	var err error
	query := fmt.Sprintf(
		"delete from \"authors_sources\" where \"source_id\" = $1 and \"author_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, len(related), 1, 1),
	)
	values := []interface{}{o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}
	removeAuthorsFromSourcesSlice(o, related)
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Authors {
			if rel != ri {
				continue
			}

			ln := len(o.R.Authors)
			if ln > 1 && i < ln-1 {
				o.R.Authors[i] = o.R.Authors[ln-1]
			}
			o.R.Authors = o.R.Authors[:ln-1]
			break
		}
	}

	return nil
}

func removeAuthorsFromSourcesSlice(o *Source, related []*Author) {
	for _, rel := range related {
		if rel.R == nil {
			continue
		}
		for i, ri := range rel.R.Sources {
			if o.ID != ri.ID {
				continue
			}

			ln := len(rel.R.Sources)
			if ln > 1 && i < ln-1 {
				rel.R.Sources[i] = rel.R.Sources[ln-1]
			}
			rel.R.Sources = rel.R.Sources[:ln-1]
			break
		}
	}
}

// AddSourceI18nsG adds the given related objects to the existing relationships
// of the source, optionally inserting them as new records.
// Appends related to o.R.SourceI18ns.
// Sets related.R.Source appropriately.
// Uses the global database handle.
func (o *Source) AddSourceI18nsG(insert bool, related ...*SourceI18n) error {
	return o.AddSourceI18ns(boil.GetDB(), insert, related...)
}

// AddSourceI18nsP adds the given related objects to the existing relationships
// of the source, optionally inserting them as new records.
// Appends related to o.R.SourceI18ns.
// Sets related.R.Source appropriately.
// Panics on error.
func (o *Source) AddSourceI18nsP(exec boil.Executor, insert bool, related ...*SourceI18n) {
	if err := o.AddSourceI18ns(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddSourceI18nsGP adds the given related objects to the existing relationships
// of the source, optionally inserting them as new records.
// Appends related to o.R.SourceI18ns.
// Sets related.R.Source appropriately.
// Uses the global database handle and panics on error.
func (o *Source) AddSourceI18nsGP(insert bool, related ...*SourceI18n) {
	if err := o.AddSourceI18ns(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddSourceI18ns adds the given related objects to the existing relationships
// of the source, optionally inserting them as new records.
// Appends related to o.R.SourceI18ns.
// Sets related.R.Source appropriately.
func (o *Source) AddSourceI18ns(exec boil.Executor, insert bool, related ...*SourceI18n) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.SourceID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"source_i18n\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"source_id"}),
				strmangle.WhereClause("\"", "\"", 2, sourceI18nPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.SourceID, rel.Language}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.SourceID = o.ID
		}
	}

	if o.R == nil {
		o.R = &sourceR{
			SourceI18ns: related,
		}
	} else {
		o.R.SourceI18ns = append(o.R.SourceI18ns, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &sourceI18nR{
				Source: o,
			}
		} else {
			rel.R.Source = o
		}
	}
	return nil
}

// SourcesG retrieves all records.
func SourcesG(mods ...qm.QueryMod) sourceQuery {
	return Sources(boil.GetDB(), mods...)
}

// Sources retrieves all the records using an executor.
func Sources(exec boil.Executor, mods ...qm.QueryMod) sourceQuery {
	mods = append(mods, qm.From("\"sources\""))
	return sourceQuery{NewQuery(exec, mods...)}
}

// FindSourceG retrieves a single record by ID.
func FindSourceG(id int64, selectCols ...string) (*Source, error) {
	return FindSource(boil.GetDB(), id, selectCols...)
}

// FindSourceGP retrieves a single record by ID, and panics on error.
func FindSourceGP(id int64, selectCols ...string) *Source {
	retobj, err := FindSource(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindSource retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSource(exec boil.Executor, id int64, selectCols ...string) (*Source, error) {
	sourceObj := &Source{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"sources\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(sourceObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from sources")
	}

	return sourceObj, nil
}

// FindSourceP retrieves a single record by ID with an executor, and panics on error.
func FindSourceP(exec boil.Executor, id int64, selectCols ...string) *Source {
	retobj, err := FindSource(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Source) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Source) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Source) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Source) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no sources provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(sourceColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	sourceInsertCacheMut.RLock()
	cache, cached := sourceInsertCache[key]
	sourceInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			sourceColumns,
			sourceColumnsWithDefault,
			sourceColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(sourceType, sourceMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(sourceType, sourceMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"sources\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into sources")
	}

	if !cached {
		sourceInsertCacheMut.Lock()
		sourceInsertCache[key] = cache
		sourceInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single Source record. See Update for
// whitelist behavior description.
func (o *Source) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Source record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Source) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Source, and panics on error.
// See Update for whitelist behavior description.
func (o *Source) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Source.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Source) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	sourceUpdateCacheMut.RLock()
	cache, cached := sourceUpdateCache[key]
	sourceUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(sourceColumns, sourcePrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update sources, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"sources\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, sourcePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(sourceType, sourceMapping, append(wl, sourcePrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update sources row")
	}

	if !cached {
		sourceUpdateCacheMut.Lock()
		sourceUpdateCache[key] = cache
		sourceUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q sourceQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q sourceQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for sources")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o SourceSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o SourceSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o SourceSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SourceSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sourcePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"sources\" SET %s WHERE (\"id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(sourcePrimaryKeyColumns), len(colNames)+1, len(sourcePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in source slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Source) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Source) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Source) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Source) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no sources provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(sourceColumnsWithDefault, o)

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

	sourceUpsertCacheMut.RLock()
	cache, cached := sourceUpsertCache[key]
	sourceUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			sourceColumns,
			sourceColumnsWithDefault,
			sourceColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			sourceColumns,
			sourcePrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert sources, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(sourcePrimaryKeyColumns))
			copy(conflict, sourcePrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"sources\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(sourceType, sourceMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(sourceType, sourceMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert sources")
	}

	if !cached {
		sourceUpsertCacheMut.Lock()
		sourceUpsertCache[key] = cache
		sourceUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single Source record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Source) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Source record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Source) DeleteG() error {
	if o == nil {
		return errors.New("models: no Source provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Source record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Source) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Source record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Source) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Source provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), sourcePrimaryKeyMapping)
	sql := "DELETE FROM \"sources\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from sources")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q sourceQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q sourceQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no sourceQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from sources")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o SourceSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o SourceSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Source slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o SourceSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SourceSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Source slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sourcePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"sources\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, sourcePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(sourcePrimaryKeyColumns), 1, len(sourcePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from source slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Source) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Source) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Source) ReloadG() error {
	if o == nil {
		return errors.New("models: no Source provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Source) Reload(exec boil.Executor) error {
	ret, err := FindSource(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *SourceSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *SourceSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SourceSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty SourceSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SourceSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	sources := SourceSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sourcePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"sources\".* FROM \"sources\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, sourcePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(sourcePrimaryKeyColumns), 1, len(sourcePrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&sources)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in SourceSlice")
	}

	*o = sources

	return nil
}

// SourceExists checks if the Source row exists.
func SourceExists(exec boil.Executor, id int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"sources\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if sources exists")
	}

	return exists, nil
}

// SourceExistsG checks if the Source row exists.
func SourceExistsG(id int64) (bool, error) {
	return SourceExists(boil.GetDB(), id)
}

// SourceExistsGP checks if the Source row exists. Panics on error.
func SourceExistsGP(id int64) bool {
	e, err := SourceExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// SourceExistsP checks if the Source row exists. Panics on error.
func SourceExistsP(exec boil.Executor, id int64) bool {
	e, err := SourceExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}