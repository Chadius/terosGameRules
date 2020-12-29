package terosbattleserver

import "math/rand"

// StringWithCharset generates a psuedo random string using the characters in charset of the desired length.
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
