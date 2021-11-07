package squaddie_test

import (
	powerEntity "github.com/chadius/terosbattleserver/entity/power"
	squaddieEntity "github.com/chadius/terosbattleserver/entity/squaddie"
	classEntity "github.com/chadius/terosbattleserver/entity/squaddieclass"
	"github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	"github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	"github.com/chadius/terosbattleserver/utility/testutility/builder/squaddieclass"
	. "gopkg.in/check.v1"
)

type SquaddieIdentificationBuilder struct{}

var _ = Suite(&SquaddieIdentificationBuilder{})

func (suite *SquaddieIdentificationBuilder) TestBuildSquaddieWithName(checker *C) {
	teros := squaddie.Builder().WithName("Teros").Build()
	checker.Assert("Teros", Equals, teros.Name())
}

func (suite *SquaddieIdentificationBuilder) TestBuildSquaddieWithID(checker *C) {
	teros := squaddie.Builder().WithID("squaddieTeros").Build()
	checker.Assert("squaddieTeros", Equals, teros.ID())
}

func (suite *SquaddieIdentificationBuilder) TestBuildIdentificationAffiliationPlayer(checker *C) {
	teros := squaddie.Builder().AsPlayer().Build()
	checker.Assert(squaddieEntity.Player, Equals, teros.Affiliation())
}

func (suite *SquaddieIdentificationBuilder) TestBuildIdentificationAffiliationEnemy(checker *C) {
	bandit := squaddie.Builder().AsEnemy().Build()
	checker.Assert(squaddieEntity.Enemy, Equals, bandit.Affiliation())
}

func (suite *SquaddieIdentificationBuilder) TestBuildIdentificationAffiliationAlly(checker *C) {
	citizen := squaddie.Builder().AsAlly().Build()
	checker.Assert(squaddieEntity.Ally, Equals, citizen.Affiliation())
}

func (suite *SquaddieIdentificationBuilder) TestBuildIdentificationAffiliationNeutral(checker *C) {
	bomb := squaddie.Builder().AsNeutral().Build()
	checker.Assert(squaddieEntity.Neutral, Equals, bomb.Affiliation())
}

type SquaddieOffenseBuilder struct{}

var _ = Suite(&SquaddieOffenseBuilder{})

func (suite *SquaddieOffenseBuilder) TestBuildSquaddieWithAim(checker *C) {
	teros := squaddie.Builder().Aim(3).Build()
	checker.Assert(3, Equals, teros.Aim())
}

func (suite *SquaddieOffenseBuilder) TestBuildSquaddieWithStrength(checker *C) {
	teros := squaddie.Builder().Strength(2).Build()
	checker.Assert(2, Equals, teros.Strength())
}

func (suite *SquaddieOffenseBuilder) TestBuildSquaddieWithMind(checker *C) {
	teros := squaddie.Builder().Mind(4).Build()
	checker.Assert(4, Equals, teros.Mind())
}

type SquaddieDefenseBuilder struct{}

var _ = Suite(&SquaddieDefenseBuilder{})

func (suite *SquaddieDefenseBuilder) TestBuildSquaddieWithHitPoints(checker *C) {
	teros := squaddie.Builder().HitPoints(9).Build()
	checker.Assert(9, Equals, teros.CurrentHitPoints())
	checker.Assert(9, Equals, teros.MaxHitPoints())
}

func (suite *SquaddieDefenseBuilder) TestBuildSquaddieWithBarrier(checker *C) {
	teros := squaddie.Builder().Barrier(3).Build()
	checker.Assert(3, Equals, teros.MaxBarrier())
}

func (suite *SquaddieDefenseBuilder) TestBuildSquaddieWithArmor(checker *C) {
	teros := squaddie.Builder().Armor(2).Build()
	checker.Assert(2, Equals, teros.Armor())
}

func (suite *SquaddieDefenseBuilder) TestBuildSquaddieWithDodge(checker *C) {
	teros := squaddie.Builder().Dodge(1).Build()
	checker.Assert(1, Equals, teros.Dodge())
}

func (suite *SquaddieDefenseBuilder) TestBuildSquaddieWithDeflect(checker *C) {
	teros := squaddie.Builder().Deflect(4).Build()
	checker.Assert(4, Equals, teros.Deflect())
}

type SquaddieMovementBuilder struct{}

var _ = Suite(&SquaddieMovementBuilder{})

