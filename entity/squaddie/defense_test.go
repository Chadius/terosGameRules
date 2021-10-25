package squaddie_test

import (
	"github.com/chadius/terosbattleserver/entity/damagedistribution"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	. "gopkg.in/check.v1"
)

type SquaddieDefenseSuite struct {
	teros *squaddie.Squaddie
}

var _ = Suite(&SquaddieDefenseSuite{})

func (suite *SquaddieDefenseSuite) SetUpTest(checker *C) {
	suite.teros = squaddieBuilder.Builder().Teros().Barrier(3).Build()
}

func (suite *SquaddieDefenseSuite) TestSetMaxHPAndMatchToCurrentHP(checker *C) {
	maxHP := suite.teros.MaxHitPoints()
	suite.teros.Defense.SetHPToMax()
	checker.Assert(suite.teros.CurrentHitPoints(), Equals, maxHP)
}

func (suite *SquaddieDefenseSuite) TestCanSetCurrentBarrierToMax(checker *C) {
	suite.teros.Defense.SquaddieMaxBarrier = 2
	suite.teros.Defense.SetBarrierToMax()
	checker.Assert(suite.teros.CurrentBarrier(), Equals, 2)
}

func (suite *SquaddieDefenseSuite) TestDefaultHitPoints(checker *C) {
	checker.Assert(suite.teros.MaxHitPoints(), Equals, 5)
	checker.Assert(suite.teros.CurrentHitPoints(), Equals, 5)
}

func (suite *SquaddieDefenseSuite) TestTakeDamageLowersHitPoints(checker *C) {
	suite.teros.Defense.ReduceHitPoints(3)
	checker.Assert(suite.teros.CurrentHitPoints(), Equals, suite.teros.Defense.SquaddieMaxHitPoints-3)
}

func (suite *SquaddieDefenseSuite) TestCannotReduceHitPointsBelowZero(checker *C) {
	suite.teros.Defense.ReduceHitPoints(suite.teros.MaxHitPoints() * 3)
	checker.Assert(suite.teros.CurrentHitPoints(), Equals, 0)
}

func (suite *SquaddieDefenseSuite) TestTakeBarrierBurnLowersBarrier(checker *C) {
	suite.teros.Defense.SquaddieMaxBarrier = 3
	suite.teros.Defense.SetBarrierToMax()
	suite.teros.Defense.ReduceBarrier(2)
	checker.Assert(suite.teros.CurrentBarrier(), Equals, suite.teros.Defense.SquaddieMaxBarrier-2)
}

func (suite *SquaddieDefenseSuite) TestCannotReduceBarrierBelowZero(checker *C) {
	suite.teros.Defense.SetBarrierToMax()
	suite.teros.Defense.ReduceBarrier(suite.teros.MaxBarrier() * 3)
	checker.Assert(suite.teros.CurrentBarrier(), Equals, 0)
}

func (suite *SquaddieDefenseSuite) TestDamageDistributionLowersBarrierAndHealth(checker *C) {
	suite.teros.Defense.SetBarrierToMax()
	suite.teros.Defense.TakeDamageDistribution(
		&damagedistribution.DamageDistribution{
			DamageAbsorbedByBarrier: suite.teros.MaxBarrier() * 3,
			RawDamageDealt:          1,
		},
	)
	checker.Assert(suite.teros.CurrentBarrier(), Equals, 0)
	checker.Assert(suite.teros.CurrentHitPoints(), Equals, suite.teros.Defense.SquaddieMaxHitPoints-1)
}

func (suite *SquaddieDefenseSuite) TestDamageDistributionShowsCappedDamage(checker *C) {
	suite.teros.Defense.SquaddieCurrentBarrier = 1
	damageTaken := &damagedistribution.DamageDistribution{
		DamageAbsorbedByBarrier: suite.teros.MaxBarrier() * 3,
		RawDamageDealt:          suite.teros.MaxHitPoints() * 3,
	}
	suite.teros.Defense.TakeDamageDistribution(damageTaken)
	checker.Assert(damageTaken.ActualBarrierBurn, Equals, 1)
	checker.Assert(damageTaken.ActualDamageTaken, Equals, suite.teros.MaxHitPoints())
}

func (suite *SquaddieDefenseSuite) TestSquaddiesAreDeadWhenAtZeroHitPoints(checker *C) {
	checker.Assert(suite.teros.Defense.IsDead(), Equals, false)
	suite.teros.Defense.ReduceHitPoints(suite.teros.MaxHitPoints())
	checker.Assert(suite.teros.Defense.IsDead(), Equals, true)
}

func (suite *SquaddieDefenseSuite) TestGainHitPoints(checker *C) {
	suite.teros.Defense.SquaddieCurrentHitPoints = 1
	healingAmount := suite.teros.Defense.GainHitPoints(suite.teros.MaxHitPoints())
	checker.Assert(healingAmount, Equals, suite.teros.MaxHitPoints() - 1)
}
