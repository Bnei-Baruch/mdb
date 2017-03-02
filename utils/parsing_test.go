package utils

import (
	"testing"
	"strings"
	"fmt"
)

func TestParseFilename(t *testing.T) {
	fn, err := ParseFileName("heb_o_rav_rb-1990-02-kishalon_2016-09-14_lesson.mp4")
	if err != nil {
		t.Error("ParseFileName should succeed.")
	}
	if fn.Type != "mp4" {
		t.Error("Expected type to be mp4")
	}
	if fn.Part != "rb-1990-02-kishalon" {
		t.Errorf("Expected different part %s", fn.Part)
	}

	_, err = ParseFileName("heb_o_rav_rb-1990-02-kishalon_201-09-14_lesson.mp4")
	if e := "could not parse date"; err == nil || !strings.Contains(err.Error(), e) {
		t.Error(fmt.Sprintf("ParseFileName should contain %s, got %s.", e, err))
	}

	// Make sure code does not crash.
	ParseFileName("2017-01-04_02-40-19")
}
