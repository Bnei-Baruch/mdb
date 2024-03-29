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

func testTagI18ns(t *testing.T) {
	t.Parallel()

	query := TagI18ns()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testTagI18nsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TagI18n{}
	if err = randomize.Struct(seed, o, tagI18nDBTypes, true, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
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

	count, err := TagI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTagI18nsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TagI18n{}
	if err = randomize.Struct(seed, o, tagI18nDBTypes, true, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := TagI18ns().DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := TagI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTagI18nsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TagI18n{}
	if err = randomize.Struct(seed, o, tagI18nDBTypes, true, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := TagI18nSlice{o}

	if rowsAff, err := slice.DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := TagI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTagI18nsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TagI18n{}
	if err = randomize.Struct(seed, o, tagI18nDBTypes, true, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := TagI18nExists(tx, o.TagID, o.Language)
	if err != nil {
		t.Errorf("Unable to check if TagI18n exists: %s", err)
	}
	if !e {
		t.Errorf("Expected TagI18nExists to return true, but got false.")
	}
}

func testTagI18nsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TagI18n{}
	if err = randomize.Struct(seed, o, tagI18nDBTypes, true, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	tagI18nFound, err := FindTagI18n(tx, o.TagID, o.Language)
	if err != nil {
		t.Error(err)
	}

	if tagI18nFound == nil {
		t.Error("want a record, got nil")
	}
}

func testTagI18nsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TagI18n{}
	if err = randomize.Struct(seed, o, tagI18nDBTypes, true, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = TagI18ns().Bind(nil, tx, o); err != nil {
		t.Error(err)
	}
}

func testTagI18nsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TagI18n{}
	if err = randomize.Struct(seed, o, tagI18nDBTypes, true, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := TagI18ns().One(tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testTagI18nsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tagI18nOne := &TagI18n{}
	tagI18nTwo := &TagI18n{}
	if err = randomize.Struct(seed, tagI18nOne, tagI18nDBTypes, false, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}
	if err = randomize.Struct(seed, tagI18nTwo, tagI18nDBTypes, false, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = tagI18nOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = tagI18nTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := TagI18ns().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testTagI18nsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	tagI18nOne := &TagI18n{}
	tagI18nTwo := &TagI18n{}
	if err = randomize.Struct(seed, tagI18nOne, tagI18nDBTypes, false, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}
	if err = randomize.Struct(seed, tagI18nTwo, tagI18nDBTypes, false, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = tagI18nOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = tagI18nTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := TagI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testTagI18nsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TagI18n{}
	if err = randomize.Struct(seed, o, tagI18nDBTypes, true, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := TagI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testTagI18nsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TagI18n{}
	if err = randomize.Struct(seed, o, tagI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Whitelist(tagI18nColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := TagI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testTagI18nToOneTagUsingTag(t *testing.T) {

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var local TagI18n
	var foreign Tag

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, tagI18nDBTypes, false, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, tagDBTypes, false, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	if err := foreign.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.TagID = foreign.ID
	if err := local.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Tag().One(tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := TagI18nSlice{&local}
	if err = local.L.LoadTag(tx, false, (*[]*TagI18n)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Tag == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Tag = nil
	if err = local.L.LoadTag(tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Tag == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testTagI18nToOneUserUsingUser(t *testing.T) {

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var local TagI18n
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, tagI18nDBTypes, true, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err := foreign.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&local.UserID, foreign.ID)
	if err := local.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.User().One(tx)
	if err != nil {
		t.Fatal(err)
	}

	if !queries.Equal(check.ID, foreign.ID) {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := TagI18nSlice{&local}
	if err = local.L.LoadUser(tx, false, (*[]*TagI18n)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.User == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.User = nil
	if err = local.L.LoadUser(tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.User == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testTagI18nToOneSetOpTagUsingTag(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a TagI18n
	var b, c Tag

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagI18nDBTypes, false, strmangle.SetComplement(tagI18nPrimaryKeyColumns, tagI18nColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Tag{&b, &c} {
		err = a.SetTag(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Tag != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.TagI18ns[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.TagID != x.ID {
			t.Error("foreign key was wrong value", a.TagID)
		}

		if exists, err := TagI18nExists(tx, a.TagID, a.Language); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}
func testTagI18nToOneSetOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a TagI18n
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagI18nDBTypes, false, strmangle.SetComplement(tagI18nPrimaryKeyColumns, tagI18nColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*User{&b, &c} {
		err = a.SetUser(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.User != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.TagI18ns[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if !queries.Equal(a.UserID, x.ID) {
			t.Error("foreign key was wrong value", a.UserID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.UserID))
		reflect.Indirect(reflect.ValueOf(&a.UserID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if !queries.Equal(a.UserID, x.ID) {
			t.Error("foreign key was wrong value", a.UserID, x.ID)
		}
	}
}

func testTagI18nToOneRemoveOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a TagI18n
	var b User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagI18nDBTypes, false, strmangle.SetComplement(tagI18nPrimaryKeyColumns, tagI18nColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = a.SetUser(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveUser(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.User().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.User != nil {
		t.Error("R struct entry should be nil")
	}

	if !queries.IsValuerNil(a.UserID) {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.TagI18ns) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testTagI18nsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TagI18n{}
	if err = randomize.Struct(seed, o, tagI18nDBTypes, true, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
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

func testTagI18nsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TagI18n{}
	if err = randomize.Struct(seed, o, tagI18nDBTypes, true, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := TagI18nSlice{o}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}

func testTagI18nsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TagI18n{}
	if err = randomize.Struct(seed, o, tagI18nDBTypes, true, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := TagI18ns().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	tagI18nDBTypes = map[string]string{`TagID`: `bigint`, `Language`: `character`, `OriginalLanguage`: `character`, `Label`: `text`, `UserID`: `bigint`, `CreatedAt`: `timestamp with time zone`}
	_              = bytes.MinRead
)

func testTagI18nsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(tagI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(tagI18nAllColumns) == len(tagI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &TagI18n{}
	if err = randomize.Struct(seed, o, tagI18nDBTypes, true, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := TagI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, tagI18nDBTypes, true, tagI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	if rowsAff, err := o.Update(tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testTagI18nsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(tagI18nAllColumns) == len(tagI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &TagI18n{}
	if err = randomize.Struct(seed, o, tagI18nDBTypes, true, tagI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := TagI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, tagI18nDBTypes, true, tagI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(tagI18nAllColumns, tagI18nPrimaryKeyColumns) {
		fields = tagI18nAllColumns
	} else {
		fields = strmangle.SetComplement(
			tagI18nAllColumns,
			tagI18nPrimaryKeyColumns,
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

	slice := TagI18nSlice{o}
	if rowsAff, err := slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testTagI18nsUpsert(t *testing.T) {
	t.Parallel()

	if len(tagI18nAllColumns) == len(tagI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := TagI18n{}
	if err = randomize.Struct(seed, &o, tagI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert TagI18n: %s", err)
	}

	count, err := TagI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, tagI18nDBTypes, false, tagI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize TagI18n struct: %s", err)
	}

	if err = o.Upsert(tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert TagI18n: %s", err)
	}

	count, err = TagI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
