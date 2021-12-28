package levelup_test

import (
	"github.com/chadius/terosbattleserver/entity/levelupbenefit"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	"github.com/chadius/terosbattleserver/entity/squaddieinterface"
	"github.com/chadius/terosbattleserver/usecase/levelup"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/utility/testutility/builder"
	. "gopkg.in/check.v1"
)

type SquaddieChoosesLevelsSuite struct {
	teros                   squaddieinterface.Interface
	mageClass               *squaddieclass.Class
	onlySmallLevelsClass    *squaddieclass.Class
	classWithInitialLevel   *squaddieclass.Class
	lotsOfSmallLevels       []*levelupbenefit.LevelUpBenefit
	lotsOfBigLevels         []*levelupbenefit.LevelUpBenefit
	classRepo               *squaddieclass.Repository
	levelRepo               *levelupbenefit.Repository
	repos                   *repositories.RepositoryCollection
	improveSquaddieStrategy levelup.ImproveSquaddieStrategy
	selectLevelUpStrategy   levelup.SelectLevelUpBasedOnSquaddieStrategy
}

var _ = Suite(&SquaddieChoosesLevelsSuite{})

func (suite *SquaddieChoosesLevelsSuite) SetUpTest(checker *C) {
	suite.mageClass = squaddieclass.ClassBuilder().WithID("class0").WithName("Mage").Build()
	suite.onlySmallLevelsClass = squaddieclass.ClassBuilder().WithID("class1").WithName("SmallLevels").Build()
	suite.classWithInitialLevel = squaddieclass.ClassBuilder().WithID("classWithInitialLevel").
		WithName("Class wants big level first").WithInitialBigLevelID("classWithInitialLevelThisIsTakenFirst").
		Build()

	suite.classRepo = squaddieclass.NewRepository()
	suite.classRepo.AddListOfClasses([]*squaddieclass.Class{suite.mageClass, suite.onlySmallLevelsClass, suite.classWithInitialLevel})

	suite.lotsOfSmallLevels = (&builder.LevelGenerator{
		Instructions: &builder.LevelGeneratorInstruction{
			NumberOfLevels: 10,
			ClassID:        suite.mageClass.ID(),
			PrefixLevelID:  "lotsLevelsSmall",
			Type:           levelupbenefit.Small,
		},
	}).Build()

	suite.lotsOfBigLevels = (&builder.LevelGenerator{
		Instructions: &builder.LevelGeneratorInstruction{
			NumberOfLevels: 4,
			ClassID:        suite.mageClass.ID(),
			PrefixLevelID:  "lotsLevelsBig",
			Type:           levelupbenefit.Big,
		},
	}).Build()

	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	suite.levelRepo.AddLevels(suite.lotsOfSmallLevels)
	suite.levelRepo.AddLevels(suite.lotsOfBigLevels)
	suite.levelRepo.AddLevels((&builder.LevelGenerator{
		Instructions: &builder.LevelGeneratorInstruction{
			NumberOfLevels: 2,
			ClassID:        suite.onlySmallLevelsClass.ID(),
			PrefixLevelID:  "smallLevel",
			Type:           levelupbenefit.Small,
		},
	}).Build())

	classWithInitialLevel0, _ := levelupbenefit.NewLevelUpBenefitBuilder().LevelID("classWithInitialLevel0").ClassID(suite.classWithInitialLevel.ID()).Build()
	classWithInitialLevel1, _ := levelupbenefit.NewLevelUpBenefitBuilder().LevelID("classWithInitialLevel1").ClassID(suite.classWithInitialLevel.ID()).Build()
	classWithInitialLevel2, _ := levelupbenefit.NewLevelUpBenefitBuilder().LevelID("classWithInitialLevel2").ClassID(suite.classWithInitialLevel.ID()).Build()
	classWithInitialLevelThisIsTakenFirst, _ := levelupbenefit.NewLevelUpBenefitBuilder().LevelID("classWithInitialLevelThisIsTakenFirst").ClassID(suite.classWithInitialLevel.ID()).BigLevel().Build()
	classWithInitialLevelThisShouldNotBeTakenFirst, _ := levelupbenefit.NewLevelUpBenefitBuilder().LevelID("classWithInitialLevelThisShouldNotBeTakenFirst").ClassID(suite.classWithInitialLevel.ID()).BigLevel().Build()

	suite.levelRepo.AddLevels([]*levelupbenefit.LevelUpBenefit{
		classWithInitialLevel0,
		classWithInitialLevel1,
		classWithInitialLevel2,
		classWithInitialLevelThisIsTakenFirst,
		classWithInitialLevelThisShouldNotBeTakenFirst,
	})

	suite.repos = &repositories.RepositoryCollection{
		LevelRepo: suite.levelRepo,
		ClassRepo: suite.classRepo,
	}

	suite.teros = squaddie.NewSquaddieBuilder().Teros().AddClassByReference(suite.mageClass.GetReference()).Build()
	suite.improveSquaddieStrategy = &levelup.ImproveSquaddieClass{}
	suite.selectLevelUpStrategy = &levelup.SelectLevelUpBasedOnSquaddieBigLevelsOnEvenLevels{}
}

