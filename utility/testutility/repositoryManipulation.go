package testutility

import (
	"github.com/chadius/terosgamerules/entity/powerinterface"
	"github.com/chadius/terosgamerules/entity/powerreference"
	"github.com/chadius/terosgamerules/entity/squaddieinterface"
	"github.com/chadius/terosgamerules/usecase/powerequip"
	"github.com/chadius/terosgamerules/usecase/repositories"
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
