package power_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	. "gopkg.in/check.v1"
	"reflect"
)

type PowerBuilder struct{}

var _ = Suite(&PowerBuilder{})

func (suite *PowerBuilder) TestBuildPowerWithName(checker *C) {
	sword := power.NewPowerBuilder().WithName("Master Sword").Build()
	checker.Assert("Master Sword", Equals, sword.Name())
}

func (suite *PowerBuilder) TestBuildPowerWithID(checker *C) {
	sword := power.NewPowerBuilder().WithID("power123").Build()
	checker.Assert("power123", Equals, sword.ID())
}

func (suite *PowerBuilder) TestBuildPowerTargetsSelf(checker *C) {
	sword := power.NewPowerBuilder().TargetsSelf().Build()
	checker.Assert(true, Equals, sword.CanPowerTargetSelf())
}

func (suite *PowerBuilder) TestBuildPowerTargetsFriend(checker *C) {
	sword := power.NewPowerBuilder().TargetsFriend().Build()
	checker.Assert(true, Equals, sword.CanPowerTargetFriend())
}

func (suite *PowerBuilder) TestBuildPowerTargetsFoe(checker *C) {
	sword := power.NewPowerBuilder().TargetsFoe().Build()
	checker.Assert(true, Equals, sword.CanPowerTargetFoe())
}

func (suite *PowerBuilder) TestBuildPowerIsPhysical(checker *C) {
	sword := power.NewPowerBuilder().IsPhysical().Build()
	checker.Assert(reflect.TypeOf(sword.PowerSourceLogic()).String(), Equals, "*powersource.Physical")
}

func (suite *PowerBuilder) TestBuildPowerIsSpell(checker *C) {
	lightning := power.NewPowerBuilder().IsSpell().Build()
	checker.Assert(reflect.TypeOf(lightning.PowerSourceLogic()).String(), Equals, "*powersource.Spell")
}

func (suite *PowerBuilder) TestHealingAdjustmentFull(checker *C) {
	bigHeals := power.NewPowerBuilder().HealingAdjustmentBasedOnUserMindFull().Build()
	checker.Assert(reflect.TypeOf(bigHeals.HealingLogic()).String(), Equals, "*healing.FullMindBonus")
}

func (suite *PowerBuilder) TestHealingAdjustmentHalf(checker *C) {
	someHeals := power.NewPowerBuilder().HealingAdjustmentBasedOnUserMindHalf().Build()
	checker.Assert(reflect.TypeOf(someHeals.HealingLogic()).String(), Equals, "*healing.HalfMindBonus")
}

func (suite *PowerBuilder) TestHealingAdjustmentZero(checker *C) {
	noHeals := power.NewPowerBuilder().HealingAdjustmentBasedOnUserMindZero().Build()
	checker.Assert(reflect.TypeOf(noHeals.HealingLogic()).String(), Equals, "*healing.ZeroMindBonus")
}

func (suite *PowerBuilder) TestHitPointsHealed(checker *C) {
	bigHeals := power.NewPowerBuilder().HitPointsHealed(5).HealingAdjustmentBasedOnUserMindFull().Build()
	checker.Assert(5, Equals, bigHeals.HitPointsHealed())
}

func (suite *PowerBuilder) TestBuildAttackEffectToHitBonus(checker *C) {
	damageEffect := power.NewPowerBuilder().ToHitBonus(2).Build()
	checker.Assert(2, Equals, damageEffect.ToHitBonus())
}

func (suite *PowerBuilder) TestBuildAttackEffectDamageBonus(checker *C) {
	damageEffect := power.NewPowerBuilder().DealsDamage(3).Build()
	checker.Assert(3, Equals, damageEffect.DamageBonus())
}

func (suite *PowerBuilder) TestBuildAttackEffectExtraBarrierBurn(checker *C) {
	damageEffect := power.NewPowerBuilder().ExtraBarrierBurn(1).Build()
	checker.Assert(1, Equals, damageEffect.ExtraBarrierBurn())
}

func (suite *PowerBuilder) TestBuildAttackEffectCounterAttackPenaltyReduction(checker *C) {
	damageEffect := power.NewPowerBuilder().CounterAttackPenaltyReduction(4).Build()
	checker.Assert(4, Equals, damageEffect.CounterAttackPenaltyReduction())
}

func (suite *PowerBuilder) TestBuildAttackEffectCanBeEquipped(checker *C) {
	sword := power.NewPowerBuilder().CanBeEquipped().Build()
	checker.Assert(true, Equals, sword.CanBeEquipped())
}

func (suite *PowerBuilder) TestBuildAttackEffectCannotBeEquipped(checker *C) {
	scroll := power.NewPowerBuilder().CanBeEquipped().CannotBeEquipped().Build()
	checker.Assert(false, Equals, scroll.CanBeEquipped())
}

