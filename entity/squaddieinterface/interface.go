package squaddieinterface

// Interface will shape how healing powers work with squaddies.
type Interface interface {
	Mind() int
	CurrentHitPoints() int
	MaxHitPoints() int
}
