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
)

type SquaddieChoosesLevelsSuite struct {
	teros                 *squaddie.Squaddie
	mageClass             *squaddieclass.Class
	onlySmallLevelsClass  *squaddieclass.Class
	classWithInitialLevel *squaddieclass.Class
	lotsOfSmallLevels     []*levelupbenefit.LevelUpBenefit
	lotsOfBigLevels       []*levelupbenefit.LevelUpBenefit
	classRepo             *squaddieclass.Repository
	levelRepo             *levelupbenefit.Repository
	repos                 *repositories.RepositoryCollection
}

var _ = Suite(&SquaddieChoosesLevelsSuite{})

func (suite *SquaddieChoosesLevelsSuite) SetUpTest(checker *C) {
	suite.mageClass = squaddieClassBuilder.ClassBuilder().WithID("class0").WithName("Mage").Build()
	suite.onlySmallLevelsClass = squaddieClassBuilder.ClassBuilder().WithID("class1").WithName("SmallLevels").Build()
	suite.classWithInitialLevel = squaddieClassBuilder.ClassBuilder().WithID("classWithInitialLevel").
		WithName("Class wants big level first").WithInitialBigLevelID("classWithInitialLevelThisIsTakenFirst").
		Build()

	suite.classRepo = squaddieclass.NewRepository()
	suite.classRepo.AddListOfClasses([]*squaddieclass.Class{suite.mageClass, suite.onlySmallLevelsClass, suite.classWithInitialLevel})

	suite.lotsOfSmallLevels = (&power.LevelGenerator{
		Instructions: &power.LevelGeneratorInstruction{
			NumberOfLevels: 10,
			ClassID:        suite.mageClass.ID(),
			PrefixLevelID:  "lotsLevelsSmall",
			Type:           levelupbenefit.Small,
		},
	}).Build()

	suite.lotsOfBigLevels = (&power.LevelGenerator{
		Instructions: &power.LevelGeneratorInstruction{
			NumberOfLevels: 4,
			ClassID:        suite.mageClass.ID(),
			PrefixLevelID:  "lotsLevelsBig",
			Type:           levelupbenefit.Big,
		},
	}).Build()

	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	suite.levelRepo.AddLevels(suite.lotsOfSmallLevels)
	suite.levelRepo.AddLevels(suite.lotsOfBigLevels)
	suite.levelRepo.AddLevels((&power.LevelGenerator{
		Instructions: &power.LevelGeneratorInstruction{
			NumberOfLevels: 2,
			ClassID:        suite.onlySmallLevelsClass.ID(),
			PrefixLevelID:  "smallLevel",
			Type:           levelupbenefit.Small,
		},
	}).Build())

	suite.levelRepo.AddLevels([]*levelupbenefit.LevelUpBenefit{
		{
			Identification: &levelupbenefit.Identification{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            suite.classWithInitialLevel.ID(),
				ID:                 "classWithInitialLevel0",
			},
		},
		{
			Identification: &levelupbenefit.Identification{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            suite.classWithInitialLevel.ID(),
				ID:                 "classWithInitialLevel1",
			},
		},
		{
			Identification: &levelupbenefit.Identification{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            suite.classWithInitialLevel.ID(),
				ID:                 "classWithInitialLevel2",
			},
		},
		{
			Identification: &levelupbenefit.Identification{
				LevelUpBenefitType: levelupbenefit.Big,
				ClassID:            suite.classWithInitialLevel.ID(),
				ID:                 "classWithInitialLevelThisIsTakenFirst",
			},
		},
		{
			Identification: &levelupbenefit.Identification{
				LevelUpBenefitType: levelupbenefit.Big,
				ClassID:            suite.classWithInitialLevel.ID(),
				ID:                 "classWithInitialLevelThisShouldNotBeTakenFirst",
			},
		},
	})

	suite.repos = &repositories.RepositoryCollection{
		LevelRepo: suite.levelRepo,
		ClassRepo: suite.classRepo,
	}

	suite.teros = squaddieBuilder.Builder().Teros().AddClassByReference(suite.mageClass.GetReference()).Build()
}

