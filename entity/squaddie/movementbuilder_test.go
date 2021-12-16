package squaddie_test

import (
	"github.com/chadius/terosbattleserver/entity/squaddie"
	. "gopkg.in/check.v1"
	"reflect"
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
	checker.Assert(reflect.TypeOf(movement.MovementLogic()).String(), Equals, "*movement.Foot")
}

func (suite *MovementBuilder) TestChangeMovementLight(checker *C) {
	movement := squaddie.MovementBuilder().Light().Build()
	checker.Assert(reflect.TypeOf(movement.MovementLogic()).String(), Equals, "*movement.Light")
}

func (suite *MovementBuilder) TestChangeMovementFly(checker *C) {
	movement := squaddie.MovementBuilder().Fly().Build()
	checker.Assert(reflect.TypeOf(movement.MovementLogic()).String(), Equals, "*movement.Fly")
}

func (suite *MovementBuilder) TestChangeMovementTeleport(checker *C) {
	movement := squaddie.MovementBuilder().Teleport().Build()
	checker.Assert(reflect.TypeOf(movement.MovementLogic()).String(), Equals, "*movement.Teleport")
}
