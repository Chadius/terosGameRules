package healing_test

import (
	"github.com/chadius/terosgamerules/entity/healing"
	"github.com/chadius/terosgamerules/entity/squaddie"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type FullHealSuite struct{}

var _ = Suite(&FullHealSuite{})

func (suite *FullHealSuite) TestFullHealAppliesUserMindStat(checker *C) {
	healer := squaddie.NewSquaddieBuilder().Mind(4).Build()
	target := squaddie.NewSquaddieBuilder().HitPoints(10).Build()
	target.ReduceHitPoints(9)
	fullHeal := &healing.FullMindBonus{}

	healingAmount := fullHeal.CalculateExpectedHeal(healer, 1, target)

	checker.Assert(healingAmount, Equals, 5)
}

func (suite *FullHealSuite) TestFullHealCapsAtMaxHP(checker *C) {
	healer := squaddie.NewSquaddieBuilder().Mind(4).Build()
	target := squaddie.NewSquaddieBuilder().HitPoints(10).Build()
	target.ReduceHitPoints(1)
	FullHeal := &healing.FullMindBonus{}

	healingAmount := FullHeal.CalculateExpectedHeal(healer, 1, target)

	checker.Assert(healingAmount, Equals, 1)
}
