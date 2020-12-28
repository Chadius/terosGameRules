package terosbattleserver

import (
	"encoding/json"
	"fmt"
)

type Squaddie struct {
	Name             string `json:"name"`
	Affiliation      string `json:"affiliation"`
	CurrentHitPoints int    `json:"current_hit_points"`
	MaxHitPoints     int    `json:"max_hit_points"`
	Aim              int    `json:"aim"`
	Strength         int    `json:"strength"`
	Mind             int    `json:"mind"`
	Dodge            int    `json:"dodge"`
	Deflect          int    `json:"deflect"`
	CurrentBarrier   int    `json:"current_barrier"`
	MaxBarrier       int    `json:"max_barrier"`
	Armor            int    `json:"armor"`
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

	if newSquaddie.Affiliation != "Player" {
		err = fmt.Errorf("Squaddie has unknown affiliation: '%s'", newSquaddie.Affiliation)
	}

	newSquaddie.SetHPToMax()
	return newSquaddie, err
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
