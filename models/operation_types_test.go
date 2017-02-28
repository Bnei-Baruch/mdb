package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testOperationTypes(t *testing.T) {
	t.Parallel()

	query := OperationTypes(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testOperationTypesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operationType := &OperationType{}
	if err = randomize.Struct(seed, operationType, operationTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationType.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = operationType.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := OperationTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testOperationTypesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operationType := &OperationType{}
	if err = randomize.Struct(seed, operationType, operationTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationType.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = OperationTypes(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := OperationTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testOperationTypesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operationType := &OperationType{}
	if err = randomize.Struct(seed, operationType, operationTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationType.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := OperationTypeSlice{operationType}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := OperationTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testOperationTypesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operationType := &OperationType{}
	if err = randomize.Struct(seed, operationType, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationType.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := OperationTypeExists(tx, operationType.ID)
	if err != nil {
		t.Errorf("Unable to check if OperationType exists: %s", err)
	}
	if !e {
		t.Errorf("Expected OperationTypeExistsG to return true, but got false.")
	}
}
func testOperationTypesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operationType := &OperationType{}
	if err = randomize.Struct(seed, operationType, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationType.Insert(tx); err != nil {
		t.Error(err)
	}

	operationTypeFound, err := FindOperationType(tx, operationType.ID)
	if err != nil {
		t.Error(err)
	}

	if operationTypeFound == nil {
		t.Error("want a record, got nil")
	}
}
func testOperationTypesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operationType := &OperationType{}
	if err = randomize.Struct(seed, operationType, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationType.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = OperationTypes(tx).Bind(operationType); err != nil {
		t.Error(err)
	}
}

func testOperationTypesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operationType := &OperationType{}
	if err = randomize.Struct(seed, operationType, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationType.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := OperationTypes(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testOperationTypesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operationTypeOne := &OperationType{}
	operationTypeTwo := &OperationType{}
	if err = randomize.Struct(seed, operationTypeOne, operationTypeDBTypes, false, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}
	if err = randomize.Struct(seed, operationTypeTwo, operationTypeDBTypes, false, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationTypeOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = operationTypeTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := OperationTypes(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testOperationTypesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	operationTypeOne := &OperationType{}
	operationTypeTwo := &OperationType{}
	if err = randomize.Struct(seed, operationTypeOne, operationTypeDBTypes, false, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}
	if err = randomize.Struct(seed, operationTypeTwo, operationTypeDBTypes, false, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationTypeOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = operationTypeTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := OperationTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testOperationTypesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operationType := &OperationType{}
	if err = randomize.Struct(seed, operationType, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationType.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := OperationTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testOperationTypesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operationType := &OperationType{}
	if err = randomize.Struct(seed, operationType, operationTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationType.Insert(tx, operationTypeColumns...); err != nil {
		t.Error(err)
	}

	count, err := OperationTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testOperationTypeToManyTypeOperations(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a OperationType
	var b, c Operation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, operationDBTypes, false, operationColumnsWithDefault...)
	randomize.Struct(seed, &c, operationDBTypes, false, operationColumnsWithDefault...)

	b.TypeID = a.ID
	c.TypeID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	operation, err := a.TypeOperations(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range operation {
		if v.TypeID == b.TypeID {
			bFound = true
		}
		if v.TypeID == c.TypeID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := OperationTypeSlice{&a}
	if err = a.L.LoadTypeOperations(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TypeOperations); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.TypeOperations = nil
	if err = a.L.LoadTypeOperations(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TypeOperations); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", operation)
	}
}

func testOperationTypeToManyAddOpTypeOperations(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a OperationType
	var b, c, d, e Operation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationTypeDBTypes, false, strmangle.SetComplement(operationTypePrimaryKeyColumns, operationTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Operation{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, operationDBTypes, false, strmangle.SetComplement(operationPrimaryKeyColumns, operationColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Operation{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddTypeOperations(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.TypeID {
			t.Error("foreign key was wrong value", a.ID, first.TypeID)
		}
		if a.ID != second.TypeID {
			t.Error("foreign key was wrong value", a.ID, second.TypeID)
		}

		if first.R.Type != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Type != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.TypeOperations[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.TypeOperations[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.TypeOperations(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testOperationTypesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operationType := &OperationType{}
	if err = randomize.Struct(seed, operationType, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationType.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = operationType.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testOperationTypesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operationType := &OperationType{}
	if err = randomize.Struct(seed, operationType, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationType.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := OperationTypeSlice{operationType}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testOperationTypesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operationType := &OperationType{}
	if err = randomize.Struct(seed, operationType, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationType.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := OperationTypes(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	operationTypeDBTypes = map[string]string{`Description`: `character varying`, `ID`: `bigint`, `Name`: `character varying`}
	_                    = bytes.MinRead
)

func testOperationTypesUpdate(t *testing.T) {
	t.Parallel()

	if len(operationTypeColumns) == len(operationTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	operationType := &OperationType{}
	if err = randomize.Struct(seed, operationType, operationTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationType.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := OperationTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, operationType, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	if err = operationType.Update(tx); err != nil {
		t.Error(err)
	}
}

func testOperationTypesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(operationTypeColumns) == len(operationTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	operationType := &OperationType{}
	if err = randomize.Struct(seed, operationType, operationTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationType.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := OperationTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, operationType, operationTypeDBTypes, true, operationTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(operationTypeColumns, operationTypePrimaryKeyColumns) {
		fields = operationTypeColumns
	} else {
		fields = strmangle.SetComplement(
			operationTypeColumns,
			operationTypePrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(operationType))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := OperationTypeSlice{operationType}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testOperationTypesUpsert(t *testing.T) {
	t.Parallel()

	if len(operationTypeColumns) == len(operationTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	operationType := OperationType{}
	if err = randomize.Struct(seed, &operationType, operationTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationType.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert OperationType: %s", err)
	}

	count, err := OperationTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &operationType, operationTypeDBTypes, false, operationTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	if err = operationType.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert OperationType: %s", err)
	}

	count, err = OperationTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
