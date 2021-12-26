package powersource_test

import (
	"github.com/chadius/terosbattleserver/entity/powersource"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	. "gopkg.in/check.v1"
)

type PhysicalPowerSourceSuite struct{}

var _ = Suite(&PhysicalPowerSourceSuite{})

func (suite *PhysicalPowerSourceSuite) TestName(checker *C) {
	source := powersource.NewPowerSourceLogic("physical")
	checker.Assert(source.Name(), Equals, "physical")
}

func (suite *PhysicalPowerSourceSuite) TestToHitPenalty(checker *C) {
	source := powersource.NewPowerSourceLogic("physical")
	soldier := squaddie.NewSquaddieBuilder().Dodge(2).Build()
	checker.Assert(source.ToHitPenalty(soldier), Equals, 2)
}

func (suite *PhysicalPowerSourceSuite) TestArmorResistance(checker *C) {
	source := powersource.NewPowerSourceLogic("physical")
	soldier := squaddie.NewSquaddieBuilder().Armor(3).Build()
	checker.Assert(source.ArmorResistance(soldier), Equals, 3)
}

func (suite *PhysicalPowerSourceSuite) TestBarrierResistance(checker *C) {
	source := powersource.NewPowerSourceLogic("physical")
	soldier := squaddie.NewSquaddieBuilder().Barrier(5).Build()
	checker.Assert(source.BarrierResistance(soldier), Equals, 0)

	soldier.SetBarrierToMax()
	checker.Assert(source.BarrierResistance(soldier), Equals, 5)
}

func (suite *PhysicalPowerSourceSuite) TestRawDamage(checker *C) {
	source := powersource.NewPowerSourceLogic("physical")
	soldier := squaddie.NewSquaddieBuilder().Strength(7).Mind(11).Build()
	checker.Assert(source.RawDamage(soldier), Equals, 7)
}
