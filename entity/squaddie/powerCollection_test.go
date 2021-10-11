package squaddie_test

import (
	"fmt"
	"github.com/cserrant/terosbattleserver/entity/power"
	"github.com/cserrant/terosbattleserver/entity/squaddie"
	"github.com/cserrant/terosbattleserver/usecase/powerequip"
	"github.com/cserrant/terosbattleserver/usecase/repositories"
	. "gopkg.in/check.v1"
)

type SquaddiePowerCollectionTests struct {
	teros *squaddie.Squaddie
}

var _ = Suite(&SquaddiePowerCollectionTests{})

func (suite *SquaddiePowerCollectionTests) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
}

func (suite *SquaddiePowerCollectionTests) TestGainInnatePowers(checker *C) {
	attackA := power.NewPower("Attack Formation A")
	suite.teros.PowerCollection.AddInnatePower(attackA)

	attackIDNamePairs := suite.teros.PowerCollection.GetInnatePowerIDNames()
	checker.Assert(attackIDNamePairs, HasLen, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "Attack Formation A")
	checker.Assert(attackIDNamePairs[0].ID, Equals, attackA.ID)
}

func (suite *SquaddiePowerCollectionTests) TestClearInnatePowers(checker *C) {
	attackA := power.NewPower("Attack Formation A")
	suite.teros.PowerCollection.AddInnatePower(attackA)
	suite.teros.PowerCollection.ClearInnatePowers()

	attackIDNamePairs := suite.teros.PowerCollection.GetInnatePowerIDNames()
	checker.Assert(attackIDNamePairs, DeepEquals, []*power.Reference{})
}

func (suite *SquaddiePowerCollectionTests) TestClearPowerReferences(checker *C) {
	suite.teros.PowerCollection.PowerReferences = []*power.Reference{{Name: "Pow pow", ID: "Power Wheels"}}
	suite.teros.PowerCollection.ClearTemporaryPowerReferences()
	checker.Assert(suite.teros.PowerCollection.PowerReferences, DeepEquals, []*power.Reference{})
}

func (suite *SquaddiePowerCollectionTests) TestRaiseErrorIfTryToRegainSamePower(checker *C) {
	attackA := power.NewPower("Attack Formation A")
	err := suite.teros.PowerCollection.AddInnatePower(attackA)
	checker.Assert(err, IsNil)
	err = suite.teros.PowerCollection.AddInnatePower(attackA)
	expectedErrorMessage := fmt.Sprintf(`squaddie already has innate power with ID "%s"`, attackA.ID)
	checker.Assert(err, ErrorMatches, expectedErrorMessage)

	attackIDNamePairs := suite.teros.PowerCollection.GetInnatePowerIDNames()
	checker.Assert(attackIDNamePairs, HasLen, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "Attack Formation A")
	checker.Assert(attackIDNamePairs[0].ID, Equals, attackA.ID)
}

func (suite *SquaddiePowerCollectionTests) TestSquaddieHasEquippedPower(checker *C) {
	spear := power.NewPower("spear")
	spear.PowerType = power.Physical
	spear.AttackEffect = &power.AttackingEffect{
		CanBeEquipped: true,
	}

	powerRepo := power.NewPowerRepository()
	powerRepo.AddSlicePowerSource([]*power.Power{spear})

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			spear.GetReference(),
		},
		&repositories.RepositoryCollection{PowerRepo: powerRepo},
	)

	checker.Assert(suite.teros.PowerCollection.HasEquippedPower(), Equals, false)

	equippedSpearPower := powerequip.SquaddieEquipPower(suite.teros, spear.ID, &repositories.RepositoryCollection{PowerRepo: powerRepo})
	checker.Assert(equippedSpearPower, Equals, true)

	checker.Assert(suite.teros.PowerCollection.HasEquippedPower(), Equals, true)
	checker.Assert(suite.teros.PowerCollection.GetEquippedPowerID(), Equals, spear.ID)

	suite.teros.PowerCollection.ClearInnatePowers()
	checker.Assert(suite.teros.PowerCollection.HasEquippedPower(), Equals, false)
}
