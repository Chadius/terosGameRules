package levelupbenefit

import "github.com/chadius/terosbattleserver/entity/power"

// PowerChanges tracks changes to the squaddie's innate powers.
type PowerChanges struct {
	gained []*power.Reference
	lost   []*power.Reference
}

// NewPowerChanges returns a new PowerChanges object.
func NewPowerChanges(gained, lost []*power.Reference) *PowerChanges {
	return &PowerChanges{
		gained: gained,
		lost:   lost,
	}
}

// Gained is a getter.
func (p PowerChanges) Gained() []*power.Reference {
	return p.gained
}

// Lost is a getter.
func (p PowerChanges) Lost() []*power.Reference {
	return p.lost
}
