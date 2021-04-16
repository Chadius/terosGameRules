package powercommit

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagecontext"
	"github.com/cserrant/terosBattleServer/entity/report"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/usecase/powerforecast"
	"github.com/cserrant/terosBattleServer/utility"
)

// UsePowerAgainstSquaddiesAndGetReport will make the actingSquaddie use the powerUsed against all targetSquaddies.
//   Returns a report indicating what happened to each target.
func UsePowerAgainstSquaddiesAndGetReport(context *powerusagecontext.PowerUsageContext, d6generator utility.SixSideGenerator) *report.PowerReport {
	powerReport := &report.PowerReport{
		AttackerID:            context.ActingSquaddieID,
		PowerID:               context.PowerID,
		AttackingPowerReports: []*report.AttackingPowerReport{},
	}

	for _, targetSquaddieID := range context.TargetSquaddieIDs {
		attackingReports := calculateAllAttackPowerReportsForThisAttack(context, targetSquaddieID, d6generator)
		powerReport.AttackingPowerReports = append(powerReport.AttackingPowerReports, attackingReports...)
	}
	return powerReport
}

// calculateAllAttackPowerReportsForThisAttack forecasts the attack from the context
//	and figures out all secondary effects.
func calculateAllAttackPowerReportsForThisAttack(context *powerusagecontext.PowerUsageContext, targetSquaddieID string, d6generator utility.SixSideGenerator) []*report.AttackingPowerReport {
	attackForecast := powerforecast.GetExpectedDamage(
		context,
		&powerusagecontext.AttackContext{
			PowerID:         context.PowerID,
			AttackerID:      context.ActingSquaddieID,
			TargetID:        targetSquaddieID,
			IsCounterAttack: false,
		},
	)

	reports := []*report.AttackingPowerReport{
		calculateAttackPowerReportFromForecast(attackForecast, d6generator),
	}

	if attackForecast.CounterAttack != nil {
		counterattackReport := calculateAttackPowerReportFromForecast(attackForecast.CounterAttack, d6generator)
		reports = append(reports, counterattackReport)
	}

	return reports
}

func calculateAttackPowerReportFromForecast(attackForecast *powerusagecontext.AttackingPowerForecast, d6generator utility.SixSideGenerator) *report.AttackingPowerReport {
	attackReport := &report.AttackingPowerReport{
		AttackerID:			attackForecast.AttackingSquaddieID,
		TargetID:			attackForecast.TargetSquaddieID,
		PowerID:			attackForecast.PowerID,
		DamageTaken:		0,
		BarrierDamage:		0,
		WasAHit:			false,
		WasACriticalHit:	false,
	}

	if !DetermineIfItHit(attackForecast, d6generator) {
		return attackReport
	}

	if !DetermineIfItWasACriticalHit(attackForecast, d6generator) {
		attackReport.WasAHit = true
		attackReport.DamageTaken = attackForecast.DamageTaken
		attackReport.BarrierDamage = attackForecast.BarrierDamageTaken
		return attackReport
	}

	attackReport.WasAHit = true
	attackReport.WasACriticalHit = true
	attackReport.DamageTaken = attackForecast.CriticalDamageTaken
	attackReport.BarrierDamage = attackForecast.CriticalBarrierDamageTaken
	return attackReport
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

// CommitPowerUse will apply the given PowerReport.
//    Squaddies will move, Targets will take damage, etc.
func CommitPowerUse(powerReport *report.PowerReport, squaddieRepo *squaddie.Repository, powerRepo *power.Repository) {
	squaddieToEquip := squaddieRepo.GetOriginalSquaddieByID(powerReport.AttackerID)
	powerequip.SquaddieEquipPower(squaddieToEquip, powerReport.PowerID, powerRepo)
}