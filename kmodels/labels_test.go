package kmodels

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testLabels(t *testing.T) {
	t.Parallel()

	query := Labels(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testLabelsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	label := &Label{}
	if err = randomize.Struct(seed, label, labelDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = label.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = label.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Labels(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testLabelsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	label := &Label{}
	if err = randomize.Struct(seed, label, labelDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = label.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Labels(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Labels(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testLabelsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	label := &Label{}
	if err = randomize.Struct(seed, label, labelDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = label.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := LabelSlice{label}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Labels(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testLabelsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	label := &Label{}
	if err = randomize.Struct(seed, label, labelDBTypes, true, labelColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = label.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := LabelExists(tx, label.ID)
	if err != nil {
		t.Errorf("Unable to check if Label exists: %s", err)
	}
	if !e {
		t.Errorf("Expected LabelExistsG to return true, but got false.")
	}
}
func testLabelsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	label := &Label{}
	if err = randomize.Struct(seed, label, labelDBTypes, true, labelColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = label.Insert(tx); err != nil {
		t.Error(err)
	}

	labelFound, err := FindLabel(tx, label.ID)
	if err != nil {
		t.Error(err)
	}

	if labelFound == nil {
		t.Error("want a record, got nil")
	}
}
func testLabelsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	label := &Label{}
	if err = randomize.Struct(seed, label, labelDBTypes, true, labelColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = label.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Labels(tx).Bind(label); err != nil {
		t.Error(err)
	}
}

func testLabelsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	label := &Label{}
	if err = randomize.Struct(seed, label, labelDBTypes, true, labelColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = label.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Labels(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testLabelsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	labelOne := &Label{}
	labelTwo := &Label{}
	if err = randomize.Struct(seed, labelOne, labelDBTypes, false, labelColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}
	if err = randomize.Struct(seed, labelTwo, labelDBTypes, false, labelColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = labelOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = labelTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Labels(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testLabelsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	labelOne := &Label{}
	labelTwo := &Label{}
	if err = randomize.Struct(seed, labelOne, labelDBTypes, false, labelColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}
	if err = randomize.Struct(seed, labelTwo, labelDBTypes, false, labelColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = labelOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = labelTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Labels(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testLabelsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	label := &Label{}
	if err = randomize.Struct(seed, label, labelDBTypes, true, labelColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = label.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Labels(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testLabelsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	label := &Label{}
	if err = randomize.Struct(seed, label, labelDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = label.Insert(tx, labelColumns...); err != nil {
		t.Error(err)
	}

	count, err := Labels(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testLabelToManyContainers(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Label
	var b, c Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, labelDBTypes, true, labelColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, containerDBTypes, false, containerColumnsWithDefault...)
	randomize.Struct(seed, &c, containerDBTypes, false, containerColumnsWithDefault...)

	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec("insert into \"containers_labels\" (\"label_id\", \"container_id\") values ($1, $2)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into \"containers_labels\" (\"label_id\", \"container_id\") values ($1, $2)", a.ID, c.ID)
	if err != nil {
		t.Fatal(err)
	}

	container, err := a.Containers(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range container {
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

	slice := LabelSlice{&a}
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

func testLabelToManyAddOpContainers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Label
	var b, c, d, e Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, labelDBTypes, false, strmangle.SetComplement(labelPrimaryKeyColumns, labelColumnsWithoutDefault)...); err != nil {
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

		if first.R.Labels[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.Labels[0] != &a {
			t.Error("relationship was not added properly to the slice")
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

func testLabelToManySetOpContainers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Label
	var b, c, d, e Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, labelDBTypes, false, strmangle.SetComplement(labelPrimaryKeyColumns, labelColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.Labels) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Labels) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Labels[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.Labels[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.Containers[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Containers[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testLabelToManyRemoveOpContainers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Label
	var b, c, d, e Container

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, labelDBTypes, false, strmangle.SetComplement(labelPrimaryKeyColumns, labelColumnsWithoutDefault)...); err != nil {
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

	if len(b.R.Labels) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Labels) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Labels[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Labels[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
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

func testLabelsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	label := &Label{}
	if err = randomize.Struct(seed, label, labelDBTypes, true, labelColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = label.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = label.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testLabelsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	label := &Label{}
	if err = randomize.Struct(seed, label, labelDBTypes, true, labelColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = label.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := LabelSlice{label}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testLabelsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	label := &Label{}
	if err = randomize.Struct(seed, label, labelDBTypes, true, labelColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = label.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Labels(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	labelDBTypes = map[string]string{`CreatedAt`: `timestamp without time zone`, `DictionaryID`: `integer`, `ID`: `integer`, `Suid`: `character varying`, `UpdatedAt`: `timestamp without time zone`}
	_            = bytes.MinRead
)

func testLabelsUpdate(t *testing.T) {
	t.Parallel()

	if len(labelColumns) == len(labelPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	label := &Label{}
	if err = randomize.Struct(seed, label, labelDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = label.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Labels(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, label, labelDBTypes, true, labelColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	if err = label.Update(tx); err != nil {
		t.Error(err)
	}
}

func testLabelsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(labelColumns) == len(labelPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	label := &Label{}
	if err = randomize.Struct(seed, label, labelDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = label.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Labels(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, label, labelDBTypes, true, labelPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(labelColumns, labelPrimaryKeyColumns) {
		fields = labelColumns
	} else {
		fields = strmangle.SetComplement(
			labelColumns,
			labelPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(label))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := LabelSlice{label}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testLabelsUpsert(t *testing.T) {
	t.Parallel()

	if len(labelColumns) == len(labelPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	label := Label{}
	if err = randomize.Struct(seed, &label, labelDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = label.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Label: %s", err)
	}

	count, err := Labels(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &label, labelDBTypes, false, labelPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Label struct: %s", err)
	}

	if err = label.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Label: %s", err)
	}

	count, err = Labels(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
