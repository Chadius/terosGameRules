package actionviewer_test

import (
	"github.com/cserrant/terosBattleServer/entity/actionviewer"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagescenario"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerattackforecast"
	"github.com/cserrant/terosBattleServer/usecase/powercommit"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
	"github.com/cserrant/terosBattleServer/utility/testutility"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ConsoleViewerSuite struct {
	teros			*squaddie.Squaddie
	bandit			*squaddie.Squaddie
	bandit2			*squaddie.Squaddie
	lini *squaddie.Squaddie

	blot    *power.Power
	axe      *power.Power
	healingStaff *power.Power

	powerRepo 		*power.Repository
	squaddieRepo 	*squaddie.Repository
	repos			*repositories.RepositoryCollection

	forecastBlotOnBandit *powerattackforecast.Forecast
	resultBlotOnBandit *powercommit.Result

	forecastBlotOnMultipleBandits *powerattackforecast.Forecast
	resultBlotOnMultipleBandits *powercommit.Result

	forecastHealingStaffOnTeros *powerattackforecast.Forecast
	resultHealingStaffOnTeros *powercommit.Result

	viewer *actionviewer.ConsoleActionViewer
}

var _ = Suite(&ConsoleViewerSuite{})

func (suite *ConsoleViewerSuite) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
	suite.teros.Identification.Name = "Teros"

	suite.bandit = squaddie.NewSquaddie("bandit")
	suite.bandit.Identification.Name = "Bandit"
	suite.bandit.Identification.ID = "banditID"

	suite.bandit2 = squaddie.NewSquaddie("bandit2")
	suite.bandit2.Identification.Name = "Bandit2"
	suite.bandit2.Identification.ID = "banditID2"

	suite.lini = squaddie.NewSquaddie("Lini")
	suite.lini.Identification.ID = "squaddie_lini"
	suite.lini.Identification.Name = "Lini"

	suite.blot = power.NewPower("Blot")
	suite.blot.PowerType = power.Spell
	suite.blot.AttackEffect = &power.AttackingEffect{}

	suite.axe = power.NewPower("axe")
	suite.axe.ID = "axe"
	suite.axe.PowerType = power.Physical
	suite.axe.AttackEffect = &power.AttackingEffect{
		ToHitBonus: 1,
		DamageBonus: 1,
		CanBeEquipped: true,
	}

	suite.healingStaff = power.NewPower("Healing Staff")
	suite.healingStaff.PowerType = power.Spell
	suite.healingStaff.HealingEffect = &power.HealingEffect{
		HitPointsHealed: 3,
	}

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{
		suite.teros,
		suite.bandit,
		suite.bandit2,
		suite.lini,
	})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{
		suite.blot,
		suite.axe,
		suite.healingStaff,
	})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.blot.GetReference(),
		},
		suite.repos,
	)

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.bandit,
		[]*power.Reference{
			suite.axe.GetReference(),
		},
		suite.repos,
	)

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.bandit2,
		[]*power.Reference{
			suite.axe.GetReference(),
		},
		suite.repos,
	)

	suite.forecastBlotOnBandit = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.blot.ID,
			Targets:         []string{suite.bandit.Identification.ID},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo:    suite.squaddieRepo,
			PowerRepo:       suite.powerRepo,
		},
	}
	suite.resultBlotOnBandit = &powercommit.Result{
		Forecast: suite.forecastBlotOnBandit,
	}

	suite.forecastBlotOnMultipleBandits = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.blot.ID,
			Targets:         []string{suite.bandit.Identification.ID, suite.bandit2.Identification.ID},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo:    suite.squaddieRepo,
			PowerRepo:       suite.powerRepo,
		},
	}
	suite.resultBlotOnMultipleBandits = &powercommit.Result{
		Forecast: suite.forecastBlotOnMultipleBandits,
	}

	suite.forecastHealingStaffOnTeros = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.lini.Identification.ID,
			PowerID:         suite.healingStaff.ID,
			Targets:         []string{suite.teros.Identification.ID},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo:    suite.squaddieRepo,
			PowerRepo:       suite.powerRepo,
		},
	}

	suite.resultHealingStaffOnTeros = &powercommit.Result{
		Forecast: suite.forecastHealingStaffOnTeros,
	}

	suite.viewer = &actionviewer.ConsoleActionViewer{IgnorePrinting: true}
}

