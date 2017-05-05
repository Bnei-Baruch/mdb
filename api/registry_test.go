package api

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vattle/sqlboiler/boil"

	"github.com/Bnei-Baruch/mdb/utils"
)

type RegistrySuite struct {
	suite.Suite
	utils.TestDBManager
	tx boil.Transactor
}

func (suite *RegistrySuite) SetupSuite() {
	suite.Require().Nil(suite.InitTestDB())
}

func (suite *RegistrySuite) TearDownSuite() {
	suite.Require().Nil(suite.DestroyTestDB())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRegistry(t *testing.T) {
	suite.Run(t, new(RegistrySuite))
}

func (suite *RegistrySuite) TestTypeRegistries() {
	suite.Require().Nil(InitTypeRegistries(boil.GetDB()))

	for _, x := range ALL_CONTENT_TYPES {
		ct, ok := CONTENT_TYPE_REGISTRY.ByName[x]
		suite.True(ok, "CT exists %s", x)
		suite.NotNil(ct, "CT not nil %s", x)
	}

	for _, x := range ALL_OPERATION_TYPES {
		ot, ok := OPERATION_TYPE_REGISTRY.ByName[x]
		suite.True(ok, "OT exists %s", x)
		suite.NotNil(ot, "OT not nil %s", x)
	}

	for _, x := range []string{CR_LECTURER} {
		cr, ok := CONTENT_ROLE_TYPE_REGISTRY.ByName[x]
		suite.True(ok, "CR exists %s", x)
		suite.NotNil(cr, "CR not nil %s", x)
	}

	for _, x := range []string{P_RAV} {
		p, ok := PERSONS_REGISTRY.ByPattern[x]
		suite.True(ok, "P exists %s", x)
		suite.NotNil(p, "P not nil %s", x)
	}
}
