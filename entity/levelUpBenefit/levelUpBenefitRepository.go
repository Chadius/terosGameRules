package levelUpBenefit

import (
	"encoding/json"
	"fmt"
	"github.com/cserrant/terosBattleServer/utility"
	"gopkg.in/yaml.v2"
)

// BenefitsByNameAndClass is the heading that explains all of the ways a Squaddie with the given SquaddieName
//   can advance.
type BenefitsByNameAndClass struct {
	SquaddieName string                `json:"squaddie_name" yaml:"squaddie_name"`
	LevelUpsByClass []*BenefitsByClass `json:"level_ups_by_class" yaml:"level_ups_by_class"`
}

// BenefitsByClass describes each class and its potential level ups.
type BenefitsByClass struct {
	ClassName        string            `json:"class_name" yaml:"class_name"`
	LevelUpBenefits  []*LevelUpBenefit `json:"level_up_benefits" yaml:"level_up_benefits"`
}

// Repository is used to load and retrieve LevelUpBenefit objects for
//   squaddies, classes and levels.
type Repository struct {
	levelUpBenefitsByNameAndClass map[string]map[string][]*LevelUpBenefit
}

// GetNumberOfLevelUpBenefits returns a total count of all of the LevelUpBenefit objects stored.
func (repository *Repository) GetNumberOfLevelUpBenefits() int {
	count := 0
	for _, levelUpBenefitsByClass := range repository.levelUpBenefitsByNameAndClass {
		for _, levelUpBenefits := range levelUpBenefitsByClass {
			count = count + len(levelUpBenefits)
		}
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

// AddSource consumes a given bytestream of the given sourceType and tries to analyze it.
func (repository *Repository) addSource(data []byte, unmarshal utility.UnmarshalFunc) (bool, error) {
	var unmarshalError error

	var allBenefits []BenefitsByNameAndClass

	unmarshalError = unmarshal(data, &allBenefits)

	if unmarshalError != nil {
		return false, unmarshalError
	}

	for _, squaddieBenefits := range allBenefits {
		squaddieName := squaddieBenefits.SquaddieName
		for _, benefitsByClass := range squaddieBenefits.LevelUpsByClass {
			className := benefitsByClass.ClassName
			for _, levelUpBenefit := range benefitsByClass.LevelUpBenefits {
				err := levelUpBenefit.CheckForErrors()
				if err != nil {
					return false, err
				}

				if repository.levelUpBenefitsByNameAndClass[squaddieName] == nil {
					repository.levelUpBenefitsByNameAndClass[squaddieName] = map[string][]*LevelUpBenefit{}
				}
				if repository.levelUpBenefitsByNameAndClass[squaddieName][className] == nil {
					repository.levelUpBenefitsByNameAndClass[squaddieName][className] = []*LevelUpBenefit{}
				}

				levelUpBenefit.SquaddieName = squaddieName
				levelUpBenefit.ClassName = className

				repository.levelUpBenefitsByNameAndClass[squaddieName][className] =
					append(repository.levelUpBenefitsByNameAndClass[squaddieName][className], levelUpBenefit)
			}
		}
	}
	return true, nil
}

// GetLevelUpBenefitsByNameAndClass uses the squaddieName and className to return a list of Level Up Benefits.
func (repository *Repository) GetLevelUpBenefitsByNameAndClass(squaddieName string, className string) ([]*LevelUpBenefit, error) {
	_, squaddieExists := repository.levelUpBenefitsByNameAndClass[squaddieName]
	if !squaddieExists {
		return nil, fmt.Errorf(`no LevelUpBenefits for this squaddie: "%s"`, squaddieName)
	}

	classBenefits, classExists := repository.levelUpBenefitsByNameAndClass[squaddieName][className]
	if !classExists {
		return nil, fmt.Errorf(`no LevelUpBenefits for this class: "%s"`, className)
	}

	return classBenefits, nil
}

// NewLevelUpBenefitRepository generates a pointer to a new LevelUpBenefitRepository.
func NewLevelUpBenefitRepository() *Repository {
	repository := Repository{
		map[string]map[string][]*LevelUpBenefit{},
	}
	return &repository
}

