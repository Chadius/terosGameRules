package powerequip_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerreference"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type SquaddieEquipPowersFromRepo struct {
	teros        *squaddie.Squaddie
	spear        *power.Power
	scimitar     *power.Power
	blot         *power.Power
	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection
	equipCheck   powerequip.Strategy
}

var _ = Suite(&SquaddieEquipPowersFromRepo{})

func (suite *SquaddieEquipPowersFromRepo) SetUpTest(checker *C) {
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
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieEquipsFirstPowerByDefault(checker *C) {
	terosPowerReferences := []*powerreference.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	suite.equipCheck.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	suite.equipCheck.EquipDefaultPower(suite.teros, suite.repos)
	checker.Assert(suite.teros.GetEquippedPowerID(), Equals, suite.spear.ID())
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieSkipsUnequippablePowersWhenDefaultEquipping(checker *C) {
	terosPowerReferences := []*powerreference.Reference{
		suite.blot.GetReference(),
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
	}
	suite.equipCheck.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	suite.equipCheck.EquipDefaultPower(suite.teros, suite.repos)
	checker.Assert(suite.teros.GetEquippedPowerID(), Equals, suite.spear.ID())
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieWillNotEquipIfNoEquippablePowers(checker *C) {
	terosPowerReferences := []*powerreference.Reference{
		suite.blot.GetReference(),
	}
	suite.equipCheck.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	suite.equipCheck.EquipDefaultPower(suite.teros, suite.repos)
	checker.Assert(suite.teros.HasEquippedPower(), Equals, false)
}

func (suite *SquaddieEquipPowersFromRepo) TestCanChangeEquippedPower(checker *C) {
	terosPowerReferences := []*powerreference.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	suite.equipCheck.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	suite.equipCheck.EquipDefaultPower(suite.teros, suite.repos)
	success := suite.equipCheck.SquaddieEquipPower(suite.teros, suite.scimitar.ID(), suite.repos)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.teros.GetEquippedPowerID(), Equals, suite.scimitar.ID())
}

func (suite *SquaddieEquipPowersFromRepo) TestFailToEquipUnequibbablePower(checker *C) {
	terosPowerReferences := []*powerreference.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	suite.equipCheck.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	suite.equipCheck.EquipDefaultPower(suite.teros, suite.repos)
	success := suite.equipCheck.SquaddieEquipPower(suite.teros, suite.blot.ID(), suite.repos)
	checker.Assert(success, Equals, false)
	checker.Assert(suite.teros.GetEquippedPowerID(), Equals, suite.spear.ID())
}

func (suite *SquaddieEquipPowersFromRepo) TestFailToEquipNonexistentPowers(checker *C) {
	success := suite.equipCheck.SquaddieEquipPower(suite.teros, "name", suite.repos)
	checker.Assert(success, Equals, false)
	checker.Assert(suite.teros.HasEquippedPower(), Equals, false)
}

func (suite *SquaddieEquipPowersFromRepo) TestFailToEquipUnownedPower(checker *C) {
	notTerosPower := power.NewPowerBuilder().CanBeEquipped().Build()

	suite.powerRepo.AddSlicePowerSource([]*power.Power{
		notTerosPower,
	})

	terosPowerReferences := []*powerreference.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	suite.equipCheck.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	suite.equipCheck.EquipDefaultPower(suite.teros, suite.repos)
	success := suite.equipCheck.SquaddieEquipPower(suite.teros, notTerosPower.ID(), suite.repos)
	checker.Assert(success, Equals, false)
	checker.Assert(suite.teros.GetEquippedPowerID(), Equals, suite.spear.ID())
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieCanCounter(checker *C) {
	terosPowerReferences := []*powerreference.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	suite.equipCheck.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	suite.equipCheck.EquipDefaultPower(suite.teros, suite.repos)
	canCounter, _ := suite.equipCheck.CanSquaddieCounterWithEquippedWeapon(suite.teros.ID(), suite.repos)
	checker.Assert(canCounter, Equals, true)
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieCannotCounterWithUncounterablePower(checker *C) {
	terosPowerReferences := []*powerreference.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	suite.equipCheck.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	suite.equipCheck.SquaddieEquipPower(suite.teros, suite.scimitar.ID(), suite.repos)
	canCounter, _ := suite.equipCheck.CanSquaddieCounterWithEquippedWeapon(suite.teros.ID(), suite.repos)
	checker.Assert(canCounter, Equals, false)
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieCannotCounterWithUnequippablePower(checker *C) {
	terosPowerReferences := []*powerreference.Reference{
		suite.blot.GetReference(),
	}
	suite.equipCheck.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	suite.equipCheck.EquipDefaultPower(suite.teros, suite.repos)
	canCounter, _ := suite.equipCheck.CanSquaddieCounterWithEquippedWeapon(suite.teros.ID(), suite.repos)
	checker.Assert(canCounter, Equals, false)
}
