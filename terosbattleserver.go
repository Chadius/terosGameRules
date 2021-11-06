package terosbattleserver

import (
	"fmt"
	"github.com/chadius/terosbattleserver/entity/actioncontroller"
	"github.com/chadius/terosbattleserver/entity/actionviewer"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	"github.com/chadius/terosbattleserver/entity/replay"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/utility"
	"io"
	"io/ioutil"
	"log"
)

// ReplayBattleScript uses the input streams to read and replay several rounds of combat,
//  writing the results to a supplied output stream.
func ReplayBattleScript(scriptFileHandle, squaddieFileHandle, powerFileHandle io.Reader, output io.Writer) error {
	utility.Logger = &utility.FileLogger{}

	squaddieRepo := createSquaddieRepo(squaddieFileHandle)
	powerRepo := createPowerRepo(powerFileHandle)
	chapterReplay := createChapterReplay(scriptFileHandle)

	repos := &repositories.RepositoryCollection{
		PowerRepo:    powerRepo,
		SquaddieRepo: squaddieRepo,
	}

	controller := actioncontroller.WhiteRoomController{}
	viewer := actionviewer.ConsoleActionViewer{}
	processSquaddieActions(chapterReplay, &viewer, &controller, repos)

	viewer.PrintMessages(output)
	return nil
}

func processSquaddieActions(
	chapterReplay *replay.ChapterReplay,
	viewer *actionviewer.ConsoleActionViewer,
	controller *actioncontroller.WhiteRoomController,
	repositories *repositories.RepositoryCollection) {
	initializeAllSquaddies(chapterReplay, repositories)
	for _, action := range chapterReplay.Actions {
		continueProcessing := processSquaddieAction(
			action,
			viewer,
			controller,
			repositories,
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
				viewer.Messages = append(viewer.Messages, description)
			}
		}
		return false
	}

	forecast := controller.GenerateForecast(powerSetup, repositories)
	viewer.PrepareForecast(forecast, repositories)

	result := controller.GenerateResult(forecast, repositories, true, action.RandomSeed)
	viewer.PrepareResult(result, repositories, &actionviewer.ConsoleActionViewerVerbosity{
		ShowTargetStatus: true,
	})
	return true
}

func loadSquaddieRepo(squaddieYamlData []byte) (repo *squaddie.Repository) {
	squaddieRepo := squaddie.NewSquaddieRepository()
	_, err := squaddieRepo.AddYAMLSource(squaddieYamlData)
	if err != nil {
		panic(err.Error())
	}
	return squaddieRepo
}

func loadPowerRepo(powerYamlData []byte) (repo *powerrepository.Repository) {
	powerRepo := powerrepository.NewPowerRepository()
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

func createSquaddieRepo(input io.Reader) *squaddie.Repository {
	squaddieData, squaddieErr := ioutil.ReadAll(input)
	if squaddieErr != nil {
		log.Fatal(squaddieErr)
	}

	repo := loadSquaddieRepo(squaddieData)
	if repo == nil {
		log.Fatal("Squaddie Repo is nil")
	}
	return repo
}

func createPowerRepo(input io.Reader) *powerrepository.Repository {
	powerData, powerErr := ioutil.ReadAll(input)
	if powerErr != nil {
		log.Fatal(powerErr)
	}

	repo := loadPowerRepo(powerData)
	if repo == nil {
		log.Fatal("Power Repo is nil")
	}
	return repo
}

func createChapterReplay(input io.Reader) *replay.ChapterReplay {
	scriptData, scriptErr := ioutil.ReadAll(input)
	if scriptErr != nil {
		log.Fatal(scriptErr)
	}

	chapterReplay, replayErr := replay.NewCreateMapReplayFromYAML(scriptData)
	if replayErr != nil {
		log.Fatal(replayErr)
	}
	return chapterReplay
}
