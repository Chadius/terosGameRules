package levelup

import (
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/repositories"
)

// SquaddieCanSwitchToClass returns true if the squaddie can use the class with the given SquaddieID.
func SquaddieCanSwitchToClass(squaddieToTest *squaddie.Squaddie, testingClassID string, repositories *repositories.RepositoryCollection) bool {
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

func squaddieHasEnoughLevelsInClassToSwitch(squaddieToTest *squaddie.Squaddie, classID string, repositories *repositories.RepositoryCollection) bool {
	levelsInClass, _ := repositories.LevelRepo.GetLevelUpBenefitsByClassID(classID)
	levelsSquaddieConsumedInThisClass := countLevelsInClassTaken(squaddieToTest, classID)
	return levelsSquaddieConsumedInThisClass >= 10 || levelsSquaddieConsumedInThisClass >= len(levelsInClass)
}

func areAllLevelsInClassTaken(squaddieToTest *squaddie.Squaddie, classID string, repositories *repositories.RepositoryCollection) bool {
	levelsInClass, _ := repositories.LevelRepo.GetLevelUpBenefitsByClassID(classID)
	levelsSquaddieConsumedInThisClass := countLevelsInClassTaken(squaddieToTest, classID)
	return levelsSquaddieConsumedInThisClass >= len(levelsInClass)
}

func countLevelsInClassTaken(squaddieToTest *squaddie.Squaddie, classID string) int {
	squaddieLevelsConsumedInClasses := squaddieToTest.GetLevelCountsByClass()
	levelsSquaddieConsumedInThisClass, squaddieConsumedAnyLevelsInClass := squaddieLevelsConsumedInClasses[classID]

	if squaddieConsumedAnyLevelsInClass != true {
		return 0
	}

	return levelsSquaddieConsumedInThisClass
}
