package metus

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"bufio"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"time"
)

type FieldsHelper struct {
	byID               map[int]*Field
	roots              []*Field
	nonMissingMetadata int
}

func (fh *FieldsHelper) Load() error {
	rows, err := metusDB.Query(`with rec_fields as (select *
                    from fields
                    where ParentID = -1
    union all
    select f.*
    from FIELDS f
           inner join rec_fields rf on f.ParentID = rf.ID)
select rf.ID,
       rf.ParentID,
       rf.GUID,
       rf.Version,
       rf.UnitID,
       rf.Type,
       rf.DefinitionPackageID,
       fl.Name_0,
       fl.Description_0
from rec_fields rf
       left join FIELD_LANG fl on rf.ID = fl.FieldID;`)

	if err != nil {
		return errors.Wrap(err, "fetch fields")
	}
	defer rows.Close()

	fh.byID = make(map[int]*Field)
	fh.roots = make([]*Field, 0)

	for rows.Next() {
		var f Field
		err := rows.Scan(&f.ID,
			&f.ParentID,
			&f.GUID,
			&f.Version,
			&f.UnitID,
			&f.Type,
			&f.DefinitionPackageID,
			&f.Name,
			&f.Description)
		if err != nil {
			return errors.Wrap(err, "rows.Scan")
		}

		fh.byID[f.ID] = &f
		if f.ParentID > 0 {
			p := fh.byID[f.ParentID]
			p.Children = append(p.Children, &f)
		} else {
			fh.roots = append(fh.roots, &f)
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Wrap(err, "rows.Err")
	}

	log.Infof("FieldsHelper.Load: len(fieldMap) %d [%d roots]", len(fh.byID), len(fh.roots))

	return nil
}

func (fh *FieldsHelper) dump() {
	s := fh.roots[:]
	var node *Field
	for len(s) > 0 {
		node, s = s[0], s[1:]
		if len(node.Children) > 0 {
			s = append(node.Children[:], s...) // DFS
			//s = append(s, node.Children...)      // BFS
		}
		log.Infof("%d\t\t%s", node.ID, strings.Join(fh.getJsonPath(node.ID), "."))
	}
}

func (fh *FieldsHelper) getJsonPath(fID int) []string {
	path := make([]string, 0)

	f, ok := fh.byID[fID]
	if !ok {
		return path
	}
	path = append(path, f.JsonKey())
	for f.ParentID > 0 {
		f = fh.byID[f.ParentID]
		path = append([]string{f.JsonKey()}, path...)
	}

	return path
}

func (fh *FieldsHelper) getMetadataAsJson(o *Object) map[string]interface{} {
	mj := make(map[string]interface{})

	for i := range o.Metadata {
		m := o.Metadata[i]

		// skip missing values
		if (!m.ValueString.Valid || m.ValueString.String == "") && !m.ValueNumber.Valid {
			continue
		}

		fh.nonMissingMetadata++

		path := fh.getJsonPath(m.FieldID)
		c := mj
		for j := range path {
			k := path[j]
			if j == len(path)-1 {
				if m.ValueString.Valid {
					c[k] = m.ValueString.String
				} else if m.ValueNumber.Valid {
					c[k] = m.ValueNumber.Float64
				}
			} else {
				if v, ok := c[k]; ok {
					if cv, ok := v.(map[string]interface{}); ok {
						c = cv
					} else {
						log.Warnf("conflicting paths: %d\t%d\t%v", o.ID, m.RowID, path)
					}
				} else {
					c[k] = make(map[string]interface{})
					c = c[k].(map[string]interface{})
				}
			}
		}

	}

	return mj
}

type ObjectsHelper struct {
	byID  map[int]*Object
	roots []*Object
}

func (oh *ObjectsHelper) LoadFromDB() error {
	rows, err := metusDB.Query(`select ObjectID,
       ParentObjectID,
       GUID,
       ObjectType,
       SubType,
       AssetType,
       AssetFormat,
       HasSecurity,
       SecurityTFU,
       Status,
       SubStatus,
       IsDeleted,
       IsLocked,
       IsProtected,
       FileSignature,
       ObjectOrder
from OBJECTS;`)

	if err != nil {
		return errors.Wrap(err, "fetch objects")
	}
	defer rows.Close()

	oh.byID = make(map[int]*Object)
	for rows.Next() {
		var o Object
		err := rows.Scan(&o.ID,
			&o.ParentID,
			&o.GUID,
			&o.Type,
			&o.SubType,
			&o.AssetType,
			&o.AssetFormat,
			&o.HasSecurity,
			&o.SecurityTFU,
			&o.Status,
			&o.SubStatus,
			&o.IsDeleted,
			&o.IsLocked,
			&o.IsProtected,
			&o.FileSignature,
			&o.ObjectOrder)
		if err != nil {
			return errors.Wrap(err, "rows.Scan")
		}

		oh.byID[o.ID] = &o
	}

	if err := rows.Err(); err != nil {
		return errors.Wrap(err, "rows.Err")
	}

	// construct hierarchy
	for _, o := range oh.byID {
		if o.ParentID > 0 {
			if p, ok := oh.byID[o.ParentID]; ok {
				p.Children = append(p.Children, o)
			}
		} else {
			oh.roots = append(oh.roots, o)
		}
	}

	log.Infof("len(oh.byID) %d [%d roots]", len(oh.byID), len(oh.roots))

	// load METADATA_0
	mrows, err := metusDB.Query(`select MetadataID,
ObjectID,
FieldID,
SubFieldID,
RowID,
Value_String,
Value_Number
-- cast(row_version as datetime)
from METADATA_0`)
	if err != nil {
		return errors.Wrap(err, "fetch metadata")
	}
	defer mrows.Close()

	mTotal := 0
	mHist := make(map[int]int)
	for mrows.Next() {
		var m MetaData
		err := mrows.Scan(&m.ID,
			&m.ObjectID,
			&m.FieldID,
			&m.SubFieldID,
			&m.RowID,
			&m.ValueString,
			&m.ValueNumber,
			//&m.RowVersion,
		)
		if err != nil {
			return errors.Wrap(err, "scan metadata")
		}

		mTotal++
		mHist[m.FieldID]++

		if o, ok := oh.byID[m.ObjectID]; ok {
			o.Metadata = append(o.Metadata, &m)
		} else {
			log.Warnf("no object found %d %d", m.ID, m.ObjectID)
		}
	}
	if err := mrows.Err(); err != nil {
		return errors.Wrap(err, "mrows.Err")
	}

	log.Infof("%d metadata key-values", mTotal)
	//for k, v := range mHist {
	//	log.Infof("%d\t%d", k, v)
	//}

	return nil
}

func (oh *ObjectsHelper) LoadFromDisk(rootDir string) error {
	oh.byID = make(map[int]*Object)

	// read json files from disk to build byID
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || info.Name() == "files-catalog-nas" {
			return nil
		}

		var o Object
		f, err := os.Open(path)
		if err != nil {
			return errors.Wrapf(err, "os.Open %s", path)
		}
		if err := json.NewDecoder(f).Decode(&o); err != nil {
			return errors.Wrapf(err, "json.Decode %s", path)
		}

		oh.byID[o.ID] = &o

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "filepath.Walk")
	}

	// resurrect hierarchy from byID
	for _, o := range oh.byID {
		if o.ParentID > 0 {
			if p, ok := oh.byID[o.ParentID]; ok {
				p.Children = append(p.Children, o)
			}
		} else {
			oh.roots = append(oh.roots, o)
		}
	}

	return nil
}

