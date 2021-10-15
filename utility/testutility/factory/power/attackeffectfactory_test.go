package power_test

import (
	"github.com/chadius/terosbattleserver/utility/testutility/factory/power"
	. "gopkg.in/check.v1"
)


type AttackEffectBuilder struct {}

var _ = Suite(&AttackEffectBuilder{})

func (suite *AttackEffectBuilder) TestBuildAttackEffectToHitBonus(checker *C) {
	damageEffect := power.AttackEffectFactory().ToHitBonus(2).Build()
	checker.Assert(2, Equals, damageEffect.ToHitBonus)
}

func (suite *AttackEffectBuilder) TestBuildAttackEffectDamageBonus(checker *C) {
	damageEffect := power.AttackEffectFactory().DealsDamage(3).Build()
	checker.Assert(3, Equals, damageEffect.DamageBonus)
}

func (suite *AttackEffectBuilder) TestBuildAttackEffectExtraBarrierBurn(checker *C) {
	damageEffect := power.AttackEffectFactory().ExtraBarrierBurn(1).Build()
	checker.Assert(1, Equals, damageEffect.ExtraBarrierBurn)
}

func (suite *AttackEffectBuilder) TestBuildAttackEffectCounterAttackPenaltyReduction(checker *C) {
	damageEffect := power.AttackEffectFactory().CounterAttackPenaltyReduction(4).Build()
	checker.Assert(4, Equals, damageEffect.CounterAttackPenaltyReduction)
}

func (suite *AttackEffectBuilder) TestBuildAttackEffectCanBeEquipped(checker *C) {
	sword := power.AttackEffectFactory().CanBeEquipped().Build()
	checker.Assert(true, Equals, sword.CanBeEquipped)
}

func (suite *AttackEffectBuilder) TestBuildAttackEffectCannotBeEquipped(checker *C) {
	scroll := power.AttackEffectFactory().CanBeEquipped().CannotBeEquipped().Build()
	checker.Assert(false, Equals, scroll.CanBeEquipped)
}

func (suite *AttackEffectBuilder) TestBuildAttackEffectCanCounterAttack(checker *C) {
	sword := power.AttackEffectFactory().CanCounterAttack().Build()
	checker.Assert(true, Equals, sword.CanCounterAttack)
}


func (suite *AttackEffectBuilder) TestBuildCriticalEffectDamage(checker *C) {
	criticalDamageEffect := power.AttackEffectFactory().CriticalDealsDamage(8).Build()
	checker.Assert(8, Equals, criticalDamageEffect.CriticalEffect.Damage)
}

func (suite *AttackEffectBuilder) TestBuildCriticalEffectThresholdBonus(checker *C) {
	criticalDamageEffect := power.AttackEffectFactory().CriticalHitThresholdBonus(-2).Build()
	checker.Assert(-2, Equals, criticalDamageEffect.CriticalEffect.CriticalHitThresholdBonus)
}
