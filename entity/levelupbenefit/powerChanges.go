package levelupbenefit

import (
	"github.com/chadius/terosgamerules/entity/powerreference"
)

// PowerChanges tracks changes to the squaddie's innate powers.
type PowerChanges struct {
	gained []*powerreference.Reference
	lost   []*powerreference.Reference
}

// NewPowerChanges returns a new PowerChanges object.
func NewPowerChanges(gained, lost []*powerreference.Reference) *PowerChanges {
	return &PowerChanges{
		gained: gained,
		lost:   lost,
	}
}

// Gained is a getter.
func (p PowerChanges) Gained() []*powerreference.Reference {
	return p.gained
}

// Lost is a getter.
func (p PowerChanges) Lost() []*powerreference.Reference {
	return p.lost
}
