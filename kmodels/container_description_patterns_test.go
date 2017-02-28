package kmodels

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testContainerDescriptionPatterns(t *testing.T) {
	t.Parallel()

	query := ContainerDescriptionPatterns(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testContainerDescriptionPatternsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescriptionPattern := &ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, containerDescriptionPattern, containerDescriptionPatternDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPattern.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = containerDescriptionPattern.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := ContainerDescriptionPatterns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContainerDescriptionPatternsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescriptionPattern := &ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, containerDescriptionPattern, containerDescriptionPatternDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPattern.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = ContainerDescriptionPatterns(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := ContainerDescriptionPatterns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContainerDescriptionPatternsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescriptionPattern := &ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, containerDescriptionPattern, containerDescriptionPatternDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPattern.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ContainerDescriptionPatternSlice{containerDescriptionPattern}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := ContainerDescriptionPatterns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testContainerDescriptionPatternsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescriptionPattern := &ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, containerDescriptionPattern, containerDescriptionPatternDBTypes, true, containerDescriptionPatternColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPattern.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := ContainerDescriptionPatternExists(tx, containerDescriptionPattern.ID)
	if err != nil {
		t.Errorf("Unable to check if ContainerDescriptionPattern exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ContainerDescriptionPatternExistsG to return true, but got false.")
	}
}
func testContainerDescriptionPatternsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescriptionPattern := &ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, containerDescriptionPattern, containerDescriptionPatternDBTypes, true, containerDescriptionPatternColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPattern.Insert(tx); err != nil {
		t.Error(err)
	}

	containerDescriptionPatternFound, err := FindContainerDescriptionPattern(tx, containerDescriptionPattern.ID)
	if err != nil {
		t.Error(err)
	}

	if containerDescriptionPatternFound == nil {
		t.Error("want a record, got nil")
	}
}
func testContainerDescriptionPatternsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescriptionPattern := &ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, containerDescriptionPattern, containerDescriptionPatternDBTypes, true, containerDescriptionPatternColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPattern.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = ContainerDescriptionPatterns(tx).Bind(containerDescriptionPattern); err != nil {
		t.Error(err)
	}
}

func testContainerDescriptionPatternsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescriptionPattern := &ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, containerDescriptionPattern, containerDescriptionPatternDBTypes, true, containerDescriptionPatternColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPattern.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := ContainerDescriptionPatterns(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testContainerDescriptionPatternsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescriptionPatternOne := &ContainerDescriptionPattern{}
	containerDescriptionPatternTwo := &ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, containerDescriptionPatternOne, containerDescriptionPatternDBTypes, false, containerDescriptionPatternColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}
	if err = randomize.Struct(seed, containerDescriptionPatternTwo, containerDescriptionPatternDBTypes, false, containerDescriptionPatternColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPatternOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = containerDescriptionPatternTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := ContainerDescriptionPatterns(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testContainerDescriptionPatternsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	containerDescriptionPatternOne := &ContainerDescriptionPattern{}
	containerDescriptionPatternTwo := &ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, containerDescriptionPatternOne, containerDescriptionPatternDBTypes, false, containerDescriptionPatternColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}
	if err = randomize.Struct(seed, containerDescriptionPatternTwo, containerDescriptionPatternDBTypes, false, containerDescriptionPatternColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPatternOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = containerDescriptionPatternTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContainerDescriptionPatterns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testContainerDescriptionPatternsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescriptionPattern := &ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, containerDescriptionPattern, containerDescriptionPatternDBTypes, true, containerDescriptionPatternColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPattern.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContainerDescriptionPatterns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testContainerDescriptionPatternsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescriptionPattern := &ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, containerDescriptionPattern, containerDescriptionPatternDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPattern.Insert(tx, containerDescriptionPatternColumns...); err != nil {
		t.Error(err)
	}

	count, err := ContainerDescriptionPatterns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testContainerDescriptionPatternToManyCatalogs(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContainerDescriptionPattern
	var b, c Catalog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDescriptionPatternDBTypes, true, containerDescriptionPatternColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, catalogDBTypes, false, catalogColumnsWithDefault...)
	randomize.Struct(seed, &c, catalogDBTypes, false, catalogColumnsWithDefault...)

	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec("insert into \"catalogs_container_description_patterns\" (\"container_description_pattern_id\", \"catalog_id\") values ($1, $2)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into \"catalogs_container_description_patterns\" (\"container_description_pattern_id\", \"catalog_id\") values ($1, $2)", a.ID, c.ID)
	if err != nil {
		t.Fatal(err)
	}

	catalog, err := a.Catalogs(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range catalog {
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

	slice := ContainerDescriptionPatternSlice{&a}
	if err = a.L.LoadCatalogs(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Catalogs); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Catalogs = nil
	if err = a.L.LoadCatalogs(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Catalogs); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", catalog)
	}
}

func testContainerDescriptionPatternToManyAddOpCatalogs(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContainerDescriptionPattern
	var b, c, d, e Catalog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDescriptionPatternDBTypes, false, strmangle.SetComplement(containerDescriptionPatternPrimaryKeyColumns, containerDescriptionPatternColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Catalog{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Catalog{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddCatalogs(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if first.R.ContainerDescriptionPatterns[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.ContainerDescriptionPatterns[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}

		if a.R.Catalogs[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Catalogs[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Catalogs(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testContainerDescriptionPatternToManySetOpCatalogs(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContainerDescriptionPattern
	var b, c, d, e Catalog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDescriptionPatternDBTypes, false, strmangle.SetComplement(containerDescriptionPatternPrimaryKeyColumns, containerDescriptionPatternColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Catalog{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
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

	err = a.SetCatalogs(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Catalogs(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetCatalogs(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Catalogs(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.ContainerDescriptionPatterns) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.ContainerDescriptionPatterns) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.ContainerDescriptionPatterns[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.ContainerDescriptionPatterns[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.Catalogs[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Catalogs[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testContainerDescriptionPatternToManyRemoveOpCatalogs(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContainerDescriptionPattern
	var b, c, d, e Catalog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDescriptionPatternDBTypes, false, strmangle.SetComplement(containerDescriptionPatternPrimaryKeyColumns, containerDescriptionPatternColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Catalog{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddCatalogs(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Catalogs(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveCatalogs(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Catalogs(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.ContainerDescriptionPatterns) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.ContainerDescriptionPatterns) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.ContainerDescriptionPatterns[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.ContainerDescriptionPatterns[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if len(a.R.Catalogs) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.Catalogs[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.Catalogs[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testContainerDescriptionPatternsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescriptionPattern := &ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, containerDescriptionPattern, containerDescriptionPatternDBTypes, true, containerDescriptionPatternColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPattern.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = containerDescriptionPattern.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testContainerDescriptionPatternsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescriptionPattern := &ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, containerDescriptionPattern, containerDescriptionPatternDBTypes, true, containerDescriptionPatternColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPattern.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ContainerDescriptionPatternSlice{containerDescriptionPattern}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testContainerDescriptionPatternsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescriptionPattern := &ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, containerDescriptionPattern, containerDescriptionPatternDBTypes, true, containerDescriptionPatternColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPattern.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := ContainerDescriptionPatterns(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	containerDescriptionPatternDBTypes = map[string]string{`CreatedAt`: `timestamp without time zone`, `Description`: `character varying`, `ID`: `integer`, `Lang`: `character varying`, `Pattern`: `character varying`, `UpdatedAt`: `timestamp without time zone`, `UserID`: `integer`}
	_                                  = bytes.MinRead
)

func testContainerDescriptionPatternsUpdate(t *testing.T) {
	t.Parallel()

	if len(containerDescriptionPatternColumns) == len(containerDescriptionPatternPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	containerDescriptionPattern := &ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, containerDescriptionPattern, containerDescriptionPatternDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPattern.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContainerDescriptionPatterns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, containerDescriptionPattern, containerDescriptionPatternDBTypes, true, containerDescriptionPatternColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	if err = containerDescriptionPattern.Update(tx); err != nil {
		t.Error(err)
	}
}

func testContainerDescriptionPatternsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(containerDescriptionPatternColumns) == len(containerDescriptionPatternPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	containerDescriptionPattern := &ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, containerDescriptionPattern, containerDescriptionPatternDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPattern.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContainerDescriptionPatterns(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, containerDescriptionPattern, containerDescriptionPatternDBTypes, true, containerDescriptionPatternPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(containerDescriptionPatternColumns, containerDescriptionPatternPrimaryKeyColumns) {
		fields = containerDescriptionPatternColumns
	} else {
		fields = strmangle.SetComplement(
			containerDescriptionPatternColumns,
			containerDescriptionPatternPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(containerDescriptionPattern))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := ContainerDescriptionPatternSlice{containerDescriptionPattern}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testContainerDescriptionPatternsUpsert(t *testing.T) {
	t.Parallel()

	if len(containerDescriptionPatternColumns) == len(containerDescriptionPatternPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	containerDescriptionPattern := ContainerDescriptionPattern{}
	if err = randomize.Struct(seed, &containerDescriptionPattern, containerDescriptionPatternDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionPattern.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert ContainerDescriptionPattern: %s", err)
	}

	count, err := ContainerDescriptionPatterns(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &containerDescriptionPattern, containerDescriptionPatternDBTypes, false, containerDescriptionPatternPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ContainerDescriptionPattern struct: %s", err)
	}

	if err = containerDescriptionPattern.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert ContainerDescriptionPattern: %s", err)
	}

	count, err = ContainerDescriptionPatterns(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
