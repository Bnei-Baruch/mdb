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
)

// SourceType is an object representing the database table.
type SourceType struct {
	ID   int64  `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name string `boil:"name" json:"name" toml:"name" yaml:"name"`

	R *sourceTypeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L sourceTypeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// sourceTypeR is where relationships are stored.
type sourceTypeR struct {
	TypeSources SourceSlice
}

// sourceTypeL is where Load methods for each relationship are stored.
type sourceTypeL struct{}

var (
	sourceTypeColumns               = []string{"id", "name"}
	sourceTypeColumnsWithoutDefault = []string{"name"}
	sourceTypeColumnsWithDefault    = []string{"id"}
	sourceTypePrimaryKeyColumns     = []string{"id"}
)

type (
	// SourceTypeSlice is an alias for a slice of pointers to SourceType.
	// This should generally be used opposed to []SourceType.
	SourceTypeSlice []*SourceType

	sourceTypeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	sourceTypeType                 = reflect.TypeOf(&SourceType{})
	sourceTypeMapping              = queries.MakeStructMapping(sourceTypeType)
	sourceTypePrimaryKeyMapping, _ = queries.BindMapping(sourceTypeType, sourceTypeMapping, sourceTypePrimaryKeyColumns)
	sourceTypeInsertCacheMut       sync.RWMutex
	sourceTypeInsertCache          = make(map[string]insertCache)
	sourceTypeUpdateCacheMut       sync.RWMutex
	sourceTypeUpdateCache          = make(map[string]updateCache)
	sourceTypeUpsertCacheMut       sync.RWMutex
	sourceTypeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single sourceType record from the query, and panics on error.
func (q sourceTypeQuery) OneP() *SourceType {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single sourceType record from the query.
func (q sourceTypeQuery) One() (*SourceType, error) {
	o := &SourceType{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for source_types")
	}

	return o, nil
}

// AllP returns all SourceType records from the query, and panics on error.
func (q sourceTypeQuery) AllP() SourceTypeSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all SourceType records from the query.
func (q sourceTypeQuery) All() (SourceTypeSlice, error) {
	var o SourceTypeSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to SourceType slice")
	}

	return o, nil
}

// CountP returns the count of all SourceType records in the query, and panics on error.
func (q sourceTypeQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all SourceType records in the query.
func (q sourceTypeQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count source_types rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q sourceTypeQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q sourceTypeQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if source_types exists")
	}

	return count > 0, nil
}

// TypeSourcesG retrieves all the source's sources via type_id column.
func (o *SourceType) TypeSourcesG(mods ...qm.QueryMod) sourceQuery {
	return o.TypeSources(boil.GetDB(), mods...)
}

// TypeSources retrieves all the source's sources with an executor via type_id column.
func (o *SourceType) TypeSources(exec boil.Executor, mods ...qm.QueryMod) sourceQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"type_id\"=?", o.ID),
	)

	query := Sources(exec, queryMods...)
	queries.SetFrom(query.Query, "\"sources\" as \"a\"")
	return query
}

// LoadTypeSources allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (sourceTypeL) LoadTypeSources(e boil.Executor, singular bool, maybeSourceType interface{}) error {
	var slice []*SourceType
	var object *SourceType

	count := 1
	if singular {
		object = maybeSourceType.(*SourceType)
	} else {
		slice = *maybeSourceType.(*SourceTypeSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &sourceTypeR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &sourceTypeR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"sources\" where \"type_id\" in (%s)",
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
		object.R.TypeSources = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.TypeID {
				local.R.TypeSources = append(local.R.TypeSources, foreign)
				break
			}
		}
	}

	return nil
}

// AddTypeSourcesG adds the given related objects to the existing relationships
// of the source_type, optionally inserting them as new records.
// Appends related to o.R.TypeSources.
// Sets related.R.Type appropriately.
// Uses the global database handle.
func (o *SourceType) AddTypeSourcesG(insert bool, related ...*Source) error {
	return o.AddTypeSources(boil.GetDB(), insert, related...)
}

// AddTypeSourcesP adds the given related objects to the existing relationships
// of the source_type, optionally inserting them as new records.
// Appends related to o.R.TypeSources.
// Sets related.R.Type appropriately.
// Panics on error.
func (o *SourceType) AddTypeSourcesP(exec boil.Executor, insert bool, related ...*Source) {
	if err := o.AddTypeSources(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddTypeSourcesGP adds the given related objects to the existing relationships
// of the source_type, optionally inserting them as new records.
// Appends related to o.R.TypeSources.
// Sets related.R.Type appropriately.
// Uses the global database handle and panics on error.
func (o *SourceType) AddTypeSourcesGP(insert bool, related ...*Source) {
	if err := o.AddTypeSources(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddTypeSources adds the given related objects to the existing relationships
// of the source_type, optionally inserting them as new records.
// Appends related to o.R.TypeSources.
// Sets related.R.Type appropriately.
func (o *SourceType) AddTypeSources(exec boil.Executor, insert bool, related ...*Source) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.TypeID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"sources\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"type_id"}),
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

			rel.TypeID = o.ID
		}
	}

	if o.R == nil {
		o.R = &sourceTypeR{
			TypeSources: related,
		}
	} else {
		o.R.TypeSources = append(o.R.TypeSources, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &sourceR{
				Type: o,
			}
		} else {
			rel.R.Type = o
		}
	}
	return nil
}

// SourceTypesG retrieves all records.
func SourceTypesG(mods ...qm.QueryMod) sourceTypeQuery {
	return SourceTypes(boil.GetDB(), mods...)
}

// SourceTypes retrieves all the records using an executor.
func SourceTypes(exec boil.Executor, mods ...qm.QueryMod) sourceTypeQuery {
	mods = append(mods, qm.From("\"source_types\""))
	return sourceTypeQuery{NewQuery(exec, mods...)}
}

// FindSourceTypeG retrieves a single record by ID.
func FindSourceTypeG(id int64, selectCols ...string) (*SourceType, error) {
	return FindSourceType(boil.GetDB(), id, selectCols...)
}

// FindSourceTypeGP retrieves a single record by ID, and panics on error.
func FindSourceTypeGP(id int64, selectCols ...string) *SourceType {
	retobj, err := FindSourceType(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindSourceType retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSourceType(exec boil.Executor, id int64, selectCols ...string) (*SourceType, error) {
	sourceTypeObj := &SourceType{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"source_types\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(sourceTypeObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from source_types")
	}

	return sourceTypeObj, nil
}

// FindSourceTypeP retrieves a single record by ID with an executor, and panics on error.
func FindSourceTypeP(exec boil.Executor, id int64, selectCols ...string) *SourceType {
	retobj, err := FindSourceType(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *SourceType) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *SourceType) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *SourceType) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *SourceType) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no source_types provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(sourceTypeColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	sourceTypeInsertCacheMut.RLock()
	cache, cached := sourceTypeInsertCache[key]
	sourceTypeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			sourceTypeColumns,
			sourceTypeColumnsWithDefault,
			sourceTypeColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(sourceTypeType, sourceTypeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(sourceTypeType, sourceTypeMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"source_types\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into source_types")
	}

	if !cached {
		sourceTypeInsertCacheMut.Lock()
		sourceTypeInsertCache[key] = cache
		sourceTypeInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single SourceType record. See Update for
// whitelist behavior description.
func (o *SourceType) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single SourceType record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *SourceType) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the SourceType, and panics on error.
// See Update for whitelist behavior description.
func (o *SourceType) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the SourceType.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *SourceType) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	sourceTypeUpdateCacheMut.RLock()
	cache, cached := sourceTypeUpdateCache[key]
	sourceTypeUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(sourceTypeColumns, sourceTypePrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update source_types, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"source_types\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, sourceTypePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(sourceTypeType, sourceTypeMapping, append(wl, sourceTypePrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update source_types row")
	}

	if !cached {
		sourceTypeUpdateCacheMut.Lock()
		sourceTypeUpdateCache[key] = cache
		sourceTypeUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q sourceTypeQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q sourceTypeQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for source_types")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o SourceTypeSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o SourceTypeSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o SourceTypeSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SourceTypeSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sourceTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"source_types\" SET %s WHERE (\"id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(sourceTypePrimaryKeyColumns), len(colNames)+1, len(sourceTypePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in sourceType slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *SourceType) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *SourceType) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *SourceType) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *SourceType) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no source_types provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(sourceTypeColumnsWithDefault, o)

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

	sourceTypeUpsertCacheMut.RLock()
	cache, cached := sourceTypeUpsertCache[key]
	sourceTypeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			sourceTypeColumns,
			sourceTypeColumnsWithDefault,
			sourceTypeColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			sourceTypeColumns,
			sourceTypePrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert source_types, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(sourceTypePrimaryKeyColumns))
			copy(conflict, sourceTypePrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"source_types\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(sourceTypeType, sourceTypeMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(sourceTypeType, sourceTypeMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert source_types")
	}

	if !cached {
		sourceTypeUpsertCacheMut.Lock()
		sourceTypeUpsertCache[key] = cache
		sourceTypeUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single SourceType record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *SourceType) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single SourceType record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *SourceType) DeleteG() error {
	if o == nil {
		return errors.New("models: no SourceType provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single SourceType record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *SourceType) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single SourceType record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *SourceType) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no SourceType provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), sourceTypePrimaryKeyMapping)
	sql := "DELETE FROM \"source_types\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from source_types")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q sourceTypeQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q sourceTypeQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no sourceTypeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from source_types")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o SourceTypeSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o SourceTypeSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no SourceType slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o SourceTypeSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SourceTypeSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no SourceType slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sourceTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"source_types\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, sourceTypePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(sourceTypePrimaryKeyColumns), 1, len(sourceTypePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from sourceType slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *SourceType) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *SourceType) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *SourceType) ReloadG() error {
	if o == nil {
		return errors.New("models: no SourceType provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *SourceType) Reload(exec boil.Executor) error {
	ret, err := FindSourceType(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *SourceTypeSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *SourceTypeSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SourceTypeSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty SourceTypeSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SourceTypeSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	sourceTypes := SourceTypeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sourceTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"source_types\".* FROM \"source_types\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, sourceTypePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(sourceTypePrimaryKeyColumns), 1, len(sourceTypePrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&sourceTypes)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in SourceTypeSlice")
	}

	*o = sourceTypes

	return nil
}

// SourceTypeExists checks if the SourceType row exists.
func SourceTypeExists(exec boil.Executor, id int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"source_types\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if source_types exists")
	}

	return exists, nil
}

// SourceTypeExistsG checks if the SourceType row exists.
func SourceTypeExistsG(id int64) (bool, error) {
	return SourceTypeExists(boil.GetDB(), id)
}

// SourceTypeExistsGP checks if the SourceType row exists. Panics on error.
func SourceTypeExistsGP(id int64) bool {
	e, err := SourceTypeExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// SourceTypeExistsP checks if the SourceType row exists. Panics on error.
func SourceTypeExistsP(exec boil.Executor, id int64) bool {
	e, err := SourceTypeExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}