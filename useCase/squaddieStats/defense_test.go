package squaddiestats_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
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

	weakerSpear *power.Power
	weakerBlot  *power.Power

	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository

	repos *repositories.RepositoryCollection

	defenseStrategy squaddiestats.CalculateSquaddieDefenseStatsStrategy
}

var _ = Suite(&squaddieDefense{})

func (suite *squaddieDefense) SetUpTest(checker *C) {
	suite.teros = squaddieBuilder.Builder().Teros().Build()

	suite.weakerSpear = powerBuilder.Builder().Spear().DealsDamage(0).ToHitBonus(0).Build()
	suite.weakerBlot = powerBuilder.Builder().Blot().DealsDamage(0).ToHitBonus(0).Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros})

	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.weakerSpear, suite.weakerBlot})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}

	checkEquip := powerequip.CheckRepositories{}
	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.weakerSpear.GetReference(),
			suite.weakerBlot.GetReference(),
		},
		suite.repos,
	)

	suite.defenseStrategy = &squaddiestats.CalculateSquaddieDefenseStats{}
}

func (suite *squaddieDefense) TestToHitPenaltyAgainstPhysicalAttacks(checker *C) {
	dodgyTeros := squaddieBuilder.Builder().Teros().Dodge(1).Build()
	suite.squaddieRepo.AddSquaddie(dodgyTeros)

	spearDodge, spearErr := suite.defenseStrategy.GetSquaddieToHitPenaltyAgainstPower(dodgyTeros.ID(), suite.weakerSpear.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearDodge, Equals, 1)
}

func (suite *squaddieDefense) TestToHitPenaltyAgainstSpellAttacks(checker *C) {
	deflectingTeros := squaddieBuilder.Builder().Teros().Deflect(2).Build()
	suite.squaddieRepo.AddSquaddie(deflectingTeros)

	blotDodge, blotErr := suite.defenseStrategy.GetSquaddieToHitPenaltyAgainstPower(deflectingTeros.ID(), suite.weakerBlot.ID(), suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotDodge, Equals, 2)
}

func (suite *squaddieDefense) TestGetDefenderArmorResistance(checker *C) {
	armoredTeros := squaddieBuilder.Builder().Teros().Armor(3).Build()
	suite.squaddieRepo.AddSquaddie(armoredTeros)

	spearArmor, spearErr := suite.defenseStrategy.GetSquaddieArmorAgainstPower(armoredTeros.ID(), suite.weakerSpear.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearArmor, Equals, 3)

	blotArmor, blotErr := suite.defenseStrategy.GetSquaddieArmorAgainstPower(armoredTeros.ID(), suite.weakerBlot.ID(), suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotArmor, Equals, 0)
}

func (suite *squaddieDefense) TestGetDefenderBarrierResistance(checker *C) {
	armoredTeros := squaddieBuilder.Builder().Teros().Armor(3).Barrier(1).Build()
	armoredTeros.Defense.SetBarrierToMax()
	suite.squaddieRepo.AddSquaddie(armoredTeros)

	spearBarrier, spearErr := suite.defenseStrategy.GetSquaddieBarrierAgainstPower(armoredTeros.ID(), suite.weakerSpear.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearBarrier, Equals, 1)
}

func (suite *squaddieDefense) TestGetDefenderCurrentHitPoints(checker *C) {
	injuredTeros := squaddieBuilder.Builder().Teros().Armor(3).Barrier(1).Build()
	injuredTeros.Defense.ReduceHitPoints(injuredTeros.MaxHitPoints() - 2)
	suite.squaddieRepo.AddSquaddie(injuredTeros)

	spearBarrier, spearErr := suite.defenseStrategy.GetSquaddieCurrentHitPoints(injuredTeros.ID(), suite.weakerSpear.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearBarrier, Equals, 2)
}
