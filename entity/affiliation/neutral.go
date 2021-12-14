package affiliation

// Neutral represents Neutral controlled units. They are unaffiliated and are foes with everything, including other Neutrals.
type Neutral struct {
}

// IsFriendsWith returns false. Neutral units have no allies.
func (a *Neutral) IsFriendsWith(other Interface) bool {
	return false
}

// IsFoesWith returns true. Neutral units oppose everything.
func (a *Neutral) IsFoesWith(other Interface) bool {
	return true
}
