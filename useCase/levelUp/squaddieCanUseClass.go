package levelup

import (
	"github.com/chadius/terosbattleserver/entity/squaddieinterface"
	"github.com/chadius/terosbattleserver/usecase/repositories"
)

// SquaddieCanSwitchClassStrategy is the shape of all classes that can determine if a squaddie can switch their class.
type SquaddieCanSwitchClassStrategy interface {
	SquaddieCanSwitchToClass(squaddieToTest squaddieinterface.Interface, testingClassID string, repositories *repositories.RepositoryCollection) bool
}

// LevelsConsumedChecker looks at the Squaddie's consumed levels to determine if they can switch.
type LevelsConsumedChecker struct{}

// SquaddieCanSwitchToClass returns true if the squaddie can use the class with the given id.
func (l *LevelsConsumedChecker) SquaddieCanSwitchToClass(squaddieToTest squaddieinterface.Interface, testingClassID string, repositories *repositories.RepositoryCollection) bool {
	classToTest, _ := repositories.ClassRepo.GetClassByID(testingClassID)

	if squaddieToTest.BaseClassID() == "" && classToTest.BaseClassRequired() != true {
		return true
	}
	if squaddieToTest.BaseClassID() == "" && classToTest.BaseClassRequired() == true {
		return false
	}

	if squaddieToTest.CurrentClassID() == testingClassID {
		return false
	}

	if squaddieHasEnoughLevelsInClassToSwitch(squaddieToTest, squaddieToTest.CurrentClassID(), repositories) == false {
		return false
	}

	testingClassCompleted := areAllLevelsInClassTaken(squaddieToTest, testingClassID, repositories)
	if testingClassCompleted == true {
		return false
	}
	return true
}

func squaddieHasEnoughLevelsInClassToSwitch(squaddieToTest squaddieinterface.Interface, classID string, repositories *repositories.RepositoryCollection) bool {
	levelsInClass, _ := repositories.LevelRepo.GetLevelUpBenefitsByClassID(classID)
	levelsSquaddieConsumedInThisClass := countLevelsInClassTaken(squaddieToTest, classID)
	return levelsSquaddieConsumedInThisClass >= 10 || levelsSquaddieConsumedInThisClass >= len(levelsInClass)
}

func areAllLevelsInClassTaken(squaddieToTest squaddieinterface.Interface, classID string, repositories *repositories.RepositoryCollection) bool {
	levelsInClass, _ := repositories.LevelRepo.GetLevelUpBenefitsByClassID(classID)
	levelsSquaddieConsumedInThisClass := countLevelsInClassTaken(squaddieToTest, classID)
	return levelsSquaddieConsumedInThisClass >= len(levelsInClass)
}

func countLevelsInClassTaken(squaddieToTest squaddieinterface.Interface, classID string) int {
	squaddieLevelsConsumedInClasses := squaddieToTest.GetLevelCountsByClass()
	levelsSquaddieConsumedInThisClass, squaddieConsumedAnyLevelsInClass := squaddieLevelsConsumedInClasses[classID]

	if squaddieConsumedAnyLevelsInClass != true {
		return 0
	}

	return levelsSquaddieConsumedInThisClass
}
