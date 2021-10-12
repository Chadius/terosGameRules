package terosbattleserver

import (
	"bufio"
	"fmt"
	"github.com/chadius/terosbattleserver/entity/actioncontroller"
	"github.com/chadius/terosbattleserver/entity/actionviewer"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/replay"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/utility"
	"log"
	"os"
)

func ReplayBattleScript(scriptFileHandle, squaddieFileHandle, powerFileHandle *os.File) error {
	utility.Logger = &utility.FileLogger{}

	var squaddieData = []byte{}
	squaddieScanner := bufio.NewScanner(squaddieFileHandle)
	squaddieScanner.Split(bufio.ScanBytes)
	for squaddieScanner.Scan() {
		squaddieData = append(squaddieData, squaddieScanner.Bytes()...)
	}
	squaddieRepo := loadSquaddieRepo(squaddieData)

	var powerData = []byte{}
	powerScanner := bufio.NewScanner(powerFileHandle)
	powerScanner.Split(bufio.ScanBytes)
	for powerScanner.Scan() {
		powerData = append(powerData, powerScanner.Bytes()...)
	}
	powerRepo := loadPowerRepo(powerData)

	repos := &repositories.RepositoryCollection{
		PowerRepo:    powerRepo,
		SquaddieRepo: squaddieRepo,
	}

	controller := actioncontroller.WhiteRoomController{}
	viewer := actionviewer.ConsoleActionViewer{}

	var scriptData = []byte{}
	scriptScanner := bufio.NewScanner(scriptFileHandle)
	scriptScanner.Split(bufio.ScanBytes)
	for scriptScanner.Scan() {
		scriptData = append(scriptData, scriptScanner.Bytes()...)
	}
	chapterReplay, replayErr := replay.NewCreateMapReplayFromYAML(scriptData)
	if replayErr != nil {
		return replayErr
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
	return nil
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

func loadSquaddieRepo(squaddieYamlData []byte) (repo *squaddie.Repository) {
	squaddieRepo := squaddie.NewSquaddieRepository()
	_, err := squaddieRepo.AddYAMLSource(squaddieYamlData)
	if err != nil {
		println (err.Error())
	}
	return squaddieRepo
}

func loadPowerRepo(powerYamlData []byte) (repo *power.Repository) {
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
