package actionviewer_test

import (
	"github.com/chadius/terosbattleserver/entity/actionviewer"
	"github.com/chadius/terosbattleserver/entity/damagedistribution"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	"github.com/chadius/terosbattleserver/entity/powerusagescenario"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerattackforecast"
	"github.com/chadius/terosbattleserver/usecase/powercommit"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/usecase/squaddiestats"
	"github.com/chadius/terosbattleserver/utility"
	"github.com/chadius/terosbattleserver/utility/testutility"
	powerBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	. "gopkg.in/check.v1"
	"strings"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MockResult struct {
	ResultsPerTargetToReturn []*powercommit.ResultPerTarget
}

func (g MockResult) Forecast() *powerattackforecast.Forecast {
	return nil
}

func (g MockResult) DieRoller() utility.SixSideGenerator {
	return testutility.AlwaysHitDieRoller{}
}

func (g MockResult) ResultPerTarget() []*powercommit.ResultPerTarget {
	return g.ResultsPerTargetToReturn
}

func (g MockResult) Commit() {}

type ConsoleViewerSuite struct {
	teros   *squaddie.Squaddie
	bandit  *squaddie.Squaddie
	bandit2 *squaddie.Squaddie
	lini    *squaddie.Squaddie

	blot         *power.Power
	axe          *power.Power
	healingStaff *power.Power

	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection

	forecastBlotOnBandit *powerattackforecast.Forecast

	forecastBlotOnMultipleBandits *powerattackforecast.Forecast
	resultBlotOnMultipleBandits   *powercommit.Result

	forecastHealingStaffOnTeros *powerattackforecast.Forecast
	resultHealingStaffOnTeros   *powercommit.Result

	viewer *actionviewer.ConsoleActionViewer

	resultBlotOnBanditHit           powercommit.ResultStrategy
	resultBlotOnBanditMissed        powercommit.ResultStrategy
	resultBlotOnBanditCriticallyHit powercommit.ResultStrategy
	resultBlotOnMultipleBanditsHit  powercommit.ResultStrategy
}

var _ = Suite(&ConsoleViewerSuite{})

func (suite *ConsoleViewerSuite) SetUpTest(checker *C) {
	suite.teros = squaddieBuilder.Builder().Teros().Build()
	suite.bandit = squaddieBuilder.Builder().Bandit().Build()
	suite.bandit2 = squaddieBuilder.Builder().Bandit().WithName("Bandit2").WithID("banditID2").Build()
	suite.lini = squaddieBuilder.Builder().Lini().Build()

	suite.blot = powerBuilder.Builder().Blot().WithName("Blot").DealsDamage(0).Build()
	suite.axe = powerBuilder.Builder().Axe().Build()
	suite.healingStaff = powerBuilder.Builder().HealingStaff().WithName("healing Staff").Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{
		suite.teros,
		suite.bandit,
		suite.bandit2,
		suite.lini,
	})

	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{
		suite.blot,
		suite.axe,
		suite.healingStaff,
	})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}

	checkEquip := powerequip.CheckRepositories{}
	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.blot.GetReference(),
		},
		suite.repos,
	)

	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.bandit,
		[]*power.Reference{
			suite.axe.GetReference(),
		},
		suite.repos,
	)

	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.bandit2,
		[]*power.Reference{
			suite.axe.GetReference(),
		},
		suite.repos,
	)

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

	suite.forecastBlotOnMultipleBandits = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.blot.ID(),
			Targets:         []string{suite.bandit.ID(), suite.bandit2.ID()},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}
	suite.resultBlotOnMultipleBandits = powercommit.NewResult(suite.forecastBlotOnMultipleBandits, nil, nil)

	suite.forecastHealingStaffOnTeros = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.lini.ID(),
			PowerID:         suite.healingStaff.PowerID,
			Targets:         []string{suite.teros.ID()},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
		OffenseStrategy: &squaddiestats.CalculateSquaddieOffenseStats{},
	}

	suite.resultHealingStaffOnTeros = powercommit.NewResult(suite.forecastHealingStaffOnTeros, nil, nil)

	suite.viewer = &actionviewer.ConsoleActionViewer{}

	suite.resultBlotOnBanditHit = MockResult{
		ResultsPerTargetToReturn: []*powercommit.ResultPerTarget{
			powercommit.NewResultPerTargetBuilder().
				User(suite.teros).
				Power(suite.blot).
				Target(suite.bandit).
				AttackResult(
					powercommit.NewAttackResultBuilder().DamageDistribution(&damagedistribution.DamageDistribution{
						RawDamageDealt:    3,
						ActualDamageTaken: 3,
					}).Build(),
				).
				Build(),
		},
	}
	suite.resultBlotOnBanditCriticallyHit = MockResult{
		ResultsPerTargetToReturn: []*powercommit.ResultPerTarget{
			powercommit.NewResultPerTargetBuilder().
				User(suite.teros).
				Power(suite.blot).
				Target(suite.bandit).
				AttackResult(
					powercommit.NewAttackResultBuilder().DamageDistribution(&damagedistribution.DamageDistribution{
						RawDamageDealt:    4,
						ActualDamageTaken: 4,
					}).CriticallyHit().Build(),
				).
				Build(),
		},
	}
	suite.resultBlotOnBanditMissed = MockResult{
		ResultsPerTargetToReturn: []*powercommit.ResultPerTarget{
			powercommit.NewResultPerTargetBuilder().
				User(suite.teros).
				Power(suite.blot).
				Target(suite.bandit).
				AttackResult(
					powercommit.NewAttackResultBuilder().Build(),
				).
				Build(),
		},
	}
	suite.resultBlotOnMultipleBanditsHit = MockResult{
		ResultsPerTargetToReturn: []*powercommit.ResultPerTarget{
			powercommit.NewResultPerTargetBuilder().
				User(suite.teros).
				Power(suite.blot).
				Target(suite.bandit).
				AttackResult(
					powercommit.NewAttackResultBuilder().DamageDistribution(&damagedistribution.DamageDistribution{
						RawDamageDealt:    3,
						ActualDamageTaken: 3,
					}).Build(),
				).
				Build(),
		},
	}
}

