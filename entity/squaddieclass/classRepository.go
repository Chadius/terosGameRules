package squaddieclass

import (
	"encoding/json"
	"fmt"
	"github.com/chadius/terosbattleserver/utility"
)

// Repository will interact with external devices to manage Squaddie Classes.
type Repository struct {
	classesByID map[string]*Class
}

// NewRepository generates a pointer to a new Repository.
func NewRepository() *Repository {
	repository := Repository{
		map[string]*Class{},
	}
	return &repository
}

// AddJSONSource consumes a given bytestream and tries to analyze it.
func (repository *Repository) AddJSONSource(data []byte) (bool, error) {
	return repository.addSource(data, json.Unmarshal)
}

// AddListOfClasses adds multiple classes directly.
func (repository *Repository) AddListOfClasses(classes []*Class) (bool, error) {
	for _, classToAdd := range classes {
		repository.classesByID[classToAdd.ID()] = classToAdd
	}

	return true, nil
}

// AddYAMLSource consumes a given bytestream and tries to analyze it.
//func (repository *Repository) AddYAMLSource(data []byte) (bool, error) {
//	return repository.addSource(data, yaml.Unmarshal)
//}

// AddSource consumes a given bytestream of the given sourceType and tries to analyze it.
func (repository *Repository) addSource(data []byte, unmarshal utility.UnmarshalFunc) (bool, error) {
	var unmarshalError error
	var classes []Class
	unmarshalError = unmarshal(data, &classes)

	if unmarshalError != nil {
		return false, unmarshalError
	}
	for _, classToAdd := range classes {
		repository.classesByID[classToAdd.ID()] = &classToAdd
	}

	return true, nil
}

// GetNumberOfClasses returns the number of Classes ready to retrieve.
func (repository *Repository) GetNumberOfClasses() int {
	return len(repository.classesByID)
}

// GetClassByID returns a Class that matches the SquaddieID.
func (repository *Repository) GetClassByID(classID string) (*Class, error) {
	class, classFound := repository.classesByID[classID]
	if classFound == false {
		newError := fmt.Errorf(`class repository: No class found with SquaddieID: "%s"`, classID)
		utility.Log(newError.Error(), 0, utility.Error)
		return nil, newError
	}

	return class, nil
}
