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

// Person is an object representing the database table.
type Person struct {
	ID            int64      `boil:"id" json:"id" toml:"id" yaml:"id"`
	UID           string     `boil:"uid" json:"uid" toml:"uid" yaml:"uid"`
	NameID        int64      `boil:"name_id" json:"name_id" toml:"name_id" yaml:"name_id"`
	DescriptionID null.Int64 `boil:"description_id" json:"description_id,omitempty" toml:"description_id" yaml:"description_id,omitempty"`

	R *personR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L personL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// personR is where relationships are stored.
type personR struct {
	Name                *StringTranslation
	Description         *StringTranslation
	ContentUnitsPersons ContentUnitsPersonSlice
}

// personL is where Load methods for each relationship are stored.
type personL struct{}

var (
	personColumns               = []string{"id", "uid", "name_id", "description_id"}
	personColumnsWithoutDefault = []string{"uid", "name_id", "description_id"}
	personColumnsWithDefault    = []string{"id"}
	personPrimaryKeyColumns     = []string{"id"}
)

type (
	// PersonSlice is an alias for a slice of pointers to Person.
	// This should generally be used opposed to []Person.
	PersonSlice []*Person

	personQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	personType                 = reflect.TypeOf(&Person{})
	personMapping              = queries.MakeStructMapping(personType)
	personPrimaryKeyMapping, _ = queries.BindMapping(personType, personMapping, personPrimaryKeyColumns)
	personInsertCacheMut       sync.RWMutex
	personInsertCache          = make(map[string]insertCache)
	personUpdateCacheMut       sync.RWMutex
	personUpdateCache          = make(map[string]updateCache)
	personUpsertCacheMut       sync.RWMutex
	personUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single person record from the query, and panics on error.
func (q personQuery) OneP() *Person {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single person record from the query.
func (q personQuery) One() (*Person, error) {
	o := &Person{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "gmodels: failed to execute a one query for persons")
	}

	return o, nil
}

// AllP returns all Person records from the query, and panics on error.
func (q personQuery) AllP() PersonSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Person records from the query.
func (q personQuery) All() (PersonSlice, error) {
	var o PersonSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "gmodels: failed to assign all query results to Person slice")
	}

	return o, nil
}

// CountP returns the count of all Person records in the query, and panics on error.
func (q personQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Person records in the query.
func (q personQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "gmodels: failed to count persons rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q personQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q personQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "gmodels: failed to check if persons exists")
	}

	return count > 0, nil
}

// NameG pointed to by the foreign key.
func (o *Person) NameG(mods ...qm.QueryMod) stringTranslationQuery {
	return o.Name(boil.GetDB(), mods...)
}

// Name pointed to by the foreign key.
func (o *Person) Name(exec boil.Executor, mods ...qm.QueryMod) stringTranslationQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.NameID),
	}

	queryMods = append(queryMods, mods...)

	query := StringTranslations(exec, queryMods...)
	queries.SetFrom(query.Query, "\"string_translations\"")

	return query
}

// DescriptionG pointed to by the foreign key.
func (o *Person) DescriptionG(mods ...qm.QueryMod) stringTranslationQuery {
	return o.Description(boil.GetDB(), mods...)
}

// Description pointed to by the foreign key.
func (o *Person) Description(exec boil.Executor, mods ...qm.QueryMod) stringTranslationQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.DescriptionID),
	}

	queryMods = append(queryMods, mods...)

	query := StringTranslations(exec, queryMods...)
	queries.SetFrom(query.Query, "\"string_translations\"")

	return query
}

// ContentUnitsPersonsG retrieves all the content_units_person's content units persons.
func (o *Person) ContentUnitsPersonsG(mods ...qm.QueryMod) contentUnitsPersonQuery {
	return o.ContentUnitsPersons(boil.GetDB(), mods...)
}

// ContentUnitsPersons retrieves all the content_units_person's content units persons with an executor.
func (o *Person) ContentUnitsPersons(exec boil.Executor, mods ...qm.QueryMod) contentUnitsPersonQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"person_id\"=?", o.ID),
	)

	query := ContentUnitsPersons(exec, queryMods...)
	queries.SetFrom(query.Query, "\"content_units_persons\" as \"a\"")
	return query
}

// LoadName allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (personL) LoadName(e boil.Executor, singular bool, maybePerson interface{}) error {
	var slice []*Person
	var object *Person

	count := 1
	if singular {
		object = maybePerson.(*Person)
	} else {
		slice = *maybePerson.(*PersonSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &personR{}
		}
		args[0] = object.NameID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &personR{}
			}
			args[i] = obj.NameID
		}
	}

	query := fmt.Sprintf(
		"select * from \"string_translations\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load StringTranslation")
	}
	defer results.Close()

	var resultSlice []*StringTranslation
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice StringTranslation")
	}

	if singular && len(resultSlice) != 0 {
		object.R.Name = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.NameID == foreign.ID {
				local.R.Name = foreign
				break
			}
		}
	}

	return nil
}

