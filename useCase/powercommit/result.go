package powercommit

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/chadius/terosbattleserver/entity/damagedistribution"
	"github.com/chadius/terosbattleserver/entity/powerusagescenario"
	"github.com/chadius/terosbattleserver/usecase/powerattackforecast"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/usecase/squaddiestats"
	"github.com/chadius/terosbattleserver/utility"
)

//counterfeiter:generate . ResultStrategy
// ResultStrategy describes the shape of Result objects.
type ResultStrategy interface {
	Forecast() *powerattackforecast.Forecast
	DieRoller() utility.SixSideGenerator
	ResultPerTarget() []*ResultPerTarget
	Commit()
}

// Result applies the forecast given to determine what actually happened. Changes are committed.
type Result struct {
	forecast        *powerattackforecast.Forecast
	dieRoller       utility.SixSideGenerator
	resultPerTarget []*ResultPerTarget
}

// NewResult returns a new Result object.
func NewResult(forecast *powerattackforecast.Forecast, dieRoller utility.SixSideGenerator, resultsPerTarget []*ResultPerTarget) *Result {
	return &Result{forecast: forecast, dieRoller: dieRoller, resultPerTarget: resultsPerTarget}
}

// CopyResultWithNewDieRoller copies this Result, overriding the fields with the given die Roller.
func (result *Result) CopyResultWithNewDieRoller(newDieRoller utility.SixSideGenerator) *Result {
	return &Result{forecast: result.forecast, dieRoller: newDieRoller, resultPerTarget: []*ResultPerTarget{}}
}

// Forecast is a getter.
func (result *Result) Forecast() *powerattackforecast.Forecast {
	return result.forecast
}

// DieRoller is a getter.
func (result *Result) DieRoller() utility.SixSideGenerator {
	return result.dieRoller
}

// ResultPerTarget is a getter.
func (result *Result) ResultPerTarget() []*ResultPerTarget {
	return result.resultPerTarget
}

// Commit tries to use the power and records the effects.
func (result *Result) Commit() {
	for _, calculation := range result.forecast.ForecastedResultPerTarget() {
		attackResultForTarget := result.getAttackResult(calculation)
		if attackResultForTarget != nil {
			result.resultPerTarget = append(result.resultPerTarget, attackResultForTarget)
		}

		healResultForTarget := result.getHealingResult(calculation)
		if healResultForTarget != nil {
			result.resultPerTarget = append(result.resultPerTarget, healResultForTarget)
		}
	}
	for _, calculation := range result.forecast.ForecastedResultPerTarget() {
		if result.isCounterAttackPossible(calculation) {
			counterAttackResultForTarget := result.calculateAttackResultForThisTarget(calculation.CounterAttackSetup(), calculation.CounterAttack(), result.forecast.Repositories())
			result.resultPerTarget = append(result.resultPerTarget, counterAttackResultForTarget)
		}
	}
}

func (result *Result) getAttackResult(calculation powerattackforecast.CalculationInterface) *ResultPerTarget {
	if calculation.Attack() == nil {
		return nil
	}
	return result.calculateAttackResultForThisTarget(calculation.Setup(), calculation.Attack(), result.forecast.Repositories())
}

func (result *Result) getHealingResult(calculation powerattackforecast.CalculationInterface) *ResultPerTarget {
	if calculation.HealingForecast() == nil {
		return nil
	}
	return result.calculateHealingResultForThisTarget(calculation.Setup(), calculation.HealingForecast(), result.forecast.Repositories())
}

func (result *Result) isCounterAttackPossible(calculation powerattackforecast.CalculationInterface) bool {
	if calculation.CounterAttack() == nil {
		return false
	}

	counterAttacker := calculation.Repositories().SquaddieRepo.GetOriginalSquaddieByID(calculation.CounterAttackSetup().UserID)
	if counterAttacker.IsDead() {
		return false
	}

	return true
}

func (result *Result) calculateAttackResultForThisTarget(setup *powerusagescenario.Setup, attack *powerattackforecast.AttackForecast, repositories *repositories.RepositoryCollection) *ResultPerTarget {
	results := &ResultPerTarget{
		userID:   setup.UserID,
		targetID: setup.Targets[0],
		powerID:  setup.PowerID,
		attack: &AttackResult{
			isCounterAttack: attack.AttackerContext.IsCounterAttack(),
		},
	}

	attackingSquaddie := repositories.SquaddieRepo.GetOriginalSquaddieByID(setup.UserID)
	checkEquip := powerequip.CheckRepositories{}
	checkEquip.SquaddieEquipPower(attackingSquaddie, setup.PowerID, repositories)

	attackRoll, defendRoll := result.dieRoller.RollTwoDice()
	results.attack.attackerToHitBonus = attack.VersusContext.ToHit().AttackerToHitBonus
	results.attack.defenderToHitPenalty = attack.VersusContext.ToHit().DefenderToHitPenalty

	results.attack.attackRoll = attackRoll
	results.attack.defendRoll = defendRoll

	results.attack.attackerTotal = results.attack.attackRoll + results.attack.attackerToHitBonus
	results.attack.defenderTotal = results.attack.defendRoll + results.attack.defenderToHitPenalty

	results.attack.hitTarget = results.attack.attackerTotal >= results.attack.defenderTotal
	results.attack.criticallyHitTarget = attack.AttackerContext.CanCritical() && results.attack.attackerTotal >= results.attack.defenderTotal+attack.AttackerContext.CriticalHitThreshold()

	if !results.attack.hitTarget {
		results.attack.damage = &damagedistribution.DamageDistribution{
			DamageAbsorbedByArmor:   0,
			DamageAbsorbedByBarrier: 0,
			RawDamageDealt:          0,
			ExtraBarrierBurnt:       0,
			TotalRawBarrierBurnt:    0,
		}
	} else if results.attack.criticallyHitTarget {
		results.attack.damage = attack.VersusContext.CriticalHitDamage()
	} else {
		results.attack.damage = attack.VersusContext.NormalDamage()
	}

	targetSquaddie := repositories.SquaddieRepo.GetOriginalSquaddieByID(results.targetID)
	targetSquaddie.TakeDamageDistribution(results.attack.damage)

	return results
}

func (result *Result) calculateHealingResultForThisTarget(setup *powerusagescenario.Setup, forecast *powerattackforecast.HealingForecast, repositories *repositories.RepositoryCollection) *ResultPerTarget {
	resultForThisTarget := &ResultPerTarget{
		userID:   setup.UserID,
		targetID: setup.Targets[0],
		powerID:  setup.PowerID,
		healing: &HealResult{
			hitPointsRestored: 0,
		},
	}

	healingSquaddie := repositories.SquaddieRepo.GetOriginalSquaddieByID(setup.UserID)
	checkEquip := powerequip.CheckRepositories{}
	checkEquip.SquaddieEquipPower(healingSquaddie, setup.PowerID, repositories)

	offenseStrategy := squaddiestats.CalculateSquaddieOffenseStats{}

	targetSquaddie := repositories.SquaddieRepo.GetOriginalSquaddieByID(resultForThisTarget.targetID)
	maximumHealing, err := offenseStrategy.GetHitPointsHealedWithPower(setup.UserID, setup.PowerID, resultForThisTarget.targetID, repositories)
	if err != nil {
		return resultForThisTarget
	}
	hitPointsRestored := targetSquaddie.GainHitPoints(maximumHealing)
	resultForThisTarget.healing.hitPointsRestored = hitPointsRestored
	return resultForThisTarget
}
