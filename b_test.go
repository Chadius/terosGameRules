package terosbattleserver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	terosbattleserver "github.com/cserrant/terosBattleServer"
)

var _ = Describe("B", func() {
	It("Math works", func() {
		g := terosbattleserver.Gub{Stuff: 4}
		Expect(g.Stuff).To(Equal(2+2))
	})
})
