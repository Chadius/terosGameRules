package testutility

import (
	"github.com/chadius/terosbattleserver/entity/power"
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
			[]*power.Reference{
				powerToAdd.GetReference(),
			},
			repos,
		)

		if equipPower {
			checkEquip.SquaddieEquipPower(squaddieToAdd, powerToAdd.ID(), repos)
		}
	}
}

func UpdateForecastWithNewTarget(squaddieToTarget *squaddie.Squaddie, squaddieRepo *squaddie.Repository, forecast *powerattackforecast.Forecast) {
	squaddieRepo.AddSquaddie(squaddieToTarget)
	forecast.Setup.Targets[0] = squaddieToTarget.ID()
}

func UpdateForecastWithNewUser(squaddieToUse *squaddie.Squaddie, squaddieRepo *squaddie.Repository, forecast *powerattackforecast.Forecast) {
	squaddieRepo.AddSquaddie(squaddieToUse)
	forecast.Setup.UserID = squaddieToUse.ID()
}
