package powerrepository

import (
	"errors"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/utility"
)

// Repository will interact with external devices to manage Powers.
type Repository struct {
	powersByID map[string]*power.Power
}

// NewPowerRepository generates a pointer to a new Repository.
func NewPowerRepository() *Repository {
	repository := Repository{
		map[string]*power.Power{},
	}
	return &repository
}

// AddJSONSource consumes a given bytestream and tries to analyze it.
func (repository *Repository) AddJSONSource(data []byte) (bool, error) {
	builderOptions := power.CreatePowerBuilderOptionsFromJSON(data)
	if builderOptions == nil {
		return false, errors.New("could not create BuilderOptions with given JSON")
	}

	powersToAdd := []*power.Power{}
	for _, option := range builderOptions {
		newPower := option.Build()
		powersToAdd = append(powersToAdd, newPower)
	}

	return repository.AddSlicePowerSource(powersToAdd)
}

// AddYAMLSource consumes a given bytestream and tries to analyze it.
func (repository *Repository) AddYAMLSource(data []byte) (bool, error) {
	builderOptions := power.CreatePowerBuilderOptionsFromYAML(data)
	if builderOptions == nil {
		return false, errors.New("could not create BuilderOptions with given YAML")
	}

	powersToAdd := []*power.Power{}
	for _, option := range builderOptions {
		newPower := option.Build()
		powersToAdd = append(powersToAdd, newPower)
	}

	return repository.AddSlicePowerSource(powersToAdd)
}

// AddSlicePowerSource tries to add the slice of powers to the repo.
func (repository *Repository) AddSlicePowerSource(powersToAdd []*power.Power) (bool, error) {
	for _, powerToAdd := range powersToAdd {
		success, err := repository.tryToAddPower(powerToAdd)
		if success == false {
			return false, err
		}
	}
	return true, nil
}

// AddPower tries to add a single power to the repository.
func (repository *Repository) AddPower(powerToAdd *power.Power) (bool, error) {
	success, err := repository.tryToAddPower(powerToAdd)
	return success, err
}

// addSource consumes a given bytestream of the given sourceType and tries to analyze it.
func (repository *Repository) addSource(data []byte, unmarshal utility.UnmarshalFunc) (bool, error) {
	var unmarshalError error
	var listOfPowers []power.Power

	unmarshalError = unmarshal(data, &listOfPowers)

	if unmarshalError != nil {
		return false, unmarshalError
	}
	for index := range listOfPowers {
		powerToAdd := listOfPowers[index]
		success, err := repository.tryToAddPower(&powerToAdd)
		if success == false {
			return false, err
		}
	}
	return true, nil
}

func (repository *Repository) tryToAddPower(powerToAdd *power.Power) (bool, error) {
	PowerErr := power.CheckPowerForErrors(powerToAdd)
	if PowerErr != nil {
		return false, PowerErr
	}
	repository.powersByID[powerToAdd.ID()] = powerToAdd
	return true, nil
}

// GetNumberOfPowers returns the number of Powers ready to retrieve.
func (repository *Repository) GetNumberOfPowers() int {
	return len(repository.powersByID)
}

// GetPowerByID returns the Power stored by SquaddieID.
func (repository *Repository) GetPowerByID(powerID string) *power.Power {
	return repository.powersByID[powerID]
}

// GetAllPowersByName returns a slice of powers based on a given name
func (repository *Repository) GetAllPowersByName(powerNameToFind string) []*power.Power {
	powersFound := []*power.Power{}
	for _, power := range repository.powersByID {
		if power.Name() == powerNameToFind {
			powersFound = append(powersFound, power)
		}
	}

	return powersFound
}
