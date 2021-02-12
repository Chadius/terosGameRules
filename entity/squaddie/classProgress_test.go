package squaddie_test

import (
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClassProgress", func() {
	It("Can tell if a class has been used", func() {
		progress := &squaddie.ClassProgress{
			ClassID:        "007",
			ClassName:      "Superspy",
			LevelsConsumed: []string{"1", "2", "3", "4", "5"},
		}
		Expect(progress.IsLevelAlreadyConsumed("1")).To(BeTrue())
		Expect(progress.IsLevelAlreadyConsumed("10")).To(BeFalse())
	})
})
