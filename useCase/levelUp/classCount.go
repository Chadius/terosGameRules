package levelUp

import (
	"github.com/cserrant/terosBattleServer/entity/levelUpBenefit"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/entity/squaddieClass"
	"github.com/cserrant/terosBattleServer/utility"
)

// GetSquaddieClassLevels returns a mapping of the squaddie's class to the number of times they leveled up.
func GetSquaddieClassLevels(squaddieToInspect *squaddie.Squaddie, levelRepo *levelUpBenefit.Repository) map[string]int {
	levels := map[string]int{}
	for classID, progress := range squaddieToInspect.ClassLevelsConsumed {
		levelsInClass, _ := levelRepo.GetLevelUpBenefitsByClassID(classID)

		smallLevelCount := 0
		for _, levelID := range progress.LevelsConsumed {
			for _, levelFromRepo := range levelsInClass {
				if levelFromRepo.ID == levelID && levelFromRepo.LevelUpBenefitType == levelUpBenefit.Small {
					smallLevelCount = smallLevelCount + 1
				}
			}
		}
		levels[classID] = smallLevelCount
	}
	return levels
}

//LevelUpSquaddieBasedOnSquaddieLevel uses game logic to determine how to level up the squaddie.
//  Player gets to specify the big level IF the squaddie's level is even.
func LevelUpSquaddieBasedOnSquaddieLevel(squaddieToLevelUp *squaddie.Squaddie, bigLevelID string, levelRepo *levelUpBenefit.Repository, classRepo *squaddieClass.Repository, powerRepo *power.Repository) error {
	classToUse, err := classRepo.GetClassByID(squaddieToLevelUp.CurrentClass)
	if err != nil {
		return err
	}

	levelsFromClass, err := levelRepo.GetLevelUpBenefitsForClassByType(classToUse.ID)
	if err != nil {
		return err
	}

	squaddieLevels := GetSquaddieClassLevels(squaddieToLevelUp, levelRepo)

	if squaddieLevels[classToUse.ID] % 2 == 0 {
		bigLevelCandidates := levelUpBenefit.FilterLevelUpBenefits(levelsFromClass[levelUpBenefit.Big], func(level *levelUpBenefit.LevelUpBenefit) bool {
			return level.ID == bigLevelID
		})
		if len(bigLevelCandidates) > 0 {
			LevelUpSquaddie(bigLevelCandidates[0], squaddieToLevelUp, powerRepo)
		}
	}

	smallLevelsToChooseFrom := levelUpBenefit.FilterLevelUpBenefits(levelsFromClass[levelUpBenefit.Small],
		func(level *levelUpBenefit.LevelUpBenefit) bool {
			if squaddieToLevelUp.ClassLevelsConsumed[squaddieToLevelUp.CurrentClass].IsLevelAlreadyConsumed(level.ID) {
				return false
			}
			return true
		},
	)

	if len(smallLevelsToChooseFrom)> 0 {
		smallLevelToConsume := smallLevelsToChooseFrom[utility.RandomInt(len(smallLevelsToChooseFrom))]
		LevelUpSquaddie(smallLevelToConsume, squaddieToLevelUp, powerRepo)
	}
	return nil
}