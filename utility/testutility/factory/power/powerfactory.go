package power

import "github.com/chadius/terosbattleserver/entity/power"

type PowerFactoryOptions struct {
	name string
	id string
	targetSelf bool
	targetFriend bool
	targetFoe bool
	powerType power.DamageType
	healingEffectOptions *HealingEffectOptions
	attackEffectOptions *AttackEffectOptions
}

// PowerFactory creates a PowerFactoryOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func PowerFactory() *PowerFactoryOptions {
	return &PowerFactoryOptions{
		name: "power with no name",
		id: "",
		targetSelf: false,
		targetFriend: false,
		targetFoe: false,
		powerType: power.Physical,
		healingEffectOptions: nil,
		attackEffectOptions: nil,
	}
}

// WithName applies the given name to the power.
func (p *PowerFactoryOptions) WithName(name string) *PowerFactoryOptions {
	p.name = name
	return p
}

// WithID applies the given ID to the power.
func (p *PowerFactoryOptions) WithID(id string) *PowerFactoryOptions {
	p.id = id
	return p
}

// IsPhysical sets the power type to physical.
func (p *PowerFactoryOptions) IsPhysical() *PowerFactoryOptions {
	p.powerType = power.Physical
	return p
}

// IsSpell sets the power type to spell.
func (p *PowerFactoryOptions) IsSpell() *PowerFactoryOptions {
	p.powerType = power.Spell
	return p
}

// TargetsSelf means the power can target the user.
func (p *PowerFactoryOptions) TargetsSelf() *PowerFactoryOptions {
	p.targetSelf = true
	return p
}

// TargetsFriend means the power can target the user's friends.
func (p *PowerFactoryOptions) TargetsFriend() *PowerFactoryOptions {
	p.targetFriend = true
	return p
}

// TargetsFoe means the power can target the user's friends.
func (p *PowerFactoryOptions) TargetsFoe() *PowerFactoryOptions {
	p.targetFoe = true
	return p
}

// HitPointsHealed delegates to the HealingEffectOptions.
func (p *PowerFactoryOptions) HitPointsHealed(heal int) *PowerFactoryOptions {
	if p.healingEffectOptions == nil {
		p.healingEffectOptions = HealingEffectFactory()
	}
	p.healingEffectOptions.HitPointsHealed(heal)
	return p
}

// HealingAdjustmentBasedOnUserMindFull delegates to the HealingEffectOptions.
func (p *PowerFactoryOptions) HealingAdjustmentBasedOnUserMindFull() *PowerFactoryOptions {
	if p.healingEffectOptions == nil {
		p.healingEffectOptions = HealingEffectFactory()
	}
	p.healingEffectOptions.HealingAdjustmentBasedOnUserMindFull()
	return p
}

// HealingAdjustmentBasedOnUserMindHalf delegates to the HealingEffectOptions.
func (p *PowerFactoryOptions) HealingAdjustmentBasedOnUserMindHalf() *PowerFactoryOptions {
	if p.healingEffectOptions == nil {
		p.healingEffectOptions = HealingEffectFactory()
	}
	p.healingEffectOptions.HealingAdjustmentBasedOnUserMindHalf()
	return p
}

// HealingAdjustmentBasedOnUserMindZero delegates to the HealingEffectOptions.
func (p *PowerFactoryOptions) HealingAdjustmentBasedOnUserMindZero() *PowerFactoryOptions {
	if p.healingEffectOptions == nil {
		p.healingEffectOptions = HealingEffectFactory()
	}
	p.healingEffectOptions.HealingAdjustmentBasedOnUserMindZero()
	return p
}

// DealsDamage delegates to the AttackEffectOptions.
func (p *PowerFactoryOptions) DealsDamage(damage int) *PowerFactoryOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectFactory()
	}
	p.attackEffectOptions.DealsDamage(damage)
	return p
}

// ToHitBonus delegates to the AttackEffectOptions.
func (p *PowerFactoryOptions) ToHitBonus(toHitBonus int) *PowerFactoryOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectFactory()
	}
	p.attackEffectOptions.ToHitBonus(toHitBonus)
	return p
}

