package squaddie_test

import (
	"bytes"
	"fmt"
	"github.com/cserrant/terosbattleserver/entity/power"
	"github.com/cserrant/terosbattleserver/entity/squaddie"
	"github.com/cserrant/terosbattleserver/entity/squaddieclass"
	. "gopkg.in/check.v1"
)

type SquaddieRepositorySuite struct {
	squaddieFactory *squaddie.Repository
	teros           *squaddie.Squaddie
}

var _ = Suite(&SquaddieRepositorySuite{})

func (suite *SquaddieRepositorySuite) SetUpTest(checker *C) {
	suite.squaddieFactory = squaddie.NewSquaddieRepository()
	suite.teros = squaddie.NewSquaddie("Teros")
	suite.teros.Defense.Armor = 2
	suite.teros.Defense.Dodge = 3
	suite.teros.Defense.Deflect = 4
	suite.teros.Defense.MaxBarrier = 1
}

func (suite *SquaddieRepositorySuite) TestUseJSONSource(checker *C) {
	checker.Assert(suite.squaddieFactory.GetNumberOfSquaddies(), Equals, 0)
	jsonByteStream := []byte(`[{
				"identification": {
					"name": "teros",
					"affiliation": "Player"
				}
			}]`)
	success, _ := suite.squaddieFactory.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.squaddieFactory.GetNumberOfSquaddies(), Equals, 1)
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
	success, _ := suite.squaddieFactory.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)

	suite.teros = suite.squaddieFactory.GetSquaddieByID("squaddieID")
	checker.Assert(suite.teros, NotNil)
	checker.Assert(suite.teros.Identification.Name, Equals, "teros")
	checker.Assert(suite.teros.Offense.Aim, Equals, 5)

	missingno := suite.squaddieFactory.GetSquaddieByID("Does not exist")
	checker.Assert(missingno, IsNil)
}

func (suite *SquaddieRepositorySuite) TestGetExistingSquaddieUsingID(checker *C) {
	jsonByteStream := []byte(`[{
				"identification": {
					"ID": "12345",
					"name": "teros",
					"affiliation": "Player"
				},
				"offense": {
					"aim": 5
				}
			}]`)
	success, _ := suite.squaddieFactory.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)

	suite.teros = suite.squaddieFactory.GetSquaddieByID("12345")
	checker.Assert(suite.teros, NotNil)
	checker.Assert(suite.teros.Identification.Name, Equals, "teros")
	checker.Assert(suite.teros.Offense.Aim, Equals, 5)

	missingno := suite.squaddieFactory.GetSquaddieByID("Does not exist")
	checker.Assert(missingno, IsNil)
}

func (suite *SquaddieRepositorySuite) TestClonedSquaddiesHaveDifferentID(checker *C) {
	jsonByteStream := []byte(`[{
				"identification": {
					"ID": "terosID",
					"name": "teros",
					"affiliation": "Player"
				},
				"offense": {
					"aim": 5
				}
			}]`)
	success, _ := suite.squaddieFactory.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)

	teros0 := suite.squaddieFactory.GetSquaddieByID("terosID")
	teros1 := suite.squaddieFactory.GetSquaddieByID("terosID")
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
	success, _ := suite.squaddieFactory.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)

	suite.teros = suite.squaddieFactory.GetSquaddieByID("terosSquaddieID")
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
	success, err := suite.squaddieFactory.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, false)
	checker.Assert(err, ErrorMatches, "squaddie has unknown affiliation: 'Unknown Affiliation'")
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

	success, _ := suite.squaddieFactory.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.squaddieFactory.GetNumberOfSquaddies(), Equals, 4)

	soldier := suite.squaddieFactory.GetSquaddieByID("Soldier")
	checker.Assert(soldier.Identification.Name, Equals, "Soldier")
	checker.Assert(soldier.Movement.GetMovementDistancePerRound(), Equals, 5)
	checker.Assert(soldier.Movement.GetMovementType(), Equals, squaddie.MovementType(squaddie.Foot))
	checker.Assert(soldier.Movement.CanHitAndRun(), Equals, false)

	scout := suite.squaddieFactory.GetSquaddieByID("Scout")
	checker.Assert(scout.Identification.Name, Equals, "Scout")
	checker.Assert(scout.Movement.GetMovementDistancePerRound(), Equals, 4)
	checker.Assert(scout.Movement.GetMovementType(), Equals, squaddie.MovementType(squaddie.Light))
	checker.Assert(scout.Movement.CanHitAndRun(), Equals, false)

	bird := suite.squaddieFactory.GetSquaddieByID("Bird")
	checker.Assert(bird.Identification.Name, Equals, "Bird")
	checker.Assert(bird.Movement.GetMovementDistancePerRound(), Equals, 3)
	checker.Assert(bird.Movement.GetMovementType(), Equals, squaddie.MovementType(squaddie.Fly))
	checker.Assert(bird.Movement.CanHitAndRun(), Equals, true)

	teleporter := suite.squaddieFactory.GetSquaddieByID("Teleporter")
	checker.Assert(teleporter.Identification.Name, Equals, "Teleporter")
	checker.Assert(teleporter.Movement.GetMovementDistancePerRound(), Equals, 2)
	checker.Assert(teleporter.Movement.GetMovementType(), Equals, squaddie.MovementType(squaddie.Teleport))
	checker.Assert(teleporter.Movement.CanHitAndRun(), Equals, false)
}

