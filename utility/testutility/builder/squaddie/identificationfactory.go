package squaddie

import "github.com/chadius/terosbattleserver/entity/squaddie"

// IdentificationBuilderOptions is used to create healing effects.
type IdentificationBuilderOptions struct {
	name        string
	id          string
	affiliation squaddie.Affiliation
}

// IdentificationBuilder creates a IdentificationBuilderOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func IdentificationBuilder() *IdentificationBuilderOptions {
	return &IdentificationBuilderOptions{
		name:        "squaddie with no name",
		id:          "",
		affiliation: squaddie.Neutral,
	}
}

// WithName applies the given name to the Identification.
func (i *IdentificationBuilderOptions) WithName(name string) *IdentificationBuilderOptions {
	i.name = name
	return i
}

// WithID applies the given PowerID to the Identification.
func (i *IdentificationBuilderOptions) WithID(id string) *IdentificationBuilderOptions {
	i.id = id
	return i
}

// AsPlayer makes the Identification as a Player.
func (i *IdentificationBuilderOptions) AsPlayer() *IdentificationBuilderOptions {
	i.affiliation = squaddie.Player
	return i
}

// AsEnemy makes the Identification as an Enemy.
func (i *IdentificationBuilderOptions) AsEnemy() *IdentificationBuilderOptions {
	i.affiliation = squaddie.Enemy
	return i
}

// AsAlly makes the Identification as an Ally.
func (i *IdentificationBuilderOptions) AsAlly() *IdentificationBuilderOptions {
	i.affiliation = squaddie.Ally
	return i
}

// AsNeutral makes the Identification as a Neutral.
func (i *IdentificationBuilderOptions) AsNeutral() *IdentificationBuilderOptions {
	i.affiliation = squaddie.Neutral
	return i
}

// Build uses the IdentificationBuilderOptions to create a Movement.
func (i *IdentificationBuilderOptions) Build() *squaddie.Identification {
	newIdentification := &squaddie.Identification{
		SquaddieName:        i.name,
		SquaddieID:          i.id,
		SquaddieAffiliation: i.affiliation,
	}
	return newIdentification
}
