package powerattackforecast_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	"github.com/chadius/terosbattleserver/entity/powerusagescenario"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerattackforecast"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	. "gopkg.in/check.v1"
)

type VersusContextTestSuite struct {
	teros  *squaddie.Squaddie
	bandit *squaddie.Squaddie
	spear  *power.Power
	blot   *power.Power
	axe    *power.Power

	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository

	forecastSpearOnBandit *powerattackforecast.Forecast
	forecastBlotOnBandit  *powerattackforecast.Forecast
}

var _ = Suite(&VersusContextTestSuite{})

func (suite *VersusContextTestSuite) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Aim(2).Strength(2).Mind(2).Build()

	suite.spear = power.NewPowerBuilder().Spear().Build()
	suite.blot = power.NewPowerBuilder().Blot().Build()

	suite.bandit = squaddie.NewSquaddieBuilder().Bandit().Barrier(3).Armor(1).Deflect(2).Dodge(1).Build()
	suite.bandit.Defense.SetBarrierToMax()

	suite.axe = power.NewPowerBuilder().Axe().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.bandit})

	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.blot, suite.axe})

	suite.forecastSpearOnBandit = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.spear.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}

	suite.forecastBlotOnBandit = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.blot.ID(),
			Targets:         []string{suite.bandit.ID()},
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
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.ToHit().ToHitBonus, Equals, 2)

	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.ToHit().ToHitBonus, Equals, 0)
}

func (suite *VersusContextTestSuite) TestTargetTakesFullDamageAgainstPhysicalWhenNoArmor(checker *C) {
	suite.bandit.Defense = *squaddie.DefenseBuilder().Armor(0).Barrier(0).Build()

	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().RawDamageDealt, Equals, 3)

	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().RawDamageDealt, Equals, 5)
}

func (suite *VersusContextTestSuite) TestTargetUsesArmorResistAgainstPhysicalOnly(checker *C) {
	suite.bandit.Defense = *squaddie.DefenseBuilder().Armor(1).Barrier(0).Build()

	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().DamageAbsorbedByArmor, Equals, 1)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().RawDamageDealt, Equals, 2)

	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().RawDamageDealt, Equals, 5)
}

func (suite *VersusContextTestSuite) TestTargetUsesBarrierToResistDamageFromAllAttacks(checker *C) {
	suite.bandit.Defense = *squaddie.DefenseBuilder().Armor(1).Barrier(3).Build()
	suite.bandit.Defense.SetBarrierToMax()

	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().RawDamageDealt, Equals, 0)

	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().RawDamageDealt, Equals, 2)
}

func (suite *VersusContextTestSuite) TestBarrierBurnCanSpillOverDamage(checker *C) {
	suite.blot = power.NewPowerBuilder().CloneOf(suite.blot).WithID(suite.blot.ID()).DealsDamage(1).ExtraBarrierBurn(2).Build()
	suite.powerRepo.AddPower(suite.blot)

	suite.forecastBlotOnBandit.CalculateForecast()

	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().DamageAbsorbedByBarrier, Equals, 1)
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().ExtraBarrierBurnt, Equals, 2)
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().TotalRawBarrierBurnt, Equals, 3)
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().RawDamageDealt, Equals, 2)
}

func (suite *VersusContextTestSuite) TestBarrierBurnCanBeTolerated(checker *C) {
	suite.blot = power.NewPowerBuilder().CloneOf(suite.blot).WithID(suite.blot.ID()).DealsDamage(0).ExtraBarrierBurn(1).Build()
	suite.powerRepo.AddPower(suite.blot)

	suite.forecastBlotOnBandit.CalculateForecast()

	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().DamageAbsorbedByBarrier, Equals, 2)
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().ExtraBarrierBurnt, Equals, 1)
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().TotalRawBarrierBurnt, Equals, 3)
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().RawDamageDealt, Equals, 0)
}

func (suite *VersusContextTestSuite) TestCriticalHitChanceIsShown(checker *C) {
	suite.spear = power.NewPowerBuilder().CloneOf(suite.spear).WithID(suite.spear.ID()).CriticalDealsDamage(3).CriticalHitThresholdBonus(0).Build()
	suite.powerRepo.AddPower(suite.spear)

	suite.forecastSpearOnBandit.CalculateForecast()

	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.CanCritical(), Equals, true)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.CriticalHitThreshold(), Equals, 6)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.CanCritical(), Equals, true)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.CriticalHitThreshold(), Equals, 6)
}

func (suite *VersusContextTestSuite) TestCriticalDamageDistributes(checker *C) {
	suite.spear = power.NewPowerBuilder().CloneOf(suite.spear).WithID(suite.spear.ID()).CriticalDealsDamage(3).CriticalHitThresholdBonus(0).Build()
	suite.powerRepo.AddPower(suite.spear)
	suite.bandit.Defense = *squaddie.DefenseBuilder().Armor(1).Barrier(3).Build()
	suite.bandit.Defense.SetBarrierToMax()

	suite.forecastSpearOnBandit.CalculateForecast()

	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.CriticalHitDamage().DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.CriticalHitDamage().DamageAbsorbedByArmor, Equals, 1)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.CriticalHitDamage().RawDamageDealt, Equals, 2)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.CriticalHitDamage().ExtraBarrierBurnt, Equals, 0)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.CriticalHitDamage().TotalRawBarrierBurnt, Equals, 3)
}

func (suite *VersusContextTestSuite) TestNoCriticalDamageDistributionIfCannotCritical(checker *C) {
	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.CriticalHitDamage(), IsNil)
}

func (suite *VersusContextTestSuite) TestKnowsIfAttackIsNotFatalToTarget(checker *C) {
	suite.bandit.Defense = *squaddie.DefenseBuilder().Armor(0).Barrier(0).Build()
	suite.teros.Offense = *squaddie.OffenseBuilder().Mind(0).Build()

	suite.blot = power.NewPowerBuilder().CloneOf(suite.blot).WithID(suite.blot.ID()).DealsDamage(0).Build()
	suite.powerRepo.AddPower(suite.blot)

	suite.forecastBlotOnBandit.CalculateForecast()

	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().IsFatalToTarget, Equals, false)
}

func (suite *VersusContextTestSuite) TestKnowsIfAttackIsFatalToTarget(checker *C) {
	suite.spear = power.NewPowerBuilder().CloneOf(suite.spear).WithID(suite.spear.ID()).DealsDamage(
		suite.bandit.MaxHitPoints() + suite.bandit.Armor() + suite.bandit.MaxBarrier(),
	).Build()
	suite.powerRepo.AddPower(suite.spear)

	suite.forecastSpearOnBandit.CalculateForecast()

	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.VersusContext.NormalDamage().IsFatalToTarget, Equals, true)
}
