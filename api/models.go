package api

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/nullbio/null.v6"

	"github.com/Bnei-Baruch/mdb/models"
)

type (
	// Common

	Operation struct {
		Station    string `json:"station" binding:"required"`
		User       string `json:"user" binding:"required,email"`
		WorkflowID string `json:"workflow_id"`
	}

	File struct {
		FileName  string     `json:"file_name" binding:"required,max=255"`
		Sha1      string     `json:"sha1" binding:"required,len=40,hexadecimal"`
		Size      int64      `json:"size" binding:"required"`
		CreatedAt *Timestamp `json:"created_at" binding:"required"`
		Type      string     `json:"type" binding:"max=16"`
		SubType   string     `json:"sub_type" binding:"max=16"`
		MimeType  string     `json:"mime_type" binding:"max=255"`
		Language  string     `json:"language" binding:"omitempty,min=2,max=3"`
	}

	AVFile struct {
		File
		Duration float64 `json:"duration"`
	}

	CITMetadataMajor struct {
		Type string `json:"type" binding:"required,eq=source|eq=tag"`
		Idx  int    `json:"idx" binding:"gte=0"`
	}
	CITMetadata struct {
		ContentType    string           `json:"content_type" binding:"required"`
		CaptureDate    Date             `json:"capture_date" binding:"required"`
		FinalName      string           `json:"final_name" binding:"required,max=255"`
		Language       string           `json:"language" binding:"required,min=2,max=3"`
		Lecturer       string           `json:"lecturer" binding:"required"`
		AutoName       string           `json:"auto_name"`
		ManualName     string           `json:"manual_name"`
		WeekDate       *Date            `json:"week_date"`
		Number         null.Int         `json:"number"`
		Part           null.Int         `json:"part"`
		Sources        []string         `json:"sources" binding:"omitempty,dive,len=8"`
		Tags           []string         `json:"tags" binding:"omitempty,dive,len=8"`
		ArtifactType   null.String      `json:"artifact_type"`
		HasTranslation bool             `json:"has_translation"`
		RequireTest    bool             `json:"require_test"`
		CollectionUID  null.String      `json:"collection_uid" binding:"omitempty,len=8"`
		Episode        null.String      `json:"episode"`
		PartType       null.Int         `json:"part_type"`
		Major          *CITMetadataMajor `json:"major" binding:"omitempty"`
	}

	Rename struct {
		Sha1     string `json:"sha1" binding:"required,len=40,hexadecimal"`
		FileName string `json:"file_name" binding:"required,max=255"`
	}

	// Operations

	CaptureStartRequest struct {
		Operation
		FileName      string `json:"file_name" binding:"max=255"`
		CaptureSource string `json:"capture_source"`
		CollectionUID string `json:"collection_uid"`
	}

	CaptureStopRequest struct {
		Operation
		File
		CaptureSource string `json:"capture_source"`
		CollectionUID string `json:"collection_uid"`
		Part          string `json:"part"`
	}

	DemuxRequest struct {
		Operation
		Sha1          string `json:"sha1" binding:"required,len=40,hexadecimal"`
		Original      AVFile `json:"original"`
		Proxy         AVFile `json:"proxy"`
		CaptureSource string `json:"capture_source"`
	}

	TrimRequest struct {
		Operation
		OriginalSha1  string    `json:"original_sha1" binding:"required,len=40,hexadecimal"`
		ProxySha1     string    `json:"proxy_sha1" binding:"required,len=40,hexadecimal"`
		Original      AVFile    `json:"original"`
		Proxy         AVFile    `json:"proxy"`
		In            []float64 `json:"in"`
		Out           []float64 `json:"out"`
		CaptureSource string    `json:"capture_source"`
	}

	SendRequest struct {
		Operation
		Original Rename      `json:"original"`
		Proxy    Rename      `json:"proxy"`
		Metadata CITMetadata `json:"metadata"`
	}

	ConvertRequest struct {
		Operation
		Sha1   string   `json:"sha1" binding:"required,len=40,hexadecimal"`
		Output []AVFile `json:"output"`
	}

	UploadRequest struct {
		Operation
		AVFile
		Url string `json:"url" binding:"required"`
	}

	// REST

	ListRequest struct {
		PageNumber int    `json:"page_no" form:"page_no" binding:"omitempty,min=1"`
		PageSize   int    `json:"page_size" form:"page_size" binding:"omitempty,min=1"`
		StartIndex int    `json:"start_index" form:"start_index" binding:"omitempty,min=1"`
		StopIndex  int    `json:"stop_index" form:"stop_index" binding:"omitempty,min=1"`
		OrderBy    string `json:"order_by" form:"order_by" binding:"omitempty"`
	}

	ListResponse struct {
		Total int64 `json:"total"`
	}

	ContentTypesFilter struct {
		ContentTypes []string `json:"content_types" form:"content_type" binding:"omitempty"`
	}

	OperationTypesFilter struct {
		OperationTypes []string `json:"operation_types" form:"operation_type" binding:"omitempty"`
	}

	SourcesFilter struct {
		Authors []string `json:"authors" form:"author" binding:"omitempty"`
		Sources []int64  `json:"sources" form:"source" binding:"omitempty"`
	}

	TagsFilter struct {
		Tags []int64 `json:"tags" form:"tag" binding:"omitempty"`
	}

	SearchTermFilter struct {
		Query string `json:"query" form:"query" binding:"omitempty"`
	}

	DateRangeFilter struct {
		StartDate string `json:"start_date" form:"start_date" binding:"omitempty"`
		EndDate   string `json:"end_date" form:"end_date" binding:"omitempty"`
	}

	SecureFilter struct {
		Levels []int16 `json:"security_levels" form:"secure" binding:"omitempty"`
	}

	PublishedFilter struct {
		Published string `json:"published" form:"published" binding:"omitempty"`
	}

	CollectionsRequest struct {
		ListRequest
		ContentTypesFilter
		DateRangeFilter
		SecureFilter
		PublishedFilter
	}

	CollectionsResponse struct {
		ListResponse
		Collections []*Collection `json:"data"`
	}

	ContentUnitsRequest struct {
		ListRequest
		ContentTypesFilter
		DateRangeFilter
		SecureFilter
		PublishedFilter
		SourcesFilter
		TagsFilter
	}

	ContentUnitsResponse struct {
		ListResponse
		ContentUnits []*ContentUnit `json:"data"`
	}

	FilesRequest struct {
		ListRequest
		DateRangeFilter
		SecureFilter
		PublishedFilter
		SearchTermFilter
	}

	FilesResponse struct {
		ListResponse
		Files []*MFile `json:"data"`
	}

	OperationsRequest struct {
		ListRequest
		DateRangeFilter
		OperationTypesFilter
	}

	OperationsResponse struct {
		ListResponse
		Operations []*models.Operation `json:"data"`
	}

	AuthorsResponse struct {
		ListResponse
		Authors []*Author `json:"data"`
	}

	SourcesRequest struct {
		ListRequest
	}

	SourcesResponse struct {
		ListResponse
		Sources []*Source `json:"data"`
	}

	CreateSourceRequest struct {
		Source
		AuthorID null.Int64 `json:"author"`
	}

	TagsRequest struct {
		ListRequest
	}

	TagsResponse struct {
		ListResponse
		Tags []*Tag `json:"data"`
	}

	HierarchyRequest struct {
		Language string `json:"language" form:"language" binding:"omitempty,len=2"`
		RootUID  string `json:"root" form:"root" binding:"omitempty,len=8"`
		Depth    int    `json:"depth" form:"depth"`
	}

	SourcesHierarchyRequest struct {
		HierarchyRequest
	}

	TagsHierarchyRequest struct {
		HierarchyRequest
	}

	Collection struct {
		models.Collection
		I18n map[string]*models.CollectionI18n `json:"i18n"`
	}

	ContentUnit struct {
		models.ContentUnit
		I18n map[string]*models.ContentUnitI18n `json:"i18n"`
	}

	CollectionContentUnit struct {
		Collection  *Collection  `json:"collection,omitempty"`
		ContentUnit *ContentUnit `json:"content_unit,omitempty"`
		Name        string       `json:"name"`
	}

	// Marshalable File
	MFile struct {
		models.File
		Sha1Str string `json:"sha1"`
	}

	Author struct {
		models.Author
		I18n    map[string]*models.AuthorI18n `json:"i18n"`
		Sources []*Source                     `json:"sources"`
	}

	Source struct {
		models.Source
		I18n map[string]*models.SourceI18n `json:"i18n"`
	}

	Tag struct {
		models.Tag
		I18n map[string]*models.TagI18n `json:"i18n"`
	}

	SourceH struct {
		ID          int64       `json:"id"`
		UID         string      `json:"uid"`
		ParentID    null.Int64  `json:"parent_id"`
		Type        string      `json:"type"`
		Position    null.Int    `json:"position"`
		Pattern     null.String `json:"pattern,omitempty"`
		Name        null.String `json:"name"`
		Description null.String `json:"description,omitempty"`
		Children    []*SourceH  `json:"children,omitempty"`
	}

	AuthorH struct {
		Code     string      `json:"code"`
		Name     string      `json:"name"`
		FullName null.String `json:"full_name,omitempty"`
		Children []*SourceH  `json:"children,omitempty"`
	}

	TagH struct {
		ID       int64       `json:"id"`
		UID      string      `json:"uid"`
		ParentID null.Int64  `json:"parent_id"`
		Pattern  null.String `json:"pattern,omitempty"`
		Label    null.String `json:"label"`
		Children []*TagH     `json:"children,omitempty"`
	}
)

