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

// CollectionI18n is an object representing the database table.
type CollectionI18n struct {
	CollectionID     int64       `boil:"collection_id" json:"collection_id" toml:"collection_id" yaml:"collection_id"`
	Language         string      `boil:"language" json:"language" toml:"language" yaml:"language"`
	OriginalLanguage null.String `boil:"original_language" json:"original_language,omitempty" toml:"original_language" yaml:"original_language,omitempty"`
	Name             null.String `boil:"name" json:"name,omitempty" toml:"name" yaml:"name,omitempty"`
	Description      null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`
	UserID           null.Int64  `boil:"user_id" json:"user_id,omitempty" toml:"user_id" yaml:"user_id,omitempty"`
	CreatedAt        time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *collectionI18nR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L collectionI18nL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// collectionI18nR is where relationships are stored.
type collectionI18nR struct {
	Collection *Collection
	User       *User
}

// collectionI18nL is where Load methods for each relationship are stored.
type collectionI18nL struct{}

var (
	collectionI18nColumns               = []string{"collection_id", "language", "original_language", "name", "description", "user_id", "created_at"}
	collectionI18nColumnsWithoutDefault = []string{"collection_id", "language", "original_language", "name", "description", "user_id"}
	collectionI18nColumnsWithDefault    = []string{"created_at"}
	collectionI18nPrimaryKeyColumns     = []string{"collection_id", "language"}
)

type (
	// CollectionI18nSlice is an alias for a slice of pointers to CollectionI18n.
	// This should generally be used opposed to []CollectionI18n.
	CollectionI18nSlice []*CollectionI18n

	collectionI18nQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	collectionI18nType                 = reflect.TypeOf(&CollectionI18n{})
	collectionI18nMapping              = queries.MakeStructMapping(collectionI18nType)
	collectionI18nPrimaryKeyMapping, _ = queries.BindMapping(collectionI18nType, collectionI18nMapping, collectionI18nPrimaryKeyColumns)
	collectionI18nInsertCacheMut       sync.RWMutex
	collectionI18nInsertCache          = make(map[string]insertCache)
	collectionI18nUpdateCacheMut       sync.RWMutex
	collectionI18nUpdateCache          = make(map[string]updateCache)
	collectionI18nUpsertCacheMut       sync.RWMutex
	collectionI18nUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single collectionI18n record from the query, and panics on error.
func (q collectionI18nQuery) OneP() *CollectionI18n {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single collectionI18n record from the query.
func (q collectionI18nQuery) One() (*CollectionI18n, error) {
	o := &CollectionI18n{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for collection_i18n")
	}

	return o, nil
}

// AllP returns all CollectionI18n records from the query, and panics on error.
func (q collectionI18nQuery) AllP() CollectionI18nSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all CollectionI18n records from the query.
func (q collectionI18nQuery) All() (CollectionI18nSlice, error) {
	var o CollectionI18nSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to CollectionI18n slice")
	}

	return o, nil
}

// CountP returns the count of all CollectionI18n records in the query, and panics on error.
func (q collectionI18nQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all CollectionI18n records in the query.
func (q collectionI18nQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count collection_i18n rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q collectionI18nQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q collectionI18nQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if collection_i18n exists")
	}

	return count > 0, nil
}

// CollectionG pointed to by the foreign key.
func (o *CollectionI18n) CollectionG(mods ...qm.QueryMod) collectionQuery {
	return o.Collection(boil.GetDB(), mods...)
}

// Collection pointed to by the foreign key.
func (o *CollectionI18n) Collection(exec boil.Executor, mods ...qm.QueryMod) collectionQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.CollectionID),
	}

	queryMods = append(queryMods, mods...)

	query := Collections(exec, queryMods...)
	queries.SetFrom(query.Query, "\"collections\"")

	return query
}

// UserG pointed to by the foreign key.
func (o *CollectionI18n) UserG(mods ...qm.QueryMod) userQuery {
	return o.User(boil.GetDB(), mods...)
}

// User pointed to by the foreign key.
func (o *CollectionI18n) User(exec boil.Executor, mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(exec, queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// LoadCollection allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (collectionI18nL) LoadCollection(e boil.Executor, singular bool, maybeCollectionI18n interface{}) error {
	var slice []*CollectionI18n
	var object *CollectionI18n

	count := 1
	if singular {
		object = maybeCollectionI18n.(*CollectionI18n)
	} else {
		slice = *maybeCollectionI18n.(*CollectionI18nSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &collectionI18nR{}
		}
		args[0] = object.CollectionID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &collectionI18nR{}
			}
			args[i] = obj.CollectionID
		}
	}

	query := fmt.Sprintf(
		"select * from \"collections\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Collection")
	}
	defer results.Close()

	var resultSlice []*Collection
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Collection")
	}

	if singular && len(resultSlice) != 0 {
		object.R.Collection = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.CollectionID == foreign.ID {
				local.R.Collection = foreign
				break
			}
		}
	}

	return nil
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (collectionI18nL) LoadUser(e boil.Executor, singular bool, maybeCollectionI18n interface{}) error {
	var slice []*CollectionI18n
	var object *CollectionI18n

	count := 1
	if singular {
		object = maybeCollectionI18n.(*CollectionI18n)
	} else {
		slice = *maybeCollectionI18n.(*CollectionI18nSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &collectionI18nR{}
		}
		args[0] = object.UserID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &collectionI18nR{}
			}
			args[i] = obj.UserID
		}
	}

	query := fmt.Sprintf(
		"select * from \"users\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load User")
	}
	defer results.Close()

	var resultSlice []*User
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice User")
	}

	if singular && len(resultSlice) != 0 {
		object.R.User = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.UserID.Int64 == foreign.ID {
				local.R.User = foreign
				break
			}
		}
	}

	return nil
}

// SetCollectionG of the collection_i18n to the related item.
// Sets o.R.Collection to related.
// Adds o to related.R.CollectionI18ns.
// Uses the global database handle.
func (o *CollectionI18n) SetCollectionG(insert bool, related *Collection) error {
	return o.SetCollection(boil.GetDB(), insert, related)
}

// SetCollectionP of the collection_i18n to the related item.
// Sets o.R.Collection to related.
// Adds o to related.R.CollectionI18ns.
// Panics on error.
func (o *CollectionI18n) SetCollectionP(exec boil.Executor, insert bool, related *Collection) {
	if err := o.SetCollection(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetCollectionGP of the collection_i18n to the related item.
// Sets o.R.Collection to related.
// Adds o to related.R.CollectionI18ns.
// Uses the global database handle and panics on error.
func (o *CollectionI18n) SetCollectionGP(insert bool, related *Collection) {
	if err := o.SetCollection(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetCollection of the collection_i18n to the related item.
// Sets o.R.Collection to related.
// Adds o to related.R.CollectionI18ns.
func (o *CollectionI18n) SetCollection(exec boil.Executor, insert bool, related *Collection) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"collection_i18n\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"collection_id"}),
		strmangle.WhereClause("\"", "\"", 2, collectionI18nPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.CollectionID, o.Language}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.CollectionID = related.ID

	if o.R == nil {
		o.R = &collectionI18nR{
			Collection: related,
		}
	} else {
		o.R.Collection = related
	}

	if related.R == nil {
		related.R = &collectionR{
			CollectionI18ns: CollectionI18nSlice{o},
		}
	} else {
		related.R.CollectionI18ns = append(related.R.CollectionI18ns, o)
	}

	return nil
}

// SetUserG of the collection_i18n to the related item.
// Sets o.R.User to related.
// Adds o to related.R.CollectionI18ns.
// Uses the global database handle.
func (o *CollectionI18n) SetUserG(insert bool, related *User) error {
	return o.SetUser(boil.GetDB(), insert, related)
}

// SetUserP of the collection_i18n to the related item.
// Sets o.R.User to related.
// Adds o to related.R.CollectionI18ns.
// Panics on error.
func (o *CollectionI18n) SetUserP(exec boil.Executor, insert bool, related *User) {
	if err := o.SetUser(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetUserGP of the collection_i18n to the related item.
// Sets o.R.User to related.
// Adds o to related.R.CollectionI18ns.
// Uses the global database handle and panics on error.
func (o *CollectionI18n) SetUserGP(insert bool, related *User) {
	if err := o.SetUser(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetUser of the collection_i18n to the related item.
// Sets o.R.User to related.
// Adds o to related.R.CollectionI18ns.
func (o *CollectionI18n) SetUser(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"collection_i18n\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, collectionI18nPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.CollectionID, o.Language}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.UserID.Int64 = related.ID
	o.UserID.Valid = true

	if o.R == nil {
		o.R = &collectionI18nR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			CollectionI18ns: CollectionI18nSlice{o},
		}
	} else {
		related.R.CollectionI18ns = append(related.R.CollectionI18ns, o)
	}

	return nil
}

// RemoveUserG relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle.
func (o *CollectionI18n) RemoveUserG(related *User) error {
	return o.RemoveUser(boil.GetDB(), related)
}

// RemoveUserP relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Panics on error.
func (o *CollectionI18n) RemoveUserP(exec boil.Executor, related *User) {
	if err := o.RemoveUser(exec, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveUserGP relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle and panics on error.
func (o *CollectionI18n) RemoveUserGP(related *User) {
	if err := o.RemoveUser(boil.GetDB(), related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveUser relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *CollectionI18n) RemoveUser(exec boil.Executor, related *User) error {
	var err error

	o.UserID.Valid = false
	if err = o.Update(exec, "user_id"); err != nil {
		o.UserID.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.User = nil
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.CollectionI18ns {
		if o.UserID.Int64 != ri.UserID.Int64 {
			continue
		}

		ln := len(related.R.CollectionI18ns)
		if ln > 1 && i < ln-1 {
			related.R.CollectionI18ns[i] = related.R.CollectionI18ns[ln-1]
		}
		related.R.CollectionI18ns = related.R.CollectionI18ns[:ln-1]
		break
	}
	return nil
}

// CollectionI18nsG retrieves all records.
func CollectionI18nsG(mods ...qm.QueryMod) collectionI18nQuery {
	return CollectionI18ns(boil.GetDB(), mods...)
}

// CollectionI18ns retrieves all the records using an executor.
func CollectionI18ns(exec boil.Executor, mods ...qm.QueryMod) collectionI18nQuery {
	mods = append(mods, qm.From("\"collection_i18n\""))
	return collectionI18nQuery{NewQuery(exec, mods...)}
}

// FindCollectionI18nG retrieves a single record by ID.
func FindCollectionI18nG(collectionID int64, language string, selectCols ...string) (*CollectionI18n, error) {
	return FindCollectionI18n(boil.GetDB(), collectionID, language, selectCols...)
}

// FindCollectionI18nGP retrieves a single record by ID, and panics on error.
func FindCollectionI18nGP(collectionID int64, language string, selectCols ...string) *CollectionI18n {
	retobj, err := FindCollectionI18n(boil.GetDB(), collectionID, language, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindCollectionI18n retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindCollectionI18n(exec boil.Executor, collectionID int64, language string, selectCols ...string) (*CollectionI18n, error) {
	collectionI18nObj := &CollectionI18n{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"collection_i18n\" where \"collection_id\"=$1 AND \"language\"=$2", sel,
	)

	q := queries.Raw(exec, query, collectionID, language)

	err := q.Bind(collectionI18nObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from collection_i18n")
	}

	return collectionI18nObj, nil
}

// FindCollectionI18nP retrieves a single record by ID with an executor, and panics on error.
func FindCollectionI18nP(exec boil.Executor, collectionID int64, language string, selectCols ...string) *CollectionI18n {
	retobj, err := FindCollectionI18n(exec, collectionID, language, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *CollectionI18n) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *CollectionI18n) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *CollectionI18n) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *CollectionI18n) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no collection_i18n provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(collectionI18nColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	collectionI18nInsertCacheMut.RLock()
	cache, cached := collectionI18nInsertCache[key]
	collectionI18nInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			collectionI18nColumns,
			collectionI18nColumnsWithDefault,
			collectionI18nColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(collectionI18nType, collectionI18nMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(collectionI18nType, collectionI18nMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"collection_i18n\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into collection_i18n")
	}

	if !cached {
		collectionI18nInsertCacheMut.Lock()
		collectionI18nInsertCache[key] = cache
		collectionI18nInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single CollectionI18n record. See Update for
// whitelist behavior description.
func (o *CollectionI18n) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single CollectionI18n record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *CollectionI18n) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the CollectionI18n, and panics on error.
// See Update for whitelist behavior description.
func (o *CollectionI18n) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the CollectionI18n.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *CollectionI18n) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	collectionI18nUpdateCacheMut.RLock()
	cache, cached := collectionI18nUpdateCache[key]
	collectionI18nUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(collectionI18nColumns, collectionI18nPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update collection_i18n, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"collection_i18n\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, collectionI18nPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(collectionI18nType, collectionI18nMapping, append(wl, collectionI18nPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update collection_i18n row")
	}

	if !cached {
		collectionI18nUpdateCacheMut.Lock()
		collectionI18nUpdateCache[key] = cache
		collectionI18nUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q collectionI18nQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q collectionI18nQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for collection_i18n")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o CollectionI18nSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o CollectionI18nSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o CollectionI18nSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CollectionI18nSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), collectionI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"collection_i18n\" SET %s WHERE (\"collection_id\",\"language\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(collectionI18nPrimaryKeyColumns), len(colNames)+1, len(collectionI18nPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in collectionI18n slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *CollectionI18n) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *CollectionI18n) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *CollectionI18n) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *CollectionI18n) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no collection_i18n provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(collectionI18nColumnsWithDefault, o)

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

	collectionI18nUpsertCacheMut.RLock()
	cache, cached := collectionI18nUpsertCache[key]
	collectionI18nUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			collectionI18nColumns,
			collectionI18nColumnsWithDefault,
			collectionI18nColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			collectionI18nColumns,
			collectionI18nPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert collection_i18n, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(collectionI18nPrimaryKeyColumns))
			copy(conflict, collectionI18nPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"collection_i18n\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(collectionI18nType, collectionI18nMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(collectionI18nType, collectionI18nMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert collection_i18n")
	}

	if !cached {
		collectionI18nUpsertCacheMut.Lock()
		collectionI18nUpsertCache[key] = cache
		collectionI18nUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single CollectionI18n record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *CollectionI18n) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single CollectionI18n record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *CollectionI18n) DeleteG() error {
	if o == nil {
		return errors.New("models: no CollectionI18n provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single CollectionI18n record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *CollectionI18n) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single CollectionI18n record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *CollectionI18n) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no CollectionI18n provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), collectionI18nPrimaryKeyMapping)
	sql := "DELETE FROM \"collection_i18n\" WHERE \"collection_id\"=$1 AND \"language\"=$2"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from collection_i18n")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q collectionI18nQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q collectionI18nQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no collectionI18nQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from collection_i18n")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o CollectionI18nSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o CollectionI18nSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no CollectionI18n slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o CollectionI18nSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CollectionI18nSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no CollectionI18n slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), collectionI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"collection_i18n\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, collectionI18nPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(collectionI18nPrimaryKeyColumns), 1, len(collectionI18nPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from collectionI18n slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *CollectionI18n) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *CollectionI18n) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *CollectionI18n) ReloadG() error {
	if o == nil {
		return errors.New("models: no CollectionI18n provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *CollectionI18n) Reload(exec boil.Executor) error {
	ret, err := FindCollectionI18n(exec, o.CollectionID, o.Language)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *CollectionI18nSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *CollectionI18nSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CollectionI18nSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty CollectionI18nSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CollectionI18nSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	collectionI18ns := CollectionI18nSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), collectionI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"collection_i18n\".* FROM \"collection_i18n\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, collectionI18nPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(collectionI18nPrimaryKeyColumns), 1, len(collectionI18nPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&collectionI18ns)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in CollectionI18nSlice")
	}

	*o = collectionI18ns

	return nil
}

// CollectionI18nExists checks if the CollectionI18n row exists.
func CollectionI18nExists(exec boil.Executor, collectionID int64, language string) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"collection_i18n\" where \"collection_id\"=$1 AND \"language\"=$2 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, collectionID, language)
	}

	row := exec.QueryRow(sql, collectionID, language)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if collection_i18n exists")
	}

	return exists, nil
}

// CollectionI18nExistsG checks if the CollectionI18n row exists.
func CollectionI18nExistsG(collectionID int64, language string) (bool, error) {
	return CollectionI18nExists(boil.GetDB(), collectionID, language)
}

// CollectionI18nExistsGP checks if the CollectionI18n row exists. Panics on error.
func CollectionI18nExistsGP(collectionID int64, language string) bool {
	e, err := CollectionI18nExists(boil.GetDB(), collectionID, language)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// CollectionI18nExistsP checks if the CollectionI18n row exists. Panics on error.
func CollectionI18nExistsP(exec boil.Executor, collectionID int64, language string) bool {
	e, err := CollectionI18nExists(exec, collectionID, language)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