func (suite *PowerBuilder) TestBuildAttackEffectCanCounterAttack(checker *C) {
	sword := power.NewPowerBuilder().CanCounterAttack().Build()
	checker.Assert(true, Equals, sword.CanCounterAttack())
}

func (suite *PowerBuilder) TestBuildCriticalEffectDamage(checker *C) {
	criticalDamageEffect := power.NewPowerBuilder().CriticalDealsDamage(8).Build()
	checker.Assert(8, Equals, criticalDamageEffect.ExtraCriticalHitDamage())
}

func (suite *PowerBuilder) TestBuildCriticalEffectThresholdBonus(checker *C) {
	criticalDamageEffect := power.NewPowerBuilder().CriticalHitThresholdBonus(-2).Build()
	checker.Assert(-2, Equals, criticalDamageEffect.CriticalHitThresholdBonus())
}

type SpecificPowerBuilder struct{}

var _ = Suite(&SpecificPowerBuilder{})

func (suite *SpecificPowerBuilder) TestAxe(checker *C) {
	axe := power.NewPowerBuilder().Axe().Build()

	checker.Assert("axe", Equals, axe.Name())
	checker.Assert("powerAxe", Equals, axe.ID())
	checker.Assert(true, Equals, axe.CanPowerTargetFoe())
	checker.Assert(reflect.TypeOf(axe.PowerSourceLogic()).String(), Equals, "*powersource.Physical")
	checker.Assert(true, Equals, axe.CanBeEquipped())
	checker.Assert(true, Equals, axe.CanCounterAttack())
	checker.Assert(1, Equals, axe.DamageBonus())
	checker.Assert(1, Equals, axe.ToHitBonus())
}

func (suite *SpecificPowerBuilder) TestSpear(checker *C) {
	spear := power.NewPowerBuilder().Spear().Build()

	checker.Assert("spear", Equals, spear.Name())
	checker.Assert("powerSpear", Equals, spear.ID())
	checker.Assert(true, Equals, spear.CanPowerTargetFoe())
	checker.Assert(reflect.TypeOf(spear.PowerSourceLogic()).String(), Equals, "*powersource.Physical")
	checker.Assert(true, Equals, spear.CanBeEquipped())
	checker.Assert(true, Equals, spear.CanCounterAttack())
	checker.Assert(1, Equals, spear.DamageBonus())
	checker.Assert(1, Equals, spear.ToHitBonus())
}

func (suite *SpecificPowerBuilder) TestBlot(checker *C) {
	blot := power.NewPowerBuilder().Blot().Build()

	checker.Assert("blot", Equals, blot.Name())
	checker.Assert("powerBlot", Equals, blot.ID())
	checker.Assert(true, Equals, blot.CanPowerTargetFoe())
	checker.Assert(reflect.TypeOf(blot.PowerSourceLogic()).String(), Equals, "*powersource.Spell")
	checker.Assert(true, Equals, blot.CanBeEquipped())
	checker.Assert(false, Equals, blot.CanCounterAttack())
	checker.Assert(3, Equals, blot.DamageBonus())
	checker.Assert(0, Equals, blot.ToHitBonus())
}

func (suite *SpecificPowerBuilder) TestHealingStaff(checker *C) {
	healingStaff := power.NewPowerBuilder().HealingStaff().Build()

	checker.Assert("healingStaff", Equals, healingStaff.Name())
	checker.Assert("powerHealingStaff", Equals, healingStaff.ID())
	checker.Assert(true, Equals, healingStaff.CanPowerTargetFriend())
	checker.Assert(reflect.TypeOf(healingStaff.PowerSourceLogic()).String(), Equals, "*powersource.Spell")
	checker.Assert(3, Equals, healingStaff.HitPointsHealed())
}

type YAMLBuilderSuite struct {
	yamlData []byte
}

var _ = Suite(&YAMLBuilderSuite{})

func (suite *YAMLBuilderSuite) SetUpTest(checker *C) {
	suite.yamlData = []byte(
		`
id: power_id
name: Power name
source: spell
target_self: true
target_foe: true
can_attack: true
to_hit_bonus: 2
damage_bonus: 3
extra_barrier_damage: 5 
can_be_equipped: true
can_counter_attack: true
counter_attack_penalty_reduction: 7
can_critical: true
critical_hit_threshold_bonus: 9
critical_damage: 11
`)
}

func (suite *YAMLBuilderSuite) TestIdentificationMatchesNewPower(checker *C) {
	yamlPower := power.NewPowerBuilder().UsingYAML(suite.yamlData).Build()

	checker.Assert(yamlPower.ID(), Equals, "power_id")
	checker.Assert(yamlPower.Name(), Equals, "Power name")
	checker.Assert(reflect.TypeOf(yamlPower.PowerSourceLogic()).String(), Equals, "*powersource.Spell")
}

