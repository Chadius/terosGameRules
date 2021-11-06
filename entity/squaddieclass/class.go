package squaddieclass

// Class contains information about a group of LevelUpBenefits.
type Class struct {
	id   string
	name                                string
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

func (c *Class) ID() string {
	return c.id
}

func (c *Class) Name() string {
	return c.name
}

func (c *Class) BaseClassRequired() bool {
	return c.baseClassRequired
}

func (c *Class) InitialBigLevelID() string {
	return c.initialBigLevelID
}

type ClassReference struct {
	ID		string
	Name	string
}

func (c *Class) GetReference() *ClassReference {
	return &ClassReference{ID: c.ID(), Name: c.Name()}
}