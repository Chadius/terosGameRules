package levelupbenefit

import "github.com/chadius/terosbattleserver/entity/power"

// PowerChanges tracks changes to the squaddie's innate powers.
type PowerChanges struct {
	PrivatizeMeGained []*power.Reference `json:"gained" yaml:"gained"`
	PrivatizeMeLost   []*power.Reference `json:"lost" yaml:"lost"`
}

// NewPowerChanges returns a new PowerChanges object.
func NewPowerChanges(gained, lost []*power.Reference) *PowerChanges {
	return &PowerChanges{
		PrivatizeMeGained: gained,
		PrivatizeMeLost:   lost,
	}
}

// Gained is a getter.
func (p PowerChanges) Gained() []*power.Reference {
	return p.PrivatizeMeGained
}

// Lost is a getter.
func (p PowerChanges) Lost() []*power.Reference {
	return p.PrivatizeMeLost
}