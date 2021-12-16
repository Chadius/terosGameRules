package movement_test

import (
	"github.com/chadius/terosbattleserver/entity/movement"
	. "gopkg.in/check.v1"
)

type FlyMovementSuite struct{}

var _ = Suite(&FlyMovementSuite{})

func (suite *FlyMovementSuite) TestFlyIsGreaterThanFoot(checker *C) {
	fly := movement.NewMovementLogic("fly")
	foot := movement.NewMovementLogic("foot")
	checker.Assert(fly.GreaterThan(foot), Equals, true)
}

func (suite *FlyMovementSuite) TestFlyIsGreaterThanLight(checker *C) {
	fly := movement.NewMovementLogic("fly")
	light := movement.NewMovementLogic("light")
	checker.Assert(fly.GreaterThan(light), Equals, true)
}

func (suite *FlyMovementSuite) TestFlyIsNotGreaterThanFly(checker *C) {
	fly := movement.NewMovementLogic("fly")
	fly2 := movement.NewMovementLogic("fly")
	checker.Assert(fly.GreaterThan(fly2), Equals, false)
}

func (suite *FlyMovementSuite) TestFlyIsNotGreaterThanTeleport(checker *C) {
	fly := movement.NewMovementLogic("fly")
	teleport := movement.NewMovementLogic("teleport")
	checker.Assert(fly.GreaterThan(teleport), Equals, false)
}
