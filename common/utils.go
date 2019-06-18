package common

import "strings"

// Return standard language or LANG_UNKNOWN
//
// 	if len(lang) = 2 we assume it's an MDB language code and check KNOWN_LANGS.
// 	if len(lang) = 3 we assume it's a workflow / kmedia lang code and check LANG_MAP.
func StdLang(lang string) string {
	switch len(lang) {
	case 2:
		if l := strings.ToLower(lang); KNOWN_LANGS.MatchString(l) {
			return l
		}
	case 3:
		if l, ok := LANG_MAP[strings.ToUpper(lang)]; ok {
			return l
		}
	}

	return LANG_UNKNOWN
}

