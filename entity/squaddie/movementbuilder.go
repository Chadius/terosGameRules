package squaddie

import "github.com/chadius/terosgamerules/entity/movement"

// MovementBuilderOptions is used to create healing effects.
type MovementBuilderOptions struct {
	distance      int
	canHitAndRun  bool
	movementLogic movement.Interface
}

// MovementBuilder creates a MovementBuilderOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func MovementBuilder() *MovementBuilderOptions {
	return &MovementBuilderOptions{
		distance:      0,
		canHitAndRun:  false,
		movementLogic: movement.NewMovementLogic("foot"),
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
	return m.WithMovementLogicKeyword("foot")
}

// Light sets the movement type to Light.
func (m *MovementBuilderOptions) Light() *MovementBuilderOptions {
	return m.WithMovementLogicKeyword("light")
}

// Fly sets the movement type to Fly.
func (m *MovementBuilderOptions) Fly() *MovementBuilderOptions {
	return m.WithMovementLogicKeyword("fly")
}

// Teleport sets the movement type to Teleport.
func (m *MovementBuilderOptions) Teleport() *MovementBuilderOptions {
	return m.WithMovementLogicKeyword("teleport")
}

// WithMovementLogicKeyword creates movement logic object using the given keyword.
func (m *MovementBuilderOptions) WithMovementLogicKeyword(keyword string) *MovementBuilderOptions {
	m.movementLogic = movement.NewMovementLogic(keyword)
	return m
}

// Build uses the MovementBuilderOptions to create a Movement.
func (m *MovementBuilderOptions) Build() *Movement {
	newMovement := NewMovement(m.distance, m.canHitAndRun, m.movementLogic)
	return newMovement
}
