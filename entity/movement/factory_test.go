package movement_test

import (
	"github.com/chadius/terosgamerules/entity/movement"
	. "gopkg.in/check.v1"
	"reflect"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type FactorySuite struct{}

var _ = Suite(&FactorySuite{})

func (suite *FactorySuite) TestLightMovement(checker *C) {
	powerSourceLogic := movement.NewMovementLogic("light")
	checker.Assert(reflect.TypeOf(powerSourceLogic).String(), Equals, "*movement.Light")
}

func (suite *FactorySuite) TestFlyMovement(checker *C) {
	powerSourceLogic := movement.NewMovementLogic("fly")
	checker.Assert(reflect.TypeOf(powerSourceLogic).String(), Equals, "*movement.Fly")
}

func (suite *FactorySuite) TestTeleportMovement(checker *C) {
	powerSourceLogic := movement.NewMovementLogic("teleport")
	checker.Assert(reflect.TypeOf(powerSourceLogic).String(), Equals, "*movement.Teleport")
}

func (suite *FactorySuite) TestWhenUnknownKeyword_ThenFactoryReturnsFoot(checker *C) {
	powerSourceLogic := movement.NewMovementLogic("kwyjibo")
	checker.Assert(reflect.TypeOf(powerSourceLogic).String(), Equals, "*movement.Foot")
}
