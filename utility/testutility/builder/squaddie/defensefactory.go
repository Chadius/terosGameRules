package squaddie

import "github.com/chadius/terosbattleserver/entity/squaddie"

// DefenseBuilderOptions is used to set a squaddie's defensive attributes.
type DefenseBuilderOptions struct {
	armor        int
	deflect      int
	dodge        int
	maxHitPoints int
	maxBarrier   int
}

// DefenseBuilder creates a DefenseBuilderOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func DefenseBuilder() *DefenseBuilderOptions {
	return &DefenseBuilderOptions{
		armor:        0,
		deflect:      0,
		dodge:        0,
		maxHitPoints: 5,
		maxBarrier:   0,
	}
}

// Armor sets the damage reduction.
func (d *DefenseBuilderOptions) Armor(damageReduction int) *DefenseBuilderOptions {
	d.armor = damageReduction
	return d
}

// Dodge sets the chance to dodge physical attacks.
func (d *DefenseBuilderOptions) Dodge(dodgeBonus int) *DefenseBuilderOptions {
	d.dodge = dodgeBonus
	return d
}

// Deflect sets the chance to dodge spell attacks.
func (d *DefenseBuilderOptions) Deflect(deflectBonus int) *DefenseBuilderOptions {
	d.deflect = deflectBonus
	return d
}

// HitPoints sets the maximum and current hit points.
func (d *DefenseBuilderOptions) HitPoints(maxHitPoints int) *DefenseBuilderOptions {
	d.maxHitPoints = maxHitPoints
	return d
}

// Barrier sets the maximum barrier.
func (d *DefenseBuilderOptions) Barrier(maxBarrier int) *DefenseBuilderOptions {
	d.maxBarrier = maxBarrier
	return d
}

// Build uses the DefenseBuilderOptions to create a Movement.
func (d *DefenseBuilderOptions) Build() *squaddie.Defense {
	newDefense := &squaddie.Defense{
		SquaddieArmor:        d.armor,
		SquaddieDodge:        d.dodge,
		SquaddieDeflect:      d.deflect,
		SquaddieMaxHitPoints: d.maxHitPoints,
		SquaddieMaxBarrier:   d.maxBarrier,
	}

	newDefense.SetHPToMax()
	return newDefense
}
