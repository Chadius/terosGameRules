package powercommit_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerinterface"
	"github.com/chadius/terosbattleserver/entity/powerreference"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	"github.com/chadius/terosbattleserver/entity/powerusagescenario"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/squaddieinterface"
	"github.com/chadius/terosbattleserver/usecase/powerattackforecast"
	"github.com/chadius/terosbattleserver/usecase/powercommit"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/usecase/squaddiestats"
	"github.com/chadius/terosbattleserver/utility/testutility"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type resultOnAttack struct {
	teros      squaddieinterface.Interface
	bandit     squaddieinterface.Interface
	bandit2    squaddieinterface.Interface
	mysticMage squaddieinterface.Interface

	spear    powerinterface.Interface
	blot     powerinterface.Interface
	fireball powerinterface.Interface
	axe      powerinterface.Interface

	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection

	forecastSpearOnBandit *powerattackforecast.Forecast
	resultSpearOnBandit   *powercommit.Result

	forecastBlotOnBandit *powerattackforecast.Forecast
	resultBlotOnBandit   *powercommit.Result

	equipCheck powerequip.Strategy
}

var _ = Suite(&resultOnAttack{})

func (suite *resultOnAttack) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Build()
	suite.mysticMage = squaddie.NewSquaddieBuilder().MysticMage().Build()
	suite.bandit = squaddie.NewSquaddieBuilder().Bandit().Build()
	suite.bandit2 = squaddie.NewSquaddieBuilder().Bandit().WithID("bandit2ID").WithName("bandit2").Build()

	suite.spear = power.NewPowerBuilder().Spear().Build()
	suite.blot = power.NewPowerBuilder().Blot().Build()
	suite.axe = power.NewPowerBuilder().Axe().Build()
	suite.fireball = power.NewPowerBuilder().IsSpell().DealsDamage(3).Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]squaddieinterface.Interface{suite.teros, suite.bandit, suite.bandit2, suite.mysticMage})

	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]powerinterface.Interface{suite.spear, suite.blot, suite.axe, suite.fireball})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}

	checkEquip := powerequip.CheckRepositories{}
	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*powerreference.Reference{
			suite.spear.GetReference(),
			suite.blot.GetReference(),
		},
		suite.repos,
	)

	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.bandit,
		[]*powerreference.Reference{
			suite.axe.GetReference(),
		},
		suite.repos,
	)

	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.bandit2,
		[]*powerreference.Reference{
			suite.axe.GetReference(),
		},
		suite.repos,
	)

	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.mysticMage,
		[]*powerreference.Reference{
			suite.fireball.GetReference(),
		},
		suite.repos,
	)

	suite.CalculateSpearOnBandit(nil)
	suite.resultSpearOnBandit = powercommit.NewResult(suite.forecastSpearOnBandit, nil, nil)

	suite.CalculateBlotOnBandit(nil)
	suite.resultBlotOnBandit = powercommit.NewResult(suite.forecastBlotOnBandit, nil, nil)

	suite.equipCheck = &powerequip.CheckRepositories{}
}

func (suite *resultOnAttack) CalculateSpearOnBandit(setup *powerusagescenario.Setup) {
	setupToUse := powerusagescenario.Setup{
		UserID:          suite.teros.ID(),
		PowerID:         suite.spear.ID(),
		Targets:         []string{suite.bandit.ID()},
		IsCounterAttack: false,
	}
	if setup != nil {
		setupToUse = *setup
	}

	repos := &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}

	suite.forecastSpearOnBandit = powerattackforecast.NewForecastBuilder().
		Setup(&setupToUse).
		Repositories(repos).
		OffenseStrategy(&squaddiestats.CalculateSquaddieOffenseStats{}).
		Build()

	suite.forecastSpearOnBandit.CalculateForecast()
}