// LoadDescription allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (personL) LoadDescription(e boil.Executor, singular bool, maybePerson interface{}) error {
	var slice []*Person
	var object *Person

	count := 1
	if singular {
		object = maybePerson.(*Person)
	} else {
		slice = *maybePerson.(*PersonSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &personR{}
		}
		args[0] = object.DescriptionID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &personR{}
			}
			args[i] = obj.DescriptionID
		}
	}

	query := fmt.Sprintf(
		"select * from \"string_translations\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load StringTranslation")
	}
	defer results.Close()

	var resultSlice []*StringTranslation
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice StringTranslation")
	}

	if singular && len(resultSlice) != 0 {
		object.R.Description = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.DescriptionID.Int64 == foreign.ID {
				local.R.Description = foreign
				break
			}
		}
	}

	return nil
}

// LoadContentUnitsPersons allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (personL) LoadContentUnitsPersons(e boil.Executor, singular bool, maybePerson interface{}) error {
	var slice []*Person
	var object *Person

	count := 1
	if singular {
		object = maybePerson.(*Person)
	} else {
		slice = *maybePerson.(*PersonSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &personR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &personR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"content_units_persons\" where \"person_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load content_units_persons")
	}
	defer results.Close()

	var resultSlice []*ContentUnitsPerson
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice content_units_persons")
	}

	if singular {
		object.R.ContentUnitsPersons = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.PersonID {
				local.R.ContentUnitsPersons = append(local.R.ContentUnitsPersons, foreign)
				break
			}
		}
	}

	return nil
}

// SetName of the person to the related item.
// Sets o.R.Name to related.
// Adds o to related.R.NamePersons.
func (o *Person) SetName(exec boil.Executor, insert bool, related *StringTranslation) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"persons\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"name_id"}),
		strmangle.WhereClause("\"", "\"", 2, personPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.NameID = related.ID

	if o.R == nil {
		o.R = &personR{
			Name: related,
		}
	} else {
		o.R.Name = related
	}

	if related.R == nil {
		related.R = &stringTranslationR{
			NamePersons: PersonSlice{o},
		}
	} else {
		related.R.NamePersons = append(related.R.NamePersons, o)
	}

	return nil
}

// SetDescription of the person to the related item.
// Sets o.R.Description to related.
// Adds o to related.R.DescriptionPersons.
func (o *Person) SetDescription(exec boil.Executor, insert bool, related *StringTranslation) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"persons\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"description_id"}),
		strmangle.WhereClause("\"", "\"", 2, personPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.DescriptionID.Int64 = related.ID
	o.DescriptionID.Valid = true

	if o.R == nil {
		o.R = &personR{
			Description: related,
		}
	} else {
		o.R.Description = related
	}

	if related.R == nil {
		related.R = &stringTranslationR{
			DescriptionPersons: PersonSlice{o},
		}
	} else {
		related.R.DescriptionPersons = append(related.R.DescriptionPersons, o)
	}

	return nil
}

// RemoveDescription relationship.
// Sets o.R.Description to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *Person) RemoveDescription(exec boil.Executor, related *StringTranslation) error {
	var err error

	o.DescriptionID.Valid = false
	if err = o.Update(exec, "description_id"); err != nil {
		o.DescriptionID.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.Description = nil
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.DescriptionPersons {
		if o.DescriptionID.Int64 != ri.DescriptionID.Int64 {
			continue
		}

		ln := len(related.R.DescriptionPersons)
		if ln > 1 && i < ln-1 {
			related.R.DescriptionPersons[i] = related.R.DescriptionPersons[ln-1]
		}
		related.R.DescriptionPersons = related.R.DescriptionPersons[:ln-1]
		break
	}
	return nil
}

// AddContentUnitsPersons adds the given related objects to the existing relationships
// of the person, optionally inserting them as new records.
// Appends related to o.R.ContentUnitsPersons.
// Sets related.R.Person appropriately.
func (o *Person) AddContentUnitsPersons(exec boil.Executor, insert bool, related ...*ContentUnitsPerson) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.PersonID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"content_units_persons\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"person_id"}),
				strmangle.WhereClause("\"", "\"", 2, contentUnitsPersonPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ContentUnitID, rel.PersonID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.PersonID = o.ID
		}
	}

	if o.R == nil {
		o.R = &personR{
			ContentUnitsPersons: related,
		}
	} else {
		o.R.ContentUnitsPersons = append(o.R.ContentUnitsPersons, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &contentUnitsPersonR{
				Person: o,
			}
		} else {
			rel.R.Person = o
		}
	}
	return nil
}

