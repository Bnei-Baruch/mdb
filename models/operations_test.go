package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testOperations(t *testing.T) {
	t.Parallel()

	query := Operations(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testOperationsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operation := &Operation{}
	if err = randomize.Struct(seed, operation, operationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operation.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = operation.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Operations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testOperationsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operation := &Operation{}
	if err = randomize.Struct(seed, operation, operationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operation.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Operations(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Operations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testOperationsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operation := &Operation{}
	if err = randomize.Struct(seed, operation, operationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operation.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := OperationSlice{operation}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Operations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testOperationsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operation := &Operation{}
	if err = randomize.Struct(seed, operation, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operation.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := OperationExists(tx, operation.ID)
	if err != nil {
		t.Errorf("Unable to check if Operation exists: %s", err)
	}
	if !e {
		t.Errorf("Expected OperationExistsG to return true, but got false.")
	}
}
func testOperationsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operation := &Operation{}
	if err = randomize.Struct(seed, operation, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operation.Insert(tx); err != nil {
		t.Error(err)
	}

	operationFound, err := FindOperation(tx, operation.ID)
	if err != nil {
		t.Error(err)
	}

	if operationFound == nil {
		t.Error("want a record, got nil")
	}
}
func testOperationsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operation := &Operation{}
	if err = randomize.Struct(seed, operation, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operation.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Operations(tx).Bind(operation); err != nil {
		t.Error(err)
	}
}

func testOperationsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operation := &Operation{}
	if err = randomize.Struct(seed, operation, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operation.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Operations(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testOperationsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operationOne := &Operation{}
	operationTwo := &Operation{}
	if err = randomize.Struct(seed, operationOne, operationDBTypes, false, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}
	if err = randomize.Struct(seed, operationTwo, operationDBTypes, false, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = operationTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Operations(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testOperationsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	operationOne := &Operation{}
	operationTwo := &Operation{}
	if err = randomize.Struct(seed, operationOne, operationDBTypes, false, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}
	if err = randomize.Struct(seed, operationTwo, operationDBTypes, false, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operationOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = operationTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Operations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testOperationsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operation := &Operation{}
	if err = randomize.Struct(seed, operation, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operation.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Operations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testOperationsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operation := &Operation{}
	if err = randomize.Struct(seed, operation, operationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operation.Insert(tx, operationColumns...); err != nil {
		t.Error(err)
	}

	count, err := Operations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testOperationToManyFiles(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Operation
	var b, c File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, fileDBTypes, false, fileColumnsWithDefault...)
	randomize.Struct(seed, &c, fileDBTypes, false, fileColumnsWithDefault...)

	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec("insert into \"files_operations\" (\"operation_id\", \"file_id\") values ($1, $2)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into \"files_operations\" (\"operation_id\", \"file_id\") values ($1, $2)", a.ID, c.ID)
	if err != nil {
		t.Fatal(err)
	}

	file, err := a.Files(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range file {
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

	slice := OperationSlice{&a}
	if err = a.L.LoadFiles(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Files); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Files = nil
	if err = a.L.LoadFiles(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Files); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", file)
	}
}

func testOperationToManyAddOpFiles(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Operation
	var b, c, d, e File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationDBTypes, false, strmangle.SetComplement(operationPrimaryKeyColumns, operationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*File{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*File{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddFiles(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if first.R.Operations[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.Operations[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}

		if a.R.Files[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Files[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Files(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testOperationToManySetOpFiles(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Operation
	var b, c, d, e File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationDBTypes, false, strmangle.SetComplement(operationPrimaryKeyColumns, operationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*File{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
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

	err = a.SetFiles(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Files(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetFiles(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Files(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Operations) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Operations) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Operations[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.Operations[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.Files[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Files[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testOperationToManyRemoveOpFiles(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Operation
	var b, c, d, e File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationDBTypes, false, strmangle.SetComplement(operationPrimaryKeyColumns, operationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*File{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddFiles(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Files(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveFiles(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Files(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Operations) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Operations) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Operations[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Operations[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if len(a.R.Files) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.Files[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.Files[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testOperationToOneOperationTypeUsingType(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Operation
	var foreign OperationType

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, operationTypeDBTypes, true, operationTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize OperationType struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.TypeID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Type(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := OperationSlice{&local}
	if err = local.L.LoadType(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Type == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Type = nil
	if err = local.L.LoadType(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Type == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testOperationToOneUserUsingUser(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Operation
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	local.UserID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.UserID.Int64 = foreign.ID
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

	slice := OperationSlice{&local}
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

func testOperationToOneSetOpOperationTypeUsingType(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Operation
	var b, c OperationType

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationDBTypes, false, strmangle.SetComplement(operationPrimaryKeyColumns, operationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, operationTypeDBTypes, false, strmangle.SetComplement(operationTypePrimaryKeyColumns, operationTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, operationTypeDBTypes, false, strmangle.SetComplement(operationTypePrimaryKeyColumns, operationTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*OperationType{&b, &c} {
		err = a.SetType(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Type != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.TypeOperations[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.TypeID != x.ID {
			t.Error("foreign key was wrong value", a.TypeID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.TypeID))
		reflect.Indirect(reflect.ValueOf(&a.TypeID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.TypeID != x.ID {
			t.Error("foreign key was wrong value", a.TypeID, x.ID)
		}
	}
}
func testOperationToOneSetOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Operation
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationDBTypes, false, strmangle.SetComplement(operationPrimaryKeyColumns, operationColumnsWithoutDefault)...); err != nil {
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

		if x.R.Operations[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.UserID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.UserID.Int64)
		}

		zero := reflect.Zero(reflect.TypeOf(a.UserID.Int64))
		reflect.Indirect(reflect.ValueOf(&a.UserID.Int64)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.UserID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.UserID.Int64, x.ID)
		}
	}
}

func testOperationToOneRemoveOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Operation
	var b User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, operationDBTypes, false, strmangle.SetComplement(operationPrimaryKeyColumns, operationColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.Operations) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testOperationsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operation := &Operation{}
	if err = randomize.Struct(seed, operation, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operation.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = operation.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testOperationsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operation := &Operation{}
	if err = randomize.Struct(seed, operation, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operation.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := OperationSlice{operation}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testOperationsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	operation := &Operation{}
	if err = randomize.Struct(seed, operation, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operation.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Operations(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	operationDBTypes = map[string]string{`CreatedAt`: `timestamp with time zone`, `Details`: `character varying`, `ID`: `bigint`, `Properties`: `jsonb`, `Station`: `character varying`, `TypeID`: `bigint`, `UID`: `character`, `UserID`: `bigint`}
	_                = bytes.MinRead
)

func testOperationsUpdate(t *testing.T) {
	t.Parallel()

	if len(operationColumns) == len(operationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	operation := &Operation{}
	if err = randomize.Struct(seed, operation, operationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operation.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Operations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, operation, operationDBTypes, true, operationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	if err = operation.Update(tx); err != nil {
		t.Error(err)
	}
}

func testOperationsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(operationColumns) == len(operationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	operation := &Operation{}
	if err = randomize.Struct(seed, operation, operationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operation.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Operations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, operation, operationDBTypes, true, operationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(operationColumns, operationPrimaryKeyColumns) {
		fields = operationColumns
	} else {
		fields = strmangle.SetComplement(
			operationColumns,
			operationPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(operation))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := OperationSlice{operation}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testOperationsUpsert(t *testing.T) {
	t.Parallel()

	if len(operationColumns) == len(operationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	operation := Operation{}
	if err = randomize.Struct(seed, &operation, operationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = operation.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Operation: %s", err)
	}

	count, err := Operations(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &operation, operationDBTypes, false, operationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Operation struct: %s", err)
	}

	if err = operation.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Operation: %s", err)
	}

	count, err = Operations(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
