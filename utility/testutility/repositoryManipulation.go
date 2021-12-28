package testutility

import (
	"github.com/chadius/terosbattleserver/entity/powerinterface"
	"github.com/chadius/terosbattleserver/entity/powerreference"
	"github.com/chadius/terosbattleserver/entity/squaddieinterface"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
)

func AddSquaddieWithInnatePowersToRepos(squaddieToAdd squaddieinterface.Interface, powerToAdd powerinterface.Interface, repos *repositories.RepositoryCollection, equipPower bool) {
	repos.SquaddieRepo.AddSquaddies([]squaddieinterface.Interface{squaddieToAdd})
	if powerToAdd != nil {
		repos.PowerRepo.AddSlicePowerSource([]powerinterface.Interface{powerToAdd})
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
