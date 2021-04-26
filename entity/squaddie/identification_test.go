package squaddie_test

import (
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	. "gopkg.in/check.v1"
)

type SquaddieIdentificationCreationTests struct{
	teros *squaddie.Squaddie
	mageClass *squaddieclass.Class
	mushroomClass *squaddieclass.Class
}

var _ = Suite(&SquaddieIdentificationCreationTests {})

func (suite *SquaddieIdentificationCreationTests) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
	suite.mageClass = &squaddieclass.Class{ID: "1", Name: "Mage"}
	suite.mushroomClass = &squaddieclass.Class{ID: "2", Name: "Mushroom"}
}

func (suite *SquaddieIdentificationCreationTests) TestNameIsSet(checker *C) {
	checker.Assert(suite.teros.Identification.Name, Equals, "teros")
}

func (suite *SquaddieIdentificationCreationTests) TestGetARandomIDUponCreation(checker *C) {
	checker.Assert(suite.teros.Identification.ID, NotNil)
	checker.Assert(suite.teros.Identification.ID, Not(Equals), "")
}

func (suite *SquaddieIdentificationCreationTests) TestGetANewID(checker *C) {
	initialID := suite.teros.Identification.ID
	suite.teros.Identification.SetNewIDToRandom()
	checker.Assert(suite.teros.Identification.ID, Not(Equals), initialID)
}

func (suite *SquaddieMovementTests) TestRaisesErrorIfSquaddieHasUnknownAffiliation(checker *C) {
	newSquaddie := squaddie.NewSquaddie("teros")
	newSquaddie.Identification.Affiliation = "Unknown Affiliation"
	err := squaddie.CheckSquaddieForErrors(newSquaddie)
	checker.Assert(err, NotNil)
	checker.Assert(err, ErrorMatches,"Squaddie has unknown affiliation: 'Unknown Affiliation'")
}
