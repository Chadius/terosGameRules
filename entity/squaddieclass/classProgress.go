package squaddieclass

import (
	"fmt"
	"github.com/chadius/terosgamerules/utility"
)

// ClassProgress tracks the ClassProgress's current class and any levels they have taken so far.
type ClassProgress struct {
	ClassProgressBaseClassID         string                          `json:"base_class" yaml:"base_class"`
	ClassProgressCurrentClassID      string                          `json:"current_class" yaml:"current_class"`
	ClassProgressClassLevelsConsumed map[string]*ClassLevelsConsumed `json:"class_levels" yaml:"class_levels"`
}

//NewClassProgress returns a new object.
func NewClassProgress(baseClassID, currentClassID string, levelsConsumedByClassID map[string]*ClassLevelsConsumed) *ClassProgress {
	newClassProgress := &ClassProgress{
		ClassProgressBaseClassID:         baseClassID,
		ClassProgressCurrentClassID:      currentClassID,
		ClassProgressClassLevelsConsumed: map[string]*ClassLevelsConsumed{},
	}

	for classID, levelsConsumed := range levelsConsumedByClassID {
		newClassProgress.AddClass(&ClassReference{
			ID:   levelsConsumed.GetClassID(),
			Name: levelsConsumed.GetClassName(),
		})
		for _, levelID := range levelsConsumed.GetLevelsConsumed() {
			newClassProgress.MarkLevelUpBenefitAsConsumed(classID, levelID)
		}
	}

	return newClassProgress
}

// AddClass gives the ClassProgress a new class it can gain levels in, if it wasn't already added.
func (classProgress *ClassProgress) AddClass(classReference *ClassReference) {
	if classProgress.ClassProgressClassLevelsConsumed[classReference.ID] != nil {
		return
	}
	classProgress.ClassProgressClassLevelsConsumed[classReference.ID] = NewClassLevelsConsumed(classReference.ID, classReference.Name, []string{})
}

// GetLevelCountsByClass returns a mapping of class names to levels gained.
func (classProgress *ClassProgress) GetLevelCountsByClass() map[string]int {
	count := map[string]int{}
	for classID, progress := range classProgress.ClassProgressClassLevelsConsumed {
		count[classID] = len(progress.GetLevelsConsumed())
	}

	return count
}

// MarkLevelUpBenefitAsConsumed makes the ClassProgress remember it used this benefit to level up already.
func (classProgress *ClassProgress) MarkLevelUpBenefitAsConsumed(benefitClassID, benefitID string) {
	classProgress.ClassProgressClassLevelsConsumed[benefitClassID].MarkLevelUpBenefitAsConsumed(benefitID)
}

// SetClass changes the ClassProgress's ClassProgressCurrentClassID to the given classID.
//   It also sets the BaseClass if it hasn't been already.
//   Raises an error if classID has not been added to the squaddie yet.
func (classProgress *ClassProgress) SetClass(classID string) error {
	if _, exists := classProgress.ClassProgressClassLevelsConsumed[classID]; !exists {
		newError := fmt.Errorf(`cannot switch to unknown class "%s"`, classID)
		utility.Log(newError.Error(), 0, utility.Error)
		return newError
	}

	if classProgress.ClassProgressBaseClassID == "" {
		classProgress.ClassProgressBaseClassID = classID
	}

	classProgress.ClassProgressCurrentClassID = classID
	return nil
}

// SetBaseClassIfNoBaseClass sets the BaseClass if it hasn't been already.
func (classProgress *ClassProgress) SetBaseClassIfNoBaseClass(classID string) {
	if classProgress.ClassProgressBaseClassID == "" {
		classProgress.ClassProgressBaseClassID = classID
	}
}

// IsClassLevelAlreadyUsed returns true if a LevelUpBenefit with the given id has already been used.
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
	for classID, progress := range classProgress.ClassProgressClassLevelsConsumed {
		if condition(classID, progress) {
			return true
		}
	}
	return false
}

// BaseClassID returns the base class ID.
func (classProgress *ClassProgress) BaseClassID() string {
	return classProgress.ClassProgressBaseClassID
}

// CurrentClassID returns the current class ID.
func (classProgress *ClassProgress) CurrentClassID() string {
	return classProgress.ClassProgressCurrentClassID
}

// ClassLevelsConsumed returns the current class ID.
func (classProgress *ClassProgress) ClassLevelsConsumed() *map[string]*ClassLevelsConsumed {
	return &classProgress.ClassProgressClassLevelsConsumed
}
