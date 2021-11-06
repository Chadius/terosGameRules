package powerattackforecast_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	"github.com/chadius/terosbattleserver/entity/powerusagescenario"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerattackforecast"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	powerBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	. "gopkg.in/check.v1"
)

type AttackContextTestSuite struct {
	teros  *squaddie.Squaddie
	bandit *squaddie.Squaddie
	spear  *power.Power
	blot   *power.Power

	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository

	forecastSpearOnBandit *powerattackforecast.Forecast
	forecastBlotOnBandit  *powerattackforecast.Forecast
}

var _ = Suite(&AttackContextTestSuite{})

func (suite *AttackContextTestSuite) SetUpTest(checker *C) {
	suite.teros = squaddieBuilder.Builder().Teros().Aim(2).Strength(2).Mind(2).Build()

	suite.spear = powerBuilder.Builder().Spear().Build()
	suite.blot = powerBuilder.Builder().Blot().Build()

	suite.bandit = squaddieBuilder.Builder().Bandit().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.bandit})

	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.blot})

	suite.forecastSpearOnBandit = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.spear.PowerID,
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
			PowerID:         suite.blot.PowerID,
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}
}

func (suite *AttackContextTestSuite) TestGetAttackerHitBonus(checker *C) {
	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.TotalToHitBonus, Equals, 3)
}

func (suite *AttackContextTestSuite) TestGetAttackerHitBonusOnCounterAttacks(checker *C) {
	forecastCounterSpearOnBandit := &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.spear.PowerID,
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: true,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}
	forecastCounterSpearOnBandit.CalculateForecast()
	checker.Assert(forecastCounterSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.TotalToHitBonus, Equals, 1)
}

func (suite *AttackContextTestSuite) TestGetAttackerPhysicalRawDamage(checker *C) {
	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.DamageType, Equals, power.DamageType(power.Physical))
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.RawDamage, Equals, 3)
}

func (suite *AttackContextTestSuite) TestGetAttackerSpellDamage(checker *C) {
	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.DamageType, Equals, power.DamageType(power.Spell))
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.RawDamage, Equals, 5)
}

func (suite *AttackContextTestSuite) TestCriticalHits(checker *C) {
	suite.spear = powerBuilder.Builder().CloneOf(suite.spear).WithID(suite.spear.ID()).CriticalDealsDamage(3).CriticalHitThresholdBonus(0).Build()
	suite.powerRepo.AddPower(suite.spear)

	suite.forecastSpearOnBandit.CalculateForecast()

	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.CanCritical, Equals, true)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.CriticalHitThreshold, Equals, 6)
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.AttackerContext.CriticalHitDamage, Equals, 6)
}
