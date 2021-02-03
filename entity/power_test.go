package entity_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cserrant/terosBattleServer/entity"
)

var _ = Describe("Calculate combination of Attacking Power and Squaddie", func() {
	It("Calculates the crit chance based on To Hit Bonus", func() {
		Expect(entity.GetChanceToHitBasedOnHitRate(9001)).To(Equal(36))
		Expect(entity.GetChanceToHitBasedOnHitRate(5)).To(Equal(36))

		Expect(entity.GetChanceToHitBasedOnHitRate(-6)).To(Equal(0))
		Expect(entity.GetChanceToHitBasedOnHitRate(-9001)).To(Equal(0))

		Expect(entity.GetChanceToHitBasedOnHitRate(4)).To(Equal(35))
		Expect(entity.GetChanceToHitBasedOnHitRate(3)).To(Equal(33))
		Expect(entity.GetChanceToHitBasedOnHitRate(2)).To(Equal(30))
		Expect(entity.GetChanceToHitBasedOnHitRate(1)).To(Equal(26))
		Expect(entity.GetChanceToHitBasedOnHitRate(0)).To(Equal(21))
		Expect(entity.GetChanceToHitBasedOnHitRate(-1)).To(Equal(15))
		Expect(entity.GetChanceToHitBasedOnHitRate(-2)).To(Equal(10))
		Expect(entity.GetChanceToHitBasedOnHitRate(-3)).To(Equal(6))
		Expect(entity.GetChanceToHitBasedOnHitRate(-4)).To(Equal(3))
		Expect(entity.GetChanceToHitBasedOnHitRate(-5)).To(Equal(1))
	})

	It("Will get a random ID if none is given", func() {
		powerWithoutID := entity.NewPower("New Attack")
		Expect(powerWithoutID.ID).NotTo(BeNil())
		Expect(powerWithoutID.ID).NotTo(Equal(""))
	})
})
