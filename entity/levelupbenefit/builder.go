package levelupbenefit

// Builder is used to create formula objects.
type Builder struct {
	levelID   string
	classID   string
	levelSize Size
}

// NewLevelUpBenefitBuilder returns a new object used to build Term objects.
func NewLevelUpBenefitBuilder() *Builder {
	return &Builder{
		levelID:   "newID",
		classID:   "newClass",
		levelSize: Small,
	}
}

// WithID sets the field.
func (b *Builder) WithID(levelID string) *Builder {
	b.levelID = levelID
	return b
}

// WithClassID sets the field.
func (b *Builder) WithClassID(classID string) *Builder {
	b.classID = classID
	return b
}

// BigLevel means this level is a Big major milestone.
func (b *Builder) BigLevel() *Builder {
	b.levelSize = Big
	return b
}

// Build creates a new LevelUpBenefit object.
func (b *Builder) Build() (*LevelUpBenefit, error) {
	return NewLevelUpBenefit(
		NewIdentification(
			b.levelID,
			b.classID,
			b.levelSize,
		),
	), nil
}
