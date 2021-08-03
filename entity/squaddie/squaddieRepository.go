package squaddie

import (
	"encoding/json"
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/utility"
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
		success, err := repository.tryToAddSquaddie(&newSquaddie)
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
	squaddieToAdd.Defense.SetHPToMax()

	if squaddieToAdd.Identification.ID == "" {
		squaddieToAdd.Identification.SetNewIDToRandom()
	}
	repository.squaddiesByID[squaddieToAdd.Identification.ID] = squaddieToAdd
	return true, nil
}

// GetNumberOfSquaddies returns the number of Squaddies ready to retrieve.
func (repository *Repository) GetNumberOfSquaddies() int {
	return len(repository.squaddiesByID)
}

// MarshalSquaddieIntoJSON converts the given Squaddie into JSON.
func (repository *Repository) MarshalSquaddieIntoJSON(squaddie *Squaddie) ([]byte, error) {
	type Alias Squaddie

	return json.Marshal(&struct {
		*Alias
		PowerIDNames []*power.Reference `json:"powers" yaml:"powers"`
	}{
		Alias:        (*Alias)(squaddie),
		PowerIDNames: squaddie.PowerCollection.GetInnatePowerIDNames(),
	})
}

//CloneSquaddieWithNewID uses the base Squaddie to create a new one.
//  All fields will be the same except the ID.
//  If newID isn't empty, the clone ID is set to that.
//  Otherwise it is randomly generated.
func (repository *Repository) CloneSquaddieWithNewID(base *Squaddie, newID string) (*Squaddie, error) {
	clone := NewSquaddie(base.Identification.Name)
	clone.Identification.Affiliation = base.Identification.Affiliation
	if newID != "" {
		clone.Identification.ID = newID
	}

	clone.Defense.CurrentHitPoints = base.Defense.CurrentHitPoints
	clone.Defense.MaxHitPoints = base.Defense.MaxHitPoints
	clone.Offense.Aim = base.Offense.Aim
	clone.Offense.Strength = base.Offense.Strength
	clone.Offense.Mind = base.Offense.Mind
	clone.Defense.Dodge = base.Defense.Dodge
	clone.Defense.Deflect = base.Defense.Deflect
	clone.Defense.CurrentBarrier = base.Defense.CurrentBarrier
	clone.Defense.MaxBarrier = base.Defense.MaxBarrier
	clone.Defense.Armor = base.Defense.Armor

	clone.Movement.Distance = base.Movement.Distance
	clone.Movement.Type = base.Movement.Type
	clone.Movement.HitAndRun = base.Movement.HitAndRun

	clone.PowerCollection.PowerReferences = append([]*power.Reference{}, base.PowerCollection.PowerReferences...)

	clone.ClassProgress.BaseClassID = base.ClassProgress.BaseClassID
	clone.ClassProgress.CurrentClass = base.ClassProgress.CurrentClass

	for classID, progress := range base.ClassProgress.ClassLevelsConsumed {
		newProgress := ClassLevelsConsumed{
			ClassID:        classID,
			ClassName:      progress.ClassName,
			LevelsConsumed: append([]string{}, progress.LevelsConsumed...),
		}

		clone.ClassProgress.ClassLevelsConsumed[classID] = &newProgress
	}
	return clone, nil
}

// CloneAndRenameSquaddie clones the base squaddie and renames them
//  newName must be non-empty or raise an error.
//  See CloneSquaddieWithNewID to see how newID is used.
func (repository *Repository) CloneAndRenameSquaddie(base *Squaddie, newName string, newID string) (*Squaddie, error) {
	clone, err := repository.CloneSquaddieWithNewID(base, newID)
	if err != nil {
		return nil, err
	}

	if newName == "" {
		newError := fmt.Errorf(`cannot clone squaddie "%s" without a name`, base.Identification.Name)
		utility.Log(newError.Error(),0, utility.Error)
		return nil, newError
	}

	clone.Identification.Name = newName
	return clone, nil
}

// GetSquaddieByID returns the Squaddie based on the one with the given ID.
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

// GetOriginalSquaddieByID returns the stored Squaddie based on the ID.
func (repository *Repository) GetOriginalSquaddieByID(squaddieID string) *Squaddie {
	squaddie, _ := repository.squaddiesByID[squaddieID]
	return squaddie
}