func (suite *SquaddieChoosesLevelsSuite) TestUseSmallLevelsForClassLevel(checker *C) {
	suite.teros.AddClass(suite.mageClass.GetReference())
	for index, _ := range [5]int{} {
		suite.improveSquaddieStrategy.ImproveSquaddie(suite.lotsOfSmallLevels[index], suite.teros)
	}

	suite.improveSquaddieStrategy.ImproveSquaddie(suite.lotsOfBigLevels[0], suite.teros)
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.lotsOfBigLevels[1], suite.teros)

	classLevels := suite.selectLevelUpStrategy.GetSquaddieClassLevels(suite.teros, suite.repos)
	checker.Assert(classLevels[suite.mageClass.ID()], Equals, 5)
}

func (suite *SquaddieChoosesLevelsSuite) TestOddClassLevelEarnsBigAndSmallLevel(checker *C) {
	suite.teros.AddClass(suite.mageClass.GetReference())
	suite.teros.SetClass(suite.mageClass.ID())
	err := suite.selectLevelUpStrategy.ImproveSquaddieBasedOnLevel(suite.teros, suite.lotsOfBigLevels[0].ID(), suite.repos)
	checker.Assert(err, IsNil)

	classLevels := suite.selectLevelUpStrategy.GetSquaddieClassLevels(suite.teros, suite.repos)
	checker.Assert(classLevels[suite.mageClass.ID()], Equals, 1)
	checker.Assert((*suite.teros.ClassLevelsConsumed())[suite.mageClass.ID()].GetLevelsConsumed(), HasLen, 2)

	hasSmallLevel := (*suite.teros.ClassLevelsConsumed())[suite.mageClass.ID()].AnyLevelsConsumed(func(consumedLevelID string) bool {
		return levelupbenefit.AnyLevelUpBenefits(suite.lotsOfSmallLevels, func(level *levelupbenefit.LevelUpBenefit) bool {
			return level.ID() == consumedLevelID
		})
	})
	checker.Assert(hasSmallLevel, Equals, true)

	hasBigLevel := (*suite.teros.ClassLevelsConsumed())[suite.mageClass.ID()].AnyLevelsConsumed(func(consumedLevelID string) bool {
		return levelupbenefit.AnyLevelUpBenefits(suite.lotsOfBigLevels, func(level *levelupbenefit.LevelUpBenefit) bool {
			return level.ID() == consumedLevelID
		})
	})
	checker.Assert(hasBigLevel, Equals, true)
}

func (suite *SquaddieChoosesLevelsSuite) TestRaisesAnErrorIfClassIsNotFound(checker *C) {
	err := suite.selectLevelUpStrategy.ImproveSquaddieBasedOnLevel(suite.teros, suite.lotsOfBigLevels[0].ID(), suite.repos)
	checker.Assert(err, ErrorMatches, `class repository: No class found with id: ""`)
}

func (suite *SquaddieChoosesLevelsSuite) TestDoesNotChooseBigLevelIfNoneAvailable(checker *C) {
	suite.teros.AddClass(suite.onlySmallLevelsClass.GetReference())
	suite.teros.SetClass(suite.onlySmallLevelsClass.ID())
	err := suite.selectLevelUpStrategy.ImproveSquaddieBasedOnLevel(suite.teros, suite.lotsOfBigLevels[0].ID(), suite.repos)
	checker.Assert(err, IsNil)

	classLevels := suite.selectLevelUpStrategy.GetSquaddieClassLevels(suite.teros, suite.repos)
	checker.Assert(classLevels[suite.onlySmallLevelsClass.ID()], Equals, 1)
	checker.Assert((*suite.teros.ClassLevelsConsumed())[suite.onlySmallLevelsClass.ID()].GetLevelsConsumed(), HasLen, 1)
}

