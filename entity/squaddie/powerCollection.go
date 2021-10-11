package squaddie

import (
	"fmt"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/utility"
)

// PowerCollection tracks what powers the squaddie has as well as what is in use.
type PowerCollection struct {
	PowerReferences          []*power.Reference `json:"powers" yaml:"powers"`
	CurrentlyEquippedPowerID string             `json:"equipped_power_id" yaml:"equipped_power_id"`
}

// AddInnatePower gives the Squaddie access to the power.
//  Raises an error if the squaddie already has the power.
func (powerCollection *PowerCollection) AddInnatePower(newPower *power.Power) error {
	if ContainsPowerID(powerCollection.PowerReferences, newPower.ID) {
		newError := fmt.Errorf(`squaddie already has innate power with ID "%s"`, newPower.ID)
		utility.Log(newError.Error(), 0, utility.Error)
		return newError
	}

	powerCollection.PowerReferences = append(powerCollection.PowerReferences, &power.Reference{Name: newPower.Name, ID: newPower.ID})
	return nil
}

// GetInnatePowerIDNames returns a list of all the powers the squaddie has access to.
func (powerCollection *PowerCollection) GetInnatePowerIDNames() []*power.Reference {
	powerIDNames := []*power.Reference{}
	for _, reference := range powerCollection.PowerReferences {
		powerIDNames = append(powerIDNames, &power.Reference{Name: reference.Name, ID: reference.ID})
	}
	return powerIDNames
}

// ClearInnatePowers removes all of the squaddie's powers.
func (powerCollection *PowerCollection) ClearInnatePowers() {
	powerCollection.PowerReferences = []*power.Reference{}
	powerCollection.CurrentlyEquippedPowerID = ""
}

// ClearTemporaryPowerReferences empties the temporary references to powers.
func (powerCollection *PowerCollection) ClearTemporaryPowerReferences() {
	powerCollection.PowerReferences = []*power.Reference{}
}

// HasPowerWithID returns a bool indicating if the squaddie has this power.
func (powerCollection *PowerCollection) HasPowerWithID(powerID string) bool {
	for _, powerReference := range powerCollection.GetInnatePowerIDNames() {
		if powerReference.ID == powerID {
			return true
		}
	}

	return false
}

// ContainsPowerID returns true if the squaddie has a reference to a power with the given ID.
func ContainsPowerID(references []*power.Reference, powerID string) bool {
	for _, reference := range references {
		if reference.ID == powerID {
			return true
		}
	}
	return false
}

// FilterPowerID returns a list of power references that satisfy the condition.
func FilterPowerID(references []*power.Reference, condition func(*power.Reference) bool) []*power.Reference {
	selectedReferences := []*power.Reference{}
	for _, reference := range references {
		if condition(reference) == true {
			selectedReferences = append(selectedReferences, reference)
		}
	}
	return selectedReferences
}

// HasEquippedPower returns true if the squaddie has already equipped a power.
func (powerCollection *PowerCollection) HasEquippedPower() bool {
	return powerCollection.CurrentlyEquippedPowerID != ""
}

// EquipPower sets the currently equipped power ID to the one given.
func (powerCollection *PowerCollection) EquipPower(powerID string) {
	powerCollection.CurrentlyEquippedPowerID = powerID
}

// GetEquippedPowerID returns the ID of the Power this squaddie has equipped.
//  returns an empty string if nothing is equipped.
func (powerCollection *PowerCollection) GetEquippedPowerID() string {
	return powerCollection.CurrentlyEquippedPowerID
}
