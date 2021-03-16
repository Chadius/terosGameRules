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
			TargetID: "SquaddieBullseye",
			DamageTaken: 5,
			BarrierDamage: 1,
			WasAHit: true,
			WasACriticalHit: false,
		}
		Expect(attackingPowerReport.DamageTaken).To(Equal(5))
		Expect(attackingPowerReport.BarrierDamage).To(Equal(1))
		Expect(attackingPowerReport.WasACriticalHit).To(BeFalse())
		Expect(attackingPowerReport.TargetID).To(Equal("SquaddieBullseye"))
	})
	It("Can create power reports and set fields", func() {
		powerReport = &report.PowerReport{
			AttackerID: "SquaddieTeros",
			PowerID: "PowerBlot",
			AttackingPowerResults: []*report.AttackingPowerReport{
				{
					TargetID: "SquaddieBandit",
					DamageTaken: 5,
					BarrierDamage: 1,
					WasAHit: true,
					WasACriticalHit: false,
				},
			},
		}
		Expect(powerReport.AttackerID).To(Equal("SquaddieTeros"))
		Expect(powerReport.PowerID).To(Equal("PowerBlot"))

		Expect(powerReport.AttackingPowerResults).To(HaveLen(1))
		Expect(powerReport.AttackingPowerResults[0].TargetID).To(Equal("SquaddieBandit"))
		Expect(powerReport.AttackingPowerResults[0].WasAHit).To(BeTrue())
		Expect(powerReport.AttackingPowerResults[0].WasACriticalHit).To(BeFalse())
	})
})