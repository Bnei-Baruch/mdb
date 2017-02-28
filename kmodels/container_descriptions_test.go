package kmodels

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testContainerDescriptions(t *testing.T) {
	t.Parallel()

	query := ContainerDescriptions(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testContainerDescriptionsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescription := &ContainerDescription{}
	if err = randomize.Struct(seed, containerDescription, containerDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = containerDescription.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := ContainerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContainerDescriptionsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescription := &ContainerDescription{}
	if err = randomize.Struct(seed, containerDescription, containerDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = ContainerDescriptions(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := ContainerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContainerDescriptionsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescription := &ContainerDescription{}
	if err = randomize.Struct(seed, containerDescription, containerDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ContainerDescriptionSlice{containerDescription}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := ContainerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testContainerDescriptionsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescription := &ContainerDescription{}
	if err = randomize.Struct(seed, containerDescription, containerDescriptionDBTypes, true, containerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := ContainerDescriptionExists(tx, containerDescription.ID)
	if err != nil {
		t.Errorf("Unable to check if ContainerDescription exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ContainerDescriptionExistsG to return true, but got false.")
	}
}
func testContainerDescriptionsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescription := &ContainerDescription{}
	if err = randomize.Struct(seed, containerDescription, containerDescriptionDBTypes, true, containerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	containerDescriptionFound, err := FindContainerDescription(tx, containerDescription.ID)
	if err != nil {
		t.Error(err)
	}

	if containerDescriptionFound == nil {
		t.Error("want a record, got nil")
	}
}
func testContainerDescriptionsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescription := &ContainerDescription{}
	if err = randomize.Struct(seed, containerDescription, containerDescriptionDBTypes, true, containerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = ContainerDescriptions(tx).Bind(containerDescription); err != nil {
		t.Error(err)
	}
}

func testContainerDescriptionsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescription := &ContainerDescription{}
	if err = randomize.Struct(seed, containerDescription, containerDescriptionDBTypes, true, containerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := ContainerDescriptions(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testContainerDescriptionsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescriptionOne := &ContainerDescription{}
	containerDescriptionTwo := &ContainerDescription{}
	if err = randomize.Struct(seed, containerDescriptionOne, containerDescriptionDBTypes, false, containerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}
	if err = randomize.Struct(seed, containerDescriptionTwo, containerDescriptionDBTypes, false, containerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = containerDescriptionTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := ContainerDescriptions(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testContainerDescriptionsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	containerDescriptionOne := &ContainerDescription{}
	containerDescriptionTwo := &ContainerDescription{}
	if err = randomize.Struct(seed, containerDescriptionOne, containerDescriptionDBTypes, false, containerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}
	if err = randomize.Struct(seed, containerDescriptionTwo, containerDescriptionDBTypes, false, containerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescriptionOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = containerDescriptionTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContainerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testContainerDescriptionsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescription := &ContainerDescription{}
	if err = randomize.Struct(seed, containerDescription, containerDescriptionDBTypes, true, containerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContainerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testContainerDescriptionsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescription := &ContainerDescription{}
	if err = randomize.Struct(seed, containerDescription, containerDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescription.Insert(tx, containerDescriptionColumns...); err != nil {
		t.Error(err)
	}

	count, err := ContainerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testContainerDescriptionToOneLanguageUsingLang(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local ContainerDescription
	var foreign Language

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, containerDescriptionDBTypes, true, containerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
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

	slice := ContainerDescriptionSlice{&local}
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

func testContainerDescriptionToOneContainerUsingContainer(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local ContainerDescription
	var foreign Container

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, containerDescriptionDBTypes, true, containerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, containerDBTypes, true, containerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Container struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.ContainerID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Container(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ContainerDescriptionSlice{&local}
	if err = local.L.LoadContainer(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Container == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Container = nil
	if err = local.L.LoadContainer(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Container == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testContainerDescriptionToOneSetOpLanguageUsingLang(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContainerDescription
	var b, c Language

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDescriptionDBTypes, false, strmangle.SetComplement(containerDescriptionPrimaryKeyColumns, containerDescriptionColumnsWithoutDefault)...); err != nil {
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

		if x.R.LangContainerDescriptions[0] != &a {
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

func testContainerDescriptionToOneRemoveOpLanguageUsingLang(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContainerDescription
	var b Language

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDescriptionDBTypes, false, strmangle.SetComplement(containerDescriptionPrimaryKeyColumns, containerDescriptionColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.LangContainerDescriptions) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testContainerDescriptionToOneSetOpContainerUsingContainer(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContainerDescription
	var b, c Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, containerDescriptionDBTypes, false, strmangle.SetComplement(containerDescriptionPrimaryKeyColumns, containerDescriptionColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, containerDBTypes, false, strmangle.SetComplement(containerPrimaryKeyColumns, containerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Container{&b, &c} {
		err = a.SetContainer(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Container != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.ContainerDescriptions[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.ContainerID != x.ID {
			t.Error("foreign key was wrong value", a.ContainerID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.ContainerID))
		reflect.Indirect(reflect.ValueOf(&a.ContainerID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.ContainerID != x.ID {
			t.Error("foreign key was wrong value", a.ContainerID, x.ID)
		}
	}
}
func testContainerDescriptionsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescription := &ContainerDescription{}
	if err = randomize.Struct(seed, containerDescription, containerDescriptionDBTypes, true, containerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = containerDescription.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testContainerDescriptionsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescription := &ContainerDescription{}
	if err = randomize.Struct(seed, containerDescription, containerDescriptionDBTypes, true, containerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ContainerDescriptionSlice{containerDescription}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testContainerDescriptionsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	containerDescription := &ContainerDescription{}
	if err = randomize.Struct(seed, containerDescription, containerDescriptionDBTypes, true, containerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := ContainerDescriptions(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	containerDescriptionDBTypes = map[string]string{`ContainerDesc`: `character varying`, `ContainerID`: `integer`, `CreatedAt`: `timestamp without time zone`, `Descr`: `text`, `ID`: `integer`, `LangID`: `character`, `UpdatedAt`: `timestamp without time zone`}
	_                           = bytes.MinRead
)

func testContainerDescriptionsUpdate(t *testing.T) {
	t.Parallel()

	if len(containerDescriptionColumns) == len(containerDescriptionPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	containerDescription := &ContainerDescription{}
	if err = randomize.Struct(seed, containerDescription, containerDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContainerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, containerDescription, containerDescriptionDBTypes, true, containerDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	if err = containerDescription.Update(tx); err != nil {
		t.Error(err)
	}
}

func testContainerDescriptionsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(containerDescriptionColumns) == len(containerDescriptionPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	containerDescription := &ContainerDescription{}
	if err = randomize.Struct(seed, containerDescription, containerDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContainerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, containerDescription, containerDescriptionDBTypes, true, containerDescriptionPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(containerDescriptionColumns, containerDescriptionPrimaryKeyColumns) {
		fields = containerDescriptionColumns
	} else {
		fields = strmangle.SetComplement(
			containerDescriptionColumns,
			containerDescriptionPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(containerDescription))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := ContainerDescriptionSlice{containerDescription}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testContainerDescriptionsUpsert(t *testing.T) {
	t.Parallel()

	if len(containerDescriptionColumns) == len(containerDescriptionPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	containerDescription := ContainerDescription{}
	if err = randomize.Struct(seed, &containerDescription, containerDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = containerDescription.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert ContainerDescription: %s", err)
	}

	count, err := ContainerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &containerDescription, containerDescriptionDBTypes, false, containerDescriptionPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ContainerDescription struct: %s", err)
	}

	if err = containerDescription.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert ContainerDescription: %s", err)
	}

	count, err = ContainerDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
