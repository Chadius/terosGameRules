package levelup

import (
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
)

// SquaddieCanSwitchToClass returns true if the squaddie can use the class with the given ID.
func SquaddieCanSwitchToClass(squaddieToTest *squaddie.Squaddie, testingClassID string, classRepo *squaddieclass.Repository, levelRepo *levelupbenefit.Repository) bool {
	classToTest, _ := classRepo.GetClassByID(testingClassID)

	if squaddieToTest.BaseClassID == "" && classToTest.BaseClassRequired != true {
		return true
	}
	if squaddieToTest.BaseClassID == "" && classToTest.BaseClassRequired == true {
		return false
	}

	if squaddieToTest.CurrentClass == testingClassID {
		return false
	}

	if squaddieHasEnoughLevelsInClassToSwitch(squaddieToTest, squaddieToTest.CurrentClass, levelRepo) == false {
		return false
	}

	testingClassCompleted := areAllLevelsInClassTaken(squaddieToTest, testingClassID, levelRepo)
	if testingClassCompleted == true {
		return false
	}
	return true
}

func squaddieHasEnoughLevelsInClassToSwitch(squaddieToTest *squaddie.Squaddie, classID string, levelRepo *levelupbenefit.Repository) bool {
	levelsInClass, _ := levelRepo.GetLevelUpBenefitsByClassID(classID)
	levelsSquaddieConsumedInThisClass := countLevelsInClassTaken(squaddieToTest, classID)
	return levelsSquaddieConsumedInThisClass >= 10 || levelsSquaddieConsumedInThisClass >= len(levelsInClass)
}

func areAllLevelsInClassTaken(squaddieToTest *squaddie.Squaddie, classID string, levelRepo *levelupbenefit.Repository) bool {
	levelsInClass, _ := levelRepo.GetLevelUpBenefitsByClassID(classID)
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