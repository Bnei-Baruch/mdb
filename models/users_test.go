package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testUsers(t *testing.T) {
	t.Parallel()

	query := Users(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testUsersDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = user.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUsersQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Users(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUsersSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := UserSlice{user}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testUsersExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := UserExists(tx, user.ID)
	if err != nil {
		t.Errorf("Unable to check if User exists: %s", err)
	}
	if !e {
		t.Errorf("Expected UserExistsG to return true, but got false.")
	}
}
func testUsersFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	userFound, err := FindUser(tx, user.ID)
	if err != nil {
		t.Error(err)
	}

	if userFound == nil {
		t.Error("want a record, got nil")
	}
}
func testUsersBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Users(tx).Bind(user); err != nil {
		t.Error(err)
	}
}

func testUsersOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Users(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testUsersAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	userOne := &User{}
	userTwo := &User{}
	if err = randomize.Struct(seed, userOne, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}
	if err = randomize.Struct(seed, userTwo, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = userOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = userTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Users(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testUsersCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	userOne := &User{}
	userTwo := &User{}
	if err = randomize.Struct(seed, userOne, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}
	if err = randomize.Struct(seed, userTwo, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = userOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = userTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testUsersInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testUsersInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx, userColumns...); err != nil {
		t.Error(err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testUserToManyOperations(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a User
	var b, c Operation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, operationDBTypes, false, operationColumnsWithDefault...)
	randomize.Struct(seed, &c, operationDBTypes, false, operationColumnsWithDefault...)

	b.UserID.Valid = true
	c.UserID.Valid = true
	b.UserID.Int64 = a.ID
	c.UserID.Int64 = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	operation, err := a.Operations(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range operation {
		if v.UserID.Int64 == b.UserID.Int64 {
			bFound = true
		}
		if v.UserID.Int64 == c.UserID.Int64 {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := UserSlice{&a}
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

func testUserToManyCollectionI18ns(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a User
	var b, c CollectionI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, collectionI18nDBTypes, false, collectionI18nColumnsWithDefault...)
	randomize.Struct(seed, &c, collectionI18nDBTypes, false, collectionI18nColumnsWithDefault...)

	b.UserID.Valid = true
	c.UserID.Valid = true
	b.UserID.Int64 = a.ID
	c.UserID.Int64 = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	collectionI18n, err := a.CollectionI18ns(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range collectionI18n {
		if v.UserID.Int64 == b.UserID.Int64 {
			bFound = true
		}
		if v.UserID.Int64 == c.UserID.Int64 {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := UserSlice{&a}
	if err = a.L.LoadCollectionI18ns(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.CollectionI18ns); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.CollectionI18ns = nil
	if err = a.L.LoadCollectionI18ns(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.CollectionI18ns); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", collectionI18n)
	}
}

func testUserToManyContentUnitI18ns(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a User
	var b, c ContentUnitI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, contentUnitI18nDBTypes, false, contentUnitI18nColumnsWithDefault...)
	randomize.Struct(seed, &c, contentUnitI18nDBTypes, false, contentUnitI18nColumnsWithDefault...)

	b.UserID.Valid = true
	c.UserID.Valid = true
	b.UserID.Int64 = a.ID
	c.UserID.Int64 = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	contentUnitI18n, err := a.ContentUnitI18ns(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range contentUnitI18n {
		if v.UserID.Int64 == b.UserID.Int64 {
			bFound = true
		}
		if v.UserID.Int64 == c.UserID.Int64 {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := UserSlice{&a}
	if err = a.L.LoadContentUnitI18ns(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ContentUnitI18ns); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.ContentUnitI18ns = nil
	if err = a.L.LoadContentUnitI18ns(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ContentUnitI18ns); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", contentUnitI18n)
	}
}

func testUserToManyTagsI18ns(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a User
	var b, c TagsI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, tagsI18nDBTypes, false, tagsI18nColumnsWithDefault...)
	randomize.Struct(seed, &c, tagsI18nDBTypes, false, tagsI18nColumnsWithDefault...)

	b.UserID.Valid = true
	c.UserID.Valid = true
	b.UserID.Int64 = a.ID
	c.UserID.Int64 = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	tagsI18n, err := a.TagsI18ns(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range tagsI18n {
		if v.UserID.Int64 == b.UserID.Int64 {
			bFound = true
		}
		if v.UserID.Int64 == c.UserID.Int64 {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := UserSlice{&a}
	if err = a.L.LoadTagsI18ns(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TagsI18ns); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.TagsI18ns = nil
	if err = a.L.LoadTagsI18ns(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TagsI18ns); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", tagsI18n)
	}
}

func testUserToManyAddOpOperations(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a User
	var b, c, d, e Operation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
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

		if a.ID != first.UserID.Int64 {
			t.Error("foreign key was wrong value", a.ID, first.UserID.Int64)
		}
		if a.ID != second.UserID.Int64 {
			t.Error("foreign key was wrong value", a.ID, second.UserID.Int64)
		}

		if first.R.User != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.User != &a {
			t.Error("relationship was not added properly to the foreign slice")
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

func testUserToManySetOpOperations(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a User
	var b, c, d, e Operation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
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

	if b.UserID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.UserID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.ID != d.UserID.Int64 {
		t.Error("foreign key was wrong value", a.ID, d.UserID.Int64)
	}
	if a.ID != e.UserID.Int64 {
		t.Error("foreign key was wrong value", a.ID, e.UserID.Int64)
	}

	if b.R.User != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.User != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.User != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.User != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.Operations[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Operations[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testUserToManyRemoveOpOperations(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a User
	var b, c, d, e Operation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
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

	if b.UserID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.UserID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.User != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.User != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.User != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.User != &a {
		t.Error("relationship to a should have been preserved")
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

func testUserToManyAddOpCollectionI18ns(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a User
	var b, c, d, e CollectionI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*CollectionI18n{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, collectionI18nDBTypes, false, strmangle.SetComplement(collectionI18nPrimaryKeyColumns, collectionI18nColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*CollectionI18n{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddCollectionI18ns(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.UserID.Int64 {
			t.Error("foreign key was wrong value", a.ID, first.UserID.Int64)
		}
		if a.ID != second.UserID.Int64 {
			t.Error("foreign key was wrong value", a.ID, second.UserID.Int64)
		}

		if first.R.User != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.User != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.CollectionI18ns[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.CollectionI18ns[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.CollectionI18ns(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testUserToManySetOpCollectionI18ns(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a User
	var b, c, d, e CollectionI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*CollectionI18n{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, collectionI18nDBTypes, false, strmangle.SetComplement(collectionI18nPrimaryKeyColumns, collectionI18nColumnsWithoutDefault)...); err != nil {
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

	err = a.SetCollectionI18ns(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.CollectionI18ns(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetCollectionI18ns(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.CollectionI18ns(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.UserID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.UserID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.ID != d.UserID.Int64 {
		t.Error("foreign key was wrong value", a.ID, d.UserID.Int64)
	}
	if a.ID != e.UserID.Int64 {
		t.Error("foreign key was wrong value", a.ID, e.UserID.Int64)
	}

	if b.R.User != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.User != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.User != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.User != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.CollectionI18ns[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.CollectionI18ns[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testUserToManyRemoveOpCollectionI18ns(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a User
	var b, c, d, e CollectionI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*CollectionI18n{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, collectionI18nDBTypes, false, strmangle.SetComplement(collectionI18nPrimaryKeyColumns, collectionI18nColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddCollectionI18ns(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.CollectionI18ns(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveCollectionI18ns(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.CollectionI18ns(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.UserID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.UserID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.User != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.User != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.User != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.User != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.CollectionI18ns) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.CollectionI18ns[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.CollectionI18ns[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testUserToManyAddOpContentUnitI18ns(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a User
	var b, c, d, e ContentUnitI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentUnitI18n{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentUnitI18nDBTypes, false, strmangle.SetComplement(contentUnitI18nPrimaryKeyColumns, contentUnitI18nColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*ContentUnitI18n{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddContentUnitI18ns(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.UserID.Int64 {
			t.Error("foreign key was wrong value", a.ID, first.UserID.Int64)
		}
		if a.ID != second.UserID.Int64 {
			t.Error("foreign key was wrong value", a.ID, second.UserID.Int64)
		}

		if first.R.User != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.User != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.ContentUnitI18ns[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.ContentUnitI18ns[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.ContentUnitI18ns(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testUserToManySetOpContentUnitI18ns(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a User
	var b, c, d, e ContentUnitI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentUnitI18n{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentUnitI18nDBTypes, false, strmangle.SetComplement(contentUnitI18nPrimaryKeyColumns, contentUnitI18nColumnsWithoutDefault)...); err != nil {
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

	err = a.SetContentUnitI18ns(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ContentUnitI18ns(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetContentUnitI18ns(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ContentUnitI18ns(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.UserID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.UserID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.ID != d.UserID.Int64 {
		t.Error("foreign key was wrong value", a.ID, d.UserID.Int64)
	}
	if a.ID != e.UserID.Int64 {
		t.Error("foreign key was wrong value", a.ID, e.UserID.Int64)
	}

	if b.R.User != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.User != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.User != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.User != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.ContentUnitI18ns[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.ContentUnitI18ns[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testUserToManyRemoveOpContentUnitI18ns(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a User
	var b, c, d, e ContentUnitI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentUnitI18n{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentUnitI18nDBTypes, false, strmangle.SetComplement(contentUnitI18nPrimaryKeyColumns, contentUnitI18nColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddContentUnitI18ns(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.ContentUnitI18ns(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveContentUnitI18ns(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.ContentUnitI18ns(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.UserID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.UserID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.User != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.User != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.User != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.User != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.ContentUnitI18ns) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.ContentUnitI18ns[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.ContentUnitI18ns[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testUserToManyAddOpTagsI18ns(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a User
	var b, c, d, e TagsI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*TagsI18n{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, tagsI18nDBTypes, false, strmangle.SetComplement(tagsI18nPrimaryKeyColumns, tagsI18nColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*TagsI18n{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddTagsI18ns(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.UserID.Int64 {
			t.Error("foreign key was wrong value", a.ID, first.UserID.Int64)
		}
		if a.ID != second.UserID.Int64 {
			t.Error("foreign key was wrong value", a.ID, second.UserID.Int64)
		}

		if first.R.User != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.User != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.TagsI18ns[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.TagsI18ns[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.TagsI18ns(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testUserToManySetOpTagsI18ns(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a User
	var b, c, d, e TagsI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*TagsI18n{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, tagsI18nDBTypes, false, strmangle.SetComplement(tagsI18nPrimaryKeyColumns, tagsI18nColumnsWithoutDefault)...); err != nil {
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

	err = a.SetTagsI18ns(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.TagsI18ns(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetTagsI18ns(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.TagsI18ns(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.UserID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.UserID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.ID != d.UserID.Int64 {
		t.Error("foreign key was wrong value", a.ID, d.UserID.Int64)
	}
	if a.ID != e.UserID.Int64 {
		t.Error("foreign key was wrong value", a.ID, e.UserID.Int64)
	}

	if b.R.User != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.User != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.User != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.User != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.TagsI18ns[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.TagsI18ns[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testUserToManyRemoveOpTagsI18ns(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a User
	var b, c, d, e TagsI18n

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*TagsI18n{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, tagsI18nDBTypes, false, strmangle.SetComplement(tagsI18nPrimaryKeyColumns, tagsI18nColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddTagsI18ns(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.TagsI18ns(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveTagsI18ns(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.TagsI18ns(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.UserID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.UserID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.User != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.User != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.User != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.User != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.TagsI18ns) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.TagsI18ns[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.TagsI18ns[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testUsersReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = user.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testUsersReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := UserSlice{user}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testUsersSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Users(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	userDBTypes = map[string]string{`Comments`: `character varying`, `CreatedAt`: `timestamp with time zone`, `DeletedAt`: `timestamp with time zone`, `Email`: `character varying`, `ID`: `bigint`, `Name`: `character`, `Phone`: `character varying`, `UpdatedAt`: `timestamp with time zone`}
	_           = bytes.MinRead
)

func testUsersUpdate(t *testing.T) {
	t.Parallel()

	if len(userColumns) == len(userPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, user, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err = user.Update(tx); err != nil {
		t.Error(err)
	}
}

func testUsersSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(userColumns) == len(userPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	user := &User{}
	if err = randomize.Struct(seed, user, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, user, userDBTypes, true, userPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(userColumns, userPrimaryKeyColumns) {
		fields = userColumns
	} else {
		fields = strmangle.SetComplement(
			userColumns,
			userPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(user))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := UserSlice{user}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testUsersUpsert(t *testing.T) {
	t.Parallel()

	if len(userColumns) == len(userPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	user := User{}
	if err = randomize.Struct(seed, &user, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = user.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert User: %s", err)
	}

	count, err := Users(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &user, userDBTypes, false, userPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err = user.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert User: %s", err)
	}

	count, err = Users(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
