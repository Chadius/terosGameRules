package squaddie_test

import (
	"github.com/chadius/terosgamerules/entity/squaddie"
	"github.com/chadius/terosgamerules/entity/squaddieclass"
	"github.com/chadius/terosgamerules/entity/squaddieinterface"
	. "gopkg.in/check.v1"
)

type SquaddieIdentificationCreationTests struct {
	teros         squaddieinterface.Interface
	mageClass     *squaddieclass.Class
	mushroomClass *squaddieclass.Class
}

var _ = Suite(&SquaddieIdentificationCreationTests{})

func (suite *SquaddieIdentificationCreationTests) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().WithName("teros").Build()
	suite.mageClass = squaddieclass.ClassBuilder().WithID("mageClassID").Build()
	suite.mushroomClass = squaddieclass.ClassBuilder().WithID("mushroomClassID").Build()
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
	suite.teros.SetNewIDToRandom()
	checker.Assert(suite.teros.ID(), Not(Equals), initialID)
}
