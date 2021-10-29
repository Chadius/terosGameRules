package power_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	powerBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type PowerBuilder struct{}

var _ = Suite(&PowerBuilder{})

func (suite *PowerBuilder) TestBuildPowerWithName(checker *C) {
	sword := powerBuilder.Builder().WithName("Master Sword").Build()
	checker.Assert("Master Sword", Equals, sword.Name())
}

func (suite *PowerBuilder) TestBuildPowerWithID(checker *C) {
	sword := powerBuilder.Builder().WithID("power123").Build()
	checker.Assert("power123", Equals, sword.ID())
}

func (suite *PowerBuilder) TestBuildPowerTargetsSelf(checker *C) {
	sword := powerBuilder.Builder().TargetsSelf().Build()
	checker.Assert(true, Equals, sword.CanPowerTargetSelf())
}

func (suite *PowerBuilder) TestBuildPowerTargetsFriend(checker *C) {
	sword := powerBuilder.Builder().TargetsFriend().Build()
	checker.Assert(true, Equals, sword.CanPowerTargetFriend())
}

func (suite *PowerBuilder) TestBuildPowerTargetsFoe(checker *C) {
	sword := powerBuilder.Builder().TargetsFoe().Build()
	checker.Assert(true, Equals, sword.CanPowerTargetFoe())
}

func (suite *PowerBuilder) TestBuildPowerIsPhysical(checker *C) {
	sword := powerBuilder.Builder().IsPhysical().Build()
	checker.Assert(power.Physical, Equals, sword.PowerType)
}

func (suite *PowerBuilder) TestBuildPowerIsSpell(checker *C) {
	lightning := powerBuilder.Builder().IsSpell().Build()
	checker.Assert(power.Spell, Equals, lightning.PowerType)
}

func (suite *PowerBuilder) TestHealingAdjustmentFull(checker *C) {
	bigHeals := powerBuilder.Builder().HealingAdjustmentBasedOnUserMindFull().Build()
	checker.Assert(power.Full, Equals, bigHeals.HealingEffect.HealingHealingAdjustmentBasedOnUserMind)
}

func (suite *PowerBuilder) TestHealingAdjustmentHalf(checker *C) {
	someHeals := powerBuilder.Builder().HealingAdjustmentBasedOnUserMindHalf().Build()
	checker.Assert(power.Half, Equals, someHeals.HealingEffect.HealingHealingAdjustmentBasedOnUserMind)
}

func (suite *PowerBuilder) TestHealingAdjustmentZero(checker *C) {
	someHeals := powerBuilder.Builder().HealingAdjustmentBasedOnUserMindZero().Build()
	checker.Assert(power.Zero, Equals, someHeals.HealingEffect.HealingHealingAdjustmentBasedOnUserMind)
}

func (suite *PowerBuilder) TestHitPointsHealed(checker *C) {
	bigHeals := powerBuilder.Builder().HitPointsHealed(5).Build()
	checker.Assert(5, Equals, bigHeals.HealingEffect.HealingHitPointsHealed)
}

func (suite *PowerBuilder) TestBuildAttackEffectToHitBonus(checker *C) {
	damageEffect := powerBuilder.Builder().ToHitBonus(2).Build()
	checker.Assert(2, Equals, damageEffect.AttackEffect.AttackToHitBonus)
}

func (suite *PowerBuilder) TestBuildAttackEffectDamageBonus(checker *C) {
	damageEffect := powerBuilder.Builder().DealsDamage(3).Build()
	checker.Assert(3, Equals, damageEffect.AttackEffect.AttackDamageBonus)
}

func (suite *PowerBuilder) TestBuildAttackEffectExtraBarrierBurn(checker *C) {
	damageEffect := powerBuilder.Builder().ExtraBarrierBurn(1).Build()
	checker.Assert(1, Equals, damageEffect.AttackEffect.AttackExtraBarrierBurn)
}

