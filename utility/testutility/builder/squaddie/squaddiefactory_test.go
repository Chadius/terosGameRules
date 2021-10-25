package squaddie_test

import (
	squaddieEntity "github.com/chadius/terosbattleserver/entity/squaddie"
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
	checker.Assert(3, Equals, soldier.Movement.Distance)
}

func (suite *SquaddieMovementBuilder) TestBuildMovementCanHitAndRun(checker *C) {
	runner := squaddie.Builder().CanHitAndRun().Build()
	checker.Assert(true, Equals, runner.Movement.HitAndRun)
}

func (suite *SquaddieMovementBuilder) TestChangeMovementFoot(checker *C) {
	soldier := squaddie.Builder().MovementFoot().Build()
	checker.Assert(squaddieEntity.Foot, Equals, soldier.Movement.Type)
}

func (suite *SquaddieMovementBuilder) TestChangeMovementLight(checker *C) {
	ninja := squaddie.Builder().MovementLight().Build()
	checker.Assert(squaddieEntity.Light, Equals, ninja.Movement.Type)
}

func (suite *SquaddieMovementBuilder) TestChangeMovementFly(checker *C) {
	bird := squaddie.Builder().MovementFly().Build()
	checker.Assert(squaddieEntity.Fly, Equals, bird.Movement.Type)
}

func (suite *SquaddieMovementBuilder) TestChangeMovementTeleport(checker *C) {
	wizard := squaddie.Builder().MovementTeleport().Build()
	checker.Assert(squaddieEntity.Teleport, Equals, wizard.Movement.Type)
}

type SquaddiePowerBuilder struct{}

var _ = Suite(&SquaddiePowerBuilder{})

func (suite *SquaddiePowerBuilder) TestBuildAddPower(checker *C) {
	spear := power.Builder().Spear().Build()
	teros := squaddie.Builder().AddPower(spear).Build()
	checker.Assert(spear.ID, Equals, teros.PowerCollection.PowerReferences[0].ID)
}

type SquaddieClassBuilder struct{}

var _ = Suite(&SquaddieClassBuilder{})

func (suite *SquaddieClassBuilder) TestBuildAddClass(checker *C) {
	mageClass := squaddieclass.ClassBuilder().WithID("A class SquaddieID").WithName("mage").WithInitialBigLevelID("level0").Build()
	teros := squaddie.Builder().AddClass(mageClass).Build()
	checker.Assert(true, Equals, teros.ClassProgress.HasAddedClass(mageClass.ID))
}

func (suite *SquaddieClassBuilder) TestBuildSetClass(checker *C) {
	mageClass := squaddieclass.ClassBuilder().WithID("A class SquaddieID").WithName("mage").WithInitialBigLevelID("level0").Build()
	teros := squaddie.Builder().AddClass(mageClass).SetClass(mageClass).Build()
	checker.Assert(mageClass.ID, Equals, teros.ClassProgress.CurrentClass)
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

type YAMLBuilderSuite struct{}

var _ = Suite(&YAMLBuilderSuite{})

func (suite *YAMLBuilderSuite) TestBuildFromYAML(checker *C) {
	yamlData := []byte(
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
movement_type: Light
hit_and_run: true
`)

	// TODO how to handle adding Squaddie Classes?

	// TODO TestBuildFromYAML could be its own test suite, divided by Squaddie subsection

	yamlSquaddie := squaddie.Builder().UsingYAML(yamlData).Build()

	checker.Assert(yamlSquaddie.ID(), Equals, "squaddie_yaml")
	checker.Assert(yamlSquaddie.Name(), Equals, "YAML squaddie")
	checker.Assert(yamlSquaddie.Affiliation(), Equals, squaddieEntity.Enemy)

	checker.Assert(yamlSquaddie.MaxHitPoints(), Equals, 2)
	checker.Assert(yamlSquaddie.Dodge(), Equals, 3)
	checker.Assert(yamlSquaddie.Deflect(), Equals, 5)
	checker.Assert(yamlSquaddie.MaxBarrier(), Equals, 7)
	checker.Assert(yamlSquaddie.Armor(), Equals, 9)

	checker.Assert(yamlSquaddie.Aim(), Equals, 11)
	checker.Assert(yamlSquaddie.Strength(), Equals, 13)
	checker.Assert(yamlSquaddie.Mind(), Equals, 17)

	//checker.Assert(yamlSquaddie.MovementDistance(), Equals, 19)
	//checker.Assert(yamlSquaddie.MovementType(), Equals, squaddieEntity.Light)
	//checker.Assert(yamlSquaddie.MovementCanHitAndRun(), Equals, true)
}