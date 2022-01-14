package powersource

import "github.com/chadius/terosgamerules/entity/squaddieinterface"

// Interface shapes the power type and how it is calculated
type Interface interface {
	Name() string
	ToHitPenalty(s squaddieinterface.Interface) int
	ArmorResistance(s squaddieinterface.Interface) int
	BarrierResistance(s squaddieinterface.Interface) int
	RawDamage(s squaddieinterface.Interface) int
}
