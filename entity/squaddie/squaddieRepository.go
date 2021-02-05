package squaddie

import (
	"encoding/json"
	"github.com/cserrant/terosBattleServer/entity/power"
	"gopkg.in/yaml.v2"
)

// Repository will interact with external devices to manage Squaddies.
type Repository struct {
	squaddiesByName map[string]Squaddie
}

// NewSquaddieRepository generates a pointer to a new Squaddie.
func NewSquaddieRepository() *Repository {
	repository := Repository{
		map[string]Squaddie{},
	}
	return &repository
}

type unmarshalFunc func([]byte, interface{}) error

// AddJSONSource consumes a given bytestream and tries to analyze it.
func (repository *Repository) AddJSONSource(data []byte) (bool, error) {
	return repository.addSource(data, json.Unmarshal)
}

// AddYAMLSource consumes a given bytestream and tries to analyze it.
func (repository *Repository) AddYAMLSource(data []byte) (bool, error) {
	return repository.addSource(data, yaml.Unmarshal)
}

// AddSource consumes a given bytestream of the given sourceType and tries to analyze it.
func (repository *Repository) addSource(data []byte, unmarshal unmarshalFunc) (bool, error) {
	var unmarshalError error
	var listOfSquaddies []Squaddie
	unmarshalError = unmarshal(data, &listOfSquaddies)

	if unmarshalError != nil {
		return false, unmarshalError
	}
	for _, squaddieToAdd := range listOfSquaddies {
		squaddieErr := CheckSquaddieForErrors(&squaddieToAdd)
		if squaddieErr != nil {
			return false, squaddieErr
		}
		squaddieToAdd.SetHPToMax()
		repository.squaddiesByName[squaddieToAdd.Name] = squaddieToAdd
	}

	return true, nil
}

// GetNumberOfSquaddies returns the number of Squaddies ready to retrieve.
func (repository *Repository) GetNumberOfSquaddies() int {
	return len(repository.squaddiesByName)
}

// GetByName retrieves a Squaddie by name
func (repository *Repository) GetByName(squaddieName string) *Squaddie {
	squaddie, squaddieExists := repository.squaddiesByName[squaddieName]
	if !squaddieExists {
		return nil
	}
	return &squaddie
}

// MarshalSquaddieIntoJSON converts the given Squaddie into JSON.
func (repository *Repository) MarshalSquaddieIntoJSON(squaddie *Squaddie) ([]byte, error) {
	type Alias Squaddie

	return json.Marshal(&struct {
		*Alias
		PowerIDNames []*power.Reference `json:"powers" yaml:"powers"`
	}{
		Alias:        (*Alias)(squaddie),
		PowerIDNames: squaddie.GetInnatePowerIDNames(),
	})
}
