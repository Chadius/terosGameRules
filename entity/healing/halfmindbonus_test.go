package healing_test

import (
	"github.com/chadius/terosbattleserver/entity/healing"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	. "gopkg.in/check.v1"
)

type HalfHealSuite struct{}

var _ = Suite(&HalfHealSuite{})

func (suite *HalfHealSuite) TestHalfHealAppliesUserMindStat(checker *C) {
	healer := squaddie.NewSquaddieBuilder().Mind(4).Build()
	healingPower := power.NewPowerBuilder().HitPointsHealed(1).HealingAdjustmentBasedOnUserMindHalf().Build()
	target := squaddie.NewSquaddieBuilder().HitPoints(10).Build()
	target.ReduceHitPoints(9)
	HalfHeal := &healing.HalfMindBonus{}

	healingAmount := HalfHeal.CalculateExpectedHeal(healer, healingPower, target)

	checker.Assert(healingAmount, Equals, 3)
}

func (suite *HalfHealSuite) TestHalfHealCapsAtMaxHP(checker *C) {
	healer := squaddie.NewSquaddieBuilder().Mind(4).Build()
	healingPower := power.NewPowerBuilder().HitPointsHealed(1).HealingAdjustmentBasedOnUserMindHalf().Build()
	target := squaddie.NewSquaddieBuilder().HitPoints(10).Build()
	target.ReduceHitPoints(1)
	HalfHeal := &healing.HalfMindBonus{}

	healingAmount := HalfHeal.CalculateExpectedHeal(healer, healingPower, target)

	checker.Assert(healingAmount, Equals, 1)
}
