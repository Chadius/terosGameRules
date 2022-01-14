package builder

import (
	"github.com/chadius/terosgamerules/entity/levelupbenefit"
	"strconv"
)

// LevelGeneratorInstruction tells the LevelGenerator how to build the levels.
type LevelGeneratorInstruction struct {
	NumberOfLevels int
	ClassID        string
	PrefixLevelID  string
	Type           levelupbenefit.Size
}

// LevelGenerator is a Builder pattern that can be programmed to generate a slice of LevelUpBenefits.
type LevelGenerator struct {
	Instructions *LevelGeneratorInstruction
}

// Build follows the instructions to create a slice of levelUpBenefits.
func (generator *LevelGenerator) Build() []*levelupbenefit.LevelUpBenefit {
	levels := []*levelupbenefit.LevelUpBenefit{}
	for i := 0; i < generator.Instructions.NumberOfLevels; i++ {
		builder := levelupbenefit.NewLevelUpBenefitBuilder().
			WithID(generator.Instructions.PrefixLevelID + strconv.Itoa(i)).
			WithClassID(generator.Instructions.ClassID)

		if generator.Instructions.Type == levelupbenefit.Big {
			builder.BigLevel()
		}

		newLevelUpBenefit, _ := builder.Build()
		levels = append(levels, newLevelUpBenefit)
	}
	return levels
}
