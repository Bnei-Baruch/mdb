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

// PersonI18n is an object representing the database table.
type PersonI18n struct {
	PersonID         int64       `boil:"person_id" json:"person_id" toml:"person_id" yaml:"person_id"`
	Language         string      `boil:"language" json:"language" toml:"language" yaml:"language"`
	OriginalLanguage null.String `boil:"original_language" json:"original_language,omitempty" toml:"original_language" yaml:"original_language,omitempty"`
	Name             null.String `boil:"name" json:"name,omitempty" toml:"name" yaml:"name,omitempty"`
	Description      null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`
	UserID           null.Int64  `boil:"user_id" json:"user_id,omitempty" toml:"user_id" yaml:"user_id,omitempty"`
	CreatedAt        time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *personI18nR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L personI18nL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// personI18nR is where relationships are stored.
type personI18nR struct {
	Person *Person
	User   *User
}

// personI18nL is where Load methods for each relationship are stored.
type personI18nL struct{}

var (
	personI18nColumns               = []string{"person_id", "language", "original_language", "name", "description", "user_id", "created_at"}
	personI18nColumnsWithoutDefault = []string{"person_id", "language", "original_language", "name", "description", "user_id"}
	personI18nColumnsWithDefault    = []string{"created_at"}
	personI18nPrimaryKeyColumns     = []string{"person_id", "language"}
)

type (
	// PersonI18nSlice is an alias for a slice of pointers to PersonI18n.
	// This should generally be used opposed to []PersonI18n.
	PersonI18nSlice []*PersonI18n

	personI18nQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	personI18nType                 = reflect.TypeOf(&PersonI18n{})
	personI18nMapping              = queries.MakeStructMapping(personI18nType)
	personI18nPrimaryKeyMapping, _ = queries.BindMapping(personI18nType, personI18nMapping, personI18nPrimaryKeyColumns)
	personI18nInsertCacheMut       sync.RWMutex
	personI18nInsertCache          = make(map[string]insertCache)
	personI18nUpdateCacheMut       sync.RWMutex
	personI18nUpdateCache          = make(map[string]updateCache)
	personI18nUpsertCacheMut       sync.RWMutex
	personI18nUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single personI18n record from the query, and panics on error.
func (q personI18nQuery) OneP() *PersonI18n {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single personI18n record from the query.
func (q personI18nQuery) One() (*PersonI18n, error) {
	o := &PersonI18n{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for person_i18n")
	}

	return o, nil
}

// AllP returns all PersonI18n records from the query, and panics on error.
func (q personI18nQuery) AllP() PersonI18nSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all PersonI18n records from the query.
func (q personI18nQuery) All() (PersonI18nSlice, error) {
	var o PersonI18nSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to PersonI18n slice")
	}

	return o, nil
}

// CountP returns the count of all PersonI18n records in the query, and panics on error.
func (q personI18nQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all PersonI18n records in the query.
func (q personI18nQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count person_i18n rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q personI18nQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q personI18nQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if person_i18n exists")
	}

	return count > 0, nil
}

// PersonG pointed to by the foreign key.
func (o *PersonI18n) PersonG(mods ...qm.QueryMod) personQuery {
	return o.Person(boil.GetDB(), mods...)
}

// Person pointed to by the foreign key.
func (o *PersonI18n) Person(exec boil.Executor, mods ...qm.QueryMod) personQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.PersonID),
	}

	queryMods = append(queryMods, mods...)

	query := Persons(exec, queryMods...)
	queries.SetFrom(query.Query, "\"persons\"")

	return query
}

// UserG pointed to by the foreign key.
func (o *PersonI18n) UserG(mods ...qm.QueryMod) userQuery {
	return o.User(boil.GetDB(), mods...)
}

// User pointed to by the foreign key.
func (o *PersonI18n) User(exec boil.Executor, mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(exec, queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// LoadPerson allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (personI18nL) LoadPerson(e boil.Executor, singular bool, maybePersonI18n interface{}) error {
	var slice []*PersonI18n
	var object *PersonI18n

	count := 1
	if singular {
		object = maybePersonI18n.(*PersonI18n)
	} else {
		slice = *maybePersonI18n.(*PersonI18nSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &personI18nR{}
		}
		args[0] = object.PersonID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &personI18nR{}
			}
			args[i] = obj.PersonID
		}
	}

	query := fmt.Sprintf(
		"select * from \"persons\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Person")
	}
	defer results.Close()

	var resultSlice []*Person
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Person")
	}

	if singular && len(resultSlice) != 0 {
		object.R.Person = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.PersonID == foreign.ID {
				local.R.Person = foreign
				break
			}
		}
	}

	return nil
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (personI18nL) LoadUser(e boil.Executor, singular bool, maybePersonI18n interface{}) error {
	var slice []*PersonI18n
	var object *PersonI18n

	count := 1
	if singular {
		object = maybePersonI18n.(*PersonI18n)
	} else {
		slice = *maybePersonI18n.(*PersonI18nSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &personI18nR{}
		}
		args[0] = object.UserID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &personI18nR{}
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

// SetPersonG of the person_i18n to the related item.
// Sets o.R.Person to related.
// Adds o to related.R.PersonI18ns.
// Uses the global database handle.
func (o *PersonI18n) SetPersonG(insert bool, related *Person) error {
	return o.SetPerson(boil.GetDB(), insert, related)
}

// SetPersonP of the person_i18n to the related item.
// Sets o.R.Person to related.
// Adds o to related.R.PersonI18ns.
// Panics on error.
func (o *PersonI18n) SetPersonP(exec boil.Executor, insert bool, related *Person) {
	if err := o.SetPerson(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetPersonGP of the person_i18n to the related item.
// Sets o.R.Person to related.
// Adds o to related.R.PersonI18ns.
// Uses the global database handle and panics on error.
func (o *PersonI18n) SetPersonGP(insert bool, related *Person) {
	if err := o.SetPerson(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetPerson of the person_i18n to the related item.
// Sets o.R.Person to related.
// Adds o to related.R.PersonI18ns.
func (o *PersonI18n) SetPerson(exec boil.Executor, insert bool, related *Person) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"person_i18n\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"person_id"}),
		strmangle.WhereClause("\"", "\"", 2, personI18nPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.PersonID, o.Language}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.PersonID = related.ID

	if o.R == nil {
		o.R = &personI18nR{
			Person: related,
		}
	} else {
		o.R.Person = related
	}

	if related.R == nil {
		related.R = &personR{
			PersonI18ns: PersonI18nSlice{o},
		}
	} else {
		related.R.PersonI18ns = append(related.R.PersonI18ns, o)
	}

	return nil
}

// SetUserG of the person_i18n to the related item.
// Sets o.R.User to related.
// Adds o to related.R.PersonI18ns.
// Uses the global database handle.
func (o *PersonI18n) SetUserG(insert bool, related *User) error {
	return o.SetUser(boil.GetDB(), insert, related)
}

// SetUserP of the person_i18n to the related item.
// Sets o.R.User to related.
// Adds o to related.R.PersonI18ns.
// Panics on error.
func (o *PersonI18n) SetUserP(exec boil.Executor, insert bool, related *User) {
	if err := o.SetUser(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetUserGP of the person_i18n to the related item.
// Sets o.R.User to related.
// Adds o to related.R.PersonI18ns.
// Uses the global database handle and panics on error.
func (o *PersonI18n) SetUserGP(insert bool, related *User) {
	if err := o.SetUser(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetUser of the person_i18n to the related item.
// Sets o.R.User to related.
// Adds o to related.R.PersonI18ns.
func (o *PersonI18n) SetUser(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"person_i18n\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, personI18nPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.PersonID, o.Language}

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
		o.R = &personI18nR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			PersonI18ns: PersonI18nSlice{o},
		}
	} else {
		related.R.PersonI18ns = append(related.R.PersonI18ns, o)
	}

	return nil
}

// RemoveUserG relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle.
func (o *PersonI18n) RemoveUserG(related *User) error {
	return o.RemoveUser(boil.GetDB(), related)
}

// RemoveUserP relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Panics on error.
func (o *PersonI18n) RemoveUserP(exec boil.Executor, related *User) {
	if err := o.RemoveUser(exec, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveUserGP relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle and panics on error.
func (o *PersonI18n) RemoveUserGP(related *User) {
	if err := o.RemoveUser(boil.GetDB(), related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveUser relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *PersonI18n) RemoveUser(exec boil.Executor, related *User) error {
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

	for i, ri := range related.R.PersonI18ns {
		if o.UserID.Int64 != ri.UserID.Int64 {
			continue
		}

		ln := len(related.R.PersonI18ns)
		if ln > 1 && i < ln-1 {
			related.R.PersonI18ns[i] = related.R.PersonI18ns[ln-1]
		}
		related.R.PersonI18ns = related.R.PersonI18ns[:ln-1]
		break
	}
	return nil
}

// PersonI18nsG retrieves all records.
func PersonI18nsG(mods ...qm.QueryMod) personI18nQuery {
	return PersonI18ns(boil.GetDB(), mods...)
}

// PersonI18ns retrieves all the records using an executor.
func PersonI18ns(exec boil.Executor, mods ...qm.QueryMod) personI18nQuery {
	mods = append(mods, qm.From("\"person_i18n\""))
	return personI18nQuery{NewQuery(exec, mods...)}
}

// FindPersonI18nG retrieves a single record by ID.
func FindPersonI18nG(personID int64, language string, selectCols ...string) (*PersonI18n, error) {
	return FindPersonI18n(boil.GetDB(), personID, language, selectCols...)
}

// FindPersonI18nGP retrieves a single record by ID, and panics on error.
func FindPersonI18nGP(personID int64, language string, selectCols ...string) *PersonI18n {
	retobj, err := FindPersonI18n(boil.GetDB(), personID, language, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindPersonI18n retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPersonI18n(exec boil.Executor, personID int64, language string, selectCols ...string) (*PersonI18n, error) {
	personI18nObj := &PersonI18n{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"person_i18n\" where \"person_id\"=$1 AND \"language\"=$2", sel,
	)

	q := queries.Raw(exec, query, personID, language)

	err := q.Bind(personI18nObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from person_i18n")
	}

	return personI18nObj, nil
}

// FindPersonI18nP retrieves a single record by ID with an executor, and panics on error.
func FindPersonI18nP(exec boil.Executor, personID int64, language string, selectCols ...string) *PersonI18n {
	retobj, err := FindPersonI18n(exec, personID, language, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *PersonI18n) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *PersonI18n) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *PersonI18n) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *PersonI18n) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no person_i18n provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(personI18nColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	personI18nInsertCacheMut.RLock()
	cache, cached := personI18nInsertCache[key]
	personI18nInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			personI18nColumns,
			personI18nColumnsWithDefault,
			personI18nColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(personI18nType, personI18nMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(personI18nType, personI18nMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"person_i18n\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into person_i18n")
	}

	if !cached {
		personI18nInsertCacheMut.Lock()
		personI18nInsertCache[key] = cache
		personI18nInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single PersonI18n record. See Update for
// whitelist behavior description.
func (o *PersonI18n) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single PersonI18n record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *PersonI18n) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the PersonI18n, and panics on error.
// See Update for whitelist behavior description.
func (o *PersonI18n) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the PersonI18n.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *PersonI18n) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	personI18nUpdateCacheMut.RLock()
	cache, cached := personI18nUpdateCache[key]
	personI18nUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(personI18nColumns, personI18nPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update person_i18n, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"person_i18n\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, personI18nPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(personI18nType, personI18nMapping, append(wl, personI18nPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update person_i18n row")
	}

	if !cached {
		personI18nUpdateCacheMut.Lock()
		personI18nUpdateCache[key] = cache
		personI18nUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q personI18nQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q personI18nQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for person_i18n")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o PersonI18nSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o PersonI18nSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o PersonI18nSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PersonI18nSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), personI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"person_i18n\" SET %s WHERE (\"person_id\",\"language\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(personI18nPrimaryKeyColumns), len(colNames)+1, len(personI18nPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in personI18n slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *PersonI18n) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *PersonI18n) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *PersonI18n) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *PersonI18n) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no person_i18n provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(personI18nColumnsWithDefault, o)

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

	personI18nUpsertCacheMut.RLock()
	cache, cached := personI18nUpsertCache[key]
	personI18nUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			personI18nColumns,
			personI18nColumnsWithDefault,
			personI18nColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			personI18nColumns,
			personI18nPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert person_i18n, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(personI18nPrimaryKeyColumns))
			copy(conflict, personI18nPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"person_i18n\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(personI18nType, personI18nMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(personI18nType, personI18nMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert person_i18n")
	}

	if !cached {
		personI18nUpsertCacheMut.Lock()
		personI18nUpsertCache[key] = cache
		personI18nUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single PersonI18n record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *PersonI18n) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single PersonI18n record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *PersonI18n) DeleteG() error {
	if o == nil {
		return errors.New("models: no PersonI18n provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single PersonI18n record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *PersonI18n) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single PersonI18n record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *PersonI18n) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no PersonI18n provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), personI18nPrimaryKeyMapping)
	sql := "DELETE FROM \"person_i18n\" WHERE \"person_id\"=$1 AND \"language\"=$2"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from person_i18n")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q personI18nQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q personI18nQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no personI18nQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from person_i18n")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o PersonI18nSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o PersonI18nSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no PersonI18n slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o PersonI18nSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PersonI18nSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no PersonI18n slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), personI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"person_i18n\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, personI18nPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(personI18nPrimaryKeyColumns), 1, len(personI18nPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from personI18n slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *PersonI18n) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *PersonI18n) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *PersonI18n) ReloadG() error {
	if o == nil {
		return errors.New("models: no PersonI18n provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *PersonI18n) Reload(exec boil.Executor) error {
	ret, err := FindPersonI18n(exec, o.PersonID, o.Language)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *PersonI18nSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *PersonI18nSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PersonI18nSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty PersonI18nSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PersonI18nSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	personI18ns := PersonI18nSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), personI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"person_i18n\".* FROM \"person_i18n\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, personI18nPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(personI18nPrimaryKeyColumns), 1, len(personI18nPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&personI18ns)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in PersonI18nSlice")
	}

	*o = personI18ns

	return nil
}

// PersonI18nExists checks if the PersonI18n row exists.
func PersonI18nExists(exec boil.Executor, personID int64, language string) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"person_i18n\" where \"person_id\"=$1 AND \"language\"=$2 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, personID, language)
	}

	row := exec.QueryRow(sql, personID, language)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if person_i18n exists")
	}

	return exists, nil
}

// PersonI18nExistsG checks if the PersonI18n row exists.
func PersonI18nExistsG(personID int64, language string) (bool, error) {
	return PersonI18nExists(boil.GetDB(), personID, language)
}

// PersonI18nExistsGP checks if the PersonI18n row exists. Panics on error.
func PersonI18nExistsGP(personID int64, language string) bool {
	e, err := PersonI18nExists(boil.GetDB(), personID, language)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// PersonI18nExistsP checks if the PersonI18n row exists. Panics on error.
func PersonI18nExistsP(exec boil.Executor, personID int64, language string) bool {
	e, err := PersonI18nExists(exec, personID, language)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
