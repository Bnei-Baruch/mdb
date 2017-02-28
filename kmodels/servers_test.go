package kmodels

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testServers(t *testing.T) {
	t.Parallel()

	query := Servers(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testServersDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	server := &Server{}
	if err = randomize.Struct(seed, server, serverDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = server.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = server.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Servers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testServersQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	server := &Server{}
	if err = randomize.Struct(seed, server, serverDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = server.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Servers(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Servers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testServersSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	server := &Server{}
	if err = randomize.Struct(seed, server, serverDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = server.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ServerSlice{server}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Servers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testServersExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	server := &Server{}
	if err = randomize.Struct(seed, server, serverDBTypes, true, serverColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = server.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := ServerExists(tx, server.ID)
	if err != nil {
		t.Errorf("Unable to check if Server exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ServerExistsG to return true, but got false.")
	}
}
func testServersFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	server := &Server{}
	if err = randomize.Struct(seed, server, serverDBTypes, true, serverColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = server.Insert(tx); err != nil {
		t.Error(err)
	}

	serverFound, err := FindServer(tx, server.ID)
	if err != nil {
		t.Error(err)
	}

	if serverFound == nil {
		t.Error("want a record, got nil")
	}
}
func testServersBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	server := &Server{}
	if err = randomize.Struct(seed, server, serverDBTypes, true, serverColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = server.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Servers(tx).Bind(server); err != nil {
		t.Error(err)
	}
}

func testServersOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	server := &Server{}
	if err = randomize.Struct(seed, server, serverDBTypes, true, serverColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = server.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Servers(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testServersAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	serverOne := &Server{}
	serverTwo := &Server{}
	if err = randomize.Struct(seed, serverOne, serverDBTypes, false, serverColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}
	if err = randomize.Struct(seed, serverTwo, serverDBTypes, false, serverColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = serverOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = serverTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Servers(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testServersCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	serverOne := &Server{}
	serverTwo := &Server{}
	if err = randomize.Struct(seed, serverOne, serverDBTypes, false, serverColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}
	if err = randomize.Struct(seed, serverTwo, serverDBTypes, false, serverColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = serverOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = serverTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Servers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testServersInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	server := &Server{}
	if err = randomize.Struct(seed, server, serverDBTypes, true, serverColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = server.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Servers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testServersInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	server := &Server{}
	if err = randomize.Struct(seed, server, serverDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = server.Insert(tx, serverColumns...); err != nil {
		t.Error(err)
	}

	count, err := Servers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testServerToManyServernameFileAssets(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Server
	var b, c FileAsset

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, serverDBTypes, true, serverColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, fileAssetDBTypes, false, fileAssetColumnsWithDefault...)
	randomize.Struct(seed, &c, fileAssetDBTypes, false, fileAssetColumnsWithDefault...)

	b.ServernameID.Valid = true
	c.ServernameID.Valid = true
	b.ServernameID.String = a.Servername
	c.ServernameID.String = a.Servername
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	fileAsset, err := a.ServernameFileAssets(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range fileAsset {
		if v.ServernameID.String == b.ServernameID.String {
			bFound = true
		}
		if v.ServernameID.String == c.ServernameID.String {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ServerSlice{&a}
	if err = a.L.LoadServernameFileAssets(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ServernameFileAssets); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.ServernameFileAssets = nil
	if err = a.L.LoadServernameFileAssets(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ServernameFileAssets); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", fileAsset)
	}
}

func testServerToManyAddOpServernameFileAssets(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Server
	var b, c, d, e FileAsset

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, serverDBTypes, false, strmangle.SetComplement(serverPrimaryKeyColumns, serverColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*FileAsset{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*FileAsset{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddServernameFileAssets(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.Servername != first.ServernameID.String {
			t.Error("foreign key was wrong value", a.Servername, first.ServernameID.String)
		}
		if a.Servername != second.ServernameID.String {
			t.Error("foreign key was wrong value", a.Servername, second.ServernameID.String)
		}

		if first.R.Servername != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Servername != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.ServernameFileAssets[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.ServernameFileAssets[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.ServernameFileAssets(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testServerToManySetOpServernameFileAssets(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Server
	var b, c, d, e FileAsset

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, serverDBTypes, false, strmangle.SetComplement(serverPrimaryKeyColumns, serverColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*FileAsset{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
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

	err = a.SetServernameFileAssets(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ServernameFileAssets(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetServernameFileAssets(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ServernameFileAssets(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.ServernameID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.ServernameID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.Servername != d.ServernameID.String {
		t.Error("foreign key was wrong value", a.Servername, d.ServernameID.String)
	}
	if a.Servername != e.ServernameID.String {
		t.Error("foreign key was wrong value", a.Servername, e.ServernameID.String)
	}

	if b.R.Servername != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Servername != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Servername != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Servername != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.ServernameFileAssets[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.ServernameFileAssets[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testServerToManyRemoveOpServernameFileAssets(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Server
	var b, c, d, e FileAsset

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, serverDBTypes, false, strmangle.SetComplement(serverPrimaryKeyColumns, serverColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*FileAsset{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, fileAssetDBTypes, false, strmangle.SetComplement(fileAssetPrimaryKeyColumns, fileAssetColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddServernameFileAssets(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ServernameFileAssets(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveServernameFileAssets(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ServernameFileAssets(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.ServernameID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.ServernameID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.Servername != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Servername != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Servername != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.Servername != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.ServernameFileAssets) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.ServernameFileAssets[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.ServernameFileAssets[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testServersReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	server := &Server{}
	if err = randomize.Struct(seed, server, serverDBTypes, true, serverColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = server.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = server.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testServersReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	server := &Server{}
	if err = randomize.Struct(seed, server, serverDBTypes, true, serverColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = server.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ServerSlice{server}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testServersSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	server := &Server{}
	if err = randomize.Struct(seed, server, serverDBTypes, true, serverColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = server.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Servers(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	serverDBTypes = map[string]string{`Created`: `timestamp without time zone`, `Httpurl`: `character varying`, `ID`: `integer`, `Lastuser`: `character varying`, `Path`: `character varying`, `Servername`: `character varying`, `Updated`: `timestamp without time zone`}
	_             = bytes.MinRead
)

func testServersUpdate(t *testing.T) {
	t.Parallel()

	if len(serverColumns) == len(serverPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	server := &Server{}
	if err = randomize.Struct(seed, server, serverDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = server.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Servers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, server, serverDBTypes, true, serverColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	if err = server.Update(tx); err != nil {
		t.Error(err)
	}
}

func testServersSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(serverColumns) == len(serverPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	server := &Server{}
	if err = randomize.Struct(seed, server, serverDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = server.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Servers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, server, serverDBTypes, true, serverPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(serverColumns, serverPrimaryKeyColumns) {
		fields = serverColumns
	} else {
		fields = strmangle.SetComplement(
			serverColumns,
			serverPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(server))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := ServerSlice{server}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testServersUpsert(t *testing.T) {
	t.Parallel()

	if len(serverColumns) == len(serverPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	server := Server{}
	if err = randomize.Struct(seed, &server, serverDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = server.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Server: %s", err)
	}

	count, err := Servers(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &server, serverDBTypes, false, serverPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Server struct: %s", err)
	}

	if err = server.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Server: %s", err)
	}

	count, err = Servers(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
