package affiliation

import "reflect"

// Enemy represents Enemy controlled units, trying to force the player to lose.
type Enemy struct {
}

// IsFriendsWith returns true if the two affiliations are friendly.
func (a *Enemy) IsFriendsWith(other Interface) bool {
	if reflect.TypeOf(other).String() == "*affiliation.Enemy" {
		return true
	}

	return false
}

// IsFoesWith returns true if the two affiliations are enemies.
func (a *Enemy) IsFoesWith(other Interface) bool {
	if reflect.TypeOf(other).String() == "*affiliation.Enemy" {
		return false
	}

	return true
}

// Name returns the name of this logic block.
func (a *Enemy) Name() string {
	return "enemy"
}
