package power

import (
	"encoding/json"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/utility"
	"gopkg.in/yaml.v2"
)

// BuilderOptions covers options used to make Power objects.
type BuilderOptions struct {
	name                 string
	id                   string
	targetSelf           bool
	targetFriend         bool
	targetFoe            bool
	powerType            power.DamageType
	healingEffectOptions *HealingEffectOptions
	attackEffectOptions  *AttackEffectOptions
}

// Builder creates a BuilderOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func Builder() *BuilderOptions {
	return &BuilderOptions{
		name:                 "power with no name",
		id:                   "",
		targetSelf:           false,
		targetFriend:         false,
		targetFoe:            false,
		powerType:            power.Physical,
		healingEffectOptions: nil,
		attackEffectOptions:  nil,
	}
}

// WithName applies the given name to the power.
func (p *BuilderOptions) WithName(name string) *BuilderOptions {
	p.name = name
	return p
}

// WithID applies the given SquaddieID to the power.
func (p *BuilderOptions) WithID(id string) *BuilderOptions {
	p.id = id
	return p
}

// IsPhysical sets the power type to physical.
func (p *BuilderOptions) IsPhysical() *BuilderOptions {
	p.powerType = power.Physical
	return p
}

// IsSpell sets the power type to spell.
func (p *BuilderOptions) IsSpell() *BuilderOptions {
	p.powerType = power.Spell
	return p
}

// TargetsSelf means the power can target the user.
func (p *BuilderOptions) TargetsSelf() *BuilderOptions {
	p.targetSelf = true
	return p
}

// TargetsFriend means the power can target the user's friends.
func (p *BuilderOptions) TargetsFriend() *BuilderOptions {
	p.targetFriend = true
	return p
}

// TargetsFoe means the power can target the user's friends.
func (p *BuilderOptions) TargetsFoe() *BuilderOptions {
	p.targetFoe = true
	return p
}

// HitPointsHealed delegates to the HealingEffectOptions.
func (p *BuilderOptions) HitPointsHealed(heal int) *BuilderOptions {
	if p.healingEffectOptions == nil {
		p.healingEffectOptions = HealingEffectBuilder()
	}
	p.healingEffectOptions.HitPointsHealed(heal)
	return p
}

// HealingAdjustmentBasedOnUserMindFull delegates to the HealingEffectOptions.
func (p *BuilderOptions) HealingAdjustmentBasedOnUserMindFull() *BuilderOptions {
	if p.healingEffectOptions == nil {
		p.healingEffectOptions = HealingEffectBuilder()
	}
	p.healingEffectOptions.HealingAdjustmentBasedOnUserMindFull()
	return p
}

// HealingAdjustmentBasedOnUserMindHalf delegates to the HealingEffectOptions.
func (p *BuilderOptions) HealingAdjustmentBasedOnUserMindHalf() *BuilderOptions {
	if p.healingEffectOptions == nil {
		p.healingEffectOptions = HealingEffectBuilder()
	}
	p.healingEffectOptions.HealingAdjustmentBasedOnUserMindHalf()
	return p
}

// HealingAdjustmentBasedOnUserMindZero delegates to the HealingEffectOptions.
func (p *BuilderOptions) HealingAdjustmentBasedOnUserMindZero() *BuilderOptions {
	if p.healingEffectOptions == nil {
		p.healingEffectOptions = HealingEffectBuilder()
	}
	p.healingEffectOptions.HealingAdjustmentBasedOnUserMindZero()
	return p
}

// DealsDamage delegates to the AttackEffectOptions.
func (p *BuilderOptions) DealsDamage(damage int) *BuilderOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.DealsDamage(damage)
	return p
}

// ToHitBonus delegates to the AttackEffectOptions.
func (p *BuilderOptions) ToHitBonus(toHitBonus int) *BuilderOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.ToHitBonus(toHitBonus)
	return p
}

// ExtraBarrierBurn delegates to the AttackEffectOptions.
func (p *BuilderOptions) ExtraBarrierBurn(extraBarrierBurn int) *BuilderOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.ExtraBarrierBurn(extraBarrierBurn)
	return p
}

// CounterAttackPenaltyReduction delegates to the AttackEffectOptions.
func (p *BuilderOptions) CounterAttackPenaltyReduction(penaltyReduction int) *BuilderOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.CounterAttackPenaltyReduction(penaltyReduction)
	return p
}

// CanBeEquipped delegates to the AttackEffectOptions.
func (p *BuilderOptions) CanBeEquipped() *BuilderOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.CanBeEquipped()
	return p
}

// CannotBeEquipped delegates to the AttackEffectOptions.
func (p *BuilderOptions) CannotBeEquipped() *BuilderOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.CannotBeEquipped()
	return p
}

// CanCounterAttack delegates to the AttackEffectOptions.
func (p *BuilderOptions) CanCounterAttack() *BuilderOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.CanCounterAttack()
	return p
}

// CriticalDealsDamage delegates to the AttackEffectOptions.
func (p *BuilderOptions) CriticalDealsDamage(damage int) *BuilderOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.CriticalDealsDamage(damage)
	return p
}

// CriticalHitThresholdBonus delegates to the AttackEffectOptions.
func (p *BuilderOptions) CriticalHitThresholdBonus(thresholdBonus int) *BuilderOptions {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.CriticalHitThresholdBonus(thresholdBonus)
	return p
}

