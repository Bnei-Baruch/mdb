package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/vattle/sqlboiler/boil"
	"gopkg.in/nullbio/null.v6"

	"fmt"
	"github.com/Bnei-Baruch/mdb/models"
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

func (suite *AutonameSuite) TestGenericDescriberCollection() {
	describer := new(GenericDescriber)
	c := new(models.Collection)

	for _, x := range ALL_CONTENT_TYPES {
		c.TypeID = CONTENT_TYPE_REGISTRY.ByName[x].ID
		names, err := GetI18ns(fmt.Sprintf("content_type.%s", x))
		suite.Require().Nil(err)
		i18ns, err := describer.DescribeCollection(c)
		suite.Require().Nil(err)
		suite.Len(i18ns, len(names), "len(i18ns)")
		for _, i18n := range i18ns {
			suite.Equal(names[i18n.Language], i18n.Name.String, "%s name", i18n.Language)
		}
	}
}

func (suite *AutonameSuite) TestGenericDescriberContentUnit() {
	describer := new(GenericDescriber)
	cu := new(models.ContentUnit)

	for _, x := range ALL_CONTENT_TYPES {
		metadata := CITMetadata{
			FinalName: "final_name",
		}
		cu.TypeID = CONTENT_TYPE_REGISTRY.ByName[x].ID
		i18ns, err := describer.DescribeContentUnit(suite.tx, cu, metadata)
		suite.Require().Nil(err)
		suite.Len(i18ns, 3, "len(i18ns)")
		for _, i18n := range i18ns {
			suite.Equal(metadata.FinalName, i18n.Name.String, "%s name", i18n.Language)
		}
	}
}

func (suite *AutonameSuite) TestLessonPartDescriber() {
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

	describer := new(LessonPartDescriber)
	i18ns, err := describer.DescribeContentUnit(suite.tx, cu, metadata)
	suite.Require().Nil(err)
	suite.NotEmpty(i18ns, "i18ns.empty")
	for _, i18n := range i18ns {
		switch i18n.Language {
		case LANG_HEBREW:
			suite.Equal("הכנה לשיעור", i18n.Name.String, "Hebrew name")
			break
		case LANG_ENGLISH:
			suite.Equal("Preparation to the Lesson", i18n.Name.String, "English name")
			break
		case LANG_RUSSIAN:
			suite.Equal("Подготовка к Уроку", i18n.Name.String, "Russian name")
			break
		}
	}
}

func (suite *AutonameSuite) TestDescribeContentUnit() {
	cu, err := CreateContentUnit(suite.tx, CT_LESSON_PART, nil)
	suite.Require().Nil(err)

	metadata := CITMetadata{
		ContentType: CT_UNKNOWN,
		FinalName:   "final_name",
	}

	err = DescribeContentUnit(suite.tx, cu, metadata)
	suite.Require().Nil(err)
	err = cu.L.LoadContentUnitI18ns(suite.tx, true, cu)
	suite.Require().Nil(err)
	i18ns := cu.R.ContentUnitI18ns
	suite.Len(i18ns, 3, "len(i18ns)")
	for _, i18n := range i18ns {
		switch i18n.Language {
		case LANG_HEBREW:
		case LANG_ENGLISH:
		case LANG_RUSSIAN:
			suite.Equal(metadata.FinalName, i18n.Name.String, "%s name", i18n.Language)
			break
		default:
			suite.Fail("Unexpected Language %s", i18n.Language)
		}
	}
}

func (suite *AutonameSuite) TestDescribeCollection() {
	c, err := CreateCollection(suite.tx, CT_UNKNOWN, nil)
	suite.Require().Nil(err)

	err = DescribeCollection(suite.tx, c)
	suite.Require().Nil(err)
	err = c.L.LoadCollectionI18ns(suite.tx, true, c)
	suite.Require().Nil(err)
	i18ns := c.R.CollectionI18ns

	names, err := GetI18ns("content_type.UNKNOWN")
	suite.Require().Nil(err)
	suite.Len(i18ns, len(names), "len(i18ns)")
	for _, i18n := range i18ns {
		suite.Equal(names[i18n.Language], i18n.Name.String, "%s name", i18n.Language)
	}
}
