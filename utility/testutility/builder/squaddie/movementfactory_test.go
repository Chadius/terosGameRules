package squaddie_test

import (
	squaddieEntity "github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MovementBuilder struct{}

var _ = Suite(&MovementBuilder{})

func (suite *MovementBuilder) TestBuildWithDistance(checker *C) {
	movement := squaddie.MovementBuilder().Distance(3).Build()
	checker.Assert(3, Equals, movement.Distance)
}

func (suite *MovementBuilder) TestBuildMovementCanHitAndRun(checker *C) {
	movement := squaddie.MovementBuilder().CanHitAndRun().Build()
	checker.Assert(true, Equals, movement.HitAndRun)
}

func (suite *MovementBuilder) TestChangeMovementFoot(checker *C) {
	movement := squaddie.MovementBuilder().Foot().Build()
	checker.Assert(squaddieEntity.Foot, Equals, movement.Type)
}

func (suite *MovementBuilder) TestChangeMovementLight(checker *C) {
	movement := squaddie.MovementBuilder().Light().Build()
	checker.Assert(squaddieEntity.Light, Equals, movement.Type)
}

func (suite *MovementBuilder) TestChangeMovementFly(checker *C) {
	movement := squaddie.MovementBuilder().Fly().Build()
	checker.Assert(squaddieEntity.Fly, Equals, movement.Type)
}

func (suite *MovementBuilder) TestChangeMovementTeleport(checker *C) {
	movement := squaddie.MovementBuilder().Teleport().Build()
	checker.Assert(squaddieEntity.Teleport, Equals, movement.Type)
}
