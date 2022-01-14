package powerattackforecast

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/chadius/terosgamerules/entity/damagedistribution"
)

//counterfeiter:generate . VersusContextStrategy
// VersusContextStrategy describes objects that compare attack and defender contexts to figure out what happens when they attack.
type VersusContextStrategy interface {
	Calculate(attackerContext AttackerContext, defenderContext DefenderContext)
	ToHit() *damagedistribution.ToHitComparison
	NormalDamage() *damagedistribution.DamageDistribution
	CriticalHitDamage() *damagedistribution.DamageDistribution
	CanCritical() bool
	CriticalHitThreshold() int
}

// VersusContext compares an AttackerContext and DefenderContext to determine the possible results.
type VersusContext struct {
	toHit                *damagedistribution.ToHitComparison
	normalDamage         *damagedistribution.DamageDistribution
	criticalHitDamage    *damagedistribution.DamageDistribution
	canCritical          bool
	criticalHitThreshold int
}

// Calculate figures out what the results are.
func (context *VersusContext) Calculate(attackerContext AttackerContext, defenderContext DefenderContext) {
	context.calculateToHitBonus(attackerContext, defenderContext)
	context.setNormalDamageBreakdown(attackerContext, defenderContext)

	context.setCriticalHitChance(attackerContext)
	context.setCriticalDamageBreakdown(attackerContext, defenderContext)
}

func (context *VersusContext) calculateToHitBonus(attackerContext AttackerContext, defenderContext DefenderContext) {
	context.toHit = &damagedistribution.ToHitComparison{}
	context.toHit.AttackerToHitBonus = attackerContext.TotalToHitBonus()
	context.toHit.DefenderToHitPenalty = defenderContext.TotalToHitPenalty()
	context.toHit.ToHitBonus = context.toHit.AttackerToHitBonus - context.toHit.DefenderToHitPenalty
}

func (context *VersusContext) setNormalDamageBreakdown(attackerContext AttackerContext, defenderContext DefenderContext) {
	context.normalDamage = context.setDamageBreakdown(attackerContext.RawDamage(), attackerContext, defenderContext)
}

func (context *VersusContext) setCriticalDamageBreakdown(attackerContext AttackerContext, defenderContext DefenderContext) {
	if context.canCritical {
		context.criticalHitDamage = context.setDamageBreakdown(attackerContext.CriticalHitDamage(), attackerContext, defenderContext)
	}
}

func (context *VersusContext) setDamageBreakdown(damageDealtToTarget int, attackerContext AttackerContext, defenderContext DefenderContext) *damagedistribution.DamageDistribution {
	distribution := &damagedistribution.DamageDistribution{}

	context.setBarrierBurntAndDamageAbsorbed(distribution, attackerContext, defenderContext, damageDealtToTarget)

	damageDealtToTarget -= distribution.DamageAbsorbedByBarrier
	distribution.TotalRawBarrierBurnt = distribution.DamageAbsorbedByBarrier + distribution.ExtraBarrierBurnt

	distribution.DamageAbsorbedByArmor = context.calculateDamageAbsorbedByArmor(attackerContext, defenderContext, damageDealtToTarget)
	damageDealtToTarget -= distribution.DamageAbsorbedByArmor

	distribution.RawDamageDealt = damageDealtToTarget
	distribution.IsFatalToTarget = distribution.RawDamageDealt >= defenderContext.HitPoints()

	return distribution
}

func (context *VersusContext) calculateDamageAbsorbedByArmor(attackerContext AttackerContext, defenderContext DefenderContext, damageDealtToTarget int) int {
	armorAbsorbsAllDamage := damageDealtToTarget <= defenderContext.ArmorResistance()
	if armorAbsorbsAllDamage {
		return damageDealtToTarget
	}
	return defenderContext.ArmorResistance()
}

func (context *VersusContext) setBarrierBurntAndDamageAbsorbed(distribution *damagedistribution.DamageDistribution, attackerContext AttackerContext, defenderContext DefenderContext, damageDealtToTarget int) {
	barrierAbsorbsAllDamageAndExtraBurn := damageDealtToTarget+attackerContext.ExtraBarrierBurn() <= defenderContext.BarrierResistance()
	if barrierAbsorbsAllDamageAndExtraBurn {
		distribution.ExtraBarrierBurnt = attackerContext.ExtraBarrierBurn()
		distribution.DamageAbsorbedByBarrier = damageDealtToTarget
		distribution.TotalRawBarrierBurnt = distribution.DamageAbsorbedByBarrier + distribution.ExtraBarrierBurnt
		return
	}

	barrierAbsorbsExtraBarrierBurn := attackerContext.ExtraBarrierBurn() <= defenderContext.BarrierResistance()
	if !barrierAbsorbsExtraBarrierBurn {
		distribution.ExtraBarrierBurnt = defenderContext.BarrierResistance()
		distribution.DamageAbsorbedByBarrier = 0
		distribution.TotalRawBarrierBurnt = distribution.ExtraBarrierBurnt
		return
	}

	barrierRemainingAfterExtraBarrierBurn := defenderContext.BarrierResistance() - attackerContext.ExtraBarrierBurn()

	remainingBarrierAbsorbsDamage := damageDealtToTarget <= barrierRemainingAfterExtraBarrierBurn
	if remainingBarrierAbsorbsDamage {
		distribution.ExtraBarrierBurnt = attackerContext.ExtraBarrierBurn()
		distribution.DamageAbsorbedByBarrier = damageDealtToTarget
		distribution.TotalRawBarrierBurnt = distribution.DamageAbsorbedByBarrier + distribution.ExtraBarrierBurnt
		return
	}

	distribution.ExtraBarrierBurnt = attackerContext.ExtraBarrierBurn()
	distribution.DamageAbsorbedByBarrier = barrierRemainingAfterExtraBarrierBurn
	distribution.TotalRawBarrierBurnt = defenderContext.BarrierResistance()
}

func (context *VersusContext) setCriticalHitChance(attackerContext AttackerContext) {
	context.canCritical = attackerContext.CanCritical()
	if context.canCritical {
		context.criticalHitThreshold = attackerContext.CriticalHitThreshold()
	} else {
		context.criticalHitThreshold = 0
	}
}

// ToHit is a getter.
func (context *VersusContext) ToHit() *damagedistribution.ToHitComparison {
	return context.toHit
}

// NormalDamage is a getter.
func (context *VersusContext) NormalDamage() *damagedistribution.DamageDistribution {
	return context.normalDamage
}

// CriticalHitDamage is a getter.
func (context *VersusContext) CriticalHitDamage() *damagedistribution.DamageDistribution {
	return context.criticalHitDamage
}

// CanCritical is a getter.
func (context *VersusContext) CanCritical() bool {
	return context.canCritical
}

// CriticalHitThreshold is a getter.
func (context *VersusContext) CriticalHitThreshold() int {
	return context.criticalHitThreshold
}
