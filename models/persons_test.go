package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testPersons(t *testing.T) {
	t.Parallel()

	query := Persons(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testPersonsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	person := &Person{}
	if err = randomize.Struct(seed, person, personDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = person.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = person.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Persons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testPersonsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	person := &Person{}
	if err = randomize.Struct(seed, person, personDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = person.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Persons(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Persons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testPersonsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	person := &Person{}
	if err = randomize.Struct(seed, person, personDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = person.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := PersonSlice{person}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Persons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testPersonsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	person := &Person{}
	if err = randomize.Struct(seed, person, personDBTypes, true, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = person.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := PersonExists(tx, person.ID)
	if err != nil {
		t.Errorf("Unable to check if Person exists: %s", err)
	}
	if !e {
		t.Errorf("Expected PersonExistsG to return true, but got false.")
	}
}
func testPersonsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	person := &Person{}
	if err = randomize.Struct(seed, person, personDBTypes, true, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = person.Insert(tx); err != nil {
		t.Error(err)
	}

	personFound, err := FindPerson(tx, person.ID)
	if err != nil {
		t.Error(err)
	}

	if personFound == nil {
		t.Error("want a record, got nil")
	}
}
func testPersonsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	person := &Person{}
	if err = randomize.Struct(seed, person, personDBTypes, true, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = person.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Persons(tx).Bind(person); err != nil {
		t.Error(err)
	}
}

func testPersonsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	person := &Person{}
	if err = randomize.Struct(seed, person, personDBTypes, true, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = person.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Persons(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testPersonsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	personOne := &Person{}
	personTwo := &Person{}
	if err = randomize.Struct(seed, personOne, personDBTypes, false, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}
	if err = randomize.Struct(seed, personTwo, personDBTypes, false, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = personOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = personTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Persons(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testPersonsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	personOne := &Person{}
	personTwo := &Person{}
	if err = randomize.Struct(seed, personOne, personDBTypes, false, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}
	if err = randomize.Struct(seed, personTwo, personDBTypes, false, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = personOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = personTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Persons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testPersonsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	person := &Person{}
	if err = randomize.Struct(seed, person, personDBTypes, true, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = person.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Persons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testPersonsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	person := &Person{}
	if err = randomize.Struct(seed, person, personDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = person.Insert(tx, personColumns...); err != nil {
		t.Error(err)
	}

	count, err := Persons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testPersonToManyContentUnitsPersons(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Person
	var b, c ContentUnitsPerson

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, personDBTypes, true, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, contentUnitsPersonDBTypes, false, contentUnitsPersonColumnsWithDefault...)
	randomize.Struct(seed, &c, contentUnitsPersonDBTypes, false, contentUnitsPersonColumnsWithDefault...)

	b.PersonID = a.ID
	c.PersonID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	contentUnitsPerson, err := a.ContentUnitsPersons(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range contentUnitsPerson {
		if v.PersonID == b.PersonID {
			bFound = true
		}
		if v.PersonID == c.PersonID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := PersonSlice{&a}
	if err = a.L.LoadContentUnitsPersons(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ContentUnitsPersons); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.ContentUnitsPersons = nil
	if err = a.L.LoadContentUnitsPersons(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ContentUnitsPersons); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", contentUnitsPerson)
	}
}

func testPersonToManyAddOpContentUnitsPersons(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Person
	var b, c, d, e ContentUnitsPerson

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, personDBTypes, false, strmangle.SetComplement(personPrimaryKeyColumns, personColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentUnitsPerson{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentUnitsPersonDBTypes, false, strmangle.SetComplement(contentUnitsPersonPrimaryKeyColumns, contentUnitsPersonColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*ContentUnitsPerson{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddContentUnitsPersons(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.PersonID {
			t.Error("foreign key was wrong value", a.ID, first.PersonID)
		}
		if a.ID != second.PersonID {
			t.Error("foreign key was wrong value", a.ID, second.PersonID)
		}

		if first.R.Person != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Person != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.ContentUnitsPersons[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.ContentUnitsPersons[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.ContentUnitsPersons(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testPersonToOneStringTranslationUsingDescription(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Person
	var foreign StringTranslation

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, personDBTypes, true, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	local.DescriptionID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.DescriptionID.Int64 = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Description(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := PersonSlice{&local}
	if err = local.L.LoadDescription(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Description == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Description = nil
	if err = local.L.LoadDescription(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Description == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testPersonToOneStringTranslationUsingName(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Person
	var foreign StringTranslation

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, personDBTypes, true, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.NameID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Name(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := PersonSlice{&local}
	if err = local.L.LoadName(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Name == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Name = nil
	if err = local.L.LoadName(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Name == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testPersonToOneSetOpStringTranslationUsingDescription(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Person
	var b, c StringTranslation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, personDBTypes, false, strmangle.SetComplement(personPrimaryKeyColumns, personColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*StringTranslation{&b, &c} {
		err = a.SetDescription(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Description != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.DescriptionPersons[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.DescriptionID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.DescriptionID.Int64)
		}

		zero := reflect.Zero(reflect.TypeOf(a.DescriptionID.Int64))
		reflect.Indirect(reflect.ValueOf(&a.DescriptionID.Int64)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.DescriptionID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.DescriptionID.Int64, x.ID)
		}
	}
}

func testPersonToOneRemoveOpStringTranslationUsingDescription(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Person
	var b StringTranslation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, personDBTypes, false, strmangle.SetComplement(personPrimaryKeyColumns, personColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.SetDescription(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveDescription(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.Description(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.Description != nil {
		t.Error("R struct entry should be nil")
	}

	if a.DescriptionID.Valid {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.DescriptionPersons) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testPersonToOneSetOpStringTranslationUsingName(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Person
	var b, c StringTranslation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, personDBTypes, false, strmangle.SetComplement(personPrimaryKeyColumns, personColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*StringTranslation{&b, &c} {
		err = a.SetName(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Name != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.NamePersons[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.NameID != x.ID {
			t.Error("foreign key was wrong value", a.NameID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.NameID))
		reflect.Indirect(reflect.ValueOf(&a.NameID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.NameID != x.ID {
			t.Error("foreign key was wrong value", a.NameID, x.ID)
		}
	}
}
func testPersonsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	person := &Person{}
	if err = randomize.Struct(seed, person, personDBTypes, true, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = person.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = person.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testPersonsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	person := &Person{}
	if err = randomize.Struct(seed, person, personDBTypes, true, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = person.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := PersonSlice{person}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testPersonsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	person := &Person{}
	if err = randomize.Struct(seed, person, personDBTypes, true, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = person.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Persons(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	personDBTypes = map[string]string{`DescriptionID`: `bigint`, `ID`: `bigint`, `NameID`: `bigint`, `UID`: `character`}
	_             = bytes.MinRead
)

func testPersonsUpdate(t *testing.T) {
	t.Parallel()

	if len(personColumns) == len(personPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	person := &Person{}
	if err = randomize.Struct(seed, person, personDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = person.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Persons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, person, personDBTypes, true, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	if err = person.Update(tx); err != nil {
		t.Error(err)
	}
}

func testPersonsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(personColumns) == len(personPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	person := &Person{}
	if err = randomize.Struct(seed, person, personDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = person.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Persons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, person, personDBTypes, true, personPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(personColumns, personPrimaryKeyColumns) {
		fields = personColumns
	} else {
		fields = strmangle.SetComplement(
			personColumns,
			personPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(person))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := PersonSlice{person}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testPersonsUpsert(t *testing.T) {
	t.Parallel()

	if len(personColumns) == len(personPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	person := Person{}
	if err = randomize.Struct(seed, &person, personDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = person.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Person: %s", err)
	}

	count, err := Persons(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &person, personDBTypes, false, personPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	if err = person.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Person: %s", err)
	}

	count, err = Persons(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
