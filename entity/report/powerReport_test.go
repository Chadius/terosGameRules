package report_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/report"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reporting Power usage", func() {
	var (
		powerReport *report.PowerReport,
		attackingPowerReport *report.AttackingPowerReport,
	)
	It("Can create Attacking Power reports and set fields", func() {
		attackingPowerReport = report.AttackingPowerReport{
			DamageTaken: 5,
			BarrierDamage: 1,
			WasACriticalHit: false,
		}

		Expect(attackingPowerReport.DamageTaken).To(Equal(5))
		Expect(attackingPowerReport.BarrierDamage).To(Equal(1))
		Expect(attackingPowerReport.WasACrititicalHit).To(BeFalse())
	})
	It("Can create power reports and set fields", func() {
		powerReport = report.PowerReport{
			AttackerID: "SquaddieTeros",
			TargetID: "SquaddieBandit",
			PowerID: "PowerBlot",
			AttackingPowerReport: report.AttackingPowerReport{
				DamageTaken: 5,
				BarrierDamage: 1,
				WasACriticalHit: false,
			}
		}
		
		Expect(repo.GetNumberOfPowers()).To(Equal(0))
		spear := power.NewPower("Spear")
		spear.PowerType = power.Physical
		newPowers := []*power.Power{spear}
		success, _ := repo.AddSlicePowerSource(newPowers)
		Expect(success).To(BeTrue())
		Expect(repo.GetNumberOfPowers()).To(Equal(1))
	})
})