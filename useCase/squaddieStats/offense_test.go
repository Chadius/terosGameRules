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
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type squaddieOffense struct {
	teros *squaddie.Squaddie

	spear *power.Power
	blot  *power.Power

	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository

	repos           *repositories.RepositoryCollection
	offenseStrategy squaddiestats.CalculateSquaddieOffenseStatsStrategy
}

var _ = Suite(&squaddieOffense{})

func (suite *squaddieOffense) SetUpTest(checker *C) {
	suite.teros = squaddieBuilder.Builder().Teros().Build()

	suite.spear = powerBuilder.Builder().Spear().Build()
	suite.blot = powerBuilder.Builder().Blot().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros})

	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.blot})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}

	checkEquip := powerequip.CheckRepositories{}
	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.spear.GetReference(),
			suite.blot.GetReference(),
		},
		suite.repos,
	)

	suite.offenseStrategy = &squaddiestats.CalculateSquaddieOffenseStats{}
}

func (suite *squaddieOffense) TestSquaddieMeasuresAim(checker *C) {
	terosWithAim := squaddieBuilder.Builder().Teros().Aim(1).Build()
	suite.squaddieRepo.AddSquaddie(terosWithAim)

	spearAim, spearErr := suite.offenseStrategy.GetSquaddieAimWithPower(terosWithAim.ID(), suite.spear.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearAim, Equals, 2)

	accurateBlot := powerBuilder.Builder().Blot().ToHitBonus(2).Build()
	suite.repos.PowerRepo.AddPower(accurateBlot)

	blotAim, blotErr := suite.offenseStrategy.GetSquaddieAimWithPower(terosWithAim.ID(), accurateBlot.ID(), suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotAim, Equals, 3)
}

func (suite *squaddieOffense) TestReturnsAnErrorIfSquaddieDoesNotExist(checker *C) {
	_, err := suite.offenseStrategy.GetSquaddieAimWithPower("does not exist", suite.spear.ID(), suite.repos)
	checker.Assert(err, ErrorMatches, "squaddie could not be found, SquaddieID: does not exist")
}

func (suite *squaddieOffense) TestReturnsAnErrorIfPowerDoesNotExist(checker *C) {
	_, err := suite.offenseStrategy.GetSquaddieAimWithPower(suite.teros.ID(), "does not exist", suite.repos)
	checker.Assert(err, ErrorMatches, "power could not be found, SquaddieID: does not exist")
}

func (suite *squaddieOffense) TestReturnsAnErrorIfPowerHasNoAttackEffect(checker *C) {
	wait := powerBuilder.Builder().WithID("powerWait").Build()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{wait})

	checkEquip := powerequip.CheckRepositories{}
	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.spear.GetReference(),
			suite.blot.GetReference(),
			wait.GetReference(),
		},
		suite.repos,
	)

	_, err := suite.offenseStrategy.GetSquaddieAimWithPower(suite.teros.ID(), wait.ID(), suite.repos)
	checker.Assert(err, ErrorMatches, "cannot attack with power, SquaddieID: powerWait")
}

func (suite *squaddieOffense) TestGetRawDamageOfPhysicalPower(checker *C) {
	strongTeros := squaddieBuilder.Builder().Teros().Strength(1).Build()
	suite.squaddieRepo.AddSquaddie(strongTeros)

	weakerSpear := powerBuilder.Builder().Spear().DealsDamage(1).Build()
	suite.repos.PowerRepo.AddPower(weakerSpear)

	spearDamage, spearErr := suite.offenseStrategy.GetSquaddieRawDamageWithPower(strongTeros.ID(), weakerSpear.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearDamage, Equals, 2)
}

func (suite *squaddieOffense) TestGetRawDamageOfSpellPower(checker *C) {
	smartTeros := squaddieBuilder.Builder().Teros().Mind(3).Build()
	suite.squaddieRepo.AddSquaddie(smartTeros)

	weakerBlot := powerBuilder.Builder().Blot().DealsDamage(0).Build()
	suite.repos.PowerRepo.AddPower(weakerBlot)

	blotDamage, blotErr := suite.offenseStrategy.GetSquaddieRawDamageWithPower(smartTeros.ID(), weakerBlot.ID(), suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotDamage, Equals, 3)
}

func (suite *squaddieOffense) TestGetCriticalThresholdOfPower(checker *C) {
	criticalSpear := powerBuilder.Builder().Spear().CriticalHitThresholdBonus(2).CriticalDealsDamage(5).Build()
	suite.repos.PowerRepo.AddPower(criticalSpear)

	spearCritThreat, critErr := suite.offenseStrategy.GetSquaddieCriticalThresholdWithPower(suite.teros.ID(), criticalSpear.ID(), suite.repos)
	checker.Assert(critErr, IsNil)
	checker.Assert(spearCritThreat, Equals, 4)
}

func (suite *squaddieOffense) TestReturnsAnErrorIfPowerDoesNotCrit(checker *C) {
	_, critErr := suite.offenseStrategy.GetSquaddieCriticalThresholdWithPower(suite.teros.ID(), suite.blot.ID(), suite.repos)
	checker.Assert(critErr, ErrorMatches, "cannot critical hit with power, SquaddieID: powerBlot")
}

