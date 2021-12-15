package powersource

// NewPowerSourceLogic returns a new power source logic object based on the keyword given. Or it returns a nohealing logic.
func NewPowerSourceLogic(keyword string) Interface {
	powersourceLogicByKeyword := map[string]string{
		"spell":              "spell",
		"Spell":              "spell",
		"*powersource.spell": "spell",
		"*powersource.Spell": "spell",
	}

	if powersourceLogicByKeyword[keyword] == "spell" {
		return &Spell{}
	}

	return &Physical{}
}
