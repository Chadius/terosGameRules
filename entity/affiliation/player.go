package affiliation

import "reflect"

// Player represents player controlled units, trying to win.
type Player struct {
}

// IsFriendsWith returns true if the two affiliations are friendly.
func (a *Player) IsFriendsWith(other Interface) bool {
	if reflect.TypeOf(other).String() == "*affiliation.Player" ||
		reflect.TypeOf(other).String() == "*affiliation.Ally" {
		return true
	}

	return false
}

// IsFoesWith returns true if the two affiliations are enemies.
func (a *Player) IsFoesWith(other Interface) bool {
	if reflect.TypeOf(other).String() == "*affiliation.Player" ||
		reflect.TypeOf(other).String() == "*affiliation.Ally" {
		return false
	}

	return true
}

// Name returns the name of this logic block.
func (a *Player) Name() string {
	return "player"
}
