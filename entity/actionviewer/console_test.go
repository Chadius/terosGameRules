package actionviewer_test

import (
	"github.com/chadius/terosgamerules/entity/actionviewer"
	"github.com/chadius/terosgamerules/entity/damagedistribution"
	"github.com/chadius/terosgamerules/entity/power"
	"github.com/chadius/terosgamerules/entity/powerinterface"
	"github.com/chadius/terosgamerules/entity/powerrepository"
	"github.com/chadius/terosgamerules/entity/powerusagescenario"
	"github.com/chadius/terosgamerules/entity/squaddie"
	"github.com/chadius/terosgamerules/entity/squaddieinterface"
	"github.com/chadius/terosgamerules/usecase/powerattackforecast"
	"github.com/chadius/terosgamerules/usecase/powerattackforecast/powerattackforecastfakes"
	"github.com/chadius/terosgamerules/usecase/powercommit"
	"github.com/chadius/terosgamerules/usecase/powercommit/powercommitfakes"
	"github.com/chadius/terosgamerules/usecase/repositories"
	"github.com/chadius/terosgamerules/utility/testutility"
	. "gopkg.in/check.v1"
	"strings"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ConsoleShowsFatalAttacksSuite struct {
	teros   squaddieinterface.Interface
	bandit  squaddieinterface.Interface
	bandit2 squaddieinterface.Interface

	blot powerinterface.Interface

	viewer       *actionviewer.ConsoleActionViewer
	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection
}

var _ = Suite(&ConsoleShowsFatalAttacksSuite{})

func (suite *ConsoleShowsFatalAttacksSuite) SetUpTest(checker *C) {
	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}
	suite.viewer = &actionviewer.ConsoleActionViewer{}

	suite.teros = squaddie.NewSquaddieBuilder().Teros().Mind(3).Build()
	suite.bandit = squaddie.NewSquaddieBuilder().Bandit().WithName("Bandit").WithID("banditID1").
		Strength(0).
		HitPoints(2).
		Barrier(1).
		Build()
	suite.bandit.SetBarrierToMax()

	suite.bandit2 = squaddie.NewSquaddieBuilder().Bandit().WithName("Bandit2").WithID("banditID2").
		HitPoints(2).
		Deflect(2).
		Build()

	suite.blot = power.NewPowerBuilder().Blot().WithName("Blot").DealsDamage(0).Build()

	testutility.AddSquaddieWithInnatePowersToRepos(suite.teros, suite.blot, suite.repos, false)
	suite.repos.SquaddieRepo.AddSquaddies([]squaddieinterface.Interface{suite.bandit, suite.bandit2})
}

func (suite *ConsoleShowsFatalAttacksSuite) TestShowForecastAttackIsFatal(checker *C) {
	killFirstBanditCalculation := &powerattackforecastfakes.FakeCalculationInterface{}
	killFirstBanditCalculation.SetupReturns(&powerusagescenario.Setup{
		UserID:          suite.teros.ID(),
		PowerID:         suite.blot.ID(),
		Targets:         []string{suite.bandit.ID(), suite.bandit2.ID()},
		IsCounterAttack: false,
	})
	killFirstBanditCalculation.AttackReturns(&powerattackforecast.AttackForecast{
		DefenderContext: *powerattackforecast.NewDefenderContext(suite.bandit.ID(), nil),
		VersusContext: &powerattackforecastfakes.FakeVersusContextStrategy{
			NormalDamageStub: func() *damagedistribution.DamageDistribution {
				return &damagedistribution.DamageDistribution{
					IsFatalToTarget: true,
				}
			},
			ToHitStub: func() *damagedistribution.ToHitComparison {
				return &damagedistribution.ToHitComparison{
					ToHitBonus: 0,
				}
			},
		},
	})
	killSecondBanditCalculation := &powerattackforecastfakes.FakeCalculationInterface{}
	killSecondBanditCalculation.SetupReturns(&powerusagescenario.Setup{
		UserID:          suite.teros.ID(),
		PowerID:         suite.blot.ID(),
		Targets:         []string{suite.bandit.ID(), suite.bandit2.ID()},
		IsCounterAttack: false,
	})
	killSecondBanditCalculation.AttackReturns(&powerattackforecast.AttackForecast{
		DefenderContext: *powerattackforecast.NewDefenderContext(suite.bandit2.ID(), nil),
		VersusContext: &powerattackforecastfakes.FakeVersusContextStrategy{
			NormalDamageStub: func() *damagedistribution.DamageDistribution {
				return &damagedistribution.DamageDistribution{
					ActualDamageTaken: 2,
					IsFatalToTarget:   true,
				}
			},
			ToHitStub: func() *damagedistribution.ToHitComparison {
				return &damagedistribution.ToHitComparison{
					ToHitBonus: -2,
				}
			},
		},
	})

	forecastBlotOnMultipleBandits := &powerattackforecastfakes.FakeForecastInterface{}
	forecastBlotOnMultipleBandits.ForecastedResultPerTargetReturns([]powerattackforecast.CalculationInterface{
		killFirstBanditCalculation,
		killSecondBanditCalculation,
	})
	forecastBlotOnMultipleBandits.RepositoriesReturns(suite.repos)

	var forecastOutput strings.Builder
	suite.viewer.PrintForecast(
		forecastBlotOnMultipleBandits,
		forecastBlotOnMultipleBandits.Repositories(),
		&forecastOutput,
	)

	checker.Assert(forecastOutput.String(), Equals, "Teros (Blot) vs Bandit: +0 (21/36), FATAL\n- also vs Bandit2: -2 (10/36), FATAL\n")
}

