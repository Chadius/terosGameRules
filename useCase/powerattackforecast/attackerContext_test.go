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
	"reflect"
)

type AttackContextTestSuite struct {
	teros  squaddieinterface.Interface
	bandit squaddieinterface.Interface
	spear  powerinterface.Interface
	blot   powerinterface.Interface

	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository

	attackerContextSpearOnBandit powerattackforecast.AttackerContextStrategy
	attackerContextBlotOnBandit  powerattackforecast.AttackerContextStrategy
}

var _ = Suite(&AttackContextTestSuite{})

func (suite *AttackContextTestSuite) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Aim(2).Strength(2).Mind(2).Build()

	suite.spear = power.NewPowerBuilder().Spear().Build()
	suite.blot = power.NewPowerBuilder().Blot().Build()

	suite.bandit = squaddie.NewSquaddieBuilder().Bandit().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]squaddieinterface.Interface{suite.teros, suite.bandit})

	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]powerinterface.Interface{suite.spear, suite.blot})

	suite.attackerContextSpearOnBandit = powerattackforecast.NewAttackerContext(&squaddiestats.CalculateSquaddieOffenseStats{})
	suite.CalculateSpearOnBandit(nil)

	suite.attackerContextBlotOnBandit = powerattackforecast.NewAttackerContext(&squaddiestats.CalculateSquaddieOffenseStats{})
	suite.CalculateBlotOnBandit(nil)
}

func (suite *AttackContextTestSuite) CalculateSpearOnBandit(setup *powerusagescenario.Setup) {
	setupToUse := powerusagescenario.Setup{
		UserID:          suite.teros.ID(),
		PowerID:         suite.spear.ID(),
		Targets:         []string{suite.bandit.ID()},
		IsCounterAttack: false,
	}
	if setup != nil {
		setupToUse = *setup
	}

	suite.attackerContextSpearOnBandit.Calculate(
		setupToUse,
		&repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	)
}

func (suite *AttackContextTestSuite) CalculateBlotOnBandit(setup *powerusagescenario.Setup) {
	setupToUse := powerusagescenario.Setup{
		UserID:          suite.teros.ID(),
		PowerID:         suite.blot.ID(),
		Targets:         []string{suite.bandit.ID()},
		IsCounterAttack: false,
	}
	if setup != nil {
		setupToUse = *setup
	}

	suite.attackerContextBlotOnBandit.Calculate(
		setupToUse,
		&repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	)
}

func (suite *AttackContextTestSuite) TestGetAttackerHitBonus(checker *C) {
	suite.CalculateSpearOnBandit(nil)
	checker.Assert(suite.attackerContextSpearOnBandit.TotalToHitBonus(), Equals, 3)
}

func (suite *AttackContextTestSuite) TestGetAttackerHitBonusOnCounterAttacks(checker *C) {
	suite.CalculateSpearOnBandit(&powerusagescenario.Setup{
		UserID:          suite.teros.ID(),
		PowerID:         suite.spear.ID(),
		Targets:         []string{suite.bandit.ID()},
		IsCounterAttack: true,
	})
	checker.Assert(suite.attackerContextSpearOnBandit.TotalToHitBonus(), Equals, 1)
}

func (suite *AttackContextTestSuite) TestGetAttackerPhysicalRawDamage(checker *C) {
	checker.Assert(
		reflect.TypeOf(suite.attackerContextSpearOnBandit.PowerSourceLogic()).String(),
		Equals,
		"*powersource.Physical",
	)
	checker.Assert(suite.attackerContextSpearOnBandit.RawDamage(), Equals, 3)
}

func (suite *AttackContextTestSuite) TestGetAttackerSpellDamage(checker *C) {
	suite.CalculateBlotOnBandit(nil)
	checker.Assert(
		reflect.TypeOf(suite.attackerContextBlotOnBandit.PowerSourceLogic()).String(),
		Equals,
		"*powersource.Spell",
	)
	checker.Assert(suite.attackerContextBlotOnBandit.RawDamage(), Equals, 5)
}

func (suite *AttackContextTestSuite) TestCriticalHits(checker *C) {
	suite.spear = power.NewPowerBuilder().CloneOf(suite.spear).WithID(suite.spear.ID()).CriticalDealsDamage(3).CriticalHitThresholdBonus(0).Build()
	suite.powerRepo.AddPower(suite.spear)

	suite.CalculateSpearOnBandit(
		&powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.spear.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: true,
		},
	)

	checker.Assert(suite.attackerContextSpearOnBandit.CanCritical(), Equals, true)
	checker.Assert(suite.attackerContextSpearOnBandit.CriticalHitThreshold(), Equals, 6)
	checker.Assert(suite.attackerContextSpearOnBandit.CriticalHitDamage(), Equals, 6)
}
