package api

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/casbin/casbin"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/permissions"
	"github.com/Bnei-Baruch/mdb/utils"
)

type RestSuite struct {
	suite.Suite
	common.TestDBManager
	tx *sql.Tx
}

func (suite *RestSuite) SetupSuite() {
	suite.Require().Nil(suite.InitTestDB())
	suite.Require().Nil(common.InitTypeRegistries(suite.DB))
}

func (suite *RestSuite) TearDownSuite() {
	suite.Require().Nil(suite.DestroyTestDB())
}

func (suite *RestSuite) SetupTest() {
	var err error
	suite.tx, err = suite.DB.Begin()
	suite.Require().Nil(err)
}

func (suite *RestSuite) TearDownTest() {
	err := suite.tx.Rollback()
	suite.Require().Nil(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRest(t *testing.T) {
	suite.Run(t, new(RestSuite))
}

func (suite *RestSuite) TestCollectionsList() {
	cp := new(DummyAuthProvider)

	req := CollectionsRequest{
		ListRequest: ListRequest{StartIndex: 1, StopIndex: 5},
	}
	resp, err := handleCollectionsList(cp, suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(0, resp.Total, "empty total")
	suite.Empty(resp.Collections, "empty data")

	collections := createDummyCollections(suite.tx, 10)

	resp, err = handleCollectionsList(cp, suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.Collections {
		suite.assertEqualDummyCollection(collections[i], x, i)
	}

	req.StartIndex = 6
	req.StopIndex = 10
	resp, err = handleCollectionsList(cp, suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.Collections {
		suite.assertEqualDummyCollection(collections[i+5], x, i+5)
	}

	// test IDsFilter, UIDsFilter
	ids := make([]int64, len(collections))
	uids := make([]string, len(collections))
	for i := range collections {
		c := collections[i]
		ids[i] = c.ID
		uids[i] = c.UID
	}

	req = CollectionsRequest{
		IDsFilter: IDsFilter{
			IDs: ids,
		},
	}
	resp, err = handleCollectionsList(cp, suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.Collections {
		suite.assertEqualDummyCollection(collections[i], x, i)
	}

	req = CollectionsRequest{
		UIDsFilter: UIDsFilter{
			UIDs: uids,
		},
	}
	resp, err = handleCollectionsList(cp, suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.Collections {
		suite.assertEqualDummyCollection(collections[i], x, i)
	}
}

func (suite *RestSuite) TestCollectionItem() {
	cp := new(DummyAuthProvider)
	c, err := handleGetCollection(cp, suite.tx, 1)
	suite.Nil(c, "Collection nil")
	suite.Require().NotNil(err, "Not Found error")
	suite.Equal(http.StatusNotFound, err.Code, "Error http status code")
	suite.Equal(gin.ErrorTypePublic, err.Type, "Error gin type")

	collections := createDummyCollections(suite.tx, 3)
	for i, c := range collections {
		x, err := handleGetCollection(cp, suite.tx, c.ID)
		suite.Require().Nil(err, "Collection item err [%d]", i)
		suite.assertEqualDummyCollection(c, x, i)
	}
}

func (suite *RestSuite) TestContentUnitsList() {
	cp := new(DummyAuthProvider)
	req := ContentUnitsRequest{
		ListRequest: ListRequest{StartIndex: 1, StopIndex: 5},
	}
	resp, err := handleContentUnitsList(cp, suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(0, resp.Total, "empty total")
	suite.Empty(resp.ContentUnits, "empty data")

	units := createDummyContentUnits(suite.tx, 10)

	resp, err = handleContentUnitsList(cp, suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.ContentUnits {
		suite.assertEqualDummyContentUnit(units[i], x, i)
	}

	req.StartIndex = 6
	req.StopIndex = 10
	resp, err = handleContentUnitsList(cp, suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.ContentUnits {
		suite.assertEqualDummyContentUnit(units[i+5], x, i+5)
	}

	// test IDsFilter, UIDsFilter
	ids := make([]int64, len(units))
	uids := make([]string, len(units))
	for i := range units {
		cu := units[i]
		ids[i] = cu.ID
		uids[i] = cu.UID
	}

	req = ContentUnitsRequest{
		IDsFilter: IDsFilter{
			IDs: ids,
		},
	}
	resp, err = handleContentUnitsList(cp, suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.ContentUnits {
		suite.assertEqualDummyContentUnit(units[i], x, i)
	}

	req = ContentUnitsRequest{
		UIDsFilter: UIDsFilter{
			UIDs: uids,
		},
	}
	resp, err = handleContentUnitsList(cp, suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.ContentUnits {
		suite.assertEqualDummyContentUnit(units[i], x, i)
	}
}

func (suite *RestSuite) TestContentUnitItem() {
	cp := new(DummyAuthProvider)
	cu, err := handleGetContentUnit(cp, suite.tx, 1)
	suite.Nil(cu, "ContentUnit nil")
	suite.Require().NotNil(err, "Not Found error")
	suite.Equal(http.StatusNotFound, err.Code, "Error http status code")
	suite.Equal(gin.ErrorTypePublic, err.Type, "Error gin type")

	units := createDummyContentUnits(suite.tx, 3)
	for i, cu := range units {
		x, err := handleGetContentUnit(cp, suite.tx, cu.ID)
		suite.Require().Nil(err, "ContentUnit item err [%d]", i)
		suite.assertEqualDummyContentUnit(cu, x, i)
	}
}

func (suite *RestSuite) TestFilesList() {
	cp := new(DummyAuthProvider)
	req := FilesRequest{
		ListRequest: ListRequest{StartIndex: 1, StopIndex: 5},
	}
	resp, err := handleFilesList(cp, suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(0, resp.Total, "empty total")
	suite.Empty(resp.Files, "empty data")

	files := createDummyFiles(suite.tx, 10)

	resp, err = handleFilesList(cp, suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.Files {
		suite.assertEqualDummyFile(files[i], x, i)
	}

	req.StartIndex = 6
	req.StopIndex = 10
	resp, err = handleFilesList(cp, suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.Files {
		suite.assertEqualDummyFile(files[i+5], x, i+5)
	}

	// test IDsFilter, UIDsFilter, SHA1sFilter
	ids := make([]int64, len(files))
	uids := make([]string, len(files))
	sha1s := make([]string, len(files))
	for i := range files {
		f := files[i]
		ids[i] = f.ID
		uids[i] = f.UID
		sha1s[i] = hex.EncodeToString(f.Sha1.Bytes)
	}

	req = FilesRequest{
		IDsFilter: IDsFilter{
			IDs: ids,
		},
	}
	resp, err = handleFilesList(cp, suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.Files {
		suite.assertEqualDummyFile(files[i], x, i)
	}

	req = FilesRequest{
		UIDsFilter: UIDsFilter{
			UIDs: uids,
		},
	}
	resp, err = handleFilesList(cp, suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.Files {
		suite.assertEqualDummyFile(files[i], x, i)
	}

	req = FilesRequest{
		SHA1sFilter: SHA1sFilter{
			SHA1s: sha1s,
		},
	}
	resp, err = handleFilesList(cp, suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.Files {
		suite.assertEqualDummyFile(files[i], x, i)
	}
}

func (suite *RestSuite) TestFileItem() {
	cp := new(DummyAuthProvider)
	f, err := handleGetFile(cp, suite.tx, 1)
	suite.Nil(f, "file nil")
	suite.Require().NotNil(err, "Not Found error")
	suite.Equal(http.StatusNotFound, err.Code, "Error http status code")
	suite.Equal(gin.ErrorTypePublic, err.Type, "Error gin type")

	files := createDummyFiles(suite.tx, 3)
	for i, f := range files {
		x, err := handleGetFile(cp, suite.tx, f.ID)
		suite.Require().Nil(err, "file item err [%d]", i)
		suite.assertEqualDummyFile(f, x, i)
	}
}

func (suite *RestSuite) TestOperationsList() {
	req := OperationsRequest{
		ListRequest: ListRequest{StartIndex: 1, StopIndex: 5},
	}
	resp, err := handleOperationsList(suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(0, resp.Total, "empty total")
	suite.Empty(resp.Operations, "empty data")

	operations := createDummyOperations(suite.tx, 10)

	resp, err = handleOperationsList(suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.Operations {
		suite.assertEqualDummyOperation(operations[i], x, i)
	}

	req.StartIndex = 6
	req.StopIndex = 10
	resp, err = handleOperationsList(suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.Operations {
		suite.assertEqualDummyOperation(operations[i+5], x, i+5)
	}
}

func (suite *RestSuite) TestOperationItem() {
	cp := new(DummyAuthProvider)
	f, err := handleGetFile(cp, suite.tx, 1)
	suite.Nil(f, "file nil")
	suite.Require().NotNil(err, "Not Found error")
	suite.Equal(http.StatusNotFound, err.Code, "Error http status code")
	suite.Equal(gin.ErrorTypePublic, err.Type, "Error gin type")

	files := createDummyFiles(suite.tx, 3)
	for i, f := range files {
		x, err := handleGetFile(cp, suite.tx, f.ID)
		suite.Require().Nil(err, "file item err [%d]", i)
		suite.assertEqualDummyFile(f, x, i)
	}
}

// custom assertions

func (suite *RestSuite) assertEqualDummyCollection(c *models.Collection, x *Collection, idx int) {
	suite.Equal(c.ID, x.ID, "collection.ID [%d]", idx)
	suite.Equal(c.UID, x.UID, "collection.UID [%d]", idx)
	suite.Equal(c.TypeID, x.TypeID, "collection.TypeID [%d]", idx)
	suite.Equal(len(c.R.CollectionI18ns), len(x.I18n), "collection i18ns length [%d]", idx)
	for _, i18n := range c.R.CollectionI18ns {
		xi18n := x.I18n[i18n.Language]
		suite.Equal(i18n.CollectionID, xi18n.CollectionID,
			"collection %s i18n.CollectionID [%d]", i18n.Language, idx)
		suite.Equal(i18n.Name, xi18n.Name,
			"collection %s i18n.Name [%d]", i18n.Language, idx)
	}
}

func (suite *RestSuite) assertEqualDummyContentUnit(cu *models.ContentUnit, x *ContentUnit, idx int) {
	suite.Equal(cu.ID, x.ID, "content_unit.ID [%d]", idx)
	suite.Equal(cu.UID, x.UID, "content_unit.UID [%d]", idx)
	suite.Equal(cu.TypeID, x.TypeID, "content_unit.TypeID [%d]", idx)
	suite.Equal(len(cu.R.ContentUnitI18ns), len(x.I18n), "content_unit i18ns length [%d]", idx)
	for _, i18n := range cu.R.ContentUnitI18ns {
		xi18n := x.I18n[i18n.Language]
		suite.Equal(i18n.ContentUnitID, xi18n.ContentUnitID,
			"content_unit %s i18n.ContentUnitID [%d]", i18n.Language, idx)
		suite.Equal(i18n.Name, xi18n.Name,
			"content_unit %s i18n.Name [%d]", i18n.Language, idx)
	}
}

func (suite *RestSuite) assertEqualDummyFile(f *models.File, x *MFile, idx int) {
	suite.Equal(f.ID, x.ID, "file.ID [%d]", idx)
	suite.Equal(f.UID, x.UID, "file.UID [%d]", idx)
	suite.Equal(f.Size, x.Size, "file.Size [%d]", idx)
	suite.Equal(hex.EncodeToString(f.Sha1.Bytes), x.Sha1Str, "file.Sha1Str [%d]", idx)
}

func (suite *RestSuite) assertEqualDummyOperation(o *models.Operation, x *models.Operation, idx int) {
	suite.Equal(o.ID, x.ID, "operation.ID [%d]", idx)
	suite.Equal(o.UID, x.UID, "operation.UID [%d]", idx)
	suite.Equal(o.Station, x.Station, "operation.Station [%d]", idx)
	suite.Equal(o.UserID, x.UserID, "operation.UserID [%d]", idx)
	suite.Equal(o.TypeID, x.TypeID, "operation.TypeID [%d]", idx)
}

// Helpers

func createDummyCollections(exec boil.Executor, n int) []*models.Collection {
	collections := make([]*models.Collection, n)
	for i := range collections {
		j := n - i - 1
		collections[j] = &models.Collection{
			UID:    common.GenerateUID(8),
			TypeID: common.CONTENT_TYPE_REGISTRY.ByName[common.ALL_CONTENT_TYPES[rand.Intn(len(common.ALL_CONTENT_TYPES))]].ID,
		}
		utils.Must(collections[j].Insert(exec, boil.Infer()))

		i18ns := []*models.CollectionI18n{
			{Language: common.LANG_HEBREW, Name: null.StringFrom("name")},
			{Language: common.LANG_ENGLISH, Name: null.StringFrom("name")},
			{Language: common.LANG_RUSSIAN, Name: null.StringFrom("name")},
		}
		collections[j].AddCollectionI18ns(exec, true, i18ns...)
	}

	return collections
}

func createDummyContentUnits(exec boil.Executor, n int) []*models.ContentUnit {
	units := make([]*models.ContentUnit, n)
	for i := range units {
		j := n - i - 1
		units[j] = &models.ContentUnit{
			UID:    common.GenerateUID(8),
			TypeID: common.CONTENT_TYPE_REGISTRY.ByName[common.ALL_CONTENT_TYPES[rand.Intn(len(common.ALL_CONTENT_TYPES))]].ID,
		}
		utils.Must(units[j].Insert(exec, boil.Infer()))

		i18ns := []*models.ContentUnitI18n{
			{Language: common.LANG_HEBREW, Name: null.StringFrom("name")},
			{Language: common.LANG_ENGLISH, Name: null.StringFrom("name")},
			{Language: common.LANG_RUSSIAN, Name: null.StringFrom("name")},
		}
		units[j].AddContentUnitI18ns(exec, true, i18ns...)
	}

	return units
}

func createDummyFiles(exec boil.Executor, n int) []*models.File {
	files := make([]*models.File, n)
	for i := range files {
		j := n - i - 1
		sha1 := make([]byte, 20)
		rand.Read(sha1)
		files[j] = &models.File{
			UID:  common.GenerateUID(8),
			Name: fmt.Sprintf("test_file_%d", j),
			Size: rand.Int63(),
			Sha1: null.BytesFrom(sha1),
		}
		utils.Must(files[j].Insert(exec, boil.Infer()))
	}

	return files
}

func createDummyOperations(exec boil.Executor, n int) []*models.Operation {
	operations := make([]*models.Operation, n)
	for i := range operations {
		j := n - i - 1
		//sha1 := make([]byte, 20)
		//rand.Read(sha1)
		operations[j] = &models.Operation{
			UID:     common.GenerateUID(8),
			Station: null.StringFrom(fmt.Sprintf("station_%d", j)),
			UserID:  null.Int64From(1),
			TypeID: common.OPERATION_TYPE_REGISTRY.
				ByName[common.ALL_OPERATION_TYPES[rand.Intn(len(common.ALL_OPERATION_TYPES))]].ID,
		}
		utils.Must(operations[j].Insert(exec, boil.Infer()))
	}

	return operations
}

type DummyAuthProvider struct {
}

func (p *DummyAuthProvider) Get(key string) (interface{}, bool) {

	switch key {
	case "ID_TOKEN_CLAIMS":
		return permissions.IDTokenClaims{
			Name: "Test User",
			Sub:  "test-user",
			RealmAccess: permissions.Roles{
				Roles: []string{"test_user"},
			},
		}, true
	default:
		return nil, false
	}
}

func (p *DummyAuthProvider) MustGet(key string) interface{} {
	switch key {
	case "PERMISSIONS_ENFORCER":
		enforcer := casbin.NewEnforcer()
		enforcer.EnableEnforce(false)
		return enforcer
	default:
		return nil
	}
}
