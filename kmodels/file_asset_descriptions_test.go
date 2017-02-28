package kmodels

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testFileAssetDescriptions(t *testing.T) {
	t.Parallel()

	query := FileAssetDescriptions(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testFileAssetDescriptionsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAssetDescription := &FileAssetDescription{}
	if err = randomize.Struct(seed, fileAssetDescription, fileAssetDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = fileAssetDescription.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := FileAssetDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testFileAssetDescriptionsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAssetDescription := &FileAssetDescription{}
	if err = randomize.Struct(seed, fileAssetDescription, fileAssetDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = FileAssetDescriptions(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := FileAssetDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testFileAssetDescriptionsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAssetDescription := &FileAssetDescription{}
	if err = randomize.Struct(seed, fileAssetDescription, fileAssetDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := FileAssetDescriptionSlice{fileAssetDescription}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := FileAssetDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testFileAssetDescriptionsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAssetDescription := &FileAssetDescription{}
	if err = randomize.Struct(seed, fileAssetDescription, fileAssetDescriptionDBTypes, true, fileAssetDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := FileAssetDescriptionExists(tx, fileAssetDescription.ID)
	if err != nil {
		t.Errorf("Unable to check if FileAssetDescription exists: %s", err)
	}
	if !e {
		t.Errorf("Expected FileAssetDescriptionExistsG to return true, but got false.")
	}
}
func testFileAssetDescriptionsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAssetDescription := &FileAssetDescription{}
	if err = randomize.Struct(seed, fileAssetDescription, fileAssetDescriptionDBTypes, true, fileAssetDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	fileAssetDescriptionFound, err := FindFileAssetDescription(tx, fileAssetDescription.ID)
	if err != nil {
		t.Error(err)
	}

	if fileAssetDescriptionFound == nil {
		t.Error("want a record, got nil")
	}
}
func testFileAssetDescriptionsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAssetDescription := &FileAssetDescription{}
	if err = randomize.Struct(seed, fileAssetDescription, fileAssetDescriptionDBTypes, true, fileAssetDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = FileAssetDescriptions(tx).Bind(fileAssetDescription); err != nil {
		t.Error(err)
	}
}

func testFileAssetDescriptionsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAssetDescription := &FileAssetDescription{}
	if err = randomize.Struct(seed, fileAssetDescription, fileAssetDescriptionDBTypes, true, fileAssetDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := FileAssetDescriptions(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testFileAssetDescriptionsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAssetDescriptionOne := &FileAssetDescription{}
	fileAssetDescriptionTwo := &FileAssetDescription{}
	if err = randomize.Struct(seed, fileAssetDescriptionOne, fileAssetDescriptionDBTypes, false, fileAssetDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}
	if err = randomize.Struct(seed, fileAssetDescriptionTwo, fileAssetDescriptionDBTypes, false, fileAssetDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescriptionOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = fileAssetDescriptionTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := FileAssetDescriptions(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testFileAssetDescriptionsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	fileAssetDescriptionOne := &FileAssetDescription{}
	fileAssetDescriptionTwo := &FileAssetDescription{}
	if err = randomize.Struct(seed, fileAssetDescriptionOne, fileAssetDescriptionDBTypes, false, fileAssetDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}
	if err = randomize.Struct(seed, fileAssetDescriptionTwo, fileAssetDescriptionDBTypes, false, fileAssetDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescriptionOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = fileAssetDescriptionTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := FileAssetDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testFileAssetDescriptionsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAssetDescription := &FileAssetDescription{}
	if err = randomize.Struct(seed, fileAssetDescription, fileAssetDescriptionDBTypes, true, fileAssetDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := FileAssetDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testFileAssetDescriptionsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAssetDescription := &FileAssetDescription{}
	if err = randomize.Struct(seed, fileAssetDescription, fileAssetDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescription.Insert(tx, fileAssetDescriptionColumns...); err != nil {
		t.Error(err)
	}

	count, err := FileAssetDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testFileAssetDescriptionToOneFileAssetUsingFile(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local FileAssetDescription
	var foreign FileAsset

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, fileAssetDescriptionDBTypes, true, fileAssetDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, fileAssetDBTypes, true, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.FileID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.File(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := FileAssetDescriptionSlice{&local}
	if err = local.L.LoadFile(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.File == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.File = nil
	if err = local.L.LoadFile(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.File == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testFileAssetDescriptionToOneSetOpFileAssetUsingFile(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a FileAssetDescription
	var b, c FileAsset

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileAssetDescriptionDBTypes, false, strmangle.SetComplement(fileAssetDescriptionPrimaryKeyColumns, fileAssetDescriptionColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*FileAsset{&b, &c} {
		err = a.SetFile(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.File != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.FileFileAssetDescriptions[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.FileID != x.ID {
			t.Error("foreign key was wrong value", a.FileID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.FileID))
		reflect.Indirect(reflect.ValueOf(&a.FileID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.FileID != x.ID {
			t.Error("foreign key was wrong value", a.FileID, x.ID)
		}
	}
}
func testFileAssetDescriptionsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAssetDescription := &FileAssetDescription{}
	if err = randomize.Struct(seed, fileAssetDescription, fileAssetDescriptionDBTypes, true, fileAssetDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = fileAssetDescription.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testFileAssetDescriptionsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAssetDescription := &FileAssetDescription{}
	if err = randomize.Struct(seed, fileAssetDescription, fileAssetDescriptionDBTypes, true, fileAssetDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := FileAssetDescriptionSlice{fileAssetDescription}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testFileAssetDescriptionsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAssetDescription := &FileAssetDescription{}
	if err = randomize.Struct(seed, fileAssetDescription, fileAssetDescriptionDBTypes, true, fileAssetDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := FileAssetDescriptions(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	fileAssetDescriptionDBTypes = map[string]string{`CreatedAt`: `timestamp without time zone`, `FileID`: `integer`, `Filedesc`: `character varying`, `ID`: `integer`, `Lang`: `character`, `UpdatedAt`: `timestamp without time zone`}
	_                           = bytes.MinRead
)

func testFileAssetDescriptionsUpdate(t *testing.T) {
	t.Parallel()

	if len(fileAssetDescriptionColumns) == len(fileAssetDescriptionPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	fileAssetDescription := &FileAssetDescription{}
	if err = randomize.Struct(seed, fileAssetDescription, fileAssetDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := FileAssetDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, fileAssetDescription, fileAssetDescriptionDBTypes, true, fileAssetDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	if err = fileAssetDescription.Update(tx); err != nil {
		t.Error(err)
	}
}

func testFileAssetDescriptionsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(fileAssetDescriptionColumns) == len(fileAssetDescriptionPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	fileAssetDescription := &FileAssetDescription{}
	if err = randomize.Struct(seed, fileAssetDescription, fileAssetDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := FileAssetDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, fileAssetDescription, fileAssetDescriptionDBTypes, true, fileAssetDescriptionPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(fileAssetDescriptionColumns, fileAssetDescriptionPrimaryKeyColumns) {
		fields = fileAssetDescriptionColumns
	} else {
		fields = strmangle.SetComplement(
			fileAssetDescriptionColumns,
			fileAssetDescriptionPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(fileAssetDescription))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := FileAssetDescriptionSlice{fileAssetDescription}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testFileAssetDescriptionsUpsert(t *testing.T) {
	t.Parallel()

	if len(fileAssetDescriptionColumns) == len(fileAssetDescriptionPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	fileAssetDescription := FileAssetDescription{}
	if err = randomize.Struct(seed, &fileAssetDescription, fileAssetDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetDescription.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert FileAssetDescription: %s", err)
	}

	count, err := FileAssetDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &fileAssetDescription, fileAssetDescriptionDBTypes, false, fileAssetDescriptionPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize FileAssetDescription struct: %s", err)
	}

	if err = fileAssetDescription.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert FileAssetDescription: %s", err)
	}

	count, err = FileAssetDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
