package power

import (
	"encoding/json"
	"github.com/chadius/terosgamerules/entity/healing"
	"github.com/chadius/terosgamerules/entity/powerinterface"
	"github.com/chadius/terosgamerules/entity/powersource"
	"github.com/chadius/terosgamerules/entity/target"
	"github.com/chadius/terosgamerules/utility"
	"gopkg.in/yaml.v2"
	"reflect"
)

// Builder covers options used to make Power objects.
type Builder struct {
	name                 string
	id                   string
	targetSelf           bool
	targetFriend         bool
	targetFoe            bool
	powerSourceLogic     powersource.Interface
	healingEffectOptions *HealingEffectOptions
	attackEffectOptions  *AttackEffectOptions
	healingLogic         healing.Interface
}

// NewPowerBuilder creates a Builder with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func NewPowerBuilder() *Builder {
	return &Builder{
		name:                 "power with no name",
		id:                   "",
		targetSelf:           false,
		targetFriend:         false,
		targetFoe:            false,
		powerSourceLogic:     powersource.NewPowerSourceLogic("physical"),
		healingEffectOptions: HealingEffectBuilder(),
		attackEffectOptions:  nil,
		healingLogic:         &healing.NoHealing{},
	}
}

// WithName applies the given name to the power.
func (p *Builder) WithName(name string) *Builder {
	p.name = name
	return p
}

// WithID applies the given ID to the power.
func (p *Builder) WithID(id string) *Builder {
	p.id = id
	return p
}

// IsPhysical sets the power type to physical.
func (p *Builder) IsPhysical() *Builder {
	p.powerSourceLogic = powersource.NewPowerSourceLogic("physical")
	return p
}

// IsSpell sets the power type to spell.
func (p *Builder) IsSpell() *Builder {
	p.powerSourceLogic = powersource.NewPowerSourceLogic("spell")
	return p
}

// TargetsSelf means the power can target the user.
func (p *Builder) TargetsSelf() *Builder {
	p.targetSelf = true
	return p
}

// TargetsFriend means the power can target the user's friends.
func (p *Builder) TargetsFriend() *Builder {
	p.targetFriend = true
	return p
}

// TargetsFoe means the power can target the user's friends.
func (p *Builder) TargetsFoe() *Builder {
	p.targetFoe = true
	return p
}

// HitPointsHealed delegates to the HealingEffectOptions.
func (p *Builder) HitPointsHealed(heal int) *Builder {
	p.healingEffectOptions.HitPointsHealed(heal)
	return p
}

// HealingAdjustmentBasedOnUserMindFull delegates to the HealingEffectOptions.
func (p *Builder) HealingAdjustmentBasedOnUserMindFull() *Builder {
	p.healingLogic = &healing.FullMindBonus{}
	return p
}

// HealingAdjustmentBasedOnUserMindHalf delegates to the HealingEffectOptions.
func (p *Builder) HealingAdjustmentBasedOnUserMindHalf() *Builder {
	p.healingLogic = &healing.HalfMindBonus{}
	return p
}

// HealingAdjustmentBasedOnUserMindZero delegates to the HealingEffectOptions.
func (p *Builder) HealingAdjustmentBasedOnUserMindZero() *Builder {
	p.healingLogic = &healing.ZeroMindBonus{}
	return p
}

// WithHealingLogic adds HealingLogic, using the given keyword
func (p *Builder) WithHealingLogic(keyword string) *Builder {
	p.healingLogic = healing.NewHealingLogic(keyword)
	return p
}

// DealsDamage delegates to the AttackEffectOptions.
func (p *Builder) DealsDamage(damage int) *Builder {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.DealsDamage(damage)
	return p
}

// ToHitBonus delegates to the AttackEffectOptions.
func (p *Builder) ToHitBonus(toHitBonus int) *Builder {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.ToHitBonus(toHitBonus)
	return p
}

// ExtraBarrierBurn delegates to the AttackEffectOptions.
func (p *Builder) ExtraBarrierBurn(extraBarrierBurn int) *Builder {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.ExtraBarrierBurn(extraBarrierBurn)
	return p
}

// CounterAttackPenaltyReduction delegates to the AttackEffectOptions.
func (p *Builder) CounterAttackPenaltyReduction(penaltyReduction int) *Builder {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.CanCounterAttack().CounterAttackPenaltyReduction(penaltyReduction)
	return p
}

// CanBeEquipped delegates to the AttackEffectOptions.
func (p *Builder) CanBeEquipped() *Builder {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.CanBeEquipped()
	return p
}

