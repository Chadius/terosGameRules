package power

import (
	"errors"
	"github.com/cserrant/terosBattleServer/utility"
)

// CounterAttackPenaltyInitialValue is the default penalty counter attacks suffer on to hit rolls.
const CounterAttackPenaltyInitialValue = -2

// AttackingEffect is a power designed to deal damage.
type AttackingEffect struct {
	ToHitBonus                    int             `json:"to_hit_bonus" yaml:"to_hit_bonus"`
	DamageBonus                   int             `json:"damage_bonus" yaml:"damage_bonus"`
	ExtraBarrierBurn              int             `json:"extra_barrier_damage" yaml:"extra_barrier_damage"`
	CanBeEquipped                 bool            `json:"can_be_equipped" yaml:"can_be_equipped"`
	CanCounterAttack              bool            `json:"can_counter_attack" yaml:"can_counter_attack"`
	CounterAttackPenaltyReduction int             `json:"counter_attack_penalty_reduction" yaml:"counter_attack_penalty_reduction"`
	CriticalEffect                *CriticalEffect `json:"critical_effect" yaml:"critical_effect"`
}

// CanCriticallyHit returns true if this power is capable of critically hitting for additional effects.
func (a *AttackingEffect) CanCriticallyHit() bool {
	return a.CriticalEffect != nil
}

// CounterAttackPenalty returns the amount the counter attack to hit check suffers.
func (a *AttackingEffect) CounterAttackPenalty() (int, error) {
	if a.CanCounterAttack != true {
		newError := errors.New("power cannot counter, cannot calculate penalty")
		utility.Log(newError.Error(),0, utility.Error)
		return 0, newError
	}

	return CounterAttackPenaltyInitialValue + a.CounterAttackPenaltyReduction, nil
}