func (suite *SquaddieRepositorySuite) TestCanGetExisitingSquaddies(checker *C) {
	originalSquaddie := squaddie.NewSquaddie("Original")
	suite.squaddieFactory.AddSquaddies([]*squaddie.Squaddie{originalSquaddie})
	referencedSquaddie := suite.squaddieFactory.GetOriginalSquaddieByID(originalSquaddie.Identification.ID)
	checker.Assert(referencedSquaddie, Equals, originalSquaddie)
}

func (suite *SquaddieRepositorySuite) TestMarshallIntoJSON(checker *C) {
	byteStream, err := suite.squaddieFactory.MarshalSquaddieIntoJSON(suite.teros)
	checker.Assert(err, IsNil)

	hasIDNameAffiliationInJSON := bytes.Contains(byteStream, []byte(fmt.Sprintf(`"id":"%s","name":"Teros","affiliation":"Player"`, suite.teros.Identification.ID)))
	checker.Assert(hasIDNameAffiliationInJSON, Equals, true)

	hasDefaultHitPointsInJSON := bytes.Contains(byteStream, []byte(`"current_hit_points":5,"max_hit_points":5`))
	checker.Assert(hasDefaultHitPointsInJSON, Equals, true)

	hasOffensiveStatsInJSON := bytes.Contains(byteStream, []byte(`"aim":0,"strength":0,"mind":0`))
	checker.Assert(hasOffensiveStatsInJSON, Equals, true)

	hasDefensiveStatsInJSON := bytes.Contains(byteStream, []byte(`"dodge":3,"deflect":4,"current_barrier":0,"max_barrier":1,"armor":2`))
	checker.Assert(hasDefensiveStatsInJSON, Equals, true)

	hasDefaultMovementInJSON := bytes.Contains(byteStream, []byte(`"movement":{"distance":3,"type":"foot","hit_and_run":false}`))
	checker.Assert(hasDefaultMovementInJSON, Equals, true)

	hasNoPowersInJSON := bytes.Contains(byteStream, []byte(`"powers":[]`))
	checker.Assert(hasNoPowersInJSON, Equals, true)

	hasNoClassesInJSON := bytes.Contains(byteStream, []byte(`"base_class":"","current_class":"","class_levels":{}`))
	checker.Assert(hasNoClassesInJSON, Equals, true)
}

func (suite *SquaddieRepositorySuite) TestMarshallSquaddieMovement(checker *C) {
	suite.teros.Movement.Type = squaddie.Teleport
	suite.teros.Movement.Distance = 8
	suite.teros.Movement.HitAndRun = true

	byteStream, err := suite.squaddieFactory.MarshalSquaddieIntoJSON(suite.teros)
	checker.Assert(err, IsNil)
	movementJSON := `"movement":{"distance":8,"type":"teleport","hit_and_run":true}`
	containsPowersJSON := bytes.Contains(byteStream, []byte(movementJSON))
	checker.Assert(containsPowersJSON, Equals, true)
}

func (suite *SquaddieRepositorySuite) TestMarshallSquaddiePowers(checker *C) {
	attackA := power.NewPower("Attack Formation A")
	suite.teros.PowerCollection.AddInnatePower(attackA)
	byteStream, err := suite.squaddieFactory.MarshalSquaddieIntoJSON(suite.teros)
	checker.Assert(err, IsNil)

	powersJSON := fmt.Sprintf(`"powers":[{"name":"Attack Formation A","id":"%s"}]`, attackA.ID)
	containsPowersJSON := bytes.Contains(byteStream, []byte(powersJSON))
	checker.Assert(containsPowersJSON, Equals, true)
}

