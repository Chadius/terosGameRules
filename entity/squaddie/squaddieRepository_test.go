package squaddie_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	powerBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	. "gopkg.in/check.v1"
)

type SquaddieRepositorySuite struct {
	squaddieRepository *squaddie.Repository
	teros              *squaddie.Squaddie
	attackA         *power.Power
}

var _ = Suite(&SquaddieRepositorySuite{})

func (suite *SquaddieRepositorySuite) SetUpTest(checker *C) {
	suite.squaddieRepository = squaddie.NewSquaddieRepository()
	suite.teros = squaddieBuilder.Builder().Teros().Armor(2).Dodge(3).Deflect(4).Barrier(1).Build()
	suite.attackA = powerBuilder.Builder().WithName("Attack Formation A").Build()
}

func (suite *SquaddieRepositorySuite) TestUseJSONSource(checker *C) {
	checker.Assert(suite.squaddieRepository.GetNumberOfSquaddies(), Equals, 0)
	jsonByteStream := []byte(`[{
				"identification": {
					"name": "teros",
					"affiliation": "Player"
				}
			}]`)
	success, _ := suite.squaddieRepository.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.squaddieRepository.GetNumberOfSquaddies(), Equals, 1)
}

func (suite *SquaddieRepositorySuite) TestCloneSquaddies(checker *C) {
	jsonByteStream := []byte(`[{
				"identification": {
					"id": "squaddieID",
					"name": "teros",
					"affiliation": "Player"
				},
				"offense": {
					"aim": 5
				}
			}]`)
	success, _ := suite.squaddieRepository.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)

	suite.teros = suite.squaddieRepository.GetSquaddieByID("squaddieID")
	checker.Assert(suite.teros, NotNil)
	checker.Assert(suite.teros.Name(), Equals, "teros")
	checker.Assert(suite.teros.Aim(), Equals, 5)

	missingno := suite.squaddieRepository.GetSquaddieByID("Does not exist")
	checker.Assert(missingno, IsNil)
}

func (suite *SquaddieRepositorySuite) TestGetExistingSquaddieUsingID(checker *C) {
	jsonByteStream := []byte(`[{
				"identification": {
					"id": "12345",
					"name": "teros",
					"affiliation": "Player"
				},
				"offense": {
					"aim": 5
				}
			}]`)
	success, _ := suite.squaddieRepository.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)

	suite.teros = suite.squaddieRepository.GetSquaddieByID("12345")
	checker.Assert(suite.teros, NotNil)
	checker.Assert(suite.teros.Name(), Equals, "teros")
	checker.Assert(suite.teros.Aim(), Equals, 5)

	missingno := suite.squaddieRepository.GetSquaddieByID("Does not exist")
	checker.Assert(missingno, IsNil)
}

func (suite *SquaddieRepositorySuite) TestClonedSquaddiesHaveDifferentID(checker *C) {
	jsonByteStream := []byte(`[{
				"identification": {
					"id": "terosID",
					"name": "teros",
					"affiliation": "Player"
				},
				"offense": {
					"aim": 5
				}
			}]`)
	success, _ := suite.squaddieRepository.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)

	teros0 := suite.squaddieRepository.GetSquaddieByID("terosID")
	teros1 := suite.squaddieRepository.GetSquaddieByID("terosID")
	checker.Assert(teros0, Not(Equals), teros1)
}

func (suite *SquaddieRepositorySuite) TestLoadClassLevels(checker *C) {
	jsonByteStream := []byte(`[{
				"identification": {
					"id": "terosSquaddieID",
					"name": "teros",
					"affiliation": "Player"
				},
				"offense": {
					"aim": 5
				},
				"class_progress": {
					"class_levels": {
					  "1": {
						"id": "1",
						"name": "Mage",
						"levels_gained": ["123"]
					  },
					  "2": {
						"id": "2",
						"name": "Dimension Walker"
					  }
				   }
				}
			}]`)
	success, _ := suite.squaddieRepository.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)

	suite.teros = suite.squaddieRepository.GetSquaddieByID("terosSquaddieID")
	checker.Assert(suite.teros.ClassProgress.GetLevelCountsByClass(), DeepEquals, map[string]int{"1": 1, "2": 0})
}

