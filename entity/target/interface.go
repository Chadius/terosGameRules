package target

import (
	"github.com/chadius/terosbattleserver/entity/squaddieinterface"
)

// Interface describes how target logic works.
type Interface interface {
	Name() string
	SquaddieCanTargetOtherSquaddie(squaddieinterface.Interface, squaddieinterface.Interface) bool
}
