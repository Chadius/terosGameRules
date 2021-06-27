package powerequip

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
	"github.com/cserrant/terosBattleServer/usecase/squaddiestats"
)

// EquipDefaultPower will automatically equip the first power the squaddie has.
//  Returns the power and a boolean.
func EquipDefaultPower(squaddie *squaddie.Squaddie, repos *repositories.RepositoryCollection) (*power.Power, bool) {
	for _, powerReference := range squaddie.PowerCollection.PowerReferences {
		powerToCheck := repos.PowerRepo.GetPowerByID(powerReference.ID)
		if powerToCheck.AttackEffect != nil && powerToCheck.AttackEffect.CanBeEquipped == true {
			equippingPowerWasSuccessful := SquaddieEquipPower(squaddie, powerToCheck.ID, repos)
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
	if powerToEquip.AttackEffect == nil || powerToEquip.AttackEffect.CanBeEquipped == false {
		return false
	}

	squaddie.PowerCollection.EquipPower(powerToEquipID)
	return true
}

// LoadAllOfSquaddieInnatePowers loads the powers from the repo the squaddie needs and gives it to them.
//  Raises an error if the PowerRepository does not have one of the squaddie's powers.
func LoadAllOfSquaddieInnatePowers(squaddie *squaddie.Squaddie, powerReferencesToLoad []*power.Reference, repos *repositories.RepositoryCollection) (int, error) {
	numberOfPowersAdded := 0

	squaddie.PowerCollection.ClearInnatePowers()
	squaddie.PowerCollection.ClearTemporaryPowerReferences()

	for _, powerIDName := range powerReferencesToLoad {
		powerToAdd := repos.PowerRepo.GetPowerByID(powerIDName.ID)
		if powerToAdd == nil {
			return numberOfPowersAdded, fmt.Errorf("squaddie '%s' tried to add Power '%s' but it does not exist", squaddie.Identification.Name, powerIDName.Name)
		}

		err := squaddie.PowerCollection.AddInnatePower(powerToAdd)
		if err == nil {
			numberOfPowersAdded = numberOfPowersAdded + 1
		}
	}

	return numberOfPowersAdded, nil
}

// CanSquaddieCounterWithEquippedWeapon returns true if the squaddie can use the currently equipped
//   weapon for counter attacks.
func CanSquaddieCounterWithEquippedWeapon(squaddieID string, repos *repositories.RepositoryCollection) (bool, error) {
	squaddie := repos.SquaddieRepo.GetOriginalSquaddieByID(squaddieID)
	equippedPowerID := squaddie.PowerCollection.GetEquippedPowerID()
	if equippedPowerID == "" {
		return false, fmt.Errorf("squaddie has no equipped power, %s", squaddieID)
	}

	canCounter, counterErr := squaddiestats.GetSquaddieCanCounterAttackWithPower(squaddieID, equippedPowerID, repos)
	return canCounter, counterErr
}
