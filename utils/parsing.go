package utils

import (
	"path/filepath"
	"strings"
	"time"
	"github.com/pkg/errors"
)

type FileName struct {
	Name     string
	Base     string
	Type     string // File extension, mp3 or mp4 or other.
	Language string
	Rav      bool
	Part     string
	Date     time.Time
	DateStr  string
}

func ParseFileName(name string) (*FileName, error) {
	format := "Expected file name is [lang]_o_[rav/norav]_[part-a]_[2006-01-02]_[anyhing else].mp4"
	fn := FileName{
		Name: name,
		Base: filepath.Base(name),
		Type: strings.Replace(filepath.Ext(name), ".", "", 1),
	}
	parts := strings.Split(strings.TrimSuffix(fn.Base, filepath.Ext(fn.Base)), "_")
	if len(parts) < 4 {
		return nil, errors.Errorf("Bad filename, expected at least 4 parts, found %d: %s. %s", len(parts), parts, format)
	}
	fn.Language = parts[0]
	if parts[2] == "rav" {
		fn.Rav = true
	} else if parts[2] == "norav" {
		fn.Rav = false
	} else {
		return nil, errors.Errorf("Bad filename, expected rav/norav got %s. %s", parts[2], format)
	}
	fn.Part = parts[3]
	var err error
	fn.Date, err = time.Parse("2006-01-02", parts[4])
	if err != nil {
		return nil, errors.Errorf("Bad filename, could not parse date (%s): %s. %s", parts[4], err.Error(), format)
	}
	fn.DateStr = parts[4]

	return &fn, nil
}
