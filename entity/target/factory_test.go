package target_test

import (
	"github.com/chadius/terosgamerules/entity/target"
	. "gopkg.in/check.v1"
	"reflect"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type FactorySuite struct{}

var _ = Suite(&FactorySuite{})

func (suite *FactorySuite) TestFriendTargeting(checker *C) {
	targetingLogic := target.NewTargetingLogic("friend")
	checker.Assert(reflect.TypeOf(targetingLogic).String(), Equals, "*target.Friend")
}

func (suite *FactorySuite) TestFoeTargeting(checker *C) {
	targetingLogic := target.NewTargetingLogic("foe")
	checker.Assert(reflect.TypeOf(targetingLogic).String(), Equals, "*target.Foe")
}

func (suite *FactorySuite) TestWhenUnknownKeyword_ThenFactoryReturnsSelf(checker *C) {
	targetingLogic := target.NewTargetingLogic("kwyjibo")
	checker.Assert(reflect.TypeOf(targetingLogic).String(), Equals, "*target.Self")
}
