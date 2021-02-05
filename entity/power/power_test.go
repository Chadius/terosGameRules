package power_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Calculate combination of Attacking Power and Squaddie", func() {
	It("Calculates the crit chance based on To Hit Bonus", func() {
		Expect(power.GetChanceToHitBasedOnHitRate(9001)).To(Equal(36))
		Expect(power.GetChanceToHitBasedOnHitRate(5)).To(Equal(36))

		Expect(power.GetChanceToHitBasedOnHitRate(-6)).To(Equal(0))
		Expect(power.GetChanceToHitBasedOnHitRate(-9001)).To(Equal(0))

		Expect(power.GetChanceToHitBasedOnHitRate(4)).To(Equal(35))
		Expect(power.GetChanceToHitBasedOnHitRate(3)).To(Equal(33))
		Expect(power.GetChanceToHitBasedOnHitRate(2)).To(Equal(30))
		Expect(power.GetChanceToHitBasedOnHitRate(1)).To(Equal(26))
		Expect(power.GetChanceToHitBasedOnHitRate(0)).To(Equal(21))
		Expect(power.GetChanceToHitBasedOnHitRate(-1)).To(Equal(15))
		Expect(power.GetChanceToHitBasedOnHitRate(-2)).To(Equal(10))
		Expect(power.GetChanceToHitBasedOnHitRate(-3)).To(Equal(6))
		Expect(power.GetChanceToHitBasedOnHitRate(-4)).To(Equal(3))
		Expect(power.GetChanceToHitBasedOnHitRate(-5)).To(Equal(1))
	})

	It("Will get a random ID if none is given", func() {
		powerWithoutID := power.NewPower("New Attack")
		Expect(powerWithoutID.ID).NotTo(BeNil())
		Expect(powerWithoutID.ID).NotTo(Equal(""))
	})
})
