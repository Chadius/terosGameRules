package squaddie

import (
	"fmt"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	"github.com/chadius/terosbattleserver/utility"
)

// ClassProgress tracks the ClassProgress's current class and any levels they have taken so far.
type ClassProgress struct {
	ClassProgressBaseClassID         string                          `json:"base_class" yaml:"base_class"`
	ClassProgressCurrentClassID      string                          `json:"current_class" yaml:"current_class"`
	ClassProgressClassLevelsConsumed map[string]*ClassLevelsConsumed `json:"class_levels" yaml:"class_levels"`
}

// AddClass gives the ClassProgress a new class it can gain levels in, if it wasn't already added.
func (classProgress *ClassProgress) AddClass(classReference *squaddieclass.ClassReference) {
	if classProgress.ClassProgressClassLevelsConsumed[classReference.ID] != nil {
		return
	}

	classProgress.ClassProgressClassLevelsConsumed[classReference.ID] = &ClassLevelsConsumed{
		ClassID:        classReference.ID,
		ClassName:      classReference.Name,
		LevelsConsumed: []string{},
	}
}

// GetLevelCountsByClass returns a mapping of class names to levels gained.
func (classProgress *ClassProgress) GetLevelCountsByClass() map[string]int {
	count := map[string]int{}
	for classID, progress := range classProgress.ClassProgressClassLevelsConsumed {
		count[classID] = len(progress.LevelsConsumed)
	}

	return count
}

// MarkLevelUpBenefitAsConsumed makes the ClassProgress remember it used this benefit to level up already.
func (classProgress *ClassProgress) MarkLevelUpBenefitAsConsumed(benefitClassID, benefitID string) {
	classProgress.ClassProgressClassLevelsConsumed[benefitClassID].LevelsConsumed = append(classProgress.ClassProgressClassLevelsConsumed[benefitClassID].LevelsConsumed, benefitID)
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

// IsClassLevelAlreadyUsed returns true if a LevelUpBenefit with the given SquaddieID has already been used.
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

// HasSameClassesAs sees if the other class progress has the same fields.
func (classProgress *ClassProgress) HasSameClassesAs(other *ClassProgress) bool {
	if classProgress.BaseClassID() != other.BaseClassID() {
		return false
	}
	if classProgress.CurrentClassID() != other.CurrentClassID() {
		return false
	}

	otherClassLevelsConsumed := *other.ClassLevelsConsumed()
	if len(*classProgress.ClassLevelsConsumed()) != len(otherClassLevelsConsumed) {
		return false
	}

	classLevelsConsumedByClassID := map[string]bool{}
	for classLevelsConsumedClassID := range *classProgress.ClassLevelsConsumed() {
		classLevelsConsumedByClassID[classLevelsConsumedClassID] = false
	}

	for classID, classLevelsConsumed := range otherClassLevelsConsumed {
		_, exists := classLevelsConsumedByClassID[classID]
		if !exists {
			return false
		}
		if !classProgress.ClassProgressClassLevelsConsumed[classID].HasSameConsumptionAs(classLevelsConsumed) {
			return false
		}
		classLevelsConsumedByClassID[classID] = true
	}

	for _, wasFound := range classLevelsConsumedByClassID {
		if wasFound == false {
			return false
		}
	}

	return true
}
