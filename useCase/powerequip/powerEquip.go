package powerequip

import (
	"fmt"
	"github.com/chadius/terosbattleserver/entity/powerinterface"
	"github.com/chadius/terosbattleserver/entity/powerreference"
	"github.com/chadius/terosbattleserver/entity/squaddieinterface"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/utility"
)

// Strategy will equip squaddies with powers.
type Strategy interface {
	EquipDefaultPower(squaddie squaddieinterface.Interface, repos *repositories.RepositoryCollection) (powerinterface.Interface, bool)
	SquaddieEquipPower(squaddie squaddieinterface.Interface, powerToEquipID string, repos *repositories.RepositoryCollection) bool
	LoadAllOfSquaddieInnatePowers(squaddie squaddieinterface.Interface, powerReferencesToLoad []*powerreference.Reference, repos *repositories.RepositoryCollection) error
}

// CheckRepositories uses the repositories to make sure powers and squaddies exist.
type CheckRepositories struct{}

// EquipDefaultPower will automatically equip the first power the squaddie has.
//  Returns the power and a boolean value.
func (p *CheckRepositories) EquipDefaultPower(squaddie squaddieinterface.Interface, repos *repositories.RepositoryCollection) (powerinterface.Interface, bool) {
	for _, powerReference := range squaddie.GetCopyOfPowerReferences() {
		powerToCheck := repos.PowerRepo.GetPowerByID(powerReference.PowerID)
		if powerToCheck.CanAttack() && powerToCheck.CanBeEquipped() == true {
			equippingPowerWasSuccessful := p.SquaddieEquipPower(squaddie, powerToCheck.ID(), repos)
			return powerToCheck, equippingPowerWasSuccessful
		}
	}
	return nil, false
}

// SquaddieEquipPower will make the Squaddie equip a different power.
//   returns true upon success
func (p *CheckRepositories) SquaddieEquipPower(squaddie squaddieinterface.Interface, powerToEquipID string, repos *repositories.RepositoryCollection) bool {
	if squaddie.HasPowerWithID(powerToEquipID) == false {
		return false
	}

	powerToEquip := repos.PowerRepo.GetPowerByID(powerToEquipID)
	if powerToEquip == nil {
		return false
	}
	if !powerToEquip.CanAttack() || powerToEquip.CanBeEquipped() == false {
		return false
	}

	squaddie.EquipPower(powerToEquipID)
	return true
}

// LoadAllOfSquaddieInnatePowers loads the powers from the repo the squaddie needs and gives it to them.
//  Raises an error if the PowerRepository does not have one of the squaddie's powers.
func (p *CheckRepositories) LoadAllOfSquaddieInnatePowers(squaddie squaddieinterface.Interface, powerReferencesToLoad []*powerreference.Reference, repos *repositories.RepositoryCollection) error {
	squaddie.ClearPowerReferences()

	for _, powerIDName := range powerReferencesToLoad {
		powerToAdd := repos.PowerRepo.GetPowerByID(powerIDName.PowerID)
		if powerToAdd == nil {
			newError := fmt.Errorf("squaddie '%s' tried to add Power '%s' but it does not exist", squaddie.Name(), powerIDName.Name)
			utility.Log(newError.Error(), 0, utility.Error)
			return newError
		}

		squaddie.AddPowerReference(powerToAdd.GetReference())
	}

	return nil
}
