package api

import (
	"math/rand"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vattle/sqlboiler/boil"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

type RestSuite struct {
	suite.Suite
	utils.TestDBManager
	tx boil.Transactor
}

func (suite *RestSuite) SetupSuite() {
	suite.Require().Nil(suite.InitTestDB())
	suite.Require().Nil(OPERATION_TYPE_REGISTRY.Init())
	suite.Require().Nil(CONTENT_TYPE_REGISTRY.Init())
}

func (suite *RestSuite) TearDownSuite() {
	suite.Require().Nil(suite.DestroyTestDB())
}

func (suite *RestSuite) SetupTest() {
	var err error
	suite.tx, err = boil.Begin()
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
	req := CollectionsRequest{
		ListRequest: ListRequest{StartIndex: 1, StopIndex: 5},
	}
	resp, err := handleCollectionsList(suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(0, resp.Total, "empty total")
	suite.Empty(resp.Collections, "empty data")

	collections := createDummyCollections(suite.tx, 10)

	resp, err = handleCollectionsList(suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.Collections {
		c := collections[i]
		suite.Equal(c.ID, x.ID, "collection.ID [%d]", i)
		suite.Equal(c.UID, x.UID, "collection.UID [%d]", i)
		suite.Equal(c.TypeID, x.TypeID, "collection.TypeID [%d]", i)
	}

	req.StartIndex = 6
	req.StopIndex = 10
	resp, err = handleCollectionsList(suite.tx, req)
	suite.Require().Nil(err)
	suite.EqualValues(10, resp.Total, "total")
	for i, x := range resp.Collections {
		c := collections[i+5]
		suite.Equal(c.ID, x.ID, "collection.ID [%d]", i+5)
		suite.Equal(c.UID, x.UID, "collection.UID [%d]", i+5)
		suite.Equal(c.TypeID, x.TypeID, "collection.TypeID [%d]", i+5)
	}
}

func (suite *RestSuite) TestCollectionsItem() {
	c, err := handleCollectionItem(suite.tx, 1)
	suite.Nil(c, "Collection nil")
	suite.Require().NotNil(err, "Not Found error")
	suite.Equal(http.StatusNotFound, err.Code, "Error http status code")
	suite.Equal(gin.ErrorTypePublic, err.Type, "Error gin type")

	collections := createDummyCollections(suite.tx, 3)
	for i, x := range collections {
		c, err := handleCollectionItem(suite.tx, x.ID)
		suite.Require().Nil(err, "Collection Item err [%d]", i)
		suite.Equal(c.ID, x.ID, "collection.ID [%d]", i)
		suite.Equal(c.UID, x.UID, "collection.UID [%d]", i)
		suite.Equal(c.TypeID, x.TypeID, "collection.TypeID [%d]", i)
	}
}

// Helpers

func createDummyCollections(exec boil.Executor, n int) []*models.Collection {
	collections := make([]*models.Collection, n)
	for i := range collections {
		collections[i] = &models.Collection{
			UID:    utils.GenerateUID(8),
			TypeID: CONTENT_TYPE_REGISTRY.ByName[ALL_CONTENT_TYPES[rand.Intn(len(ALL_CONTENT_TYPES))]].ID,
		}
		utils.Must(collections[i].Insert(exec))
	}

	return collections
}

