package powercantarget

import (
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/repositories"
)

// InvalidTargetReason explains why the target is invalid
type InvalidTargetReason string

// InvalidTargetReason constants. If a target is invalid it should fall into one of these categories.
const (
	TargetIsValid                InvalidTargetReason = "TargetIsValid"
	PowerCannotTargetAffiliation InvalidTargetReason = "PowerCannotTargetAffiliation"
	TargetIsDead                 InvalidTargetReason = "TargetIsDead"
	UserIsDead                   InvalidTargetReason = "UserIsDead"
)

// IsValidTarget checks to see if the user can apply the power against the target.
//   returns a bool and a InvalidTargetReason.
//   If the action is valid, the bool is true and the InvalidTargetReason is TargetIsValid.
func IsValidTarget(userID string, powerID string, targetID string, repos *repositories.RepositoryCollection) (bool, InvalidTargetReason) {
	if !(TargetIsStillAlive(userID, repos)) {
		return false, UserIsDead
	}

	if !(TargetIsStillAlive(targetID, repos) || UserCanTargetDead()) {
		return false, TargetIsDead
	}

	if CanTargetTargetAffiliationWithPower(userID, powerID, targetID, repos) == false {
		return false, PowerCannotTargetAffiliation
	}
	return true, TargetIsValid
}

// CanTargetTargetAffiliationWithPower sees if the power can be used on the target because of the affiliation.
//    Returns true if so, false otherwise.
func CanTargetTargetAffiliationWithPower(userID string, powerID string, targetID string, repos *repositories.RepositoryCollection) bool {
	user := repos.SquaddieRepo.GetSquaddieByID(userID)
	target := repos.SquaddieRepo.GetSquaddieByID(targetID)
	powerUsed := repos.PowerRepo.GetPowerByID(powerID)

	if powerUsed.CanPowerTargetSelf() && userID == targetID {
		return true
	}

	areFriendsBecauseAffiliationsAreTheSame := user.Affiliation() != squaddie.Neutral && user.Affiliation() == target.Affiliation()
	areFriendsBecausePlayerAndAlly := (user.Affiliation() == squaddie.Player && target.Affiliation() == squaddie.Ally) || (user.Affiliation() == squaddie.Ally && target.Affiliation() == squaddie.Player)
	if powerUsed.CanPowerTargetFriend() && (areFriendsBecauseAffiliationsAreTheSame || areFriendsBecausePlayerAndAlly) {
		return true
	}

	areFoesBecauseNeutral := user.Affiliation() == squaddie.Neutral || target.Affiliation() == squaddie.Neutral
	areFoesBecauseExactlyOneIsEnemy := (user.Affiliation() == squaddie.Enemy && target.Affiliation() != squaddie.Enemy) || (user.Affiliation() != squaddie.Enemy && target.Affiliation() == squaddie.Enemy)
	if powerUsed.CanPowerTargetFoe() && (areFoesBecauseNeutral || areFoesBecauseExactlyOneIsEnemy) {
		return true
	}

	return false
}

// TargetIsStillAlive returns true if the target is alive.
func TargetIsStillAlive(targetID string, repos *repositories.RepositoryCollection) bool {
	target := repos.SquaddieRepo.GetSquaddieByID(targetID)
	return !target.Defense.IsDead()
}

// UserCanTargetDead returns true if the target is dead and the power can target dead.
func UserCanTargetDead() bool {
	return false
}
