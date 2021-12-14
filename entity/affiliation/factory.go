package affiliation

// NewAffiliationLogic returns a new healing logic object based on the keyword given. Or it returns a nohealing logic.
func NewAffiliationLogic(keyword string) Interface {
	affiliationLogicByKeyword := map[string]string{
		"Player":              "Player",
		"player":              "Player",
		"*affiliation.Player": "Player",
		"*affiliation.player": "Player",

		"Enemy":              "Enemy",
		"enemy":              "Enemy",
		"*affiliation.Enemy": "Enemy",
		"*affiliation.enemy": "Enemy",

		"ally":              "Ally",
		"Ally":              "Ally",
		"*affiliation.ally": "Ally",
		"*affiliation.Ally": "Ally",
	}

	if affiliationLogicByKeyword[keyword] == "Player" {
		return &Player{}
	}

	if affiliationLogicByKeyword[keyword] == "Enemy" {
		return &Enemy{}
	}

	if affiliationLogicByKeyword[keyword] == "Ally" {
		return &Ally{}
	}

	return &Neutral{}
}