func NewCollectionsResponse() *CollectionsResponse {
	return &CollectionsResponse{Collections: make([]*Collection, 0)}
}

func NewContentUnitsResponse() *ContentUnitsResponse {
	return &ContentUnitsResponse{ContentUnits: make([]*ContentUnit, 0)}
}

func NewFilesResponse() *FilesResponse {
	return &FilesResponse{Files: make([]*MFile, 0)}
}

func NewMFile(f *models.File) *MFile {
	x := &MFile{File: *f}
	if f.Sha1.Valid {
		x.Sha1Str = hex.EncodeToString(f.Sha1.Bytes)
	}
	return x
}

func NewOperationsResponse() *OperationsResponse {
	return &OperationsResponse{Operations: make([]*models.Operation, 0)}
}

func NewSourcesResponse() *SourcesResponse {
	return &SourcesResponse{Sources: make([]*Source, 0)}
}

func NewTagsResponse() *TagsResponse {
	return &TagsResponse{Tags: make([]*Tag, 0)}
}

func (drf *DateRangeFilter) Range() (time.Time, time.Time, error) {
	var err error
	var s, e time.Time

	if drf.StartDate != "" {
		s, err = time.Parse("2006-01-02", drf.StartDate)
	}
	if err == nil && drf.EndDate != "" {
		e, err = time.Parse("2006-01-02", drf.EndDate)
	}

	return s, e, err
}

// A time.Time like structure with Unix timestamp JSON marshalling
type Timestamp struct {
	time.Time
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	ts := t.Time.Unix()
	stamp := fmt.Sprint(ts)
	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	t.Time = time.Unix(int64(ts), 0)
	return nil
}

// A time.Time like structure with date part only JSON marshalling
type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.Time.Format("2006-01-02"))), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	var err error
	d.Time, err = time.Parse("2006-01-02", strings.Trim(string(b), "\""))
	if err != nil {
		fmt.Println(err)
	}
	return err
}
