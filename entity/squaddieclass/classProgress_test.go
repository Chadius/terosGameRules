package squaddieclass_test

import (
	"github.com/chadius/terosgamerules/entity/squaddie"
	"github.com/chadius/terosgamerules/entity/squaddieclass"
	"github.com/chadius/terosgamerules/entity/squaddieinterface"
	. "gopkg.in/check.v1"
)

type ClassProgressTests struct {
	teros         squaddieinterface.Interface
	mageClass     *squaddieclass.Class
	mushroomClass *squaddieclass.Class
}

var _ = Suite(&ClassProgressTests{})

func (suite *ClassProgressTests) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Build()
	suite.mageClass = squaddieclass.ClassBuilder().WithID("mageClassID").Build()
	suite.mushroomClass = squaddieclass.ClassBuilder().WithID("mushroomClassID").Build()
}

func (suite *ClassProgressTests) TestNewSquaddieHasNoClassesOrLevels(checker *C) {
	checker.Assert(suite.teros.CurrentClassID(), Equals, "")
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{})
}

func (suite *ClassProgressTests) TestAddClassDoesNotOverwriteExistingClass(checker *C) {
	suite.teros.AddClass(suite.mageClass.GetReference())
	suite.teros.SetBaseClassIfNoBaseClass(suite.mageClass.ID())
	suite.teros.MarkLevelUpBenefitAsConsumed(suite.mageClass.ID(), "classLevelAlreadyConsumed")
	checker.Assert(suite.teros.IsClassLevelAlreadyUsed("classLevelAlreadyConsumed"), Equals, true)

	suite.teros.AddClass(suite.mageClass.GetReference())

	checker.Assert(suite.teros.IsClassLevelAlreadyUsed("classLevelAlreadyConsumed"), Equals, true)
}

func (suite *ClassProgressTests) TestChangeCurrentClass(checker *C) {
	suite.teros.AddClass(suite.mageClass.GetReference())
	checker.Assert(suite.teros.CurrentClassID(), Equals, "")
	err := suite.teros.SetClass(suite.mageClass.ID())
	checker.Assert(err, IsNil)
	checker.Assert(suite.teros.CurrentClassID(), Equals, suite.mageClass.ID())
}

func (suite *ClassProgressTests) TestCanSetBaseClass(checker *C) {
	suite.teros.AddClass(suite.mageClass.GetReference())
	checker.Assert(suite.teros.BaseClassID(), Equals, "")
	suite.teros.SetBaseClassIfNoBaseClass(suite.mageClass.ID())
	checker.Assert(suite.teros.BaseClassID(), Equals, suite.mageClass.ID())
}

func (suite *ClassProgressTests) TestRaiseErrorIfClassDoesNotExist(checker *C) {
	suite.teros.AddClass(suite.mageClass.GetReference())
	checker.Assert(suite.teros.CurrentClassID(), Equals, "")
	err := suite.teros.SetClass(suite.mushroomClass.ID())
	checker.Assert(err.Error(), Equals, `cannot switch to unknown class "mushroomClassID"`)
}

func (suite *ClassProgressTests) TestAddClassToSquaddie(checker *C) {
	suite.teros.AddClass(suite.mageClass.GetReference())
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{suite.mageClass.ID(): 0})
}

func (suite *ClassProgressTests) TestCanTellIfSquaddieAddedClass(checker *C) {
	suite.teros.AddClass(suite.mageClass.GetReference())
	checker.Assert(suite.teros.HasAddedClass(suite.mageClass.ID()), Equals, true)
	checker.Assert(suite.teros.HasAddedClass(suite.mushroomClass.ID()), Equals, false)
}
