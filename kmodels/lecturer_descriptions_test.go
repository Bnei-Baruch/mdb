package kmodels

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testLecturerDescriptions(t *testing.T) {
	t.Parallel()

	query := LecturerDescriptions(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testLecturerDescriptionsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturerDescription := &LecturerDescription{}
	if err = randomize.Struct(seed, lecturerDescription, lecturerDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = lecturerDescription.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := LecturerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testLecturerDescriptionsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturerDescription := &LecturerDescription{}
	if err = randomize.Struct(seed, lecturerDescription, lecturerDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = LecturerDescriptions(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := LecturerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testLecturerDescriptionsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturerDescription := &LecturerDescription{}
	if err = randomize.Struct(seed, lecturerDescription, lecturerDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := LecturerDescriptionSlice{lecturerDescription}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := LecturerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testLecturerDescriptionsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturerDescription := &LecturerDescription{}
	if err = randomize.Struct(seed, lecturerDescription, lecturerDescriptionDBTypes, true, lecturerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := LecturerDescriptionExists(tx, lecturerDescription.ID)
	if err != nil {
		t.Errorf("Unable to check if LecturerDescription exists: %s", err)
	}
	if !e {
		t.Errorf("Expected LecturerDescriptionExistsG to return true, but got false.")
	}
}
func testLecturerDescriptionsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturerDescription := &LecturerDescription{}
	if err = randomize.Struct(seed, lecturerDescription, lecturerDescriptionDBTypes, true, lecturerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	lecturerDescriptionFound, err := FindLecturerDescription(tx, lecturerDescription.ID)
	if err != nil {
		t.Error(err)
	}

	if lecturerDescriptionFound == nil {
		t.Error("want a record, got nil")
	}
}
func testLecturerDescriptionsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturerDescription := &LecturerDescription{}
	if err = randomize.Struct(seed, lecturerDescription, lecturerDescriptionDBTypes, true, lecturerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = LecturerDescriptions(tx).Bind(lecturerDescription); err != nil {
		t.Error(err)
	}
}

func testLecturerDescriptionsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturerDescription := &LecturerDescription{}
	if err = randomize.Struct(seed, lecturerDescription, lecturerDescriptionDBTypes, true, lecturerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := LecturerDescriptions(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testLecturerDescriptionsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturerDescriptionOne := &LecturerDescription{}
	lecturerDescriptionTwo := &LecturerDescription{}
	if err = randomize.Struct(seed, lecturerDescriptionOne, lecturerDescriptionDBTypes, false, lecturerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}
	if err = randomize.Struct(seed, lecturerDescriptionTwo, lecturerDescriptionDBTypes, false, lecturerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescriptionOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = lecturerDescriptionTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := LecturerDescriptions(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testLecturerDescriptionsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	lecturerDescriptionOne := &LecturerDescription{}
	lecturerDescriptionTwo := &LecturerDescription{}
	if err = randomize.Struct(seed, lecturerDescriptionOne, lecturerDescriptionDBTypes, false, lecturerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}
	if err = randomize.Struct(seed, lecturerDescriptionTwo, lecturerDescriptionDBTypes, false, lecturerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescriptionOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = lecturerDescriptionTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := LecturerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testLecturerDescriptionsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturerDescription := &LecturerDescription{}
	if err = randomize.Struct(seed, lecturerDescription, lecturerDescriptionDBTypes, true, lecturerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := LecturerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testLecturerDescriptionsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturerDescription := &LecturerDescription{}
	if err = randomize.Struct(seed, lecturerDescription, lecturerDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescription.Insert(tx, lecturerDescriptionColumns...); err != nil {
		t.Error(err)
	}

	count, err := LecturerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testLecturerDescriptionsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturerDescription := &LecturerDescription{}
	if err = randomize.Struct(seed, lecturerDescription, lecturerDescriptionDBTypes, true, lecturerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = lecturerDescription.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testLecturerDescriptionsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturerDescription := &LecturerDescription{}
	if err = randomize.Struct(seed, lecturerDescription, lecturerDescriptionDBTypes, true, lecturerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := LecturerDescriptionSlice{lecturerDescription}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testLecturerDescriptionsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturerDescription := &LecturerDescription{}
	if err = randomize.Struct(seed, lecturerDescription, lecturerDescriptionDBTypes, true, lecturerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := LecturerDescriptions(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	lecturerDescriptionDBTypes = map[string]string{`CreatedAt`: `timestamp without time zone`, `Desc`: `character`, `ID`: `integer`, `Lang`: `character`, `LecturerID`: `integer`, `UpdatedAt`: `timestamp without time zone`}
	_                          = bytes.MinRead
)

func testLecturerDescriptionsUpdate(t *testing.T) {
	t.Parallel()

	if len(lecturerDescriptionColumns) == len(lecturerDescriptionPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	lecturerDescription := &LecturerDescription{}
	if err = randomize.Struct(seed, lecturerDescription, lecturerDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := LecturerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, lecturerDescription, lecturerDescriptionDBTypes, true, lecturerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	if err = lecturerDescription.Update(tx); err != nil {
		t.Error(err)
	}
}

func testLecturerDescriptionsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(lecturerDescriptionColumns) == len(lecturerDescriptionPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	lecturerDescription := &LecturerDescription{}
	if err = randomize.Struct(seed, lecturerDescription, lecturerDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := LecturerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, lecturerDescription, lecturerDescriptionDBTypes, true, lecturerDescriptionPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(lecturerDescriptionColumns, lecturerDescriptionPrimaryKeyColumns) {
		fields = lecturerDescriptionColumns
	} else {
		fields = strmangle.SetComplement(
			lecturerDescriptionColumns,
			lecturerDescriptionPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(lecturerDescription))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := LecturerDescriptionSlice{lecturerDescription}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testLecturerDescriptionsUpsert(t *testing.T) {
	t.Parallel()

	if len(lecturerDescriptionColumns) == len(lecturerDescriptionPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	lecturerDescription := LecturerDescription{}
	if err = randomize.Struct(seed, &lecturerDescription, lecturerDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerDescription.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert LecturerDescription: %s", err)
	}

	count, err := LecturerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &lecturerDescription, lecturerDescriptionDBTypes, false, lecturerDescriptionPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize LecturerDescription struct: %s", err)
	}

	if err = lecturerDescription.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert LecturerDescription: %s", err)
	}

	count, err = LecturerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
