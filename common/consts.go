package common

const (
	// Collection Types
	CT_DAILY_LESSON       = "DAILY_LESSON"
	CT_SPECIAL_LESSON     = "SPECIAL_LESSON"
	CT_FRIENDS_GATHERINGS = "FRIENDS_GATHERINGS"
	CT_CONGRESS           = "CONGRESS"
	CT_VIDEO_PROGRAM      = "VIDEO_PROGRAM"
	CT_LECTURE_SERIES     = "LECTURE_SERIES"
	CT_VIRTUAL_LESSONS    = "VIRTUAL_LESSONS"
	CT_CHILDREN_LESSONS   = "CHILDREN_LESSONS"
	CT_WOMEN_LESSONS      = "WOMEN_LESSONS"
	CT_MEALS              = "MEALS"
	CT_HOLIDAY            = "HOLIDAY"
	CT_PICNIC             = "PICNIC"
	CT_UNITY_DAY          = "UNITY_DAY"
	CT_CLIPS              = "CLIPS"
	CT_ARTICLES           = "ARTICLES"
	CT_LESSONS_SERIES     = "LESSONS_SERIES"
	CT_SONGS              = "SONGS"
	CT_BOOKS              = "BOOKS"
	CT_PUBLIC_EVENTS      = "PUBLIC_EVENTS"

	// Content Unit Types
	CT_LESSON_PART           = "LESSON_PART"
	CT_LECTURE               = "LECTURE"
	CT_VIRTUAL_LESSON        = "VIRTUAL_LESSON"
	CT_CHILDREN_LESSON       = "CHILDREN_LESSON"
	CT_WOMEN_LESSON          = "WOMEN_LESSON"
	CT_FRIENDS_GATHERING     = "FRIENDS_GATHERING"
	CT_MEAL                  = "MEAL"
	CT_VIDEO_PROGRAM_CHAPTER = "VIDEO_PROGRAM_CHAPTER"
	CT_FULL_LESSON           = "FULL_LESSON"
	CT_ARTICLE               = "ARTICLE"
	CT_EVENT_PART            = "EVENT_PART"
	CT_UNKNOWN               = "UNKNOWN"
	CT_CLIP                  = "CLIP"
	CT_TRAINING              = "TRAINING"
	CT_KITEI_MAKOR           = "KITEI_MAKOR"
	CT_PUBLICATION           = "PUBLICATION"
	CT_LELO_MIKUD            = "LELO_MIKUD"
	CT_SONG                  = "SONG"
	CT_BOOK                  = "BOOK"
	CT_BLOG_POST             = "BLOG_POST"
	CT_RESEARCH_MATERIAL     = "RESEARCH_MATERIAL"
	CT_KTAIM_NIVCHARIM       = "KTAIM_NIVCHARIM"
	CT_SOURCE                = "SOURCE"
	CT_LIKUTIM               = "LIKUTIM"

	// Operation Types
	OP_CAPTURE_START = "capture_start"
	OP_CAPTURE_STOP  = "capture_stop"
	OP_DEMUX         = "demux"
	OP_TRIM          = "trim"
	OP_SEND          = "send"
	OP_CONVERT       = "convert"
	OP_UPLOAD        = "upload"
	OP_IMPORT_KMEDIA = "import_kmedia"
	OP_SIRTUTIM      = "sirtutim"
	OP_INSERT        = "insert"
	OP_TRANSCODE     = "transcode"
	OP_JOIN          = "join"
	OP_REPLACE       = "replace"

	// Source Types
	SRC_COLLECTION = "COLLECTION"
	SRC_BOOK       = "BOOK"
	SRC_VOLUME     = "VOLUME"
	SRC_PART       = "PART"
	SRC_PARASHA    = "PARASHA"
	SRC_CHAPTER    = "CHAPTER"
	SRC_ARTICLE    = "ARTICLE"
	SRC_TITLE      = "TITLE"
	SRC_LETTER     = "LETTER"
	SRC_ITEM       = "ITEM"

	// Content Role types
	CR_LECTURER = "LECTURER"

	// Persons patterns
	P_RAV = "rav"

	// Security levels
	SEC_PUBLIC    = int16(0)
	SEC_SENSITIVE = int16(1)
	SEC_PRIVATE   = int16(2)

	// Permissions
	PERM_READ             = "read"
	PERM_WRITE            = "write"
	PERM_I18N_WRITE       = "i18n_write"
	PERM_METADATA_WRITE   = "metadata_write"
	PERM_LABEL_WRITE      = "label_write"
	PERM_LABEL_I18N_WRITE = "label_i18n_write"
	PERM_LABEL_READ       = "label_read"
	PERM_LABEL_MODERATE   = "label_moderate"

	// Approve state levels
	APR_NONE     = int16(0)
	APR_APPROVED = int16(1)
	APR_DECLINED = int16(2)

	// Languages
	LANG_ENGLISH    = "en"
	LANG_HEBREW     = "he"
	LANG_RUSSIAN    = "ru"
	LANG_SPANISH    = "es"
	LANG_ITALIAN    = "it"
	LANG_GERMAN     = "de"
	LANG_DUTCH      = "nl"
	LANG_FRENCH     = "fr"
	LANG_PORTUGUESE = "pt"
	LANG_TURKISH    = "tr"
	LANG_POLISH     = "pl"
	LANG_ARABIC     = "ar"
	LANG_HUNGARIAN  = "hu"
	LANG_FINNISH    = "fi"
	LANG_LITHUANIAN = "lt"
	LANG_JAPANESE   = "ja"
	LANG_BULGARIAN  = "bg"
	LANG_GEORGIAN   = "ka"
	LANG_NORWEGIAN  = "no"
	LANG_SWEDISH    = "sv"
	LANG_CROATIAN   = "hr"
	LANG_CHINESE    = "zh"
	LANG_PERSIAN    = "fa"
	LANG_ROMANIAN   = "ro"
	LANG_HINDI      = "hi"
	LANG_UKRAINIAN  = "ua"
	LANG_MACEDONIAN = "mk"
	LANG_LATVIAN    = "lv"
	LANG_AMHARIC    = "am"
	LANG_INDONESIAN = "id"
	LANG_ARMENIAN   = "hy"
	LANG_MULTI      = "zz"
	LANG_UNKNOWN    = "xx"
	LANG_ORIGINAL   = "or"
	LANG_SLOVENIAN  = "sl"
	LANG_CZECH      = "cs"
)
