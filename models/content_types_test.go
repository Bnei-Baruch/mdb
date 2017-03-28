package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testContentTypes(t *testing.T) {
	t.Parallel()

	query := ContentTypes(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testContentTypesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentType := &ContentType{}
	if err = randomize.Struct(seed, contentType, contentTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentType.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = contentType.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContentTypesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentType := &ContentType{}
	if err = randomize.Struct(seed, contentType, contentTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentType.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = ContentTypes(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := ContentTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContentTypesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentType := &ContentType{}
	if err = randomize.Struct(seed, contentType, contentTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentType.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ContentTypeSlice{contentType}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testContentTypesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentType := &ContentType{}
	if err = randomize.Struct(seed, contentType, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentType.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := ContentTypeExists(tx, contentType.ID)
	if err != nil {
		t.Errorf("Unable to check if ContentType exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ContentTypeExistsG to return true, but got false.")
	}
}
func testContentTypesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentType := &ContentType{}
	if err = randomize.Struct(seed, contentType, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentType.Insert(tx); err != nil {
		t.Error(err)
	}

	contentTypeFound, err := FindContentType(tx, contentType.ID)
	if err != nil {
		t.Error(err)
	}

	if contentTypeFound == nil {
		t.Error("want a record, got nil")
	}
}
func testContentTypesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentType := &ContentType{}
	if err = randomize.Struct(seed, contentType, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentType.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = ContentTypes(tx).Bind(contentType); err != nil {
		t.Error(err)
	}
}

func testContentTypesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentType := &ContentType{}
	if err = randomize.Struct(seed, contentType, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentType.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := ContentTypes(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testContentTypesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentTypeOne := &ContentType{}
	contentTypeTwo := &ContentType{}
	if err = randomize.Struct(seed, contentTypeOne, contentTypeDBTypes, false, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}
	if err = randomize.Struct(seed, contentTypeTwo, contentTypeDBTypes, false, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentTypeOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = contentTypeTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := ContentTypes(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testContentTypesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	contentTypeOne := &ContentType{}
	contentTypeTwo := &ContentType{}
	if err = randomize.Struct(seed, contentTypeOne, contentTypeDBTypes, false, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}
	if err = randomize.Struct(seed, contentTypeTwo, contentTypeDBTypes, false, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentTypeOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = contentTypeTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testContentTypesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentType := &ContentType{}
	if err = randomize.Struct(seed, contentType, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentType.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testContentTypesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentType := &ContentType{}
	if err = randomize.Struct(seed, contentType, contentTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentType.Insert(tx, contentTypeColumns...); err != nil {
		t.Error(err)
	}

	count, err := ContentTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testContentTypeToManyTypeContentUnits(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentType
	var b, c ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, contentUnitDBTypes, false, contentUnitColumnsWithDefault...)
	randomize.Struct(seed, &c, contentUnitDBTypes, false, contentUnitColumnsWithDefault...)

	b.TypeID = a.ID
	c.TypeID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	contentUnit, err := a.TypeContentUnits(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range contentUnit {
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

	slice := ContentTypeSlice{&a}
	if err = a.L.LoadTypeContentUnits(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TypeContentUnits); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.TypeContentUnits = nil
	if err = a.L.LoadTypeContentUnits(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TypeContentUnits); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", contentUnit)
	}
}

func testContentTypeToManyTypeCollections(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentType
	var b, c Collection

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, collectionDBTypes, false, collectionColumnsWithDefault...)
	randomize.Struct(seed, &c, collectionDBTypes, false, collectionColumnsWithDefault...)

	b.TypeID = a.ID
	c.TypeID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	collection, err := a.TypeCollections(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range collection {
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

	slice := ContentTypeSlice{&a}
	if err = a.L.LoadTypeCollections(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TypeCollections); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.TypeCollections = nil
	if err = a.L.LoadTypeCollections(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TypeCollections); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", collection)
	}
}

func testContentTypeToManyAddOpTypeContentUnits(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentType
	var b, c, d, e ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentTypeDBTypes, false, strmangle.SetComplement(contentTypePrimaryKeyColumns, contentTypeColumnsWithoutDefault)...); err != nil {
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
		err = a.AddTypeContentUnits(tx, i != 0, x...)
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

		if a.R.TypeContentUnits[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.TypeContentUnits[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.TypeContentUnits(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testContentTypeToManyAddOpTypeCollections(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentType
	var b, c, d, e Collection

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentTypeDBTypes, false, strmangle.SetComplement(contentTypePrimaryKeyColumns, contentTypeColumnsWithoutDefault)...); err != nil {
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
		err = a.AddTypeCollections(tx, i != 0, x...)
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

		if a.R.TypeCollections[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.TypeCollections[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.TypeCollections(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testContentTypesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentType := &ContentType{}
	if err = randomize.Struct(seed, contentType, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentType.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = contentType.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testContentTypesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentType := &ContentType{}
	if err = randomize.Struct(seed, contentType, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentType.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ContentTypeSlice{contentType}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testContentTypesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentType := &ContentType{}
	if err = randomize.Struct(seed, contentType, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentType.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := ContentTypes(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	contentTypeDBTypes = map[string]string{`Description`: `character varying`, `ID`: `bigint`, `Name`: `character varying`}
	_                  = bytes.MinRead
)

func testContentTypesUpdate(t *testing.T) {
	t.Parallel()

	if len(contentTypeColumns) == len(contentTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	contentType := &ContentType{}
	if err = randomize.Struct(seed, contentType, contentTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentType.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, contentType, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	if err = contentType.Update(tx); err != nil {
		t.Error(err)
	}
}

func testContentTypesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(contentTypeColumns) == len(contentTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	contentType := &ContentType{}
	if err = randomize.Struct(seed, contentType, contentTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentType.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, contentType, contentTypeDBTypes, true, contentTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(contentTypeColumns, contentTypePrimaryKeyColumns) {
		fields = contentTypeColumns
	} else {
		fields = strmangle.SetComplement(
			contentTypeColumns,
			contentTypePrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(contentType))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := ContentTypeSlice{contentType}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testContentTypesUpsert(t *testing.T) {
	t.Parallel()

	if len(contentTypeColumns) == len(contentTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	contentType := ContentType{}
	if err = randomize.Struct(seed, &contentType, contentTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentType.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert ContentType: %s", err)
	}

	count, err := ContentTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &contentType, contentTypeDBTypes, false, contentTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	if err = contentType.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert ContentType: %s", err)
	}

	count, err = ContentTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
