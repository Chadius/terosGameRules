package squaddie

import (
	"github.com/chadius/terosbattleserver/entity/power"
)

// PowerCollection tracks what powers the squaddie has as well as what is in use.
type PowerCollection struct {
	powerReferences          []*power.Reference
	currentlyEquippedPowerID string
}

// GetCopyOfPowerReferences returns a list of all the powers the squaddie has access to.
func (powerCollection *PowerCollection) GetCopyOfPowerReferences() []*power.Reference {
	powerIDNames := []*power.Reference{}
	for _, reference := range powerCollection.powerReferences {
		powerIDNames = append(powerIDNames, &power.Reference{Name: reference.Name, PowerID: reference.PowerID})
	}
	return powerIDNames
}

// ClearPowerReferences removes all squaddie's powers.
func (powerCollection *PowerCollection) ClearPowerReferences() {
	powerCollection.powerReferences = []*power.Reference{}
	powerCollection.currentlyEquippedPowerID = ""
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
	return powerCollection.currentlyEquippedPowerID != ""
}

// EquipPower sets the currently equipped power powerID to the one given.
func (powerCollection *PowerCollection) EquipPower(powerID string) {
	powerCollection.currentlyEquippedPowerID = powerID
}

// GetEquippedPowerID returns the powerID of the Power this squaddie has equipped.
//  returns an empty string if nothing is equipped.
func (powerCollection *PowerCollection) GetEquippedPowerID() string {
	return powerCollection.currentlyEquippedPowerID
}

// AddPowerReference adds a power reference, assuming it doesn't exist
func (powerCollection *PowerCollection) AddPowerReference(reference *power.Reference) {
	for _, ref := range powerCollection.powerReferences {
		if ref.PowerID == reference.PowerID {
			return
		}
	}

	if powerCollection.powerReferences == nil {
		powerCollection.powerReferences = []*power.Reference{}
	}

	powerCollection.powerReferences = append(powerCollection.powerReferences, reference)
}

// RemovePowerReferenceByPowerID removes a power reference.
func (powerCollection *PowerCollection) RemovePowerReferenceByPowerID(powerID string) {
	foundPowerToRemove := false
	indexToRemove := 0
	for index, ref := range powerCollection.powerReferences {
		if ref.PowerID == powerID {
			foundPowerToRemove = true
			indexToRemove = index
		}
	}

	if foundPowerToRemove {
		powerCollection.powerReferences = append(
			powerCollection.powerReferences[:indexToRemove],
			powerCollection.powerReferences[indexToRemove+1:]...,
		)
	}
}

// Len returns the number of powers in the collection
func (powerCollection *PowerCollection) Len() int {
	return len(powerCollection.powerReferences)
}

// HasSamePowersAs sees if the other collection has the same power references
func (powerCollection *PowerCollection) HasSamePowersAs(other *PowerCollection) bool {
	if powerCollection.Len() != other.Len() {
		return false
	}

	powersByID := map[string]bool{}
	for _, reference := range powerCollection.powerReferences {
		powersByID[reference.PowerID] = false
	}

	for _, reference := range other.GetCopyOfPowerReferences() {
		alreadyFound, exists := powersByID[reference.PowerID]
		if !exists {
			return false
		}
		if alreadyFound {
			return false
		}
		powersByID[reference.PowerID] = true
	}

	for _, wasFound := range powersByID {
		if wasFound == false {
			return false
		}
	}

	return true
}
