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

func (suite *AutonameSuite) TestSourceNamers() {
	author := &models.Author{ID: -1}
	author.L.LoadAuthorI18ns(suite.tx, true, author) // dummy call to initialize R
	author.R.AuthorI18ns = make(models.AuthorI18nSlice, 0)
	author.R.AuthorI18ns = append(author.R.AuthorI18ns,
		&models.AuthorI18n{
			Language: LANG_HEBREW,
			Name:     null.StringFrom("author"),
		},
		&models.AuthorI18n{
			Language: LANG_ENGLISH,
			Name:     null.StringFrom("author"),
		},
	)

	path := make([]*models.Source, 4)
	for i := 1; i < 5; i++ {
		s := &models.Source{ID: -1}
		s.L.LoadSourceI18ns(suite.tx, true, s) // dummy call to initialize R
		s.R.SourceI18ns = make(models.SourceI18nSlice, 0)
		s.R.SourceI18ns = append(s.R.SourceI18ns,
			&models.SourceI18n{
				Language: LANG_HEBREW,
				Name:     null.StringFrom(fmt.Sprintf("source %d", i)),
			})
		path[i-1] = s
	}

	var namer sourceNamer
	namer = new(PlainNamer)
	names, err := namer.GetName(author, path)
	suite.Require().Nil(err)
	suite.Len(names, 1, "len(names)")
	name := names[LANG_HEBREW]
	suite.Equal("author. source 1. source 2. source 3. source 4", name, "name")

	namer = new(PrefaceNamer)
	names, err = namer.GetName(author, path)
	suite.Require().Nil(err)
	suite.Len(names, 1, "len(names)")
	name = names[LANG_HEBREW]
	suite.Equal("author. source 4", name, "name")

	namer = new(LettersNamer)
	path[len(path) -1].R.SourceI18ns[0].Name = null.StringFrom("source 4 (1920)")
	names, err = namer.GetName(author, path)
	suite.Require().Nil(err)
	suite.Len(names, 1, "len(names)")
	name = names[LANG_HEBREW]
	suite.Equal("author. source 4", name, "name")
	path[len(path) -1].R.SourceI18ns[0].Name = null.StringFrom("source 4")

	namer = new(RBRecordsNamer)
	path[len(path) -1].Position = null.IntFrom(137)
	names, err = namer.GetName(author, path)
	suite.Require().Nil(err)
	suite.Len(names, 1, "len(names)")
	name = names[LANG_HEBREW]
	suite.Equal("author. רשומה 137. source 4", name, "name")

	namer = new(RBArticlesNamer)
	path[len(path) -1].Name = "(1984-01-2) Matarat Hevra 2"
	names, err = namer.GetName(author, path)
	suite.Require().Nil(err)
	suite.Len(names, 1, "len(names)")
	name = names[LANG_HEBREW]
	suite.Equal("author. source 4. 1-2 (1984)", name, "name")

	namer = new(ShamatiNamer)
	path[len(path) -1].Name = "015 some name"
	names, err = namer.GetName(author, path)
	suite.Require().Nil(err)
	suite.Len(names, 1, "len(names)")
	name = names[LANG_HEBREW]
	suite.Equal("author. source 1, טו. source 4", name, "name")

	path[1].R.SourceI18ns[0].Description = null.StringFrom("source 2 description")
	namer = new(ZoharNamer)
	names, err = namer.GetName(author, path)
	suite.Require().Nil(err)
	suite.Len(names, 1, "len(names)")
	name = names[LANG_HEBREW]
	suite.Equal("source 1. source 2 description. source 3. source 4", name, "name")
}
