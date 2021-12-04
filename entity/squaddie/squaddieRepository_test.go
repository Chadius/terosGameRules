package squaddie_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	squaddieClassBuilder "github.com/chadius/terosbattleserver/entity/squaddieclass"
	. "gopkg.in/check.v1"
)

type SquaddieRepositorySuite struct {
	squaddieRepository *squaddie.Repository
	teros              *squaddie.Squaddie
	attackA            *power.Power
}

var _ = Suite(&SquaddieRepositorySuite{})

func (suite *SquaddieRepositorySuite) SetUpTest(checker *C) {
	suite.squaddieRepository = squaddie.NewSquaddieRepository()
	suite.teros = squaddie.Builder().Teros().Armor(2).Dodge(3).Deflect(4).Barrier(1).Build()
	suite.attackA = power.Builder().WithName("attack Formation A").Build()
}

// TODO will I need this function once I'm done refactoring?
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
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{"1": 1, "2": 0})
}

func (suite *SquaddieRepositorySuite) TestStopLoadingSquaddiesWhenInvalid(checker *C) {
	jsonByteStream := []byte(`[{
				"identification": {
					"name": "teros",
					"Affiliation": "Player"
				}
			},{
				"identification": {
					"name": "teros2",
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
	originalSquaddie := squaddie.Builder().WithName("Original").AsAlly().Build()
	suite.squaddieRepository.AddSquaddies([]*squaddie.Squaddie{originalSquaddie})
	referencedSquaddie := suite.squaddieRepository.GetOriginalSquaddieByID(originalSquaddie.ID())
	checker.Assert(referencedSquaddie, Equals, originalSquaddie)
}

// TODO all of these may be deleted soon
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

// TODO
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
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{"1": 1, "2": 0})
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
		squaddie.Builder().Build(),
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
}

var _ = Suite(&SquaddieCloneSuite{})

func (suite *SquaddieCloneSuite) SetUpTest(checker *C) {
	suite.base = squaddie.Builder().WithName("Base").Build()
	suite.squaddieRepository = squaddie.NewSquaddieRepository()
	suite.squaddieRepository.AddSquaddies([]*squaddie.Squaddie{suite.base})
}

func (suite *SquaddieCloneSuite) TestCloneHasAffiliationAndNameNotID(checker *C) {
	originalSquaddie := squaddie.Builder().WithName("Base").AsEnemy().Build()
	clone, err := suite.squaddieRepository.CloneSquaddieWithNewID(originalSquaddie, "")
	checker.Assert(err, IsNil)
	checker.Assert(clone.Name(), Equals, originalSquaddie.Name())
	checker.Assert(clone.Affiliation(), Equals, originalSquaddie.Affiliation())
	checker.Assert(clone.ID(), Not(Equals), originalSquaddie.ID())
}

func (suite *SquaddieCloneSuite) TestCloneUsesGivenID(checker *C) {
	originalSquaddie := squaddie.Builder().WithName("Base").AsEnemy().Build()
	clone, _ := suite.squaddieRepository.CloneSquaddieWithNewID(originalSquaddie, "12345")
	checker.Assert(clone.ID(), Equals, "12345")
}

func (suite *SquaddieCloneSuite) TestCloneCopiesBasicStats(checker *C) {
	originalSquaddie := squaddie.Builder().WithName("Base").
		HitPoints(9).Barrier(7).
		Aim(2).Strength(3).Mind(5).Dodge(11).Deflect(13).Armor(17).
		Build()
	originalSquaddie.Defense.ReduceHitPoints(originalSquaddie.MaxHitPoints() - 1)
	originalSquaddie.Defense.ReduceBarrier(originalSquaddie.MaxBarrier() - 2)

	clone, _ := suite.squaddieRepository.CloneSquaddieWithNewID(originalSquaddie, "")
	checker.Assert(clone.CurrentHitPoints(), Equals, originalSquaddie.CurrentHitPoints())
	checker.Assert(clone.MaxHitPoints(), Equals, originalSquaddie.MaxHitPoints())
	checker.Assert(clone.Aim(), Equals, originalSquaddie.Aim())
	checker.Assert(clone.Strength(), Equals, originalSquaddie.Strength())
	checker.Assert(clone.Mind(), Equals, originalSquaddie.Mind())
	checker.Assert(clone.Dodge(), Equals, originalSquaddie.Dodge())
	checker.Assert(clone.Deflect(), Equals, originalSquaddie.Deflect())
	checker.Assert(clone.CurrentBarrier(), Equals, originalSquaddie.CurrentBarrier())
	checker.Assert(clone.MaxBarrier(), Equals, originalSquaddie.MaxBarrier())
	checker.Assert(clone.Armor(), Equals, originalSquaddie.Armor())
}

func (suite *SquaddieCloneSuite) TestCloneCopiesMovement(checker *C) {
	originalSquaddie := squaddie.Builder().WithName("Base").
		MoveDistance(2).MovementFly().CanHitAndRun().Build()

	clone, _ := suite.squaddieRepository.CloneSquaddieWithNewID(originalSquaddie, "")

	checker.Assert(clone.MovementDistance(), Equals, originalSquaddie.MovementDistance())
	checker.Assert(clone.MovementType(), Equals, originalSquaddie.MovementType())
	checker.Assert(clone.MovementCanHitAndRun(), Equals, originalSquaddie.MovementCanHitAndRun())
}

func (suite *SquaddieCloneSuite) TestCloneCopiesPowers(checker *C) {
	attackA := power.Builder().WithName("attack Formation A").Build()
	suite.base.AddPowerReference(attackA.GetReference())
	clone, _ := suite.squaddieRepository.CloneSquaddieWithNewID(suite.base, "")

	attackIDNamePairs := clone.PowerCollection.GetCopyOfPowerReferences()
	checker.Assert(len(attackIDNamePairs), Equals, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, attackA.Name())
	checker.Assert(attackIDNamePairs[0].PowerID, Equals, attackA.ID())
}

func (suite *SquaddieCloneSuite) TestCloneCopiesClasses(checker *C) {
	initialClass := squaddieClassBuilder.ClassBuilder().WithID("initial").Build()
	advancedClass := squaddieClassBuilder.ClassBuilder().WithID("advanced").RequiresBaseClass().Build()

	suite.base.AddClass(initialClass.GetReference())
	suite.base.AddClass(advancedClass.GetReference())
	suite.base.SetBaseClassIfNoBaseClass(initialClass.ID())
	suite.base.MarkLevelUpBenefitAsConsumed(initialClass.ID(), "initialLevel0")
	suite.base.MarkLevelUpBenefitAsConsumed(initialClass.ID(), "initialLevel1")
	suite.base.MarkLevelUpBenefitAsConsumed(initialClass.ID(), "initialLevel2")

	clone, _ := suite.squaddieRepository.CloneSquaddieWithNewID(suite.base, "")
	checker.Assert(clone.BaseClassID(), Equals, suite.base.BaseClassID())
	checker.Assert(clone.CurrentClassID(), Equals, suite.base.CurrentClassID())
	for classID, levelsConsumed := range *suite.base.ClassLevelsConsumed() {
		checker.Assert((*clone.ClassLevelsConsumed())[classID], NotNil)

		cloneLevelsConsumed := (*clone.ClassLevelsConsumed())[classID]
		checker.Assert(cloneLevelsConsumed, Not(Equals), levelsConsumed)
		checker.Assert(cloneLevelsConsumed, DeepEquals, levelsConsumed)
	}
}

type SquaddieLoadDataStreamUsingBuilders struct{}

var _ = Suite(&SquaddieLoadDataStreamUsingBuilders{})

func (suite *SquaddieLoadDataStreamUsingBuilders) TestLoadSquaddieByYAMLUsingSquaddieBuilder(checker *C) {
	yamlByteStream := []byte(`-
  id: squaddie_teros
  name: teros
  affiliation: Player
`)
	squaddieRepository := squaddie.NewSquaddieRepository()
	err := squaddieRepository.AddYAMLSourceUsingSquaddieBuilder(yamlByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(squaddieRepository.GetNumberOfSquaddies(), Equals, 1)

	teros := squaddieRepository.GetSquaddieByID("squaddie_teros")
	checker.Assert(teros.Name(), Equals, "teros")
}

func (suite *SquaddieLoadDataStreamUsingBuilders) TestSquaddieHasExpectedIdentification(checker *C) {
	yamlByteStream := []byte(`-
  id: squaddie_bandit
  name: Bandit
  affiliation: Enemy
`)
	squaddieRepository := squaddie.NewSquaddieRepository()
	err := squaddieRepository.AddYAMLSourceUsingSquaddieBuilder(yamlByteStream)
	checker.Assert(err, IsNil)

	loadedSquaddie := squaddieRepository.GetSquaddieByID("squaddie_bandit")
	checker.Assert(loadedSquaddie.Name(), Equals, "Bandit")
	checker.Assert(loadedSquaddie.Affiliation(), Equals, squaddie.Enemy)
}

func (suite *SquaddieLoadDataStreamUsingBuilders) TestSquaddieHasExpectedOffense(checker *C) {
	yamlByteStream := []byte(`-
  id: squaddie_offense
  aim: 2
  strength: 3
  mind: 5
`)
	squaddieRepository := squaddie.NewSquaddieRepository()
	err := squaddieRepository.AddYAMLSourceUsingSquaddieBuilder(yamlByteStream)
	checker.Assert(err, IsNil)

	loadedSquaddie := squaddieRepository.GetSquaddieByID("squaddie_offense")
	checker.Assert(loadedSquaddie.Aim(), Equals, 2)
	checker.Assert(loadedSquaddie.Strength(), Equals, 3)
	checker.Assert(loadedSquaddie.Mind(), Equals, 5)
}

func (suite *SquaddieLoadDataStreamUsingBuilders) TestSquaddieHasExpectedDefense(checker *C) {
	yamlByteStream := []byte(`-
  id: squaddie_defense
  max_hit_points: 2
  max_barrier: 3
  armor: 5
  dodge: 7
  deflect: 11
`)
	squaddieRepository := squaddie.NewSquaddieRepository()
	err := squaddieRepository.AddYAMLSourceUsingSquaddieBuilder(yamlByteStream)
	checker.Assert(err, IsNil)

	loadedSquaddie := squaddieRepository.GetSquaddieByID("squaddie_defense")
	checker.Assert(loadedSquaddie.MaxHitPoints(), Equals, 2)
	checker.Assert(loadedSquaddie.CurrentHitPoints(), Equals, 2)
	checker.Assert(loadedSquaddie.MaxBarrier(), Equals, 3)
	checker.Assert(loadedSquaddie.Armor(), Equals, 5)
	checker.Assert(loadedSquaddie.Dodge(), Equals, 7)
	checker.Assert(loadedSquaddie.Deflect(), Equals, 11)
}

func (suite *SquaddieLoadDataStreamUsingBuilders) TestSquaddieHasExpectedMovement(checker *C) {
	yamlByteStream := []byte(`-
  id: squaddie_on_the_move
  hit_and_run: true
  movement_type: fly
  movement_distance: 2
`)
	squaddieRepository := squaddie.NewSquaddieRepository()
	err := squaddieRepository.AddYAMLSourceUsingSquaddieBuilder(yamlByteStream)
	checker.Assert(err, IsNil)

	loadedSquaddie := squaddieRepository.GetSquaddieByID("squaddie_on_the_move")
	checker.Assert(loadedSquaddie.MovementDistance(), Equals, 2)
	checker.Assert(loadedSquaddie.MovementCanHitAndRun(), Equals, true)
	checker.Assert(loadedSquaddie.MovementType(), Equals, squaddie.Fly)
}

func (suite *SquaddieLoadDataStreamUsingBuilders) TestSquaddieHasExpectedPowers(checker *C) {
	yamlByteStream := []byte(`-
  id: squaddie_with_powers
  powers:
  -
    name: Baseball Bat
    id: power_baseball_bat
  -
    name: BasketBrawl
    id: power_basket_brawl
`)
	squaddieRepository := squaddie.NewSquaddieRepository()
	err := squaddieRepository.AddYAMLSourceUsingSquaddieBuilder(yamlByteStream)
	checker.Assert(err, IsNil)
	loadedSquaddie := squaddieRepository.GetSquaddieByID("squaddie_with_powers")
	squaddiePowers := loadedSquaddie.GetCopyOfPowerReferences()

	checker.Assert(squaddiePowers, HasLen, 2)
	checker.Assert(squaddiePowers[0].PowerID, Equals, "power_baseball_bat")
	checker.Assert(squaddiePowers[0].Name, Equals, "Baseball Bat")

	checker.Assert(squaddiePowers[1].PowerID, Equals, "power_basket_brawl")
	checker.Assert(squaddiePowers[1].Name, Equals, "BasketBrawl")
}

func (suite *SquaddieLoadDataStreamUsingBuilders) TestSquaddieHasExpectedClasses(checker *C) {
	yamlByteStream := []byte(`-
  id: squaddie_with_class
  class_progress:
  -
    is_base_class: true
    class_id: class_dirt_farmer
    class_name: Dirt Farmer
    levels_gained: ["levelDirtFarmer0", "levelDirtFarmer1", "levelDirtFarmer2"]
  -
    is_base_class: false
    is_current_class: true
    class_id: class_jedi_knight
    class_name: Jedi Knight
    levels_gained: ["levelJediKnight0"]
`)
	squaddieRepository := squaddie.NewSquaddieRepository()
	err := squaddieRepository.AddYAMLSourceUsingSquaddieBuilder(yamlByteStream)
	checker.Assert(err, IsNil)
	loadedSquaddie := squaddieRepository.GetSquaddieByID("squaddie_with_class")
	levelCountsByClass := loadedSquaddie.GetLevelCountsByClass()

	checker.Assert(loadedSquaddie.BaseClassID(), Equals, "class_dirt_farmer")
	checker.Assert(loadedSquaddie.CurrentClassID(), Equals, "class_jedi_knight")

	checker.Assert(levelCountsByClass, HasLen, 2)
	checker.Assert(levelCountsByClass["class_dirt_farmer"], Equals, 3)
	checker.Assert(levelCountsByClass["class_jedi_knight"], Equals, 1)

	checker.Assert(loadedSquaddie.IsClassLevelAlreadyUsed("levelDirtFarmer0"), Equals, true)
	checker.Assert(loadedSquaddie.IsClassLevelAlreadyUsed("levelDirtFarmer1"), Equals, true)
	checker.Assert(loadedSquaddie.IsClassLevelAlreadyUsed("levelDirtFarmer2"), Equals, true)
	checker.Assert(loadedSquaddie.IsClassLevelAlreadyUsed("levelDirtFarmer3"), Equals, false)

	checker.Assert(loadedSquaddie.IsClassLevelAlreadyUsed("levelJediKnight0"), Equals, true)
	checker.Assert(loadedSquaddie.IsClassLevelAlreadyUsed("levelJediKnight1"), Equals, false)
}
