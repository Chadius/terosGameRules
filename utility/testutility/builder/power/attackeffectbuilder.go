package power

import "github.com/chadius/terosbattleserver/entity/power"

// AttackEffectOptions is used to create healing effects.
type AttackEffectOptions struct {
	damage                        int
	toHitBonus                    int
	extraBarrierBurn              int
	counterAttackPenaltyReduction int
	canBeEquipped                 bool
	canCounterAttack              bool
	criticalEffectOptions         *CriticalEffectOptions
}

// AttackEffectBuilder creates a AttackEffectOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func AttackEffectBuilder() *AttackEffectOptions {
	return &AttackEffectOptions{
		damage:                        0,
		toHitBonus:                    0,
		extraBarrierBurn:              0,
		counterAttackPenaltyReduction: 0,
		canBeEquipped:                 false,
		canCounterAttack:              false,
		criticalEffectOptions:         nil,
	}
}

// DealsDamage sets the amount of base damage.
func (a *AttackEffectOptions) DealsDamage(damage int) *AttackEffectOptions {
	a.damage = damage
	return a
}

// ToHitBonus sets the to hit bonus.
func (a *AttackEffectOptions) ToHitBonus(toHitBonus int) *AttackEffectOptions {
	a.toHitBonus = toHitBonus
	return a
}

// ExtraBarrierBurn makes the attack deal more damage to squaddies with barriers.
func (a *AttackEffectOptions) ExtraBarrierBurn(extraBarrierBurn int) *AttackEffectOptions {
	a.extraBarrierBurn = extraBarrierBurn
	return a
}

// CounterAttackPenaltyReduction makes the counterattack more accurate.
func (a *AttackEffectOptions) CounterAttackPenaltyReduction(penaltyReduction int) *AttackEffectOptions {
	a.counterAttackPenaltyReduction = penaltyReduction
	return a
}

// CanBeEquipped means the attack can be equipped.
func (a *AttackEffectOptions) CanBeEquipped() *AttackEffectOptions {
	a.canBeEquipped = true
	return a
}

// CannotBeEquipped means the attack cannot be equipped.
func (a *AttackEffectOptions) CannotBeEquipped() *AttackEffectOptions {
	a.canBeEquipped = false
	return a
}

// CanCounterAttack means the attack can counterattack.
func (a *AttackEffectOptions) CanCounterAttack() *AttackEffectOptions {
	a.canCounterAttack = true
	return a
}

// CriticalDealsDamage delegates to the CriticalEffectOptions.
func (a *AttackEffectOptions) CriticalDealsDamage(damage int) *AttackEffectOptions {
	if a.criticalEffectOptions == nil {
		a.criticalEffectOptions = CriticalEffectBuilder()
	}
	a.criticalEffectOptions.DealsDamage(damage)
	return a
}

// CriticalHitThresholdBonus delegates to the CriticalEffectOptions.
func (a *AttackEffectOptions) CriticalHitThresholdBonus(thresholdBonus int) *AttackEffectOptions {
	if a.criticalEffectOptions == nil {
		a.criticalEffectOptions = CriticalEffectBuilder()
	}
	a.criticalEffectOptions.CriticalHitThresholdBonus(thresholdBonus)
	return a
}

// Build uses the AttackEffectOptions to create an AttackingEffect.
func (a *AttackEffectOptions) Build() *power.AttackingEffect {
	var criticalEffect *power.CriticalEffect = nil
	if a.criticalEffectOptions != nil {
		criticalEffect = a.criticalEffectOptions.Build()
	}
	newAttackingEffect := power.NewAttackingEffect(
		a.toHitBonus,
		a.damage,
		a.extraBarrierBurn,
		a.canBeEquipped,
		a.canCounterAttack,
		a.counterAttackPenaltyReduction,
		criticalEffect,
	)
	return newAttackingEffect
}
