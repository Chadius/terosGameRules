package squaddie_test

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	. "gopkg.in/check.v1"
)

type SquaddiePowerCollectionTests struct{
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
