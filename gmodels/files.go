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

// File is an object representing the database table.
type File struct {
	ID              int64       `boil:"id" json:"id" toml:"id" yaml:"id"`
	UID             string      `boil:"uid" json:"uid" toml:"uid" yaml:"uid"`
	Name            string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	Size            int64       `boil:"size" json:"size" toml:"size" yaml:"size"`
	Type            string      `boil:"type" json:"type" toml:"type" yaml:"type"`
	SubType         string      `boil:"sub_type" json:"sub_type" toml:"sub_type" yaml:"sub_type"`
	MimeType        null.String `boil:"mime_type" json:"mime_type,omitempty" toml:"mime_type" yaml:"mime_type,omitempty"`
	Sha1            null.Bytes  `boil:"sha1" json:"sha1,omitempty" toml:"sha1" yaml:"sha1,omitempty"`
	OperationID     null.Int64  `boil:"operation_id" json:"operation_id,omitempty" toml:"operation_id" yaml:"operation_id,omitempty"`
	ContentUnitID   null.Int64  `boil:"content_unit_id" json:"content_unit_id,omitempty" toml:"content_unit_id" yaml:"content_unit_id,omitempty"`
	CreatedAt       time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	Language        null.String `boil:"language" json:"language,omitempty" toml:"language" yaml:"language,omitempty"`
	BackupCount     null.Int16  `boil:"backup_count" json:"backup_count,omitempty" toml:"backup_count" yaml:"backup_count,omitempty"`
	FirstBackupTime null.Time   `boil:"first_backup_time" json:"first_backup_time,omitempty" toml:"first_backup_time" yaml:"first_backup_time,omitempty"`
	Properties      null.JSON   `boil:"properties" json:"properties,omitempty" toml:"properties" yaml:"properties,omitempty"`
	ParentID        null.Int64  `boil:"parent_id" json:"parent_id,omitempty" toml:"parent_id" yaml:"parent_id,omitempty"`

	R *fileR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L fileL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// fileR is where relationships are stored.
type fileR struct {
	Operation   *Operation
	ContentUnit *ContentUnit
	Parent      *File
	ParentFiles FileSlice
}

// fileL is where Load methods for each relationship are stored.
type fileL struct{}

var (
	fileColumns               = []string{"id", "uid", "name", "size", "type", "sub_type", "mime_type", "sha1", "operation_id", "content_unit_id", "created_at", "language", "backup_count", "first_backup_time", "properties", "parent_id"}
	fileColumnsWithoutDefault = []string{"uid", "name", "size", "type", "sub_type", "mime_type", "sha1", "operation_id", "content_unit_id", "language", "first_backup_time", "properties", "parent_id"}
	fileColumnsWithDefault    = []string{"id", "created_at", "backup_count"}
	filePrimaryKeyColumns     = []string{"id"}
)

type (
	// FileSlice is an alias for a slice of pointers to File.
	// This should generally be used opposed to []File.
	FileSlice []*File

	fileQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	fileType                 = reflect.TypeOf(&File{})
	fileMapping              = queries.MakeStructMapping(fileType)
	filePrimaryKeyMapping, _ = queries.BindMapping(fileType, fileMapping, filePrimaryKeyColumns)
	fileInsertCacheMut       sync.RWMutex
	fileInsertCache          = make(map[string]insertCache)
	fileUpdateCacheMut       sync.RWMutex
	fileUpdateCache          = make(map[string]updateCache)
	fileUpsertCacheMut       sync.RWMutex
	fileUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single file record from the query, and panics on error.
func (q fileQuery) OneP() *File {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single file record from the query.
func (q fileQuery) One() (*File, error) {
	o := &File{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "gmodels: failed to execute a one query for files")
	}

	return o, nil
}

// AllP returns all File records from the query, and panics on error.
func (q fileQuery) AllP() FileSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all File records from the query.
func (q fileQuery) All() (FileSlice, error) {
	var o FileSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "gmodels: failed to assign all query results to File slice")
	}

	return o, nil
}

