package cusource

import (
	"database/sql"
	"fmt"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"testing"
)

type CUSourceSuite struct {
	suite.Suite
	utils.TestDBManager
	tx *sql.Tx
}

func (suite *CUSourceSuite) rootPath() string {
	return "../../"
}

func (suite *CUSourceSuite) SetupSuite() {
	suite.Require().Nil(suite.InitTestDB("../../"))
	suite.Require().Nil(common.InitTypeRegistries(suite.DB))
}

func (suite *CUSourceSuite) TearDownSuite() {
	suite.Require().Nil(suite.DestroyTestDB())
}

func (suite *CUSourceSuite) SetupTest() {
	var err error
	suite.tx, err = suite.DB.Begin()
	suite.Require().Nil(err)
}

func (suite *CUSourceSuite) TearDownTest() {
	err := suite.tx.Rollback()
	suite.Require().Nil(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCUSource(t *testing.T) {
	suite.Run(t, new(CUSourceSuite))
}

func (suite *CUSourceSuite) TestCreateSourceCU() {
	sources, cus := BuildCUSources(suite.DB)
	suite.Require().Equal(len(sources), len(cus))
	fmt.Println("after insert units:", len(cus), len(sources))

	suite.testResults(sources)
}
func (suite *CUSourceSuite) testResults(sources []*models.Source) {
	for _, s := range sources {
		has, err := models.ContentUnits(suite.DB,
			qm.InnerJoin("content_units_sources cus ON cus.source_id = ?", s.ID),
		).Exists()
		suite.Require().Nil(err)
		suite.Require().True(has)
	}
}
