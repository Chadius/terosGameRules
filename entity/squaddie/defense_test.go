package squaddie_test

import (
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	. "gopkg.in/check.v1"
)

type SquaddieDefenstSuite struct{
	teros *squaddie.Squaddie
}

var _ = Suite(&SquaddieDefenstSuite{})

func (suite *SquaddieDefenstSuite) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
}

func (suite *SquaddieDefenstSuite) TestSetMaxHPAndMatchToCurrentHP(checker *C) {
	maxHP := suite.teros.Defense.MaxHitPoints
	suite.teros.Defense.SetHPToMax()
	checker.Assert(suite.teros.Defense.CurrentHitPoints, Equals, maxHP)
}

func (suite *SquaddieDefenstSuite) TestCanSetCurrentBarrierToMax(checker *C) {
	suite.teros.Defense.MaxBarrier = 2
	suite.teros.Defense.SetBarrierToMax()
	checker.Assert(suite.teros.Defense.CurrentBarrier, Equals, 2)
}

func (suite *SquaddieDefenstSuite) TestDefaultHitPoints(checker *C) {
	checker.Assert(suite.teros.Defense.MaxHitPoints, Equals, 5)
	checker.Assert(suite.teros.Defense.CurrentHitPoints, Equals, 5)
}
