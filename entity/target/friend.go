package target

import "github.com/chadius/terosgamerules/entity/squaddieinterface"

// Friend means the squaddie can target friends.
type Friend struct{}

// Name returns a human-readable name
func (t *Friend) Name() string {
	return "friend"
}

// SquaddieCanTargetOtherSquaddie returns true if the user can target the squaddie.
func (t *Friend) SquaddieCanTargetOtherSquaddie(user squaddieinterface.Interface, target squaddieinterface.Interface) bool {
	return user.AffiliationLogic().IsFriendsWith(target.AffiliationLogic())
}