func (suite *ConsoleViewerSuite) SetUpTerosAttacksBanditsAndSuffersCounterAttack() *powercommit.Result {
	resultBlotOnMultipleBanditsAlwaysHits := suite.resultBlotOnMultipleBandits.CopyResultWithNewDieRoller(&testutility.AlwaysHitDieRoller{})

	suite.teros.Offense = *squaddieBuilder.OffenseBuilder().Aim(suite.teros.Aim()).Mind(3).Build()

	suite.bandit.Defense = *squaddieBuilder.DefenseBuilder().
		Deflect(suite.bandit.Deflect()).
		HitPoints(suite.bandit.MaxHitPoints()).
		Barrier(1).
		Build()
	suite.bandit.Defense.SetBarrierToMax()

	suite.axe = powerBuilder.Builder().CloneOf(suite.axe).WithID(suite.axe.ID()).CanCounterAttack().DealsDamage(3).Build()
	suite.powerRepo.AddPower(suite.axe)

	suite.bandit.Offense = *squaddieBuilder.OffenseBuilder().Strength(0).Build()
	checkEquip := powerequip.CheckRepositories{}
	checkEquip.SquaddieEquipPower(suite.bandit2, suite.axe.PowerID, suite.repos)

	suite.forecastBlotOnMultipleBandits.CalculateForecast()
	resultBlotOnMultipleBanditsAlwaysHits.Commit()
	return resultBlotOnMultipleBanditsAlwaysHits
}

func (suite *ConsoleViewerSuite) SetUpLiniHealsTeros() *powercommit.Result {
	resultHealingStaffOnTerosAlwaysHits := suite.resultHealingStaffOnTeros.CopyResultWithNewDieRoller(&testutility.AlwaysHitDieRoller{})

	suite.teros.Defense.SetHPToMax()
	suite.teros.Defense.ReduceHitPoints(suite.teros.MaxHitPoints() - 1)
	suite.lini.Offense = *squaddieBuilder.OffenseBuilder().Mind(1).Build()

	suite.forecastHealingStaffOnTeros.CalculateForecast()
	resultHealingStaffOnTerosAlwaysHits.Commit()
	return resultHealingStaffOnTerosAlwaysHits
}

