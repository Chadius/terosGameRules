package repository

import (
	"encoding/json"
	"github.com/cserrant/terosBattleServer/entity"
	"gopkg.in/yaml.v2"
)

// PowerRepository will interact with external devices to manage Powers.
type PowerRepository struct {
	powersByName map[string]entity.Power
	powersByID map[string]*entity.Power
}

// NewPowerRepository generates a pointer to a new PowerRepository.
func NewPowerRepository() *PowerRepository {
	repository := PowerRepository{
		map[string]entity.Power{},
		map[string]*entity.Power{},
	}
	return &repository
}

// AddJSONSource consumes a given bytestream and tries to analyze it.
func (repository *PowerRepository) AddJSONSource(data []byte) (bool, error) {
	return repository.addSource(data, json.Unmarshal)
}

// AddYAMLSource consumes a given bytestream and tries to analyze it.
func (repository *PowerRepository) AddYAMLSource(data []byte) (bool, error) {
	return repository.addSource(data, yaml.Unmarshal)
}

// AddSlicePowerSource tries to add the slice of powers to the repo.
func (repository *PowerRepository) AddSlicePowerSource(powersToAdd []*entity.Power) (bool, error) {
	for _, PowerToAdd := range powersToAdd {
		PowerErr := entity.CheckPowerForErrors(PowerToAdd)
		if PowerErr != nil {
			return false, PowerErr
		}
		repository.powersByName[PowerToAdd.Name] = *PowerToAdd
		repository.powersByID[PowerToAdd.ID] = PowerToAdd
	}
	return true, nil
}

// AddSource consumes a given bytestream of the given sourceType and tries to analyze it.
func (repository *PowerRepository) addSource(data []byte, unmarshal unmarshalFunc) (bool, error) {
	var unmarshalError error
	var listOfPowers []entity.Power

	unmarshalError = unmarshal(data, &listOfPowers)

	if unmarshalError != nil {
		return false, unmarshalError
	}
	for _, PowerToAdd := range listOfPowers {
		PowerErr := entity.CheckPowerForErrors(&PowerToAdd)
		if PowerErr != nil {
			return false, PowerErr
		}
		repository.powersByName[PowerToAdd.Name] = PowerToAdd
		repository.powersByID[PowerToAdd.ID] = &PowerToAdd
	}
	return true, nil
}

// GetNumberOfPowers returns the number of Powers ready to retrieve.
func (repository *PowerRepository) GetNumberOfPowers() int {
	return len(repository.powersByID)
}

// GetPowerByID returns the Power stored by ID.
func (repository *PowerRepository) GetPowerByID(powerID string) *entity.Power {
	return repository.powersByID[powerID]
}

// GetAllPowersByName returns a slice of powers based on a given name
func (repository *PowerRepository) GetAllPowersByName(powerNameToFind string) []*entity.Power {
	powersFound := []*entity.Power{}

	for _, power := range repository.powersByID {
		if power.Name == powerNameToFind {
			powersFound = append(powersFound, power)
		}
	}

	return powersFound
}
