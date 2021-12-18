package target_test

import (
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/target"
	. "gopkg.in/check.v1"
)

type TargetSelfSuite struct{}

var _ = Suite(&TargetSelfSuite{})

func (suite *TargetSelfSuite) TestTargetSelfDoesNotWorkOnOtherSquaddie(checker *C) {
	targeting := target.NewTargetingLogic("self")
	user := squaddie.NewSquaddieBuilder().WithID("user").AsPlayer().Build()
	teammate := squaddie.NewSquaddieBuilder().WithID("teammate").AsPlayer().Build()

	checker.Assert(targeting.SquaddieCanTargetOtherSquaddie(user, teammate), Equals, false)
}

func (suite *TargetSelfSuite) TestTargetSelfWorksOnUser(checker *C) {
	targeting := target.NewTargetingLogic("self")
	user := squaddie.NewSquaddieBuilder().AsPlayer().Build()

	checker.Assert(targeting.SquaddieCanTargetOtherSquaddie(user, user), Equals, true)
}