// CountP returns the count of all File records in the query, and panics on error.
func (q fileQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all File records in the query.
func (q fileQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "gmodels: failed to count files rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q fileQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q fileQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "gmodels: failed to check if files exists")
	}

	return count > 0, nil
}

// OperationG pointed to by the foreign key.
func (o *File) OperationG(mods ...qm.QueryMod) operationQuery {
	return o.Operation(boil.GetDB(), mods...)
}

// Operation pointed to by the foreign key.
func (o *File) Operation(exec boil.Executor, mods ...qm.QueryMod) operationQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.OperationID),
	}

	queryMods = append(queryMods, mods...)

	query := Operations(exec, queryMods...)
	queries.SetFrom(query.Query, "\"operations\"")

	return query
}

// ContentUnitG pointed to by the foreign key.
func (o *File) ContentUnitG(mods ...qm.QueryMod) contentUnitQuery {
	return o.ContentUnit(boil.GetDB(), mods...)
}

// ContentUnit pointed to by the foreign key.
func (o *File) ContentUnit(exec boil.Executor, mods ...qm.QueryMod) contentUnitQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.ContentUnitID),
	}

	queryMods = append(queryMods, mods...)

	query := ContentUnits(exec, queryMods...)
	queries.SetFrom(query.Query, "\"content_units\"")

	return query
}

// ParentG pointed to by the foreign key.
func (o *File) ParentG(mods ...qm.QueryMod) fileQuery {
	return o.Parent(boil.GetDB(), mods...)
}

// Parent pointed to by the foreign key.
func (o *File) Parent(exec boil.Executor, mods ...qm.QueryMod) fileQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.ParentID),
	}

	queryMods = append(queryMods, mods...)

	query := Files(exec, queryMods...)
	queries.SetFrom(query.Query, "\"files\"")

	return query
}

// ParentFilesG retrieves all the file's files via parent_id column.
func (o *File) ParentFilesG(mods ...qm.QueryMod) fileQuery {
	return o.ParentFiles(boil.GetDB(), mods...)
}

// ParentFiles retrieves all the file's files with an executor via parent_id column.
func (o *File) ParentFiles(exec boil.Executor, mods ...qm.QueryMod) fileQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"parent_id\"=?", o.ID),
	)

	query := Files(exec, queryMods...)
	queries.SetFrom(query.Query, "\"files\" as \"a\"")
	return query
}

// LoadOperation allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (fileL) LoadOperation(e boil.Executor, singular bool, maybeFile interface{}) error {
	var slice []*File
	var object *File

	count := 1
	if singular {
		object = maybeFile.(*File)
	} else {
		slice = *maybeFile.(*FileSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &fileR{}
		}
		args[0] = object.OperationID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &fileR{}
			}
			args[i] = obj.OperationID
		}
	}

	query := fmt.Sprintf(
		"select * from \"operations\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Operation")
	}
	defer results.Close()

	var resultSlice []*Operation
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Operation")
	}

	if singular && len(resultSlice) != 0 {
		object.R.Operation = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.OperationID.Int64 == foreign.ID {
				local.R.Operation = foreign
				break
			}
		}
	}

	return nil
}

// LoadContentUnit allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (fileL) LoadContentUnit(e boil.Executor, singular bool, maybeFile interface{}) error {
	var slice []*File
	var object *File

	count := 1
	if singular {
		object = maybeFile.(*File)
	} else {
		slice = *maybeFile.(*FileSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &fileR{}
		}
		args[0] = object.ContentUnitID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &fileR{}
			}
			args[i] = obj.ContentUnitID
		}
	}

	query := fmt.Sprintf(
		"select * from \"content_units\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load ContentUnit")
	}
	defer results.Close()

	var resultSlice []*ContentUnit
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice ContentUnit")
	}

	if singular && len(resultSlice) != 0 {
		object.R.ContentUnit = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ContentUnitID.Int64 == foreign.ID {
				local.R.ContentUnit = foreign
				break
			}
		}
	}

	return nil
}

