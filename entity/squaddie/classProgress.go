package squaddie

// ClassProgress tracks information about how a Squaddie uses a class as well
//    as the progress made in a given class.
type ClassProgress struct {
	ClassID string `json:"id" yaml:"id"`
	ClassName string `json:"name" yaml:"name"`
	LevelsConsumed []string `json:"levels_gained" yaml:"levels_gained"`
}

// IsLevelAlreadyConsumed returns true if the level ID has already been used.
func (progress *ClassProgress) IsLevelAlreadyConsumed(levelUpBenefitLevelID string) bool {
	for _, levelID := range progress.LevelsConsumed {
		if levelID == levelUpBenefitLevelID {
			return true
		}
	}
	return false
}
