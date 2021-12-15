package affiliation

// Interface will shape how affiliations work with squaddies.
type Interface interface {
	IsFriendsWith(other Interface) bool
	IsFoesWith(other Interface) bool
	Name() string
}

// TODO Add functions that compare this to a squaddie
