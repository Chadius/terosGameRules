package powersource

import (
	"github.com/chadius/terosgamerules/entity/squaddieinterface"
)

// Physical measures physical attacks.
type Physical struct{}

// Name returns a human-readable name of the source.
func (p *Physical) Name() string {
	return "physical"
}

// ToHitPenalty returns the squaddie's effectiveness to avoid getting hit in the first place.
func (p *Physical) ToHitPenalty(s squaddieinterface.Interface) int {
	return s.Dodge()
}

// ArmorResistance measures how well the squaddie can reduce damage due to their armor.
func (p *Physical) ArmorResistance(s squaddieinterface.Interface) int {
	return s.Armor()
}

// BarrierResistance measures how well the squaddie can use the barrier to reduce damage.
func (p *Physical) BarrierResistance(s squaddieinterface.Interface) int {
	return s.CurrentBarrier()
}

// RawDamage returns how much damage the squaddie can make with these power sources.
func (p *Physical) RawDamage(s squaddieinterface.Interface) int {
	return s.Strength()
}
