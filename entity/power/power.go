package power

import (
	"fmt"

	"github.com/cserrant/terosBattleServer/utility"
)

// PowerReference is used to identify a power and is used to quickly identify a power.
type PowerReference struct {
	Name string `json:"name" yaml:"name"`
	ID   string `json:"id" yaml:"id"`
}

// PowerType defines the expected sources the power could be conjured from.
type PowerType string

const (
	// PowerTypePhysical powers use martial training and cunning. Examples: Swords, Bows, Pushing
	PowerTypePhysical = "Physical"
	// PowerTypeSpell powers are magical in nature and conjured without tools. Examples: Fireball, Mindread
	PowerTypeSpell = "Spell"
)

// Power are the abilities every Squaddie can use. These range from dealing damage, to opening doors, to healing.
type Power struct {
	PowerReference `yaml:",inline"`
	PowerType      PowerType `json:"power_type" yaml:"power_type"`
	*AttackingEffect
}

// AttackingEffect is a power designed to deal damage.
type AttackingEffect struct {
	ToHitBonus           int `json:"to_hit_bonus" yaml:"to_hit_bonus"`
	DamageBonus          int `json:"damage_bonus" yaml:"damage_bonus"`
	ExtraBarrierDamage   int `json:"extra_barrier_damage" yaml:"extra_barrier_damage"`
	CriticalHitThreshold int `json:"critical_hit_threshold" yaml:"critical_hit_threshold"`
}

// NewPower generates a Power with default values.
func NewPower(name string) *Power {
	newAttackingPower := Power{
		PowerReference: PowerReference{
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
	return &newAttackingPower
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

// GetChanceToCriticalBasedOnThreshold is a smarter look up table.
func GetChanceToCriticalBasedOnThreshold(critThreshold int) (chanceOutOf36 int) {
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
