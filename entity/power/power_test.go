package power_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	powerBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type PowerChanceCheckSuite struct{}

var _ = Suite(&PowerChanceCheckSuite{})

func (suite *PowerChanceCheckSuite) TestPowerCanCrit(checker *C) {
	staticPower := &power.Power{
		Reference: power.Reference{
			Name:    "Static",
			PowerID: "power0",
		},
		PowerType: power.Physical,
		AttackEffect: &power.AttackingEffect{
			AttackToHitBonus:                    0,
			AttackDamageBonus:                   0,
			AttackCanCounterAttack:              false,
			AttackCounterAttackPenaltyReduction: 0,
			CriticalEffect:                      nil,
		},
	}
	checker.Assert(staticPower.CanCriticallyHit(), Equals, false)

	criticalPower := &power.Power{
		Reference: power.Reference{
			Name:    "Critical",
			PowerID: "power1",
		},
		PowerType: power.Physical,
		AttackEffect: &power.AttackingEffect{
			AttackToHitBonus:                    0,
			AttackDamageBonus:                   0,
			AttackCanCounterAttack:              false,
			AttackCounterAttackPenaltyReduction: 0,
			CriticalEffect: powerBuilder.CriticalEffectBuilder().CriticalHitThresholdBonus(0).DealsDamage(1).Build(),
		},
	}
	checker.Assert(criticalPower.CanCriticallyHit(), Equals, true)
}
