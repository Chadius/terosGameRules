package target

// NewTargetingLogic returns a new target logic object based on the keyword given.
//  Defaults to Self target.
func NewTargetingLogic(keyword string) Interface {
	targetLogicByKeyword := map[string]string{
		"friend":           "friend",
		"Friend":           "friend",
		"*movement.friend": "friend",
		"*movement.Friend": "friend",

		"foe":           "foe",
		"Foe":           "foe",
		"*movement.foe": "foe",
		"*movement.Foe": "foe",
	}

	if targetLogicByKeyword[keyword] == "friend" {
		return &Friend{}
	}

	if targetLogicByKeyword[keyword] == "foe" {
		return &Foe{}
	}

	return &Self{}
}
