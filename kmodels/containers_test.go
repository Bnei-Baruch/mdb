package kmodels

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testContainers(t *testing.T) {
	t.Parallel()

	query := Containers(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testContainersDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	container := &Container{}
	if err = randomize.Struct(seed, container, containerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = container.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = container.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Containers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContainersQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	container := &Container{}
	if err = randomize.Struct(seed, container, containerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = container.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Containers(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Containers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContainersSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	container := &Container{}
	if err = randomize.Struct(seed, container, containerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = container.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ContainerSlice{container}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Containers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testContainersExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	container := &Container{}
	if err = randomize.Struct(seed, container, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = container.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := ContainerExists(tx, container.ID)
	if err != nil {
		t.Errorf("Unable to check if Container exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ContainerExistsG to return true, but got false.")
	}
}
func testContainersFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	container := &Container{}
	if err = randomize.Struct(seed, container, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = container.Insert(tx); err != nil {
		t.Error(err)
	}

	containerFound, err := FindContainer(tx, container.ID)
	if err != nil {
		t.Error(err)
	}

	if containerFound == nil {
		t.Error("want a record, got nil")
	}
}
func testContainersBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	container := &Container{}
	if err = randomize.Struct(seed, container, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = container.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Containers(tx).Bind(container); err != nil {
		t.Error(err)
	}
}

func testContainersOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	container := &Container{}
	if err = randomize.Struct(seed, container, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = container.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Containers(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testContainersAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerOne := &Container{}
	containerTwo := &Container{}
	if err = randomize.Struct(seed, containerOne, containerDBTypes, false, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}
	if err = randomize.Struct(seed, containerTwo, containerDBTypes, false, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = containerTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Containers(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testContainersCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	containerOne := &Container{}
	containerTwo := &Container{}
	if err = randomize.Struct(seed, containerOne, containerDBTypes, false, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}
	if err = randomize.Struct(seed, containerTwo, containerDBTypes, false, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = containerTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Containers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testContainersInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	container := &Container{}
	if err = randomize.Struct(seed, container, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = container.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Containers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testContainersInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	container := &Container{}
	if err = randomize.Struct(seed, container, containerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = container.Insert(tx, containerColumns...); err != nil {
		t.Error(err)
	}

	count, err := Containers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testContainerToManyCatalogs(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c Catalog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
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

	_, err = tx.Exec("insert into \"catalogs_containers\" (\"container_id\", \"catalog_id\") values ($1, $2)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into \"catalogs_containers\" (\"container_id\", \"catalog_id\") values ($1, $2)", a.ID, c.ID)
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

	slice := ContainerSlice{&a}
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

func testContainerToManyFileAssets(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c FileAsset

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, fileAssetDBTypes, false, fileAssetColumnsWithDefault...)
	randomize.Struct(seed, &c, fileAssetDBTypes, false, fileAssetColumnsWithDefault...)

	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec("insert into \"containers_file_assets\" (\"container_id\", \"file_asset_id\") values ($1, $2)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into \"containers_file_assets\" (\"container_id\", \"file_asset_id\") values ($1, $2)", a.ID, c.ID)
	if err != nil {
		t.Fatal(err)
	}

	fileAsset, err := a.FileAssets(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range fileAsset {
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

	slice := ContainerSlice{&a}
	if err = a.L.LoadFileAssets(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.FileAssets); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.FileAssets = nil
	if err = a.L.LoadFileAssets(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.FileAssets); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", fileAsset)
	}
}

func testContainerToManyContainerDescriptions(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c ContainerDescription

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, containerDescriptionDBTypes, false, containerDescriptionColumnsWithDefault...)
	randomize.Struct(seed, &c, containerDescriptionDBTypes, false, containerDescriptionColumnsWithDefault...)

	b.ContainerID = a.ID
	c.ContainerID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	containerDescription, err := a.ContainerDescriptions(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range containerDescription {
		if v.ContainerID == b.ContainerID {
			bFound = true
		}
		if v.ContainerID == c.ContainerID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ContainerSlice{&a}
	if err = a.L.LoadContainerDescriptions(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ContainerDescriptions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.ContainerDescriptions = nil
	if err = a.L.LoadContainerDescriptions(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ContainerDescriptions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", containerDescription)
	}
}

func testContainerToManyLabels(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c Label

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, labelDBTypes, false, labelColumnsWithDefault...)
	randomize.Struct(seed, &c, labelDBTypes, false, labelColumnsWithDefault...)

	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec("insert into \"containers_labels\" (\"container_id\", \"label_id\") values ($1, $2)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into \"containers_labels\" (\"container_id\", \"label_id\") values ($1, $2)", a.ID, c.ID)
	if err != nil {
		t.Fatal(err)
	}

	label, err := a.Labels(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range label {
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

	slice := ContainerSlice{&a}
	if err = a.L.LoadLabels(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Labels); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Labels = nil
	if err = a.L.LoadLabels(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Labels); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", label)
	}
}

func testContainerToManyAddOpCatalogs(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c, d, e Catalog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
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

		if first.R.Containers[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.Containers[0] != &a {
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

func testContainerToManySetOpCatalogs(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c, d, e Catalog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.Containers) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Containers) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Containers[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.Containers[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.Catalogs[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Catalogs[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testContainerToManyRemoveOpCatalogs(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c, d, e Catalog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.Containers) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Containers) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Containers[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Containers[0] != &a {
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

func testContainerToManyAddOpFileAssets(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c, d, e FileAsset

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*FileAsset{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*FileAsset{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddFileAssets(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if first.R.Containers[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.Containers[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}

		if a.R.FileAssets[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.FileAssets[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.FileAssets(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testContainerToManySetOpFileAssets(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c, d, e FileAsset

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*FileAsset{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
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

	err = a.SetFileAssets(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.FileAssets(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetFileAssets(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.FileAssets(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Containers) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Containers) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Containers[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.Containers[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.FileAssets[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.FileAssets[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testContainerToManyRemoveOpFileAssets(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c, d, e FileAsset

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*FileAsset{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddFileAssets(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.FileAssets(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveFileAssets(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.FileAssets(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Containers) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Containers) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Containers[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Containers[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if len(a.R.FileAssets) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.FileAssets[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.FileAssets[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testContainerToManyAddOpContainerDescriptions(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c, d, e ContainerDescription

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContainerDescription{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, containerDescriptionDBTypes, false, strmangle.SetComplement(containerDescriptionPrimaryKeyColumns, containerDescriptionColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*ContainerDescription{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddContainerDescriptions(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.ContainerID {
			t.Error("foreign key was wrong value", a.ID, first.ContainerID)
		}
		if a.ID != second.ContainerID {
			t.Error("foreign key was wrong value", a.ID, second.ContainerID)
		}

		if first.R.Container != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Container != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.ContainerDescriptions[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.ContainerDescriptions[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.ContainerDescriptions(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testContainerToManyAddOpLabels(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c, d, e Label

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Label{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, labelDBTypes, false, strmangle.SetComplement(labelPrimaryKeyColumns, labelColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Label{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddLabels(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if first.R.Containers[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.Containers[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}

		if a.R.Labels[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Labels[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Labels(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testContainerToManySetOpLabels(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c, d, e Label

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Label{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, labelDBTypes, false, strmangle.SetComplement(labelPrimaryKeyColumns, labelColumnsWithoutDefault)...); err != nil {
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

	err = a.SetLabels(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Labels(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetLabels(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Labels(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Containers) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Containers) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Containers[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.Containers[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.Labels[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Labels[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testContainerToManyRemoveOpLabels(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c, d, e Label

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Label{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, labelDBTypes, false, strmangle.SetComplement(labelPrimaryKeyColumns, labelColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddLabels(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Labels(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveLabels(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Labels(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Containers) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Containers) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Containers[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Containers[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if len(a.R.Labels) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.Labels[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.Labels[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testContainerToOneLanguageUsingLang(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Container
	var foreign Language

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, languageDBTypes, true, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	local.LangID.Valid = true
	foreign.Code3.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.LangID.String = foreign.Code3.String
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Lang(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.Code3.String != foreign.Code3.String {
		t.Errorf("want: %v, got %v", foreign.Code3.String, check.Code3.String)
	}

	slice := ContainerSlice{&local}
	if err = local.L.LoadLang(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Lang == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Lang = nil
	if err = local.L.LoadLang(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Lang == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testContainerToOneContentTypeUsingContentType(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Container
	var foreign ContentType

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	local.ContentTypeID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.ContentTypeID.Int = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.ContentType(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ContainerSlice{&local}
	if err = local.L.LoadContentType(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.ContentType == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.ContentType = nil
	if err = local.L.LoadContentType(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.ContentType == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testContainerToOneVirtualLessonUsingVirtualLesson(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Container
	var foreign VirtualLesson

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, virtualLessonDBTypes, true, virtualLessonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	local.VirtualLessonID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.VirtualLessonID.Int = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.VirtualLesson(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ContainerSlice{&local}
	if err = local.L.LoadVirtualLesson(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.VirtualLesson == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.VirtualLesson = nil
	if err = local.L.LoadVirtualLesson(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.VirtualLesson == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testContainerToOneSetOpLanguageUsingLang(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c Language

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, languageDBTypes, false, strmangle.SetComplement(languagePrimaryKeyColumns, languageColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, languageDBTypes, false, strmangle.SetComplement(languagePrimaryKeyColumns, languageColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Language{&b, &c} {
		err = a.SetLang(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Lang != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.LangContainers[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.LangID.String != x.Code3.String {
			t.Error("foreign key was wrong value", a.LangID.String)
		}

		zero := reflect.Zero(reflect.TypeOf(a.LangID.String))
		reflect.Indirect(reflect.ValueOf(&a.LangID.String)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.LangID.String != x.Code3.String {
			t.Error("foreign key was wrong value", a.LangID.String, x.Code3.String)
		}
	}
}

func testContainerToOneRemoveOpLanguageUsingLang(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b Language

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, languageDBTypes, false, strmangle.SetComplement(languagePrimaryKeyColumns, languageColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.SetLang(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveLang(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.Lang(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.Lang != nil {
		t.Error("R struct entry should be nil")
	}

	if a.LangID.Valid {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.LangContainers) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testContainerToOneSetOpContentTypeUsingContentType(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c ContentType

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, contentTypeDBTypes, false, strmangle.SetComplement(contentTypePrimaryKeyColumns, contentTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, contentTypeDBTypes, false, strmangle.SetComplement(contentTypePrimaryKeyColumns, contentTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*ContentType{&b, &c} {
		err = a.SetContentType(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.ContentType != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.Containers[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.ContentTypeID.Int != x.ID {
			t.Error("foreign key was wrong value", a.ContentTypeID.Int)
		}

		zero := reflect.Zero(reflect.TypeOf(a.ContentTypeID.Int))
		reflect.Indirect(reflect.ValueOf(&a.ContentTypeID.Int)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.ContentTypeID.Int != x.ID {
			t.Error("foreign key was wrong value", a.ContentTypeID.Int, x.ID)
		}
	}
}

func testContainerToOneRemoveOpContentTypeUsingContentType(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b ContentType

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, contentTypeDBTypes, false, strmangle.SetComplement(contentTypePrimaryKeyColumns, contentTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.SetContentType(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveContentType(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.ContentType(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.ContentType != nil {
		t.Error("R struct entry should be nil")
	}

	if a.ContentTypeID.Valid {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.Containers) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testContainerToOneSetOpVirtualLessonUsingVirtualLesson(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b, c VirtualLesson

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, virtualLessonDBTypes, false, strmangle.SetComplement(virtualLessonPrimaryKeyColumns, virtualLessonColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, virtualLessonDBTypes, false, strmangle.SetComplement(virtualLessonPrimaryKeyColumns, virtualLessonColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*VirtualLesson{&b, &c} {
		err = a.SetVirtualLesson(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.VirtualLesson != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.Containers[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.VirtualLessonID.Int != x.ID {
			t.Error("foreign key was wrong value", a.VirtualLessonID.Int)
		}

		zero := reflect.Zero(reflect.TypeOf(a.VirtualLessonID.Int))
		reflect.Indirect(reflect.ValueOf(&a.VirtualLessonID.Int)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.VirtualLessonID.Int != x.ID {
			t.Error("foreign key was wrong value", a.VirtualLessonID.Int, x.ID)
		}
	}
}

func testContainerToOneRemoveOpVirtualLessonUsingVirtualLesson(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Container
	var b VirtualLesson

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, virtualLessonDBTypes, false, strmangle.SetComplement(virtualLessonPrimaryKeyColumns, virtualLessonColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.SetVirtualLesson(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveVirtualLesson(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.VirtualLesson(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.VirtualLesson != nil {
		t.Error("R struct entry should be nil")
	}

	if a.VirtualLessonID.Valid {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.Containers) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testContainersReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	container := &Container{}
	if err = randomize.Struct(seed, container, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = container.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = container.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testContainersReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	container := &Container{}
	if err = randomize.Struct(seed, container, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = container.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ContainerSlice{container}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testContainersSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	container := &Container{}
	if err = randomize.Struct(seed, container, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = container.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Containers(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	containerDBTypes = map[string]string{`AutoParsed`: `boolean`, `CensorID`: `integer`, `ClosedByCensor`: `boolean`, `ContentTypeID`: `integer`, `CreatedAt`: `timestamp without time zone`, `Filmdate`: `date`, `ForCensorship`: `boolean`, `ID`: `integer`, `LangID`: `character`, `LecturerID`: `integer`, `MarkedForMerge`: `boolean`, `Name`: `character varying`, `OpenedByCensor`: `boolean`, `PlaytimeSecs`: `integer`, `Position`: `integer`, `Secure`: `integer`, `SecureChanged`: `boolean`, `UpdatedAt`: `timestamp without time zone`, `UserID`: `integer`, `VirtualLessonID`: `integer`}
	_                = bytes.MinRead
)

func testContainersUpdate(t *testing.T) {
	t.Parallel()

	if len(containerColumns) == len(containerPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	container := &Container{}
	if err = randomize.Struct(seed, container, containerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = container.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Containers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, container, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	if err = container.Update(tx); err != nil {
		t.Error(err)
	}
}

func testContainersSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(containerColumns) == len(containerPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	container := &Container{}
	if err = randomize.Struct(seed, container, containerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = container.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Containers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, container, containerDBTypes, true, containerPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(containerColumns, containerPrimaryKeyColumns) {
		fields = containerColumns
	} else {
		fields = strmangle.SetComplement(
			containerColumns,
			containerPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(container))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := ContainerSlice{container}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testContainersUpsert(t *testing.T) {
	t.Parallel()

	if len(containerColumns) == len(containerPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	container := Container{}
	if err = randomize.Struct(seed, &container, containerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = container.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Container: %s", err)
	}

	count, err := Containers(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &container, containerDBTypes, false, containerPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	if err = container.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Container: %s", err)
	}

	count, err = Containers(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
