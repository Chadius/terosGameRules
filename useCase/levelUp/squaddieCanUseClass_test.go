package levelup_test

import (
	"github.com/chadius/terosgamerules/entity/levelupbenefit"
	"github.com/chadius/terosgamerules/entity/squaddie"
	"github.com/chadius/terosgamerules/entity/squaddieclass"
	"github.com/chadius/terosgamerules/entity/squaddieinterface"
	"github.com/chadius/terosgamerules/usecase/levelup"
	"github.com/chadius/terosgamerules/usecase/repositories"
	"github.com/chadius/terosgamerules/utility/testutility/builder"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type SquaddieQualifiesForClassSuite struct {
	classRepo *squaddieclass.Repository
	levelRepo *levelupbenefit.Repository
	repos     *repositories.RepositoryCollection

	mageClass  *squaddieclass.Class
	mageLevel0 *levelupbenefit.LevelUpBenefit
	mageLevel1 *levelupbenefit.LevelUpBenefit

	dimensionWalkerClass  *squaddieclass.Class
	dimensionWalkerLevel0 *levelupbenefit.LevelUpBenefit
	dimensionWalkerLevel1 *levelupbenefit.LevelUpBenefit

	ancientTomeClass       *squaddieclass.Class
	ancientTomeClassLevel0 *levelupbenefit.LevelUpBenefit

	atLeastTenLevelsBaseClass *squaddieclass.Class
	lotsOfLevels              []*levelupbenefit.LevelUpBenefit

	teros                   squaddieinterface.Interface
	improveSquaddieStrategy levelup.ImproveSquaddieStrategy
	levelUpCheck            levelup.SquaddieCanSwitchClassStrategy
}

var _ = Suite(&SquaddieQualifiesForClassSuite{})

func (suite *SquaddieQualifiesForClassSuite) SetUpTest(checker *C) {
	suite.mageClass = squaddieclass.ClassBuilder().WithID("class1").WithName("Mage").Build()
	suite.dimensionWalkerClass = squaddieclass.ClassBuilder().WithID("class2").WithName("Dimension Walker").RequiresBaseClass().Build()
	suite.ancientTomeClass = squaddieclass.ClassBuilder().WithID("class3").WithName("Ancient Tome").RequiresBaseClass().Build()
	suite.atLeastTenLevelsBaseClass = squaddieclass.ClassBuilder().WithID("class4").WithName("Base with many levels").Build()

	suite.classRepo = squaddieclass.NewRepository()
	suite.classRepo.AddListOfClasses([]*squaddieclass.Class{suite.mageClass, suite.dimensionWalkerClass, suite.ancientTomeClass})

	suite.mageLevel0, _ = levelupbenefit.NewLevelUpBenefitBuilder().
		WithID("mageLevel0").
		WithClassID(suite.mageClass.ID()).
		HitPoints(1).
		Dodge(1).
		Deflect(1).
		Barrier(1).
		Armor(1).
		Aim(1).
		Strength(1).
		Mind(1).
		Build()

	suite.mageLevel1, _ = levelupbenefit.NewLevelUpBenefitBuilder().
		WithID("mageLevel1").
		WithClassID(suite.mageClass.ID()).
		BigLevel().
		MovementDistance(1).
		Build()

	suite.dimensionWalkerLevel0, _ = levelupbenefit.NewLevelUpBenefitBuilder().
		WithID("dwLevel0").
		WithClassID(suite.dimensionWalkerClass.ID()).
		BigLevel().
		MovementLogic("light").
		MovementDistance(1).
		Build()

	suite.dimensionWalkerLevel1, _ = levelupbenefit.NewLevelUpBenefitBuilder().
		WithID("dwLevel1").
		WithClassID(suite.dimensionWalkerClass.ID()).
		BigLevel().
		HitPoints(1).
		Build()

	suite.ancientTomeClassLevel0, _ = levelupbenefit.NewLevelUpBenefitBuilder().
		WithID("ancientTomeLevel0").
		WithClassID(suite.ancientTomeClass.ID()).
		Dodge(1).
		Deflect(1).
		Mind(1).
		Build()

	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	suite.levelRepo.AddLevels([]*levelupbenefit.LevelUpBenefit{
		suite.mageLevel0, suite.mageLevel1, suite.dimensionWalkerLevel0, suite.dimensionWalkerLevel1, suite.ancientTomeClassLevel0,
	})

	suite.lotsOfLevels = append(
		(&builder.LevelGenerator{
			Instructions: &builder.LevelGeneratorInstruction{
				NumberOfLevels: 11,
				ClassID:        suite.atLeastTenLevelsBaseClass.ID(),
				PrefixLevelID:  "lotsLevelsSmall",
				Type:           levelupbenefit.Small,
			},
		}).Build(),
		(&builder.LevelGenerator{
			Instructions: &builder.LevelGeneratorInstruction{
				NumberOfLevels: 11,
				ClassID:        suite.atLeastTenLevelsBaseClass.ID(),
				PrefixLevelID:  "lotsLevelsBig",
				Type:           levelupbenefit.Big,
			},
		}).Build()...,
	)

	suite.levelRepo.AddLevels(suite.lotsOfLevels)
	suite.repos = &repositories.RepositoryCollection{
		ClassRepo: suite.classRepo,
		LevelRepo: suite.levelRepo,
	}

	suite.teros = squaddie.NewSquaddieBuilder().Teros().AddClassByReference(suite.mageClass.GetReference()).AddClassByReference(suite.dimensionWalkerClass.GetReference()).AddClassByReference(suite.ancientTomeClass.GetReference()).AddClassByReference(suite.atLeastTenLevelsBaseClass.GetReference()).Build()
	suite.improveSquaddieStrategy = &levelup.ImproveSquaddieClass{}
	suite.levelUpCheck = &levelup.LevelsConsumedChecker{}
}

func (suite *SquaddieQualifiesForClassSuite) TestNewSquaddieCanSwitchToBaseClass(checker *C) {
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID(), suite.repos), Equals, true)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.dimensionWalkerClass.ID(), suite.repos), Equals, false)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.ancientTomeClass.ID(), suite.repos), Equals, false)
}

