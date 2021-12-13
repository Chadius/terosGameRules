package healing

// NewHealingLogic returns a new healing logic object based on the keyword given. Or it returns a nohealing logic.
func NewHealingLogic(keyword string) Interface {
	healingLogicByKeyword := map[string]string{
		"Full":                   "FullMindBonus",
		"full":                   "FullMindBonus",
		"*healing.FullMindBonus": "FullMindBonus",
		"healing.FullMindBonus":  "FullMindBonus",

		"Half":                   "HalfMindBonus",
		"half":                   "HalfMindBonus",
		"*healing.HalfMindBonus": "HalfMindBonus",
		"healing.HalfMindBonus":  "HalfMindBonus",

		"Zero":                   "ZeroMindBonus",
		"zero":                   "ZeroMindBonus",
		"*healing.ZeroMindBonus": "ZeroMindBonus",
		"healing.ZeroMindBonus":  "ZeroMindBonus",
	}

	if healingLogicByKeyword[keyword] == "FullMindBonus" {
		return &FullMindBonus{}
	}

	if healingLogicByKeyword[keyword] == "HalfMindBonus" {
		return &HalfMindBonus{}
	}

	if healingLogicByKeyword[keyword] == "ZeroMindBonus" {
		return &ZeroMindBonus{}
	}

	return &NoHealing{}
}
