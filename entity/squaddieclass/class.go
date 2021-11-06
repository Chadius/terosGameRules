package squaddieclass

// Class contains information about a group of LevelUpBenefits.
type Class struct {
	id                string
	name              string
	baseClassRequired bool
	initialBigLevelID string
}

// NewClass returns a new class object.
func NewClass(classID, className string, baseClassRequired bool, classInitialBigLevelID string) *Class {
	return &Class{
		id:                classID,
		name:              className,
		baseClassRequired: baseClassRequired,
		initialBigLevelID: classInitialBigLevelID,
	}
}

// ID returns the class ID.
func (c *Class) ID() string {
	return c.id
}

// Name returns the class Name.
func (c *Class) Name() string {
	return c.name
}

// BaseClassRequired returns true if this class requires the squaddie to satisfy their base class before using it.
func (c *Class) BaseClassRequired() bool {
	return c.baseClassRequired
}

// InitialBigLevelID returns the class's first level ID, if needed.
func (c *Class) InitialBigLevelID() string {
	return c.initialBigLevelID
}

// ClassReference is a lightweight way to refer to classes
type ClassReference struct {
	ID   string
	Name string
}

// GetReference returns the ClassReference for a given class.
func (c *Class) GetReference() *ClassReference {
	return &ClassReference{ID: c.ID(), Name: c.Name()}
}
