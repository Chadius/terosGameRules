package squaddiestats_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/usecase/squaddiestats"
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
	suite.teros = squaddie.NewSquaddie("teros")
	suite.teros.Identification.Name = "teros"

	suite.spear = power.NewPower("spear")
	suite.spear.PowerType = power.Physical
	suite.spear.ID = "powerSpear"

	suite.blot = power.NewPower("blot")
	suite.blot.PowerType = power.Spell
	suite.blot.ID = "powerBlot"

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
	suite.teros.Defense.Dodge = 1

	suite.spear.AttackEffect = &power.AttackingEffect{}

	spearDodge, spearErr := squaddiestats.GetSquaddieToHitPenaltyAgainstPower(suite.teros.Identification.ID, suite.spear.ID, suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearDodge, Equals, 1)
}

func (suite *squaddieDefense) TestToHitPenaltyAgainstSpellAttacks(checker *C) {
	suite.teros.Defense.Deflect = 2

	suite.blot.AttackEffect = &power.AttackingEffect{}

	blotDodge, blotErr := squaddiestats.GetSquaddieToHitPenaltyAgainstPower(suite.teros.Identification.ID, suite.blot.ID, suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotDodge, Equals, 2)
}

func (suite *squaddieDefense) TestGetDefenderArmorResistance(checker *C) {
	suite.teros.Defense.Armor = 3
	suite.teros.Defense.CurrentBarrier = 0

	suite.spear.AttackEffect = &power.AttackingEffect{}
	spearArmor, spearErr := squaddiestats.GetSquaddieArmorAgainstPower(suite.teros.Identification.ID, suite.spear.ID, suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearArmor, Equals, 3)

	suite.blot.AttackEffect = &power.AttackingEffect{}
	blotArmor, blotErr := squaddiestats.GetSquaddieArmorAgainstPower(suite.teros.Identification.ID, suite.blot.ID, suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotArmor, Equals, 0)
}

func (suite *squaddieDefense) TestGetDefenderBarrierResistance(checker *C) {
	suite.teros.Defense.Armor = 3
	suite.teros.Defense.CurrentBarrier = 1
	suite.spear.AttackEffect = &power.AttackingEffect{}

	spearBarrier, spearErr := squaddiestats.GetSquaddieBarrierAgainstPower(suite.teros.Identification.ID, suite.spear.ID, suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearBarrier, Equals, 1)
}

func (suite *squaddieDefense) TestGetDefenderCurrentHitPoints(checker *C) {
	suite.teros.Defense.MaxHitPoints = 5
	suite.teros.Defense.CurrentHitPoints = 2

	suite.spear.AttackEffect = &power.AttackingEffect{}

	spearBarrier, spearErr := squaddiestats.GetSquaddieCurrentHitPoints(suite.teros.Identification.ID, suite.spear.ID, suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearBarrier, Equals, 2)
}