// ExtraBarrierBurn delegates to the AttackEffectOptions.
func (p *PowerFactoryOptions) ExtraBarrierBurn(extraBarrierBurn int) *PowerFactoryOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectFactory()
	}
	p.attackEffectOptions.ExtraBarrierBurn(extraBarrierBurn)
	return p
}

// CounterAttackPenaltyReduction delegates to the AttackEffectOptions.
func (p *PowerFactoryOptions) CounterAttackPenaltyReduction(penaltyReduction int) *PowerFactoryOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectFactory()
	}
	p.attackEffectOptions.CounterAttackPenaltyReduction(penaltyReduction)
	return p
}

// CanBeEquipped delegates to the AttackEffectOptions.
func (p *PowerFactoryOptions) CanBeEquipped() *PowerFactoryOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectFactory()
	}
	p.attackEffectOptions.CanBeEquipped()
	return p
}

// CannotBeEquipped delegates to the AttackEffectOptions.
func (p *PowerFactoryOptions) CannotBeEquipped() *PowerFactoryOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectFactory()
	}
	p.attackEffectOptions.CannotBeEquipped()
	return p
}

// CanCounterAttack delegates to the AttackEffectOptions.
func (p *PowerFactoryOptions) CanCounterAttack() *PowerFactoryOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectFactory()
	}
	p.attackEffectOptions.CanCounterAttack()
	return p
}

// CriticalDealsDamage delegates to the AttackEffectOptions.
func (p *PowerFactoryOptions) CriticalDealsDamage(damage int) *PowerFactoryOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectFactory()
	}
	p.attackEffectOptions.CriticalDealsDamage(damage)
	return p
}

// CriticalHitThresholdBonus delegates to the AttackEffectOptions.
func (p *PowerFactoryOptions) CriticalHitThresholdBonus(thresholdBonus int) *PowerFactoryOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectFactory()
	}
	p.attackEffectOptions.CriticalHitThresholdBonus(thresholdBonus)
	return p
}


// Build uses the PowerFactoryOptions to create a power.
func (p *PowerFactoryOptions) Build() *power.Power {
	newPower := power.NewPower(p.name)
	if p.id != "" {
		newPower.Reference.ID = p.id
	}

	newPower.Targeting.TargetSelf = p.targetSelf
	newPower.Targeting.TargetFriend = p.targetFriend
	newPower.Targeting.TargetFoe = p.targetFoe
	newPower.PowerType = p.powerType

	if p.healingEffectOptions != nil {
		newPower.HealingEffect = p.healingEffectOptions.Build()
	}

	if p.attackEffectOptions != nil {
		newPower.AttackEffect = p.attackEffectOptions.Build()
	}
	return newPower
}


//Axe creates a Specific example of a physical attack power.
func (p *PowerFactoryOptions) Axe() *PowerFactoryOptions {
	p.WithName("axe").WithID("powerAxe").TargetsFoe().CanBeEquipped().CanCounterAttack().DealsDamage(1).ToHitBonus(1).Build()
	return p
}

//Spear creates a Specific example of a physical attack power.
func (p *PowerFactoryOptions) Spear() *PowerFactoryOptions {
	p.WithName("spear").WithID("powerSpear").TargetsFoe().CanBeEquipped().CanCounterAttack().DealsDamage(1).ToHitBonus(1).Build()
	return p
}

//Blot creates a Specific example of a spell attack power.
func (p *PowerFactoryOptions) Blot() *PowerFactoryOptions {
	p.WithName("blot").WithID("powerBlot").TargetsFoe().IsSpell().CanBeEquipped().DealsDamage(3).Build()
	return p
}

//HealingStaff creates a Specific example of a spell healing power.
func (p *PowerFactoryOptions) HealingStaff() *PowerFactoryOptions {
	p.WithName("healingStaff").WithID("powerHealingStaff").TargetsFriend().IsSpell().HitPointsHealed(3)
	return p
}