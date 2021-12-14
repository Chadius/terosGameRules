package squaddie_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	squaddieClassBuilder "github.com/chadius/terosbattleserver/entity/squaddieclass"
	. "gopkg.in/check.v1"
	"reflect"
)

type SquaddieRepositorySuite struct {
	squaddieRepository *squaddie.Repository
	teros              *squaddie.Squaddie
	attackA            *power.Power
}

var _ = Suite(&SquaddieRepositorySuite{})

func (suite *SquaddieRepositorySuite) SetUpTest(checker *C) {
	suite.squaddieRepository = squaddie.NewSquaddieRepository()
	suite.teros = squaddie.NewSquaddieBuilder().Teros().Armor(2).Dodge(3).Deflect(4).Barrier(1).Build()
	suite.attackA = power.NewPowerBuilder().WithName("attack Formation A").Build()
}

func (suite *SquaddieRepositorySuite) TestCloneSquaddies(checker *C) {
	jsonByteStream := []byte(`[{
		"id": "squaddieID",
		"name": "teros",
		"affiliation": "player",
		"aim": 5
	}]`)
	err := suite.squaddieRepository.AddSquaddiesUsingJSON(jsonByteStream)
	checker.Assert(err, IsNil)

	suite.teros = suite.squaddieRepository.GetSquaddieByID("squaddieID")
	checker.Assert(suite.teros, NotNil)
	checker.Assert(suite.teros.Name(), Equals, "teros")
	checker.Assert(suite.teros.Aim(), Equals, 5)

	missingno := suite.squaddieRepository.GetSquaddieByID("Does not exist")
	checker.Assert(missingno, IsNil)
}

func (suite *SquaddieRepositorySuite) TestGetExistingSquaddieUsingID(checker *C) {
	jsonByteStream := []byte(`[{
		"id": "12345",
		"name": "teros",
		"affiliation": "player",
		"aim": 5
	}]`)
	err := suite.squaddieRepository.AddSquaddiesUsingJSON(jsonByteStream)
	checker.Assert(err, IsNil)

	suite.teros = suite.squaddieRepository.GetSquaddieByID("12345")
	checker.Assert(suite.teros, NotNil)
	checker.Assert(suite.teros.Name(), Equals, "teros")
	checker.Assert(suite.teros.Aim(), Equals, 5)

	missingno := suite.squaddieRepository.GetSquaddieByID("Does not exist")
	checker.Assert(missingno, IsNil)
}

func (suite *SquaddieRepositorySuite) TestClonedSquaddiesHaveDifferentID(checker *C) {
	jsonByteStream := []byte(`[{
		"id": "terosID",
		"name": "teros",
		"affiliation": "player",
		"aim": 5
	}]`)
	err := suite.squaddieRepository.AddSquaddiesUsingJSON(jsonByteStream)
	checker.Assert(err, IsNil)

	teros0 := suite.squaddieRepository.GetSquaddieByID("terosID")
	teros1 := suite.squaddieRepository.GetSquaddieByID("terosID")
	checker.Assert(teros0, Not(Equals), teros1)
}

func (suite *SquaddieRepositorySuite) TestLoadClassLevels(checker *C) {
	jsonByteStream := []byte(`[{
		"id": "terosSquaddieID",
		"name": "teros",
		"affiliation": "player",
		"aim": 5,
		"class_progress": [
			{
				"class_id": "1",
				"class_name":"Mage",
				"levels_gained": ["123"]
			},
			{
				"class_id": "2",
				"class_name":"Dimension Walker"
			}
		]
	}]`)
	err := suite.squaddieRepository.AddSquaddiesUsingJSON(jsonByteStream)
	checker.Assert(err, IsNil)

	suite.teros = suite.squaddieRepository.GetSquaddieByID("terosSquaddieID")
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{"1": 1, "2": 0})
}

