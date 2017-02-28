package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testContentRoles(t *testing.T) {
	t.Parallel()

	query := ContentRoles(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testContentRolesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentRole := &ContentRole{}
	if err = randomize.Struct(seed, contentRole, contentRoleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRole.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = contentRole.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentRoles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContentRolesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentRole := &ContentRole{}
	if err = randomize.Struct(seed, contentRole, contentRoleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRole.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = ContentRoles(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := ContentRoles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContentRolesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentRole := &ContentRole{}
	if err = randomize.Struct(seed, contentRole, contentRoleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRole.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ContentRoleSlice{contentRole}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentRoles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testContentRolesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentRole := &ContentRole{}
	if err = randomize.Struct(seed, contentRole, contentRoleDBTypes, true, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRole.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := ContentRoleExists(tx, contentRole.ID)
	if err != nil {
		t.Errorf("Unable to check if ContentRole exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ContentRoleExistsG to return true, but got false.")
	}
}
func testContentRolesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentRole := &ContentRole{}
	if err = randomize.Struct(seed, contentRole, contentRoleDBTypes, true, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRole.Insert(tx); err != nil {
		t.Error(err)
	}

	contentRoleFound, err := FindContentRole(tx, contentRole.ID)
	if err != nil {
		t.Error(err)
	}

	if contentRoleFound == nil {
		t.Error("want a record, got nil")
	}
}
func testContentRolesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentRole := &ContentRole{}
	if err = randomize.Struct(seed, contentRole, contentRoleDBTypes, true, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRole.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = ContentRoles(tx).Bind(contentRole); err != nil {
		t.Error(err)
	}
}

func testContentRolesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentRole := &ContentRole{}
	if err = randomize.Struct(seed, contentRole, contentRoleDBTypes, true, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRole.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := ContentRoles(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testContentRolesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentRoleOne := &ContentRole{}
	contentRoleTwo := &ContentRole{}
	if err = randomize.Struct(seed, contentRoleOne, contentRoleDBTypes, false, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}
	if err = randomize.Struct(seed, contentRoleTwo, contentRoleDBTypes, false, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRoleOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = contentRoleTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := ContentRoles(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testContentRolesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	contentRoleOne := &ContentRole{}
	contentRoleTwo := &ContentRole{}
	if err = randomize.Struct(seed, contentRoleOne, contentRoleDBTypes, false, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}
	if err = randomize.Struct(seed, contentRoleTwo, contentRoleDBTypes, false, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRoleOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = contentRoleTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentRoles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testContentRolesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentRole := &ContentRole{}
	if err = randomize.Struct(seed, contentRole, contentRoleDBTypes, true, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRole.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentRoles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testContentRolesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentRole := &ContentRole{}
	if err = randomize.Struct(seed, contentRole, contentRoleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRole.Insert(tx, contentRoleColumns...); err != nil {
		t.Error(err)
	}

	count, err := ContentRoles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testContentRoleToManyRoleContentUnitsPersons(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentRole
	var b, c ContentUnitsPerson

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentRoleDBTypes, true, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, contentUnitsPersonDBTypes, false, contentUnitsPersonColumnsWithDefault...)
	randomize.Struct(seed, &c, contentUnitsPersonDBTypes, false, contentUnitsPersonColumnsWithDefault...)

	b.RoleID = a.ID
	c.RoleID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	contentUnitsPerson, err := a.RoleContentUnitsPersons(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range contentUnitsPerson {
		if v.RoleID == b.RoleID {
			bFound = true
		}
		if v.RoleID == c.RoleID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ContentRoleSlice{&a}
	if err = a.L.LoadRoleContentUnitsPersons(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.RoleContentUnitsPersons); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.RoleContentUnitsPersons = nil
	if err = a.L.LoadRoleContentUnitsPersons(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.RoleContentUnitsPersons); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", contentUnitsPerson)
	}
}

func testContentRoleToManyAddOpRoleContentUnitsPersons(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentRole
	var b, c, d, e ContentUnitsPerson

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentRoleDBTypes, false, strmangle.SetComplement(contentRolePrimaryKeyColumns, contentRoleColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentUnitsPerson{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentUnitsPersonDBTypes, false, strmangle.SetComplement(contentUnitsPersonPrimaryKeyColumns, contentUnitsPersonColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*ContentUnitsPerson{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddRoleContentUnitsPersons(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.RoleID {
			t.Error("foreign key was wrong value", a.ID, first.RoleID)
		}
		if a.ID != second.RoleID {
			t.Error("foreign key was wrong value", a.ID, second.RoleID)
		}

		if first.R.Role != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Role != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.RoleContentUnitsPersons[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.RoleContentUnitsPersons[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.RoleContentUnitsPersons(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testContentRoleToOneStringTranslationUsingDescription(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local ContentRole
	var foreign StringTranslation

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, contentRoleDBTypes, true, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	local.DescriptionID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.DescriptionID.Int64 = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Description(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ContentRoleSlice{&local}
	if err = local.L.LoadDescription(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Description == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Description = nil
	if err = local.L.LoadDescription(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Description == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testContentRoleToOneStringTranslationUsingName(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local ContentRole
	var foreign StringTranslation

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, contentRoleDBTypes, true, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.NameID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Name(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ContentRoleSlice{&local}
	if err = local.L.LoadName(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Name == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Name = nil
	if err = local.L.LoadName(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Name == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testContentRoleToOneSetOpStringTranslationUsingDescription(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentRole
	var b, c StringTranslation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentRoleDBTypes, false, strmangle.SetComplement(contentRolePrimaryKeyColumns, contentRoleColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*StringTranslation{&b, &c} {
		err = a.SetDescription(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Description != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.DescriptionContentRoles[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.DescriptionID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.DescriptionID.Int64)
		}

		zero := reflect.Zero(reflect.TypeOf(a.DescriptionID.Int64))
		reflect.Indirect(reflect.ValueOf(&a.DescriptionID.Int64)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.DescriptionID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.DescriptionID.Int64, x.ID)
		}
	}
}

func testContentRoleToOneRemoveOpStringTranslationUsingDescription(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentRole
	var b StringTranslation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentRoleDBTypes, false, strmangle.SetComplement(contentRolePrimaryKeyColumns, contentRoleColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.SetDescription(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveDescription(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.Description(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.Description != nil {
		t.Error("R struct entry should be nil")
	}

	if a.DescriptionID.Valid {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.DescriptionContentRoles) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testContentRoleToOneSetOpStringTranslationUsingName(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentRole
	var b, c StringTranslation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentRoleDBTypes, false, strmangle.SetComplement(contentRolePrimaryKeyColumns, contentRoleColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*StringTranslation{&b, &c} {
		err = a.SetName(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Name != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.NameContentRoles[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.NameID != x.ID {
			t.Error("foreign key was wrong value", a.NameID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.NameID))
		reflect.Indirect(reflect.ValueOf(&a.NameID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.NameID != x.ID {
			t.Error("foreign key was wrong value", a.NameID, x.ID)
		}
	}
}
func testContentRolesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentRole := &ContentRole{}
	if err = randomize.Struct(seed, contentRole, contentRoleDBTypes, true, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRole.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = contentRole.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testContentRolesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentRole := &ContentRole{}
	if err = randomize.Struct(seed, contentRole, contentRoleDBTypes, true, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRole.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ContentRoleSlice{contentRole}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testContentRolesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentRole := &ContentRole{}
	if err = randomize.Struct(seed, contentRole, contentRoleDBTypes, true, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRole.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := ContentRoles(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	contentRoleDBTypes = map[string]string{`DescriptionID`: `bigint`, `ID`: `bigint`, `NameID`: `bigint`}
	_                  = bytes.MinRead
)

func testContentRolesUpdate(t *testing.T) {
	t.Parallel()

	if len(contentRoleColumns) == len(contentRolePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	contentRole := &ContentRole{}
	if err = randomize.Struct(seed, contentRole, contentRoleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRole.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentRoles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, contentRole, contentRoleDBTypes, true, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	if err = contentRole.Update(tx); err != nil {
		t.Error(err)
	}
}

func testContentRolesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(contentRoleColumns) == len(contentRolePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	contentRole := &ContentRole{}
	if err = randomize.Struct(seed, contentRole, contentRoleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRole.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentRoles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, contentRole, contentRoleDBTypes, true, contentRolePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(contentRoleColumns, contentRolePrimaryKeyColumns) {
		fields = contentRoleColumns
	} else {
		fields = strmangle.SetComplement(
			contentRoleColumns,
			contentRolePrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(contentRole))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := ContentRoleSlice{contentRole}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testContentRolesUpsert(t *testing.T) {
	t.Parallel()

	if len(contentRoleColumns) == len(contentRolePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	contentRole := ContentRole{}
	if err = randomize.Struct(seed, &contentRole, contentRoleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentRole.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert ContentRole: %s", err)
	}

	count, err := ContentRoles(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &contentRole, contentRoleDBTypes, false, contentRolePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	if err = contentRole.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert ContentRole: %s", err)
	}

	count, err = ContentRoles(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
