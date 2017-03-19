package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testCollectionI18ns(t *testing.T) {
	t.Parallel()

	query := CollectionI18ns(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testCollectionI18nsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionI18n := &CollectionI18n{}
	if err = randomize.Struct(seed, collectionI18n, collectionI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = collectionI18n.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := CollectionI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCollectionI18nsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionI18n := &CollectionI18n{}
	if err = randomize.Struct(seed, collectionI18n, collectionI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = CollectionI18ns(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := CollectionI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCollectionI18nsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionI18n := &CollectionI18n{}
	if err = randomize.Struct(seed, collectionI18n, collectionI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := CollectionI18nSlice{collectionI18n}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := CollectionI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testCollectionI18nsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionI18n := &CollectionI18n{}
	if err = randomize.Struct(seed, collectionI18n, collectionI18nDBTypes, true, collectionI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := CollectionI18nExists(tx, collectionI18n.CollectionID, collectionI18n.Language)
	if err != nil {
		t.Errorf("Unable to check if CollectionI18n exists: %s", err)
	}
	if !e {
		t.Errorf("Expected CollectionI18nExistsG to return true, but got false.")
	}
}
func testCollectionI18nsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionI18n := &CollectionI18n{}
	if err = randomize.Struct(seed, collectionI18n, collectionI18nDBTypes, true, collectionI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	collectionI18nFound, err := FindCollectionI18n(tx, collectionI18n.CollectionID, collectionI18n.Language)
	if err != nil {
		t.Error(err)
	}

	if collectionI18nFound == nil {
		t.Error("want a record, got nil")
	}
}
func testCollectionI18nsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionI18n := &CollectionI18n{}
	if err = randomize.Struct(seed, collectionI18n, collectionI18nDBTypes, true, collectionI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = CollectionI18ns(tx).Bind(collectionI18n); err != nil {
		t.Error(err)
	}
}

func testCollectionI18nsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionI18n := &CollectionI18n{}
	if err = randomize.Struct(seed, collectionI18n, collectionI18nDBTypes, true, collectionI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := CollectionI18ns(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testCollectionI18nsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionI18nOne := &CollectionI18n{}
	collectionI18nTwo := &CollectionI18n{}
	if err = randomize.Struct(seed, collectionI18nOne, collectionI18nDBTypes, false, collectionI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}
	if err = randomize.Struct(seed, collectionI18nTwo, collectionI18nDBTypes, false, collectionI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18nOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = collectionI18nTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := CollectionI18ns(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testCollectionI18nsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	collectionI18nOne := &CollectionI18n{}
	collectionI18nTwo := &CollectionI18n{}
	if err = randomize.Struct(seed, collectionI18nOne, collectionI18nDBTypes, false, collectionI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}
	if err = randomize.Struct(seed, collectionI18nTwo, collectionI18nDBTypes, false, collectionI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18nOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = collectionI18nTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := CollectionI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testCollectionI18nsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionI18n := &CollectionI18n{}
	if err = randomize.Struct(seed, collectionI18n, collectionI18nDBTypes, true, collectionI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := CollectionI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testCollectionI18nsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionI18n := &CollectionI18n{}
	if err = randomize.Struct(seed, collectionI18n, collectionI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18n.Insert(tx, collectionI18nColumns...); err != nil {
		t.Error(err)
	}

	count, err := CollectionI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testCollectionI18nToOneCollectionUsingCollection(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local CollectionI18n
	var foreign Collection

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, collectionI18nDBTypes, true, collectionI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, collectionDBTypes, true, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.CollectionID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Collection(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := CollectionI18nSlice{&local}
	if err = local.L.LoadCollection(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Collection == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Collection = nil
	if err = local.L.LoadCollection(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Collection == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testCollectionI18nToOneUserUsingUser(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local CollectionI18n
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, collectionI18nDBTypes, true, collectionI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	local.UserID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.UserID.Int64 = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.User(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := CollectionI18nSlice{&local}
	if err = local.L.LoadUser(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.User == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.User = nil
	if err = local.L.LoadUser(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.User == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testCollectionI18nToOneSetOpCollectionUsingCollection(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a CollectionI18n
	var b, c Collection

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, collectionI18nDBTypes, false, strmangle.SetComplement(collectionI18nPrimaryKeyColumns, collectionI18nColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, collectionDBTypes, false, strmangle.SetComplement(collectionPrimaryKeyColumns, collectionColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, collectionDBTypes, false, strmangle.SetComplement(collectionPrimaryKeyColumns, collectionColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Collection{&b, &c} {
		err = a.SetCollection(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Collection != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.CollectionI18ns[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.CollectionID != x.ID {
			t.Error("foreign key was wrong value", a.CollectionID)
		}

		if exists, err := CollectionI18nExists(tx, a.CollectionID, a.Language); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}
func testCollectionI18nToOneSetOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a CollectionI18n
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, collectionI18nDBTypes, false, strmangle.SetComplement(collectionI18nPrimaryKeyColumns, collectionI18nColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
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

		if x.R.CollectionI18ns[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.UserID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.UserID.Int64)
		}

		zero := reflect.Zero(reflect.TypeOf(a.UserID.Int64))
		reflect.Indirect(reflect.ValueOf(&a.UserID.Int64)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.UserID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.UserID.Int64, x.ID)
		}
	}
}

func testCollectionI18nToOneRemoveOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a CollectionI18n
	var b User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, collectionI18nDBTypes, false, strmangle.SetComplement(collectionI18nPrimaryKeyColumns, collectionI18nColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.SetUser(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveUser(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.User(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.User != nil {
		t.Error("R struct entry should be nil")
	}

	if a.UserID.Valid {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.CollectionI18ns) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testCollectionI18nsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionI18n := &CollectionI18n{}
	if err = randomize.Struct(seed, collectionI18n, collectionI18nDBTypes, true, collectionI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = collectionI18n.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testCollectionI18nsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionI18n := &CollectionI18n{}
	if err = randomize.Struct(seed, collectionI18n, collectionI18nDBTypes, true, collectionI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := CollectionI18nSlice{collectionI18n}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testCollectionI18nsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionI18n := &CollectionI18n{}
	if err = randomize.Struct(seed, collectionI18n, collectionI18nDBTypes, true, collectionI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := CollectionI18ns(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	collectionI18nDBTypes = map[string]string{`CollectionID`: `bigint`, `CreatedAt`: `timestamp with time zone`, `Description`: `text`, `Language`: `character`, `Name`: `text`, `OriginalLanguage`: `character`, `UserID`: `bigint`}
	_                     = bytes.MinRead
)

func testCollectionI18nsUpdate(t *testing.T) {
	t.Parallel()

	if len(collectionI18nColumns) == len(collectionI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	collectionI18n := &CollectionI18n{}
	if err = randomize.Struct(seed, collectionI18n, collectionI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := CollectionI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, collectionI18n, collectionI18nDBTypes, true, collectionI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	if err = collectionI18n.Update(tx); err != nil {
		t.Error(err)
	}
}

func testCollectionI18nsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(collectionI18nColumns) == len(collectionI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	collectionI18n := &CollectionI18n{}
	if err = randomize.Struct(seed, collectionI18n, collectionI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := CollectionI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, collectionI18n, collectionI18nDBTypes, true, collectionI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(collectionI18nColumns, collectionI18nPrimaryKeyColumns) {
		fields = collectionI18nColumns
	} else {
		fields = strmangle.SetComplement(
			collectionI18nColumns,
			collectionI18nPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(collectionI18n))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := CollectionI18nSlice{collectionI18n}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testCollectionI18nsUpsert(t *testing.T) {
	t.Parallel()

	if len(collectionI18nColumns) == len(collectionI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	collectionI18n := CollectionI18n{}
	if err = randomize.Struct(seed, &collectionI18n, collectionI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionI18n.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert CollectionI18n: %s", err)
	}

	count, err := CollectionI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &collectionI18n, collectionI18nDBTypes, false, collectionI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize CollectionI18n struct: %s", err)
	}

	if err = collectionI18n.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert CollectionI18n: %s", err)
	}

	count, err = CollectionI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
