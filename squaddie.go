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
	innatePowers     []*AttackingPower
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

// GetInnatePowersFromRepository uses a list of power names and retrieves them from the repository.
func (squaddie *Squaddie) GetInnatePowersFromRepository(powerNames []string, powerRepository []*AttackingPower) {
	for _, name := range powerNames {
		for _, power := range powerRepository {
			if name == power.Name {
				squaddie.GainInnatePower(power)
			}
		}
	}
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

// MarshalJSON from json package is overridden.
func (squaddie *Squaddie) MarshalJSON() ([]byte, error) {
	type Alias Squaddie

	return json.Marshal(&struct {
		*Alias
		AttackingPowerIDNames []*PowerIDName `json:"attacking_powers" yaml:"attacking_powers"`
	}{
		Alias:                 (*Alias)(squaddie),
		AttackingPowerIDNames: squaddie.GetInnatePowerIDNames(),
	})
}

// UnmarshalJSON from json package is overridden.
func (squaddie *Squaddie) UnmarshalJSON(byteStream []byte) error {
	type Alias Squaddie
	aux := &struct {
		AttackingPowerIDNames []*PowerIDName `json:"attacking_powers" yaml:"attacking_powers"`
		*Alias
	}{
		Alias: (*Alias)(squaddie),
	}

	err := json.Unmarshal(byteStream, &aux)
	if err != nil {
		return err
	}

	return nil
}

// GetAttackingPowerNamesFromJSON gets the names of the powers from the JSON bytestream.
func GetAttackingPowerNamesFromJSON(byteStream []byte) ([]string, error) {
	return getAttackingPowerNamesFromByteStream(byteStream, json.Unmarshal)
}

// GetAttackingPowerNamesFromYAML gets the names of the powers from the YAML bytestream.
func GetAttackingPowerNamesFromYAML(byteStream []byte) ([]string, error) {
	return getAttackingPowerNamesFromByteStream(byteStream, yaml.Unmarshal)
}

type unmarshalFunc func([]byte, interface{}) error

func getAttackingPowerNamesFromByteStream(byteStream []byte, unmarshal unmarshalFunc) ([]string, error) {
	aux := &struct {
		AttackingPowerNames []string `json:"innate_powers" yaml:"innate_powers"`
	}{}
	err := unmarshal(byteStream, &aux)
	if err != nil {
		return []string{}, err
	}
	return aux.AttackingPowerNames, nil
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

// GainInnatePower gives the Squaddie access to the power.
func (squaddie *Squaddie) GainInnatePower(newPower *AttackingPower) {
	squaddie.innatePowers = append(squaddie.innatePowers, newPower)
}

// GetInnatePowerIDNames returns a list of all the powers the squaddie has access to.
func (squaddie *Squaddie) GetInnatePowerIDNames() []*PowerIDName {
	powerIDNames := []*PowerIDName{}
	for _, power := range squaddie.innatePowers {
		powerIDNames = append(powerIDNames, &PowerIDName{Name: power.Name, ID: power.ID})
	}
	return powerIDNames
}
