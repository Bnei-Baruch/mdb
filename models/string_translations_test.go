package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testStringTranslations(t *testing.T) {
	t.Parallel()

	query := StringTranslations(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testStringTranslationsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	stringTranslation := &StringTranslation{}
	if err = randomize.Struct(seed, stringTranslation, stringTranslationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslation.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = stringTranslation.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := StringTranslations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testStringTranslationsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	stringTranslation := &StringTranslation{}
	if err = randomize.Struct(seed, stringTranslation, stringTranslationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslation.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = StringTranslations(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := StringTranslations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testStringTranslationsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	stringTranslation := &StringTranslation{}
	if err = randomize.Struct(seed, stringTranslation, stringTranslationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslation.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := StringTranslationSlice{stringTranslation}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := StringTranslations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testStringTranslationsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	stringTranslation := &StringTranslation{}
	if err = randomize.Struct(seed, stringTranslation, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslation.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := StringTranslationExists(tx, stringTranslation.ID, stringTranslation.Language)
	if err != nil {
		t.Errorf("Unable to check if StringTranslation exists: %s", err)
	}
	if !e {
		t.Errorf("Expected StringTranslationExistsG to return true, but got false.")
	}
}
func testStringTranslationsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	stringTranslation := &StringTranslation{}
	if err = randomize.Struct(seed, stringTranslation, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslation.Insert(tx); err != nil {
		t.Error(err)
	}

	stringTranslationFound, err := FindStringTranslation(tx, stringTranslation.ID, stringTranslation.Language)
	if err != nil {
		t.Error(err)
	}

	if stringTranslationFound == nil {
		t.Error("want a record, got nil")
	}
}
func testStringTranslationsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	stringTranslation := &StringTranslation{}
	if err = randomize.Struct(seed, stringTranslation, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslation.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = StringTranslations(tx).Bind(stringTranslation); err != nil {
		t.Error(err)
	}
}

func testStringTranslationsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	stringTranslation := &StringTranslation{}
	if err = randomize.Struct(seed, stringTranslation, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslation.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := StringTranslations(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testStringTranslationsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	stringTranslationOne := &StringTranslation{}
	stringTranslationTwo := &StringTranslation{}
	if err = randomize.Struct(seed, stringTranslationOne, stringTranslationDBTypes, false, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}
	if err = randomize.Struct(seed, stringTranslationTwo, stringTranslationDBTypes, false, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslationOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = stringTranslationTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := StringTranslations(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testStringTranslationsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	stringTranslationOne := &StringTranslation{}
	stringTranslationTwo := &StringTranslation{}
	if err = randomize.Struct(seed, stringTranslationOne, stringTranslationDBTypes, false, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}
	if err = randomize.Struct(seed, stringTranslationTwo, stringTranslationDBTypes, false, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslationOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = stringTranslationTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := StringTranslations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testStringTranslationsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	stringTranslation := &StringTranslation{}
	if err = randomize.Struct(seed, stringTranslation, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslation.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := StringTranslations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testStringTranslationsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	stringTranslation := &StringTranslation{}
	if err = randomize.Struct(seed, stringTranslation, stringTranslationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslation.Insert(tx, stringTranslationColumns...); err != nil {
		t.Error(err)
	}

	count, err := StringTranslations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testStringTranslationToManyDescriptionPersons(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c Person

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, personDBTypes, false, personColumnsWithDefault...)
	randomize.Struct(seed, &c, personDBTypes, false, personColumnsWithDefault...)

	b.DescriptionID.Valid = true
	c.DescriptionID.Valid = true
	b.DescriptionID.Int64 = a.ID
	c.DescriptionID.Int64 = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	person, err := a.DescriptionPersons(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range person {
		if v.DescriptionID.Int64 == b.DescriptionID.Int64 {
			bFound = true
		}
		if v.DescriptionID.Int64 == c.DescriptionID.Int64 {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := StringTranslationSlice{&a}
	if err = a.L.LoadDescriptionPersons(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.DescriptionPersons); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.DescriptionPersons = nil
	if err = a.L.LoadDescriptionPersons(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.DescriptionPersons); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", person)
	}
}

func testStringTranslationToManyNamePersons(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c Person

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, personDBTypes, false, personColumnsWithDefault...)
	randomize.Struct(seed, &c, personDBTypes, false, personColumnsWithDefault...)

	b.NameID = a.ID
	c.NameID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	person, err := a.NamePersons(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range person {
		if v.NameID == b.NameID {
			bFound = true
		}
		if v.NameID == c.NameID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := StringTranslationSlice{&a}
	if err = a.L.LoadNamePersons(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.NamePersons); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.NamePersons = nil
	if err = a.L.LoadNamePersons(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.NamePersons); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", person)
	}
}

func testStringTranslationToManyLabelTags(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c Tag

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, tagDBTypes, false, tagColumnsWithDefault...)
	randomize.Struct(seed, &c, tagDBTypes, false, tagColumnsWithDefault...)

	b.LabelID = a.ID
	c.LabelID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	tag, err := a.LabelTags(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range tag {
		if v.LabelID == b.LabelID {
			bFound = true
		}
		if v.LabelID == c.LabelID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := StringTranslationSlice{&a}
	if err = a.L.LoadLabelTags(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.LabelTags); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.LabelTags = nil
	if err = a.L.LoadLabelTags(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.LabelTags); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", tag)
	}
}

func testStringTranslationToManyDescriptionContentUnits(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, contentUnitDBTypes, false, contentUnitColumnsWithDefault...)
	randomize.Struct(seed, &c, contentUnitDBTypes, false, contentUnitColumnsWithDefault...)

	b.DescriptionID.Valid = true
	c.DescriptionID.Valid = true
	b.DescriptionID.Int64 = a.ID
	c.DescriptionID.Int64 = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	contentUnit, err := a.DescriptionContentUnits(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range contentUnit {
		if v.DescriptionID.Int64 == b.DescriptionID.Int64 {
			bFound = true
		}
		if v.DescriptionID.Int64 == c.DescriptionID.Int64 {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := StringTranslationSlice{&a}
	if err = a.L.LoadDescriptionContentUnits(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.DescriptionContentUnits); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.DescriptionContentUnits = nil
	if err = a.L.LoadDescriptionContentUnits(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.DescriptionContentUnits); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", contentUnit)
	}
}

func testStringTranslationToManyNameContentUnits(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, contentUnitDBTypes, false, contentUnitColumnsWithDefault...)
	randomize.Struct(seed, &c, contentUnitDBTypes, false, contentUnitColumnsWithDefault...)

	b.NameID = a.ID
	c.NameID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	contentUnit, err := a.NameContentUnits(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range contentUnit {
		if v.NameID == b.NameID {
			bFound = true
		}
		if v.NameID == c.NameID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := StringTranslationSlice{&a}
	if err = a.L.LoadNameContentUnits(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.NameContentUnits); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.NameContentUnits = nil
	if err = a.L.LoadNameContentUnits(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.NameContentUnits); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", contentUnit)
	}
}

func testStringTranslationToManyDescriptionCollections(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c Collection

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, collectionDBTypes, false, collectionColumnsWithDefault...)
	randomize.Struct(seed, &c, collectionDBTypes, false, collectionColumnsWithDefault...)

	b.DescriptionID.Valid = true
	c.DescriptionID.Valid = true
	b.DescriptionID.Int64 = a.ID
	c.DescriptionID.Int64 = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	collection, err := a.DescriptionCollections(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range collection {
		if v.DescriptionID.Int64 == b.DescriptionID.Int64 {
			bFound = true
		}
		if v.DescriptionID.Int64 == c.DescriptionID.Int64 {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := StringTranslationSlice{&a}
	if err = a.L.LoadDescriptionCollections(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.DescriptionCollections); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.DescriptionCollections = nil
	if err = a.L.LoadDescriptionCollections(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.DescriptionCollections); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", collection)
	}
}

func testStringTranslationToManyNameCollections(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c Collection

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, collectionDBTypes, false, collectionColumnsWithDefault...)
	randomize.Struct(seed, &c, collectionDBTypes, false, collectionColumnsWithDefault...)

	b.NameID = a.ID
	c.NameID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	collection, err := a.NameCollections(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range collection {
		if v.NameID == b.NameID {
			bFound = true
		}
		if v.NameID == c.NameID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := StringTranslationSlice{&a}
	if err = a.L.LoadNameCollections(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.NameCollections); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.NameCollections = nil
	if err = a.L.LoadNameCollections(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.NameCollections); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", collection)
	}
}

func testStringTranslationToManyDescriptionContentRoles(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c ContentRole

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, contentRoleDBTypes, false, contentRoleColumnsWithDefault...)
	randomize.Struct(seed, &c, contentRoleDBTypes, false, contentRoleColumnsWithDefault...)

	b.DescriptionID.Valid = true
	c.DescriptionID.Valid = true
	b.DescriptionID.Int64 = a.ID
	c.DescriptionID.Int64 = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	contentRole, err := a.DescriptionContentRoles(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range contentRole {
		if v.DescriptionID.Int64 == b.DescriptionID.Int64 {
			bFound = true
		}
		if v.DescriptionID.Int64 == c.DescriptionID.Int64 {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := StringTranslationSlice{&a}
	if err = a.L.LoadDescriptionContentRoles(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.DescriptionContentRoles); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.DescriptionContentRoles = nil
	if err = a.L.LoadDescriptionContentRoles(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.DescriptionContentRoles); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", contentRole)
	}
}

func testStringTranslationToManyNameContentRoles(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c ContentRole

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, contentRoleDBTypes, false, contentRoleColumnsWithDefault...)
	randomize.Struct(seed, &c, contentRoleDBTypes, false, contentRoleColumnsWithDefault...)

	b.NameID = a.ID
	c.NameID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	contentRole, err := a.NameContentRoles(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range contentRole {
		if v.NameID == b.NameID {
			bFound = true
		}
		if v.NameID == c.NameID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := StringTranslationSlice{&a}
	if err = a.L.LoadNameContentRoles(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.NameContentRoles); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.NameContentRoles = nil
	if err = a.L.LoadNameContentRoles(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.NameContentRoles); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", contentRole)
	}
}

func testStringTranslationToManyAddOpDescriptionPersons(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e Person

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Person{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, personDBTypes, false, strmangle.SetComplement(personPrimaryKeyColumns, personColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Person{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddDescriptionPersons(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.DescriptionID.Int64 {
			t.Error("foreign key was wrong value", a.ID, first.DescriptionID.Int64)
		}
		if a.ID != second.DescriptionID.Int64 {
			t.Error("foreign key was wrong value", a.ID, second.DescriptionID.Int64)
		}

		if first.R.Description != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Description != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.DescriptionPersons[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.DescriptionPersons[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.DescriptionPersons(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testStringTranslationToManySetOpDescriptionPersons(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e Person

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Person{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, personDBTypes, false, strmangle.SetComplement(personPrimaryKeyColumns, personColumnsWithoutDefault)...); err != nil {
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

	err = a.SetDescriptionPersons(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.DescriptionPersons(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetDescriptionPersons(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.DescriptionPersons(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.DescriptionID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.DescriptionID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.ID != d.DescriptionID.Int64 {
		t.Error("foreign key was wrong value", a.ID, d.DescriptionID.Int64)
	}
	if a.ID != e.DescriptionID.Int64 {
		t.Error("foreign key was wrong value", a.ID, e.DescriptionID.Int64)
	}

	if b.R.Description != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Description != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Description != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Description != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.DescriptionPersons[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.DescriptionPersons[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testStringTranslationToManyRemoveOpDescriptionPersons(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e Person

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Person{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, personDBTypes, false, strmangle.SetComplement(personPrimaryKeyColumns, personColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddDescriptionPersons(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.DescriptionPersons(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveDescriptionPersons(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.DescriptionPersons(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.DescriptionID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.DescriptionID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.Description != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Description != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Description != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.Description != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.DescriptionPersons) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.DescriptionPersons[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.DescriptionPersons[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testStringTranslationToManyAddOpNamePersons(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e Person

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Person{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, personDBTypes, false, strmangle.SetComplement(personPrimaryKeyColumns, personColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Person{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddNamePersons(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.NameID {
			t.Error("foreign key was wrong value", a.ID, first.NameID)
		}
		if a.ID != second.NameID {
			t.Error("foreign key was wrong value", a.ID, second.NameID)
		}

		if first.R.Name != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Name != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.NamePersons[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.NamePersons[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.NamePersons(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testStringTranslationToManyAddOpLabelTags(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e Tag

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
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
		err = a.AddLabelTags(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.LabelID {
			t.Error("foreign key was wrong value", a.ID, first.LabelID)
		}
		if a.ID != second.LabelID {
			t.Error("foreign key was wrong value", a.ID, second.LabelID)
		}

		if first.R.Label != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Label != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.LabelTags[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.LabelTags[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.LabelTags(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testStringTranslationToManyAddOpDescriptionContentUnits(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
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
		err = a.AddDescriptionContentUnits(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.DescriptionID.Int64 {
			t.Error("foreign key was wrong value", a.ID, first.DescriptionID.Int64)
		}
		if a.ID != second.DescriptionID.Int64 {
			t.Error("foreign key was wrong value", a.ID, second.DescriptionID.Int64)
		}

		if first.R.Description != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Description != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.DescriptionContentUnits[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.DescriptionContentUnits[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.DescriptionContentUnits(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testStringTranslationToManySetOpDescriptionContentUnits(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
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

	err = a.SetDescriptionContentUnits(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.DescriptionContentUnits(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetDescriptionContentUnits(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.DescriptionContentUnits(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.DescriptionID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.DescriptionID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.ID != d.DescriptionID.Int64 {
		t.Error("foreign key was wrong value", a.ID, d.DescriptionID.Int64)
	}
	if a.ID != e.DescriptionID.Int64 {
		t.Error("foreign key was wrong value", a.ID, e.DescriptionID.Int64)
	}

	if b.R.Description != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Description != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Description != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Description != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.DescriptionContentUnits[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.DescriptionContentUnits[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testStringTranslationToManyRemoveOpDescriptionContentUnits(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
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

	err = a.AddDescriptionContentUnits(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.DescriptionContentUnits(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveDescriptionContentUnits(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.DescriptionContentUnits(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.DescriptionID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.DescriptionID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.Description != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Description != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Description != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.Description != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.DescriptionContentUnits) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.DescriptionContentUnits[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.DescriptionContentUnits[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testStringTranslationToManyAddOpNameContentUnits(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
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
		err = a.AddNameContentUnits(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.NameID {
			t.Error("foreign key was wrong value", a.ID, first.NameID)
		}
		if a.ID != second.NameID {
			t.Error("foreign key was wrong value", a.ID, second.NameID)
		}

		if first.R.Name != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Name != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.NameContentUnits[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.NameContentUnits[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.NameContentUnits(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testStringTranslationToManyAddOpDescriptionCollections(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e Collection

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Collection{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, collectionDBTypes, false, strmangle.SetComplement(collectionPrimaryKeyColumns, collectionColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Collection{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddDescriptionCollections(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.DescriptionID.Int64 {
			t.Error("foreign key was wrong value", a.ID, first.DescriptionID.Int64)
		}
		if a.ID != second.DescriptionID.Int64 {
			t.Error("foreign key was wrong value", a.ID, second.DescriptionID.Int64)
		}

		if first.R.Description != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Description != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.DescriptionCollections[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.DescriptionCollections[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.DescriptionCollections(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testStringTranslationToManySetOpDescriptionCollections(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e Collection

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Collection{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, collectionDBTypes, false, strmangle.SetComplement(collectionPrimaryKeyColumns, collectionColumnsWithoutDefault)...); err != nil {
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

	err = a.SetDescriptionCollections(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.DescriptionCollections(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetDescriptionCollections(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.DescriptionCollections(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.DescriptionID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.DescriptionID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.ID != d.DescriptionID.Int64 {
		t.Error("foreign key was wrong value", a.ID, d.DescriptionID.Int64)
	}
	if a.ID != e.DescriptionID.Int64 {
		t.Error("foreign key was wrong value", a.ID, e.DescriptionID.Int64)
	}

	if b.R.Description != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Description != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Description != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Description != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.DescriptionCollections[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.DescriptionCollections[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testStringTranslationToManyRemoveOpDescriptionCollections(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e Collection

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Collection{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, collectionDBTypes, false, strmangle.SetComplement(collectionPrimaryKeyColumns, collectionColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddDescriptionCollections(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.DescriptionCollections(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveDescriptionCollections(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.DescriptionCollections(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.DescriptionID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.DescriptionID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.Description != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Description != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Description != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.Description != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.DescriptionCollections) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.DescriptionCollections[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.DescriptionCollections[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testStringTranslationToManyAddOpNameCollections(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e Collection

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Collection{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, collectionDBTypes, false, strmangle.SetComplement(collectionPrimaryKeyColumns, collectionColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Collection{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddNameCollections(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.NameID {
			t.Error("foreign key was wrong value", a.ID, first.NameID)
		}
		if a.ID != second.NameID {
			t.Error("foreign key was wrong value", a.ID, second.NameID)
		}

		if first.R.Name != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Name != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.NameCollections[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.NameCollections[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.NameCollections(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testStringTranslationToManyAddOpDescriptionContentRoles(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e ContentRole

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentRole{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentRoleDBTypes, false, strmangle.SetComplement(contentRolePrimaryKeyColumns, contentRoleColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*ContentRole{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddDescriptionContentRoles(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.DescriptionID.Int64 {
			t.Error("foreign key was wrong value", a.ID, first.DescriptionID.Int64)
		}
		if a.ID != second.DescriptionID.Int64 {
			t.Error("foreign key was wrong value", a.ID, second.DescriptionID.Int64)
		}

		if first.R.Description != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Description != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.DescriptionContentRoles[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.DescriptionContentRoles[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.DescriptionContentRoles(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testStringTranslationToManySetOpDescriptionContentRoles(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e ContentRole

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentRole{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentRoleDBTypes, false, strmangle.SetComplement(contentRolePrimaryKeyColumns, contentRoleColumnsWithoutDefault)...); err != nil {
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

	err = a.SetDescriptionContentRoles(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.DescriptionContentRoles(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetDescriptionContentRoles(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.DescriptionContentRoles(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.DescriptionID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.DescriptionID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.ID != d.DescriptionID.Int64 {
		t.Error("foreign key was wrong value", a.ID, d.DescriptionID.Int64)
	}
	if a.ID != e.DescriptionID.Int64 {
		t.Error("foreign key was wrong value", a.ID, e.DescriptionID.Int64)
	}

	if b.R.Description != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Description != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Description != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Description != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.DescriptionContentRoles[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.DescriptionContentRoles[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testStringTranslationToManyRemoveOpDescriptionContentRoles(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e ContentRole

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentRole{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentRoleDBTypes, false, strmangle.SetComplement(contentRolePrimaryKeyColumns, contentRoleColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddDescriptionContentRoles(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.DescriptionContentRoles(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveDescriptionContentRoles(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.DescriptionContentRoles(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.DescriptionID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.DescriptionID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.Description != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Description != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Description != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.Description != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.DescriptionContentRoles) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.DescriptionContentRoles[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.DescriptionContentRoles[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testStringTranslationToManyAddOpNameContentRoles(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a StringTranslation
	var b, c, d, e ContentRole

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentRole{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentRoleDBTypes, false, strmangle.SetComplement(contentRolePrimaryKeyColumns, contentRoleColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*ContentRole{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddNameContentRoles(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.NameID {
			t.Error("foreign key was wrong value", a.ID, first.NameID)
		}
		if a.ID != second.NameID {
			t.Error("foreign key was wrong value", a.ID, second.NameID)
		}

		if first.R.Name != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Name != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.NameContentRoles[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.NameContentRoles[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.NameContentRoles(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testStringTranslationsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	stringTranslation := &StringTranslation{}
	if err = randomize.Struct(seed, stringTranslation, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslation.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = stringTranslation.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testStringTranslationsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	stringTranslation := &StringTranslation{}
	if err = randomize.Struct(seed, stringTranslation, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslation.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := StringTranslationSlice{stringTranslation}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testStringTranslationsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	stringTranslation := &StringTranslation{}
	if err = randomize.Struct(seed, stringTranslation, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslation.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := StringTranslations(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	stringTranslationDBTypes = map[string]string{`CreatedAt`: `timestamp with time zone`, `ID`: `bigint`, `Language`: `character`, `OriginalLanguage`: `character`, `Text`: `text`, `UserID`: `bigint`}
	_                        = bytes.MinRead
)

func testStringTranslationsUpdate(t *testing.T) {
	t.Parallel()

	if len(stringTranslationColumns) == len(stringTranslationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	stringTranslation := &StringTranslation{}
	if err = randomize.Struct(seed, stringTranslation, stringTranslationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslation.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := StringTranslations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, stringTranslation, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	if err = stringTranslation.Update(tx); err != nil {
		t.Error(err)
	}
}

func testStringTranslationsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(stringTranslationColumns) == len(stringTranslationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	stringTranslation := &StringTranslation{}
	if err = randomize.Struct(seed, stringTranslation, stringTranslationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslation.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := StringTranslations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, stringTranslation, stringTranslationDBTypes, true, stringTranslationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(stringTranslationColumns, stringTranslationPrimaryKeyColumns) {
		fields = stringTranslationColumns
	} else {
		fields = strmangle.SetComplement(
			stringTranslationColumns,
			stringTranslationPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(stringTranslation))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := StringTranslationSlice{stringTranslation}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testStringTranslationsUpsert(t *testing.T) {
	t.Parallel()

	if len(stringTranslationColumns) == len(stringTranslationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	stringTranslation := StringTranslation{}
	if err = randomize.Struct(seed, &stringTranslation, stringTranslationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = stringTranslation.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert StringTranslation: %s", err)
	}

	count, err := StringTranslations(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &stringTranslation, stringTranslationDBTypes, false, stringTranslationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	if err = stringTranslation.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert StringTranslation: %s", err)
	}

	count, err = StringTranslations(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
