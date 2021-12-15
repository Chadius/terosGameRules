package powercantarget

import (
	"github.com/chadius/terosbattleserver/usecase/repositories"
)

// ValidTargetStrategy describes the shape of classes that check for valid attacks.
type ValidTargetStrategy interface {
	IsValidTarget(userID string, powerID string, targetID string, repos *repositories.RepositoryCollection) (bool, InvalidTargetReason)
	CanTargetTargetAffiliationWithPower(userID string, powerID string, targetID string, repos *repositories.RepositoryCollection) bool
}

// ValidTargetChecker applies business logic to figure out if the user squaddie can target another squaddie with a given power.
type ValidTargetChecker struct{}

// IsValidTarget checks to see if the user can apply the power against the target.
//   returns a bool and a InvalidTargetReason.
//   If the action is valid, the bool is true and the InvalidTargetReason is TargetIsValid.
func (v *ValidTargetChecker) IsValidTarget(userID string, powerID string, targetID string, repos *repositories.RepositoryCollection) (bool, InvalidTargetReason) {
	if !(v.targetIsStillAlive(userID, repos)) {
		return false, UserIsDead
	}

	if !(v.targetIsStillAlive(targetID, repos) || v.userCanTargetDead()) {
		return false, TargetIsDead
	}

	if v.CanTargetTargetAffiliationWithPower(userID, powerID, targetID, repos) == false {
		return false, PowerCannotTargetAffiliation
	}
	return true, TargetIsValid
}

// CanTargetTargetAffiliationWithPower sees if the power can be used on the target because of the affiliation.
//    Returns true if so, false otherwise.
func (v *ValidTargetChecker) CanTargetTargetAffiliationWithPower(userID string, powerID string, targetID string, repos *repositories.RepositoryCollection) bool {
	user := repos.SquaddieRepo.GetSquaddieByID(userID)
	target := repos.SquaddieRepo.GetSquaddieByID(targetID)
	powerUsed := repos.PowerRepo.GetPowerByID(powerID)

	if powerUsed.CanPowerTargetSelf() && userID == targetID {
		return true
	}

	// TODO Demeter violation - AffiliationLogic should be able to work on a squaddie, not affiliation logic
	areFriends := user.AffiliationLogic().IsFriendsWith(target.AffiliationLogic())
	if powerUsed.CanPowerTargetFriend() && areFriends {
		return true
	}

	areFoes := user.AffiliationLogic().IsFoesWith(target.AffiliationLogic())
	if powerUsed.CanPowerTargetFoe() && areFoes {
		return true
	}

	return false
}

// userCanTargetDead returns true if the target is dead and the power can target dead.
func (v *ValidTargetChecker) userCanTargetDead() bool {
	return false
}

// targetIsStillAlive returns true if the target is alive.
func (v *ValidTargetChecker) targetIsStillAlive(targetID string, repos *repositories.RepositoryCollection) bool {
	target := repos.SquaddieRepo.GetSquaddieByID(targetID)
	return !target.Defense.IsDead()
}

// InvalidTargetReason explains why the target is invalid
type InvalidTargetReason string

// InvalidTargetReason constants. If a target is invalid it should fall into one of these categories.
const (
	TargetIsValid                InvalidTargetReason = "TargetIsValid"
	PowerCannotTargetAffiliation InvalidTargetReason = "PowerCannotTargetAffiliation"
	TargetIsDead                 InvalidTargetReason = "TargetIsDead"
	UserIsDead                   InvalidTargetReason = "UserIsDead"
)
