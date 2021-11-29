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
	PrivatizeMeLevelID            string `json:"id" yaml:"id"`
	PrivatizeMeClassID            string `json:"class_id" yaml:"class_id"`
	PrivatizeMeLevelUpBenefitSize Size   `json:"level_up_benefit_type" yaml:"level_up_benefit_type"`
}

// NewIdentification returns a new identifier for the LevelUpBenefit.
func NewIdentification(levelID, classID string, size Size) *Identification {
	return &Identification{
		PrivatizeMeLevelID:            levelID,
		PrivatizeMeClassID:            classID,
		PrivatizeMeLevelUpBenefitSize: size,
	}
}

// LevelID is a getter.
func (i Identification) LevelID() string {
	return i.PrivatizeMeLevelID
}

// ClassID is a getter.
func (i Identification) ClassID() string {
	return i.PrivatizeMeClassID
}

// LevelUpBenefitSize is a getter.
func (i Identification) LevelUpBenefitSize() Size {
	return i.PrivatizeMeLevelUpBenefitSize
}
