package kmodels

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testVirtualLessons(t *testing.T) {
	t.Parallel()

	query := VirtualLessons(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testVirtualLessonsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	virtualLesson := &VirtualLesson{}
	if err = randomize.Struct(seed, virtualLesson, virtualLessonDBTypes, true); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLesson.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = virtualLesson.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := VirtualLessons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testVirtualLessonsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	virtualLesson := &VirtualLesson{}
	if err = randomize.Struct(seed, virtualLesson, virtualLessonDBTypes, true); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLesson.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = VirtualLessons(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := VirtualLessons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testVirtualLessonsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	virtualLesson := &VirtualLesson{}
	if err = randomize.Struct(seed, virtualLesson, virtualLessonDBTypes, true); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLesson.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := VirtualLessonSlice{virtualLesson}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := VirtualLessons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testVirtualLessonsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	virtualLesson := &VirtualLesson{}
	if err = randomize.Struct(seed, virtualLesson, virtualLessonDBTypes, true, virtualLessonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLesson.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := VirtualLessonExists(tx, virtualLesson.ID)
	if err != nil {
		t.Errorf("Unable to check if VirtualLesson exists: %s", err)
	}
	if !e {
		t.Errorf("Expected VirtualLessonExistsG to return true, but got false.")
	}
}
func testVirtualLessonsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	virtualLesson := &VirtualLesson{}
	if err = randomize.Struct(seed, virtualLesson, virtualLessonDBTypes, true, virtualLessonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLesson.Insert(tx); err != nil {
		t.Error(err)
	}

	virtualLessonFound, err := FindVirtualLesson(tx, virtualLesson.ID)
	if err != nil {
		t.Error(err)
	}

	if virtualLessonFound == nil {
		t.Error("want a record, got nil")
	}
}
func testVirtualLessonsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	virtualLesson := &VirtualLesson{}
	if err = randomize.Struct(seed, virtualLesson, virtualLessonDBTypes, true, virtualLessonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLesson.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = VirtualLessons(tx).Bind(virtualLesson); err != nil {
		t.Error(err)
	}
}

func testVirtualLessonsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	virtualLesson := &VirtualLesson{}
	if err = randomize.Struct(seed, virtualLesson, virtualLessonDBTypes, true, virtualLessonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLesson.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := VirtualLessons(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testVirtualLessonsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	virtualLessonOne := &VirtualLesson{}
	virtualLessonTwo := &VirtualLesson{}
	if err = randomize.Struct(seed, virtualLessonOne, virtualLessonDBTypes, false, virtualLessonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}
	if err = randomize.Struct(seed, virtualLessonTwo, virtualLessonDBTypes, false, virtualLessonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLessonOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = virtualLessonTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := VirtualLessons(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testVirtualLessonsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	virtualLessonOne := &VirtualLesson{}
	virtualLessonTwo := &VirtualLesson{}
	if err = randomize.Struct(seed, virtualLessonOne, virtualLessonDBTypes, false, virtualLessonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}
	if err = randomize.Struct(seed, virtualLessonTwo, virtualLessonDBTypes, false, virtualLessonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLessonOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = virtualLessonTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := VirtualLessons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testVirtualLessonsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	virtualLesson := &VirtualLesson{}
	if err = randomize.Struct(seed, virtualLesson, virtualLessonDBTypes, true, virtualLessonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLesson.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := VirtualLessons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testVirtualLessonsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	virtualLesson := &VirtualLesson{}
	if err = randomize.Struct(seed, virtualLesson, virtualLessonDBTypes, true); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLesson.Insert(tx, virtualLessonColumns...); err != nil {
		t.Error(err)
	}

	count, err := VirtualLessons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testVirtualLessonToManyContainers(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a VirtualLesson
	var b, c Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, virtualLessonDBTypes, true, virtualLessonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, containerDBTypes, false, containerColumnsWithDefault...)
	randomize.Struct(seed, &c, containerDBTypes, false, containerColumnsWithDefault...)

	b.VirtualLessonID.Valid = true
	c.VirtualLessonID.Valid = true
	b.VirtualLessonID.Int = a.ID
	c.VirtualLessonID.Int = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	container, err := a.Containers(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range container {
		if v.VirtualLessonID.Int == b.VirtualLessonID.Int {
			bFound = true
		}
		if v.VirtualLessonID.Int == c.VirtualLessonID.Int {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := VirtualLessonSlice{&a}
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

func testVirtualLessonToManyAddOpContainers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a VirtualLesson
	var b, c, d, e Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, virtualLessonDBTypes, false, strmangle.SetComplement(virtualLessonPrimaryKeyColumns, virtualLessonColumnsWithoutDefault)...); err != nil {
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

		if a.ID != first.VirtualLessonID.Int {
			t.Error("foreign key was wrong value", a.ID, first.VirtualLessonID.Int)
		}
		if a.ID != second.VirtualLessonID.Int {
			t.Error("foreign key was wrong value", a.ID, second.VirtualLessonID.Int)
		}

		if first.R.VirtualLesson != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.VirtualLesson != &a {
			t.Error("relationship was not added properly to the foreign slice")
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

func testVirtualLessonToManySetOpContainers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a VirtualLesson
	var b, c, d, e Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, virtualLessonDBTypes, false, strmangle.SetComplement(virtualLessonPrimaryKeyColumns, virtualLessonColumnsWithoutDefault)...); err != nil {
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

	if b.VirtualLessonID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.VirtualLessonID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.ID != d.VirtualLessonID.Int {
		t.Error("foreign key was wrong value", a.ID, d.VirtualLessonID.Int)
	}
	if a.ID != e.VirtualLessonID.Int {
		t.Error("foreign key was wrong value", a.ID, e.VirtualLessonID.Int)
	}

	if b.R.VirtualLesson != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.VirtualLesson != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.VirtualLesson != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.VirtualLesson != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.Containers[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Containers[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testVirtualLessonToManyRemoveOpContainers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a VirtualLesson
	var b, c, d, e Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, virtualLessonDBTypes, false, strmangle.SetComplement(virtualLessonPrimaryKeyColumns, virtualLessonColumnsWithoutDefault)...); err != nil {
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

	if b.VirtualLessonID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.VirtualLessonID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.VirtualLesson != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.VirtualLesson != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.VirtualLesson != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.VirtualLesson != &a {
		t.Error("relationship to a should have been preserved")
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

func testVirtualLessonsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	virtualLesson := &VirtualLesson{}
	if err = randomize.Struct(seed, virtualLesson, virtualLessonDBTypes, true, virtualLessonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLesson.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = virtualLesson.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testVirtualLessonsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	virtualLesson := &VirtualLesson{}
	if err = randomize.Struct(seed, virtualLesson, virtualLessonDBTypes, true, virtualLessonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLesson.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := VirtualLessonSlice{virtualLesson}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testVirtualLessonsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	virtualLesson := &VirtualLesson{}
	if err = randomize.Struct(seed, virtualLesson, virtualLessonDBTypes, true, virtualLessonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLesson.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := VirtualLessons(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	virtualLessonDBTypes = map[string]string{`CreatedAt`: `timestamp without time zone`, `FilmDate`: `date`, `ID`: `integer`, `UpdatedAt`: `timestamp without time zone`, `UserID`: `integer`}
	_                    = bytes.MinRead
)

func testVirtualLessonsUpdate(t *testing.T) {
	t.Parallel()

	if len(virtualLessonColumns) == len(virtualLessonPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	virtualLesson := &VirtualLesson{}
	if err = randomize.Struct(seed, virtualLesson, virtualLessonDBTypes, true); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLesson.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := VirtualLessons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, virtualLesson, virtualLessonDBTypes, true, virtualLessonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	if err = virtualLesson.Update(tx); err != nil {
		t.Error(err)
	}
}

func testVirtualLessonsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(virtualLessonColumns) == len(virtualLessonPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	virtualLesson := &VirtualLesson{}
	if err = randomize.Struct(seed, virtualLesson, virtualLessonDBTypes, true); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLesson.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := VirtualLessons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, virtualLesson, virtualLessonDBTypes, true, virtualLessonPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(virtualLessonColumns, virtualLessonPrimaryKeyColumns) {
		fields = virtualLessonColumns
	} else {
		fields = strmangle.SetComplement(
			virtualLessonColumns,
			virtualLessonPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(virtualLesson))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := VirtualLessonSlice{virtualLesson}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testVirtualLessonsUpsert(t *testing.T) {
	t.Parallel()

	if len(virtualLessonColumns) == len(virtualLessonPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	virtualLesson := VirtualLesson{}
	if err = randomize.Struct(seed, &virtualLesson, virtualLessonDBTypes, true); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = virtualLesson.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert VirtualLesson: %s", err)
	}

	count, err := VirtualLessons(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &virtualLesson, virtualLessonDBTypes, false, virtualLessonPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize VirtualLesson struct: %s", err)
	}

	if err = virtualLesson.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert VirtualLesson: %s", err)
	}

	count, err = VirtualLessons(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
