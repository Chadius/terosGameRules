package healing_test

import (
	"github.com/chadius/terosbattleserver/entity/healing"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	. "gopkg.in/check.v1"
)

type NoHealingTestSuite struct{}

var _ = Suite(&NoHealingTestSuite{})

func (suite *NoHealingTestSuite) TestNoHealingAlwaysHealsZero(checker *C) {
	healer := squaddie.NewSquaddieBuilder().Mind(4).Build()
	healingPower := power.NewPowerBuilder().HitPointsHealed(1).Build()
	target := squaddie.NewSquaddieBuilder().HitPoints(10).Build()
	target.ReduceHitPoints(9)
	noHeal := &healing.NoHealing{}

	healingAmount := noHeal.CalculateExpectedHeal(healer, healingPower, target)

	checker.Assert(healingAmount, Equals, 0)
}
