package powerinterface

// Interface shapes the power.
type Interface interface {
	ID() string
	HitPointsHealed() int
}