// CannotBeEquipped delegates to the AttackEffectOptions.
func (p *Builder) CannotBeEquipped() *Builder {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.CannotBeEquipped()
	return p
}

// CanCounterAttack delegates to the AttackEffectOptions.
func (p *Builder) CanCounterAttack() *Builder {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.CanCounterAttack()
	return p
}

// CriticalDealsDamage delegates to the AttackEffectOptions.
func (p *Builder) CriticalDealsDamage(damage int) *Builder {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.CriticalDealsDamage(damage)
	return p
}

// CriticalHitThresholdBonus delegates to the AttackEffectOptions.
func (p *Builder) CriticalHitThresholdBonus(thresholdBonus int) *Builder {
	if p.attackEffectOptions == nil {
		p.attackEffectOptions = AttackEffectBuilder()
	}
	p.attackEffectOptions.CriticalHitThresholdBonus(thresholdBonus)
	return p
}

// Build uses the Builder to create a power.
func (p *Builder) Build() *Power {
	var attackEffect *AttackingEffect = nil
	if p.attackEffectOptions != nil {
		attackEffect = p.attackEffectOptions.Build()
	}
	var healingEffect = p.healingEffectOptions.Build()

	targetOptions := []target.Interface{}
	if p.targetSelf {
		targetOptions = append(targetOptions, target.NewTargetingLogic("self"))
	}
	if p.targetFriend {
		targetOptions = append(targetOptions, target.NewTargetingLogic("friend"))
	}
	if p.targetFoe {
		targetOptions = append(targetOptions, target.NewTargetingLogic("foe"))
	}

	newPower := NewPower(
		p.name,
		p.id,
		p.powerSourceLogic,
		attackEffect,
		healingEffect,
		p.healingLogic,
		targetOptions,
	)
	return newPower
}

//Axe creates a Specific example of a physical attack power.
func (p *Builder) Axe() *Builder {
	p.WithName("axe").WithID("powerAxe").TargetsFoe().CanBeEquipped().CanCounterAttack().DealsDamage(1).ToHitBonus(1).Build()
	return p
}

//Spear creates a Specific example of a physical attack power.
func (p *Builder) Spear() *Builder {
	p.WithName("spear").WithID("powerSpear").TargetsFoe().CanBeEquipped().CanCounterAttack().DealsDamage(1).ToHitBonus(1).Build()
	return p
}

//Blot creates a Specific example of a spell attack power.
func (p *Builder) Blot() *Builder {
	p.WithName("blot").WithID("powerBlot").TargetsFoe().IsSpell().CanBeEquipped().DealsDamage(3).Build()
	return p
}

//HealingStaff creates a Specific example of a spell healing power.
func (p *Builder) HealingStaff() *Builder {
	p.WithName("healingStaff").WithID("powerHealingStaff").TargetsFriend().IsSpell().HitPointsHealed(3).HealingAdjustmentBasedOnUserMindFull()
	return p
}

