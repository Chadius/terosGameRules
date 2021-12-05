package squaddie

import (
	"encoding/json"
	"github.com/chadius/terosbattleserver/utility"
	"gopkg.in/yaml.v2"
)

// Repository will interact with external devices to manage Squaddies.
type Repository struct {
	squaddiesByID map[string]*Squaddie
}

// NewSquaddieRepository generates a pointer to a new Repository.
func NewSquaddieRepository() *Repository {
	repository := Repository{
		map[string]*Squaddie{},
	}
	return &repository
}

// AddSquaddies adds a slice of Squaddie to the repository.
func (repository *Repository) AddSquaddies(squaddies []*Squaddie) (bool, error) {
	for _, squaddieToAdd := range squaddies {
		_, err := repository.tryToAddSquaddie(squaddieToAdd)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// AddSquaddie adds a Squaddie to the repository.
func (repository *Repository) AddSquaddie(squaddieToAdd *Squaddie) (bool, error) {
	_, err := repository.tryToAddSquaddie(squaddieToAdd)
	if err != nil {
		return false, err
	}
	return true, nil
}

// AddSquaddiesUsingYAML adds multiple squaddies using builder objects and a YAML data stream.
func (repository *Repository) AddSquaddiesUsingYAML(data []byte) error {
	_, err := repository.unmarshalDataAndAddSquaddies(data, yaml.Unmarshal)
	return err
}

// AddSquaddiesUsingJSON adds multiple squaddies using builder objects and a JSON data stream.
func (repository *Repository) AddSquaddiesUsingJSON(data []byte) error {
	_, err := repository.unmarshalDataAndAddSquaddies(data, json.Unmarshal)
	return err
}

// unmarshalDataAndAddSquaddies reads the byte stream to create new squaddies.
func (repository *Repository) unmarshalDataAndAddSquaddies(data []byte, unmarshal utility.UnmarshalFunc) (bool, error) {
	var unmarshalError error

	var builderInstructions []BuilderOptionMarshal

	unmarshalError = unmarshal(data, &builderInstructions)

	if unmarshalError != nil {
		return false, unmarshalError
	}

	for _, instruction := range builderInstructions {
		newSquaddie := NewSquaddieFromMarshal(instruction).Build()

		newSquaddie.Defense.SetHPToMax()
		success, err := repository.tryToAddSquaddie(newSquaddie)
		if success == false {
			return false, err
		}
	}

	return true, nil
}

func (repository *Repository) tryToAddSquaddie(squaddieToAdd *Squaddie) (bool, error) {
	squaddieErr := CheckSquaddieForErrors(squaddieToAdd)
	if squaddieErr != nil {
		return false, squaddieErr
	}

	if squaddieToAdd.ID() == "" {
		squaddieToAdd.Identification.SetNewIDToRandom() // TODO ask the squaddie to set its ID
	}
	repository.squaddiesByID[squaddieToAdd.ID()] = squaddieToAdd
	return true, nil
}

// GetNumberOfSquaddies returns the number of Squaddies ready to retrieve.
func (repository *Repository) GetNumberOfSquaddies() int {
	return len(repository.squaddiesByID)
}

//CloneSquaddieWithNewID uses the base Squaddie to create a new one.
//  All fields will be the same except the squaddieID.
//  If newID isn't empty, the clone squaddieID is set to that.
//  Otherwise, it is randomly generated.
func (repository *Repository) CloneSquaddieWithNewID(base *Squaddie, newID string) (*Squaddie, error) {
	cloneBuilder := NewSquaddieBuilder().CloneOf(base)
	if newID != "" {
		cloneBuilder.WithID(newID)
	} else {
		cloneBuilder.WithID(utility.StringWithCharset(8, "abcdefgh0123456789"))
	}
	clone := cloneBuilder.Build()

	clone.ReduceHitPoints(clone.MaxHitPoints() - base.CurrentHitPoints())
	clone.ReduceBarrier(clone.MaxBarrier() - base.CurrentBarrier())
	return clone, nil
}

// GetSquaddieByID returns the Squaddie based on the one with the given squaddieID.
func (repository *Repository) GetSquaddieByID(squaddieID string) *Squaddie {
	squaddie, squaddieExists := repository.squaddiesByID[squaddieID]
	if !squaddieExists {
		return nil
	}

	clonedSquaddie, cloneErr := repository.CloneSquaddieWithNewID(squaddie, "")
	if cloneErr != nil {
		return nil
	}
	return clonedSquaddie
}

// GetOriginalSquaddieByID returns the stored Squaddie based on the squaddieID.
func (repository *Repository) GetOriginalSquaddieByID(squaddieID string) *Squaddie {
	squaddie, _ := repository.squaddiesByID[squaddieID]
	return squaddie
}
