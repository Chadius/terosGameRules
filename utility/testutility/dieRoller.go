package testutility

// AlwaysMissDieRoller generates rolls that are guaranteed to miss.
type AlwaysMissDieRoller struct{}

// RollTwoDice rolls two dice so the attacker will certainly miss.
func (a AlwaysMissDieRoller) RollTwoDice() (int, int) {
	return -999, 999
}

// AlwaysHitDieRoller generates rolls that are guaranteed to hit.
type AlwaysHitDieRoller struct{}

// RollTwoDice rolls two dice so the attacker will certainly hit.
func (a AlwaysHitDieRoller) RollTwoDice() (int, int) {
	return 999, -999
}

// ReplayDiceRoller will replay the list of rolls given as requested.
//   Throws an error if it's out of dice to replay.
type ReplayDiceRoller struct {
	RollHistory [][]int
}

// RollTwoDice consumes the 0th record in the history and returns it as a roll.
func (r ReplayDiceRoller) RollTwoDice() (int, int) {
	currentRoll := r.RollHistory[0]
	r.RollHistory = r.RollHistory[0:]
	return currentRoll[0], currentRoll[1]
}
