package target

import (
	"github.com/chadius/terosgamerules/entity/squaddieinterface"
)

// Interface describes how target logic works.
type Interface interface {
	Name() string
	SquaddieCanTargetOtherSquaddie(squaddieinterface.Interface, squaddieinterface.Interface) bool
}
