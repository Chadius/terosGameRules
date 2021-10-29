package squaddiestats_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/usecase/squaddiestats"
	powerBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	. "gopkg.in/check.v1"
)

type squaddieDefense struct {
	teros *squaddie.Squaddie

	spear *power.Power
	blot  *power.Power

	powerRepo    *power.Repository
	squaddieRepo *squaddie.Repository

	repos *repositories.RepositoryCollection
}

var _ = Suite(&squaddieDefense{})

func (suite *squaddieDefense) SetUpTest(checker *C) {
	suite.teros = squaddieBuilder.Builder().Teros().Build()

	suite.spear = powerBuilder.Builder().Spear().Build()
	suite.blot = powerBuilder.Builder().Blot().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.blot})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.spear.GetReference(),
			suite.blot.GetReference(),
		},
		suite.repos,
	)
}

func (suite *squaddieDefense) TestToHitPenaltyAgainstPhysicalAttacks(checker *C) {
	suite.teros.Defense.SquaddieDodge = 1

	suite.spear.AttackEffect = &power.AttackingEffect{}

	spearDodge, spearErr := squaddiestats.GetSquaddieToHitPenaltyAgainstPower(suite.teros.ID(), suite.spear.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearDodge, Equals, 1)
}

func (suite *squaddieDefense) TestToHitPenaltyAgainstSpellAttacks(checker *C) {
	suite.teros.Defense.SquaddieDeflect = 2

	suite.blot.AttackEffect = &power.AttackingEffect{}

	blotDodge, blotErr := squaddiestats.GetSquaddieToHitPenaltyAgainstPower(suite.teros.ID(), suite.blot.ID(), suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotDodge, Equals, 2)
}

func (suite *squaddieDefense) TestGetDefenderArmorResistance(checker *C) {
	suite.teros.Defense.SquaddieArmor = 3
	suite.teros.Defense.SquaddieCurrentBarrier = 0

	suite.spear.AttackEffect = &power.AttackingEffect{}
	spearArmor, spearErr := squaddiestats.GetSquaddieArmorAgainstPower(suite.teros.ID(), suite.spear.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearArmor, Equals, 3)

	suite.blot.AttackEffect = &power.AttackingEffect{}
	blotArmor, blotErr := squaddiestats.GetSquaddieArmorAgainstPower(suite.teros.ID(), suite.blot.ID(), suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotArmor, Equals, 0)
}

func (suite *squaddieDefense) TestGetDefenderBarrierResistance(checker *C) {
	suite.teros.Defense.SquaddieArmor = 3
	suite.teros.Defense.SquaddieCurrentBarrier = 1
	suite.spear.AttackEffect = &power.AttackingEffect{}

	spearBarrier, spearErr := squaddiestats.GetSquaddieBarrierAgainstPower(suite.teros.ID(), suite.spear.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearBarrier, Equals, 1)
}

func (suite *squaddieDefense) TestGetDefenderCurrentHitPoints(checker *C) {
	suite.teros.Defense.SquaddieMaxHitPoints = 5
	suite.teros.Defense.SquaddieCurrentHitPoints = 2

	suite.spear.AttackEffect = &power.AttackingEffect{}

	spearBarrier, spearErr := squaddiestats.GetSquaddieCurrentHitPoints(suite.teros.ID(), suite.spear.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearBarrier, Equals, 2)
}
