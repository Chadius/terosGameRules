package squaddie_test

import (
	"github.com/chadius/terosbattleserver/entity/movement"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	. "gopkg.in/check.v1"
	"reflect"
)

type SquaddieMovementTests struct {
	teros *squaddie.Squaddie
}

var _ = Suite(&SquaddieMovementTests{})

func (suite *SquaddieMovementTests) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Build()
}

func (suite *SquaddieMovementTests) TestDefaultMovement(checker *C) {
	checker.Assert(suite.teros.MovementDistance(), Equals, 3)
	checker.Assert(reflect.TypeOf(suite.teros.MovementLogic()).String(), Equals, "*movement.Foot")
}

type improveMovement struct {
	initialMovement *squaddie.Movement
}

var _ = Suite(&improveMovement{})

func (suite *improveMovement) SetUpTest(checker *C) {
	suite.initialMovement = squaddie.NewMovement(2, false, movement.NewMovementLogic("foot"))
}

func (suite *improveMovement) TestWhenImproveIsCalled_ThenMovementImproves(checker *C) {
	suite.initialMovement.Improve(3, true, movement.NewMovementLogic("fly"))

	checker.Assert(suite.initialMovement.MovementDistance(), Equals, 5)
	checker.Assert(reflect.TypeOf(suite.initialMovement.MovementLogic()).String(), Equals, "*movement.Fly")
	checker.Assert(suite.initialMovement.CanHitAndRun(), Equals, true)
}
