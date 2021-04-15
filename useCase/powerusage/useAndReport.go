package powerusage

import (
	"github.com/cserrant/terosBattleServer/entity/powerusagecontext"
	"github.com/cserrant/terosBattleServer/entity/report"
	"github.com/cserrant/terosBattleServer/utility"
)

// UsePowerAgainstSquaddiesAndGetResults will make the actingSquaddie use the powerUsed against all targetSquaddies.
//   Returns a report indicating what happened to each target.
func UsePowerAgainstSquaddiesAndGetResults(context *powerusagecontext.PowerUsageContext, d6generator utility.SixSideGenerator) *report.PowerReport {
	powerResults := &report.PowerReport{
		AttackerID:            context.ActingSquaddieID,
		PowerID:               context.PowerID,
		AttackingPowerResults: []*report.AttackingPowerReport{},
	}

	for _, targetSquaddieID := range context.TargetSquaddieIDs {
		attackingResult := GetAttackEffectResults(context, targetSquaddieID, d6generator)
		powerResults.AttackingPowerResults = append(powerResults.AttackingPowerResults, attackingResult)
	}
	return powerResults
}

// GetAttackEffectResults looks at the actingSquaddie's powerUsed's AttackingEffect to figure out what happened to the targetSquaddie.
func GetAttackEffectResults(context *powerusagecontext.PowerUsageContext, targetSquaddieID string, d6generator utility.SixSideGenerator) *report.AttackingPowerReport {
	attackForecast := GetExpectedDamage(
		context,
		&powerusagecontext.AttackContext{
			PowerID:			context.PowerID,
			AttackerID:			context.ActingSquaddieID,
			TargetID:			targetSquaddieID,
			IsCounterAttack: 	false,
		},
	)

	if !DetermineIfItHit(attackForecast, d6generator) {
		return &report.AttackingPowerReport{
			TargetID:        targetSquaddieID,
			DamageTaken:     0,
			BarrierDamage:   0,
			WasAHit:         false,
			WasACriticalHit: false,
		}
	}

	if !DetermineIfItWasACriticalHit(attackForecast, d6generator) {
		return &report.AttackingPowerReport{
			TargetID:        targetSquaddieID,
			DamageTaken:     attackForecast.DamageTaken,
			BarrierDamage:   attackForecast.BarrierDamageTaken,
			WasAHit:         true,
			WasACriticalHit: false,
		}
	}

	return &report.AttackingPowerReport{
		TargetID:        targetSquaddieID,
		DamageTaken:     attackForecast.CriticalDamageTaken,
		BarrierDamage:   attackForecast.CriticalBarrierDamageTaken,
		WasAHit:         true,
		WasACriticalHit: true,
	}
}

// DetermineIfItHit rolls attacks and determines if the attack hit.
func DetermineIfItHit(summary *powerusagecontext.AttackingPowerForecast, d6generator utility.SixSideGenerator) bool {
	hitRate := summary.HitRate
	attackRoll, defendRoll := d6generator.RollTwoDice()
	return attackRoll + hitRate >= defendRoll
}

// DetermineIfItWasACriticalHit rolls and determines if the attack was a critical hit.
func DetermineIfItWasACriticalHit(summary *powerusagecontext.AttackingPowerForecast, d6generator utility.SixSideGenerator) bool {
	criticalHitThreshold := summary.CriticalHitThreshold
	roll1, roll2 := d6generator.RollTwoDice()
	return roll1 + roll2 < criticalHitThreshold
}