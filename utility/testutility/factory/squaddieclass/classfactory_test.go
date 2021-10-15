package squaddieclass_test

import (
	"github.com/chadius/terosbattleserver/utility/testutility/factory/squaddieclass"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ClassBuilder struct {}

var _ = Suite(&ClassBuilder{})

func (suite *ClassBuilder) TestBuildClassWithID(checker *C) {
	dimensionalMage := squaddieclass.ClassFactory().WithID("classDimensionalMage").Build()
	checker.Assert("classDimensionalMage", Equals, dimensionalMage.ID)
}

func (suite *ClassBuilder) TestBuildClassWithName(checker *C) {
	superSoldier := squaddieclass.ClassFactory().WithName("super soldier").Build()
	checker.Assert("super soldier", Equals, superSoldier.Name)
}

func (suite *ClassBuilder) TestBuildClassThatRequiresABaseClass(checker *C) {
	advancedMage := squaddieclass.ClassFactory().RequiresBaseClass().Build()
	checker.Assert(true, Equals, advancedMage.BaseClassRequired)
}

func (suite *ClassBuilder) TestBuildClassWithInitialBigLevel(checker *C) {
	multipleLevelClass := squaddieclass.ClassFactory().WithInitialBigLevelID("class0").Build()
	checker.Assert("class0", Equals, multipleLevelClass.InitialBigLevelID)
}
