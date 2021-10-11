package terosbattleserver

import (
	"fmt"
	"github.com/chadius/terosbattleserver/entity/actioncontroller"
	"github.com/chadius/terosbattleserver/entity/actionviewer"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/replay"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/utility"
	"io/ioutil"
	"log"
)

func ReplayBattleScript(scriptFilename, squaddieRepositoryFilename, powerRepositoryFilename string) {
	utility.Logger = &utility.FileLogger{}
	squaddieRepo := loadSquaddieRepo(squaddieRepositoryFilename)
	powerRepo := loadPowerRepo(powerRepositoryFilename)
	repos := &repositories.RepositoryCollection{
		PowerRepo:    powerRepo,
		SquaddieRepo: squaddieRepo,
	}

	controller := actioncontroller.WhiteRoomController{}
	viewer := actionviewer.ConsoleActionViewer{}

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

func loadSquaddieRepo(squaddieRepositoryFilename string) (repo *squaddie.Repository) {
	squaddieYamlData, err := ioutil.ReadFile(squaddieRepositoryFilename)
	if err != nil {
		log.Fatal(err)
	}

	squaddieRepo := squaddie.NewSquaddieRepository()
	squaddieRepo.AddYAMLSource(squaddieYamlData)
	return squaddieRepo
}

func loadPowerRepo(powerRepositoryFilename string) (repo *power.Repository) {
	powerYamlData, err := ioutil.ReadFile(powerRepositoryFilename)
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
