package target_test

import (
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/target"
	. "gopkg.in/check.v1"
)

type TargetFriendSuite struct{}

var _ = Suite(&TargetFriendSuite{})

func (suite *TargetFriendSuite) TestTargetFriendDoesWorkOnFriends(checker *C) {
	targeting := target.NewTargetingLogic("friend")
	user := squaddie.NewSquaddieBuilder().AsPlayer().Build()
	teammate := squaddie.NewSquaddieBuilder().AsPlayer().Build()

	checker.Assert(targeting.SquaddieCanTargetOtherSquaddie(user, teammate), Equals, true)
}

func (suite *TargetFriendSuite) TestTargetFriendWorksOnFriends(checker *C) {
	targeting := target.NewTargetingLogic("friend")
	user := squaddie.NewSquaddieBuilder().AsPlayer().Build()
	enemy := squaddie.NewSquaddieBuilder().AsEnemy().Build()

	checker.Assert(targeting.SquaddieCanTargetOtherSquaddie(user, enemy), Equals, false)
}
