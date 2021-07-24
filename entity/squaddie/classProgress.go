package squaddie

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	"github.com/cserrant/terosBattleServer/utility"
)

// ClassProgress tracks the ClassProgress's current class and any levels they have taken so far.
type ClassProgress struct {
	BaseClassID         		string                          `json:"base_class" yaml:"base_class"`
	CurrentClass        		string                         `json:"current_class" yaml:"current_class"`
	ClassLevelsConsumed 		map[string]*ClassLevelsConsumed `json:"class_levels" yaml:"class_levels"`
}

// AddClass gives the ClassProgress a new class it can gain levels in.
func (classProgress *ClassProgress) AddClass(class *squaddieclass.Class) {
	classProgress.ClassLevelsConsumed[class.ID] = &ClassLevelsConsumed{
		ClassID:        class.ID,
		ClassName:      class.Name,
		LevelsConsumed: []string{},
	}
}

// GetLevelCountsByClass returns a mapping of class names to levels gained.
func (classProgress *ClassProgress) GetLevelCountsByClass() map[string]int {
	count := map[string]int{}
	for classID, progress := range classProgress.ClassLevelsConsumed {
		count[classID] = len(progress.LevelsConsumed)
	}

	return count
}

// MarkLevelUpBenefitAsConsumed makes the ClassProgress remember it used this benefit to level up already.
func (classProgress *ClassProgress) MarkLevelUpBenefitAsConsumed(benefitClassID, benefitID string)  {
	classProgress.ClassLevelsConsumed[benefitClassID].LevelsConsumed = append(classProgress.ClassLevelsConsumed[benefitClassID].LevelsConsumed, benefitID)
}

// SetClass changes the ClassProgress's CurrentClass to the given classID.
//   It also sets the BaseClass if it hasn't been already.
//   Raises an error if classID has not been added to the squaddie yet.
func (classProgress *ClassProgress) SetClass(classID string) error {
	if _, exists := classProgress.ClassLevelsConsumed[classID]; !exists {
		newError := fmt.Errorf(`cannot switch to unknown class "%s"`, classID)
		utility.Log(newError.Error(),0, utility.Error)
		return newError
	}

	if classProgress.BaseClassID == "" {
		classProgress.BaseClassID = classID
	}

	classProgress.CurrentClass = classID
	return nil
}

// SetBaseClassIfNoBaseClass sets the BaseClass if it hasn't been already.
func (classProgress *ClassProgress) SetBaseClassIfNoBaseClass(classID string) {
	if classProgress.BaseClassID == "" {
		classProgress.BaseClassID = classID
	}
}

// IsClassLevelAlreadyUsed returns true if a LevelUpBenefit with the given ID has already been used.
func (classProgress *ClassProgress) IsClassLevelAlreadyUsed(benefitID string) bool {
	return classProgress.anyClassLevelsConsumed(func(classID string, progress *ClassLevelsConsumed) bool {
		return progress.IsLevelAlreadyConsumed(benefitID)
	})
}

// HasAddedClass returns true if the ClassProgress has already added a class with the name classIDToFind
func (classProgress *ClassProgress) HasAddedClass(classIDToFind string) bool {
	return classProgress.anyClassLevelsConsumed(func(classID string, progress *ClassLevelsConsumed) bool {
		return classID == classIDToFind
	})
}

// anyClassLevelsConsumed returns true if any of the squaddie's class levels consumed satisfies a given condition.
func (classProgress *ClassProgress) anyClassLevelsConsumed(condition func(classID string, progress *ClassLevelsConsumed) bool) bool {
	for classID, progress := range classProgress.ClassLevelsConsumed {
		if condition(classID, progress) {
			return true
		}
	}
	return false
}