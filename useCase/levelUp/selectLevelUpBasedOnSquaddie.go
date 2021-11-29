package levelup

import (
	"github.com/chadius/terosbattleserver/entity/levelupbenefit"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/utility"
)

// SelectLevelUpBasedOnSquaddieStrategy describes how to select the Level Up Benefit needed to improve a squaddie.
type SelectLevelUpBasedOnSquaddieStrategy interface {
	GetSquaddieClassLevels(
		squaddieToInspect *squaddie.Squaddie,
		repos *repositories.RepositoryCollection,
	) map[string]int

	ImproveSquaddieBasedOnLevel(
		squaddieToLevelUp *squaddie.Squaddie,
		bigLevelID string,
		repos *repositories.RepositoryCollection,
	) error
}

// SelectLevelUpBasedOnSquaddieBigLevelsOnEvenLevels will select a random small level every level and a selected big level at every even level.
type SelectLevelUpBasedOnSquaddieBigLevelsOnEvenLevels struct{}

// GetSquaddieClassLevels counts the levels for each class.
func (s *SelectLevelUpBasedOnSquaddieBigLevelsOnEvenLevels) GetSquaddieClassLevels(
	squaddieToInspect *squaddie.Squaddie,
	repos *repositories.RepositoryCollection,
) map[string]int {
	levels := map[string]int{}
	for classID, progress := range *squaddieToInspect.ClassLevelsConsumed() {
		levelsInClass, _ := repos.LevelRepo.GetLevelUpBenefitsByClassID(classID)

		smallLevelCount := progress.AccumulateLevelsConsumed(func(consumedLevelID string) int {
			return levelupbenefit.CountLevelUpBenefits(levelsInClass, func(benefit *levelupbenefit.LevelUpBenefit) bool {
				return benefit.ID() == consumedLevelID && benefit.LevelUpBenefitType() == levelupbenefit.Small
			})
		})
		levels[classID] = smallLevelCount
	}
	return levels
}

// ImproveSquaddieBasedOnLevel selects the levels the squaddie should get and then applies them.
func (s *SelectLevelUpBasedOnSquaddieBigLevelsOnEvenLevels) ImproveSquaddieBasedOnLevel(
	squaddieToLevelUp *squaddie.Squaddie,
	bigLevelID string,
	repos *repositories.RepositoryCollection,
) error {
	classToUse, err := repos.ClassRepo.GetClassByID(squaddieToLevelUp.CurrentClassID())
	if err != nil {
		return err
	}

	levelsFromClass, err := repos.LevelRepo.GetLevelUpBenefitsForClassByType(classToUse.ID())
	if err != nil {
		return err
	}

	squaddieLevels := s.GetSquaddieClassLevels(squaddieToLevelUp, repos)

	levelUpStrategy := ImproveSquaddieClass{}

	bigLevelToConsume := s.selectBigLevelUpForSquaddie(squaddieToLevelUp, bigLevelID, squaddieLevels, classToUse, levelsFromClass)
	if bigLevelToConsume != nil {
		levelUpStrategy.ImproveSquaddie(bigLevelToConsume, squaddieToLevelUp)
	}

	smallLevelToConsume := s.selectSmallLevelUpForSquaddie(squaddieToLevelUp, levelsFromClass)
	if smallLevelToConsume != nil {
		levelUpStrategy.ImproveSquaddie(smallLevelToConsume, squaddieToLevelUp)
	}
	return nil
}

// selectBigLevelUpForSquaddie chooses a Big LevelUpBenefit for the squaddie
//    and returns a pointer to it.
//    bigLevelSelectedID specifies the level to use,
//    but the classToUse can override this choice with an initial big level.
//    If the squaddie does not qualify for a Big level up, this will return nil.
//    If there are no big levels, return nil
func (s *SelectLevelUpBasedOnSquaddieBigLevelsOnEvenLevels) selectBigLevelUpForSquaddie(
	squaddieToLevelUp *squaddie.Squaddie,
	bigLevelSelectedID string,
	squaddieLevels map[string]int,
	classToUse *squaddieclass.Class,
	levelsFromClass map[levelupbenefit.Size][]*levelupbenefit.LevelUpBenefit,
) *levelupbenefit.LevelUpBenefit {

	squaddieClassIsEven := squaddieLevels[classToUse.ID()]%2 == 0
	if !squaddieClassIsEven {
		return nil
	}

	bigLevelIDToRetrieve := bigLevelSelectedID
	if classToUse.InitialBigLevelID() != "" &&
		squaddieToLevelUp.IsClassLevelAlreadyUsed(classToUse.InitialBigLevelID()) == false {
		bigLevelIDToRetrieve = classToUse.InitialBigLevelID()
	}

	bigLevelCandidates := levelupbenefit.FilterLevelUpBenefits(levelsFromClass[levelupbenefit.Big], func(level *levelupbenefit.LevelUpBenefit) bool {
		return level.ID() == bigLevelIDToRetrieve
	})

	if len(bigLevelCandidates) == 0 {
		return nil
	}

	return bigLevelCandidates[0]
}

// selectSmallLevelUpForSquaddie chooses a Small LevelUpBenefit for the squaddie
//    and returns a pointer to it. It is selected randomly.
//    If there are no small levels to choose from, return nil
func (s *SelectLevelUpBasedOnSquaddieBigLevelsOnEvenLevels) selectSmallLevelUpForSquaddie(
	squaddieToLevelUp *squaddie.Squaddie,
	levelsFromClass map[levelupbenefit.Size][]*levelupbenefit.LevelUpBenefit,
) *levelupbenefit.LevelUpBenefit {
	smallLevelsToChooseFrom := levelupbenefit.FilterLevelUpBenefits(levelsFromClass[levelupbenefit.Small],
		func(level *levelupbenefit.LevelUpBenefit) bool {
			if squaddieToLevelUp.IsClassLevelAlreadyUsed(level.ID()) {
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
