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

func testTags(t *testing.T) {
	t.Parallel()

	query := Tags()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testTagsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Tag{}
	if err = randomize.Struct(seed, o, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
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

	count, err := Tags().Count(tx)
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
	o := &Tag{}
	if err = randomize.Struct(seed, o, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Tags().DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Tags().Count(tx)
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
	o := &Tag{}
	if err = randomize.Struct(seed, o, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := TagSlice{o}

	if rowsAff, err := slice.DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Tags().Count(tx)
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
	o := &Tag{}
	if err = randomize.Struct(seed, o, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := TagExists(tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Tag exists: %s", err)
	}
	if !e {
		t.Errorf("Expected TagExists to return true, but got false.")
	}
}

func testTagsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Tag{}
	if err = randomize.Struct(seed, o, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	tagFound, err := FindTag(tx, o.ID)
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
	o := &Tag{}
	if err = randomize.Struct(seed, o, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Tags().Bind(nil, tx, o); err != nil {
		t.Error(err)
	}
}

func testTagsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Tag{}
	if err = randomize.Struct(seed, o, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Tags().One(tx); err != nil {
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
	defer func() { _ = tx.Rollback() }()
	if err = tagOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = tagTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Tags().All(tx)
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
	defer func() { _ = tx.Rollback() }()
	if err = tagOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = tagTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Tags().Count(tx)
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
	o := &Tag{}
	if err = randomize.Struct(seed, o, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Tags().Count(tx)
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
	o := &Tag{}
	if err = randomize.Struct(seed, o, tagDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Whitelist(tagColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Tags().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testTagToManyContentUnits(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a Tag
	var b, c ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, contentUnitDBTypes, false, contentUnitColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, contentUnitDBTypes, false, contentUnitColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec("insert into \"content_units_tags\" (\"tag_id\", \"content_unit_id\") values ($1, $2)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into \"content_units_tags\" (\"tag_id\", \"content_unit_id\") values ($1, $2)", a.ID, c.ID)
	if err != nil {
		t.Fatal(err)
	}

	check, err := a.ContentUnits().All(tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.ID == b.ID {
			bFound = true
		}
		if v.ID == c.ID {
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
	if err = a.L.LoadContentUnits(tx, false, (*[]*Tag)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ContentUnits); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.ContentUnits = nil
	if err = a.L.LoadContentUnits(tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ContentUnits); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testTagToManyTagI18ns(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a Tag
	var b, c TagI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, tagI18nDBTypes, false, tagI18nColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, tagI18nDBTypes, false, tagI18nColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.TagID = a.ID
	c.TagID = a.ID

	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.TagI18ns().All(tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.TagID == b.TagID {
			bFound = true
		}
		if v.TagID == c.TagID {
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
	if err = a.L.LoadTagI18ns(tx, false, (*[]*Tag)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TagI18ns); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.TagI18ns = nil
	if err = a.L.LoadTagI18ns(tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TagI18ns); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testTagToManyParentTags(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a Tag
	var b, c Tag

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, tagDBTypes, false, tagColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, tagDBTypes, false, tagColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&b.ParentID, a.ID)
	queries.Assign(&c.ParentID, a.ID)
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.ParentTags().All(tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if queries.Equal(v.ParentID, b.ParentID) {
			bFound = true
		}
		if queries.Equal(v.ParentID, c.ParentID) {
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
	if err = a.L.LoadParentTags(tx, false, (*[]*Tag)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ParentTags); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.ParentTags = nil
	if err = a.L.LoadParentTags(tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ParentTags); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testTagToManyAddOpContentUnits(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a Tag
	var b, c, d, e ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentUnit{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*ContentUnit{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddContentUnits(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if first.R.Tags[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.Tags[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}

		if a.R.ContentUnits[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.ContentUnits[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.ContentUnits().Count(tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testTagToManySetOpContentUnits(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a Tag
	var b, c, d, e ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentUnit{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err = a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.SetContentUnits(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ContentUnits().Count(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetContentUnits(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ContentUnits().Count(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	// The following checks cannot be implemented since we have no handle
	// to these when we call Set(). Leaving them here as wishful thinking
	// and to let people know there's dragons.
	//
	// if len(b.R.Tags) != 0 {
	// 	t.Error("relationship was not removed properly from the slice")
	// }
	// if len(c.R.Tags) != 0 {
	// 	t.Error("relationship was not removed properly from the slice")
	// }
	if d.R.Tags[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.Tags[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.ContentUnits[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.ContentUnits[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testTagToManyRemoveOpContentUnits(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a Tag
	var b, c, d, e ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentUnit{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.AddContentUnits(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ContentUnits().Count(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveContentUnits(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ContentUnits().Count(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Tags) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Tags) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Tags[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Tags[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if len(a.R.ContentUnits) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.ContentUnits[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.ContentUnits[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testTagToManyAddOpTagI18ns(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a Tag
	var b, c, d, e TagI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*TagI18n{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, tagI18nDBTypes, false, strmangle.SetComplement(tagI18nPrimaryKeyColumns, tagI18nColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*TagI18n{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddTagI18ns(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.TagID {
			t.Error("foreign key was wrong value", a.ID, first.TagID)
		}
		if a.ID != second.TagID {
			t.Error("foreign key was wrong value", a.ID, second.TagID)
		}

		if first.R.Tag != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Tag != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.TagI18ns[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.TagI18ns[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.TagI18ns().Count(tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testTagToManyAddOpParentTags(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

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

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
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

		if !queries.Equal(a.ID, first.ParentID) {
			t.Error("foreign key was wrong value", a.ID, first.ParentID)
		}
		if !queries.Equal(a.ID, second.ParentID) {
			t.Error("foreign key was wrong value", a.ID, second.ParentID)
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

		count, err := a.ParentTags().Count(tx)
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
	defer func() { _ = tx.Rollback() }()

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

	if err = a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.SetParentTags(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ParentTags().Count(tx)
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

	count, err = a.ParentTags().Count(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if !queries.IsValuerNil(b.ParentID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.ParentID) {
		t.Error("want c's foreign key value to be nil")
	}
	if !queries.Equal(a.ID, d.ParentID) {
		t.Error("foreign key was wrong value", a.ID, d.ParentID)
	}
	if !queries.Equal(a.ID, e.ParentID) {
		t.Error("foreign key was wrong value", a.ID, e.ParentID)
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
	defer func() { _ = tx.Rollback() }()

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

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.AddParentTags(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ParentTags().Count(tx)
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

	count, err = a.ParentTags().Count(tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if !queries.IsValuerNil(b.ParentID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.ParentID) {
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

func testTagToOneTagUsingParent(t *testing.T) {

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var local Tag
	var foreign Tag

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, tagDBTypes, false, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	if err := foreign.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&local.ParentID, foreign.ID)
	if err := local.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Parent().One(tx)
	if err != nil {
		t.Fatal(err)
	}

	if !queries.Equal(check.ID, foreign.ID) {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := TagSlice{&local}
	if err = local.L.LoadParent(tx, false, (*[]*Tag)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Parent == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Parent = nil
	if err = local.L.LoadParent(tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Parent == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testTagToOneSetOpTagUsingParent(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

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

	if err := a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx, boil.Infer()); err != nil {
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
		if !queries.Equal(a.ParentID, x.ID) {
			t.Error("foreign key was wrong value", a.ParentID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.ParentID))
		reflect.Indirect(reflect.ValueOf(&a.ParentID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if !queries.Equal(a.ParentID, x.ID) {
			t.Error("foreign key was wrong value", a.ParentID, x.ID)
		}
	}
}

func testTagToOneRemoveOpTagUsingParent(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()

	var a Tag
	var b Tag

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, tagDBTypes, false, strmangle.SetComplement(tagPrimaryKeyColumns, tagColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = a.SetParent(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveParent(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.Parent().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.Parent != nil {
		t.Error("R struct entry should be nil")
	}

	if !queries.IsValuerNil(a.ParentID) {
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
	o := &Tag{}
	if err = randomize.Struct(seed, o, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
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

func testTagsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Tag{}
	if err = randomize.Struct(seed, o, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := TagSlice{o}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}

func testTagsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Tag{}
	if err = randomize.Struct(seed, o, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Tags().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	tagDBTypes = map[string]string{`ID`: `bigint`, `Description`: `character varying`, `ParentID`: `bigint`, `UID`: `character`, `Pattern`: `character varying`}
	_          = bytes.MinRead
)

func testTagsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(tagPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(tagAllColumns) == len(tagPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Tag{}
	if err = randomize.Struct(seed, o, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Tags().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, tagDBTypes, true, tagPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	if rowsAff, err := o.Update(tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testTagsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(tagAllColumns) == len(tagPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Tag{}
	if err = randomize.Struct(seed, o, tagDBTypes, true, tagColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Tags().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, tagDBTypes, true, tagPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(tagAllColumns, tagPrimaryKeyColumns) {
		fields = tagAllColumns
	} else {
		fields = strmangle.SetComplement(
			tagAllColumns,
			tagPrimaryKeyColumns,
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

	slice := TagSlice{o}
	if rowsAff, err := slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testTagsUpsert(t *testing.T) {
	t.Parallel()

	if len(tagAllColumns) == len(tagPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Tag{}
	if err = randomize.Struct(seed, &o, tagDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Tag: %s", err)
	}

	count, err := Tags().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, tagDBTypes, false, tagPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Tag struct: %s", err)
	}

	if err = o.Upsert(tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Tag: %s", err)
	}

	count, err = Tags().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
