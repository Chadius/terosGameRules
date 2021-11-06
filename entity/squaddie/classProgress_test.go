package squaddie_test

import (
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	squaddieClassBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddieclass"
	. "gopkg.in/check.v1"
)

type ClassProgressTests struct {
	teros         *squaddie.Squaddie
	mageClass     *squaddieclass.Class
	mushroomClass *squaddieclass.Class
}

var _ = Suite(&ClassProgressTests{})

func (suite *ClassProgressTests) SetUpTest(checker *C) {
	suite.teros = squaddieBuilder.Builder().Teros().Build()
	suite.mageClass = squaddieClassBuilder.ClassBuilder().WithID("mageClassID").Build()
	suite.mushroomClass = squaddieClassBuilder.ClassBuilder().WithID("mushroomClassID").Build()
}

func (suite *ClassProgressTests) TestNewSquaddieHasNoClassesOrLevels(checker *C) {
	checker.Assert(suite.teros.ClassProgress.CurrentClassID, Equals, "")
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{})
}

func (suite *ClassProgressTests) TestAddClassDoesNotOverwriteExistingClass(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.mageClass.GetReference())
	suite.teros.SetBaseClassIfNoBaseClass(suite.mageClass.ID())
	suite.teros.MarkLevelUpBenefitAsConsumed(suite.mageClass.ID(), "classLevelAlreadyConsumed")
	checker.Assert(suite.teros.IsClassLevelAlreadyUsed("classLevelAlreadyConsumed"), Equals, true)

	suite.teros.ClassProgress.AddClass(suite.mageClass.GetReference())

	checker.Assert(suite.teros.IsClassLevelAlreadyUsed("classLevelAlreadyConsumed"), Equals, true)
}

func (suite *ClassProgressTests) TestChangeCurrentClass(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.mageClass.GetReference())
	checker.Assert(suite.teros.ClassProgress.CurrentClassID, Equals, "")
	err := suite.teros.SetClass(suite.mageClass.ID())
	checker.Assert(err, IsNil)
	checker.Assert(suite.teros.ClassProgress.CurrentClassID, Equals, suite.mageClass.ID())
}

func (suite *ClassProgressTests) TestCanSetBaseClass(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.mageClass.GetReference())
	checker.Assert(suite.teros.ClassProgress.BaseClassID, Equals, "")
	suite.teros.SetBaseClassIfNoBaseClass(suite.mageClass.ID())
	checker.Assert(suite.teros.ClassProgress.BaseClassID, Equals, suite.mageClass.ID())
}

func (suite *ClassProgressTests) TestRaiseErrorIfClassDoesNotExist(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.mageClass.GetReference())
	checker.Assert(suite.teros.ClassProgress.CurrentClassID, Equals, "")
	err := suite.teros.SetClass(suite.mushroomClass.ID())
	checker.Assert(err.Error(), Equals, `cannot switch to unknown class "mushroomClassID"`)
}

func (suite *ClassProgressTests) TestAddClassToSquaddie(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.mageClass.GetReference())
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{suite.mageClass.ID(): 0})
}

func (suite *ClassProgressTests) TestCanTellIfSquaddieAddedClass(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.mageClass.GetReference())
	checker.Assert(suite.teros.HasAddedClass(suite.mageClass.ID()), Equals, true)
	checker.Assert(suite.teros.HasAddedClass(suite.mushroomClass.ID()), Equals, false)
}
