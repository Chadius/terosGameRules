package healing

import (
	"github.com/chadius/terosgamerules/entity/squaddieinterface"
)

// Interface will shape how healing powers work with squaddies.
type Interface interface {
	CalculateExpectedHeal(healer squaddieinterface.Interface, powerHealingAmount int, target squaddieinterface.Interface) int
}
