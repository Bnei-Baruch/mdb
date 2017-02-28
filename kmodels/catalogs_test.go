package kmodels

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testCatalogs(t *testing.T) {
	t.Parallel()

	query := Catalogs(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testCatalogsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalog := &Catalog{}
	if err = randomize.Struct(seed, catalog, catalogDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalog.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = catalog.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Catalogs(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCatalogsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalog := &Catalog{}
	if err = randomize.Struct(seed, catalog, catalogDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalog.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Catalogs(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Catalogs(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCatalogsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalog := &Catalog{}
	if err = randomize.Struct(seed, catalog, catalogDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalog.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := CatalogSlice{catalog}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Catalogs(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testCatalogsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalog := &Catalog{}
	if err = randomize.Struct(seed, catalog, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalog.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := CatalogExists(tx, catalog.ID)
	if err != nil {
		t.Errorf("Unable to check if Catalog exists: %s", err)
	}
	if !e {
		t.Errorf("Expected CatalogExistsG to return true, but got false.")
	}
}
func testCatalogsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalog := &Catalog{}
	if err = randomize.Struct(seed, catalog, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalog.Insert(tx); err != nil {
		t.Error(err)
	}

	catalogFound, err := FindCatalog(tx, catalog.ID)
	if err != nil {
		t.Error(err)
	}

	if catalogFound == nil {
		t.Error("want a record, got nil")
	}
}
func testCatalogsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalog := &Catalog{}
	if err = randomize.Struct(seed, catalog, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalog.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Catalogs(tx).Bind(catalog); err != nil {
		t.Error(err)
	}
}

func testCatalogsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalog := &Catalog{}
	if err = randomize.Struct(seed, catalog, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalog.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Catalogs(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testCatalogsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalogOne := &Catalog{}
	catalogTwo := &Catalog{}
	if err = randomize.Struct(seed, catalogOne, catalogDBTypes, false, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}
	if err = randomize.Struct(seed, catalogTwo, catalogDBTypes, false, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = catalogTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Catalogs(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testCatalogsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	catalogOne := &Catalog{}
	catalogTwo := &Catalog{}
	if err = randomize.Struct(seed, catalogOne, catalogDBTypes, false, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}
	if err = randomize.Struct(seed, catalogTwo, catalogDBTypes, false, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = catalogTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Catalogs(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testCatalogsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalog := &Catalog{}
	if err = randomize.Struct(seed, catalog, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalog.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Catalogs(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testCatalogsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalog := &Catalog{}
	if err = randomize.Struct(seed, catalog, catalogDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalog.Insert(tx, catalogColumns...); err != nil {
		t.Error(err)
	}

	count, err := Catalogs(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testCatalogToManyContainers(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b, c Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, containerDBTypes, false, containerColumnsWithDefault...)
	randomize.Struct(seed, &c, containerDBTypes, false, containerColumnsWithDefault...)

	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec("insert into \"catalogs_containers\" (\"catalog_id\", \"container_id\") values ($1, $2)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into \"catalogs_containers\" (\"catalog_id\", \"container_id\") values ($1, $2)", a.ID, c.ID)
	if err != nil {
		t.Fatal(err)
	}

	container, err := a.Containers(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range container {
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

	slice := CatalogSlice{&a}
	if err = a.L.LoadContainers(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Containers); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Containers = nil
	if err = a.L.LoadContainers(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Containers); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", container)
	}
}

func testCatalogToManyCatalogDescriptions(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b, c CatalogDescription

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, catalogDescriptionDBTypes, false, catalogDescriptionColumnsWithDefault...)
	randomize.Struct(seed, &c, catalogDescriptionDBTypes, false, catalogDescriptionColumnsWithDefault...)

	b.CatalogID = a.ID
	c.CatalogID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	catalogDescription, err := a.CatalogDescriptions(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range catalogDescription {
		if v.CatalogID == b.CatalogID {
			bFound = true
		}
		if v.CatalogID == c.CatalogID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := CatalogSlice{&a}
	if err = a.L.LoadCatalogDescriptions(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.CatalogDescriptions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.CatalogDescriptions = nil
	if err = a.L.LoadCatalogDescriptions(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.CatalogDescriptions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", catalogDescription)
	}
}

func testCatalogToManyParentCatalogs(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b, c Catalog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, catalogDBTypes, false, catalogColumnsWithDefault...)
	randomize.Struct(seed, &c, catalogDBTypes, false, catalogColumnsWithDefault...)

	b.ParentID.Valid = true
	c.ParentID.Valid = true
	b.ParentID.Int = a.ID
	c.ParentID.Int = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	catalog, err := a.ParentCatalogs(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range catalog {
		if v.ParentID.Int == b.ParentID.Int {
			bFound = true
		}
		if v.ParentID.Int == c.ParentID.Int {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := CatalogSlice{&a}
	if err = a.L.LoadParentCatalogs(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ParentCatalogs); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.ParentCatalogs = nil
	if err = a.L.LoadParentCatalogs(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ParentCatalogs); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", catalog)
	}
}

func testCatalogToManyContainerDescriptionPatterns(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b, c ContainerDescriptionPattern

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, containerDescriptionPatternDBTypes, false, containerDescriptionPatternColumnsWithDefault...)
	randomize.Struct(seed, &c, containerDescriptionPatternDBTypes, false, containerDescriptionPatternColumnsWithDefault...)

	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec("insert into \"catalogs_container_description_patterns\" (\"catalog_id\", \"container_description_pattern_id\") values ($1, $2)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into \"catalogs_container_description_patterns\" (\"catalog_id\", \"container_description_pattern_id\") values ($1, $2)", a.ID, c.ID)
	if err != nil {
		t.Fatal(err)
	}

	containerDescriptionPattern, err := a.ContainerDescriptionPatterns(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range containerDescriptionPattern {
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

	slice := CatalogSlice{&a}
	if err = a.L.LoadContainerDescriptionPatterns(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ContainerDescriptionPatterns); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.ContainerDescriptionPatterns = nil
	if err = a.L.LoadContainerDescriptionPatterns(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ContainerDescriptionPatterns); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", containerDescriptionPattern)
	}
}

func testCatalogToManyAddOpContainers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b, c, d, e Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Container{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Container{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddContainers(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if first.R.Catalogs[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.Catalogs[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}

		if a.R.Containers[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Containers[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Containers(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testCatalogToManySetOpContainers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b, c, d, e Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Container{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
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

	err = a.SetContainers(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Containers(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetContainers(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Containers(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Catalogs) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Catalogs) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Catalogs[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.Catalogs[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.Containers[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Containers[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testCatalogToManyRemoveOpContainers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b, c, d, e Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Container{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddContainers(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Containers(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveContainers(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Containers(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Catalogs) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Catalogs) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Catalogs[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Catalogs[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if len(a.R.Containers) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.Containers[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.Containers[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testCatalogToManyAddOpCatalogDescriptions(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b, c, d, e CatalogDescription

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*CatalogDescription{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, catalogDescriptionDBTypes, false, strmangle.SetComplement(catalogDescriptionPrimaryKeyColumns, catalogDescriptionColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*CatalogDescription{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddCatalogDescriptions(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.CatalogID {
			t.Error("foreign key was wrong value", a.ID, first.CatalogID)
		}
		if a.ID != second.CatalogID {
			t.Error("foreign key was wrong value", a.ID, second.CatalogID)
		}

		if first.R.Catalog != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Catalog != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.CatalogDescriptions[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.CatalogDescriptions[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.CatalogDescriptions(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testCatalogToManyAddOpParentCatalogs(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b, c, d, e Catalog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
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
		err = a.AddParentCatalogs(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.ParentID.Int {
			t.Error("foreign key was wrong value", a.ID, first.ParentID.Int)
		}
		if a.ID != second.ParentID.Int {
			t.Error("foreign key was wrong value", a.ID, second.ParentID.Int)
		}

		if first.R.Parent != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Parent != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.ParentCatalogs[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.ParentCatalogs[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.ParentCatalogs(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testCatalogToManySetOpParentCatalogs(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b, c, d, e Catalog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
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

	err = a.SetParentCatalogs(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ParentCatalogs(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetParentCatalogs(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ParentCatalogs(tx).Count()
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
	if a.ID != d.ParentID.Int {
		t.Error("foreign key was wrong value", a.ID, d.ParentID.Int)
	}
	if a.ID != e.ParentID.Int {
		t.Error("foreign key was wrong value", a.ID, e.ParentID.Int)
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

	if a.R.ParentCatalogs[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.ParentCatalogs[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testCatalogToManyRemoveOpParentCatalogs(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b, c, d, e Catalog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
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

	err = a.AddParentCatalogs(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ParentCatalogs(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveParentCatalogs(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ParentCatalogs(tx).Count()
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

	if len(a.R.ParentCatalogs) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.ParentCatalogs[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.ParentCatalogs[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testCatalogToManyAddOpContainerDescriptionPatterns(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b, c, d, e ContainerDescriptionPattern

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContainerDescriptionPattern{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, containerDescriptionPatternDBTypes, false, strmangle.SetComplement(containerDescriptionPatternPrimaryKeyColumns, containerDescriptionPatternColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*ContainerDescriptionPattern{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddContainerDescriptionPatterns(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if first.R.Catalogs[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.Catalogs[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}

		if a.R.ContainerDescriptionPatterns[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.ContainerDescriptionPatterns[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.ContainerDescriptionPatterns(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testCatalogToManySetOpContainerDescriptionPatterns(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b, c, d, e ContainerDescriptionPattern

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContainerDescriptionPattern{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, containerDescriptionPatternDBTypes, false, strmangle.SetComplement(containerDescriptionPatternPrimaryKeyColumns, containerDescriptionPatternColumnsWithoutDefault)...); err != nil {
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

	err = a.SetContainerDescriptionPatterns(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ContainerDescriptionPatterns(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetContainerDescriptionPatterns(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ContainerDescriptionPatterns(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Catalogs) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Catalogs) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Catalogs[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.Catalogs[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.ContainerDescriptionPatterns[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.ContainerDescriptionPatterns[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testCatalogToManyRemoveOpContainerDescriptionPatterns(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b, c, d, e ContainerDescriptionPattern

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContainerDescriptionPattern{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, containerDescriptionPatternDBTypes, false, strmangle.SetComplement(containerDescriptionPatternPrimaryKeyColumns, containerDescriptionPatternColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddContainerDescriptionPatterns(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ContainerDescriptionPatterns(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveContainerDescriptionPatterns(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ContainerDescriptionPatterns(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Catalogs) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Catalogs) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Catalogs[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Catalogs[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if len(a.R.ContainerDescriptionPatterns) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.ContainerDescriptionPatterns[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.ContainerDescriptionPatterns[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testCatalogToOneCatalogUsingParent(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Catalog
	var foreign Catalog

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	local.ParentID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.ParentID.Int = foreign.ID
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

	slice := CatalogSlice{&local}
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

func testCatalogToOneUserUsingUser(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Catalog
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	local.UserID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.UserID.Int = foreign.ID
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

	slice := CatalogSlice{&local}
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

func testCatalogToOneSetOpCatalogUsingParent(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b, c Catalog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Catalog{&b, &c} {
		err = a.SetParent(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Parent != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.ParentCatalogs[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.ParentID.Int != x.ID {
			t.Error("foreign key was wrong value", a.ParentID.Int)
		}

		zero := reflect.Zero(reflect.TypeOf(a.ParentID.Int))
		reflect.Indirect(reflect.ValueOf(&a.ParentID.Int)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.ParentID.Int != x.ID {
			t.Error("foreign key was wrong value", a.ParentID.Int, x.ID)
		}
	}
}

func testCatalogToOneRemoveOpCatalogUsingParent(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b Catalog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.ParentCatalogs) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testCatalogToOneSetOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
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

		if x.R.Catalogs[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.UserID.Int != x.ID {
			t.Error("foreign key was wrong value", a.UserID.Int)
		}

		zero := reflect.Zero(reflect.TypeOf(a.UserID.Int))
		reflect.Indirect(reflect.ValueOf(&a.UserID.Int)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.UserID.Int != x.ID {
			t.Error("foreign key was wrong value", a.UserID.Int, x.ID)
		}
	}
}

func testCatalogToOneRemoveOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Catalog
	var b User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDBTypes, false, strmangle.SetComplement(catalogPrimaryKeyColumns, catalogColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.Catalogs) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testCatalogsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalog := &Catalog{}
	if err = randomize.Struct(seed, catalog, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalog.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = catalog.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testCatalogsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalog := &Catalog{}
	if err = randomize.Struct(seed, catalog, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalog.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := CatalogSlice{catalog}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testCatalogsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalog := &Catalog{}
	if err = randomize.Struct(seed, catalog, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalog.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Catalogs(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	catalogDBTypes = map[string]string{`BooksCatalog`: `boolean`, `Catorder`: `integer`, `CreatedAt`: `timestamp without time zone`, `ID`: `integer`, `Label`: `character varying`, `Name`: `character varying`, `Open`: `boolean`, `ParentID`: `integer`, `Secure`: `integer`, `SelectedCatalog`: `integer`, `UpdatedAt`: `timestamp without time zone`, `UserID`: `integer`, `Visible`: `boolean`}
	_              = bytes.MinRead
)

func testCatalogsUpdate(t *testing.T) {
	t.Parallel()

	if len(catalogColumns) == len(catalogPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	catalog := &Catalog{}
	if err = randomize.Struct(seed, catalog, catalogDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalog.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Catalogs(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, catalog, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	if err = catalog.Update(tx); err != nil {
		t.Error(err)
	}
}

func testCatalogsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(catalogColumns) == len(catalogPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	catalog := &Catalog{}
	if err = randomize.Struct(seed, catalog, catalogDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalog.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Catalogs(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, catalog, catalogDBTypes, true, catalogPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(catalogColumns, catalogPrimaryKeyColumns) {
		fields = catalogColumns
	} else {
		fields = strmangle.SetComplement(
			catalogColumns,
			catalogPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(catalog))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := CatalogSlice{catalog}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testCatalogsUpsert(t *testing.T) {
	t.Parallel()

	if len(catalogColumns) == len(catalogPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	catalog := Catalog{}
	if err = randomize.Struct(seed, &catalog, catalogDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalog.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Catalog: %s", err)
	}

	count, err := Catalogs(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &catalog, catalogDBTypes, false, catalogPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	if err = catalog.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Catalog: %s", err)
	}

	count, err = Catalogs(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