// PersonsG retrieves all records.
func PersonsG(mods ...qm.QueryMod) personQuery {
	return Persons(boil.GetDB(), mods...)
}

// Persons retrieves all the records using an executor.
func Persons(exec boil.Executor, mods ...qm.QueryMod) personQuery {
	mods = append(mods, qm.From("\"persons\""))
	return personQuery{NewQuery(exec, mods...)}
}

// FindPersonG retrieves a single record by ID.
func FindPersonG(id int64, selectCols ...string) (*Person, error) {
	return FindPerson(boil.GetDB(), id, selectCols...)
}

// FindPersonGP retrieves a single record by ID, and panics on error.
func FindPersonGP(id int64, selectCols ...string) *Person {
	retobj, err := FindPerson(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindPerson retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPerson(exec boil.Executor, id int64, selectCols ...string) (*Person, error) {
	personObj := &Person{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"persons\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(personObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "gmodels: unable to select from persons")
	}

	return personObj, nil
}

// FindPersonP retrieves a single record by ID with an executor, and panics on error.
func FindPersonP(exec boil.Executor, id int64, selectCols ...string) *Person {
	retobj, err := FindPerson(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Person) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Person) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Person) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Person) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("gmodels: no persons provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(personColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	personInsertCacheMut.RLock()
	cache, cached := personInsertCache[key]
	personInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			personColumns,
			personColumnsWithDefault,
			personColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(personType, personMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(personType, personMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"persons\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "gmodels: unable to insert into persons")
	}

	if !cached {
		personInsertCacheMut.Lock()
		personInsertCache[key] = cache
		personInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single Person record. See Update for
// whitelist behavior description.
func (o *Person) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Person record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Person) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Person, and panics on error.
// See Update for whitelist behavior description.
func (o *Person) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Person.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Person) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	personUpdateCacheMut.RLock()
	cache, cached := personUpdateCache[key]
	personUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(personColumns, personPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("gmodels: unable to update persons, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"persons\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, personPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(personType, personMapping, append(wl, personPrimaryKeyColumns...))
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
		return errors.Wrap(err, "gmodels: unable to update persons row")
	}

	if !cached {
		personUpdateCacheMut.Lock()
		personUpdateCache[key] = cache
		personUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q personQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q personQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to update all for persons")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o PersonSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o PersonSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o PersonSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PersonSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), personPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"persons\" SET %s WHERE (\"id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(personPrimaryKeyColumns), len(colNames)+1, len(personPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to update all in person slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Person) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Person) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Person) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Person) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("gmodels: no persons provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(personColumnsWithDefault, o)

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

	personUpsertCacheMut.RLock()
	cache, cached := personUpsertCache[key]
	personUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			personColumns,
			personColumnsWithDefault,
			personColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			personColumns,
			personPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("gmodels: unable to upsert persons, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(personPrimaryKeyColumns))
			copy(conflict, personPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"persons\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(personType, personMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(personType, personMapping, ret)
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
		return errors.Wrap(err, "gmodels: unable to upsert persons")
	}

	if !cached {
		personUpsertCacheMut.Lock()
		personUpsertCache[key] = cache
		personUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single Person record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Person) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Person record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Person) DeleteG() error {
	if o == nil {
		return errors.New("gmodels: no Person provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Person record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Person) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Person record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Person) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("gmodels: no Person provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), personPrimaryKeyMapping)
	sql := "DELETE FROM \"persons\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete from persons")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q personQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q personQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("gmodels: no personQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete all from persons")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o PersonSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o PersonSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("gmodels: no Person slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o PersonSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PersonSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("gmodels: no Person slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), personPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"persons\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, personPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(personPrimaryKeyColumns), 1, len(personPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete all from person slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Person) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Person) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Person) ReloadG() error {
	if o == nil {
		return errors.New("gmodels: no Person provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Person) Reload(exec boil.Executor) error {
	ret, err := FindPerson(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *PersonSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *PersonSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PersonSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("gmodels: empty PersonSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PersonSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	persons := PersonSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), personPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"persons\".* FROM \"persons\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, personPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(personPrimaryKeyColumns), 1, len(personPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&persons)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to reload all in PersonSlice")
	}

	*o = persons

	return nil
}

// PersonExists checks if the Person row exists.
func PersonExists(exec boil.Executor, id int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"persons\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "gmodels: unable to check if persons exists")
	}

	return exists, nil
}

// PersonExistsG checks if the Person row exists.
func PersonExistsG(id int64) (bool, error) {
	return PersonExists(boil.GetDB(), id)
}

// PersonExistsGP checks if the Person row exists. Panics on error.
func PersonExistsGP(id int64) bool {
	e, err := PersonExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// PersonExistsP checks if the Person row exists. Panics on error.
func PersonExistsP(exec boil.Executor, id int64) bool {
	e, err := PersonExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
