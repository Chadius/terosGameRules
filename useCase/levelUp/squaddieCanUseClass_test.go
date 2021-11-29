package levelup_test

import (
	"github.com/chadius/terosbattleserver/entity/levelupbenefit"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	"github.com/chadius/terosbattleserver/usecase/levelup"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	squaddieClassBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddieclass"
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

	teros                   *squaddie.Squaddie
	improveSquaddieStrategy levelup.ImproveSquaddieStrategy
	levelUpCheck            levelup.SquaddieCanSwitchClassStrategy
}

var _ = Suite(&SquaddieQualifiesForClassSuite{})

func (suite *SquaddieQualifiesForClassSuite) SetUpTest(checker *C) {
	suite.mageClass = squaddieClassBuilder.ClassBuilder().WithID("class1").WithName("Mage").Build()
	suite.dimensionWalkerClass = squaddieClassBuilder.ClassBuilder().WithID("class2").WithName("Dimension Walker").RequiresBaseClass().Build()
	suite.ancientTomeClass = squaddieClassBuilder.ClassBuilder().WithID("class3").WithName("Ancient Tome").RequiresBaseClass().Build()
	suite.atLeastTenLevelsBaseClass = squaddieClassBuilder.ClassBuilder().WithID("class4").WithName("Base with many levels").Build()

	suite.classRepo = squaddieclass.NewRepository()
	suite.classRepo.AddListOfClasses([]*squaddieclass.Class{suite.mageClass, suite.dimensionWalkerClass, suite.ancientTomeClass})

	suite.mageLevel0 = &levelupbenefit.LevelUpBenefit{
		Identification: levelupbenefit.NewIdentification("mageLevel0", suite.mageClass.ID(), levelupbenefit.Small),
		Defense: levelupbenefit.NewDefense(1,1,1,1,1),
		Offense: levelupbenefit.NewOffense(1,1,1),
		PowerChanges: &levelupbenefit.PowerChanges{
			Gained: nil,
			Lost:   nil,
		},
		Movement: nil,
	}

	suite.mageLevel1 = &levelupbenefit.LevelUpBenefit{
		Identification: levelupbenefit.NewIdentification("mageLevel1", suite.mageClass.ID(), levelupbenefit.Big),
		Defense: levelupbenefit.NewDefense(0,0,0,0,0),
		Offense: levelupbenefit.NewOffense(0,0,0),
		PowerChanges: &levelupbenefit.PowerChanges{
			Gained: nil,
			Lost:   nil,
		},
		Movement: squaddieBuilder.MovementBuilder().Distance(1).Build(),
	}

	suite.dimensionWalkerLevel0 = &levelupbenefit.LevelUpBenefit{
		Identification: levelupbenefit.NewIdentification("dwLevel0", suite.dimensionWalkerClass.ID(), levelupbenefit.Big),
		Movement:       squaddieBuilder.MovementBuilder().Light().Distance(1).Build(),
	}

	suite.dimensionWalkerLevel1 = &levelupbenefit.LevelUpBenefit{
		Identification: levelupbenefit.NewIdentification("dwLevel1", suite.dimensionWalkerClass.ID(), levelupbenefit.Big),
		Defense: levelupbenefit.NewDefense(1,0,0,0,0),
		Offense: levelupbenefit.NewOffense(0,0,0),
	}

	suite.ancientTomeClassLevel0 = &levelupbenefit.LevelUpBenefit{
		Identification: levelupbenefit.NewIdentification("ancientTomeLevel0", suite.ancientTomeClass.ID(), levelupbenefit.Small),
		Defense: levelupbenefit.NewDefense(0,1,1,0,0),
		Offense: levelupbenefit.NewOffense(0,0,1),
	}

	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	suite.levelRepo.AddLevels([]*levelupbenefit.LevelUpBenefit{
		suite.mageLevel0, suite.mageLevel1, suite.dimensionWalkerLevel0, suite.dimensionWalkerLevel1, suite.ancientTomeClassLevel0,
	})

	suite.lotsOfLevels = append(
		(&power.LevelGenerator{
			Instructions: &power.LevelGeneratorInstruction{
				NumberOfLevels: 11,
				ClassID:        suite.atLeastTenLevelsBaseClass.ID(),
				PrefixLevelID:  "lotsLevelsSmall",
				Type:           levelupbenefit.Small,
			},
		}).Build(),
		(&power.LevelGenerator{
			Instructions: &power.LevelGeneratorInstruction{
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

	suite.teros = squaddieBuilder.Builder().Teros().AddClassByReference(suite.mageClass.GetReference()).AddClassByReference(suite.dimensionWalkerClass.GetReference()).AddClassByReference(suite.ancientTomeClass.GetReference()).AddClassByReference(suite.atLeastTenLevelsBaseClass.GetReference()).Build()
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