func (suite *SquaddieRepositorySuite) TestStopLoadingSquaddiesWhenInvalid(checker *C) {
	jsonByteStream := []byte(`[{
				"identification": {
					"Name": "teros",
					"Affiliation": "Player"
				}
			},{
				"identification": {
					"Name": "teros2",
					"Affiliation": "Unknown Affiliation"
				}
			}]`)
	success, err := suite.squaddieRepository.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, false)
	checker.Assert(err, ErrorMatches, "squaddie  has unknown affiliation: 'Unknown Affiliation'")
}

func (suite *SquaddieRepositorySuite) TestCreateSquaddiesWithMovement(checker *C) {
	jsonByteStream := []byte(`[
			{
				"identification": {
					"id": "Soldier",
					"name": "Soldier",
					"affiliation": "Player"
				},
				"movement": { "distance": 5, "type": "foot"}
			},
			{
				"identification": {
					"id": "Scout",
					"name": "Scout",
					"affiliation": "Player"
				},
				"movement": { "distance": 4, "type": "light"}
			},
			{
				"identification": {
					"id": "Bird",
					"name": "Bird",
					"affiliation": "Player"
				},
				"movement": { "distance": 3, "type": "fly", "hit_and_run": true}
			},
			{
				"identification": {
					"id": "Teleporter",
					"name": "Teleporter",
					"affiliation": "Player"
				},
				"movement": { "distance": 2, "type": "teleport"}
			}
			]`)

	success, _ := suite.squaddieRepository.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.squaddieRepository.GetNumberOfSquaddies(), Equals, 4)

	soldier := suite.squaddieRepository.GetSquaddieByID("Soldier")
	checker.Assert(soldier.Name(), Equals, "Soldier")
	checker.Assert(soldier.Movement.MovementDistance(), Equals, 5)
	checker.Assert(soldier.Movement.MovementType(), Equals, squaddie.MovementType(squaddie.Foot))
	checker.Assert(soldier.Movement.CanHitAndRun(), Equals, false)

	scout := suite.squaddieRepository.GetSquaddieByID("Scout")
	checker.Assert(scout.Name(), Equals, "Scout")
	checker.Assert(scout.Movement.MovementDistance(), Equals, 4)
	checker.Assert(scout.Movement.MovementType(), Equals, squaddie.MovementType(squaddie.Light))
	checker.Assert(scout.Movement.CanHitAndRun(), Equals, false)

	bird := suite.squaddieRepository.GetSquaddieByID("Bird")
	checker.Assert(bird.Name(), Equals, "Bird")
	checker.Assert(bird.Movement.MovementDistance(), Equals, 3)
	checker.Assert(bird.Movement.MovementType(), Equals, squaddie.MovementType(squaddie.Fly))
	checker.Assert(bird.Movement.CanHitAndRun(), Equals, true)

	teleporter := suite.squaddieRepository.GetSquaddieByID("Teleporter")
	checker.Assert(teleporter.Name(), Equals, "Teleporter")
	checker.Assert(teleporter.Movement.MovementDistance(), Equals, 2)
	checker.Assert(teleporter.Movement.MovementType(), Equals, squaddie.MovementType(squaddie.Teleport))
	checker.Assert(teleporter.Movement.CanHitAndRun(), Equals, false)
}

func (suite *SquaddieRepositorySuite) TestCanGetExistingSquaddies(checker *C) {
	originalSquaddie := squaddieBuilder.Builder().WithName("Original").AsAlly().Build()
	suite.squaddieRepository.AddSquaddies([]*squaddie.Squaddie{originalSquaddie})
	referencedSquaddie := suite.squaddieRepository.GetOriginalSquaddieByID(originalSquaddie.ID())
	checker.Assert(referencedSquaddie, Equals, originalSquaddie)
}

