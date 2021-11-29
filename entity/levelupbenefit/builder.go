package levelupbenefit

// Builder is used to create formula objects.
type Builder struct {
	levelID   string
	classID   string
	levelSize Size

	hitPoints int
	deflect int
	dodge int
	barrier int
	armor int

	aim int
	strength int
	mind int
}

// NewLevelUpBenefitBuilder returns a new object used to build Term objects.
func NewLevelUpBenefitBuilder() *Builder {
	return &Builder{
		levelID:   "newID",
		classID:   "newClass",
		levelSize: Small,

		hitPoints: 0,
		deflect: 0,
		dodge: 0,
		barrier: 0,
		armor: 0,

		aim: 0,
		strength: 0,
		mind: 0,
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


// HitPoints increases the defensive parameter
func (b * Builder) HitPoints(hitPoints int) *Builder {
	b.hitPoints = hitPoints
	return b
}

// Deflect increases the defensive parameter
func (b * Builder) Deflect(deflect int) *Builder {
	b.deflect = deflect
	return b
}

// Dodge increases the defensive parameter
func (b * Builder) Dodge(dodge int) *Builder {
	b.dodge = dodge
	return b
}

// Barrier increases the defensive parameter
func (b * Builder) Barrier(barrier int) *Builder {
	b.barrier = barrier
	return b
}

// Armor increases the defensive parameter
func (b * Builder) Armor(armor int) *Builder {
	b.armor = armor
	return b
}

// Aim increases the offensive parameter.
func (b * Builder) Aim(aim int) *Builder {
	b.aim = aim
	return b
}

// Strength increases the offensive parameter.
func (b * Builder) Strength(strength int) *Builder {
	b.strength = strength
	return b
}

// Mind increases the offensive parameter.
func (b * Builder) Mind(mind int) *Builder {
	b.mind = mind
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
		NewDefense(
			b.hitPoints,
			b.dodge,
			b.deflect,
			b.barrier,
			b.armor,
		),
		NewOffense(
			b.aim,
			b.strength,
			b.mind,
		),
	), nil
}