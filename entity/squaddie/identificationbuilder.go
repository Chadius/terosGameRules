package squaddie

import (
	"github.com/chadius/terosbattleserver/entity/affiliation"
)

// IdentificationBuilderOptions is used to create healing effects.
type IdentificationBuilderOptions struct {
	name             string
	id               string
	affiliationLogic affiliation.Interface
}

// IdentificationBuilder creates a IdentificationBuilderOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func IdentificationBuilder() *IdentificationBuilderOptions {
	return &IdentificationBuilderOptions{
		name:             "squaddie with no name",
		id:               "",
		affiliationLogic: &affiliation.Neutral{},
	}
}

// WithName applies the given name to the Identification.
func (i *IdentificationBuilderOptions) WithName(name string) *IdentificationBuilderOptions {
	i.name = name
	return i
}

// WithID applies the given powerID to the Identification.
func (i *IdentificationBuilderOptions) WithID(id string) *IdentificationBuilderOptions {
	i.id = id
	return i
}

// AsPlayer makes the Identification as a Player.
func (i *IdentificationBuilderOptions) AsPlayer() *IdentificationBuilderOptions {
	i.affiliationLogic = &affiliation.Player{}
	return i
}

// AsEnemy makes the Identification as an Enemy.
func (i *IdentificationBuilderOptions) AsEnemy() *IdentificationBuilderOptions {
	i.affiliationLogic = &affiliation.Enemy{}
	return i
}

// AsAlly makes the Identification as an Ally.
func (i *IdentificationBuilderOptions) AsAlly() *IdentificationBuilderOptions {
	i.affiliationLogic = &affiliation.Ally{}
	return i
}

// AsNeutral makes the Identification as a Neutral.
func (i *IdentificationBuilderOptions) AsNeutral() *IdentificationBuilderOptions {
	i.affiliationLogic = &affiliation.Neutral{}
	return i
}

// Build uses the IdentificationBuilderOptions to create a Movement.
func (i *IdentificationBuilderOptions) Build() *Identification {
	newIdentification := NewIdentification(i.id, i.name, i.affiliationLogic)
	return newIdentification
}

// WithAffiliationLogic sets the affiliation logic.
func (i *IdentificationBuilderOptions) WithAffiliationLogic(keyword string) {
	i.affiliationLogic = affiliation.NewAffiliationLogic(keyword)
}
