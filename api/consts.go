package api

const (
	// Collection Types

	CT_DAILY_LESSON             = "DAILY_LESSON"
	CT_SATURDAY_LESSON          = "SATURDAY_LESSON"
	CT_WEEKLY_FRIENDS_GATHERING = "WEEKLY_FRIENDS_GATHERING"
	CT_CONGRESS                 = "CONGRESS"
	CT_VIDEO_PROGRAM            = "VIDEO_PROGRAM"
	CT_LECTURE_SERIES           = "LECTURE_SERIES"
	CT_MEALS                    = "MEALS"
	CT_HOLIDAY                  = "HOLIDAY"
	CT_PICNIC                   = "PICNIC"
	CT_UNITY_DAY                = "UNITY_DAY"

	// Content Unit Types

	CT_LESSON_PART           = "LESSON_PART"
	CT_LECTURE               = "LECTURE"
	CT_CHILDREN_LESSON_PART  = "CHILDREN_LESSON_PART"
	CT_WOMEN_LESSON_PART     = "WOMEN_LESSON_PART"
	CT_CAMPUS_LESSON         = "CAMPUS_LESSON"
	CT_LC_LESSON             = "LC_LESSON"
	CT_VIRTUAL_LESSON        = "VIRTUAL_LESSON"
	CT_FRIENDS_GATHERING     = "FRIENDS_GATHERING"
	CT_MEAL                  = "MEAL"
	CT_VIDEO_PROGRAM_CHAPTER = "VIDEO_PROGRAM_CHAPTER"
	CT_FULL_LESSON           = "FULL_LESSON"
	CT_TEXT                  = "TEXT"

	// Operation Types

	OP_CAPTURE_START = "capture_start"
	OP_CAPTURE_STOP  = "capture_stop"
	OP_DEMUX         = "demux"
	OP_SEND          = "send"
	OP_UPLOAD        = "upload"
)

var DEFAULT_NAMES = map[string]string{
	CT_DAILY_LESSON:             "Daily Lesson",
	CT_SATURDAY_LESSON:          "Saturday Lesson",
	CT_WEEKLY_FRIENDS_GATHERING: "Weekly Friends Gathering",
	CT_CONGRESS:                 "Congress",
	CT_VIDEO_PROGRAM:            "Video Program",
	CT_LECTURE_SERIES:           "Lecture Series",
	CT_MEALS:                    "Meals",
	CT_HOLIDAY:                  "Holiday",
	CT_PICNIC:                   "Picnic",
	CT_UNITY_DAY:                "Unity Day",

	CT_LESSON_PART:           "Morning Lesson",
	CT_LECTURE:               "Lecture",
	CT_CHILDREN_LESSON_PART:  "Children Lesson",
	CT_WOMEN_LESSON_PART:     "Women Lesson",
	CT_CAMPUS_LESSON:         "Campus Lesson",
	CT_LC_LESSON:             "Learning Center  Lesson",
	CT_VIRTUAL_LESSON:        "Virtual Lesson",
	CT_FRIENDS_GATHERING:     "Friends Gathering",
	CT_MEAL:                  "Meal",
	CT_VIDEO_PROGRAM_CHAPTER: "Video Program",
	CT_FULL_LESSON:           "Full Lesson",
	CT_TEXT:                  "text",
}
