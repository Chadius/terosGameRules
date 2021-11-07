package squaddie

import "github.com/chadius/terosbattleserver/entity/squaddie"

// OffenseBuilderOptions is used to set a squaddie's offensive attributes.
type OffenseBuilderOptions struct {
	aim      int
	strength int
	mind     int
}

// OffenseBuilder creates a OffenseBuilderOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func OffenseBuilder() *OffenseBuilderOptions {
	return &OffenseBuilderOptions{
		aim:      0,
		strength: 0,
		mind:     0,
	}
}

// Aim sets the squaddie's chance to hit.
func (o *OffenseBuilderOptions) Aim(aim int) *OffenseBuilderOptions {
	o.aim = aim
	return o
}

// Strength sets the squaddie's damage bonus with Physical attacks.
func (o *OffenseBuilderOptions) Strength(strength int) *OffenseBuilderOptions {
	o.strength = strength
	return o
}

// Mind sets the squaddie's damage bonus with Physical attacks.
func (o *OffenseBuilderOptions) Mind(mind int) *OffenseBuilderOptions {
	o.mind = mind
	return o
}

// Build uses the OffenseBuilderOptions to create a Movement.
func (o *OffenseBuilderOptions) Build() *squaddie.Offense {
	newOffense := squaddie.NewOffense(o.aim, o.strength, o.mind)
	return newOffense
}
