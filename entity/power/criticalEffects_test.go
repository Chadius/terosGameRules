package power_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	. "gopkg.in/check.v1"
)

type CriticalEffectActivation struct{}

var _ = Suite(&CriticalEffectActivation{})

func (suite *CriticalEffectActivation) SetUpTest(checker *C) {}

func (suite *CriticalEffectActivation) TestCanSetCriticalHitThreshold(checker *C) {
	defaultCriticalHitEffect := &power.CriticalEffect{}

	checker.Assert(defaultCriticalHitEffect.CriticalHitThreshold(), Equals, 6)

	moreLikelyToCrit := &power.CriticalEffect{
		CriticalHitThresholdBonus: 2,
	}

	checker.Assert(moreLikelyToCrit.CriticalHitThreshold(), Equals, 4)
}

func (suite *CriticalEffectActivation) TestExtraCriticalHitDamage(checker *C) {
	extraCriticalDamage := &power.CriticalEffect{
		Damage: 6,
	}

	checker.Assert(extraCriticalDamage.ExtraCriticalHitDamage(), Equals, 6)
}

func (suite *CriticalEffectActivation) TestCriticalEffectsWithJSON(checker *C) {
	jsonByteStream := []byte(`[{
					"name": "Scimitar",
					"id": "deadbeef",
					"damage_bonus": 2,
					"power_type": "Physical",
					"attack_effect": {
						"can_counter_attack": true,
						"counter_attack_penalty": -2,
						"critical_effect": {
							"critical_hit_threshold_bonus": 2,
							"damage": 3
						}
					}
				}]`)
	newRepo := power.NewPowerRepository()
	success, _ := newRepo.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(newRepo.GetNumberOfPowers(), Equals, 1)

	scimitar := newRepo.GetPowerByID("deadbeef")
	checker.Assert(scimitar.AttackEffect.CriticalEffect.CriticalHitThresholdBonus, Equals, 2)
	checker.Assert(scimitar.AttackEffect.CriticalEffect.Damage, Equals, 3)
}

func (suite *CriticalEffectActivation) TestCriticalEffectsWithYAML(checker *C) {
	yamlByteStream := []byte(`-
  name: Scimitar
  id: deadbeef
  power_type: Physical
  attack_effect:
    damage_bonus: 2
    can_counter_attack: true
    counter_attack_penalty: -2
    critical_effect:
      critical_hit_threshold_bonus: 2
      damage: 3
`)
	newRepo := power.NewPowerRepository()
	success, _ := newRepo.AddYAMLSource(yamlByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(newRepo.GetNumberOfPowers(), Equals, 1)

	scimitar := newRepo.GetPowerByID("deadbeef")
	checker.Assert(scimitar.AttackEffect.CriticalEffect.CriticalHitThresholdBonus, Equals, 2)
	checker.Assert(scimitar.AttackEffect.CriticalEffect.Damage, Equals, 3)
}
