package actioncontroller

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagescenario"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerattackforecast"
	"github.com/cserrant/terosBattleServer/usecase/powercommit"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
	"github.com/cserrant/terosBattleServer/utility"
)

// WhiteRoomController assumes all Squaddies are within range and can attack each other.
type WhiteRoomController struct {}

// SetupAction creates a record of the next action.
func (controller *WhiteRoomController) SetupAction(attacker, target *squaddie.Squaddie, power *power.Power) *powerusagescenario.Setup {
	powerSetup := &powerusagescenario.Setup{
		UserID:          attacker.Identification.ID,
		PowerID:         power.ID,
		Targets:         []string{target.Identification.ID},
		IsCounterAttack: false,
	}
	return powerSetup
}

// GenerateForecast uses the action to predict results.
func (controller *WhiteRoomController) GenerateForecast(action *powerusagescenario.Setup, repos *repositories.RepositoryCollection) *powerattackforecast.Forecast {
	powerForecast := &powerattackforecast.Forecast{
		Setup: *action,
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo:    repos.SquaddieRepo,
			PowerRepo:       repos.PowerRepo,
		},
	}
	powerForecast.CalculateForecast()
	return powerForecast
}

// GenerateResult uses the forecast to create results.
func (controller *WhiteRoomController) GenerateResult(forecast *powerattackforecast.Forecast, repos *repositories.RepositoryCollection) *powercommit.Result {
	powerResult := &powercommit.Result{
		Forecast: forecast,
		DieRoller: &utility.RandomDieRoller{},
	}
	powerResult.Commit()
	return powerResult
}
