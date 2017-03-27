package kmodels

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

// ContainerTranscript is an object representing the database table.
type ContainerTranscript struct {
	ID          int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	ContainerID null.Int    `boil:"container_id" json:"container_id,omitempty" toml:"container_id" yaml:"container_id,omitempty"`
	Toc         null.String `boil:"toc" json:"toc,omitempty" toml:"toc" yaml:"toc,omitempty"`
	Transcript  null.String `boil:"transcript" json:"transcript,omitempty" toml:"transcript" yaml:"transcript,omitempty"`
	Lang        null.String `boil:"lang" json:"lang,omitempty" toml:"lang" yaml:"lang,omitempty"`
	CreatedAt   time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time   `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *containerTranscriptR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L containerTranscriptL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// containerTranscriptR is where relationships are stored.
type containerTranscriptR struct {
}

// containerTranscriptL is where Load methods for each relationship are stored.
type containerTranscriptL struct{}

var (
	containerTranscriptColumns               = []string{"id", "container_id", "toc", "transcript", "lang", "created_at", "updated_at"}
	containerTranscriptColumnsWithoutDefault = []string{"container_id", "toc", "transcript", "lang", "created_at", "updated_at"}
	containerTranscriptColumnsWithDefault    = []string{"id"}
	containerTranscriptPrimaryKeyColumns     = []string{"id"}
)

type (
	// ContainerTranscriptSlice is an alias for a slice of pointers to ContainerTranscript.
	// This should generally be used opposed to []ContainerTranscript.
	ContainerTranscriptSlice []*ContainerTranscript

	containerTranscriptQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	containerTranscriptType                 = reflect.TypeOf(&ContainerTranscript{})
	containerTranscriptMapping              = queries.MakeStructMapping(containerTranscriptType)
	containerTranscriptPrimaryKeyMapping, _ = queries.BindMapping(containerTranscriptType, containerTranscriptMapping, containerTranscriptPrimaryKeyColumns)
	containerTranscriptInsertCacheMut       sync.RWMutex
	containerTranscriptInsertCache          = make(map[string]insertCache)
	containerTranscriptUpdateCacheMut       sync.RWMutex
	containerTranscriptUpdateCache          = make(map[string]updateCache)
	containerTranscriptUpsertCacheMut       sync.RWMutex
	containerTranscriptUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single containerTranscript record from the query, and panics on error.
func (q containerTranscriptQuery) OneP() *ContainerTranscript {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single containerTranscript record from the query.
func (q containerTranscriptQuery) One() (*ContainerTranscript, error) {
	o := &ContainerTranscript{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "kmodels: failed to execute a one query for container_transcripts")
	}

	return o, nil
}

// AllP returns all ContainerTranscript records from the query, and panics on error.
func (q containerTranscriptQuery) AllP() ContainerTranscriptSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all ContainerTranscript records from the query.
func (q containerTranscriptQuery) All() (ContainerTranscriptSlice, error) {
	var o ContainerTranscriptSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "kmodels: failed to assign all query results to ContainerTranscript slice")
	}

	return o, nil
}

// CountP returns the count of all ContainerTranscript records in the query, and panics on error.
func (q containerTranscriptQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all ContainerTranscript records in the query.
func (q containerTranscriptQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "kmodels: failed to count container_transcripts rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q containerTranscriptQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q containerTranscriptQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "kmodels: failed to check if container_transcripts exists")
	}

	return count > 0, nil
}

// ContainerTranscriptsG retrieves all records.
func ContainerTranscriptsG(mods ...qm.QueryMod) containerTranscriptQuery {
	return ContainerTranscripts(boil.GetDB(), mods...)
}

// ContainerTranscripts retrieves all the records using an executor.
func ContainerTranscripts(exec boil.Executor, mods ...qm.QueryMod) containerTranscriptQuery {
	mods = append(mods, qm.From("\"container_transcripts\""))
	return containerTranscriptQuery{NewQuery(exec, mods...)}
}

// FindContainerTranscriptG retrieves a single record by ID.
func FindContainerTranscriptG(id int, selectCols ...string) (*ContainerTranscript, error) {
	return FindContainerTranscript(boil.GetDB(), id, selectCols...)
}

// FindContainerTranscriptGP retrieves a single record by ID, and panics on error.
func FindContainerTranscriptGP(id int, selectCols ...string) *ContainerTranscript {
	retobj, err := FindContainerTranscript(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindContainerTranscript retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindContainerTranscript(exec boil.Executor, id int, selectCols ...string) (*ContainerTranscript, error) {
	containerTranscriptObj := &ContainerTranscript{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"container_transcripts\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(containerTranscriptObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "kmodels: unable to select from container_transcripts")
	}

	return containerTranscriptObj, nil
}

// FindContainerTranscriptP retrieves a single record by ID with an executor, and panics on error.
func FindContainerTranscriptP(exec boil.Executor, id int, selectCols ...string) *ContainerTranscript {
	retobj, err := FindContainerTranscript(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *ContainerTranscript) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *ContainerTranscript) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *ContainerTranscript) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *ContainerTranscript) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("kmodels: no container_transcripts provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}
	if o.UpdatedAt.IsZero() {
		o.UpdatedAt = currTime
	}

	nzDefaults := queries.NonZeroDefaultSet(containerTranscriptColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	containerTranscriptInsertCacheMut.RLock()
	cache, cached := containerTranscriptInsertCache[key]
	containerTranscriptInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			containerTranscriptColumns,
			containerTranscriptColumnsWithDefault,
			containerTranscriptColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(containerTranscriptType, containerTranscriptMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(containerTranscriptType, containerTranscriptMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"container_transcripts\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "kmodels: unable to insert into container_transcripts")
	}

	if !cached {
		containerTranscriptInsertCacheMut.Lock()
		containerTranscriptInsertCache[key] = cache
		containerTranscriptInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single ContainerTranscript record. See Update for
// whitelist behavior description.
func (o *ContainerTranscript) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single ContainerTranscript record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *ContainerTranscript) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the ContainerTranscript, and panics on error.
// See Update for whitelist behavior description.
func (o *ContainerTranscript) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the ContainerTranscript.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *ContainerTranscript) Update(exec boil.Executor, whitelist ...string) error {
	currTime := time.Now().In(boil.GetLocation())

	o.UpdatedAt = currTime

	var err error
	key := makeCacheKey(whitelist, nil)
	containerTranscriptUpdateCacheMut.RLock()
	cache, cached := containerTranscriptUpdateCache[key]
	containerTranscriptUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(containerTranscriptColumns, containerTranscriptPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("kmodels: unable to update container_transcripts, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"container_transcripts\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, containerTranscriptPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(containerTranscriptType, containerTranscriptMapping, append(wl, containerTranscriptPrimaryKeyColumns...))
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
		return errors.Wrap(err, "kmodels: unable to update container_transcripts row")
	}

	if !cached {
		containerTranscriptUpdateCacheMut.Lock()
		containerTranscriptUpdateCache[key] = cache
		containerTranscriptUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q containerTranscriptQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q containerTranscriptQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to update all for container_transcripts")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o ContainerTranscriptSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o ContainerTranscriptSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o ContainerTranscriptSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ContainerTranscriptSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("kmodels: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), containerTranscriptPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"container_transcripts\" SET %s WHERE (\"id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(containerTranscriptPrimaryKeyColumns), len(colNames)+1, len(containerTranscriptPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to update all in containerTranscript slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *ContainerTranscript) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *ContainerTranscript) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *ContainerTranscript) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *ContainerTranscript) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("kmodels: no container_transcripts provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}
	o.UpdatedAt = currTime

	nzDefaults := queries.NonZeroDefaultSet(containerTranscriptColumnsWithDefault, o)

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

	containerTranscriptUpsertCacheMut.RLock()
	cache, cached := containerTranscriptUpsertCache[key]
	containerTranscriptUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			containerTranscriptColumns,
			containerTranscriptColumnsWithDefault,
			containerTranscriptColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			containerTranscriptColumns,
			containerTranscriptPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("kmodels: unable to upsert container_transcripts, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(containerTranscriptPrimaryKeyColumns))
			copy(conflict, containerTranscriptPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"container_transcripts\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(containerTranscriptType, containerTranscriptMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(containerTranscriptType, containerTranscriptMapping, ret)
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
		return errors.Wrap(err, "kmodels: unable to upsert container_transcripts")
	}

	if !cached {
		containerTranscriptUpsertCacheMut.Lock()
		containerTranscriptUpsertCache[key] = cache
		containerTranscriptUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single ContainerTranscript record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *ContainerTranscript) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single ContainerTranscript record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *ContainerTranscript) DeleteG() error {
	if o == nil {
		return errors.New("kmodels: no ContainerTranscript provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single ContainerTranscript record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *ContainerTranscript) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single ContainerTranscript record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ContainerTranscript) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("kmodels: no ContainerTranscript provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), containerTranscriptPrimaryKeyMapping)
	sql := "DELETE FROM \"container_transcripts\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to delete from container_transcripts")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q containerTranscriptQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q containerTranscriptQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("kmodels: no containerTranscriptQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to delete all from container_transcripts")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o ContainerTranscriptSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o ContainerTranscriptSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("kmodels: no ContainerTranscript slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o ContainerTranscriptSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ContainerTranscriptSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("kmodels: no ContainerTranscript slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), containerTranscriptPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"container_transcripts\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, containerTranscriptPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(containerTranscriptPrimaryKeyColumns), 1, len(containerTranscriptPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to delete all from containerTranscript slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *ContainerTranscript) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *ContainerTranscript) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *ContainerTranscript) ReloadG() error {
	if o == nil {
		return errors.New("kmodels: no ContainerTranscript provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ContainerTranscript) Reload(exec boil.Executor) error {
	ret, err := FindContainerTranscript(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ContainerTranscriptSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ContainerTranscriptSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ContainerTranscriptSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("kmodels: empty ContainerTranscriptSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ContainerTranscriptSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	containerTranscripts := ContainerTranscriptSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), containerTranscriptPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"container_transcripts\".* FROM \"container_transcripts\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, containerTranscriptPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(containerTranscriptPrimaryKeyColumns), 1, len(containerTranscriptPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&containerTranscripts)
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to reload all in ContainerTranscriptSlice")
	}

	*o = containerTranscripts

	return nil
}

// ContainerTranscriptExists checks if the ContainerTranscript row exists.
func ContainerTranscriptExists(exec boil.Executor, id int) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"container_transcripts\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "kmodels: unable to check if container_transcripts exists")
	}

	return exists, nil
}

// ContainerTranscriptExistsG checks if the ContainerTranscript row exists.
func ContainerTranscriptExistsG(id int) (bool, error) {
	return ContainerTranscriptExists(boil.GetDB(), id)
}

// ContainerTranscriptExistsGP checks if the ContainerTranscript row exists. Panics on error.
func ContainerTranscriptExistsGP(id int) bool {
	e, err := ContainerTranscriptExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// ContainerTranscriptExistsP checks if the ContainerTranscript row exists. Panics on error.
func ContainerTranscriptExistsP(exec boil.Executor, id int) bool {
	e, err := ContainerTranscriptExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}