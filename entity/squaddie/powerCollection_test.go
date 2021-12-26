package squaddie_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	. "gopkg.in/check.v1"
)

type SquaddiePowerCollectionTests struct {
	teros   *squaddie.Squaddie
	attackA *power.Power
}

var _ = Suite(&SquaddiePowerCollectionTests{})

func (suite *SquaddiePowerCollectionTests) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Build()
	suite.attackA = power.NewPowerBuilder().WithName("attack Formation A").Build()
}

func (suite *SquaddiePowerCollectionTests) TestAddPowerReference(checker *C) {
	suite.teros.AddPowerReference(suite.attackA.GetReference())

	attackIDNamePairs := suite.teros.GetCopyOfPowerReferences()
	checker.Assert(attackIDNamePairs, HasLen, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "attack Formation A")
	checker.Assert(attackIDNamePairs[0].PowerID, Equals, suite.attackA.ID())
}

func (suite *SquaddiePowerCollectionTests) TestAddPowerReferenceIsIdempotent(checker *C) {
	suite.teros.AddPowerReference(suite.attackA.GetReference())
	suite.teros.AddPowerReference(suite.attackA.GetReference())

	attackIDNamePairs := suite.teros.GetCopyOfPowerReferences()
	checker.Assert(attackIDNamePairs, HasLen, 1)
}

func (suite *SquaddiePowerCollectionTests) TestRemovePowerReference(checker *C) {
	suite.teros.AddPowerReference(suite.attackA.GetReference())
	suite.teros.RemovePowerReferenceByPowerID(suite.attackA.ID())

	attackIDNamePairs := suite.teros.GetCopyOfPowerReferences()
	checker.Assert(attackIDNamePairs, HasLen, 0)
}

func (suite *SquaddiePowerCollectionTests) TestRemovePowerReferenceIsIdempotent(checker *C) {
	suite.teros.AddPowerReference(suite.attackA.GetReference())
	suite.teros.RemovePowerReferenceByPowerID(suite.attackA.ID())
	suite.teros.RemovePowerReferenceByPowerID(suite.attackA.ID())

	attackIDNamePairs := suite.teros.GetCopyOfPowerReferences()
	checker.Assert(attackIDNamePairs, HasLen, 0)
}

func (suite *SquaddiePowerCollectionTests) TestClearInnatePowers(checker *C) {
	suite.teros.AddPowerReference(suite.attackA.GetReference())
	suite.teros.ClearPowerReferences()

	attackIDNamePairs := suite.teros.GetCopyOfPowerReferences()
	checker.Assert(attackIDNamePairs, DeepEquals, []*power.Reference{})
}

func (suite *SquaddiePowerCollectionTests) TestSquaddieHasEquippedPower(checker *C) {
	spear := power.NewPowerBuilder().Spear().Build()

	powerRepo := powerrepository.NewPowerRepository()
	powerRepo.AddSlicePowerSource([]*power.Power{spear})
	checkEquip := powerequip.CheckRepositories{}
	checkEquip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			spear.GetReference(),
		},
		&repositories.RepositoryCollection{PowerRepo: powerRepo},
	)

	checker.Assert(suite.teros.HasEquippedPower(), Equals, false)

	equippedSpearPower := checkEquip.SquaddieEquipPower(suite.teros, spear.ID(), &repositories.RepositoryCollection{PowerRepo: powerRepo})
	checker.Assert(equippedSpearPower, Equals, true)

	checker.Assert(suite.teros.HasEquippedPower(), Equals, true)
	checker.Assert(suite.teros.GetEquippedPowerID(), Equals, spear.ID())

	suite.teros.ClearPowerReferences()
	checker.Assert(suite.teros.HasEquippedPower(), Equals, false)
}
