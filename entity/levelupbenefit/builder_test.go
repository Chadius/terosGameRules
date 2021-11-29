package levelupbenefit_test

import (
	"github.com/chadius/terosbattleserver/entity/levelupbenefit"
	"github.com/chadius/terosbattleserver/entity/squaddie"
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

func (l *LevelUpBuilderSuite) TestBuildWithMovement(checker *C) {
	levelUpBenefit, err := levelupbenefit.NewLevelUpBenefitBuilder().
		MovementDistance(2).
		MovementType(squaddie.Fly).
		CanHitAndRun().
		Build()

	checker.Assert(err, IsNil)
	checker.Assert(levelUpBenefit.MovementDistance(), Equals, 2)
	checker.Assert(levelUpBenefit.MovementType(), Equals, squaddie.Fly)
	checker.Assert(levelUpBenefit.CanHitAndRun(), Equals, true)
}

func (l *LevelUpBuilderSuite) TestBuildWithPowerChanges(checker *C) {
	levelUpBenefit, err := levelupbenefit.NewLevelUpBenefitBuilder().
		GainPower("spearLevel1", "Gold spear").
		GainPower("healingKitLevel0", "Healing kit").
		LosePower("spearLevel0").
		Build()

	checker.Assert(err, IsNil)
	checker.Assert(levelUpBenefit.PowersGained(), HasLen, 2)
	checker.Assert(levelUpBenefit.PowersGained()[0].Name, Equals, "Gold spear")
	checker.Assert(levelUpBenefit.PowersGained()[0].PowerID, Equals, "spearLevel1")
	checker.Assert(levelUpBenefit.PowersGained()[1].Name, Equals, "Healing kit")
	checker.Assert(levelUpBenefit.PowersGained()[1].PowerID, Equals, "healingKitLevel0")

	checker.Assert(levelUpBenefit.PowersLost(), HasLen, 1)
	checker.Assert(levelUpBenefit.PowersLost()[0].PowerID, Equals, "spearLevel0")
}
