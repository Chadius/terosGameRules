package healing_test

import (
	"github.com/chadius/terosgamerules/entity/healing"
	"github.com/chadius/terosgamerules/entity/squaddie"
	. "gopkg.in/check.v1"
)

type NoHealingTestSuite struct{}

var _ = Suite(&NoHealingTestSuite{})

func (suite *NoHealingTestSuite) TestNoHealingAlwaysHealsZero(checker *C) {
	healer := squaddie.NewSquaddieBuilder().Mind(4).Build()
	target := squaddie.NewSquaddieBuilder().HitPoints(10).Build()
	target.ReduceHitPoints(9)
	noHeal := &healing.NoHealing{}

	healingAmount := noHeal.CalculateExpectedHeal(healer, 1, target)

	checker.Assert(healingAmount, Equals, 0)
}
