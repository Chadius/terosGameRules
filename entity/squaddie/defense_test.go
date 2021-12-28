package squaddie_test

import (
	"github.com/chadius/terosbattleserver/entity/damagedistribution"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/squaddieinterface"
	. "gopkg.in/check.v1"
)

type SquaddieDefenseSuite struct {
	teros squaddieinterface.Interface
}

var _ = Suite(&SquaddieDefenseSuite{})

func (suite *SquaddieDefenseSuite) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Barrier(3).Build()
}

func (suite *SquaddieDefenseSuite) TestSetMaxHPAndMatchToCurrentHP(checker *C) {
	maxHP := suite.teros.MaxHitPoints()
	suite.teros.SetHPToMax()
	checker.Assert(suite.teros.CurrentHitPoints(), Equals, maxHP)
}

func (suite *SquaddieDefenseSuite) TestCanSetCurrentBarrierToMax(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Barrier(2).Build()
	suite.teros.SetBarrierToMax()
	checker.Assert(suite.teros.CurrentBarrier(), Equals, 2)
}

func (suite *SquaddieDefenseSuite) TestDefaultHitPoints(checker *C) {
	checker.Assert(suite.teros.MaxHitPoints(), Equals, 5)
	checker.Assert(suite.teros.CurrentHitPoints(), Equals, 5)
}

func (suite *SquaddieDefenseSuite) TestTakeDamageLowersHitPoints(checker *C) {
	suite.teros.ReduceHitPoints(3)
	checker.Assert(suite.teros.CurrentHitPoints(), Equals, suite.teros.MaxHitPoints()-3)
}

func (suite *SquaddieDefenseSuite) TestCannotReduceHitPointsBelowZero(checker *C) {
	suite.teros.ReduceHitPoints(suite.teros.MaxHitPoints() * 3)
	checker.Assert(suite.teros.CurrentHitPoints(), Equals, 0)
}

func (suite *SquaddieDefenseSuite) TestTakeBarrierBurnLowersBarrier(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Barrier(3).Build()
	suite.teros.SetBarrierToMax()
	suite.teros.ReduceBarrier(2)
	checker.Assert(suite.teros.CurrentBarrier(), Equals, suite.teros.MaxBarrier()-2)
}

func (suite *SquaddieDefenseSuite) TestCannotReduceBarrierBelowZero(checker *C) {
	suite.teros.SetBarrierToMax()
	suite.teros.ReduceBarrier(suite.teros.MaxBarrier() * 3)
	checker.Assert(suite.teros.CurrentBarrier(), Equals, 0)
}

func (suite *SquaddieDefenseSuite) TestDamageDistributionLowersBarrierAndHealth(checker *C) {
	suite.teros.SetBarrierToMax()
	suite.teros.TakeDamageDistribution(
		&damagedistribution.DamageDistribution{
			DamageAbsorbedByBarrier: suite.teros.MaxBarrier() * 3,
			RawDamageDealt:          1,
		},
	)
	checker.Assert(suite.teros.CurrentBarrier(), Equals, 0)
	checker.Assert(suite.teros.CurrentHitPoints(), Equals, suite.teros.MaxHitPoints()-1)
}

func (suite *SquaddieDefenseSuite) TestDamageDistributionShowsCappedDamage(checker *C) {
	suite.teros.SetBarrierToMax()
	suite.teros.ReduceBarrier(suite.teros.MaxBarrier() - 1)
	damageTaken := &damagedistribution.DamageDistribution{
		DamageAbsorbedByBarrier: suite.teros.MaxBarrier() * 3,
		RawDamageDealt:          suite.teros.MaxHitPoints() * 3,
	}
	suite.teros.TakeDamageDistribution(damageTaken)
	checker.Assert(damageTaken.ActualBarrierBurn, Equals, 1)
	checker.Assert(damageTaken.ActualDamageTaken, Equals, suite.teros.MaxHitPoints())
}

func (suite *SquaddieDefenseSuite) TestSquaddiesAreDeadWhenAtZeroHitPoints(checker *C) {
	checker.Assert(suite.teros.IsDead(), Equals, false)
	suite.teros.ReduceHitPoints(suite.teros.MaxHitPoints())
	checker.Assert(suite.teros.IsDead(), Equals, true)
}

func (suite *SquaddieDefenseSuite) TestGainHitPoints(checker *C) {
	suite.teros.SetHPToMax()
	suite.teros.ReduceHitPoints(suite.teros.MaxHitPoints() - 1)
	healingAmount := suite.teros.GainHitPoints(suite.teros.MaxHitPoints())
	checker.Assert(healingAmount, Equals, suite.teros.MaxHitPoints()-1)
}

type improveDefense struct {
	initialDefense *squaddie.Defense
}

var _ = Suite(&improveDefense{})

func (suite *improveDefense) SetUpTest(checker *C) {
	suite.initialDefense = squaddie.NewDefense(0, 2, 3, 5, 0, 7, 11)
}

func (suite *improveDefense) TestWhenImproveIsCalled_ThenAimStrengthMindIncrease(checker *C) {
	suite.initialDefense.Improve(2, 3, 5, 7, 11)

	checker.Assert(suite.initialDefense.MaxHitPoints(), Equals, 4)
	checker.Assert(suite.initialDefense.Dodge(), Equals, 6)
	checker.Assert(suite.initialDefense.Deflect(), Equals, 10)
	checker.Assert(suite.initialDefense.MaxBarrier(), Equals, 14)
	checker.Assert(suite.initialDefense.Armor(), Equals, 22)
}
