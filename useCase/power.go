package usecase

import (
	"github.com/cserrant/terosBattleServer/entity"
)

// GetPowerToHitBonusWhenUsedBySquaddie calculates the total to hit bonus for the attacking squaddie and attacking power
func GetPowerToHitBonusWhenUsedBySquaddie(power *entity.Power, squaddie *entity.Squaddie) (toHit int) {
	return power.AttackingEffect.ToHitBonus + squaddie.Aim
}

// GetPowerDamageBonusWhenUsedBySquaddie calculates the total Damage bonus for the attacking squaddie and attacking power
func GetPowerDamageBonusWhenUsedBySquaddie(power *entity.Power, squaddie *entity.Squaddie) (damageBonus int) {
	if power.PowerType == entity.PowerTypePhysical {
		return power.AttackingEffect.DamageBonus + squaddie.Strength
	}
	return power.AttackingEffect.DamageBonus + squaddie.Mind
}

// GetPowerCriticalDamageBonusWhenUsedBySquaddie calculates the total Critical Hit Damage bonus for the attacking squaddie and attacking power
func GetPowerCriticalDamageBonusWhenUsedBySquaddie(power *entity.Power, squaddie *entity.Squaddie) (damageBonus int) {
	return 2 * GetPowerDamageBonusWhenUsedBySquaddie(power, squaddie)
}

// GetHowTargetDistributesDamage factors the attacker's damage bonuses and target's damage reduction to figure out the base damage and barrier damage.
func GetHowTargetDistributesDamage(power *entity.Power, attacker *entity.Squaddie, target *entity.Squaddie) (healthDamage, barrierDamage, extraBarrierDamage int) {
	damageToAbsorb := GetPowerDamageBonusWhenUsedBySquaddie(power, attacker)
	return calculateHowTargetTakesDamage(power, target, damageToAbsorb)
}

// GetHowTargetDistributesCriticalDamage factors the attacker's damage bonuses and target's damage reduction to figure out the base damage and barrier damage.
func GetHowTargetDistributesCriticalDamage(power *entity.Power, attacker *entity.Squaddie, target *entity.Squaddie) (healthDamage, barrierDamage, extraBarrierDamage int) {
	damageToAbsorb := GetPowerCriticalDamageBonusWhenUsedBySquaddie(power, attacker)
	return calculateHowTargetTakesDamage(power, target, damageToAbsorb)
}

// calculateHowTargetTakesDamage factors the target's damage reduction to figure out how the damage is split between barrier, armor and health.
func calculateHowTargetTakesDamage(power *entity.Power, target *entity.Squaddie, damageToAbsorb int) (healthDamage, barrierDamage, extraBarrierDamage int) {
	remainingBarrier := target.CurrentBarrier

	damageToAbsorb, barrierDamage, remainingBarrier = calculateDamageAfterInitialBarrierAbsorption(target, damageToAbsorb, barrierDamage, remainingBarrier)

	extraBarrierDamage = calculateDamageAfterExtraBarrierDamage(power, remainingBarrier, extraBarrierDamage)

	healthDamage = calculateDamageAfterArmorAbsorption(power, target, damageToAbsorb, healthDamage)

	return healthDamage, barrierDamage, extraBarrierDamage
}

func calculateDamageAfterArmorAbsorption(power *entity.Power, target *entity.Squaddie, damageToAbsorb int, healthDamage int) int {
	var armorCanAbsorbDamage bool = power.PowerType == entity.PowerTypePhysical
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
	return healthDamage
}

func calculateDamageAfterExtraBarrierDamage(power *entity.Power, remainingBarrier int, extraBarrierDamage int) int {
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
	return extraBarrierDamage
}

func calculateDamageAfterInitialBarrierAbsorption(target *entity.Squaddie, damageToAbsorb int, barrierDamage int, remainingBarrier int) (int, int, int) {
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
	return damageToAbsorb, barrierDamage, remainingBarrier
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
func GetExpectedDamage(power *entity.Power, attacker *entity.Squaddie, target *entity.Squaddie) (battleSummary *AttackingPowerSummary) {
	toHitBonus := GetPowerToHitBonusWhenUsedBySquaddie(power, attacker)
	toHitPenalty := GetPowerToHitPenaltyAgainstSquaddie(power, target)
	totalChanceToHit := entity.GetChanceToHitBasedOnHitRate(toHitBonus - toHitPenalty)

	healthDamage, barrierDamage, extraBarrierDamage := GetHowTargetDistributesDamage(power, attacker, target)

	chanceToCritical := entity.GetChanceToCriticalBasedOnThreshold(power.CriticalHitThreshold)
	var criticalHealthDamage, criticalBarrierDamage, criticalExtraBarrierDamage int
	if chanceToCritical > 0 {
		criticalHealthDamage, criticalBarrierDamage, criticalExtraBarrierDamage = GetHowTargetDistributesCriticalDamage(power, attacker, target)
	} else {
		criticalHealthDamage, criticalBarrierDamage, criticalExtraBarrierDamage = 0, 0, 0
	}

	return &AttackingPowerSummary{
		ChanceToHit:                   totalChanceToHit,
		DamageTaken:                   healthDamage,
		ExpectedDamage:                totalChanceToHit * healthDamage,
		BarrierDamageTaken:            barrierDamage + extraBarrierDamage,
		ExpectedBarrierDamage:         totalChanceToHit * (barrierDamage + extraBarrierDamage),
		ChanceToCrit:                  chanceToCritical,
		CriticalDamageTaken:           criticalHealthDamage,
		CriticalBarrierDamageTaken:    criticalBarrierDamage + criticalExtraBarrierDamage,
		CriticalExpectedDamage:        totalChanceToHit * criticalHealthDamage,
		CriticalExpectedBarrierDamage: totalChanceToHit * (criticalBarrierDamage + criticalExtraBarrierDamage),
	}
}

// GetPowerToHitPenaltyAgainstSquaddie calculates how much the target can reduce the chance of getting hit by the attacking power.
func GetPowerToHitPenaltyAgainstSquaddie(power *entity.Power, target *entity.Squaddie) (toHitPenalty int) {
	if power.PowerType == entity.PowerTypePhysical {
		return target.Dodge
	}
	return target.Deflect
}