func (suite *SquaddieChoosesLevelsSuite) TestChooseSmallLevelAtMostOnce(checker *C) {
	suite.teros.AddClass(suite.onlySmallLevelsClass.GetReference())
	suite.teros.SetClass(suite.onlySmallLevelsClass.ID())
	suite.selectLevelUpStrategy.ImproveSquaddieBasedOnLevel(suite.teros, "", suite.repos)

	err := suite.selectLevelUpStrategy.ImproveSquaddieBasedOnLevel(suite.teros, "", suite.repos)
	checker.Assert(err, IsNil)
	classLevels := suite.selectLevelUpStrategy.GetSquaddieClassLevels(suite.teros, suite.repos)
	checker.Assert(classLevels[suite.onlySmallLevelsClass.ID()], Equals, 2)
	checker.Assert((*suite.teros.ClassLevelsConsumed())[suite.onlySmallLevelsClass.ID()].GetLevelsConsumed(), HasLen, 2)
}

func (suite *SquaddieChoosesLevelsSuite) TestDoesNotChooseSmallLevelIfNoneAvailable(checker *C) {
	suite.teros.AddClass(suite.onlySmallLevelsClass.GetReference())
	suite.teros.SetClass(suite.onlySmallLevelsClass.ID())
	suite.selectLevelUpStrategy.ImproveSquaddieBasedOnLevel(suite.teros, "", suite.repos)
	suite.selectLevelUpStrategy.ImproveSquaddieBasedOnLevel(suite.teros, "", suite.repos)
	err := suite.selectLevelUpStrategy.ImproveSquaddieBasedOnLevel(suite.teros, "", suite.repos)
	checker.Assert(err, IsNil)

	classLevels := suite.selectLevelUpStrategy.GetSquaddieClassLevels(suite.teros, suite.repos)
	checker.Assert(classLevels[suite.onlySmallLevelsClass.ID()], Equals, 2)
	checker.Assert((*suite.teros.ClassLevelsConsumed())[suite.onlySmallLevelsClass.ID()].GetLevelsConsumed(), HasLen, 2)
}

func (suite *SquaddieChoosesLevelsSuite) TestSquaddieMustChooseInitialLevel(checker *C) {
	suite.teros.AddClass(suite.classWithInitialLevel.GetReference())
	suite.teros.SetClass(suite.classWithInitialLevel.ID())
	err := suite.selectLevelUpStrategy.ImproveSquaddieBasedOnLevel(suite.teros, "classWithInitialLevelThisShouldNotBeTakenFirst", suite.repos)
	checker.Assert(err, IsNil)

	classLevels := suite.selectLevelUpStrategy.GetSquaddieClassLevels(suite.teros, suite.repos)
	checker.Assert(classLevels[suite.classWithInitialLevel.ID()], Equals, 1)
	checker.Assert((*suite.teros.ClassLevelsConsumed())[suite.classWithInitialLevel.ID()].GetLevelsConsumed(), HasLen, 2)
	checker.Assert(suite.teros.IsClassLevelAlreadyUsed("classWithInitialLevelThisIsTakenFirst"), Equals, true)
	checker.Assert(suite.teros.IsClassLevelAlreadyUsed("classWithInitialLevelThisShouldNotBeTakenFirst"), Equals, false)

	suite.selectLevelUpStrategy.ImproveSquaddieBasedOnLevel(suite.teros, "classWithInitialLevelThisShouldNotBeTakenFirst", suite.repos)
	checker.Assert((*suite.teros.ClassLevelsConsumed())[suite.classWithInitialLevel.ID()].GetLevelsConsumed(), HasLen, 3)

	suite.selectLevelUpStrategy.ImproveSquaddieBasedOnLevel(suite.teros, "classWithInitialLevelThisShouldNotBeTakenFirst", suite.repos)
	checker.Assert((*suite.teros.ClassLevelsConsumed())[suite.classWithInitialLevel.ID()].GetLevelsConsumed(), HasLen, 5)
	checker.Assert(suite.teros.IsClassLevelAlreadyUsed("classWithInitialLevelThisShouldNotBeTakenFirst"), Equals, true)
}
