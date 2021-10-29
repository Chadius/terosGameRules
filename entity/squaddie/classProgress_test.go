package squaddie_test

import (
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
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
	suite.mageClass = &squaddieclass.Class{ID: "1", Name: "Mage"}
	suite.mushroomClass = &squaddieclass.Class{ID: "2", Name: "Mushroom"}
}

func (suite *ClassProgressTests) TestNewSquaddieHasNoClassesOrLevels(checker *C) {
	checker.Assert(suite.teros.ClassProgress.CurrentClassID, Equals, "")
	checker.Assert(suite.teros.ClassProgress.GetLevelCountsByClass(), DeepEquals, map[string]int{})
}

func (suite *ClassProgressTests) TestChangeCurrentClass(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.mageClass)
	checker.Assert(suite.teros.ClassProgress.CurrentClassID, Equals, "")
	err := suite.teros.ClassProgress.SetClass(suite.mageClass.ID)
	checker.Assert(err, IsNil)
	checker.Assert(suite.teros.ClassProgress.CurrentClassID, Equals, suite.mageClass.ID)
}

func (suite *ClassProgressTests) TestCanSetBaseClass(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.mageClass)
	checker.Assert(suite.teros.ClassProgress.BaseClassID, Equals, "")
	suite.teros.ClassProgress.SetBaseClassIfNoBaseClass(suite.mageClass.ID)
	checker.Assert(suite.teros.ClassProgress.BaseClassID, Equals, suite.mageClass.ID)
}

func (suite *ClassProgressTests) TestRaiseErrorIfClassDoesNotExist(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.mageClass)
	checker.Assert(suite.teros.ClassProgress.CurrentClassID, Equals, "")
	err := suite.teros.ClassProgress.SetClass(suite.mushroomClass.ID)
	checker.Assert(err.Error(), Equals, `cannot switch to unknown class "2"`)
}

func (suite *ClassProgressTests) TestAddClassToSquaddie(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.mageClass)
	checker.Assert(suite.teros.ClassProgress.GetLevelCountsByClass(), DeepEquals, map[string]int{suite.mageClass.ID: 0})
}

func (suite *ClassProgressTests) TestCanTellIfSquaddieAddedClass(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.mageClass)
	checker.Assert(suite.teros.ClassProgress.HasAddedClass(suite.mageClass.ID), Equals, true)
	checker.Assert(suite.teros.ClassProgress.HasAddedClass(suite.mushroomClass.ID), Equals, false)
}
