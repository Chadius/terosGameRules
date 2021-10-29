package squaddie

import (
	"github.com/chadius/terosbattleserver/entity/squaddie"
)

// MovementBuilderOptions is used to create healing effects.
type MovementBuilderOptions struct {
	distance     int
	canHitAndRun bool
	movementType squaddie.MovementType
}

// MovementBuilder creates a MovementBuilderOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func MovementBuilder() *MovementBuilderOptions {
	return &MovementBuilderOptions{
		distance:     0,
		canHitAndRun: false,
		movementType: squaddie.Foot,
	}
}

// Distance sets the distance the squaddie can travel.
func (m *MovementBuilderOptions) Distance(distance int) *MovementBuilderOptions {
	m.distance = distance
	return m
}

// CanHitAndRun means you can move after acting
func (m *MovementBuilderOptions) CanHitAndRun() *MovementBuilderOptions {
	m.canHitAndRun = true
	return m
}

// Foot sets the movement type to Foot.
func (m *MovementBuilderOptions) Foot() *MovementBuilderOptions {
	m.movementType = squaddie.Foot
	return m
}

// Light sets the movement type to Light.
func (m *MovementBuilderOptions) Light() *MovementBuilderOptions {
	m.movementType = squaddie.Light
	return m
}

// Fly sets the movement type to Fly.
func (m *MovementBuilderOptions) Fly() *MovementBuilderOptions {
	m.movementType = squaddie.Fly
	return m
}

// Teleport sets the movement type to Teleport.
func (m *MovementBuilderOptions) Teleport() *MovementBuilderOptions {
	m.movementType = squaddie.Teleport
	return m
}

// Build uses the MovementBuilderOptions to create a Movement.
func (m *MovementBuilderOptions) Build() *squaddie.Movement {
	newMovement := &squaddie.Movement{
		SquaddieMovementDistance:     m.distance,
		SquaddieMovementCanHitAndRun: m.canHitAndRun,
		SquaddieMovementType:         m.movementType,
	}
	return newMovement
}
