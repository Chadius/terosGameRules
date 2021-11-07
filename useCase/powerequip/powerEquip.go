package powerequip

import (
	"fmt"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/usecase/squaddiestats"
	"github.com/chadius/terosbattleserver/utility"
)

// EquipDefaultPower will automatically equip the first power the squaddie has.
//  Returns the power and a boolean.
func EquipDefaultPower(squaddie *squaddie.Squaddie, repos *repositories.RepositoryCollection) (*power.Power, bool) {
	for _, powerReference := range squaddie.GetCopyOfPowerReferences() {
		powerToCheck := repos.PowerRepo.GetPowerByID(powerReference.PowerID)
		if powerToCheck.CanAttack() && powerToCheck.CanBeEquipped() == true {
			equippingPowerWasSuccessful := SquaddieEquipPower(squaddie, powerToCheck.ID(), repos)
			return powerToCheck, equippingPowerWasSuccessful
		}
	}
	return nil, false
}

// SquaddieEquipPower will make the Squaddie equip a different power.
//   returns true upon success
func SquaddieEquipPower(squaddie *squaddie.Squaddie, powerToEquipID string, repos *repositories.RepositoryCollection) bool {
	if squaddie.PowerCollection.HasPowerWithID(powerToEquipID) == false {
		return false
	}

	powerToEquip := repos.PowerRepo.GetPowerByID(powerToEquipID)
	if powerToEquip == nil {
		return false
	}
	if !powerToEquip.CanAttack() || powerToEquip.CanBeEquipped() == false {
		return false
	}

	squaddie.PowerCollection.EquipPower(powerToEquipID)
	return true
}

// LoadAllOfSquaddieInnatePowers loads the powers from the repo the squaddie needs and gives it to them.
//  Raises an error if the PowerRepository does not have one of the squaddie's powers.
func LoadAllOfSquaddieInnatePowers(squaddie *squaddie.Squaddie, powerReferencesToLoad []*power.Reference, repos *repositories.RepositoryCollection) error {
	squaddie.PowerCollection.ClearPowerReferences()

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

// CanSquaddieCounterWithEquippedWeapon returns true if the squaddie can use the currently equipped
//   weapon for counter attacks.
func CanSquaddieCounterWithEquippedWeapon(squaddieID string, repos *repositories.RepositoryCollection) (bool, error) {
	squaddie := repos.SquaddieRepo.GetOriginalSquaddieByID(squaddieID)
	equippedPowerID := squaddie.PowerCollection.GetEquippedPowerID()
	if equippedPowerID == "" {
		newError := fmt.Errorf("squaddie has no equipped power, %s", squaddieID)
		utility.Log(newError.Error(), 0, utility.Error)
		return false, newError
	}

	canCounter, counterErr := squaddiestats.GetSquaddieCanCounterAttackWithPower(squaddieID, equippedPowerID, repos)
	return canCounter, counterErr
}