func (suite *YAMLBuilderSuite) TestTargetingMatchesNewPower(checker *C) {
	yamlPower := power.NewPowerBuilder().UsingYAML(suite.yamlData).Build()
	checker.Assert(yamlPower.CanPowerTargetSelf(), Equals, true)
	checker.Assert(yamlPower.CanPowerTargetFoe(), Equals, true)
	checker.Assert(yamlPower.CanPowerTargetFriend(), Equals, false)
}

func (suite *YAMLBuilderSuite) TestAttackEffectMatchesNewPower(checker *C) {
	yamlPower := power.NewPowerBuilder().UsingYAML(suite.yamlData).Build()

	checker.Assert(yamlPower.ToHitBonus(), Equals, 2)
	checker.Assert(yamlPower.DamageBonus(), Equals, 3)
	checker.Assert(yamlPower.ExtraBarrierBurn(), Equals, 5)
	checker.Assert(yamlPower.CanBeEquipped(), Equals, true)
	checker.Assert(yamlPower.CanCounterAttack(), Equals, true)
	checker.Assert(yamlPower.CounterAttackPenaltyReduction(), Equals, 7)
}

func (suite *YAMLBuilderSuite) TestCriticalEffectMatchesNewPower(checker *C) {
	yamlPower := power.NewPowerBuilder().UsingYAML(suite.yamlData).Build()

	checker.Assert(yamlPower.CriticalHitThreshold(), Equals, power.CriticalHitThresholdInitialValue-9)
	checker.Assert(yamlPower.ExtraCriticalHitDamage(), Equals, 11)
}

func (suite *YAMLBuilderSuite) TestPowersThatCannotHealHaveNoHealingLogic(checker *C) {
	yamlPower := power.NewPowerBuilder().UsingYAML(suite.yamlData).Build()
	checker.Assert(reflect.TypeOf(yamlPower.HealingLogic()).String(), Equals, "*healing.NoHealing")
	checker.Assert(yamlPower.CanHeal(), Equals, false)
}

type JSONBuilderSuite struct {
	jsonData []byte
}

var _ = Suite(&JSONBuilderSuite{})

func (suite *JSONBuilderSuite) SetUpTest(checker *C) {
	suite.jsonData = []byte(
		`
{
   "id": "power_id",
   "name": "Power name",
   "source": "physical",
   "can_heal": true,
   "healing_logic": "half",
   "hit_points_healed": 2
}
`)
}

func (suite *JSONBuilderSuite) TestIdentificationMatchesNewPower(checker *C) {
	jsonPower := power.NewPowerBuilder().UsingJSON(suite.jsonData).Build()

	checker.Assert(jsonPower.ID(), Equals, "power_id")
	checker.Assert(jsonPower.Name(), Equals, "Power name")
	checker.Assert(reflect.TypeOf(jsonPower.PowerSourceLogic()).String(), Equals, "*powersource.Physical")
}

func (suite *JSONBuilderSuite) TestHealingMatchesNewPower(checker *C) {
	jsonPower := power.NewPowerBuilder().UsingJSON(suite.jsonData).Build()

	checker.Assert(reflect.TypeOf(jsonPower.HealingLogic()).String(), Equals, "*healing.HalfMindBonus")
	checker.Assert(jsonPower.HitPointsHealed(), Equals, 2)
}

type BuildCopySuite struct {
	spear        *power.Power
	healingStaff *power.Power
}

var _ = Suite(&BuildCopySuite{})

func (suite *BuildCopySuite) SetUpTest(checker *C) {
	suite.spear = power.NewPowerBuilder().Spear().Build()
	suite.healingStaff = power.NewPowerBuilder().HealingStaff().Build()
}

func (suite *BuildCopySuite) TestCopyAttackPower(checker *C) {
	copySpear := power.NewPowerBuilder().CloneOf(suite.spear).Build()
	checker.Assert(copySpear.HasSameStatsAs(suite.spear), Equals, true)
}

func (suite *BuildCopySuite) TestCopyHealingPower(checker *C) {
	copyHealingStaff := power.NewPowerBuilder().CloneOf(suite.healingStaff).Build()
	checker.Assert(copyHealingStaff.HasSameStatsAs(suite.healingStaff), Equals, true)
}

func (suite *BuildCopySuite) TestCopyCriticalAttackPower(checker *C) {
	criticalSpear := power.NewPowerBuilder().CloneOf(suite.spear).CriticalDealsDamage(10).CriticalHitThresholdBonus(2).Build()
	copyCriticalSpear := power.NewPowerBuilder().CloneOf(criticalSpear).Build()
	checker.Assert(copyCriticalSpear.HasSameStatsAs(criticalSpear), Equals, true)
}
