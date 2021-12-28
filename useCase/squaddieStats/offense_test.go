package squaddiestats_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerreference"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/usecase/squaddiestats"
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
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Build()

	suite.spear = power.NewPowerBuilder().Spear().Build()
	suite.blot = power.NewPowerBuilder().Blot().Build()

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
		[]*powerreference.Reference{
			suite.spear.GetReference(),
			suite.blot.GetReference(),
		},
		suite.repos,
	)

	suite.offenseStrategy = &squaddiestats.CalculateSquaddieOffenseStats{}
}

func (suite *squaddieOffense) TestSquaddieMeasuresAim(checker *C) {
	terosWithAim := squaddie.NewSquaddieBuilder().Teros().Aim(1).Build()
	suite.squaddieRepo.AddSquaddie(terosWithAim)

	spearAim, spearErr := suite.offenseStrategy.GetSquaddieAimWithPower(terosWithAim.ID(), suite.spear.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearAim, Equals, 2)

	accurateBlot := power.NewPowerBuilder().Blot().ToHitBonus(2).Build()
	suite.repos.PowerRepo.AddPower(accurateBlot)

	blotAim, blotErr := suite.offenseStrategy.GetSquaddieAimWithPower(terosWithAim.ID(), accurateBlot.ID(), suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotAim, Equals, 3)
}

func (suite *squaddieOffense) TestReturnsAnErrorIfSquaddieDoesNotExist(checker *C) {
	_, err := suite.offenseStrategy.GetSquaddieAimWithPower("does not exist", suite.spear.ID(), suite.repos)
	checker.Assert(err, ErrorMatches, "squaddie could not be found, id: does not exist")
}

func (suite *squaddieOffense) TestReturnsAnErrorIfPowerDoesNotExist(checker *C) {
	_, err := suite.offenseStrategy.GetSquaddieAimWithPower(suite.teros.ID(), "does not exist", suite.repos)
	checker.Assert(err, ErrorMatches, "power could not be found, id: does not exist")
}

func (suite *squaddieOffense) TestReturnsAnErrorIfPowerHasNoAttackEffect(checker *C) {
	wait := power.NewPowerBuilder().WithID("powerWait").Build()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{wait})

	checkEquip := powerequip.CheckRepositories{}
	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*powerreference.Reference{
			suite.spear.GetReference(),
			suite.blot.GetReference(),
			wait.GetReference(),
		},
		suite.repos,
	)

	_, err := suite.offenseStrategy.GetSquaddieAimWithPower(suite.teros.ID(), wait.ID(), suite.repos)
	checker.Assert(err, ErrorMatches, "cannot attack with power, id: powerWait")
}

func (suite *squaddieOffense) TestGetRawDamageOfPhysicalPower(checker *C) {
	strongTeros := squaddie.NewSquaddieBuilder().Teros().Strength(1).Build()
	suite.squaddieRepo.AddSquaddie(strongTeros)

	weakerSpear := power.NewPowerBuilder().Spear().DealsDamage(1).Build()
	suite.repos.PowerRepo.AddPower(weakerSpear)

	spearDamage, spearErr := suite.offenseStrategy.GetSquaddieRawDamageWithPower(strongTeros.ID(), weakerSpear.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearDamage, Equals, 2)
}

func (suite *squaddieOffense) TestGetRawDamageOfSpellPower(checker *C) {
	smartTeros := squaddie.NewSquaddieBuilder().Teros().Mind(3).Build()
	suite.squaddieRepo.AddSquaddie(smartTeros)

	weakerBlot := power.NewPowerBuilder().Blot().DealsDamage(0).Build()
	suite.repos.PowerRepo.AddPower(weakerBlot)

	blotDamage, blotErr := suite.offenseStrategy.GetSquaddieRawDamageWithPower(smartTeros.ID(), weakerBlot.ID(), suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotDamage, Equals, 3)
}

func (suite *squaddieOffense) TestGetCriticalThresholdOfPower(checker *C) {
	criticalSpear := power.NewPowerBuilder().Spear().CriticalHitThresholdBonus(2).CriticalDealsDamage(5).Build()
	suite.repos.PowerRepo.AddPower(criticalSpear)

	spearCritThreat, critErr := suite.offenseStrategy.GetSquaddieCriticalThresholdWithPower(suite.teros.ID(), criticalSpear.ID(), suite.repos)
	checker.Assert(critErr, IsNil)
	checker.Assert(spearCritThreat, Equals, 4)
}

func (suite *squaddieOffense) TestReturnsAnErrorIfPowerDoesNotCrit(checker *C) {
	_, critErr := suite.offenseStrategy.GetSquaddieCriticalThresholdWithPower(suite.teros.ID(), suite.blot.ID(), suite.repos)
	checker.Assert(critErr, ErrorMatches, "cannot critical hit with power, id: powerBlot")
}

