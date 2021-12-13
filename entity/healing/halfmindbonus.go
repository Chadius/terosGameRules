package healing

import (
	"github.com/chadius/terosbattleserver/entity/powerinterface"
	"github.com/chadius/terosbattleserver/entity/squaddieinterface"
)

type HalfMindBonus struct{}

// CalculateExpectedHeal determines how much the healer can heal the target using the given healing Power.
func (f *HalfMindBonus) CalculateExpectedHeal(healer squaddieinterface.Interface, healingPower powerinterface.Interface, target squaddieinterface.Interface) int {
	squaddieMindBonus := healer.Mind() / 2
	maximumHealing := healingPower.HitPointsHealed() + squaddieMindBonus
	missingHitPoints := target.MaxHitPoints() - target.CurrentHitPoints()
	if missingHitPoints < maximumHealing {
		return missingHitPoints
	}
	return maximumHealing
}
