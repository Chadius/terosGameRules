package power_test

import (
	"github.com/cserrant/terosbattleserver/entity/power"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type PowerChanceCheckSuite struct{}

var _ = Suite(&PowerChanceCheckSuite{})

func (suite *PowerChanceCheckSuite) TestPowerCanCrit(checker *C) {
	staticPower := &power.Power{
		Reference: power.Reference{
			Name: "Static",
			ID:   "power0",
		},
		PowerType: power.Physical,
		AttackEffect: &power.AttackingEffect{
			ToHitBonus:                    0,
			DamageBonus:                   0,
			CanCounterAttack:              false,
			CounterAttackPenaltyReduction: 0,
			CriticalEffect:                nil,
		},
	}
	checker.Assert(staticPower.AttackEffect.CanCriticallyHit(), Equals, false)

	criticalPower := &power.Power{
		Reference: power.Reference{
			Name: "Critical",
			ID:   "power1",
		},
		PowerType: power.Physical,
		AttackEffect: &power.AttackingEffect{
			ToHitBonus:                    0,
			DamageBonus:                   0,
			CanCounterAttack:              false,
			CounterAttackPenaltyReduction: 0,
			CriticalEffect: &power.CriticalEffect{
				CriticalHitThresholdBonus: 0,
				Damage:                    1,
			},
		},
	}
	checker.Assert(criticalPower.AttackEffect.CanCriticallyHit(), Equals, true)
}
