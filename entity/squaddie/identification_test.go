package squaddie_test

import (
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	squaddieClassBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddieclass"
	. "gopkg.in/check.v1"
)

type SquaddieIdentificationCreationTests struct {
	teros         *squaddie.Squaddie
	mageClass     *squaddieclass.Class
	mushroomClass *squaddieclass.Class
}

var _ = Suite(&SquaddieIdentificationCreationTests{})

func (suite *SquaddieIdentificationCreationTests) SetUpTest(checker *C) {
	suite.teros = squaddieBuilder.Builder().Teros().WithName("teros").Build()
	suite.mageClass = squaddieClassBuilder.ClassBuilder().WithID("mageClassID").Build()
	suite.mushroomClass = squaddieClassBuilder.ClassBuilder().WithID("mushroomClassID").Build()
}

func (suite *SquaddieIdentificationCreationTests) TestNameIsSet(checker *C) {
	checker.Assert(suite.teros.Name(), Equals, "teros")
}

func (suite *SquaddieIdentificationCreationTests) TestGetARandomIDUponCreation(checker *C) {
	checker.Assert(suite.teros.ID(), NotNil)
	checker.Assert(suite.teros.ID(), Not(Equals), "")
}

func (suite *SquaddieIdentificationCreationTests) TestGetANewID(checker *C) {
	initialID := suite.teros.ID()
	suite.teros.Identification.SetNewIDToRandom()
	checker.Assert(suite.teros.ID(), Not(Equals), initialID)
}

func (suite *SquaddieMovementTests) TestRaisesErrorIfSquaddieHasUnknownAffiliation(checker *C) {
	newSquaddie := squaddie.NewSquaddie("teros")
	newSquaddie.Identification.SquaddieAffiliation = "Unknown SquaddieAffiliation"
	newSquaddie.Identification.SquaddieID = "squaddieTeros"
	err := squaddie.CheckSquaddieForErrors(newSquaddie)
	checker.Assert(err, NotNil)
	checker.Assert(err, ErrorMatches, "squaddie squaddieTeros has unknown affiliation: 'Unknown SquaddieAffiliation'")
}