func (suite *SquaddieMovementBuilder) TestBuildWithDistance(checker *C) {
	soldier := squaddie.Builder().MoveDistance(3).Build()
	checker.Assert(3, Equals, soldier.Movement.SquaddieMovementDistance)
}

func (suite *SquaddieMovementBuilder) TestBuildMovementCanHitAndRun(checker *C) {
	runner := squaddie.Builder().CanHitAndRun().Build()
	checker.Assert(true, Equals, runner.Movement.SquaddieMovementCanHitAndRun)
}

func (suite *SquaddieMovementBuilder) TestChangeMovementFoot(checker *C) {
	soldier := squaddie.Builder().MovementFoot().Build()
	checker.Assert(squaddieEntity.Foot, Equals, soldier.Movement.SquaddieMovementType)
}

func (suite *SquaddieMovementBuilder) TestChangeMovementLight(checker *C) {
	ninja := squaddie.Builder().MovementLight().Build()
	checker.Assert(squaddieEntity.Light, Equals, ninja.Movement.SquaddieMovementType)
}

func (suite *SquaddieMovementBuilder) TestChangeMovementFly(checker *C) {
	bird := squaddie.Builder().MovementFly().Build()
	checker.Assert(squaddieEntity.Fly, Equals, bird.Movement.SquaddieMovementType)
}

func (suite *SquaddieMovementBuilder) TestChangeMovementTeleport(checker *C) {
	wizard := squaddie.Builder().MovementTeleport().Build()
	checker.Assert(squaddieEntity.Teleport, Equals, wizard.Movement.SquaddieMovementType)
}

type SquaddiePowerBuilder struct{}

var _ = Suite(&SquaddiePowerBuilder{})

func (suite *SquaddiePowerBuilder) TestBuildAddPowerReference(checker *C) {
	spear := power.Builder().Spear().Build()

	teros := squaddie.Builder().AddPowerByReference(spear.GetReference()).Build()

	checker.Assert(teros.HasPowerWithID(spear.ID()), Equals, true)
}

type SquaddieClassBuilder struct{}

var _ = Suite(&SquaddieClassBuilder{})

func (suite *SquaddieClassBuilder) TestBuildAddClass(checker *C) {
	mageClass := squaddieclass.ClassBuilder().WithID("A class SquaddieID").WithName("mage").WithInitialBigLevelID("level0").Build()
	teros := squaddie.Builder().AddClassByReference(mageClass.GetReference()).Build()
	checker.Assert(true, Equals, teros.HasAddedClass(mageClass.ID()))
}

func (suite *SquaddieClassBuilder) TestBuildSetClass(checker *C) {
	mageClass := squaddieclass.ClassBuilder().WithID("A class SquaddieID").WithName("mage").WithInitialBigLevelID("level0").Build()
	teros := squaddie.Builder().AddClassByReference(mageClass.GetReference()).SetClassByID(mageClass.ID()).Build()
	checker.Assert(mageClass.ID(), Equals, teros.CurrentClassID())
}

type SpecificSquaddieBuilder struct{}

var _ = Suite(&SpecificSquaddieBuilder{})

func (suite *SpecificSquaddieBuilder) TestBuildTeros(checker *C) {
	teros := squaddie.Builder().Teros().Build()
	checker.Assert("Teros", Equals, teros.Name())
}

func (suite *SpecificSquaddieBuilder) TestBuildBandit(checker *C) {
	bandit := squaddie.Builder().Bandit().Build()
	checker.Assert("Bandit", Equals, bandit.Name())
}

func (suite *SpecificSquaddieBuilder) TestBuildLini(checker *C) {
	lini := squaddie.Builder().Lini().Build()
	checker.Assert("Lini", Equals, lini.Name())
}

func (suite *SpecificSquaddieBuilder) TestBuildMysticMage(checker *C) {
	mysticMage := squaddie.Builder().MysticMage().Build()
	checker.Assert("Mystic Mage", Equals, mysticMage.Name())
}

type YAMLBuilderSuite struct {
	yamlData []byte
}

var _ = Suite(&YAMLBuilderSuite{})

func (suite *YAMLBuilderSuite) SetUpTest(checker *C) {
	suite.yamlData = []byte(
		`
id: squaddie_yaml
name: YAML squaddie
affiliation: Enemy
max_hit_points: 2
dodge: 3
deflect: 5
max_barrier: 7
armor: 9
aim: 11
strength: 13
mind: 17
movement_distance: 19
movement_type: light
hit_and_run: true
powers:
- id: shove_id
  name: Shove
- id: bandage_id
  name: Bandage
class_progress:
- class_id: baseClassID
  class_name: Introductory Class
  is_base_class: true
  levels_gained: ["level0", "level1"]
- class_id: currentClassID
  class_name: Advanced Class
  is_current_class: true
  levels_gained: 
  - advanced level0
  - advanced level1
`)
}