func (suite *ConsoleShowsFatalAttacksSuite) TestIndicateIfItIsAKillingBlow(checker *C) {
	suite.bandit.ReduceHitPoints(suite.bandit.MaxHitPoints())
	resultBlotOnBanditKills := &powercommitfakes.FakeResultStrategy{}
	resultBlotOnBanditKills.ResultPerTargetReturns([]*powercommit.ResultPerTarget{
		powercommit.NewResultPerTargetBuilder().
			User(suite.teros).
			Power(suite.blot).
			Target(suite.bandit).
			AttackResult(
				powercommit.NewAttackResultBuilder().HitTarget().Build(),
			).
			Build(),
	})

	var output strings.Builder
	suite.viewer.PrintResult(resultBlotOnBanditKills, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Teros (Blot) hits Bandit, felling\n---\n")
}

type ConsoleShowsHealingAttempts struct {
	teros squaddieinterface.Interface
	lini  squaddieinterface.Interface

	healingStaff powerinterface.Interface

	viewer       *actionviewer.ConsoleActionViewer
	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection

	forecastHealingStaffOnTeros *powerattackforecast.Forecast
	resultHealingStaffOnTeros   *powercommit.Result
}

var _ = Suite(&ConsoleShowsHealingAttempts{})

func (suite *ConsoleShowsHealingAttempts) SetUpTest(checker *C) {
	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}
	suite.viewer = &actionviewer.ConsoleActionViewer{}

	suite.teros = squaddie.NewSquaddieBuilder().Teros().Build()

	suite.lini = squaddie.NewSquaddieBuilder().Lini().Mind(1).Build()
	suite.healingStaff = power.NewPowerBuilder().HealingStaff().WithName("healing Staff").Build()

	testutility.AddSquaddieWithInnatePowersToRepos(suite.teros, nil, suite.repos, false)
	testutility.AddSquaddieWithInnatePowersToRepos(suite.lini, suite.healingStaff, suite.repos, false)
}

