package actioncontroller

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagescenario"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerattackforecast"
	"github.com/cserrant/terosBattleServer/usecase/powercantarget"
	"github.com/cserrant/terosBattleServer/usecase/powercommit"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
	"github.com/cserrant/terosBattleServer/utility"
	"strings"
)

// WhiteRoomController assumes all Squaddies are within range and can attack each other.
type WhiteRoomController struct {}

// SetupAction creates a record of the next action.
func (controller *WhiteRoomController) SetupAction(attacker, target *squaddie.Squaddie, power *power.Power) *powerusagescenario.Setup {
	powerSetup := &powerusagescenario.Setup{
		UserID:          attacker.Identification.ID,
		PowerID:         power.ID,
		Targets:         []string{target.Identification.ID},
		IsCounterAttack: false,
	}
	return powerSetup
}

// GenerateForecast uses the action to predict results.
func (controller *WhiteRoomController) GenerateForecast(action *powerusagescenario.Setup, repos *repositories.RepositoryCollection) *powerattackforecast.Forecast {
	powerForecast := &powerattackforecast.Forecast{
		Setup: *action,
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo:    repos.SquaddieRepo,
			PowerRepo:       repos.PowerRepo,
		},
	}
	powerForecast.CalculateForecast()
	return powerForecast
}

// GenerateResult uses the forecast to create results.
func (controller *WhiteRoomController) GenerateResult(forecast *powerattackforecast.Forecast, repos *repositories.RepositoryCollection) *powercommit.Result {
	powerResult := &powercommit.Result{
		Forecast: forecast,
		DieRoller: &utility.RandomDieRoller{},
	}
	powerResult.Commit()
	return powerResult
}

//InvalidAttackDescription gives more detail on why an attack is invalid.
type InvalidAttackDescription struct {
	Reason powercantarget.InvalidTargetReason
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