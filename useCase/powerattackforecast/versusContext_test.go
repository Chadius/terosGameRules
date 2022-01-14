package powerattackforecast_test

import (
	"github.com/chadius/terosgamerules/entity/power"
	"github.com/chadius/terosgamerules/entity/powerinterface"
	"github.com/chadius/terosgamerules/entity/powerrepository"
	"github.com/chadius/terosgamerules/entity/powerusagescenario"
	"github.com/chadius/terosgamerules/entity/squaddie"
	"github.com/chadius/terosgamerules/entity/squaddieinterface"
	"github.com/chadius/terosgamerules/usecase/powerattackforecast"
	"github.com/chadius/terosgamerules/usecase/repositories"
	"github.com/chadius/terosgamerules/usecase/squaddiestats"
	. "gopkg.in/check.v1"
)

type VersusContextTestSuite struct {
	teros  squaddieinterface.Interface
	bandit squaddieinterface.Interface
	spear  powerinterface.Interface
	blot   powerinterface.Interface
	axe    powerinterface.Interface

	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository

	versusContextSpearOnBandit powerattackforecast.VersusContextStrategy
	versusContextBlotOnBandit  powerattackforecast.VersusContextStrategy
}

var _ = Suite(&VersusContextTestSuite{})

func (suite *VersusContextTestSuite) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Aim(2).Strength(2).Mind(2).Build()

	suite.spear = power.NewPowerBuilder().Spear().Build()
	suite.blot = power.NewPowerBuilder().Blot().Build()

	suite.bandit = squaddie.NewSquaddieBuilder().Bandit().Barrier(3).Armor(1).Deflect(2).Dodge(1).Build()
	suite.bandit.SetBarrierToMax()

	suite.axe = power.NewPowerBuilder().Axe().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]squaddieinterface.Interface{suite.teros, suite.bandit})

	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]powerinterface.Interface{suite.spear, suite.blot, suite.axe})

	suite.versusContextSpearOnBandit = &powerattackforecast.VersusContext{}
	suite.CalculateSpearOnBandit(nil)

	suite.versusContextBlotOnBandit = &powerattackforecast.VersusContext{}
	suite.CalculateBlotOnBandit(nil)
}

func (suite *VersusContextTestSuite) CalculateSpearOnBandit(setup *powerusagescenario.Setup) {
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

	attackerContext := powerattackforecast.NewAttackerContext(&squaddiestats.CalculateSquaddieOffenseStats{})
	attackerContext.Calculate(
		setupToUse,
		repos,
	)

	defenderContext := powerattackforecast.NewDefenderContext(suite.bandit.ID(), &squaddiestats.CalculateSquaddieDefenseStats{})
	defenderContext.Calculate(
		&setupToUse,
		&repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	)

	suite.versusContextSpearOnBandit.Calculate(*attackerContext, *defenderContext)
}

func (suite *VersusContextTestSuite) CalculateBlotOnBandit(setup *powerusagescenario.Setup) {
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

	attackerContext := powerattackforecast.NewAttackerContext(&squaddiestats.CalculateSquaddieOffenseStats{})
	attackerContext.Calculate(
		setupToUse,
		repos,
	)

	defenderContext := powerattackforecast.NewDefenderContext(suite.bandit.ID(), &squaddiestats.CalculateSquaddieDefenseStats{})
	defenderContext.Calculate(
		&setupToUse,
		&repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	)

	suite.versusContextBlotOnBandit.Calculate(*attackerContext, *defenderContext)
}

func (suite *VersusContextTestSuite) TestNetToHitReliesOnToHitMinusDodgeOrDeflect(checker *C) {
	checker.Assert(suite.versusContextSpearOnBandit.ToHit().ToHitBonus, Equals, 2)
	checker.Assert(suite.versusContextBlotOnBandit.ToHit().ToHitBonus, Equals, 0)
}

