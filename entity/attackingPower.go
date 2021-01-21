package entity

import (
	"encoding/json"
	"fmt"

	"github.com/cserrant/terosBattleServer/utility"
	"gopkg.in/yaml.v2"
)

// PowerIDName is a pair of ID and Name items used to quickly identify
//   a PowerHeader.
type PowerIDName struct {
	Name string `json:"name" yaml:"name"`
	ID   string `json:"id" yaml:"id"`
}

// PowerType defines the expected sources the power could be conjured from.
type PowerType string

const (
	// PowerTypePhysical powers use martial training and cunning. Examples: Swords, Bows, Pushing
	PowerTypePhysical PowerType = "Physical"
	// PowerTypeSpell powers are magical in nature and conjured without tools. Examples: Fireball, Mindread
	PowerTypeSpell = "Spell"
)

// Power are the abilities every Squaddie can use. These range from dealing damage, to opening doors, to healing.
type Power struct {
	PowerIDName `yaml:",inline"`
	PowerType   PowerType `json:"power_type" yaml:"power_type"`
	*AttackingEffect
}

// AttackingEffect is a power designed to deal damage.
type AttackingEffect struct {
	ToHitBonus           int `json:"to_hit_bonus" yaml:"to_hit_bonus"`
	DamageBonus          int `json:"damage_bonus" yaml:"damage_bonus"`
	ExtraBarrierDamage   int `json:"extra_barrier_damage" yaml:"extra_barrier_damage"`
	CriticalHitThreshold int `json:"critical_hit_threshold" yaml:"critical_hit_threshold"`
}

// NewAttackingPower generates a Power with default values.
func NewAttackingPower(name string) Power {
	newAttackingPower := Power{
		PowerIDName: PowerIDName{
			Name: name,
			ID:   utility.StringWithCharset(8, "abcdefgh0123456789"),
		},
		PowerType: PowerTypePhysical,
		AttackingEffect: &AttackingEffect{
			ToHitBonus:           0,
			DamageBonus:          0,
			ExtraBarrierDamage:   0,
			CriticalHitThreshold: 0,
		},
	}
	return newAttackingPower
}

// NewAttackingPowerFromJSON reads the JSON byte stream to create a new Squaddie.
// 	Defaults to NewAttackingPower.
func NewAttackingPowerFromJSON(data []byte) (newPower Power, err error) {
	newAttackingPower := NewAttackingPower("NewAttackingPowerFromJSON")
	err = json.Unmarshal(data, &newAttackingPower)
	if err != nil {
		return newAttackingPower, err
	}

	err = checkAttackingEffectForErrors(&newAttackingPower)
	return newAttackingPower, err
}

// NewAttackingPowerFromYAML reads the JSON byte stream to create a new Squaddie.
// 	Defaults to NewAttackingPower.
func NewAttackingPowerFromYAML(data []byte) (newPower Power, err error) {
	newAttackingPower := NewAttackingPower("NewAttackingPowerFromYAML")
	err = yaml.Unmarshal(data, &newAttackingPower)
	if err != nil {
		return newAttackingPower, err
	}

	attackingEffect := &AttackingEffect{}
	err = yaml.Unmarshal(data, &attackingEffect)
	if err != nil {
		return newAttackingPower, err
	}
	newAttackingPower.AttackingEffect = attackingEffect

	err = checkAttackingEffectForErrors(&newAttackingPower)
	return newAttackingPower, err
}

// CheckPowerForErrors verifies the Power's fields and raises an error if it's invalid.
func CheckPowerForErrors(newPower *Power) (newError error) {
	return checkAttackingEffectForErrors(newPower)
}

func checkAttackingEffectForErrors(newAttackingPower *Power) (newError error) {
	if newAttackingPower.PowerType != PowerTypePhysical &&
		newAttackingPower.PowerType != PowerTypeSpell {
		return fmt.Errorf("AttackingPower '%s' has unknown power_type: '%s'", newAttackingPower.Name, newAttackingPower.PowerType)
	}

	return nil
}