func (suite *ConsoleViewerSuite) SetUpTerosAttacksBanditsAndSuffersCounterAttack() {
	suite.resultBlotOnMultipleBandits.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Offense.Mind = 3

	suite.bandit.Defense.MaxBarrier = 1
	suite.bandit.Defense.SetBarrierToMax()

	suite.axe.AttackEffect.CanCounterAttack = true
	suite.axe.AttackEffect.DamageBonus = 3
	suite.bandit.Offense.Strength = 0
	powerequip.SquaddieEquipPower(suite.bandit2, suite.axe.ID, suite.repos)

	suite.forecastBlotOnMultipleBandits.CalculateForecast()
	suite.resultBlotOnMultipleBandits.Commit()
}

func (suite *ConsoleViewerSuite) SetUpLiniHealsTeros() {
	suite.resultHealingStaffOnTeros.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Defense.CurrentHitPoints = 1
	suite.lini.Offense.Mind = 1

	suite.forecastHealingStaffOnTeros.CalculateForecast()
	suite.resultHealingStaffOnTeros.Commit()
}

func (suite *ConsoleViewerSuite) TestShowPowerHitTargetAndDamage(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Offense.Mind = 3

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	suite.viewer.PrintResult(suite.resultBlotOnBandit, suite.repos, nil)

	checker.Assert(suite.viewer.Messages, HasLen, 1)
	checker.Assert(suite.viewer.Messages[0], Equals, "Teros (Blot) hits Bandit, for 3 damage")
}

func (suite *ConsoleViewerSuite) TestShowWhenPowerMisses(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysMissDieRoller{}

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	suite.viewer.PrintResult(suite.resultBlotOnBandit, suite.repos, nil)

	checker.Assert(suite.viewer.Messages, HasLen, 1)
	checker.Assert(suite.viewer.Messages[0], Equals, "Teros (Blot) misses Bandit")
}

func (suite *ConsoleViewerSuite) TestShowWhenPowerCriticallyHits(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}
	suite.teros.Offense.Mind = 3
	suite.blot.AttackEffect = &power.AttackingEffect{
		CriticalEffect: &power.CriticalEffect{
			CriticalHitThresholdBonus: 9000,
			Damage:                    1,
		},
	}
	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	suite.viewer.PrintResult(suite.resultBlotOnBandit, suite.repos, nil)

	checker.Assert(suite.viewer.Messages, HasLen, 1)
	checker.Assert(suite.viewer.Messages[0], Equals, "Teros (Blot) CRITICALLY hits Bandit, for 4 damage")
}

func (suite *ConsoleViewerSuite) TestShowCounterattacks(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.axe.AttackEffect.CanCounterAttack = true
	suite.axe.AttackEffect.DamageBonus = 2
	suite.bandit.Offense.Strength = 0
	powerequip.SquaddieEquipPower(suite.bandit, suite.axe.ID, suite.repos)

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	suite.viewer.PrintResult(suite.resultBlotOnBandit, suite.repos, nil)

	checker.Assert(suite.viewer.Messages, HasLen, 2)
	checker.Assert(suite.viewer.Messages[1], Equals, "Bandit (axe) counters Teros, for 2 damage")
}

func (suite *ConsoleViewerSuite) TestIndicateIfItIsAKillingBlow(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Offense.Mind = suite.bandit.Defense.MaxHitPoints * 2

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	suite.viewer.PrintResult(suite.resultBlotOnBandit, suite.repos, nil)

	checker.Assert(suite.viewer.Messages, HasLen, 1)
	checker.Assert(suite.viewer.Messages[0], Equals, "Teros (Blot) hits Bandit, felling")
}

func (suite *ConsoleViewerSuite) TestShowPowerBarrierBurn(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Offense.Mind = 3
	suite.bandit.Defense.MaxBarrier = 1
	suite.bandit.Defense.SetBarrierToMax()

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	suite.viewer.PrintResult(suite.resultBlotOnBandit, suite.repos, nil)

	checker.Assert(suite.viewer.Messages, HasLen, 1)
	checker.Assert(suite.viewer.Messages[0], Equals, "Teros (Blot) hits Bandit, for 2 damage + 1 barrier burn")
}