func (suite *resultOnAttack) CalculateBlotOnBandit(setup *powerusagescenario.Setup) {
	setupToUse := powerusagescenario.Setup{
		UserID:          suite.teros.ID(),
		PowerID:         suite.blot.ID(),
		Targets:         []string{suite.bandit.ID()},
		IsCounterAttack: false,
	}
	if setup != nil {
		setupToUse = *setup
	}

	repos := &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}

	suite.forecastBlotOnBandit = powerattackforecast.NewForecastBuilder().
		Setup(&setupToUse).
		Repositories(repos).
		OffenseStrategy(&squaddiestats.CalculateSquaddieOffenseStats{}).
		Build()

	suite.forecastBlotOnBandit.CalculateForecast()
}

func (suite *resultOnAttack) TestAttackCanMiss(checker *C) {
	suite.CalculateSpearOnBandit(nil)

	resultSpearOnBanditAlwaysMisses := suite.resultSpearOnBandit.CopyResultWithNewDieRoller(&testutility.AlwaysMissDieRoller{})
	resultSpearOnBanditAlwaysMisses.Commit()

	checker.Assert(resultSpearOnBanditAlwaysMisses.ResultPerTarget(), HasLen, 1)
	checker.Assert(resultSpearOnBanditAlwaysMisses.ResultPerTarget()[0].UserID(), Equals, suite.teros.ID())
	checker.Assert(resultSpearOnBanditAlwaysMisses.ResultPerTarget()[0].PowerID(), Equals, suite.spear.ID())
	checker.Assert(resultSpearOnBanditAlwaysMisses.ResultPerTarget()[0].TargetID(), Equals, suite.bandit.ID())
	checker.Assert(resultSpearOnBanditAlwaysMisses.ResultPerTarget()[0].Attack().HitTarget(), Equals, false)
	checker.Assert(resultSpearOnBanditAlwaysMisses.ResultPerTarget()[0].Attack().CriticallyHitTarget(), Equals, false)
	checker.Assert(resultSpearOnBanditAlwaysMisses.ResultPerTarget()[0].Attack().Damage().DamageAbsorbedByBarrier, Equals, 0)
	checker.Assert(resultSpearOnBanditAlwaysMisses.ResultPerTarget()[0].Attack().Damage().DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(resultSpearOnBanditAlwaysMisses.ResultPerTarget()[0].Attack().Damage().RawDamageDealt, Equals, 0)
}

func (suite *resultOnAttack) TestAttackCanHitButNotCritically(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Mind(2).Build()
	suite.squaddieRepo.AddSquaddie(suite.teros)

	suite.bandit = squaddie.NewSquaddieBuilder().Bandit().Armor(1).Barrier(3).Build()
	suite.bandit.SetBarrierToMax()
	suite.squaddieRepo.AddSquaddie(suite.bandit)

	suite.CalculateBlotOnBandit(&powerusagescenario.Setup{
		UserID:          suite.teros.ID(),
		PowerID:         suite.blot.ID(),
		Targets:         []string{suite.bandit.ID()},
		IsCounterAttack: false,
	})
	resultBlotOnBanditAlwaysHits := powercommit.NewResult(suite.forecastBlotOnBandit, testutility.AlwaysHitDieRoller{}, nil)
	resultBlotOnBanditAlwaysHits.Commit()

	checker.Assert(resultBlotOnBanditAlwaysHits.ResultPerTarget()[0].PowerID(), Equals, suite.blot.ID())
	checker.Assert(resultBlotOnBanditAlwaysHits.ResultPerTarget()[0].Attack().HitTarget(), Equals, true)
	checker.Assert(resultBlotOnBanditAlwaysHits.ResultPerTarget()[0].Attack().CriticallyHitTarget(), Equals, false)
	checker.Assert(resultBlotOnBanditAlwaysHits.ResultPerTarget()[0].Attack().Damage().DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(resultBlotOnBanditAlwaysHits.ResultPerTarget()[0].Attack().Damage().DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(resultBlotOnBanditAlwaysHits.ResultPerTarget()[0].Attack().Damage().RawDamageDealt, Equals, 2)

	checker.Assert(
		suite.bandit.CurrentHitPoints(),
		Equals,
		suite.bandit.MaxHitPoints()-resultBlotOnBanditAlwaysHits.ResultPerTarget()[0].Attack().Damage().RawDamageDealt,
	)
}

func (suite *resultOnAttack) TestAttackCanHitCritically(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Mind(2).Build()
	suite.squaddieRepo.AddSquaddie(suite.teros)
	suite.blot = power.NewPowerBuilder().CloneOf(suite.blot).WithID(suite.blot.ID()).DealsDamage(3).CriticalDealsDamage(3).CriticalHitThresholdBonus(9000).Build()
	suite.powerRepo.AddPower(suite.blot)
	checkEquip := powerequip.CheckRepositories{}
	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*powerreference.Reference{
			suite.blot.GetReference(),
		},
		suite.repos,
	)

	suite.bandit = squaddie.NewSquaddieBuilder().HitPoints(1).Armor(1).Barrier(3).Build()
	suite.bandit.SetBarrierToMax()
	suite.bandit.SetHPToMax()
	suite.squaddieRepo.AddSquaddie(suite.bandit)

	suite.CalculateBlotOnBandit(&powerusagescenario.Setup{
		UserID:          suite.teros.ID(),
		PowerID:         suite.blot.ID(),
		Targets:         []string{suite.bandit.ID()},
		IsCounterAttack: false,
	})

	resultBlotOnBanditAlwaysHits := powercommit.NewResult(suite.forecastBlotOnBandit, testutility.AlwaysHitDieRoller{}, nil)
	resultBlotOnBanditAlwaysHits.Commit()

	checker.Assert(resultBlotOnBanditAlwaysHits.ResultPerTarget()[0].PowerID(), Equals, suite.blot.ID())
	checker.Assert(resultBlotOnBanditAlwaysHits.ResultPerTarget()[0].Attack().HitTarget(), Equals, true)
	checker.Assert(resultBlotOnBanditAlwaysHits.ResultPerTarget()[0].Attack().CriticallyHitTarget(), Equals, true)
	checker.Assert(resultBlotOnBanditAlwaysHits.ResultPerTarget()[0].Attack().Damage().DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(resultBlotOnBanditAlwaysHits.ResultPerTarget()[0].Attack().Damage().DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(resultBlotOnBanditAlwaysHits.ResultPerTarget()[0].Attack().Damage().RawDamageDealt, Equals, 5)

	checker.Assert(
		suite.bandit.CurrentHitPoints(),
		Equals,
		0,
	)
}

func (suite *resultOnAttack) TestCounterAttacks(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Armor(0).Barrier(0).Strength(2).Build()
	suite.squaddieRepo.AddSquaddie(suite.teros)

	suite.spear = power.NewPowerBuilder().CloneOf(suite.spear).WithID(suite.spear.ID()).DealsDamage(3).Build()
	suite.axe = power.NewPowerBuilder().CloneOf(suite.axe).WithID(suite.axe.ID()).CanCounterAttack().DealsDamage(3).Build()
	suite.powerRepo.AddSlicePowerSource([]powerinterface.Interface{suite.spear, suite.axe})

	suite.bandit = squaddie.NewSquaddieBuilder().Bandit().Strength(0).Armor(1).Build()
	suite.squaddieRepo.AddSquaddie(suite.bandit)
	checkEquip := powerequip.CheckRepositories{}
	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.bandit,
		[]*powerreference.Reference{
			suite.axe.GetReference(),
		},
		suite.repos,
	)
	checkEquip.SquaddieEquipPower(suite.bandit, suite.axe.ID(), suite.repos)

	suite.CalculateSpearOnBandit(&powerusagescenario.Setup{
		UserID:          suite.teros.ID(),
		PowerID:         suite.spear.ID(),
		Targets:         []string{suite.bandit.ID()},
		IsCounterAttack: false,
	})

	suite.resultSpearOnBandit = powercommit.NewResult(suite.forecastSpearOnBandit, nil, nil)
	resultSpearOnBanditAlwaysHits := suite.resultSpearOnBandit.CopyResultWithNewDieRoller(&testutility.AlwaysHitDieRoller{})
	resultSpearOnBanditAlwaysHits.Commit()

	checker.Assert(resultSpearOnBanditAlwaysHits.ResultPerTarget()[0].PowerID(), Equals, suite.spear.ID())
	checker.Assert(resultSpearOnBanditAlwaysHits.ResultPerTarget()[0].UserID(), Equals, suite.teros.ID())
	checker.Assert(resultSpearOnBanditAlwaysHits.ResultPerTarget()[0].TargetID(), Equals, suite.bandit.ID())

	checker.Assert(resultSpearOnBanditAlwaysHits.ResultPerTarget()[1].PowerID(), Equals, suite.axe.ID())
	checker.Assert(resultSpearOnBanditAlwaysHits.ResultPerTarget()[1].UserID(), Equals, suite.bandit.ID())
	checker.Assert(resultSpearOnBanditAlwaysHits.ResultPerTarget()[1].TargetID(), Equals, suite.teros.ID())
	checker.Assert(resultSpearOnBanditAlwaysHits.ResultPerTarget()[1].Attack().HitTarget(), Equals, true)
	checker.Assert(resultSpearOnBanditAlwaysHits.ResultPerTarget()[1].Attack().IsCounterAttack(), Equals, true)
	checker.Assert(resultSpearOnBanditAlwaysHits.ResultPerTarget()[1].Attack().CriticallyHitTarget(), Equals, false)
	checker.Assert(resultSpearOnBanditAlwaysHits.ResultPerTarget()[1].Attack().Damage().DamageAbsorbedByBarrier, Equals, 0)
	checker.Assert(resultSpearOnBanditAlwaysHits.ResultPerTarget()[1].Attack().Damage().DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(resultSpearOnBanditAlwaysHits.ResultPerTarget()[1].Attack().Damage().RawDamageDealt, Equals, 3)

	checker.Assert(
		suite.teros.CurrentHitPoints(),
		Equals,
		suite.teros.MaxHitPoints()-resultSpearOnBanditAlwaysHits.ResultPerTarget()[1].Attack().Damage().RawDamageDealt,
	)
}

func (suite *resultOnAttack) TestCounterAttacksApplyLast(checker *C) {
	suite.bandit = squaddie.NewSquaddieBuilder().HitPoints(suite.fireball.DamageBonus() + suite.mysticMage.Mind() + 1).Build()
	suite.bandit.SetHPToMax()
	suite.squaddieRepo.AddSquaddie(suite.bandit)
	checkEquip := powerequip.CheckRepositories{}
	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.bandit,
		[]*powerreference.Reference{
			suite.axe.GetReference(),
		},
		suite.repos,
	)
	checkEquip.SquaddieEquipPower(suite.bandit, suite.axe.ID(), suite.repos)

	suite.bandit2 = squaddie.NewSquaddieBuilder().HitPoints(suite.fireball.DamageBonus() + suite.mysticMage.Mind() + 1).Build()
	suite.squaddieRepo.AddSquaddie(suite.bandit2)
	suite.bandit2.SetHPToMax()
	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.bandit2,
		[]*powerreference.Reference{
			suite.axe.GetReference(),
		},
		suite.repos,
	)
	checkEquip.SquaddieEquipPower(suite.bandit2, suite.axe.ID(), suite.repos)

	forecastFireballOnBandits := powerattackforecast.NewForecastBuilder().
		Setup(
			&powerusagescenario.Setup{
				UserID:          suite.mysticMage.ID(),
				PowerID:         suite.fireball.ID(),
				Targets:         []string{suite.bandit.ID(), suite.bandit2.ID()},
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

	forecastFireballOnBandits.CalculateForecast()
	resultFireballOnBanditsAlwaysHits := powercommit.NewResult(forecastFireballOnBandits, &testutility.AlwaysHitDieRoller{}, nil)
	resultFireballOnBanditsAlwaysHits.Commit()

	checker.Assert(resultFireballOnBanditsAlwaysHits.ResultPerTarget()[0].PowerID(), Equals, suite.fireball.ID())
	checker.Assert(resultFireballOnBanditsAlwaysHits.ResultPerTarget()[0].UserID(), Equals, suite.mysticMage.ID())
	checker.Assert(resultFireballOnBanditsAlwaysHits.ResultPerTarget()[0].TargetID(), Equals, suite.bandit.ID())

	checker.Assert(resultFireballOnBanditsAlwaysHits.ResultPerTarget()[1].PowerID(), Equals, suite.fireball.ID())
	checker.Assert(resultFireballOnBanditsAlwaysHits.ResultPerTarget()[1].UserID(), Equals, suite.mysticMage.ID())
	checker.Assert(resultFireballOnBanditsAlwaysHits.ResultPerTarget()[1].TargetID(), Equals, suite.bandit2.ID())

	checker.Assert(resultFireballOnBanditsAlwaysHits.ResultPerTarget()[2].PowerID(), Equals, suite.axe.ID())
	checker.Assert(resultFireballOnBanditsAlwaysHits.ResultPerTarget()[2].UserID(), Equals, suite.bandit.ID())
	checker.Assert(resultFireballOnBanditsAlwaysHits.ResultPerTarget()[2].TargetID(), Equals, suite.mysticMage.ID())

	checker.Assert(resultFireballOnBanditsAlwaysHits.ResultPerTarget()[3].PowerID(), Equals, suite.axe.ID())
	checker.Assert(resultFireballOnBanditsAlwaysHits.ResultPerTarget()[3].UserID(), Equals, suite.bandit2.ID())
	checker.Assert(resultFireballOnBanditsAlwaysHits.ResultPerTarget()[3].TargetID(), Equals, suite.mysticMage.ID())
}

func (suite *resultOnAttack) TestDeadSquaddiesCannotCounterAttack(checker *C) {
	checkEquip := powerequip.CheckRepositories{}
	suite.teros = squaddie.NewSquaddieBuilder().WithID("squaddieTeros").Armor(0).Barrier(0).Strength(suite.bandit.MaxHitPoints()).Build()
	suite.squaddieRepo.AddSquaddie(suite.teros)
	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*powerreference.Reference{
			suite.spear.GetReference(),
		},
		suite.repos,
	)
	checkEquip.SquaddieEquipPower(suite.teros, suite.spear.ID(), suite.repos)

	suite.bandit = squaddie.NewSquaddieBuilder().WithID("squaddieBandit").Armor(0).Strength(0).Build()
	suite.squaddieRepo.AddSquaddie(suite.bandit)
	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.bandit,
		[]*powerreference.Reference{
			suite.axe.GetReference(),
		},
		suite.repos,
	)

	suite.CalculateSpearOnBandit(&powerusagescenario.Setup{
		UserID:          suite.teros.ID(),
		PowerID:         suite.spear.ID(),
		Targets:         []string{suite.bandit.ID()},
		IsCounterAttack: false,
	})

	resultSpearOnBanditAlwaysHits := powercommit.NewResult(suite.forecastSpearOnBandit, &testutility.AlwaysHitDieRoller{}, nil)
	resultSpearOnBanditAlwaysHits.Commit()

	checker.Assert(resultSpearOnBanditAlwaysHits.ResultPerTarget()[0].PowerID(), Equals, suite.spear.ID())
	checker.Assert(resultSpearOnBanditAlwaysHits.ResultPerTarget()[0].UserID(), Equals, suite.teros.ID())
	checker.Assert(resultSpearOnBanditAlwaysHits.ResultPerTarget()[0].TargetID(), Equals, suite.bandit.ID())

	checker.Assert(suite.bandit.IsDead(), Equals, true)
	checker.Assert(resultSpearOnBanditAlwaysHits.ResultPerTarget(), HasLen, 1)
}

type EquipPowerWhenCommitting struct {
	teros      squaddieinterface.Interface
	bandit     squaddieinterface.Interface
	mysticMage squaddieinterface.Interface

	spear    powerinterface.Interface
	blot     powerinterface.Interface
	fireball powerinterface.Interface

	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection

	forecastSpearOnBandit *powerattackforecast.Forecast
	resultSpearOnBandit   *powercommit.Result

	forecastBlotOnBandit *powerattackforecast.Forecast
	resultBlotOnBandit   *powercommit.Result

	forecastFireballOnBandit *powerattackforecast.Forecast
	resultFireballOnBandit   *powercommit.Result

	equipCheck powerequip.Strategy
}

var _ = Suite(&EquipPowerWhenCommitting{})

func (suite *EquipPowerWhenCommitting) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Build()
	suite.mysticMage = squaddie.NewSquaddieBuilder().MysticMage().Build()
	suite.bandit = squaddie.NewSquaddieBuilder().Bandit().Build()

	suite.spear = power.NewPowerBuilder().Spear().Build()
	suite.blot = power.NewPowerBuilder().Blot().CannotBeEquipped().Build()
	suite.fireball = power.NewPowerBuilder().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]squaddieinterface.Interface{suite.teros, suite.bandit, suite.mysticMage})

	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]powerinterface.Interface{suite.spear, suite.blot, suite.fireball})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}

	checkEquip := powerequip.CheckRepositories{}
	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*powerreference.Reference{
			suite.spear.GetReference(),
			suite.blot.GetReference(),
		},
		suite.repos,
	)

	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.mysticMage,
		[]*powerreference.Reference{
			suite.fireball.GetReference(),
		},
		suite.repos,
	)

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

	suite.resultSpearOnBandit = powercommit.NewResult(suite.forecastSpearOnBandit, nil, nil)

	suite.forecastBlotOnBandit = powerattackforecast.NewForecastBuilder().
		Setup(
			&powerusagescenario.Setup{
				UserID:          suite.teros.ID(),
				PowerID:         suite.blot.ID(),
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

	suite.resultBlotOnBandit = powercommit.NewResult(suite.forecastBlotOnBandit, nil, nil)

	suite.forecastFireballOnBandit = powerattackforecast.NewForecastBuilder().
		Setup(
			&powerusagescenario.Setup{
				UserID:          suite.mysticMage.ID(),
				PowerID:         suite.fireball.ID(),
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

	suite.resultFireballOnBandit = powercommit.NewResult(suite.forecastFireballOnBandit, nil, nil)
	suite.equipCheck = &powerequip.CheckRepositories{}
}

func (suite *EquipPowerWhenCommitting) TestCommitWillTryToEquipPower(checker *C) {
	resultSpearOnBanditAlwaysMisses := suite.resultSpearOnBandit.CopyResultWithNewDieRoller(&testutility.AlwaysMissDieRoller{})
	suite.forecastSpearOnBandit.CalculateForecast()

	resultSpearOnBanditAlwaysMisses.Commit()

	checker.Assert(suite.teros.GetEquippedPowerID(), Equals, suite.spear.ID())
}

func (suite *EquipPowerWhenCommitting) TestSquaddieWillKeepPreviousPowerIfCommitPowerCannotBeEquipped(checker *C) {
	suite.equipCheck.SquaddieEquipPower(suite.teros, suite.spear.ID(), suite.repos)
	resultBlotOnBanditAlwaysMisses := suite.resultBlotOnBandit.CopyResultWithNewDieRoller(&testutility.AlwaysMissDieRoller{})

	suite.forecastBlotOnBandit.CalculateForecast()
	resultBlotOnBanditAlwaysMisses.Commit()

	checker.Assert(suite.teros.HasEquippedPower(), Equals, true)
	checker.Assert(suite.teros.GetEquippedPowerID(), Equals, suite.spear.ID())
}

func (suite *EquipPowerWhenCommitting) TestSquaddieWillNotEquipPowerIfNoneExistAfterCommitting(checker *C) {
	resultFireballOnBanditAlwaysMisses := suite.resultFireballOnBandit.CopyResultWithNewDieRoller(&testutility.AlwaysMissDieRoller{})

	suite.forecastFireballOnBandit.CalculateForecast()
	resultFireballOnBanditAlwaysMisses.Commit()

	checker.Assert(suite.mysticMage.HasEquippedPower(), Equals, false)
}

type ResultOnHealing struct {
	lini  squaddieinterface.Interface
	teros squaddieinterface.Interface
	vale  squaddieinterface.Interface

	healingStaff powerinterface.Interface

	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection

	forecastHealingStaffOnTeros *powerattackforecast.Forecast
	resultHealingStaffOnTeros   *powercommit.Result

	forecastHealingStaffOnTerosAndVale *powerattackforecast.Forecast
	resultHealingStaffOnTerosAndVale   *powercommit.Result

	equipCheck powerequip.Strategy
}

var _ = Suite(&ResultOnHealing{})

func (suite *ResultOnHealing) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Build()
	suite.lini = squaddie.NewSquaddieBuilder().Lini().Build()
	suite.vale = squaddie.NewSquaddieBuilder().WithName("Vale").AsPlayer().Build()

	suite.healingStaff = power.NewPowerBuilder().HealingStaff().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]squaddieinterface.Interface{suite.teros, suite.lini, suite.vale})

	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]powerinterface.Interface{suite.healingStaff})

	suite.repos = &repositories.RepositoryCollection{PowerRepo: suite.powerRepo, SquaddieRepo: suite.squaddieRepo}
}

