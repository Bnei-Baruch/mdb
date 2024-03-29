package api

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Bnei-Baruch/mdb/common"
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
	suite.Require().Nil(common.InitTypeRegistries(suite.DB))
}

func (suite *MetadataProcessorSuite) TearDownSuite() {
	suite.Require().Nil(suite.DestroyTestDB())
}

func (suite *MetadataProcessorSuite) SetupTest() {
	var err error
	suite.tx, err = suite.DB.Begin()
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

func (suite *MetadataProcessorSuite) TestDailyLesson() {
	likutim, err := SomeLikutim(suite.tx)
	suite.Require().Nil(err)
	likutimUIDs := make([]string, len(likutim))
	for i, l := range likutim {
		likutimUIDs[i] = l.UID
	}

	chain := suite.simulateLessonChain()
	// send parts
	// send full
	// send kitei makor of part 1
	// send ktaim nivcharim from full
	metadata := CITMetadata{
		ContentType:    common.CT_LESSON_PART,
		AutoName:       "auto_name",
		FinalName:      "final_name",
		CaptureDate:    Date{time.Now()},
		Language:       common.LANG_HEBREW,
		HasTranslation: true,
		Lecturer:       "rav",
		Number:         null.IntFrom(1),
		Part:           null.IntFrom(0),
		Sources:        suite.someSources(),
		Tags:           suite.someTags(),
		Likutim:        likutimUIDs,
		RequireTest:    false,
	}
	original, proxy := chain["part0"].Original, chain["part0"].Proxy

	evnts, err := ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(evnts)

	err = original.Reload(suite.tx)
	suite.Require().Nil(err)
	err = proxy.Reload(suite.tx)
	suite.Require().Nil(err)

	suite.assertFiles(metadata, original, proxy)
	suite.assertContentUnit(metadata, original, proxy, false)

	// collection association
	err = original.L.LoadContentUnit(suite.tx, true, original, nil)
	suite.Require().Nil(err)
	cu := original.R.ContentUnit
	err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Equal(1, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")
	ccu := cu.R.CollectionsContentUnits[0]
	suite.Equal("0", ccu.Name, "ccu.Name")
	suite.Equal(0, ccu.Position, "ccu.Position")

	// collection
	err = ccu.L.LoadCollection(suite.tx, true, ccu, nil)
	suite.Require().Nil(err)
	c := ccu.R.Collection
	suite.Equal(common.CONTENT_TYPE_REGISTRY.ByName[common.CT_DAILY_LESSON].ID, c.TypeID, "c.TypeID")
	suite.True(c.Properties.Valid, "c.Properties.Valid")
	var props map[string]interface{}
	err = json.Unmarshal(c.Properties.JSON, &props)
	suite.Require().Nil(err)
	suite.Equal(metadata.CaptureDate.Format("2006-01-02"), props["capture_date"], "c.Properties[\"capture_date\"]")
	suite.Equal(metadata.CaptureDate.Format("2006-01-02"), props["film_date"], "c.Properties[\"film_date\"]")
	suite.Equal("c12356789", props["capture_id"], "c.Properties[\"capture_id\"]")
	suite.EqualValues(metadata.Number.Int, props["number"], "c.Properties[\"number\"]")

	// process other parts
	for i := 1; i < 4; i++ {
		metadata.Part = null.IntFrom(i)
		metadata.Sources = suite.someSources()
		metadata.Tags = suite.someTags()
		tf := chain[fmt.Sprintf("part%d", i)]
		original, proxy := tf.Original, tf.Proxy

		evnts, err := ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
		suite.Require().Nil(err)
		suite.Require().NotNil(evnts)

		err = original.Reload(suite.tx)
		suite.Require().Nil(err)
		err = proxy.Reload(suite.tx)
		suite.Require().Nil(err)

		suite.assertFiles(metadata, original, proxy)
		suite.assertContentUnit(metadata, original, proxy, false)

		// collection association
		err = original.L.LoadContentUnit(suite.tx, true, original, nil)
		suite.Require().Nil(err)
		cu := original.R.ContentUnit
		err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
		suite.Require().Nil(err)
		suite.Equal(1, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")
		ccu := cu.R.CollectionsContentUnits[0]
		suite.Equal(strconv.Itoa(i), ccu.Name, "ccu.Name")
		suite.Equal(i, ccu.Position, "ccu.Position")
		suite.Equal(c.ID, ccu.CollectionID, "ccu.CollectionID")
	}

	// process full
	metadata.ContentType = common.CT_FULL_LESSON
	metadata.Part = null.IntFrom(-1)
	metadata.Sources = nil
	metadata.Tags = nil
	tf := chain["full"]
	original, proxy = tf.Original, tf.Proxy

	evnts, err = ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(evnts)

	err = original.Reload(suite.tx)
	suite.Require().Nil(err)
	err = proxy.Reload(suite.tx)
	suite.Require().Nil(err)

	suite.assertFiles(metadata, original, proxy)
	suite.assertContentUnit(metadata, original, proxy, false)

	// collection association
	err = original.L.LoadContentUnit(suite.tx, true, original, nil)
	suite.Require().Nil(err)
	cu = original.R.ContentUnit
	err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Equal(1, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")
	ccu = cu.R.CollectionsContentUnits[0]
	suite.Equal("full", ccu.Name, "ccu.Name")
	suite.Equal(4, ccu.Position, "ccu.Position")
	suite.Equal(c.ID, ccu.CollectionID, "ccu.CollectionID")

	// full with week_date different from capture_date
	metadata.WeekDate = &Date{Time: time.Now().AddDate(1, 0, 0)}
	evnts, err = ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(evnts)
	err = c.Reload(suite.tx)
	suite.Require().Nil(err)
	suite.Equal(common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SPECIAL_LESSON].ID, c.TypeID, "c.TypeID")
	suite.True(c.Properties.Valid, "c.Properties.Valid")
	err = json.Unmarshal(c.Properties.JSON, &props)
	suite.Require().Nil(err)
	suite.Equal(metadata.CaptureDate.Format("2006-01-02"), props["capture_date"], "c.Properties[\"capture_date\"]")
	suite.Equal(metadata.CaptureDate.Format("2006-01-02"), props["film_date"], "c.Properties[\"film_date\"]")
	suite.Equal("c12356789", props["capture_id"], "c.Properties[\"capture_id\"]")
	suite.EqualValues(metadata.Number.Int, props["number"], "c.Properties[\"number\"]")

	// process kitei makor for part 1
	metadata.ContentType = common.CT_LESSON_PART
	metadata.Part = null.IntFrom(1)
	metadata.ArtifactType = null.StringFrom("kitei_makor")
	metadata.WeekDate = nil
	tf = chain["part1_kitei-makor"]
	original, proxy = tf.Original, tf.Proxy
	evnts, err = ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(evnts)

	err = original.Reload(suite.tx)
	suite.Require().Nil(err)
	err = proxy.Reload(suite.tx)
	suite.Require().Nil(err)

	suite.assertFiles(metadata, original, proxy)
	suite.assertContentUnit(metadata, original, proxy, false)

	// associated to "main" content unit
	err = original.L.LoadContentUnit(suite.tx, true, original, nil)
	suite.Require().Nil(err)
	cu = original.R.ContentUnit
	err = cu.L.LoadDerivedContentUnitDerivations(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Require().Len(cu.R.DerivedContentUnitDerivations, 1, "cu.R.DerivationContentUnitDerivations length")
	cud := cu.R.DerivedContentUnitDerivations[0]
	suite.Equal(chain["part1"].Original.ContentUnitID.Int64, cud.SourceID, "cud.SourceID")
	suite.Equal("kitei_makor", cud.Name, "cud.Name")
	err = cu.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	_, ok := props["artifact_type"]
	suite.False(ok, "cu.properties[\"artifact_type\"]")

	// not associated with collection
	err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Empty(cu.R.CollectionsContentUnits, "cu.R.CollectionsContentUnits empty")

	// no changes to collection
	err = c.Reload(suite.tx)
	suite.Require().Nil(err)
	suite.Equal(common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SPECIAL_LESSON].ID, c.TypeID, "c.TypeID")
	suite.True(c.Properties.Valid, "c.Properties.Valid")
	err = json.Unmarshal(c.Properties.JSON, &props)
	suite.Require().Nil(err)
	suite.Equal(metadata.CaptureDate.Format("2006-01-02"), props["capture_date"], "c.Properties[\"capture_date\"]")
	suite.Equal(metadata.CaptureDate.Format("2006-01-02"), props["film_date"], "c.Properties[\"film_date\"]")
	suite.Equal("c12356789", props["capture_id"], "c.Properties[\"capture_id\"]")
	suite.EqualValues(metadata.Number.Int, props["number"], "c.Properties[\"number\"]")

	// process ktaim nivcharim from full
	metadata.ContentType = common.CT_FULL_LESSON
	metadata.Part = null.IntFrom(-1)
	metadata.ArtifactType = null.StringFrom("KTAIM_NIVCHARIM")
	metadata.WeekDate = nil
	tf = chain["ktaim-nivcharim"]
	original, proxy = tf.Original, tf.Proxy
	evnts, err = ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(evnts)

	err = original.Reload(suite.tx)
	suite.Require().Nil(err)
	err = proxy.Reload(suite.tx)
	suite.Require().Nil(err)

	suite.assertFiles(metadata, original, proxy)
	suite.assertContentUnit(metadata, original, proxy, false)

	// associated to "main" content unit
	err = original.L.LoadContentUnit(suite.tx, true, original, nil)
	suite.Require().Nil(err)
	cu = original.R.ContentUnit
	err = cu.L.LoadDerivedContentUnitDerivations(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Require().Len(cu.R.DerivedContentUnitDerivations, 1, "cu.R.DerivationContentUnitDerivations length")
	cud = cu.R.DerivedContentUnitDerivations[0]
	suite.Equal(chain["full"].Original.ContentUnitID.Int64, cud.SourceID, "cud.SourceID")
	suite.Equal("KTAIM_NIVCHARIM", cud.Name, "cud.Name")
	err = cu.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	_, ok = props["artifact_type"]
	suite.False(ok, "cu.properties[\"artifact_type\"]")

	// associated with collection
	err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Equal(1, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")
	ccu = cu.R.CollectionsContentUnits[0]
	suite.Equal("KTAIM_NIVCHARIM_1", ccu.Name, "ccu.Name")
	suite.Equal(6, ccu.Position, "ccu.Position")
	suite.Equal(c.ID, ccu.CollectionID, "ccu.CollectionID")
}

func (suite *MetadataProcessorSuite) TestSpecialLesson() {
	chain := suite.simulateSpecialLessonChain()

	// send parts
	// send full
	// send kitei makor of all parts
	// send lelo mikud of all parts

	metadata := CITMetadata{
		ContentType:    common.CT_LESSON_PART,
		AutoName:       "auto_name",
		FinalName:      "final_name",
		CaptureDate:    Date{time.Now()},
		Language:       common.LANG_HEBREW,
		HasTranslation: true,
		Lecturer:       "rav",
		Number:         null.IntFrom(1),
		Part:           null.IntFrom(0),
		Sources:        suite.someSources(),
		Tags:           suite.someTags(),
		RequireTest:    false,
	}
	original, proxy := chain["part0"].Original, chain["part0"].Proxy

	evnts, err := ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(evnts)

	err = original.Reload(suite.tx)
	suite.Require().Nil(err)
	err = proxy.Reload(suite.tx)
	suite.Require().Nil(err)

	suite.assertFiles(metadata, original, proxy)
	suite.assertContentUnit(metadata, original, proxy, false)

	// collection association
	err = original.L.LoadContentUnit(suite.tx, true, original, nil)
	suite.Require().Nil(err)
	cu := original.R.ContentUnit
	err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Equal(1, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")
	ccu := cu.R.CollectionsContentUnits[0]
	suite.Equal("0", ccu.Name, "ccu.Name")
	suite.Equal(0, ccu.Position, "ccu.Position")

	// collection
	err = ccu.L.LoadCollection(suite.tx, true, ccu, nil)
	suite.Require().Nil(err)
	c := ccu.R.Collection
	suite.Equal(common.CONTENT_TYPE_REGISTRY.ByName[common.CT_DAILY_LESSON].ID, c.TypeID, "c.TypeID")
	suite.True(c.Properties.Valid, "c.Properties.Valid")
	var props map[string]interface{}
	err = json.Unmarshal(c.Properties.JSON, &props)
	suite.Require().Nil(err)
	suite.Equal(metadata.CaptureDate.Format("2006-01-02"), props["capture_date"], "c.Properties[\"capture_date\"]")
	suite.Equal(metadata.CaptureDate.Format("2006-01-02"), props["film_date"], "c.Properties[\"film_date\"]")
	suite.Equal("c12356789", props["capture_id"], "c.Properties[\"capture_id\"]")
	suite.EqualValues(metadata.Number.Int, props["number"], "c.Properties[\"number\"]")

	// process other parts
	for i := 1; i < 4; i++ {
		metadata.Part = null.IntFrom(i)
		metadata.Sources = suite.someSources()
		metadata.Tags = suite.someTags()
		tf := chain[fmt.Sprintf("part%d", i)]
		original, proxy := tf.Original, tf.Proxy

		evnts, err := ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
		suite.Require().Nil(err)
		suite.Require().NotNil(evnts)

		err = original.Reload(suite.tx)
		suite.Require().Nil(err)
		err = proxy.Reload(suite.tx)
		suite.Require().Nil(err)

		suite.assertFiles(metadata, original, proxy)
		suite.assertContentUnit(metadata, original, proxy, false)

		// collection association
		err = original.L.LoadContentUnit(suite.tx, true, original, nil)
		suite.Require().Nil(err)
		cu := original.R.ContentUnit
		err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
		suite.Require().Nil(err)
		suite.Equal(1, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")
		ccu := cu.R.CollectionsContentUnits[0]
		suite.Equal(strconv.Itoa(i), ccu.Name, "ccu.Name")
		suite.Equal(i, ccu.Position, "ccu.Position")
		suite.Equal(c.ID, ccu.CollectionID, "ccu.CollectionID")
	}

	// process full
	metadata.ContentType = common.CT_FULL_LESSON
	metadata.Part = null.IntFrom(-1)
	metadata.Sources = nil
	metadata.Tags = nil
	tf := chain["full"]
	original, proxy = tf.Original, tf.Proxy

	evnts, err = ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(evnts)

	err = original.Reload(suite.tx)
	suite.Require().Nil(err)
	err = proxy.Reload(suite.tx)
	suite.Require().Nil(err)

	suite.assertFiles(metadata, original, proxy)
	suite.assertContentUnit(metadata, original, proxy, false)

	// collection association
	err = original.L.LoadContentUnit(suite.tx, true, original, nil)
	suite.Require().Nil(err)
	cu = original.R.ContentUnit
	err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Equal(1, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")
	ccu = cu.R.CollectionsContentUnits[0]
	suite.Equal("full", ccu.Name, "ccu.Name")
	suite.Equal(4, ccu.Position, "ccu.Position")
	suite.Equal(c.ID, ccu.CollectionID, "ccu.CollectionID")

	// process kitei makor for all parts
	for i := 0; i < 4; i++ {
		metadata.ContentType = common.CT_LESSON_PART
		metadata.Part = null.IntFrom(i)
		metadata.ArtifactType = null.StringFrom("kitei_makor")
		metadata.WeekDate = nil
		tf = chain[fmt.Sprintf("part%d_kitei-makor", i)]
		original, proxy = tf.Original, tf.Proxy
		evnts, err = ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
		suite.Require().Nil(err)
		suite.Require().NotNil(evnts)

		err = original.Reload(suite.tx)
		suite.Require().Nil(err)
		err = proxy.Reload(suite.tx)
		suite.Require().Nil(err)

		suite.assertFiles(metadata, original, proxy)
		suite.assertContentUnit(metadata, original, proxy, false)

		// associated to "main" content unit
		err = original.L.LoadContentUnit(suite.tx, true, original, nil)
		suite.Require().Nil(err)
		cu = original.R.ContentUnit
		err = cu.L.LoadDerivedContentUnitDerivations(suite.tx, true, cu, nil)
		suite.Require().Nil(err)
		suite.Require().Len(cu.R.DerivedContentUnitDerivations, 1, "cu.R.DerivationContentUnitDerivations length")
		cud := cu.R.DerivedContentUnitDerivations[0]
		suite.Equal(chain[fmt.Sprintf("part%d", i)].Original.ContentUnitID.Int64, cud.SourceID, "cud.SourceID")
		suite.Equal("kitei_makor", cud.Name, "cud.Name")
		err = cu.Properties.Unmarshal(&props)
		suite.Require().Nil(err)
		_, ok := props["artifact_type"]
		suite.False(ok, "cu.properties[\"artifact_type\"]")

		// not associated with collection
		err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
		suite.Require().Nil(err)
		suite.Empty(cu.R.CollectionsContentUnits, "cu.R.CollectionsContentUnits empty")

		// no changes to collection
		err = c.Reload(suite.tx)
		suite.Require().Nil(err)
		suite.Equal(common.CONTENT_TYPE_REGISTRY.ByName[common.CT_DAILY_LESSON].ID, c.TypeID, "c.TypeID")
		suite.True(c.Properties.Valid, "c.Properties.Valid")
		err = json.Unmarshal(c.Properties.JSON, &props)
		suite.Require().Nil(err)
		suite.Equal(metadata.CaptureDate.Format("2006-01-02"), props["capture_date"], "c.Properties[\"capture_date\"]")
		suite.Equal(metadata.CaptureDate.Format("2006-01-02"), props["film_date"], "c.Properties[\"film_date\"]")
		suite.Equal("c12356789", props["capture_id"], "c.Properties[\"capture_id\"]")
		suite.EqualValues(metadata.Number.Int, props["number"], "c.Properties[\"number\"]")
	}

	// process lelo mikud for all parts
	for i := 0; i < 4; i++ {
		metadata.ContentType = common.CT_LESSON_PART
		metadata.Part = null.IntFrom(i)
		metadata.ArtifactType = null.StringFrom("lelo_mikud")
		metadata.WeekDate = nil
		tf = chain[fmt.Sprintf("part%d_lelo-mikud", i)]
		original, proxy = tf.Original, tf.Proxy
		evnts, err = ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
		suite.Require().Nil(err)
		suite.Require().NotNil(evnts)

		err = original.Reload(suite.tx)
		suite.Require().Nil(err)
		err = proxy.Reload(suite.tx)
		suite.Require().Nil(err)

		suite.assertFiles(metadata, original, proxy)
		suite.assertContentUnit(metadata, original, proxy, false)

		// associated to "main" content unit
		err = original.L.LoadContentUnit(suite.tx, true, original, nil)
		suite.Require().Nil(err)
		cu = original.R.ContentUnit
		err = cu.L.LoadDerivedContentUnitDerivations(suite.tx, true, cu, nil)
		suite.Require().Nil(err)
		suite.Require().Len(cu.R.DerivedContentUnitDerivations, 1, "cu.R.DerivationContentUnitDerivations length")
		cud := cu.R.DerivedContentUnitDerivations[0]
		suite.Equal(chain[fmt.Sprintf("part%d", i)].Original.ContentUnitID.Int64, cud.SourceID, "cud.SourceID")
		suite.Equal("lelo_mikud", cud.Name, "cud.Name")
		err = cu.Properties.Unmarshal(&props)
		suite.Require().Nil(err)
		_, ok := props["artifact_type"]
		suite.False(ok, "cu.properties[\"artifact_type\"]")

		// not associated with collection
		err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
		suite.Require().Nil(err)
		suite.Empty(cu.R.CollectionsContentUnits, "cu.R.CollectionsContentUnits empty")

		// no changes to collection
		err = c.Reload(suite.tx)
		suite.Require().Nil(err)
		suite.Equal(common.CONTENT_TYPE_REGISTRY.ByName[common.CT_DAILY_LESSON].ID, c.TypeID, "c.TypeID")
		suite.True(c.Properties.Valid, "c.Properties.Valid")
		err = json.Unmarshal(c.Properties.JSON, &props)
		suite.Require().Nil(err)
		suite.Equal(metadata.CaptureDate.Format("2006-01-02"), props["capture_date"], "c.Properties[\"capture_date\"]")
		suite.Equal(metadata.CaptureDate.Format("2006-01-02"), props["film_date"], "c.Properties[\"film_date\"]")
		suite.Equal("c12356789", props["capture_id"], "c.Properties[\"capture_id\"]")
		suite.EqualValues(metadata.Number.Int, props["number"], "c.Properties[\"number\"]")
	}
}

func (suite *MetadataProcessorSuite) TestDerivedBeforeMain() {
	chain := suite.simulateLessonChain()

	// send kitei makor of part 1
	// send part 1

	metadata := CITMetadata{
		ContentType:    common.CT_LESSON_PART,
		AutoName:       "auto_name",
		FinalName:      "final_name",
		CaptureDate:    Date{time.Now()},
		Language:       common.LANG_HEBREW,
		HasTranslation: true,
		Lecturer:       "rav",
		Number:         null.IntFrom(1),
		Part:           null.IntFrom(1),
		ArtifactType:   null.StringFrom("kitei_makor"),
		RequireTest:    false,
	}

	// process kitei makor for part 1
	tf := chain["part1_kitei-makor"]
	original, proxy := tf.Original, tf.Proxy
	evnts, err := ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(evnts)

	err = original.Reload(suite.tx)
	suite.Require().Nil(err)
	err = proxy.Reload(suite.tx)
	suite.Require().Nil(err)

	suite.assertFiles(metadata, original, proxy)
	suite.assertContentUnit(metadata, original, proxy, false)

	// not associated with collection
	err = original.L.LoadContentUnit(suite.tx, true, original, nil)
	suite.Require().Nil(err)
	cu := original.R.ContentUnit
	err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Empty(cu.R.CollectionsContentUnits, "cu.R.CollectionsContentUnits empty")

	// not associated to "main" content unit
	err = cu.L.LoadDerivedContentUnitDerivations(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Empty(cu.R.DerivedContentUnitDerivations, "cu.R.DerivationContentUnitDerivations empty")
	var props map[string]interface{}
	err = json.Unmarshal(cu.Properties.JSON, &props)
	suite.Require().Nil(err)
	suite.Equal("kitei_makor", props["artifact_type"], "cu.propeties[\"artifact_type\"]")

	// process main part1
	original, proxy = chain["part1"].Original, chain["part1"].Proxy
	metadata.ArtifactType = null.NewString("", false)
	evnts, err = ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(evnts)

	err = original.Reload(suite.tx)
	suite.Require().Nil(err)
	err = proxy.Reload(suite.tx)
	suite.Require().Nil(err)

	suite.assertFiles(metadata, original, proxy)
	suite.assertContentUnit(metadata, original, proxy, false)

	// reload cu cu association
	err = cu.L.LoadDerivedContentUnitDerivations(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Require().Len(cu.R.DerivedContentUnitDerivations, 1, "cu.R.DerivationContentUnitDerivations length")
	cud := cu.R.DerivedContentUnitDerivations[0]
	suite.Equal(chain["part1"].Original.ContentUnitID.Int64, cud.SourceID, "cud.SourceID")
	suite.Equal("kitei_makor", cud.Name, "cud.Name")

	err = cu.Reload(suite.tx)
	suite.Require().Nil(err)
	props = make(map[string]interface{})
	err = json.Unmarshal(cu.Properties.JSON, &props)
	suite.Require().Nil(err)
	_, ok := props["artifact_type"]
	suite.False(ok, "cu.propeties[\"artifact_type\"] presence")

	// main cu collection association
	err = original.L.LoadContentUnit(suite.tx, true, original, nil)
	suite.Require().Nil(err)
	cu = original.R.ContentUnit
	err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Equal(1, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")
	ccu := cu.R.CollectionsContentUnits[0]
	suite.Equal("1", ccu.Name, "ccu.Name")
	suite.Equal(0, ccu.Position, "ccu.Position")

	// collection
	err = ccu.L.LoadCollection(suite.tx, true, ccu, nil)
	suite.Require().Nil(err)
	c := ccu.R.Collection
	suite.Equal(common.CONTENT_TYPE_REGISTRY.ByName[common.CT_DAILY_LESSON].ID, c.TypeID, "c.TypeID")
	suite.True(c.Properties.Valid, "c.Properties.Valid")
	err = json.Unmarshal(c.Properties.JSON, &props)
	suite.Require().Nil(err)
	suite.Equal(metadata.CaptureDate.Format("2006-01-02"), props["capture_date"], "c.Properties[\"capture_date\"]")
	suite.Equal(metadata.CaptureDate.Format("2006-01-02"), props["film_date"], "c.Properties[\"film_date\"]")
	suite.Equal("c12356789", props["capture_id"], "c.Properties[\"capture_id\"]")
	suite.EqualValues(metadata.Number.Int, props["number"], "c.Properties[\"number\"]")

}

func (suite *MetadataProcessorSuite) TestVideoProgram() {
	tf, _ := suite.simulateSimpleChain()
	original, proxy := tf.Original, tf.Proxy

	c, err := CreateCollection(suite.tx, common.CT_VIDEO_PROGRAM, nil)
	suite.Require().Nil(err)

	metadata := CITMetadata{
		ContentType:    common.CT_VIDEO_PROGRAM_CHAPTER,
		AutoName:       "auto_name",
		FinalName:      "final_name",
		CaptureDate:    Date{time.Now()},
		Language:       common.LANG_HEBREW,
		HasTranslation: false,
		Lecturer:       "norav",
		CollectionUID:  null.StringFrom(c.UID),
		Episode:        null.StringFrom("827"),
		RequireTest:    true,
	}

	evnts, err := ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(evnts)

	err = original.Reload(suite.tx)
	suite.Require().Nil(err)
	err = proxy.Reload(suite.tx)
	suite.Require().Nil(err)

	suite.assertFiles(metadata, original, proxy)
	suite.assertContentUnit(metadata, original, proxy, false)

	// collection association
	err = original.L.LoadContentUnit(suite.tx, true, original, nil)
	suite.Require().Nil(err)
	cu := original.R.ContentUnit
	err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Equal(1, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")
	ccu := cu.R.CollectionsContentUnits[0]
	suite.Equal(c.ID, ccu.CollectionID, "ccu.CollectionID")
	suite.Equal(metadata.Episode.String, ccu.Name, "ccu.Name")
	suite.Equal(0, ccu.Position, "ccu.Position")
}

func (suite *MetadataProcessorSuite) TestEventPart() {
	tf, _ := suite.simulateSimpleChain()
	original, proxy := tf.Original, tf.Proxy

	EVENT_TYPES := [...]string{common.CT_CONGRESS, common.CT_HOLIDAY, common.CT_PICNIC, common.CT_UNITY_DAY}
	EVENT_PART_TYPES := [...]string{common.CT_FRIENDS_GATHERING, common.CT_MEAL,
		common.CT_EVENT_PART, common.CT_EVENT_PART, common.CT_EVENT_PART, common.CT_EVENT_PART,
		common.CT_EVENT_PART, common.CT_EVENT_PART, common.CT_EVENT_PART, common.CT_EVENT_PART}

	for _, eventType := range EVENT_TYPES {
		c, err := CreateCollection(suite.tx, eventType, nil)
		suite.Require().Nil(err)

		for i, partType := range EVENT_PART_TYPES {
			metadata := CITMetadata{
				ContentType:    partType,
				AutoName:       "auto_name",
				FinalName:      "final_name",
				CaptureDate:    Date{time.Now()},
				Language:       common.LANG_HEBREW,
				HasTranslation: true,
				CollectionUID:  null.StringFrom(c.UID),
				Number:         null.IntFrom(i + 1),
				RequireTest:    true,
				PartType:       null.IntFrom(i),
				Lecturer:       "norav",
			}

			//if partType == CT_FULL_LESSON {
			//	metadata.Lecturer = "rav"
			//	metadata.Sources = suite.someSources()
			//	metadata.Tags = suite.someTags()
			//} else {
			//	metadata.Lecturer = "norav"
			//}

			evnts, err := ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
			suite.Require().Nil(err)
			suite.Require().NotNil(evnts)

			err = original.Reload(suite.tx)
			suite.Require().Nil(err)
			err = proxy.Reload(suite.tx)
			suite.Require().Nil(err)

			suite.assertFiles(metadata, original, proxy)
			suite.assertContentUnit(metadata, original, proxy, false)

			// collection association
			err = original.L.LoadContentUnit(suite.tx, true, original, nil)
			suite.Require().Nil(err)
			cu := original.R.ContentUnit
			err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
			suite.Require().Nil(err)
			suite.Equal(1, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")
			ccu := cu.R.CollectionsContentUnits[0]
			suite.Equal(c.ID, ccu.CollectionID, "ccu.CollectionID")
			if i < 3 {
				suite.Equal(strconv.Itoa(metadata.Number.Int), ccu.Name, "ccu.Name")
			} else {
				suite.Equal(common.MISC_EVENT_PART_TYPES[i-3]+strconv.Itoa(metadata.Number.Int),
					ccu.Name, "ccu.Name")
			}
			suite.Equal(i, ccu.Position, "ccu.Position")
		}
	}
}

func (suite *MetadataProcessorSuite) TestEventPartLesson() {
	chain := suite.simulateLessonChain()

	eventType := common.CT_CONGRESS
	eventCollection, err := CreateCollection(suite.tx, eventType, nil)
	suite.Require().Nil(err)

	// send parts
	// send full

	metadata := CITMetadata{
		ContentType:    common.CT_LESSON_PART,
		AutoName:       "auto_name",
		FinalName:      "final_name",
		CaptureDate:    Date{time.Now()},
		Language:       common.LANG_HEBREW,
		HasTranslation: true,
		Lecturer:       "rav",
		CollectionUID:  null.StringFrom(eventCollection.UID),
		Number:         null.IntFrom(1),
		Part:           null.IntFrom(0),
		Sources:        suite.someSources(),
		Tags:           suite.someTags(),
		RequireTest:    false,
	}
	original, proxy := chain["part0"].Original, chain["part0"].Proxy

	evnts, err := ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(evnts)

	err = original.Reload(suite.tx)
	suite.Require().Nil(err)
	err = proxy.Reload(suite.tx)
	suite.Require().Nil(err)

	suite.assertFiles(metadata, original, proxy)
	suite.assertContentUnit(metadata, original, proxy, false)

	// collection association
	err = original.L.LoadContentUnit(suite.tx, true, original, nil)
	suite.Require().Nil(err)
	cu := original.R.ContentUnit
	err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Equal(2, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")

	// load ccu collections
	var lccu, eccu *models.CollectionsContentUnit
	for i := range cu.R.CollectionsContentUnits {
		ccu := cu.R.CollectionsContentUnits[i]
		err = ccu.L.LoadCollection(suite.tx, true, ccu, nil)
		suite.Require().Nil(err)
		switch common.CONTENT_TYPE_REGISTRY.ByID[ccu.R.Collection.TypeID].Name {
		case common.CT_DAILY_LESSON:
			lccu = ccu
		case eventType:
			eccu = ccu
		default:
			suite.FailNow("ccu.collection type %s", common.CONTENT_TYPE_REGISTRY.ByID[ccu.R.Collection.TypeID].Name)
		}
	}

	// lesson collection association
	suite.Equal("0", lccu.Name, "lccu.Name")
	suite.Equal(0, lccu.Position, "lccu.Position")
	lessonCollection := lccu.R.Collection
	suite.Equal(common.CONTENT_TYPE_REGISTRY.ByName[common.CT_DAILY_LESSON].ID, lessonCollection.TypeID, "c.TypeID")
	suite.True(lessonCollection.Properties.Valid, "c.Properties.Valid")
	var props map[string]interface{}
	err = json.Unmarshal(lessonCollection.Properties.JSON, &props)
	suite.Require().Nil(err)
	suite.Equal(metadata.CaptureDate.Format("2006-01-02"), props["capture_date"], "c.Properties[\"capture_date\"]")
	suite.Equal(metadata.CaptureDate.Format("2006-01-02"), props["film_date"], "c.Properties[\"film_date\"]")
	suite.Equal("c12356789", props["capture_id"], "c.Properties[\"capture_id\"]")
	suite.EqualValues(metadata.Number.Int, props["number"], "c.Properties[\"number\"]")

	// event collection association
	suite.Equal("0", eccu.Name, "eccu.Name")
	suite.Equal(eventCollection.ID, eccu.CollectionID, "eccu.CollectionID")

	// process other parts
	for i := 1; i < 4; i++ {
		metadata.Part = null.IntFrom(i)
		metadata.Sources = suite.someSources()
		metadata.Tags = suite.someTags()
		tf := chain[fmt.Sprintf("part%d", i)]
		original, proxy := tf.Original, tf.Proxy

		evnts, err := ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
		suite.Require().Nil(err)
		suite.Require().NotNil(evnts)

		err = original.Reload(suite.tx)
		suite.Require().Nil(err)
		err = proxy.Reload(suite.tx)
		suite.Require().Nil(err)

		suite.assertFiles(metadata, original, proxy)
		suite.assertContentUnit(metadata, original, proxy, false)

		// collection association
		err = original.L.LoadContentUnit(suite.tx, true, original, nil)
		suite.Require().Nil(err)
		cu := original.R.ContentUnit
		err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
		suite.Require().Nil(err)
		suite.Equal(2, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")

		// load ccu collections
		var lccu, eccu *models.CollectionsContentUnit
		for j := range cu.R.CollectionsContentUnits {
			ccu := cu.R.CollectionsContentUnits[j]
			err = ccu.L.LoadCollection(suite.tx, true, ccu, nil)
			suite.Require().Nil(err)
			switch common.CONTENT_TYPE_REGISTRY.ByID[ccu.R.Collection.TypeID].Name {
			case common.CT_DAILY_LESSON:
				lccu = ccu
			case eventType:
				eccu = ccu
			default:
				suite.FailNow("ccu.collection type %s", common.CONTENT_TYPE_REGISTRY.ByID[ccu.R.Collection.TypeID].Name)
			}
		}

		// lesson collection association
		suite.Equal(strconv.Itoa(i), lccu.Name, "lccu.Name")
		suite.NotZero(lccu.Position, "lccu.Position")
		suite.Equal(lessonCollection.ID, lccu.CollectionID, "lccu.CollectionID")

		// event collection association
		suite.Equal(strconv.Itoa(i), eccu.Name, "eccu.Name")
		suite.Equal(eventCollection.ID, eccu.CollectionID, "eccu.CollectionID")
	}

	// process full
	metadata.ContentType = common.CT_FULL_LESSON
	metadata.Part = null.IntFrom(-1)
	metadata.Sources = nil
	metadata.Tags = nil
	tf := chain["full"]
	original, proxy = tf.Original, tf.Proxy

	evnts, err = ProcessCITMetadata(suite.tx, metadata, original, proxy, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(evnts)

	err = original.Reload(suite.tx)
	suite.Require().Nil(err)
	err = proxy.Reload(suite.tx)
	suite.Require().Nil(err)

	suite.assertFiles(metadata, original, proxy)
	suite.assertContentUnit(metadata, original, proxy, false)

	// collection association
	err = original.L.LoadContentUnit(suite.tx, true, original, nil)
	suite.Require().Nil(err)
	cu = original.R.ContentUnit
	err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Equal(2, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")

	// load ccu collections
	lccu = nil
	eccu = nil
	for j := range cu.R.CollectionsContentUnits {
		ccu := cu.R.CollectionsContentUnits[j]
		err = ccu.L.LoadCollection(suite.tx, true, ccu, nil)
		suite.Require().Nil(err)
		switch common.CONTENT_TYPE_REGISTRY.ByID[ccu.R.Collection.TypeID].Name {
		case common.CT_DAILY_LESSON:
			lccu = ccu
		case eventType:
			eccu = ccu
		default:
			suite.FailNow("ccu.collection type %s", common.CONTENT_TYPE_REGISTRY.ByID[ccu.R.Collection.TypeID].Name)
		}
	}

	// lesson collection association
	suite.Equal("full", lccu.Name, "lccu.Name")
	suite.Equal(4, lccu.Position, "lccu.Position")
	suite.Equal(lessonCollection.ID, lccu.CollectionID, "lccu.CollectionID")

	// event collection association
	suite.Equal(strconv.Itoa(metadata.Number.Int), eccu.Name, "eccu.Name")
	suite.Equal(eventCollection.ID, eccu.CollectionID, "eccu.CollectionID")

}

func (suite *MetadataProcessorSuite) TestFixUnit() {
	chain := suite.simulateSpecialLessonChain()
	originals := make(map[string]TrimFiles)

	// send parts
	// send full
	// send kitei makor of all parts
	// send lelo mikud of all parts

	metadata := CITMetadata{
		ContentType:    common.CT_LESSON_PART,
		AutoName:       "auto_name",
		FinalName:      "final_name",
		CaptureDate:    Date{time.Now()},
		Language:       common.LANG_HEBREW,
		HasTranslation: true,
		Lecturer:       "rav",
		Number:         null.IntFrom(1),
		Part:           null.IntFrom(0),
		Sources:        suite.someSources(),
		Tags:           suite.someTags(),
		RequireTest:    false,
	}
	tf := chain["part0"]
	originals["part0"] = tf

	_, err := ProcessCITMetadata(suite.tx, metadata, tf.Original, tf.Proxy, nil)
	suite.Require().Nil(err)

	// process other parts
	for i := 1; i < 4; i++ {
		metadata.Part = null.IntFrom(i)
		metadata.Sources = suite.someSources()
		metadata.Tags = suite.someTags()
		tf := chain[fmt.Sprintf("part%d", i)]
		originals[fmt.Sprintf("part%d", i)] = tf

		_, err := ProcessCITMetadata(suite.tx, metadata, tf.Original, tf.Proxy, nil)
		suite.Require().Nil(err)
	}

	// process full
	metadata.ContentType = common.CT_FULL_LESSON
	metadata.Part = null.IntFrom(-1)
	metadata.Sources = nil
	metadata.Tags = nil
	tf = chain["full"]
	originals["full"] = tf

	_, err = ProcessCITMetadata(suite.tx, metadata, tf.Original, tf.Proxy, nil)
	suite.Require().Nil(err)

	// process kitei makor for all parts
	for i := 0; i < 4; i++ {
		metadata.ContentType = common.CT_LESSON_PART
		metadata.Part = null.IntFrom(i)
		metadata.ArtifactType = null.StringFrom("kitei_makor")
		metadata.WeekDate = nil
		tf = chain[fmt.Sprintf("part%d_kitei-makor", i)]
		originals[fmt.Sprintf("part%d_kitei-makor", i)] = tf

		_, err = ProcessCITMetadata(suite.tx, metadata, tf.Original, tf.Proxy, nil)
		suite.Require().Nil(err)
	}

	// process lelo mikud for all parts
	for i := 0; i < 4; i++ {
		metadata.ContentType = common.CT_LESSON_PART
		metadata.Part = null.IntFrom(i)
		metadata.ArtifactType = null.StringFrom("lelo_mikud")
		metadata.WeekDate = nil
		tf = chain[fmt.Sprintf("part%d_lelo-mikud", i)]
		originals[fmt.Sprintf("part%d_lelo-mikud", i)] = tf

		_, err = ProcessCITMetadata(suite.tx, metadata, tf.Original, tf.Proxy, nil)
		suite.Require().Nil(err)
	}

	// process updated metadata
	tf = originals["part2"]
	convertFiles := suite.simulateConvertUpload(tf.Original)

	cu, err := tf.Original.ContentUnit().One(suite.tx)
	suite.Require().Nil(err)

	metadata.UnitToFixUID = null.StringFrom(cu.UID)
	metadata.ContentType = common.CT_CLIP
	metadata.ArtifactType = null.NewString("", false)
	metadata.AutoName = "auto_name_update"
	metadata.FinalName = "final_name_update"
	metadata.Language = common.LANG_ENGLISH
	metadata.HasTranslation = false
	metadata.Lecturer = "norav"
	metadata.Sources = suite.someSources()
	metadata.Tags = suite.someTags()

	suite.True(cu.Published, "cu.Published before fix")

	evnts, err := ProcessCITMetadataUpdate(suite.tx, metadata, tf.Original, tf.Proxy, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(evnts)

	// verify CU changed metadata
	err = cu.Reload(suite.tx)
	suite.Require().Nil(err)
	suite.assertFiles(metadata, tf.Original, tf.Proxy)
	suite.assertContentUnit(metadata, tf.Original, tf.Proxy, true)
	suite.False(cu.Published, "cu.Published after fix")

	// verify removed files
	for i, cf := range convertFiles {
		err = cf.Reload(suite.tx)
		suite.Require().Nil(err)
		suite.True(cf.RemovedAt.Valid, "cf.RemovedAt.Valid %d", i)
	}
}

func (suite *MetadataProcessorSuite) TestWithAdditionalCapture() {
	originalTf, _ := suite.simulateSimpleChain()

	metadata := CITMetadata{
		ContentType:    common.CT_LESSON_PART,
		AutoName:       "auto_name",
		FinalName:      "final_name",
		CaptureDate:    Date{time.Now()},
		Language:       common.LANG_HEBREW,
		HasTranslation: true,
		Lecturer:       "rav",
		Number:         null.IntFrom(1),
		Part:           null.IntFrom(0),
		Sources:        suite.someSources(),
		Tags:           suite.someTags(),
		RequireTest:    false,
	}

	metadata.Sources = suite.someSources()
	metadata.Tags = suite.someTags()
	original, proxy := originalTf.Original, originalTf.Proxy

	souce := suite.simulateAdditionalCapture()
	evnts, err := ProcessCITMetadata(suite.tx, metadata, original, proxy, souce)
	suite.Require().Nil(err)
	suite.Require().NotNil(evnts)

	err = original.Reload(suite.tx)
	suite.Require().Nil(err)

	// verify add source
	err = souce.Reload(suite.tx)
	suite.Require().Nil(err)
	suite.Equal(originalTf.Original.ContentUnitID.Int64, souce.ContentUnitID.Int64, "original and source unit")

	ancestors, err := FindFileAncestors(suite.tx, souce.ID)
	suite.Require().Nil(err)
	for _, ancestor := range ancestors {
		suite.Equal(originalTf.Original.ContentUnitID.Int64, ancestor.ContentUnitID.Int64, "source ancestor files")
	}
}

func (suite *MetadataProcessorSuite) TestDoubleTrimFromOneStart() {
	CS_SHA1 := utils.RandomSHA1()
	DMX_SHA1 := utils.RandomSHA1()
	TRM_1_SHA1 := utils.RandomSHA1()
	TRM_2_SHA1 := utils.RandomSHA1()
	WorkflowID := "c12356788"

	// capture_start
	_, evnts, err := handleCaptureStart(suite.tx, CaptureStartRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: WorkflowID,
		},
		FileName:      "capture_start_simple",
		CaptureSource: "mltcap",
		CollectionUID: "c12356788",
	})
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	// capture_stop
	_, evnts, err = handleCaptureStop(suite.tx, CaptureStopRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: WorkflowID,
		},
		File: File{
			FileName:  "capture_stop_double_trim.mp4",
			Sha1:      CS_SHA1,
			Size:      98737,
			CreatedAt: &Timestamp{Time: time.Now()},
			Language:  common.LANG_HEBREW,
		},
		CaptureSource: "mltcap",
		CollectionUID: "c12356789",
		Part:          "false",
	})
	suite.Require().Nil(err)

	// demux
	_, evnts, err = handleDemux(suite.tx, DemuxRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		Sha1: CS_SHA1,
		Original: AVFile{
			File: File{
				FileName:  "demux_double_trim.mp4",
				Sha1:      DMX_SHA1,
				Size:      98737,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		CaptureSource: "mltcap",
	})
	suite.Require().Nil(err)

	// trim
	op_1, _, err := handleTrim(suite.tx, TrimRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		OriginalSha1: DMX_SHA1,
		Original: AVFile{
			File: File{
				FileName:  "double_trim_1.mp4",
				Sha1:      TRM_1_SHA1,
				Size:      98000,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		CaptureSource: "mltcap",
		In:            []float64{10.05, 249.43},
		Out:           []float64{240.51, 899.27},
	})
	suite.Require().Nil(err)

	op_2, _, err := handleTrim(suite.tx, TrimRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		OriginalSha1: DMX_SHA1,
		Original: AVFile{
			File: File{
				FileName:  "double_trim_2.mp4",
				Sha1:      TRM_2_SHA1,
				Size:      98000,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		CaptureSource: "mltcap",
		In:            []float64{10.05, 249.43},
		Out:           []float64{240.51, 899.27},
	})
	suite.Require().Nil(err)

	original_1 := suite.opFilesBySHA1(op_1)[TRM_1_SHA1]
	original_2 := suite.opFilesBySHA1(op_2)[TRM_2_SHA1]

	metadata := CITMetadata{
		ContentType:    common.CT_LESSON_PART,
		AutoName:       "auto_name",
		FinalName:      "final_name",
		CaptureDate:    Date{time.Now()},
		Language:       common.LANG_HEBREW,
		HasTranslation: true,
		Lecturer:       "rav",
		Number:         null.IntFrom(1),
		Part:           null.IntFrom(0),
		Sources:        suite.someSources(),
		Tags:           suite.someTags(),
		RequireTest:    false,
	}

	metadata.Sources = suite.someSources()
	metadata.Tags = suite.someTags()

	souce := suite.simulateAdditionalCapture()
	_, err = ProcessCITMetadata(suite.tx, metadata, original_1, nil, souce)
	suite.Require().Nil(err)

	// verify add source
	suite.Require().Nil(original_1.Reload(suite.tx))
	suite.Require().Nil(souce.Reload(suite.tx))
	suite.Equal(original_1.ContentUnitID.Int64, souce.ContentUnitID.Int64, "first original and source unit")

	_, err = ProcessCITMetadata(suite.tx, metadata, original_2, nil, nil)
	suite.Require().Nil(err)
	suite.Require().Nil(original_1.Reload(suite.tx))
	suite.Require().Nil(original_2.Reload(suite.tx))
	suite.Require().Nil(souce.Reload(suite.tx))
	suite.Equal(original_1.ContentUnitID.Int64, souce.ContentUnitID.Int64, "first original and source unit")
	suite.NotEqual(original_2.ContentUnitID.Int64, souce.ContentUnitID.Int64, "second original and source unit")
}

func (suite *MetadataProcessorSuite) TestReplaceProcess() {
	tfMain, WorkflowID := suite.simulateSimpleChain()
	original, proxy := tfMain.Original, tfMain.Proxy

	c, err := CreateCollection(suite.tx, common.CT_VIDEO_PROGRAM, nil)
	suite.Require().Nil(err)

	metadata := CITMetadata{
		ContentType:    common.CT_VIDEO_PROGRAM_CHAPTER,
		AutoName:       "auto_name",
		FinalName:      "final_name",
		CaptureDate:    Date{time.Now()},
		Language:       common.LANG_HEBREW,
		HasTranslation: false,
		Lecturer:       "norav",
		CollectionUID:  null.StringFrom(c.UID),
		Episode:        null.StringFrom("827"),
		RequireTest:    true,
	}

	// HLS capture_stop
	CS_SHA1 := utils.RandomSHA1()
	_, _, err = handleCaptureStop(suite.tx, CaptureStopRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: WorkflowID,
		},
		File: File{
			FileName:  "capture_stop_hls.mp4",
			Sha1:      CS_SHA1,
			Size:      98737,
			CreatedAt: &Timestamp{Time: time.Now()},
			Language:  common.LANG_HEBREW,
		},
		CaptureSource: "archcap",
		CollectionUID: "c12356789",
		Part:          "false",
	})
	suite.Require().Nil(err)

	// HLS trim
	HLS1_TRM_SHA1 := utils.RandomSHA1()
	opTrim, _, err := handleTrim(suite.tx, TrimRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		OriginalSha1: CS_SHA1,
		Original: AVFile{
			File: File{
				FileName:  "trim_hls_original.mp4",
				Sha1:      HLS1_TRM_SHA1,
				Size:      98000,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.0,
		},
		CaptureSource: "archcap",
		In:            []float64{10.05, 249.43},
		Out:           []float64{240.51, 899.27},
	})
	suite.Require().Nil(err)

	trimFiles := suite.opFilesBySHA1(opTrim)
	suite.Require().Nil(opTrim.L.LoadUser(suite.tx, true, opTrim, nil))
	hlsTrimFile := trimFiles[HLS1_TRM_SHA1]

	//HLS convert old
	HLS1_SHA1 := utils.RandomSHA1()
	OLD_HLS_LANGS := []string{"he", "ru"}
	OLD_HLS_QUALITIES := []string{"HD"}
	inputOldConvert := ConvertRequest{
		Operation: Operation{
			Station: "Convert station",
			User:    "operator@dev.com",
		},
		Sha1: HLS1_TRM_SHA1,
		Output: []HLSFile{
			{
				Languages: OLD_HLS_LANGS,
				Qualities: OLD_HLS_QUALITIES,
				AVFile: AVFile{
					File: File{
						FileName:  "test_file.mp4",
						Sha1:      HLS1_SHA1,
						Size:      694,
						CreatedAt: &Timestamp{Time: time.Now()},
						Type:      "video",
						MimeType:  "video/mp4",
					},
					Duration: 871,
				},
			},
		},
	}

	op, _, err := handleConvert(suite.tx, inputOldConvert)
	suite.Require().Nil(err)

	files := suite.opFilesBySHA1(op)
	suite.Require().Len(files, 2)
	suite.Require().Nil(op.L.LoadUser(suite.tx, true, op, nil))

	hlsFile1 := files[HLS1_SHA1]
	//HLS send old
	req := SendRequest{
		Operation: Operation{
			Station:    op.Station.String,
			User:       op.R.User.Email,
			WorkflowID: WorkflowID,
		},
		Original: Rename{
			Sha1:     hex.EncodeToString(original.Sha1.Bytes),
			FileName: original.Name,
		},
		Proxy: &Rename{
			Sha1:     hex.EncodeToString(proxy.Sha1.Bytes),
			FileName: proxy.Name,
		},
		Source: &Rename{
			Sha1:     hex.EncodeToString(hlsFile1.Sha1.Bytes),
			FileName: hlsFile1.Name,
		},
		Metadata: metadata,
		Mode:     null.String{},
	}

	_, _, err = handleSend(suite.tx, req)
	suite.Require().Nil(err)
	suite.Require().Nil(original.Reload(suite.tx))
	suite.Require().Nil(hlsFile1.Reload(suite.tx))
	suite.Require().Equal(hlsTrimFile.ID, hlsFile1.ParentID.Int64)
	suite.Require().Equal(original.ContentUnitID, hlsFile1.ContentUnitID)

	suite.Require().Nil(hlsFile1.L.LoadContentUnit(suite.tx, true, hlsFile1, nil))
	hlsFile1.Published = true
	_, err = hlsFile1.Update(suite.tx, boil.Whitelist("published"))
	suite.Require().Nil(err)

	cu := hlsFile1.R.ContentUnit
	var propsCu map[string]interface{}
	suite.Require().Nil(cu.Properties.Unmarshal(&propsCu))
	var props1 map[string]interface{}
	suite.Require().Nil(hlsFile1.Properties.Unmarshal(&props1))
	suite.Require().EqualValues(utils.ConvertArgsString(OLD_HLS_LANGS), props1["languages"])
	suite.Require().EqualValues(utils.ConvertArgsString(OLD_HLS_QUALITIES), props1["video_qualities"])

	// new HLS replace
	HLS2_SHA1 := utils.RandomSHA1()
	NEW_HLS_QUALITIES := []string{"HD", "fHD"}
	NEW_HLS_LANGS := []string{"he", "ru", "en"}
	hls2FileReq := HLSFile{
		AVFile: AVFile{
			File: File{
				FileName:  "Replaced HLS file",
				Sha1:      HLS2_SHA1,
				Size:      98000,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 820.0,
		},
		Languages: NEW_HLS_LANGS,
		Qualities: NEW_HLS_QUALITIES,
	}
	reqReplace := ReplaceRequest{
		Operation: Operation{
			Station: "Replace station",
			User:    "operator@dev.com",
		},
		HLSFile: hls2FileReq,
		OldSha1: HLS1_SHA1,
	}
	opHLS, _, err := handleReplace(suite.tx, reqReplace)
	suite.Require().Nil(err)
	hlsFile2 := suite.opFilesBySHA1(opHLS)[HLS2_SHA1]
	suite.Require().Equal(hlsTrimFile.ID, hlsFile2.ParentID.Int64)
	suite.Require().Equal(original.ContentUnitID, hlsFile2.ContentUnitID)
	suite.Require().False(hlsFile2.Published)
	var props2 map[string]interface{}
	suite.Require().Nil(hlsFile2.Properties.Unmarshal(&props2))
	suite.Require().EqualValues(utils.ConvertArgsString(NEW_HLS_LANGS), props2["languages"])
	suite.Require().EqualValues(utils.ConvertArgsString(NEW_HLS_QUALITIES), props2["video_qualities"])

	reqUpload := UploadRequest{
		Operation: Operation{
			Station: "Upload station",
			User:    "operator@dev.com",
		},
		AVFile: hls2FileReq.AVFile,
		Url:    "http://example.com/some/url/to/file.mp4",
	}
	_, _, err = handleUpload(suite.tx, reqUpload)
	suite.Require().Nil(err)

	suite.Require().Nil(hlsFile2.Reload(suite.tx))
	suite.Require().Nil(hlsFile1.Reload(suite.tx))
	suite.Require().True(hlsFile2.Published)
	suite.Require().NotNil(hlsFile1.RemovedAt)
}

func (suite *MetadataProcessorSuite) TestReplaceNotPublishedProcess() {
	// Create dummy content unit
	cu, err := CreateContentUnit(suite.tx, common.CT_LESSON_PART, nil)
	suite.Require().Nil(err)

	SHA_INS := utils.RandomSHA1()
	// Do insert operation
	input := InsertRequest{
		Operation: Operation{
			Station:    "Some station",
			User:       "operator@dev.com",
			WorkflowID: "workflow_id",
		},
		InsertType:     "akladot",
		ContentUnitUID: cu.UID,
		AVFile: AVFile{
			File: File{
				FileName:  "akladot.doc",
				Sha1:      SHA_INS,
				Size:      98737,
				CreatedAt: &Timestamp{Time: time.Now()},
				MimeType:  "application/msword",
				Language:  common.LANG_HEBREW,
			},
			Duration: 123.4,
		},
		Mode: "new",
	}

	_, _, err = handleInsert(suite.tx, input)
	suite.Require().Nil(err)
	insFile, _, err := FindFileBySHA1(suite.tx, SHA_INS)
	suite.Require().Nil(err)
	insFile.Published = false
	_, err = insFile.Update(suite.tx, boil.Whitelist("published"))
	suite.Require().Nil(err)

	SHA_NEW := utils.RandomSHA1()
	reqReplace := ReplaceRequest{
		Operation: Operation{
			Station: "Replace station",
			User:    "operator@dev.com",
		},
		HLSFile: HLSFile{
			AVFile: AVFile{
				File: File{
					FileName:  "Replaced file",
					Sha1:      SHA_NEW,
					Size:      98000,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 820.0,
			},
		},
		OldSha1: SHA_INS,
	}
	_, _, err = handleReplace(suite.tx, reqReplace)
	suite.Require().Nil(err)

	newFile, _, err := FindFileBySHA1(suite.tx, SHA_NEW)
	suite.Require().Nil(err)
	suite.Require().Nil(insFile.Reload(suite.tx))
	suite.Require().NotNil(insFile.RemovedAt)
	suite.Require().Equal(newFile.ContentUnitID, insFile.ContentUnitID)
	suite.Require().Equal(newFile.ParentID, insFile.ParentID)
}

func (suite *MetadataProcessorSuite) TestDailyLesson_SourcesAttachLessonsSeries() {
	tf, _ := suite.simulateSimpleChain()
	sUids := createDummySources(suite.tx, nil)
	metadata := CITMetadata{
		ContentType:    common.CT_LESSON_PART,
		AutoName:       "auto_name",
		FinalName:      "final_name",
		CaptureDate:    Date{time.Now()},
		Language:       common.LANG_HEBREW,
		HasTranslation: true,
		Lecturer:       "rav",
		Number:         null.IntFrom(1),
		Part:           null.IntFrom(0),
		Sources:        sUids,
		Tags:           suite.someTags(),
		RequireTest:    false,
	}

	props := map[string]interface{}{"source": sUids[0]}
	c, err := CreateCollection(suite.tx, common.CT_LESSONS_SERIES, props)
	utils.Must(err)
	c.Published = true
	_, err = c.Update(suite.tx, boil.Infer())
	s, err := models.Sources(models.SourceWhere.UID.EQ(metadata.Sources[0])).One(suite.tx)
	utils.Must(err)

	for i := 0; i < MinCuNumberForNewLessonSeries; i++ {
		_ = createCUWithSourceForLessonsSeries(suite.tx, s, c, i)
	}

	_, err = ProcessCITMetadata(suite.tx, metadata, tf.Original, tf.Proxy, nil)
	suite.Require().Nil(err)
	_, err = makePublishedLast(suite.tx)
	suite.Require().Nil(err)
	countCcu, err := models.CollectionsContentUnits(models.CollectionsContentUnitWhere.CollectionID.EQ(c.ID)).Count(suite.tx)
	suite.Require().Nil(err)
	suite.EqualValues(MinCuNumberForNewLessonSeries+1, countCcu)

	for i := 0; i <= MinCuNumberForNewLessonSeries; i++ {
		_, err = ProcessCITMetadata(suite.tx, metadata, tf.Original, tf.Proxy, nil)
		suite.Require().Nil(err)
		last, err := makePublishedLast(suite.tx)
		utils.Must(c.Reload(suite.tx))
		_countCcu, err := models.Collections(qm.WhereIn(`properties->>'source' IN ?`, utils.ConvertArgsString(sUids)...)).Count(suite.tx)
		utils.Must(err)
		if i < MinCuNumberForNewLessonSeries-1 {
			suite.EqualValues(1, _countCcu)
		} else {
			suite.EqualValues(len(sUids), _countCcu)
		}
		ccus, err := models.CollectionsContentUnits(models.CollectionsContentUnitWhere.CollectionID.EQ(c.ID)).All(suite.tx)
		utils.Must(err)
		suite.EqualValues(i+int(countCcu)+1, len(ccus))

		lastPosition := 0
		for _, ccu := range ccus {
			if ccu.ContentUnitID == last.ID {
				lastPosition = ccu.Position
			}
		}
		suite.EqualValues(lastPosition, len(ccus))
	}
}

func (suite *MetadataProcessorSuite) TestDailyLesson_SourcesTESAttachLessonsSeries() {
	s, err := models.Sources(qm.OrderBy("id DESC")).One(suite.tx)
	suite.Require().Nil(err)
	//prepare TES root
	rootUID := TES_PARTS_UIDS[rand.Intn(len(TES_PARTS_UIDS)-1)]
	author := &models.Author{Code: "bs", Name: "TES Author", FullName: null.StringFrom("test author")}
	suite.Require().Nil(author.Insert(suite.tx, boil.Infer()))
	tesRoot := &models.Source{ID: s.ID + 1, UID: rootUID, TypeID: 1, Name: "TES root source"}
	suite.Require().Nil(tesRoot.Insert(suite.tx, boil.Infer()))
	suite.Require().Nil((tesRoot.AddAuthors(suite.tx, false, author)))

	sUids := createDummySources(suite.tx, author)
	sources, err := models.Sources(models.SourceWhere.UID.IN(sUids)).All(suite.tx)
	suite.Require().Nil(err)
	_, err = sources.UpdateAll(suite.tx, models.M{"parent_id": tesRoot.ID})
	suite.Require().Nil(err)

	metadata := CITMetadata{
		ContentType:    common.CT_LESSON_PART,
		AutoName:       "auto_name",
		FinalName:      "final_name",
		CaptureDate:    Date{time.Now()},
		Language:       common.LANG_HEBREW,
		HasTranslation: true,
		Lecturer:       "rav",
		Number:         null.IntFrom(1),
		Part:           null.IntFrom(0),
		Tags:           suite.someTags(),
		RequireTest:    false,
	}

	for i := 0; i <= MinCuNumberForNewLessonSeries+1; i++ {
		tf, _ := suite.simulateSimpleChain()
		_sUid := sUids[i%len(sUids)]
		metadata.Sources = []string{_sUid}
		_, err = ProcessCITMetadata(suite.tx, metadata, tf.Original, tf.Proxy, nil)
		suite.Require().Nil(err)
		_, err = makePublishedLast(suite.tx)
		suite.Require().Nil(err)
	}

	c, err := models.Collections(
		models.CollectionWhere.TypeID.EQ(common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSONS_SERIES].ID),
		qm.OrderBy("id DESC"),
	).One(suite.tx)
	var props map[string]interface{}
	err = json.Unmarshal(c.Properties.JSON, &props)
	suite.Require().Nil(err)
	suite.Equal(rootUID, props["source"])
}

func (suite *MetadataProcessorSuite) TestDailyLesson_LikutimsAttachLessonsSeries() {
	tf, _ := suite.simulateSimpleChain()
	likutim, err := createDummyLikutim(suite.tx)
	utils.Must(err)
	lUids := make([]string, len(likutim))
	for i, l := range likutim {
		lUids[i] = l.UID
	}
	metadata := CITMetadata{
		ContentType:    common.CT_LESSON_PART,
		AutoName:       "auto_name",
		FinalName:      "final_name",
		CaptureDate:    Date{time.Now()},
		Language:       common.LANG_HEBREW,
		HasTranslation: true,
		Lecturer:       "rav",
		Number:         null.IntFrom(1),
		Part:           null.IntFrom(0),
		Likutim:        lUids,
		Tags:           suite.someTags(),
		RequireTest:    false,
	}

	props := map[string]interface{}{"likutim": []string{lUids[0]}}
	c, err := CreateCollection(suite.tx, common.CT_LESSONS_SERIES, props)
	utils.Must(err)
	c.Published = true
	_, err = c.Update(suite.tx, boil.Infer())

	for i := 0; i < MinCuNumberForNewLessonSeries; i++ {
		_ = createCUWithLikutForLessonsSeries(suite.tx, likutim[0], c, i)
	}

	countCcu, err := models.CollectionsContentUnits(models.CollectionsContentUnitWhere.CollectionID.EQ(c.ID)).Count(suite.tx)
	utils.Must(err)

	suite.EqualValues(MinCuNumberForNewLessonSeries, countCcu)
	_, err = ProcessCITMetadata(suite.tx, metadata, tf.Original, tf.Proxy, nil)
	suite.Require().Nil(err)

	_, err = makePublishedLast(suite.tx)
	suite.Require().Nil(err)
	countCcu, err = models.CollectionsContentUnits(models.CollectionsContentUnitWhere.CollectionID.EQ(c.ID)).Count(suite.tx)
	suite.Require().Nil(err)
	suite.EqualValues(MinCuNumberForNewLessonSeries+1, countCcu)
	countC, err := models.Collections(models.CollectionWhere.TypeID.EQ(common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSONS_SERIES].ID)).Count(suite.tx)
	suite.Require().Nil(err)

	for i := 0; i < MinCuNumberForNewLessonSeries; i++ {
		_, err = ProcessCITMetadata(suite.tx, metadata, tf.Original, tf.Proxy, nil)
		suite.Require().Nil(err)

		_, err = makePublishedLast(suite.tx)
		suite.Require().Nil(err)
		utils.Must(c.Reload(suite.tx))
		err = c.L.LoadCollectionsContentUnits(suite.tx, true, c, nil)
		suite.Require().Nil(err)
		_countC, err := models.Collections(
			models.CollectionWhere.TypeID.EQ(common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSONS_SERIES].ID),
		).Count(suite.tx)
		suite.Require().Nil(err)
		if i < MinCuNumberForNewLessonSeries-1 {
			suite.EqualValues(int(countCcu)+i+1, len(c.R.CollectionsContentUnits), fmt.Sprintf("iteration no new series %d", i))
			suite.EqualValues(countC, _countC, fmt.Sprintf("iteration %d", i))
		} else {
			suite.EqualValues(int(countC)+len(lUids)-1, _countC, fmt.Sprintf("iteration %d", i))
		}
	}
}

func createCUWithSourceForLessonsSeries(tx boil.Transactor, s *models.Source, c *models.Collection, i int) *models.ContentUnit {
	fd := time.Now().Add(-time.Hour * 24 * time.Duration(i))
	props := map[string]interface{}{"film_date": fd.Format("2006-01-02")}
	cu, err := CreateContentUnit(tx, common.CT_LESSON_PART, props)
	utils.Must(err)
	cu.Published = true
	_, err = cu.Update(tx, boil.Infer())
	utils.Must(err)
	utils.Must(cu.AddSources(tx, false, s))
	if c != nil {
		utils.Must(c.AddCollectionsContentUnits(tx, true, &models.CollectionsContentUnit{ContentUnitID: cu.ID, Position: i + 1}))
	}
	return cu
}

func createCUWithLikutForLessonsSeries(tx boil.Transactor, l *models.ContentUnit, c *models.Collection, i int) *models.ContentUnit {
	fd := time.Now().Add(-time.Hour * 24 * time.Duration(i))
	props := map[string]interface{}{"film_date": fd.Format("2006-01-02")}
	cu, err := CreateContentUnit(tx, common.CT_LESSON_PART, props)
	utils.Must(err)
	cu.Published = true
	_, err = cu.Update(tx, boil.Infer())
	utils.Must(err)
	utils.Must(cu.AddSourceContentUnitDerivations(tx, true, &models.ContentUnitDerivation{
		SourceID:  cu.ID,
		DerivedID: l.ID,
		Name:      "name",
	}))
	utils.Must(c.AddCollectionsContentUnits(tx, true, &models.CollectionsContentUnit{ContentUnitID: cu.ID, Position: i}))
	return cu
}

func makePublishedLast(tx boil.Transactor) (*models.ContentUnit, error) {
	prev, err := models.ContentUnits(
		models.ContentUnitWhere.TypeID.EQ(common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSON_PART].ID),
		qm.OrderBy("id DESC"),
	).One(tx)
	if err != nil {
		return nil, err
	}
	prev.Published = true
	_, err = prev.Update(tx, boil.Infer())

	if err != nil {
		return nil, err
	}
	return prev, nil
}

// Helpers

func createDummySources(exec boil.Executor, author *models.Author) []string {
	uids := make([]string, rand.Intn(5)+5)
	s, err := models.Sources(qm.OrderBy("id DESC")).One(exec)
	utils.Must(err)
	if author == nil {
		author, err = models.Authors().One(exec)
		if err == sql.ErrNoRows {
			author = &models.Author{Code: "tt", Name: "test", FullName: null.StringFrom("test author")}
			utils.Must(author.Insert(exec, boil.Infer()))
		} else if err != nil {
			utils.Must(err)
		}
	}
	for i := range uids {
		s := &models.Source{
			ID:  s.ID + int64(i) + 1,
			UID: utils.GenerateUID(8), TypeID: 1,
			Name: fmt.Sprintf("Dummy source %d", i),
		}
		utils.Must(s.Insert(exec, boil.Infer()))
		utils.Must(s.AddAuthors(exec, false, author))
		uids[i] = s.UID
	}
	return uids
}

func createDummyLikutim(exec boil.Executor) ([]*models.ContentUnit, error) {
	likutim := make([]*models.ContentUnit, rand.Intn(5)+2)
	for i, _ := range likutim {
		likutim[i] = &models.ContentUnit{
			UID:       utils.GenerateUID(8),
			TypeID:    common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LIKUTIM].ID,
			Published: true,
		}
		utils.Must(likutim[i].Insert(exec, boil.Infer()))

		i18ns := []*models.ContentUnitI18n{{Language: common.LANG_HEBREW, Name: null.StringFrom("name")}}
		utils.Must(likutim[i].AddContentUnitI18ns(exec, false, i18ns...))
	}
	return likutim, nil
}

type TrimFiles struct {
	Original *models.File
	Proxy    *models.File
}

func (suite *MetadataProcessorSuite) simulateSimpleChain() (TrimFiles, string) {
	CS_SHA1 := utils.RandomSHA1()
	DMX_O_SHA1 := utils.RandomSHA1()
	DMX_P_SHA1 := utils.RandomSHA1()
	TRM_O_SHA1 := utils.RandomSHA1()
	TRM_P_SHA1 := utils.RandomSHA1()
	WorkflowID := "c12356788"

	// capture_start
	_, evnts, err := handleCaptureStart(suite.tx, CaptureStartRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: WorkflowID,
		},
		FileName:      "capture_start_simple",
		CaptureSource: "mltcap",
		CollectionUID: "c12356788",
	})
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	// capture_stop
	_, evnts, err = handleCaptureStop(suite.tx, CaptureStopRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: WorkflowID,
		},
		File: File{
			FileName:  "capture_stop_simple.mp4",
			Sha1:      CS_SHA1,
			Size:      98737,
			CreatedAt: &Timestamp{Time: time.Now()},
			Language:  common.LANG_HEBREW,
		},
		CaptureSource: "mltcap",
		CollectionUID: "c12356789",
		Part:          "false",
	})
	suite.Require().Nil(err)

	// demux
	_, evnts, err = handleDemux(suite.tx, DemuxRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		Sha1: CS_SHA1,
		Original: AVFile{
			File: File{
				FileName:  "demux_simple_original.mp4",
				Sha1:      DMX_O_SHA1,
				Size:      98737,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		Proxy: &AVFile{
			File: File{
				FileName:  "demux_simple_proxy.mp4",
				Sha1:      DMX_P_SHA1,
				Size:      9878,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 891.8800,
		},
		CaptureSource: "mltcap",
	})
	suite.Require().Nil(err)

	// trim
	op, evnts, err := handleTrim(suite.tx, TrimRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		OriginalSha1: DMX_O_SHA1,
		ProxySha1:    DMX_P_SHA1,
		Original: AVFile{
			File: File{
				FileName:  "trim_simple_original.mp4",
				Sha1:      TRM_O_SHA1,
				Size:      98000,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		Proxy: &AVFile{
			File: File{
				FileName:  "trim_simple_proxy.mp4",
				Sha1:      TRM_P_SHA1,
				Size:      9800,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 891.8800,
		},
		CaptureSource: "mltcap",
		In:            []float64{10.05, 249.43},
		Out:           []float64{240.51, 899.27},
	})
	suite.Require().Nil(err)

	files := suite.opFilesBySHA1(op)
	return TrimFiles{Original: files[TRM_O_SHA1], Proxy: files[TRM_P_SHA1]}, WorkflowID
}

func (suite *MetadataProcessorSuite) simulateLessonChain() map[string]TrimFiles {
	CS_SHA1 := [5]string{}
	DMX_O_SHA1 := [5]string{}
	DMX_P_SHA1 := [5]string{}
	TRM_O_SHA1 := [7]string{}
	TRM_P_SHA1 := [7]string{}
	for i := range CS_SHA1 {
		CS_SHA1[i] = utils.RandomSHA1()
		DMX_O_SHA1[i] = utils.RandomSHA1()
		DMX_P_SHA1[i] = utils.RandomSHA1()
		TRM_O_SHA1[i] = utils.RandomSHA1()
		TRM_P_SHA1[i] = utils.RandomSHA1()
	}
	TRM_O_SHA1[5] = utils.RandomSHA1()
	TRM_P_SHA1[5] = utils.RandomSHA1()
	TRM_O_SHA1[6] = utils.RandomSHA1()
	TRM_P_SHA1[6] = utils.RandomSHA1()

	// capture_start
	_, evnts, err := handleCaptureStart(suite.tx, CaptureStartRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: "c12356789",
		},
		FileName:      "capture_start_full",
		CaptureSource: "mltbackup",
		CollectionUID: "c12356789",
	})
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	for i := 0; i < 4; i++ {
		part := strconv.Itoa(i)
		_, evnts, err := handleCaptureStart(suite.tx, CaptureStartRequest{
			Operation: Operation{
				Station:    "Capture station",
				User:       "operator@dev.com",
				WorkflowID: "c" + strings.Repeat(part, 8),
			},
			FileName:      "capture_start_part" + part,
			CaptureSource: "mltcap",
			CollectionUID: "c12356789",
		})
		suite.Require().Nil(err)
		suite.Require().Nil(evnts)
	}

	// capture_stop
	_, evnts, err = handleCaptureStop(suite.tx, CaptureStopRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: "c12356789",
		},
		File: File{
			FileName:  "capture_stop_full.mp4",
			Sha1:      CS_SHA1[4],
			Size:      98737,
			CreatedAt: &Timestamp{Time: time.Now()},
			Language:  common.LANG_MULTI,
		},
		CaptureSource: "mltbackup",
		CollectionUID: "c12356789",
		Part:          "full",
	})
	suite.Require().Nil(err)

	for i := 0; i < 4; i++ {
		part := strconv.Itoa(i)
		_, evnts, err := handleCaptureStop(suite.tx, CaptureStopRequest{
			Operation: Operation{
				Station:    "Capture station",
				User:       "operator@dev.com",
				WorkflowID: "c" + strings.Repeat(part, 8),
			},
			File: File{
				FileName:  "capture_stop_part" + part + ".mp4",
				Sha1:      CS_SHA1[i],
				Size:      int64(i),
				CreatedAt: &Timestamp{Time: time.Now()},
				Language:  common.LANG_MULTI,
			},
			CaptureSource: "mltcap",
			CollectionUID: "c12356789",
			Part:          part,
		})
		suite.Require().Nil(err)
		suite.Require().Nil(evnts)
	}

	// demux
	_, evnts, err = handleDemux(suite.tx, DemuxRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		Sha1: CS_SHA1[4],
		Original: AVFile{
			File: File{
				FileName:  "demux_full_original.mp4",
				Sha1:      DMX_O_SHA1[4],
				Size:      98737,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		Proxy: &AVFile{
			File: File{
				FileName:  "demux_full_proxy.mp4",
				Sha1:      DMX_P_SHA1[4],
				Size:      9878,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 891.8800,
		},
		CaptureSource: "mltbackup",
	})
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	for i := 0; i < 4; i++ {
		part := strconv.Itoa(i)
		_, evnts, err := handleDemux(suite.tx, DemuxRequest{
			Operation: Operation{
				Station: "Trimmer station",
				User:    "operator@dev.com",
			},
			Sha1: CS_SHA1[i],
			Original: AVFile{
				File: File{
					FileName:  "demux_part" + part + "_original.mp4",
					Sha1:      DMX_O_SHA1[i],
					Size:      98737,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 892.1900,
			},
			Proxy: &AVFile{
				File: File{
					FileName:  "demux_part" + part + "_proxy.mp4",
					Sha1:      DMX_P_SHA1[i],
					Size:      9878,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 891.8800,
			},
			CaptureSource: "mltcap",
		})
		suite.Require().Nil(err)
		suite.Require().Nil(evnts)
	}

	trimFiles := make(map[string]TrimFiles)

	// trim
	op, evnts, err := handleTrim(suite.tx, TrimRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		OriginalSha1: DMX_O_SHA1[4],
		ProxySha1:    DMX_P_SHA1[4],
		Original: AVFile{
			File: File{
				FileName:  "trim_full_original.mp4",
				Sha1:      TRM_O_SHA1[4],
				Size:      98000,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		Proxy: &AVFile{
			File: File{
				FileName:  "trim_full_proxy.mp4",
				Sha1:      TRM_P_SHA1[4],
				Size:      9800,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 891.8800,
		},
		CaptureSource: "mltbackup",
		In:            []float64{10.05, 249.43},
		Out:           []float64{240.51, 899.27},
	})
	suite.Require().Nil(err)
	files := suite.opFilesBySHA1(op)
	trimFiles["full"] = TrimFiles{
		Original: files[TRM_O_SHA1[4]],
		Proxy:    files[TRM_P_SHA1[4]],
	}

	for i := 0; i < 4; i++ {
		part := strconv.Itoa(i)
		op, evnts, err := handleTrim(suite.tx, TrimRequest{
			Operation: Operation{
				Station: "Trimmer station",
				User:    "operator@dev.com",
			},
			OriginalSha1: DMX_O_SHA1[i],
			ProxySha1:    DMX_P_SHA1[i],
			Original: AVFile{
				File: File{
					FileName:  "trim_part" + part + "_original.mp4",
					Sha1:      TRM_O_SHA1[i],
					Size:      98000,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 892.1900,
			},
			Proxy: &AVFile{
				File: File{
					FileName:  "trim_part" + part + "_proxy.mp4",
					Sha1:      TRM_P_SHA1[i],
					Size:      9800,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 891.8800,
			},
			CaptureSource: "mltbackup",
			In:            []float64{10.05, 249.43},
			Out:           []float64{240.51, 899.27},
		})
		suite.Require().Nil(err)
		suite.Require().Nil(evnts)

		files := suite.opFilesBySHA1(op)
		trimFiles["part"+part] = TrimFiles{
			Original: files[TRM_O_SHA1[i]],
			Proxy:    files[TRM_P_SHA1[i]],
		}
	}

	// trim kitei makor from part1
	op, evnts, err = handleTrim(suite.tx, TrimRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		OriginalSha1: DMX_O_SHA1[1],
		ProxySha1:    DMX_P_SHA1[1],
		Original: AVFile{
			File: File{
				FileName:  "trim_part1_kitei_makor_original.mp4",
				Sha1:      TRM_O_SHA1[5],
				Size:      6700,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 92.1900,
		},
		Proxy: &AVFile{
			File: File{
				FileName:  "trim_part1_kitei_makor_proxy.mp4",
				Sha1:      TRM_P_SHA1[5],
				Size:      6700,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 91.8800,
		},
		CaptureSource: "mltcap",
		In:            []float64{10.05, 249.43, 253.83, 312.23, 463.3, 512.3},
		Out:           []float64{240.51, 250.31, 282.13, 441.03, 483.39, 899.81},
	})
	suite.Require().Nil(err)
	files = suite.opFilesBySHA1(op)
	trimFiles["part1_kitei-makor"] = TrimFiles{
		Original: files[TRM_O_SHA1[5]],
		Proxy:    files[TRM_P_SHA1[5]],
	}

	// trim ktaim nivcharim from full
	op, evnts, err = handleTrim(suite.tx, TrimRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		OriginalSha1: DMX_O_SHA1[4],
		ProxySha1:    DMX_P_SHA1[4],
		Original: AVFile{
			File: File{
				FileName:  "heb_o_rav_2019-12-23_ktaim-nivharim_n1_original.mp4",
				Sha1:      TRM_O_SHA1[6],
				Size:      6700,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 92.1900,
		},
		Proxy: &AVFile{
			File: File{
				FileName:  "heb_o_rav_2019-12-23_ktaim-nivharim_n1_proxy.mp4",
				Sha1:      TRM_P_SHA1[6],
				Size:      6700,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 91.8800,
		},
		CaptureSource: "mltcap",
		In:            []float64{10.05, 249.43, 253.83, 312.23, 463.3, 512.3},
		Out:           []float64{240.51, 250.31, 282.13, 441.03, 483.39, 899.81},
	})
	suite.Require().Nil(err)
	files = suite.opFilesBySHA1(op)
	trimFiles["ktaim-nivcharim"] = TrimFiles{
		Original: files[TRM_O_SHA1[6]],
		Proxy:    files[TRM_P_SHA1[6]],
	}

	return trimFiles
}

func (suite *MetadataProcessorSuite) simulateAdditionalCapture() *models.File {
	CsSha1 := utils.RandomSHA1()
	TrmSSha1 := utils.RandomSHA1()

	// capture_start
	_, evnts, err := handleCaptureStart(suite.tx, CaptureStartRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: "c987654321",
		},
		FileName:      "capture_start_source",
		CaptureSource: "archcap",
		CollectionUID: "c987654321",
	})
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	// capture_stop
	_, evnts, err = handleCaptureStop(suite.tx, CaptureStopRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: "c987654321",
		},
		File: File{
			FileName:  "capture_stop_source.mp4",
			Sha1:      CsSha1,
			Size:      98737,
			CreatedAt: &Timestamp{Time: time.Now()},
			Language:  common.LANG_MULTI,
		},
		CaptureSource: "archcap",
		CollectionUID: "c987654321",
		Part:          "source",
	})
	suite.Require().Nil(err)

	// trim
	op, evnts, err := handleTrim(suite.tx, TrimRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		OriginalSha1: CsSha1,
		ProxySha1:    "",
		Original: AVFile{
			File: File{
				FileName:  "trim_full_original.mp4",
				Sha1:      TrmSSha1,
				Size:      98000,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		Proxy:         nil,
		CaptureSource: "archcap",
		In:            []float64{10.05, 249.43},
		Out:           []float64{240.51, 899.27},
	})
	suite.Require().Nil(err)
	files := suite.opFilesBySHA1(op)
	return files[TrmSSha1]
}

func (suite *MetadataProcessorSuite) simulateSpecialLessonChain() map[string]TrimFiles {
	CS_SHA1 := [2]string{}
	DMX_O_SHA1 := [2]string{}
	DMX_P_SHA1 := [2]string{}
	TRM_O_SHA1 := [13]string{}
	TRM_P_SHA1 := [13]string{}

	CS_SHA1[0] = utils.RandomSHA1()
	CS_SHA1[1] = utils.RandomSHA1()
	DMX_O_SHA1[0] = utils.RandomSHA1()
	DMX_O_SHA1[1] = utils.RandomSHA1()
	DMX_P_SHA1[0] = utils.RandomSHA1()
	DMX_P_SHA1[1] = utils.RandomSHA1()
	for i := range TRM_O_SHA1 {
		TRM_O_SHA1[i] = utils.RandomSHA1()
		TRM_P_SHA1[i] = utils.RandomSHA1()
	}

	// capture_start
	_, evnts, err := handleCaptureStart(suite.tx, CaptureStartRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: "c12356789",
		},
		FileName:      "capture_start_full",
		CaptureSource: "mltbackup",
		CollectionUID: "c12356789",
	})
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	part := strconv.Itoa(0)
	_, evnts, err = handleCaptureStart(suite.tx, CaptureStartRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: "c" + strings.Repeat(part, 8),
		},
		FileName:      "capture_start_part" + part,
		CaptureSource: "mltcap",
		CollectionUID: "c12356789",
	})
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	// capture_stop
	_, evnts, err = handleCaptureStop(suite.tx, CaptureStopRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: "c12356789",
		},
		File: File{
			FileName:  "capture_stop_full.mp4",
			Sha1:      CS_SHA1[0],
			Size:      98737,
			CreatedAt: &Timestamp{Time: time.Now()},
			Language:  common.LANG_MULTI,
		},
		CaptureSource: "mltbackup",
		CollectionUID: "c12356789",
		Part:          "full",
	})
	suite.Require().Nil(err)

	_, evnts, err = handleCaptureStop(suite.tx, CaptureStopRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: "c" + strings.Repeat(part, 8),
		},
		File: File{
			FileName:  "capture_stop_part" + part + ".mp4",
			Sha1:      CS_SHA1[1],
			Size:      int64(1),
			CreatedAt: &Timestamp{Time: time.Now()},
			Language:  common.LANG_MULTI,
		},
		CaptureSource: "mltcap",
		CollectionUID: "c12356789",
		Part:          part,
	})
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	// demux
	_, evnts, err = handleDemux(suite.tx, DemuxRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		Sha1: CS_SHA1[0],
		Original: AVFile{
			File: File{
				FileName:  "demux_full_original.mp4",
				Sha1:      DMX_O_SHA1[0],
				Size:      98737,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		Proxy: &AVFile{
			File: File{
				FileName:  "demux_full_proxy.mp4",
				Sha1:      DMX_P_SHA1[0],
				Size:      9878,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 891.8800,
		},
		CaptureSource: "mltbackup",
	})
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	_, evnts, err = handleDemux(suite.tx, DemuxRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		Sha1: CS_SHA1[1],
		Original: AVFile{
			File: File{
				FileName:  "demux_part" + part + "_original.mp4",
				Sha1:      DMX_O_SHA1[1],
				Size:      98737,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		Proxy: &AVFile{
			File: File{
				FileName:  "demux_part" + part + "_proxy.mp4",
				Sha1:      DMX_P_SHA1[1],
				Size:      9878,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 891.8800,
		},
		CaptureSource: "mltcap",
	})
	suite.Require().Nil(err)
	suite.Require().Nil(evnts)

	trimFiles := make(map[string]TrimFiles)

	// trim
	op, evnts, err := handleTrim(suite.tx, TrimRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		OriginalSha1: DMX_O_SHA1[0],
		ProxySha1:    DMX_P_SHA1[0],
		Original: AVFile{
			File: File{
				FileName:  "trim_full_original.mp4",
				Sha1:      TRM_O_SHA1[4],
				Size:      98000,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		Proxy: &AVFile{
			File: File{
				FileName:  "trim_full_proxy.mp4",
				Sha1:      TRM_P_SHA1[4],
				Size:      9800,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 891.8800,
		},
		CaptureSource: "mltbackup",
		In:            []float64{10.05, 249.43},
		Out:           []float64{240.51, 899.27},
	})
	suite.Require().Nil(err)
	files := suite.opFilesBySHA1(op)
	trimFiles["full"] = TrimFiles{
		Original: files[TRM_O_SHA1[4]],
		Proxy:    files[TRM_P_SHA1[4]],
	}

	for i := 0; i < 4; i++ {
		part := strconv.Itoa(i)
		op, evnts, err := handleTrim(suite.tx, TrimRequest{
			Operation: Operation{
				Station: "Trimmer station",
				User:    "operator@dev.com",
			},
			OriginalSha1: DMX_O_SHA1[1],
			ProxySha1:    DMX_P_SHA1[1],
			Original: AVFile{
				File: File{
					FileName:  "trim_part" + part + "_original.mp4",
					Sha1:      TRM_O_SHA1[i],
					Size:      98000,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 892.1900,
			},
			Proxy: &AVFile{
				File: File{
					FileName:  "trim_part" + part + "_proxy.mp4",
					Sha1:      TRM_P_SHA1[i],
					Size:      9800,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 891.8800,
			},
			CaptureSource: "mltbackup",
			In:            []float64{10.05, 249.43},
			Out:           []float64{240.51, 899.27},
		})
		suite.Require().Nil(err)
		suite.Require().Nil(evnts)

		files := suite.opFilesBySHA1(op)
		trimFiles["part"+part] = TrimFiles{
			Original: files[TRM_O_SHA1[i]],
			Proxy:    files[TRM_P_SHA1[i]],
		}
	}

	// trim kitei makor from all parts
	for i := 0; i < 4; i++ {
		part := strconv.Itoa(i)
		op, evnts, err = handleTrim(suite.tx, TrimRequest{
			Operation: Operation{
				Station: "Trimmer station",
				User:    "operator@dev.com",
			},
			OriginalSha1: DMX_O_SHA1[1],
			ProxySha1:    DMX_P_SHA1[1],
			Original: AVFile{
				File: File{
					FileName:  fmt.Sprintf("trim_part%d_kitei_makor_original.mp4", i),
					Sha1:      TRM_O_SHA1[5+i],
					Size:      6700,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 92.1900,
			},
			Proxy: &AVFile{
				File: File{
					FileName:  fmt.Sprintf("trim_part%d_kitei_makor_proxy.mp4", i),
					Sha1:      TRM_P_SHA1[5+i],
					Size:      6700,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 91.8800,
			},
			CaptureSource: "mltcap",
			In:            []float64{10.05, 249.43, 253.83, 312.23, 463.3, 512.3},
			Out:           []float64{240.51, 250.31, 282.13, 441.03, 483.39, 899.81},
		})
		suite.Require().Nil(err)
		files = suite.opFilesBySHA1(op)
		trimFiles["part"+part+"_kitei-makor"] = TrimFiles{
			Original: files[TRM_O_SHA1[5+i]],
			Proxy:    files[TRM_P_SHA1[5+i]],
		}
	}

	// trim lelo mikud from all parts
	for i := 0; i < 4; i++ {
		part := strconv.Itoa(i)
		op, evnts, err = handleTrim(suite.tx, TrimRequest{
			Operation: Operation{
				Station: "Trimmer station",
				User:    "operator@dev.com",
			},
			OriginalSha1: DMX_O_SHA1[1],
			ProxySha1:    DMX_P_SHA1[1],
			Original: AVFile{
				File: File{
					FileName:  fmt.Sprintf("trim_part%d_lelo_mikud_original.mp4", i),
					Sha1:      TRM_O_SHA1[9+i],
					Size:      6700,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 92.1900,
			},
			Proxy: &AVFile{
				File: File{
					FileName:  fmt.Sprintf("trim_part%d_lelo_mikud_proxy.mp4", i),
					Sha1:      TRM_P_SHA1[9+i],
					Size:      6700,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 91.8800,
			},
			CaptureSource: "mltcap",
			In:            []float64{10.05, 249.43, 253.83, 312.23, 463.3, 512.3},
			Out:           []float64{240.51, 250.31, 282.13, 441.03, 483.39, 899.81},
		})
		suite.Require().Nil(err)
		files = suite.opFilesBySHA1(op)
		trimFiles["part"+part+"_lelo-mikud"] = TrimFiles{
			Original: files[TRM_O_SHA1[9+i]],
			Proxy:    files[TRM_P_SHA1[9+i]],
		}
	}

	return trimFiles
}

func (suite *MetadataProcessorSuite) simulateConvertUpload(original *models.File) map[string]*models.File {
	files := make(map[string]*models.File)

	originalSha1 := hex.EncodeToString(original.Sha1.Bytes)
	for _, lang := range common.ALL_LANGS {
		if lang == common.LANG_UNKNOWN || lang == common.LANG_MULTI {
			continue
		}

		input := ConvertRequest{
			Operation: Operation{
				Station: "Trimmer station",
				User:    "operator@dev.com",
			},
			Sha1: originalSha1,
			Output: []HLSFile{
				{
					AVFile: AVFile{
						File: File{
							FileName:  fmt.Sprintf("%s_test_file.mp4", lang),
							Sha1:      utils.RandomSHA1(),
							Size:      694,
							CreatedAt: &Timestamp{Time: time.Now()},
							Type:      "video",
							MimeType:  "video/mp4",
							Language:  lang,
						},
						Duration: 871,
					},
				},

				{
					AVFile: AVFile{
						File: File{
							FileName:  fmt.Sprintf("%s_test_file.mp3", lang),
							Sha1:      utils.RandomSHA1(),
							Size:      694,
							CreatedAt: &Timestamp{Time: time.Now()},
							Type:      "audio",
							MimeType:  "audio/mpeg",
							Language:  lang,
						},
						Duration: 871,
					},
				},
			},
		}

		op, _, err := handleConvert(suite.tx, input)
		suite.Require().Nil(err)
		err = op.L.LoadFiles(suite.tx, true, op, nil)
		suite.Require().Nil(err)

		for _, f := range op.R.Files {
			// This is the trimmed file, not converted...
			if f.ID == original.ID {
				continue
			}

			sha1Str := hex.EncodeToString(f.Sha1.Bytes)
			files[sha1Str] = f

			// upload
			input := UploadRequest{
				Operation: Operation{
					Station: "Upload station",
					User:    "operator@dev.com",
				},
				AVFile: AVFile{
					File: File{
						FileName:  f.Name,
						Sha1:      sha1Str,
						Size:      f.Size,
						CreatedAt: &Timestamp{f.CreatedAt},
					},
					Duration: 871,
				},
				Url: "http://example.com/some/url/to/file.mp4",
			}

			_, _, err = handleUpload(suite.tx, input)
			suite.Require().Nil(err)
		}
	}

	return files
}

func (suite *MetadataProcessorSuite) simulateLessonChainWithSource() map[string]TrimFiles {
	trimFiles := suite.simulateLessonChain()

	CS_SHA1 := [4]string{}
	DMX_O_SHA1 := [4]string{}
	TRM_O_SHA1 := [4]string{}
	for i := range CS_SHA1 {
		CS_SHA1[i] = utils.RandomSHA1()
		DMX_O_SHA1[i] = utils.RandomSHA1()
		TRM_O_SHA1[i] = utils.RandomSHA1()
	}

	// capture_start
	for i := 0; i < 4; i++ {
		part := strconv.Itoa(i)
		_, evnts, err := handleCaptureStart(suite.tx, CaptureStartRequest{
			Operation: Operation{
				Station:    "Capture station",
				User:       "operator@dev.com",
				WorkflowID: "c" + strings.Repeat(part, 8),
			},
			FileName:      "capture_start_source_part" + part,
			CaptureSource: "capture_of_source",
			CollectionUID: "c12356789",
		})
		suite.Require().Nil(err)
		suite.Require().Nil(evnts)
	}

	// capture_stop
	for i := 0; i < 4; i++ {
		part := strconv.Itoa(i)
		_, evnts, err := handleCaptureStop(suite.tx, CaptureStopRequest{
			Operation: Operation{
				Station:    "Capture station",
				User:       "operator@dev.com",
				WorkflowID: "c" + strings.Repeat(part, 8),
			},
			File: File{
				FileName:  "capture_stop_source_part" + part + ".mp4",
				Sha1:      CS_SHA1[i],
				Size:      int64(i),
				CreatedAt: &Timestamp{Time: time.Now()},
				Language:  common.LANG_MULTI,
			},
			CaptureSource: "capture_of_source",
			CollectionUID: "c12356789",
			Part:          part,
		})
		suite.Require().Nil(err)
		suite.Require().Nil(evnts)
	}

	// demux
	for i := 0; i < 4; i++ {
		part := strconv.Itoa(i)
		_, evnts, err := handleDemux(suite.tx, DemuxRequest{
			Operation: Operation{
				Station: "Trimmer station",
				User:    "operator@dev.com",
			},
			Sha1: CS_SHA1[i],
			Original: AVFile{
				File: File{
					FileName:  "demux_source_part" + part + "_original.mp4",
					Sha1:      DMX_O_SHA1[i],
					Size:      98737,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 892.1900,
			},
			CaptureSource: "capture_of_source",
		})
		suite.Require().Nil(err)
		suite.Require().Nil(evnts)
	}

	// trim
	for i := 0; i < 4; i++ {
		part := strconv.Itoa(i)
		op, evnts, err := handleTrim(suite.tx, TrimRequest{
			Operation: Operation{
				Station: "Trimmer station",
				User:    "operator@dev.com",
			},
			OriginalSha1: DMX_O_SHA1[i],
			Original: AVFile{
				File: File{
					FileName:  "trim_source_part" + part + "_original.mp4",
					Sha1:      TRM_O_SHA1[i],
					Size:      98000,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 892.1900,
			},
			CaptureSource: "capture_of_source",
			In:            []float64{10.05, 249.43},
			Out:           []float64{240.51, 899.27},
		})
		suite.Require().Nil(err)
		suite.Require().Nil(evnts)

		files := suite.opFilesBySHA1(op)
		trimFiles["source_part"+part] = TrimFiles{
			Original: files[TRM_O_SHA1[i]],
		}
	}

	return trimFiles
}

func (suite *MetadataProcessorSuite) opFilesBySHA1(o *models.Operation) map[string]*models.File {
	files := make(map[string]*models.File)
	for _, f := range o.R.Files {
		files[hex.EncodeToString(f.Sha1.Bytes)] = f
	}
	return files
}

func (suite *MetadataProcessorSuite) someSources() []string {
	items, err := models.Sources(qm.Limit(1 + rand.Intn(10))).All(suite.tx)
	suite.Require().Nil(err)
	uids := make([]string, len(items))
	for i, x := range items {
		uids[i] = x.UID
	}
	return uids
}

func (suite *MetadataProcessorSuite) someTags() []string {
	items, err := models.Tags(qm.Limit(1 + rand.Intn(10))).All(suite.tx)
	suite.Require().Nil(err)
	uids := make([]string, len(items))
	for i, x := range items {
		uids[i] = x.UID
	}
	return uids
}

func (suite *MetadataProcessorSuite) assertFiles(metadata CITMetadata, original, proxy *models.File) {
	capDate := metadata.CaptureDate
	filmDate := metadata.CaptureDate
	//if metadata.WeekDate != nil {
	//	filmDate = *metadata.WeekDate
	//}
	var lang string
	if metadata.HasTranslation {
		lang = common.LANG_MULTI
	} else {
		lang = common.StdLang(metadata.Language)
	}

	var props map[string]interface{}

	// original properties
	suite.Require().True(original.Properties.Valid)
	err := original.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(capDate.Format("2006-01-02"), props["capture_date"], "original.Properties[\"capture_date\"]")
	suite.Equal(filmDate.Format("2006-01-02"), props["film_date"], "original.Properties[\"film_date\"]")

	// proxy properties
	suite.Require().True(proxy.Properties.Valid)
	err = proxy.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(capDate.Format("2006-01-02"), props["capture_date"], "proxy.Properties[\"capture_date\"]")
	suite.Equal(filmDate.Format("2006-01-02"), props["film_date"], "proxy.Properties[\"film_date\"]")

	// original language
	suite.True(original.Language.Valid, "original.Language.Valid")
	suite.Equal(lang, original.Language.String, "original.Language")
}

func (suite *MetadataProcessorSuite) assertContentUnit(metadata CITMetadata, original, proxy *models.File, isUpdate bool) {
	// reload unit
	err := original.L.LoadContentUnit(suite.tx, true, original, nil)
	suite.Require().Nil(err)
	cu := original.R.ContentUnit

	isDerived := metadata.ArtifactType.Valid && metadata.ArtifactType.String != "main"
	ct := metadata.ContentType
	if isDerived {
		ct = strings.ToUpper(metadata.ArtifactType.String)
	}

	// properties
	suite.Equal(cu.TypeID, common.CONTENT_TYPE_REGISTRY.ByName[ct].ID, "cu.type_id")

	capDate := metadata.CaptureDate
	filmDate := metadata.CaptureDate
	if metadata.WeekDate != nil {
		filmDate = *metadata.WeekDate
	}
	suite.Require().True(original.Properties.Valid)
	var originalProps map[string]interface{}
	err = original.Properties.Unmarshal(&originalProps)
	suite.Require().Nil(err)

	suite.Require().True(cu.Properties.Valid)
	var props map[string]interface{}
	err = cu.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(capDate.Format("2006-01-02"), props["capture_date"], "cu.Properties[\"capture_date\"]")
	suite.Equal(filmDate.Format("2006-01-02"), props["film_date"], "cu.Properties[\"film_date\"]")
	suite.Equal(common.StdLang(metadata.Language), props["original_language"], "cu.Properties[\"original_language\"]")
	suite.EqualValues(int(originalProps["duration"].(float64)), props["duration"], "cu.Properties[\"duration\"]")

	// files in unit
	err = cu.L.LoadFiles(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	filesInUnit := []*models.File{original, proxy}

	if !isUpdate {
		if !isDerived {
			// ancestors
			ancestors, err := FindFileAncestors(suite.tx, original.ID)
			suite.Require().Nil(err)
			filesInUnit = append(filesInUnit, ancestors...)
			proxy.L.LoadParent(suite.tx, true, proxy, nil)
			suite.Require().Nil(err)
			filesInUnit = append(filesInUnit, proxy.R.Parent)
		}

		suite.Equal(len(filesInUnit), len(cu.R.Files), "len(cu.R.Files)")
		for i, f := range filesInUnit {
			suite.True(f.ContentUnitID.Valid, "Ancestor[%d].ContentUnitID.Valid", i)
			suite.Equal(original.ContentUnitID.Int64, f.ContentUnitID.Int64, "Ancestor[%d]ContentUnitID.Int64", i)
		}
	}

	// sources
	err = cu.L.LoadSources(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Equal(len(metadata.Sources), len(cu.R.Sources), "len(cu.R.Sources)")
	for _, x := range metadata.Sources {
		missing := true
		for _, y := range cu.R.Sources {
			if x == y.UID {
				missing = false
				break
			}
		}
		suite.False(missing, "Missing source %s", x)
	}

	// tags
	err = cu.L.LoadTags(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	suite.Equal(len(metadata.Tags), len(cu.R.Tags), "len(cu.R.Tags)")
	for _, x := range metadata.Tags {
		missing := true
		for _, y := range cu.R.Tags {
			if x == y.UID {
				missing = false
				break
			}
		}
		suite.False(missing, "Missing tag %s", x)
	}

	// likutim
	likutim, err := models.ContentUnits(
		qm.InnerJoin("content_unit_derivations cud ON cud.derived_id = \"content_units\".id"),
		qm.Where("cud.source_id = ? AND \"content_units\".type_id = ? AND published IS TRUE",
			cu.ID, common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LIKUTIM].ID)).
		All(suite.tx)
	suite.Require().Nil(err)
	suite.Equal(len(metadata.Likutim), len(likutim), "len(likutim)")
	for _, x := range metadata.Likutim {
		missing := true
		for _, y := range likutim {
			if x == y.UID {
				missing = false
				break
			}
		}
		suite.False(missing, "Missing Likutim %s", x)
	}

	// persons
	err = cu.L.LoadContentUnitsPersons(suite.tx, true, cu, nil)
	suite.Require().Nil(err)
	if metadata.Lecturer == "rav" {
		suite.Require().Len(cu.R.ContentUnitsPersons, 1, "cu.R.ContentUnitsPersons Length")
		cup := cu.R.ContentUnitsPersons[0]
		suite.Equal(common.PERSON_REGISTRY.ByPattern[common.P_RAV].ID, cup.PersonID, "cup.PersonID")
		suite.Equal(common.CONTENT_ROLE_TYPE_REGISTRY.ByName[common.CR_LECTURER].ID, cup.RoleID, "cup.PersonID")
	} else {
		suite.Empty(cu.R.ContentUnitsPersons, "Empty cu.R.ContentUnitsPersons")
	}
}

func SomeLikutim(exec boil.Executor) ([]*models.ContentUnit, error) {
	likutim := make([]*models.ContentUnit, 1+rand.Intn(10))
	for i, _ := range likutim {
		l := &models.ContentUnit{
			UID:       utils.GenerateUID(8),
			TypeID:    common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LIKUTIM].ID,
			Secure:    0,
			Published: true,
		}
		if err := l.Insert(exec, boil.Infer()); err != nil {
			return nil, err
		}
		i18ns := []*models.ContentUnitI18n{{Language: common.LANG_HEBREW, Name: null.StringFrom("name")}}
		if err := l.AddContentUnitI18ns(exec, true, i18ns...); err != nil {
			return nil, err
		}
		likutim[i] = l
	}
	return likutim, nil
}