// LoadParent allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (fileL) LoadParent(e boil.Executor, singular bool, maybeFile interface{}) error {
	var slice []*File
	var object *File

	count := 1
	if singular {
		object = maybeFile.(*File)
	} else {
		slice = *maybeFile.(*FileSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &fileR{}
		}
		args[0] = object.ParentID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &fileR{}
			}
			args[i] = obj.ParentID
		}
	}

	query := fmt.Sprintf(
		"select * from \"files\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load File")
	}
	defer results.Close()

	var resultSlice []*File
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice File")
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

// LoadParentFiles allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (fileL) LoadParentFiles(e boil.Executor, singular bool, maybeFile interface{}) error {
	var slice []*File
	var object *File

	count := 1
	if singular {
		object = maybeFile.(*File)
	} else {
		slice = *maybeFile.(*FileSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &fileR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &fileR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"files\" where \"parent_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load files")
	}
	defer results.Close()

	var resultSlice []*File
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice files")
	}

	if singular {
		object.R.ParentFiles = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.ParentID.Int64 {
				local.R.ParentFiles = append(local.R.ParentFiles, foreign)
				break
			}
		}
	}

	return nil
}

// SetOperation of the file to the related item.
// Sets o.R.Operation to related.
// Adds o to related.R.Files.
func (o *File) SetOperation(exec boil.Executor, insert bool, related *Operation) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"files\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"operation_id"}),
		strmangle.WhereClause("\"", "\"", 2, filePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.OperationID.Int64 = related.ID
	o.OperationID.Valid = true

	if o.R == nil {
		o.R = &fileR{
			Operation: related,
		}
	} else {
		o.R.Operation = related
	}

	if related.R == nil {
		related.R = &operationR{
			Files: FileSlice{o},
		}
	} else {
		related.R.Files = append(related.R.Files, o)
	}

	return nil
}

// RemoveOperation relationship.
// Sets o.R.Operation to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *File) RemoveOperation(exec boil.Executor, related *Operation) error {
	var err error

	o.OperationID.Valid = false
	if err = o.Update(exec, "operation_id"); err != nil {
		o.OperationID.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.Operation = nil
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.Files {
		if o.OperationID.Int64 != ri.OperationID.Int64 {
			continue
		}

		ln := len(related.R.Files)
		if ln > 1 && i < ln-1 {
			related.R.Files[i] = related.R.Files[ln-1]
		}
		related.R.Files = related.R.Files[:ln-1]
		break
	}
	return nil
}

// SetContentUnit of the file to the related item.
// Sets o.R.ContentUnit to related.
// Adds o to related.R.Files.
func (o *File) SetContentUnit(exec boil.Executor, insert bool, related *ContentUnit) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"files\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"content_unit_id"}),
		strmangle.WhereClause("\"", "\"", 2, filePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ContentUnitID.Int64 = related.ID
	o.ContentUnitID.Valid = true

	if o.R == nil {
		o.R = &fileR{
			ContentUnit: related,
		}
	} else {
		o.R.ContentUnit = related
	}

	if related.R == nil {
		related.R = &contentUnitR{
			Files: FileSlice{o},
		}
	} else {
		related.R.Files = append(related.R.Files, o)
	}

	return nil
}

// RemoveContentUnit relationship.
// Sets o.R.ContentUnit to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *File) RemoveContentUnit(exec boil.Executor, related *ContentUnit) error {
	var err error

	o.ContentUnitID.Valid = false
	if err = o.Update(exec, "content_unit_id"); err != nil {
		o.ContentUnitID.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.ContentUnit = nil
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.Files {
		if o.ContentUnitID.Int64 != ri.ContentUnitID.Int64 {
			continue
		}

		ln := len(related.R.Files)
		if ln > 1 && i < ln-1 {
			related.R.Files[i] = related.R.Files[ln-1]
		}
		related.R.Files = related.R.Files[:ln-1]
		break
	}
	return nil
}

