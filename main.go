package main

import (
	"github.com/cserrant/terosBattleServer/entity/actioncontroller"
	"github.com/cserrant/terosBattleServer/entity/actionviewer"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
	"github.com/cserrant/terosBattleServer/utility"
	"io/ioutil"
	"log"
)

func main() {
	utility.Logger = &utility.FileLogger{}
	squaddieRepo := loadSquaddieRepo()
	powerRepo := loadPowerRepo()
	repos := &repositories.RepositoryCollection{
		PowerRepo: powerRepo,
		SquaddieRepo: squaddieRepo,
	}

	attacker, target, power := loadActors(
		"squaddieTeros",
		"squaddieBandit",
		"powerSpear",
		repos,
	)

	controller := actioncontroller.WhiteRoomController{}
	viewer := actionviewer.ConsoleActionViewer{}
	powerSetup := controller.SetupAction(attacker, target, power)

	reasonsForInvalidAction := controller.CheckForValidAction(powerSetup, repos)
	if len(reasonsForInvalidAction) > 0 {
		for _, reason := range reasonsForInvalidAction {
			for _, description := range reason.Description {
				println(description)
			}
		}
		return
	}

	forecast := controller.GenerateForecast(powerSetup, repos)
	viewer.PrintForecast(forecast, repos)

	println()
	result := controller.GenerateResult(forecast, repos)
	viewer.PrintResult(result, repos, &actionviewer.ConsoleActionViewerVerbosity{
		ShowTargetStatus: true,
	})
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

func loadActors (attackerID, targetID, powerID string, repositories *repositories.RepositoryCollection) (*squaddie.Squaddie, *squaddie.Squaddie, *power.Power) {
	squaddieRepo := repositories.SquaddieRepo
	powerRepo := repositories.PowerRepo
	attacker := squaddieRepo.GetOriginalSquaddieByID(attackerID)
	attacker.Defense.SetBarrierToMax()
	powerequip.LoadAllOfSquaddieInnatePowers(attacker, attacker.PowerCollection.PowerReferences, repositories)
	powerequip.EquipDefaultPower(attacker, repositories)

	target := squaddieRepo.GetOriginalSquaddieByID(targetID)
	target.Defense.SetBarrierToMax()
	powerequip.LoadAllOfSquaddieInnatePowers(target, target.PowerCollection.PowerReferences, repositories)
	powerequip.EquipDefaultPower(target, repositories)

	power := powerRepo.GetPowerByID(powerID)
	return attacker, target, power
}