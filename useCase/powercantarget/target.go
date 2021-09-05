package powercantarget

import (
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
)

// CanTargetTargetAffiliationWithPower sees if the power applies
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
