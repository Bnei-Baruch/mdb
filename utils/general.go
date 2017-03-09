package utils

import (
	"math/rand"
	"time"
)

const uidBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const lettersBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// panic if err != nil
func Must(err error) {
	if err != nil {
		panic(err)
	}
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
