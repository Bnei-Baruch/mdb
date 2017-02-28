package kmodels

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testFileTypes(t *testing.T) {
	t.Parallel()

	query := FileTypes(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testFileTypesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileType := &FileType{}
	if err = randomize.Struct(seed, fileType, fileTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileType.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = fileType.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := FileTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testFileTypesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileType := &FileType{}
	if err = randomize.Struct(seed, fileType, fileTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileType.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = FileTypes(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := FileTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testFileTypesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileType := &FileType{}
	if err = randomize.Struct(seed, fileType, fileTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileType.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := FileTypeSlice{fileType}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := FileTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testFileTypesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileType := &FileType{}
	if err = randomize.Struct(seed, fileType, fileTypeDBTypes, true, fileTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileType.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := FileTypeExists(tx, fileType.Name)
	if err != nil {
		t.Errorf("Unable to check if FileType exists: %s", err)
	}
	if !e {
		t.Errorf("Expected FileTypeExistsG to return true, but got false.")
	}
}
func testFileTypesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileType := &FileType{}
	if err = randomize.Struct(seed, fileType, fileTypeDBTypes, true, fileTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileType.Insert(tx); err != nil {
		t.Error(err)
	}

	fileTypeFound, err := FindFileType(tx, fileType.Name)
	if err != nil {
		t.Error(err)
	}

	if fileTypeFound == nil {
		t.Error("want a record, got nil")
	}
}
func testFileTypesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileType := &FileType{}
	if err = randomize.Struct(seed, fileType, fileTypeDBTypes, true, fileTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileType.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = FileTypes(tx).Bind(fileType); err != nil {
		t.Error(err)
	}
}

func testFileTypesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileType := &FileType{}
	if err = randomize.Struct(seed, fileType, fileTypeDBTypes, true, fileTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileType.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := FileTypes(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testFileTypesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileTypeOne := &FileType{}
	fileTypeTwo := &FileType{}
	if err = randomize.Struct(seed, fileTypeOne, fileTypeDBTypes, false, fileTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}
	if err = randomize.Struct(seed, fileTypeTwo, fileTypeDBTypes, false, fileTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileTypeOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = fileTypeTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := FileTypes(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testFileTypesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	fileTypeOne := &FileType{}
	fileTypeTwo := &FileType{}
	if err = randomize.Struct(seed, fileTypeOne, fileTypeDBTypes, false, fileTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}
	if err = randomize.Struct(seed, fileTypeTwo, fileTypeDBTypes, false, fileTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileTypeOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = fileTypeTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := FileTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testFileTypesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileType := &FileType{}
	if err = randomize.Struct(seed, fileType, fileTypeDBTypes, true, fileTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileType.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := FileTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testFileTypesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileType := &FileType{}
	if err = randomize.Struct(seed, fileType, fileTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileType.Insert(tx, fileTypeColumns...); err != nil {
		t.Error(err)
	}

	count, err := FileTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testFileTypesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileType := &FileType{}
	if err = randomize.Struct(seed, fileType, fileTypeDBTypes, true, fileTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileType.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = fileType.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testFileTypesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileType := &FileType{}
	if err = randomize.Struct(seed, fileType, fileTypeDBTypes, true, fileTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileType.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := FileTypeSlice{fileType}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testFileTypesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileType := &FileType{}
	if err = randomize.Struct(seed, fileType, fileTypeDBTypes, true, fileTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileType.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := FileTypes(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	fileTypeDBTypes = map[string]string{`Extlist`: `character varying`, `Name`: `character varying`, `Pic`: `character varying`}
	_               = bytes.MinRead
)

func testFileTypesUpdate(t *testing.T) {
	t.Parallel()

	if len(fileTypeColumns) == len(fileTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	fileType := &FileType{}
	if err = randomize.Struct(seed, fileType, fileTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileType.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := FileTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, fileType, fileTypeDBTypes, true, fileTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	if err = fileType.Update(tx); err != nil {
		t.Error(err)
	}
}

func testFileTypesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(fileTypeColumns) == len(fileTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	fileType := &FileType{}
	if err = randomize.Struct(seed, fileType, fileTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileType.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := FileTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, fileType, fileTypeDBTypes, true, fileTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(fileTypeColumns, fileTypePrimaryKeyColumns) {
		fields = fileTypeColumns
	} else {
		fields = strmangle.SetComplement(
			fileTypeColumns,
			fileTypePrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(fileType))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := FileTypeSlice{fileType}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testFileTypesUpsert(t *testing.T) {
	t.Parallel()

	if len(fileTypeColumns) == len(fileTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	fileType := FileType{}
	if err = randomize.Struct(seed, &fileType, fileTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileType.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert FileType: %s", err)
	}

	count, err := FileTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &fileType, fileTypeDBTypes, false, fileTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize FileType struct: %s", err)
	}

	if err = fileType.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert FileType: %s", err)
	}

	count, err = FileTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
