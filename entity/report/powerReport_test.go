package report_test

import (
	"github.com/cserrant/terosBattleServer/entity/report"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reporting Power usage", func() {
	var (
		powerReport *report.PowerReport
		attackingPowerReport *report.AttackingPowerReport
	)
	It("Can create Attacking Power reports and set fields", func() {
		attackingPowerReport = &report.AttackingPowerReport{
			DamageTaken: 5,
			BarrierDamage: 1,
			WasACriticalHit: false,
		}
		Expect(attackingPowerReport.DamageTaken).To(Equal(5))
		Expect(attackingPowerReport.BarrierDamage).To(Equal(1))
		Expect(attackingPowerReport.WasACriticalHit).To(BeFalse())
	})
	It("Can create power reports and set fields", func() {
		powerReport = &report.PowerReport{
			AttackerID: "SquaddieTeros",
			TargetID: "SquaddieBandit",
			PowerID: "PowerBlot",
			AttackingPowerResults: []*report.AttackingPowerReport{
				{
					DamageTaken: 5,
					BarrierDamage: 1,
					WasACriticalHit: false,
				},
			},
		}
		Expect(powerReport.AttackerID).To(Equal("SquaddieTeros"))
		Expect(powerReport.TargetID).To(Equal("SquaddieBandit"))
		Expect(powerReport.PowerID).To(Equal("PowerBlot"))
	})
})