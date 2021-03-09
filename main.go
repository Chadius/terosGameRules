package main

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerattack"
)

func main() {
	attacker, target, power := loadActors()

	attackingPowerSummary := powerattack.GetExpectedDamage(power, attacker, target)
	println(attacker.Name, "will attack", target.Name, "with", power.Name)
	println("Chance to hit (out of 36) ", attackingPowerSummary.ChanceToHit)
	println("Damage taken              ", attackingPowerSummary.DamageTaken)
	println("Barrier damage            ", attackingPowerSummary.BarrierDamageTaken)
	println("---")
	println("Expected damage (36 = 1HP)", attackingPowerSummary.ExpectedDamage)
	println("Expected barrier damage   ", attackingPowerSummary.ExpectedBarrierDamage)
}

func loadActors () (*squaddie.Squaddie, *squaddie.Squaddie, *power.Power) {
	teros := squaddie.NewSquaddie("Teros")
	teros.Strength = 1
	teros.Armor = 2
	teros.Dodge = 3
	teros.Deflect = 4
	teros.MaxBarrier = 1
	teros.SetBarrierToMax()

	bandit := squaddie.NewSquaddie("Bandit")
	bandit.Armor = 2
	bandit.Dodge = 0
	bandit.Deflect = 0
	bandit.MaxBarrier = 0
	bandit.SetBarrierToMax()

	powerRepository := power.NewPowerRepository()
	spear := power.NewPower("Spear")
	spear.PowerType = power.Physical
	spear.ID = "deadbeef"
	spear.DamageBonus = 2
	spear.ToHitBonus = 1
	newPowers := []*power.Power{spear}
	powerRepository.AddSlicePowerSource(newPowers)

	temporaryPowerReferences := []*power.Reference{{Name: "Spear", ID: spear.ID}}
	powerattack.LoadAllOfSquaddieInnatePowers(teros, temporaryPowerReferences, powerRepository)

	attacker := teros
	power := spear
	target := bandit

	return attacker, target, power
}