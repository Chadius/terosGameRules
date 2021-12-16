package levelupbenefit_test

import (
	"github.com/chadius/terosbattleserver/entity/levelupbenefit"
	. "gopkg.in/check.v1"
	"reflect"
)

type LevelUpBuilderSuite struct{}

var _ = Suite(&LevelUpBuilderSuite{})

func (l *LevelUpBuilderSuite) TestBuildDefaultLevelUpBenefit(checker *C) {
	levelUpBenefit, err := levelupbenefit.NewLevelUpBenefitBuilder().Build()

	checker.Assert(err, IsNil)
	checker.Assert(levelUpBenefit.ID(), Not(Equals), "")
	checker.Assert(levelUpBenefit.ClassID(), Not(Equals), "")
	checker.Assert(levelUpBenefit.LevelUpBenefitType(), Equals, levelupbenefit.Small)
}

func (l *LevelUpBuilderSuite) TestBuildWithIdentification(checker *C) {
	levelUpBenefit, err := levelupbenefit.NewLevelUpBenefitBuilder().
		WithID("levelID").
		WithClassID("classID").
		BigLevel().
		Build()

	checker.Assert(err, IsNil)
	checker.Assert(levelUpBenefit.ID(), Equals, "levelID")
	checker.Assert(levelUpBenefit.ClassID(), Equals, "classID")
	checker.Assert(levelUpBenefit.LevelUpBenefitType(), Equals, levelupbenefit.Big)
}

func (l *LevelUpBuilderSuite) TestBuildWithDefense(checker *C) {
	levelUpBenefit, err := levelupbenefit.NewLevelUpBenefitBuilder().
		HitPoints(5).
		Deflect(4).
		Dodge(3).
		Barrier(2).
		Armor(1).
		Build()

	checker.Assert(err, IsNil)
	checker.Assert(levelUpBenefit.MaxHitPoints(), Equals, 5)
	checker.Assert(levelUpBenefit.Deflect(), Equals, 4)
	checker.Assert(levelUpBenefit.Dodge(), Equals, 3)
	checker.Assert(levelUpBenefit.MaxBarrier(), Equals, 2)
	checker.Assert(levelUpBenefit.Armor(), Equals, 1)
}

func (l *LevelUpBuilderSuite) TestBuildWithOffense(checker *C) {
	levelUpBenefit, err := levelupbenefit.NewLevelUpBenefitBuilder().
		Aim(2).
		Strength(3).
		Mind(5).
		Build()

	checker.Assert(err, IsNil)
	checker.Assert(levelUpBenefit.Aim(), Equals, 2)
	checker.Assert(levelUpBenefit.Strength(), Equals, 3)
	checker.Assert(levelUpBenefit.Mind(), Equals, 5)
}

func (l *LevelUpBuilderSuite) TestBuildWithMovement(checker *C) {
	levelUpBenefit, err := levelupbenefit.NewLevelUpBenefitBuilder().
		MovementDistance(2).
		MovementLogic("fly").
		CanHitAndRun().
		Build()

	checker.Assert(err, IsNil)
	checker.Assert(levelUpBenefit.MovementDistance(), Equals, 2)
	checker.Assert(reflect.TypeOf(levelUpBenefit.MovementLogic()).String(), Equals, "*movement.Fly")
	checker.Assert(levelUpBenefit.CanHitAndRun(), Equals, true)
}

func (l *LevelUpBuilderSuite) TestBuildWithDefaultMovement(checker *C) {
	levelUpBenefit, err := levelupbenefit.NewLevelUpBenefitBuilder().Build()

	checker.Assert(err, IsNil)
	checker.Assert(reflect.TypeOf(levelUpBenefit.MovementLogic()).String(), Equals, "*movement.Foot")
}

func (l *LevelUpBuilderSuite) TestBuildWithPowerChanges(checker *C) {
	levelUpBenefit, err := levelupbenefit.NewLevelUpBenefitBuilder().
		GainPower("spearLevel1", "Gold spear").
		GainPower("healingKitLevel0", "Healing kit").
		LosePower("spearLevel0").
		Build()

	checker.Assert(err, IsNil)
	checker.Assert(levelUpBenefit.PowersGained(), HasLen, 2)
	checker.Assert(levelUpBenefit.PowersGained()[0].Name, Equals, "Gold spear")
	checker.Assert(levelUpBenefit.PowersGained()[0].PowerID, Equals, "spearLevel1")
	checker.Assert(levelUpBenefit.PowersGained()[1].Name, Equals, "Healing kit")
	checker.Assert(levelUpBenefit.PowersGained()[1].PowerID, Equals, "healingKitLevel0")

	checker.Assert(levelUpBenefit.PowersLost(), HasLen, 1)
	checker.Assert(levelUpBenefit.PowersLost()[0].PowerID, Equals, "spearLevel0")
}

