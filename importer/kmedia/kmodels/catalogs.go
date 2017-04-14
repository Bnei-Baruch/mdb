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

// Catalog is an object representing the database table.
type Catalog struct {
	ID              int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name            string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	ParentID        null.Int    `boil:"parent_id" json:"parent_id,omitempty" toml:"parent_id" yaml:"parent_id,omitempty"`
	CreatedAt       null.Time   `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`
	UpdatedAt       null.Time   `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
	Catorder        int         `boil:"catorder" json:"catorder" toml:"catorder" yaml:"catorder"`
	Secure          int         `boil:"secure" json:"secure" toml:"secure" yaml:"secure"`
	Visible         null.Bool   `boil:"visible" json:"visible,omitempty" toml:"visible" yaml:"visible,omitempty"`
	Open            null.Bool   `boil:"open" json:"open,omitempty" toml:"open" yaml:"open,omitempty"`
	Label           null.String `boil:"label" json:"label,omitempty" toml:"label" yaml:"label,omitempty"`
	SelectedCatalog null.Int    `boil:"selected_catalog" json:"selected_catalog,omitempty" toml:"selected_catalog" yaml:"selected_catalog,omitempty"`
	UserID          null.Int    `boil:"user_id" json:"user_id,omitempty" toml:"user_id" yaml:"user_id,omitempty"`
	BooksCatalog    null.Bool   `boil:"books_catalog" json:"books_catalog,omitempty" toml:"books_catalog" yaml:"books_catalog,omitempty"`

	R *catalogR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L catalogL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// catalogR is where relationships are stored.
type catalogR struct {
	Parent                       *Catalog
	User                         *User
	Containers                   ContainerSlice
	CatalogDescriptions          CatalogDescriptionSlice
	ParentCatalogs               CatalogSlice
	ContainerDescriptionPatterns ContainerDescriptionPatternSlice
}

// catalogL is where Load methods for each relationship are stored.
type catalogL struct{}

var (
	catalogColumns               = []string{"id", "name", "parent_id", "created_at", "updated_at", "catorder", "secure", "visible", "open", "label", "selected_catalog", "user_id", "books_catalog"}
	catalogColumnsWithoutDefault = []string{"parent_id", "created_at", "updated_at", "label", "selected_catalog", "user_id", "books_catalog"}
	catalogColumnsWithDefault    = []string{"id", "name", "catorder", "secure", "visible", "open"}
	catalogPrimaryKeyColumns     = []string{"id"}
)

type (
	// CatalogSlice is an alias for a slice of pointers to Catalog.
	// This should generally be used opposed to []Catalog.
	CatalogSlice []*Catalog

	catalogQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	catalogType                 = reflect.TypeOf(&Catalog{})
	catalogMapping              = queries.MakeStructMapping(catalogType)
	catalogPrimaryKeyMapping, _ = queries.BindMapping(catalogType, catalogMapping, catalogPrimaryKeyColumns)
	catalogInsertCacheMut       sync.RWMutex
	catalogInsertCache          = make(map[string]insertCache)
	catalogUpdateCacheMut       sync.RWMutex
	catalogUpdateCache          = make(map[string]updateCache)
	catalogUpsertCacheMut       sync.RWMutex
	catalogUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single catalog record from the query, and panics on error.
func (q catalogQuery) OneP() *Catalog {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single catalog record from the query.
func (q catalogQuery) One() (*Catalog, error) {
	o := &Catalog{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "kmodels: failed to execute a one query for catalogs")
	}

	return o, nil
}

// AllP returns all Catalog records from the query, and panics on error.
func (q catalogQuery) AllP() CatalogSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Catalog records from the query.
func (q catalogQuery) All() (CatalogSlice, error) {
	var o CatalogSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "kmodels: failed to assign all query results to Catalog slice")
	}

	return o, nil
}

// CountP returns the count of all Catalog records in the query, and panics on error.
func (q catalogQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Catalog records in the query.
func (q catalogQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "kmodels: failed to count catalogs rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q catalogQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q catalogQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "kmodels: failed to check if catalogs exists")
	}

	return count > 0, nil
}

// ParentG pointed to by the foreign key.
func (o *Catalog) ParentG(mods ...qm.QueryMod) catalogQuery {
	return o.Parent(boil.GetDB(), mods...)
}

// Parent pointed to by the foreign key.
func (o *Catalog) Parent(exec boil.Executor, mods ...qm.QueryMod) catalogQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.ParentID),
	}

	queryMods = append(queryMods, mods...)

	query := Catalogs(exec, queryMods...)
	queries.SetFrom(query.Query, "\"catalogs\"")

	return query
}

// UserG pointed to by the foreign key.
func (o *Catalog) UserG(mods ...qm.QueryMod) userQuery {
	return o.User(boil.GetDB(), mods...)
}

// User pointed to by the foreign key.
func (o *Catalog) User(exec boil.Executor, mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(exec, queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// ContainersG retrieves all the container's containers.
func (o *Catalog) ContainersG(mods ...qm.QueryMod) containerQuery {
	return o.Containers(boil.GetDB(), mods...)
}

// Containers retrieves all the container's containers with an executor.
func (o *Catalog) Containers(exec boil.Executor, mods ...qm.QueryMod) containerQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.InnerJoin("\"catalogs_containers\" as \"b\" on \"a\".\"id\" = \"b\".\"container_id\""),
		qm.Where("\"b\".\"catalog_id\"=?", o.ID),
	)

	query := Containers(exec, queryMods...)
	queries.SetFrom(query.Query, "\"containers\" as \"a\"")
	return query
}

// CatalogDescriptionsG retrieves all the catalog_description's catalog descriptions.
func (o *Catalog) CatalogDescriptionsG(mods ...qm.QueryMod) catalogDescriptionQuery {
	return o.CatalogDescriptions(boil.GetDB(), mods...)
}

// CatalogDescriptions retrieves all the catalog_description's catalog descriptions with an executor.
func (o *Catalog) CatalogDescriptions(exec boil.Executor, mods ...qm.QueryMod) catalogDescriptionQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"catalog_id\"=?", o.ID),
	)

	query := CatalogDescriptions(exec, queryMods...)
	queries.SetFrom(query.Query, "\"catalog_descriptions\" as \"a\"")
	return query
}

// ParentCatalogsG retrieves all the catalog's catalogs via parent_id column.
func (o *Catalog) ParentCatalogsG(mods ...qm.QueryMod) catalogQuery {
	return o.ParentCatalogs(boil.GetDB(), mods...)
}

// ParentCatalogs retrieves all the catalog's catalogs with an executor via parent_id column.
func (o *Catalog) ParentCatalogs(exec boil.Executor, mods ...qm.QueryMod) catalogQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"parent_id\"=?", o.ID),
	)

	query := Catalogs(exec, queryMods...)
	queries.SetFrom(query.Query, "\"catalogs\" as \"a\"")
	return query
}

// ContainerDescriptionPatternsG retrieves all the container_description_pattern's container description patterns.
func (o *Catalog) ContainerDescriptionPatternsG(mods ...qm.QueryMod) containerDescriptionPatternQuery {
	return o.ContainerDescriptionPatterns(boil.GetDB(), mods...)
}

// ContainerDescriptionPatterns retrieves all the container_description_pattern's container description patterns with an executor.
func (o *Catalog) ContainerDescriptionPatterns(exec boil.Executor, mods ...qm.QueryMod) containerDescriptionPatternQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.InnerJoin("\"catalogs_container_description_patterns\" as \"b\" on \"a\".\"id\" = \"b\".\"container_description_pattern_id\""),
		qm.Where("\"b\".\"catalog_id\"=?", o.ID),
	)

	query := ContainerDescriptionPatterns(exec, queryMods...)
	queries.SetFrom(query.Query, "\"container_description_patterns\" as \"a\"")
	return query
}

// LoadParent allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (catalogL) LoadParent(e boil.Executor, singular bool, maybeCatalog interface{}) error {
	var slice []*Catalog
	var object *Catalog

	count := 1
	if singular {
		object = maybeCatalog.(*Catalog)
	} else {
		slice = *maybeCatalog.(*CatalogSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &catalogR{}
		}
		args[0] = object.ParentID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &catalogR{}
			}
			args[i] = obj.ParentID
		}
	}

	query := fmt.Sprintf(
		"select * from \"catalogs\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Catalog")
	}
	defer results.Close()

	var resultSlice []*Catalog
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Catalog")
	}

	if singular && len(resultSlice) != 0 {
		object.R.Parent = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ParentID.Int == foreign.ID {
				local.R.Parent = foreign
				break
			}
		}
	}

	return nil
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (catalogL) LoadUser(e boil.Executor, singular bool, maybeCatalog interface{}) error {
	var slice []*Catalog
	var object *Catalog

	count := 1
	if singular {
		object = maybeCatalog.(*Catalog)
	} else {
		slice = *maybeCatalog.(*CatalogSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &catalogR{}
		}
		args[0] = object.UserID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &catalogR{}
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
			if local.UserID.Int == foreign.ID {
				local.R.User = foreign
				break
			}
		}
	}

	return nil
}

// LoadContainers allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (catalogL) LoadContainers(e boil.Executor, singular bool, maybeCatalog interface{}) error {
	var slice []*Catalog
	var object *Catalog

	count := 1
	if singular {
		object = maybeCatalog.(*Catalog)
	} else {
		slice = *maybeCatalog.(*CatalogSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &catalogR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &catalogR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select \"a\".*, \"b\".\"catalog_id\" from \"containers\" as \"a\" inner join \"catalogs_containers\" as \"b\" on \"a\".\"id\" = \"b\".\"container_id\" where \"b\".\"catalog_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load containers")
	}
	defer results.Close()

	var resultSlice []*Container

	var localJoinCols []int
	for results.Next() {
		one := new(Container)
		var localJoinCol int

		err = results.Scan(&one.ID, &one.Name, &one.CreatedAt, &one.UpdatedAt, &one.Filmdate, &one.LangID, &one.LecturerID, &one.Secure, &one.ContentTypeID, &one.MarkedForMerge, &one.SecureChanged, &one.AutoParsed, &one.VirtualLessonID, &one.PlaytimeSecs, &one.UserID, &one.ForCensorship, &one.OpenedByCensor, &one.ClosedByCensor, &one.CensorID, &one.Position, &localJoinCol)
		if err = results.Err(); err != nil {
			return errors.Wrap(err, "failed to plebian-bind eager loaded slice containers")
		}

		resultSlice = append(resultSlice, one)
		localJoinCols = append(localJoinCols, localJoinCol)
	}

	if err = results.Err(); err != nil {
		return errors.Wrap(err, "failed to plebian-bind eager loaded slice containers")
	}

	if singular {
		object.R.Containers = resultSlice
		return nil
	}

	for i, foreign := range resultSlice {
		localJoinCol := localJoinCols[i]
		for _, local := range slice {
			if local.ID == localJoinCol {
				local.R.Containers = append(local.R.Containers, foreign)
				break
			}
		}
	}

	return nil
}

// LoadCatalogDescriptions allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (catalogL) LoadCatalogDescriptions(e boil.Executor, singular bool, maybeCatalog interface{}) error {
	var slice []*Catalog
	var object *Catalog

	count := 1
	if singular {
		object = maybeCatalog.(*Catalog)
	} else {
		slice = *maybeCatalog.(*CatalogSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &catalogR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &catalogR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"catalog_descriptions\" where \"catalog_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load catalog_descriptions")
	}
	defer results.Close()

	var resultSlice []*CatalogDescription
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice catalog_descriptions")
	}

	if singular {
		object.R.CatalogDescriptions = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.CatalogID {
				local.R.CatalogDescriptions = append(local.R.CatalogDescriptions, foreign)
				break
			}
		}
	}

	return nil
}

// LoadParentCatalogs allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (catalogL) LoadParentCatalogs(e boil.Executor, singular bool, maybeCatalog interface{}) error {
	var slice []*Catalog
	var object *Catalog

	count := 1
	if singular {
		object = maybeCatalog.(*Catalog)
	} else {
		slice = *maybeCatalog.(*CatalogSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &catalogR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &catalogR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"catalogs\" where \"parent_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load catalogs")
	}
	defer results.Close()

	var resultSlice []*Catalog
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice catalogs")
	}

	if singular {
		object.R.ParentCatalogs = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.ParentID.Int {
				local.R.ParentCatalogs = append(local.R.ParentCatalogs, foreign)
				break
			}
		}
	}

	return nil
}

// LoadContainerDescriptionPatterns allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (catalogL) LoadContainerDescriptionPatterns(e boil.Executor, singular bool, maybeCatalog interface{}) error {
	var slice []*Catalog
	var object *Catalog

	count := 1
	if singular {
		object = maybeCatalog.(*Catalog)
	} else {
		slice = *maybeCatalog.(*CatalogSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &catalogR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &catalogR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select \"a\".*, \"b\".\"catalog_id\" from \"container_description_patterns\" as \"a\" inner join \"catalogs_container_description_patterns\" as \"b\" on \"a\".\"id\" = \"b\".\"container_description_pattern_id\" where \"b\".\"catalog_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load container_description_patterns")
	}
	defer results.Close()

	var resultSlice []*ContainerDescriptionPattern

	var localJoinCols []int
	for results.Next() {
		one := new(ContainerDescriptionPattern)
		var localJoinCol int

		err = results.Scan(&one.ID, &one.Pattern, &one.Description, &one.Lang, &one.CreatedAt, &one.UpdatedAt, &one.UserID, &localJoinCol)
		if err = results.Err(); err != nil {
			return errors.Wrap(err, "failed to plebian-bind eager loaded slice container_description_patterns")
		}

		resultSlice = append(resultSlice, one)
		localJoinCols = append(localJoinCols, localJoinCol)
	}

	if err = results.Err(); err != nil {
		return errors.Wrap(err, "failed to plebian-bind eager loaded slice container_description_patterns")
	}

	if singular {
		object.R.ContainerDescriptionPatterns = resultSlice
		return nil
	}

	for i, foreign := range resultSlice {
		localJoinCol := localJoinCols[i]
		for _, local := range slice {
			if local.ID == localJoinCol {
				local.R.ContainerDescriptionPatterns = append(local.R.ContainerDescriptionPatterns, foreign)
				break
			}
		}
	}

	return nil
}

// SetParentG of the catalog to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentCatalogs.
// Uses the global database handle.
func (o *Catalog) SetParentG(insert bool, related *Catalog) error {
	return o.SetParent(boil.GetDB(), insert, related)
}

// SetParentP of the catalog to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentCatalogs.
// Panics on error.
func (o *Catalog) SetParentP(exec boil.Executor, insert bool, related *Catalog) {
	if err := o.SetParent(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParentGP of the catalog to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentCatalogs.
// Uses the global database handle and panics on error.
func (o *Catalog) SetParentGP(insert bool, related *Catalog) {
	if err := o.SetParent(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParent of the catalog to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentCatalogs.
func (o *Catalog) SetParent(exec boil.Executor, insert bool, related *Catalog) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"catalogs\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"parent_id"}),
		strmangle.WhereClause("\"", "\"", 2, catalogPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ParentID.Int = related.ID
	o.ParentID.Valid = true

	if o.R == nil {
		o.R = &catalogR{
			Parent: related,
		}
	} else {
		o.R.Parent = related
	}

	if related.R == nil {
		related.R = &catalogR{
			ParentCatalogs: CatalogSlice{o},
		}
	} else {
		related.R.ParentCatalogs = append(related.R.ParentCatalogs, o)
	}

	return nil
}

// RemoveParentG relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle.
func (o *Catalog) RemoveParentG(related *Catalog) error {
	return o.RemoveParent(boil.GetDB(), related)
}

// RemoveParentP relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Panics on error.
func (o *Catalog) RemoveParentP(exec boil.Executor, related *Catalog) {
	if err := o.RemoveParent(exec, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParentGP relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle and panics on error.
func (o *Catalog) RemoveParentGP(related *Catalog) {
	if err := o.RemoveParent(boil.GetDB(), related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParent relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *Catalog) RemoveParent(exec boil.Executor, related *Catalog) error {
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

	for i, ri := range related.R.ParentCatalogs {
		if o.ParentID.Int != ri.ParentID.Int {
			continue
		}

		ln := len(related.R.ParentCatalogs)
		if ln > 1 && i < ln-1 {
			related.R.ParentCatalogs[i] = related.R.ParentCatalogs[ln-1]
		}
		related.R.ParentCatalogs = related.R.ParentCatalogs[:ln-1]
		break
	}
	return nil
}

// SetUserG of the catalog to the related item.
// Sets o.R.User to related.
// Adds o to related.R.Catalogs.
// Uses the global database handle.
func (o *Catalog) SetUserG(insert bool, related *User) error {
	return o.SetUser(boil.GetDB(), insert, related)
}

// SetUserP of the catalog to the related item.
// Sets o.R.User to related.
// Adds o to related.R.Catalogs.
// Panics on error.
func (o *Catalog) SetUserP(exec boil.Executor, insert bool, related *User) {
	if err := o.SetUser(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetUserGP of the catalog to the related item.
// Sets o.R.User to related.
// Adds o to related.R.Catalogs.
// Uses the global database handle and panics on error.
func (o *Catalog) SetUserGP(insert bool, related *User) {
	if err := o.SetUser(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetUser of the catalog to the related item.
// Sets o.R.User to related.
// Adds o to related.R.Catalogs.
func (o *Catalog) SetUser(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"catalogs\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, catalogPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.UserID.Int = related.ID
	o.UserID.Valid = true

	if o.R == nil {
		o.R = &catalogR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			Catalogs: CatalogSlice{o},
		}
	} else {
		related.R.Catalogs = append(related.R.Catalogs, o)
	}

	return nil
}

// RemoveUserG relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle.
func (o *Catalog) RemoveUserG(related *User) error {
	return o.RemoveUser(boil.GetDB(), related)
}

// RemoveUserP relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Panics on error.
func (o *Catalog) RemoveUserP(exec boil.Executor, related *User) {
	if err := o.RemoveUser(exec, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveUserGP relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle and panics on error.
func (o *Catalog) RemoveUserGP(related *User) {
	if err := o.RemoveUser(boil.GetDB(), related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveUser relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *Catalog) RemoveUser(exec boil.Executor, related *User) error {
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

	for i, ri := range related.R.Catalogs {
		if o.UserID.Int != ri.UserID.Int {
			continue
		}

		ln := len(related.R.Catalogs)
		if ln > 1 && i < ln-1 {
			related.R.Catalogs[i] = related.R.Catalogs[ln-1]
		}
		related.R.Catalogs = related.R.Catalogs[:ln-1]
		break
	}
	return nil
}

// AddContainersG adds the given related objects to the existing relationships
// of the catalog, optionally inserting them as new records.
// Appends related to o.R.Containers.
// Sets related.R.Catalogs appropriately.
// Uses the global database handle.
func (o *Catalog) AddContainersG(insert bool, related ...*Container) error {
	return o.AddContainers(boil.GetDB(), insert, related...)
}

// AddContainersP adds the given related objects to the existing relationships
// of the catalog, optionally inserting them as new records.
// Appends related to o.R.Containers.
// Sets related.R.Catalogs appropriately.
// Panics on error.
func (o *Catalog) AddContainersP(exec boil.Executor, insert bool, related ...*Container) {
	if err := o.AddContainers(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddContainersGP adds the given related objects to the existing relationships
// of the catalog, optionally inserting them as new records.
// Appends related to o.R.Containers.
// Sets related.R.Catalogs appropriately.
// Uses the global database handle and panics on error.
func (o *Catalog) AddContainersGP(insert bool, related ...*Container) {
	if err := o.AddContainers(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddContainers adds the given related objects to the existing relationships
// of the catalog, optionally inserting them as new records.
// Appends related to o.R.Containers.
// Sets related.R.Catalogs appropriately.
func (o *Catalog) AddContainers(exec boil.Executor, insert bool, related ...*Container) error {
	var err error
	for _, rel := range related {
		if insert {
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		}
	}

	for _, rel := range related {
		query := "insert into \"catalogs_containers\" (\"catalog_id\", \"container_id\") values ($1, $2)"
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
		o.R = &catalogR{
			Containers: related,
		}
	} else {
		o.R.Containers = append(o.R.Containers, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &containerR{
				Catalogs: CatalogSlice{o},
			}
		} else {
			rel.R.Catalogs = append(rel.R.Catalogs, o)
		}
	}
	return nil
}

// SetContainersG removes all previously related items of the
// catalog replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Catalogs's Containers accordingly.
// Replaces o.R.Containers with related.
// Sets related.R.Catalogs's Containers accordingly.
// Uses the global database handle.
func (o *Catalog) SetContainersG(insert bool, related ...*Container) error {
	return o.SetContainers(boil.GetDB(), insert, related...)
}

// SetContainersP removes all previously related items of the
// catalog replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Catalogs's Containers accordingly.
// Replaces o.R.Containers with related.
// Sets related.R.Catalogs's Containers accordingly.
// Panics on error.
func (o *Catalog) SetContainersP(exec boil.Executor, insert bool, related ...*Container) {
	if err := o.SetContainers(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetContainersGP removes all previously related items of the
// catalog replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Catalogs's Containers accordingly.
// Replaces o.R.Containers with related.
// Sets related.R.Catalogs's Containers accordingly.
// Uses the global database handle and panics on error.
func (o *Catalog) SetContainersGP(insert bool, related ...*Container) {
	if err := o.SetContainers(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetContainers removes all previously related items of the
// catalog replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Catalogs's Containers accordingly.
// Replaces o.R.Containers with related.
// Sets related.R.Catalogs's Containers accordingly.
func (o *Catalog) SetContainers(exec boil.Executor, insert bool, related ...*Container) error {
	query := "delete from \"catalogs_containers\" where \"catalog_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	removeContainersFromCatalogsSlice(o, related)
	o.R.Containers = nil
	return o.AddContainers(exec, insert, related...)
}

// RemoveContainersG relationships from objects passed in.
// Removes related items from R.Containers (uses pointer comparison, removal does not keep order)
// Sets related.R.Catalogs.
// Uses the global database handle.
func (o *Catalog) RemoveContainersG(related ...*Container) error {
	return o.RemoveContainers(boil.GetDB(), related...)
}

// RemoveContainersP relationships from objects passed in.
// Removes related items from R.Containers (uses pointer comparison, removal does not keep order)
// Sets related.R.Catalogs.
// Panics on error.
func (o *Catalog) RemoveContainersP(exec boil.Executor, related ...*Container) {
	if err := o.RemoveContainers(exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveContainersGP relationships from objects passed in.
// Removes related items from R.Containers (uses pointer comparison, removal does not keep order)
// Sets related.R.Catalogs.
// Uses the global database handle and panics on error.
func (o *Catalog) RemoveContainersGP(related ...*Container) {
	if err := o.RemoveContainers(boil.GetDB(), related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveContainers relationships from objects passed in.
// Removes related items from R.Containers (uses pointer comparison, removal does not keep order)
// Sets related.R.Catalogs.
func (o *Catalog) RemoveContainers(exec boil.Executor, related ...*Container) error {
	var err error
	query := fmt.Sprintf(
		"delete from \"catalogs_containers\" where \"catalog_id\" = $1 and \"container_id\" in (%s)",
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
	removeContainersFromCatalogsSlice(o, related)
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Containers {
			if rel != ri {
				continue
			}

			ln := len(o.R.Containers)
			if ln > 1 && i < ln-1 {
				o.R.Containers[i] = o.R.Containers[ln-1]
			}
			o.R.Containers = o.R.Containers[:ln-1]
			break
		}
	}

	return nil
}

func removeContainersFromCatalogsSlice(o *Catalog, related []*Container) {
	for _, rel := range related {
		if rel.R == nil {
			continue
		}
		for i, ri := range rel.R.Catalogs {
			if o.ID != ri.ID {
				continue
			}

			ln := len(rel.R.Catalogs)
			if ln > 1 && i < ln-1 {
				rel.R.Catalogs[i] = rel.R.Catalogs[ln-1]
			}
			rel.R.Catalogs = rel.R.Catalogs[:ln-1]
			break
		}
	}
}

// AddCatalogDescriptionsG adds the given related objects to the existing relationships
// of the catalog, optionally inserting them as new records.
// Appends related to o.R.CatalogDescriptions.
// Sets related.R.Catalog appropriately.
// Uses the global database handle.
func (o *Catalog) AddCatalogDescriptionsG(insert bool, related ...*CatalogDescription) error {
	return o.AddCatalogDescriptions(boil.GetDB(), insert, related...)
}

// AddCatalogDescriptionsP adds the given related objects to the existing relationships
// of the catalog, optionally inserting them as new records.
// Appends related to o.R.CatalogDescriptions.
// Sets related.R.Catalog appropriately.
// Panics on error.
func (o *Catalog) AddCatalogDescriptionsP(exec boil.Executor, insert bool, related ...*CatalogDescription) {
	if err := o.AddCatalogDescriptions(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddCatalogDescriptionsGP adds the given related objects to the existing relationships
// of the catalog, optionally inserting them as new records.
// Appends related to o.R.CatalogDescriptions.
// Sets related.R.Catalog appropriately.
// Uses the global database handle and panics on error.
func (o *Catalog) AddCatalogDescriptionsGP(insert bool, related ...*CatalogDescription) {
	if err := o.AddCatalogDescriptions(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddCatalogDescriptions adds the given related objects to the existing relationships
// of the catalog, optionally inserting them as new records.
// Appends related to o.R.CatalogDescriptions.
// Sets related.R.Catalog appropriately.
func (o *Catalog) AddCatalogDescriptions(exec boil.Executor, insert bool, related ...*CatalogDescription) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.CatalogID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"catalog_descriptions\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"catalog_id"}),
				strmangle.WhereClause("\"", "\"", 2, catalogDescriptionPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.CatalogID = o.ID
		}
	}

	if o.R == nil {
		o.R = &catalogR{
			CatalogDescriptions: related,
		}
	} else {
		o.R.CatalogDescriptions = append(o.R.CatalogDescriptions, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &catalogDescriptionR{
				Catalog: o,
			}
		} else {
			rel.R.Catalog = o
		}
	}
	return nil
}

// AddParentCatalogsG adds the given related objects to the existing relationships
// of the catalog, optionally inserting them as new records.
// Appends related to o.R.ParentCatalogs.
// Sets related.R.Parent appropriately.
// Uses the global database handle.
func (o *Catalog) AddParentCatalogsG(insert bool, related ...*Catalog) error {
	return o.AddParentCatalogs(boil.GetDB(), insert, related...)
}

// AddParentCatalogsP adds the given related objects to the existing relationships
// of the catalog, optionally inserting them as new records.
// Appends related to o.R.ParentCatalogs.
// Sets related.R.Parent appropriately.
// Panics on error.
func (o *Catalog) AddParentCatalogsP(exec boil.Executor, insert bool, related ...*Catalog) {
	if err := o.AddParentCatalogs(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddParentCatalogsGP adds the given related objects to the existing relationships
// of the catalog, optionally inserting them as new records.
// Appends related to o.R.ParentCatalogs.
// Sets related.R.Parent appropriately.
// Uses the global database handle and panics on error.
func (o *Catalog) AddParentCatalogsGP(insert bool, related ...*Catalog) {
	if err := o.AddParentCatalogs(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddParentCatalogs adds the given related objects to the existing relationships
// of the catalog, optionally inserting them as new records.
// Appends related to o.R.ParentCatalogs.
// Sets related.R.Parent appropriately.
func (o *Catalog) AddParentCatalogs(exec boil.Executor, insert bool, related ...*Catalog) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.ParentID.Int = o.ID
			rel.ParentID.Valid = true
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"catalogs\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"parent_id"}),
				strmangle.WhereClause("\"", "\"", 2, catalogPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.ParentID.Int = o.ID
			rel.ParentID.Valid = true
		}
	}

	if o.R == nil {
		o.R = &catalogR{
			ParentCatalogs: related,
		}
	} else {
		o.R.ParentCatalogs = append(o.R.ParentCatalogs, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &catalogR{
				Parent: o,
			}
		} else {
			rel.R.Parent = o
		}
	}
	return nil
}

// SetParentCatalogsG removes all previously related items of the
// catalog replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentCatalogs accordingly.
// Replaces o.R.ParentCatalogs with related.
// Sets related.R.Parent's ParentCatalogs accordingly.
// Uses the global database handle.
func (o *Catalog) SetParentCatalogsG(insert bool, related ...*Catalog) error {
	return o.SetParentCatalogs(boil.GetDB(), insert, related...)
}

// SetParentCatalogsP removes all previously related items of the
// catalog replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentCatalogs accordingly.
// Replaces o.R.ParentCatalogs with related.
// Sets related.R.Parent's ParentCatalogs accordingly.
// Panics on error.
func (o *Catalog) SetParentCatalogsP(exec boil.Executor, insert bool, related ...*Catalog) {
	if err := o.SetParentCatalogs(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParentCatalogsGP removes all previously related items of the
// catalog replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentCatalogs accordingly.
// Replaces o.R.ParentCatalogs with related.
// Sets related.R.Parent's ParentCatalogs accordingly.
// Uses the global database handle and panics on error.
func (o *Catalog) SetParentCatalogsGP(insert bool, related ...*Catalog) {
	if err := o.SetParentCatalogs(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParentCatalogs removes all previously related items of the
// catalog replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentCatalogs accordingly.
// Replaces o.R.ParentCatalogs with related.
// Sets related.R.Parent's ParentCatalogs accordingly.
func (o *Catalog) SetParentCatalogs(exec boil.Executor, insert bool, related ...*Catalog) error {
	query := "update \"catalogs\" set \"parent_id\" = null where \"parent_id\" = $1"
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
		for _, rel := range o.R.ParentCatalogs {
			rel.ParentID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Parent = nil
		}

		o.R.ParentCatalogs = nil
	}
	return o.AddParentCatalogs(exec, insert, related...)
}

// RemoveParentCatalogsG relationships from objects passed in.
// Removes related items from R.ParentCatalogs (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
// Uses the global database handle.
func (o *Catalog) RemoveParentCatalogsG(related ...*Catalog) error {
	return o.RemoveParentCatalogs(boil.GetDB(), related...)
}

// RemoveParentCatalogsP relationships from objects passed in.
// Removes related items from R.ParentCatalogs (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
// Panics on error.
func (o *Catalog) RemoveParentCatalogsP(exec boil.Executor, related ...*Catalog) {
	if err := o.RemoveParentCatalogs(exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParentCatalogsGP relationships from objects passed in.
// Removes related items from R.ParentCatalogs (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
// Uses the global database handle and panics on error.
func (o *Catalog) RemoveParentCatalogsGP(related ...*Catalog) {
	if err := o.RemoveParentCatalogs(boil.GetDB(), related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParentCatalogs relationships from objects passed in.
// Removes related items from R.ParentCatalogs (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
func (o *Catalog) RemoveParentCatalogs(exec boil.Executor, related ...*Catalog) error {
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
		for i, ri := range o.R.ParentCatalogs {
			if rel != ri {
				continue
			}

			ln := len(o.R.ParentCatalogs)
			if ln > 1 && i < ln-1 {
				o.R.ParentCatalogs[i] = o.R.ParentCatalogs[ln-1]
			}
			o.R.ParentCatalogs = o.R.ParentCatalogs[:ln-1]
			break
		}
	}

	return nil
}

// AddContainerDescriptionPatternsG adds the given related objects to the existing relationships
// of the catalog, optionally inserting them as new records.
// Appends related to o.R.ContainerDescriptionPatterns.
// Sets related.R.Catalogs appropriately.
// Uses the global database handle.
func (o *Catalog) AddContainerDescriptionPatternsG(insert bool, related ...*ContainerDescriptionPattern) error {
	return o.AddContainerDescriptionPatterns(boil.GetDB(), insert, related...)
}

// AddContainerDescriptionPatternsP adds the given related objects to the existing relationships
// of the catalog, optionally inserting them as new records.
// Appends related to o.R.ContainerDescriptionPatterns.
// Sets related.R.Catalogs appropriately.
// Panics on error.
func (o *Catalog) AddContainerDescriptionPatternsP(exec boil.Executor, insert bool, related ...*ContainerDescriptionPattern) {
	if err := o.AddContainerDescriptionPatterns(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddContainerDescriptionPatternsGP adds the given related objects to the existing relationships
// of the catalog, optionally inserting them as new records.
// Appends related to o.R.ContainerDescriptionPatterns.
// Sets related.R.Catalogs appropriately.
// Uses the global database handle and panics on error.
func (o *Catalog) AddContainerDescriptionPatternsGP(insert bool, related ...*ContainerDescriptionPattern) {
	if err := o.AddContainerDescriptionPatterns(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddContainerDescriptionPatterns adds the given related objects to the existing relationships
// of the catalog, optionally inserting them as new records.
// Appends related to o.R.ContainerDescriptionPatterns.
// Sets related.R.Catalogs appropriately.
func (o *Catalog) AddContainerDescriptionPatterns(exec boil.Executor, insert bool, related ...*ContainerDescriptionPattern) error {
	var err error
	for _, rel := range related {
		if insert {
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		}
	}

	for _, rel := range related {
		query := "insert into \"catalogs_container_description_patterns\" (\"catalog_id\", \"container_description_pattern_id\") values ($1, $2)"
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
		o.R = &catalogR{
			ContainerDescriptionPatterns: related,
		}
	} else {
		o.R.ContainerDescriptionPatterns = append(o.R.ContainerDescriptionPatterns, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &containerDescriptionPatternR{
				Catalogs: CatalogSlice{o},
			}
		} else {
			rel.R.Catalogs = append(rel.R.Catalogs, o)
		}
	}
	return nil
}

// SetContainerDescriptionPatternsG removes all previously related items of the
// catalog replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Catalogs's ContainerDescriptionPatterns accordingly.
// Replaces o.R.ContainerDescriptionPatterns with related.
// Sets related.R.Catalogs's ContainerDescriptionPatterns accordingly.
// Uses the global database handle.
func (o *Catalog) SetContainerDescriptionPatternsG(insert bool, related ...*ContainerDescriptionPattern) error {
	return o.SetContainerDescriptionPatterns(boil.GetDB(), insert, related...)
}

// SetContainerDescriptionPatternsP removes all previously related items of the
// catalog replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Catalogs's ContainerDescriptionPatterns accordingly.
// Replaces o.R.ContainerDescriptionPatterns with related.
// Sets related.R.Catalogs's ContainerDescriptionPatterns accordingly.
// Panics on error.
func (o *Catalog) SetContainerDescriptionPatternsP(exec boil.Executor, insert bool, related ...*ContainerDescriptionPattern) {
	if err := o.SetContainerDescriptionPatterns(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetContainerDescriptionPatternsGP removes all previously related items of the
// catalog replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Catalogs's ContainerDescriptionPatterns accordingly.
// Replaces o.R.ContainerDescriptionPatterns with related.
// Sets related.R.Catalogs's ContainerDescriptionPatterns accordingly.
// Uses the global database handle and panics on error.
func (o *Catalog) SetContainerDescriptionPatternsGP(insert bool, related ...*ContainerDescriptionPattern) {
	if err := o.SetContainerDescriptionPatterns(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetContainerDescriptionPatterns removes all previously related items of the
// catalog replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Catalogs's ContainerDescriptionPatterns accordingly.
// Replaces o.R.ContainerDescriptionPatterns with related.
// Sets related.R.Catalogs's ContainerDescriptionPatterns accordingly.
func (o *Catalog) SetContainerDescriptionPatterns(exec boil.Executor, insert bool, related ...*ContainerDescriptionPattern) error {
	query := "delete from \"catalogs_container_description_patterns\" where \"catalog_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	removeContainerDescriptionPatternsFromCatalogsSlice(o, related)
	o.R.ContainerDescriptionPatterns = nil
	return o.AddContainerDescriptionPatterns(exec, insert, related...)
}

// RemoveContainerDescriptionPatternsG relationships from objects passed in.
// Removes related items from R.ContainerDescriptionPatterns (uses pointer comparison, removal does not keep order)
// Sets related.R.Catalogs.
// Uses the global database handle.
func (o *Catalog) RemoveContainerDescriptionPatternsG(related ...*ContainerDescriptionPattern) error {
	return o.RemoveContainerDescriptionPatterns(boil.GetDB(), related...)
}

// RemoveContainerDescriptionPatternsP relationships from objects passed in.
// Removes related items from R.ContainerDescriptionPatterns (uses pointer comparison, removal does not keep order)
// Sets related.R.Catalogs.
// Panics on error.
func (o *Catalog) RemoveContainerDescriptionPatternsP(exec boil.Executor, related ...*ContainerDescriptionPattern) {
	if err := o.RemoveContainerDescriptionPatterns(exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveContainerDescriptionPatternsGP relationships from objects passed in.
// Removes related items from R.ContainerDescriptionPatterns (uses pointer comparison, removal does not keep order)
// Sets related.R.Catalogs.
// Uses the global database handle and panics on error.
func (o *Catalog) RemoveContainerDescriptionPatternsGP(related ...*ContainerDescriptionPattern) {
	if err := o.RemoveContainerDescriptionPatterns(boil.GetDB(), related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveContainerDescriptionPatterns relationships from objects passed in.
// Removes related items from R.ContainerDescriptionPatterns (uses pointer comparison, removal does not keep order)
// Sets related.R.Catalogs.
func (o *Catalog) RemoveContainerDescriptionPatterns(exec boil.Executor, related ...*ContainerDescriptionPattern) error {
	var err error
	query := fmt.Sprintf(
		"delete from \"catalogs_container_description_patterns\" where \"catalog_id\" = $1 and \"container_description_pattern_id\" in (%s)",
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
	removeContainerDescriptionPatternsFromCatalogsSlice(o, related)
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.ContainerDescriptionPatterns {
			if rel != ri {
				continue
			}

			ln := len(o.R.ContainerDescriptionPatterns)
			if ln > 1 && i < ln-1 {
				o.R.ContainerDescriptionPatterns[i] = o.R.ContainerDescriptionPatterns[ln-1]
			}
			o.R.ContainerDescriptionPatterns = o.R.ContainerDescriptionPatterns[:ln-1]
			break
		}
	}

	return nil
}

func removeContainerDescriptionPatternsFromCatalogsSlice(o *Catalog, related []*ContainerDescriptionPattern) {
	for _, rel := range related {
		if rel.R == nil {
			continue
		}
		for i, ri := range rel.R.Catalogs {
			if o.ID != ri.ID {
				continue
			}

			ln := len(rel.R.Catalogs)
			if ln > 1 && i < ln-1 {
				rel.R.Catalogs[i] = rel.R.Catalogs[ln-1]
			}
			rel.R.Catalogs = rel.R.Catalogs[:ln-1]
			break
		}
	}
}

// CatalogsG retrieves all records.
func CatalogsG(mods ...qm.QueryMod) catalogQuery {
	return Catalogs(boil.GetDB(), mods...)
}

// Catalogs retrieves all the records using an executor.
func Catalogs(exec boil.Executor, mods ...qm.QueryMod) catalogQuery {
	mods = append(mods, qm.From("\"catalogs\""))
	return catalogQuery{NewQuery(exec, mods...)}
}

// FindCatalogG retrieves a single record by ID.
func FindCatalogG(id int, selectCols ...string) (*Catalog, error) {
	return FindCatalog(boil.GetDB(), id, selectCols...)
}

// FindCatalogGP retrieves a single record by ID, and panics on error.
func FindCatalogGP(id int, selectCols ...string) *Catalog {
	retobj, err := FindCatalog(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindCatalog retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindCatalog(exec boil.Executor, id int, selectCols ...string) (*Catalog, error) {
	catalogObj := &Catalog{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"catalogs\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(catalogObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "kmodels: unable to select from catalogs")
	}

	return catalogObj, nil
}

// FindCatalogP retrieves a single record by ID with an executor, and panics on error.
func FindCatalogP(exec boil.Executor, id int, selectCols ...string) *Catalog {
	retobj, err := FindCatalog(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Catalog) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Catalog) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Catalog) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Catalog) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("kmodels: no catalogs provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.Time.IsZero() {
		o.CreatedAt.Time = currTime
		o.CreatedAt.Valid = true
	}
	if o.UpdatedAt.Time.IsZero() {
		o.UpdatedAt.Time = currTime
		o.UpdatedAt.Valid = true
	}

	nzDefaults := queries.NonZeroDefaultSet(catalogColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	catalogInsertCacheMut.RLock()
	cache, cached := catalogInsertCache[key]
	catalogInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			catalogColumns,
			catalogColumnsWithDefault,
			catalogColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(catalogType, catalogMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(catalogType, catalogMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"catalogs\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "kmodels: unable to insert into catalogs")
	}

	if !cached {
		catalogInsertCacheMut.Lock()
		catalogInsertCache[key] = cache
		catalogInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single Catalog record. See Update for
// whitelist behavior description.
func (o *Catalog) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Catalog record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Catalog) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Catalog, and panics on error.
// See Update for whitelist behavior description.
func (o *Catalog) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Catalog.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Catalog) Update(exec boil.Executor, whitelist ...string) error {
	currTime := time.Now().In(boil.GetLocation())

	o.UpdatedAt.Time = currTime
	o.UpdatedAt.Valid = true

	var err error
	key := makeCacheKey(whitelist, nil)
	catalogUpdateCacheMut.RLock()
	cache, cached := catalogUpdateCache[key]
	catalogUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(catalogColumns, catalogPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("kmodels: unable to update catalogs, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"catalogs\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, catalogPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(catalogType, catalogMapping, append(wl, catalogPrimaryKeyColumns...))
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
		return errors.Wrap(err, "kmodels: unable to update catalogs row")
	}

	if !cached {
		catalogUpdateCacheMut.Lock()
		catalogUpdateCache[key] = cache
		catalogUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q catalogQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q catalogQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to update all for catalogs")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o CatalogSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o CatalogSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o CatalogSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CatalogSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), catalogPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"catalogs\" SET %s WHERE (\"id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(catalogPrimaryKeyColumns), len(colNames)+1, len(catalogPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to update all in catalog slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Catalog) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Catalog) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Catalog) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Catalog) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("kmodels: no catalogs provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.Time.IsZero() {
		o.CreatedAt.Time = currTime
		o.CreatedAt.Valid = true
	}
	o.UpdatedAt.Time = currTime
	o.UpdatedAt.Valid = true

	nzDefaults := queries.NonZeroDefaultSet(catalogColumnsWithDefault, o)

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

	catalogUpsertCacheMut.RLock()
	cache, cached := catalogUpsertCache[key]
	catalogUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			catalogColumns,
			catalogColumnsWithDefault,
			catalogColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			catalogColumns,
			catalogPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("kmodels: unable to upsert catalogs, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(catalogPrimaryKeyColumns))
			copy(conflict, catalogPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"catalogs\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(catalogType, catalogMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(catalogType, catalogMapping, ret)
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
		return errors.Wrap(err, "kmodels: unable to upsert catalogs")
	}

	if !cached {
		catalogUpsertCacheMut.Lock()
		catalogUpsertCache[key] = cache
		catalogUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single Catalog record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Catalog) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Catalog record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Catalog) DeleteG() error {
	if o == nil {
		return errors.New("kmodels: no Catalog provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Catalog record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Catalog) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Catalog record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Catalog) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("kmodels: no Catalog provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), catalogPrimaryKeyMapping)
	sql := "DELETE FROM \"catalogs\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to delete from catalogs")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q catalogQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q catalogQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("kmodels: no catalogQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to delete all from catalogs")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o CatalogSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o CatalogSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("kmodels: no Catalog slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o CatalogSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CatalogSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("kmodels: no Catalog slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), catalogPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"catalogs\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, catalogPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(catalogPrimaryKeyColumns), 1, len(catalogPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to delete all from catalog slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Catalog) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Catalog) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Catalog) ReloadG() error {
	if o == nil {
		return errors.New("kmodels: no Catalog provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Catalog) Reload(exec boil.Executor) error {
	ret, err := FindCatalog(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *CatalogSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *CatalogSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CatalogSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("kmodels: empty CatalogSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CatalogSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	catalogs := CatalogSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), catalogPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"catalogs\".* FROM \"catalogs\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, catalogPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(catalogPrimaryKeyColumns), 1, len(catalogPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&catalogs)
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to reload all in CatalogSlice")
	}

	*o = catalogs

	return nil
}

// CatalogExists checks if the Catalog row exists.
func CatalogExists(exec boil.Executor, id int) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"catalogs\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "kmodels: unable to check if catalogs exists")
	}

	return exists, nil
}

// CatalogExistsG checks if the Catalog row exists.
func CatalogExistsG(id int) (bool, error) {
	return CatalogExists(boil.GetDB(), id)
}

// CatalogExistsGP checks if the Catalog row exists. Panics on error.
func CatalogExistsGP(id int) bool {
	e, err := CatalogExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// CatalogExistsP checks if the Catalog row exists. Panics on error.
func CatalogExistsP(exec boil.Executor, id int) bool {
	e, err := CatalogExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
