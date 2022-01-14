package squaddie

import (
	"github.com/chadius/terosgamerules/entity/affiliation"
	"github.com/chadius/terosgamerules/utility"
)

// Identification is akin to an squaddieID card for each Squaddie. Each Squaddie carries a unique Identification.
type Identification struct {
	squaddieID       string
	squaddieName     string
	affiliationLogic affiliation.Interface
}

// NewIdentification creates a new Identification object.
func NewIdentification(id, name string, newAffiliationLogic affiliation.Interface) *Identification {
	return &Identification{
		squaddieID:       id,
		squaddieName:     name,
		affiliationLogic: newAffiliationLogic,
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

// AffiliationLogic shows what affiliation the squaddie is a part of.
func (identification *Identification) AffiliationLogic() affiliation.Interface {
	return identification.affiliationLogic
}