// BuilderOptionMarshal is a flattened representation of all Squaddie NewPowerBuilder options.
type BuilderOptionMarshal struct {
	ID          string `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	PowerSource string `json:"source" yaml:"source"`

	TargetSelf   bool `json:"target_self" yaml:"target_self"`
	TargetFoe    bool `json:"target_foe" yaml:"target_foe"`
	TargetFriend bool `json:"target_friend" yaml:"target_friend"`

	CanAttack                     bool `json:"can_attack" yaml:"can_attack"`
	ToHitBonus                    int  `json:"to_hit_bonus" yaml:"to_hit_bonus"`
	DamageBonus                   int  `json:"damage_bonus" yaml:"damage_bonus"`
	ExtraBarrierBurn              int  `json:"extra_barrier_damage" yaml:"extra_barrier_damage"`
	CanBeEquipped                 bool `json:"can_be_equipped" yaml:"can_be_equipped"`
	CanCounterAttack              bool `json:"can_counter_attack" yaml:"can_counter_attack"`
	CounterAttackPenaltyReduction int  `json:"counter_attack_penalty_reduction" yaml:"counter_attack_penalty_reduction"`

	CanCritical               bool `json:"can_critical" yaml:"can_critical"`
	CriticalHitThresholdBonus int  `json:"critical_hit_threshold_bonus" yaml:"critical_hit_threshold_bonus"`
	CriticalDamage            int  `json:"critical_damage" yaml:"critical_damage"`

	HealingLogic    string `json:"healing_logic" yaml:"healing_logic"`
	HitPointsHealed int    `json:"hit_points_healed" yaml:"hit_points_healed"`
}

// UsingYAML uses the yaml data to generate Builder.
func (p *Builder) UsingYAML(yamlData []byte) *Builder {
	return p.usingByteStreamForOneOption(yamlData, yaml.Unmarshal)
}

// UsingJSON uses the yaml data to generate Builder.
func (p *Builder) UsingJSON(jsonData []byte) *Builder {
	return p.usingByteStreamForOneOption(jsonData, json.Unmarshal)
}

func (p *Builder) usingByteStreamForOneOption(data []byte, unmarshal utility.UnmarshalFunc) *Builder {
	var unmarshalError error
	var marshaledOptions BuilderOptionMarshal
	unmarshalError = unmarshal(data, &marshaledOptions)

	if unmarshalError != nil {
		return p
	}

	return p.usingMarshaledOptions(&marshaledOptions)
}

// CreatePowerBuilderOptionsFromYAML takes a YAML stream and converts them to a list of Builder.
func CreatePowerBuilderOptionsFromYAML(yamlData []byte) []*Builder {
	return usingByteStreamForMultipleOptions(yamlData, yaml.Unmarshal)
}

// CreatePowerBuilderOptionsFromJSON takes a JSON stream and converts them to a list of Builder.
func CreatePowerBuilderOptionsFromJSON(jsonData []byte) []*Builder {
	return usingByteStreamForMultipleOptions(jsonData, json.Unmarshal)
}

func usingByteStreamForMultipleOptions(data []byte, unmarshal utility.UnmarshalFunc) []*Builder {
	var unmarshalError error
	var allMarshaledOptions []BuilderOptionMarshal
	unmarshalError = unmarshal(data, &allMarshaledOptions)

	if unmarshalError != nil {
		return nil
	}

	builderOptions := []*Builder{}
	for _, marshaledOptions := range allMarshaledOptions {
		newOption := NewPowerBuilder().usingMarshaledOptions(&marshaledOptions)
		builderOptions = append(builderOptions, newOption)
	}

	return builderOptions
}

func (p *Builder) usingMarshaledOptions(marshaledOptions *BuilderOptionMarshal) *Builder {
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

	p.HitPointsHealed(marshaledOptions.HitPointsHealed)
	p.WithHealingLogic(marshaledOptions.HealingLogic)

	p.powerSourceLogic = powersource.NewPowerSourceLogic(marshaledOptions.PowerSource)

	if marshaledOptions.TargetSelf == true {
		p.TargetsSelf()
	}
	if marshaledOptions.TargetFoe == true {
		p.TargetsFoe()
	}
	if marshaledOptions.TargetFriend == true {
		p.TargetsFriend()
	}

	return p
}

// CloneOf modifies the Builder based on the source, except for the classID.
func (p *Builder) CloneOf(source powerinterface.Interface) *Builder {
	p.WithName(source.Name())

	p.clonePowerType(source)
	p.cloneTargeting(source)
	p.cloneAttackEffect(source)
	p.cloneHealingEffect(source)

	return p
}

func (p *Builder) cloneHealingEffect(source powerinterface.Interface) {
	p.HitPointsHealed(source.HitPointsHealed())
	p.WithHealingLogic(reflect.TypeOf(source.HealingLogic()).String())
}

func (p *Builder) cloneAttackEffect(source powerinterface.Interface) {
	if source.CanAttack() {
		p.ToHitBonus(source.ToHitBonus()).DealsDamage(source.DamageBonus()).ExtraBarrierBurn(source.ExtraBarrierBurn()).
			CounterAttackPenaltyReduction(source.CounterAttackPenaltyReduction())

		if source.CanCritical() {
			p.CriticalHitThresholdBonus(source.CriticalHitThresholdBonus()).CriticalDealsDamage(source.ExtraCriticalHitDamage())
		}

		if source.CanBeEquipped() {
			p.CanBeEquipped()
		}
		if source.CanCounterAttack() {
			p.CanCounterAttack()
		}

		if source.CanPowerTargetFriend() {
			p.TargetsFriend()
		}
		if source.CanPowerTargetFoe() {
			p.TargetsFoe()
		}
		if source.CanPowerTargetSelf() {
			p.TargetsSelf()
		}
	}
}

func (p *Builder) cloneTargeting(source powerinterface.Interface) {
	if source.CanPowerTargetFoe() {
		p.TargetsFoe()
	}

	if source.CanPowerTargetFriend() {
		p.TargetsFriend()
	}

	if source.CanPowerTargetSelf() {
		p.TargetsSelf()
	}
}

func (p *Builder) clonePowerType(source powerinterface.Interface) {
	p.powerSourceLogic = powersource.NewPowerSourceLogic(
		source.PowerSourceLogic().Name(),
	)
}
