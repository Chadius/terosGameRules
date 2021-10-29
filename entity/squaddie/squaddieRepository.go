package squaddie

import (
	"encoding/json"
	"github.com/chadius/terosbattleserver/entity/power"
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
	clone.Identification = Identification{
		SquaddieID: newID,
		SquaddieName: base.Name(),
		SquaddieAffiliation: base.Affiliation(),
	}

	clone.Offense = Offense{
		SquaddieAim: base.Aim(),
		SquaddieStrength: base.Strength(),
		SquaddieMind: base.Mind(),
	}

	clone.Defense = Defense{
		SquaddieCurrentHitPoints: base.CurrentHitPoints(),
		SquaddieMaxHitPoints: base.MaxHitPoints(),
		SquaddieDodge: base.Dodge(),
		SquaddieDeflect: base.Deflect(),
		SquaddieCurrentBarrier: base.CurrentBarrier(),
		SquaddieMaxBarrier: base.MaxBarrier(),
		SquaddieArmor: base.Armor(),
	}

	clone.Movement = Movement{
		SquaddieMovementDistance: base.MovementDistance(),
		SquaddieMovementType: base.MovementType(),
		SquaddieMovementCanHitAndRun: base.MovementCanHitAndRun(),
	}

	clone.PowerCollection.PowerReferences = append([]*power.Reference{}, base.PowerCollection.PowerReferences...)

	clone.ClassProgress.BaseClassID = base.ClassProgress.BaseClassID
	clone.ClassProgress.CurrentClassID = base.ClassProgress.CurrentClassID

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
