package healing_test

import (
	"github.com/chadius/terosgamerules/entity/healing"
	. "gopkg.in/check.v1"
	"reflect"
)

type FactorySuite struct{}

var _ = Suite(&FactorySuite{})

func (suite *FactorySuite) TestFactoryReturnsFullHeal(checker *C) {
	healingLogic := healing.NewHealingLogic("Full")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.FullMindBonus")

	healingLogic = healing.NewHealingLogic("full")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.FullMindBonus")

	healingLogic = healing.NewHealingLogic("*healing.FullMindBonus")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.FullMindBonus")

	healingLogic = healing.NewHealingLogic("healing.FullMindBonus")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.FullMindBonus")
}

func (suite *FactorySuite) TestFactoryReturnsHalfHeal(checker *C) {
	healingLogic := healing.NewHealingLogic("Half")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.HalfMindBonus")

	healingLogic = healing.NewHealingLogic("half")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.HalfMindBonus")

	healingLogic = healing.NewHealingLogic("*healing.HalfMindBonus")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.HalfMindBonus")

	healingLogic = healing.NewHealingLogic("healing.HalfMindBonus")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.HalfMindBonus")
}

func (suite *FactorySuite) TestFactoryReturnsZeroHeal(checker *C) {
	healingLogic := healing.NewHealingLogic("Zero")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.ZeroMindBonus")

	healingLogic = healing.NewHealingLogic("zero")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.ZeroMindBonus")

	healingLogic = healing.NewHealingLogic("*healing.ZeroMindBonus")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.ZeroMindBonus")

	healingLogic = healing.NewHealingLogic("healing.ZeroMindBonus")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.ZeroMindBonus")
}

func (suite *FactorySuite) TestFactoryReturnsNoHealing(checker *C) {
	healingLogic := healing.NewHealingLogic("NoHealing")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.NoHealing")

	healingLogic = healing.NewHealingLogic("nohealing")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.NoHealing")

	healingLogic = healing.NewHealingLogic("No healing")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.NoHealing")

	healingLogic = healing.NewHealingLogic("No Healing")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.NoHealing")

	healingLogic = healing.NewHealingLogic("*healing.NoHealing")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.NoHealing")

	healingLogic = healing.NewHealingLogic("healing.NoHealing")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.NoHealing")
}

func (suite *FactorySuite) TestWhenUnknownKeyword_ThenFactoryReturnsNoHealing(checker *C) {
	healingLogic := healing.NewHealingLogic("kwyjibo")
	checker.Assert(reflect.TypeOf(healingLogic).String(), Equals, "*healing.NoHealing")
}
