package api

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/nullbio/null.v6"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

type MetadataProcessorSuite struct {
	suite.Suite
	utils.TestDBManager
	tx boil.Transactor
}

func (suite *MetadataProcessorSuite) SetupSuite() {
	suite.Require().Nil(suite.InitTestDB())
	suite.Require().Nil(OPERATION_TYPE_REGISTRY.Init())
	suite.Require().Nil(CONTENT_TYPE_REGISTRY.Init())
}

func (suite *MetadataProcessorSuite) TearDownSuite() {
	suite.Require().Nil(suite.DestroyTestDB())
}

func (suite *MetadataProcessorSuite) SetupTest() {
	var err error
	suite.tx, err = boil.Begin()
	suite.Require().Nil(err)
}

func (suite *MetadataProcessorSuite) TearDownTest() {
	err := suite.tx.Rollback()
	suite.Require().Nil(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMetadataProcessor(t *testing.T) {
	suite.Run(t, new(MetadataProcessorSuite))
}

func (suite *MetadataProcessorSuite) TestProcessCITMetadata() {
	// Create dummy original and proxy trimmed files
	ofi := File{
		FileName:  "dummy original trimmed file",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef9876543210",
		Size:      math.MaxInt64,
	}
	original, err := CreateFile(suite.tx, nil, ofi, nil)
	suite.Require().Nil(err)

	pfi := File{
		FileName:  "dummy proxy trimmed file",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "987653210abcdef012356789abcdef9876543210",
		Size:      math.MaxInt64,
	}
	proxy, err := CreateFile(suite.tx, nil, pfi, nil)
	suite.Require().Nil(err)

	metadata := CITMetadata{
		ContentType:    CT_LESSON_PART,
		FinalName:      "auto_name",
		CaptureDate:    Date{time.Now()},
		Language:       LANG_HEBREW,
		Lecturer:       "rav",
		WeekDate:       nil,
		Number:         null.IntFrom(2),
		Part:           null.IntFrom(1),
		Sources:        []string{},
		Tags:           []string{},
		ArtifactType:   null.NewString("", false),
		HasTranslation: true,
		RequireTest:    false,
		CollectionUID:  null.NewString("", false),
		Episode:        null.NewString("", false),
	}

	// add some sources and tags
	sources, err := models.Sources(suite.tx, qm.Limit(3)).All()
	suite.Require().Nil(err)
	sourcesUIDs := make([]string, len(sources))
	for i, x := range sources {
		sourcesUIDs[i] = x.UID
	}
	metadata.Sources = sourcesUIDs

	tags, err := models.Tags(suite.tx, qm.Limit(3)).All()
	suite.Require().Nil(err)
	tagsUIDs := make([]string, len(tags))
	for i, x := range tags {
		tagsUIDs[i] = x.UID
	}
	metadata.Tags = tagsUIDs

	err = ProcessCITMetadata(suite.tx, metadata, original, proxy)
	suite.Require().Nil(err)

	// check files
	err = original.Reload(suite.tx)
	suite.Require().Nil(err)
	err = proxy.Reload(suite.tx)
	suite.Require().Nil(err)

	var props map[string]interface{}

	// original properties
	suite.Require().True(original.Properties.Valid)
	err = original.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(metadata.CaptureDate.Format(time.RFC3339Nano), props["capture_date"], "original.Properties[\"capture_date\"]")
	suite.Equal(metadata.CaptureDate.Format(time.RFC3339Nano), props["film_date"], "original.Properties[\"film_date\"]")

	// proxy properties
	suite.Require().True(proxy.Properties.Valid)
	err = proxy.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(metadata.CaptureDate.Format(time.RFC3339Nano), props["capture_date"], "proxy.Properties[\"capture_date\"]")
	suite.Equal(metadata.CaptureDate.Format(time.RFC3339Nano), props["film_date"], "proxy.Properties[\"film_date\"]")

	// original language
	suite.True(original.Language.Valid, "original.Language.Valid")
	suite.Equal(LANG_MULTI, original.Language.String, "original.Language")

	// new content unit
	suite.True(original.ContentUnitID.Valid, "original.ContentUnitID.Valid")
	suite.True(proxy.ContentUnitID.Valid, "proxy.ContentUnitID.Valid")
	suite.Equal(original.ContentUnitID.Int64, proxy.ContentUnitID.Int64, "original.cuid = proxy.cuid")

	original.L.LoadContentUnit(suite.tx, true, original)
	cu := original.R.ContentUnit
	suite.Equal(8, len(cu.UID), "cu.UID")
	// unit properties
	suite.Require().True(cu.Properties.Valid)
	err = cu.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(metadata.CaptureDate.Format(time.RFC3339Nano), props["capture_date"], "cu.Properties[\"capture_date\"]")
	suite.Equal(metadata.CaptureDate.Format(time.RFC3339Nano), props["film_date"], "cu.Properties[\"film_date\"]")
	// unit sources
	err = cu.L.LoadSources(suite.tx, true, cu)
	suite.Require().Nil(err)
	suite.Equal(len(sources), len(cu.R.Sources), "len(cu.R.Sources)")
	for _, x := range sources {
		missing := true
		for _, y := range cu.R.Sources {
			if x.UID == y.UID {
				missing = false
				break
			}
		}
		suite.False(missing, "Missing source %s", x.ID)
	}
	// unit tags
	err = cu.L.LoadTags(suite.tx, true, cu)
	suite.Require().Nil(err)
	suite.Equal(len(tags), len(cu.R.Tags), "len(cu.R.Tags)")
	for _, x := range tags {
		missing := true
		for _, y := range cu.R.Tags {
			if x.ID == y.ID {
				missing = false
				break
			}
		}
		suite.False(missing, "Missing tag %s", x.ID)
	}

	// collection
	err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu)
	suite.Require().Nil(err)
	suite.Equal(1, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")
}
