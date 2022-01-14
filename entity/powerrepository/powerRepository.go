package powerrepository

import (
	"errors"
	"github.com/chadius/terosgamerules/entity/power"
	"github.com/chadius/terosgamerules/entity/powerinterface"
	"github.com/chadius/terosgamerules/utility"
)

// Repository will interact with external devices to manage Powers.
type Repository struct {
	powersByID map[string]powerinterface.Interface
}

// NewPowerRepository generates a pointer to a new Repository.
func NewPowerRepository() *Repository {
	repository := Repository{
		map[string]powerinterface.Interface{},
	}
	return &repository
}

// AddJSONSource consumes a given bytestream and tries to analyze it.
func (repository *Repository) AddJSONSource(data []byte) (bool, error) {
	builderOptions := power.CreatePowerBuilderOptionsFromJSON(data)
	if builderOptions == nil {
		return false, errors.New("could not create Builder with given JSON")
	}

	powersToAdd := []powerinterface.Interface{}
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
		return false, errors.New("could not create Builder with given YAML")
	}

	powersToAdd := []powerinterface.Interface{}
	for _, option := range builderOptions {
		newPower := option.Build()
		powersToAdd = append(powersToAdd, newPower)
	}

	return repository.AddSlicePowerSource(powersToAdd)
}

// AddSlicePowerSource tries to add the slice of powers to the repo.
func (repository *Repository) AddSlicePowerSource(powersToAdd []powerinterface.Interface) (bool, error) {
	for _, powerToAdd := range powersToAdd {
		success, err := repository.tryToAddPower(powerToAdd)
		if success == false {
			return false, err
		}
	}
	return true, nil
}

// AddPower tries to add a single power to the repository.
func (repository *Repository) AddPower(powerToAdd powerinterface.Interface) (bool, error) {
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

func (repository *Repository) tryToAddPower(powerToAdd powerinterface.Interface) (bool, error) {
	repository.powersByID[powerToAdd.ID()] = powerToAdd
	return true, nil
}

// GetNumberOfPowers returns the number of Powers ready to retrieve.
func (repository *Repository) GetNumberOfPowers() int {
	return len(repository.powersByID)
}

// GetPowerByID returns the Power stored by powerID.
func (repository *Repository) GetPowerByID(powerID string) powerinterface.Interface {
	return repository.powersByID[powerID]
}

// GetAllPowersByName returns a slice of powers based on a given name
func (repository *Repository) GetAllPowersByName(powerNameToFind string) []powerinterface.Interface {
	powersFound := []powerinterface.Interface{}
	for _, power := range repository.powersByID {
		if power.Name() == powerNameToFind {
			powersFound = append(powersFound, power)
		}
	}

	return powersFound
}
