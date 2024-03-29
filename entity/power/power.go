package power

import (
	"errors"
	"github.com/chadius/terosgamerules/entity/healing"
	"github.com/chadius/terosgamerules/entity/powerinterface"
	"github.com/chadius/terosgamerules/entity/powerreference"
	"github.com/chadius/terosgamerules/entity/powersource"
	"github.com/chadius/terosgamerules/entity/target"
	"github.com/chadius/terosgamerules/utility"
	"reflect"
)

// Power are the abilities every Squaddie can use. These range from dealing damage, to opening doors, to healing.
type Power struct {
	powerreference.Reference
	powerSourceLogic powersource.Interface
	attackEffect     *AttackingEffect
	healingEffect    *HealingEffect
	healingLogic     healing.Interface
	targetLogic      []target.Interface
}

// GetReference returns a new PowerReference.
func (p Power) GetReference() *powerreference.Reference {
	return &powerreference.Reference{
		Name:    p.Name(),
		PowerID: p.ID(),
	}
}

// NewPower generates a Power.
func NewPower(name, id string, powerSourceLogic powersource.Interface, attackEffect *AttackingEffect, healingEffect *HealingEffect, healingLogic healing.Interface, targetLogicObjects []target.Interface) *Power {
	powerID := "power_" + utility.StringWithCharset(8, "abcdefgh0123456789")
	if id != "" {
		powerID = id
	}
	newAttackingPower := Power{
		Reference: powerreference.Reference{
			Name:    name,
			PowerID: powerID,
		},
		powerSourceLogic: powerSourceLogic,
		attackEffect:     attackEffect,
		healingEffect:    healingEffect,
		healingLogic:     healingLogic,
		targetLogic:      targetLogicObjects,
	}
	return &newAttackingPower
}

// ID returns the power's ID.
func (p *Power) ID() string {
	return p.Reference.PowerID
}

// Name returns the power's Name.
func (p *Power) Name() string {
	return p.Reference.Name
}

// PowerSourceLogic returns the power's source logic.
func (p *Power) PowerSourceLogic() powersource.Interface {
	return p.powerSourceLogic
}

// CanPowerTargetSelf checks to see if the power can be used on the user.
func (p *Power) CanPowerTargetSelf() bool {
	for _, targetObject := range p.targetLogic {
		if targetObject.Name() == "self" {
			return true
		}
	}

	return false
}

// CanPowerTargetFriend checks to see if the power can be used on allies and teammates.
func (p *Power) CanPowerTargetFriend() bool {
	for _, targetObject := range p.targetLogic {
		if targetObject.Name() == "friend" {
			return true
		}
	}

	return false
}

// CanPowerTargetFoe checks to see if the power can be used on enemies.
func (p *Power) CanPowerTargetFoe() bool {
	for _, targetObject := range p.targetLogic {
		if targetObject.Name() == "foe" {
			return true
		}
	}

	return false
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
func (p *Power) HasSameStatsAs(other powerinterface.Interface) bool {
	if p.Name() != other.Name() {
		return false
	}
	if reflect.TypeOf(p.PowerSourceLogic()).String() != reflect.TypeOf(other.PowerSourceLogic()).String() {
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

func (p *Power) hasSameHealingEffectAs(other powerinterface.Interface) bool {
	if p.HitPointsHealed() != other.HitPointsHealed() {
		return false
	}
	if reflect.TypeOf(p.HealingLogic()).String() != reflect.TypeOf(other.HealingLogic()).String() {
		return false
	}
	return true
}

func (p *Power) hasSameAttackEffectAs(other powerinterface.Interface) bool {
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

func (p *Power) hasSameTargetingAs(other powerinterface.Interface) bool {
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
