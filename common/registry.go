package common

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Bnei-Baruch/mdb/models"
)

var (
	CONTENT_TYPE_REGISTRY      = &ContentTypeRegistry{}
	OPERATION_TYPE_REGISTRY    = &OperationTypeRegistry{}
	CONTENT_ROLE_TYPE_REGISTRY = &ContentRoleTypeRegistry{}
	PERSON_REGISTRY            = &PersonRegistry{}
	AUTHOR_REGISTRY            = &AuthorRegistry{}
	SOURCE_TYPE_REGISTRY       = &SourceTypeRegistry{}
	MEDIA_TYPE_REGISTRY        = &MediaTypeRegistry{}
	TWITTER_USERS_REGISTRY     = &TwitterUsersRegistry{}

	ALL_LANGS = []string{
		LANG_UNKNOWN, LANG_MULTI, LANG_ENGLISH, LANG_HEBREW, LANG_RUSSIAN, LANG_SPANISH, LANG_ITALIAN,
		LANG_GERMAN, LANG_DUTCH, LANG_FRENCH, LANG_PORTUGUESE, LANG_TURKISH, LANG_POLISH, LANG_ARABIC,
		LANG_HUNGARIAN, LANG_FINNISH, LANG_LITHUANIAN, LANG_JAPANESE, LANG_BULGARIAN, LANG_GEORGIAN,
		LANG_NORWEGIAN, LANG_SWEDISH, LANG_CROATIAN, LANG_CHINESE, LANG_PERSIAN, LANG_ROMANIAN, LANG_HINDI,
		LANG_MACEDONIAN, LANG_SLOVENIAN, LANG_LATVIAN, LANG_CZECH, LANG_UKRAINIAN, LANG_AMHARIC,
		LANG_INDONESIAN, LANG_ARMENIAN, LANG_ORIGINAL,
	}

	KNOWN_LANGS = regexp.MustCompile(strings.Join(ALL_LANGS, "|"))

	// kmedia - select concat('"',code3,'": "',locale,'",') from languages;
	LANG_MAP = map[string]string{
		"":    LANG_UNKNOWN,
		"MLT": LANG_MULTI,
		"ENG": LANG_ENGLISH,
		"HEB": LANG_HEBREW,
		"RUS": LANG_RUSSIAN,
		"SPA": LANG_SPANISH,
		"ITA": LANG_ITALIAN,
		"GER": LANG_GERMAN,
		"DUT": LANG_DUTCH,
		"FRE": LANG_FRENCH,
		"POR": LANG_PORTUGUESE,
		"TRK": LANG_TURKISH,
		"TUR": LANG_TURKISH,
		"POL": LANG_POLISH,
		"ARB": LANG_ARABIC,
		"ARA": LANG_ARABIC,
		"HUN": LANG_HUNGARIAN,
		"FIN": LANG_FINNISH,
		"LIT": LANG_LITHUANIAN,
		"JPN": LANG_JAPANESE,
		"BUL": LANG_BULGARIAN,
		"GEO": LANG_GEORGIAN,
		"NOR": LANG_NORWEGIAN,
		"SWE": LANG_SWEDISH,
		"HRV": LANG_CROATIAN,
		"CHN": LANG_CHINESE,
		"CHI": LANG_CHINESE,
		"PER": LANG_PERSIAN,
		"RON": LANG_ROMANIAN,
		"HIN": LANG_HINDI,
		"MKD": LANG_MACEDONIAN,
		"LAV": LANG_LATVIAN,
		"UKR": LANG_UKRAINIAN,
		"AMH": LANG_AMHARIC,
		"IND": LANG_INDONESIAN,
		"ARM": LANG_ARMENIAN,
		"ORI": LANG_ORIGINAL,
		"SLV": LANG_CZECH,
		"CZE": LANG_CZECH,
	}

	ALL_CONTENT_TYPES = []string{
		CT_DAILY_LESSON, CT_SPECIAL_LESSON, CT_FRIENDS_GATHERINGS, CT_CONGRESS, CT_VIDEO_PROGRAM,
		CT_LECTURE_SERIES, CT_VIRTUAL_LESSONS, CT_CHILDREN_LESSONS, CT_WOMEN_LESSONS, CT_MEALS, CT_HOLIDAY, CT_PICNIC,
		CT_UNITY_DAY, CT_CLIPS, CT_ARTICLES, CT_LESSONS_SERIES, CT_SONGS, CT_BOOKS, CT_LESSON_PART, CT_LECTURE,
		CT_CHILDREN_LESSON, CT_WOMEN_LESSON, CT_VIRTUAL_LESSON, CT_FRIENDS_GATHERING, CT_MEAL, CT_VIDEO_PROGRAM_CHAPTER,
		CT_FULL_LESSON, CT_ARTICLE, CT_EVENT_PART, CT_UNKNOWN, CT_CLIP, CT_TRAINING, CT_KITEI_MAKOR, CT_PUBLICATION,
		CT_LELO_MIKUD, CT_SONG, CT_BOOK, CT_BLOG_POST, CT_RESEARCH_MATERIAL, CT_KTAIM_NIVCHARIM, CT_SOURCE,
	}

	ALL_OPERATION_TYPES = []string{
		OP_CAPTURE_START, OP_CAPTURE_STOP, OP_DEMUX, OP_TRIM, OP_SEND, OP_CONVERT, OP_UPLOAD, OP_IMPORT_KMEDIA,
		OP_SIRTUTIM, OP_INSERT, OP_TRANSCODE, OP_JOIN,
	}

	UNIT_CONTENT_TYPE_CAN_CHANGE = []string{
		CT_LESSON_PART,
		CT_LECTURE,
		CT_VIRTUAL_LESSON,
		CT_CHILDREN_LESSON,
		CT_WOMEN_LESSON,
		CT_FRIENDS_GATHERING,
		CT_MEAL,
		CT_VIDEO_PROGRAM_CHAPTER,
		CT_EVENT_PART,
		CT_UNKNOWN,
		CT_CLIP,
		CT_TRAINING,
		CT_LELO_MIKUD,
		CT_KTAIM_NIVCHARIM,
	}

	// Types of various, secondary, content slots in big events like congress, unity day, etc...
	// This list is not part of content_types to prevent explosion of that list.
	// This came to life for mdb-cit UI only Ease of Use. (prevent typing errors and keep consistency)
	// We keep it here so CCU's would have some information.
	// This list should be kept in sync with mdb-cit (consts.js)
	MISC_EVENT_PART_TYPES = [8]string{
		"TEKES_PTIHA",
		"TEKES_SIYUM",
		"EREV_PATUAH",
		"EREV_TARBUT",
		"ATZAGAT_PROEKT",
		"HAANAKAT_TEUDOT",
		"HATIMAT_SFARIM",
		"EVENT",
	}

	// kmedia - select asset_type, count(*) from file_assets group by asset_type order by count(*) desc;
	ALL_MEDIA_TYPES = []*MediaType{
		{Extension: "mp4", Type: "video", SubType: "", MimeType: "video/mp4"},
		{Extension: "wmv", Type: "video", SubType: "", MimeType: "video/x-ms-wmv"},
		{Extension: "flv", Type: "video", SubType: "", MimeType: "video/x-flv"},
		{Extension: "mov", Type: "video", SubType: "", MimeType: "video/quicktime"},
		{Extension: "asf", Type: "video", SubType: "", MimeType: "video/x-ms-asf"},
		{Extension: "mpg", Type: "video", SubType: "", MimeType: "video/mpeg"},
		{Extension: "avi", Type: "video", SubType: "", MimeType: "video/x-msvideo"},
		{Extension: "mp3", Type: "audio", SubType: "", MimeType: "audio/mpeg"},
		{Extension: "mp3", Type: "audio", SubType: "", MimeType: "audio/mp3"},
		{Extension: "wma", Type: "audio", SubType: "", MimeType: "audio/x-ms-wma"},
		{Extension: "mid", Type: "audio", SubType: "", MimeType: "audio/midi"},
		{Extension: "wav", Type: "audio", SubType: "", MimeType: "audio/x-wav"},
		{Extension: "aac", Type: "audio", SubType: "", MimeType: "audio/aac"},
		{Extension: "jpg", Type: "image", SubType: "", MimeType: "image/jpeg"},
		{Extension: "png", Type: "image", SubType: "", MimeType: "image/png"},
		{Extension: "gif", Type: "image", SubType: "", MimeType: "image/gif"},
		{Extension: "bmp", Type: "image", SubType: "", MimeType: "image/bmp"},
		{Extension: "tif", Type: "image", SubType: "", MimeType: "image/tiff"},
		{Extension: "zip", Type: "image", SubType: "", MimeType: "application/zip"},
		{Extension: "zip", Type: "image", SubType: "", MimeType: "application/x-zip-compressed"},
		{Extension: "7z", Type: "image", SubType: "", MimeType: "application/x-7z-compressed"},
		{Extension: "rar", Type: "image", SubType: "", MimeType: "application/x-rar-compressed"},
		{Extension: "sfk", Type: "image", SubType: "", MimeType: ""},
		{Extension: "doc", Type: "text", SubType: "", MimeType: "application/msword"},
		{Extension: "docx", Type: "text", SubType: "", MimeType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
		{Extension: "htm", Type: "text", SubType: "", MimeType: "text/html"},
		{Extension: "html", Type: "text", SubType: "", MimeType: "text/html"},
		{Extension: "pdf", Type: "text", SubType: "", MimeType: "application/pdf"},
		{Extension: "epub", Type: "text", SubType: "", MimeType: "application/epub+zip"},
		{Extension: "rtf", Type: "text", SubType: "", MimeType: "text/rtf"},
		{Extension: "txt", Type: "text", SubType: "", MimeType: "text/plain"},
		{Extension: "fb2", Type: "text", SubType: "", MimeType: "text/xml"},
		{Extension: "rb", Type: "text", SubType: "", MimeType: "application/x-rocketbook"},
		{Extension: "xls", Type: "sheet", SubType: "", MimeType: "application/vnd.ms-excel"},
		{Extension: "swf", Type: "banner", SubType: "", MimeType: "application/x-shockwave-flash"},
		{Extension: "ppt", Type: "presentation", SubType: "", MimeType: "application/vnd.ms-powerpoint"},
		{Extension: "pptx", Type: "presentation", SubType: "", MimeType: "application/vnd.openxmlformats-officedocument.presentationml.presentation"},
		{Extension: "pps", Type: "presentation", SubType: "", MimeType: "application/vnd.ms-powerpoint"},
		{Extension: "vtt", Type: "subtitles", SubType: "", MimeType: "text/vtt"},
	}
)

type MediaType struct {
	Extension string
	Type      string
	SubType   string
	MimeType  string
}

func InitTypeRegistries(exec boil.Executor) error {
	registries := []TypeRegistry{CONTENT_TYPE_REGISTRY,
		OPERATION_TYPE_REGISTRY,
		CONTENT_ROLE_TYPE_REGISTRY,
		PERSON_REGISTRY,
		AUTHOR_REGISTRY,
		SOURCE_TYPE_REGISTRY,
		MEDIA_TYPE_REGISTRY,
		TWITTER_USERS_REGISTRY,
	}

	for _, x := range registries {
		if err := x.Init(exec); err != nil {
			return err
		}
	}

	return nil
}

type TypeRegistry interface {
	Init(exec boil.Executor) error
}

type ContentTypeRegistry struct {
	ByName map[string]*models.ContentType
	ByID   map[int64]*models.ContentType
}

func (r *ContentTypeRegistry) Init(exec boil.Executor) error {
	types, err := models.ContentTypes().All(exec)
	if err != nil {
		return errors.Wrap(err, "Load content_types from DB")
	}

	r.ByName = make(map[string]*models.ContentType)
	r.ByID = make(map[int64]*models.ContentType)
	for _, t := range types {
		r.ByName[t.Name] = t
		r.ByID[t.ID] = t
	}

	return nil
}

type OperationTypeRegistry struct {
	ByName map[string]*models.OperationType
}

func (r *OperationTypeRegistry) Init(exec boil.Executor) error {
	types, err := models.OperationTypes().All(exec)
	if err != nil {
		return errors.Wrap(err, "Load operation_types from DB")
	}

	r.ByName = make(map[string]*models.OperationType)
	for _, t := range types {
		r.ByName[t.Name] = t
	}

	return nil
}

type ContentRoleTypeRegistry struct {
	ByName map[string]*models.ContentRoleType
}

func (r *ContentRoleTypeRegistry) Init(exec boil.Executor) error {
	types, err := models.ContentRoleTypes().All(exec)
	if err != nil {
		return errors.Wrap(err, "Load content_role_types from DB")
	}

	r.ByName = make(map[string]*models.ContentRoleType)
	for _, t := range types {
		r.ByName[t.Name] = t
	}

	return nil
}

type PersonRegistry struct {
	ByPattern map[string]*models.Person
}

func (r *PersonRegistry) Init(exec boil.Executor) error {
	types, err := models.Persons(qm.Where("pattern is not null")).All(exec)
	if err != nil {
		return errors.Wrap(err, "Load persons from DB")
	}

	r.ByPattern = make(map[string]*models.Person)
	for _, t := range types {
		r.ByPattern[t.Pattern.String] = t
	}

	return nil
}

type AuthorRegistry struct {
	ByCode map[string]*models.Author
}

func (r *AuthorRegistry) Init(exec boil.Executor) error {
	authors, err := models.Authors().All(exec)
	if err != nil {
		return errors.Wrap(err, "Load authors from DB")
	}

	r.ByCode = make(map[string]*models.Author)
	for _, a := range authors {
		r.ByCode[a.Code] = a
	}

	return nil
}

type SourceTypeRegistry struct {
	ByName map[string]*models.SourceType
	ByID   map[int64]*models.SourceType
}

func (r *SourceTypeRegistry) Init(exec boil.Executor) error {
	types, err := models.SourceTypes().All(exec)
	if err != nil {
		return errors.Wrap(err, "Load source_types from DB")
	}

	r.ByName = make(map[string]*models.SourceType)
	r.ByID = make(map[int64]*models.SourceType)
	for _, t := range types {
		r.ByName[t.Name] = t
		r.ByID[t.ID] = t
	}

	return nil
}

type MediaTypeRegistry struct {
	ByExtension map[string]*MediaType
	ByMime      map[string]*MediaType
}

func (r *MediaTypeRegistry) Init(exec boil.Executor) error {
	r.ByExtension = make(map[string]*MediaType, len(ALL_MEDIA_TYPES))
	r.ByMime = make(map[string]*MediaType, len(ALL_MEDIA_TYPES))

	for _, x := range ALL_MEDIA_TYPES {
		r.ByExtension[x.Extension] = x
		r.ByMime[x.MimeType] = x
	}

	return nil
}

type TwitterUsersRegistry struct {
	ByUsername map[string]*models.TwitterUser
	ByID       map[int64]*models.TwitterUser
}

func (r *TwitterUsersRegistry) Init(exec boil.Executor) error {
	users, err := models.TwitterUsers().All(exec)
	if err != nil {
		return errors.Wrap(err, "Load twitter users from DB")
	}

	r.ByUsername = make(map[string]*models.TwitterUser)
	r.ByID = make(map[int64]*models.TwitterUser)
	for _, t := range users {
		r.ByUsername[t.Username] = t
		r.ByID[t.ID] = t
	}

	return nil
}
