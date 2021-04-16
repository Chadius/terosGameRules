package powerequip

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
)

// GetEquippedPower returns the power the squaddie has equipped.
//   May return nil if the squaddie hasn't equipped anything.
func GetEquippedPower (squaddie *squaddie.Squaddie, repo *power.Repository) *power.Power {
	if squaddie.CurrentlyEquippedPowerID != "" {
		return repo.GetPowerByID(squaddie.CurrentlyEquippedPowerID)
	}

	for _, powerReference := range squaddie.PowerReferences {
		powerToCheck := repo.GetPowerByID(powerReference.ID)
		if powerToCheck.AttackEffect != nil && powerToCheck.AttackEffect.CanBeEquipped == true {
			return powerToCheck
		}
	}
	return nil
}

// SquaddieEquipPower will make the Squaddie equip a different power.
//   If the power is invalid, will return nil
func SquaddieEquipPower(squaddie *squaddie.Squaddie, powerToEquipID string, repo *power.Repository) bool {
	if squaddie.HasPowerWithID(powerToEquipID) == false {
		return false
	}

	powerToEquip := repo.GetPowerByID(powerToEquipID)
	if powerToEquip == nil {
		return false
	}
	if powerToEquip.AttackEffect.CanBeEquipped == false {
		return false
	}

	squaddie.CurrentlyEquippedPowerID = powerToEquipID
	return true
}

// LoadAllOfSquaddieInnatePowers loads the powers from the repo the squaddie needs and gives it to them.
//  Raises an error if the PowerRepository does not have one of the squaddie's powers.
func LoadAllOfSquaddieInnatePowers(squaddie *squaddie.Squaddie, powerReferencesToLoad []*power.Reference, repo *power.Repository) (int, error) {
	numberOfPowersAdded := 0

	squaddie.ClearInnatePowers()
	squaddie.ClearTemporaryPowerReferences()

	for _, powerIDName := range powerReferencesToLoad {
		powerToAdd := repo.GetPowerByID(powerIDName.ID)
		if powerToAdd == nil {
			return numberOfPowersAdded, fmt.Errorf("squaddie '%s' tried to add Power '%s' but it does not exist", squaddie.Name, powerIDName.Name)
		}

		err := squaddie.AddInnatePower(powerToAdd)
		if err == nil {
			numberOfPowersAdded = numberOfPowersAdded + 1
		}
	}

	return numberOfPowersAdded, nil
}