func (suite *ConsoleViewerSuite) TestShowPowerHitTargetAndDamageNew(checker *C) {
	var output strings.Builder
	suite.viewer.PrintResult(suite.resultBlotOnBanditHit, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Teros (Blot) hits Bandit, for 3 damage\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowPowerHitTargetAndDamage(checker *C) {
	var output strings.Builder
	suite.viewer.PrintResult(suite.resultBlotOnBanditHit, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Teros (Blot) hits Bandit, for 3 damage\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowWhenPowerMisses(checker *C) {
	var output strings.Builder
	suite.viewer.PrintResult(suite.resultBlotOnBanditMissed, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Teros (Blot) misses Bandit\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowWhenPowerCriticallyHits(checker *C) {
	var output strings.Builder
	suite.viewer.PrintResult(suite.resultBlotOnBanditCriticallyHit, suite.repos, nil, &output)
	checker.Assert(output.String(), Equals, "Teros (Blot) CRITICALLY hits Bandit, for 4 damage\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowCounterattacks(checker *C) {
	resultBlotOnBanditAndCounterAxeOnTeros := MockResult{
		ResultsPerTargetToReturn: []*powercommit.ResultPerTarget{
			powercommit.NewResultPerTargetBuilder().
				User(suite.teros).
				Power(suite.blot).
				Target(suite.bandit).
				AttackResult(
					powercommit.NewAttackResultBuilder().DamageDistribution(&damagedistribution.DamageDistribution{
						RawDamageDealt:    0,
						ActualDamageTaken: 0,
					}).Build(),
				).
				Build(),
			powercommit.NewResultPerTargetBuilder().
				User(suite.bandit).
				Power(suite.axe).
				Target(suite.teros).
				AttackResult(
					powercommit.NewAttackResultBuilder().DamageDistribution(&damagedistribution.DamageDistribution{
						RawDamageDealt:    2,
						ActualDamageTaken: 2,
					}).CounterAttack().Build(),
				).
				Build(),
		},
	}

	var output strings.Builder
	suite.viewer.PrintResult(resultBlotOnBanditAndCounterAxeOnTeros, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Teros (Blot) hits Bandit, for 0 damage\nBandit (axe) counters Teros, for 2 damage\n---\n")
}

func (suite *ConsoleViewerSuite) TestIndicateIfItIsAKillingBlow(checker *C) {
	suite.bandit.Defense.ReduceHitPoints(suite.bandit.Defense.MaxHitPoints())
	resultBlotOnBanditKills := MockResult{
		ResultsPerTargetToReturn: []*powercommit.ResultPerTarget{
			powercommit.NewResultPerTargetBuilder().
				User(suite.teros).
				Power(suite.blot).
				Target(suite.bandit).
				AttackResult(
					powercommit.NewAttackResultBuilder().HitTarget().Build(),
				).
				Build(),
		},
	}

	var output strings.Builder
	suite.viewer.PrintResult(resultBlotOnBanditKills, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Teros (Blot) hits Bandit, felling\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowPowerBarrierBurn(checker *C) {
	resultBlotOnBanditBurnsBarrier := MockResult{
		ResultsPerTargetToReturn: []*powercommit.ResultPerTarget{
			powercommit.NewResultPerTargetBuilder().
				User(suite.teros).
				Power(suite.blot).
				Target(suite.bandit).
				AttackResult(
					powercommit.NewAttackResultBuilder().DamageDistribution(&damagedistribution.DamageDistribution{
						RawDamageDealt:       3,
						ActualDamageTaken:    2,
						TotalRawBarrierBurnt: 1,
					}).Build(),
				).
				Build(),
		},
	}

	var output strings.Builder
	suite.viewer.PrintResult(resultBlotOnBanditBurnsBarrier, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Teros (Blot) hits Bandit, for 2 damage + 1 barrier burn\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowMultipleTargets(checker *C) {
	resultBlotOnBanditsAndBandit2Counters := MockResult{
		ResultsPerTargetToReturn: []*powercommit.ResultPerTarget{
			powercommit.NewResultPerTargetBuilder().
				User(suite.teros).
				Power(suite.blot).
				Target(suite.bandit).
				AttackResult(
					powercommit.NewAttackResultBuilder().DamageDistribution(&damagedistribution.DamageDistribution{
						RawDamageDealt:       3,
						ActualDamageTaken:    2,
						TotalRawBarrierBurnt: 1,
					}).Build(),
				).
				Build(),
			powercommit.NewResultPerTargetBuilder().
				User(suite.teros).
				Power(suite.blot).
				Target(suite.bandit2).
				AttackResult(
					powercommit.NewAttackResultBuilder().DamageDistribution(&damagedistribution.DamageDistribution{
						ActualDamageTaken: 3,
					}).Build(),
				).
				Build(),
			powercommit.NewResultPerTargetBuilder().
				User(suite.bandit2).
				Power(suite.axe).
				Target(suite.teros).
				AttackResult(
					powercommit.NewAttackResultBuilder().DamageDistribution(&damagedistribution.DamageDistribution{
						ActualDamageTaken: 3,
					}).CounterAttack().Build(),
				).
				Build(),
		},
	}

	var output strings.Builder
	suite.viewer.PrintResult(resultBlotOnBanditsAndBandit2Counters, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals,
		"Teros (Blot) hits Bandit, for 2 damage + 1 barrier burn\n"+
			"- also hits Bandit2, for 3 damage\n"+
			"Bandit2 (axe) counters Teros, for 3 damage\n"+
			"---\n",
	)
}

func (suite *ConsoleViewerSuite) TestShowPowerHealingEffects(checker *C) {
	resultLiniHealsTeros := MockResult{
		ResultsPerTargetToReturn: []*powercommit.ResultPerTarget{
			powercommit.NewResultPerTargetBuilder().
				User(suite.lini).
				Power(suite.healingStaff).
				Target(suite.teros).
				HealResult(
					powercommit.NewHealResultBuilder().HitPointsRestored(4).Build(),
				).
				Build(),
		},
	}

	var output strings.Builder
	suite.viewer.PrintResult(resultLiniHealsTeros, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Lini (healing Staff) heals Teros, for 4 healing\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowTargetStatusVerbosity(checker *C) {
	banditWithBarrier := squaddieBuilder.Builder().WithName("Bandit").Barrier(2).Build()
	suite.repos.SquaddieRepo.AddSquaddie(banditWithBarrier)
	resultBlotOnBanditsAndBandit2Counters := MockResult{
		ResultsPerTargetToReturn: []*powercommit.ResultPerTarget{
			powercommit.NewResultPerTargetBuilder().
				User(suite.teros).
				Power(suite.blot).
				Target(banditWithBarrier).
				AttackResult(
					powercommit.NewAttackResultBuilder().DamageDistribution(&damagedistribution.DamageDistribution{
						ActualDamageTaken:    2,
						TotalRawBarrierBurnt: 1,
					}).Build(),
				).
				Build(),
			powercommit.NewResultPerTargetBuilder().
				User(suite.teros).
				Power(suite.blot).
				Target(suite.bandit2).
				AttackResult(
					powercommit.NewAttackResultBuilder().DamageDistribution(&damagedistribution.DamageDistribution{
						ActualDamageTaken: 3,
					}).Build(),
				).
				Build(),
			powercommit.NewResultPerTargetBuilder().
				User(suite.bandit2).
				Power(suite.axe).
				Target(suite.teros).
				AttackResult(
					powercommit.NewAttackResultBuilder().DamageDistribution(&damagedistribution.DamageDistribution{
						ActualDamageTaken: 3,
					}).CounterAttack().Build(),
				).
				Build(),
		},
	}
	banditWithBarrier.Defense.ReduceHitPoints(2)
	banditWithBarrier.Defense.ReduceBarrier(1)
	suite.bandit2.Defense.ReduceHitPoints(3)
	suite.teros.Defense.ReduceHitPoints(3)

	var counterAttackOutput strings.Builder
	suite.viewer.PrintResult(
		resultBlotOnBanditsAndBandit2Counters,
		suite.repos,
		&actionviewer.ConsoleActionViewerVerbosity{
			ShowTargetStatus: true,
		},
		&counterAttackOutput,
	)

	checker.Assert(counterAttackOutput.String(), Equals,
		"Teros (Blot) hits Bandit, for 2 damage + 1 barrier burn\n"+
			"- also hits Bandit2, for 3 damage\n"+
			"   Bandit: 3/5 HP, 0 barrier\n"+
			"   Bandit2: 2/5 HP\n"+
			"Bandit2 (axe) counters Teros, for 3 damage\n"+
			"   Teros: 2/5 HP\n"+
			"---\n",
	)

	healResult := MockResult{
		ResultsPerTargetToReturn: []*powercommit.ResultPerTarget{
			powercommit.NewResultPerTargetBuilder().
				User(suite.lini).
				Power(suite.healingStaff).
				Target(suite.teros).
				HealResult(
					powercommit.NewHealResultBuilder().HitPointsRestored(4).Build(),
				).
				Build(),
		},
	}
	suite.teros.Defense.SetHPToMax()

	var output strings.Builder
	suite.viewer.PrintResult(
		healResult,
		suite.repos,
		&actionviewer.ConsoleActionViewerVerbosity{
			ShowTargetStatus: true,
		},
		&output,
	)

	checker.Assert(output.String(), Equals, "Lini (healing Staff) heals Teros, for 4 healing\n   Teros: 5/5 HP\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowRollsVerbosity(checker *C) {
	banditWithBarrier := squaddieBuilder.Builder().WithName("Bandit").Barrier(2).Build()
	suite.repos.SquaddieRepo.AddSquaddie(banditWithBarrier)
	resultBlotOnBanditsAndBandit2Counters := MockResult{
		ResultsPerTargetToReturn: []*powercommit.ResultPerTarget{
			powercommit.NewResultPerTargetBuilder().
				User(suite.teros).
				Power(suite.blot).
				Target(banditWithBarrier).
				AttackResult(
					powercommit.NewAttackResultBuilder().DamageDistribution(&damagedistribution.DamageDistribution{
						ActualDamageTaken:    2,
						TotalRawBarrierBurnt: 1,
					}).
						AttackRoll(999).
						AttackerToHitBonus(0).
						AttackerTotal(999).
						DefendRoll(-999).
						DefenderToHitPenalty(0).
						DefenderTotal(-999).
						Build(),
				).
				Build(),
			powercommit.NewResultPerTargetBuilder().
				User(suite.teros).
				Power(suite.blot).
				Target(suite.bandit2).
				AttackResult(
					powercommit.NewAttackResultBuilder().DamageDistribution(&damagedistribution.DamageDistribution{
						ActualDamageTaken: 3,
					}).
						AttackRoll(999).
						AttackerToHitBonus(0).
						AttackerTotal(999).
						DefendRoll(-999).
						DefenderToHitPenalty(0).
						DefenderTotal(-999).
						Build(),
				).
				Build(),
			powercommit.NewResultPerTargetBuilder().
				User(suite.bandit2).
				Power(suite.axe).
				Target(suite.teros).
				AttackResult(
					powercommit.NewAttackResultBuilder().DamageDistribution(&damagedistribution.DamageDistribution{
						ActualDamageTaken: 3,
					}).
						CounterAttack().
						AttackRoll(999).
						AttackerToHitBonus(-1).
						AttackerTotal(998).
						DefendRoll(-999).
						DefenderToHitPenalty(0).
						DefenderTotal(-999).
						Build(),
				).
				Build(),
		},
	}
	banditWithBarrier.Defense.ReduceHitPoints(2)
	banditWithBarrier.Defense.ReduceBarrier(1)
	suite.bandit2.Defense.ReduceHitPoints(3)
	suite.teros.Defense.ReduceHitPoints(3)

	var counterAttackOutput strings.Builder
	suite.viewer.PrintResult(
		resultBlotOnBanditsAndBandit2Counters,
		suite.repos,
		&actionviewer.ConsoleActionViewerVerbosity{
			ShowRolls: true,
		},
		&counterAttackOutput,
	)

	checker.Assert(counterAttackOutput.String(), Equals,
		"Teros (Blot) hits Bandit, for 2 damage + 1 barrier burn\n"+
			"- also hits Bandit2, for 3 damage\n"+
			"   Teros rolls 999 + 0 = 999, Bandit rolls -999 + 0 = -999\n"+
			"   Teros rolls 999 + 0 = 999, Bandit2 rolls -999 + 0 = -999\n"+
			"Bandit2 (axe) counters Teros, for 3 damage\n"+
			"   Bandit2 rolls 999 + -1 = 998, Teros rolls -999 + 0 = -999\n"+
			"---\n",
	)

	healResult := MockResult{
		ResultsPerTargetToReturn: []*powercommit.ResultPerTarget{
			powercommit.NewResultPerTargetBuilder().
				User(suite.lini).
				Power(suite.healingStaff).
				Target(suite.teros).
				HealResult(
					powercommit.NewHealResultBuilder().HitPointsRestored(4).Build(),
				).
				Build(),
		},
	}
	suite.teros.Defense.SetHPToMax()
	var healingOutput strings.Builder
	suite.viewer.PrintResult(
		healResult,
		suite.repos,
		&actionviewer.ConsoleActionViewerVerbosity{
			ShowRolls: true,
		},
		&healingOutput,
	)

	checker.Assert(healingOutput.String(), Equals, "Lini (healing Staff) heals Teros, for 4 healing\n   Auto-hit\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowForecastChanceToHitAndHealing(checker *C) {
	suite.teros.Offense = *squaddieBuilder.OffenseBuilder().
		Aim(2).
		Build()

	suite.bandit2.Defense = *squaddieBuilder.DefenseBuilder().
		Deflect(2).
		Barrier(20).
		Build()
	suite.bandit2.Defense.SetBarrierToMax()

	suite.SetUpTerosAttacksBanditsAndSuffersCounterAttack()
	var forecastOutput strings.Builder
	suite.viewer.PrintForecast(
		suite.forecastBlotOnMultipleBandits,
		suite.repos,
		&forecastOutput,
	)

	checker.Assert(forecastOutput.String(), Equals,
		"Teros (Blot) vs Bandit: +2 (30/36), for 2 damage + 1 barrier burn\n"+
			"- also vs Bandit2: +0 (21/36) for NO DAMAGE + 3 barrier burn\n"+
			"Bandit2 (axe) counters Teros: -1 (15/36), for 3 damage\n",
	)

	suite.SetUpLiniHealsTeros()
	var healingOutput strings.Builder
	suite.viewer.PrintForecast(
		suite.forecastHealingStaffOnTeros,
		suite.repos,
		&healingOutput,
	)

	checker.Assert(healingOutput.String(), Equals, "Lini (healing Staff) heals Teros, for 4 healing\n")
}

func (suite *ConsoleViewerSuite) TestShowForecastChanceToCriticallyHitAndGuaranteedMiss(checker *C) {
	suite.teros.Offense = *squaddieBuilder.OffenseBuilder().Aim(2).Build()
	suite.blot = powerBuilder.Builder().CloneOf(suite.blot).WithID(suite.blot.ID()).CriticalDealsDamage(1).CriticalHitThresholdBonus(1).Build()
	suite.powerRepo.AddPower(suite.blot)
	suite.bandit.Defense = *squaddieBuilder.DefenseBuilder().Deflect(-200).Build()
	suite.bandit2.Defense = *squaddieBuilder.DefenseBuilder().Deflect(2).Barrier(20).Build()
	suite.bandit2.Defense.SetBarrierToMax()
	suite.SetUpTerosAttacksBanditsAndSuffersCounterAttack()

	var forecastOutput strings.Builder
	suite.viewer.PrintForecast(
		suite.forecastBlotOnMultipleBandits,
		suite.repos,
		&forecastOutput,
	)

	checker.Assert(forecastOutput.String(), Equals, "Teros (Blot) vs Bandit: +202 (36/36), for 2 damage + 1 barrier burn\n"+
		" crit: 36/36, for 3 damage + 1 barrier burn\n"+
		"- also vs Bandit2: +0 (21/36) for NO DAMAGE + 3 barrier burn\n"+
		" crit: 1/36 for NO DAMAGE + 4 barrier burn\n"+
		"Bandit2 (axe) counters Teros: -1 (15/36), for 3 damage\n",
	)
}

func (suite *ConsoleViewerSuite) TestShowForecastAttackIsFatal(checker *C) {

	suite.bandit.Defense = *squaddieBuilder.DefenseBuilder().HitPoints(2).Build()
	suite.bandit.Defense.SetHPToMax()

	suite.bandit2.Defense = *squaddieBuilder.DefenseBuilder().Deflect(2).HitPoints(2).Build()
	suite.bandit2.Defense.SetHPToMax()

	suite.SetUpTerosAttacksBanditsAndSuffersCounterAttack()

	var forecastOutput strings.Builder
	suite.viewer.PrintForecast(
		suite.forecastBlotOnMultipleBandits,
		suite.repos,
		&forecastOutput,
	)

	checker.Assert(forecastOutput.String(), Equals, "Teros (Blot) vs Bandit: +0 (21/36), FATAL\n- also vs Bandit2: -2 (10/36), FATAL\nBandit2 (axe) counters Teros: -1 (15/36), for 3 damage\n")
}
