package squaddie

import (
	"github.com/chadius/terosgamerules/entity/damagedistribution"
)

// Defense holds everything needed to prevent the squaddie from getting hindered.
type Defense struct {
	currentHitPoints int
	maxHitPoints     int
	dodge            int
	deflect          int
	currentBarrier   int
	maxBarrier       int
	armor            int
}

// NewDefense returns a new Defense object.
func NewDefense(currentHitPoints, maxHitPoints, dodge, deflect, currentBarrier, maxBarrier, armor int) *Defense {
	return &Defense{
		currentHitPoints: currentHitPoints,
		maxHitPoints:     maxHitPoints,
		dodge:            dodge,
		deflect:          deflect,
		currentBarrier:   currentBarrier,
		maxBarrier:       maxBarrier,
		armor:            armor,
	}
}

// SetHPToMax restores the Squaddie's HitPoints.
func (defense *Defense) SetHPToMax() {
	defense.currentHitPoints = defense.maxHitPoints
}

// SetBarrierToMax restores the Squaddie's Barrier.
func (defense *Defense) SetBarrierToMax() {
	defense.currentBarrier = defense.maxBarrier
}

// ReduceHitPoints reduces the squaddie's HP, possibly killing them.
//   Hit Points cannot be reduced below 0.
func (defense *Defense) ReduceHitPoints(damage int) int {
	actualDamageTaken := damage
	if defense.currentHitPoints < damage {
		actualDamageTaken = defense.currentHitPoints
	}

	defense.currentHitPoints -= damage

	if defense.currentHitPoints < 0 {
		defense.currentHitPoints = 0
	}
	return actualDamageTaken
}

// ReduceBarrier reduces the squaddie's Barrier.
//   Barrier cannot be reduced below 0.
func (defense *Defense) ReduceBarrier(burn int) int {
	actualBarrierBurn := burn
	if defense.currentBarrier < burn {
		actualBarrierBurn = defense.currentBarrier
	}

	defense.currentBarrier -= burn

	if defense.currentBarrier < 0 {
		defense.currentBarrier = 0
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
	return defense.currentHitPoints <= 0
}

// GainHitPoints heals the squaddie and returns the number of hit points healed.
func (defense *Defense) GainHitPoints(hitPoints int) int {
	actualHealingReceived := hitPoints
	if defense.currentHitPoints+actualHealingReceived >= defense.maxHitPoints {
		actualHealingReceived = defense.maxHitPoints - defense.currentHitPoints
	}
	defense.currentHitPoints += actualHealingReceived
	return actualHealingReceived
}

// MaxHitPoints returns the value.
func (defense *Defense) MaxHitPoints() int {
	return defense.maxHitPoints
}

// CurrentHitPoints returns the value.
func (defense *Defense) CurrentHitPoints() int {
	return defense.currentHitPoints
}

// Dodge returns the value.
func (defense *Defense) Dodge() int {
	return defense.dodge
}

// Deflect returns the value.
func (defense *Defense) Deflect() int {
	return defense.deflect
}

// MaxBarrier returns the value.
func (defense *Defense) MaxBarrier() int {
	return defense.maxBarrier
}

// CurrentBarrier returns the value.
func (defense *Defense) CurrentBarrier() int {
	return defense.currentBarrier
}

// Armor returns the value.
func (defense *Defense) Armor() int {
	return defense.armor
}

// Improve improves the defensive stats.
func (defense *Defense) Improve(maxHitPoints, dodge, deflect, maxBarrier, armor int) {
	defense.maxHitPoints += maxHitPoints
	defense.dodge += dodge
	defense.deflect += deflect
	defense.maxBarrier += maxBarrier
	defense.armor += armor
}
