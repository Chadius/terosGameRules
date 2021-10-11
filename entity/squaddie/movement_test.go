package squaddie_test

import (
	"github.com/cserrant/terosbattleserver/entity/squaddie"
	. "gopkg.in/check.v1"
)

type SquaddieMovementTests struct {
	teros *squaddie.Squaddie
}

var _ = Suite(&SquaddieMovementTests{})

func (suite *SquaddieMovementTests) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
}

func (suite *SquaddieMovementTests) TestDefaultMovement(checker *C) {
	checker.Assert(suite.teros.Movement.GetMovementDistancePerRound(), Equals, 3)
	checker.Assert(suite.teros.Movement.GetMovementType(), Equals, squaddie.MovementType(squaddie.Foot))
}
