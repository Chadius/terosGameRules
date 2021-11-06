package power_test

import (
	powerBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	. "gopkg.in/check.v1"
)

type AttackingEffectCounterAttackPenaltyTest struct{}

var _ = Suite(&AttackingEffectCounterAttackPenaltyTest{})

func (suite *AttackingEffectCounterAttackPenaltyTest) SetUpTest(checker *C) {}

func (suite *AttackingEffectCounterAttackPenaltyTest) TestDefaultPenalty(checker *C) {
	counterAttackingPower := powerBuilder.Builder().DealsDamage(1).CanCounterAttack().Build()
	counterAttackPenalty, err := counterAttackingPower.CounterAttackPenalty()
	checker.Assert(err, IsNil)
	checker.Assert(counterAttackPenalty, Equals, -2)
}

func (suite *AttackingEffectCounterAttackPenaltyTest) TestRaisesErrorIfPowerCannotCounterAttack(checker *C) {
	cannotCounterWithThisPower := powerBuilder.Builder().DealsDamage(1).Build()
	_, err := cannotCounterWithThisPower.CounterAttackPenalty()
	checker.Assert(err, ErrorMatches, "power cannot counter, cannot calculate penalty")
}

func (suite *AttackingEffectCounterAttackPenaltyTest) TestAppliesPenaltyReduction(checker *C) {
	counterAttackingPower := powerBuilder.Builder().DealsDamage(1).CanCounterAttack().CounterAttackPenaltyReduction(2).Build()
	counterAttackPenalty, err := counterAttackingPower.CounterAttackPenalty()
	checker.Assert(err, IsNil)
	checker.Assert(counterAttackPenalty, Equals, 0)
}
