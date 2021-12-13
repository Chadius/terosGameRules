package power_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	. "gopkg.in/check.v1"
)

type HealingEffectBuilder struct{}

var _ = Suite(&HealingEffectBuilder{})

func (suite *HealingEffectBuilder) TestHitPointsHealed(checker *C) {
	bigHeals := power.HealingEffectBuilder().HitPointsHealed(5).Build()
	checker.Assert(5, Equals, bigHeals.HitPointsHealed())
}
