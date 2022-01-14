package healing

import (
	"github.com/chadius/terosgamerules/entity/squaddieinterface"
)

// NoHealing does no healing to the target.
type NoHealing struct{}

// CalculateExpectedHeal always returns 0.
func (f *NoHealing) CalculateExpectedHeal(healer squaddieinterface.Interface, powerHealingAmount int, target squaddieinterface.Interface) int {
	return 0
}