func (suite *ConsoleShowsHealingAttempts) TestShowPowerHealingEffects(checker *C) {
	resultLiniHealsTeros := &powercommitfakes.FakeResultStrategy{}
	resultLiniHealsTeros.ResultPerTargetReturns([]*powercommit.ResultPerTarget{
		powercommit.NewResultPerTargetBuilder().
			User(suite.lini).
			Power(suite.healingStaff).
			Target(suite.teros).
			HealResult(
				powercommit.NewHealResultBuilder().HitPointsRestored(4).Build(),
			).
			Build(),
	})

	var output strings.Builder
	suite.viewer.PrintResult(resultLiniHealsTeros, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Lini (healing Staff) heals Teros, for 4 healing\n---\n")
}

type ConsoleShowsCounterAttackSuite struct {
	teros   squaddieinterface.Interface
	bandit  squaddieinterface.Interface
	bandit2 squaddieinterface.Interface
	bandit3 squaddieinterface.Interface

	blot         powerinterface.Interface
	criticalBlot powerinterface.Interface
	axe          powerinterface.Interface

	viewer       *actionviewer.ConsoleActionViewer
	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection
}

var _ = Suite(&ConsoleShowsCounterAttackSuite{})

func (suite *ConsoleShowsCounterAttackSuite) SetUpTest(checker *C) {
	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}
	suite.viewer = &actionviewer.ConsoleActionViewer{}

	suite.teros = squaddie.NewSquaddieBuilder().Teros().Mind(3).Aim(2).Build()
	suite.bandit = squaddie.NewSquaddieBuilder().Bandit().WithName("Bandit").WithID("banditID1").
		Strength(0).
		Barrier(1).
		Build()
	suite.bandit.SetBarrierToMax()

	suite.bandit2 = squaddie.NewSquaddieBuilder().Bandit().WithName("Bandit2").WithID("banditID2").
		HitPoints(2).
		Deflect(2).
		Barrier(20).
		Build()

	suite.bandit3 = squaddie.NewSquaddieBuilder().Bandit().WithID("banditID3").
		Barrier(1).
		Deflect(-200).
		Build()
	suite.bandit3.SetBarrierToMax()

	suite.axe = power.NewPowerBuilder().Axe().CanCounterAttack().DealsDamage(3).Build()
	suite.blot = power.NewPowerBuilder().Blot().WithName("Blot").DealsDamage(0).Build()
	suite.criticalBlot = power.NewPowerBuilder().CloneOf(suite.blot).WithID("Critical Blot").CriticalDealsDamage(1).CriticalHitThresholdBonus(1).Build()

	testutility.AddSquaddieWithInnatePowersToRepos(suite.teros, suite.blot, suite.repos, false)
	testutility.AddSquaddieWithInnatePowersToRepos(suite.teros, suite.criticalBlot, suite.repos, false)
	testutility.AddSquaddieWithInnatePowersToRepos(suite.bandit, suite.axe, suite.repos, false)
	testutility.AddSquaddieWithInnatePowersToRepos(suite.bandit2, suite.axe, suite.repos, true)
	testutility.AddSquaddieWithInnatePowersToRepos(suite.bandit3, suite.axe, suite.repos, false)
}

func (suite *ConsoleShowsCounterAttackSuite) TestShowForecastChanceToHit(checker *C) {
	attackSetup := &powerusagescenario.Setup{
		UserID:          suite.teros.ID(),
		PowerID:         suite.blot.ID(),
		Targets:         []string{suite.bandit.ID(), suite.bandit2.ID()},
		IsCounterAttack: false,
	}
	hitFirstBanditCalculation := &powerattackforecastfakes.FakeCalculationInterface{}
	hitFirstBanditCalculation.SetupReturns(attackSetup)
	hitFirstBanditCalculation.AttackReturns(&powerattackforecast.AttackForecast{
		DefenderContext: *powerattackforecast.NewDefenderContext(suite.bandit.ID(), nil),
		VersusContext: &powerattackforecastfakes.FakeVersusContextStrategy{
			NormalDamageStub: func() *damagedistribution.DamageDistribution {
				return &damagedistribution.DamageDistribution{
					RawDamageDealt:       2,
					TotalRawBarrierBurnt: 1,
				}
			},
			ToHitStub: func() *damagedistribution.ToHitComparison {
				return &damagedistribution.ToHitComparison{
					ToHitBonus: 2,
				}
			},
		},
	})
	hitSecondBanditCalculation := &powerattackforecastfakes.FakeCalculationInterface{}
	hitSecondBanditCalculation.SetupReturns(attackSetup)
	hitSecondBanditCalculation.AttackReturns(&powerattackforecast.AttackForecast{
		DefenderContext: *powerattackforecast.NewDefenderContext(suite.bandit2.ID(), nil),
		VersusContext: &powerattackforecastfakes.FakeVersusContextStrategy{
			NormalDamageStub: func() *damagedistribution.DamageDistribution {
				return &damagedistribution.DamageDistribution{
					TotalRawBarrierBurnt: 3,
				}
			},
			ToHitStub: func() *damagedistribution.ToHitComparison {
				return &damagedistribution.ToHitComparison{
					ToHitBonus: 0,
				}
			},
		},
	})

	secondBanditCountersCalculation := &powerattackforecastfakes.FakeCalculationInterface{}
	secondBanditCountersCalculation.CounterAttackSetupReturns(&powerusagescenario.Setup{
		UserID:          suite.bandit.ID(),
		PowerID:         suite.axe.ID(),
		Targets:         []string{suite.teros.ID()},
		IsCounterAttack: true,
	})
	secondBanditCountersCalculation.CounterAttackReturns(&powerattackforecast.AttackForecast{
		DefenderContext: *powerattackforecast.NewDefenderContext(suite.teros.ID(), nil),
		VersusContext: &powerattackforecastfakes.FakeVersusContextStrategy{
			NormalDamageStub: func() *damagedistribution.DamageDistribution {
				return &damagedistribution.DamageDistribution{
					RawDamageDealt: 3,
				}
			},
			ToHitStub: func() *damagedistribution.ToHitComparison {
				return &damagedistribution.ToHitComparison{
					ToHitBonus: -1,
				}
			},
		},
	})

	forecastBlotOnMultipleBandits := &powerattackforecastfakes.FakeForecastInterface{}
	forecastBlotOnMultipleBandits.ForecastedResultPerTargetReturns([]powerattackforecast.CalculationInterface{
		hitFirstBanditCalculation,
		hitSecondBanditCalculation,
		secondBanditCountersCalculation,
	})
	forecastBlotOnMultipleBandits.RepositoriesReturns(suite.repos)
	forecastBlotOnMultipleBandits.CalculateForecast()

	var forecastOutput strings.Builder
	suite.viewer.PrintForecast(
		forecastBlotOnMultipleBandits,
		forecastBlotOnMultipleBandits.Repositories(),
		&forecastOutput,
	)

	checker.Assert(forecastOutput.String(), Equals,
		"Teros (Blot) vs Bandit: +2 (30/36), for 2 damage + 1 barrier burn\n"+
			"- also vs Bandit2: +0 (21/36) for NO DAMAGE + 3 barrier burn\n"+
			"Bandit (axe) counters Teros: -1 (15/36), for 3 damage\n",
	)
}

