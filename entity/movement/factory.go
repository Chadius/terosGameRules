package movement

// NewMovementLogic returns a new movement logic object based on the keyword given.
//  Defaults to foot based movement.
func NewMovementLogic(keyword string) Interface {
	movementLogicByKeyword := map[string]string{
		"light":           "light",
		"Light":           "light",
		"*movement.light": "light",
		"*movement.Light": "light",

		"fly":           "fly",
		"Fly":           "fly",
		"*movement.fly": "fly",
		"*movement.Fly": "fly",

		"teleport":           "teleport",
		"Teleport":           "teleport",
		"*movement.teleport": "teleport",
		"*movement.Teleport": "teleport",
	}

	if movementLogicByKeyword[keyword] == "light" {
		return &Light{}
	}

	if movementLogicByKeyword[keyword] == "fly" {
		return &Fly{}
	}

	if movementLogicByKeyword[keyword] == "teleport" {
		return &Teleport{}
	}

	return &Foot{}
}
