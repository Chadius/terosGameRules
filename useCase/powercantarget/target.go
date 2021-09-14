package powercantarget

import (
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
)

// InvalidTargetReason explains why the target is invalid
type InvalidTargetReason string

// InvalidTargetReason constants. If a target is invalid it should fall into one of these categories.
const (
	TargetIsValid InvalidTargetReason = "TargetIsValid"
	PowerCannotTargetAffiliation InvalidTargetReason = "PowerCannotTargetAffiliation"
	TargetIsDead InvalidTargetReason = "TargetIsDead"
)

// IsValidTarget checks to see if the user can apply the power against the target.
//   returns a bool and a InvalidTargetReason.
//   If the action is valid, the bool is true and the InvalidTargetReason is TargetIsValid.
func IsValidTarget(userID string, powerID string, targetID string, repos *repositories.RepositoryCollection) (bool, InvalidTargetReason) {
	if !(TargetIsStillAlive(targetID, repos) || UserCanTargetDead())  {
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

	if powerUsed.Targeting.TargetSelf && userID == targetID {
		return true
	}

	areFriendsBecauseAffiliationsAreTheSame := user.Identification.Affiliation != squaddie.Neutral && user.Identification.Affiliation == target.Identification.Affiliation
	areFriendsBecausePlayerAndAlly := (user.Identification.Affiliation == squaddie.Player && target.Identification.Affiliation == squaddie.Ally) || (user.Identification.Affiliation == squaddie.Ally && target.Identification.Affiliation == squaddie.Player)
	if powerUsed.Targeting.TargetFriend && (areFriendsBecauseAffiliationsAreTheSame || areFriendsBecausePlayerAndAlly) {
		return true
	}

	areFoesBecauseNeutral := user.Identification.Affiliation == squaddie.Neutral || target.Identification.Affiliation == squaddie.Neutral
	areFoesBecauseExactlyOneIsEnemy := (user.Identification.Affiliation == squaddie.Enemy && target.Identification.Affiliation != squaddie.Enemy) || (user.Identification.Affiliation != squaddie.Enemy && target.Identification.Affiliation == squaddie.Enemy)
	if powerUsed.Targeting.TargetFoe && (areFoesBecauseNeutral || areFoesBecauseExactlyOneIsEnemy) {
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