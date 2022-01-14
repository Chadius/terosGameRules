package healing_test

import (
	"github.com/chadius/terosgamerules/entity/healing"
	"github.com/chadius/terosgamerules/entity/squaddie"
	. "gopkg.in/check.v1"
)

type ZeroHealSuite struct{}

var _ = Suite(&ZeroHealSuite{})

func (suite *ZeroHealSuite) TestZeroHealAppliesUserMindStat(checker *C) {
	healer := squaddie.NewSquaddieBuilder().Mind(4).Build()
	target := squaddie.NewSquaddieBuilder().HitPoints(10).Build()
	target.ReduceHitPoints(9)
	ZeroHeal := &healing.ZeroMindBonus{}

	healingAmount := ZeroHeal.CalculateExpectedHeal(healer, 1, target)

	checker.Assert(healingAmount, Equals, 1)
}

func (suite *ZeroHealSuite) TestZeroHealCapsAtMaxHP(checker *C) {
	healer := squaddie.NewSquaddieBuilder().Mind(4).Build()
	target := squaddie.NewSquaddieBuilder().HitPoints(10).Build()
	ZeroHeal := &healing.ZeroMindBonus{}

	healingAmount := ZeroHeal.CalculateExpectedHeal(healer, 10, target)

	checker.Assert(healingAmount, Equals, 0)
}
