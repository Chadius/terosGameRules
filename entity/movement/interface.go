package movement

// Interface will shape how movement works with squaddies.
type Interface interface {
	Name() string
	GreaterThan(Interface) bool
}
