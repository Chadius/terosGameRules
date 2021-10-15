package powerattackforecast_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerusagescenario"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerattackforecast"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	powerFactory "github.com/chadius/terosbattleserver/utility/testutility/factory/power"
	squaddieFactory "github.com/chadius/terosbattleserver/utility/testutility/factory/squaddie"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type CounterAttackCalculate struct {
	teros      *squaddie.Squaddie
	bandit     *squaddie.Squaddie
	mysticMage *squaddie.Squaddie

	spear    *power.Power
	fireball *power.Power
	axe      *power.Power

	powerRepo    *power.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection

	forecastSpearOnBandit     *powerattackforecast.Forecast
	forecastSpearOnMysticMage *powerattackforecast.Forecast
}

var _ = Suite(&CounterAttackCalculate{})

func (suite *CounterAttackCalculate) SetUpTest(checker *C) {
	suite.teros = squaddieFactory.SquaddieFactory().Teros().Build()
	suite.mysticMage = squaddieFactory.SquaddieFactory().MysticMage().Build()
	suite.bandit = squaddieFactory.SquaddieFactory().Bandit().Build()

	suite.spear = powerFactory.PowerFactory().Spear().Build()
	suite.axe = powerFactory.PowerFactory().Axe().Build()
	suite.fireball = powerFactory.PowerFactory().IsSpell().CanBeEquipped().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.bandit, suite.mysticMage})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.axe, suite.fireball})

	suite.repos = &repositories.RepositoryCollection{PowerRepo: suite.powerRepo, SquaddieRepo: suite.squaddieRepo}

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

	suite.forecastSpearOnMysticMage = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.spear.ID,
			Targets:         []string{suite.mysticMage.Identification.ID},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}
}

func (suite *CounterAttackCalculate) TestNoCounterAttackHappensIfNoEquippedPower(checker *C) {
	suite.forecastSpearOnMysticMage.CalculateForecast()

	checker.Assert(suite.forecastSpearOnMysticMage.ForecastedResultPerTarget[0].CounterAttack, IsNil)
}

func (suite *CounterAttackCalculate) TestNoCounterAttackHappensIfEquippedPowerCannotCounter(checker *C) {
	powerAddedErrors := suite.mysticMage.PowerCollection.AddInnatePower(suite.fireball)
	checker.Assert(powerAddedErrors, IsNil)

	mysticMageEquipsFireball := powerequip.SquaddieEquipPower(suite.mysticMage, suite.fireball.ID, suite.repos)
	checker.Assert(mysticMageEquipsFireball, Equals, true)

	suite.forecastSpearOnMysticMage.CalculateForecast()

	checker.Assert(suite.forecastSpearOnMysticMage.ForecastedResultPerTarget[0].CounterAttack, IsNil)
}

func (suite *CounterAttackCalculate) TestCounterAttackHappensIfPossible(checker *C) {
	suite.axe.AttackEffect.CanCounterAttack = true
	powerAddedErrors := suite.bandit.PowerCollection.AddInnatePower(suite.axe)
	checker.Assert(powerAddedErrors, IsNil)

	banditEquipsAxe := powerequip.SquaddieEquipPower(suite.bandit, suite.axe.ID, suite.repos)
	checker.Assert(banditEquipsAxe, Equals, true)

	suite.forecastSpearOnBandit.CalculateForecast()

	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget[0].CounterAttack.VersusContext.ToHit.ToHitBonus, Equals, -1)
}

type HealingEffectForecast struct {
	lini  *squaddie.Squaddie
	teros *squaddie.Squaddie
	vale  *squaddie.Squaddie

	healingStaff *power.Power

	powerRepo    *power.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection

	forecastHealingStaffOnTeros        *powerattackforecast.Forecast
	forecastHealingStaffOnTerosAndVale *powerattackforecast.Forecast
}

var _ = Suite(&HealingEffectForecast{})

