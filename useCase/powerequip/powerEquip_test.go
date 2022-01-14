package powerequip_test

import (
	"github.com/chadius/terosgamerules/entity/power"
	"github.com/chadius/terosgamerules/entity/powerinterface"
	"github.com/chadius/terosgamerules/entity/powerreference"
	"github.com/chadius/terosgamerules/entity/powerrepository"
	"github.com/chadius/terosgamerules/entity/squaddie"
	"github.com/chadius/terosgamerules/entity/squaddieinterface"
	"github.com/chadius/terosgamerules/usecase/powerequip"
	"github.com/chadius/terosgamerules/usecase/repositories"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type SquaddieEquipPowersFromRepo struct {
	teros        squaddieinterface.Interface
	spear        powerinterface.Interface
	scimitar     powerinterface.Interface
	blot         powerinterface.Interface
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
	suite.powerRepo.AddSlicePowerSource([]powerinterface.Interface{
		suite.spear,
		suite.scimitar,
		suite.blot,
	})

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]squaddieinterface.Interface{suite.teros})

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

	suite.powerRepo.AddSlicePowerSource([]powerinterface.Interface{
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
