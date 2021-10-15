package power_test

import (
	"github.com/chadius/terosbattleserver/utility/testutility/factory/power"
	. "gopkg.in/check.v1"
)

type CriticalEffectBuilder struct {}

var _ = Suite(&CriticalEffectBuilder{})

func (suite *CriticalEffectBuilder) TestBuildCriticalEffectDamage(checker *C) {
	criticalDamageEffect := power.CriticalEffectFactory().DealsDamage(8).Build()
	checker.Assert(8, Equals, criticalDamageEffect.Damage)
}

func (suite *CriticalEffectBuilder) TestBuildCriticalEffectThresholdBonus(checker *C) {
	criticalDamageEffect := power.CriticalEffectFactory().CriticalHitThresholdBonus(-2).Build()
	checker.Assert(-2, Equals, criticalDamageEffect.CriticalHitThresholdBonus)
}
