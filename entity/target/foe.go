package target

import "github.com/chadius/terosgamerules/entity/squaddieinterface"

// Foe means the squaddie can target Foes.
type Foe struct{}

// Name returns a human-readable name
func (t *Foe) Name() string {
	return "foe"
}

// SquaddieCanTargetOtherSquaddie returns true if the user can target the squaddie.
func (t *Foe) SquaddieCanTargetOtherSquaddie(user squaddieinterface.Interface, target squaddieinterface.Interface) bool {
	return user.AffiliationLogic().IsFoesWith(target.AffiliationLogic())
}