func (suite *ConsoleViewerSuite) TestShowMultipleTargets(checker *C) {
	suite.SetUpTerosAttacksBanditsAndSuffersCounterAttack()

	suite.viewer.PrintResult(suite.resultBlotOnMultipleBandits, suite.repos, nil)

	checker.Assert(suite.viewer.Messages, HasLen, 3)
	checker.Assert(suite.viewer.Messages[0], Equals, "Teros (Blot) hits Bandit, for 2 damage + 1 barrier burn")
	checker.Assert(suite.viewer.Messages[1], Equals, "- also hits Bandit2, for 3 damage")
	checker.Assert(suite.viewer.Messages[2], Equals, "Bandit2 (axe) counters Teros, for 3 damage")
}

func (suite *ConsoleViewerSuite) TestShowPowerHealingEffects(checker *C) {
	suite.resultHealingStaffOnTeros.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Defense.CurrentHitPoints = 1
	suite.lini.Offense.Mind = 1

	suite.forecastHealingStaffOnTeros.CalculateForecast()
	suite.resultHealingStaffOnTeros.Commit()

	suite.viewer.PrintResult(suite.resultHealingStaffOnTeros, suite.repos, nil)

	checker.Assert(suite.viewer.Messages, HasLen, 1)
	checker.Assert(suite.viewer.Messages[0], Equals, "Lini (Healing Staff) heals Teros, for 4 healing")
}

func (suite *ConsoleViewerSuite) TestShowTargetStatusVerbosity(checker *C) {
	suite.SetUpTerosAttacksBanditsAndSuffersCounterAttack()

	suite.viewer.PrintResult(
		suite.resultBlotOnMultipleBandits,
		suite.repos,
		&actionviewer.ConsoleActionViewerVerbosity{
			ShowTargetStatus: true,
		},
	)

	checker.Assert(suite.viewer.Messages, HasLen, 6)
	checker.Assert(suite.viewer.Messages[0], Equals, "Teros (Blot) hits Bandit, for 2 damage + 1 barrier burn")
	checker.Assert(suite.viewer.Messages[1], Equals, "- also hits Bandit2, for 3 damage")
	checker.Assert(suite.viewer.Messages[2], Equals, "   Bandit: 3/5 HP, 0 barrier")
	checker.Assert(suite.viewer.Messages[3], Equals, "   Bandit2: 2/5 HP")
	checker.Assert(suite.viewer.Messages[4], Equals, "Bandit2 (axe) counters Teros, for 3 damage")
	checker.Assert(suite.viewer.Messages[5], Equals, "   Teros: 2/5 HP")

	suite.SetUpLiniHealsTeros()

	suite.viewer.PrintResult(
		suite.resultHealingStaffOnTeros,
		suite.repos,
		&actionviewer.ConsoleActionViewerVerbosity{
			ShowTargetStatus: true,
		},
	)

	checker.Assert(suite.viewer.Messages, HasLen, 8)
	checker.Assert(suite.viewer.Messages[6], Equals, "Lini (Healing Staff) heals Teros, for 4 healing")
	checker.Assert(suite.viewer.Messages[7], Equals, "   Teros: 5/5 HP")
}

func (suite *ConsoleViewerSuite) TestShowRollsVerbosity(checker *C) {
	suite.SetUpTerosAttacksBanditsAndSuffersCounterAttack()

	suite.viewer.PrintResult(
		suite.resultBlotOnMultipleBandits,
		suite.repos,
		&actionviewer.ConsoleActionViewerVerbosity{
			ShowRolls: true,
		},
	)

	checker.Assert(suite.viewer.Messages, HasLen, 6)
	checker.Assert(suite.viewer.Messages[0], Equals, "Teros (Blot) hits Bandit, for 2 damage + 1 barrier burn")
	checker.Assert(suite.viewer.Messages[1], Equals, "- also hits Bandit2, for 3 damage")
	checker.Assert(suite.viewer.Messages[2], Equals, "   Teros rolls 999 + 0 = 999, Bandit rolls -999 + 0 = -999")
	checker.Assert(suite.viewer.Messages[3], Equals, "   Teros rolls 999 + 0 = 999, Bandit2 rolls -999 + 0 = -999")
	checker.Assert(suite.viewer.Messages[4], Equals, "Bandit2 (axe) counters Teros, for 3 damage")
	checker.Assert(suite.viewer.Messages[5], Equals, "   Bandit2 rolls 999 + -1 = 998, Teros rolls -999 + 0 = -999")

	suite.SetUpLiniHealsTeros()

	suite.viewer.PrintResult(
		suite.resultHealingStaffOnTeros,
		suite.repos,
		&actionviewer.ConsoleActionViewerVerbosity{
			ShowRolls: true,
		},
	)

	checker.Assert(suite.viewer.Messages, HasLen, 8)
	checker.Assert(suite.viewer.Messages[6], Equals, "Lini (Healing Staff) heals Teros, for 4 healing")
	checker.Assert(suite.viewer.Messages[7], Equals, "   Auto-hit")
}