func (suite *VersusContextTestSuite) TestTargetTakesFullDamageAgainstPhysicalWhenNoArmor(checker *C) {
	suite.bandit = squaddie.NewSquaddieBuilder().Bandit().Armor(0).Barrier(0).Build()
	suite.squaddieRepo.AddSquaddie(suite.bandit)

	suite.CalculateSpearOnBandit(
		&powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.spear.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
	)

	suite.CalculateBlotOnBandit(
		&powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.blot.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
	)
	checker.Assert(suite.versusContextSpearOnBandit.NormalDamage().RawDamageDealt, Equals, 3)
	checker.Assert(suite.versusContextBlotOnBandit.NormalDamage().RawDamageDealt, Equals, 5)
}

func (suite *VersusContextTestSuite) TestTargetUsesArmorResistAgainstPhysicalOnly(checker *C) {
	suite.bandit = squaddie.NewSquaddieBuilder().Armor(1).Barrier(0).Build()
	suite.squaddieRepo.AddSquaddie(suite.bandit)

	suite.CalculateSpearOnBandit(
		&powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.spear.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
	)

	suite.CalculateBlotOnBandit(
		&powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.blot.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
	)

	checker.Assert(suite.versusContextSpearOnBandit.NormalDamage().DamageAbsorbedByArmor, Equals, 1)
	checker.Assert(suite.versusContextSpearOnBandit.NormalDamage().RawDamageDealt, Equals, 2)

	checker.Assert(suite.versusContextBlotOnBandit.NormalDamage().DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(suite.versusContextBlotOnBandit.NormalDamage().RawDamageDealt, Equals, 5)
}

func (suite *VersusContextTestSuite) TestTargetUsesBarrierToResistDamageFromAllAttacks(checker *C) {
	suite.bandit = squaddie.NewSquaddieBuilder().Armor(1).Barrier(3).Build()
	suite.bandit.SetBarrierToMax()
	suite.squaddieRepo.AddSquaddie(suite.bandit)

	suite.CalculateSpearOnBandit(
		&powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.spear.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
	)

	suite.CalculateBlotOnBandit(
		&powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.blot.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
	)

	checker.Assert(suite.versusContextSpearOnBandit.NormalDamage().DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(suite.versusContextSpearOnBandit.NormalDamage().RawDamageDealt, Equals, 0)

	checker.Assert(suite.versusContextBlotOnBandit.NormalDamage().DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(suite.versusContextBlotOnBandit.NormalDamage().RawDamageDealt, Equals, 2)
}

func (suite *VersusContextTestSuite) TestBarrierBurnCanSpillOverDamage(checker *C) {
	suite.blot = power.NewPowerBuilder().CloneOf(suite.blot).WithID(suite.blot.ID()).DealsDamage(1).ExtraBarrierBurn(2).Build()
	suite.powerRepo.AddPower(suite.blot)

	suite.CalculateBlotOnBandit(
		&powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.blot.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
	)

	checker.Assert(suite.versusContextBlotOnBandit.NormalDamage().DamageAbsorbedByBarrier, Equals, 1)
	checker.Assert(suite.versusContextBlotOnBandit.NormalDamage().ExtraBarrierBurnt, Equals, 2)
	checker.Assert(suite.versusContextBlotOnBandit.NormalDamage().TotalRawBarrierBurnt, Equals, 3)
	checker.Assert(suite.versusContextBlotOnBandit.NormalDamage().RawDamageDealt, Equals, 2)
}

func (suite *VersusContextTestSuite) TestBarrierBurnCanBeTolerated(checker *C) {
	suite.blot = power.NewPowerBuilder().CloneOf(suite.blot).WithID(suite.blot.ID()).DealsDamage(0).ExtraBarrierBurn(1).Build()
	suite.powerRepo.AddPower(suite.blot)

	suite.CalculateBlotOnBandit(
		&powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.blot.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
	)

	checker.Assert(suite.versusContextBlotOnBandit.NormalDamage().DamageAbsorbedByBarrier, Equals, 2)
	checker.Assert(suite.versusContextBlotOnBandit.NormalDamage().ExtraBarrierBurnt, Equals, 1)
	checker.Assert(suite.versusContextBlotOnBandit.NormalDamage().TotalRawBarrierBurnt, Equals, 3)
	checker.Assert(suite.versusContextBlotOnBandit.NormalDamage().RawDamageDealt, Equals, 0)
}

func (suite *VersusContextTestSuite) TestCriticalHitChanceIsShown(checker *C) {
	suite.spear = power.NewPowerBuilder().CloneOf(suite.spear).WithID(suite.spear.ID()).CriticalDealsDamage(3).CriticalHitThresholdBonus(0).Build()
	suite.powerRepo.AddPower(suite.spear)

	suite.CalculateSpearOnBandit(
		&powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.spear.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
	)

	checker.Assert(suite.versusContextSpearOnBandit.CanCritical(), Equals, true)
	checker.Assert(suite.versusContextSpearOnBandit.CriticalHitThreshold(), Equals, 6)
}

func (suite *VersusContextTestSuite) TestCriticalDamageDistributes(checker *C) {
	suite.spear = power.NewPowerBuilder().CloneOf(suite.spear).WithID(suite.spear.ID()).CriticalDealsDamage(3).CriticalHitThresholdBonus(0).Build()
	suite.powerRepo.AddPower(suite.spear)

	suite.bandit = squaddie.NewSquaddieBuilder().Armor(1).Barrier(3).Build()
	suite.squaddieRepo.AddSquaddie(suite.bandit)
	suite.bandit.SetBarrierToMax()

	suite.CalculateSpearOnBandit(
		&powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.spear.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
	)

	checker.Assert(suite.versusContextSpearOnBandit.CriticalHitDamage().DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(suite.versusContextSpearOnBandit.CriticalHitDamage().DamageAbsorbedByArmor, Equals, 1)
	checker.Assert(suite.versusContextSpearOnBandit.CriticalHitDamage().RawDamageDealt, Equals, 2)
	checker.Assert(suite.versusContextSpearOnBandit.CriticalHitDamage().ExtraBarrierBurnt, Equals, 0)
	checker.Assert(suite.versusContextSpearOnBandit.CriticalHitDamage().TotalRawBarrierBurnt, Equals, 3)
}

func (suite *VersusContextTestSuite) TestNoCriticalDamageDistributionIfCannotCritical(checker *C) {
	checker.Assert(suite.versusContextSpearOnBandit.CriticalHitDamage(), IsNil)
}

func (suite *VersusContextTestSuite) TestKnowsIfAttackIsNotFatalToTarget(checker *C) {
	suite.bandit = squaddie.NewSquaddieBuilder().Armor(0).Barrier(0).Build()
	suite.squaddieRepo.AddSquaddie(suite.bandit)

	suite.teros = squaddie.NewSquaddieBuilder().Teros().Mind(0).Build()
	suite.squaddieRepo.AddSquaddie(suite.teros)

	suite.blot = power.NewPowerBuilder().CloneOf(suite.blot).WithID(suite.blot.ID()).DealsDamage(0).Build()
	suite.powerRepo.AddPower(suite.blot)

	suite.CalculateSpearOnBandit(
		&powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.spear.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
	)

	checker.Assert(suite.versusContextSpearOnBandit.NormalDamage().IsFatalToTarget, Equals, false)
}

func (suite *VersusContextTestSuite) TestKnowsIfAttackIsFatalToTarget(checker *C) {
	suite.spear = power.NewPowerBuilder().CloneOf(suite.spear).WithID(suite.spear.ID()).DealsDamage(
		suite.bandit.MaxHitPoints() + suite.bandit.Armor() + suite.bandit.MaxBarrier(),
	).Build()
	suite.powerRepo.AddPower(suite.spear)

	suite.CalculateSpearOnBandit(
		&powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.spear.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
	)
	checker.Assert(suite.versusContextSpearOnBandit.NormalDamage().IsFatalToTarget, Equals, true)
}
