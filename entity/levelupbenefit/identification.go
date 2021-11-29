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
	levelID   string
	classID   string
	levelSize Size
}

// NewIdentification returns a new identifier for the LevelUpBenefit.
func NewIdentification(levelID, classID string, size Size) *Identification {
	return &Identification{
		levelID:   levelID,
		classID:   classID,
		levelSize: size,
	}
}

// LevelID is a getter.
func (i Identification) LevelID() string {
	return i.levelID
}

// ClassID is a getter.
func (i Identification) ClassID() string {
	return i.classID
}

// LevelUpBenefitSize is a getter.
func (i Identification) LevelUpBenefitSize() Size {
	return i.levelSize
}
