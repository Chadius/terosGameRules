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
	checker.Assert("Master Sword", Equals, sword.Reference.Name)
}

func (suite *PowerBuilder) TestBuildPowerWithID(checker *C) {
	sword := powerBuilder.Builder().WithID("power123").Build()
	checker.Assert("power123", Equals, sword.Reference.ID)
}

func (suite *PowerBuilder) TestBuildPowerTargetsSelf(checker *C) {
	sword := powerBuilder.Builder().TargetsSelf().Build()
	checker.Assert(true, Equals, sword.Targeting.TargetSelf)
}

func (suite *PowerBuilder) TestBuildPowerTargetsFriend(checker *C) {
	sword := powerBuilder.Builder().TargetsFriend().Build()
	checker.Assert(true, Equals, sword.Targeting.TargetFriend)
}

func (suite *PowerBuilder) TestBuildPowerTargetsFoe(checker *C) {
	sword := powerBuilder.Builder().TargetsFoe().Build()
	checker.Assert(true, Equals, sword.Targeting.TargetFoe)
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
	checker.Assert(power.Full, Equals, bigHeals.HealingEffect.HealingAdjustmentBasedOnUserMind)
}

func (suite *PowerBuilder) TestHealingAdjustmentHalf(checker *C) {
	someHeals := powerBuilder.Builder().HealingAdjustmentBasedOnUserMindHalf().Build()
	checker.Assert(power.Half, Equals, someHeals.HealingEffect.HealingAdjustmentBasedOnUserMind)
}

func (suite *PowerBuilder) TestHealingAdjustmentZero(checker *C) {
	someHeals := powerBuilder.Builder().HealingAdjustmentBasedOnUserMindZero().Build()
	checker.Assert(power.Zero, Equals, someHeals.HealingEffect.HealingAdjustmentBasedOnUserMind)
}

func (suite *PowerBuilder) TestHitPointsHealed(checker *C) {
	bigHeals := powerBuilder.Builder().HitPointsHealed(5).Build()
	checker.Assert(5, Equals, bigHeals.HealingEffect.HitPointsHealed)
}

func (suite *PowerBuilder) TestBuildAttackEffectToHitBonus(checker *C) {
	damageEffect := powerBuilder.Builder().ToHitBonus(2).Build()
	checker.Assert(2, Equals, damageEffect.AttackEffect.ToHitBonus)
}

func (suite *PowerBuilder) TestBuildAttackEffectDamageBonus(checker *C) {
	damageEffect := powerBuilder.Builder().DealsDamage(3).Build()
	checker.Assert(3, Equals, damageEffect.AttackEffect.DamageBonus)
}

func (suite *PowerBuilder) TestBuildAttackEffectExtraBarrierBurn(checker *C) {
	damageEffect := powerBuilder.Builder().ExtraBarrierBurn(1).Build()
	checker.Assert(1, Equals, damageEffect.AttackEffect.ExtraBarrierBurn)
}

func (suite *PowerBuilder) TestBuildAttackEffectCounterAttackPenaltyReduction(checker *C) {
	damageEffect := powerBuilder.Builder().CounterAttackPenaltyReduction(4).Build()
	checker.Assert(4, Equals, damageEffect.AttackEffect.CounterAttackPenaltyReduction)
}

func (suite *PowerBuilder) TestBuildAttackEffectCanBeEquipped(checker *C) {
	sword := powerBuilder.Builder().CanBeEquipped().Build()
	checker.Assert(true, Equals, sword.AttackEffect.CanBeEquipped)
}

func (suite *PowerBuilder) TestBuildAttackEffectCannotBeEquipped(checker *C) {
	scroll := powerBuilder.Builder().CanBeEquipped().CannotBeEquipped().Build()
	checker.Assert(false, Equals, scroll.AttackEffect.CanBeEquipped)
}

func (suite *PowerBuilder) TestBuildAttackEffectCanCounterAttack(checker *C) {
	sword := powerBuilder.Builder().CanCounterAttack().Build()
	checker.Assert(true, Equals, sword.AttackEffect.CanCounterAttack)
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

	checker.Assert("axe", Equals, axe.Name)
	checker.Assert("powerAxe", Equals, axe.ID)
	checker.Assert(true, Equals, axe.Targeting.TargetFoe)
	checker.Assert(power.Physical, Equals, axe.PowerType)
	checker.Assert(true, Equals, axe.AttackEffect.CanBeEquipped)
	checker.Assert(true, Equals, axe.AttackEffect.CanCounterAttack)
	checker.Assert(1, Equals, axe.AttackEffect.DamageBonus)
	checker.Assert(1, Equals, axe.AttackEffect.ToHitBonus)
}

func (suite *SpecificPowerBuilder) TestSpear(checker *C) {
	spear := powerBuilder.Builder().Spear().Build()

	checker.Assert("spear", Equals, spear.Name)
	checker.Assert("powerSpear", Equals, spear.ID)
	checker.Assert(true, Equals, spear.Targeting.TargetFoe)
	checker.Assert(power.Physical, Equals, spear.PowerType)
	checker.Assert(true, Equals, spear.AttackEffect.CanBeEquipped)
	checker.Assert(true, Equals, spear.AttackEffect.CanCounterAttack)
	checker.Assert(1, Equals, spear.AttackEffect.DamageBonus)
	checker.Assert(1, Equals, spear.AttackEffect.ToHitBonus)
}

func (suite *SpecificPowerBuilder) TestBlot(checker *C) {
	blot := powerBuilder.Builder().Blot().Build()

	checker.Assert("blot", Equals, blot.Name)
	checker.Assert("powerBlot", Equals, blot.ID)
	checker.Assert(true, Equals, blot.Targeting.TargetFoe)
	checker.Assert(power.Spell, Equals, blot.PowerType)
	checker.Assert(true, Equals, blot.AttackEffect.CanBeEquipped)
	checker.Assert(false, Equals, blot.AttackEffect.CanCounterAttack)
	checker.Assert(3, Equals, blot.AttackEffect.DamageBonus)
	checker.Assert(0, Equals, blot.AttackEffect.ToHitBonus)
}

func (suite *SpecificPowerBuilder) TestHealingStaff(checker *C) {
	healingStaff := powerBuilder.Builder().HealingStaff().Build()

	checker.Assert("healingStaff", Equals, healingStaff.Name)
	checker.Assert("powerHealingStaff", Equals, healingStaff.ID)
	checker.Assert(true, Equals, healingStaff.Targeting.TargetFriend)
	checker.Assert(power.Spell, Equals, healingStaff.PowerType)
	checker.Assert(3, Equals, healingStaff.HealingEffect.HitPointsHealed)
}
