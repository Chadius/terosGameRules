package powersource_test

import (
	"github.com/chadius/terosbattleserver/entity/powersource"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	. "gopkg.in/check.v1"
)

type SpellPowerSourceSuite struct{}

var _ = Suite(&SpellPowerSourceSuite{})

func (suite *SpellPowerSourceSuite) TestName(checker *C) {
	source := powersource.NewPowerSourceLogic("Spell")
	checker.Assert(source.Name(), Equals, "spell")
}

func (suite *SpellPowerSourceSuite) TestToHitPenalty(checker *C) {
	source := powersource.NewPowerSourceLogic("spell")
	soldier := squaddie.NewSquaddieBuilder().Deflect(2).Build()
	checker.Assert(source.ToHitPenalty(soldier), Equals, 2)
}

func (suite *SpellPowerSourceSuite) TestArmorResistance(checker *C) {
	source := powersource.NewPowerSourceLogic("spell")
	soldier := squaddie.NewSquaddieBuilder().Armor(3).Build()
	checker.Assert(source.ArmorResistance(soldier), Equals, 0)
}

func (suite *SpellPowerSourceSuite) TestBarrierResistance(checker *C) {
	source := powersource.NewPowerSourceLogic("spell")
	soldier := squaddie.NewSquaddieBuilder().Barrier(5).Build()
	checker.Assert(source.BarrierResistance(soldier), Equals, 0)

	soldier.Defense.SetBarrierToMax()
	checker.Assert(source.BarrierResistance(soldier), Equals, 5)
}

func (suite *SpellPowerSourceSuite) TestRawDamage(checker *C) {
	source := powersource.NewPowerSourceLogic("spell")
	soldier := squaddie.NewSquaddieBuilder().Strength(7).Mind(11).Build()
	checker.Assert(source.RawDamage(soldier), Equals, 11)
}
