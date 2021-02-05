package squaddie

// SquaddieMovement contains all of the information needed to describe a Squaddie's movement.
type SquaddieMovement struct {
	Distance  int                  `json:"distance" yaml:"distance"`
	Type      SquaddieMovementType `json:"type" yaml:"type"`
	HitAndRun bool                 `json:"hit_and_run" yaml:"hit_and_run"`
}

// SquaddieMovementType describes how Squaddies traverse terrain. This affects
//   movement costs and crossing pits
type SquaddieMovementType string

const (
	// SquaddieMovementTypeFoot Squaddies take full terrain penalties and cannot cross pits.
	SquaddieMovementTypeFoot = "foot"
	// SquaddieMovementTypeLight Squaddies have no terrain penalties and cannot cross pits.
	SquaddieMovementTypeLight = "light"
	// SquaddieMovementTypeFly Squaddies have no terrain penalties and can cross pits.
	SquaddieMovementTypeFly = "fly"
	// SquaddieMovementTypeTeleport Squaddies have no terrain penalties and can cross pits.
	//   They also ignore walls and other barriers.
	SquaddieMovementTypeTeleport = "teleport"
)

var MovementValueByType = map[SquaddieMovementType]int{
	SquaddieMovementTypeFoot:     0,
	SquaddieMovementTypeLight:    1,
	SquaddieMovementTypeFly:      2,
	SquaddieMovementTypeTeleport: 3,
}