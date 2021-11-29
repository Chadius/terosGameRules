package levelupbenefit_test

import (
	"github.com/chadius/terosbattleserver/entity/levelupbenefit"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	"github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddieclass"
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

	suite.mageClass = squaddieBuilder.ClassBuilder().WithID("class1").Build()
	suite.lotsOfSmallLevels = (&power.LevelGenerator{
		Instructions: &power.LevelGeneratorInstruction{
			NumberOfLevels: 11,
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
}

func (suite *LevelUpBenefitRepositorySuite) TestCreateLevelUpBenefitsFromJSON(checker *C) {
	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	suite.jsonByteStream = []byte(`[
          {
           "identification": {
            "id":"abcdefg0",
            "level_up_benefit_type": "small",
            "class_id": "class0"
          },
           "defense": {
            "dodge": 4,
            "deflect": 5,
            "max_barrier": 6,
            "max_hit_points": 1,
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
	checker.Assert(suite.levelRepo.GetNumberOfLevelUpBenefits(), Equals, 0)
	success, err := suite.levelRepo.AddJSONSource(suite.jsonByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.levelRepo.GetNumberOfLevelUpBenefits(), Equals, 1)
}

func (suite *LevelUpBenefitRepositorySuite) TestCreateLevelUpBenefitsFromYAML(checker *C) {
	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	suite.yamlByteStream = []byte(
		`
- identification:
    id: abcdefg0
    class_id: class0
    level_up_benefit_type: small
  defense:
    max_hit_points: 1
    dodge: 4
    deflect: 5
    max_barrier: 6
    armor: 7
  offense:
    aim: 0
    strength: 2
    mind: 3
  powers:
    gained:
    - name: Scimitar
      id: deadbeef
  movement:
    distance: 1,
    type: teleport
    hit_and_run": true
`)
	checker.Assert(suite.levelRepo.GetNumberOfLevelUpBenefits(), Equals, 0)
	success, _ := suite.levelRepo.AddYAMLSource(suite.yamlByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.levelRepo.GetNumberOfLevelUpBenefits(), Equals, 1)
}

func (suite *LevelUpBenefitRepositorySuite) TestCreateLevelUpBenefitsFromASlice(checker *C) {
	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	checker.Assert(suite.levelRepo.GetNumberOfLevelUpBenefits(), Equals, 0)
	success, _ := suite.levelRepo.AddLevels([]*levelupbenefit.LevelUpBenefit{
		{
			Identification: levelupbenefit.NewIdentification("level0", "class0", levelupbenefit.Small),
		},
		{
			Identification: levelupbenefit.NewIdentification("level1", "class0", levelupbenefit.Small),
		},
	})
	checker.Assert(success, Equals, true)
	checker.Assert(suite.levelRepo.GetNumberOfLevelUpBenefits(), Equals, 2)
}

func (suite *LevelUpBenefitRepositorySuite) TestStopLoadingOnFirstInvalidLevelUpBenefit(checker *C) {
	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	byteStream := []byte(`[
          {
           "identification": {
            "id":"abcdefg0",
            "class_id": "class0",
            "level_up_benefit_type": "small"
          },
           "defense": {
            "max_hit_points": 1,
            "dodge": 4,
            "deflect": 5,
            "max_barrier": 6,
            "armor": 7
           },
            "aim": 0,
            "strength": 2,
            "mind": 3,
            "powers": {
              "gained": [
                {
                  "name": "Scimitar",
                  "id": "deadbeef"
                }
              ]
            }
          },
		  {
           "identification": {
				"level_up_benefit_type": "unknown",
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
             "gained": [{"name": "Scimitar", "id": "deadbeef"}]
           }
        }
]`)
	success, err := suite.levelRepo.AddJSONSource(byteStream)
	checker.Assert(success, Equals, false)
	checker.Assert(err.Error(), Equals, `unknown level up benefit type`)
}

func (suite *LevelUpBenefitRepositorySuite) TestCanSearchLevelUpBenefits(checker *C) {
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
             "gained":[{
               "name": "Scimitar",
               "id": "deadbeef"
             }]
           },
           "movement": {
             "distance": 1,
             "type": "teleport",
             "hit_and_run": true
           }
     }
]`)
	suite.levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
	success, _ := suite.levelRepo.AddJSONSource(suite.jsonByteStream)
	checker.Assert(success, Equals, true)

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

	checker.Assert(firstBenefit.PowerChanges.Gained(), HasLen, 1)
	checker.Assert(firstBenefit.PowerChanges.Gained()[0].Name, Equals, "Scimitar")
	checker.Assert(firstBenefit.PowerChanges.Gained()[0].PowerID, Equals, "deadbeef")
}

func (suite *LevelUpBenefitRepositorySuite) TestRaisesAnErrorWithNonexistentClassID(checker *C) {
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
	suite.levelRepo.AddJSONSource(suite.jsonByteStream)

	benefits, err := suite.levelRepo.GetLevelUpBenefitsByClassID("Class not found")
	checker.Assert(err, ErrorMatches, `no LevelUpBenefits for this class SquaddieID: "Class not found"`)
	checker.Assert(benefits, HasLen, 0)
}

func (suite *LevelUpBenefitRepositorySuite) TestGetBigAndSmallLevelsForAGivenClass(checker *C) {
	levelsByBenefitType, err := suite.levelRepo.GetLevelUpBenefitsForClassByType(suite.mageClass.ID())
	checker.Assert(err, IsNil)
	checker.Assert(levelsByBenefitType[levelupbenefit.Small], HasLen, 11)
	checker.Assert(levelsByBenefitType[levelupbenefit.Big], HasLen, 4)
}

func (suite *LevelUpBenefitRepositorySuite) TestRaiseErrorIfClassDoesNotExist(checker *C) {
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
            "aim": 0,
            "strength": 2,
            "mind": 3,
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
	suite.levelRepo.AddJSONSource(suite.jsonByteStream)
	levelsByBenefitType, err := suite.levelRepo.GetLevelUpBenefitsForClassByType("bad SquaddieID")
	checker.Assert(err, ErrorMatches, `no LevelUpBenefits for this class SquaddieID: "bad SquaddieID"`)
	checker.Assert(levelsByBenefitType[levelupbenefit.Small], HasLen, 0)
	checker.Assert(levelsByBenefitType[levelupbenefit.Big], HasLen, 0)
}