func (suite *SquaddieChoosesLevelsSuite) TestUseSmallLevelsForClassLevel(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.mageClass.GetReference())
	for index, _ := range [5]int{} {
		levelup.ImproveSquaddie(suite.lotsOfSmallLevels[index], suite.teros, nil)
	}

	levelup.ImproveSquaddie(suite.lotsOfBigLevels[0], suite.teros, nil)
	levelup.ImproveSquaddie(suite.lotsOfBigLevels[1], suite.teros, nil)

	classLevels := levelup.GetSquaddieClassLevels(suite.teros, suite.repos)
	checker.Assert(classLevels[suite.mageClass.ID()], Equals, 5)
}

func (suite *SquaddieChoosesLevelsSuite) TestOddClassLevelEarnsBigAndSmallLevel(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.mageClass.GetReference())
	suite.teros.SetClass(suite.mageClass.ID())
	err := levelup.ImproveSquaddieBasedOnLevel(suite.teros, suite.lotsOfBigLevels[0].Identification.ID, suite.repos)
	checker.Assert(err, IsNil)

	classLevels := levelup.GetSquaddieClassLevels(suite.teros, suite.repos)
	checker.Assert(classLevels[suite.mageClass.ID()], Equals, 1)
	checker.Assert(suite.teros.ClassProgress.ClassLevelsConsumed[suite.mageClass.ID()].LevelsConsumed, HasLen, 2)

	hasSmallLevel := suite.teros.ClassProgress.ClassLevelsConsumed[suite.mageClass.ID()].AnyLevelsConsumed(func(consumedLevelID string) bool {
		return levelupbenefit.AnyLevelUpBenefits(suite.lotsOfSmallLevels, func(level *levelupbenefit.LevelUpBenefit) bool {
			return level.Identification.ID == consumedLevelID
		})
	})
	checker.Assert(hasSmallLevel, Equals, true)

	hasBigLevel := suite.teros.ClassProgress.ClassLevelsConsumed[suite.mageClass.ID()].AnyLevelsConsumed(func(consumedLevelID string) bool {
		return levelupbenefit.AnyLevelUpBenefits(suite.lotsOfBigLevels, func(level *levelupbenefit.LevelUpBenefit) bool {
			return level.Identification.ID == consumedLevelID
		})
	})
	checker.Assert(hasBigLevel, Equals, true)
}

func (suite *SquaddieChoosesLevelsSuite) TestRaisesAnErrorIfClassIsNotFound(checker *C) {
	err := levelup.ImproveSquaddieBasedOnLevel(suite.teros, suite.lotsOfBigLevels[0].Identification.ID, suite.repos)
	checker.Assert(err, ErrorMatches, `class repository: No class found with SquaddieID: ""`)
}

func (suite *SquaddieChoosesLevelsSuite) TestDoesNotChooseBigLevelIfNoneAvailable(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.onlySmallLevelsClass.GetReference())
	suite.teros.SetClass(suite.onlySmallLevelsClass.ID())
	err := levelup.ImproveSquaddieBasedOnLevel(suite.teros, suite.lotsOfBigLevels[0].Identification.ID, suite.repos)
	checker.Assert(err, IsNil)

	classLevels := levelup.GetSquaddieClassLevels(suite.teros, suite.repos)
	checker.Assert(classLevels[suite.onlySmallLevelsClass.ID()], Equals, 1)
	checker.Assert(suite.teros.ClassProgress.ClassLevelsConsumed[suite.onlySmallLevelsClass.ID()].LevelsConsumed, HasLen, 1)
}

