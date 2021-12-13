package healing

import (
	"github.com/chadius/terosbattleserver/entity/powerinterface"
	"github.com/chadius/terosbattleserver/entity/squaddieinterface"
)

// FullMindBonus applies the healer's Mind as a bonus to the healing amount.
type FullMindBonus struct{}

// CalculateExpectedHeal determines how much the healer can heal the target using the given healing Power.
func (f *FullMindBonus) CalculateExpectedHeal(healer squaddieinterface.Interface, healingPower powerinterface.Interface, target squaddieinterface.Interface) int {
	squaddieMindBonus := healer.Mind()
	maximumHealing := healingPower.HitPointsHealed() + squaddieMindBonus
	missingHitPoints := target.MaxHitPoints() - target.CurrentHitPoints()
	if missingHitPoints < maximumHealing {
		return missingHitPoints
	}
	return maximumHealing
}
