package levelupbenefit_test

import (
	"github.com/chadius/terosbattleserver/entity/levelupbenefit"
	. "gopkg.in/check.v1"
)

type LevelUpBenefitSuite struct{}

var _ = Suite(&LevelUpBenefitSuite{})

func (s *LevelUpBenefitSuite) TestFiltersAList(checker *C) {
	listToTest := []*levelupbenefit.LevelUpBenefit{}
	var newLevel *levelupbenefit.LevelUpBenefit
	newLevel, _ = levelupbenefit.NewLevelUpBenefitBuilder().
		WithID("level0").
		WithClassID("class0").
		Aim(1).
		Build()
	listToTest = append(listToTest, newLevel)

	newLevel, _ = levelupbenefit.NewLevelUpBenefitBuilder().
		WithID("level1").
		WithClassID("class0").
		HitPoints(1).
		Build()
	listToTest = append(listToTest, newLevel)

	newLevel, _ = levelupbenefit.NewLevelUpBenefitBuilder().
		WithID("level2").
		WithClassID("class0").
		Aim(1).
		BigLevel().
		Build()
	listToTest = append(listToTest, newLevel)

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
		return benefit.Aim() > 0
	})
	checker.Assert(increasesAimLevels, HasLen, 2)
}
