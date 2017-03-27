package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testSources(t *testing.T) {
	t.Parallel()

	query := Sources(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testSourcesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	source := &Source{}
	if err = randomize.Struct(seed, source, sourceDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = source.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = source.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Sources(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testSourcesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	source := &Source{}
	if err = randomize.Struct(seed, source, sourceDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = source.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Sources(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Sources(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testSourcesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	source := &Source{}
	if err = randomize.Struct(seed, source, sourceDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = source.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := SourceSlice{source}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Sources(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testSourcesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	source := &Source{}
	if err = randomize.Struct(seed, source, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = source.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := SourceExists(tx, source.ID)
	if err != nil {
		t.Errorf("Unable to check if Source exists: %s", err)
	}
	if !e {
		t.Errorf("Expected SourceExistsG to return true, but got false.")
	}
}
func testSourcesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	source := &Source{}
	if err = randomize.Struct(seed, source, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = source.Insert(tx); err != nil {
		t.Error(err)
	}

	sourceFound, err := FindSource(tx, source.ID)
	if err != nil {
		t.Error(err)
	}

	if sourceFound == nil {
		t.Error("want a record, got nil")
	}
}
func testSourcesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	source := &Source{}
	if err = randomize.Struct(seed, source, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = source.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Sources(tx).Bind(source); err != nil {
		t.Error(err)
	}
}

func testSourcesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	source := &Source{}
	if err = randomize.Struct(seed, source, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = source.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Sources(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testSourcesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceOne := &Source{}
	sourceTwo := &Source{}
	if err = randomize.Struct(seed, sourceOne, sourceDBTypes, false, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}
	if err = randomize.Struct(seed, sourceTwo, sourceDBTypes, false, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = sourceTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Sources(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testSourcesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	sourceOne := &Source{}
	sourceTwo := &Source{}
	if err = randomize.Struct(seed, sourceOne, sourceDBTypes, false, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}
	if err = randomize.Struct(seed, sourceTwo, sourceDBTypes, false, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = sourceTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Sources(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testSourcesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	source := &Source{}
	if err = randomize.Struct(seed, source, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = source.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Sources(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testSourcesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	source := &Source{}
	if err = randomize.Struct(seed, source, sourceDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = source.Insert(tx, sourceColumns...); err != nil {
		t.Error(err)
	}

	count, err := Sources(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testSourceToManyParentSources(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b, c Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, sourceDBTypes, false, sourceColumnsWithDefault...)
	randomize.Struct(seed, &c, sourceDBTypes, false, sourceColumnsWithDefault...)

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

	source, err := a.ParentSources(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range source {
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

	slice := SourceSlice{&a}
	if err = a.L.LoadParentSources(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ParentSources); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.ParentSources = nil
	if err = a.L.LoadParentSources(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ParentSources); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", source)
	}
}

func testSourceToManySourceI18ns(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b, c SourceI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, sourceI18nDBTypes, false, sourceI18nColumnsWithDefault...)
	randomize.Struct(seed, &c, sourceI18nDBTypes, false, sourceI18nColumnsWithDefault...)

	b.SourceID = a.ID
	c.SourceID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	sourceI18n, err := a.SourceI18ns(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range sourceI18n {
		if v.SourceID == b.SourceID {
			bFound = true
		}
		if v.SourceID == c.SourceID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := SourceSlice{&a}
	if err = a.L.LoadSourceI18ns(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.SourceI18ns); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.SourceI18ns = nil
	if err = a.L.LoadSourceI18ns(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.SourceI18ns); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", sourceI18n)
	}
}

func testSourceToManyContentUnits(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b, c ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, contentUnitDBTypes, false, contentUnitColumnsWithDefault...)
	randomize.Struct(seed, &c, contentUnitDBTypes, false, contentUnitColumnsWithDefault...)

	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec("insert into \"content_units_sources\" (\"source_id\", \"content_unit_id\") values ($1, $2)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into \"content_units_sources\" (\"source_id\", \"content_unit_id\") values ($1, $2)", a.ID, c.ID)
	if err != nil {
		t.Fatal(err)
	}

	contentUnit, err := a.ContentUnits(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range contentUnit {
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

	slice := SourceSlice{&a}
	if err = a.L.LoadContentUnits(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ContentUnits); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.ContentUnits = nil
	if err = a.L.LoadContentUnits(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ContentUnits); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", contentUnit)
	}
}

func testSourceToManyAuthors(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b, c Author

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, authorDBTypes, false, authorColumnsWithDefault...)
	randomize.Struct(seed, &c, authorDBTypes, false, authorColumnsWithDefault...)

	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec("insert into \"authors_sources\" (\"source_id\", \"author_id\") values ($1, $2)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into \"authors_sources\" (\"source_id\", \"author_id\") values ($1, $2)", a.ID, c.ID)
	if err != nil {
		t.Fatal(err)
	}

	author, err := a.Authors(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range author {
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

	slice := SourceSlice{&a}
	if err = a.L.LoadAuthors(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Authors); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Authors = nil
	if err = a.L.LoadAuthors(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Authors); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", author)
	}
}

func testSourceToManyAddOpParentSources(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b, c, d, e Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Source{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Source{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddParentSources(tx, i != 0, x...)
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

		if a.R.ParentSources[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.ParentSources[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.ParentSources(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testSourceToManySetOpParentSources(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b, c, d, e Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Source{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
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

	err = a.SetParentSources(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ParentSources(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetParentSources(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ParentSources(tx).Count()
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

	if a.R.ParentSources[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.ParentSources[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testSourceToManyRemoveOpParentSources(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b, c, d, e Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Source{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddParentSources(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ParentSources(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveParentSources(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ParentSources(tx).Count()
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

	if len(a.R.ParentSources) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.ParentSources[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.ParentSources[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testSourceToManyAddOpSourceI18ns(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b, c, d, e SourceI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*SourceI18n{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, sourceI18nDBTypes, false, strmangle.SetComplement(sourceI18nPrimaryKeyColumns, sourceI18nColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*SourceI18n{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddSourceI18ns(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.SourceID {
			t.Error("foreign key was wrong value", a.ID, first.SourceID)
		}
		if a.ID != second.SourceID {
			t.Error("foreign key was wrong value", a.ID, second.SourceID)
		}

		if first.R.Source != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Source != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.SourceI18ns[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.SourceI18ns[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.SourceI18ns(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testSourceToManyAddOpContentUnits(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b, c, d, e ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentUnit{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
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

		if first.R.Sources[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.Sources[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}

		if a.R.ContentUnits[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.ContentUnits[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.ContentUnits(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testSourceToManySetOpContentUnits(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b, c, d, e ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentUnit{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
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

	err = a.SetContentUnits(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ContentUnits(tx).Count()
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

	count, err = a.ContentUnits(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Sources) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Sources) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Sources[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.Sources[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.ContentUnits[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.ContentUnits[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testSourceToManyRemoveOpContentUnits(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b, c, d, e ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentUnit{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddContentUnits(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ContentUnits(tx).Count()
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

	count, err = a.ContentUnits(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Sources) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Sources) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Sources[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Sources[0] != &a {
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

func testSourceToManyAddOpAuthors(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b, c, d, e Author

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Author{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, authorDBTypes, false, strmangle.SetComplement(authorPrimaryKeyColumns, authorColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Author{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddAuthors(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if first.R.Sources[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.Sources[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}

		if a.R.Authors[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Authors[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Authors(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testSourceToManySetOpAuthors(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b, c, d, e Author

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Author{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, authorDBTypes, false, strmangle.SetComplement(authorPrimaryKeyColumns, authorColumnsWithoutDefault)...); err != nil {
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

	err = a.SetAuthors(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Authors(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetAuthors(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Authors(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Sources) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Sources) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Sources[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.Sources[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.Authors[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Authors[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testSourceToManyRemoveOpAuthors(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b, c, d, e Author

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Author{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, authorDBTypes, false, strmangle.SetComplement(authorPrimaryKeyColumns, authorColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddAuthors(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Authors(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveAuthors(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Authors(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Sources) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Sources) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Sources[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Sources[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if len(a.R.Authors) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.Authors[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.Authors[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testSourceToOneSourceUsingParent(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Source
	var foreign Source

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
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

	slice := SourceSlice{&local}
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

func testSourceToOneSourceTypeUsingType(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Source
	var foreign SourceType

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, sourceTypeDBTypes, true, sourceTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.TypeID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Type(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := SourceSlice{&local}
	if err = local.L.LoadType(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Type == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Type = nil
	if err = local.L.LoadType(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Type == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testSourceToOneSetOpSourceUsingParent(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b, c Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
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
		err = a.SetParent(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Parent != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.ParentSources[0] != &a {
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

func testSourceToOneRemoveOpSourceUsingParent(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.ParentSources) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testSourceToOneSetOpSourceTypeUsingType(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Source
	var b, c SourceType

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, sourceTypeDBTypes, false, strmangle.SetComplement(sourceTypePrimaryKeyColumns, sourceTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, sourceTypeDBTypes, false, strmangle.SetComplement(sourceTypePrimaryKeyColumns, sourceTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*SourceType{&b, &c} {
		err = a.SetType(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Type != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.TypeSources[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.TypeID != x.ID {
			t.Error("foreign key was wrong value", a.TypeID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.TypeID))
		reflect.Indirect(reflect.ValueOf(&a.TypeID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.TypeID != x.ID {
			t.Error("foreign key was wrong value", a.TypeID, x.ID)
		}
	}
}
func testSourcesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	source := &Source{}
	if err = randomize.Struct(seed, source, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = source.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = source.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testSourcesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	source := &Source{}
	if err = randomize.Struct(seed, source, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = source.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := SourceSlice{source}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testSourcesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	source := &Source{}
	if err = randomize.Struct(seed, source, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = source.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Sources(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	sourceDBTypes = map[string]string{`CreatedAt`: `timestamp with time zone`, `Description`: `text`, `ID`: `bigint`, `Name`: `character varying`, `ParentID`: `bigint`, `Pattern`: `character varying`, `Position`: `integer`, `Properties`: `jsonb`, `TypeID`: `bigint`, `UID`: `character`}
	_             = bytes.MinRead
)

func testSourcesUpdate(t *testing.T) {
	t.Parallel()

	if len(sourceColumns) == len(sourcePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	source := &Source{}
	if err = randomize.Struct(seed, source, sourceDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = source.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Sources(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, source, sourceDBTypes, true, sourceColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	if err = source.Update(tx); err != nil {
		t.Error(err)
	}
}

func testSourcesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(sourceColumns) == len(sourcePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	source := &Source{}
	if err = randomize.Struct(seed, source, sourceDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = source.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Sources(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, source, sourceDBTypes, true, sourcePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(sourceColumns, sourcePrimaryKeyColumns) {
		fields = sourceColumns
	} else {
		fields = strmangle.SetComplement(
			sourceColumns,
			sourcePrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(source))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := SourceSlice{source}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testSourcesUpsert(t *testing.T) {
	t.Parallel()

	if len(sourceColumns) == len(sourcePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	source := Source{}
	if err = randomize.Struct(seed, &source, sourceDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = source.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Source: %s", err)
	}

	count, err := Sources(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &source, sourceDBTypes, false, sourcePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Source struct: %s", err)
	}

	if err = source.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Source: %s", err)
	}

	count, err = Sources(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}