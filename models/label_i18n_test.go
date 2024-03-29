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

func testLabelI18ns(t *testing.T) {
	t.Parallel()

	query := LabelI18ns()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testLabelI18nsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LabelI18n{}
	if err = randomize.Struct(seed, o, labelI18nDBTypes, true, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
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

	count, err := LabelI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testLabelI18nsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LabelI18n{}
	if err = randomize.Struct(seed, o, labelI18nDBTypes, true, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := LabelI18ns().DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := LabelI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testLabelI18nsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LabelI18n{}
	if err = randomize.Struct(seed, o, labelI18nDBTypes, true, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := LabelI18nSlice{o}

	if rowsAff, err := slice.DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := LabelI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testLabelI18nsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LabelI18n{}
	if err = randomize.Struct(seed, o, labelI18nDBTypes, true, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := LabelI18nExists(tx, o.LabelID, o.Language)
	if err != nil {
		t.Errorf("Unable to check if LabelI18n exists: %s", err)
	}
	if !e {
		t.Errorf("Expected LabelI18nExists to return true, but got false.")
	}
}

func testLabelI18nsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LabelI18n{}
	if err = randomize.Struct(seed, o, labelI18nDBTypes, true, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	labelI18nFound, err := FindLabelI18n(tx, o.LabelID, o.Language)
	if err != nil {
		t.Error(err)
	}

	if labelI18nFound == nil {
		t.Error("want a record, got nil")
	}
}

func testLabelI18nsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LabelI18n{}
	if err = randomize.Struct(seed, o, labelI18nDBTypes, true, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = LabelI18ns().Bind(nil, tx, o); err != nil {
		t.Error(err)
	}
}

func testLabelI18nsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LabelI18n{}
	if err = randomize.Struct(seed, o, labelI18nDBTypes, true, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := LabelI18ns().One(tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testLabelI18nsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	labelI18nOne := &LabelI18n{}
	labelI18nTwo := &LabelI18n{}
	if err = randomize.Struct(seed, labelI18nOne, labelI18nDBTypes, false, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}
	if err = randomize.Struct(seed, labelI18nTwo, labelI18nDBTypes, false, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = labelI18nOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = labelI18nTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := LabelI18ns().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testLabelI18nsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	labelI18nOne := &LabelI18n{}
	labelI18nTwo := &LabelI18n{}
	if err = randomize.Struct(seed, labelI18nOne, labelI18nDBTypes, false, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}
	if err = randomize.Struct(seed, labelI18nTwo, labelI18nDBTypes, false, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = labelI18nOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = labelI18nTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := LabelI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testLabelI18nsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LabelI18n{}
	if err = randomize.Struct(seed, o, labelI18nDBTypes, true, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := LabelI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testLabelI18nsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LabelI18n{}
	if err = randomize.Struct(seed, o, labelI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Whitelist(labelI18nColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := LabelI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testLabelI18nToOneLabelUsingLabel(t *testing.T) {

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var local LabelI18n
	var foreign Label

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, labelI18nDBTypes, false, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, labelDBTypes, false, labelColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	if err := foreign.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.LabelID = foreign.ID
	if err := local.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Label().One(tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := LabelI18nSlice{&local}
	if err = local.L.LoadLabel(tx, false, (*[]*LabelI18n)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Label == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Label = nil
	if err = local.L.LoadLabel(tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Label == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testLabelI18nToOneUserUsingUser(t *testing.T) {

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var local LabelI18n
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, labelI18nDBTypes, true, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
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

	slice := LabelI18nSlice{&local}
	if err = local.L.LoadUser(tx, false, (*[]*LabelI18n)(&slice), nil); err != nil {
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

func testLabelI18nToOneSetOpLabelUsingLabel(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a LabelI18n
	var b, c Label

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, labelI18nDBTypes, false, strmangle.SetComplement(labelI18nPrimaryKeyColumns, labelI18nColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, labelDBTypes, false, strmangle.SetComplement(labelPrimaryKeyColumns, labelColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, labelDBTypes, false, strmangle.SetComplement(labelPrimaryKeyColumns, labelColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Label{&b, &c} {
		err = a.SetLabel(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Label != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.LabelI18ns[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.LabelID != x.ID {
			t.Error("foreign key was wrong value", a.LabelID)
		}

		if exists, err := LabelI18nExists(tx, a.LabelID, a.Language); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}
func testLabelI18nToOneSetOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a LabelI18n
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, labelI18nDBTypes, false, strmangle.SetComplement(labelI18nPrimaryKeyColumns, labelI18nColumnsWithoutDefault)...); err != nil {
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

		if x.R.LabelI18ns[0] != &a {
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

func testLabelI18nToOneRemoveOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a LabelI18n
	var b User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, labelI18nDBTypes, false, strmangle.SetComplement(labelI18nPrimaryKeyColumns, labelI18nColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.LabelI18ns) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testLabelI18nsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LabelI18n{}
	if err = randomize.Struct(seed, o, labelI18nDBTypes, true, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
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

func testLabelI18nsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LabelI18n{}
	if err = randomize.Struct(seed, o, labelI18nDBTypes, true, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := LabelI18nSlice{o}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}

func testLabelI18nsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &LabelI18n{}
	if err = randomize.Struct(seed, o, labelI18nDBTypes, true, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := LabelI18ns().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	labelI18nDBTypes = map[string]string{`LabelID`: `bigint`, `Language`: `character`, `Name`: `text`, `UserID`: `bigint`, `CreatedAt`: `timestamp with time zone`}
	_                = bytes.MinRead
)

func testLabelI18nsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(labelI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(labelI18nAllColumns) == len(labelI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &LabelI18n{}
	if err = randomize.Struct(seed, o, labelI18nDBTypes, true, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := LabelI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, labelI18nDBTypes, true, labelI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	if rowsAff, err := o.Update(tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testLabelI18nsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(labelI18nAllColumns) == len(labelI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &LabelI18n{}
	if err = randomize.Struct(seed, o, labelI18nDBTypes, true, labelI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := LabelI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, labelI18nDBTypes, true, labelI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(labelI18nAllColumns, labelI18nPrimaryKeyColumns) {
		fields = labelI18nAllColumns
	} else {
		fields = strmangle.SetComplement(
			labelI18nAllColumns,
			labelI18nPrimaryKeyColumns,
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

	slice := LabelI18nSlice{o}
	if rowsAff, err := slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testLabelI18nsUpsert(t *testing.T) {
	t.Parallel()

	if len(labelI18nAllColumns) == len(labelI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := LabelI18n{}
	if err = randomize.Struct(seed, &o, labelI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert LabelI18n: %s", err)
	}

	count, err := LabelI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, labelI18nDBTypes, false, labelI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize LabelI18n struct: %s", err)
	}

	if err = o.Upsert(tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert LabelI18n: %s", err)
	}

	count, err = LabelI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
