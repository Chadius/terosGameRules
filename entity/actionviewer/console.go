package actionviewer

import (
	"github.com/cserrant/terosBattleServer/entity/powerusagescenario"
	"github.com/cserrant/terosBattleServer/usecase/powerattackforecast"
	"github.com/cserrant/terosBattleServer/usecase/powercommit"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
	"strconv"
)

// ConsoleActionViewer prints the results of actions to the console
type ConsoleActionViewer struct {}

// PrintActionForecast outputs the calculated forecast.
func (controller *ConsoleActionViewer) PrintActionForecast(action *powerattackforecast.Forecast, repos *repositories.RepositoryCollection ) {
	powerToUse := repos.PowerRepo.GetPowerByID(action.Setup.PowerID)
	if powerToUse.AttackEffect != nil {
		controller.printAllAttackForecasts(action, repos)
	}
	if powerToUse.HealingEffect != nil {
		controller.printAllHealingForecasts(action, repos)
	}
}

// PrintActionResults outputs the calculated results.
func (controller *ConsoleActionViewer) PrintActionResults(result *powercommit.Result, repos *repositories.RepositoryCollection ) {
	powerToUse := repos.PowerRepo.GetPowerByID(result.Forecast.Setup.PowerID)
	if powerToUse.AttackEffect != nil {
		controller.printAllAttackResults(result, repos)
	}
	if powerToUse.HealingEffect != nil {
		controller.printAllHealingResults(result, repos)
	}
}

func (controller *ConsoleActionViewer) printAllAttackForecasts(action *powerattackforecast.Forecast, repos *repositories.RepositoryCollection ) {
	for _, forecast := range action.ForecastedResultPerTarget {
		printAttackForecastPerTarget(&forecast)
	}
}

func printAttackForecastPerTarget(forecast *powerattackforecast.Calculation) {
	printAttackIterationForecast(forecast.Attack, forecast.Setup, forecast.Repositories)
	if forecast.CounterAttack != nil {
		println("")
		println("then Counterattack:")
		printAttackIterationForecast(forecast.CounterAttack, forecast.CounterAttackSetup, forecast.Repositories)
	}
}

func printAttackIterationForecast(forecast *powerattackforecast.AttackForecast, setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) {
	squaddieRepo := repositories.SquaddieRepo
	powerRepo := repositories.PowerRepo

	attacker := squaddieRepo.GetOriginalSquaddieByID(setup.UserID)
	target := squaddieRepo.GetOriginalSquaddieByID(setup.Targets[0])
	attackingPower := powerRepo.GetPowerByID(setup.PowerID)

	println(attacker.Identification.Name, "will attack", target.Identification.Name, "with", attackingPower.Name)
	println("Attacker ToHit bonus", forecast.VersusContext.ToHit.ToHitBonus)

	if forecast.VersusContext.NormalDamage.IsFatalToTarget {
		println("will kill if it hits")
	}

	println("Forecasted Damage              ", forecast.VersusContext.NormalDamage.RawDamageDealt)
}

func (controller *ConsoleActionViewer) printAllAttackResults(result *powercommit.Result, repos *repositories.RepositoryCollection ) {
	for _, attackReport := range result.ResultPerTarget {
		controller.printAttackPerTarget(attackReport, repos)
		println()
	}
}

func (controller *ConsoleActionViewer) printAttackPerTarget(result *powercommit.ResultPerTarget, repositories *repositories.RepositoryCollection) {
	squaddieRepo := repositories.SquaddieRepo
	powerRepo := repositories.PowerRepo

	attacker := squaddieRepo.GetOriginalSquaddieByID(result.UserID)
	target := squaddieRepo.GetOriginalSquaddieByID(result.TargetID)
	attackingPower := powerRepo.GetPowerByID(result.PowerID)

	println(attacker.Identification.Name, "attacks", target.Identification.Name, "with", attackingPower.Name)

	println(attacker.Identification.Name, "attacks with a", result.Attack.AttackRoll, "+", result.Attack.AttackerToHitBonus, "=", result.Attack.AttackerTotal)
	println(target.Identification.Name, "defends with a", result.Attack.DefendRoll, "+", result.Attack.DefenderToHitPenalty, "=", result.Attack.DefenderTotal)
	if !result.Attack.HitTarget {
		println("Missed")
		return
	}

	if result.Attack.CriticallyHitTarget {
		println("Critical Hit")
	} else {
		println("Hit")
	}
	damageTaken := "  deals " + strconv.Itoa(result.Attack.Damage.RawDamageDealt)
	if result.Attack.Damage.TotalRawBarrierBurnt > 0 {
		damageTaken += " damage, " + strconv.Itoa(result.Attack.Damage.TotalRawBarrierBurnt) + " barrier burn"
	}
	println(damageTaken)

	healthStatus := target.Identification.Name + " HP: " + strconv.Itoa(target.Defense.CurrentHitPoints) + "/" + strconv.Itoa(target.Defense.MaxHitPoints)
	if target.Defense.CurrentBarrier > 0  {
		healthStatus += "Barrier" + strconv.Itoa(target.Defense.CurrentBarrier)
	}
	println(healthStatus)

	if target.Defense.IsDead() {
		println(target.Identification.Name, "falls!")
	}
}

func (controller *ConsoleActionViewer) printAllHealingForecasts(action *powerattackforecast.Forecast, repos *repositories.RepositoryCollection) {
	for _, calculation := range action.ForecastedResultPerTarget {
		controller.printHealingForecastPerTarget(calculation.HealingForecast, &action.Setup, repos)
	}
}

func (controller *ConsoleActionViewer) printHealingForecastPerTarget(forecast *powerattackforecast.HealingForecast, setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) {
	squaddieRepo := repositories.SquaddieRepo
	powerRepo := repositories.PowerRepo

	attacker := squaddieRepo.GetOriginalSquaddieByID(setup.UserID)
	target := squaddieRepo.GetOriginalSquaddieByID(setup.Targets[0])
	healingPower := powerRepo.GetPowerByID(setup.PowerID)

	println(attacker.Identification.Name, "will heal", target.Identification.Name, "with", healingPower.Name)

	println("Forecasted Healing              ", forecast.RawHitPointsRestored)
}

func (controller *ConsoleActionViewer) printAllHealingResults(powerResult *powercommit.Result, repos *repositories.RepositoryCollection) {
	for _, healingReport := range powerResult.ResultPerTarget {
		controller.printHealingReportPerTarget(healingReport, repos)
		println()
	}
}

func (controller *ConsoleActionViewer) printHealingReportPerTarget(result *powercommit.ResultPerTarget, repositories *repositories.RepositoryCollection) {
	squaddieRepo := repositories.SquaddieRepo
	powerRepo := repositories.PowerRepo

	healer := squaddieRepo.GetOriginalSquaddieByID(result.UserID)
	target := squaddieRepo.GetOriginalSquaddieByID(result.TargetID)
	healingPower := powerRepo.GetPowerByID(result.PowerID)

	println("---")

	println(healer.Identification.Name, "heals", target.Identification.Name, "with", healingPower.Name)

	println("  heals ", result.Healing.HitPointsRestored)
}