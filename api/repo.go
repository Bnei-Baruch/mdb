package api

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/events"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
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

const FILE_DESCENDANTS_SQL = `
WITH RECURSIVE rf AS (
  SELECT f.*
  FROM files f
  WHERE f.id = $1
  UNION
  SELECT f.*
  FROM files f INNER JOIN rf ON f.parent_id = rf.id
) SELECT *
  FROM rf
  WHERE id != $1
`

const SOURCE_PATH_SQL = `
WITH RECURSIVE rs AS (
  SELECT s.*
  FROM sources s
  WHERE s.id = $1
  UNION
  SELECT s.*
  FROM sources s INNER JOIN rs ON s.id = rs.parent_id
) SELECT *
  FROM rs;
`

const TAG_PATH_SQL = `
WITH RECURSIVE rt AS (
  SELECT t.*
  FROM tags t
  WHERE t.id = $1
  UNION
  SELECT t.*
  FROM tags t INNER JOIN rt ON t.id = rt.parent_id
) SELECT *
  FROM rt;
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

const FILES_TREE_WITH_OPERATIONS = `
-- find all ancestors of a file
with ids as ((WITH RECURSIVE rfa AS (
  SELECT f.*
  FROM files f
  WHERE f.id = $1
  UNION
  SELECT f.*
  FROM files f INNER JOIN rfa ON f.id = rfa.parent_id
) 
SELECT id
  FROM rfa
  WHERE id != $1)

  UNION

-- find all descendants of a file
(WITH RECURSIVE rfd AS (
  SELECT f.*
  FROM files f
  WHERE f.id = $1
  UNION
  SELECT f.*
  FROM files f INNER JOIN rfd ON f.parent_id = rfd.id
) SELECT id
  FROM rfd))
  
 select f.id, f.uid, f.sha1, f.name, f.size, f.type, f.sub_type, f.mime_type, f.created_at, f.language, f.file_created_at, f.parent_id, f.published,
 array_agg(fop.operation_id) as OperationIds from ids
 join files f on f.id=ids.id
 join files_operations fop on fop.file_id = ids.id
 group by f.id
 `

func CreateOperation(exec boil.Executor, name string, o Operation, properties map[string]interface{}) (*models.Operation, error) {
	uid, err := GetFreeUID(exec, new(OperationUIDChecker))
	if err != nil {
		return nil, err
	}

	operation := models.Operation{
		TypeID:  common.OPERATION_TYPE_REGISTRY.ByName[name].ID,
		UID:     uid,
		Station: null.StringFrom(o.Station),
	}

	// Lookup user
	user, err := models.Users(qm.Where("email=?", o.User)).One(exec)
	if err == nil {
		operation.UserID = null.Int64From(user.ID)
	} else {
		return nil, errors.Wrap(err, "Check user exists")
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
			return nil, errors.Wrap(err, "json.Marshal")
		}
		operation.Properties = null.JSONFrom(props)
	}

	return &operation, operation.Insert(exec, boil.Infer())
}

func FindUpChainOperation(exec boil.Executor, fileID int64, opType string) (*models.Operation, error) {
	var op models.Operation

	opTypeID := common.OPERATION_TYPE_REGISTRY.ByName[opType].ID

	err := queries.Raw(UPCHAIN_OPERATION_SQL, fileID, opTypeID).Bind(nil, exec, &op)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UpChainOperationNotFound{FileID: fileID, opType: opType}
		} else {
			return nil, errors.Wrap(err, "DB lookup")
		}
	}

	return &op, nil
}

func FindOperationsByWorkflowID(exec boil.Executor, workflowID interface{}, opType string) ([]*models.Operation, error) {
	opTypeID := common.OPERATION_TYPE_REGISTRY.ByName[opType].ID
	return models.Operations(
		qm.Where("properties->>'workflow_id' = ?", workflowID),
		qm.Where("type_id = ?", opTypeID),
	).All(exec)
}

func CreateCollection(exec boil.Executor, contentType string, properties map[string]interface{}) (*models.Collection, error) {
	ct, ok := common.CONTENT_TYPE_REGISTRY.ByName[contentType]
	if !ok {
		return nil, errors.Errorf("Unknown content type %s", contentType)
	}

	uid, err := GetFreeUID(exec, new(CollectionUIDChecker))
	if err != nil {
		return nil, err
	}

	collection := &models.Collection{
		UID:    uid,
		TypeID: ct.ID,
	}

	if properties != nil {
		props, err := json.Marshal(properties)
		if err != nil {
			return nil, errors.Wrap(err, "json Marshal")
		}
		collection.Properties = null.JSONFrom(props)
	}

	err = collection.Insert(exec, boil.Infer())
	if err != nil {
		return nil, errors.Wrap(err, "Save to DB")
	}

	return collection, err
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
	_, err = collection.Update(exec, boil.Whitelist("properties"))
	if err != nil {
		return errors.Wrap(err, "Save properties to DB")
	}

	return nil
}

func FindCollectionByUID(exec boil.Executor, uid string) (*models.Collection, error) {
	return models.Collections(qm.Where("uid = ?", uid)).One(exec)
}

func FindCollectionByCaptureID(exec boil.Executor, cid interface{}) (*models.Collection, error) {
	var c models.Collection

	err := queries.Raw(
		`SELECT * FROM collections WHERE properties -> 'capture_id' ? $1`,
		cid).Bind(nil, exec, &c)
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
	ct, ok := common.CONTENT_TYPE_REGISTRY.ByName[contentType]
	if !ok {
		return nil, errors.Errorf("Unknown content type %s", contentType)
	}

	uid, err := GetFreeUID(exec, new(ContentUnitUIDChecker))
	if err != nil {
		return nil, err
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

	err = unit.Insert(exec, boil.Infer())
	if err != nil {
		return nil, errors.Wrap(err, "Save to DB")
	}

	return unit, err
}

func DeleteContentUnit(exec boil.Executor, unit *models.ContentUnit) error {
	log.Infof("Removing content_unit %d", unit.ID)

	tables := [...]string{
		"collections_content_units",
		"content_units_persons",
		"content_units_sources",
		"content_units_tags",
		"content_units_publishers",
		"content_unit_i18n",
	}
	for i := range tables {
		q := fmt.Sprintf("DELETE FROM %s WHERE content_unit_id = $1", tables[i])
		_, err := queries.Raw(q, unit.ID).Exec(exec)
		if err != nil {
			return errors.Wrapf(err, "Delete %s", tables[i])
		}
	}

	_, err := unit.Delete(exec)
	return err
}

func GetNextPositionInCollection(exec boil.Executor, id int64) (position int, err error) {
	err = queries.Raw(
		"SELECT COALESCE(MAX(position), -1) + 1 FROM collections_content_units WHERE collection_id = $1", id).
		QueryRow(exec).Scan(&position)
	return
}

func UpdateContentUnitProperties(exec boil.Executor, unit *models.ContentUnit, props map[string]interface{}) error {
	if len(props) == 0 {
		return nil
	}

	var p map[string]interface{}
	if unit.Properties.Valid {
		err := unit.Properties.Unmarshal(&p)
		if err != nil {
			return errors.Wrap(err, "json.Unmarshal")
		}
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

	unit.Properties = null.JSONFrom(fpa)
	_, err = unit.Update(exec, boil.Whitelist("properties"))
	if err != nil {
		return errors.Wrap(err, "Save properties to DB")
	}

	return nil
}

func CreateFile(exec boil.Executor, parent *models.File, f File, properties map[string]interface{}) (*models.File, error) {
	file, err := makeFile(parent, f, properties)
	if err != nil {
		return nil, errors.Wrap(err, "Make file")
	}

	uid, err := GetFreeUID(exec, new(FileUIDChecker))
	if err != nil {
		return nil, err
	}
	file.UID = uid

	err = file.Insert(exec, boil.Infer())
	if err != nil {
		return nil, errors.Wrap(err, "Save to DB")
	}

	return file, nil
}

func UpdateFile(exec boil.Executor, obj *models.File, parent *models.File, f File, properties map[string]interface{}) error {
	tmp, err := makeFile(parent, f, properties)
	if err != nil {
		return errors.Wrap(err, "Make file")
	}

	obj.Name = tmp.Name
	obj.Type = tmp.Type
	obj.SubType = tmp.SubType
	obj.MimeType = tmp.MimeType
	obj.ContentUnitID = tmp.ContentUnitID
	obj.Language = tmp.Language
	obj.ParentID = tmp.ParentID
	obj.FileCreatedAt = tmp.FileCreatedAt

	_, err = obj.Update(exec, boil.Whitelist("name", "type", "sub_type", "mime_type", "content_unit_id",
		"language", "parent_id", "file_created_at"))
	if err != nil {
		return errors.Wrap(err, "update file")
	}

	err = UpdateFileProperties(exec, obj, properties)
	if err != nil {
		return errors.Wrap(err, "update properties")
	}

	return nil
}

func makeFile(parent *models.File, f File, properties map[string]interface{}) (*models.File, error) {
	sha1, err := hex.DecodeString(f.Sha1)
	if err != nil {
		return nil, errors.Wrap(err, "hex Decode")
	}

	// Standardize and validate language
	var mdbLang = ""
	if f.Language != "" {
		mdbLang = common.StdLang(f.Language)
		if mdbLang == common.LANG_UNKNOWN && f.Language != common.LANG_UNKNOWN {
			return nil, errors.Errorf("Unknown language %s", f.Language)
		}
	}

	file := &models.File{
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
			if mt, ok := common.MEDIA_TYPE_REGISTRY.ByMime[strings.ToLower(f.MimeType)]; ok {
				file.Type = mt.Type
				file.SubType = mt.SubType
			}
		}
	}

	if parent != nil {
		file.ParentID = null.Int64From(parent.ID)
		file.ContentUnitID = parent.ContentUnitID
	}

	// Handle properties
	if properties != nil {
		props, err := json.Marshal(properties)
		if err != nil {
			return nil, errors.Wrap(err, "json Marshal")
		}
		file.Properties = null.JSONFrom(props)
	}

	return file, nil
}

func UpdateFileProperties(exec boil.Executor, file *models.File, props map[string]interface{}) error {
	if len(props) == 0 {
		return nil
	}

	var p map[string]interface{}
	if file.Properties.Valid {
		err := file.Properties.Unmarshal(&p)
		if err != nil {
			return errors.Wrap(err, "json.Unmarshal")
		}
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
	_, err = file.Update(exec, boil.Whitelist("properties"))
	if err != nil {
		return errors.Wrap(err, "Save properties to DB")
	}

	return nil
}

type PublishedChangeImpact struct {
	ChangedContentUnit *models.ContentUnit
	ChangedCollections []*models.Collection
}

func (p *PublishedChangeImpact) Events() []events.Event {
	evnts := make([]events.Event, 0)

	if p.ChangedContentUnit != nil {
		evnts = append(evnts, events.ContentUnitPublishedChangeEvent(p.ChangedContentUnit))
	}
	if p.ChangedCollections != nil {
		for i := range p.ChangedCollections {
			evnts = append(evnts, events.CollectionPublishedChangeEvent(p.ChangedCollections[i]))
		}
	}

	return evnts
}

func PublishFile(exec boil.Executor, file *models.File) (*PublishedChangeImpact, error) {
	log.Infof("Publishing file [%d]", file.ID)
	file.Published = true
	_, err := file.Update(exec, boil.Whitelist("published"))
	if err != nil {
		return nil, errors.Wrap(err, "Save file to DB")
	}

	if !file.ContentUnitID.Valid {
		return new(PublishedChangeImpact), nil
	}

	return FileAddedUnitImpact(exec, file.Published, file.ContentUnitID.Int64)
}

func RemoveFile(exec boil.Executor, file *models.File) (*PublishedChangeImpact, error) {
	log.Infof("Removing file [%d]", file.ID)
	file.RemovedAt = null.TimeFrom(time.Now().UTC())
	_, err := file.Update(exec, boil.Whitelist("removed_at"))
	if err != nil {
		return nil, errors.Wrap(err, "Save file to DB")
	}

	if !file.ContentUnitID.Valid {
		return new(PublishedChangeImpact), nil
	}

	return FileLeftUnitImpact(exec, file.Published, file.ContentUnitID.Int64)
}

func FileAddedUnitImpact(exec boil.Executor, fileIsPublished bool, cuID int64) (*PublishedChangeImpact, error) {
	impact := new(PublishedChangeImpact)

	if !fileIsPublished {
		return impact, nil
	}

	// Load CU
	cu, err := models.ContentUnits(
		qm.Where("id=?", cuID),
		qm.Load("CollectionsContentUnits"),
		qm.Load("CollectionsContentUnits.Collection"),
	).One(exec)
	if err != nil {
		return nil, errors.Wrapf(err, "Load content_unit %d", cuID)
	}

	// Publish CU and associated collections if necessary
	if !cu.Published {
		cu.Published = true
		if _, err := cu.Update(exec, boil.Whitelist("published")); err != nil {
			return nil, errors.Wrapf(err, "Update content_unit %d", cuID)
		}
		impact.ChangedContentUnit = cu

		// handle associated collections
		if len(cu.R.CollectionsContentUnits) > 0 {
			for i := range cu.R.CollectionsContentUnits {
				c := cu.R.CollectionsContentUnits[i].R.Collection
				if !c.Published {
					c.Published = true
					if _, err := c.Update(exec, boil.Whitelist("published")); err != nil {
						return nil, errors.Wrapf(err, "Update collection %d", cuID)
					}
					impact.ChangedCollections = append(impact.ChangedCollections, c)
				}
			}
		}
	}

	return impact, nil
}

func FileLeftUnitImpact(exec boil.Executor, fileIsPublished bool, cuID int64) (*PublishedChangeImpact, error) {
	impact := new(PublishedChangeImpact)

	if !fileIsPublished {
		return impact, nil
	}

	// Load CU
	cu, err := models.ContentUnits(
		qm.Where("id=?", cuID),
		qm.Load("Files"),
		qm.Load("CollectionsContentUnits"),
	).One(exec)
	if err != nil {
		return nil, errors.Wrapf(err, "Load content_unit %d", cuID)
	}

	if !cu.Published {
		return impact, nil
	}

	// Check if any other file in CU is published
	unpublishCU := true
	for i := range cu.R.Files {
		f := cu.R.Files[i]
		if f.Published && !f.RemovedAt.Valid {
			unpublishCU = false
			break
		}
	}

	// cu has other published files so no change
	if !unpublishCU {
		return impact, nil
	}

	// unpublish content unit
	cu.Published = false
	if _, err := cu.Update(exec, boil.Whitelist("published")); err != nil {
		return nil, errors.Wrapf(err, "Update [published=false] content_unit %d", cuID)
	}
	impact.ChangedContentUnit = cu

	// Load all collections associated with this CU and do the same for them
	if len(cu.R.CollectionsContentUnits) > 0 {

		// Load collections
		cIDs := make([]int64, len(cu.R.CollectionsContentUnits))
		for i := range cu.R.CollectionsContentUnits {
			cIDs[i] = cu.R.CollectionsContentUnits[i].CollectionID
		}
		cs, err := models.Collections(
			qm.WhereIn("id in ?", utils.ConvertArgsInt64(cIDs)...),
			qm.Load("CollectionsContentUnits"),
			qm.Load("CollectionsContentUnits.ContentUnit")).
			All(exec)
		if err != nil {
			return nil, errors.Wrapf(err, "Load collections CCU's %v", cIDs)
		}

		// Check if each collection has any other published CU and unpublish if not
		for i := range cs {
			c := cs[i]
			if c.Published {
				unpublishC := true
				for i := range c.R.CollectionsContentUnits {
					cu := c.R.CollectionsContentUnits[i].R.ContentUnit
					if cu.Published {
						unpublishC = false
						break
					}
				}

				if unpublishC {
					c.Published = false
					if _, err := c.Update(exec, boil.Whitelist("published")); err != nil {
						return nil, errors.Wrapf(err, "Update [published=false] collection %d", cuID)
					}
					impact.ChangedCollections = append(impact.ChangedCollections, c)
				}
			}
		}
	}

	return impact, nil
}

func FindFileBySHA1(exec boil.Executor, sha1 string) (*models.File, []byte, error) {
	s, err := hex.DecodeString(sha1)
	if err != nil {
		return nil, nil, errors.Wrap(err, "hex decode")
	}

	f, err := models.Files(qm.Where("sha1=?", s)).One(exec)
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

func FindFileAncestors(exec boil.Executor, id int64) ([]*models.File, error) {
	var ancestors []*models.File

	err := queries.Raw(FILE_ANCESTORS_SQL, id).Bind(nil, exec, &ancestors)
	if err != nil {
		return nil, errors.Wrap(err, "DB lookup")
	}

	return ancestors, nil
}

func FindFileDescendants(exec boil.Executor, id int64) ([]*models.File, error) {
	var descendants []*models.File

	err := queries.Raw(FILE_DESCENDANTS_SQL, id).Bind(nil, exec, &descendants)
	if err != nil {
		return nil, errors.Wrap(err, "DB lookup")
	}

	return descendants, nil
}

func FindFileTreeWithOperations(exec boil.Executor, fileID int64) ([]*MFile, error) {
	files := make([]*MFile, 0)

	rows, err := queries.Raw(FILES_TREE_WITH_OPERATIONS, fileID).Query(exec)
	if err != nil {
		return nil, NewInternalError(err)
	}
	defer rows.Close()

	for rows.Next() {
		f := new(MFile)
		err := rows.Scan(&f.ID, &f.UID, &f.Sha1, &f.Name, &f.Size, &f.Type, &f.SubType, &f.MimeType, &f.CreatedAt,
			&f.Language, &f.FileCreatedAt, &f.ParentID, &f.Published, &f.OperationIds)
		if err != nil {
			return nil, NewInternalError(err)
		}
		if f.Sha1.Valid {
			f.Sha1Str = hex.EncodeToString(f.Sha1.Bytes)
		}
		files = append(files, f)
	}

	err = rows.Err()
	if err != nil {
		return nil, NewInternalError(err)
	}

	return files, nil
}

func FindSourceByUID(exec boil.Executor, uid string) (*models.Source, error) {
	return models.Sources(qm.Where("uid = ?", uid)).One(exec)
}

func FindSourcePath(exec boil.Executor, id int64) ([]*models.Source, error) {
	var ancestors []*models.Source

	err := queries.Raw(SOURCE_PATH_SQL, id).Bind(nil, exec, &ancestors)
	if err != nil {
		return nil, errors.Wrap(err, "DB lookup")
	}

	return ancestors, nil
}

func FindAuthorBySourceID(exec boil.Executor, id int64) (*models.Author, error) {
	return models.Authors(
		qm.InnerJoin("authors_sources as x on x.author_id=id and x.source_id = ?", id),
		qm.Load("AuthorI18ns")).
		One(exec)
}

func FindTagByUID(exec boil.Executor, uid string) (*models.Tag, error) {
	return models.Tags(qm.Where("uid = ?", uid)).One(exec)
}

func FindTagPath(exec boil.Executor, id int64) ([]*models.Tag, error) {
	var ancestors []*models.Tag

	err := queries.Raw(TAG_PATH_SQL, id).Bind(nil, exec, &ancestors)
	if err != nil {
		return nil, errors.Wrap(err, "DB lookup")
	}

	return ancestors, nil
}

type UIDChecker interface {
	Check(exec boil.Executor, uid string) (exists bool, err error)
}

type CollectionUIDChecker struct{}

func (c *CollectionUIDChecker) Check(exec boil.Executor, uid string) (exists bool, err error) {
	return models.Collections(qm.Where("uid = ?", uid)).Exists(exec)
}

type ContentUnitUIDChecker struct{}

func (c *ContentUnitUIDChecker) Check(exec boil.Executor, uid string) (exists bool, err error) {
	return models.ContentUnits(qm.Where("uid = ?", uid)).Exists(exec)
}

type FileUIDChecker struct{}

func (c *FileUIDChecker) Check(exec boil.Executor, uid string) (exists bool, err error) {
	return models.Files(qm.Where("uid = ?", uid)).Exists(exec)
}

type OperationUIDChecker struct{}

func (c *OperationUIDChecker) Check(exec boil.Executor, uid string) (exists bool, err error) {
	return models.Operations(qm.Where("uid = ?", uid)).Exists(exec)
}

type SourceUIDChecker struct{}

func (c *SourceUIDChecker) Check(exec boil.Executor, uid string) (exists bool, err error) {
	return models.Sources(qm.Where("uid = ?", uid)).Exists(exec)
}

type TagUIDChecker struct{}

func (c *TagUIDChecker) Check(exec boil.Executor, uid string) (exists bool, err error) {
	return models.Tags(qm.Where("uid = ?", uid)).Exists(exec)
}

type PersonUIDChecker struct{}

func (c *PersonUIDChecker) Check(exec boil.Executor, uid string) (exists bool, err error) {
	return models.Persons(qm.Where("uid = ?", uid)).Exists(exec)
}

type PublisherUIDChecker struct{}

func (c *PublisherUIDChecker) Check(exec boil.Executor, uid string) (exists bool, err error) {
	return models.Publishers(qm.Where("uid = ?", uid)).Exists(exec)
}

type LabelUIDChecker struct{}

func (c *LabelUIDChecker) Check(exec boil.Executor, uid string) (exists bool, err error) {
	return models.Labels(qm.Where("uid = ?", uid)).Exists(exec)
}

func GetFreeUID(exec boil.Executor, checker UIDChecker) (uid string, err error) {
	for {
		uid = utils.GenerateUID(8)
		exists, ex := checker.Check(exec, uid)
		if ex != nil {
			err = errors.Wrap(ex, "Check UID exists")
			break
		}
		if !exists {
			break
		}
	}

	return
}
