package api

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"strings"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/nullbio/null.v6"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/vattle/sqlboiler/queries"
)

const FILE_ANCESTORS_SQL = `
WITH RECURSIVE rf AS (
  SELECT f.*
  FROM files f
  WHERE f.id = $1
  UNION
  SELECT f.*
  FROM files f INNER JOIN rf ON f.id = rf.parent_id
) SELECT *
  FROM rf
  WHERE id != $1
`

const UPCHAIN_OPERATION_SQL = `
WITH RECURSIVE
    rf AS (
    SELECT
      f.id,
      f.parent_id,
      NULL :: BIGINT "o_id",
      NULL :: BIGINT "o_type"
    FROM files f
    WHERE f.id = $1
    UNION
    SELECT
      f.id,
      f.parent_id,
      o.id      "o_id",
      o.type_id "o_type"
    FROM files f INNER JOIN rf ON f.id = rf.parent_id
      ,
      operations o
    WHERE o.id = (SELECT min(operation_id)
                  FROM files_operations
                  WHERE file_id = f.id)
  ) SELECT *
    FROM operations
    WHERE id = (SELECT o_id
                FROM rf
                WHERE o_type = $2);
`

func CreateOperation(exec boil.Executor, name string, o Operation, properties map[string]interface{}) (*models.Operation, error) {
	operation := models.Operation{
		TypeID:  OPERATION_TYPE_REGISTRY.ByName[name].ID,
		UID:     utils.GenerateUID(8),
		Station: null.StringFrom(o.Station),
	}

	// Lookup user, skip if doesn't exist
	user, err := models.Users(exec, qm.Where("email=?", o.User)).One()
	if err == nil {
		operation.UserID = null.Int64From(user.ID)
	} else {
		if err == sql.ErrNoRows {
			log.Debugf("Unknown User [%s]. Skipping.", o.User)
		} else {
			return nil, err
		}
	}

	// Handle properties
	if o.WorkflowID != "" {
		if properties == nil {
			properties = make(map[string]interface{})
		}
		properties["workflow_id"] = o.WorkflowID
	}
	if properties != nil {
		props, err := json.Marshal(properties)
		if err != nil {
			return nil, err
		}
		operation.Properties = null.JSONFrom(props)
	}

	return &operation, operation.Insert(exec)
}

func FindUpChainOperation(exec boil.Executor, fileID int64, opType int64) (*models.Operation, error) {
	var op models.Operation

	err := queries.Raw(exec, UPCHAIN_OPERATION_SQL, fileID, opType).Bind(&op)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UpChainOperationNotFound{FileID: fileID, OperationType: opType}
		} else {
			return nil, errors.Wrap(err, "DB lookup")
		}
	}

	return &op, nil
}

func CreateCollection(exec boil.Executor, contentType string, properties map[string]interface{}) (*models.Collection, error) {
	ct, ok := CONTENT_TYPE_REGISTRY.ByName[contentType]
	if !ok {
		return nil, errors.Errorf("Unknown content type %s", contentType)
	}

	var uid string
	for {
		uid = utils.GenerateUID(8)
		exists, err := models.Collections(exec, qm.Where("uid = ?", uid)).Exists()
		if err != nil {
			return nil, errors.Wrap(err, "Check UID exists")
		}
		if !exists {
			break
		}
	}

	unit := &models.Collection{
		UID:    uid,
		TypeID: ct.ID,
	}

	if properties != nil {
		props, err := json.Marshal(properties)
		if err != nil {
			return nil, errors.Wrap(err, "json Marshal")
		}
		unit.Properties = null.JSONFrom(props)
	}

	err := unit.Insert(exec)
	if err != nil {
		return nil, errors.Wrap(err, "Save to DB")
	}

	return unit, err
}

func UpdateCollectionProperties(exec boil.Executor, collection *models.Collection, props map[string]interface{}) error {
	if len(props) == 0 {
		return nil
	}

	var p map[string]interface{}
	if collection.Properties.Valid {
		collection.Properties.Unmarshal(&p)
		for k, v := range props {
			p[k] = v
		}
	} else {
		p = props
	}

	fpa, err := json.Marshal(p)
	if err != nil {
		return errors.Wrap(err, "json Marshal")
	}

	collection.Properties = null.JSONFrom(fpa)
	err = collection.Update(exec, "properties")
	if err != nil {
		return errors.Wrap(err, "Save properties to DB")
	}

	return nil
}

func FindCollectionByCaptureID(exec boil.Executor, cid interface{}) (*models.Collection, error) {
	var c models.Collection

	err := queries.Raw(exec,
		`SELECT * FROM collections WHERE properties -> 'capture_id' ? $1`,
		cid).Bind(&c)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, CollectionNotFound{CaptureID: cid}
		} else {
			return nil, errors.Wrap(err, "DB lookup")
		}
	}

	return &c, nil
}

