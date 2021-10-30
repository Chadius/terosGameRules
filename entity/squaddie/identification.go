package squaddie

import (
	"github.com/chadius/terosbattleserver/utility"
)

// Identification is akin to an SquaddieID card for each Squaddie. Each Squaddie carries a unique Identification.
type Identification struct {
	SquaddieID          string      `json:"id" yaml:"id"`
	SquaddieName        string      `json:"name" yaml:"name"`
	SquaddieAffiliation Affiliation `json:"affiliation" yaml:"affiliation"`
}

// SetNewIDToRandom changes the SquaddieID to a random value.
func (identification *Identification) SetNewIDToRandom() {
	identification.SquaddieID = utility.StringWithCharset(8, "abcdefgh0123456789")
}

// ID returns the squaddieID.
func (identification *Identification) ID() string {
	return identification.SquaddieID
}

// Name returns the squaddie's name.
func (identification *Identification) Name() string {
	return identification.SquaddieName
}

// Affiliation shows what affiliation the squaddie is a part of.
func (identification *Identification) Affiliation() Affiliation {
	return identification.SquaddieAffiliation
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
