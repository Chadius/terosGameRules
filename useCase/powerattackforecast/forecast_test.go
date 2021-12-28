package powerattackforecast_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerinterface"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	"github.com/chadius/terosbattleserver/entity/powerusagescenario"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/squaddieinterface"
	"github.com/chadius/terosbattleserver/usecase/powerattackforecast"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/usecase/squaddiestats"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type CounterAttackCalculate struct {
	teros      squaddieinterface.Interface
	bandit     squaddieinterface.Interface
	mysticMage squaddieinterface.Interface

	spear    powerinterface.Interface
	fireball powerinterface.Interface
	axe      powerinterface.Interface

	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection

	forecastSpearOnBandit     *powerattackforecast.Forecast
	forecastSpearOnMysticMage *powerattackforecast.Forecast
}

var _ = Suite(&CounterAttackCalculate{})

func (suite *CounterAttackCalculate) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Build()
	suite.mysticMage = squaddie.NewSquaddieBuilder().MysticMage().Build()
	suite.bandit = squaddie.NewSquaddieBuilder().Bandit().Build()

	suite.spear = power.NewPowerBuilder().Spear().Build()
	suite.axe = power.NewPowerBuilder().Axe().Build()
	suite.fireball = power.NewPowerBuilder().IsSpell().CanBeEquipped().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]squaddieinterface.Interface{suite.teros, suite.bandit, suite.mysticMage})

	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]powerinterface.Interface{suite.spear, suite.axe, suite.fireball})

	suite.repos = &repositories.RepositoryCollection{PowerRepo: suite.powerRepo, SquaddieRepo: suite.squaddieRepo}

	suite.forecastSpearOnBandit = powerattackforecast.NewForecastBuilder().
		Setup(
			&powerusagescenario.Setup{
				UserID:          suite.teros.ID(),
				PowerID:         suite.spear.ID(),
				Targets:         []string{suite.bandit.ID()},
				IsCounterAttack: false,
			},
		).
		Repositories(
			&repositories.RepositoryCollection{
				SquaddieRepo: suite.squaddieRepo,
				PowerRepo:    suite.powerRepo,
			},
		).
		OffenseStrategy(&squaddiestats.CalculateSquaddieOffenseStats{}).
		Build()

	suite.forecastSpearOnMysticMage = powerattackforecast.NewForecastBuilder().
		Setup(
			&powerusagescenario.Setup{
				UserID:          suite.teros.ID(),
				PowerID:         suite.spear.ID(),
				Targets:         []string{suite.mysticMage.ID()},
				IsCounterAttack: false,
			},
		).
		Repositories(
			&repositories.RepositoryCollection{
				SquaddieRepo: suite.squaddieRepo,
				PowerRepo:    suite.powerRepo,
			},
		).
		OffenseStrategy(&squaddiestats.CalculateSquaddieOffenseStats{}).
		Build()
}

func (suite *CounterAttackCalculate) TestNoCounterAttackHappensIfNoEquippedPower(checker *C) {
	suite.forecastSpearOnMysticMage.CalculateForecast()

	checker.Assert(suite.forecastSpearOnMysticMage.ForecastedResultPerTarget()[0].CounterAttack(), IsNil)
}

func (suite *CounterAttackCalculate) TestNoCounterAttackHappensIfEquippedPowerCannotCounter(checker *C) {
	suite.mysticMage.AddPowerReference(suite.fireball.GetReference())
	checkEquip := powerequip.CheckRepositories{}

	mysticMageEquipsFireball := checkEquip.SquaddieEquipPower(suite.mysticMage, suite.fireball.ID(), suite.repos)
	checker.Assert(mysticMageEquipsFireball, Equals, true)

	suite.forecastSpearOnMysticMage.CalculateForecast()

	checker.Assert(suite.forecastSpearOnMysticMage.ForecastedResultPerTarget()[0].CounterAttack(), IsNil)
}

func (suite *CounterAttackCalculate) TestCounterAttackHappensIfPossible(checker *C) {
	suite.bandit.AddPowerReference(suite.axe.GetReference())
	checkEquip := powerequip.CheckRepositories{}

	banditEquipsAxe := checkEquip.SquaddieEquipPower(suite.bandit, suite.axe.ID(), suite.repos)
	checker.Assert(banditEquipsAxe, Equals, true)

	suite.forecastSpearOnBandit.CalculateForecast()

	checker.Assert(suite.forecastSpearOnBandit.ForecastedResultPerTarget()[0].CounterAttack().VersusContext.ToHit().ToHitBonus, Equals, -1)
}

type HealingEffectForecast struct {
	lini  squaddieinterface.Interface
	teros squaddieinterface.Interface
	vale  squaddieinterface.Interface

	healingStaff powerinterface.Interface

	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection

	forecastHealingStaffOnTeros        *powerattackforecast.Forecast
	forecastHealingStaffOnTerosAndVale *powerattackforecast.Forecast
}

var _ = Suite(&HealingEffectForecast{})

func (suite *HealingEffectForecast) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Build()
	suite.lini = squaddie.NewSquaddieBuilder().Lini().Build()
	suite.vale = squaddie.NewSquaddieBuilder().WithName("Vale").AsPlayer().Build()

	suite.healingStaff = power.NewPowerBuilder().HealingStaff().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]squaddieinterface.Interface{suite.teros, suite.lini, suite.vale})

	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]powerinterface.Interface{suite.healingStaff})

	suite.repos = &repositories.RepositoryCollection{PowerRepo: suite.powerRepo, SquaddieRepo: suite.squaddieRepo}

	suite.forecastHealingStaffOnTerosAndVale = powerattackforecast.NewForecastBuilder().
		Setup(
			&powerusagescenario.Setup{
				UserID:          suite.lini.ID(),
				PowerID:         suite.healingStaff.ID(),
				Targets:         []string{suite.teros.ID(), suite.vale.ID()},
				IsCounterAttack: false,
			},
		).
		Repositories(
			&repositories.RepositoryCollection{
				SquaddieRepo: suite.squaddieRepo,
				PowerRepo:    suite.powerRepo,
			},
		).
		OffenseStrategy(&squaddiestats.CalculateSquaddieOffenseStats{}).
		Build()
}

