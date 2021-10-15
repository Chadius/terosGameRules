package squaddieclass

import "github.com/chadius/terosbattleserver/entity/squaddieclass"

// ClassBuilderOptions is used to set a squaddie's defensive attributes.
type ClassBuilderOptions struct {
	id                string
	name              string
	baseClassRequired bool
	initialBigLevelID string
}

// ClassFactory creates a ClassBuilderOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func ClassFactory() *ClassBuilderOptions {
	return &ClassBuilderOptions{
		id: "",
		name: "",
		baseClassRequired: false,
		initialBigLevelID: "",
	}
}

// WithID sets the class ID.
func (c *ClassBuilderOptions) WithID(id string) *ClassBuilderOptions {
	c.id = id
	return c
}

// WithName sets the class name.
func (c *ClassBuilderOptions) WithName(name string) *ClassBuilderOptions {
	c.name = name
	return c
}

// WithInitialBigLevelID sets the first big level the class gives.
func (c *ClassBuilderOptions) WithInitialBigLevelID(levelID string) *ClassBuilderOptions {
	c.initialBigLevelID = levelID
	return c
}

// RequiresBaseClass says this class cannot be selected without a base class.
func (c *ClassBuilderOptions) RequiresBaseClass() *ClassBuilderOptions {
	c.baseClassRequired = true
	return c
}



// Build uses the ClassBuilderOptions to create a Movement.
func (c *ClassBuilderOptions) Build() *squaddieclass.Class {
	newClass := &squaddieclass.Class{
		ID: c.id,
		Name: c.name,
		BaseClassRequired: c.baseClassRequired,
		InitialBigLevelID: c.initialBigLevelID,
	}
	return newClass
}