func (l *LevelUpBuilderSuite) TestBuildWithAlternateMovementTypes(checker *C) {
	onFootLevelUp, _ := levelupbenefit.NewLevelUpBenefitBuilder().
		MovementLogic("Fly").
		FootMovement().
		Build()

	checker.Assert(reflect.TypeOf(onFootLevelUp.MovementLogic()).String(), Equals, "*movement.Foot")

	onLightLevelUp, _ := levelupbenefit.NewLevelUpBenefitBuilder().
		LightMovement().
		Build()

	checker.Assert(reflect.TypeOf(onLightLevelUp.MovementLogic()).String(), Equals, "*movement.Light")

	onFlyLevelUp, _ := levelupbenefit.NewLevelUpBenefitBuilder().
		FlyMovement().
		Build()

	checker.Assert(reflect.TypeOf(onFlyLevelUp.MovementLogic()).String(), Equals, "*movement.Fly")

	onTeleportLevelUp, _ := levelupbenefit.NewLevelUpBenefitBuilder().
		TeleportMovement().
		Build()

	checker.Assert(reflect.TypeOf(onTeleportLevelUp.MovementLogic()).String(), Equals, "*movement.Teleport")
}

type LevelUpBuilderDataSuite struct{}

var _ = Suite(&LevelUpBuilderDataSuite{})

func (l LevelUpBuilderDataSuite) TestUseYAMLToCreateBuilder(checker *C) {
	yamlByteStream := []byte(
		`
id: abcdefg0
class_id: class0
hit_points: 2
dodge: 3
deflect: 5
barrier: 7
armor: 11
aim: 13
strength: 15
mind: 17
powers_gained:
  - name: Scimitar
    id: deadbeef
movement_distance: 19
movement_type: teleport
can_hit_and_run: true
`)
	levelUp, _ := levelupbenefit.NewLevelUpBenefitBuilderFromYAML(yamlByteStream).Build()

	checker.Assert(levelUp.ID(), Equals, "abcdefg0")
	checker.Assert(levelUp.ClassID(), Equals, "class0")
	checker.Assert(levelUp.LevelUpBenefitType(), Equals, levelupbenefit.Small)

	checker.Assert(levelUp.MaxHitPoints(), Equals, 2)
	checker.Assert(levelUp.Dodge(), Equals, 3)
	checker.Assert(levelUp.Deflect(), Equals, 5)
	checker.Assert(levelUp.MaxBarrier(), Equals, 7)
	checker.Assert(levelUp.Armor(), Equals, 11)

	checker.Assert(levelUp.Aim(), Equals, 13)
	checker.Assert(levelUp.Strength(), Equals, 15)
	checker.Assert(levelUp.Mind(), Equals, 17)

	checker.Assert(levelUp.MovementDistance(), Equals, 19)
	checker.Assert(reflect.TypeOf(levelUp.MovementLogic()).String(), Equals, "*movement.Teleport")
	checker.Assert(levelUp.CanHitAndRun(), Equals, true)

	checker.Assert(levelUp.PowersGained(), HasLen, 1)
	checker.Assert(levelUp.PowersGained()[0].Name, Equals, "Scimitar")
	checker.Assert(levelUp.PowersGained()[0].PowerID, Equals, "deadbeef")
}

func (l LevelUpBuilderDataSuite) TestUseJSONToCreateBuilder(checker *C) {
	jsonByteStream := []byte(
		`{
"id": "abcdefg0",
"class_id": "class0",
"is_a_big_level": true,
"hit_points": 2,
"dodge": 3,
"deflect": 5,
"barrier": 7,
"armor": 11,
"aim": 13,
"strength": 15,
"mind": 17,
"powers_lost": ["deadbeef"],
"movement_distance": 19,
"movement_type": "light"
}`)
	levelUp, _ := levelupbenefit.NewLevelUpBenefitBuilderFromJSON(jsonByteStream).Build()

	checker.Assert(levelUp.ID(), Equals, "abcdefg0")
	checker.Assert(levelUp.ClassID(), Equals, "class0")
	checker.Assert(levelUp.LevelUpBenefitType(), Equals, levelupbenefit.Big)

	checker.Assert(levelUp.MaxHitPoints(), Equals, 2)
	checker.Assert(levelUp.Dodge(), Equals, 3)
	checker.Assert(levelUp.Deflect(), Equals, 5)
	checker.Assert(levelUp.MaxBarrier(), Equals, 7)
	checker.Assert(levelUp.Armor(), Equals, 11)

	checker.Assert(levelUp.Aim(), Equals, 13)
	checker.Assert(levelUp.Strength(), Equals, 15)
	checker.Assert(levelUp.Mind(), Equals, 17)

	checker.Assert(levelUp.MovementDistance(), Equals, 19)
	checker.Assert(reflect.TypeOf(levelUp.MovementLogic()).String(), Equals, "*movement.Light")
	checker.Assert(levelUp.CanHitAndRun(), Equals, false)

	checker.Assert(levelUp.PowersGained(), HasLen, 0)

	checker.Assert(levelUp.PowersLost(), HasLen, 1)
	checker.Assert(levelUp.PowersLost()[0].PowerID, Equals, "deadbeef")
}
