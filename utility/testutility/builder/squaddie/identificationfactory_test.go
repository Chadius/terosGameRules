package squaddie_test

import (
	squaddieEntity "github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	. "gopkg.in/check.v1"
)

type IdentificationBuilder struct{}

var _ = Suite(&IdentificationBuilder{})

func (suite *IdentificationBuilder) TestBuildIdentificationWithName(checker *C) {
	teros := squaddie.IdentificationBuilder().WithName("Teros").Build()
	checker.Assert("Teros", Equals, teros.Name())
}

func (suite *IdentificationBuilder) TestBuildIdentificationWithID(checker *C) {
	teros := squaddie.IdentificationBuilder().WithID("squaddieTeros").Build()
	checker.Assert(teros.ID(), Equals, "squaddieTeros")
}

func (suite *IdentificationBuilder) TestBuildIdentificationAffiliationPlayer(checker *C) {
	teros := squaddie.IdentificationBuilder().AsPlayer().Build()
	checker.Assert(squaddieEntity.Player, Equals, teros.Affiliation())
}

func (suite *IdentificationBuilder) TestBuildIdentificationAffiliationEnemy(checker *C) {
	bandit := squaddie.IdentificationBuilder().AsEnemy().Build()
	checker.Assert(squaddieEntity.Enemy, Equals, bandit.Affiliation())
}

func (suite *IdentificationBuilder) TestBuildIdentificationAffiliationAlly(checker *C) {
	citizen := squaddie.IdentificationBuilder().AsAlly().Build()
	checker.Assert(squaddieEntity.Ally, Equals, citizen.Affiliation())
}

func (suite *IdentificationBuilder) TestBuildIdentificationAffiliationNeutral(checker *C) {
	bomb := squaddie.IdentificationBuilder().AsNeutral().Build()
	checker.Assert(squaddieEntity.Neutral, Equals, bomb.Affiliation())
}