func (suite *SquaddieQualifiesForClassSuite) TestSquaddieCannotSwitchToCurrentClass(checker *C) {
	suite.teros.SetBaseClassIfNoBaseClass(suite.mageClass.ID())
	suite.teros.SetClass(suite.mageClass.ID())
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID(), suite.repos), Equals, false)
}

func (suite *SquaddieQualifiesForClassSuite) TestSquaddieMustFinishBaseClassBeforeSwitching(checker *C) {
	suite.teros.SetBaseClassIfNoBaseClass(suite.mageClass.ID())
	suite.teros.SetClass(suite.mageClass.ID())
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.mageLevel0, suite.teros)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID(), suite.repos), Equals, false)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.dimensionWalkerClass.ID(), suite.repos), Equals, false)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.ancientTomeClass.ID(), suite.repos), Equals, false)
}

func (suite *SquaddieQualifiesForClassSuite) TestSquaddieCanSwitchFromFinishedBaseToAdvancedClass(checker *C) {
	suite.teros.SetBaseClassIfNoBaseClass(suite.mageClass.ID())
	suite.teros.SetClass(suite.mageClass.ID())
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.mageLevel0, suite.teros)
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.mageLevel1, suite.teros)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID(), suite.repos), Equals, false)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.dimensionWalkerClass.ID(), suite.repos), Equals, true)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.ancientTomeClass.ID(), suite.repos), Equals, true)
}

func (suite *SquaddieQualifiesForClassSuite) TestSquaddieCanSwitchToNewAdvancedClassAtLevelTen(checker *C) {
	suite.teros.SetBaseClassIfNoBaseClass(suite.mageClass.ID())
	suite.teros.SetClass(suite.mageClass.ID())
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.mageLevel0, suite.teros)
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.mageLevel1, suite.teros)
	suite.teros.SetClass(suite.dimensionWalkerClass.ID())
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID(), suite.repos), Equals, false)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.dimensionWalkerClass.ID(), suite.repos), Equals, false)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.ancientTomeClass.ID(), suite.repos), Equals, false)
}

func (suite *SquaddieQualifiesForClassSuite) TestCannotSwitchToCompletedClass(checker *C) {
	suite.teros.SetBaseClassIfNoBaseClass(suite.mageClass.ID())
	suite.teros.SetClass(suite.mageClass.ID())
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.mageLevel0, suite.teros)
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.mageLevel1, suite.teros)
	suite.teros.SetClass(suite.dimensionWalkerClass.ID())
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.dimensionWalkerLevel0, suite.teros)
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.dimensionWalkerLevel1, suite.teros)
	suite.teros.SetClass(suite.ancientTomeClass.ID())
	checker.Assert(suite.teros.CurrentClassID(), Equals, suite.ancientTomeClass.ID())
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID(), suite.repos), Equals, false)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.dimensionWalkerClass.ID(), suite.repos), Equals, false)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.ancientTomeClass.ID(), suite.repos), Equals, false)
}

func (suite *SquaddieQualifiesForClassSuite) TestQualifyForAdvancedClassesWhenBaseClassCompletes(checker *C) {
	suite.teros.SetBaseClassIfNoBaseClass(suite.mageClass.ID())
	suite.teros.SetClass(suite.mageClass.ID())
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.mageLevel0, suite.teros)
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.mageLevel1, suite.teros)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID(), suite.repos), Equals, false)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.dimensionWalkerClass.ID(), suite.repos), Equals, true)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.ancientTomeClass.ID(), suite.repos), Equals, true)
}

func (suite *SquaddieQualifiesForClassSuite) TestSquaddieStaysInAdvancedClassUntilCompletion(checker *C) {
	suite.teros.SetBaseClassIfNoBaseClass(suite.mageClass.ID())
	suite.teros.SetClass(suite.mageClass.ID())
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.mageLevel0, suite.teros)
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.mageLevel1, suite.teros)
	suite.teros.SetClass(suite.dimensionWalkerClass.ID())
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.dimensionWalkerLevel0, suite.teros)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID(), suite.repos), Equals, false)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.dimensionWalkerClass.ID(), suite.repos), Equals, false)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.ancientTomeClass.ID(), suite.repos), Equals, false)
}

func (suite *SquaddieQualifiesForClassSuite) TestCanSwitchClassAfterTenLevels(checker *C) {
	suite.teros.SetBaseClassIfNoBaseClass(suite.atLeastTenLevelsBaseClass.ID())
	suite.teros.SetClass(suite.atLeastTenLevelsBaseClass.ID())
	for index, _ := range [10]int{} {
		suite.improveSquaddieStrategy.ImproveSquaddie(suite.lotsOfLevels[index], suite.teros)
	}
	checker.Assert(suite.teros.GetLevelCountsByClass()[suite.atLeastTenLevelsBaseClass.ID()], Equals, 10)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.mageClass.ID(), suite.repos), Equals, true)
	checker.Assert(suite.levelUpCheck.SquaddieCanSwitchToClass(suite.teros, suite.dimensionWalkerClass.ID(), suite.repos), Equals, true)
}