func (suite *ConsoleViewerSuite) TestShowForecastChanceToHitAndHealing(checker *C) {
	suite.teros.Offense.Aim = 2

	suite.bandit2.Defense.Deflect = 2
	suite.bandit2.Defense.MaxBarrier = 20
	suite.bandit2.Defense.SetBarrierToMax()

	suite.SetUpTerosAttacksBanditsAndSuffersCounterAttack()

	suite.viewer.PrintForecast(
		suite.forecastBlotOnMultipleBandits,
		suite.repos,
	)

	checker.Assert(suite.viewer.Messages, HasLen, 3)
	checker.Assert(suite.viewer.Messages[0], Equals, "Teros (Blot) vs Bandit: +2 (30/36), for 2 damage + 1 barrier burn")
	checker.Assert(suite.viewer.Messages[1], Equals, "- also vs Bandit2: +0 (21/36) for NO DAMAGE + 3 barrier burn")
	checker.Assert(suite.viewer.Messages[2], Equals, "Bandit2 (axe) counters Teros: -1 (15/36), for 3 damage")

	suite.SetUpLiniHealsTeros()
	suite.viewer.PrintForecast(
		suite.forecastHealingStaffOnTeros,
		suite.repos,
	)

	checker.Assert(suite.viewer.Messages[3], Equals, "Lini (Healing Staff) heals Teros, for 4 healing")
}

func (suite *ConsoleViewerSuite) TestShowForecastChanceToCriticallyHitAndGuaranteedMiss(checker *C) {
	suite.teros.Offense.Aim = 2

	suite.blot.AttackEffect = &power.AttackingEffect{
		CriticalEffect: &power.CriticalEffect{
			CriticalHitThresholdBonus: 1,
			Damage:                    1,
		},
	}

	suite.bandit.Defense.Deflect = -200

	suite.bandit2.Defense.Deflect = 2
	suite.bandit2.Defense.MaxBarrier = 20
	suite.bandit2.Defense.SetBarrierToMax()

	suite.SetUpTerosAttacksBanditsAndSuffersCounterAttack()

	suite.viewer.PrintForecast(
		suite.forecastBlotOnMultipleBandits,
		suite.repos,
	)

	checker.Assert(suite.viewer.Messages[0], Equals, "Teros (Blot) vs Bandit: +202 (36/36), for 2 damage + 1 barrier burn")
	checker.Assert(suite.viewer.Messages[1], Equals, " crit: 36/36, for 3 damage + 1 barrier burn")
	checker.Assert(suite.viewer.Messages[2], Equals, "- also vs Bandit2: +0 (21/36) for NO DAMAGE + 3 barrier burn")
	checker.Assert(suite.viewer.Messages[3], Equals, " crit: 1/36 for NO DAMAGE + 4 barrier burn")
}

func (suite *ConsoleViewerSuite) TestShowForecastAttackIsFatal(checker *C) {

	suite.bandit.Defense.MaxHitPoints = 2
	suite.bandit.Defense.SetHPToMax()

	suite.bandit2.Defense.Deflect = 2
	suite.bandit2.Defense.MaxHitPoints = 2
	suite.bandit2.Defense.SetHPToMax()

	suite.SetUpTerosAttacksBanditsAndSuffersCounterAttack()

	suite.viewer.PrintForecast(
		suite.forecastBlotOnMultipleBandits,
		suite.repos,
	)

	checker.Assert(suite.viewer.Messages[0], Equals, "Teros (Blot) vs Bandit: +0 (21/36), FATAL")
	checker.Assert(suite.viewer.Messages[1], Equals, "- also vs Bandit2: -2 (10/36), FATAL")
}
