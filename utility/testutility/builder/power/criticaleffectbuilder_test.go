package power_test

import (
	"github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	. "gopkg.in/check.v1"
)

type CriticalEffectBuilder struct{}

var _ = Suite(&CriticalEffectBuilder{})

func (suite *CriticalEffectBuilder) TestBuildCriticalEffectDamage(checker *C) {
	criticalDamageEffect := power.CriticalEffectBuilder().DealsDamage(8).Build()
	checker.Assert(8, Equals, criticalDamageEffect.ExtraCriticalHitDamage())
}

func (suite *CriticalEffectBuilder) TestBuildCriticalEffectThresholdBonus(checker *C) {
	criticalDamageEffect := power.CriticalEffectBuilder().CriticalHitThresholdBonus(-2).Build()
	checker.Assert(-2, Equals, criticalDamageEffect.CriticalHitThresholdBonus())
}
