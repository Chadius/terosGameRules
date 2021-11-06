package squaddie_test

import (
	"fmt"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	powerBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	. "gopkg.in/check.v1"
)

type SquaddiePowerCollectionTests struct {
	teros   *squaddie.Squaddie
	attackA *power.Power
}

var _ = Suite(&SquaddiePowerCollectionTests{})

func (suite *SquaddiePowerCollectionTests) SetUpTest(checker *C) {
	suite.teros = squaddieBuilder.Builder().Teros().Build()
	suite.attackA = powerBuilder.Builder().WithName("Attack Formation A").Build()
}

func (suite *SquaddiePowerCollectionTests) TestGainInnatePowers(checker *C) {
	suite.teros.PowerCollection.AddInnatePower(suite.attackA)

	attackIDNamePairs := suite.teros.PowerCollection.GetInnatePowerIDNames()
	checker.Assert(attackIDNamePairs, HasLen, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "Attack Formation A")
	checker.Assert(attackIDNamePairs[0].PowerID, Equals, suite.attackA.ID())
}

func (suite *SquaddiePowerCollectionTests) TestClearInnatePowers(checker *C) {
	suite.teros.PowerCollection.AddInnatePower(suite.attackA)
	suite.teros.PowerCollection.ClearInnatePowers()

	attackIDNamePairs := suite.teros.PowerCollection.GetInnatePowerIDNames()
	checker.Assert(attackIDNamePairs, DeepEquals, []*power.Reference{})
}

func (suite *SquaddiePowerCollectionTests) TestClearPowerReferences(checker *C) {
	suite.teros.PowerCollection.PowerReferences = []*power.Reference{{Name: "Pow pow", PowerID: "Power Wheels"}}
	suite.teros.PowerCollection.ClearTemporaryPowerReferences()
	checker.Assert(suite.teros.PowerCollection.PowerReferences, DeepEquals, []*power.Reference{})
}

func (suite *SquaddiePowerCollectionTests) TestRaiseErrorIfTryToRegainSamePower(checker *C) {
	err := suite.teros.PowerCollection.AddInnatePower(suite.attackA)
	checker.Assert(err, IsNil)
	err = suite.teros.PowerCollection.AddInnatePower(suite.attackA)
	expectedErrorMessage := fmt.Sprintf(`squaddie already has innate power with SquaddieID "%s"`, suite.attackA.ID())
	checker.Assert(err, ErrorMatches, expectedErrorMessage)

	attackIDNamePairs := suite.teros.PowerCollection.GetInnatePowerIDNames()
	checker.Assert(attackIDNamePairs, HasLen, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, suite.attackA.Name())
	checker.Assert(attackIDNamePairs[0].PowerID, Equals, suite.attackA.ID())
}

func (suite *SquaddiePowerCollectionTests) TestSquaddieHasEquippedPower(checker *C) {
	spear := powerBuilder.Builder().Spear().Build()

	powerRepo := powerrepository.NewPowerRepository()
	powerRepo.AddSlicePowerSource([]*power.Power{spear})

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			spear.GetReference(),
		},
		&repositories.RepositoryCollection{PowerRepo: powerRepo},
	)

	checker.Assert(suite.teros.PowerCollection.HasEquippedPower(), Equals, false)

	equippedSpearPower := powerequip.SquaddieEquipPower(suite.teros, spear.ID(), &repositories.RepositoryCollection{PowerRepo: powerRepo})
	checker.Assert(equippedSpearPower, Equals, true)

	checker.Assert(suite.teros.PowerCollection.HasEquippedPower(), Equals, true)
	checker.Assert(suite.teros.PowerCollection.GetEquippedPowerID(), Equals, spear.ID())

	suite.teros.PowerCollection.ClearInnatePowers()
	checker.Assert(suite.teros.PowerCollection.HasEquippedPower(), Equals, false)
}