// SetParent of the file to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentFiles.
func (o *File) SetParent(exec boil.Executor, insert bool, related *File) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"files\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"parent_id"}),
		strmangle.WhereClause("\"", "\"", 2, filePrimaryKeyColumns),
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
		o.R = &fileR{
			Parent: related,
		}
	} else {
		o.R.Parent = related
	}

	if related.R == nil {
		related.R = &fileR{
			ParentFiles: FileSlice{o},
		}
	} else {
		related.R.ParentFiles = append(related.R.ParentFiles, o)
	}

	return nil
}

// RemoveParent relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *File) RemoveParent(exec boil.Executor, related *File) error {
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

	for i, ri := range related.R.ParentFiles {
		if o.ParentID.Int64 != ri.ParentID.Int64 {
			continue
		}

		ln := len(related.R.ParentFiles)
		if ln > 1 && i < ln-1 {
			related.R.ParentFiles[i] = related.R.ParentFiles[ln-1]
		}
		related.R.ParentFiles = related.R.ParentFiles[:ln-1]
		break
	}
	return nil
}

// AddParentFiles adds the given related objects to the existing relationships
// of the file, optionally inserting them as new records.
// Appends related to o.R.ParentFiles.
// Sets related.R.Parent appropriately.
func (o *File) AddParentFiles(exec boil.Executor, insert bool, related ...*File) error {
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
				"UPDATE \"files\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"parent_id"}),
				strmangle.WhereClause("\"", "\"", 2, filePrimaryKeyColumns),
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
		o.R = &fileR{
			ParentFiles: related,
		}
	} else {
		o.R.ParentFiles = append(o.R.ParentFiles, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &fileR{
				Parent: o,
			}
		} else {
			rel.R.Parent = o
		}
	}
	return nil
}

// SetParentFiles removes all previously related items of the
// file replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentFiles accordingly.
// Replaces o.R.ParentFiles with related.
// Sets related.R.Parent's ParentFiles accordingly.
func (o *File) SetParentFiles(exec boil.Executor, insert bool, related ...*File) error {
	query := "update \"files\" set \"parent_id\" = null where \"parent_id\" = $1"
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
		for _, rel := range o.R.ParentFiles {
			rel.ParentID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Parent = nil
		}

		o.R.ParentFiles = nil
	}
	return o.AddParentFiles(exec, insert, related...)
}

// RemoveParentFiles relationships from objects passed in.
// Removes related items from R.ParentFiles (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
func (o *File) RemoveParentFiles(exec boil.Executor, related ...*File) error {
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
		for i, ri := range o.R.ParentFiles {
			if rel != ri {
				continue
			}

			ln := len(o.R.ParentFiles)
			if ln > 1 && i < ln-1 {
				o.R.ParentFiles[i] = o.R.ParentFiles[ln-1]
			}
			o.R.ParentFiles = o.R.ParentFiles[:ln-1]
			break
		}
	}

	return nil
}

// FilesG retrieves all records.
func FilesG(mods ...qm.QueryMod) fileQuery {
	return Files(boil.GetDB(), mods...)
}

// Files retrieves all the records using an executor.
func Files(exec boil.Executor, mods ...qm.QueryMod) fileQuery {
	mods = append(mods, qm.From("\"files\""))
	return fileQuery{NewQuery(exec, mods...)}
}

// FindFileG retrieves a single record by ID.
func FindFileG(id int64, selectCols ...string) (*File, error) {
	return FindFile(boil.GetDB(), id, selectCols...)
}

