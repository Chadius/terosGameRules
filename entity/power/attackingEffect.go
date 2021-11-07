package power

import (
	"errors"
	"github.com/chadius/terosbattleserver/utility"
)

// CounterAttackPenaltyInitialValue is the default penalty counterattacks suffer on to hit rolls.
const CounterAttackPenaltyInitialValue = -2

// AttackingEffect is a power designed to deal damage.
type AttackingEffect struct {
	toHitBonus                    int
	damageBonus                   int
	extraBarrierBurn              int
	canBeEquipped                 bool
	canCounterAttack              bool
	counterAttackPenaltyReduction int
	criticalEffect                *CriticalEffect
}

// NewAttackingEffect returns a new AttackingEffect with the given options.
func NewAttackingEffect(toHitBonus, damageBonus, extraBarrierBurn int, canBeEquipped, canCounterAttack bool, counterAttackPenaltyReduction int, criticalEffect *CriticalEffect) *AttackingEffect {
	return &AttackingEffect{
		toHitBonus:                    toHitBonus,
		damageBonus:                   damageBonus,
		extraBarrierBurn:              extraBarrierBurn,
		canBeEquipped:                 canBeEquipped,
		canCounterAttack:              canCounterAttack,
		counterAttackPenaltyReduction: counterAttackPenaltyReduction,
		criticalEffect:                criticalEffect,
	}
}

// CanCriticallyHit returns true if this power is capable of critically hitting for additional effects.
func (a *AttackingEffect) CanCriticallyHit() bool {
	return a.criticalEffect != nil
}

// CounterAttackPenalty returns the amount the counterattack to hit check suffers.
func (a *AttackingEffect) CounterAttackPenalty() (int, error) {
	if a.CanCounterAttack() != true {
		newError := errors.New("power cannot counter, cannot calculate penalty")
		utility.Log(newError.Error(), 0, utility.Error)
		return 0, newError
	}

	return CounterAttackPenaltyInitialValue + a.CounterAttackPenaltyReduction(), nil
}

// ToHitBonus returns the value.
func (a *AttackingEffect) ToHitBonus() int {
	return a.toHitBonus
}

// DamageBonus returns the value.
func (a *AttackingEffect) DamageBonus() int {
	return a.damageBonus
}

// ExtraBarrierBurn returns the value.
func (a *AttackingEffect) ExtraBarrierBurn() int {
	return a.extraBarrierBurn
}

// CanBeEquipped returns the value.
func (a *AttackingEffect) CanBeEquipped() bool {
	return a.canBeEquipped
}

// CanCounterAttack returns the value.
func (a *AttackingEffect) CanCounterAttack() bool {
	return a.canCounterAttack
}

// CounterAttackPenaltyReduction returns the value.
func (a *AttackingEffect) CounterAttackPenaltyReduction() int {
	return a.counterAttackPenaltyReduction
}

// CriticalHitThreshold delegates.
func (a *AttackingEffect) CriticalHitThreshold() int {
	return a.criticalEffect.CriticalHitThreshold()
}

// CriticalHitThresholdBonus delegates.
func (a *AttackingEffect) CriticalHitThresholdBonus() int {
	return a.criticalEffect.CriticalHitThresholdBonus()
}

// ExtraCriticalHitDamage delegates.
func (a *AttackingEffect) ExtraCriticalHitDamage() int {
	return a.criticalEffect.ExtraCriticalHitDamage()
}
