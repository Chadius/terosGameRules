package healing

import (
	"github.com/chadius/terosbattleserver/entity/powerinterface"
	"github.com/chadius/terosbattleserver/entity/squaddieinterface"
)

// ZeroMindBonus applies none of the healer's Mind as a bonus to the healing amount.
type ZeroMindBonus struct{}

// CalculateExpectedHeal determines how much the healer can heal the target using the given healing Power.
func (f *ZeroMindBonus) CalculateExpectedHeal(healer squaddieinterface.Interface, healingPower powerinterface.Interface, target squaddieinterface.Interface) int {
	maximumHealing := healingPower.HitPointsHealed()
	missingHitPoints := target.MaxHitPoints() - target.CurrentHitPoints()
	if missingHitPoints < maximumHealing {
		return missingHitPoints
	}
	return maximumHealing
}