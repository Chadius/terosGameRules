package terosgamerules

import (
	"errors"
	"fmt"
	"github.com/chadius/terosgamerules/entity/actioncontroller"
	"github.com/chadius/terosgamerules/entity/actionviewer"
	"github.com/chadius/terosgamerules/entity/powerrepository"
	"github.com/chadius/terosgamerules/entity/replay"
	"github.com/chadius/terosgamerules/entity/squaddie"
	"github.com/chadius/terosgamerules/usecase/powerequip"
	"github.com/chadius/terosgamerules/usecase/repositories"
	"github.com/chadius/terosgamerules/utility"
	"io"
	"io/ioutil"
	"reflect"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . RulesStrategy
// RulesStrategy shapes the expected messages and the expected responses when running the rules.
type RulesStrategy interface {
	ReplayBattleScript(scriptFileHandle, squaddieFileHandle, powerFileHandle io.Reader, output io.Writer) error
}

type GameRules struct{}

// ReplayBattleScript uses the input streams to read and replay several rounds of combat,
//  writing the results to a supplied output stream.
func (g *GameRules) ReplayBattleScript(scriptFileHandle, squaddieFileHandle, powerFileHandle io.Reader, output io.Writer) error {
	utility.Logger = &utility.FileLogger{}

	squaddieRepo, squaddieErr := g.createSquaddieRepo(squaddieFileHandle)
	if squaddieErr != nil {
		return squaddieErr
	}

	powerRepo, powerErr := g.createPowerRepo(powerFileHandle)
	if powerErr != nil {
		return powerErr
	}

	chapterReplay, scriptErr := g.createChapterReplay(scriptFileHandle)
	if scriptErr != nil {
		return scriptErr
	}

	repos := &repositories.RepositoryCollection{
		PowerRepo:    powerRepo,
		SquaddieRepo: squaddieRepo,
	}

	controller := actioncontroller.WhiteRoomController{}
	viewer := actionviewer.ConsoleActionViewer{}
	g.processSquaddieActions(chapterReplay, &viewer, &controller, repos)

	viewer.PrintMessages(output)
	return nil
}

func (g *GameRules) processSquaddieActions(
	chapterReplay *replay.ChapterReplay,
	viewer *actionviewer.ConsoleActionViewer,
	controller *actioncontroller.WhiteRoomController,
	repositories *repositories.RepositoryCollection) {
	g.initializeAllSquaddies(chapterReplay, repositories)
	for _, action := range chapterReplay.Actions {
		continueProcessing := g.processSquaddieAction(
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

func (g *GameRules) processSquaddieAction(
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

func (g *GameRules) loadSquaddieRepo(squaddieYamlData []byte) (repo *squaddie.Repository) {
	squaddieRepo := squaddie.NewSquaddieRepository()
	err := squaddieRepo.AddSquaddiesUsingYAML(squaddieYamlData)
	if err != nil {
		return nil
	}
	return squaddieRepo
}

func (g *GameRules) loadPowerRepo(powerYamlData []byte) (*powerrepository.Repository, error) {
	powerRepo := powerrepository.NewPowerRepository()
	_, err := powerRepo.AddYAMLSource(powerYamlData)
	if err != nil {
		return nil, err
	}
	return powerRepo, nil
}

func (g *GameRules) initializeAllSquaddies(replay *replay.ChapterReplay, repositories *repositories.RepositoryCollection) {
	squaddiesFound := map[string]bool{}

	for _, action := range replay.Actions {
		if squaddiesFound[action.UserID] != true {
			g.loadAndInitializeSquaddie(action.UserID, repositories)
			squaddiesFound[action.UserID] = true
		}
		for _, targetID := range action.TargetIDs {
			if squaddiesFound[targetID] != true {
				g.loadAndInitializeSquaddie(targetID, repositories)
				squaddiesFound[targetID] = true
			}
		}
	}
}

func (g *GameRules) loadAndInitializeSquaddie(squaddieID string, repositories *repositories.RepositoryCollection) {
	squaddieRepo := repositories.SquaddieRepo
	squaddie := squaddieRepo.GetOriginalSquaddieByID(squaddieID)
	if squaddie == nil {
		fmt.Sprintf("Squaddie %s does not exist, exiting", squaddieID)
	}
	squaddie.SetBarrierToMax()

	equipCheck := powerequip.CheckRepositories{}
	equipCheck.LoadAllOfSquaddieInnatePowers(squaddie, squaddie.GetCopyOfPowerReferences(), repositories)
	equipCheck.EquipDefaultPower(squaddie, repositories)
}

func (g *GameRules) createSquaddieRepo(input io.Reader) (*squaddie.Repository, error) {
	if input == nil || reflect.ValueOf(input).IsNil() {
		return nil, errors.New("no squaddie data found")
	}

	squaddieData, squaddieErr := io.ReadAll(input)
	if squaddieErr != nil {
		return nil, squaddieErr
	}

	repo := g.loadSquaddieRepo(squaddieData)
	if repo == nil {
		return nil, errors.New("squaddie data is invalid")
	}
	return repo, nil
}

func (g *GameRules) createPowerRepo(input io.Reader) (*powerrepository.Repository, error) {
	if input == nil || reflect.ValueOf(input).IsNil() {
		return nil, errors.New("no power data found")
	}

	powerData, powerErr := ioutil.ReadAll(input)
	if powerErr != nil {
		return nil, powerErr
	}

	repo, loadErr := g.loadPowerRepo(powerData)
	if loadErr != nil {
		return nil, errors.New("power data is invalid")
	}
	return repo, nil
}

func (g *GameRules) createChapterReplay(input io.Reader) (*replay.ChapterReplay, error) {
	if input == nil || reflect.ValueOf(input).IsNil() {
		return nil, errors.New("no script data found")
	}

	scriptData, scriptErr := ioutil.ReadAll(input)
	if scriptErr != nil {
		return nil, scriptErr
	}

	chapterReplay, replayErr := replay.NewCreateMapReplayFromYAML(scriptData)
	if replayErr != nil {
		return nil, errors.New("script data is invalid")
	}

	return chapterReplay, nil
}
