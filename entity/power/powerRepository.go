package power

import (
	"encoding/json"
	"github.com/cserrant/terosBattleServer/utility"
	"gopkg.in/yaml.v2"
)

// Repository will interact with external devices to manage Powers.
type Repository struct {
	powersByName map[string]Power
	powersByID map[string]*Power
}

// NewPowerRepository generates a pointer to a new PowerRepository.
func NewPowerRepository() *Repository {
	repository := Repository{
		map[string]Power{},
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
	for _, PowerToAdd := range powersToAdd {
		PowerErr := CheckPowerForErrors(PowerToAdd)
		if PowerErr != nil {
			return false, PowerErr
		}
		repository.powersByName[PowerToAdd.Name] = *PowerToAdd
		repository.powersByID[PowerToAdd.ID] = PowerToAdd
	}
	return true, nil
}

// AddSource consumes a given bytestream of the given sourceType and tries to analyze it.
func (repository *Repository) addSource(data []byte, unmarshal utility.UnmarshalFunc) (bool, error) {
	var unmarshalError error
	var listOfPowers []Power

	unmarshalError = unmarshal(data, &listOfPowers)

	if unmarshalError != nil {
		return false, unmarshalError
	}
	for _, PowerToAdd := range listOfPowers {
		PowerErr := CheckPowerForErrors(&PowerToAdd)
		if PowerErr != nil {
			return false, PowerErr
		}
		repository.powersByName[PowerToAdd.Name] = PowerToAdd
		repository.powersByID[PowerToAdd.ID] = &PowerToAdd
	}
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
