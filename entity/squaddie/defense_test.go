package squaddie_test

import (
	"github.com/cserrant/terosBattleServer/entity/damagedistribution"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	. "gopkg.in/check.v1"
)

type SquaddieDefenseSuite struct{
	teros *squaddie.Squaddie
}

var _ = Suite(&SquaddieDefenseSuite{})

func (suite *SquaddieDefenseSuite) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
	suite.teros.Defense.MaxBarrier = 3
}

func (suite *SquaddieDefenseSuite) TestSetMaxHPAndMatchToCurrentHP(checker *C) {
	maxHP := suite.teros.Defense.MaxHitPoints
	suite.teros.Defense.SetHPToMax()
	checker.Assert(suite.teros.Defense.CurrentHitPoints, Equals, maxHP)
}

func (suite *SquaddieDefenseSuite) TestCanSetCurrentBarrierToMax(checker *C) {
	suite.teros.Defense.MaxBarrier = 2
	suite.teros.Defense.SetBarrierToMax()
	checker.Assert(suite.teros.Defense.CurrentBarrier, Equals, 2)
}

func (suite *SquaddieDefenseSuite) TestDefaultHitPoints(checker *C) {
	checker.Assert(suite.teros.Defense.MaxHitPoints, Equals, 5)
	checker.Assert(suite.teros.Defense.CurrentHitPoints, Equals, 5)
}

func (suite *SquaddieDefenseSuite) TestTakeDamageLowersHitPoints(checker *C) {
	suite.teros.Defense.ReduceHitPoints(3)
	checker.Assert(suite.teros.Defense.CurrentHitPoints, Equals, suite.teros.Defense.MaxHitPoints - 3)
}

func (suite *SquaddieDefenseSuite) TestCannotReduceHitPointsBelowZero(checker *C) {
	suite.teros.Defense.ReduceHitPoints(suite.teros.Defense.MaxHitPoints * 3)
	checker.Assert(suite.teros.Defense.CurrentHitPoints, Equals, 0)
}

func (suite *SquaddieDefenseSuite) TestTakeBarrierBurnLowersBarrier(checker *C) {
	suite.teros.Defense.MaxBarrier = 3
	suite.teros.Defense.SetBarrierToMax()
	suite.teros.Defense.ReduceBarrier(2)
	checker.Assert(suite.teros.Defense.CurrentBarrier, Equals, suite.teros.Defense.MaxBarrier - 2)
}

func (suite *SquaddieDefenseSuite) TestCannotReduceBarrierBelowZero(checker *C) {
	suite.teros.Defense.SetBarrierToMax()
	suite.teros.Defense.ReduceBarrier(suite.teros.Defense.MaxBarrier * 3)
	checker.Assert(suite.teros.Defense.CurrentBarrier, Equals, 0)
}

func (suite *SquaddieDefenseSuite) TestDamageDistributionLowersBarrierAndHealth(checker *C) {
	suite.teros.Defense.SetBarrierToMax()
	suite.teros.Defense.TakeDamageDistribution(
		&damagedistribution.DamageDistribution{
			DamageAbsorbedByBarrier: suite.teros.Defense.MaxBarrier * 3,
			RawDamageDealt:          1,
		},
	)
	checker.Assert(suite.teros.Defense.CurrentBarrier, Equals, 0)
	checker.Assert(suite.teros.Defense.CurrentHitPoints, Equals, suite.teros.Defense.MaxHitPoints - 1)
}

func (suite *SquaddieDefenseSuite) TestDamageDistributionShowsCappedDamage(checker *C) {
	suite.teros.Defense.CurrentBarrier = 1
	damageTaken := &damagedistribution.DamageDistribution{
		DamageAbsorbedByBarrier: suite.teros.Defense.MaxBarrier * 3,
		RawDamageDealt:          suite.teros.Defense.MaxHitPoints * 3,
	}
	suite.teros.Defense.TakeDamageDistribution(damageTaken)
	checker.Assert(damageTaken.ActualBarrierBurn, Equals, 1)
	checker.Assert(damageTaken.ActualDamageTaken, Equals, suite.teros.Defense.MaxHitPoints)
}

func (suite *SquaddieDefenseSuite) TestSquaddiesAreDeadWhenAtZeroHitPoints(checker *C) {
	checker.Assert(suite.teros.Defense.IsDead(), Equals, false)
	suite.teros.Defense.ReduceHitPoints(suite.teros.Defense.MaxHitPoints)
	checker.Assert(suite.teros.Defense.IsDead(), Equals, true)
}

func (suite *SquaddieDefenseSuite) TestGainHitPoints(checker *C) {
	suite.teros.Defense.CurrentHitPoints = 1
	healingAmount := suite.teros.Defense.GainHitPoints(suite.teros.Defense.MaxHitPoints)
	checker.Assert(healingAmount, Equals, suite.teros.Defense.MaxHitPoints - 1)
}