func CreateContentUnit(exec boil.Executor, contentType string, properties map[string]interface{}) (*models.ContentUnit, error) {
	ct, ok := CONTENT_TYPE_REGISTRY.ByName[contentType]
	if !ok {
		return nil, errors.Errorf("Unknown content type %s", contentType)
	}

	var uid string
	for {
		uid = utils.GenerateUID(8)
		exists, err := models.ContentUnits(exec, qm.Where("uid = ?", uid)).Exists()
		if err != nil {
			return nil, errors.Wrap(err, "Check UID exists")
		}
		if !exists {
			break
		}
	}

	unit := &models.ContentUnit{
		UID:    uid,
		TypeID: ct.ID,
	}

	if properties != nil {
		props, err := json.Marshal(properties)
		if err != nil {
			return nil, errors.Wrap(err, "json Marshal")
		}
		unit.Properties = null.JSONFrom(props)
	}

	err := unit.Insert(exec)
	if err != nil {
		return nil, errors.Wrap(err, "Save to DB")
	}

	return unit, err
}

func CreateFile(exec boil.Executor, parent *models.File, f File, properties map[string]interface{}) (*models.File, error) {
	sha1, err := hex.DecodeString(f.Sha1)
	if err != nil {
		return nil, errors.Wrap(err, "hex Decode")
	}

	// Standardize and validate language
	var mdbLang = ""
	if f.Language != "" {
		mdbLang = StdLang(f.Language)
		if mdbLang == LANG_UNKNOWN && f.Language != LANG_UNKNOWN {
			return nil, errors.Errorf("Unknown language %s", f.Language)
		}
	}

	file := &models.File{
		UID:           utils.GenerateUID(8),
		Name:          f.FileName,
		Sha1:          null.BytesFrom(sha1),
		Size:          f.Size,
		FileCreatedAt: null.TimeFrom(f.CreatedAt.Time),
		Type:          f.Type,
		SubType:       f.SubType,
		Language:      null.NewString(mdbLang, mdbLang != ""),
	}

	if f.MimeType != "" {
		file.MimeType = null.StringFrom(f.MimeType)

		// Try to complement missing type and subtype
		if file.Type == "" && file.SubType == "" {
			if mt, ok := MEDIA_TYPE_REGISTRY.ByMime[strings.ToLower(f.MimeType)]; ok {
				file.Type = mt.Type
				file.SubType = mt.SubType
			}
		}
	}

	if parent != nil {
		file.ParentID = null.Int64From(parent.ID)
	}

	// Handle properties
	if properties != nil {
		props, err := json.Marshal(properties)
		if err != nil {
			return nil, errors.Wrap(err, "json Marshal")
		}
		file.Properties = null.JSONFrom(props)
	}

	err = file.Insert(exec)
	if err != nil {
		return nil, errors.Wrap(err, "Save to DB")
	}

	return file, nil
}

func UpdateFileProperties(exec boil.Executor, file *models.File, props map[string]interface{}) error {
	if len(props) == 0 {
		return nil
	}

	var p map[string]interface{}
	if file.Properties.Valid {
		file.Properties.Unmarshal(&p)
		for k, v := range props {
			p[k] = v
		}
	} else {
		p = props
	}

	fpa, err := json.Marshal(p)
	if err != nil {
		return errors.Wrap(err, "json Marshal")
	}

	file.Properties = null.JSONFrom(fpa)
	err = file.Update(exec, "properties")
	if err != nil {
		return errors.Wrap(err, "Save properties to DB")
	}

	return nil
}

func FindFileBySHA1(exec boil.Executor, sha1 string) (*models.File, []byte, error) {
	s, err := hex.DecodeString(sha1)
	if err != nil {
		return nil, nil, errors.Wrap(err, "hex decode")
	}

	f, err := models.Files(exec, qm.Where("sha1=?", s)).One()
	if err == nil {
		return f, s, nil
	} else {
		if err == sql.ErrNoRows {
			return nil, s, FileNotFound{Sha1: sha1}
		} else {
			return nil, s, errors.Wrap(err, "DB lookup")
		}
	}
}

func FindFileAncestors(exec boil.Executor, fileID int64) ([]*models.File, error) {
	var ancestors []*models.File

	err := queries.Raw(exec, FILE_ANCESTORS_SQL, fileID).Bind(&ancestors)
	if err != nil {
		return nil, errors.Wrap(err, "DB lookup")
	}

	return ancestors, nil
}

// Return standard language or LANG_UNKNOWN
//
// 	if len(lang) = 2 we assume it's an MDB language code and check KNOWN_LANGS.
// 	if len(lang) = 3 we assume it's a workflow / kmedia lang code and check LANG_MAP.
func StdLang(lang string) string {
	switch len(lang) {
	case 2:
		if l := strings.ToLower(lang); KNOWN_LANGS.MatchString(l) {
			return l
		}
	case 3:
		if l, ok := LANG_MAP[strings.ToUpper(lang)]; ok {
			return l
		}
	}

	return LANG_UNKNOWN
}
