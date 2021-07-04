package squaddie

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/utility"
)

// Affiliation describes the Squaddie's allegiance
type Affiliation string

const (
	// Player Squaddies are controlled by the player and try to succeed at missions.
	Player = "Player"
	// Enemy Squaddies oppose the player
	Enemy = "Enemy"
)

// Squaddie represents a person, creature or thing that can take actions on a battlefield.
type Squaddie struct {
	Identification				Identification              `json:"identification" yaml:"identification"`
	ClassProgress				ClassProgress				`json:"class_progress" yaml:"class_progress"`
	Defense						Defense						`json:"defense" yaml:"defense"`
	Offense						Offense						`json:"offense" yaml:"offense"`
	Movement					Movement					`json:"movement" yaml:"movement"`
	PowerCollection				PowerCollection				`json:"powers" yaml:"powers"`
}

// NewSquaddie generates a squaddie with maxed out health.
func NewSquaddie(name string) *Squaddie {
	newSquaddie := Squaddie{
		Identification: Identification{
			ID:                  "squaddie_" + utility.StringWithCharset(8, "abcdefgh0123456789"),
			Name:                name,
			Affiliation:         Player,
		},
		ClassProgress: ClassProgress{
			ClassLevelsConsumed: map[string]*ClassLevelsConsumed{},
		},
		Defense: Defense{
			CurrentHitPoints:    0,
			MaxHitPoints:        5,
			Dodge:               0,
			Deflect:             0,
			CurrentBarrier:      0,
			MaxBarrier:          0,
			Armor:               0,
		},
		Offense: Offense{
			Aim:      0,
			Strength: 0,
			Mind:     0,
		},
		Movement:            Movement{
			Distance:        3,
			Type:            Foot,
			HitAndRun:       false,
		},
	}
	newSquaddie.Defense.SetHPToMax()
	return &newSquaddie
}

// CheckSquaddieForErrors makes sure the created squaddie doesn't have an error.
func CheckSquaddieForErrors(newSquaddie *Squaddie) (newError error) {
	if newSquaddie.Identification.Affiliation != Player &&
		newSquaddie.Identification.Affiliation != Enemy {
		return fmt.Errorf("Squaddie has unknown affiliation: '%s'", newSquaddie.Identification.Affiliation)
	}

	return nil
}
