package healing_test

import (
	"github.com/chadius/terosbattleserver/entity/healing"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	. "gopkg.in/check.v1"
)

type ZeroHealSuite struct{}

var _ = Suite(&ZeroHealSuite{})

func (suite *ZeroHealSuite) TestZeroHealAppliesUserMindStat(checker *C) {
	healer := squaddie.NewSquaddieBuilder().Mind(4).Build()
	healingPower := power.NewPowerBuilder().HitPointsHealed(1).HealingAdjustmentBasedOnUserMindZero().Build()
	target := squaddie.NewSquaddieBuilder().HitPoints(10).Build()
	target.ReduceHitPoints(9)
	ZeroHeal := &healing.ZeroMindBonus{}

	healingAmount := ZeroHeal.CalculateExpectedHeal(healer, healingPower, target)

	checker.Assert(healingAmount, Equals, 1)
}

func (suite *ZeroHealSuite) TestZeroHealCapsAtMaxHP(checker *C) {
	healer := squaddie.NewSquaddieBuilder().Mind(4).Build()
	healingPower := power.NewPowerBuilder().HitPointsHealed(10).HealingAdjustmentBasedOnUserMindZero().Build()
	target := squaddie.NewSquaddieBuilder().HitPoints(10).Build()
	ZeroHeal := &healing.ZeroMindBonus{}

	healingAmount := ZeroHeal.CalculateExpectedHeal(healer, healingPower, target)

	checker.Assert(healingAmount, Equals, 0)
}
