package repository

import (
	"encoding/json"

	"github.com/cserrant/terosBattleServer/entity"
	"gopkg.in/yaml.v2"
)

// PowerRepository will interact with external devices to manage Powers.
type PowerRepository struct {
	powersByName map[string]entity.Power
}

// NewPowerRepository generates a pointer to a new PowerRepository.
func NewPowerRepository() *PowerRepository {
	repository := PowerRepository{
		map[string]entity.Power{},
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
	}
	return true, nil
}

// GetNumberOfPowers returns the number of Powers ready to retrieve.
func (repository *PowerRepository) GetNumberOfPowers() int {
	return len(repository.powersByName)
}

// GetByName retrieves a Power by name
func (repository *PowerRepository) GetByName(PowerName string) *entity.Power {
	Power, PowerExists := repository.powersByName[PowerName]
	if !PowerExists {
		return nil
	}
	return &Power
}
