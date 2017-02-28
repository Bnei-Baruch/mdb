package kmodels

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testFileAssets(t *testing.T) {
	t.Parallel()

	query := FileAssets(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testFileAssetsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAsset := &FileAsset{}
	if err = randomize.Struct(seed, fileAsset, fileAssetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAsset.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = fileAsset.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := FileAssets(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testFileAssetsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAsset := &FileAsset{}
	if err = randomize.Struct(seed, fileAsset, fileAssetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAsset.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = FileAssets(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := FileAssets(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testFileAssetsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAsset := &FileAsset{}
	if err = randomize.Struct(seed, fileAsset, fileAssetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAsset.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := FileAssetSlice{fileAsset}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := FileAssets(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testFileAssetsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAsset := &FileAsset{}
	if err = randomize.Struct(seed, fileAsset, fileAssetDBTypes, true, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAsset.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := FileAssetExists(tx, fileAsset.ID)
	if err != nil {
		t.Errorf("Unable to check if FileAsset exists: %s", err)
	}
	if !e {
		t.Errorf("Expected FileAssetExistsG to return true, but got false.")
	}
}
func testFileAssetsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAsset := &FileAsset{}
	if err = randomize.Struct(seed, fileAsset, fileAssetDBTypes, true, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAsset.Insert(tx); err != nil {
		t.Error(err)
	}

	fileAssetFound, err := FindFileAsset(tx, fileAsset.ID)
	if err != nil {
		t.Error(err)
	}

	if fileAssetFound == nil {
		t.Error("want a record, got nil")
	}
}
func testFileAssetsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAsset := &FileAsset{}
	if err = randomize.Struct(seed, fileAsset, fileAssetDBTypes, true, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAsset.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = FileAssets(tx).Bind(fileAsset); err != nil {
		t.Error(err)
	}
}

func testFileAssetsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAsset := &FileAsset{}
	if err = randomize.Struct(seed, fileAsset, fileAssetDBTypes, true, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAsset.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := FileAssets(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testFileAssetsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAssetOne := &FileAsset{}
	fileAssetTwo := &FileAsset{}
	if err = randomize.Struct(seed, fileAssetOne, fileAssetDBTypes, false, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}
	if err = randomize.Struct(seed, fileAssetTwo, fileAssetDBTypes, false, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = fileAssetTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := FileAssets(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testFileAssetsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	fileAssetOne := &FileAsset{}
	fileAssetTwo := &FileAsset{}
	if err = randomize.Struct(seed, fileAssetOne, fileAssetDBTypes, false, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}
	if err = randomize.Struct(seed, fileAssetTwo, fileAssetDBTypes, false, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAssetOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = fileAssetTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := FileAssets(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testFileAssetsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAsset := &FileAsset{}
	if err = randomize.Struct(seed, fileAsset, fileAssetDBTypes, true, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAsset.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := FileAssets(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testFileAssetsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAsset := &FileAsset{}
	if err = randomize.Struct(seed, fileAsset, fileAssetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAsset.Insert(tx, fileAssetColumns...); err != nil {
		t.Error(err)
	}

	count, err := FileAssets(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testFileAssetToManyContainers(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a FileAsset
	var b, c Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileAssetDBTypes, true, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
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

	_, err = tx.Exec("insert into \"containers_file_assets\" (\"file_asset_id\", \"container_id\") values ($1, $2)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into \"containers_file_assets\" (\"file_asset_id\", \"container_id\") values ($1, $2)", a.ID, c.ID)
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

	slice := FileAssetSlice{&a}
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

func testFileAssetToManyFileFileAssetDescriptions(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a FileAsset
	var b, c FileAssetDescription

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileAssetDBTypes, true, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, fileAssetDescriptionDBTypes, false, fileAssetDescriptionColumnsWithDefault...)
	randomize.Struct(seed, &c, fileAssetDescriptionDBTypes, false, fileAssetDescriptionColumnsWithDefault...)

	b.FileID = a.ID
	c.FileID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	fileAssetDescription, err := a.FileFileAssetDescriptions(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range fileAssetDescription {
		if v.FileID == b.FileID {
			bFound = true
		}
		if v.FileID == c.FileID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := FileAssetSlice{&a}
	if err = a.L.LoadFileFileAssetDescriptions(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.FileFileAssetDescriptions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.FileFileAssetDescriptions = nil
	if err = a.L.LoadFileFileAssetDescriptions(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.FileFileAssetDescriptions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", fileAssetDescription)
	}
}

func testFileAssetToManyAddOpContainers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a FileAsset
	var b, c, d, e Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
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

		if first.R.FileAssets[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.FileAssets[0] != &a {
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

func testFileAssetToManySetOpContainers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a FileAsset
	var b, c, d, e Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.FileAssets) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.FileAssets) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.FileAssets[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.FileAssets[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.Containers[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Containers[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testFileAssetToManyRemoveOpContainers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a FileAsset
	var b, c, d, e Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.FileAssets) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.FileAssets) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.FileAssets[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.FileAssets[0] != &a {
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

func testFileAssetToManyAddOpFileFileAssetDescriptions(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a FileAsset
	var b, c, d, e FileAssetDescription

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*FileAssetDescription{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, fileAssetDescriptionDBTypes, false, strmangle.SetComplement(fileAssetDescriptionPrimaryKeyColumns, fileAssetDescriptionColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*FileAssetDescription{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddFileFileAssetDescriptions(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.FileID {
			t.Error("foreign key was wrong value", a.ID, first.FileID)
		}
		if a.ID != second.FileID {
			t.Error("foreign key was wrong value", a.ID, second.FileID)
		}

		if first.R.File != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.File != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.FileFileAssetDescriptions[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.FileFileAssetDescriptions[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.FileFileAssetDescriptions(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testFileAssetToOneLanguageUsingLang(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local FileAsset
	var foreign Language

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, fileAssetDBTypes, true, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
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

	slice := FileAssetSlice{&local}
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

func testFileAssetToOneUserUsingUser(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local FileAsset
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, fileAssetDBTypes, true, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
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

	slice := FileAssetSlice{&local}
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

func testFileAssetToOneServerUsingServername(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local FileAsset
	var foreign Server

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, fileAssetDBTypes, true, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, serverDBTypes, true, serverColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	local.ServernameID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.ServernameID.String = foreign.Servername
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Servername(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.Servername != foreign.Servername {
		t.Errorf("want: %v, got %v", foreign.Servername, check.Servername)
	}

	slice := FileAssetSlice{&local}
	if err = local.L.LoadServername(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Servername == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Servername = nil
	if err = local.L.LoadServername(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Servername == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testFileAssetToOneSetOpLanguageUsingLang(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a FileAsset
	var b, c Language

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
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

		if x.R.LangFileAssets[0] != &a {
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

func testFileAssetToOneRemoveOpLanguageUsingLang(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a FileAsset
	var b Language

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.LangFileAssets) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testFileAssetToOneSetOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a FileAsset
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
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

		if x.R.FileAssets[0] != &a {
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

func testFileAssetToOneRemoveOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a FileAsset
	var b User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.FileAssets) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testFileAssetToOneSetOpServerUsingServername(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a FileAsset
	var b, c Server

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, serverDBTypes, false, strmangle.SetComplement(serverPrimaryKeyColumns, serverColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, serverDBTypes, false, strmangle.SetComplement(serverPrimaryKeyColumns, serverColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Server{&b, &c} {
		err = a.SetServername(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Servername != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.ServernameFileAssets[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.ServernameID.String != x.Servername {
			t.Error("foreign key was wrong value", a.ServernameID.String)
		}

		zero := reflect.Zero(reflect.TypeOf(a.ServernameID.String))
		reflect.Indirect(reflect.ValueOf(&a.ServernameID.String)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.ServernameID.String != x.Servername {
			t.Error("foreign key was wrong value", a.ServernameID.String, x.Servername)
		}
	}
}

func testFileAssetToOneRemoveOpServerUsingServername(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a FileAsset
	var b Server

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, serverDBTypes, false, strmangle.SetComplement(serverPrimaryKeyColumns, serverColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.SetServername(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveServername(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.Servername(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.Servername != nil {
		t.Error("R struct entry should be nil")
	}

	if a.ServernameID.Valid {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.ServernameFileAssets) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testFileAssetsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAsset := &FileAsset{}
	if err = randomize.Struct(seed, fileAsset, fileAssetDBTypes, true, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAsset.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = fileAsset.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testFileAssetsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAsset := &FileAsset{}
	if err = randomize.Struct(seed, fileAsset, fileAssetDBTypes, true, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAsset.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := FileAssetSlice{fileAsset}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testFileAssetsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileAsset := &FileAsset{}
	if err = randomize.Struct(seed, fileAsset, fileAssetDBTypes, true, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAsset.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := FileAssets(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	fileAssetDBTypes = map[string]string{`AssetType`: `character varying`, `Clicks`: `integer`, `CreatedAt`: `timestamp without time zone`, `Date`: `timestamp without time zone`, `ID`: `integer`, `LangID`: `character`, `Lastuser`: `character varying`, `Name`: `character varying`, `PlaytimeSecs`: `integer`, `Secure`: `integer`, `ServernameID`: `character varying`, `Size`: `integer`, `Status`: `character varying`, `UpdatedAt`: `timestamp without time zone`, `UserID`: `integer`}
	_                = bytes.MinRead
)

func testFileAssetsUpdate(t *testing.T) {
	t.Parallel()

	if len(fileAssetColumns) == len(fileAssetPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	fileAsset := &FileAsset{}
	if err = randomize.Struct(seed, fileAsset, fileAssetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAsset.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := FileAssets(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, fileAsset, fileAssetDBTypes, true, fileAssetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	if err = fileAsset.Update(tx); err != nil {
		t.Error(err)
	}
}

func testFileAssetsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(fileAssetColumns) == len(fileAssetPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	fileAsset := &FileAsset{}
	if err = randomize.Struct(seed, fileAsset, fileAssetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAsset.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := FileAssets(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, fileAsset, fileAssetDBTypes, true, fileAssetPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(fileAssetColumns, fileAssetPrimaryKeyColumns) {
		fields = fileAssetColumns
	} else {
		fields = strmangle.SetComplement(
			fileAssetColumns,
			fileAssetPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(fileAsset))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := FileAssetSlice{fileAsset}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testFileAssetsUpsert(t *testing.T) {
	t.Parallel()

	if len(fileAssetColumns) == len(fileAssetPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	fileAsset := FileAsset{}
	if err = randomize.Struct(seed, &fileAsset, fileAssetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileAsset.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert FileAsset: %s", err)
	}

	count, err := FileAssets(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &fileAsset, fileAssetDBTypes, false, fileAssetPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize FileAsset struct: %s", err)
	}

	if err = fileAsset.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert FileAsset: %s", err)
	}

	count, err = FileAssets(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