func (suite *YAMLBuilderSuite) TestIdentificationMatchesNewSquaddie(checker *C) {
	yamlSquaddie := squaddie.Builder().UsingYAML(suite.yamlData).Build()

	checker.Assert(yamlSquaddie.ID(), Equals, "squaddie_yaml")
	checker.Assert(yamlSquaddie.Name(), Equals, "YAML squaddie")
	checker.Assert(yamlSquaddie.Affiliation(), Equals, squaddieEntity.Enemy)
}

func (suite *YAMLBuilderSuite) TestDefenseMatchesNewSquaddie(checker *C) {
	yamlSquaddie := squaddie.Builder().UsingYAML(suite.yamlData).Build()

	checker.Assert(yamlSquaddie.MaxHitPoints(), Equals, 2)
	checker.Assert(yamlSquaddie.Dodge(), Equals, 3)
	checker.Assert(yamlSquaddie.Deflect(), Equals, 5)
	checker.Assert(yamlSquaddie.MaxBarrier(), Equals, 7)
	checker.Assert(yamlSquaddie.Armor(), Equals, 9)
}

func (suite *YAMLBuilderSuite) TestOffenseMatchesNewSquaddie(checker *C) {
	yamlSquaddie := squaddie.Builder().UsingYAML(suite.yamlData).Build()

	checker.Assert(yamlSquaddie.Aim(), Equals, 11)
	checker.Assert(yamlSquaddie.Strength(), Equals, 13)
	checker.Assert(yamlSquaddie.Mind(), Equals, 17)
}

func (suite *YAMLBuilderSuite) TestMovementMatchesNewSquaddie(checker *C) {
	yamlSquaddie := squaddie.Builder().UsingYAML(suite.yamlData).Build()

	checker.Assert(yamlSquaddie.MovementDistance(), Equals, 19)
	checker.Assert(yamlSquaddie.MovementType(), Equals, squaddieEntity.Light)
	checker.Assert(yamlSquaddie.MovementCanHitAndRun(), Equals, true)
}

func (suite *YAMLBuilderSuite) TestPowersMatchesNewSquaddie(checker *C) {
	yamlSquaddie := squaddie.Builder().UsingYAML(suite.yamlData).Build()

	powerReferences := yamlSquaddie.GetCopyOfPowerReferences()
	checker.Assert(powerReferences, HasLen, 2)
	checker.Assert(powerReferences[0].PowerID, Equals, "shove_id")
	checker.Assert(powerReferences[0].Name, Equals, "Shove")
	checker.Assert(powerReferences[1].PowerID, Equals, "bandage_id")
	checker.Assert(powerReferences[1].Name, Equals, "Bandage")
}

func (suite *YAMLBuilderSuite) TestClassesMatchesNewSquaddie(checker *C) {
	yamlSquaddie := squaddie.Builder().UsingYAML(suite.yamlData).Build()

	checker.Assert(yamlSquaddie.BaseClassID(), Equals, "baseClassID")
	checker.Assert(yamlSquaddie.CurrentClassID(), Equals, "currentClassID")
	classLevelsConsumed := *yamlSquaddie.ClassLevelsConsumed()
	checker.Assert(classLevelsConsumed["baseClassID"].ClassID, Equals, "baseClassID")
	checker.Assert(classLevelsConsumed["baseClassID"].ClassName, Equals, "Introductory Class")
	checker.Assert(classLevelsConsumed["baseClassID"].LevelsConsumed, HasLen, 2)
	checker.Assert(classLevelsConsumed["baseClassID"].LevelsConsumed[0], Equals, "level0")
	checker.Assert(classLevelsConsumed["baseClassID"].LevelsConsumed[1], Equals, "level1")
	checker.Assert(classLevelsConsumed["currentClassID"].ClassID, Equals, "currentClassID")
	checker.Assert(classLevelsConsumed["currentClassID"].ClassName, Equals, "Advanced Class")
	checker.Assert(classLevelsConsumed["currentClassID"].LevelsConsumed, HasLen, 2)
	checker.Assert(classLevelsConsumed["currentClassID"].LevelsConsumed[0], Equals, "advanced level0")
	checker.Assert(classLevelsConsumed["currentClassID"].LevelsConsumed[1], Equals, "advanced level1")
}

