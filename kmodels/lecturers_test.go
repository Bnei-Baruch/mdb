package kmodels

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testLecturers(t *testing.T) {
	t.Parallel()

	query := Lecturers(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testLecturersDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturer := &Lecturer{}
	if err = randomize.Struct(seed, lecturer, lecturerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturer.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = lecturer.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Lecturers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testLecturersQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturer := &Lecturer{}
	if err = randomize.Struct(seed, lecturer, lecturerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturer.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Lecturers(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Lecturers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testLecturersSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturer := &Lecturer{}
	if err = randomize.Struct(seed, lecturer, lecturerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturer.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := LecturerSlice{lecturer}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Lecturers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testLecturersExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturer := &Lecturer{}
	if err = randomize.Struct(seed, lecturer, lecturerDBTypes, true, lecturerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturer.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := LecturerExists(tx, lecturer.ID)
	if err != nil {
		t.Errorf("Unable to check if Lecturer exists: %s", err)
	}
	if !e {
		t.Errorf("Expected LecturerExistsG to return true, but got false.")
	}
}
func testLecturersFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturer := &Lecturer{}
	if err = randomize.Struct(seed, lecturer, lecturerDBTypes, true, lecturerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturer.Insert(tx); err != nil {
		t.Error(err)
	}

	lecturerFound, err := FindLecturer(tx, lecturer.ID)
	if err != nil {
		t.Error(err)
	}

	if lecturerFound == nil {
		t.Error("want a record, got nil")
	}
}
func testLecturersBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturer := &Lecturer{}
	if err = randomize.Struct(seed, lecturer, lecturerDBTypes, true, lecturerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturer.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Lecturers(tx).Bind(lecturer); err != nil {
		t.Error(err)
	}
}

func testLecturersOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturer := &Lecturer{}
	if err = randomize.Struct(seed, lecturer, lecturerDBTypes, true, lecturerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturer.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Lecturers(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testLecturersAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturerOne := &Lecturer{}
	lecturerTwo := &Lecturer{}
	if err = randomize.Struct(seed, lecturerOne, lecturerDBTypes, false, lecturerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}
	if err = randomize.Struct(seed, lecturerTwo, lecturerDBTypes, false, lecturerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = lecturerTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Lecturers(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testLecturersCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	lecturerOne := &Lecturer{}
	lecturerTwo := &Lecturer{}
	if err = randomize.Struct(seed, lecturerOne, lecturerDBTypes, false, lecturerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}
	if err = randomize.Struct(seed, lecturerTwo, lecturerDBTypes, false, lecturerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturerOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = lecturerTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Lecturers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testLecturersInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturer := &Lecturer{}
	if err = randomize.Struct(seed, lecturer, lecturerDBTypes, true, lecturerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturer.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Lecturers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testLecturersInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturer := &Lecturer{}
	if err = randomize.Struct(seed, lecturer, lecturerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturer.Insert(tx, lecturerColumns...); err != nil {
		t.Error(err)
	}

	count, err := Lecturers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testLecturersReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturer := &Lecturer{}
	if err = randomize.Struct(seed, lecturer, lecturerDBTypes, true, lecturerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturer.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = lecturer.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testLecturersReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturer := &Lecturer{}
	if err = randomize.Struct(seed, lecturer, lecturerDBTypes, true, lecturerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturer.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := LecturerSlice{lecturer}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testLecturersSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	lecturer := &Lecturer{}
	if err = randomize.Struct(seed, lecturer, lecturerDBTypes, true, lecturerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturer.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Lecturers(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	lecturerDBTypes = map[string]string{`CreatedAt`: `timestamp without time zone`, `ID`: `integer`, `Name`: `character`, `Ordnum`: `integer`, `UpdatedAt`: `timestamp without time zone`}
	_               = bytes.MinRead
)

func testLecturersUpdate(t *testing.T) {
	t.Parallel()

	if len(lecturerColumns) == len(lecturerPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	lecturer := &Lecturer{}
	if err = randomize.Struct(seed, lecturer, lecturerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturer.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Lecturers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, lecturer, lecturerDBTypes, true, lecturerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	if err = lecturer.Update(tx); err != nil {
		t.Error(err)
	}
}

func testLecturersSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(lecturerColumns) == len(lecturerPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	lecturer := &Lecturer{}
	if err = randomize.Struct(seed, lecturer, lecturerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturer.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Lecturers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, lecturer, lecturerDBTypes, true, lecturerPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(lecturerColumns, lecturerPrimaryKeyColumns) {
		fields = lecturerColumns
	} else {
		fields = strmangle.SetComplement(
			lecturerColumns,
			lecturerPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(lecturer))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := LecturerSlice{lecturer}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testLecturersUpsert(t *testing.T) {
	t.Parallel()

	if len(lecturerColumns) == len(lecturerPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	lecturer := Lecturer{}
	if err = randomize.Struct(seed, &lecturer, lecturerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = lecturer.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Lecturer: %s", err)
	}

	count, err := Lecturers(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &lecturer, lecturerDBTypes, false, lecturerPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Lecturer struct: %s", err)
	}

	if err = lecturer.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Lecturer: %s", err)
	}

	count, err = Lecturers(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
