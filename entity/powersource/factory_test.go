package powersource_test

import (
	"github.com/chadius/terosgamerules/entity/powersource"
	. "gopkg.in/check.v1"
	"reflect"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type FactorySuite struct{}

var _ = Suite(&FactorySuite{})

func (suite *FactorySuite) TestSpellSource(checker *C) {
	powerSourceLogic := powersource.NewPowerSourceLogic("spell")
	checker.Assert(reflect.TypeOf(powerSourceLogic).String(), Equals, "*powersource.Spell")
}

func (suite *FactorySuite) TestWhenUnknownKeyword_ThenFactoryReturnsPhysical(checker *C) {
	powerSourceLogic := powersource.NewPowerSourceLogic("kwyjibo")
	checker.Assert(reflect.TypeOf(powerSourceLogic).String(), Equals, "*powersource.Physical")
}