func (suite *SquaddieRepositorySuite) TestLoadSquaddieByYAML(checker *C) {
	checker.Assert(suite.squaddieFactory.GetNumberOfSquaddies(), Equals, 0)
	yamlByteStream := []byte(`-
  identification:
    name: teros
    affiliation: Player
`)
	suite.squaddieFactory.AddYAMLSource(yamlByteStream)
	checker.Assert(suite.squaddieFactory.GetNumberOfSquaddies(), Equals, 1)
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
	success, _ := suite.squaddieFactory.AddYAMLSource(yamlByteStream)
	checker.Assert(success, Equals, true)

	suite.teros = suite.squaddieFactory.GetSquaddieByID("terosSquaddieID")
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
    affiliation: Unknown Affiliation`)
	success, err := suite.squaddieFactory.AddYAMLSource(yamlByteStream)
	checker.Assert(success, Equals, false)
	checker.Assert(err, NotNil)
	checker.Assert(err.Error(), Equals, "squaddie has unknown affiliation: 'Unknown Affiliation'")
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

	success, _ := suite.squaddieFactory.AddYAMLSource(yamlByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(suite.squaddieFactory.GetNumberOfSquaddies(), Equals, 4)

	soldier := suite.squaddieFactory.GetSquaddieByID("Soldier")
	checker.Assert(soldier.Identification.Name, Equals, "Soldier")
	checker.Assert(soldier.Movement.GetMovementDistancePerRound(), Equals, 5)
	checker.Assert(soldier.Movement.GetMovementType(), Equals, squaddie.MovementType(squaddie.Foot))
	checker.Assert(soldier.Movement.CanHitAndRun(), Equals, false)

	scout := suite.squaddieFactory.GetSquaddieByID("Scout")
	checker.Assert(scout.Identification.Name, Equals, "Scout")
	checker.Assert(scout.Movement.GetMovementDistancePerRound(), Equals, 4)
	checker.Assert(scout.Movement.GetMovementType(), Equals, squaddie.MovementType(squaddie.Light))
	checker.Assert(scout.Movement.CanHitAndRun(), Equals, false)

	bird := suite.squaddieFactory.GetSquaddieByID("Bird")
	checker.Assert(bird.Identification.Name, Equals, "Bird")
	checker.Assert(bird.Movement.GetMovementDistancePerRound(), Equals, 3)
	checker.Assert(bird.Movement.GetMovementType(), Equals, squaddie.MovementType(squaddie.Fly))
	checker.Assert(bird.Movement.CanHitAndRun(), Equals, true)

	teleporter := suite.squaddieFactory.GetSquaddieByID("Teleporter")
	checker.Assert(teleporter.Identification.Name, Equals, "Teleporter")
	checker.Assert(teleporter.Movement.GetMovementDistancePerRound(), Equals, 2)
	checker.Assert(teleporter.Movement.GetMovementType(), Equals, squaddie.MovementType(squaddie.Teleport))
	checker.Assert(teleporter.Movement.CanHitAndRun(), Equals, false)
}

func (suite *SquaddieRepositorySuite) TestAddSquaddieDirectly(checker *C) {
	success, err := suite.squaddieFactory.AddSquaddies([]*squaddie.Squaddie{squaddie.NewSquaddie("Generic")})
	checker.Assert(success, Equals, true)
	checker.Assert(err, IsNil)
	checker.Assert(suite.squaddieFactory.GetNumberOfSquaddies(), Equals, 1)
}

func (suite *SquaddieRepositorySuite) TestCloneSquaddie(checker *C) {

}

type SquaddieCloneSuite struct {
	squaddieFactory *squaddie.Repository
	base            *squaddie.Squaddie
}

var _ = Suite(&SquaddieCloneSuite{})

func (suite *SquaddieCloneSuite) SetUpTest(checker *C) {
	suite.base = squaddie.NewSquaddie("Base")
	suite.squaddieFactory = squaddie.NewSquaddieRepository()
	suite.squaddieFactory.AddSquaddies([]*squaddie.Squaddie{suite.base})
}

func (suite *SquaddieCloneSuite) TestCloneHasAffiliationAndNameNotID(checker *C) {
	suite.base.Identification.Affiliation = squaddie.Enemy
	clone, err := suite.squaddieFactory.CloneSquaddieWithNewID(suite.base, "")
	checker.Assert(err, IsNil)
	checker.Assert(clone.Identification.Name, Equals, suite.base.Identification.Name)
	checker.Assert(clone.Identification.Affiliation, Equals, suite.base.Identification.Affiliation)
	checker.Assert(clone.Identification.ID, Not(Equals), suite.base.Identification.ID)
}

func (suite *SquaddieCloneSuite) TestCloneUsesGivenID(checker *C) {
	clone, _ := suite.squaddieFactory.CloneSquaddieWithNewID(suite.base, "12345")
	checker.Assert(clone.Identification.ID, Equals, "12345")
}

func (suite *SquaddieCloneSuite) TestCloneUsesName(checker *C) {
	clone, err := suite.squaddieFactory.CloneAndRenameSquaddie(suite.base, "ClonedSquaddie", "12345")
	checker.Assert(err, IsNil)
	checker.Assert(suite.base.Identification.Name, Equals, "Base")
	checker.Assert(clone.Identification.Name, Equals, "ClonedSquaddie")
	checker.Assert(clone.Identification.ID, Not(Equals), suite.base.Identification.ID)
}

func (suite *SquaddieCloneSuite) TestRenameCloneNeedsAName(checker *C) {
	clone, err := suite.squaddieFactory.CloneAndRenameSquaddie(suite.base, "", "12345")
	checker.Assert(err, ErrorMatches, `cannot clone squaddie "Base" without a name`)
	checker.Assert(clone, IsNil)
}

func (suite *SquaddieCloneSuite) TestCloneCopiesBasicStats(checker *C) {
	suite.base.Defense.CurrentHitPoints = 1
	suite.base.Defense.MaxHitPoints += 5
	suite.base.Defense.CurrentBarrier = 2
	suite.base.Defense.MaxBarrier += 5

	suite.base.Offense.Aim = 2
	suite.base.Offense.Strength = 3
	suite.base.Offense.Mind = 4
	suite.base.Defense.Dodge = 5
	suite.base.Defense.Deflect = 6
	suite.base.Defense.Armor = 7

	clone, _ := suite.squaddieFactory.CloneSquaddieWithNewID(suite.base, "")
	checker.Assert(clone.Defense.CurrentHitPoints, Equals, suite.base.Defense.CurrentHitPoints)
	checker.Assert(clone.Defense.MaxHitPoints, Equals, suite.base.Defense.MaxHitPoints)
	checker.Assert(clone.Offense.Aim, Equals, suite.base.Offense.Aim)
	checker.Assert(clone.Offense.Strength, Equals, suite.base.Offense.Strength)
	checker.Assert(clone.Offense.Mind, Equals, suite.base.Offense.Mind)
	checker.Assert(clone.Defense.Dodge, Equals, suite.base.Defense.Dodge)
	checker.Assert(clone.Defense.Deflect, Equals, suite.base.Defense.Deflect)
	checker.Assert(clone.Defense.CurrentBarrier, Equals, suite.base.Defense.CurrentBarrier)
	checker.Assert(clone.Defense.MaxBarrier, Equals, suite.base.Defense.MaxBarrier)
	checker.Assert(clone.Defense.Armor, Equals, suite.base.Defense.Armor)
}

func (suite *SquaddieCloneSuite) TestCloneCopiesMovement(checker *C) {
	suite.base.Movement = squaddie.Movement{
		Distance:  suite.base.Movement.Distance + 2,
		Type:      squaddie.Fly,
		HitAndRun: true,
	}

	clone, _ := suite.squaddieFactory.CloneSquaddieWithNewID(suite.base, "")
	checker.Assert(clone.Movement.Distance, Equals, suite.base.Movement.Distance)
	checker.Assert(clone.Movement.Type, Equals, suite.base.Movement.Type)
	checker.Assert(clone.Movement.HitAndRun, Equals, suite.base.Movement.HitAndRun)
}

func (suite *SquaddieCloneSuite) TestCloneCopiesPowers(checker *C) {
	attackA := power.NewPower("Attack Formation A")
	suite.base.PowerCollection.AddInnatePower(attackA)
	clone, _ := suite.squaddieFactory.CloneSquaddieWithNewID(suite.base, "")

	attackIDNamePairs := clone.PowerCollection.GetInnatePowerIDNames()
	checker.Assert(len(attackIDNamePairs), Equals, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "Attack Formation A")
	checker.Assert(attackIDNamePairs[0].ID, Equals, attackA.ID)
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

	clone, _ := suite.squaddieFactory.CloneSquaddieWithNewID(suite.base, "")
	checker.Assert(clone.ClassProgress.BaseClassID, Equals, suite.base.ClassProgress.BaseClassID)
	checker.Assert(clone.ClassProgress.CurrentClass, Equals, suite.base.ClassProgress.CurrentClass)
	for classID, levelsConsumed := range suite.base.ClassProgress.ClassLevelsConsumed {
		checker.Assert(clone.ClassProgress.ClassLevelsConsumed[classID], NotNil)

		cloneLevelsConsumed := clone.ClassProgress.ClassLevelsConsumed[classID]
		checker.Assert(cloneLevelsConsumed, Not(Equals), levelsConsumed)
		checker.Assert(cloneLevelsConsumed, DeepEquals, levelsConsumed)
	}
}
