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

// ContentRole is an object representing the database table.
type ContentRole struct {
	ID            int64      `boil:"id" json:"id" toml:"id" yaml:"id"`
	NameID        int64      `boil:"name_id" json:"name_id" toml:"name_id" yaml:"name_id"`
	DescriptionID null.Int64 `boil:"description_id" json:"description_id,omitempty" toml:"description_id" yaml:"description_id,omitempty"`

	R *contentRoleR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L contentRoleL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// contentRoleR is where relationships are stored.
type contentRoleR struct {
	Name                    *StringTranslation
	Description             *StringTranslation
	RoleContentUnitsPersons ContentUnitsPersonSlice
}

// contentRoleL is where Load methods for each relationship are stored.
type contentRoleL struct{}

var (
	contentRoleColumns               = []string{"id", "name_id", "description_id"}
	contentRoleColumnsWithoutDefault = []string{"name_id", "description_id"}
	contentRoleColumnsWithDefault    = []string{"id"}
	contentRolePrimaryKeyColumns     = []string{"id"}
)

type (
	// ContentRoleSlice is an alias for a slice of pointers to ContentRole.
	// This should generally be used opposed to []ContentRole.
	ContentRoleSlice []*ContentRole

	contentRoleQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	contentRoleType                 = reflect.TypeOf(&ContentRole{})
	contentRoleMapping              = queries.MakeStructMapping(contentRoleType)
	contentRolePrimaryKeyMapping, _ = queries.BindMapping(contentRoleType, contentRoleMapping, contentRolePrimaryKeyColumns)
	contentRoleInsertCacheMut       sync.RWMutex
	contentRoleInsertCache          = make(map[string]insertCache)
	contentRoleUpdateCacheMut       sync.RWMutex
	contentRoleUpdateCache          = make(map[string]updateCache)
	contentRoleUpsertCacheMut       sync.RWMutex
	contentRoleUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single contentRole record from the query, and panics on error.
func (q contentRoleQuery) OneP() *ContentRole {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single contentRole record from the query.
func (q contentRoleQuery) One() (*ContentRole, error) {
	o := &ContentRole{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "gmodels: failed to execute a one query for content_roles")
	}

	return o, nil
}

// AllP returns all ContentRole records from the query, and panics on error.
func (q contentRoleQuery) AllP() ContentRoleSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all ContentRole records from the query.
func (q contentRoleQuery) All() (ContentRoleSlice, error) {
	var o ContentRoleSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "gmodels: failed to assign all query results to ContentRole slice")
	}

	return o, nil
}

// CountP returns the count of all ContentRole records in the query, and panics on error.
func (q contentRoleQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all ContentRole records in the query.
func (q contentRoleQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "gmodels: failed to count content_roles rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q contentRoleQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q contentRoleQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "gmodels: failed to check if content_roles exists")
	}

	return count > 0, nil
}

// NameG pointed to by the foreign key.
func (o *ContentRole) NameG(mods ...qm.QueryMod) stringTranslationQuery {
	return o.Name(boil.GetDB(), mods...)
}

// Name pointed to by the foreign key.
func (o *ContentRole) Name(exec boil.Executor, mods ...qm.QueryMod) stringTranslationQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.NameID),
	}

	queryMods = append(queryMods, mods...)

	query := StringTranslations(exec, queryMods...)
	queries.SetFrom(query.Query, "\"string_translations\"")

	return query
}

// DescriptionG pointed to by the foreign key.
func (o *ContentRole) DescriptionG(mods ...qm.QueryMod) stringTranslationQuery {
	return o.Description(boil.GetDB(), mods...)
}

// Description pointed to by the foreign key.
func (o *ContentRole) Description(exec boil.Executor, mods ...qm.QueryMod) stringTranslationQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.DescriptionID),
	}

	queryMods = append(queryMods, mods...)

	query := StringTranslations(exec, queryMods...)
	queries.SetFrom(query.Query, "\"string_translations\"")

	return query
}

// RoleContentUnitsPersonsG retrieves all the content_units_person's content units persons via role_id column.
func (o *ContentRole) RoleContentUnitsPersonsG(mods ...qm.QueryMod) contentUnitsPersonQuery {
	return o.RoleContentUnitsPersons(boil.GetDB(), mods...)
}

// RoleContentUnitsPersons retrieves all the content_units_person's content units persons with an executor via role_id column.
func (o *ContentRole) RoleContentUnitsPersons(exec boil.Executor, mods ...qm.QueryMod) contentUnitsPersonQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"role_id\"=?", o.ID),
	)

	query := ContentUnitsPersons(exec, queryMods...)
	queries.SetFrom(query.Query, "\"content_units_persons\" as \"a\"")
	return query
}

