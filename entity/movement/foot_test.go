package movement_test

import (
	"github.com/chadius/terosbattleserver/entity/movement"
	. "gopkg.in/check.v1"
)

type FootMovementSuite struct{}

var _ = Suite(&FootMovementSuite{})

func (suite *FootMovementSuite) TestFootIsNotGreaterThanFoot(checker *C) {
	foot := movement.NewMovementLogic("foot")
	foot2 := movement.NewMovementLogic("foot")
	checker.Assert(foot.GreaterThan(foot2), Equals, false)
}

func (suite *FootMovementSuite) TestFootIsNotGreaterThanLight(checker *C) {
	foot := movement.NewMovementLogic("foot")
	light := movement.NewMovementLogic("light")
	checker.Assert(foot.GreaterThan(light), Equals, false)
}

func (suite *FootMovementSuite) TestFootIsNotGreaterThanFly(checker *C) {
	foot := movement.NewMovementLogic("foot")
	fly := movement.NewMovementLogic("fly")
	checker.Assert(foot.GreaterThan(fly), Equals, false)
}

func (suite *FootMovementSuite) TestFootIsNotGreaterThanTeleport(checker *C) {
	foot := movement.NewMovementLogic("foot")
	teleport := movement.NewMovementLogic("teleport")
	checker.Assert(foot.GreaterThan(teleport), Equals, false)
}