func (suite *SquaddieChoosesLevelsSuite) TestChooseSmallLevelAtMostOnce(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.onlySmallLevelsClass.GetReference())
	suite.teros.SetClass(suite.onlySmallLevelsClass.ID())
	levelup.ImproveSquaddieBasedOnLevel(suite.teros, "", suite.repos)

	err := levelup.ImproveSquaddieBasedOnLevel(suite.teros, "", suite.repos)
	checker.Assert(err, IsNil)
	classLevels := levelup.GetSquaddieClassLevels(suite.teros, suite.repos)
	checker.Assert(classLevels[suite.onlySmallLevelsClass.ID()], Equals, 2)
	checker.Assert(suite.teros.ClassProgress.ClassLevelsConsumed[suite.onlySmallLevelsClass.ID()].LevelsConsumed, HasLen, 2)
}

func (suite *SquaddieChoosesLevelsSuite) TestDoesNotChooseSmallLevelIfNoneAvailable(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.onlySmallLevelsClass.GetReference())
	suite.teros.SetClass(suite.onlySmallLevelsClass.ID())
	levelup.ImproveSquaddieBasedOnLevel(suite.teros, "", suite.repos)
	levelup.ImproveSquaddieBasedOnLevel(suite.teros, "", suite.repos)
	err := levelup.ImproveSquaddieBasedOnLevel(suite.teros, "", suite.repos)
	checker.Assert(err, IsNil)

	classLevels := levelup.GetSquaddieClassLevels(suite.teros, suite.repos)
	checker.Assert(classLevels[suite.onlySmallLevelsClass.ID()], Equals, 2)
	checker.Assert(suite.teros.ClassProgress.ClassLevelsConsumed[suite.onlySmallLevelsClass.ID()].LevelsConsumed, HasLen, 2)
}

func (suite *SquaddieChoosesLevelsSuite) TestSquaddieMustChooseInitialLevel(checker *C) {
	suite.teros.ClassProgress.AddClass(suite.classWithInitialLevel.GetReference())
	suite.teros.SetClass(suite.classWithInitialLevel.ID())
	err := levelup.ImproveSquaddieBasedOnLevel(suite.teros, "classWithInitialLevelThisShouldNotBeTakenFirst", suite.repos)
	checker.Assert(err, IsNil)

	classLevels := levelup.GetSquaddieClassLevels(suite.teros, suite.repos)
	checker.Assert(classLevels[suite.classWithInitialLevel.ID()], Equals, 1)
	checker.Assert(suite.teros.ClassProgress.ClassLevelsConsumed[suite.classWithInitialLevel.ID()].LevelsConsumed, HasLen, 2)
	checker.Assert(suite.teros.ClassProgress.ClassLevelsConsumed[suite.classWithInitialLevel.ID()].IsLevelAlreadyConsumed("classWithInitialLevelThisIsTakenFirst"), Equals, true)
	checker.Assert(suite.teros.ClassProgress.ClassLevelsConsumed[suite.classWithInitialLevel.ID()].IsLevelAlreadyConsumed("classWithInitialLevelThisShouldNotBeTakenFirst"), Equals, false)

	levelup.ImproveSquaddieBasedOnLevel(suite.teros, "classWithInitialLevelThisShouldNotBeTakenFirst", suite.repos)
	checker.Assert(suite.teros.ClassProgress.ClassLevelsConsumed[suite.classWithInitialLevel.ID()].LevelsConsumed, HasLen, 3)

	levelup.ImproveSquaddieBasedOnLevel(suite.teros, "classWithInitialLevelThisShouldNotBeTakenFirst", suite.repos)
	checker.Assert(suite.teros.ClassProgress.ClassLevelsConsumed[suite.classWithInitialLevel.ID()].LevelsConsumed, HasLen, 5)
	checker.Assert(suite.teros.ClassProgress.ClassLevelsConsumed[suite.classWithInitialLevel.ID()].IsLevelAlreadyConsumed("classWithInitialLevelThisShouldNotBeTakenFirst"), Equals, true)
}
