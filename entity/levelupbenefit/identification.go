package levelupbenefit

// Size defines the expected sources the LevelUpBenefit could be conjured from.
type Size string

const (
	// Small is for small improvements (stats mostly)
	Small Size = "small"
	// Big is for substantial changes to character (new powers, movement changes)
	Big Size = "big"
)

// Identification is used to uniquely identify this LevelUpBenefit.
type Identification struct {
	ClassID            string `json:"class_id" yaml:"class_id"`
	LevelUpBenefitType Size   `json:"level_up_benefit_type" yaml:"level_up_benefit_type"`
	ID                 string `json:"id" yaml:"id"`
}
