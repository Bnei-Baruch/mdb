package api

import (
	"database/sql"
	"encoding/hex"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

type RepoSuite struct {
	suite.Suite
	utils.TestDBManager
	tx *sql.Tx
}

func (suite *RepoSuite) SetupSuite() {
	suite.Require().Nil(suite.InitTestDB())
	suite.Require().Nil(common.InitTypeRegistries(suite.DB))
}

func (suite *RepoSuite) TearDownSuite() {
	suite.Require().Nil(suite.DestroyTestDB())
}

func (suite *RepoSuite) SetupTest() {
	var err error
	suite.tx, err = suite.DB.Begin()
	suite.Require().Nil(err)
}

func (suite *RepoSuite) TearDownTest() {
	err := suite.tx.Rollback()
	suite.Require().Nil(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRepo(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}

func (suite *RepoSuite) TestCreateOperation() {
	// test minimal input
	o := Operation{
		Station: "station",
		User:    "operator@dev.com",
	}
	op, err := CreateOperation(suite.tx, common.OP_CAPTURE_START, o, nil)
	suite.Require().Nil(err)
	suite.Require().Nil(op.Reload(suite.tx))

	suite.Equal(common.OPERATION_TYPE_REGISTRY.ByName[common.OP_CAPTURE_START].ID, op.TypeID, "TypeID")
	suite.Regexp(utils.UID_REGEX, op.UID, "UID")
	suite.True(op.Station.Valid, "Station.Valid")
	suite.Equal(o.Station, op.Station.String, "Station.String")
	user, err := models.Users(qm.Where("email=?", o.User)).One(suite.tx)
	suite.Equal(user.ID, op.UserID.Int64, "User")

	// test with unknown user
	o.User = "unknown@example.com"
	op, err = CreateOperation(suite.tx, common.OP_CAPTURE_START, o, nil)
	suite.Require().Nil(err)
	suite.Require().Nil(op.Reload(suite.tx))
	suite.False(op.UserID.Valid)
	o.User = "operator@dev.com"

	// test with workflow_id
	o.WorkflowID = "workflow_id"
	op, err = CreateOperation(suite.tx, common.OP_CAPTURE_START, o, nil)
	suite.Require().Nil(err)
	suite.Require().Nil(op.Reload(suite.tx))

	suite.True(op.Properties.Valid)
	var props map[string]interface{}
	err = op.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(o.WorkflowID, props["workflow_id"], "props: workflow_id")

	// test with custom props
	customProps := suite.getCustomProps()
	op, err = CreateOperation(suite.tx, common.OP_CAPTURE_START, o, customProps)
	suite.Require().Nil(err)
	suite.Require().Nil(op.Reload(suite.tx))

	err = op.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.assertCustomProps(customProps, props)
}

func (suite *RepoSuite) TestCreateFile() {
	// test minimal input
	f := File{
		FileName:  "file_name",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef0123456789",
		Size:      math.MaxInt64,
	}
	file, err := CreateFile(suite.tx, nil, f, nil)
	suite.Require().Nil(err)
	suite.Require().Nil(file.Reload(suite.tx))

	suite.Regexp(utils.UID_REGEX, file.UID, "UID")
	suite.Equal(f.FileName, file.Name, "Name")
	suite.True(file.FileCreatedAt.Valid, "FileCreatedAt.Valid")
	suite.Equal(f.CreatedAt.Unix(), file.FileCreatedAt.Time.Unix(), "FileCreatedAt.Time")
	suite.True(file.Sha1.Valid, "Sha1.Valid")
	suite.Equal(f.Sha1, hex.EncodeToString(file.Sha1.Bytes), "Sha1.Bytes")
	suite.Equal(f.Size, file.Size, "Size")
	suite.Empty(file.Type, "Type")
	suite.Empty(file.SubType, "SubType")
	suite.False(file.MimeType.Valid, "MimeType.Valid")
	suite.False(file.Language.Valid, "Language.Valid")
	suite.False(file.ParentID.Valid, "ParentID.Valid")
	suite.False(file.Published, "Published")
	suite.Equal(common.SEC_PUBLIC, file.Secure, "Secure")
	suite.False(file.RemovedAt.Valid, "RemovedAt")

	// test with optional attributes
	f2 := File{
		FileName:  "file_name",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef9876543210",
		Size:      1,
		Type:      "type",
		SubType:   "subtype",
		MimeType:  "mimetype",
		Language:  common.LANG_RUSSIAN,
	}
	cu, err := CreateContentUnit(suite.tx, common.CT_LESSON_PART, nil)
	suite.Require().Nil(err)
	err = file.SetContentUnit(suite.tx, false, cu)
	suite.Require().Nil(err)
	file2, err := CreateFile(suite.tx, file, f2, nil)
	suite.Require().Nil(err)
	suite.Require().Nil(file2.Reload(suite.tx))

	suite.Equal(f2.Sha1, hex.EncodeToString(file2.Sha1.Bytes), "file2.Sha1.Bytes")
	suite.Equal(f2.Size, file2.Size, "file2.Size")
	suite.Equal(f2.Type, file2.Type, "file2.Type")
	suite.Equal(f2.SubType, file2.SubType, "file2.SubType")
	suite.True(file2.MimeType.Valid, "file2.MimeType.Valid")
	suite.Equal(f2.MimeType, file2.MimeType.String, "file2.MimeType.String")
	suite.True(file2.Language.Valid, "file2.Language.Valid")
	suite.Equal(f2.Language, file2.Language.String, "file2.Language.String")
	suite.True(file2.ParentID.Valid, "file2.ParentID.Valid")
	suite.Equal(file.ID, file2.ParentID.Int64, "file2.ParentID.Int64")
	suite.True(file2.ContentUnitID.Valid, "file2.ContentUnitID.Valid")
	suite.Equal(file.ContentUnitID.Int64, file2.ContentUnitID.Int64, "file2.ContentUnitID.Int64")

	// test with custom props
	f3 := File{
		FileName:  "file_name",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "abcdef012356789abcdef0123456789abcdef012",
		Size:      1,
	}
	customProps := suite.getCustomProps()
	file3, err := CreateFile(suite.tx, nil, f3, customProps)
	suite.Require().Nil(err)
	suite.Require().Nil(file3.Reload(suite.tx))

	var props map[string]interface{}
	err = file3.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.assertCustomProps(customProps, props)

	// test invalid languages
	f4 := File{Sha1: "abcdef012356789abcdef0123456789abcdef012"}
	for _, x := range []string{"a", "aa", "aaa", "aaaa"} {
		f4.Language = x
		_, err = CreateFile(suite.tx, nil, f4, nil)
		suite.Require().EqualError(err, "Make file: Unknown language "+f4.Language, "Invalid language "+f4.Language)
	}

	// test mime type complements type & subtype
	f5 := File{
		FileName:  "file_name",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef1111111111",
		Size:      math.MaxInt64,
		MimeType:  common.ALL_MEDIA_TYPES[0].MimeType,
	}
	file5, err := CreateFile(suite.tx, nil, f5, nil)
	suite.Require().Nil(err)
	suite.Require().Nil(file.Reload(suite.tx))
	suite.Equal(common.ALL_MEDIA_TYPES[0].Type, file5.Type, "file5.Type")
	suite.Equal(common.ALL_MEDIA_TYPES[0].SubType, file5.SubType, "file5.SubType")
	suite.True(file5.MimeType.Valid, "file5.MimeType.Valid")
	suite.Equal(common.ALL_MEDIA_TYPES[0].MimeType, file5.MimeType.String, "file5.MimeType.String")
}

func (suite *RepoSuite) TestPublishFile() {
	f := File{
		FileName:  "file_name",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef0123456789",
		Size:      math.MaxInt64,
	}
	file, err := CreateFile(suite.tx, nil, f, nil)
	suite.Require().Nil(err)
	cu, err := CreateContentUnit(suite.tx, common.CT_LESSON_PART, nil)
	suite.Require().Nil(err)
	err = file.SetContentUnit(suite.tx, false, cu)
	suite.Require().Nil(err)
	c, err := CreateCollection(suite.tx, common.CT_DAILY_LESSON, nil)
	suite.Require().Nil(err)
	err = c.AddCollectionsContentUnits(suite.tx, true, &models.CollectionsContentUnit{ContentUnitID: cu.ID})
	suite.Require().Nil(err)

	_, err = PublishFile(suite.tx, file)
	suite.Require().Nil(err)
	err = file.Reload(suite.tx)
	suite.Require().Nil(err)
	suite.True(file.Published)
	err = cu.Reload(suite.tx)
	suite.Require().Nil(err)
	suite.True(cu.Published)
	err = c.Reload(suite.tx)
	suite.Require().Nil(err)
	suite.True(c.Published)
}

func (suite *RepoSuite) TestRemoveFile() {
	f := File{
		FileName:  "file_name",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      "012356789abcdef012356789abcdef0123456789",
		Size:      math.MaxInt64,
	}
	file, err := CreateFile(suite.tx, nil, f, nil)
	suite.Require().Nil(err)
	file.Published = true
	_, err = file.Update(suite.tx, boil.Whitelist("published"))
	suite.Require().Nil(err)
	cu, err := CreateContentUnit(suite.tx, common.CT_LESSON_PART, nil)
	suite.Require().Nil(err)
	cu.Published = true
	_, err = cu.Update(suite.tx, boil.Whitelist("published"))
	suite.Require().Nil(err)
	err = file.SetContentUnit(suite.tx, false, cu)
	suite.Require().Nil(err)
	c, err := CreateCollection(suite.tx, common.CT_DAILY_LESSON, nil)
	suite.Require().Nil(err)
	c.Published = true
	_, err = c.Update(suite.tx, boil.Whitelist("published"))
	suite.Require().Nil(err)
	err = c.AddCollectionsContentUnits(suite.tx, true, &models.CollectionsContentUnit{ContentUnitID: cu.ID})
	suite.Require().Nil(err)

	_, err = RemoveFile(suite.tx, file)
	suite.Require().Nil(err)
	err = file.Reload(suite.tx)
	suite.Require().Nil(err)
	suite.True(file.RemovedAt.Valid)
	err = cu.Reload(suite.tx)
	suite.Require().Nil(err)
	suite.False(cu.Published)
	err = c.Reload(suite.tx)
	suite.Require().Nil(err)
	suite.False(c.Published)
}

func (suite *RepoSuite) TestFindFileBySHA1() {
	const (
		sha1         = "012356789abcdef012356789abcdef0123456789"
		unknown_sha1 = "012356789abcdef012356789abcdef9876543210"
	)

	// test with empty table
	file, sha1b, err := FindFileBySHA1(suite.tx, sha1)
	suite.Exactly(FileNotFound{Sha1: sha1}, err, "empty error")
	suite.Equal(sha1, hex.EncodeToString(sha1b), "empty.sha1 bytes")
	suite.Nil(file, "empty.file")

	// test match with non empty table
	fm := File{
		FileName:  "file_name",
		CreatedAt: &Timestamp{time.Now()},
		Sha1:      sha1,
		Size:      math.MaxInt64,
	}
	newFile, err := CreateFile(suite.tx, nil, fm, nil)
	suite.Require().Nil(err)
	file, _, err = FindFileBySHA1(suite.tx, sha1)
	suite.Require().Nil(err)
	suite.NotNil(file)
	suite.Equal(newFile.ID, file.ID)

	// test miss with non empty table
	file, sha1b, err = FindFileBySHA1(suite.tx, unknown_sha1)
	suite.Exactly(FileNotFound{Sha1: unknown_sha1}, err, "miss error")
	suite.Equal(unknown_sha1, hex.EncodeToString(sha1b), "miss.sha1 bytes")
	suite.Nil(file)
}

func (suite *RepoSuite) TestGetFreeUID() {
	checkers := [6]UIDChecker{
		new(CollectionUIDChecker),
		new(ContentUnitUIDChecker),
		new(FileUIDChecker),
		new(OperationUIDChecker),
		new(SourceUIDChecker),
		new(TagUIDChecker),
	}

	for i := range checkers {
		checker := checkers[i]
		var uids = make(map[string]bool, 100)
		for i := 0; i < 100; i++ {
			uid, err := GetFreeUID(suite.tx, checker)
			suite.Require().Nil(err)
			suite.Regexp(utils.UID_REGEX, uid, "UID")
			if _, ok := uids[uid]; ok {
				suite.Fail("UID not unique")
			} else {
				uids[uid] = true
			}
		}
	}
}

// Helpers

func (suite *RepoSuite) getCustomProps() map[string]interface{} {
	return map[string]interface{}{
		"a": 1,
		"b": "2",
		"c": true,
		"d": []float64{1.2, 2.3, 3.4},
	}
}

func (suite *RepoSuite) assertCustomProps(expected map[string]interface{}, actual map[string]interface{}) {
	suite.EqualValues(expected["a"], actual["a"], "props: a")
	suite.Equal(expected["b"], actual["b"], "props: b")
	suite.Equal(expected["c"], actual["c"], "props: c")
	suite.Len(actual["d"].([]interface{}), len(expected["d"].([]float64)), "props: d length")
	for i, v := range expected["d"].([]float64) {
		suite.Equal(v, actual["d"].([]interface{})[i], "props: d[%d]", i)
	}
}