func (suite *ConsoleShowsCounterAttackSuite) TestShowForecastChanceToCriticallyHit(checker *C) {
	attackSetup := &powerusagescenario.Setup{
		UserID:          suite.teros.ID(),
		PowerID:         suite.blot.ID(),
		Targets:         []string{suite.bandit3.ID()},
		IsCounterAttack: false,
	}
	hitFirstBanditCalculation := &powerattackforecastfakes.FakeCalculationInterface{}
	hitFirstBanditCalculation.SetupReturns(attackSetup)
	hitFirstBanditCalculation.AttackReturns(&powerattackforecast.AttackForecast{
		DefenderContext: *powerattackforecast.NewDefenderContext(suite.bandit3.ID(), nil),
		VersusContext: &powerattackforecastfakes.FakeVersusContextStrategy{
			NormalDamageStub: func() *damagedistribution.DamageDistribution {
				return &damagedistribution.DamageDistribution{
					RawDamageDealt:       2,
					TotalRawBarrierBurnt: 1,
				}
			},
			ToHitStub: func() *damagedistribution.ToHitComparison {
				return &damagedistribution.ToHitComparison{
					ToHitBonus: 202,
				}
			},
			CriticalHitThresholdStub: func() int { return 12 },
			CanCriticalStub:          func() bool { return true },
			CriticalHitDamageStub: func() *damagedistribution.DamageDistribution {
				return &damagedistribution.DamageDistribution{
					RawDamageDealt:       3,
					TotalRawBarrierBurnt: 1,
				}
			},
		},
	})

	forecastBlotOnMultipleBandits := &powerattackforecastfakes.FakeForecastInterface{}
	forecastBlotOnMultipleBandits.ForecastedResultPerTargetReturns([]powerattackforecast.CalculationInterface{
		hitFirstBanditCalculation,
	})
	forecastBlotOnMultipleBandits.RepositoriesReturns(suite.repos)
	forecastBlotOnMultipleBandits.CalculateForecast()

	var forecastOutput strings.Builder
	suite.viewer.PrintForecast(
		forecastBlotOnMultipleBandits,
		forecastBlotOnMultipleBandits.Repositories(),
		&forecastOutput,
	)

	checker.Assert(forecastOutput.String(), Equals, "Teros (Blot) vs Bandit: +202 (36/36), for 2 damage + 1 barrier burn\n"+" crit: 36/36, for 3 damage + 1 barrier burn\n")
}

