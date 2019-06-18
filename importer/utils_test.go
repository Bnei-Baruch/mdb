package importer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	line := ParseLine("heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4")
	assert.Equal(t, "mp4", line.Format, "Format")
	assert.Equal(t, "heb", line.Language, "Language")
	assert.Equal(t, "o", line.OT, "OT")
	assert.Equal(t, "2016-09-14", line.FilmDate, "FilmDate")
	assert.Equal(t, "", line.Bitrate, "Bitrate")

	line = ParseLine("rus_t_rav_rb-1990-02-kishalon_1995-12-30_lesson_96k.mp3")
	assert.Equal(t, "mp3", line.Format, "Format")
	assert.Equal(t, "rus", line.Language, "Language")
	assert.Equal(t, "t", line.OT, "OT")
	assert.Equal(t, "1995-12-30", line.FilmDate, "FilmDate")
	assert.Equal(t, "96k", line.Bitrate, "Bitrate")

	line = ParseLine("bad_date_2019_13_01")
	assert.Equal(t, "", line.FilmDate, "FilmDate")
	line = ParseLine("bad_date_2019_12_32")
	assert.Equal(t, "", line.FilmDate, "FilmDate")

}
