package healing

import (
	"github.com/chadius/terosgamerules/entity/squaddieinterface"
)

// HalfMindBonus applies half of the healer's Mind as a bonus to the healing amount.
type HalfMindBonus struct{}

// CalculateExpectedHeal determines how much the healer can heal the target using the given healing Power.
func (f *HalfMindBonus) CalculateExpectedHeal(healer squaddieinterface.Interface, powerHealingAmount int, target squaddieinterface.Interface) int {
	squaddieMindBonus := healer.Mind() / 2
	maximumHealing := powerHealingAmount + squaddieMindBonus
	missingHitPoints := target.MaxHitPoints() - target.CurrentHitPoints()
	if missingHitPoints < maximumHealing {
		return missingHitPoints
	}
	return maximumHealing
}