func (suite *HealingEffectForecast) CalculateHealing(setup *powerusagescenario.Setup) {
	setupToUse := &powerusagescenario.Setup{
		UserID:          suite.lini.ID(),
		PowerID:         suite.healingStaff.ID(),
		Targets:         []string{suite.teros.ID()},
		IsCounterAttack: false,
	}
	if setup != nil {
		setupToUse = setup
	}

	suite.forecastHealingStaffOnTeros = powerattackforecast.NewForecastBuilder().
		Setup(setupToUse).
		Repositories(
			&repositories.RepositoryCollection{
				SquaddieRepo: suite.squaddieRepo,
				PowerRepo:    suite.powerRepo,
			},
		).
		OffenseStrategy(&squaddiestats.CalculateSquaddieOffenseStats{}).
		Build()

	suite.forecastHealingStaffOnTeros.CalculateForecast()
}

func (suite *HealingEffectForecast) TestForecastedHealingUsesHealingEffect(checker *C) {
	suite.teros.SetHPToMax()
	suite.teros.ReduceHitPoints(suite.teros.MaxHitPoints() - 1)
	suite.CalculateHealing(nil)

	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget()[0].HealingForecast(), NotNil)
	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget()[0].HealingForecast().RawHitPointsRestored, Equals, suite.healingStaff.HitPointsHealed())
}

func (suite *HealingEffectForecast) TestForecastedHealingAppliesMindStat(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().HitPoints(10).Build()
	suite.squaddieRepo.AddSquaddie(suite.teros)
	suite.teros.SetHPToMax()
	suite.teros.ReduceHitPoints(suite.teros.MaxHitPoints() - 1)
	suite.lini = squaddie.NewSquaddieBuilder().Mind(3).Build()
	suite.squaddieRepo.AddSquaddie(suite.lini)
	suite.CalculateHealing(nil)

	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget()[0].HealingForecast().RawHitPointsRestored, Equals, suite.healingStaff.HitPointsHealed()+suite.lini.Mind())
}

func (suite *HealingEffectForecast) TestForecastedHealingCanBeHalved(checker *C) {
	suite.teros.SetHPToMax()
	suite.teros.ReduceHitPoints(suite.teros.MaxHitPoints() - 1)
	suite.lini = squaddie.NewSquaddieBuilder().Mind(3).Build()
	suite.squaddieRepo.AddSquaddie(suite.lini)
	suite.CalculateHealing(nil)

	suite.healingStaff = power.NewPowerBuilder().CloneOf(suite.healingStaff).WithID(suite.healingStaff.ID()).HealingAdjustmentBasedOnUserMindHalf().Build()
	suite.powerRepo.AddPower(suite.healingStaff)

	suite.CalculateHealing(nil)

	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget()[0].HealingForecast().RawHitPointsRestored, Equals, suite.healingStaff.HitPointsHealed()+(suite.lini.Mind())/2)
}

func (suite *HealingEffectForecast) TestForecastedHealingCanBeZeroed(checker *C) {
	suite.teros.SetHPToMax()
	suite.teros.ReduceHitPoints(suite.teros.MaxHitPoints() - 1)
	suite.lini = squaddie.NewSquaddieBuilder().Mind(3).Build()
	suite.squaddieRepo.AddSquaddie(suite.lini)
	suite.CalculateHealing(nil)

	suite.healingStaff = power.NewPowerBuilder().CloneOf(suite.healingStaff).WithID(suite.healingStaff.ID()).HealingAdjustmentBasedOnUserMindZero().Build()
	suite.powerRepo.AddPower(suite.healingStaff)

	suite.CalculateHealing(nil)

	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget()[0].HealingForecast().RawHitPointsRestored, Equals, suite.healingStaff.HitPointsHealed())
}

func (suite *HealingEffectForecast) TestForecastedHealingCapsAtMaxHP(checker *C) {
	suite.teros.ReduceHitPoints(1)
	suite.CalculateHealing(nil)

	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget()[0].HealingForecast(), NotNil)
	checker.Assert(suite.forecastHealingStaffOnTeros.ForecastedResultPerTarget()[0].HealingForecast().RawHitPointsRestored, Equals, 1)
}

func (suite *HealingEffectForecast) TestHealMultipleTargets(checker *C) {
	suite.forecastHealingStaffOnTerosAndVale.CalculateForecast()

	checker.Assert(suite.forecastHealingStaffOnTerosAndVale.ForecastedResultPerTarget(), HasLen, 2)
	checker.Assert(suite.forecastHealingStaffOnTerosAndVale.ForecastedResultPerTarget()[0].HealingForecast(), NotNil)
	checker.Assert(suite.forecastHealingStaffOnTerosAndVale.ForecastedResultPerTarget()[0].HealingForecast().TargetID, Equals, suite.teros.ID())
	checker.Assert(suite.forecastHealingStaffOnTerosAndVale.ForecastedResultPerTarget()[1].HealingForecast(), NotNil)
	checker.Assert(suite.forecastHealingStaffOnTerosAndVale.ForecastedResultPerTarget()[1].HealingForecast().TargetID, Equals, suite.vale.ID())
}
