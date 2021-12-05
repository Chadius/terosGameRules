package squaddie

import (
	"encoding/json"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
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

// AddJSONSource consumes a given bytestream and tries to analyze it.
func (repository *Repository) AddJSONSource(data []byte) (bool, error) {
	return repository.addSource(data, json.Unmarshal)
}

// AddYAMLSource consumes a given bytestream and tries to analyze it.
func (repository *Repository) AddYAMLSource(data []byte) (bool, error) {
	return repository.addSource(data, yaml.Unmarshal)
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

// AddSource consumes a given bytestream of the given sourceType and tries to analyze it.
func (repository *Repository) addSource(data []byte, unmarshal utility.UnmarshalFunc) (bool, error) {
	var unmarshalError error
	var listOfSquaddies []Squaddie
	unmarshalError = unmarshal(data, &listOfSquaddies)

	if unmarshalError != nil {
		return false, unmarshalError
	}
	for index := range listOfSquaddies {
		newSquaddie := listOfSquaddies[index]
		newSquaddie.Defense.SetHPToMax()
		success, err := repository.tryToAddSquaddie(&newSquaddie)
		if success == false {
			return false, err
		}
	}
	return true, nil
}

// AddYAMLSourceUsingSquaddieBuilder adds multiple squaddies using builder objects and a YAML data stream.
func (repository *Repository) AddYAMLSourceUsingSquaddieBuilder(data []byte) error {
	_, err := repository.addSourceUsingSquaddieBuilder(data, yaml.Unmarshal)
	return err
}

// AddJSONSourceUsingSquaddieBuilder adds multiple squaddies using builder objects and a JSON data stream.
func (repository *Repository) AddJSONSourceUsingSquaddieBuilder(data []byte) error {
	_, err := repository.addSourceUsingSquaddieBuilder(data, json.Unmarshal)
	return err
}

// addSourceUsingSquaddieBuilder reads the byte stream to create new squaddies.
func (repository *Repository) addSourceUsingSquaddieBuilder(data []byte, unmarshal utility.UnmarshalFunc) (bool, error) {
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
		squaddieToAdd.Identification.SetNewIDToRandom()
	}
	repository.squaddiesByID[squaddieToAdd.ID()] = squaddieToAdd
	return true, nil
}

// GetNumberOfSquaddies returns the number of Squaddies ready to retrieve.
func (repository *Repository) GetNumberOfSquaddies() int {
	return len(repository.squaddiesByID)
}

//CloneSquaddieWithNewID uses the base Squaddie to create a new one.
//  All fields will be the same except the SquaddieID.
//  If newID isn't empty, the clone SquaddieID is set to that.
//  Otherwise, it is randomly generated.
func (repository *Repository) CloneSquaddieWithNewID(base *Squaddie, newID string) (*Squaddie, error) {
	clone := NewSquaddie(base.Name())

	cloneSquaddieID := clone.ID()
	if newID != "" {
		cloneSquaddieID = newID
	}

	clone.Identification = *NewIdentification(cloneSquaddieID, base.Name(), base.Affiliation())
	clone.Offense = *NewOffense(base.Aim(), base.Strength(), base.Mind())
	clone.Defense = *NewDefense(base.CurrentHitPoints(), base.MaxHitPoints(), base.Dodge(), base.Deflect(), base.CurrentBarrier(), base.MaxBarrier(), base.Armor())
	clone.Movement = *NewMovement(base.MovementDistance(), base.MovementType(), base.MovementCanHitAndRun())

	for _, reference := range base.PowerCollection.GetCopyOfPowerReferences() {
		clone.AddPowerReference(reference)
	}

	clone.ClassProgress = *squaddieclass.NewClassProgress(
		base.BaseClassID(),
		base.CurrentClassID(),
		*base.ClassLevelsConsumed(),
	)
	return clone, nil
}

// GetSquaddieByID returns the Squaddie based on the one with the given SquaddieID.
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

// GetOriginalSquaddieByID returns the stored Squaddie based on the SquaddieID.
func (repository *Repository) GetOriginalSquaddieByID(squaddieID string) *Squaddie {
	squaddie, _ := repository.squaddiesByID[squaddieID]
	return squaddie
}
