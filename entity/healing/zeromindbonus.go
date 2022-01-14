package healing

import (
	"github.com/chadius/terosgamerules/entity/squaddieinterface"
)

// ZeroMindBonus applies none of the healer's Mind as a bonus to the healing amount.
type ZeroMindBonus struct{}

// CalculateExpectedHeal determines how much the healer can heal the target using the given healing Power.
func (f *ZeroMindBonus) CalculateExpectedHeal(healer squaddieinterface.Interface, powerHealingAmount int, target squaddieinterface.Interface) int {
	maximumHealing := powerHealingAmount
	missingHitPoints := target.MaxHitPoints() - target.CurrentHitPoints()
	if missingHitPoints < maximumHealing {
		return missingHitPoints
	}
	return maximumHealing
}
