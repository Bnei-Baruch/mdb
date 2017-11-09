package utils

import (
	"math/rand"
	"time"
	"regexp"
	"unicode/utf8"
)

const uidBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const lettersBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var SHA1_RE = regexp.MustCompile("^[0-9a-f]{40}$")

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func GenerateUID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = uidBytes[rand.Intn(len(uidBytes))]
	}
	return string(b)
}

func GenerateName(n int) string {
	b := make([]byte, n)
	b[0] = lettersBytes[rand.Intn(len(lettersBytes))]
	for i := range b[1:] {
		b[i+1] = uidBytes[rand.Intn(len(uidBytes))]
	}
	return string(b)
}

// true if every string in given slice is empty
func IsEmpty(s []string) bool {
	for _, x := range s {
		if x != "" {
			return false
		}
	}
	return true
}

// Like math.Min for int
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func ConvertArgsInt(args ...int) []interface{} {
	c := make([]interface{}, len(args))
	for i, a := range args {
		c[i] = a
	}
	return c
}

func ConvertArgsInt64(args []int64) []interface{} {
	c := make([]interface{}, len(args))
	for i, a := range args {
		c[i] = a
	}
	return c
}

func ConvertArgsString(args []string) []interface{} {
	c := make([]interface{}, len(args))
	for i, a := range args {
		c[i] = a
	}
	return c
}

func ConvertArgsBytes(args [][]byte) []interface{} {
	c := make([]interface{}, len(args))
	for i, a := range args {
		c[i] = a
	}
	return c
}

// Taken AS IS from
// https://stackoverflow.com/a/34521190
// Note that this implementation DOES NOT handle combining marks correctly
func Reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}