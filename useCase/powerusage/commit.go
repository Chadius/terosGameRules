package powerusage

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/report"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
)

// CommitPowerUse will apply the given PowerReport.
//    Squaddies will move, Targets will take damage, etc.
func CommitPowerUse(powerReport *report.PowerReport, squaddieRepo *squaddie.Repository, powerRepo *power.Repository) {
	squaddieToEquip := squaddieRepo.GetOriginalSquaddieByID(powerReport.AttackerID)
	SquaddieEquipPower(squaddieToEquip, powerReport.PowerID, powerRepo)
}