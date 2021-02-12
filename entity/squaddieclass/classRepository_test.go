package squaddieclass_test

import (
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Managing Class objects", func() {
	var (
		repo *squaddieclass.Repository
		jsonByteStream []byte
	)
	BeforeEach(func() {
		repo = squaddieclass.NewRepository()
		jsonByteStream = []byte(`[{
			"ID": "aaaaaaaa",
			"name": "Mage"
		}]`)
	})
	Context("Load Class using different sources", func() {
		It("Can add a JSON source", func() {
			Expect(repo.GetNumberOfClasses()).To(Equal(0))
			success, _ := repo.AddJSONSource(jsonByteStream)
			Expect(success).To(BeTrue())
			Expect(repo.GetNumberOfClasses()).To(Equal(1))
		})
		It("Can add classes directly", func() {
			listOfClasses := []*squaddieclass.Class{
				{
					ID:                "class1",
					Name:              "Mage",
					BaseClassRequired: false,
				},
				{
					ID:                "class2",
					Name:              "Dimension Walker",
					BaseClassRequired: true,
				},
			}
			Expect(repo.GetNumberOfClasses()).To(Equal(0))
			success, _ := repo.AddListOfClasses(listOfClasses)
			Expect(success).To(BeTrue())
			Expect(repo.GetNumberOfClasses()).To(Equal(2))
		})
	})
	Context("Can retrieve classes by ID", func() {
		var (
			mageClass *squaddieclass.Class
			dimensionWalkerClass *squaddieclass.Class
		)
		BeforeEach(func() {
			mageClass = &squaddieclass.Class{
				ID:                "class0",
				Name:              "Mage",
				BaseClassRequired: false,
			}
			dimensionWalkerClass = &squaddieclass.Class{
				ID:                "class1",
				Name:              "Dimension Walker",
				BaseClassRequired: true,
			}
			repo.AddListOfClasses([]*squaddieclass.Class{mageClass, dimensionWalkerClass})
		})
		It("Can retrieve classes by ID", func() {
			foundClass, err := repo.GetClassByID(mageClass.ID)
			Expect(err).To(BeNil())
			Expect(foundClass.ID).To(Equal(mageClass.ID))
			Expect(foundClass.Name).To(Equal(mageClass.Name))
		})
		It("Raises an error when class ID does not exist", func() {
			_, err := repo.GetClassByID("bad ID")
			Expect(err.Error()).To(Equal(`class repository: No class found with ID: "bad ID"`))
		})
	})
})
