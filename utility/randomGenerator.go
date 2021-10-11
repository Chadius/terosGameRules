package utility

import "math/rand"

// SixSideGenerator is an interface that generates numbers from 1 to 6.
type SixSideGenerator interface {
	RollTwoDice() (int, int)
}

// RandomInt returns a random integer from 0 to maxInt.
func RandomInt(maxInt int) int {
	return rand.Intn(maxInt)
}

// RandomDieRoller generates a random number from 1 to 6,
//   to simulate the common six sided die.
type RandomDieRoller struct{}

// RollTwoDice rolls two dice.
func (r RandomDieRoller) RollTwoDice() (int, int) {
	return 1 + RandomInt(5), 1 + RandomInt(5)
}
