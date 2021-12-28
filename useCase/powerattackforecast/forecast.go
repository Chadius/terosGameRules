package powerattackforecast

import (
	"github.com/chadius/terosbattleserver/entity/powerusagescenario"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/usecase/squaddiestats"
)

// ForecastBuilder is used to make new Forecast objects.
type ForecastBuilder struct {
	setup                     *powerusagescenario.Setup
	repositories              *repositories.RepositoryCollection
	forecastedResultPerTarget []Calculation
	offenseStrategy           squaddiestats.CalculateSquaddieOffenseStatsStrategy
}

// NewForecastBuilder returns a new object.
func NewForecastBuilder() *ForecastBuilder {
	return &ForecastBuilder{
		setup:                     nil,
		repositories:              nil,
		forecastedResultPerTarget: nil,
		offenseStrategy:           nil,
	}
}

// Setup sets the field.
func (fb *ForecastBuilder) Setup(setup *powerusagescenario.Setup) *ForecastBuilder {
	fb.setup = setup
	return fb
}

// Repositories sets the field.
func (fb *ForecastBuilder) Repositories(repositories *repositories.RepositoryCollection) *ForecastBuilder {
	fb.repositories = repositories
	return fb
}

// ForecastedResultPerTarget sets the field.
func (fb *ForecastBuilder) ForecastedResultPerTarget(forecastedResultPerTarget []Calculation) *ForecastBuilder {
	fb.forecastedResultPerTarget = forecastedResultPerTarget
	return fb
}

// OffenseStrategy sets the field.
func (fb *ForecastBuilder) OffenseStrategy(offenseStrategy squaddiestats.CalculateSquaddieOffenseStatsStrategy) *ForecastBuilder {
	fb.offenseStrategy = offenseStrategy
	return fb
}

// Build returns the builder object
func (fb *ForecastBuilder) Build() *Forecast {
	return NewForecast(
		fb.setup,
		fb.repositories,
		fb.forecastedResultPerTarget,
		fb.offenseStrategy,
	)
}

// Forecast will store the information needed to explain what will happen when a squaddie
//  uses a given power. It can be asked multiple questions.
type Forecast struct {
	setup                     powerusagescenario.Setup
	repositories              *repositories.RepositoryCollection
	forecastedResultPerTarget []Calculation
	offenseStrategy           squaddiestats.CalculateSquaddieOffenseStatsStrategy
}

// NewForecast returns a new Forecast object.
func NewForecast(
	setup *powerusagescenario.Setup,
	collection *repositories.RepositoryCollection,
	forecastedResults []Calculation,
	offenseStrategy squaddiestats.CalculateSquaddieOffenseStatsStrategy,
) *Forecast {
	return &Forecast{
		setup:                     *setup,
		repositories:              collection,
		forecastedResultPerTarget: forecastedResults,
		offenseStrategy:           offenseStrategy,
	}
}

// Repositories gets the object
func (forecast *Forecast) Repositories() *repositories.RepositoryCollection {
	return forecast.repositories
}

// ForecastedResultPerTarget gets the object
func (forecast *Forecast) ForecastedResultPerTarget() []Calculation {
	return forecast.forecastedResultPerTarget
}

// CalculationInterface describes what all calculations will subscribe to
type CalculationInterface interface {
	Setup() *powerusagescenario.Setup
	Attack() *AttackForecast
	CounterAttackSetup() *powerusagescenario.Setup
	CounterAttack() *AttackForecast
	Repositories() *repositories.RepositoryCollection
	HealingForecast() *HealingForecast
}

// Calculation holds the results of the forecast.
type Calculation struct {
	repositories *repositories.RepositoryCollection
	setup        *powerusagescenario.Setup

	attack              *AttackForecast
	counterAttackSetup2 *powerusagescenario.Setup
	counterAttack       *AttackForecast

	healingForecast *HealingForecast
}

// Attack is a getter.
func (c *Calculation) Attack() *AttackForecast {
	return c.attack
}

// CounterAttack is a getter.
func (c *Calculation) CounterAttack() *AttackForecast {
	return c.counterAttack
}

// CounterAttackSetup is a getter.
func (c *Calculation) CounterAttackSetup() *powerusagescenario.Setup {
	return c.counterAttackSetup2
}

// Setup is a getter.
func (c *Calculation) Setup() *powerusagescenario.Setup {
	return c.setup
}

// Repositories is a getter.
func (c *Calculation) Repositories() *repositories.RepositoryCollection {
	return c.repositories
}

//HealingForecast is a getter.
func (c *Calculation) HealingForecast() *HealingForecast {
	return c.healingForecast
}

// AttackForecast shows what will happen if the power used is offensive.
type AttackForecast struct {
	AttackerContext AttackerContext
	DefenderContext DefenderContext
	VersusContext   VersusContext
}

