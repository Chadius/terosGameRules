package levelup

import (
	"github.com/cserrant/terosbattleserver/entity/levelupbenefit"
	"github.com/cserrant/terosbattleserver/entity/squaddie"
	"github.com/cserrant/terosbattleserver/entity/squaddieclass"
	"github.com/cserrant/terosbattleserver/usecase/repositories"
	"github.com/cserrant/terosbattleserver/utility"
)

// GetSquaddieClassLevels returns a mapping of the squaddie's class to the number of times they leveled up.
func GetSquaddieClassLevels(
	squaddieToInspect *squaddie.Squaddie,
	repos *repositories.RepositoryCollection,
) map[string]int {
	levels := map[string]int{}
	for classID, progress := range squaddieToInspect.ClassProgress.ClassLevelsConsumed {
		levelsInClass, _ := repos.LevelRepo.GetLevelUpBenefitsByClassID(classID)

		smallLevelCount := progress.AccumulateLevelsConsumed(func(consumedLevelID string) int {
			return levelupbenefit.CountLevelUpBenefits(levelsInClass, func(benefit *levelupbenefit.LevelUpBenefit) bool {
				return benefit.Identification.ID == consumedLevelID && benefit.Identification.LevelUpBenefitType == levelupbenefit.Small
			})
		})
		levels[classID] = smallLevelCount
	}
	return levels
}

//ImproveSquaddieBasedOnLevel uses game logic to determine how to level up the squaddie.
//  Player gets to specify the big level to consume if the squaddie qualifies.
func ImproveSquaddieBasedOnLevel(
	squaddieToLevelUp *squaddie.Squaddie,
	bigLevelID string,
	repos *repositories.RepositoryCollection,
) error {
	classToUse, err := repos.ClassRepo.GetClassByID(squaddieToLevelUp.ClassProgress.CurrentClass)
	if err != nil {
		return err
	}

	levelsFromClass, err := repos.LevelRepo.GetLevelUpBenefitsForClassByType(classToUse.ID)
	if err != nil {
		return err
	}

	squaddieLevels := GetSquaddieClassLevels(squaddieToLevelUp, repos)

	bigLevelToConsume := selectBigLevelUpForSquaddie(squaddieToLevelUp, bigLevelID, squaddieLevels, classToUse, levelsFromClass)
	if bigLevelToConsume != nil {
		ImproveSquaddie(bigLevelToConsume, squaddieToLevelUp, repos)
	}

	smallLevelToConsume := selectSmallLevelUpForSquaddie(squaddieToLevelUp, levelsFromClass)
	if smallLevelToConsume != nil {
		ImproveSquaddie(smallLevelToConsume, squaddieToLevelUp, repos)
	}
	return nil
}

// selectBigLevelUpForSquaddie chooses a Big LevelUpBenefit for the squaddie
//    and returns a pointer to it.
//    bigLevelSelectedID specifies the level to use,
//    but the classToUse can override this choice with an initial big level.
//    If the squaddie does not qualify for a Big level up, this will return nil.
//    If there are no big levels, return nil
func selectBigLevelUpForSquaddie(
	squaddieToLevelUp *squaddie.Squaddie,
	bigLevelSelectedID string,
	squaddieLevels map[string]int,
	classToUse *squaddieclass.Class,
	levelsFromClass map[levelupbenefit.Size][]*levelupbenefit.LevelUpBenefit,
) *levelupbenefit.LevelUpBenefit {

	squaddieClassIsEven := squaddieLevels[classToUse.ID]%2 == 0
	if !squaddieClassIsEven {
		return nil
	}

	bigLevelIDToRetrieve := bigLevelSelectedID
	if classToUse.InitialBigLevelID != "" &&
		squaddieToLevelUp.ClassProgress.ClassLevelsConsumed[classToUse.ID].IsLevelAlreadyConsumed(classToUse.InitialBigLevelID) == false {
		bigLevelIDToRetrieve = classToUse.InitialBigLevelID
	}

	bigLevelCandidates := levelupbenefit.FilterLevelUpBenefits(levelsFromClass[levelupbenefit.Big], func(level *levelupbenefit.LevelUpBenefit) bool {
		return level.Identification.ID == bigLevelIDToRetrieve
	})

	if len(bigLevelCandidates) == 0 {
		return nil
	}

	return bigLevelCandidates[0]
}

// selectSmallLevelUpForSquaddie chooses a Small LevelUpBenefit for the squaddie
//    and returns a pointer to it. It is selected randomly.
//    If there are no small levels to choose from, return nil
func selectSmallLevelUpForSquaddie(
	squaddieToLevelUp *squaddie.Squaddie,
	levelsFromClass map[levelupbenefit.Size][]*levelupbenefit.LevelUpBenefit,
) *levelupbenefit.LevelUpBenefit {
	smallLevelsToChooseFrom := levelupbenefit.FilterLevelUpBenefits(levelsFromClass[levelupbenefit.Small],
		func(level *levelupbenefit.LevelUpBenefit) bool {
			if squaddieToLevelUp.ClassProgress.ClassLevelsConsumed[squaddieToLevelUp.ClassProgress.CurrentClass].IsLevelAlreadyConsumed(level.Identification.ID) {
				return false
			}
			return true
		},
	)

	if len(smallLevelsToChooseFrom) > 0 {
		return smallLevelsToChooseFrom[utility.RandomInt(len(smallLevelsToChooseFrom))]
	}
	return nil
}
