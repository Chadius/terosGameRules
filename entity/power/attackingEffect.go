package power

import (
	"errors"
	"github.com/chadius/terosbattleserver/utility"
)

// CounterAttackPenaltyInitialValue is the default penalty counter attacks suffer on to hit rolls.
const CounterAttackPenaltyInitialValue = -2

// AttackingEffect is a power designed to deal damage.
type AttackingEffect struct {
	AttackToHitBonus                    int             `json:"to_hit_bonus" yaml:"to_hit_bonus"`
	AttackDamageBonus                   int             `json:"damage_bonus" yaml:"damage_bonus"`
	AttackExtraBarrierBurn              int             `json:"extra_barrier_damage" yaml:"extra_barrier_damage"`
	AttackCanBeEquipped                 bool            `json:"can_be_equipped" yaml:"can_be_equipped"`
	AttackCanCounterAttack              bool            `json:"can_counter_attack" yaml:"can_counter_attack"`
	AttackCounterAttackPenaltyReduction int             `json:"counter_attack_penalty_reduction" yaml:"counter_attack_penalty_reduction"`
	CriticalEffect                      *CriticalEffect `json:"critical_effect" yaml:"critical_effect"`
}

// CanCriticallyHit returns true if this power is capable of critically hitting for additional effects.
func (a *AttackingEffect) CanCriticallyHit() bool {
	return a.CriticalEffect != nil
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
	return a.AttackToHitBonus
}

// DamageBonus returns the value.
func (a *AttackingEffect) DamageBonus() int {
	return a.AttackDamageBonus
}

// ExtraBarrierBurn returns the value.
func (a *AttackingEffect) ExtraBarrierBurn() int {
	return a.AttackExtraBarrierBurn
}

// CanBeEquipped returns the value.
func (a *AttackingEffect) CanBeEquipped() bool {
	return a.AttackCanBeEquipped
}

// CanCounterAttack returns the value.
func (a *AttackingEffect) CanCounterAttack() bool {
	return a.AttackCanCounterAttack
}

// CounterAttackPenaltyReduction returns the value.
func (a *AttackingEffect) CounterAttackPenaltyReduction() int {
	return a.AttackCounterAttackPenaltyReduction
}

// CriticalHitThreshold delegates.
func (a *AttackingEffect) CriticalHitThreshold() int {
	return a.CriticalEffect.CriticalHitThreshold()
}

// CriticalHitThresholdBonus delegates.
func (a *AttackingEffect) CriticalHitThresholdBonus() int {
	return a.CriticalEffect.CriticalHitThresholdBonus()
}

// ExtraCriticalHitDamage delegates.
func (a *AttackingEffect) ExtraCriticalHitDamage() int {
	return a.CriticalEffect.ExtraCriticalHitDamage()
}
