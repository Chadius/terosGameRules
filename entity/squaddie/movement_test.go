package squaddie_test

import (
	"github.com/chadius/terosbattleserver/entity/squaddie"
	squaddieFactory "github.com/chadius/terosbattleserver/utility/testutility/factory/squaddie"
	. "gopkg.in/check.v1"
)

type SquaddieMovementTests struct {
	teros *squaddie.Squaddie
}

var _ = Suite(&SquaddieMovementTests{})

func (suite *SquaddieMovementTests) SetUpTest(checker *C) {
	suite.teros = squaddieFactory.SquaddieFactory().Teros().Build()
}

func (suite *SquaddieMovementTests) TestDefaultMovement(checker *C) {
	checker.Assert(suite.teros.Movement.GetMovementDistancePerRound(), Equals, 3)
	checker.Assert(suite.teros.Movement.GetMovementType(), Equals, squaddie.MovementType(squaddie.Foot))
}
