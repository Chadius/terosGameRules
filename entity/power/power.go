package power

import (
	"errors"
	"fmt"
	"github.com/chadius/terosbattleserver/entity/healing"
	"github.com/chadius/terosbattleserver/utility"
	"reflect"
)

// Reference is used to identify a power and is used to quickly identify a power.
type Reference struct {
	Name    string `json:"name" yaml:"name"`
	PowerID string `json:"id" yaml:"id"`
}

// DamageType defines the expected sources the power could be conjured from.
type DamageType string

const (
	// Physical powers use martial training and cunning. Examples: Swords, Bows, Pushing
	Physical DamageType = "physical"
	// Spell powers are magical in nature and conjured without tools. Examples: Fireball, Mindread
	Spell DamageType = "spell"
)

// Targeting notes how the power can be targeted.
type Targeting struct {
	TargetSelf   bool
	TargetFoe    bool
	TargetFriend bool
}

// Power are the abilities every Squaddie can use. These range from dealing damage, to opening doors, to healing.
type Power struct {
	Reference
	damageType    DamageType
	targeting     Targeting
	attackEffect  *AttackingEffect
	healingEffect *HealingEffect
	healingLogic  healing.Interface
}

// GetReference returns a new PowerReference.
func (p Power) GetReference() *Reference {
	return &Reference{
		Name:    p.Name(),
		PowerID: p.ID(),
	}
}

// NewPower generates a Power.
func NewPower(name, id string, damageType *DamageType, targeting *Targeting, attackEffect *AttackingEffect, healingEffect *HealingEffect, healingLogic healing.Interface) *Power {
	powerID := "power_" + utility.StringWithCharset(8, "abcdefgh0123456789")
	if id != "" {
		powerID = id
	}
	newAttackingPower := Power{
		Reference: Reference{
			Name:    name,
			PowerID: powerID,
		},
		damageType:    *damageType,
		targeting:     *targeting,
		attackEffect:  attackEffect,
		healingEffect: healingEffect,
		healingLogic:  healingLogic,
	}
	return &newAttackingPower
}

// CheckPowerForErrors verifies the Power's fields and raises an error if it's invalid.
func CheckPowerForErrors(newPower *Power) (newError error) {
	if newPower.Type() != Physical &&
		newPower.Type() != Spell {
		newError := fmt.Errorf("AttackingPower '%s' has unknown power_type: '%s'", newPower.Name(), newPower.Type())
		utility.Log(newError.Error(), 0, utility.Error)
		return newError
	}

	return nil
}

// ID returns the power's ID.
func (p *Power) ID() string {
	return p.Reference.PowerID
}

// Name returns the power's Name.
func (p *Power) Name() string {
	return p.Reference.Name
}

// Type returns the power's damage type.
func (p *Power) Type() DamageType {
	return p.damageType
}

// CanPowerTargetSelf checks to see if the power can be used on the user.
func (p *Power) CanPowerTargetSelf() bool {
	return p.targeting.TargetSelf
}

// CanPowerTargetFriend checks to see if the power can be used on allies and teammates.
func (p *Power) CanPowerTargetFriend() bool {
	return p.targeting.TargetFriend
}

// CanPowerTargetFoe checks to see if the power can be used on enemies.
func (p *Power) CanPowerTargetFoe() bool {
	return p.targeting.TargetFoe
}

// CanAttack returns true if this power can be used to attack.
func (p *Power) CanAttack() bool {
	return p.attackEffect != nil
}

// ToHitBonus delegates.
func (p *Power) ToHitBonus() int {
	if !p.CanAttack() {
		return 0
	}
	return p.attackEffect.ToHitBonus()
}

// DamageBonus delegates.
func (p *Power) DamageBonus() int {
	if !p.CanAttack() {
		return 0
	}
	return p.attackEffect.DamageBonus()
}

// ExtraBarrierBurn delegates.
func (p *Power) ExtraBarrierBurn() int {
	if !p.CanAttack() {
		return 0
	}
	return p.attackEffect.ExtraBarrierBurn()
}

// CanBeEquipped delegates.
func (p *Power) CanBeEquipped() bool {
	if !p.CanAttack() {
		return false
	}
	return p.attackEffect.CanBeEquipped()
}

// CanCounterAttack delegates.
func (p *Power) CanCounterAttack() bool {
	if !p.CanAttack() {
		return false
	}
	return p.attackEffect.CanCounterAttack()
}

