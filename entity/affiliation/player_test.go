package affiliation_test

import (
	"github.com/chadius/terosgamerules/entity/affiliation"
	. "gopkg.in/check.v1"
)

type PlayerAffiliationSuite struct{}

var _ = Suite(&PlayerAffiliationSuite{})

func (suite *PlayerAffiliationSuite) TestPlayersAreFriendsWithOtherPlayers(checker *C) {
	player := affiliation.NewAffiliationLogic("player")
	player2 := affiliation.NewAffiliationLogic("player")
	checker.Assert(player.IsFriendsWith(player2), Equals, true)
	checker.Assert(player.IsFoesWith(player2), Equals, false)
}

func (suite *PlayerAffiliationSuite) TestPlayersAreFriendsWithAllies(checker *C) {
	player := affiliation.NewAffiliationLogic("player")
	ally := affiliation.NewAffiliationLogic("ally")
	checker.Assert(player.IsFriendsWith(ally), Equals, true)
	checker.Assert(player.IsFoesWith(ally), Equals, false)
}

func (suite *PlayerAffiliationSuite) TestPlayersAreFoesOfEnemies(checker *C) {
	player := affiliation.NewAffiliationLogic("player")
	enemy := affiliation.NewAffiliationLogic("enemy")
	checker.Assert(player.IsFriendsWith(enemy), Equals, false)
	checker.Assert(player.IsFoesWith(enemy), Equals, true)
}

func (suite *PlayerAffiliationSuite) TestPlayersAreFoesOfNeutrals(checker *C) {
	player := affiliation.NewAffiliationLogic("player")
	neutral := affiliation.NewAffiliationLogic("neutral")
	checker.Assert(player.IsFriendsWith(neutral), Equals, false)
	checker.Assert(player.IsFoesWith(neutral), Equals, true)
}
