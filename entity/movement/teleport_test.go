package movement_test

import (
	"github.com/chadius/terosgamerules/entity/movement"
	. "gopkg.in/check.v1"
)

type TeleportMovementSuite struct{}

var _ = Suite(&TeleportMovementSuite{})

func (suite *TeleportMovementSuite) TestTeleportIsGreaterThanFoot(checker *C) {
	teleport := movement.NewMovementLogic("teleport")
	foot := movement.NewMovementLogic("foot")
	checker.Assert(teleport.GreaterThan(foot), Equals, true)
}
func (suite *TeleportMovementSuite) TestTeleportIsGreaterThanLight(checker *C) {
	teleport := movement.NewMovementLogic("teleport")
	light := movement.NewMovementLogic("light")
	checker.Assert(teleport.GreaterThan(light), Equals, true)
}
func (suite *TeleportMovementSuite) TestTeleportIsGreaterThanFly(checker *C) {
	teleport := movement.NewMovementLogic("teleport")
	fly := movement.NewMovementLogic("fly")
	checker.Assert(teleport.GreaterThan(fly), Equals, true)
}

func (suite *TeleportMovementSuite) TestTeleportIsNotGreaterThanTeleport(checker *C) {
	teleport := movement.NewMovementLogic("teleport")
	teleport2 := movement.NewMovementLogic("teleport")
	checker.Assert(teleport.GreaterThan(teleport2), Equals, false)
}
