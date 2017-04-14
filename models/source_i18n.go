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

// SourceI18n is an object representing the database table.
type SourceI18n struct {
	SourceID    int64       `boil:"source_id" json:"source_id" toml:"source_id" yaml:"source_id"`
	Language    string      `boil:"language" json:"language" toml:"language" yaml:"language"`
	Name        null.String `boil:"name" json:"name,omitempty" toml:"name" yaml:"name,omitempty"`
	Description null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`
	CreatedAt   time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *sourceI18nR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L sourceI18nL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// sourceI18nR is where relationships are stored.
type sourceI18nR struct {
	Source *Source
}

// sourceI18nL is where Load methods for each relationship are stored.
type sourceI18nL struct{}

var (
	sourceI18nColumns               = []string{"source_id", "language", "name", "description", "created_at"}
	sourceI18nColumnsWithoutDefault = []string{"source_id", "language", "name", "description"}
	sourceI18nColumnsWithDefault    = []string{"created_at"}
	sourceI18nPrimaryKeyColumns     = []string{"source_id", "language"}
)

type (
	// SourceI18nSlice is an alias for a slice of pointers to SourceI18n.
	// This should generally be used opposed to []SourceI18n.
	SourceI18nSlice []*SourceI18n

	sourceI18nQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	sourceI18nType                 = reflect.TypeOf(&SourceI18n{})
	sourceI18nMapping              = queries.MakeStructMapping(sourceI18nType)
	sourceI18nPrimaryKeyMapping, _ = queries.BindMapping(sourceI18nType, sourceI18nMapping, sourceI18nPrimaryKeyColumns)
	sourceI18nInsertCacheMut       sync.RWMutex
	sourceI18nInsertCache          = make(map[string]insertCache)
	sourceI18nUpdateCacheMut       sync.RWMutex
	sourceI18nUpdateCache          = make(map[string]updateCache)
	sourceI18nUpsertCacheMut       sync.RWMutex
	sourceI18nUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single sourceI18n record from the query, and panics on error.
func (q sourceI18nQuery) OneP() *SourceI18n {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single sourceI18n record from the query.
func (q sourceI18nQuery) One() (*SourceI18n, error) {
	o := &SourceI18n{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for source_i18n")
	}

	return o, nil
}

// AllP returns all SourceI18n records from the query, and panics on error.
func (q sourceI18nQuery) AllP() SourceI18nSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all SourceI18n records from the query.
func (q sourceI18nQuery) All() (SourceI18nSlice, error) {
	var o SourceI18nSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to SourceI18n slice")
	}

	return o, nil
}

// CountP returns the count of all SourceI18n records in the query, and panics on error.
func (q sourceI18nQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all SourceI18n records in the query.
func (q sourceI18nQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count source_i18n rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q sourceI18nQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q sourceI18nQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if source_i18n exists")
	}

	return count > 0, nil
}

// SourceG pointed to by the foreign key.
func (o *SourceI18n) SourceG(mods ...qm.QueryMod) sourceQuery {
	return o.Source(boil.GetDB(), mods...)
}

// Source pointed to by the foreign key.
func (o *SourceI18n) Source(exec boil.Executor, mods ...qm.QueryMod) sourceQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.SourceID),
	}

	queryMods = append(queryMods, mods...)

	query := Sources(exec, queryMods...)
	queries.SetFrom(query.Query, "\"sources\"")

	return query
}

// LoadSource allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (sourceI18nL) LoadSource(e boil.Executor, singular bool, maybeSourceI18n interface{}) error {
	var slice []*SourceI18n
	var object *SourceI18n

	count := 1
	if singular {
		object = maybeSourceI18n.(*SourceI18n)
	} else {
		slice = *maybeSourceI18n.(*SourceI18nSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &sourceI18nR{}
		}
		args[0] = object.SourceID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &sourceI18nR{}
			}
			args[i] = obj.SourceID
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
		object.R.Source = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.SourceID == foreign.ID {
				local.R.Source = foreign
				break
			}
		}
	}

	return nil
}

// SetSourceG of the source_i18n to the related item.
// Sets o.R.Source to related.
// Adds o to related.R.SourceI18ns.
// Uses the global database handle.
func (o *SourceI18n) SetSourceG(insert bool, related *Source) error {
	return o.SetSource(boil.GetDB(), insert, related)
}

// SetSourceP of the source_i18n to the related item.
// Sets o.R.Source to related.
// Adds o to related.R.SourceI18ns.
// Panics on error.
func (o *SourceI18n) SetSourceP(exec boil.Executor, insert bool, related *Source) {
	if err := o.SetSource(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetSourceGP of the source_i18n to the related item.
// Sets o.R.Source to related.
// Adds o to related.R.SourceI18ns.
// Uses the global database handle and panics on error.
func (o *SourceI18n) SetSourceGP(insert bool, related *Source) {
	if err := o.SetSource(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetSource of the source_i18n to the related item.
// Sets o.R.Source to related.
// Adds o to related.R.SourceI18ns.
func (o *SourceI18n) SetSource(exec boil.Executor, insert bool, related *Source) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"source_i18n\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"source_id"}),
		strmangle.WhereClause("\"", "\"", 2, sourceI18nPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.SourceID, o.Language}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.SourceID = related.ID

	if o.R == nil {
		o.R = &sourceI18nR{
			Source: related,
		}
	} else {
		o.R.Source = related
	}

	if related.R == nil {
		related.R = &sourceR{
			SourceI18ns: SourceI18nSlice{o},
		}
	} else {
		related.R.SourceI18ns = append(related.R.SourceI18ns, o)
	}

	return nil
}

// SourceI18nsG retrieves all records.
func SourceI18nsG(mods ...qm.QueryMod) sourceI18nQuery {
	return SourceI18ns(boil.GetDB(), mods...)
}

// SourceI18ns retrieves all the records using an executor.
func SourceI18ns(exec boil.Executor, mods ...qm.QueryMod) sourceI18nQuery {
	mods = append(mods, qm.From("\"source_i18n\""))
	return sourceI18nQuery{NewQuery(exec, mods...)}
}

// FindSourceI18nG retrieves a single record by ID.
func FindSourceI18nG(sourceID int64, language string, selectCols ...string) (*SourceI18n, error) {
	return FindSourceI18n(boil.GetDB(), sourceID, language, selectCols...)
}

// FindSourceI18nGP retrieves a single record by ID, and panics on error.
func FindSourceI18nGP(sourceID int64, language string, selectCols ...string) *SourceI18n {
	retobj, err := FindSourceI18n(boil.GetDB(), sourceID, language, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindSourceI18n retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSourceI18n(exec boil.Executor, sourceID int64, language string, selectCols ...string) (*SourceI18n, error) {
	sourceI18nObj := &SourceI18n{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"source_i18n\" where \"source_id\"=$1 AND \"language\"=$2", sel,
	)

	q := queries.Raw(exec, query, sourceID, language)

	err := q.Bind(sourceI18nObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from source_i18n")
	}

	return sourceI18nObj, nil
}

// FindSourceI18nP retrieves a single record by ID with an executor, and panics on error.
func FindSourceI18nP(exec boil.Executor, sourceID int64, language string, selectCols ...string) *SourceI18n {
	retobj, err := FindSourceI18n(exec, sourceID, language, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *SourceI18n) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *SourceI18n) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *SourceI18n) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *SourceI18n) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no source_i18n provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(sourceI18nColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	sourceI18nInsertCacheMut.RLock()
	cache, cached := sourceI18nInsertCache[key]
	sourceI18nInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			sourceI18nColumns,
			sourceI18nColumnsWithDefault,
			sourceI18nColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(sourceI18nType, sourceI18nMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(sourceI18nType, sourceI18nMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"source_i18n\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into source_i18n")
	}

	if !cached {
		sourceI18nInsertCacheMut.Lock()
		sourceI18nInsertCache[key] = cache
		sourceI18nInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single SourceI18n record. See Update for
// whitelist behavior description.
func (o *SourceI18n) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single SourceI18n record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *SourceI18n) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the SourceI18n, and panics on error.
// See Update for whitelist behavior description.
func (o *SourceI18n) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the SourceI18n.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *SourceI18n) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	sourceI18nUpdateCacheMut.RLock()
	cache, cached := sourceI18nUpdateCache[key]
	sourceI18nUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(sourceI18nColumns, sourceI18nPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update source_i18n, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"source_i18n\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, sourceI18nPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(sourceI18nType, sourceI18nMapping, append(wl, sourceI18nPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update source_i18n row")
	}

	if !cached {
		sourceI18nUpdateCacheMut.Lock()
		sourceI18nUpdateCache[key] = cache
		sourceI18nUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q sourceI18nQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q sourceI18nQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for source_i18n")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o SourceI18nSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o SourceI18nSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o SourceI18nSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SourceI18nSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sourceI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"source_i18n\" SET %s WHERE (\"source_id\",\"language\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(sourceI18nPrimaryKeyColumns), len(colNames)+1, len(sourceI18nPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in sourceI18n slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *SourceI18n) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *SourceI18n) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *SourceI18n) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *SourceI18n) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no source_i18n provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(sourceI18nColumnsWithDefault, o)

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

	sourceI18nUpsertCacheMut.RLock()
	cache, cached := sourceI18nUpsertCache[key]
	sourceI18nUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			sourceI18nColumns,
			sourceI18nColumnsWithDefault,
			sourceI18nColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			sourceI18nColumns,
			sourceI18nPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert source_i18n, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(sourceI18nPrimaryKeyColumns))
			copy(conflict, sourceI18nPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"source_i18n\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(sourceI18nType, sourceI18nMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(sourceI18nType, sourceI18nMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert source_i18n")
	}

	if !cached {
		sourceI18nUpsertCacheMut.Lock()
		sourceI18nUpsertCache[key] = cache
		sourceI18nUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single SourceI18n record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *SourceI18n) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single SourceI18n record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *SourceI18n) DeleteG() error {
	if o == nil {
		return errors.New("models: no SourceI18n provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single SourceI18n record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *SourceI18n) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single SourceI18n record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *SourceI18n) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no SourceI18n provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), sourceI18nPrimaryKeyMapping)
	sql := "DELETE FROM \"source_i18n\" WHERE \"source_id\"=$1 AND \"language\"=$2"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from source_i18n")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q sourceI18nQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q sourceI18nQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no sourceI18nQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from source_i18n")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o SourceI18nSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o SourceI18nSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no SourceI18n slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o SourceI18nSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SourceI18nSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no SourceI18n slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sourceI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"source_i18n\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, sourceI18nPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(sourceI18nPrimaryKeyColumns), 1, len(sourceI18nPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from sourceI18n slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *SourceI18n) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *SourceI18n) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *SourceI18n) ReloadG() error {
	if o == nil {
		return errors.New("models: no SourceI18n provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *SourceI18n) Reload(exec boil.Executor) error {
	ret, err := FindSourceI18n(exec, o.SourceID, o.Language)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *SourceI18nSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *SourceI18nSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SourceI18nSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty SourceI18nSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SourceI18nSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	sourceI18ns := SourceI18nSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sourceI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"source_i18n\".* FROM \"source_i18n\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, sourceI18nPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(sourceI18nPrimaryKeyColumns), 1, len(sourceI18nPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&sourceI18ns)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in SourceI18nSlice")
	}

	*o = sourceI18ns

	return nil
}

// SourceI18nExists checks if the SourceI18n row exists.
func SourceI18nExists(exec boil.Executor, sourceID int64, language string) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"source_i18n\" where \"source_id\"=$1 AND \"language\"=$2 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, sourceID, language)
	}

	row := exec.QueryRow(sql, sourceID, language)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if source_i18n exists")
	}

	return exists, nil
}

// SourceI18nExistsG checks if the SourceI18n row exists.
func SourceI18nExistsG(sourceID int64, language string) (bool, error) {
	return SourceI18nExists(boil.GetDB(), sourceID, language)
}

// SourceI18nExistsGP checks if the SourceI18n row exists. Panics on error.
func SourceI18nExistsGP(sourceID int64, language string) bool {
	e, err := SourceI18nExists(boil.GetDB(), sourceID, language)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// SourceI18nExistsP checks if the SourceI18n row exists. Panics on error.
func SourceI18nExistsP(exec boil.Executor, sourceID int64, language string) bool {
	e, err := SourceI18nExists(exec, sourceID, language)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
