package target_test

import (
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/target"
	. "gopkg.in/check.v1"
)

type TargetFoeSuite struct{}

var _ = Suite(&TargetFoeSuite{})

func (suite *TargetFoeSuite) TestTargetFoeTargetsFoes(checker *C) {
	targeting := target.NewTargetingLogic("foe")
	user := squaddie.NewSquaddieBuilder().AsPlayer().Build()
	enemy := squaddie.NewSquaddieBuilder().AsEnemy().Build()

	checker.Assert(targeting.SquaddieCanTargetOtherSquaddie(user, enemy), Equals, true)
}

func (suite *TargetFoeSuite) TestTargetFoeDoesNotWorkOnAllies(checker *C) {
	targeting := target.NewTargetingLogic("foe")
	user := squaddie.NewSquaddieBuilder().AsPlayer().Build()
	teammate := squaddie.NewSquaddieBuilder().AsPlayer().Build()

	checker.Assert(targeting.SquaddieCanTargetOtherSquaddie(user, teammate), Equals, false)
}