type JSONBuilderSuite struct {
	jsonData []byte
}

var _ = Suite(&JSONBuilderSuite{})

func (suite *JSONBuilderSuite) SetUpTest(checker *C) {
	suite.jsonData = []byte(
		`
{
	"id": "squaddie_json",
	"name": "JSON squaddie",
	"affiliation": "Ally",
	"max_hit_points": 23,
	"dodge": 19,
	"deflect": 17,
	"max_barrier": 13,
	"armor": 11,
	"aim": 7,
	"strength": 5,
	"mind": 3,
	"movement_distance": 2,
	"movement_type": "teleport",
	"hit_and_run": true,
	"powers": [
		{
			"id": "shove_id",
		  	"name": "Shove"
		},
		{
			"id": "bandage_id",
		  	"name": "Bandage"
		}
	],
	"class_progress": [
		{
			"class_id": "baseClassID",
			"class_name": "Introductory Class",
			"is_base_class": true,
			"levels_gained": ["level0", "level1"]
		},
		{
			"class_id": "currentClassID",
			"class_name": "Advanced Class",
			"is_current_class": true,
			"levels_gained": [
			  "advanced level0",
			  "advanced level1"
			]
		}
	]
}
`)
}

func (suite *JSONBuilderSuite) TestIdentificationMatchesNewSquaddie(checker *C) {
	jsonSquaddie := squaddie.Builder().UsingJSON(suite.jsonData).Build()

	checker.Assert(jsonSquaddie.ID(), Equals, "squaddie_json")
	checker.Assert(jsonSquaddie.Name(), Equals, "JSON squaddie")
	checker.Assert(jsonSquaddie.Affiliation(), Equals, squaddieEntity.Ally)
}

func (suite *JSONBuilderSuite) TestDefenseMatchesNewSquaddie(checker *C) {
	jsonSquaddie := squaddie.Builder().UsingJSON(suite.jsonData).Build()

	checker.Assert(jsonSquaddie.MaxHitPoints(), Equals, 23)
	checker.Assert(jsonSquaddie.Dodge(), Equals, 19)
	checker.Assert(jsonSquaddie.Deflect(), Equals, 17)
	checker.Assert(jsonSquaddie.MaxBarrier(), Equals, 13)
	checker.Assert(jsonSquaddie.Armor(), Equals, 11)
}

func (suite *JSONBuilderSuite) TestOffenseMatchesNewSquaddie(checker *C) {
	jsonSquaddie := squaddie.Builder().UsingJSON(suite.jsonData).Build()

	checker.Assert(jsonSquaddie.Aim(), Equals, 7)
	checker.Assert(jsonSquaddie.Strength(), Equals, 5)
	checker.Assert(jsonSquaddie.Mind(), Equals, 3)

}

func (suite *JSONBuilderSuite) TestMovementMatchesNewSquaddie(checker *C) {
	jsonSquaddie := squaddie.Builder().UsingJSON(suite.jsonData).Build()

	checker.Assert(jsonSquaddie.MovementDistance(), Equals, 2)
	checker.Assert(jsonSquaddie.MovementType(), Equals, squaddieEntity.Teleport)
	checker.Assert(jsonSquaddie.MovementCanHitAndRun(), Equals, true)
}

func (suite *JSONBuilderSuite) TestPowersMatchesNewSquaddie(checker *C) {
	jsonSquaddie := squaddie.Builder().UsingYAML(suite.jsonData).Build()

	powerReferences := jsonSquaddie.GetCopyOfPowerReferences()
	checker.Assert(powerReferences, HasLen, 2)
	checker.Assert(powerReferences[0].PowerID, Equals, "shove_id")
	checker.Assert(powerReferences[0].Name, Equals, "Shove")
	checker.Assert(powerReferences[1].PowerID, Equals, "bandage_id")
	checker.Assert(powerReferences[1].Name, Equals, "Bandage")
}

