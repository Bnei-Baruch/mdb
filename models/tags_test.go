package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testTags(t *testing.T) {
	t.Parallel()

	query := Tags(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testTagsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tag := &Tag{}
	if err = randomize.Struct(seed, tag, tagDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tag.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = tag.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Tags(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTagsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tag := &Tag{}
	if err = randomize.Struct(seed, tag, tagDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tag.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Tags(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Tags(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTagsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tag := &Tag{}
	if err = randomize.Struct(seed, tag, tagDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tag.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := TagSlice{tag}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Tags(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testTagsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tag := &Tag{}
	if err = randomize.Struct(seed, tag, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tag.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := TagExists(tx, tag.ID)
	if err != nil {
		t.Errorf("Unable to check if Tag exists: %s", err)
	}
	if !e {
		t.Errorf("Expected TagExistsG to return true, but got false.")
	}
}
func testTagsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tag := &Tag{}
	if err = randomize.Struct(seed, tag, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tag.Insert(tx); err != nil {
		t.Error(err)
	}

	tagFound, err := FindTag(tx, tag.ID)
	if err != nil {
		t.Error(err)
	}

	if tagFound == nil {
		t.Error("want a record, got nil")
	}
}
func testTagsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tag := &Tag{}
	if err = randomize.Struct(seed, tag, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tag.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Tags(tx).Bind(tag); err != nil {
		t.Error(err)
	}
}

func testTagsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tag := &Tag{}
	if err = randomize.Struct(seed, tag, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tag.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Tags(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testTagsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tagOne := &Tag{}
	tagTwo := &Tag{}
	if err = randomize.Struct(seed, tagOne, tagDBTypes, false, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}
	if err = randomize.Struct(seed, tagTwo, tagDBTypes, false, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = tagTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Tags(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testTagsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	tagOne := &Tag{}
	tagTwo := &Tag{}
	if err = randomize.Struct(seed, tagOne, tagDBTypes, false, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}
	if err = randomize.Struct(seed, tagTwo, tagDBTypes, false, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tagOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = tagTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Tags(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testTagsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tag := &Tag{}
	if err = randomize.Struct(seed, tag, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tag.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Tags(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testTagsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tag := &Tag{}
	if err = randomize.Struct(seed, tag, tagDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tag.Insert(tx, tagColumns...); err != nil {
		t.Error(err)
	}

	count, err := Tags(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testTagToManyParentTags(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Tag
	var b, c Tag

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, tagDBTypes, false, tagColumnsWithDefault...)
	randomize.Struct(seed, &c, tagDBTypes, false, tagColumnsWithDefault...)

	b.ParentID.Valid = true
	c.ParentID.Valid = true
	b.ParentID.Int64 = a.ID
	c.ParentID.Int64 = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	tag, err := a.ParentTags(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range tag {
		if v.ParentID.Int64 == b.ParentID.Int64 {
			bFound = true
		}
		if v.ParentID.Int64 == c.ParentID.Int64 {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := TagSlice{&a}
	if err = a.L.LoadParentTags(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ParentTags); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.ParentTags = nil
	if err = a.L.LoadParentTags(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ParentTags); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", tag)
	}
}

func testTagToManyAddOpParentTags(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Tag
	var b, c, d, e Tag

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Tag{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*Tag{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddParentTags(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.ParentID.Int64 {
			t.Error("foreign key was wrong value", a.ID, first.ParentID.Int64)
		}
		if a.ID != second.ParentID.Int64 {
			t.Error("foreign key was wrong value", a.ID, second.ParentID.Int64)
		}

		if first.R.Parent != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Parent != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.ParentTags[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.ParentTags[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.ParentTags(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testTagToManySetOpParentTags(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Tag
	var b, c, d, e Tag

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Tag{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.SetParentTags(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ParentTags(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetParentTags(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ParentTags(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.ParentID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.ParentID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.ID != d.ParentID.Int64 {
		t.Error("foreign key was wrong value", a.ID, d.ParentID.Int64)
	}
	if a.ID != e.ParentID.Int64 {
		t.Error("foreign key was wrong value", a.ID, e.ParentID.Int64)
	}

	if b.R.Parent != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Parent != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Parent != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Parent != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.ParentTags[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.ParentTags[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testTagToManyRemoveOpParentTags(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Tag
	var b, c, d, e Tag

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Tag{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddParentTags(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ParentTags(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveParentTags(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ParentTags(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.ParentID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.ParentID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.Parent != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Parent != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Parent != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.Parent != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.ParentTags) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.ParentTags[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.ParentTags[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testTagToOneStringTranslationUsingLabel(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Tag
	var foreign StringTranslation

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.LabelID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Label(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := TagSlice{&local}
	if err = local.L.LoadLabel(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Label == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Label = nil
	if err = local.L.LoadLabel(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Label == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testTagToOneTagUsingParent(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Tag
	var foreign Tag

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	local.ParentID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.ParentID.Int64 = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Parent(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := TagSlice{&local}
	if err = local.L.LoadParent(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Parent == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Parent = nil
	if err = local.L.LoadParent(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Parent == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testTagToOneSetOpStringTranslationUsingLabel(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Tag
	var b, c StringTranslation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*StringTranslation{&b, &c} {
		err = a.SetLabel(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Label != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.LabelTags[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.LabelID != x.ID {
			t.Error("foreign key was wrong value", a.LabelID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.LabelID))
		reflect.Indirect(reflect.ValueOf(&a.LabelID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.LabelID != x.ID {
			t.Error("foreign key was wrong value", a.LabelID, x.ID)
		}
	}
}
func testTagToOneSetOpTagUsingParent(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Tag
	var b, c Tag

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
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
		err = a.SetParent(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Parent != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.ParentTags[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.ParentID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.ParentID.Int64)
		}

		zero := reflect.Zero(reflect.TypeOf(a.ParentID.Int64))
		reflect.Indirect(reflect.ValueOf(&a.ParentID.Int64)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.ParentID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.ParentID.Int64, x.ID)
		}
	}
}

func testTagToOneRemoveOpTagUsingParent(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Tag
	var b Tag

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.SetParent(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveParent(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.Parent(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.Parent != nil {
		t.Error("R struct entry should be nil")
	}

	if a.ParentID.Valid {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.ParentTags) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testTagsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tag := &Tag{}
	if err = randomize.Struct(seed, tag, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tag.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = tag.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testTagsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tag := &Tag{}
	if err = randomize.Struct(seed, tag, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tag.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := TagSlice{tag}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testTagsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	tag := &Tag{}
	if err = randomize.Struct(seed, tag, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tag.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Tags(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	tagDBTypes = map[string]string{`Description`: `character varying`, `ID`: `bigint`, `LabelID`: `bigint`, `ParentID`: `bigint`}
	_          = bytes.MinRead
)

func testTagsUpdate(t *testing.T) {
	t.Parallel()

	if len(tagColumns) == len(tagPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	tag := &Tag{}
	if err = randomize.Struct(seed, tag, tagDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tag.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Tags(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, tag, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	if err = tag.Update(tx); err != nil {
		t.Error(err)
	}
}

func testTagsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(tagColumns) == len(tagPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	tag := &Tag{}
	if err = randomize.Struct(seed, tag, tagDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tag.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Tags(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, tag, tagDBTypes, true, tagPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(tagColumns, tagPrimaryKeyColumns) {
		fields = tagColumns
	} else {
		fields = strmangle.SetComplement(
			tagColumns,
			tagPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(tag))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := TagSlice{tag}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testTagsUpsert(t *testing.T) {
	t.Parallel()

	if len(tagColumns) == len(tagPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	tag := Tag{}
	if err = randomize.Struct(seed, &tag, tagDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = tag.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Tag: %s", err)
	}

	count, err := Tags(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &tag, tagDBTypes, false, tagPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	if err = tag.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Tag: %s", err)
	}

	count, err = Tags(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
