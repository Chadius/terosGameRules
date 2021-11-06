package power_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	powerBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	. "gopkg.in/check.v1"
)

type CriticalEffectActivation struct{}

var _ = Suite(&CriticalEffectActivation{})

func (suite *CriticalEffectActivation) SetUpTest(checker *C) {}

func (suite *CriticalEffectActivation) TestCanSetCriticalHitThreshold(checker *C) {
	defaultCriticalHitEffect := &power.CriticalEffect{}

	checker.Assert(defaultCriticalHitEffect.CriticalHitThreshold(), Equals, 6)

	moreLikelyToCrit := powerBuilder.CriticalEffectBuilder().CriticalHitThresholdBonus(2).Build()

	checker.Assert(moreLikelyToCrit.CriticalHitThreshold(), Equals, 4)
}

func (suite *CriticalEffectActivation) TestExtraCriticalHitDamage(checker *C) {
	extraCriticalDamage := powerBuilder.CriticalEffectBuilder().DealsDamage(6).Build()
	checker.Assert(extraCriticalDamage.ExtraCriticalHitDamage(), Equals, 6)
}

func (suite *CriticalEffectActivation) TestCriticalEffectsWithJSON(checker *C) {
	jsonByteStream := []byte(`[{
					"name": "Scimitar",
					"id": "deadbeef",
					"damage_bonus": 2,
					"power_type": "physical",
					"can_attack": true,
					"can_critical": true,
					"can_counter_attack": true,
					"counter_attack_penalty": -3,
					"critical_hit_threshold_bonus": 5,
					"critical_damage": 7
				}]`)
	newRepo := powerrepository.NewPowerRepository()
	success, _ := newRepo.AddJSONSource(jsonByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(newRepo.GetNumberOfPowers(), Equals, 1)

	scimitar := newRepo.GetPowerByID("deadbeef")
	checker.Assert(scimitar.CriticalHitThresholdBonus(), Equals, 5)
	checker.Assert(scimitar.ExtraCriticalHitDamage(), Equals, 7)
}

func (suite *CriticalEffectActivation) TestCriticalEffectsWithYAML(checker *C) {
	yamlByteStream := []byte(`-
  name: Scimitar
  id: deadbeef
  power_type: physical
  can_attack: true
  damage_bonus: 7
  can_counter_attack: true
  counter_attack_penalty: -5
  can_critical: true
  critical_hit_threshold_bonus: 3
  critical_damage: 2
`)
	newRepo := powerrepository.NewPowerRepository()
	success, _ := newRepo.AddYAMLSource(yamlByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(newRepo.GetNumberOfPowers(), Equals, 1)

	scimitar := newRepo.GetPowerByID("deadbeef")
	checker.Assert(scimitar.CriticalHitThresholdBonus(), Equals, 3)
	checker.Assert(scimitar.ExtraCriticalHitDamage(), Equals, 2)
}
