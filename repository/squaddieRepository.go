package repository

import (
	"encoding/json"
	"github.com/cserrant/terosBattleServer/entity"
	"gopkg.in/yaml.v2"
)

// SquaddieRepository will interact with external devices to manage Squaddies.
type SquaddieRepository struct {
	squaddiesByName map[string]entity.Squaddie
}

// NewSquaddieRepository generates a pointer to a new Squaddie.
func NewSquaddieRepository() *SquaddieRepository {
	repository := SquaddieRepository{
		map[string]entity.Squaddie{},
	}
	return &repository
}

type unmarshalFunc func([]byte, interface{}) error

// AddJSONSource consumes a given bytestream and tries to analyze it.
func (repository *SquaddieRepository) AddJSONSource(data []byte) (bool, error) {
	return repository.addSource(data, json.Unmarshal)
}

// AddYAMLSource consumes a given bytestream and tries to analyze it.
func (repository *SquaddieRepository) AddYAMLSource(data []byte) (bool, error) {
	return repository.addSource(data, yaml.Unmarshal)
}

// AddSource consumes a given bytestream of the given sourceType and tries to analyze it.
func (repository *SquaddieRepository) addSource(data []byte, unmarshal unmarshalFunc) (bool, error) {
	var unmarshalError error
	var listOfSquaddies []entity.Squaddie
	unmarshalError = unmarshal(data, &listOfSquaddies)

	if unmarshalError != nil {
		return false, unmarshalError
	}
	for _, squaddieToAdd := range listOfSquaddies {
		squaddieErr := entity.CheckSquaddieForErrors(&squaddieToAdd)
		if squaddieErr != nil {
			return false, squaddieErr
		}
		squaddieToAdd.SetHPToMax()
		repository.squaddiesByName[squaddieToAdd.Name] = squaddieToAdd
	}

	return true, nil
}

// GetNumberOfSquaddies returns the number of Squaddies ready to retrieve.
func (repository *SquaddieRepository) GetNumberOfSquaddies() int {
	return len(repository.squaddiesByName)
}

// GetByName retrieves a Squaddie by name
func (repository *SquaddieRepository) GetByName(squaddieName string) *entity.Squaddie {
	squaddie, squaddieExists := repository.squaddiesByName[squaddieName]
	if !squaddieExists {
		return nil
	}
	return &squaddie
}

// MarshalSquaddieIntoJSON converts the given Squaddie into JSON.
func (repository *SquaddieRepository) MarshalSquaddieIntoJSON(squaddie *entity.Squaddie) ([]byte, error) {
	type Alias entity.Squaddie

	return json.Marshal(&struct {
		*Alias
		PowerIDNames []*entity.PowerReference `json:"powers" yaml:"powers"`
	}{
		Alias:        (*Alias)(squaddie),
		PowerIDNames: squaddie.GetInnatePowerIDNames(),
	})
}
