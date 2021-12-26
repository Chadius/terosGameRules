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

type DefenderContextTestSuite struct {
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

var _ = Suite(&DefenderContextTestSuite{})

func (suite *DefenderContextTestSuite) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Build()

	suite.spear = power.NewPowerBuilder().Spear().Build()
	suite.blot = power.NewPowerBuilder().Blot().Build()
	suite.bandit = squaddie.NewSquaddieBuilder().Bandit().Barrier(3).Armor(1).Deflect(2).Dodge(1).Build()
	suite.bandit.SetBarrierToMax()

	suite.axe = power.NewPowerBuilder().Axe().Build()

	suite.bandit.AddPowerReference(suite.axe.GetReference())

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

func (suite *DefenderContextTestSuite) TestGetDefenderDodge(checker *C) {
	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.DefenderContext.TotalToHitPenalty(), Equals, 1)
}

func (suite *DefenderContextTestSuite) TestGetDefenderArmorResistance(checker *C) {
	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.DefenderContext.ArmorResistance(), Equals, 1)

	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.DefenderContext.ArmorResistance(), Equals, 0)
}

func (suite *DefenderContextTestSuite) TestGetDefenderDeflect(checker *C) {
	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.DefenderContext.TotalToHitPenalty(), Equals, 2)
}

func (suite *DefenderContextTestSuite) TestGetDefenderBarrierAbsorb(checker *C) {
	suite.forecastBlotOnBandit.CalculateForecast()
	checker.Assert(suite.forecastBlotOnBandit.ForecastedResultPerTarget[0].Attack.DefenderContext.BarrierResistance(), Equals, 3)

	suite.forecastSpearOnBandit.CalculateForecast()
	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].Attack.DefenderContext.BarrierResistance(), Equals, 3)
}
