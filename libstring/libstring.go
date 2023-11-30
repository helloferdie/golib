package libstring

import (
	"math/rand"
	"strings"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

const numset = "0123456789"
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" + numset

// Ucfirst - Upper case first letter
func Ucfirst(s string) string {
	return strings.ToUpper(s[0:1]) + strings.ToLower(s[1:])
}

// RandWithCharset - Return random string with charset
func RandWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Rand - Return random string
func Rand(length int) string {
	return RandWithCharset(length, charset)
}

// RandNum - Return random numeric string
func RandNum(length int) string {
	return RandWithCharset(length, numset)
}
