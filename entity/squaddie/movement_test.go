package squaddie_test

import (
	"github.com/chadius/terosbattleserver/entity/squaddie"
	. "gopkg.in/check.v1"
)

type SquaddieMovementTests struct {
	teros *squaddie.Squaddie
}

var _ = Suite(&SquaddieMovementTests{})

func (suite *SquaddieMovementTests) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Build()
}

func (suite *SquaddieMovementTests) TestDefaultMovement(checker *C) {
	checker.Assert(suite.teros.Movement.MovementDistance(), Equals, 3)
	checker.Assert(suite.teros.Movement.MovementType(), Equals, squaddie.MovementType(squaddie.Foot))
}

type improveMovement struct {
	initialMovement *squaddie.Movement
}

var _ = Suite(&improveMovement{})

func (suite *improveMovement) SetUpTest(checker *C) {
	suite.initialMovement = squaddie.NewMovement(2, squaddie.Foot, false)
}

func (suite *improveMovement) TestWhenImproveIsCalled_ThenMovementImproves(checker *C) {
	suite.initialMovement.Improve(3, squaddie.Fly, true)

	checker.Assert(suite.initialMovement.MovementDistance(), Equals, 5)
	checker.Assert(suite.initialMovement.MovementType(), Equals, squaddie.Fly)
	checker.Assert(suite.initialMovement.CanHitAndRun(), Equals, true)
}
