package squaddie_test

import (
	"github.com/chadius/terosbattleserver/entity/squaddie"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MovementBuilder struct{}

var _ = Suite(&MovementBuilder{})

func (suite *MovementBuilder) TestBuildWithDistance(checker *C) {
	movement := squaddie.MovementBuilder().Distance(3).Build()
	checker.Assert(3, Equals, movement.MovementDistance())
}

func (suite *MovementBuilder) TestBuildMovementCanHitAndRun(checker *C) {
	movement := squaddie.MovementBuilder().CanHitAndRun().Build()
	checker.Assert(true, Equals, movement.CanHitAndRun())
}

func (suite *MovementBuilder) TestChangeMovementFoot(checker *C) {
	movement := squaddie.MovementBuilder().Foot().Build()
	checker.Assert(squaddie.Foot, Equals, movement.MovementType())
}

func (suite *MovementBuilder) TestChangeMovementLight(checker *C) {
	movement := squaddie.MovementBuilder().Light().Build()
	checker.Assert(squaddie.Light, Equals, movement.MovementType())
}

func (suite *MovementBuilder) TestChangeMovementFly(checker *C) {
	movement := squaddie.MovementBuilder().Fly().Build()
	checker.Assert(squaddie.Fly, Equals, movement.MovementType())
}

func (suite *MovementBuilder) TestChangeMovementTeleport(checker *C) {
	movement := squaddie.MovementBuilder().Teleport().Build()
	checker.Assert(squaddie.Teleport, Equals, movement.MovementType())
}
