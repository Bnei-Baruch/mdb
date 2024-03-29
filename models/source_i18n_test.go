// Code generated by SQLBoiler 4.8.6 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testSourceI18ns(t *testing.T) {
	t.Parallel()

	query := SourceI18ns()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testSourceI18nsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SourceI18n{}
	if err = randomize.Struct(seed, o, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := SourceI18ns().Count(tx)
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
	o := &SourceI18n{}
	if err = randomize.Struct(seed, o, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := SourceI18ns().DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := SourceI18ns().Count(tx)
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
	o := &SourceI18n{}
	if err = randomize.Struct(seed, o, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := SourceI18nSlice{o}

	if rowsAff, err := slice.DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := SourceI18ns().Count(tx)
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
	o := &SourceI18n{}
	if err = randomize.Struct(seed, o, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := SourceI18nExists(tx, o.SourceID, o.Language)
	if err != nil {
		t.Errorf("Unable to check if SourceI18n exists: %s", err)
	}
	if !e {
		t.Errorf("Expected SourceI18nExists to return true, but got false.")
	}
}

func testSourceI18nsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SourceI18n{}
	if err = randomize.Struct(seed, o, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	sourceI18nFound, err := FindSourceI18n(tx, o.SourceID, o.Language)
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
	o := &SourceI18n{}
	if err = randomize.Struct(seed, o, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = SourceI18ns().Bind(nil, tx, o); err != nil {
		t.Error(err)
	}
}

func testSourceI18nsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SourceI18n{}
	if err = randomize.Struct(seed, o, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := SourceI18ns().One(tx); err != nil {
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
	defer func() { _ = tx.Rollback() }()
	if err = sourceI18nOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = sourceI18nTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := SourceI18ns().All(tx)
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
	defer func() { _ = tx.Rollback() }()
	if err = sourceI18nOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = sourceI18nTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := SourceI18ns().Count(tx)
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
	o := &SourceI18n{}
	if err = randomize.Struct(seed, o, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := SourceI18ns().Count(tx)
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
	o := &SourceI18n{}
	if err = randomize.Struct(seed, o, sourceI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Whitelist(sourceI18nColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := SourceI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testSourceI18nToOneSourceUsingSource(t *testing.T) {

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var local SourceI18n
	var foreign Source

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, sourceI18nDBTypes, false, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, sourceDBTypes, false, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	if err := foreign.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.SourceID = foreign.ID
	if err := local.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Source().One(tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := SourceI18nSlice{&local}
	if err = local.L.LoadSource(tx, false, (*[]*SourceI18n)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Source == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Source = nil
	if err = local.L.LoadSource(tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Source == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testSourceI18nToOneSetOpSourceUsingSource(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

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

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
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
	o := &SourceI18n{}
	if err = randomize.Struct(seed, o, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testSourceI18nsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SourceI18n{}
	if err = randomize.Struct(seed, o, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := SourceI18nSlice{o}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}

func testSourceI18nsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &SourceI18n{}
	if err = randomize.Struct(seed, o, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := SourceI18ns().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	sourceI18nDBTypes = map[string]string{`SourceID`: `bigint`, `Language`: `character`, `Name`: `character varying`, `Description`: `text`, `CreatedAt`: `timestamp with time zone`}
	_                 = bytes.MinRead
)

func testSourceI18nsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(sourceI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(sourceI18nAllColumns) == len(sourceI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &SourceI18n{}
	if err = randomize.Struct(seed, o, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := SourceI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, sourceI18nDBTypes, true, sourceI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	if rowsAff, err := o.Update(tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testSourceI18nsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(sourceI18nAllColumns) == len(sourceI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &SourceI18n{}
	if err = randomize.Struct(seed, o, sourceI18nDBTypes, true, sourceI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := SourceI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, sourceI18nDBTypes, true, sourceI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(sourceI18nAllColumns, sourceI18nPrimaryKeyColumns) {
		fields = sourceI18nAllColumns
	} else {
		fields = strmangle.SetComplement(
			sourceI18nAllColumns,
			sourceI18nPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := SourceI18nSlice{o}
	if rowsAff, err := slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testSourceI18nsUpsert(t *testing.T) {
	t.Parallel()

	if len(sourceI18nAllColumns) == len(sourceI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := SourceI18n{}
	if err = randomize.Struct(seed, &o, sourceI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert SourceI18n: %s", err)
	}

	count, err := SourceI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, sourceI18nDBTypes, false, sourceI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize SourceI18n struct: %s", err)
	}

	if err = o.Upsert(tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert SourceI18n: %s", err)
	}

	count, err = SourceI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