// CalculateForecast gives a numerical prediction of the power's effect.
func (forecast *Forecast) CalculateForecast() {
	powerToUse := forecast.repositories.PowerRepo.GetPowerByID(forecast.setup.PowerID)

	for _, targetID := range forecast.setup.Targets {
		calculation := Calculation{
			setup: &powerusagescenario.Setup{
				UserID:          forecast.setup.UserID,
				PowerID:         forecast.setup.PowerID,
				Targets:         []string{targetID},
				IsCounterAttack: false,
			},
			repositories: &repositories.RepositoryCollection{
				SquaddieRepo: forecast.repositories.SquaddieRepo,
				PowerRepo:    forecast.repositories.PowerRepo,
			},
		}

		if powerToUse.CanAttack() {
			forecast.addAttackAndCounterAttackToCalculation(targetID, &calculation)
		}
		if powerToUse.CanHeal() {
			forecast.addHealingEffectToCalculation(targetID, &calculation)
		}

		forecast.forecastedResultPerTarget = append(forecast.forecastedResultPerTarget, calculation)
	}
}

func (forecast *Forecast) addAttackAndCounterAttackToCalculation(targetID string, calculation *Calculation) {
	attack := forecast.CalculateAttackForecast(targetID)
	var counterAttack *AttackForecast
	var counterAttackSetup *powerusagescenario.Setup
	if forecast.IsCounterattackPossible(targetID, forecast.repositories) {
		counterAttackSetup, counterAttack = forecast.createCounterAttackForecast(targetID)
	}

	calculation.attack = attack
	calculation.counterAttackSetup2 = counterAttackSetup
	calculation.counterAttack = counterAttack
}

func (forecast *Forecast) addHealingEffectToCalculation(targetID string, calculation *Calculation) {
	calculation.healingForecast = forecast.CalculateHealingForecast(targetID)
}

// IsCounterattackPossible returns true if the squaddie with the targetID can currently counterattack.
func (forecast *Forecast) IsCounterattackPossible(targetID string, collection *repositories.RepositoryCollection) bool {
	if forecast.setup.IsCounterAttack == false {
		canCounter, _ := forecast.offenseStrategy.CanSquaddieCounterWithEquippedWeapon(targetID, collection)
		if canCounter {
			return true
		}
	}
	return false
}

func (forecast *Forecast) createCounterAttackForecast(counterAttackingSquaddieID string) (*powerusagescenario.Setup, *AttackForecast) {
	counterAttackingSquaddie := forecast.repositories.SquaddieRepo.GetOriginalSquaddieByID(counterAttackingSquaddieID)
	counterAttackingPowerID := counterAttackingSquaddie.GetEquippedPowerID()
	counterAttackingTargetID := forecast.setup.UserID

	counterForecastSetup := powerusagescenario.Setup{
		UserID:          counterAttackingSquaddieID,
		PowerID:         counterAttackingPowerID,
		Targets:         []string{counterAttackingTargetID},
		IsCounterAttack: true,
	}

	counterAttackForecast := Forecast{
		setup: counterForecastSetup,
		repositories: &repositories.RepositoryCollection{
			SquaddieRepo: forecast.repositories.SquaddieRepo,
			PowerRepo:    forecast.repositories.PowerRepo,
		},
	}

	counterAttackForecast.CalculateForecast()

	return &counterForecastSetup, counterAttackForecast.CalculateAttackForecast(counterAttackingTargetID)
}

// CalculateAttackForecast figures out what will happen when this attack power is used.
func (forecast *Forecast) CalculateAttackForecast(targetID string) *AttackForecast {
	attackerContext := *NewAttackerContext(&squaddiestats.CalculateSquaddieOffenseStats{})
	attackerContext.Calculate(forecast.setup, forecast.repositories)

	defenderContext := *NewDefenderContext(targetID, &squaddiestats.CalculateSquaddieDefenseStats{})
	defenderContext.Calculate(&forecast.setup, forecast.repositories)

	versusContext := VersusContext{}
	versusContext.Calculate(attackerContext, defenderContext)

	return &AttackForecast{
		AttackerContext: attackerContext,
		DefenderContext: defenderContext,
		VersusContext:   versusContext,
	}
}

// HealingForecast showcases beneficial abilities
type HealingForecast struct {
	RawHitPointsRestored int
	TargetID             string
}

// CalculateHealingForecast figures out what will happen when this attack power is used.
func (forecast *Forecast) CalculateHealingForecast(targetID string) *HealingForecast {
	maximumHealing, err := forecast.offenseStrategy.GetHitPointsHealedWithPower(
		forecast.setup.UserID,
		forecast.setup.PowerID,
		targetID,
		forecast.repositories,
	)
	if err != nil {
		return &HealingForecast{
			RawHitPointsRestored: 0,
			TargetID:             targetID,
		}
	}

	return &HealingForecast{
		RawHitPointsRestored: maximumHealing,
		TargetID:             targetID,
	}
}

// TODO Delete test hook functions

// UpdateForecastWithNewTarget is a test hook and should be deleted.
func (forecast *Forecast) UpdateForecastWithNewTarget(i int, squaddieID string) {
	forecast.setup.Targets[i] = squaddieID
}

// UpdateForecastWithNewUser is a test hook and should be deleted.
func (forecast *Forecast) UpdateForecastWithNewUser(squaddieID string) {
	forecast.setup.UserID = squaddieID
}
