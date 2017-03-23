package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testAuthors(t *testing.T) {
	t.Parallel()

	query := Authors(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testAuthorsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = author.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuthorsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Authors(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuthorsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := AuthorSlice{author}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testAuthorsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := AuthorExists(tx, author.ID)
	if err != nil {
		t.Errorf("Unable to check if Author exists: %s", err)
	}
	if !e {
		t.Errorf("Expected AuthorExistsG to return true, but got false.")
	}
}
func testAuthorsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	authorFound, err := FindAuthor(tx, author.ID)
	if err != nil {
		t.Error(err)
	}

	if authorFound == nil {
		t.Error("want a record, got nil")
	}
}
func testAuthorsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Authors(tx).Bind(author); err != nil {
		t.Error(err)
	}
}

func testAuthorsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Authors(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testAuthorsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorOne := &Author{}
	authorTwo := &Author{}
	if err = randomize.Struct(seed, authorOne, authorDBTypes, false, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}
	if err = randomize.Struct(seed, authorTwo, authorDBTypes, false, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = authorTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Authors(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testAuthorsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	authorOne := &Author{}
	authorTwo := &Author{}
	if err = randomize.Struct(seed, authorOne, authorDBTypes, false, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}
	if err = randomize.Struct(seed, authorTwo, authorDBTypes, false, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = authorTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testAuthorsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAuthorsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx, authorColumns...); err != nil {
		t.Error(err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAuthorToManyAuthorI18ns(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Author
	var b, c AuthorI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, authorI18nDBTypes, false, authorI18nColumnsWithDefault...)
	randomize.Struct(seed, &c, authorI18nDBTypes, false, authorI18nColumnsWithDefault...)

	b.AuthorID = a.ID
	c.AuthorID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	authorI18n, err := a.AuthorI18ns(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range authorI18n {
		if v.AuthorID == b.AuthorID {
			bFound = true
		}
		if v.AuthorID == c.AuthorID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := AuthorSlice{&a}
	if err = a.L.LoadAuthorI18ns(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.AuthorI18ns); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.AuthorI18ns = nil
	if err = a.L.LoadAuthorI18ns(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.AuthorI18ns); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", authorI18n)
	}
}

func testAuthorToManySources(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Author
	var b, c Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, sourceDBTypes, false, sourceColumnsWithDefault...)
	randomize.Struct(seed, &c, sourceDBTypes, false, sourceColumnsWithDefault...)

	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	_, err = tx.Exec("insert into \"authors_sources\" (\"author_id\", \"source_id\") values ($1, $2)", a.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into \"authors_sources\" (\"author_id\", \"source_id\") values ($1, $2)", a.ID, c.ID)
	if err != nil {
		t.Fatal(err)
	}

	source, err := a.Sources(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range source {
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

	slice := AuthorSlice{&a}
	if err = a.L.LoadSources(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Sources); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Sources = nil
	if err = a.L.LoadSources(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Sources); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", source)
	}
}

func testAuthorToManyAddOpAuthorI18ns(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Author
	var b, c, d, e AuthorI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorDBTypes, false, strmangle.SetComplement(authorPrimaryKeyColumns, authorColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*AuthorI18n{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, authorI18nDBTypes, false, strmangle.SetComplement(authorI18nPrimaryKeyColumns, authorI18nColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*AuthorI18n{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddAuthorI18ns(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.AuthorID {
			t.Error("foreign key was wrong value", a.ID, first.AuthorID)
		}
		if a.ID != second.AuthorID {
			t.Error("foreign key was wrong value", a.ID, second.AuthorID)
		}

		if first.R.Author != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Author != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.AuthorI18ns[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.AuthorI18ns[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.AuthorI18ns(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testAuthorToManyAddOpSources(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Author
	var b, c, d, e Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorDBTypes, false, strmangle.SetComplement(authorPrimaryKeyColumns, authorColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Source{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Source{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddSources(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if first.R.Authors[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.Authors[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}

		if a.R.Sources[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Sources[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Sources(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testAuthorToManySetOpSources(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Author
	var b, c, d, e Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorDBTypes, false, strmangle.SetComplement(authorPrimaryKeyColumns, authorColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Source{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
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

	err = a.SetSources(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Sources(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetSources(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Sources(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Authors) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Authors) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Authors[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.Authors[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}

	if a.R.Sources[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Sources[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testAuthorToManyRemoveOpSources(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Author
	var b, c, d, e Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, authorDBTypes, false, strmangle.SetComplement(authorPrimaryKeyColumns, authorColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Source{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddSources(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Sources(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveSources(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Sources(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if len(b.R.Authors) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.Authors) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.Authors[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Authors[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if len(a.R.Sources) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.Sources[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.Sources[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testAuthorsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = author.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testAuthorsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := AuthorSlice{author}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testAuthorsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Authors(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	authorDBTypes = map[string]string{`Code`: `character`, `CreatedAt`: `timestamp with time zone`, `FullName`: `character varying`, `ID`: `bigint`, `Name`: `character varying`}
	_             = bytes.MinRead
)

func testAuthorsUpdate(t *testing.T) {
	t.Parallel()

	if len(authorColumns) == len(authorPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	if err = author.Update(tx); err != nil {
		t.Error(err)
	}
}

func testAuthorsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(authorColumns) == len(authorPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, author, authorDBTypes, true, authorPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(authorColumns, authorPrimaryKeyColumns) {
		fields = authorColumns
	} else {
		fields = strmangle.SetComplement(
			authorColumns,
			authorPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(author))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := AuthorSlice{author}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testAuthorsUpsert(t *testing.T) {
	t.Parallel()

	if len(authorColumns) == len(authorPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	author := Author{}
	if err = randomize.Struct(seed, &author, authorDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Author: %s", err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &author, authorDBTypes, false, authorPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	if err = author.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Author: %s", err)
	}

	count, err = Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