// Build uses the BuilderOptions to create a power.
func (p *BuilderOptions) Build() *power.Power {
	newPower := power.NewPower(p.name)
	if p.id != "" {
		newPower.Reference.PowerID = p.id
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
func (p *BuilderOptions) Axe() *BuilderOptions {
	p.WithName("axe").WithID("powerAxe").TargetsFoe().CanBeEquipped().CanCounterAttack().DealsDamage(1).ToHitBonus(1).Build()
	return p
}

//Spear creates a Specific example of a physical attack power.
func (p *BuilderOptions) Spear() *BuilderOptions {
	p.WithName("spear").WithID("powerSpear").TargetsFoe().CanBeEquipped().CanCounterAttack().DealsDamage(1).ToHitBonus(1).Build()
	return p
}

//Blot creates a Specific example of a spell attack power.
func (p *BuilderOptions) Blot() *BuilderOptions {
	p.WithName("blot").WithID("powerBlot").TargetsFoe().IsSpell().CanBeEquipped().DealsDamage(3).Build()
	return p
}

//HealingStaff creates a Specific example of a spell healing power.
func (p *BuilderOptions) HealingStaff() *BuilderOptions {
	p.WithName("healingStaff").WithID("powerHealingStaff").TargetsFriend().IsSpell().HitPointsHealed(3)
	return p
}

// BuilderOptionMarshal is a flattened representation of all Squaddie Builder options.
type BuilderOptionMarshal struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
	PowerType     power.DamageType       `json:"power_type" yaml:"power_type"`

	TargetSelf   bool `json:"target_self" yaml:"target_self"`
	TargetFoe    bool `json:"target_foe" yaml:"target_foe"`
	TargetFriend bool `json:"target_friend" yaml:"target_friend"`

	CanAttack bool `json:"can_attack" yaml:"can_attack"`
	ToHitBonus                    int             `json:"to_hit_bonus" yaml:"to_hit_bonus"`
	DamageBonus                   int             `json:"damage_bonus" yaml:"damage_bonus"`
	ExtraBarrierBurn              int             `json:"extra_barrier_damage" yaml:"extra_barrier_damage"`
	CanBeEquipped                 bool            `json:"can_be_equipped" yaml:"can_be_equipped"`
	CanCounterAttack              bool            `json:"can_counter_attack" yaml:"can_counter_attack"`
	CounterAttackPenaltyReduction int             `json:"counter_attack_penalty_reduction" yaml:"counter_attack_penalty_reduction"`

	CanCritical bool `json:"can_critical" yaml:"can_critical"`
	CriticalHitThresholdBonus int `json:"critical_hit_threshold_bonus" yaml:"critical_hit_threshold_bonus"`
	CriticalDamage                    int `json:"critical_damage" yaml:"critical_damage"`

	CanHeal bool `json:"can_heal" yaml:"can_heal"`
	HealingAdjustmentBasedOnUserMind power.HealingAdjustmentBasedOnUserMind `json:"healing_adjustment_based_on_user_mind" yaml:"healing_adjustment_based_on_user_mind"`
	HitPointsHealed                  int                              `json:"hit_points_healed" yaml:"hit_points_healed"`
}

// UsingYAML uses the yaml data to generate BuilderOptions.
func (p *BuilderOptions) UsingYAML(yamlData []byte) *BuilderOptions {
	return p.usingByteStream(yamlData, yaml.Unmarshal)
}

// UsingJSON uses the yaml data to generate BuilderOptions.
func (p *BuilderOptions) UsingJSON(jsonData []byte) *BuilderOptions {
	return p.usingByteStream(jsonData, json.Unmarshal)
}

func (p *BuilderOptions) usingByteStream(data []byte, unmarshal utility.UnmarshalFunc) *BuilderOptions {
	var unmarshalError error
	var marshaledOptions BuilderOptionMarshal

	unmarshalError = unmarshal(data, &marshaledOptions)

	if unmarshalError != nil {
		return p
	}

	p.WithID(marshaledOptions.ID).WithName(marshaledOptions.Name)

	if marshaledOptions.CanAttack {
		p.ToHitBonus(marshaledOptions.ToHitBonus).DealsDamage(marshaledOptions.DamageBonus).
			ExtraBarrierBurn(marshaledOptions.ExtraBarrierBurn).CounterAttackPenaltyReduction(marshaledOptions.CounterAttackPenaltyReduction)

		if marshaledOptions.CanBeEquipped {
			p.CanBeEquipped()
		}

		if marshaledOptions.CanCounterAttack {
			p.CanCounterAttack()
		}

		if marshaledOptions.CanCritical {
			p.CriticalHitThresholdBonus(marshaledOptions.CriticalHitThresholdBonus).CriticalDealsDamage(marshaledOptions.CriticalDamage)
		}
	}

	if marshaledOptions.CanHeal {
		p.HitPointsHealed(marshaledOptions.HitPointsHealed)

		if marshaledOptions.HealingAdjustmentBasedOnUserMind == power.Full {
			p.HealingAdjustmentBasedOnUserMindFull()
		}
		if marshaledOptions.HealingAdjustmentBasedOnUserMind == power.Half {
			p.HealingAdjustmentBasedOnUserMindHalf()
		}
		if marshaledOptions.HealingAdjustmentBasedOnUserMind == power.Zero {
			p.HealingAdjustmentBasedOnUserMindZero()
		}
	}

	if marshaledOptions.PowerType == power.Physical {
		p.IsPhysical()
	}

	if marshaledOptions.PowerType == power.Spell {
		p.IsSpell()
	}

	if marshaledOptions.TargetSelf == true { p.TargetsSelf() }
	if marshaledOptions.TargetFoe == true { p.TargetsFoe() }
	if marshaledOptions.TargetFriend == true { p.TargetsFriend() }

	return p
}