func (suite *squaddieOffense) TestGetCriticalDamageOfPower(checker *C) {
	strongTeros := squaddieBuilder.Builder().Teros().Strength(1).Build()
	suite.squaddieRepo.AddSquaddie(strongTeros)

	criticalSpear := powerBuilder.Builder().Spear().CriticalDealsDamage(5).Build()
	suite.repos.PowerRepo.AddPower(criticalSpear)

	spearDamage, damageErr := suite.offenseStrategy.GetSquaddieCriticalRawDamageWithPower(strongTeros.ID(), criticalSpear.ID(), suite.repos)
	checker.Assert(damageErr, IsNil)
	checker.Assert(spearDamage, Equals, 7)
}

func (suite *squaddieOffense) TestSquaddieCanCounterAttackWithPower(checker *C) {
	spearCanCounter, spearErr := suite.offenseStrategy.GetSquaddieCanCounterAttackWithPower(suite.teros.ID(), suite.spear.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearCanCounter, Equals, true)

	blotCanCounter, blotErr := suite.offenseStrategy.GetSquaddieCanCounterAttackWithPower(suite.teros.ID(), suite.blot.ID(), suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotCanCounter, Equals, false)
}

func (suite *squaddieOffense) TestSquaddieShowsCounterAttackToHit(checker *C) {
	accurateTeros := squaddieBuilder.Builder().Teros().Aim(2).Build()
	suite.squaddieRepo.AddSquaddie(accurateTeros)

	spearThatIsEasyToCounter := powerBuilder.Builder().Spear().CounterAttackPenaltyReduction(1).Build()
	suite.repos.PowerRepo.AddPower(spearThatIsEasyToCounter)

	spearAim, spearErr := suite.offenseStrategy.GetSquaddieCounterAttackAimWithPower(accurateTeros.ID(), spearThatIsEasyToCounter.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearAim, Equals, 2)
}

func (suite *squaddieOffense) TestGetTotalBarrierBurnOfAttacks(checker *C) {
	smartTeros := squaddieBuilder.Builder().Teros().Mind(3).Build()
	suite.squaddieRepo.AddSquaddie(smartTeros)

	blotWithBarrierBurn := powerBuilder.Builder().Blot().ExtraBarrierBurn(2).Build()
	suite.repos.PowerRepo.AddPower(blotWithBarrierBurn)

	blotDamage, blotErr := suite.offenseStrategy.GetSquaddieExtraBarrierBurnWithPower(smartTeros.ID(), blotWithBarrierBurn.ID(), suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotDamage, Equals, 2)
}

func (suite *squaddieOffense) TestCanCriticallyHitWithPower(checker *C) {
	criticalSpear := powerBuilder.Builder().Spear().CriticalDealsDamage(1).Build()
	suite.repos.PowerRepo.AddPower(criticalSpear)

	spearCanCrit, spearCanCritErr := suite.offenseStrategy.GetSquaddieCanCriticallyHitWithPower(suite.teros.ID(), criticalSpear.ID(), suite.repos)
	checker.Assert(spearCanCritErr, IsNil)
	checker.Assert(spearCanCrit, Equals, true)

	weakerBlot := powerBuilder.Builder().Blot().ToHitBonus(0).DealsDamage(0).Build()
	suite.repos.PowerRepo.AddPower(weakerBlot)

	blotCanCrit, blotCanCritErr := suite.offenseStrategy.GetSquaddieCanCriticallyHitWithPower(suite.teros.ID(), weakerBlot.ID(), suite.repos)
	checker.Assert(blotCanCritErr, IsNil)
	checker.Assert(blotCanCrit, Equals, false)
}

type healingPower struct {
	lini         *squaddie.Squaddie
	healingStaff *power.Power

	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository

	repos           *repositories.RepositoryCollection
	offenseStrategy squaddiestats.CalculateSquaddieOffenseStatsStrategy
}

var _ = Suite(&healingPower{})

func (suite *healingPower) SetUpTest(checker *C) {
	suite.lini = squaddieBuilder.Builder().Lini().Mind(1).Build()

	suite.healingStaff = powerBuilder.Builder().HealingStaff().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.lini})

	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.healingStaff})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}

	checkEquip := powerequip.CheckRepositories{}
	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.lini,
		[]*power.Reference{
			suite.healingStaff.GetReference(),
		},
		suite.repos,
	)

	suite.offenseStrategy = &squaddiestats.CalculateSquaddieOffenseStats{}
}

func (suite *healingPower) TestSquaddieKnowsHealingPotential(checker *C) {
	suite.lini.Defense.ReduceHitPoints(suite.lini.MaxHitPoints() - 1)

	staffHeal, staffErr := suite.offenseStrategy.GetHitPointsHealedWithPower(suite.lini.ID(), suite.healingStaff.ID(), suite.lini.ID(), suite.repos)
	checker.Assert(staffErr, IsNil)
	checker.Assert(staffHeal, Equals, 4)
}

type improveOffense struct {
	initialOffense *squaddie.Offense
}

var _ = Suite(&improveOffense{})

func (suite *improveOffense) SetUpTest(checker *C) {
	suite.initialOffense = squaddie.NewOffense(2, 3, 5)
}

func (suite *improveOffense) TestWhenImproveIsCalled_ThenAimStrengthMindIncrease(checker *C) {
	suite.initialOffense.Improve(7, 11, 13)

	checker.Assert(suite.initialOffense.Aim(), Equals, 9)
	checker.Assert(suite.initialOffense.Strength(), Equals, 14)
	checker.Assert(suite.initialOffense.Mind(), Equals, 18)
}