// FindFileGP retrieves a single record by ID, and panics on error.
func FindFileGP(id int64, selectCols ...string) *File {
	retobj, err := FindFile(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindFile retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindFile(exec boil.Executor, id int64, selectCols ...string) (*File, error) {
	fileObj := &File{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"files\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(fileObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "gmodels: unable to select from files")
	}

	return fileObj, nil
}

// FindFileP retrieves a single record by ID with an executor, and panics on error.
func FindFileP(exec boil.Executor, id int64, selectCols ...string) *File {
	retobj, err := FindFile(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *File) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *File) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *File) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *File) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("gmodels: no files provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	nzDefaults := queries.NonZeroDefaultSet(fileColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	fileInsertCacheMut.RLock()
	cache, cached := fileInsertCache[key]
	fileInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			fileColumns,
			fileColumnsWithDefault,
			fileColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(fileType, fileMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(fileType, fileMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"files\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "gmodels: unable to insert into files")
	}

	if !cached {
		fileInsertCacheMut.Lock()
		fileInsertCache[key] = cache
		fileInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single File record. See Update for
// whitelist behavior description.
func (o *File) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single File record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *File) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the File, and panics on error.
// See Update for whitelist behavior description.
func (o *File) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the File.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *File) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	fileUpdateCacheMut.RLock()
	cache, cached := fileUpdateCache[key]
	fileUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(fileColumns, filePrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("gmodels: unable to update files, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"files\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, filePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(fileType, fileMapping, append(wl, filePrimaryKeyColumns...))
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
		return errors.Wrap(err, "gmodels: unable to update files row")
	}

	if !cached {
		fileUpdateCacheMut.Lock()
		fileUpdateCache[key] = cache
		fileUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q fileQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q fileQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to update all for files")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o FileSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o FileSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o FileSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o FileSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), filePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"files\" SET %s WHERE (\"id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(filePrimaryKeyColumns), len(colNames)+1, len(filePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to update all in file slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *File) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *File) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *File) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *File) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("gmodels: no files provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	nzDefaults := queries.NonZeroDefaultSet(fileColumnsWithDefault, o)

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

	fileUpsertCacheMut.RLock()
	cache, cached := fileUpsertCache[key]
	fileUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			fileColumns,
			fileColumnsWithDefault,
			fileColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			fileColumns,
			filePrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("gmodels: unable to upsert files, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(filePrimaryKeyColumns))
			copy(conflict, filePrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"files\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(fileType, fileMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(fileType, fileMapping, ret)
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
		return errors.Wrap(err, "gmodels: unable to upsert files")
	}

	if !cached {
		fileUpsertCacheMut.Lock()
		fileUpsertCache[key] = cache
		fileUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single File record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *File) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single File record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *File) DeleteG() error {
	if o == nil {
		return errors.New("gmodels: no File provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single File record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *File) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single File record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *File) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("gmodels: no File provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), filePrimaryKeyMapping)
	sql := "DELETE FROM \"files\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete from files")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q fileQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q fileQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("gmodels: no fileQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete all from files")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o FileSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o FileSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("gmodels: no File slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o FileSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o FileSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("gmodels: no File slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), filePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"files\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, filePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(filePrimaryKeyColumns), 1, len(filePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete all from file slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *File) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *File) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *File) ReloadG() error {
	if o == nil {
		return errors.New("gmodels: no File provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *File) Reload(exec boil.Executor) error {
	ret, err := FindFile(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *FileSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *FileSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *FileSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("gmodels: empty FileSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *FileSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	files := FileSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), filePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"files\".* FROM \"files\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, filePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(filePrimaryKeyColumns), 1, len(filePrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&files)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to reload all in FileSlice")
	}

	*o = files

	return nil
}

// FileExists checks if the File row exists.
func FileExists(exec boil.Executor, id int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"files\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "gmodels: unable to check if files exists")
	}

	return exists, nil
}

// FileExistsG checks if the File row exists.
func FileExistsG(id int64) (bool, error) {
	return FileExists(boil.GetDB(), id)
}

// FileExistsGP checks if the File row exists. Panics on error.
func FileExistsGP(id int64) bool {
	e, err := FileExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// FileExistsP checks if the File row exists. Panics on error.
func FileExistsP(exec boil.Executor, id int64) bool {
	e, err := FileExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
