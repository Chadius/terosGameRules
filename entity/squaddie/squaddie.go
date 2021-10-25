package squaddie

import (
	"fmt"
	"github.com/chadius/terosbattleserver/utility"
)

// Affiliation describes the Squaddie's allegiance
type Affiliation string

// Squaddie Affiliation constants
const (
	Player  Affiliation = "Player"
	Enemy   Affiliation = "Enemy"
	Ally    Affiliation = "Ally"
	Neutral Affiliation = "Neutral"
)

// Squaddie represents a person, creature or thing that can take actions on a battlefield.
type Squaddie struct {
	Identification  Identification  `json:"identification" yaml:"identification"`
	ClassProgress   ClassProgress   `json:"class_progress" yaml:"class_progress"`
	Defense         Defense         `json:"defense" yaml:"defense"`
	Offense         Offense         `json:"offense" yaml:"offense"`
	Movement        Movement        `json:"movement" yaml:"movement"`
	PowerCollection PowerCollection `json:"powers" yaml:"powers"`
}

// NewSquaddie generates a squaddie with maxed out health.
func NewSquaddie(name string) *Squaddie {
	newSquaddie := Squaddie{
		Identification: Identification{
			SquaddieID:          "squaddie_" + utility.StringWithCharset(8, "abcdefgh0123456789"),
			SquaddieName:        name,
			SquaddieAffiliation: Player,
		},
		ClassProgress: ClassProgress{
			ClassLevelsConsumed: map[string]*ClassLevelsConsumed{},
		},
		Defense: Defense{
			SquaddieCurrentHitPoints: 0,
			SquaddieMaxHitPoints:     5,
			SquaddieDodge:            0,
			SquaddieDeflect:          0,
			SquaddieCurrentBarrier:   0,
			SquaddieMaxBarrier:       0,
			SquaddieArmor:            0,
		},
		Offense: Offense{
			SquaddieAim:      0,
			SquaddieStrength: 0,
			SquaddieMind:     0,
		},
		Movement: Movement{
			Distance:  3,
			Type:      Foot,
			HitAndRun: false,
		},
	}
	newSquaddie.Defense.SetHPToMax()
	return &newSquaddie
}

// TODO Ask Identification if it has this error

// CheckSquaddieForErrors makes sure the created squaddie doesn't have an error.
func CheckSquaddieForErrors(newSquaddie *Squaddie) (newError error) {
	if newSquaddie.Affiliation() != Player &&
		newSquaddie.Affiliation() != Enemy &&
		newSquaddie.Affiliation() != Ally &&
		newSquaddie.Affiliation() != Neutral {
		newError := fmt.Errorf("squaddie %s has unknown affiliation: '%s'", newSquaddie.ID(), newSquaddie.Affiliation())
		utility.Log(newError.Error(), 0, utility.Error)
		return newError
	}

	return nil
}

// ID delegates.
func (s *Squaddie) ID() string {
	return s.Identification.ID()
}

// Affiliation delegates.
func (s *Squaddie) Affiliation() Affiliation {
	return s.Identification.Affiliation()
}

// Name delegates.
func (s *Squaddie) Name() string {
	return s.Identification.Name()
}

// MaxHitPoints delegates.
func (s *Squaddie) MaxHitPoints() int {
	return s.Defense.MaxHitPoints()
}

// Dodge delegates.
func (s *Squaddie) Dodge() int {
	return s.Defense.Dodge()
}

// Deflect delegates.
func (s *Squaddie) Deflect() int {
	return s.Defense.Deflect()
}

// MaxBarrier delegates.
func (s *Squaddie) MaxBarrier() int {
	return s.Defense.MaxBarrier()
}

// Armor delegates.
func (s *Squaddie) Armor() int {
	return s.Defense.Armor()
}

// CurrentHitPoints delegates.
func (s *Squaddie) CurrentHitPoints() int {
	return s.Defense.CurrentHitPoints()
}

// CurrentBarrier delegates.
func (s *Squaddie) CurrentBarrier() int {
	return s.Defense.CurrentBarrier()
}

// Aim delegates.
func (s *Squaddie) Aim() int {
	return s.Offense.Aim()
}

// Strength delegates.
func (s *Squaddie) Strength() int {
	return s.Offense.Strength()
}

// Mind delegates.
func (s *Squaddie) Mind() int {
	return s.Offense.Mind()
}