// CanCritical returns true if this power critically hit.
func (p *Power) CanCritical() bool {
	return p.attackEffect.CanCriticallyHit()
}

// CanCriticallyHit is an alias of CanCritical.
func (p *Power) CanCriticallyHit() bool {
	return p.CanCritical()
}

// CounterAttackPenaltyReduction delegates.
func (p *Power) CounterAttackPenaltyReduction() int {
	if !p.CanCounterAttack() {
		return 0
	}
	return p.attackEffect.CounterAttackPenaltyReduction()
}

// CounterAttackPenalty delegates.
func (p *Power) CounterAttackPenalty() (int, error) {
	if !p.CanCounterAttack() {
		newError := errors.New("power cannot counter, cannot calculate penalty")
		utility.Log(newError.Error(), 0, utility.Error)
		return 0, newError
	}

	penalty, err := p.attackEffect.CounterAttackPenalty()
	return penalty, err
}

// CriticalHitThreshold delegates.
func (p *Power) CriticalHitThreshold() int {
	if !p.CanCritical() {
		return 0
	}
	return p.attackEffect.CriticalHitThreshold()
}

// CriticalHitThresholdBonus delegates.
func (p *Power) CriticalHitThresholdBonus() int {
	if !p.CanCritical() {
		return 0
	}
	return p.attackEffect.CriticalHitThresholdBonus()
}

// ExtraCriticalHitDamage delegates.
func (p *Power) ExtraCriticalHitDamage() int {
	if !p.CanCritical() {
		return 0
	}
	return p.attackEffect.ExtraCriticalHitDamage()
}

// CanHeal returns true if this power can be used to heal.
func (p *Power) CanHeal() bool {
	return reflect.TypeOf(p.HealingLogic()).String() != "*healing.NoHealing"
}

// HitPointsHealed delegates.
func (p *Power) HitPointsHealed() int {
	return p.healingEffect.HitPointsHealed()
}

// HealingLogic returns the module used for healing.
func (p *Power) HealingLogic() healing.Interface {
	return p.healingLogic
}

// HasSameStatsAs returns true if other's stats matches this one.
//   The comparison ignores the ID.
func (p *Power) HasSameStatsAs(other *Power) bool {
	if p.Name() != other.Name() {
		return false
	}
	if p.Type() != other.Type() {
		return false
	}

	if !p.hasSameTargetingAs(other) {
		return false
	}

	if !p.hasSameAttackEffectAs(other) {
		return false
	}

	if !p.hasSameHealingEffectAs(other) {
		return false
	}

	return true
}

func (p *Power) hasSameHealingEffectAs(other *Power) bool {
	if p.HitPointsHealed() != other.HitPointsHealed() {
		return false
	}
	if reflect.TypeOf(p.HealingLogic()).String() != reflect.TypeOf(other.HealingLogic()).String() {
		return false
	}
	return true
}

func (p *Power) hasSameAttackEffectAs(other *Power) bool {
	if p.CanAttack() != other.CanAttack() {
		return false
	}
	if p.CanAttack() {
		if p.CanBeEquipped() != other.CanBeEquipped() {
			return false
		}
		if p.CanCounterAttack() != other.CanCounterAttack() {
			return false
		}
		if p.ToHitBonus() != other.ToHitBonus() {
			return false
		}
		if p.DamageBonus() != other.DamageBonus() {
			return false
		}
		if p.ExtraBarrierBurn() != other.ExtraBarrierBurn() {
			return false
		}
		if p.CounterAttackPenaltyReduction() != other.CounterAttackPenaltyReduction() {
			return false
		}

		if p.CanCritical() != other.CanCritical() {
			return false
		}
		if p.CanCritical() {
			if p.CriticalHitThresholdBonus() != other.CriticalHitThresholdBonus() {
				return false
			}
			if p.ExtraCriticalHitDamage() != other.ExtraCriticalHitDamage() {
				return false
			}
		}
	}
	return true
}

func (p *Power) hasSameTargetingAs(other *Power) bool {
	if p.CanPowerTargetFriend() != other.CanPowerTargetFriend() {
		return false
	}
	if p.CanPowerTargetFoe() != other.CanPowerTargetFoe() {
		return false
	}
	if p.CanPowerTargetSelf() != other.CanPowerTargetSelf() {
		return false
	}
	return true
}
