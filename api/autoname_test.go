package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/vattle/sqlboiler/boil"
	"gopkg.in/nullbio/null.v6"

	"github.com/Bnei-Baruch/mdb/utils"
)

type AutonameSuite struct {
	suite.Suite
	utils.TestDBManager
	tx boil.Transactor
}

func (suite *AutonameSuite) SetupSuite() {
	suite.Require().Nil(suite.InitTestDB())
	suite.Require().Nil(InitTypeRegistries(boil.GetDB()))
}

func (suite *AutonameSuite) TearDownSuite() {
	suite.Require().Nil(suite.DestroyTestDB())
}

func (suite *AutonameSuite) SetupTest() {
	var err error
	suite.tx, err = boil.Begin()
	suite.Require().Nil(err)
}

func (suite *AutonameSuite) TearDownTest() {
	err := suite.tx.Rollback()
	suite.Require().Nil(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAutoname(t *testing.T) {
	suite.Run(t, new(AutonameSuite))
}

func (suite *AutonameSuite) TestLessonPart() {
	cu, err := CreateContentUnit(suite.tx, CT_LESSON_PART, nil)
	suite.Require().Nil(err)

	metadata := CITMetadata{
		ContentType:    CT_LESSON_PART,
		AutoName:       "auto_name",
		FinalName:      "final_name",
		CaptureDate:    Date{time.Now()},
		Language:       LANG_HEBREW,
		HasTranslation: true,
		Lecturer:       "rav",
		Number:         null.IntFrom(1),
		Part:           null.IntFrom(0),
	}

	err = AutonameContentUnit(suite.tx, cu, metadata)
	suite.Require().Nil(err)
	err = cu.L.LoadContentUnitI18ns(suite.tx, true, cu)
	suite.Require().Nil(err)
	i18ns := cu.R.ContentUnitI18ns
	suite.NotEmpty(i18ns, "i18ns.empty")
	for _, i18n := range i18ns {
		switch i18n.Language {
		case LANG_HEBREW:
			suite.Equal("הכנה לשיעור", i18n.Name.String, "Hebrew name")
			break
		case LANG_ENGLISH:
			suite.Equal("lesson preparation", i18n.Name.String, "English name")
			break
		}
	}
}

// Helpers