func (suite *JSONBuilderSuite) TestClassesMatchesNewSquaddie(checker *C) {
	jsonSquaddie := squaddie.Builder().UsingJSON(suite.jsonData).Build()

	checker.Assert(jsonSquaddie.BaseClassID(), Equals, "baseClassID")
	checker.Assert(jsonSquaddie.CurrentClassID(), Equals, "currentClassID")
	classLevelsConsumed := *jsonSquaddie.ClassLevelsConsumed()
	checker.Assert(classLevelsConsumed["baseClassID"].ClassID, Equals, "baseClassID")
	checker.Assert(classLevelsConsumed["baseClassID"].ClassName, Equals, "Introductory Class")
	checker.Assert(classLevelsConsumed["baseClassID"].LevelsConsumed, HasLen, 2)
	checker.Assert(classLevelsConsumed["baseClassID"].LevelsConsumed[0], Equals, "level0")
	checker.Assert(classLevelsConsumed["baseClassID"].LevelsConsumed[1], Equals, "level1")
	checker.Assert(classLevelsConsumed["currentClassID"].ClassID, Equals, "currentClassID")
	checker.Assert(classLevelsConsumed["currentClassID"].ClassName, Equals, "Advanced Class")
	checker.Assert(classLevelsConsumed["currentClassID"].LevelsConsumed, HasLen, 2)
	checker.Assert(classLevelsConsumed["currentClassID"].LevelsConsumed[0], Equals, "advanced level0")
	checker.Assert(classLevelsConsumed["currentClassID"].LevelsConsumed[1], Equals, "advanced level1")
}

type BuildCopySuite struct {
	teros *squaddieEntity.Squaddie
}

var _ = Suite(&BuildCopySuite{})

func (suite *BuildCopySuite) SetUpTest(checker *C) {
	suite.teros = squaddie.Builder().Teros().Build()
}

func (suite *BuildCopySuite) TestCopySquaddieIdentification(checker *C) {
	cloneTeros := squaddie.Builder().CloneOf(suite.teros).Build()
	checker.Assert(cloneTeros.HasSameStatsAs(suite.teros), Equals, true)
}

func (suite *BuildCopySuite) TestCopySquaddieDefense(checker *C) {
	defensiveTeros := squaddie.Builder().CloneOf(suite.teros).HitPoints(2).Dodge(3).Deflect(5).Barrier(7).Armor(11).Build()
	cloneTeros := squaddie.Builder().CloneOf(defensiveTeros).Build()
	checker.Assert(cloneTeros.HasSameStatsAs(defensiveTeros), Equals, true)
}

func (suite *BuildCopySuite) TestCopySquaddieOffense(checker *C) {
	offensiveTeros := squaddie.Builder().CloneOf(suite.teros).Aim(2).Strength(3).Mind(5).Build()
	cloneTeros := squaddie.Builder().CloneOf(offensiveTeros).Build()
	checker.Assert(cloneTeros.HasSameStatsAs(offensiveTeros), Equals, true)
}

func (suite *BuildCopySuite) TestCopySquaddieMovement(checker *C) {
	mobileTeros := squaddie.Builder().CloneOf(suite.teros).MovementTeleport().MoveDistance(5).CanHitAndRun().Build()
	cloneTeros := squaddie.Builder().CloneOf(mobileTeros).Build()
	checker.Assert(cloneTeros.HasSameStatsAs(mobileTeros), Equals, true)
}

func (suite *BuildCopySuite) TestCopySquaddiePowers(checker *C) {
	armedTeros := squaddie.Builder().CloneOf(suite.teros).AddPowerByReference(&powerEntity.Reference{
		Name:    "Spear",
		PowerID: "powerIDForSpear",
	}).AddPowerByReference(&powerEntity.Reference{
		Name:    "Blot",
		PowerID: "powerIDForBlot",
	}).Build()
	cloneTeros := squaddie.Builder().CloneOf(armedTeros).Build()
	checker.Assert(cloneTeros.HasSameStatsAs(armedTeros), Equals, true)
}

func (suite *BuildCopySuite) TestCopyClasses(checker *C) {
	experiencedTeros := squaddie.Builder().CloneOf(suite.teros).
		AddClassByReference(&classEntity.ClassReference{ID: "scholarID", Name: "Scholar"}).
		AddClassByReference(&classEntity.ClassReference{ID: "advancedScholarID", Name: "Advanced Scholar"}).
		Build()
	experiencedTeros.SetBaseClassIfNoBaseClass("scholarID")
	experiencedTeros.SetClass("scholarID")
	experiencedTeros.MarkLevelUpBenefitAsConsumed("scholarID", "scholarLevel1")
	cloneTeros := squaddie.Builder().CloneOf(experiencedTeros).Build()
	checker.Assert(cloneTeros.HasSameStatsAs(experiencedTeros), Equals, true)
}