func (suite *SquaddieRepositorySuite) TestLoadSquaddieByYAML(checker *C) {
	checker.Assert(suite.squaddieRepository.GetNumberOfSquaddies(), Equals, 0)
	yamlByteStream := []byte(`-
  identification:
    name: teros
    affiliation: Player
`)
	suite.squaddieRepository.AddYAMLSource(yamlByteStream)
	checker.Assert(suite.squaddieRepository.GetNumberOfSquaddies(), Equals, 1)
}

func (suite *SquaddieRepositorySuite) TestLoadClassLevelsYAML(checker *C) {
	yamlByteStream := []byte(`-
  identification:
    name: teros
    id: terosSquaddieID
    affiliation: Player
  class_progress:
    class_levels:
      '1':
        id: '1'
        name: Mage
        levels_gained:
        - '123'
      '2':
        id: '2'
        name: Dimension Walker
`)
	success, _ := suite.squaddieRepository.AddYAMLSource(yamlByteStream)
	checker.Assert(success, Equals, true)

	suite.teros = suite.squaddieRepository.GetSquaddieByID("terosSquaddieID")
	checker.Assert(suite.teros.ClassProgress.GetLevelCountsByClass(), DeepEquals, map[string]int{"1": 1, "2": 0})
}

func (suite *SquaddieRepositorySuite) TestStopLoadingSquaddiesUponFirstInvalid(checker *C) {
	yamlByteStream := []byte(`-
  identification:
    name: teros
    affiliation: Player
-
  identification:
    name: teros2
    affiliation: Unknown SquaddieAffiliation`)
	success, err := suite.squaddieRepository.AddYAMLSource(yamlByteStream)
	checker.Assert(success, Equals, false)
	checker.Assert(err, NotNil)
	checker.Assert(err.Error(), Equals, "squaddie  has unknown affiliation: 'Unknown SquaddieAffiliation'")
}

func (suite *SquaddieRepositorySuite) TestLoadSquaddiesWithDifferentMovementYAML(checker *C) {
	yamlByteStream := []byte(`-
  identification:
    id: Soldier
    name: Soldier
    affiliation: Player
  movement:
    distance: 5
    type: foot
-
  identification:
    id: Scout
    name: Scout
    affiliation: Player
  movement:
    distance: 4
    type: light
-
  identification:
    id: Bird
    name: Bird
    affiliation: Player
  movement:
    distance: 3
    type: fly
    hit_and_run: true
-
  identification:
    id: Teleporter
    name: Teleporter
    affiliation: Player
  movement:
    distance: 2
    type: teleport
    hit_and_run: false
`)

	success, _ := suite.squaddieRepository.AddYAMLSource(yamlByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.squaddieRepository.GetNumberOfSquaddies(), Equals, 4)

	soldier := suite.squaddieRepository.GetSquaddieByID("Soldier")
	checker.Assert(soldier.Name(), Equals, "Soldier")
	checker.Assert(soldier.Movement.MovementDistance(), Equals, 5)
	checker.Assert(soldier.Movement.MovementType(), Equals, squaddie.MovementType(squaddie.Foot))
	checker.Assert(soldier.Movement.CanHitAndRun(), Equals, false)

	scout := suite.squaddieRepository.GetSquaddieByID("Scout")
	checker.Assert(scout.Name(), Equals, "Scout")
	checker.Assert(scout.Movement.MovementDistance(), Equals, 4)
	checker.Assert(scout.Movement.MovementType(), Equals, squaddie.MovementType(squaddie.Light))
	checker.Assert(scout.Movement.CanHitAndRun(), Equals, false)

	bird := suite.squaddieRepository.GetSquaddieByID("Bird")
	checker.Assert(bird.Name(), Equals, "Bird")
	checker.Assert(bird.Movement.MovementDistance(), Equals, 3)
	checker.Assert(bird.Movement.MovementType(), Equals, squaddie.MovementType(squaddie.Fly))
	checker.Assert(bird.Movement.CanHitAndRun(), Equals, true)

	teleporter := suite.squaddieRepository.GetSquaddieByID("Teleporter")
	checker.Assert(teleporter.Name(), Equals, "Teleporter")
	checker.Assert(teleporter.Movement.MovementDistance(), Equals, 2)
	checker.Assert(teleporter.Movement.MovementType(), Equals, squaddie.MovementType(squaddie.Teleport))
	checker.Assert(teleporter.Movement.CanHitAndRun(), Equals, false)
}

