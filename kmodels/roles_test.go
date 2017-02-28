package kmodels

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testRoles(t *testing.T) {
	t.Parallel()

	query := Roles(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testRolesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	role := &Role{}
	if err = randomize.Struct(seed, role, roleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = role.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = role.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Roles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testRolesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	role := &Role{}
	if err = randomize.Struct(seed, role, roleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = role.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Roles(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Roles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testRolesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	role := &Role{}
	if err = randomize.Struct(seed, role, roleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = role.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := RoleSlice{role}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Roles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testRolesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	role := &Role{}
	if err = randomize.Struct(seed, role, roleDBTypes, true, roleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = role.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := RoleExists(tx, role.ID)
	if err != nil {
		t.Errorf("Unable to check if Role exists: %s", err)
	}
	if !e {
		t.Errorf("Expected RoleExistsG to return true, but got false.")
	}
}
func testRolesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	role := &Role{}
	if err = randomize.Struct(seed, role, roleDBTypes, true, roleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = role.Insert(tx); err != nil {
		t.Error(err)
	}

	roleFound, err := FindRole(tx, role.ID)
	if err != nil {
		t.Error(err)
	}

	if roleFound == nil {
		t.Error("want a record, got nil")
	}
}
func testRolesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	role := &Role{}
	if err = randomize.Struct(seed, role, roleDBTypes, true, roleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = role.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Roles(tx).Bind(role); err != nil {
		t.Error(err)
	}
}

func testRolesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	role := &Role{}
	if err = randomize.Struct(seed, role, roleDBTypes, true, roleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = role.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Roles(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testRolesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	roleOne := &Role{}
	roleTwo := &Role{}
	if err = randomize.Struct(seed, roleOne, roleDBTypes, false, roleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}
	if err = randomize.Struct(seed, roleTwo, roleDBTypes, false, roleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = roleOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = roleTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Roles(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testRolesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	roleOne := &Role{}
	roleTwo := &Role{}
	if err = randomize.Struct(seed, roleOne, roleDBTypes, false, roleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}
	if err = randomize.Struct(seed, roleTwo, roleDBTypes, false, roleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = roleOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = roleTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Roles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testRolesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	role := &Role{}
	if err = randomize.Struct(seed, role, roleDBTypes, true, roleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = role.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Roles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testRolesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	role := &Role{}
	if err = randomize.Struct(seed, role, roleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = role.Insert(tx, roleColumns...); err != nil {
		t.Error(err)
	}

	count, err := Roles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testRolesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	role := &Role{}
	if err = randomize.Struct(seed, role, roleDBTypes, true, roleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = role.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = role.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testRolesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	role := &Role{}
	if err = randomize.Struct(seed, role, roleDBTypes, true, roleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = role.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := RoleSlice{role}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testRolesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	role := &Role{}
	if err = randomize.Struct(seed, role, roleDBTypes, true, roleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = role.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Roles(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	roleDBTypes = map[string]string{`CreatedAt`: `timestamp without time zone`, `Description`: `character varying`, `ID`: `integer`, `Name`: `character varying`, `UpdatedAt`: `timestamp without time zone`}
	_           = bytes.MinRead
)

func testRolesUpdate(t *testing.T) {
	t.Parallel()

	if len(roleColumns) == len(rolePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	role := &Role{}
	if err = randomize.Struct(seed, role, roleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = role.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Roles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, role, roleDBTypes, true, roleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	if err = role.Update(tx); err != nil {
		t.Error(err)
	}
}

func testRolesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(roleColumns) == len(rolePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	role := &Role{}
	if err = randomize.Struct(seed, role, roleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = role.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Roles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, role, roleDBTypes, true, rolePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(roleColumns, rolePrimaryKeyColumns) {
		fields = roleColumns
	} else {
		fields = strmangle.SetComplement(
			roleColumns,
			rolePrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(role))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := RoleSlice{role}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testRolesUpsert(t *testing.T) {
	t.Parallel()

	if len(roleColumns) == len(rolePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	role := Role{}
	if err = randomize.Struct(seed, &role, roleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = role.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Role: %s", err)
	}

	count, err := Roles(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &role, roleDBTypes, false, rolePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Role struct: %s", err)
	}

	if err = role.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Role: %s", err)
	}

	count, err = Roles(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
