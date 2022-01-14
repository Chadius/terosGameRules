package affiliation_test

import (
	"github.com/chadius/terosgamerules/entity/affiliation"
	. "gopkg.in/check.v1"
)

type AllyAffiliationSuite struct{}

var _ = Suite(&AllyAffiliationSuite{})

func (suite *AllyAffiliationSuite) TestAlliesAreFriendsWithPlayers(checker *C) {
	ally := affiliation.NewAffiliationLogic("ally")
	player := affiliation.NewAffiliationLogic("player")
	checker.Assert(ally.IsFriendsWith(player), Equals, true)
	checker.Assert(ally.IsFoesWith(player), Equals, false)
}

func (suite *AllyAffiliationSuite) TestAlliesAreFriendsWithOtherAllies(checker *C) {
	ally := affiliation.NewAffiliationLogic("ally")
	ally2 := affiliation.NewAffiliationLogic("ally")
	checker.Assert(ally.IsFriendsWith(ally2), Equals, true)
	checker.Assert(ally.IsFoesWith(ally2), Equals, false)
}

func (suite *AllyAffiliationSuite) TestAlliesAreFoesOfEnemies(checker *C) {
	ally := affiliation.NewAffiliationLogic("ally")
	enemy := affiliation.NewAffiliationLogic("enemy")
	checker.Assert(ally.IsFriendsWith(enemy), Equals, false)
	checker.Assert(ally.IsFoesWith(enemy), Equals, true)
}

func (suite *AllyAffiliationSuite) TestAlliesAreFoesOfNeutrals(checker *C) {
	ally := affiliation.NewAffiliationLogic("ally")
	neutral := affiliation.NewAffiliationLogic("neutral")
	checker.Assert(ally.IsFriendsWith(neutral), Equals, false)
	checker.Assert(ally.IsFoesWith(neutral), Equals, true)
}
