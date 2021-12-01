package levelupbenefit

import (
	"encoding/json"
	"fmt"
	"github.com/chadius/terosbattleserver/utility"
	"gopkg.in/yaml.v2"
)

// Repository is used to load and retrieve LevelUpBenefit objects for
//   squaddies, classes and levels.
type Repository struct {
	levelUpBenefitsByClassID map[string][]*LevelUpBenefit
}

// GetNumberOfLevelUpBenefits returns a total count of the LevelUpBenefit objects stored.
func (repository *Repository) GetNumberOfLevelUpBenefits() int {
	count := 0
	for _, levelUpBenefits := range repository.levelUpBenefitsByClassID {
		count = count + len(levelUpBenefits)
	}
	return count
}

// AddJSONSource consumes a given bytestream and tries to analyze it.
func (repository *Repository) AddJSONSource(data []byte) (bool, error) {
	return repository.addSource(data, json.Unmarshal)
}

// AddYAMLSource consumes a given bytestream and tries to analyze it.
func (repository *Repository) AddYAMLSource(data []byte) (bool, error) {
	return repository.addSource(data, yaml.Unmarshal)
}

// AddLevels adds a list of LevelUpBenefits to the repository.
func (repository *Repository) AddLevels(allBenefits []*LevelUpBenefit) (bool, error) {
	for _, levelUpBenefit := range allBenefits {
		success, err := repository.tryToAddLevelUpBenefitToSource(levelUpBenefit)
		if success == false {
			return false, err
		}
	}
	return true, nil
}

// AddSource consumes a given bytestream of the given sourceType and tries to analyze it.
func (repository *Repository) addSource(data []byte, unmarshal utility.UnmarshalFunc) (bool, error) {
	var unmarshalError error

	var allBenefits []LevelUpBenefit

	unmarshalError = unmarshal(data, &allBenefits)

	if unmarshalError != nil {
		return false, unmarshalError
	}

	for _, levelUpBenefit := range allBenefits {
		success, err := repository.tryToAddLevelUpBenefitToSource(&levelUpBenefit)
		if success == false {
			return false, err
		}
	}

	return true, nil
}

func (repository *Repository) tryToAddLevelUpBenefitToSource(levelUpBenefit *LevelUpBenefit) (bool, error) {
	err := levelUpBenefit.CheckForErrors()
	if err != nil {
		return false, err
	}

	classID := levelUpBenefit.ClassID()

	if repository.levelUpBenefitsByClassID[classID] == nil {
		repository.levelUpBenefitsByClassID[classID] = []*LevelUpBenefit{}
	}

	repository.levelUpBenefitsByClassID[classID] =
		append(repository.levelUpBenefitsByClassID[classID], levelUpBenefit)
	return true, nil
}

// GetLevelUpBenefitsByClassID uses the squaddieName and className to return a list of Level Up Benefits.
func (repository *Repository) GetLevelUpBenefitsByClassID(classID string) ([]*LevelUpBenefit, error) {

	classBenefits, classExists := repository.levelUpBenefitsByClassID[classID]
	if !classExists {
		newError := fmt.Errorf(`no LevelUpBenefits for this class SquaddieID: "%s"`, classID)
		utility.Log(newError.Error(), 0, utility.Error)
		return nil, newError
	}

	return classBenefits, nil
}

// GetLevelUpBenefitsForClassByType returns the levels of a class by type (Small or Big).
func (repository *Repository) GetLevelUpBenefitsForClassByType(classID string) (map[Size][]*LevelUpBenefit, error) {
	levelsInClassByType := map[Size][]*LevelUpBenefit{Small: []*LevelUpBenefit{}, Big: []*LevelUpBenefit{}}
	levelsInClass, err := repository.GetLevelUpBenefitsByClassID(classID)

	if err != nil {
		return levelsInClassByType, err
	}
	for _, level := range levelsInClass {
		levelsInClassByType[level.LevelUpBenefitType()] = append(levelsInClassByType[level.LevelUpBenefitType()], level)
	}
	return levelsInClassByType, nil
}

// AddBuilderYAML accepts a YAML data stream to create levelUpBenefits, using the LevelUpBenefitBuilder.
func (repository *Repository) AddBuilderYAML(stream []byte) error {
	_, err := repository.addSourceForBuilders(stream, yaml.Unmarshal)
	return err
}

// AddBuilderJSON accepts a JSON data stream to create levelUpBenefits, using the LevelUpBenefitBuilder.
func (repository *Repository) AddBuilderJSON(stream []byte) error {
	_, err := repository.addSourceForBuilders(stream, json.Unmarshal)
	return err
}

// addSourceForBuilders consumes a given stream of the given sourceType and tries to analyze it.
func (repository *Repository) addSourceForBuilders(data []byte, unmarshal utility.UnmarshalFunc) (bool, error) {
	var unmarshalError error

	var builderInstructions []BuilderMarshal

	unmarshalError = unmarshal(data, &builderInstructions)

	if unmarshalError != nil {
		return false, unmarshalError
	}

	for _, instruction := range builderInstructions {
		levelUpBenefit, buildError := NewBuilderFromMarshal(instruction).Build()
		if buildError != nil {
			return false, buildError
		}

		success, err := repository.tryToAddLevelUpBenefitToSource(levelUpBenefit)
		if success == false {
			return false, err
		}
	}

	return true, nil
}

// NewLevelUpBenefitRepository generates a pointer to a new LevelUpBenefitRepository.
func NewLevelUpBenefitRepository() *Repository {
	repository := Repository{
		map[string][]*LevelUpBenefit{},
	}
	return &repository
}