func (suite *squaddieOffense) TestGetCriticalDamageOfPower(checker *C) {
	strongTeros := squaddie.NewSquaddieBuilder().Teros().Strength(1).Build()
	suite.squaddieRepo.AddSquaddie(strongTeros)

	criticalSpear := power.NewPowerBuilder().Spear().CriticalDealsDamage(5).Build()
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
	accurateTeros := squaddie.NewSquaddieBuilder().Teros().Aim(2).Build()
	suite.squaddieRepo.AddSquaddie(accurateTeros)

	spearThatIsEasyToCounter := power.NewPowerBuilder().Spear().CounterAttackPenaltyReduction(1).Build()
	suite.repos.PowerRepo.AddPower(spearThatIsEasyToCounter)

	spearAim, spearErr := suite.offenseStrategy.GetSquaddieCounterAttackAimWithPower(accurateTeros.ID(), spearThatIsEasyToCounter.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearAim, Equals, 2)
}

func (suite *squaddieOffense) TestGetTotalBarrierBurnOfAttacks(checker *C) {
	smartTeros := squaddie.NewSquaddieBuilder().Teros().Mind(3).Build()
	suite.squaddieRepo.AddSquaddie(smartTeros)

	blotWithBarrierBurn := power.NewPowerBuilder().Blot().ExtraBarrierBurn(2).Build()
	suite.repos.PowerRepo.AddPower(blotWithBarrierBurn)

	blotDamage, blotErr := suite.offenseStrategy.GetSquaddieExtraBarrierBurnWithPower(smartTeros.ID(), blotWithBarrierBurn.ID(), suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotDamage, Equals, 2)
}

func (suite *squaddieOffense) TestCanCriticallyHitWithPower(checker *C) {
	criticalSpear := power.NewPowerBuilder().Spear().CriticalDealsDamage(1).Build()
	suite.repos.PowerRepo.AddPower(criticalSpear)

	spearCanCrit, spearCanCritErr := suite.offenseStrategy.GetSquaddieCanCriticallyHitWithPower(suite.teros.ID(), criticalSpear.ID(), suite.repos)
	checker.Assert(spearCanCritErr, IsNil)
	checker.Assert(spearCanCrit, Equals, true)

	weakerBlot := power.NewPowerBuilder().Blot().ToHitBonus(0).DealsDamage(0).Build()
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
	suite.lini = squaddie.NewSquaddieBuilder().Lini().Mind(1).Build()

	suite.healingStaff = power.NewPowerBuilder().HealingStaff().Build()

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
		[]*powerreference.Reference{
			suite.healingStaff.GetReference(),
		},
		suite.repos,
	)

	suite.offenseStrategy = &squaddiestats.CalculateSquaddieOffenseStats{}
}

func (suite *healingPower) TestSquaddieKnowsHealingPotential(checker *C) {
	suite.lini.ReduceHitPoints(suite.lini.MaxHitPoints() - 1)

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

type CounterAttackEquipmentCheck struct {
	teros           *squaddie.Squaddie
	spear           *power.Power
	scimitar        *power.Power
	blot            *power.Power
	powerRepo       *powerrepository.Repository
	squaddieRepo    *squaddie.Repository
	repos           *repositories.RepositoryCollection
	equipCheck      powerequip.Strategy
	offenseStrategy squaddiestats.CalculateSquaddieOffenseStatsStrategy
}

var _ = Suite(&CounterAttackEquipmentCheck{})

func (suite *CounterAttackEquipmentCheck) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Build()
	suite.spear = power.NewPowerBuilder().Spear().Build()
	suite.scimitar = power.NewPowerBuilder().WithName("scimitar the second").CanBeEquipped().Build()
	suite.blot = power.NewPowerBuilder().Blot().CannotBeEquipped().Build()

	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{
		suite.spear,
		suite.scimitar,
		suite.blot,
	})

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}

	suite.equipCheck = &powerequip.CheckRepositories{}

	suite.offenseStrategy = &squaddiestats.CalculateSquaddieOffenseStats{}
}

func (suite *CounterAttackEquipmentCheck) TestSquaddieCanCounter(checker *C) {
	terosPowerReferences := []*powerreference.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	suite.equipCheck.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	suite.equipCheck.EquipDefaultPower(suite.teros, suite.repos)
	canCounter, _ := suite.offenseStrategy.CanSquaddieCounterWithEquippedWeapon(suite.teros.ID(), suite.repos)
	checker.Assert(canCounter, Equals, true)
}

func (suite *CounterAttackEquipmentCheck) TestSquaddieCannotCounterWithUncounterablePower(checker *C) {
	terosPowerReferences := []*powerreference.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	suite.equipCheck.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	suite.equipCheck.SquaddieEquipPower(suite.teros, suite.scimitar.ID(), suite.repos)
	canCounter, _ := suite.offenseStrategy.CanSquaddieCounterWithEquippedWeapon(suite.teros.ID(), suite.repos)
	checker.Assert(canCounter, Equals, false)
}

func (suite *CounterAttackEquipmentCheck) TestSquaddieCannotCounterWithUnequippablePower(checker *C) {
	terosPowerReferences := []*powerreference.Reference{
		suite.blot.GetReference(),
	}
	suite.equipCheck.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	suite.equipCheck.EquipDefaultPower(suite.teros, suite.repos)
	canCounter, _ := suite.offenseStrategy.CanSquaddieCounterWithEquippedWeapon(suite.teros.ID(), suite.repos)
	checker.Assert(canCounter, Equals, false)
}
