package affiliation_test

import (
	"github.com/chadius/terosbattleserver/entity/affiliation"
	. "gopkg.in/check.v1"
)

type EnemyAffiliationSuite struct{}

var _ = Suite(&EnemyAffiliationSuite{})

func (suite *EnemyAffiliationSuite) TestEnemiesAreFriendsWithOtherEnemies(checker *C) {
	enemy := affiliation.NewAffiliationLogic("enemy")
	enemy2 := affiliation.NewAffiliationLogic("enemy")
	checker.Assert(enemy.IsFriendsWith(enemy2), Equals, true)
	checker.Assert(enemy.IsFoesWith(enemy2), Equals, false)
}

func (suite *EnemyAffiliationSuite) TestEnemiesAreFoesOfPlayers(checker *C) {
	enemy := affiliation.NewAffiliationLogic("enemy")
	player := affiliation.NewAffiliationLogic("player")
	checker.Assert(enemy.IsFriendsWith(player), Equals, false)
	checker.Assert(enemy.IsFoesWith(player), Equals, true)
}

func (suite *EnemyAffiliationSuite) TestEnemiesAreFoesOfAllies(checker *C) {
	enemy := affiliation.NewAffiliationLogic("enemy")
	ally := affiliation.NewAffiliationLogic("ally")
	checker.Assert(enemy.IsFriendsWith(ally), Equals, false)
	checker.Assert(enemy.IsFoesWith(ally), Equals, true)
}

func (suite *EnemyAffiliationSuite) TestEnemiesAreFoesOfNeutrals(checker *C) {
	enemy := affiliation.NewAffiliationLogic("enemy")
	neutral := affiliation.NewAffiliationLogic("neutral")
	checker.Assert(enemy.IsFriendsWith(neutral), Equals, false)
	checker.Assert(enemy.IsFoesWith(neutral), Equals, true)
}
