package squaddie_test

import (
	"github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	. "gopkg.in/check.v1"
)

type OffenseBuilder struct{}

var _ = Suite(&OffenseBuilder{})

func (suite *OffenseBuilder) TestBuildOffenseWithAim(checker *C) {
	sniperRifle := squaddie.OffenseBuilder().Aim(2).Build()
	checker.Assert(2, Equals, sniperRifle.Aim())
}

func (suite *OffenseBuilder) TestBuildOffenseWithStrength(checker *C) {
	falconPunch := squaddie.OffenseBuilder().Strength(5).Build()
	checker.Assert(5, Equals, falconPunch.Strength())
}

func (suite *OffenseBuilder) TestBuildOffenseWithMind(checker *C) {
	mindCrush := squaddie.OffenseBuilder().Mind(3).Build()
	checker.Assert(3, Equals, mindCrush.Mind())
}
