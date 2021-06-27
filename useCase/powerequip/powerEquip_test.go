package powerequip_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type SquaddieEquipPowersFromRepo struct {
	teros *squaddie.Squaddie
	spear *power.Power
	scimitar *power.Power
	blot *power.Power
	powerRepo *power.Repository
	squaddieRepo *squaddie.Repository
	repos *repositories.RepositoryCollection
}

var _ = Suite(&SquaddieEquipPowersFromRepo{})

func (suite *SquaddieEquipPowersFromRepo) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("Teros")
	suite.spear = power.NewPower("Spear")
	suite.spear.AttackEffect = &power.AttackingEffect{
		CanBeEquipped: true,
		CanCounterAttack: true,
	}

	suite.scimitar = power.NewPower("scimitar the second")
	suite.scimitar.AttackEffect = &power.AttackingEffect{
		CanBeEquipped: true,
	}

	suite.blot = power.NewPower("Magic spell Blot")

	suite.powerRepo = power.NewPowerRepository()
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
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieEquipsFirstPowerByDefault (checker *C) {
	terosPowerReferences := []*power.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	powerequip.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	powerequip.EquipDefaultPower(suite.teros, suite.repos)
	checker.Assert(suite.teros.PowerCollection.GetEquippedPowerID(), Equals, suite.spear.ID)
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieSkipsUnequippablePowersWhenDefaultEquipping (checker *C) {
	terosPowerReferences := []*power.Reference{
		suite.blot.GetReference(),
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
	}
	powerequip.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	powerequip.EquipDefaultPower(suite.teros, suite.repos)
	checker.Assert(suite.teros.PowerCollection.GetEquippedPowerID(), Equals, suite.spear.ID)
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieWillNotEquipIfNoEquippablePowers (checker *C) {
	terosPowerReferences := []*power.Reference{
		suite.blot.GetReference(),
	}
	powerequip.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	powerequip.EquipDefaultPower(suite.teros, suite.repos)
	checker.Assert(suite.teros.PowerCollection.HasEquippedPower(), Equals, false)
}

func (suite *SquaddieEquipPowersFromRepo) TestCanChangeEquippedPower (checker *C) {
	terosPowerReferences := []*power.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	powerequip.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	powerequip.EquipDefaultPower(suite.teros, suite.repos)
	success := powerequip.SquaddieEquipPower(suite.teros, suite.scimitar.ID, suite.repos)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.teros.PowerCollection.GetEquippedPowerID(), Equals, suite.scimitar.ID)
}

func (suite *SquaddieEquipPowersFromRepo) TestFailToEquipUnequibbablePower (checker *C) {
	terosPowerReferences := []*power.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	powerequip.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	powerequip.EquipDefaultPower(suite.teros, suite.repos)
	success := powerequip.SquaddieEquipPower(suite.teros, suite.blot.ID, suite.repos)
	checker.Assert(success, Equals, false)
	checker.Assert(suite.teros.PowerCollection.GetEquippedPowerID(), Equals, suite.spear.ID)
}

func (suite *SquaddieEquipPowersFromRepo) TestFailToEquipNonexistentPowers (checker *C) {
	success := powerequip.SquaddieEquipPower(suite.teros, "kwyjibo", suite.repos)
	checker.Assert(success, Equals, false)
	checker.Assert(suite.teros.PowerCollection.HasEquippedPower(), Equals, false)
}

func (suite *SquaddieEquipPowersFromRepo) TestFailToEquipUnownedPower (checker *C) {
	notTerosPower := power.NewPower("Does not belong to Teros")
	notTerosPower.AttackEffect = &power.AttackingEffect{
		CanBeEquipped: true,
	}
	suite.powerRepo.AddSlicePowerSource([]*power.Power{
		notTerosPower,
	})

	terosPowerReferences := []*power.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	powerequip.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	powerequip.EquipDefaultPower(suite.teros, suite.repos)
	success := powerequip.SquaddieEquipPower(suite.teros, notTerosPower.ID, suite.repos)
	checker.Assert(success, Equals, false)
	checker.Assert(suite.teros.PowerCollection.GetEquippedPowerID(), Equals, suite.spear.ID)
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieCanCounter (checker *C) {
	terosPowerReferences := []*power.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	powerequip.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	powerequip.EquipDefaultPower(suite.teros, suite.repos)
	canCounter, _ := powerequip.CanSquaddieCounterWithEquippedWeapon(suite.teros.Identification.ID, suite.repos)
	checker.Assert(canCounter, Equals, true)
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieCannotCounterWithUncounterablePower (checker *C) {
	terosPowerReferences := []*power.Reference{
		suite.spear.GetReference(),
		suite.scimitar.GetReference(),
		suite.blot.GetReference(),
	}
	powerequip.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	powerequip.SquaddieEquipPower(suite.teros, suite.scimitar.ID, suite.repos)
	canCounter, _ := powerequip.CanSquaddieCounterWithEquippedWeapon(suite.teros.Identification.ID, suite.repos)
	checker.Assert(canCounter, Equals, false)
}

func (suite *SquaddieEquipPowersFromRepo) TestSquaddieCannotCounterWithUnequippablePower (checker *C) {
	terosPowerReferences := []*power.Reference{
		suite.blot.GetReference(),
	}
	powerequip.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.repos)
	powerequip.EquipDefaultPower(suite.teros, suite.repos)
	canCounter, _ := powerequip.CanSquaddieCounterWithEquippedWeapon(suite.teros.Identification.ID, suite.repos)
	checker.Assert(canCounter, Equals, false)
}
