package affiliation_test

import (
	"github.com/chadius/terosgamerules/entity/affiliation"
	. "gopkg.in/check.v1"
)

type NeutralAffiliationSuite struct{}

var _ = Suite(&NeutralAffiliationSuite{})

func (suite *NeutralAffiliationSuite) TestNeutralsAreFoesOfPlayers(checker *C) {
	neutral := affiliation.NewAffiliationLogic("neutral")
	player := affiliation.NewAffiliationLogic("player")
	checker.Assert(neutral.IsFriendsWith(player), Equals, false)
	checker.Assert(neutral.IsFoesWith(player), Equals, true)
}

func (suite *NeutralAffiliationSuite) TestNeutralsAreFoesOfAllies(checker *C) {
	neutral := affiliation.NewAffiliationLogic("neutral")
	ally := affiliation.NewAffiliationLogic("ally")
	checker.Assert(neutral.IsFriendsWith(ally), Equals, false)
	checker.Assert(neutral.IsFoesWith(ally), Equals, true)
}

func (suite *NeutralAffiliationSuite) TestNeutralsAreFoesWithOtherNeutrals(checker *C) {
	neutral := affiliation.NewAffiliationLogic("neutral")
	neutral2 := affiliation.NewAffiliationLogic("neutral")
	checker.Assert(neutral.IsFriendsWith(neutral2), Equals, false)
	checker.Assert(neutral.IsFoesWith(neutral2), Equals, true)
}
