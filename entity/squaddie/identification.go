package squaddie

import (
	"github.com/chadius/terosbattleserver/utility"
)

// Identification is akin to an squaddieID card for each Squaddie. Each Squaddie carries a unique Identification.
type Identification struct {
	squaddieID   string
	squaddieName string
	affiliation  Affiliation
}

// NewIdentification creates a new Identification object.
func NewIdentification(id, name string, affiliation Affiliation) *Identification {
	return &Identification{
		squaddieID:   id,
		squaddieName: name,
		affiliation:  affiliation,
	}
}

// SetNewIDToRandom changes the squaddieID to a random value.
func (identification *Identification) SetNewIDToRandom() {
	identification.squaddieID = utility.StringWithCharset(8, "abcdefgh0123456789")
}

// ID returns the squaddieID.
func (identification *Identification) ID() string {
	return identification.squaddieID
}

// Name returns the squaddie's name.
func (identification *Identification) Name() string {
	return identification.squaddieName
}

// Affiliation shows what affiliation the squaddie is a part of.
func (identification *Identification) Affiliation() Affiliation {
	return identification.affiliation
}

// HasValidAffiliation makes sure the created squaddie doesn't have an error.
func (identification *Identification) HasValidAffiliation() bool {
	if identification.Affiliation() != Player &&
		identification.Affiliation() != Enemy &&
		identification.Affiliation() != Ally &&
		identification.Affiliation() != Neutral {
		return false
	}

	return true
}
