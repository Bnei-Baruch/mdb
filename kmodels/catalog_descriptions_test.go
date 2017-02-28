package kmodels

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testCatalogDescriptions(t *testing.T) {
	t.Parallel()

	query := CatalogDescriptions(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testCatalogDescriptionsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalogDescription := &CatalogDescription{}
	if err = randomize.Struct(seed, catalogDescription, catalogDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = catalogDescription.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := CatalogDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCatalogDescriptionsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalogDescription := &CatalogDescription{}
	if err = randomize.Struct(seed, catalogDescription, catalogDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = CatalogDescriptions(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := CatalogDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCatalogDescriptionsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalogDescription := &CatalogDescription{}
	if err = randomize.Struct(seed, catalogDescription, catalogDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := CatalogDescriptionSlice{catalogDescription}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := CatalogDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testCatalogDescriptionsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalogDescription := &CatalogDescription{}
	if err = randomize.Struct(seed, catalogDescription, catalogDescriptionDBTypes, true, catalogDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := CatalogDescriptionExists(tx, catalogDescription.ID)
	if err != nil {
		t.Errorf("Unable to check if CatalogDescription exists: %s", err)
	}
	if !e {
		t.Errorf("Expected CatalogDescriptionExistsG to return true, but got false.")
	}
}
func testCatalogDescriptionsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalogDescription := &CatalogDescription{}
	if err = randomize.Struct(seed, catalogDescription, catalogDescriptionDBTypes, true, catalogDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	catalogDescriptionFound, err := FindCatalogDescription(tx, catalogDescription.ID)
	if err != nil {
		t.Error(err)
	}

	if catalogDescriptionFound == nil {
		t.Error("want a record, got nil")
	}
}
func testCatalogDescriptionsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalogDescription := &CatalogDescription{}
	if err = randomize.Struct(seed, catalogDescription, catalogDescriptionDBTypes, true, catalogDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = CatalogDescriptions(tx).Bind(catalogDescription); err != nil {
		t.Error(err)
	}
}

func testCatalogDescriptionsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalogDescription := &CatalogDescription{}
	if err = randomize.Struct(seed, catalogDescription, catalogDescriptionDBTypes, true, catalogDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := CatalogDescriptions(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testCatalogDescriptionsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalogDescriptionOne := &CatalogDescription{}
	catalogDescriptionTwo := &CatalogDescription{}
	if err = randomize.Struct(seed, catalogDescriptionOne, catalogDescriptionDBTypes, false, catalogDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}
	if err = randomize.Struct(seed, catalogDescriptionTwo, catalogDescriptionDBTypes, false, catalogDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescriptionOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = catalogDescriptionTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := CatalogDescriptions(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testCatalogDescriptionsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	catalogDescriptionOne := &CatalogDescription{}
	catalogDescriptionTwo := &CatalogDescription{}
	if err = randomize.Struct(seed, catalogDescriptionOne, catalogDescriptionDBTypes, false, catalogDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}
	if err = randomize.Struct(seed, catalogDescriptionTwo, catalogDescriptionDBTypes, false, catalogDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescriptionOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = catalogDescriptionTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := CatalogDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testCatalogDescriptionsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalogDescription := &CatalogDescription{}
	if err = randomize.Struct(seed, catalogDescription, catalogDescriptionDBTypes, true, catalogDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := CatalogDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testCatalogDescriptionsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalogDescription := &CatalogDescription{}
	if err = randomize.Struct(seed, catalogDescription, catalogDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescription.Insert(tx, catalogDescriptionColumns...); err != nil {
		t.Error(err)
	}

	count, err := CatalogDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testCatalogDescriptionToOneCatalogUsingCatalog(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local CatalogDescription
	var foreign Catalog

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, catalogDescriptionDBTypes, true, catalogDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, catalogDBTypes, true, catalogColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Catalog struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.CatalogID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Catalog(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := CatalogDescriptionSlice{&local}
	if err = local.L.LoadCatalog(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Catalog == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Catalog = nil
	if err = local.L.LoadCatalog(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Catalog == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testCatalogDescriptionToOneLanguageUsingLang(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local CatalogDescription
	var foreign Language

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, catalogDescriptionDBTypes, true, catalogDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
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

	slice := CatalogDescriptionSlice{&local}
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

func testCatalogDescriptionToOneSetOpCatalogUsingCatalog(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a CatalogDescription
	var b, c Catalog

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDescriptionDBTypes, false, strmangle.SetComplement(catalogDescriptionPrimaryKeyColumns, catalogDescriptionColumnsWithoutDefault)...); err != nil {
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
		err = a.SetCatalog(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Catalog != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.CatalogDescriptions[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.CatalogID != x.ID {
			t.Error("foreign key was wrong value", a.CatalogID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.CatalogID))
		reflect.Indirect(reflect.ValueOf(&a.CatalogID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.CatalogID != x.ID {
			t.Error("foreign key was wrong value", a.CatalogID, x.ID)
		}
	}
}
func testCatalogDescriptionToOneSetOpLanguageUsingLang(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a CatalogDescription
	var b, c Language

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDescriptionDBTypes, false, strmangle.SetComplement(catalogDescriptionPrimaryKeyColumns, catalogDescriptionColumnsWithoutDefault)...); err != nil {
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

		if x.R.LangCatalogDescriptions[0] != &a {
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

func testCatalogDescriptionToOneRemoveOpLanguageUsingLang(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a CatalogDescription
	var b Language

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, catalogDescriptionDBTypes, false, strmangle.SetComplement(catalogDescriptionPrimaryKeyColumns, catalogDescriptionColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.LangCatalogDescriptions) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testCatalogDescriptionsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalogDescription := &CatalogDescription{}
	if err = randomize.Struct(seed, catalogDescription, catalogDescriptionDBTypes, true, catalogDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = catalogDescription.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testCatalogDescriptionsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalogDescription := &CatalogDescription{}
	if err = randomize.Struct(seed, catalogDescription, catalogDescriptionDBTypes, true, catalogDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := CatalogDescriptionSlice{catalogDescription}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testCatalogDescriptionsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	catalogDescription := &CatalogDescription{}
	if err = randomize.Struct(seed, catalogDescription, catalogDescriptionDBTypes, true, catalogDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := CatalogDescriptions(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	catalogDescriptionDBTypes = map[string]string{`CatalogID`: `integer`, `CreatedAt`: `timestamp without time zone`, `ID`: `integer`, `LangID`: `character`, `Name`: `character varying`, `UpdatedAt`: `timestamp without time zone`}
	_                         = bytes.MinRead
)

func testCatalogDescriptionsUpdate(t *testing.T) {
	t.Parallel()

	if len(catalogDescriptionColumns) == len(catalogDescriptionPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	catalogDescription := &CatalogDescription{}
	if err = randomize.Struct(seed, catalogDescription, catalogDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := CatalogDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, catalogDescription, catalogDescriptionDBTypes, true, catalogDescriptionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	if err = catalogDescription.Update(tx); err != nil {
		t.Error(err)
	}
}

func testCatalogDescriptionsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(catalogDescriptionColumns) == len(catalogDescriptionPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	catalogDescription := &CatalogDescription{}
	if err = randomize.Struct(seed, catalogDescription, catalogDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescription.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := CatalogDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, catalogDescription, catalogDescriptionDBTypes, true, catalogDescriptionPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(catalogDescriptionColumns, catalogDescriptionPrimaryKeyColumns) {
		fields = catalogDescriptionColumns
	} else {
		fields = strmangle.SetComplement(
			catalogDescriptionColumns,
			catalogDescriptionPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(catalogDescription))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := CatalogDescriptionSlice{catalogDescription}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testCatalogDescriptionsUpsert(t *testing.T) {
	t.Parallel()

	if len(catalogDescriptionColumns) == len(catalogDescriptionPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	catalogDescription := CatalogDescription{}
	if err = randomize.Struct(seed, &catalogDescription, catalogDescriptionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = catalogDescription.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert CatalogDescription: %s", err)
	}

	count, err := CatalogDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &catalogDescription, catalogDescriptionDBTypes, false, catalogDescriptionPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize CatalogDescription struct: %s", err)
	}

	if err = catalogDescription.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert CatalogDescription: %s", err)
	}

	count, err = CatalogDescriptions(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