// GetTotalToHitBonus calculates the total to hit bonus for the attacking squaddie and attacking power
func (power *Power) GetTotalToHitBonus(squaddie *Squaddie) (toHit int) {
	return power.AttackingEffect.ToHitBonus + squaddie.Aim
}

// GetTotalDamageBonus calculates the total Damage bonus for the attacking squaddie and attacking power
func (power *Power) GetTotalDamageBonus(squaddie *Squaddie) (damageBonus int) {
	if power.PowerType == PowerTypePhysical {
		return power.AttackingEffect.DamageBonus + squaddie.Strength
	}
	return power.AttackingEffect.DamageBonus + squaddie.Mind
}

// GetCriticalDamageBonus calculates the total Critical Hit Damage bonus for the attacking squaddie and attacking power
func (power *Power) GetCriticalDamageBonus(squaddie *Squaddie) (damageBonus int) {
	return 2 * power.GetTotalDamageBonus(squaddie)
}

// GetToHitPenalty calculates how much the target can reduce the chance of getting hit by the attacking power.
func (power *Power) GetToHitPenalty(target *Squaddie) (toHitPenalty int) {
	if power.PowerType == PowerTypePhysical {
		return target.Dodge
	}
	return target.Deflect
}

// GetHowTargetDistributesDamage factors the attacker's damage bonuses and target's damage reduction to figure out the base damage and barrier damage.
func (power *Power) GetHowTargetDistributesDamage(attacker *Squaddie, target *Squaddie) (healthDamage, barrierDamage, extraBarrierDamage int) {
	damageToAbsorb := power.GetTotalDamageBonus(attacker)
	return power.calculateHowTargetTakesDamage(attacker, target, damageToAbsorb)
}

// GetHowTargetDistributesCriticalDamage factors the attacker's damage bonuses and target's damage reduction to figure out the base damage and barrier damage.
func (power *Power) GetHowTargetDistributesCriticalDamage(attacker *Squaddie, target *Squaddie) (healthDamage, barrierDamage, extraBarrierDamage int) {
	damageToAbsorb := power.GetCriticalDamageBonus(attacker)
	return power.calculateHowTargetTakesDamage(attacker, target, damageToAbsorb)
}

// calculateHowTargetTakesDamage factors the attacker's damage bonuses and target's damage reduction to figure out the base damage and barrier damage.
func (power *Power) calculateHowTargetTakesDamage(attacker *Squaddie, target *Squaddie, damageToAbsorb int) (healthDamage, barrierDamage, extraBarrierDamage int) {
	remainingBarrier := target.CurrentBarrier

	var barrierFullyAbsorbsDamage bool = target.CurrentBarrier > damageToAbsorb
	if barrierFullyAbsorbsDamage {
		barrierDamage = damageToAbsorb
		remainingBarrier = remainingBarrier - barrierDamage
		damageToAbsorb = 0
	} else {
		barrierDamage = target.CurrentBarrier
		remainingBarrier = 0
		damageToAbsorb = damageToAbsorb - target.CurrentBarrier
	}

	if power.ExtraBarrierDamage > 0 {
		var barrierFullyAbsorbsExtraBarrierDamage bool = remainingBarrier > power.ExtraBarrierDamage
		if barrierFullyAbsorbsExtraBarrierDamage {
			extraBarrierDamage = power.ExtraBarrierDamage
			remainingBarrier = remainingBarrier - power.ExtraBarrierDamage
		} else {
			extraBarrierDamage = remainingBarrier
			remainingBarrier = 0
		}
	}

	var armorCanAbsorbDamage bool = power.PowerType == PowerTypePhysical
	if armorCanAbsorbDamage {

		var armorFullyAbsorbsDamage bool = target.Armor > damageToAbsorb
		if armorFullyAbsorbsDamage {
			healthDamage = 0
		} else {
			healthDamage = damageToAbsorb - target.Armor
		}
	} else {
		healthDamage = damageToAbsorb
	}

	return healthDamage, barrierDamage, extraBarrierDamage
}