func (suite *PowerBuilder) TestBuildAttackEffectCounterAttackPenaltyReduction(checker *C) {
	damageEffect := powerBuilder.Builder().CounterAttackPenaltyReduction(4).Build()
	checker.Assert(4, Equals, damageEffect.AttackEffect.AttackCounterAttackPenaltyReduction)
}

func (suite *PowerBuilder) TestBuildAttackEffectCanBeEquipped(checker *C) {
	sword := powerBuilder.Builder().CanBeEquipped().Build()
	checker.Assert(true, Equals, sword.AttackEffect.AttackCanBeEquipped)
}

func (suite *PowerBuilder) TestBuildAttackEffectCannotBeEquipped(checker *C) {
	scroll := powerBuilder.Builder().CanBeEquipped().CannotBeEquipped().Build()
	checker.Assert(false, Equals, scroll.AttackEffect.AttackCanBeEquipped)
}

func (suite *PowerBuilder) TestBuildAttackEffectCanCounterAttack(checker *C) {
	sword := powerBuilder.Builder().CanCounterAttack().Build()
	checker.Assert(true, Equals, sword.AttackEffect.AttackCanCounterAttack)
}

func (suite *PowerBuilder) TestBuildCriticalEffectDamage(checker *C) {
	criticalDamageEffect := powerBuilder.Builder().CriticalDealsDamage(8).Build()
	checker.Assert(8, Equals, criticalDamageEffect.AttackEffect.CriticalEffect.Damage)
}

func (suite *PowerBuilder) TestBuildCriticalEffectThresholdBonus(checker *C) {
	criticalDamageEffect := powerBuilder.Builder().CriticalHitThresholdBonus(-2).Build()
	checker.Assert(-2, Equals, criticalDamageEffect.AttackEffect.CriticalEffect.CriticalHitThresholdBonus)
}

type SpecificPowerBuilder struct{}

var _ = Suite(&SpecificPowerBuilder{})

func (suite *SpecificPowerBuilder) TestAxe(checker *C) {
	axe := powerBuilder.Builder().Axe().Build()

	checker.Assert("axe", Equals, axe.Name())
	checker.Assert("powerAxe", Equals, axe.ID())
	checker.Assert(true, Equals, axe.CanPowerTargetFoe())
	checker.Assert(power.Physical, Equals, axe.PowerType)
	checker.Assert(true, Equals, axe.AttackEffect.AttackCanBeEquipped)
	checker.Assert(true, Equals, axe.AttackEffect.AttackCanCounterAttack)
	checker.Assert(1, Equals, axe.AttackEffect.AttackDamageBonus)
	checker.Assert(1, Equals, axe.AttackEffect.AttackToHitBonus)
}

func (suite *SpecificPowerBuilder) TestSpear(checker *C) {
	spear := powerBuilder.Builder().Spear().Build()

	checker.Assert("spear", Equals, spear.Name())
	checker.Assert("powerSpear", Equals, spear.ID())
	checker.Assert(true, Equals, spear.CanPowerTargetFoe())
	checker.Assert(power.Physical, Equals, spear.PowerType)
	checker.Assert(true, Equals, spear.AttackEffect.AttackCanBeEquipped)
	checker.Assert(true, Equals, spear.AttackEffect.AttackCanCounterAttack)
	checker.Assert(1, Equals, spear.AttackEffect.AttackDamageBonus)
	checker.Assert(1, Equals, spear.AttackEffect.AttackToHitBonus)
}

func (suite *SpecificPowerBuilder) TestBlot(checker *C) {
	blot := powerBuilder.Builder().Blot().Build()

	checker.Assert("blot", Equals, blot.Name())
	checker.Assert("powerBlot", Equals, blot.ID())
	checker.Assert(true, Equals, blot.CanPowerTargetFoe())
	checker.Assert(power.Spell, Equals, blot.PowerType)
	checker.Assert(true, Equals, blot.AttackEffect.AttackCanBeEquipped)
	checker.Assert(false, Equals, blot.AttackEffect.AttackCanCounterAttack)
	checker.Assert(3, Equals, blot.AttackEffect.AttackDamageBonus)
	checker.Assert(0, Equals, blot.AttackEffect.AttackToHitBonus)
}

