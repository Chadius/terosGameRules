package main

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerattack"
	"io/ioutil"
	"log"
)

func main() {
	attacker, target, power := loadActors(
		"squaddieTeros",
		"squaddieBandit",
		"powerSpear",
	)

	attackingPowerSummary := powerattack.GetExpectedDamage(power, attacker, target)
	println(attacker.Name, "will attack", target.Name, "with", power.Name)
	println("Chance to hit (out of 36) ", attackingPowerSummary.ChanceToHit)
	println("Damage taken              ", attackingPowerSummary.DamageTaken)
	println("Barrier damage            ", attackingPowerSummary.BarrierDamageTaken)
	println("---")
	println("Expected damage (36 = 1HP)", attackingPowerSummary.ExpectedDamage)
	println("Expected barrier damage   ", attackingPowerSummary.ExpectedBarrierDamage)
}

func loadActors (attackerID, targetID, powerID string) (*squaddie.Squaddie, *squaddie.Squaddie, *power.Power) {
	squaddieYamlData, err := ioutil.ReadFile("data/squaddieDatabase.yml")
	if err != nil {
		log.Fatal(err)
	}

	squaddieRepo := squaddie.NewSquaddieRepository()
	squaddieRepo.AddYAMLSource(squaddieYamlData)
	attacker := squaddieRepo.GetByID(attackerID)
	attacker.SetBarrierToMax()

	target := squaddieRepo.GetByID(targetID)
	target.SetBarrierToMax()

	powerYamlData, err := ioutil.ReadFile("data/powerDatabase.yml")
	if err != nil {
		log.Fatal(err)
	}
	powerRepo := power.NewPowerRepository()
	powerRepo.AddYAMLSource(powerYamlData)

	powerattack.LoadAllOfSquaddieInnatePowers(attacker, attacker.PowerReferences, powerRepo)

	power := powerRepo.GetPowerByID(powerID)

	return attacker, target, power
}