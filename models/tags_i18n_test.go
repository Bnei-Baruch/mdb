package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testTagsI18ns(t *testing.T) {
	t.Parallel()

	query := TagsI18ns(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testTagsI18nsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tagsI18n := &TagsI18n{}
	if err = randomize.Struct(seed, tagsI18n, tagsI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = tagsI18n.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := TagsI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTagsI18nsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tagsI18n := &TagsI18n{}
	if err = randomize.Struct(seed, tagsI18n, tagsI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = TagsI18ns(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := TagsI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTagsI18nsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tagsI18n := &TagsI18n{}
	if err = randomize.Struct(seed, tagsI18n, tagsI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := TagsI18nSlice{tagsI18n}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := TagsI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testTagsI18nsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tagsI18n := &TagsI18n{}
	if err = randomize.Struct(seed, tagsI18n, tagsI18nDBTypes, true, tagsI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := TagsI18nExists(tx, tagsI18n.TagID, tagsI18n.Language)
	if err != nil {
		t.Errorf("Unable to check if TagsI18n exists: %s", err)
	}
	if !e {
		t.Errorf("Expected TagsI18nExistsG to return true, but got false.")
	}
}
func testTagsI18nsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tagsI18n := &TagsI18n{}
	if err = randomize.Struct(seed, tagsI18n, tagsI18nDBTypes, true, tagsI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	tagsI18nFound, err := FindTagsI18n(tx, tagsI18n.TagID, tagsI18n.Language)
	if err != nil {
		t.Error(err)
	}

	if tagsI18nFound == nil {
		t.Error("want a record, got nil")
	}
}
func testTagsI18nsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tagsI18n := &TagsI18n{}
	if err = randomize.Struct(seed, tagsI18n, tagsI18nDBTypes, true, tagsI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = TagsI18ns(tx).Bind(tagsI18n); err != nil {
		t.Error(err)
	}
}

func testTagsI18nsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tagsI18n := &TagsI18n{}
	if err = randomize.Struct(seed, tagsI18n, tagsI18nDBTypes, true, tagsI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := TagsI18ns(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testTagsI18nsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tagsI18nOne := &TagsI18n{}
	tagsI18nTwo := &TagsI18n{}
	if err = randomize.Struct(seed, tagsI18nOne, tagsI18nDBTypes, false, tagsI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}
	if err = randomize.Struct(seed, tagsI18nTwo, tagsI18nDBTypes, false, tagsI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18nOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = tagsI18nTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := TagsI18ns(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testTagsI18nsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	tagsI18nOne := &TagsI18n{}
	tagsI18nTwo := &TagsI18n{}
	if err = randomize.Struct(seed, tagsI18nOne, tagsI18nDBTypes, false, tagsI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}
	if err = randomize.Struct(seed, tagsI18nTwo, tagsI18nDBTypes, false, tagsI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18nOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = tagsI18nTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := TagsI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testTagsI18nsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tagsI18n := &TagsI18n{}
	if err = randomize.Struct(seed, tagsI18n, tagsI18nDBTypes, true, tagsI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := TagsI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testTagsI18nsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tagsI18n := &TagsI18n{}
	if err = randomize.Struct(seed, tagsI18n, tagsI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18n.Insert(tx, tagsI18nColumns...); err != nil {
		t.Error(err)
	}

	count, err := TagsI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testTagsI18nToOneTagUsingTag(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local TagsI18n
	var foreign Tag

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, tagsI18nDBTypes, true, tagsI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.TagID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Tag(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := TagsI18nSlice{&local}
	if err = local.L.LoadTag(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Tag == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Tag = nil
	if err = local.L.LoadTag(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Tag == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testTagsI18nToOneUserUsingUser(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local TagsI18n
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, tagsI18nDBTypes, true, tagsI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
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

	slice := TagsI18nSlice{&local}
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

func testTagsI18nToOneSetOpTagUsingTag(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a TagsI18n
	var b, c Tag

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagsI18nDBTypes, false, strmangle.SetComplement(tagsI18nPrimaryKeyColumns, tagsI18nColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
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

		if x.R.TagsI18ns[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.TagID != x.ID {
			t.Error("foreign key was wrong value", a.TagID)
		}

		if exists, err := TagsI18nExists(tx, a.TagID, a.Language); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}
func testTagsI18nToOneSetOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a TagsI18n
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagsI18nDBTypes, false, strmangle.SetComplement(tagsI18nPrimaryKeyColumns, tagsI18nColumnsWithoutDefault)...); err != nil {
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

		if x.R.TagsI18ns[0] != &a {
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

func testTagsI18nToOneRemoveOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a TagsI18n
	var b User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagsI18nDBTypes, false, strmangle.SetComplement(tagsI18nPrimaryKeyColumns, tagsI18nColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.TagsI18ns) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testTagsI18nsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tagsI18n := &TagsI18n{}
	if err = randomize.Struct(seed, tagsI18n, tagsI18nDBTypes, true, tagsI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = tagsI18n.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testTagsI18nsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tagsI18n := &TagsI18n{}
	if err = randomize.Struct(seed, tagsI18n, tagsI18nDBTypes, true, tagsI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := TagsI18nSlice{tagsI18n}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testTagsI18nsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tagsI18n := &TagsI18n{}
	if err = randomize.Struct(seed, tagsI18n, tagsI18nDBTypes, true, tagsI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := TagsI18ns(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	tagsI18nDBTypes = map[string]string{`CreatedAt`: `timestamp with time zone`, `Label`: `text`, `Language`: `character`, `OriginalLanguage`: `character`, `TagID`: `bigint`, `UserID`: `bigint`}
	_               = bytes.MinRead
)

func testTagsI18nsUpdate(t *testing.T) {
	t.Parallel()

	if len(tagsI18nColumns) == len(tagsI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	tagsI18n := &TagsI18n{}
	if err = randomize.Struct(seed, tagsI18n, tagsI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := TagsI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, tagsI18n, tagsI18nDBTypes, true, tagsI18nColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	if err = tagsI18n.Update(tx); err != nil {
		t.Error(err)
	}
}

func testTagsI18nsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(tagsI18nColumns) == len(tagsI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	tagsI18n := &TagsI18n{}
	if err = randomize.Struct(seed, tagsI18n, tagsI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18n.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := TagsI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, tagsI18n, tagsI18nDBTypes, true, tagsI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(tagsI18nColumns, tagsI18nPrimaryKeyColumns) {
		fields = tagsI18nColumns
	} else {
		fields = strmangle.SetComplement(
			tagsI18nColumns,
			tagsI18nPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(tagsI18n))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := TagsI18nSlice{tagsI18n}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testTagsI18nsUpsert(t *testing.T) {
	t.Parallel()

	if len(tagsI18nColumns) == len(tagsI18nPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	tagsI18n := TagsI18n{}
	if err = randomize.Struct(seed, &tagsI18n, tagsI18nDBTypes, true); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagsI18n.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert TagsI18n: %s", err)
	}

	count, err := TagsI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &tagsI18n, tagsI18nDBTypes, false, tagsI18nPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize TagsI18n struct: %s", err)
	}

	if err = tagsI18n.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert TagsI18n: %s", err)
	}

	count, err = TagsI18ns(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
