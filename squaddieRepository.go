package terosbattleserver

import "encoding/json"

// SquaddieRepository comment
type SquaddieRepository interface {
	getByName(string) *Squaddie
	getNumberOfSquaddies() int
}

// SquaddieRetriever comment
type SquaddieRetriever struct {
	squaddiesByName map[string]Squaddie
}

// NewSquaddieRetriever generates a pointer to a new Squaddie.
func NewSquaddieRetriever() *SquaddieRetriever {
	repository := SquaddieRetriever{
		map[string]Squaddie{},
	}
	return &repository
}

// AddJSONSource consumes a given bytestream and tries to analyze it.
func (repository *SquaddieRetriever) AddJSONSource(data []byte) (bool, error) {
	var listOfSquaddies []Squaddie
	err := json.Unmarshal(data, &listOfSquaddies)
	if err != nil {
		return false, err
	}
	for _, squaddieToAdd := range listOfSquaddies {
		repository.squaddiesByName[squaddieToAdd.Name] = squaddieToAdd
	}
	return true, nil
}

// GetNumberOfSquaddies returns the number of Squaddies ready to retrieve.
func (repository *SquaddieRetriever) GetNumberOfSquaddies() int {
	return len(repository.squaddiesByName)
}

// GetByName retrieves a Squaddie by name
func (repository *SquaddieRetriever) GetByName(squaddieName string) *Squaddie {
	squaddie, squaddieExists := repository.squaddiesByName[squaddieName]
	if !squaddieExists {
		return nil
	}
	return &squaddie
}
