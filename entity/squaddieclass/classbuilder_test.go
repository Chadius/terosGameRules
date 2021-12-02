package squaddieclass_test

import (
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	. "gopkg.in/check.v1"
)

type ClassBuilder struct{}

var _ = Suite(&ClassBuilder{})

func (suite *ClassBuilder) TestBuildClassWithID(checker *C) {
	dimensionalMage := squaddieclass.ClassBuilder().WithID("classDimensionalMage").Build()
	checker.Assert("classDimensionalMage", Equals, dimensionalMage.ID())
}

func (suite *ClassBuilder) TestBuildClassWithName(checker *C) {
	superSoldier := squaddieclass.ClassBuilder().WithName("super soldier").Build()
	checker.Assert("super soldier", Equals, superSoldier.Name())
}

func (suite *ClassBuilder) TestBuildClassThatRequiresABaseClass(checker *C) {
	advancedMage := squaddieclass.ClassBuilder().RequiresBaseClass().Build()
	checker.Assert(true, Equals, advancedMage.BaseClassRequired())
}

func (suite *ClassBuilder) TestBuildClassWithInitialBigLevel(checker *C) {
	multipleLevelClass := squaddieclass.ClassBuilder().WithInitialBigLevelID("class0").Build()
	checker.Assert("class0", Equals, multipleLevelClass.InitialBigLevelID())
}
