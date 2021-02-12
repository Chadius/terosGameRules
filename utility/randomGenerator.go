package utility

import "math/rand"

// RandomInt returns a random integer from 0 to maxInt.
func RandomInt(maxInt int) int {
	return rand.Intn(maxInt)
}
