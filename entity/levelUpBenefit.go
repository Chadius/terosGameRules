package entity

import (
	"fmt"
)

// LevelUpBenefitType defines the expected sources the LevelUpBenefit could be conjured from.
type LevelUpBenefitType string

const (
	// LevelUpBenefitTypeSmall is for small improvements (stats mostly)
	LevelUpBenefitTypeSmall LevelUpBenefitType = "small"
	// LevelUpBenefitTypeBig is for substantial changes to character (new powers, movement changes)
	LevelUpBenefitTypeBig = "big"
)

// LevelUpBenefit describes how a Squaddie improves upon levelling up.
type LevelUpBenefit struct {
	LevelUpBenefitType  LevelUpBenefitType `json:"level_up_benefit_type" yaml:"level_up_benefit_type"`
	SquaddieName      string
	ClassName         string
	ID                string            `json:"id" yaml:"id"`
	MaxHitPoints      int               `json:"max_hit_points" yaml:"max_hit_points"`
	Aim               int               `json:"aim" yaml:"aim"`
	Strength          int               `json:"strength" yaml:"strength"`
	Mind              int               `json:"mind" yaml:"mind"`
	Dodge             int               `json:"dodge" yaml:"dodge"`
	Deflect           int               `json:"deflect" yaml:"deflect"`
	MaxBarrier        int               `json:"max_barrier" yaml:"max_barrier"`
	Armor             int               `json:"armor" yaml:"armor"`
	PowerIDGained     []*PowerReference `json:"powers" yaml:"powers"`
	PowerIDLost       []*PowerReference `json:"powers_lost" yaml:"powers_lost"`
	Movement          *SquaddieMovement `json:"Movement" yaml:"Movement"`
}

// CheckForErrors ensures the LevelUpBenefit has valid fields
func (benefit *LevelUpBenefit) CheckForErrors() error {
	if benefit.LevelUpBenefitType != LevelUpBenefitTypeSmall && benefit.LevelUpBenefitType != LevelUpBenefitTypeBig {
		return fmt.Errorf(`unknown level up benefit type: "%s"`, benefit.LevelUpBenefitType)
	}
	return nil
}

