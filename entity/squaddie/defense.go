package squaddie

import (
	"github.com/chadius/terosbattleserver/entity/damagedistribution"
)

// Defense holds everything needed to prevent the squaddie from getting hindered.
type Defense struct {
	SquaddieCurrentHitPoints int `json:"current_hit_points" yaml:"current_hit_points"`
	SquaddieMaxHitPoints     int `json:"max_hit_points" yaml:"max_hit_points"`
	SquaddieDodge            int `json:"dodge" yaml:"dodge"`
	SquaddieDeflect          int `json:"deflect" yaml:"deflect"`
	SquaddieCurrentBarrier   int `json:"current_barrier" yaml:"current_barrier"`
	SquaddieMaxBarrier       int `json:"max_barrier" yaml:"max_barrier"`
	SquaddieArmor            int `json:"armor" yaml:"armor"`
}

// NewDefense returns a new Defense object.
func NewDefense(currentHitPoints, maxHitPoints, dodge, deflect, currentBarrier, maxBarrier, armor int) *Defense {
	return &Defense{
		SquaddieCurrentHitPoints: currentHitPoints,
		SquaddieMaxHitPoints:     maxHitPoints,
		SquaddieDodge:            dodge,
		SquaddieDeflect:          deflect,
		SquaddieCurrentBarrier:   currentBarrier,
		SquaddieMaxBarrier:       maxBarrier,
		SquaddieArmor:            armor,
	}
}

// SetHPToMax restores the Squaddie's HitPoints.
func (defense *Defense) SetHPToMax() {
	defense.SquaddieCurrentHitPoints = defense.SquaddieMaxHitPoints
}

// SetBarrierToMax restores the Squaddie's Barrier.
func (defense *Defense) SetBarrierToMax() {
	defense.SquaddieCurrentBarrier = defense.SquaddieMaxBarrier
}

// ReduceHitPoints reduces the squaddie's HP, possibly killing them.
//   Hit Points cannot e reduced below 0.
func (defense *Defense) ReduceHitPoints(damage int) int {
	actualDamageTaken := damage
	if defense.SquaddieCurrentHitPoints < damage {
		actualDamageTaken = defense.SquaddieCurrentHitPoints
	}

	defense.SquaddieCurrentHitPoints -= damage

	if defense.SquaddieCurrentHitPoints < 0 {
		defense.SquaddieCurrentHitPoints = 0
	}
	return actualDamageTaken
}

// ReduceBarrier reduces the squaddie's Barrier.
//   Barrier cannot e reduced below 0.
func (defense *Defense) ReduceBarrier(burn int) int {
	actualBarrierBurn := burn
	if defense.SquaddieCurrentBarrier < burn {
		actualBarrierBurn = defense.SquaddieCurrentBarrier
	}

	defense.SquaddieCurrentBarrier -= burn

	if defense.SquaddieCurrentBarrier < 0 {
		defense.SquaddieCurrentBarrier = 0
	}

	return actualBarrierBurn
}

// TakeDamageDistribution reduces the Squaddie's Barrier and Hit Points based on the distribution.
func (defense *Defense) TakeDamageDistribution(distribution *damagedistribution.DamageDistribution) {
	actualBarrierBurn := defense.ReduceBarrier(distribution.DamageAbsorbedByBarrier)
	actualDamageTaken := defense.ReduceHitPoints(distribution.RawDamageDealt)

	distribution.ActualBarrierBurn = actualBarrierBurn
	distribution.ActualDamageTaken = actualDamageTaken
}

// IsDead returns true if the squaddie has died
func (defense *Defense) IsDead() bool {
	return defense.SquaddieCurrentHitPoints <= 0
}

// GainHitPoints heals the squaddie and returns the number of hit points healed.
func (defense *Defense) GainHitPoints(hitPoints int) int {
	actualHealingReceived := hitPoints
	if defense.SquaddieCurrentHitPoints+actualHealingReceived >= defense.SquaddieMaxHitPoints {
		actualHealingReceived = defense.SquaddieMaxHitPoints - defense.SquaddieCurrentHitPoints
	}
	defense.SquaddieCurrentHitPoints += actualHealingReceived
	return actualHealingReceived
}

// MaxHitPoints returns the value.
func (defense *Defense) MaxHitPoints() int {
	return defense.SquaddieMaxHitPoints
}

// CurrentHitPoints returns the value.
func (defense *Defense) CurrentHitPoints() int {
	return defense.SquaddieCurrentHitPoints
}

// Dodge returns the value.
func (defense *Defense) Dodge() int {
	return defense.SquaddieDodge
}

// Deflect returns the value.
func (defense *Defense) Deflect() int {
	return defense.SquaddieDeflect
}

// MaxBarrier returns the value.
func (defense *Defense) MaxBarrier() int {
	return defense.SquaddieMaxBarrier
}

// CurrentBarrier returns the value.
func (defense *Defense) CurrentBarrier() int {
	return defense.SquaddieCurrentBarrier
}

// Armor returns the value.
func (defense *Defense) Armor() int {
	return defense.SquaddieArmor
}

// Improve improves the defensive stats.
func (defense *Defense) Improve(maxHitPoints, dodge, deflect, maxBarrier, armor int) {
	defense.SquaddieMaxHitPoints += maxHitPoints
	defense.SquaddieDodge += dodge
	defense.SquaddieDeflect += deflect
	defense.SquaddieMaxBarrier += maxBarrier
	defense.SquaddieArmor += armor
}
