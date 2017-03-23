package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testAuthorI18ns(t *testing.T) {
	t.Parallel()

	query := AuthorI18ns(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testAuthorI18nsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorI18n := &AuthorI18n{}
	if err = randomize.Struct(seed, authorI18n, authorI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = authorI18n.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := AuthorI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuthorI18nsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorI18n := &AuthorI18n{}
	if err = randomize.Struct(seed, authorI18n, authorI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = AuthorI18ns(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := AuthorI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuthorI18nsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorI18n := &AuthorI18n{}
	if err = randomize.Struct(seed, authorI18n, authorI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := AuthorI18nSlice{authorI18n}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := AuthorI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testAuthorI18nsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorI18n := &AuthorI18n{}
	if err = randomize.Struct(seed, authorI18n, authorI18nDBTypes, true, authorI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := AuthorI18nExists(tx, authorI18n.AuthorID, authorI18n.Language)
	if err != nil {
		t.Errorf("Unable to check if AuthorI18n exists: %s", err)
	}
	if !e {
		t.Errorf("Expected AuthorI18nExistsG to return true, but got false.")
	}
}
func testAuthorI18nsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorI18n := &AuthorI18n{}
	if err = randomize.Struct(seed, authorI18n, authorI18nDBTypes, true, authorI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	authorI18nFound, err := FindAuthorI18n(tx, authorI18n.AuthorID, authorI18n.Language)
	if err != nil {
		t.Error(err)
	}

	if authorI18nFound == nil {
		t.Error("want a record, got nil")
	}
}
func testAuthorI18nsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorI18n := &AuthorI18n{}
	if err = randomize.Struct(seed, authorI18n, authorI18nDBTypes, true, authorI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = AuthorI18ns(tx).Bind(authorI18n); err != nil {
		t.Error(err)
	}
}

func testAuthorI18nsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorI18n := &AuthorI18n{}
	if err = randomize.Struct(seed, authorI18n, authorI18nDBTypes, true, authorI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := AuthorI18ns(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testAuthorI18nsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorI18nOne := &AuthorI18n{}
	authorI18nTwo := &AuthorI18n{}
	if err = randomize.Struct(seed, authorI18nOne, authorI18nDBTypes, false, authorI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}
	if err = randomize.Struct(seed, authorI18nTwo, authorI18nDBTypes, false, authorI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18nOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = authorI18nTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := AuthorI18ns(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testAuthorI18nsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	authorI18nOne := &AuthorI18n{}
	authorI18nTwo := &AuthorI18n{}
	if err = randomize.Struct(seed, authorI18nOne, authorI18nDBTypes, false, authorI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}
	if err = randomize.Struct(seed, authorI18nTwo, authorI18nDBTypes, false, authorI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18nOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = authorI18nTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := AuthorI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testAuthorI18nsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorI18n := &AuthorI18n{}
	if err = randomize.Struct(seed, authorI18n, authorI18nDBTypes, true, authorI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := AuthorI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAuthorI18nsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorI18n := &AuthorI18n{}
	if err = randomize.Struct(seed, authorI18n, authorI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18n.Insert(tx, authorI18nColumns...); err != nil {
		t.Error(err)
	}

	count, err := AuthorI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAuthorI18nToOneAuthorUsingAuthor(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local AuthorI18n
	var foreign Author

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, authorI18nDBTypes, true, authorI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.AuthorID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Author(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := AuthorI18nSlice{&local}
	if err = local.L.LoadAuthor(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Author == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Author = nil
	if err = local.L.LoadAuthor(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Author == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testAuthorI18nToOneSetOpAuthorUsingAuthor(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a AuthorI18n
	var b, c Author

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorI18nDBTypes, false, strmangle.SetComplement(authorI18nPrimaryKeyColumns, authorI18nColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, authorDBTypes, false, strmangle.SetComplement(authorPrimaryKeyColumns, authorColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, authorDBTypes, false, strmangle.SetComplement(authorPrimaryKeyColumns, authorColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Author{&b, &c} {
		err = a.SetAuthor(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Author != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.AuthorI18ns[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.AuthorID != x.ID {
			t.Error("foreign key was wrong value", a.AuthorID)
		}

		if exists, err := AuthorI18nExists(tx, a.AuthorID, a.Language); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}
func testAuthorI18nsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorI18n := &AuthorI18n{}
	if err = randomize.Struct(seed, authorI18n, authorI18nDBTypes, true, authorI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = authorI18n.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testAuthorI18nsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorI18n := &AuthorI18n{}
	if err = randomize.Struct(seed, authorI18n, authorI18nDBTypes, true, authorI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := AuthorI18nSlice{authorI18n}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testAuthorI18nsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorI18n := &AuthorI18n{}
	if err = randomize.Struct(seed, authorI18n, authorI18nDBTypes, true, authorI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := AuthorI18ns(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	authorI18nDBTypes = map[string]string{`AuthorID`: `bigint`, `CreatedAt`: `timestamp with time zone`, `FullName`: `character varying`, `Language`: `character`, `Name`: `character varying`}
	_                 = bytes.MinRead
)

func testAuthorI18nsUpdate(t *testing.T) {
	t.Parallel()

	if len(authorI18nColumns) == len(authorI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	authorI18n := &AuthorI18n{}
	if err = randomize.Struct(seed, authorI18n, authorI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := AuthorI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, authorI18n, authorI18nDBTypes, true, authorI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	if err = authorI18n.Update(tx); err != nil {
		t.Error(err)
	}
}

func testAuthorI18nsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(authorI18nColumns) == len(authorI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	authorI18n := &AuthorI18n{}
	if err = randomize.Struct(seed, authorI18n, authorI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := AuthorI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, authorI18n, authorI18nDBTypes, true, authorI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(authorI18nColumns, authorI18nPrimaryKeyColumns) {
		fields = authorI18nColumns
	} else {
		fields = strmangle.SetComplement(
			authorI18nColumns,
			authorI18nPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(authorI18n))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := AuthorI18nSlice{authorI18n}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testAuthorI18nsUpsert(t *testing.T) {
	t.Parallel()

	if len(authorI18nColumns) == len(authorI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	authorI18n := AuthorI18n{}
	if err = randomize.Struct(seed, &authorI18n, authorI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorI18n.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert AuthorI18n: %s", err)
	}

	count, err := AuthorI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &authorI18n, authorI18nDBTypes, false, authorI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize AuthorI18n struct: %s", err)
	}

	if err = authorI18n.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert AuthorI18n: %s", err)
	}

	count, err = AuthorI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
