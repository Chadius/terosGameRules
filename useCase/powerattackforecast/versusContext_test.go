package powerattackforecast_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerusagescenario"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerattackforecast"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	powerFactory "github.com/chadius/terosbattleserver/utility/testutility/factory/power"
	squaddieFactory "github.com/chadius/terosbattleserver/utility/testutility/factory/squaddie"
	. "gopkg.in/check.v1"
)

type VersusContextTestSuite struct {
	teros  *squaddie.Squaddie
	bandit *squaddie.Squaddie
	spear  *power.Power
	blot   *power.Power
	axe    *power.Power

	powerRepo    *power.Repository
	squaddieRepo *squaddie.Repository

	forecastSpearOnBandit *powerattackforecast.Forecast
	forecastBlotOnBandit  *powerattackforecast.Forecast
}

var _ = Suite(&VersusContextTestSuite{})

func (suite *VersusContextTestSuite) SetUpTest(checker *C) {
	suite.teros = squaddieFactory.SquaddieFactory().Teros().Aim(2).Strength(2).Mind(2).Build()

	suite.spear = powerFactory.PowerFactory().Spear().Build()
	suite.blot = powerFactory.PowerFactory().Blot().Build()

	suite.bandit = squaddieFactory.SquaddieFactory().Bandit().Barrier(3).Armor(1).Deflect(2).Dodge(1).Build()
	suite.bandit.Defense.SetBarrierToMax()

	suite.axe = powerFactory.PowerFactory().Axe().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.bandit})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.blot, suite.axe})

	suite.forecastSpearOnBandit = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.spear.ID,
			Targets:         []string{suite.bandit.Identification.ID},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}

	suite.forecastBlotOnBandit = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.blot.ID,
			Targets:         []string{suite.bandit.Identification.ID},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}
}

func (suite *VersusContextTestSuite) TestNetToHitReliesOnToHitMinusDodgeOrDeflect(checker *C) {
	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.ToHit.ToHitBonus, Equals, 2)

	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.ToHit.ToHitBonus, Equals, 0)
}

func (suite *VersusContextTestSuite) TestTargetTakesFullDamageAgainstPhysicalWhenNoArmor(checker *C) {
	suite.bandit.Defense.Armor = 0
	suite.bandit.Defense.CurrentBarrier = 0

	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.RawDamageDealt, Equals, 3)

	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.RawDamageDealt, Equals, 5)
}

func (suite *VersusContextTestSuite) TestTargetUsesArmorResistAgainstPhysicalOnly(checker *C) {
	suite.bandit.Defense.Armor = 1
	suite.bandit.Defense.CurrentBarrier = 0

	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.DamageAbsorbedByArmor, Equals, 1)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.RawDamageDealt, Equals, 2)

	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.RawDamageDealt, Equals, 5)
}

func (suite *VersusContextTestSuite) TestTargetUsesBarrierToResistDamageFromAllAttacks(checker *C) {
	suite.bandit.Defense.Armor = 1
	suite.bandit.Defense.CurrentBarrier = 3

	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.RawDamageDealt, Equals, 0)

	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.RawDamageDealt, Equals, 2)
}

func (suite *VersusContextTestSuite) TestBarrierBurnCanSpillOverDamage(checker *C) {
	suite.blot.AttackEffect.DamageBonus = 1
	suite.blot.AttackEffect.ExtraBarrierBurn = 2
	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.DamageAbsorbedByBarrier, Equals, 1)

	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.ExtraBarrierBurnt, Equals, 2)
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.TotalRawBarrierBurnt, Equals, 3)

	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.RawDamageDealt, Equals, 2)
}

func (suite *VersusContextTestSuite) TestBarrierBurnCanBeTolerated(checker *C) {
	suite.blot.AttackEffect.DamageBonus = 0
	suite.blot.AttackEffect.ExtraBarrierBurn = 1
	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.DamageAbsorbedByBarrier, Equals, 2)

	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.ExtraBarrierBurnt, Equals, 1)
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.TotalRawBarrierBurnt, Equals, 3)

	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.RawDamageDealt, Equals, 0)
}

func (suite *VersusContextTestSuite) TestCriticalHitChanceIsShown(checker *C) {
	suite.spear.AttackEffect.CriticalEffect = &power.CriticalEffect{
		CriticalHitThresholdBonus: 0,
		Damage:                    3,
	}
	suite.forecastSpearOnBandit.CalculateForecast()

	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.CanCritical, Equals, true)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.CriticalHitThreshold, Equals, 6)

	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.CanCritical, Equals, true)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.CriticalHitThreshold, Equals, 6)
}

func (suite *VersusContextTestSuite) TestCriticalDamageDistributes(checker *C) {
	suite.spear.AttackEffect.CriticalEffect = &power.CriticalEffect{
		CriticalHitThresholdBonus: 0,
		Damage:                    3,
	}
	suite.bandit.Defense.Armor = 1
	suite.bandit.Defense.CurrentBarrier = 3

	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.CriticalHitDamage.DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.CriticalHitDamage.DamageAbsorbedByArmor, Equals, 1)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.CriticalHitDamage.RawDamageDealt, Equals, 2)

	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.CriticalHitDamage.ExtraBarrierBurnt, Equals, 0)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.CriticalHitDamage.TotalRawBarrierBurnt, Equals, 3)
}

func (suite *VersusContextTestSuite) TestNoCriticalDamageDistributionIfCannotCritical(checker *C) {
	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.CriticalHitDamage, IsNil)
}

func (suite *VersusContextTestSuite) TestKnowsIfAttackIsFatalToTarget(checker *C) {
	suite.bandit.Defense.Armor = 0
	suite.bandit.Defense.CurrentBarrier = 0

	suite.teros.Offense.Mind = 0
	suite.blot.AttackEffect.DamageBonus = 0
	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.IsFatalToTarget, Equals, false)

	suite.spear.AttackEffect.DamageBonus = suite.bandit.Defense.MaxHitPoints
	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage.IsFatalToTarget, Equals, true)
}
