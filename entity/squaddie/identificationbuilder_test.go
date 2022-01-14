package squaddie_test

import (
	"github.com/chadius/terosgamerules/entity/squaddie"
	. "gopkg.in/check.v1"
	"reflect"
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
	checker.Assert(reflect.TypeOf(teros.AffiliationLogic()).String(), Equals, "*affiliation.Player")
}

func (suite *IdentificationBuilder) TestBuildIdentificationAffiliationEnemy(checker *C) {
	bandit := squaddie.IdentificationBuilder().AsEnemy().Build()
	checker.Assert(reflect.TypeOf(bandit.AffiliationLogic()).String(), Equals, "*affiliation.Enemy")
}

func (suite *IdentificationBuilder) TestBuildIdentificationAffiliationAlly(checker *C) {
	citizen := squaddie.IdentificationBuilder().AsAlly().Build()
	checker.Assert(reflect.TypeOf(citizen.AffiliationLogic()).String(), Equals, "*affiliation.Ally")
}

func (suite *IdentificationBuilder) TestBuildIdentificationAffiliationNeutral(checker *C) {
	bomb := squaddie.IdentificationBuilder().AsNeutral().Build()
	checker.Assert(reflect.TypeOf(bomb.AffiliationLogic()).String(), Equals, "*affiliation.Neutral")
}