func (suite *SquaddieRepositorySuite) TestCreateSquaddiesWithMovement(checker *C) {
	jsonByteStream := []byte(`[
			{
				"id": "Soldier",
				"name": "Soldier",
				"affiliation": "player",
				"movement_distance": 5,
				"movement_type": "foot"
			},
			{
				"id": "Scout",
				"name": "Scout",
				"affiliation": "player",
				"movement_distance": 4,
				"movement_type": "light"
			},
			{
				"id": "Bird",
				"name": "Bird",
				"affiliation": "player",
				"movement_distance": 3,
				"movement_type": "fly",
				"hit_and_run": true
			},
			{
				"id": "Teleporter",
				"name": "Teleporter",
				"affiliation": "player",
				"movement_distance": 2,
				"movement_type": "teleport"
			}
		]`)

	err := suite.squaddieRepository.AddSquaddiesUsingJSON(jsonByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(suite.squaddieRepository.GetNumberOfSquaddies(), Equals, 4)

	soldier := suite.squaddieRepository.GetSquaddieByID("Soldier")
	checker.Assert(soldier.Name(), Equals, "Soldier")
	checker.Assert(soldier.Movement.MovementDistance(), Equals, 5)
	checker.Assert(soldier.Movement.MovementType(), Equals, squaddie.Foot)
	checker.Assert(soldier.Movement.CanHitAndRun(), Equals, false)

	scout := suite.squaddieRepository.GetSquaddieByID("Scout")
	checker.Assert(scout.Name(), Equals, "Scout")
	checker.Assert(scout.Movement.MovementDistance(), Equals, 4)
	checker.Assert(scout.Movement.MovementType(), Equals, squaddie.Light)
	checker.Assert(scout.Movement.CanHitAndRun(), Equals, false)

	bird := suite.squaddieRepository.GetSquaddieByID("Bird")
	checker.Assert(bird.Name(), Equals, "Bird")
	checker.Assert(bird.Movement.MovementDistance(), Equals, 3)
	checker.Assert(bird.Movement.MovementType(), Equals, squaddie.Fly)
	checker.Assert(bird.Movement.CanHitAndRun(), Equals, true)

	teleporter := suite.squaddieRepository.GetSquaddieByID("Teleporter")
	checker.Assert(teleporter.Name(), Equals, "Teleporter")
	checker.Assert(teleporter.Movement.MovementDistance(), Equals, 2)
	checker.Assert(teleporter.Movement.MovementType(), Equals, squaddie.Teleport)
	checker.Assert(teleporter.Movement.CanHitAndRun(), Equals, false)
}

func (suite *SquaddieRepositorySuite) TestCanGetExistingSquaddies(checker *C) {
	originalSquaddie := squaddie.NewSquaddieBuilder().WithName("Original").AsAlly().Build()
	suite.squaddieRepository.AddSquaddies([]*squaddie.Squaddie{originalSquaddie})
	referencedSquaddie := suite.squaddieRepository.GetOriginalSquaddieByID(originalSquaddie.ID())
	checker.Assert(referencedSquaddie, Equals, originalSquaddie)
}

func (suite *SquaddieRepositorySuite) TestAddSquaddieDirectly(checker *C) {
	success, err := suite.squaddieRepository.AddSquaddies([]*squaddie.Squaddie{
		squaddie.NewSquaddieBuilder().Build(),
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
	suite.base = squaddie.NewSquaddieBuilder().WithName("Base").Build()
	suite.squaddieRepository = squaddie.NewSquaddieRepository()
	suite.squaddieRepository.AddSquaddies([]*squaddie.Squaddie{suite.base})
}

func (suite *SquaddieCloneSuite) TestCloneHasAffiliationAndNameNotID(checker *C) {
	originalSquaddie := squaddie.NewSquaddieBuilder().WithName("Base").AsEnemy().Build()
	clone, err := suite.squaddieRepository.CloneSquaddieWithNewID(originalSquaddie, "")
	checker.Assert(err, IsNil)
	checker.Assert(clone.Name(), Equals, originalSquaddie.Name())
	checker.Assert(
		reflect.TypeOf(clone.AffiliationLogic()).String(),
		Equals,
		reflect.TypeOf(originalSquaddie.AffiliationLogic()).String(),
	)
	checker.Assert(clone.ID(), Not(Equals), originalSquaddie.ID())
}

func (suite *SquaddieCloneSuite) TestCloneUsesGivenID(checker *C) {
	originalSquaddie := squaddie.NewSquaddieBuilder().WithName("Base").AsEnemy().Build()
	clone, _ := suite.squaddieRepository.CloneSquaddieWithNewID(originalSquaddie, "12345")
	checker.Assert(clone.ID(), Equals, "12345")
}

func (suite *SquaddieCloneSuite) TestCloneCopiesBasicStats(checker *C) {
	originalSquaddie := squaddie.NewSquaddieBuilder().WithName("Base").
		HitPoints(9).Barrier(7).
		Aim(2).Strength(3).Mind(5).Dodge(11).Deflect(13).Armor(17).
		Build()
	originalSquaddie.ReduceHitPoints(originalSquaddie.MaxHitPoints() - 1)
	originalSquaddie.ReduceBarrier(originalSquaddie.MaxBarrier() - 2)

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
	originalSquaddie := squaddie.NewSquaddieBuilder().WithName("Base").
		MoveDistance(2).MovementFly().CanHitAndRun().Build()

	clone, _ := suite.squaddieRepository.CloneSquaddieWithNewID(originalSquaddie, "")

	checker.Assert(clone.MovementDistance(), Equals, originalSquaddie.MovementDistance())
	checker.Assert(clone.MovementType(), Equals, originalSquaddie.MovementType())
	checker.Assert(clone.MovementCanHitAndRun(), Equals, originalSquaddie.MovementCanHitAndRun())
}

func (suite *SquaddieCloneSuite) TestCloneCopiesPowers(checker *C) {
	attackA := power.NewPowerBuilder().WithName("attack Formation A").Build()
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
  affiliation: player
`)
	squaddieRepository := squaddie.NewSquaddieRepository()
	err := squaddieRepository.AddSquaddiesUsingYAML(yamlByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(squaddieRepository.GetNumberOfSquaddies(), Equals, 1)

	teros := squaddieRepository.GetSquaddieByID("squaddie_teros")
	checker.Assert(teros.Name(), Equals, "teros")
}

func (suite *SquaddieLoadDataStreamUsingBuilders) TestSquaddieHasExpectedIdentification(checker *C) {
	yamlByteStream := []byte(`-
  id: squaddie_bandit
  name: Bandit
  affiliation: enemy
`)
	squaddieRepository := squaddie.NewSquaddieRepository()
	err := squaddieRepository.AddSquaddiesUsingYAML(yamlByteStream)
	checker.Assert(err, IsNil)

	loadedSquaddie := squaddieRepository.GetSquaddieByID("squaddie_bandit")
	checker.Assert(loadedSquaddie.Name(), Equals, "Bandit")
	checker.Assert(reflect.TypeOf(loadedSquaddie.AffiliationLogic()).String(), Equals, "*affiliation.Enemy")
}

func (suite *SquaddieLoadDataStreamUsingBuilders) TestSquaddieHasExpectedOffense(checker *C) {
	yamlByteStream := []byte(`-
  id: squaddie_offense
  aim: 2
  strength: 3
  mind: 5
`)
	squaddieRepository := squaddie.NewSquaddieRepository()
	err := squaddieRepository.AddSquaddiesUsingYAML(yamlByteStream)
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
	err := squaddieRepository.AddSquaddiesUsingYAML(yamlByteStream)
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
	err := squaddieRepository.AddSquaddiesUsingYAML(yamlByteStream)
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
	err := squaddieRepository.AddSquaddiesUsingYAML(yamlByteStream)
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
	err := squaddieRepository.AddSquaddiesUsingYAML(yamlByteStream)
	checker.Assert(err, IsNil)
	loadedSquaddie := squaddieRepository.GetSquaddieByID("squaddie_with_class")

	checker.Assert(loadedSquaddie.BaseClassID(), Equals, "class_dirt_farmer")
	checker.Assert(loadedSquaddie.CurrentClassID(), Equals, "class_jedi_knight")

	levelCountsByClass := loadedSquaddie.GetLevelCountsByClass()
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

func (suite *SquaddieLoadDataStreamUsingBuilders) TestSquaddieCanBeBuiltWithJSON(checker *C) {
	jsonByteStream := []byte(`[
{
	"id": "squaddie_id",
	"name": "Buff Citizen",
	"affiliation": "ally",
	"aim": 2,
	"strength": 3,
	"mind": 5,
	"max_hit_points": 7,
	"max_barrier": 11,
	"armor": 13,
	"dodge": 17,
	"deflect": 19,
	"movement_distance": 23,
	"movement_type": "teleport",
	"hit_and_run": true,
	"powers": [
		{
			"name": "Citizen Punch",
			"id": "power_citizen_punch"
		},
		{
			"name": "Citizen Kick",
			"id": "power_citizen_kick"
		}
	],
	"class_progress": [
		{
			"is_current_class": true,
			"class_id": "class_action_bystander",
			"class_name": "Action Bystander",
			"levels_gained": ["levelActionBystander0", "levelActionBystander1"]
		},
		{
			"is_base_class": true,
			"class_id": "class_surviving_bystander",
			"class_name": "Surviving Bystander",
			"levels_gained": ["levelSurvivingBystander0"]
		}
	]
}
]`)
	squaddieRepository := squaddie.NewSquaddieRepository()
	err := squaddieRepository.AddSquaddiesUsingJSON(jsonByteStream)
	checker.Assert(err, IsNil)

	checker.Assert(squaddieRepository.GetNumberOfSquaddies(), Equals, 1)
	loadedSquaddie := squaddieRepository.GetSquaddieByID("squaddie_id")

	checker.Assert(loadedSquaddie.Name(), Equals, "Buff Citizen")
	checker.Assert(reflect.TypeOf(loadedSquaddie.AffiliationLogic()).String(), Equals, "*affiliation.Ally")

	checker.Assert(loadedSquaddie.Aim(), Equals, 2)
	checker.Assert(loadedSquaddie.Strength(), Equals, 3)
	checker.Assert(loadedSquaddie.Mind(), Equals, 5)

	checker.Assert(loadedSquaddie.MaxHitPoints(), Equals, 7)
	checker.Assert(loadedSquaddie.CurrentHitPoints(), Equals, 7)
	checker.Assert(loadedSquaddie.MaxBarrier(), Equals, 11)
	checker.Assert(loadedSquaddie.Armor(), Equals, 13)
	checker.Assert(loadedSquaddie.Dodge(), Equals, 17)
	checker.Assert(loadedSquaddie.Deflect(), Equals, 19)

	checker.Assert(loadedSquaddie.MovementDistance(), Equals, 23)
	checker.Assert(loadedSquaddie.MovementCanHitAndRun(), Equals, true)
	checker.Assert(loadedSquaddie.MovementType(), Equals, squaddie.Teleport)

	squaddiePowers := loadedSquaddie.GetCopyOfPowerReferences()

	checker.Assert(squaddiePowers, HasLen, 2)
	checker.Assert(squaddiePowers[0].PowerID, Equals, "power_citizen_punch")
	checker.Assert(squaddiePowers[0].Name, Equals, "Citizen Punch")

	checker.Assert(squaddiePowers[1].PowerID, Equals, "power_citizen_kick")
	checker.Assert(squaddiePowers[1].Name, Equals, "Citizen Kick")

	checker.Assert(loadedSquaddie.BaseClassID(), Equals, "class_surviving_bystander")
	checker.Assert(loadedSquaddie.CurrentClassID(), Equals, "class_action_bystander")

	levelCountsByClass := loadedSquaddie.GetLevelCountsByClass()
	checker.Assert(levelCountsByClass, HasLen, 2)
	checker.Assert(levelCountsByClass["class_surviving_bystander"], Equals, 1)
	checker.Assert(levelCountsByClass["class_action_bystander"], Equals, 2)

	checker.Assert(loadedSquaddie.IsClassLevelAlreadyUsed("levelActionBystander0"), Equals, true)
	checker.Assert(loadedSquaddie.IsClassLevelAlreadyUsed("levelActionBystander1"), Equals, true)
	checker.Assert(loadedSquaddie.IsClassLevelAlreadyUsed("levelActionBystander2"), Equals, false)

	checker.Assert(loadedSquaddie.IsClassLevelAlreadyUsed("levelSurvivingBystander0"), Equals, true)
	checker.Assert(loadedSquaddie.IsClassLevelAlreadyUsed("levelSurvivingBystander1"), Equals, false)
}
