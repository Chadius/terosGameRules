package squaddieClass

// Class contains information about a group of LevelUpBenefits.
type Class struct {
	ID                string `json:"id" yaml:"id"`
	Name              string `json:"name" yaml:"name"`
	BaseClassRequired bool   `json:"base_class_required" yaml:"base_class_required"`
}

