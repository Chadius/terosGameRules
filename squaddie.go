package terosbattleserver

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

type Squaddie struct {
	Name             string `json:"name" yaml:"name"`
	Affiliation      string `json:"affiliation" yaml:"affiliation"`
	CurrentHitPoints int    `json:"current_hit_points" yaml:"current_hit_points"`
	MaxHitPoints     int    `json:"max_hit_points" yaml:"max_hit_points"`
	Aim              int    `json:"aim" yaml:"aim"`
	Strength         int    `json:"strength" yaml:"strength"`
	Mind             int    `json:"mind" yaml:"mind"`
	Dodge            int    `json:"dodge" yaml:"dodge"`
	Deflect          int    `json:"deflect" yaml:"deflect"`
	CurrentBarrier   int    `json:"current_barrier" yaml:"current_barrier"`
	MaxBarrier       int    `json:"max_barrier" yaml:"max_barrier"`
	Armor            int    `json:"armor" yaml:"armor"`
}

// NewSquaddie generates a squaddie with default values.
func NewSquaddie(name string) Squaddie {
	newSquaddie := Squaddie{
		Name:             name,
		Affiliation:      "Player",
		CurrentHitPoints: 0,
		MaxHitPoints:     5,
		Aim:              0,
		Strength:         0,
		Mind:             0,
		Dodge:            0,
		Deflect:          0,
		CurrentBarrier:   0,
		MaxBarrier:       0,
		Armor:            0,
	}
	newSquaddie.SetHPToMax()
	return newSquaddie
}

// NewSquaddieFromJSON reads the JSON byte stream to create a new Squaddie.
// 	Defaults to NewSquaddie.
func NewSquaddieFromJSON(data []byte) (newSquaddie Squaddie, err error) {
	newSquaddie = NewSquaddie("NewSquaddieFromJSON")
	err = json.Unmarshal(data, &newSquaddie)
	if err != nil {
		return newSquaddie, err
	}

	err = checkSquaddieForErrors(newSquaddie)
	newSquaddie.SetHPToMax()
	return newSquaddie, err
}

// NewSquaddieFromYAML reads the YAML byte stream to create a new Squaddie.
// 	Defaults to NewSquaddie.
func NewSquaddieFromYAML(data []byte) (newSquaddie Squaddie, err error) {
	newSquaddie = NewSquaddie("NewSquaddieFromYAML")
	err = yaml.Unmarshal(data, &newSquaddie)
	if err != nil {
		return newSquaddie, err
	}

	err = checkSquaddieForErrors(newSquaddie)
	newSquaddie.SetHPToMax()
	return newSquaddie, err
}

// checkSquaddieForErrors makes sure the created squaddie doesn't have an error.
func checkSquaddieForErrors(newSquaddie Squaddie) (newError error) {
	if newSquaddie.Affiliation != "Player" {
		return fmt.Errorf("Squaddie has unknown affiliation: '%s'", newSquaddie.Affiliation)
	}

	return nil
}

// SetHPToMax restores the Squaddie's HitPoints.
func (squaddie *Squaddie) SetHPToMax() {
	squaddie.CurrentHitPoints = squaddie.MaxHitPoints
}

// SetBarrierToMax restores the Squaddie's Barrier.
func (squaddie *Squaddie) SetBarrierToMax() {
	squaddie.CurrentBarrier = squaddie.MaxBarrier
}

// GetDefensiveStatsAgainstPhysical calculates how this squaddie can defend against physical attacks.
func (squaddie *Squaddie) GetDefensiveStatsAgainstPhysical() (evasion, barrierDamageReduction, armorDamageReduction int) {
	return squaddie.Dodge, squaddie.CurrentBarrier, squaddie.Armor
}

// GetDefensiveStatsAgainstSpell calculates how this squaddie can defend against spell attacks.
func (squaddie *Squaddie) GetDefensiveStatsAgainstSpell() (evasion, barrierDamageReduction, armorDamageReduction int) {
	return squaddie.Deflect, squaddie.CurrentBarrier, 0
}

// GetOffensiveStatsWithPhysical calculates the squaddie's bonuses with physical attacks.
func (squaddie *Squaddie) GetOffensiveStatsWithPhysical() (toHitBonus, damageBonus int) {
	return squaddie.Aim, squaddie.Strength
}

// GetOffensiveStatsWithSpell calculates the squaddie's bonuses with Spell attacks.
func (squaddie *Squaddie) GetOffensiveStatsWithSpell() (toHitBonus, damageBonus int) {
	return squaddie.Aim, squaddie.Mind
}
