package terosbattleserver

// PowerType defines the expected sources the power could be conjured from.
type PowerType string

const (
	// PowerTypePhysical powers use martial training and cunning. Examples: Swords, Bows, Pushing
	PowerTypePhysical PowerType = "Physical"
	// PowerTypeSpell powers are magical in nature and conjured without tools. Examples: Fireball, Mindread
	PowerTypeSpell = "Spell"
)

// AttackingPower is a power designed to deal damage.
type AttackingPower struct {
	Name               string
	PowerType          PowerType
	ToHitBonus         int
	DamageBonus        int
	ExtraBarrierDamage int
}

// NewAttackingPower generates a AttackingPower with default values.
func NewAttackingPower(name string) AttackingPower {
	newAttackingPower := AttackingPower{
		Name:               name,
		PowerType:          PowerTypePhysical,
		ToHitBonus:         0,
		DamageBonus:        0,
		ExtraBarrierDamage: 0,
	}
	return newAttackingPower
}

// GetTotalToHitBonus calculates the total to hit bonus for the attacking squaddie and attacking power
func (power *AttackingPower) GetTotalToHitBonus(squaddie *Squaddie) (toHit int) {
	return power.ToHitBonus + squaddie.Aim
}

// GetTotalDamageBonus calculates the total tDamage bonus for the attacking squaddie and attacking power
func (power *AttackingPower) GetTotalDamageBonus(squaddie *Squaddie) (damageBonus int) {
	if power.PowerType == PowerTypePhysical {
		return power.DamageBonus + squaddie.Strength
	}
	return power.DamageBonus + squaddie.Mind
}

// GetToHitPenalty calculates how much the target can reduce the chance of getting hit by the attacking power.
func (power *AttackingPower) GetToHitPenalty(target *Squaddie) (toHitPenalty int) {
	if power.PowerType == PowerTypePhysical {
		return target.Dodge
	}
	return target.Deflect
}

// GetDamageAgainstTarget factors the attacker's damage bonuses and target's damage reduction to figure out the base damage and barrier damage.
func (power *AttackingPower) GetDamageAgainstTarget(attacker *Squaddie, target *Squaddie) (healthDamage, barrierDamage, extraBarrierDamage int) {
	damageToAbsorb := power.GetTotalDamageBonus(attacker)
	remainingBarrier := target.CurrentBarrier

	var barrierFullyAbsorbsDamage bool = (target.CurrentBarrier > damageToAbsorb)
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
		var barrierFullyAbsorbsExtraBarrierDamage bool = (remainingBarrier > power.ExtraBarrierDamage)
		if barrierFullyAbsorbsExtraBarrierDamage {
			extraBarrierDamage = power.ExtraBarrierDamage
			remainingBarrier = remainingBarrier - power.ExtraBarrierDamage
		} else {
			extraBarrierDamage = remainingBarrier
			remainingBarrier = 0
		}
	}

	var armorCanAbsorbDamage bool = (power.PowerType == PowerTypePhysical)
	if armorCanAbsorbDamage {

		var armorFullyAbsorbsDamage bool = (target.Armor > damageToAbsorb)
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
