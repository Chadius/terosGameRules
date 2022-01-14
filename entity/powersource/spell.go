package powersource

import "github.com/chadius/terosgamerules/entity/squaddieinterface"

// Spell measures Spell attacks.
type Spell struct{}

// Name returns a human-readable name of the source.
func (p *Spell) Name() string {
	return "spell"
}

// ToHitPenalty returns the squaddie's effectiveness to avoid getting hit in the first place.
func (p *Spell) ToHitPenalty(s squaddieinterface.Interface) int {
	return s.Deflect()
}

// ArmorResistance measures how well the squaddie can reduce damage due to their armor.
func (p *Spell) ArmorResistance(s squaddieinterface.Interface) int {
	return 0
}

// BarrierResistance measures how well the squaddie can use the barrier to reduce damage.
func (p *Spell) BarrierResistance(s squaddieinterface.Interface) int {
	return s.CurrentBarrier()
}

// RawDamage returns how much damage the squaddie can make with these power sources.
func (p *Spell) RawDamage(s squaddieinterface.Interface) int {
	return s.Mind()
}