func (oh *ObjectsHelper) getIDPath(oID int) []string {
	path := make([]string, 0)

	o, ok := oh.byID[oID]
	if !ok {
		return path
	}
	for o != nil && o.ParentID > 0 {
		path = append([]string{strconv.Itoa(o.ParentID)}, path...)
		o = oh.byID[o.ParentID]
	}

	return path
}

func (oh *ObjectsHelper) getNamePath(oID int) []string {
	path := make([]string, 0)

	o, ok := oh.byID[oID]
	if !ok {
		return []string{"NO_NAME"}
	}
	for o != nil && o.ParentID > 0 {
		name := o.getDeepValueFallback(
			"metadata-fields.metus.file.file-name",
			"metadata-fields.metus.general.former-file-name",
			"metadata-fields.metus.general.file-name")
		v, ok := name.(string)
		if !ok {
			v = fmt.Sprintf("NO_NAME [%d]", o.ID)
		}
		path = append([]string{v}, path...)
		o = oh.byID[o.ParentID]
	}

	return path
}

func (oh *ObjectsHelper) WalkDFS(walkFn func(*Object) error) error {
	return oh.Walk(walkFn, true)
}

func (oh *ObjectsHelper) WalkBFS(walkFn func(*Object) error) error {
	return oh.Walk(walkFn, false)
}

