package affiliation

import "reflect"

// Ally represents Ally controlled units, controlled by the computer to help the player.
type Ally struct {
}

// IsFriendsWith returns true if the two affiliations are friendly.
func (a *Ally) IsFriendsWith(other Interface) bool {
	if reflect.TypeOf(other).String() == "*affiliation.Ally" ||
		reflect.TypeOf(other).String() == "*affiliation.Player" {
		return true
	}

	return false
}

// IsFoesWith returns true if the two affiliations are enemies.
func (a *Ally) IsFoesWith(other Interface) bool {
	if reflect.TypeOf(other).String() == "*affiliation.Ally" ||
		reflect.TypeOf(other).String() == "*affiliation.Player" {
		return false
	}

	return true
}