// AttackingPowerSummary gives a summary of the chance to hit and damage dealt by attacks. Expected damage counts the number of 36ths so we can use ints for fractional math.
type AttackingPowerSummary struct {
	ChanceToHit                   int
	DamageTaken                   int
	ExpectedDamage                int
	BarrierDamageTaken            int
	ExpectedBarrierDamage         int
	ChanceToCrit                  int
	CriticalDamageTaken           int
	CriticalBarrierDamageTaken    int
	CriticalExpectedDamage        int
	CriticalExpectedBarrierDamage int
}

// GetExpectedDamage provides a quick summary of an attack as well as the multiplied estimate
func (power *Power) GetExpectedDamage(attacker *Squaddie, target *Squaddie) (battleSummary *AttackingPowerSummary) {
	toHitBonus := power.GetTotalToHitBonus(attacker)
	toHitPenalty := power.GetToHitPenalty(target)
	totalChanceToHit := GetChanceToHitBasedOnHitRate(toHitBonus - toHitPenalty)

	healthDamage, barrierDamage, extraBarrierDamage := power.GetHowTargetDistributesDamage(attacker, target)

	chanceToCrit := GetChanceToCritBasedOnThreshold(power.CriticalHitThreshold)
	var criticalHealthDamage, criticalBarrierDamage, criticalExtraBarrierDamage int
	if chanceToCrit > 0 {
		criticalHealthDamage, criticalBarrierDamage, criticalExtraBarrierDamage = power.GetHowTargetDistributesCriticalDamage(attacker, target)
	} else {
		criticalHealthDamage, criticalBarrierDamage, criticalExtraBarrierDamage = 0, 0, 0
	}

	return &AttackingPowerSummary{
		ChanceToHit:                   totalChanceToHit,
		DamageTaken:                   healthDamage,
		ExpectedDamage:                totalChanceToHit * healthDamage,
		BarrierDamageTaken:            barrierDamage + extraBarrierDamage,
		ExpectedBarrierDamage:         totalChanceToHit * (barrierDamage + extraBarrierDamage),
		ChanceToCrit:                  chanceToCrit,
		CriticalDamageTaken:           criticalHealthDamage,
		CriticalBarrierDamageTaken:    criticalBarrierDamage + criticalExtraBarrierDamage,
		CriticalExpectedDamage:        totalChanceToHit * criticalHealthDamage,
		CriticalExpectedBarrierDamage: totalChanceToHit * (criticalBarrierDamage + criticalExtraBarrierDamage),
	}
}

// GetChanceToHitBasedOnHitRate is a smarter look up table.
func GetChanceToHitBasedOnHitRate(toHitBonus int) (chanceOutOf36 int) {
	if toHitBonus > 4 {
		return 36
	}

	if toHitBonus < -5 {
		return 0
	}

	toHitChanceReference := map[int]int{
		4:  35,
		3:  33,
		2:  30,
		1:  26,
		0:  21,
		-1: 15,
		-2: 10,
		-3: 6,
		-4: 3,
		-5: 1,
	}

	return toHitChanceReference[toHitBonus]
}

// GetChanceToCritBasedOnThreshold is a smarter look up table.
func GetChanceToCritBasedOnThreshold(critThreshold int) (chanceOutOf36 int) {
	if critThreshold > 11 {
		return 36
	}

	if critThreshold < 2 {
		return 0
	}

	critChanceReference := map[int]int{
		11: 35,
		10: 33,
		9:  30,
		8:  26,
		7:  21,
		6:  15,
		5:  10,
		4:  6,
		3:  3,
		2:  1,
	}

	return critChanceReference[critThreshold]
}
