package power_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	powerFactory "github.com/chadius/terosbattleserver/utility/testutility/factory/power"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type PowerBuilder struct {}

var _ = Suite(&PowerBuilder{})

func (suite *PowerBuilder) TestBuildPowerWithName(checker *C) {
	sword := powerFactory.PowerFactory().WithName("Master Sword").Build()
	checker.Assert("Master Sword", Equals, sword.Reference.Name)
}

func (suite *PowerBuilder) TestBuildPowerWithID(checker *C) {
	sword := powerFactory.PowerFactory().WithID("power123").Build()
	checker.Assert("power123", Equals, sword.Reference.ID)
}

func (suite *PowerBuilder) TestBuildPowerTargetsSelf(checker *C) {
	sword := powerFactory.PowerFactory().TargetsSelf().Build()
	checker.Assert(true, Equals, sword.Targeting.TargetSelf)
}

func (suite *PowerBuilder) TestBuildPowerTargetsFriend(checker *C) {
	sword := powerFactory.PowerFactory().TargetsFriend().Build()
	checker.Assert(true, Equals, sword.Targeting.TargetFriend)
}

func (suite *PowerBuilder) TestBuildPowerTargetsFoe(checker *C) {
	sword := powerFactory.PowerFactory().TargetsFoe().Build()
	checker.Assert(true, Equals, sword.Targeting.TargetFoe)
}

func (suite *PowerBuilder) TestBuildPowerIsPhysical(checker *C) {
	sword := powerFactory.PowerFactory().IsPhysical().Build()
	checker.Assert(power.Physical, Equals, sword.PowerType)
}

func (suite *PowerBuilder) TestBuildPowerIsSpell(checker *C) {
	lightning := powerFactory.PowerFactory().IsSpell().Build()
	checker.Assert(power.Spell, Equals, lightning.PowerType)
}

func (suite *PowerBuilder) TestHealingAdjustmentFull(checker *C) {
	bigHeals := powerFactory.PowerFactory().HealingAdjustmentBasedOnUserMindFull().Build()
	checker.Assert(power.Full, Equals, bigHeals.HealingEffect.HealingAdjustmentBasedOnUserMind)
}

func (suite *PowerBuilder) TestHealingAdjustmentHalf(checker *C) {
	someHeals := powerFactory.PowerFactory().HealingAdjustmentBasedOnUserMindHalf().Build()
	checker.Assert(power.Half, Equals, someHeals.HealingEffect.HealingAdjustmentBasedOnUserMind)
}

func (suite *PowerBuilder) TestHealingAdjustmentZero(checker *C) {
	someHeals := powerFactory.PowerFactory().HealingAdjustmentBasedOnUserMindZero().Build()
	checker.Assert(power.Zero, Equals, someHeals.HealingEffect.HealingAdjustmentBasedOnUserMind)
}

func (suite *PowerBuilder) TestHitPointsHealed(checker *C) {
	bigHeals := powerFactory.PowerFactory().HitPointsHealed(5).Build()
	checker.Assert(5, Equals, bigHeals.HealingEffect.HitPointsHealed)
}



func (suite *PowerBuilder) TestBuildAttackEffectToHitBonus(checker *C) {
	damageEffect := powerFactory.PowerFactory().ToHitBonus(2).Build()
	checker.Assert(2, Equals, damageEffect.AttackEffect.ToHitBonus)
}

func (suite *PowerBuilder) TestBuildAttackEffectDamageBonus(checker *C) {
	damageEffect := powerFactory.PowerFactory().DealsDamage(3).Build()
	checker.Assert(3, Equals, damageEffect.AttackEffect.DamageBonus)
}

func (suite *PowerBuilder) TestBuildAttackEffectExtraBarrierBurn(checker *C) {
	damageEffect := powerFactory.PowerFactory().ExtraBarrierBurn(1).Build()
	checker.Assert(1, Equals, damageEffect.AttackEffect.ExtraBarrierBurn)
}

func (suite *PowerBuilder) TestBuildAttackEffectCounterAttackPenaltyReduction(checker *C) {
	damageEffect := powerFactory.PowerFactory().CounterAttackPenaltyReduction(4).Build()
	checker.Assert(4, Equals, damageEffect.AttackEffect.CounterAttackPenaltyReduction)
}

func (suite *PowerBuilder) TestBuildAttackEffectCanBeEquipped(checker *C) {
	sword := powerFactory.PowerFactory().CanBeEquipped().Build()
	checker.Assert(true, Equals, sword.AttackEffect.CanBeEquipped)
}

func (suite *PowerBuilder) TestBuildAttackEffectCannotBeEquipped(checker *C) {
	scroll := powerFactory.PowerFactory().CanBeEquipped().CannotBeEquipped().Build()
	checker.Assert(false, Equals, scroll.AttackEffect.CanBeEquipped)
}

func (suite *PowerBuilder) TestBuildAttackEffectCanCounterAttack(checker *C) {
	sword := powerFactory.PowerFactory().CanCounterAttack().Build()
	checker.Assert(true, Equals, sword.AttackEffect.CanCounterAttack)
}


func (suite *PowerBuilder) TestBuildCriticalEffectDamage(checker *C) {
	criticalDamageEffect := powerFactory.PowerFactory().CriticalDealsDamage(8).Build()
	checker.Assert(8, Equals, criticalDamageEffect.AttackEffect.CriticalEffect.Damage)
}

func (suite *PowerBuilder) TestBuildCriticalEffectThresholdBonus(checker *C) {
	criticalDamageEffect := powerFactory.PowerFactory().CriticalHitThresholdBonus(-2).Build()
	checker.Assert(-2, Equals, criticalDamageEffect.AttackEffect.CriticalEffect.CriticalHitThresholdBonus)
}

type SpecificPowerBuilder struct {}

var _ = Suite(&SpecificPowerBuilder{})

func (suite *SpecificPowerBuilder) TestAxe(checker *C) {
	axe := powerFactory.PowerFactory().Axe().Build()

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
	spear := powerFactory.PowerFactory().Spear().Build()

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
	blot := powerFactory.PowerFactory().Blot().Build()

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
	healingStaff := powerFactory.PowerFactory().HealingStaff().Build()

	checker.Assert("healingStaff", Equals, healingStaff.Name)
	checker.Assert("powerHealingStaff", Equals, healingStaff.ID)
	checker.Assert(true, Equals, healingStaff.Targeting.TargetFriend)
	checker.Assert(power.Spell, Equals, healingStaff.PowerType)
	checker.Assert(3, Equals, healingStaff.HealingEffect.HitPointsHealed)
}
