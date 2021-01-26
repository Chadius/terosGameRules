package repository

import (
	"encoding/json"
	"fmt"
	"github.com/cserrant/terosBattleServer/entity"
	"gopkg.in/yaml.v2"
)

// BenefitsByNameAndClass is the heading that explains all of the ways a Squaddie with the given SquaddieName
//   can advance.
type BenefitsByNameAndClass struct {
	SquaddieName string                         `json:"squaddie_name" yaml:"squaddie_name"`
	LevelUpsByClass []*BenefitsByClass 			`json:"level_ups_by_class" yaml:"level_ups_by_class"`
}

// BenefitsByClass describes each class and its potential level ups.
type BenefitsByClass struct {
	ClassName        string                      `json:"class_name" yaml:"class_name"`
	LevelUpBenefits  []*entity.LevelUpBenefit    `json:"level_up_benefits" yaml:"level_up_benefits"`
}

// LevelUpBenefitRepository is used to load and retrieve LevelUpBenefit objects for
//   squaddies, classes and levels.
type LevelUpBenefitRepository struct {
	levelUpBenefitsByNameAndClass map[string]map[string][]*entity.LevelUpBenefit
}

// GetNumberOfLevelUpBenefits returns a total count of all of the LevelUpBenefit objects stored.
func (repository *LevelUpBenefitRepository) GetNumberOfLevelUpBenefits() (int) {
	count := 0
	for _, levelUpBenefitsByClass := range repository.levelUpBenefitsByNameAndClass {
		for _, levelUpBenefits := range levelUpBenefitsByClass {
			count = count + len(levelUpBenefits)
		}
	}
	return count
}

// AddJSONSource consumes a given bytestream and tries to analyze it.
func (repository *LevelUpBenefitRepository) AddJSONSource(data []byte) (bool, error) {
	return repository.addSource(data, json.Unmarshal)
}

// AddYAMLSource consumes a given bytestream and tries to analyze it.
func (repository *LevelUpBenefitRepository) AddYAMLSource(data []byte) (bool, error) {
	return repository.addSource(data, yaml.Unmarshal)
}

// AddSource consumes a given bytestream of the given sourceType and tries to analyze it.
func (repository *LevelUpBenefitRepository) addSource(data []byte, unmarshal unmarshalFunc) (bool, error) {
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
					repository.levelUpBenefitsByNameAndClass[squaddieName] = map[string][]*entity.LevelUpBenefit{}
				}
				if repository.levelUpBenefitsByNameAndClass[squaddieName][className] == nil {
					repository.levelUpBenefitsByNameAndClass[squaddieName][className] = []*entity.LevelUpBenefit{}
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
func (repository *LevelUpBenefitRepository) GetLevelUpBenefitsByNameAndClass(squaddieName string, className string) ([]*entity.LevelUpBenefit, error) {
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
func NewLevelUpBenefitRepository() *LevelUpBenefitRepository {
	repository := LevelUpBenefitRepository{
		map[string]map[string][]*entity.LevelUpBenefit{},
	}
	return &repository
}
