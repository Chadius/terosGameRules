package terosbattleserver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	terosbattleserver "github.com/cserrant/terosBattleServer"
)

var _ = Describe("CRUD Squaddies", func() {
	var (
		repo *terosbattleserver.SquaddieRetriever
	)
	BeforeEach(func() {
		repo = terosbattleserver.NewSquaddieRetriever()
	})
	Context("Using JSON sources", func() {
		It("Can add a JSON source", func() {
			Expect(repo.GetNumberOfSquaddies()).To(Equal(0))
			jsonByteStream := []byte(`[{
				"name": "Teros",
				"aim": 5
			}]`)
			repo.AddJSONSource(jsonByteStream)
			Expect(repo.GetNumberOfSquaddies()).To(Equal(1))
		})
		It("Can get a Squaddie by name", func() {
			jsonByteStream := []byte(`[{
				"name": "Teros",
				"aim": 5
			}]`)
			repo.AddJSONSource(jsonByteStream)

			teros := repo.GetByName("Teros")
			Expect(teros).NotTo(BeNil())
			Expect(teros.Name).To(Equal("Teros"))

			missingno := repo.GetByName("Does not exist")
			Expect(missingno).To(BeNil())
		})
	})
})
