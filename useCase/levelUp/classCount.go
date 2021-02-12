package levelup

import (
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	"github.com/cserrant/terosBattleServer/utility"
)

// GetSquaddieClassLevels returns a mapping of the squaddie's class to the number of times they leveled up.
func GetSquaddieClassLevels(squaddieToInspect *squaddie.Squaddie, levelRepo *levelupbenefit.Repository) map[string]int {
	levels := map[string]int{}
	for classID, progress := range squaddieToInspect.ClassLevelsConsumed {
		levelsInClass, _ := levelRepo.GetLevelUpBenefitsByClassID(classID)

		smallLevelCount := progress.AccumulateLevelsConsumed(func(consumedLevelID string) int {
			return levelupbenefit.CountLevelUpBenefits(levelsInClass, func(benefit *levelupbenefit.LevelUpBenefit) bool {
				return benefit.ID == consumedLevelID && benefit.LevelUpBenefitType == levelupbenefit.Small
			})
		})
		levels[classID] = smallLevelCount
	}
	return levels
}

//ImproveSquaddieBasedOnLevel uses game logic to determine how to level up the squaddie.
//  Player gets to specify the big level IF the squaddie's level is even.
func ImproveSquaddieBasedOnLevel(squaddieToLevelUp *squaddie.Squaddie, bigLevelID string, levelRepo *levelupbenefit.Repository, classRepo *squaddieclass.Repository, powerRepo *power.Repository) error {
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
		bigLevelCandidates := levelupbenefit.FilterLevelUpBenefits(levelsFromClass[levelupbenefit.Big], func(level *levelupbenefit.LevelUpBenefit) bool {
			return level.ID == bigLevelID
		})
		if len(bigLevelCandidates) > 0 {
			ImproveSquaddie(bigLevelCandidates[0], squaddieToLevelUp, powerRepo)
		}
	}

	smallLevelsToChooseFrom := levelupbenefit.FilterLevelUpBenefits(levelsFromClass[levelupbenefit.Small],
		func(level *levelupbenefit.LevelUpBenefit) bool {
			if squaddieToLevelUp.ClassLevelsConsumed[squaddieToLevelUp.CurrentClass].IsLevelAlreadyConsumed(level.ID) {
				return false
			}
			return true
		},
	)

	if len(smallLevelsToChooseFrom)> 0 {
		smallLevelToConsume := smallLevelsToChooseFrom[utility.RandomInt(len(smallLevelsToChooseFrom))]
		ImproveSquaddie(smallLevelToConsume, squaddieToLevelUp, powerRepo)
	}
	return nil
}