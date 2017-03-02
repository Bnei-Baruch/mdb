package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testFiles(t *testing.T) {
	t.Parallel()

	query := Files(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testFilesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = file.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testFilesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Files(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testFilesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := FileSlice{file}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testFilesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := FileExists(tx, file.ID)
	if err != nil {
		t.Errorf("Unable to check if File exists: %s", err)
	}
	if !e {
		t.Errorf("Expected FileExistsG to return true, but got false.")
	}
}
func testFilesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	fileFound, err := FindFile(tx, file.ID)
	if err != nil {
		t.Error(err)
	}

	if fileFound == nil {
		t.Error("want a record, got nil")
	}
}
func testFilesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Files(tx).Bind(file); err != nil {
		t.Error(err)
	}
}

func testFilesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Files(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testFilesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	fileOne := &File{}
	fileTwo := &File{}
	if err = randomize.Struct(seed, fileOne, fileDBTypes, false, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}
	if err = randomize.Struct(seed, fileTwo, fileDBTypes, false, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = fileTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Files(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testFilesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	fileOne := &File{}
	fileTwo := &File{}
	if err = randomize.Struct(seed, fileOne, fileDBTypes, false, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}
	if err = randomize.Struct(seed, fileTwo, fileDBTypes, false, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = fileOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = fileTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testFilesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testFilesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx, fileColumns...); err != nil {
		t.Error(err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testFileToManyOperations(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c Operation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, operationDBTypes, false, operationColumnsWithDefault...)
	randomize.Struct(seed, &c, operationDBTypes, false, operationColumnsWithDefault...)

	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec("insert into \"files_operations\" (\"file_id\", \"operation_id\") values ($1, $2)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into \"files_operations\" (\"file_id\", \"operation_id\") values ($1, $2)", a.ID, c.ID)
	if err != nil {
		t.Fatal(err)
	}

	operation, err := a.Operations(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range operation {
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

	slice := FileSlice{&a}
	if err = a.L.LoadOperations(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Operations); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Operations = nil
	if err = a.L.LoadOperations(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Operations); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", operation)
	}
}

func testFileToManyParentFiles(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, fileDBTypes, false, fileColumnsWithDefault...)
	randomize.Struct(seed, &c, fileDBTypes, false, fileColumnsWithDefault...)

	b.ParentID.Valid = true
	c.ParentID.Valid = true
	b.ParentID.Int64 = a.ID
	c.ParentID.Int64 = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	file, err := a.ParentFiles(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range file {
		if v.ParentID.Int64 == b.ParentID.Int64 {
			bFound = true
		}
		if v.ParentID.Int64 == c.ParentID.Int64 {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := FileSlice{&a}
	if err = a.L.LoadParentFiles(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ParentFiles); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.ParentFiles = nil
	if err = a.L.LoadParentFiles(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ParentFiles); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", file)
	}
}

func testFileToManyAddOpOperations(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c, d, e Operation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
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
		err = a.AddOperations(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if first.R.Files[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.Files[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}

		if a.R.Operations[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Operations[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Operations(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testFileToManySetOpOperations(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c, d, e Operation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Operation{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, operationDBTypes, false, strmangle.SetComplement(operationPrimaryKeyColumns, operationColumnsWithoutDefault)...); err != nil {
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

	err = a.SetOperations(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Operations(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetOperations(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Operations(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Files) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Files) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Files[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.Files[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.Operations[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Operations[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testFileToManyRemoveOpOperations(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c, d, e Operation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
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

	err = a.AddOperations(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Operations(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveOperations(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Operations(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Files) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Files) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Files[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Files[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if len(a.R.Operations) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.Operations[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.Operations[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testFileToManyAddOpParentFiles(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c, d, e File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
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
		err = a.AddParentFiles(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.ParentID.Int64 {
			t.Error("foreign key was wrong value", a.ID, first.ParentID.Int64)
		}
		if a.ID != second.ParentID.Int64 {
			t.Error("foreign key was wrong value", a.ID, second.ParentID.Int64)
		}

		if first.R.Parent != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Parent != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.ParentFiles[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.ParentFiles[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.ParentFiles(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testFileToManySetOpParentFiles(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c, d, e File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
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

	err = a.SetParentFiles(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ParentFiles(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetParentFiles(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ParentFiles(tx).Count()
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
	if a.ID != d.ParentID.Int64 {
		t.Error("foreign key was wrong value", a.ID, d.ParentID.Int64)
	}
	if a.ID != e.ParentID.Int64 {
		t.Error("foreign key was wrong value", a.ID, e.ParentID.Int64)
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

	if a.R.ParentFiles[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.ParentFiles[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testFileToManyRemoveOpParentFiles(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c, d, e File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
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

	err = a.AddParentFiles(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ParentFiles(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveParentFiles(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ParentFiles(tx).Count()
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

	if len(a.R.ParentFiles) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.ParentFiles[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.ParentFiles[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testFileToOneContentUnitUsingContentUnit(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local File
	var foreign ContentUnit

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	local.ContentUnitID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.ContentUnitID.Int64 = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.ContentUnit(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := FileSlice{&local}
	if err = local.L.LoadContentUnit(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.ContentUnit == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.ContentUnit = nil
	if err = local.L.LoadContentUnit(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.ContentUnit == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testFileToOneFileUsingParent(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local File
	var foreign File

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	local.ParentID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.ParentID.Int64 = foreign.ID
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

	slice := FileSlice{&local}
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

func testFileToOneSetOpContentUnitUsingContentUnit(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*ContentUnit{&b, &c} {
		err = a.SetContentUnit(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.ContentUnit != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.Files[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.ContentUnitID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.ContentUnitID.Int64)
		}

		zero := reflect.Zero(reflect.TypeOf(a.ContentUnitID.Int64))
		reflect.Indirect(reflect.ValueOf(&a.ContentUnitID.Int64)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.ContentUnitID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.ContentUnitID.Int64, x.ID)
		}
	}
}

func testFileToOneRemoveOpContentUnitUsingContentUnit(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.SetContentUnit(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveContentUnit(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.ContentUnit(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.ContentUnit != nil {
		t.Error("R struct entry should be nil")
	}

	if a.ContentUnitID.Valid {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.Files) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testFileToOneSetOpFileUsingParent(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b, c File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*File{&b, &c} {
		err = a.SetParent(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Parent != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.ParentFiles[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.ParentID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.ParentID.Int64)
		}

		zero := reflect.Zero(reflect.TypeOf(a.ParentID.Int64))
		reflect.Indirect(reflect.ValueOf(&a.ParentID.Int64)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.ParentID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.ParentID.Int64, x.ID)
		}
	}
}

func testFileToOneRemoveOpFileUsingParent(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a File
	var b File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.ParentFiles) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testFilesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = file.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testFilesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := FileSlice{file}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testFilesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Files(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	fileDBTypes = map[string]string{`BackupCount`: `smallint`, `ContentUnitID`: `bigint`, `CreatedAt`: `timestamp with time zone`, `FileCreatedAt`: `timestamp with time zone`, `FirstBackupTime`: `timestamp with time zone`, `ID`: `bigint`, `Language`: `character`, `MimeType`: `character varying`, `Name`: `character varying`, `ParentID`: `bigint`, `Properties`: `jsonb`, `Sha1`: `bytea`, `Size`: `bigint`, `SubType`: `character varying`, `Type`: `character varying`, `UID`: `character`}
	_           = bytes.MinRead
)

func testFilesUpdate(t *testing.T) {
	t.Parallel()

	if len(fileColumns) == len(filePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, file, fileDBTypes, true, fileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	if err = file.Update(tx); err != nil {
		t.Error(err)
	}
}

func testFilesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(fileColumns) == len(filePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	file := &File{}
	if err = randomize.Struct(seed, file, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, file, fileDBTypes, true, filePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(fileColumns, filePrimaryKeyColumns) {
		fields = fileColumns
	} else {
		fields = strmangle.SetComplement(
			fileColumns,
			filePrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(file))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := FileSlice{file}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testFilesUpsert(t *testing.T) {
	t.Parallel()

	if len(fileColumns) == len(filePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	file := File{}
	if err = randomize.Struct(seed, &file, fileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = file.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert File: %s", err)
	}

	count, err := Files(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &file, fileDBTypes, false, filePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize File struct: %s", err)
	}

	if err = file.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert File: %s", err)
	}

	count, err = Files(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}