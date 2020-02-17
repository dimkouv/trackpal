package cryptoutils

import (
	"math/rand"
	"strings"
	"time"
)

// RandomString generates a random string of fixed length by using only chars that reside in choices array
func RandomString(length int, choices []rune) string {
	rand.Seed(time.Now().UnixNano())

	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(choices[rand.Intn(len(choices))])
	}

	return b.String()
}
