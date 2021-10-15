package squaddie_test

import (
	squaddieEntity "github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/utility/testutility/factory/squaddie"
	. "gopkg.in/check.v1"
)

type IdentificationBuilder struct {}

var _ = Suite(&IdentificationBuilder{})

func (suite *IdentificationBuilder) TestBuildIdentificationWithName(checker *C) {
	teros := squaddie.IdentificationFactory().WithName("Teros").Build()
	checker.Assert("Teros", Equals, teros.Name)
}

func (suite *IdentificationBuilder) TestBuildIdentificationWithID(checker *C) {
	teros := squaddie.IdentificationFactory().WithID("squaddieTeros").Build()
	checker.Assert("squaddieTeros", Equals, teros.ID)
}

func (suite *IdentificationBuilder) TestBuildIdentificationAffiliationPlayer(checker *C) {
	teros := squaddie.IdentificationFactory().AsPlayer().Build()
	checker.Assert(squaddieEntity.Player, Equals, teros.Affiliation)
}

func (suite *IdentificationBuilder) TestBuildIdentificationAffiliationEnemy(checker *C) {
	bandit := squaddie.IdentificationFactory().AsEnemy().Build()
	checker.Assert(squaddieEntity.Enemy, Equals, bandit.Affiliation)
}

func (suite *IdentificationBuilder) TestBuildIdentificationAffiliationAlly(checker *C) {
	citizen := squaddie.IdentificationFactory().AsAlly().Build()
	checker.Assert(squaddieEntity.Ally, Equals, citizen.Affiliation)
}

func (suite *IdentificationBuilder) TestBuildIdentificationAffiliationNeutral(checker *C) {
	bomb := squaddie.IdentificationFactory().AsNeutral().Build()
	checker.Assert(squaddieEntity.Neutral, Equals, bomb.Affiliation)
}
