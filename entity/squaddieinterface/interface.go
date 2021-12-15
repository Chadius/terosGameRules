package squaddieinterface

// Interface will shape how healing powers work with squaddies.
type Interface interface {
	CurrentHitPoints() int
	MaxHitPoints() int
	CurrentBarrier() int
	Dodge() int
	Deflect() int
	Armor() int

	Mind() int
	Strength() int
}