func (suite *SpecificPowerBuilder) TestHealingStaff(checker *C) {
	healingStaff := powerBuilder.Builder().HealingStaff().Build()

	checker.Assert("healingStaff", Equals, healingStaff.Name())
	checker.Assert("powerHealingStaff", Equals, healingStaff.ID())
	checker.Assert(true, Equals, healingStaff.CanPowerTargetFriend())
	checker.Assert(power.Spell, Equals, healingStaff.PowerType)
	checker.Assert(3, Equals, healingStaff.HealingEffect.HealingHitPointsHealed)
}

type YAMLBuilderSuite struct{
	yamlData []byte
}

var _ = Suite(&YAMLBuilderSuite{})

func (suite *YAMLBuilderSuite) SetUpTest(checker *C) {
	suite.yamlData = []byte(
		`
id: power_id
name: Power Name
power_type: spell
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
	yamlPower := powerBuilder.Builder().UsingYAML(suite.yamlData).Build()

	checker.Assert(yamlPower.ID(), Equals, "power_id")
	checker.Assert(yamlPower.Name(), Equals, "Power Name")
	checker.Assert(yamlPower.Type(), Equals, power.Spell)
}

func (suite *YAMLBuilderSuite) TestTargetingMatchesNewPower(checker *C) {
	yamlPower := powerBuilder.Builder().UsingYAML(suite.yamlData).Build()
	checker.Assert(yamlPower.CanPowerTargetSelf(), Equals, true)
	checker.Assert(yamlPower.CanPowerTargetFoe(), Equals, true)
	checker.Assert(yamlPower.CanPowerTargetFriend(), Equals, false)
}

func (suite *YAMLBuilderSuite) TestAttackEffectMatchesNewPower(checker *C) {
	yamlPower := powerBuilder.Builder().UsingYAML(suite.yamlData).Build()

	checker.Assert(yamlPower.ToHitBonus(), Equals, 2)
	checker.Assert(yamlPower.DamageBonus(), Equals, 3)
	checker.Assert(yamlPower.ExtraBarrierBurn(), Equals, 5)
	checker.Assert(yamlPower.CanBeEquipped(), Equals, true)
	checker.Assert(yamlPower.CanCounterAttack(), Equals, true)
	checker.Assert(yamlPower.CounterAttackPenaltyReduction(), Equals, 7)
}

func (suite *YAMLBuilderSuite) TestCriticalEffectMatchesNewPower(checker *C) {
	yamlPower := powerBuilder.Builder().UsingYAML(suite.yamlData).Build()

	checker.Assert(yamlPower.CriticalHitThreshold(), Equals, power.CriticalHitThresholdInitialValue - 9)
	checker.Assert(yamlPower.ExtraCriticalHitDamage(), Equals, 11)
}

type JSONBuilderSuite struct{
	jsonData []byte
}

var _ = Suite(&JSONBuilderSuite{})

func (suite *JSONBuilderSuite) SetUpTest(checker *C) {
	suite.jsonData = []byte(
		`
{
   "id": "power_id",
   "name": "Power Name",
   "power_type": "physical",
   "can_heal": true,
   "healing_adjustment_based_on_user_mind": "half",
   "hit_points_healed": 2
}
`)
}

func (suite *JSONBuilderSuite) TestIdentificationMatchesNewPower(checker *C) {
	jsonPower := powerBuilder.Builder().UsingJSON(suite.jsonData).Build()

	checker.Assert(jsonPower.ID(), Equals, "power_id")
	checker.Assert(jsonPower.Name(), Equals, "Power Name")
	checker.Assert(jsonPower.Type(), Equals, power.Physical)
}

func (suite *JSONBuilderSuite) TestHealingMatchesNewPower(checker *C) {
	jsonPower := powerBuilder.Builder().UsingJSON(suite.jsonData).Build()

	checker.Assert(jsonPower.HealingAdjustmentBasedOnUserMind(), Equals, power.Half)
	checker.Assert(jsonPower.HitPointsHealed(), Equals, 2)
}
