package healing

import (
	"github.com/chadius/terosbattleserver/entity/powerinterface"
	"github.com/chadius/terosbattleserver/entity/squaddieinterface"
)

// NoHealing does no healing to the target.
type NoHealing struct{}

// CalculateExpectedHeal always returns 0.
func (f *NoHealing) CalculateExpectedHeal(healer squaddieinterface.Interface, healingPower powerinterface.Interface, target squaddieinterface.Interface) int {
	return 0
}
