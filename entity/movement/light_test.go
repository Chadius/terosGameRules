package movement_test

import (
	"github.com/chadius/terosgamerules/entity/movement"
	. "gopkg.in/check.v1"
)

type LightMovementSuite struct{}

var _ = Suite(&LightMovementSuite{})

func (suite *LightMovementSuite) TestLightIsGreaterThanFoot(checker *C) {
	light := movement.NewMovementLogic("light")
	foot := movement.NewMovementLogic("foot")
	checker.Assert(light.GreaterThan(foot), Equals, true)
}

func (suite *LightMovementSuite) TestLightIsNotGreaterThanLight(checker *C) {
	light := movement.NewMovementLogic("light")
	light2 := movement.NewMovementLogic("light")
	checker.Assert(light.GreaterThan(light2), Equals, false)
}

func (suite *LightMovementSuite) TestLightIsNotGreaterThanFly(checker *C) {
	light := movement.NewMovementLogic("light")
	fly := movement.NewMovementLogic("fly")
	checker.Assert(light.GreaterThan(fly), Equals, false)
}

func (suite *LightMovementSuite) TestLightIsNotGreaterThanTeleport(checker *C) {
	light := movement.NewMovementLogic("light")
	teleport := movement.NewMovementLogic("teleport")
	checker.Assert(light.GreaterThan(teleport), Equals, false)
}