func (suite *ResultOnHealing) TestHealResultShowsHitPointsRestored(checker *C) {
	suite.teros.SetHPToMax()
	suite.teros.ReduceHitPoints(suite.teros.MaxHitPoints() - 1)
	suite.lini = squaddie.NewSquaddieBuilder().Mind(1).Build()
	suite.squaddieRepo.AddSquaddie(suite.lini)

	suite.forecastHealingStaffOnTeros = powerattackforecast.NewForecastBuilder().
		Setup(
			&powerusagescenario.Setup{
				UserID:          suite.lini.ID(),
				PowerID:         suite.healingStaff.ID(),
				Targets:         []string{suite.teros.ID()},
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
	suite.forecastHealingStaffOnTeros.CalculateForecast()

	suite.resultHealingStaffOnTeros = powercommit.NewResult(suite.forecastHealingStaffOnTeros, nil, nil)
	suite.resultHealingStaffOnTeros.Commit()

	checker.Assert(suite.resultHealingStaffOnTeros.ResultPerTarget()[0].PowerID(), Equals, suite.healingStaff.ID())
	checker.Assert(suite.resultHealingStaffOnTeros.ResultPerTarget()[0].Healing().HitPointsRestored(), Equals, 4)

	checker.Assert(
		suite.teros.CurrentHitPoints(),
		Equals,
		1+suite.resultHealingStaffOnTeros.ResultPerTarget()[0].Healing().HitPointsRestored(),
	)
}

func (suite *ResultOnHealing) TestHealResultShowsForEachTarget(checker *C) {
	suite.teros.SetHPToMax()
	suite.teros.ReduceHitPoints(suite.teros.MaxHitPoints() - 1)
	suite.lini = squaddie.NewSquaddieBuilder().Mind(1).Build()
	suite.squaddieRepo.AddSquaddie(suite.lini)
	suite.vale.SetHPToMax()

	suite.vale.ReduceHitPoints(suite.teros.MaxHitPoints() - 2)

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
	suite.forecastHealingStaffOnTerosAndVale.CalculateForecast()

	suite.resultHealingStaffOnTerosAndVale = powercommit.NewResult(suite.forecastHealingStaffOnTerosAndVale, nil, nil)
	suite.resultHealingStaffOnTerosAndVale.Commit()

	checker.Assert(suite.resultHealingStaffOnTerosAndVale.ResultPerTarget(), HasLen, 2)
	checker.Assert(suite.resultHealingStaffOnTerosAndVale.ResultPerTarget()[0].PowerID(), Equals, suite.healingStaff.ID())
	checker.Assert(suite.resultHealingStaffOnTerosAndVale.ResultPerTarget()[0].Healing().HitPointsRestored(), Equals, 4)
	checker.Assert(suite.resultHealingStaffOnTerosAndVale.ResultPerTarget()[1].PowerID(), Equals, suite.healingStaff.ID())
	checker.Assert(suite.resultHealingStaffOnTerosAndVale.ResultPerTarget()[1].Healing().HitPointsRestored(), Equals, 3)

	checker.Assert(
		suite.teros.CurrentHitPoints(),
		Equals,
		1+suite.resultHealingStaffOnTerosAndVale.ResultPerTarget()[0].Healing().HitPointsRestored(),
	)

	checker.Assert(
		suite.vale.CurrentHitPoints(),
		Equals,
		2+suite.resultHealingStaffOnTerosAndVale.ResultPerTarget()[1].Healing().HitPointsRestored(),
	)
}