func (suite *ConsoleShowsCounterAttackSuite) TestShowCounterattackResults(checker *C) {
	resultBlotOnBanditAndCounterAxeOnTeros := &powercommitfakes.FakeResultStrategy{}
	resultBlotOnBanditAndCounterAxeOnTeros.ResultPerTargetReturns([]*powercommit.ResultPerTarget{
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
	})

	var output strings.Builder
	suite.viewer.PrintResult(resultBlotOnBanditAndCounterAxeOnTeros, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Teros (Blot) hits Bandit, for 0 damage\nBandit (axe) counters Teros, for 2 damage\n---\n")
}

func (suite *ConsoleShowsCounterAttackSuite) TestShowMultipleTargets(checker *C) {
	resultBlotOnBanditsAndBandit2Counters := &powercommitfakes.FakeResultStrategy{}
	resultBlotOnBanditsAndBandit2Counters.ResultPerTargetReturns([]*powercommit.ResultPerTarget{
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
	})

	var output strings.Builder
	suite.viewer.PrintResult(resultBlotOnBanditsAndBandit2Counters, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals,
		"Teros (Blot) hits Bandit, for 2 damage + 1 barrier burn\n"+
			"- also hits Bandit2, for 3 damage\n"+
			"Bandit2 (axe) counters Teros, for 3 damage\n"+
			"---\n",
	)
}

type ConsoleShowsHitsAndMisses struct {
	teros  squaddieinterface.Interface
	bandit squaddieinterface.Interface

	blot powerinterface.Interface

	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection

	viewer *actionviewer.ConsoleActionViewer
}

var _ = Suite(&ConsoleShowsHitsAndMisses{})

func (suite *ConsoleShowsHitsAndMisses) SetUpTest(checker *C) {
	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}
	suite.viewer = &actionviewer.ConsoleActionViewer{}

	suite.teros = squaddie.NewSquaddieBuilder().Teros().Build()
	suite.bandit = squaddie.NewSquaddieBuilder().Bandit().Build()

	suite.blot = power.NewPowerBuilder().Blot().WithName("Blot").DealsDamage(0).Build()

	testutility.AddSquaddieWithInnatePowersToRepos(suite.teros, suite.blot, suite.repos, false)
	testutility.AddSquaddieWithInnatePowersToRepos(suite.bandit, nil, suite.repos, false)
}