func (suite *HealingEffectForecast) SetUpTest(checker *C) {
	suite.teros = squaddieFactory.SquaddieFactory().Teros().Build()
	suite.lini = squaddieFactory.SquaddieFactory().Lini().Build()
	suite.vale = squaddieFactory.SquaddieFactory().WithName("Vale").AsPlayer().Build()

	suite.healingStaff = powerFactory.PowerFactory().HealingStaff().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.lini, suite.vale})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.healingStaff})

	suite.repos = &repositories.RepositoryCollection{PowerRepo: suite.powerRepo, SquaddieRepo: suite.squaddieRepo}

	suite.forecastHealingStaffOnTeros = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.lini.Identification.ID,
			PowerID:         suite.healingStaff.ID,
			Targets:         []string{suite.teros.Identification.ID},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}

	suite.forecastHealingStaffOnTerosAndVale = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.lini.Identification.ID,
			PowerID:         suite.healingStaff.ID,
			Targets:         []string{suite.teros.Identification.ID, suite.vale.Identification.ID},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}
}

func (suite *HealingEffectForecast) TestForecastedHealingUsesHealingEffect(checker *C) {
	suite.teros.Defense.CurrentHitPoints = 1
	suite.forecastHealingStaffOnTeros.CalculateForecast()

	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget[0].HealingForecast, NotNil)
	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget[0].HealingForecast.RawHitPointsRestored, Equals, suite.healingStaff.HealingEffect.HitPointsHealed)
}

func (suite *HealingEffectForecast) TestForecastedHealingAppliesMindStat(checker *C) {
	suite.teros.Defense.MaxHitPoints = 10
	suite.teros.Defense.CurrentHitPoints = 1
	suite.lini.Offense.Mind = 3
	suite.forecastHealingStaffOnTeros.CalculateForecast()

	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget[0].HealingForecast.RawHitPointsRestored, Equals, suite.healingStaff.HealingEffect.HitPointsHealed+suite.lini.Offense.Mind)
}

func (suite *HealingEffectForecast) TestForecastedHealingCanBeHalved(checker *C) {
	suite.teros.Defense.CurrentHitPoints = 1
	suite.lini.Offense.Mind = 3
	suite.healingStaff.HealingEffect.HealingAdjustmentBasedOnUserMind = power.Half
	suite.forecastHealingStaffOnTeros.CalculateForecast()

	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget[0].HealingForecast.RawHitPointsRestored, Equals, suite.healingStaff.HealingEffect.HitPointsHealed+(suite.lini.Offense.Mind)/2)
}

func (suite *HealingEffectForecast) TestForecastedHealingCanBeZeroed(checker *C) {
	suite.teros.Defense.CurrentHitPoints = 1
	suite.lini.Offense.Mind = 3
	suite.healingStaff.HealingEffect.HealingAdjustmentBasedOnUserMind = power.Zero
	suite.forecastHealingStaffOnTeros.CalculateForecast()

	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget[0].HealingForecast.RawHitPointsRestored, Equals, suite.healingStaff.HealingEffect.HitPointsHealed)
}

func (suite *HealingEffectForecast) TestForecastedHealingCapsAtMaxHP(checker *C) {
	suite.teros.Defense.ReduceHitPoints(1)
	suite.forecastHealingStaffOnTeros.CalculateForecast()

	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget[0].HealingForecast, NotNil)
	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget[0].HealingForecast.RawHitPointsRestored, Equals, 1)
}

func (suite *HealingEffectForecast) TestHealMultipleTargets(checker *C) {
	suite.forecastHealingStaffOnTerosAndVale.CalculateForecast()

	checker.Assert(suite.forecastHealingStaffOnTerosAndVale.ForecastedResultPerTarget, HasLen, 2)
	checker.Assert(suite.forecastHealingStaffOnTerosAndVale.ForecastedResultPerTarget[0].HealingForecast, NotNil)
	checker.Assert(suite.forecastHealingStaffOnTerosAndVale.ForecastedResultPerTarget[0].HealingForecast.TargetID, Equals, suite.teros.Identification.ID)
	checker.Assert(suite.forecastHealingStaffOnTerosAndVale.ForecastedResultPerTarget[1].HealingForecast, NotNil)
	checker.Assert(suite.forecastHealingStaffOnTerosAndVale.ForecastedResultPerTarget[1].HealingForecast.TargetID, Equals, suite.vale.Identification.ID)
}
