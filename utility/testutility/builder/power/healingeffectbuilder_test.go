package power_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	powerBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	. "gopkg.in/check.v1"
)

type HealingEffectBuilder struct{}

var _ = Suite(&HealingEffectBuilder{})

func (suite *HealingEffectBuilder) TestHealingAdjustmentFull(checker *C) {
	bigHeals := powerBuilder.HealingEffectBuilder().HealingAdjustmentBasedOnUserMindFull().Build()
	checker.Assert(power.Full, Equals, bigHeals.HealingAdjustmentBasedOnUserMind())
}

func (suite *HealingEffectBuilder) TestHealingAdjustmentHalf(checker *C) {
	someHeals := powerBuilder.HealingEffectBuilder().HealingAdjustmentBasedOnUserMindHalf().Build()
	checker.Assert(power.Half, Equals, someHeals.HealingAdjustmentBasedOnUserMind())
}

func (suite *HealingEffectBuilder) TestHealingAdjustmentZero(checker *C) {
	someHeals := powerBuilder.HealingEffectBuilder().HealingAdjustmentBasedOnUserMindZero().Build()
	checker.Assert(power.Zero, Equals, someHeals.HealingAdjustmentBasedOnUserMind())
}

func (suite *HealingEffectBuilder) TestHitPointsHealed(checker *C) {
	bigHeals := powerBuilder.HealingEffectBuilder().HitPointsHealed(5).Build()
	checker.Assert(5, Equals, bigHeals.HitPointsHealed())
}
