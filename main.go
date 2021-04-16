package main

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagecontext"
	"github.com/cserrant/terosBattleServer/entity/report"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powercommit"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/usecase/powerforecast"
	"github.com/cserrant/terosBattleServer/utility"
	"io/ioutil"
	"log"
)

func main() {
	squaddieRepo := loadSquaddieRepo()
	powerRepo := loadPowerRepo()

	attacker, target, power := loadActors(
		"squaddieTeros",
		"squaddieBandit",
		"powerSpear",
		squaddieRepo,
		powerRepo,
	)

	powerContext := &powerusagecontext.PowerUsageContext{
		SquaddieRepo:      squaddieRepo,
		ActingSquaddieID:  attacker.ID,
		TargetSquaddieIDs: []string{target.ID},
		PowerID:           power.ID,
		PowerRepo:         powerRepo,
	}

	powerForecast := powerforecast.CalculatePowerForecast(powerContext)

	for _, attackingPowerForecast := range powerForecast.AttackPowerForecast {
		printAttackForecast(attackingPowerForecast, squaddieRepo, powerRepo)
	}

	println("---")
	dieRoller := &utility.RandomDieRoller{}
	attackReports := powercommit.UsePowerAgainstSquaddiesAndGetReport(powerContext, dieRoller)
	for _, attackReport := range attackReports.AttackingPowerReports {
		printAttackReport(attackReport)
	}
}

func printAttackForecast(forecast *powerusagecontext.AttackingPowerForecast,
	squaddieRepo *squaddie.Repository,
	powerRepo *power.Repository) {
	printPartOfAttackForecast(forecast, squaddieRepo, powerRepo)
	if forecast.CounterAttack != nil {
		println("")
		println("Counterattack:")
		printPartOfAttackForecast(forecast.CounterAttack, squaddieRepo, powerRepo)
	}
}


func printPartOfAttackForecast(forecast *powerusagecontext.AttackingPowerForecast,
		squaddieRepo *squaddie.Repository,
		powerRepo *power.Repository) {
	attacker := squaddieRepo.GetOriginalSquaddieByID(forecast.AttackingSquaddieID)
	target := squaddieRepo.GetOriginalSquaddieByID(forecast.TargetSquaddieID)
	power := powerRepo.GetPowerByID(forecast.PowerID)
	println(attacker.Name, "will attack", target.Name, "with", power.Name)
	println("Chance to hit (out of 36) ", forecast.ChanceToHit)
	println("Damage taken              ", forecast.DamageTaken)
	println("Barrier damage            ", forecast.BarrierDamageTaken)
	println("---")
	println("Expected damage (36 = 1HP)", forecast.ExpectedDamage)
	println("Expected barrier damage   ", forecast.ExpectedBarrierDamage)
}

func printAttackReport(report *report.AttackingPowerReport) {
	if !report.WasAHit {
		println("Missed")
	} else if report.WasACriticalHit {
		println("Critical Hit")
		println(report.DamageTaken, "damage taken")
		println(report.BarrierDamage, "barrier damage")
	} else {
		println("Hit")
		println(report.DamageTaken, "damage taken")
		println(report.BarrierDamage, "barrier damage")
	}
}

func loadSquaddieRepo() (repo *squaddie.Repository) {
	squaddieYamlData, err := ioutil.ReadFile("data/squaddieDatabase.yml")
	if err != nil {
		log.Fatal(err)
	}

	squaddieRepo := squaddie.NewSquaddieRepository()
	squaddieRepo.AddYAMLSource(squaddieYamlData)
	return squaddieRepo
}

func loadPowerRepo() (repo *power.Repository) {
	powerYamlData, err := ioutil.ReadFile("data/powerDatabase.yml")
	if err != nil {
		log.Fatal(err)
	}
	powerRepo := power.NewPowerRepository()
	powerRepo.AddYAMLSource(powerYamlData)
	return powerRepo
}

func loadActors (attackerID, targetID, powerID string, squaddieRepo *squaddie.Repository, powerRepo *power.Repository) (*squaddie.Squaddie, *squaddie.Squaddie, *power.Power) {
	attacker := squaddieRepo.GetOriginalSquaddieByID(attackerID)
	attacker.SetBarrierToMax()

	target := squaddieRepo.GetOriginalSquaddieByID(targetID)
	target.SetBarrierToMax()

	powerequip.LoadAllOfSquaddieInnatePowers(attacker, attacker.PowerReferences, powerRepo)

	power := powerRepo.GetPowerByID(powerID)

	return attacker, target, power
}