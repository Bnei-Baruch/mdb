package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testSourceI18ns(t *testing.T) {
	t.Parallel()

	query := SourceI18ns(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testSourceI18nsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceI18n := &SourceI18n{}
	if err = randomize.Struct(seed, sourceI18n, sourceI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = sourceI18n.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := SourceI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testSourceI18nsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceI18n := &SourceI18n{}
	if err = randomize.Struct(seed, sourceI18n, sourceI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = SourceI18ns(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := SourceI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testSourceI18nsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceI18n := &SourceI18n{}
	if err = randomize.Struct(seed, sourceI18n, sourceI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := SourceI18nSlice{sourceI18n}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := SourceI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testSourceI18nsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceI18n := &SourceI18n{}
	if err = randomize.Struct(seed, sourceI18n, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := SourceI18nExists(tx, sourceI18n.SourceID, sourceI18n.Language)
	if err != nil {
		t.Errorf("Unable to check if SourceI18n exists: %s", err)
	}
	if !e {
		t.Errorf("Expected SourceI18nExistsG to return true, but got false.")
	}
}
func testSourceI18nsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceI18n := &SourceI18n{}
	if err = randomize.Struct(seed, sourceI18n, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	sourceI18nFound, err := FindSourceI18n(tx, sourceI18n.SourceID, sourceI18n.Language)
	if err != nil {
		t.Error(err)
	}

	if sourceI18nFound == nil {
		t.Error("want a record, got nil")
	}
}
func testSourceI18nsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceI18n := &SourceI18n{}
	if err = randomize.Struct(seed, sourceI18n, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = SourceI18ns(tx).Bind(sourceI18n); err != nil {
		t.Error(err)
	}
}

func testSourceI18nsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceI18n := &SourceI18n{}
	if err = randomize.Struct(seed, sourceI18n, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := SourceI18ns(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testSourceI18nsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceI18nOne := &SourceI18n{}
	sourceI18nTwo := &SourceI18n{}
	if err = randomize.Struct(seed, sourceI18nOne, sourceI18nDBTypes, false, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}
	if err = randomize.Struct(seed, sourceI18nTwo, sourceI18nDBTypes, false, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18nOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = sourceI18nTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := SourceI18ns(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testSourceI18nsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	sourceI18nOne := &SourceI18n{}
	sourceI18nTwo := &SourceI18n{}
	if err = randomize.Struct(seed, sourceI18nOne, sourceI18nDBTypes, false, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}
	if err = randomize.Struct(seed, sourceI18nTwo, sourceI18nDBTypes, false, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18nOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = sourceI18nTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := SourceI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testSourceI18nsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceI18n := &SourceI18n{}
	if err = randomize.Struct(seed, sourceI18n, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := SourceI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testSourceI18nsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceI18n := &SourceI18n{}
	if err = randomize.Struct(seed, sourceI18n, sourceI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18n.Insert(tx, sourceI18nColumns...); err != nil {
		t.Error(err)
	}

	count, err := SourceI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testSourceI18nToOneSourceUsingSource(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local SourceI18n
	var foreign Source

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.SourceID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Source(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := SourceI18nSlice{&local}
	if err = local.L.LoadSource(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Source == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Source = nil
	if err = local.L.LoadSource(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Source == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testSourceI18nToOneSetOpSourceUsingSource(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a SourceI18n
	var b, c Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceI18nDBTypes, false, strmangle.SetComplement(sourceI18nPrimaryKeyColumns, sourceI18nColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Source{&b, &c} {
		err = a.SetSource(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Source != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.SourceI18ns[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.SourceID != x.ID {
			t.Error("foreign key was wrong value", a.SourceID)
		}

		if exists, err := SourceI18nExists(tx, a.SourceID, a.Language); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}
func testSourceI18nsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceI18n := &SourceI18n{}
	if err = randomize.Struct(seed, sourceI18n, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = sourceI18n.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testSourceI18nsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceI18n := &SourceI18n{}
	if err = randomize.Struct(seed, sourceI18n, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := SourceI18nSlice{sourceI18n}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testSourceI18nsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceI18n := &SourceI18n{}
	if err = randomize.Struct(seed, sourceI18n, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := SourceI18ns(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	sourceI18nDBTypes = map[string]string{`CreatedAt`: `timestamp with time zone`, `Description`: `text`, `Language`: `character`, `Name`: `character varying`, `SourceID`: `bigint`}
	_                 = bytes.MinRead
)

func testSourceI18nsUpdate(t *testing.T) {
	t.Parallel()

	if len(sourceI18nColumns) == len(sourceI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	sourceI18n := &SourceI18n{}
	if err = randomize.Struct(seed, sourceI18n, sourceI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := SourceI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, sourceI18n, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	if err = sourceI18n.Update(tx); err != nil {
		t.Error(err)
	}
}

func testSourceI18nsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(sourceI18nColumns) == len(sourceI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	sourceI18n := &SourceI18n{}
	if err = randomize.Struct(seed, sourceI18n, sourceI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := SourceI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, sourceI18n, sourceI18nDBTypes, true, sourceI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(sourceI18nColumns, sourceI18nPrimaryKeyColumns) {
		fields = sourceI18nColumns
	} else {
		fields = strmangle.SetComplement(
			sourceI18nColumns,
			sourceI18nPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(sourceI18n))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := SourceI18nSlice{sourceI18n}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testSourceI18nsUpsert(t *testing.T) {
	t.Parallel()

	if len(sourceI18nColumns) == len(sourceI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	sourceI18n := SourceI18n{}
	if err = randomize.Struct(seed, &sourceI18n, sourceI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceI18n.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert SourceI18n: %s", err)
	}

	count, err := SourceI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &sourceI18n, sourceI18nDBTypes, false, sourceI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	if err = sourceI18n.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert SourceI18n: %s", err)
	}

	count, err = SourceI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
