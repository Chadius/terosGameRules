package healing_test

import (
	"github.com/chadius/terosgamerules/entity/healing"
	"github.com/chadius/terosgamerules/entity/squaddie"
	. "gopkg.in/check.v1"
)

type HalfHealSuite struct{}

var _ = Suite(&HalfHealSuite{})

func (suite *HalfHealSuite) TestHalfHealAppliesUserMindStat(checker *C) {
	healer := squaddie.NewSquaddieBuilder().Mind(4).Build()
	target := squaddie.NewSquaddieBuilder().HitPoints(10).Build()
	target.ReduceHitPoints(9)
	HalfHeal := &healing.HalfMindBonus{}

	healingAmount := HalfHeal.CalculateExpectedHeal(healer, 1, target)

	checker.Assert(healingAmount, Equals, 3)
}

func (suite *HalfHealSuite) TestHalfHealCapsAtMaxHP(checker *C) {
	healer := squaddie.NewSquaddieBuilder().Mind(4).Build()
	target := squaddie.NewSquaddieBuilder().HitPoints(10).Build()
	target.ReduceHitPoints(1)
	HalfHeal := &healing.HalfMindBonus{}

	healingAmount := HalfHeal.CalculateExpectedHeal(healer, 1, target)

	checker.Assert(healingAmount, Equals, 1)
}
