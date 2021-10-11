package power

import (
	"encoding/json"
	"github.com/chadius/terosbattleserver/utility"
	"gopkg.in/yaml.v2"
)

// Repository will interact with external devices to manage Powers.
type Repository struct {
	powersByID map[string]*Power
}

// NewPowerRepository generates a pointer to a new Repository.
func NewPowerRepository() *Repository {
	repository := Repository{
		map[string]*Power{},
	}
	return &repository
}

// AddJSONSource consumes a given bytestream and tries to analyze it.
func (repository *Repository) AddJSONSource(data []byte) (bool, error) {
	return repository.addSource(data, json.Unmarshal)
}

// AddYAMLSource consumes a given bytestream and tries to analyze it.
func (repository *Repository) AddYAMLSource(data []byte) (bool, error) {
	return repository.addSource(data, yaml.Unmarshal)
}

// AddSlicePowerSource tries to add the slice of powers to the repo.
func (repository *Repository) AddSlicePowerSource(powersToAdd []*Power) (bool, error) {
	for _, powerToAdd := range powersToAdd {
		success, err := repository.tryToAddPower(powerToAdd)
		if success == false {
			return false, err
		}
	}
	return true, nil
}

// addSource consumes a given bytestream of the given sourceType and tries to analyze it.
func (repository *Repository) addSource(data []byte, unmarshal utility.UnmarshalFunc) (bool, error) {
	var unmarshalError error
	var listOfPowers []Power

	unmarshalError = unmarshal(data, &listOfPowers)

	if unmarshalError != nil {
		return false, unmarshalError
	}
	for index := range listOfPowers {
		powerToAdd := listOfPowers[index]
		success, err := repository.tryToAddPower(&powerToAdd)
		if success == false {
			return false, err
		}
	}
	return true, nil
}

func (repository *Repository) tryToAddPower(powerToAdd *Power) (bool, error) {
	PowerErr := CheckPowerForErrors(powerToAdd)
	if PowerErr != nil {
		return false, PowerErr
	}
	repository.powersByID[powerToAdd.ID] = powerToAdd
	return true, nil
}

// GetNumberOfPowers returns the number of Powers ready to retrieve.
func (repository *Repository) GetNumberOfPowers() int {
	return len(repository.powersByID)
}

// GetPowerByID returns the Power stored by ID.
func (repository *Repository) GetPowerByID(powerID string) *Power {
	return repository.powersByID[powerID]
}

// GetAllPowersByName returns a slice of powers based on a given name
func (repository *Repository) GetAllPowersByName(powerNameToFind string) []*Power {
	powersFound := []*Power{}
	for _, power := range repository.powersByID {
		if power.Name == powerNameToFind {
			powersFound = append(powersFound, power)
		}
	}

	return powersFound
}
