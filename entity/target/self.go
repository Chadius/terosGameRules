package target

import "github.com/chadius/terosgamerules/entity/squaddieinterface"

// Self means the squaddie can target itself.
type Self struct{}

// Name returns a human-readable name
func (t *Self) Name() string {
	return "self"
}

// SquaddieCanTargetOtherSquaddie returns true if the user can target the squaddie.
func (t *Self) SquaddieCanTargetOtherSquaddie(user squaddieinterface.Interface, target squaddieinterface.Interface) bool {
	return user.ID() == target.ID()
}
