package importer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Bnei-Baruch/mdb/common"
)

var LANG_RE *regexp.Regexp
var OT_RE = regexp.MustCompile("(?i)^[ot]$")
var RAV_NORAV_RE = regexp.MustCompile("(?i)^(rav|norav)$")
var BITRATE_RE = regexp.MustCompile("(?i)^(24k|96k|128k|hd)$")
var FILMDATE_RE = regexp.MustCompile("^((19|20)\\d\\d)-(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])$") // YYYY-MM-DD

func init() {
	i := 0
	keys := make([]string, len(common.LANG_MAP)-1)
	for k := range common.LANG_MAP {
		if k == "" {
			continue
		}
		keys[i] = k
		i++
	}
	LANG_RE = regexp.MustCompile(fmt.Sprintf("(?i)^(%s)$", strings.Join(keys, "|")))
}

type Line struct {
	Language string
	OT       string
	Rav      string
	FilmDate string
	Bitrate  string
	Format   string
}

func ParseLine(name string) Line {
	line := Line{}

	if idx := strings.LastIndex(name, "."); idx > 0 {
		line.Format = name[idx+1:]
		name = name[:idx]
	}

	parts := strings.Split(name, "_")
	for i := range parts {
		part := strings.TrimSpace(parts[i])
		if v := OT_RE.FindString(part); v != "" {
			line.OT = v
		} else if v := BITRATE_RE.FindString(part); v != "" {
			line.Bitrate = v
		} else if v := RAV_NORAV_RE.FindString(part); v != "" {
			line.Rav = v
		} else if v := LANG_RE.FindString(part); v != "" {
			line.Language = v
		} else if v := FILMDATE_RE.FindString(part); v != "" {
			line.FilmDate = v
		}
	}

	return line
}

func NormalizedFileName(name string) string {
	if idx := strings.LastIndex(name, "."); idx > 0 {
		name = name[:idx]
	}

	parts := strings.Split(name, "_")
	for i := range parts {
		part := strings.TrimSpace(parts[i])
		if OT_RE.MatchString(part) {
			parts[i] = "ot"
		} else if BITRATE_RE.MatchString(part) {
			parts[i] = "br"
		} else if LANG_RE.MatchString(part) {
			parts[i] = "lang"
		}
	}

	name = strings.Join(parts, "_")
	return strings.Replace(name, "_br", "", -1)
}
