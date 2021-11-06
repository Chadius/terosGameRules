package squaddie

import (
	"github.com/chadius/terosbattleserver/entity/power"
)

// PowerCollection tracks what powers the squaddie has as well as what is in use.
type PowerCollection struct {
	PowerReferences          []*power.Reference `json:"powers" yaml:"powers"`
	CurrentlyEquippedPowerID string             `json:"equipped_power_id" yaml:"equipped_power_id"`
}

// GetCopyOfPowerReferences returns a list of all the powers the squaddie has access to.
func (powerCollection *PowerCollection) GetCopyOfPowerReferences() []*power.Reference {
	powerIDNames := []*power.Reference{}
	for _, reference := range powerCollection.PowerReferences {
		powerIDNames = append(powerIDNames, &power.Reference{Name: reference.Name, PowerID: reference.PowerID})
	}
	return powerIDNames
}

// ClearPowerReferences removes all squaddie's powers.
func (powerCollection *PowerCollection) ClearPowerReferences() {
	powerCollection.PowerReferences = []*power.Reference{}
	powerCollection.CurrentlyEquippedPowerID = ""
}

// HasPowerWithID returns a bool indicating if the squaddie has this power.
func (powerCollection *PowerCollection) HasPowerWithID(powerID string) bool {
	for _, powerReference := range powerCollection.GetCopyOfPowerReferences() {
		if powerReference.PowerID == powerID {
			return true
		}
	}

	return false
}

// HasEquippedPower returns true if the squaddie has already equipped a power.
func (powerCollection *PowerCollection) HasEquippedPower() bool {
	return powerCollection.CurrentlyEquippedPowerID != ""
}

// EquipPower sets the currently equipped power SquaddieID to the one given.
func (powerCollection *PowerCollection) EquipPower(powerID string) {
	powerCollection.CurrentlyEquippedPowerID = powerID
}

// GetEquippedPowerID returns the SquaddieID of the Power this squaddie has equipped.
//  returns an empty string if nothing is equipped.
func (powerCollection *PowerCollection) GetEquippedPowerID() string {
	return powerCollection.CurrentlyEquippedPowerID
}

// AddPowerReference adds a power reference, assuming it doesn't exist
func (powerCollection *PowerCollection) AddPowerReference (reference *power.Reference) {
	for _, ref := range powerCollection.PowerReferences {
		if ref.PowerID == reference.PowerID {
			return
		}
	}

	if powerCollection.PowerReferences == nil {
		powerCollection.PowerReferences = []*power.Reference{}
	}

	powerCollection.PowerReferences = append(powerCollection.PowerReferences, reference)
}

// RemovePowerReferenceByPowerID removes a power reference.
func (powerCollection *PowerCollection) RemovePowerReferenceByPowerID (powerID string) {
	foundPowerToRemove := false
	indexToRemove := 0
	for index, ref := range powerCollection.PowerReferences {
		if ref.PowerID == powerID {
			foundPowerToRemove = true
			indexToRemove = index
		}
	}

	if foundPowerToRemove {
		powerCollection.PowerReferences = append(
			powerCollection.PowerReferences[:indexToRemove],
			powerCollection.PowerReferences[indexToRemove+1:]...
		)
	}
}