// LoadName allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (contentRoleL) LoadName(e boil.Executor, singular bool, maybeContentRole interface{}) error {
	var slice []*ContentRole
	var object *ContentRole

	count := 1
	if singular {
		object = maybeContentRole.(*ContentRole)
	} else {
		slice = *maybeContentRole.(*ContentRoleSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &contentRoleR{}
		}
		args[0] = object.NameID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &contentRoleR{}
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
func (contentRoleL) LoadDescription(e boil.Executor, singular bool, maybeContentRole interface{}) error {
	var slice []*ContentRole
	var object *ContentRole

	count := 1
	if singular {
		object = maybeContentRole.(*ContentRole)
	} else {
		slice = *maybeContentRole.(*ContentRoleSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &contentRoleR{}
		}
		args[0] = object.DescriptionID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &contentRoleR{}
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

// LoadRoleContentUnitsPersons allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (contentRoleL) LoadRoleContentUnitsPersons(e boil.Executor, singular bool, maybeContentRole interface{}) error {
	var slice []*ContentRole
	var object *ContentRole

	count := 1
	if singular {
		object = maybeContentRole.(*ContentRole)
	} else {
		slice = *maybeContentRole.(*ContentRoleSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &contentRoleR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &contentRoleR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"content_units_persons\" where \"role_id\" in (%s)",
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
		object.R.RoleContentUnitsPersons = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.RoleID {
				local.R.RoleContentUnitsPersons = append(local.R.RoleContentUnitsPersons, foreign)
				break
			}
		}
	}

	return nil
}

// SetName of the content_role to the related item.
// Sets o.R.Name to related.
// Adds o to related.R.NameContentRoles.
func (o *ContentRole) SetName(exec boil.Executor, insert bool, related *StringTranslation) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"content_roles\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"name_id"}),
		strmangle.WhereClause("\"", "\"", 2, contentRolePrimaryKeyColumns),
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
		o.R = &contentRoleR{
			Name: related,
		}
	} else {
		o.R.Name = related
	}

	if related.R == nil {
		related.R = &stringTranslationR{
			NameContentRoles: ContentRoleSlice{o},
		}
	} else {
		related.R.NameContentRoles = append(related.R.NameContentRoles, o)
	}

	return nil
}

// SetDescription of the content_role to the related item.
// Sets o.R.Description to related.
// Adds o to related.R.DescriptionContentRoles.
func (o *ContentRole) SetDescription(exec boil.Executor, insert bool, related *StringTranslation) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"content_roles\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"description_id"}),
		strmangle.WhereClause("\"", "\"", 2, contentRolePrimaryKeyColumns),
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
		o.R = &contentRoleR{
			Description: related,
		}
	} else {
		o.R.Description = related
	}

	if related.R == nil {
		related.R = &stringTranslationR{
			DescriptionContentRoles: ContentRoleSlice{o},
		}
	} else {
		related.R.DescriptionContentRoles = append(related.R.DescriptionContentRoles, o)
	}

	return nil
}

// RemoveDescription relationship.
// Sets o.R.Description to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *ContentRole) RemoveDescription(exec boil.Executor, related *StringTranslation) error {
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

	for i, ri := range related.R.DescriptionContentRoles {
		if o.DescriptionID.Int64 != ri.DescriptionID.Int64 {
			continue
		}

		ln := len(related.R.DescriptionContentRoles)
		if ln > 1 && i < ln-1 {
			related.R.DescriptionContentRoles[i] = related.R.DescriptionContentRoles[ln-1]
		}
		related.R.DescriptionContentRoles = related.R.DescriptionContentRoles[:ln-1]
		break
	}
	return nil
}

// AddRoleContentUnitsPersons adds the given related objects to the existing relationships
// of the content_role, optionally inserting them as new records.
// Appends related to o.R.RoleContentUnitsPersons.
// Sets related.R.Role appropriately.
func (o *ContentRole) AddRoleContentUnitsPersons(exec boil.Executor, insert bool, related ...*ContentUnitsPerson) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.RoleID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"content_units_persons\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"role_id"}),
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

			rel.RoleID = o.ID
		}
	}

	if o.R == nil {
		o.R = &contentRoleR{
			RoleContentUnitsPersons: related,
		}
	} else {
		o.R.RoleContentUnitsPersons = append(o.R.RoleContentUnitsPersons, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &contentUnitsPersonR{
				Role: o,
			}
		} else {
			rel.R.Role = o
		}
	}
	return nil
}

// ContentRolesG retrieves all records.
func ContentRolesG(mods ...qm.QueryMod) contentRoleQuery {
	return ContentRoles(boil.GetDB(), mods...)
}

// ContentRoles retrieves all the records using an executor.
func ContentRoles(exec boil.Executor, mods ...qm.QueryMod) contentRoleQuery {
	mods = append(mods, qm.From("\"content_roles\""))
	return contentRoleQuery{NewQuery(exec, mods...)}
}

// FindContentRoleG retrieves a single record by ID.
func FindContentRoleG(id int64, selectCols ...string) (*ContentRole, error) {
	return FindContentRole(boil.GetDB(), id, selectCols...)
}

