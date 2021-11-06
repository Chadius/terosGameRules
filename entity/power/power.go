package power

import (
	"fmt"
	"github.com/chadius/terosbattleserver/utility"
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
	TargetSelf   bool `json:"target_self" yaml:"target_self"`
	TargetFoe    bool `json:"target_foe" yaml:"target_foe"`
	TargetFriend bool `json:"target_friend" yaml:"target_friend"`
}

// Power are the abilities every Squaddie can use. These range from dealing damage, to opening doors, to healing.
type Power struct {
	Reference
	PowerType     DamageType
	Targeting     Targeting
	AttackEffect  *AttackingEffect
	HealingEffect *HealingEffect
}

// GetReference returns a new PowerReference.
func (p Power) GetReference() *Reference {
	return &Reference{
		Name:    p.Name(),
		PowerID: p.ID(),
	}
}

// NewPower generates a Power with default values.
func NewPower(name string) *Power {
	newAttackingPower := Power{
		Reference: Reference{
			Name:    name,
			PowerID: "power_" + utility.StringWithCharset(8, "abcdefgh0123456789"),
		},
		PowerType: Physical,
	}
	return &newAttackingPower
}

// CheckPowerForErrors verifies the Power's fields and raises an error if it's invalid.
func CheckPowerForErrors(newPower *Power) (newError error) {
	if newPower.PowerType != Physical &&
		newPower.PowerType != Spell {
		newError := fmt.Errorf("AttackingPower '%s' has unknown power_type: '%s'", newPower.Name(), newPower.PowerType)
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
	return p.PowerType
}

// CanPowerTargetSelf checks to see if the power can be used on the user.
func (p *Power) CanPowerTargetSelf() bool {
	return p.Targeting.TargetSelf
}

// CanPowerTargetFriend checks to see if the power can be used on allies and teammates.
func (p *Power) CanPowerTargetFriend() bool {
	return p.Targeting.TargetFriend
}

// CanPowerTargetFoe checks to see if the power can be used on enemies.
func (p *Power) CanPowerTargetFoe() bool {
	return p.Targeting.TargetFoe
}

// ToHitBonus delegates.
func (p *Power) ToHitBonus() int {
	return p.AttackEffect.ToHitBonus()
}

// DamageBonus delegates.
func (p *Power) DamageBonus() int {
	return p.AttackEffect.DamageBonus()
}

// ExtraBarrierBurn delegates.
func (p *Power) ExtraBarrierBurn() int {
	return p.AttackEffect.ExtraBarrierBurn()
}

// CanBeEquipped delegates.
func (p *Power) CanBeEquipped() bool {
	return p.AttackEffect.CanBeEquipped()
}

// CanCounterAttack delegates.
func (p *Power) CanCounterAttack() bool {
	return p.AttackEffect.CanCounterAttack()
}

// CounterAttackPenaltyReduction delegates.
func (p *Power) CounterAttackPenaltyReduction() int {
	return p.AttackEffect.CounterAttackPenaltyReduction()
}

// CriticalHitThreshold delegates.
func (p *Power) CriticalHitThreshold() int {
	return p.AttackEffect.CriticalHitThreshold()
}

// CriticalHitThresholdBonus delegates.
func (p *Power) CriticalHitThresholdBonus() int {
	return p.AttackEffect.CriticalHitThresholdBonus()
}

// ExtraCriticalHitDamage delegates.
func (p *Power) ExtraCriticalHitDamage() int {
	return p.AttackEffect.ExtraCriticalHitDamage()
}

// HitPointsHealed delegates.
func (p *Power) HitPointsHealed() int {
	return p.HealingEffect.HitPointsHealed()
}

// HealingAdjustmentBasedOnUserMind delegates.
func (p *Power) HealingAdjustmentBasedOnUserMind() HealingAdjustmentBasedOnUserMind {
	return p.HealingEffect.HealingAdjustmentBasedOnUserMind()
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
	if p.CanHeal() != other.CanHeal() {
		return false
	}
	if p.CanHeal() {
		if p.HitPointsHealed() != other.HitPointsHealed() {
			return false
		}
		if p.HealingAdjustmentBasedOnUserMind() != other.HealingAdjustmentBasedOnUserMind() {
			return false
		}
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

// CanAttack returns true if this power can be used to attack.
func (p *Power) CanAttack() bool {
	return p.AttackEffect != nil
}

// CanCritical returns true if this power critically hit.
func (p *Power) CanCritical() bool {
	return p.AttackEffect.CanCriticallyHit()
}

// CanCriticallyHit is an alias of CanCritical.
func (p *Power) CanCriticallyHit() bool {
	return p.CanCritical()
}

// CanHeal returns true if this power can be used to heal.
func (p *Power) CanHeal() bool {
	return p.HealingEffect != nil
}

// CounterAttackPenalty delegates.
func (p *Power) CounterAttackPenalty() (int, error) {
	penalty, err := p.AttackEffect.CounterAttackPenalty()
	return penalty, err
}
