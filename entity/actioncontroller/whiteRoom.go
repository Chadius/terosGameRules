package actioncontroller

import (
	"fmt"
	"github.com/cserrant/terosbattleserver/entity/powerusagescenario"
	"github.com/cserrant/terosbattleserver/usecase/powerattackforecast"
	"github.com/cserrant/terosbattleserver/usecase/powercantarget"
	"github.com/cserrant/terosbattleserver/usecase/powercommit"
	"github.com/cserrant/terosbattleserver/usecase/repositories"
	"github.com/cserrant/terosbattleserver/utility"
	"math/rand"
	"strings"
)

// WhiteRoomController assumes all Squaddies are within range and can attack each other.
type WhiteRoomController struct{}

// SetupAction creates a record of the next action.
func (controller *WhiteRoomController) SetupAction(userID string, targetIDs []string, powerID string) *powerusagescenario.Setup {
	powerSetup := &powerusagescenario.Setup{
		UserID:          userID,
		PowerID:         powerID,
		Targets:         targetIDs,
		IsCounterAttack: false,
	}
	return powerSetup
}

// GenerateForecast uses the action to predict results.
func (controller *WhiteRoomController) GenerateForecast(action *powerusagescenario.Setup, repos *repositories.RepositoryCollection) *powerattackforecast.Forecast {
	powerForecast := &powerattackforecast.Forecast{
		Setup: *action,
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: repos.SquaddieRepo,
			PowerRepo:    repos.PowerRepo,
		},
	}
	powerForecast.CalculateForecast()
	return powerForecast
}

// GenerateResult uses the forecast to create results.
func (controller *WhiteRoomController) GenerateResult(
	forecast *powerattackforecast.Forecast,
	repos *repositories.RepositoryCollection,
	useRandomSeed bool,
	randomSeed int64) *powercommit.Result {

	powerResult := &powercommit.Result{
		Forecast:  forecast,
		DieRoller: &utility.RandomDieRoller{},
	}
	if useRandomSeed == true {
		rand.Seed(randomSeed)
	}

	powerResult.Commit()
	return powerResult
}

//InvalidAttackDescription gives more detail on why an attack is invalid.
type InvalidAttackDescription struct {
	Reason      powercantarget.InvalidTargetReason
	Description []string
}

// CheckForValidAction makes sure the action is valid. Otherwise, it returns an error.
func (controller *WhiteRoomController) CheckForValidAction(action *powerusagescenario.Setup, repos *repositories.RepositoryCollection) []InvalidAttackDescription {
	descriptions := []InvalidAttackDescription{}

	for _, targetID := range action.Targets {

		_, reasonForInvalidTarget := powercantarget.IsValidTarget(
			action.UserID,
			action.PowerID,
			targetID,
			repos,
		)

		user := repos.SquaddieRepo.GetSquaddieByID(action.UserID)
		powerUsed := repos.PowerRepo.GetPowerByID(action.PowerID)
		target := repos.SquaddieRepo.GetSquaddieByID(targetID)

		if reasonForInvalidTarget == powercantarget.UserIsDead {
			descriptions = append(
				descriptions,
				InvalidAttackDescription{
					reasonForInvalidTarget,
					[]string{
						"User is dead, cannot use power",
						fmt.Sprintf("  %s[%s] is dead", user.Identification.Name, user.Identification.ID),
					},
				},
			)
			continue
		}

		if reasonForInvalidTarget == powercantarget.TargetIsDead {
			descriptions = append(
				descriptions,
				InvalidAttackDescription{
					reasonForInvalidTarget,
					[]string{
						"Target is dead, cannot use power",
						fmt.Sprintf("  %s[%s] is dead", target.Identification.Name, target.Identification.ID),
					},
				},
			)
			continue
		}

		if reasonForInvalidTarget == powercantarget.PowerCannotTargetAffiliation {
			affiliationRelationsTargeted := []string{}
			if powerUsed.Targeting.TargetSelf {
				affiliationRelationsTargeted = append(affiliationRelationsTargeted, "self")
			}
			if powerUsed.Targeting.TargetFriend {
				affiliationRelationsTargeted = append(affiliationRelationsTargeted, "friend")
			}
			if powerUsed.Targeting.TargetFoe {
				affiliationRelationsTargeted = append(affiliationRelationsTargeted, "foe")
			}

			descriptions = append(
				descriptions,
				InvalidAttackDescription{
					reasonForInvalidTarget,
					[]string{
						"Target is not compatible with affiliation",
						fmt.Sprintf("  %s[%s] is a %s", user.Identification.Name, user.Identification.ID, user.Identification.Affiliation),
						fmt.Sprintf("    uses %s[%s] that targets %s", powerUsed.Name, powerUsed.ID, strings.Join(affiliationRelationsTargeted, ",")),
						fmt.Sprintf("  %s[%s] is a %s", target.Identification.Name, target.Identification.ID, target.Identification.Affiliation),
					},
				},
			)
			continue
		}
	}
	return descriptions
}
