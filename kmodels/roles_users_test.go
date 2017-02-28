package kmodels

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testRolesUsers(t *testing.T) {
	t.Parallel()

	query := RolesUsers(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testRolesUsersDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	rolesUser := &RolesUser{}
	if err = randomize.Struct(seed, rolesUser, rolesUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUser.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = rolesUser.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := RolesUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testRolesUsersQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	rolesUser := &RolesUser{}
	if err = randomize.Struct(seed, rolesUser, rolesUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUser.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = RolesUsers(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := RolesUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testRolesUsersSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	rolesUser := &RolesUser{}
	if err = randomize.Struct(seed, rolesUser, rolesUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUser.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := RolesUserSlice{rolesUser}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := RolesUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testRolesUsersExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	rolesUser := &RolesUser{}
	if err = randomize.Struct(seed, rolesUser, rolesUserDBTypes, true, rolesUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUser.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := RolesUserExists(tx, rolesUser.RoleID, rolesUser.UserID)
	if err != nil {
		t.Errorf("Unable to check if RolesUser exists: %s", err)
	}
	if !e {
		t.Errorf("Expected RolesUserExistsG to return true, but got false.")
	}
}
func testRolesUsersFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	rolesUser := &RolesUser{}
	if err = randomize.Struct(seed, rolesUser, rolesUserDBTypes, true, rolesUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUser.Insert(tx); err != nil {
		t.Error(err)
	}

	rolesUserFound, err := FindRolesUser(tx, rolesUser.RoleID, rolesUser.UserID)
	if err != nil {
		t.Error(err)
	}

	if rolesUserFound == nil {
		t.Error("want a record, got nil")
	}
}
func testRolesUsersBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	rolesUser := &RolesUser{}
	if err = randomize.Struct(seed, rolesUser, rolesUserDBTypes, true, rolesUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUser.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = RolesUsers(tx).Bind(rolesUser); err != nil {
		t.Error(err)
	}
}

func testRolesUsersOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	rolesUser := &RolesUser{}
	if err = randomize.Struct(seed, rolesUser, rolesUserDBTypes, true, rolesUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUser.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := RolesUsers(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testRolesUsersAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	rolesUserOne := &RolesUser{}
	rolesUserTwo := &RolesUser{}
	if err = randomize.Struct(seed, rolesUserOne, rolesUserDBTypes, false, rolesUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}
	if err = randomize.Struct(seed, rolesUserTwo, rolesUserDBTypes, false, rolesUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUserOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = rolesUserTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := RolesUsers(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testRolesUsersCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	rolesUserOne := &RolesUser{}
	rolesUserTwo := &RolesUser{}
	if err = randomize.Struct(seed, rolesUserOne, rolesUserDBTypes, false, rolesUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}
	if err = randomize.Struct(seed, rolesUserTwo, rolesUserDBTypes, false, rolesUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUserOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = rolesUserTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := RolesUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testRolesUsersInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	rolesUser := &RolesUser{}
	if err = randomize.Struct(seed, rolesUser, rolesUserDBTypes, true, rolesUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUser.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := RolesUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testRolesUsersInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	rolesUser := &RolesUser{}
	if err = randomize.Struct(seed, rolesUser, rolesUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUser.Insert(tx, rolesUserColumns...); err != nil {
		t.Error(err)
	}

	count, err := RolesUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testRolesUserToOneUserUsingUser(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local RolesUser
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, rolesUserDBTypes, true, rolesUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.UserID = foreign.ID
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

	slice := RolesUserSlice{&local}
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

func testRolesUserToOneSetOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a RolesUser
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, rolesUserDBTypes, false, strmangle.SetComplement(rolesUserPrimaryKeyColumns, rolesUserColumnsWithoutDefault)...); err != nil {
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

		if x.R.RolesUsers[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.UserID != x.ID {
			t.Error("foreign key was wrong value", a.UserID)
		}

		if exists, err := RolesUserExists(tx, a.RoleID, a.UserID); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}
func testRolesUsersReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	rolesUser := &RolesUser{}
	if err = randomize.Struct(seed, rolesUser, rolesUserDBTypes, true, rolesUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUser.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = rolesUser.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testRolesUsersReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	rolesUser := &RolesUser{}
	if err = randomize.Struct(seed, rolesUser, rolesUserDBTypes, true, rolesUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUser.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := RolesUserSlice{rolesUser}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testRolesUsersSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	rolesUser := &RolesUser{}
	if err = randomize.Struct(seed, rolesUser, rolesUserDBTypes, true, rolesUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUser.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := RolesUsers(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	rolesUserDBTypes = map[string]string{`RoleID`: `integer`, `UserID`: `integer`}
	_                = bytes.MinRead
)

func testRolesUsersUpdate(t *testing.T) {
	t.Parallel()

	if len(rolesUserColumns) == len(rolesUserPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	rolesUser := &RolesUser{}
	if err = randomize.Struct(seed, rolesUser, rolesUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUser.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := RolesUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, rolesUser, rolesUserDBTypes, true, rolesUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	if err = rolesUser.Update(tx); err != nil {
		t.Error(err)
	}
}

func testRolesUsersSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(rolesUserColumns) == len(rolesUserPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	rolesUser := &RolesUser{}
	if err = randomize.Struct(seed, rolesUser, rolesUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUser.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := RolesUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, rolesUser, rolesUserDBTypes, true, rolesUserPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(rolesUserColumns, rolesUserPrimaryKeyColumns) {
		fields = rolesUserColumns
	} else {
		fields = strmangle.SetComplement(
			rolesUserColumns,
			rolesUserPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(rolesUser))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := RolesUserSlice{rolesUser}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testRolesUsersUpsert(t *testing.T) {
	t.Parallel()

	if len(rolesUserColumns) == len(rolesUserPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	rolesUser := RolesUser{}
	if err = randomize.Struct(seed, &rolesUser, rolesUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = rolesUser.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert RolesUser: %s", err)
	}

	count, err := RolesUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &rolesUser, rolesUserDBTypes, false, rolesUserPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize RolesUser struct: %s", err)
	}

	if err = rolesUser.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert RolesUser: %s", err)
	}

	count, err = RolesUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
