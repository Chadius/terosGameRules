package levelupbenefit_test

import (
	"github.com/chadius/terosbattleserver/entity/levelupbenefit"
	. "gopkg.in/check.v1"
)

type LevelUpBenefitSuite struct{}

var _ = Suite(&LevelUpBenefitSuite{})

func (s *LevelUpBenefitSuite) TestFiltersAList(checker *C) {
	listToTest := []*levelupbenefit.LevelUpBenefit{
		{
			Identification: levelupbenefit.NewIdentification("level0", "class0", levelupbenefit.Small),
			Offense: &levelupbenefit.Offense{
				Aim: 1,
			},
		},
		{
			Identification: levelupbenefit.NewIdentification("level1", "class0", levelupbenefit.Small),
			Defense: &levelupbenefit.Defense{
				MaxHitPoints: 1,
			},
		},
		{
			Identification: levelupbenefit.NewIdentification("level2", "class0", levelupbenefit.Big),
			Offense: &levelupbenefit.Offense{
				Aim: 1,
			},
		},
	}

	noLevelsFound := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
		return false
	})
	checker.Assert(noLevelsFound, HasLen, 0)

	allLevelsFound := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
		return true
	})
	checker.Assert(allLevelsFound, HasLen, 3)

	onlySmallLevels := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
		return benefit.LevelUpBenefitType() == levelupbenefit.Small
	})
	checker.Assert(onlySmallLevels, HasLen, 2)

	onlyBigLevels := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
		return benefit.LevelUpBenefitType() == levelupbenefit.Big
	})
	checker.Assert(onlyBigLevels, HasLen, 1)

	increasesAimLevels := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
		return benefit.Offense != nil && benefit.Offense.Aim > 0
	})
	checker.Assert(increasesAimLevels, HasLen, 2)
}