// FindContentRoleGP retrieves a single record by ID, and panics on error.
func FindContentRoleGP(id int64, selectCols ...string) *ContentRole {
	retobj, err := FindContentRole(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindContentRole retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindContentRole(exec boil.Executor, id int64, selectCols ...string) (*ContentRole, error) {
	contentRoleObj := &ContentRole{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"content_roles\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(contentRoleObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "gmodels: unable to select from content_roles")
	}

	return contentRoleObj, nil
}

// FindContentRoleP retrieves a single record by ID with an executor, and panics on error.
func FindContentRoleP(exec boil.Executor, id int64, selectCols ...string) *ContentRole {
	retobj, err := FindContentRole(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *ContentRole) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *ContentRole) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *ContentRole) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *ContentRole) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("gmodels: no content_roles provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(contentRoleColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	contentRoleInsertCacheMut.RLock()
	cache, cached := contentRoleInsertCache[key]
	contentRoleInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			contentRoleColumns,
			contentRoleColumnsWithDefault,
			contentRoleColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(contentRoleType, contentRoleMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(contentRoleType, contentRoleMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"content_roles\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "gmodels: unable to insert into content_roles")
	}

	if !cached {
		contentRoleInsertCacheMut.Lock()
		contentRoleInsertCache[key] = cache
		contentRoleInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single ContentRole record. See Update for
// whitelist behavior description.
func (o *ContentRole) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single ContentRole record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *ContentRole) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the ContentRole, and panics on error.
// See Update for whitelist behavior description.
func (o *ContentRole) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the ContentRole.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *ContentRole) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	contentRoleUpdateCacheMut.RLock()
	cache, cached := contentRoleUpdateCache[key]
	contentRoleUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(contentRoleColumns, contentRolePrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("gmodels: unable to update content_roles, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"content_roles\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, contentRolePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(contentRoleType, contentRoleMapping, append(wl, contentRolePrimaryKeyColumns...))
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
		return errors.Wrap(err, "gmodels: unable to update content_roles row")
	}

	if !cached {
		contentRoleUpdateCacheMut.Lock()
		contentRoleUpdateCache[key] = cache
		contentRoleUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q contentRoleQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q contentRoleQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to update all for content_roles")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o ContentRoleSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o ContentRoleSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o ContentRoleSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ContentRoleSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contentRolePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"content_roles\" SET %s WHERE (\"id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(contentRolePrimaryKeyColumns), len(colNames)+1, len(contentRolePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to update all in contentRole slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *ContentRole) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *ContentRole) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *ContentRole) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *ContentRole) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("gmodels: no content_roles provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(contentRoleColumnsWithDefault, o)

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

	contentRoleUpsertCacheMut.RLock()
	cache, cached := contentRoleUpsertCache[key]
	contentRoleUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			contentRoleColumns,
			contentRoleColumnsWithDefault,
			contentRoleColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			contentRoleColumns,
			contentRolePrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("gmodels: unable to upsert content_roles, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(contentRolePrimaryKeyColumns))
			copy(conflict, contentRolePrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"content_roles\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(contentRoleType, contentRoleMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(contentRoleType, contentRoleMapping, ret)
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
		return errors.Wrap(err, "gmodels: unable to upsert content_roles")
	}

	if !cached {
		contentRoleUpsertCacheMut.Lock()
		contentRoleUpsertCache[key] = cache
		contentRoleUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single ContentRole record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *ContentRole) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single ContentRole record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *ContentRole) DeleteG() error {
	if o == nil {
		return errors.New("gmodels: no ContentRole provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single ContentRole record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *ContentRole) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single ContentRole record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ContentRole) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("gmodels: no ContentRole provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), contentRolePrimaryKeyMapping)
	sql := "DELETE FROM \"content_roles\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete from content_roles")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q contentRoleQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q contentRoleQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("gmodels: no contentRoleQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete all from content_roles")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o ContentRoleSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o ContentRoleSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("gmodels: no ContentRole slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o ContentRoleSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ContentRoleSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("gmodels: no ContentRole slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contentRolePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"content_roles\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, contentRolePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(contentRolePrimaryKeyColumns), 1, len(contentRolePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete all from contentRole slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *ContentRole) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *ContentRole) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *ContentRole) ReloadG() error {
	if o == nil {
		return errors.New("gmodels: no ContentRole provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ContentRole) Reload(exec boil.Executor) error {
	ret, err := FindContentRole(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ContentRoleSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ContentRoleSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ContentRoleSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("gmodels: empty ContentRoleSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ContentRoleSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	contentRoles := ContentRoleSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contentRolePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"content_roles\".* FROM \"content_roles\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, contentRolePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(contentRolePrimaryKeyColumns), 1, len(contentRolePrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&contentRoles)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to reload all in ContentRoleSlice")
	}

	*o = contentRoles

	return nil
}

// ContentRoleExists checks if the ContentRole row exists.
func ContentRoleExists(exec boil.Executor, id int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"content_roles\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "gmodels: unable to check if content_roles exists")
	}

	return exists, nil
}

// ContentRoleExistsG checks if the ContentRole row exists.
func ContentRoleExistsG(id int64) (bool, error) {
	return ContentRoleExists(boil.GetDB(), id)
}

// ContentRoleExistsGP checks if the ContentRole row exists. Panics on error.
func ContentRoleExistsGP(id int64) bool {
	e, err := ContentRoleExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// ContentRoleExistsP checks if the ContentRole row exists. Panics on error.
func ContentRoleExistsP(exec boil.Executor, id int64) bool {
	e, err := ContentRoleExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
