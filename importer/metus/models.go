package metus

import (
	"regexp"
	"strings"
	"time"

	"gopkg.in/volatiletech/null.v6"
)

var LEGAL_JSON_KEY_RE = regexp.MustCompile("[^a-zA-Z0-9]+")
var MD5_RE = regexp.MustCompile("[^a-f0-9]{32}")

type Field struct {
	ID                  int
	ParentID            int
	GUID                string
	Version             string
	UnitID              int
	Type                int
	DefinitionPackageID int

	Name        null.String
	Description null.String

	Children []*Field
}

func (f *Field) JsonKey() string {
	c := LEGAL_JSON_KEY_RE.ReplaceAllString(f.Name.String, " ")
	return strings.ToLower(strings.Join(strings.Fields(c), "-"))
}

type MetaData struct {
	ID          int64
	ObjectID    int
	FieldID     int
	SubFieldID  int
	RowID       int
	ValueString null.String
	ValueNumber null.Float64
	RowVersion  time.Time
}

type Object struct {
	ID            int
	ParentID      int
	GUID          string
	Type          int
	SubType       int
	AssetType     int
	AssetFormat   int
	HasSecurity   null.Bool
	SecurityTFU   null.Bool
	Status        int
	SubStatus     int
	IsDeleted     null.Bool
	IsLocked      null.Bool
	IsProtected   null.Bool
	FileSignature null.String
	ObjectOrder   null.Float64

	Metadata     []*MetaData `json:"-"`
	MetadataJson map[string]interface{}

	Children        []*Object   `json:"-"`
	FileRecord      *FileRecord `json:"-"`
	IsDuplicate     bool        `json:"-"`
	IsInFilteredBin bool        `json:"-"`
	matchedCU       int64       `json:"-"`
}

func (o *Object) getDeepValue(key string) interface{} {
	s := strings.Split(key, ".")
	var x interface{}
	x = o.MetadataJson
	for i := range s {
		if xm, ok := x.(map[string]interface{}); ok {
			x = xm[s[i]]
		} else {
			return x
		}
	}
	return x
}

func (o *Object) getDeepValueFallback(keys ...string) interface{} {
	for i := range keys {
		v := o.getDeepValue(keys[i])
		if v != nil {
			return v
		}
	}
	return nil
}

var filenameKeys = []string{
	"metadata-fields.metus.general.former-file-name",
	"metadata-fields.metus.file.file-name",
	"metadata-fields.metus.general.file-name",
}

func (o *Object) getPhysicalFilename() string {
	v := o.getDeepValue("metadata-fields.metus.archive.former-path")
	if v != nil {
		s := strings.Split(v.(string), "\\")
		return s[len(s)-1]
	}

	return o.getDeepValueFallback(filenameKeys...).(string)
}

func (o *Object) getPhysicalFilepath() string {
	v := o.getDeepValue("metadata-fields.metus.file.file-local-path")
	if v != nil {
		s := strings.Split(v.(string), "\\")

		for i := range s {
			if s[i] == "4dcc9fad-c769-45c1-ae46-a7bc93764af8" {
				s[i+1] = strings.ToUpper(s[i+1])
				break
			}
		}

		return strings.Join(s[1:], "/")
	}
	return ""
}

//
//func (o *Object) isDeleted() bool {
//	v := o.getDeepValue("metadata-fields.metus.general.general-metus-not-used.deleted-date")
//	_, ok := v.(float64)
//	return ok
//}