func (suite *SquaddieRepositorySuite) TestAddSquaddieDirectly(checker *C) {
	success, err := suite.squaddieRepository.AddSquaddies([]*squaddie.Squaddie{
		squaddieBuilder.Builder().Build(),
	})
	checker.Assert(success, Equals, true)
	checker.Assert(err, IsNil)
	checker.Assert(suite.squaddieRepository.GetNumberOfSquaddies(), Equals, 1)
}

func (suite *SquaddieRepositorySuite) TestCloneSquaddie(checker *C) {

}

type SquaddieCloneSuite struct {
	squaddieRepository *squaddie.Repository
	base               *squaddie.Squaddie
	attackA         *power.Power
}

var _ = Suite(&SquaddieCloneSuite{})

func (suite *SquaddieCloneSuite) SetUpTest(checker *C) {
	suite.base = squaddieBuilder.Builder().WithName("Base").Build()
	suite.squaddieRepository = squaddie.NewSquaddieRepository()
	suite.squaddieRepository.AddSquaddies([]*squaddie.Squaddie{suite.base})

	suite.attackA = powerBuilder.Builder().WithName("Attack Formation A").Build()
}

func (suite *SquaddieCloneSuite) TestCloneHasAffiliationAndNameNotID(checker *C) {
	suite.base.Identification.SquaddieAffiliation = squaddie.Enemy
	clone, err := suite.squaddieRepository.CloneSquaddieWithNewID(suite.base, "")
	checker.Assert(err, IsNil)
	checker.Assert(clone.Name(), Equals, suite.base.Name())
	checker.Assert(clone.Affiliation(), Equals, suite.base.Affiliation())
	checker.Assert(clone.ID(), Not(Equals), suite.base.ID())
}

func (suite *SquaddieCloneSuite) TestCloneUsesGivenID(checker *C) {
	clone, _ := suite.squaddieRepository.CloneSquaddieWithNewID(suite.base, "12345")
	checker.Assert(clone.ID(), Equals, "12345")
}

func (suite *SquaddieCloneSuite) TestCloneCopiesBasicStats(checker *C) {
	suite.base.Defense.SquaddieCurrentHitPoints = 1
	suite.base.Defense.SquaddieMaxHitPoints += 5
	suite.base.Defense.SquaddieCurrentBarrier = 2
	suite.base.Defense.SquaddieMaxBarrier += 5

	suite.base.Offense.SquaddieAim = 2
	suite.base.Offense.SquaddieStrength = 3
	suite.base.Offense.SquaddieMind = 4
	suite.base.Defense.SquaddieDodge = 5
	suite.base.Defense.SquaddieDeflect = 6
	suite.base.Defense.SquaddieArmor = 7

	clone, _ := suite.squaddieRepository.CloneSquaddieWithNewID(suite.base, "")
	checker.Assert(clone.CurrentHitPoints(), Equals, suite.base.CurrentHitPoints())
	checker.Assert(clone.MaxHitPoints(), Equals, suite.base.MaxHitPoints())
	checker.Assert(clone.Aim(), Equals, suite.base.Aim())
	checker.Assert(clone.Strength(), Equals, suite.base.Strength())
	checker.Assert(clone.Mind(), Equals, suite.base.Mind())
	checker.Assert(clone.Dodge(), Equals, suite.base.Dodge())
	checker.Assert(clone.Deflect(), Equals, suite.base.Deflect())
	checker.Assert(clone.CurrentBarrier(), Equals, suite.base.CurrentBarrier())
	checker.Assert(clone.MaxBarrier(), Equals, suite.base.MaxBarrier())
	checker.Assert(clone.Armor(), Equals, suite.base.Armor())
}

