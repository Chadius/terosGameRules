package main

import (
	"flag"
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/actioncontroller"
	"github.com/cserrant/terosBattleServer/entity/actionviewer"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/replay"
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
		PowerRepo:    powerRepo,
		SquaddieRepo: squaddieRepo,
	}

	controller := actioncontroller.WhiteRoomController{}
	viewer := actionviewer.ConsoleActionViewer{}

	scriptFilename := "scripts/battle.yml"
	flag.StringVar(&scriptFilename, "f", "scripts/battle.yml", "The filename of the script file. Defaults to scripts/battle.yml")
	flag.Parse()

	scriptYAML, err := ioutil.ReadFile(scriptFilename)
	if err != nil {
		log.Fatal(err)
	}
	chapterReplay, replayErr := replay.NewCreateMapReplayFromYAML(scriptYAML)
	if replayErr != nil {
		log.Fatal(replayErr)
	}

	initializeAllSquaddies(chapterReplay, repos)
	for _, action := range chapterReplay.Actions {
		continueProcessing := processSquaddieAction(
			action,
			&viewer,
			&controller,
			repos,
		)

		if continueProcessing == false {
			break
		}
	}
}

func processSquaddieAction(
	action *replay.SquaddieAction,
	viewer *actionviewer.ConsoleActionViewer,
	controller *actioncontroller.WhiteRoomController,
	repositories *repositories.RepositoryCollection) bool {

	powerSetup := controller.SetupAction(action.UserID, action.TargetIDs, action.PowerID)

	reasonsForInvalidAction := controller.CheckForValidAction(powerSetup, repositories)
	if len(reasonsForInvalidAction) > 0 {
		for _, reason := range reasonsForInvalidAction {
			for _, description := range reason.Description {
				println(description)
			}
		}
		return false
	}

	forecast := controller.GenerateForecast(powerSetup, repositories)
	viewer.PrintForecast(forecast, repositories)

	println()
	result := controller.GenerateResult(forecast, repositories, true, action.RandomSeed)
	viewer.PrintResult(result, repositories, &actionviewer.ConsoleActionViewerVerbosity{
		ShowTargetStatus: true,
	})

	return true
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

func initializeAllSquaddies(replay *replay.ChapterReplay, repositories *repositories.RepositoryCollection) {
	squaddiesFound := map[string]bool{}

	for _, action := range replay.Actions {
		if squaddiesFound[action.UserID] != true {
			loadAndInitializeSquaddie(action.UserID, repositories)
			squaddiesFound[action.UserID] = true
		}
		for _, targetID := range action.TargetIDs {
			if squaddiesFound[targetID] != true {
				loadAndInitializeSquaddie(targetID, repositories)
				squaddiesFound[targetID] = true
			}
		}
	}
}

func loadAndInitializeSquaddie(squaddieID string, repositories *repositories.RepositoryCollection) {
	squaddieRepo := repositories.SquaddieRepo
	squaddie := squaddieRepo.GetOriginalSquaddieByID(squaddieID)
	if squaddie == nil {
		fmt.Sprintf("Squaddie %s does not exist, exiting", squaddieID)
		log.Panicf("Squaddie %s does not exist, exiting", squaddieID)
	}
	squaddie.Defense.SetBarrierToMax()
	powerequip.LoadAllOfSquaddieInnatePowers(squaddie, squaddie.PowerCollection.PowerReferences, repositories)
	powerequip.EquipDefaultPower(squaddie, repositories)
}