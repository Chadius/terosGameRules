package squaddie_test

import (
	"github.com/chadius/terosbattleserver/entity/squaddie"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	. "gopkg.in/check.v1"
)

type SquaddieMovementTests struct {
	teros *squaddie.Squaddie
}

var _ = Suite(&SquaddieMovementTests{})

func (suite *SquaddieMovementTests) SetUpTest(checker *C) {
	suite.teros = squaddieBuilder.Builder().Teros().Build()
}

func (suite *SquaddieMovementTests) TestDefaultMovement(checker *C) {
	checker.Assert(suite.teros.Movement.GetMovementDistancePerRound(), Equals, 3)
	checker.Assert(suite.teros.Movement.GetMovementType(), Equals, squaddie.MovementType(squaddie.Foot))
}