func (suite *SquaddieCloneSuite) TestCloneCopiesMovement(checker *C) {
	suite.base.Movement = squaddie.Movement{
		SquaddieMovementDistance:     suite.base.Movement.SquaddieMovementDistance + 2,
		SquaddieMovementType:         squaddie.Fly,
		SquaddieMovementCanHitAndRun: true,
	}

	clone, _ := suite.squaddieRepository.CloneSquaddieWithNewID(suite.base, "")
	checker.Assert(clone.Movement.SquaddieMovementDistance, Equals, suite.base.Movement.SquaddieMovementDistance)
	checker.Assert(clone.Movement.SquaddieMovementType, Equals, suite.base.Movement.SquaddieMovementType)
	checker.Assert(clone.Movement.SquaddieMovementCanHitAndRun, Equals, suite.base.Movement.SquaddieMovementCanHitAndRun)
}

func (suite *SquaddieCloneSuite) TestCloneCopiesPowers(checker *C) {
	suite.base.PowerCollection.AddInnatePower(suite.attackA)
	clone, _ := suite.squaddieRepository.CloneSquaddieWithNewID(suite.base, "")

	attackIDNamePairs := clone.PowerCollection.GetInnatePowerIDNames()
	checker.Assert(len(attackIDNamePairs), Equals, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, suite.attackA.Name)
	checker.Assert(attackIDNamePairs[0].ID, Equals, suite.attackA.ID)
}

func (suite *SquaddieCloneSuite) TestCloneCopiesClasses(checker *C) {
	initialClass := &squaddieclass.Class{
		ID:                "initial",
		Name:              "Initial Class",
		BaseClassRequired: false,
		InitialBigLevelID: "",
	}
	advancedClass := &squaddieclass.Class{
		ID:                "advanced",
		Name:              "Advanced Class",
		BaseClassRequired: true,
		InitialBigLevelID: "advanceLevel0",
	}

	suite.base.ClassProgress.AddClass(initialClass)
	suite.base.ClassProgress.AddClass(advancedClass)
	suite.base.ClassProgress.SetBaseClassIfNoBaseClass(initialClass.ID)
	suite.base.ClassProgress.MarkLevelUpBenefitAsConsumed(initialClass.ID, "initialLevel0")
	suite.base.ClassProgress.MarkLevelUpBenefitAsConsumed(initialClass.ID, "initialLevel1")
	suite.base.ClassProgress.MarkLevelUpBenefitAsConsumed(initialClass.ID, "initialLevel2")

	clone, _ := suite.squaddieRepository.CloneSquaddieWithNewID(suite.base, "")
	checker.Assert(clone.ClassProgress.BaseClassID, Equals, suite.base.ClassProgress.BaseClassID)
	checker.Assert(clone.ClassProgress.CurrentClassID, Equals, suite.base.ClassProgress.CurrentClassID)
	for classID, levelsConsumed := range suite.base.ClassProgress.ClassLevelsConsumed {
		checker.Assert(clone.ClassProgress.ClassLevelsConsumed[classID], NotNil)

		cloneLevelsConsumed := clone.ClassProgress.ClassLevelsConsumed[classID]
		checker.Assert(cloneLevelsConsumed, Not(Equals), levelsConsumed)
		checker.Assert(cloneLevelsConsumed, DeepEquals, levelsConsumed)
	}
}
