package kmodels

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testLanguages(t *testing.T) {
	t.Parallel()

	query := Languages(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testLanguagesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	language := &Language{}
	if err = randomize.Struct(seed, language, languageDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = language.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = language.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Languages(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testLanguagesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	language := &Language{}
	if err = randomize.Struct(seed, language, languageDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = language.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Languages(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Languages(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testLanguagesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	language := &Language{}
	if err = randomize.Struct(seed, language, languageDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = language.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := LanguageSlice{language}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Languages(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testLanguagesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	language := &Language{}
	if err = randomize.Struct(seed, language, languageDBTypes, true, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = language.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := LanguageExists(tx, language.ID)
	if err != nil {
		t.Errorf("Unable to check if Language exists: %s", err)
	}
	if !e {
		t.Errorf("Expected LanguageExistsG to return true, but got false.")
	}
}
func testLanguagesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	language := &Language{}
	if err = randomize.Struct(seed, language, languageDBTypes, true, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = language.Insert(tx); err != nil {
		t.Error(err)
	}

	languageFound, err := FindLanguage(tx, language.ID)
	if err != nil {
		t.Error(err)
	}

	if languageFound == nil {
		t.Error("want a record, got nil")
	}
}
func testLanguagesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	language := &Language{}
	if err = randomize.Struct(seed, language, languageDBTypes, true, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = language.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Languages(tx).Bind(language); err != nil {
		t.Error(err)
	}
}

func testLanguagesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	language := &Language{}
	if err = randomize.Struct(seed, language, languageDBTypes, true, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = language.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Languages(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testLanguagesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	languageOne := &Language{}
	languageTwo := &Language{}
	if err = randomize.Struct(seed, languageOne, languageDBTypes, false, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}
	if err = randomize.Struct(seed, languageTwo, languageDBTypes, false, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = languageOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = languageTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Languages(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testLanguagesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	languageOne := &Language{}
	languageTwo := &Language{}
	if err = randomize.Struct(seed, languageOne, languageDBTypes, false, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}
	if err = randomize.Struct(seed, languageTwo, languageDBTypes, false, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = languageOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = languageTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Languages(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testLanguagesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	language := &Language{}
	if err = randomize.Struct(seed, language, languageDBTypes, true, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = language.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Languages(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testLanguagesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	language := &Language{}
	if err = randomize.Struct(seed, language, languageDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = language.Insert(tx, languageColumns...); err != nil {
		t.Error(err)
	}

	count, err := Languages(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testLanguageToManyLangCatalogDescriptions(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Language
	var b, c CatalogDescription

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, languageDBTypes, true, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, catalogDescriptionDBTypes, false, catalogDescriptionColumnsWithDefault...)
	randomize.Struct(seed, &c, catalogDescriptionDBTypes, false, catalogDescriptionColumnsWithDefault...)
	a.Code3.Valid = true
	b.LangID.Valid = true
	c.LangID.Valid = true
	b.LangID.String = a.Code3.String
	c.LangID.String = a.Code3.String
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	catalogDescription, err := a.LangCatalogDescriptions(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range catalogDescription {
		if v.LangID.String == b.LangID.String {
			bFound = true
		}
		if v.LangID.String == c.LangID.String {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := LanguageSlice{&a}
	if err = a.L.LoadLangCatalogDescriptions(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.LangCatalogDescriptions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.LangCatalogDescriptions = nil
	if err = a.L.LoadLangCatalogDescriptions(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.LangCatalogDescriptions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", catalogDescription)
	}
}

func testLanguageToManyLangContainerDescriptions(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Language
	var b, c ContainerDescription

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, languageDBTypes, true, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, containerDescriptionDBTypes, false, containerDescriptionColumnsWithDefault...)
	randomize.Struct(seed, &c, containerDescriptionDBTypes, false, containerDescriptionColumnsWithDefault...)
	a.Code3.Valid = true
	b.LangID.Valid = true
	c.LangID.Valid = true
	b.LangID.String = a.Code3.String
	c.LangID.String = a.Code3.String
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	containerDescription, err := a.LangContainerDescriptions(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range containerDescription {
		if v.LangID.String == b.LangID.String {
			bFound = true
		}
		if v.LangID.String == c.LangID.String {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := LanguageSlice{&a}
	if err = a.L.LoadLangContainerDescriptions(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.LangContainerDescriptions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.LangContainerDescriptions = nil
	if err = a.L.LoadLangContainerDescriptions(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.LangContainerDescriptions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", containerDescription)
	}
}

func testLanguageToManyLangContainers(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Language
	var b, c Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, languageDBTypes, true, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, containerDBTypes, false, containerColumnsWithDefault...)
	randomize.Struct(seed, &c, containerDBTypes, false, containerColumnsWithDefault...)
	a.Code3.Valid = true
	b.LangID.Valid = true
	c.LangID.Valid = true
	b.LangID.String = a.Code3.String
	c.LangID.String = a.Code3.String
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	container, err := a.LangContainers(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range container {
		if v.LangID.String == b.LangID.String {
			bFound = true
		}
		if v.LangID.String == c.LangID.String {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := LanguageSlice{&a}
	if err = a.L.LoadLangContainers(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.LangContainers); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.LangContainers = nil
	if err = a.L.LoadLangContainers(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.LangContainers); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", container)
	}
}

func testLanguageToManyLangFileAssets(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Language
	var b, c FileAsset

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, languageDBTypes, true, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, fileAssetDBTypes, false, fileAssetColumnsWithDefault...)
	randomize.Struct(seed, &c, fileAssetDBTypes, false, fileAssetColumnsWithDefault...)
	a.Code3.Valid = true
	b.LangID.Valid = true
	c.LangID.Valid = true
	b.LangID.String = a.Code3.String
	c.LangID.String = a.Code3.String
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	fileAsset, err := a.LangFileAssets(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range fileAsset {
		if v.LangID.String == b.LangID.String {
			bFound = true
		}
		if v.LangID.String == c.LangID.String {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := LanguageSlice{&a}
	if err = a.L.LoadLangFileAssets(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.LangFileAssets); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.LangFileAssets = nil
	if err = a.L.LoadLangFileAssets(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.LangFileAssets); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", fileAsset)
	}
}

func testLanguageToManyAddOpLangCatalogDescriptions(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Language
	var b, c, d, e CatalogDescription

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, languageDBTypes, false, strmangle.SetComplement(languagePrimaryKeyColumns, languageColumnsWithoutDefault)...); err != nil {
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
		err = a.AddLangCatalogDescriptions(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.Code3.String != first.LangID.String {
			t.Error("foreign key was wrong value", a.Code3.String, first.LangID.String)
		}
		if a.Code3.String != second.LangID.String {
			t.Error("foreign key was wrong value", a.Code3.String, second.LangID.String)
		}

		if first.R.Lang != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Lang != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.LangCatalogDescriptions[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.LangCatalogDescriptions[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.LangCatalogDescriptions(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testLanguageToManySetOpLangCatalogDescriptions(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Language
	var b, c, d, e CatalogDescription

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, languageDBTypes, false, strmangle.SetComplement(languagePrimaryKeyColumns, languageColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*CatalogDescription{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, catalogDescriptionDBTypes, false, strmangle.SetComplement(catalogDescriptionPrimaryKeyColumns, catalogDescriptionColumnsWithoutDefault)...); err != nil {
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

	err = a.SetLangCatalogDescriptions(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.LangCatalogDescriptions(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetLangCatalogDescriptions(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.LangCatalogDescriptions(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.LangID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.LangID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.Code3.String != d.LangID.String {
		t.Error("foreign key was wrong value", a.Code3.String, d.LangID.String)
	}
	if a.Code3.String != e.LangID.String {
		t.Error("foreign key was wrong value", a.Code3.String, e.LangID.String)
	}

	if b.R.Lang != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Lang != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Lang != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Lang != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.LangCatalogDescriptions[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.LangCatalogDescriptions[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testLanguageToManyRemoveOpLangCatalogDescriptions(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Language
	var b, c, d, e CatalogDescription

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, languageDBTypes, false, strmangle.SetComplement(languagePrimaryKeyColumns, languageColumnsWithoutDefault)...); err != nil {
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

	err = a.AddLangCatalogDescriptions(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.LangCatalogDescriptions(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveLangCatalogDescriptions(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.LangCatalogDescriptions(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.LangID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.LangID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.Lang != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Lang != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Lang != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.Lang != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.LangCatalogDescriptions) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.LangCatalogDescriptions[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.LangCatalogDescriptions[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testLanguageToManyAddOpLangContainerDescriptions(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Language
	var b, c, d, e ContainerDescription

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, languageDBTypes, false, strmangle.SetComplement(languagePrimaryKeyColumns, languageColumnsWithoutDefault)...); err != nil {
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
		err = a.AddLangContainerDescriptions(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.Code3.String != first.LangID.String {
			t.Error("foreign key was wrong value", a.Code3.String, first.LangID.String)
		}
		if a.Code3.String != second.LangID.String {
			t.Error("foreign key was wrong value", a.Code3.String, second.LangID.String)
		}

		if first.R.Lang != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Lang != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.LangContainerDescriptions[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.LangContainerDescriptions[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.LangContainerDescriptions(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testLanguageToManySetOpLangContainerDescriptions(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Language
	var b, c, d, e ContainerDescription

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, languageDBTypes, false, strmangle.SetComplement(languagePrimaryKeyColumns, languageColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContainerDescription{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, containerDescriptionDBTypes, false, strmangle.SetComplement(containerDescriptionPrimaryKeyColumns, containerDescriptionColumnsWithoutDefault)...); err != nil {
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

	err = a.SetLangContainerDescriptions(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.LangContainerDescriptions(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetLangContainerDescriptions(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.LangContainerDescriptions(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.LangID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.LangID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.Code3.String != d.LangID.String {
		t.Error("foreign key was wrong value", a.Code3.String, d.LangID.String)
	}
	if a.Code3.String != e.LangID.String {
		t.Error("foreign key was wrong value", a.Code3.String, e.LangID.String)
	}

	if b.R.Lang != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Lang != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Lang != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Lang != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.LangContainerDescriptions[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.LangContainerDescriptions[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testLanguageToManyRemoveOpLangContainerDescriptions(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Language
	var b, c, d, e ContainerDescription

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, languageDBTypes, false, strmangle.SetComplement(languagePrimaryKeyColumns, languageColumnsWithoutDefault)...); err != nil {
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

	err = a.AddLangContainerDescriptions(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.LangContainerDescriptions(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveLangContainerDescriptions(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.LangContainerDescriptions(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.LangID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.LangID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.Lang != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Lang != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Lang != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.Lang != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.LangContainerDescriptions) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.LangContainerDescriptions[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.LangContainerDescriptions[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testLanguageToManyAddOpLangContainers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Language
	var b, c, d, e Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, languageDBTypes, false, strmangle.SetComplement(languagePrimaryKeyColumns, languageColumnsWithoutDefault)...); err != nil {
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
		err = a.AddLangContainers(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.Code3.String != first.LangID.String {
			t.Error("foreign key was wrong value", a.Code3.String, first.LangID.String)
		}
		if a.Code3.String != second.LangID.String {
			t.Error("foreign key was wrong value", a.Code3.String, second.LangID.String)
		}

		if first.R.Lang != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Lang != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.LangContainers[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.LangContainers[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.LangContainers(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testLanguageToManySetOpLangContainers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Language
	var b, c, d, e Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, languageDBTypes, false, strmangle.SetComplement(languagePrimaryKeyColumns, languageColumnsWithoutDefault)...); err != nil {
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

	err = a.SetLangContainers(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.LangContainers(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetLangContainers(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.LangContainers(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.LangID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.LangID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.Code3.String != d.LangID.String {
		t.Error("foreign key was wrong value", a.Code3.String, d.LangID.String)
	}
	if a.Code3.String != e.LangID.String {
		t.Error("foreign key was wrong value", a.Code3.String, e.LangID.String)
	}

	if b.R.Lang != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Lang != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Lang != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Lang != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.LangContainers[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.LangContainers[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testLanguageToManyRemoveOpLangContainers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Language
	var b, c, d, e Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, languageDBTypes, false, strmangle.SetComplement(languagePrimaryKeyColumns, languageColumnsWithoutDefault)...); err != nil {
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

	err = a.AddLangContainers(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.LangContainers(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveLangContainers(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.LangContainers(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.LangID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.LangID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.Lang != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Lang != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Lang != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.Lang != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.LangContainers) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.LangContainers[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.LangContainers[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testLanguageToManyAddOpLangFileAssets(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Language
	var b, c, d, e FileAsset

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, languageDBTypes, false, strmangle.SetComplement(languagePrimaryKeyColumns, languageColumnsWithoutDefault)...); err != nil {
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
		err = a.AddLangFileAssets(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.Code3.String != first.LangID.String {
			t.Error("foreign key was wrong value", a.Code3.String, first.LangID.String)
		}
		if a.Code3.String != second.LangID.String {
			t.Error("foreign key was wrong value", a.Code3.String, second.LangID.String)
		}

		if first.R.Lang != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Lang != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.LangFileAssets[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.LangFileAssets[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.LangFileAssets(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testLanguageToManySetOpLangFileAssets(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Language
	var b, c, d, e FileAsset

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, languageDBTypes, false, strmangle.SetComplement(languagePrimaryKeyColumns, languageColumnsWithoutDefault)...); err != nil {
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

	err = a.SetLangFileAssets(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.LangFileAssets(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetLangFileAssets(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.LangFileAssets(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.LangID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.LangID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.Code3.String != d.LangID.String {
		t.Error("foreign key was wrong value", a.Code3.String, d.LangID.String)
	}
	if a.Code3.String != e.LangID.String {
		t.Error("foreign key was wrong value", a.Code3.String, e.LangID.String)
	}

	if b.R.Lang != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Lang != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Lang != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Lang != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.LangFileAssets[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.LangFileAssets[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testLanguageToManyRemoveOpLangFileAssets(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Language
	var b, c, d, e FileAsset

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, languageDBTypes, false, strmangle.SetComplement(languagePrimaryKeyColumns, languageColumnsWithoutDefault)...); err != nil {
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

	err = a.AddLangFileAssets(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.LangFileAssets(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveLangFileAssets(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.LangFileAssets(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.LangID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.LangID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.Lang != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Lang != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Lang != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.Lang != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.LangFileAssets) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.LangFileAssets[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.LangFileAssets[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testLanguagesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	language := &Language{}
	if err = randomize.Struct(seed, language, languageDBTypes, true, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = language.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = language.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testLanguagesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	language := &Language{}
	if err = randomize.Struct(seed, language, languageDBTypes, true, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = language.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := LanguageSlice{language}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testLanguagesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	language := &Language{}
	if err = randomize.Struct(seed, language, languageDBTypes, true, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = language.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Languages(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	languageDBTypes = map[string]string{`Code3`: `character varying`, `ID`: `integer`, `Language`: `character varying`, `Locale`: `character varying`}
	_               = bytes.MinRead
)

func testLanguagesUpdate(t *testing.T) {
	t.Parallel()

	if len(languageColumns) == len(languagePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	language := &Language{}
	if err = randomize.Struct(seed, language, languageDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = language.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Languages(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, language, languageDBTypes, true, languageColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	if err = language.Update(tx); err != nil {
		t.Error(err)
	}
}

func testLanguagesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(languageColumns) == len(languagePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	language := &Language{}
	if err = randomize.Struct(seed, language, languageDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = language.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Languages(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, language, languageDBTypes, true, languagePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(languageColumns, languagePrimaryKeyColumns) {
		fields = languageColumns
	} else {
		fields = strmangle.SetComplement(
			languageColumns,
			languagePrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(language))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := LanguageSlice{language}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testLanguagesUpsert(t *testing.T) {
	t.Parallel()

	if len(languageColumns) == len(languagePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	language := Language{}
	if err = randomize.Struct(seed, &language, languageDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = language.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Language: %s", err)
	}

	count, err := Languages(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &language, languageDBTypes, false, languagePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Language struct: %s", err)
	}

	if err = language.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Language: %s", err)
	}

	count, err = Languages(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
