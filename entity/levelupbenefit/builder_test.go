package levelupbenefit_test

import (
	"github.com/chadius/terosbattleserver/entity/levelupbenefit"
	. "gopkg.in/check.v1"
)

type LevelUpBuilderSuite struct{}

var _ = Suite(&LevelUpBuilderSuite{})

func (l *LevelUpBuilderSuite) TestBuildDefaultLevelUpBenefit(checker *C) {
	levelUpBenefit, err := levelupbenefit.NewLevelUpBenefitBuilder().Build()

	checker.Assert(err, IsNil)
	checker.Assert(levelUpBenefit.ID(), Not(Equals), "")
	checker.Assert(levelUpBenefit.ClassID(), Not(Equals), "")
	checker.Assert(levelUpBenefit.LevelUpBenefitType(), Equals, levelupbenefit.Small)
}

func (l *LevelUpBuilderSuite) TestBuildWithIdentification(checker *C) {
	levelUpBenefit, err := levelupbenefit.NewLevelUpBenefitBuilder().
		WithID("levelID").
		WithClassID("classID").
		BigLevel().
		Build()

	checker.Assert(err, IsNil)
	checker.Assert(levelUpBenefit.ID(), Equals, "levelID")
	checker.Assert(levelUpBenefit.ClassID(), Equals, "classID")
	checker.Assert(levelUpBenefit.LevelUpBenefitType(), Equals, levelupbenefit.Big)
}

func (l *LevelUpBuilderSuite) TestBuildWithDefense(checker *C) {
	levelUpBenefit, err := levelupbenefit.NewLevelUpBenefitBuilder().
		HitPoints(5).
		Deflect(4).
		Dodge(3).
		Barrier(2).
		Armor(1).
		Build()

	checker.Assert(err, IsNil)
	checker.Assert(levelUpBenefit.MaxHitPoints(), Equals, 5)
	checker.Assert(levelUpBenefit.Deflect(), Equals, 4)
	checker.Assert(levelUpBenefit.Dodge(), Equals, 3)
	checker.Assert(levelUpBenefit.MaxBarrier(), Equals, 2)
	checker.Assert(levelUpBenefit.Armor(), Equals, 1)
}

func (l *LevelUpBuilderSuite) TestBuildWithOffense(checker *C) {
	levelUpBenefit, err := levelupbenefit.NewLevelUpBenefitBuilder().
		Aim(2).
		Strength(3).
		Mind(5).
		Build()

	checker.Assert(err, IsNil)
	checker.Assert(levelUpBenefit.Aim(), Equals, 2)
	checker.Assert(levelUpBenefit.Strength(), Equals, 3)
	checker.Assert(levelUpBenefit.Mind(), Equals, 5)
}
