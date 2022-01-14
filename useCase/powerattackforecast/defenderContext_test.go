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

type DefenderContextTestSuite struct {
	teros  squaddieinterface.Interface
	bandit squaddieinterface.Interface
	spear  powerinterface.Interface
	blot   powerinterface.Interface
	axe    powerinterface.Interface

	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository

	defenderContextSpearOnBandit powerattackforecast.DefenderContextStrategy
	defenderContextBlotOnBandit  powerattackforecast.DefenderContextStrategy
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
	suite.squaddieRepo.AddSquaddies([]squaddieinterface.Interface{suite.teros, suite.bandit})

	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]powerinterface.Interface{suite.spear, suite.blot, suite.axe})

	suite.defenderContextSpearOnBandit = powerattackforecast.NewDefenderContext(suite.bandit.ID(), &squaddiestats.CalculateSquaddieDefenseStats{})
	suite.CalculateSpearOnBandit(nil)

	suite.defenderContextBlotOnBandit = powerattackforecast.NewDefenderContext(suite.bandit.ID(), &squaddiestats.CalculateSquaddieDefenseStats{})
	suite.CalculateBlotOnBandit(nil)
}

func (suite *DefenderContextTestSuite) CalculateSpearOnBandit(setup *powerusagescenario.Setup) {
	setupToUse := powerusagescenario.Setup{
		UserID:          suite.teros.ID(),
		PowerID:         suite.spear.ID(),
		Targets:         []string{suite.bandit.ID()},
		IsCounterAttack: false,
	}
	if setup != nil {
		setupToUse = *setup
	}

	suite.defenderContextSpearOnBandit.Calculate(
		&setupToUse,
		&repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	)
}

func (suite *DefenderContextTestSuite) CalculateBlotOnBandit(setup *powerusagescenario.Setup) {
	setupToUse := powerusagescenario.Setup{
		UserID:          suite.teros.ID(),
		PowerID:         suite.blot.ID(),
		Targets:         []string{suite.bandit.ID()},
		IsCounterAttack: false,
	}
	if setup != nil {
		setupToUse = *setup
	}

	suite.defenderContextBlotOnBandit.Calculate(
		&setupToUse,
		&repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	)
}

func (suite *DefenderContextTestSuite) TestGetDefenderDodge(checker *C) {
	checker.Assert(suite.defenderContextSpearOnBandit.TotalToHitPenalty(), Equals, 1)
}

func (suite *DefenderContextTestSuite) TestGetDefenderArmorResistance(checker *C) {
	checker.Assert(suite.defenderContextSpearOnBandit.ArmorResistance(), Equals, 1)
	checker.Assert(suite.defenderContextBlotOnBandit.ArmorResistance(), Equals, 0)
}

func (suite *DefenderContextTestSuite) TestGetDefenderDeflect(checker *C) {
	checker.Assert(suite.defenderContextBlotOnBandit.TotalToHitPenalty(), Equals, 2)
}

func (suite *DefenderContextTestSuite) TestGetDefenderBarrierAbsorb(checker *C) {
	checker.Assert(suite.defenderContextBlotOnBandit.BarrierResistance(), Equals, 3)
	checker.Assert(suite.defenderContextSpearOnBandit.BarrierResistance(), Equals, 3)
}