func (suite *ConsoleShowsHitsAndMisses) TestShowPowerHitTargetAndDamage(checker *C) {
	resultBlotOnBanditHit := &powercommitfakes.FakeResultStrategy{}
	resultBlotOnBanditHit.ResultPerTargetReturns([]*powercommit.ResultPerTarget{
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
	})

	var output strings.Builder
	suite.viewer.PrintResult(resultBlotOnBanditHit, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Teros (Blot) hits Bandit, for 3 damage\n---\n")
}

func (suite *ConsoleShowsHitsAndMisses) TestShowWhenPowerMisses(checker *C) {
	resultBlotOnBanditMissed := &powercommitfakes.FakeResultStrategy{}
	resultBlotOnBanditMissed.ResultPerTargetReturns([]*powercommit.ResultPerTarget{
		powercommit.NewResultPerTargetBuilder().
			User(suite.teros).
			Power(suite.blot).
			Target(suite.bandit).
			AttackResult(
				powercommit.NewAttackResultBuilder().Build(),
			).
			Build(),
	})

	var output strings.Builder
	suite.viewer.PrintResult(resultBlotOnBanditMissed, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Teros (Blot) misses Bandit\n---\n")
}

func (suite *ConsoleShowsHitsAndMisses) TestShowWhenPowerCriticallyHits(checker *C) {
	resultBlotOnBanditCriticallyHit := &powercommitfakes.FakeResultStrategy{}
	resultBlotOnBanditCriticallyHit.ResultPerTargetReturns([]*powercommit.ResultPerTarget{
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
	})

	var output strings.Builder
	suite.viewer.PrintResult(resultBlotOnBanditCriticallyHit, suite.repos, nil, &output)
	checker.Assert(output.String(), Equals, "Teros (Blot) CRITICALLY hits Bandit, for 4 damage\n---\n")
}

func (suite *ConsoleShowsHitsAndMisses) TestShowPowerBarrierBurn(checker *C) {
	resultBlotOnBanditBurnsBarrier := &powercommitfakes.FakeResultStrategy{}
	resultBlotOnBanditBurnsBarrier.ResultPerTargetReturns([]*powercommit.ResultPerTarget{
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
	})

	var output strings.Builder
	suite.viewer.PrintResult(resultBlotOnBanditBurnsBarrier, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Teros (Blot) hits Bandit, for 2 damage + 1 barrier burn\n---\n")
}

type ConsoleShowsVerbosity struct {
	teros             squaddieinterface.Interface
	lini              squaddieinterface.Interface
	bandit            squaddieinterface.Interface
	banditWithBarrier squaddieinterface.Interface
	bandit2           squaddieinterface.Interface

	blot         powerinterface.Interface
	axe          powerinterface.Interface
	healingStaff powerinterface.Interface

	viewer       *actionviewer.ConsoleActionViewer
	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection
}

var _ = Suite(&ConsoleShowsVerbosity{})

func (suite *ConsoleShowsVerbosity) SetUpTest(checker *C) {
	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}
	suite.viewer = &actionviewer.ConsoleActionViewer{}

	suite.teros = squaddie.NewSquaddieBuilder().Teros().Build()
	suite.blot = power.NewPowerBuilder().Blot().WithName("Blot").DealsDamage(0).Build()

	suite.bandit = squaddie.NewSquaddieBuilder().Bandit().Build()
	suite.bandit2 = squaddie.NewSquaddieBuilder().Bandit().WithName("Bandit2").WithID("banditID2").Build()
	suite.banditWithBarrier = squaddie.NewSquaddieBuilder().WithName("Bandit").Barrier(2).Build()
	suite.axe = power.NewPowerBuilder().Axe().Build()

	suite.lini = squaddie.NewSquaddieBuilder().Lini().Mind(1).Build()
	suite.healingStaff = power.NewPowerBuilder().HealingStaff().WithName("healing Staff").Build()

	testutility.AddSquaddieWithInnatePowersToRepos(suite.teros, suite.blot, suite.repos, false)
	testutility.AddSquaddieWithInnatePowersToRepos(suite.lini, suite.healingStaff, suite.repos, false)
	testutility.AddSquaddieWithInnatePowersToRepos(suite.bandit, suite.axe, suite.repos, false)
	testutility.AddSquaddieWithInnatePowersToRepos(suite.bandit2, suite.axe, suite.repos, false)
	testutility.AddSquaddieWithInnatePowersToRepos(suite.banditWithBarrier, suite.axe, suite.repos, false)
}

func (suite *ConsoleShowsVerbosity) TestShowTargetStatusVerbosity(checker *C) {
	resultBlotOnBanditsAndBandit2Counters := &powercommitfakes.FakeResultStrategy{}
	resultBlotOnBanditsAndBandit2Counters.ResultPerTargetReturns([]*powercommit.ResultPerTarget{
		powercommit.NewResultPerTargetBuilder().
			User(suite.teros).
			Power(suite.blot).
			Target(suite.banditWithBarrier).
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
	})

	suite.banditWithBarrier.ReduceHitPoints(2)
	suite.banditWithBarrier.ReduceBarrier(1)
	suite.bandit2.ReduceHitPoints(3)
	suite.teros.ReduceHitPoints(3)

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
}

func (suite *ConsoleShowsVerbosity) TestShowHealingVerbosity(checker *C) {
	healResult := &powercommitfakes.FakeResultStrategy{}
	healResult.ResultPerTargetReturns([]*powercommit.ResultPerTarget{
		powercommit.NewResultPerTargetBuilder().
			User(suite.lini).
			Power(suite.healingStaff).
			Target(suite.teros).
			HealResult(
				powercommit.NewHealResultBuilder().HitPointsRestored(4).Build(),
			).
			Build(),
	})
	suite.teros.SetHPToMax()

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

func (suite *ConsoleShowsVerbosity) TestShowRollsVerbosity(checker *C) {
	resultBlotOnBanditsAndBandit2Counters := &powercommitfakes.FakeResultStrategy{}
	resultBlotOnBanditsAndBandit2Counters.ResultPerTargetReturns([]*powercommit.ResultPerTarget{
		powercommit.NewResultPerTargetBuilder().
			User(suite.teros).
			Power(suite.blot).
			Target(suite.banditWithBarrier).
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
	})

	suite.banditWithBarrier.ReduceHitPoints(2)
	suite.banditWithBarrier.ReduceBarrier(1)
	suite.bandit2.ReduceHitPoints(3)
	suite.teros.ReduceHitPoints(3)

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
}

func (suite *ConsoleShowsVerbosity) TestShowHealingAutoHitVerbosity(checker *C) {
	healResult := &powercommitfakes.FakeResultStrategy{}
	healResult.ResultPerTargetReturns([]*powercommit.ResultPerTarget{
		powercommit.NewResultPerTargetBuilder().
			User(suite.lini).
			Power(suite.healingStaff).
			Target(suite.teros).
			HealResult(
				powercommit.NewHealResultBuilder().HitPointsRestored(4).Build(),
			).
			Build(),
	})
	suite.teros.SetHPToMax()
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
