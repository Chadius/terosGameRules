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
