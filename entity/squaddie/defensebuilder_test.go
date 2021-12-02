package squaddie_test

import (
	"github.com/chadius/terosbattleserver/entity/squaddie"
	. "gopkg.in/check.v1"
)

type DefenseBuilder struct{}

var _ = Suite(&DefenseBuilder{})

func (suite *DefenseBuilder) TestBuildDefenseWithArmor(checker *C) {
	plateMail := squaddie.DefenseBuilder().Armor(5).Build()
	checker.Assert(5, Equals, plateMail.Armor())
}

func (suite *DefenseBuilder) TestBuildDefenseWithDodge(checker *C) {
	ninjaGarb := squaddie.DefenseBuilder().Dodge(2).Build()
	checker.Assert(2, Equals, ninjaGarb.Dodge())
}

func (suite *DefenseBuilder) TestBuildDefenseWithDeflect(checker *C) {
	wardedRobes := squaddie.DefenseBuilder().Deflect(3).Build()
	checker.Assert(3, Equals, wardedRobes.Deflect())
}

func (suite *DefenseBuilder) TestBuildDefenseWithHitPoints(checker *C) {
	hearty := squaddie.DefenseBuilder().HitPoints(8).Build()
	checker.Assert(8, Equals, hearty.CurrentHitPoints())
	checker.Assert(8, Equals, hearty.MaxHitPoints())
}

func (suite *DefenseBuilder) TestBuildDefenseWithBarrier(checker *C) {
	glyph := squaddie.DefenseBuilder().Barrier(3).Build()
	checker.Assert(3, Equals, glyph.MaxBarrier())
}
