package levelupbenefit_test

import (
	"github.com/chadius/terosgamerules/entity/levelupbenefit"
	"github.com/chadius/terosgamerules/entity/squaddieclass"
	"github.com/chadius/terosgamerules/utility/testutility/builder"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type LevelUpBenefitRepositorySuite struct {
	levelRepo         *levelupbenefit.Repository
	jsonByteStream    []byte
	yamlByteStream    []byte
	mageClass         *squaddieclass.Class
	lotsOfSmallLevels []*levelupbenefit.LevelUpBenefit
	lotsOfBigLevels   []*levelupbenefit.LevelUpBenefit
}

var _ = Suite(&LevelUpBenefitRepositorySuite{})

func (suite *LevelUpBenefitRepositorySuite) SetUpTest(c *C) {
	suite.jsonByteStream = []byte(`[
          {
           "identification": {
            "id":"abcdefg0",
            "level_up_benefit_type": "small",
            "class_id": "class0"
          },
          "defense": {
            "max_hit_points": 1,
            "dodge": 4,
            "deflect": 5,
            "max_barrier": 6,
            "armor": 7
          },
          "offense": {
            "aim": 0,
            "strength": 2,
            "mind": 3
          },
          "powers": {
            "gained": [
              {
                "name": "Scimitar",
                "id": "deadbeef"
              }
            ]
          },
          "movement": {
            "distance": 1,
            "type": "teleport",
            "hit_and_run": true
          }
      }
]`)

	suite.mageClass = squaddieclass.ClassBuilder().WithID("class1").Build()
	suite.lotsOfSmallLevels = (&builder.LevelGenerator{
		Instructions: &builder.LevelGeneratorInstruction{
			NumberOfLevels: 11,
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
}

func (suite *LevelUpBenefitRepositorySuite) TestCreateLevelUpBenefitsFromASlice(checker *C) {
	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	checker.Assert(suite.levelRepo.GetNumberOfLevelUpBenefits(), Equals, 0)

	level0, _ := levelupbenefit.NewLevelUpBenefitBuilder().LevelID("level0").ClassID("class0").Build()
	level1, _ := levelupbenefit.NewLevelUpBenefitBuilder().LevelID("level1").ClassID("class0").Build()

	success, _ := suite.levelRepo.AddLevels([]*levelupbenefit.LevelUpBenefit{
		level0,
		level1,
	})
	checker.Assert(success, Equals, true)
	checker.Assert(suite.levelRepo.GetNumberOfLevelUpBenefits(), Equals, 2)
}

func (suite *LevelUpBenefitRepositorySuite) TestCanSearchLevelUpBenefits(checker *C) {
	suite.jsonByteStream = []byte(
		`[
	{
        "id":"abcdefg0",
        "class_id": "class0",
        "hit_points": 1,
        "dodge": 4,
        "deflect": 5,
        "barrier": 6,
        "armor": 7,
        "aim": 0,
        "strength": 2,
        "mind": 3,
        "powers_gained": [{
           "name": "Scimitar",
           "id": "deadbeef"
        }],
        "movement_distance": 1,
        "movement_type": "teleport",
        "can_hit_and_run": true
	}
]`)

	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	loadErr := suite.levelRepo.AddJSON(suite.jsonByteStream)
	checker.Assert(loadErr, IsNil)

	benefits, err := suite.levelRepo.GetLevelUpBenefitsByClassID("class0")
	checker.Assert(err, IsNil)
	checker.Assert(benefits, HasLen, 1)

	firstBenefit := benefits[0]
	checker.Assert(firstBenefit.LevelUpBenefitType(), Equals, levelupbenefit.Small)
	checker.Assert(firstBenefit.ClassID(), Equals, "class0")

	checker.Assert(firstBenefit.MaxHitPoints(), Equals, 1)
	checker.Assert(firstBenefit.Dodge(), Equals, 4)
	checker.Assert(firstBenefit.Deflect(), Equals, 5)
	checker.Assert(firstBenefit.MaxBarrier(), Equals, 6)
	checker.Assert(firstBenefit.Armor(), Equals, 7)

	checker.Assert(firstBenefit.Aim(), Equals, 0)
	checker.Assert(firstBenefit.Strength(), Equals, 2)
	checker.Assert(firstBenefit.Mind(), Equals, 3)

	checker.Assert(firstBenefit.PowersGained(), HasLen, 1)
	checker.Assert(firstBenefit.PowersGained()[0].Name, Equals, "Scimitar")
	checker.Assert(firstBenefit.PowersGained()[0].PowerID, Equals, "deadbeef")
}

func (suite *LevelUpBenefitRepositorySuite) TestRaisesAnErrorWithNonexistentClassID(checker *C) {
	suite.jsonByteStream = []byte(
		`[
	{
        "class_id": "class0",
	}
]`)

	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	suite.levelRepo.AddJSON(suite.jsonByteStream)

	benefits, err := suite.levelRepo.GetLevelUpBenefitsByClassID("Class not found")

	checker.Assert(err, ErrorMatches, `no LevelUpBenefits for this class id: "Class not found"`)
	checker.Assert(benefits, HasLen, 0)
}

func (suite *LevelUpBenefitRepositorySuite) TestGetBigAndSmallLevelsForAGivenClass(checker *C) {
	levelsByBenefitType, err := suite.levelRepo.GetLevelUpBenefitsForClassByType(suite.mageClass.ID())
	checker.Assert(err, IsNil)
	checker.Assert(levelsByBenefitType[levelupbenefit.Small], HasLen, 11)
	checker.Assert(levelsByBenefitType[levelupbenefit.Big], HasLen, 4)
}

func (suite *LevelUpBenefitRepositorySuite) TestRaiseErrorIfClassDoesNotExist(checker *C) {
	suite.jsonByteStream = []byte(
		`[
	{
        "class_id": "class0",
	}
]`)
	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	suite.levelRepo.AddJSON(suite.jsonByteStream)

	levelsByBenefitType, err := suite.levelRepo.GetLevelUpBenefitsForClassByType("bad classID")

	checker.Assert(err, ErrorMatches, `no LevelUpBenefits for this class id: "bad classID"`)
	checker.Assert(levelsByBenefitType[levelupbenefit.Small], HasLen, 0)
	checker.Assert(levelsByBenefitType[levelupbenefit.Big], HasLen, 0)
}

type BuilderFormatSuite struct{}

var _ = Suite(&BuilderFormatSuite{})

func (suite *BuilderFormatSuite) TestCreateWithYAML(checker *C) {
	yamlByteStream := []byte(
		`-
  id: abcdefg0
  class_id: class0
  hit_points: 2
-
  id: abcdefg1
  class_id: class0
  hit_points: 3
`)
	levelRepo := levelupbenefit.NewLevelUpBenefitRepository()
	err := levelRepo.AddYAML(yamlByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(levelRepo.GetNumberOfLevelUpBenefits(), Equals, 2)
}

func (suite *BuilderFormatSuite) TestCreateWithJSON(checker *C) {
	jsonByteStream := []byte(
		`[
	{
	  "id": "abcdefg0",
	  "class_id": "class0",
	  "hit_points": 2
	},
	{
	  "id": "abcdefg1",
	  "class_id": "class0",
	  "hit_points": 3
	}
]`)
	levelRepo := levelupbenefit.NewLevelUpBenefitRepository()
	err := levelRepo.AddJSON(jsonByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(levelRepo.GetNumberOfLevelUpBenefits(), Equals, 2)
}
