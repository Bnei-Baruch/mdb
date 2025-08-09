package common

import (
	"math/rand"
	"strings"
)

const uidBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const lettersBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateUID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = uidBytes[rand.Intn(len(uidBytes))]
	}
	return string(b)
}

// Return standard language or LANG_UNKNOWN
//
//	if len(lang) = 2 we assume it's an MDB language code and check KNOWN_LANGS.
//	if len(lang) = 3 we assume it's a workflow / kmedia lang code and check LANG_MAP.
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

func GenerateName(n int) string {
	b := make([]byte, n)
	b[0] = lettersBytes[rand.Intn(len(lettersBytes))]
	for i := range b[1:] {
		b[i+1] = uidBytes[rand.Intn(len(uidBytes))]
	}
	return string(b)
}
