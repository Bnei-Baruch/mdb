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

func testPublisherI18ns(t *testing.T) {
	t.Parallel()

	query := PublisherI18ns()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testPublisherI18nsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PublisherI18n{}
	if err = randomize.Struct(seed, o, publisherI18nDBTypes, true, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
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

	count, err := PublisherI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testPublisherI18nsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PublisherI18n{}
	if err = randomize.Struct(seed, o, publisherI18nDBTypes, true, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := PublisherI18ns().DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := PublisherI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testPublisherI18nsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PublisherI18n{}
	if err = randomize.Struct(seed, o, publisherI18nDBTypes, true, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := PublisherI18nSlice{o}

	if rowsAff, err := slice.DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := PublisherI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testPublisherI18nsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PublisherI18n{}
	if err = randomize.Struct(seed, o, publisherI18nDBTypes, true, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := PublisherI18nExists(tx, o.PublisherID, o.Language)
	if err != nil {
		t.Errorf("Unable to check if PublisherI18n exists: %s", err)
	}
	if !e {
		t.Errorf("Expected PublisherI18nExists to return true, but got false.")
	}
}

func testPublisherI18nsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PublisherI18n{}
	if err = randomize.Struct(seed, o, publisherI18nDBTypes, true, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	publisherI18nFound, err := FindPublisherI18n(tx, o.PublisherID, o.Language)
	if err != nil {
		t.Error(err)
	}

	if publisherI18nFound == nil {
		t.Error("want a record, got nil")
	}
}

func testPublisherI18nsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PublisherI18n{}
	if err = randomize.Struct(seed, o, publisherI18nDBTypes, true, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = PublisherI18ns().Bind(nil, tx, o); err != nil {
		t.Error(err)
	}
}

func testPublisherI18nsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PublisherI18n{}
	if err = randomize.Struct(seed, o, publisherI18nDBTypes, true, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := PublisherI18ns().One(tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testPublisherI18nsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	publisherI18nOne := &PublisherI18n{}
	publisherI18nTwo := &PublisherI18n{}
	if err = randomize.Struct(seed, publisherI18nOne, publisherI18nDBTypes, false, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}
	if err = randomize.Struct(seed, publisherI18nTwo, publisherI18nDBTypes, false, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = publisherI18nOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = publisherI18nTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := PublisherI18ns().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testPublisherI18nsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	publisherI18nOne := &PublisherI18n{}
	publisherI18nTwo := &PublisherI18n{}
	if err = randomize.Struct(seed, publisherI18nOne, publisherI18nDBTypes, false, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}
	if err = randomize.Struct(seed, publisherI18nTwo, publisherI18nDBTypes, false, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = publisherI18nOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = publisherI18nTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := PublisherI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testPublisherI18nsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PublisherI18n{}
	if err = randomize.Struct(seed, o, publisherI18nDBTypes, true, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := PublisherI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testPublisherI18nsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PublisherI18n{}
	if err = randomize.Struct(seed, o, publisherI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Whitelist(publisherI18nColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := PublisherI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testPublisherI18nToOnePublisherUsingPublisher(t *testing.T) {

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var local PublisherI18n
	var foreign Publisher

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, publisherI18nDBTypes, false, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, publisherDBTypes, false, publisherColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Publisher struct: %s", err)
	}

	if err := foreign.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.PublisherID = foreign.ID
	if err := local.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Publisher().One(tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := PublisherI18nSlice{&local}
	if err = local.L.LoadPublisher(tx, false, (*[]*PublisherI18n)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Publisher == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Publisher = nil
	if err = local.L.LoadPublisher(tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Publisher == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testPublisherI18nToOneUserUsingUser(t *testing.T) {

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var local PublisherI18n
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, publisherI18nDBTypes, true, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
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

	slice := PublisherI18nSlice{&local}
	if err = local.L.LoadUser(tx, false, (*[]*PublisherI18n)(&slice), nil); err != nil {
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

func testPublisherI18nToOneSetOpPublisherUsingPublisher(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a PublisherI18n
	var b, c Publisher

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, publisherI18nDBTypes, false, strmangle.SetComplement(publisherI18nPrimaryKeyColumns, publisherI18nColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, publisherDBTypes, false, strmangle.SetComplement(publisherPrimaryKeyColumns, publisherColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, publisherDBTypes, false, strmangle.SetComplement(publisherPrimaryKeyColumns, publisherColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Publisher{&b, &c} {
		err = a.SetPublisher(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Publisher != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.PublisherI18ns[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.PublisherID != x.ID {
			t.Error("foreign key was wrong value", a.PublisherID)
		}

		if exists, err := PublisherI18nExists(tx, a.PublisherID, a.Language); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}
func testPublisherI18nToOneSetOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a PublisherI18n
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, publisherI18nDBTypes, false, strmangle.SetComplement(publisherI18nPrimaryKeyColumns, publisherI18nColumnsWithoutDefault)...); err != nil {
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

		if x.R.PublisherI18ns[0] != &a {
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

func testPublisherI18nToOneRemoveOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a PublisherI18n
	var b User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, publisherI18nDBTypes, false, strmangle.SetComplement(publisherI18nPrimaryKeyColumns, publisherI18nColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.PublisherI18ns) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testPublisherI18nsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PublisherI18n{}
	if err = randomize.Struct(seed, o, publisherI18nDBTypes, true, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
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

func testPublisherI18nsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PublisherI18n{}
	if err = randomize.Struct(seed, o, publisherI18nDBTypes, true, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := PublisherI18nSlice{o}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}

func testPublisherI18nsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PublisherI18n{}
	if err = randomize.Struct(seed, o, publisherI18nDBTypes, true, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := PublisherI18ns().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	publisherI18nDBTypes = map[string]string{`PublisherID`: `bigint`, `Language`: `character`, `OriginalLanguage`: `character`, `Name`: `text`, `Description`: `text`, `UserID`: `bigint`, `CreatedAt`: `timestamp with time zone`}
	_                    = bytes.MinRead
)

func testPublisherI18nsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(publisherI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(publisherI18nAllColumns) == len(publisherI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &PublisherI18n{}
	if err = randomize.Struct(seed, o, publisherI18nDBTypes, true, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := PublisherI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, publisherI18nDBTypes, true, publisherI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	if rowsAff, err := o.Update(tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testPublisherI18nsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(publisherI18nAllColumns) == len(publisherI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &PublisherI18n{}
	if err = randomize.Struct(seed, o, publisherI18nDBTypes, true, publisherI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := PublisherI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, publisherI18nDBTypes, true, publisherI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(publisherI18nAllColumns, publisherI18nPrimaryKeyColumns) {
		fields = publisherI18nAllColumns
	} else {
		fields = strmangle.SetComplement(
			publisherI18nAllColumns,
			publisherI18nPrimaryKeyColumns,
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

	slice := PublisherI18nSlice{o}
	if rowsAff, err := slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testPublisherI18nsUpsert(t *testing.T) {
	t.Parallel()

	if len(publisherI18nAllColumns) == len(publisherI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := PublisherI18n{}
	if err = randomize.Struct(seed, &o, publisherI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert PublisherI18n: %s", err)
	}

	count, err := PublisherI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, publisherI18nDBTypes, false, publisherI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize PublisherI18n struct: %s", err)
	}

	if err = o.Upsert(tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert PublisherI18n: %s", err)
	}

	count, err = PublisherI18ns().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