func (oh *ObjectsHelper) Walk(walkFn func(*Object) error, dfs bool) error {
	s := oh.roots[:]
	var node *Object
	for len(s) > 0 {
		node, s = s[0], s[1:]
		if len(node.Children) > 0 {
			if dfs {
				s = append(node.Children[:], s...) // DFS
			} else {
				s = append(s, node.Children...) // BFS
			}
		}
		if err := walkFn(node); err != nil {
			return err
		}
	}

	return nil
}

type FileRecord struct {
	Path    string
	Sha1    string
	Size    int64
	ModTime time.Time
}

type PhysicalIndex struct {
	idx map[string]*FileRecord
}

func (pi *PhysicalIndex) Load() error {
	f, err := os.Open("importer/metus/data/files-catalog-nas")
	if err != nil {
		return errors.Wrap(err, "os.Open")
	}
	defer f.Close()

	pi.idx = make(map[string]*FileRecord)

	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		i++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var v []interface{}
		err := json.NewDecoder(strings.NewReader(line)).Decode(&v)
		if err != nil {
			return errors.Wrapf(err, "json.Decode [line %d]", i)
		}

		path := v[0].(string)
		if strings.HasPrefix(path, "/net/nas/F-") ||
			strings.HasPrefix(path, "/net/nas/H-") {
			path = path[11:]
		} else if strings.HasPrefix(path, "/net/nas/H/F/") {
			path = fmt.Sprintf("4dcc9fad-c769-45c1-ae46-a7bc93764af8%s", path[10:])
		}

		fr := &FileRecord{
			Path:    path,
			Sha1:    v[1].(string),
			Size:    int64(v[2].(float64)),
			ModTime: time.Unix(int64(v[3].(float64)), 0),
		}

		pi.idx[fr.Path] = fr
	}

	return nil
}

func (pi *PhysicalIndex) Lookup(o *Object) *FileRecord {
	path := o.getPhysicalFilepath()
	if path == "" {
		return nil
	}

	return pi.idx[path]
}

func (pi *PhysicalIndex) Match(fileObjects []*Object) {
	// first iteration find by simple lookup
	// or terminate if file metus' IsDelete is true

	noPhys := make([]*Object, 0)
	physByName := make(map[string]*FileRecord)
	for i := range fileObjects {
		f := fileObjects[i]
		if fr := pi.Lookup(f); fr != nil {
			f.FileRecord = fr
			physByName[f.getPhysicalFilename()] = fr
		} else if !f.IsDeleted.Valid || !f.IsDeleted.Bool {
			noPhys = append(noPhys, f)
		}
	}

	// second iteration we try by duplicate physical file name
	// or matching file-name suffix

	for i := range noPhys {
		f := noPhys[i]
		if v, ok := physByName[f.getPhysicalFilename()]; ok {
			f.FileRecord = v
			f.IsDuplicate = true
		} else {
			if fnv := f.getDeepValue("metadata-fields.metus.file.file-name"); fnv != nil {
				n := fnv.(string)
				for k, fr := range pi.idx {
					if strings.HasSuffix(k, n) {
						f.FileRecord = fr
						break
					}
				}

			}
		}
	}

}
