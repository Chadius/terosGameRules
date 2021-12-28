package testutility

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerreference"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerattackforecast"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
)

func AddSquaddieWithInnatePowersToRepos(squaddieToAdd *squaddie.Squaddie, powerToAdd *power.Power, repos *repositories.RepositoryCollection, equipPower bool) {
	repos.SquaddieRepo.AddSquaddies([]*squaddie.Squaddie{squaddieToAdd})
	if powerToAdd != nil {
		repos.PowerRepo.AddSlicePowerSource([]*power.Power{powerToAdd})
		checkEquip := powerequip.CheckRepositories{}
		checkEquip.LoadAllOfSquaddieInnatePowers(
			squaddieToAdd,
			[]*powerreference.Reference{
				powerToAdd.GetReference(),
			},
			repos,
		)

		if equipPower {
			checkEquip.SquaddieEquipPower(squaddieToAdd, powerToAdd.ID(), repos)
		}
	}
}

// TODO Either Forecast should have these permanently or we should work on getting rid of these

// UpdateForecastWithNewTarget updates the forecast setup with a new target
func UpdateForecastWithNewTarget(squaddieToTarget *squaddie.Squaddie, squaddieRepo *squaddie.Repository, forecast *powerattackforecast.Forecast) {
	squaddieRepo.AddSquaddie(squaddieToTarget)
	forecast.UpdateForecastWithNewTarget(0, squaddieToTarget.ID())
}

// UpdateForecastWithNewUser updates the forecast setup with a new user
func UpdateForecastWithNewUser(squaddieToUse *squaddie.Squaddie, squaddieRepo *squaddie.Repository, forecast *powerattackforecast.Forecast) {
	squaddieRepo.AddSquaddie(squaddieToUse)
	forecast.UpdateForecastWithNewUser(squaddieToUse.ID())
}
