package main_test

import (
	"time"

	"github.com/GlobalNamesArchitecture/bhlindex/models"
	"github.com/GlobalNamesArchitecture/bhlindex/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = AfterEach(func() {
	models.Truncate(db, "titles")
})

var _ = Describe("Models", func() {
	Describe("Title", func() {
		It("creates defaults", func() {
			t := models.Title{InternetArchiveID: "test"}
			t.Defaults()
			Expect(t.InternetArchiveID).To(Equal("test"))
			Expect(t.Status).To(Equal(0))
		})

		Describe("Title.CreateOrSelect()", func() {
			It("Saves the title to the database", func() {
				iaID := util.UUID4()
				t := models.Title{InternetArchiveID: iaID}
				t.Defaults()
				t.CreateOrSelect(db)

				Expect(t.ID).To(BeNumerically(">", 0))
				Expect(t.UpdatedAt.Unix()).To(BeNumerically("~", time.Now().Unix(), 5))

				t2 := models.Title{InternetArchiveID: iaID, EnglishDetected: true}
				t2.CreateOrSelect(db)

				Expect(t2.ID).To(Equal(t.ID))
				Expect(t2.EnglishDetected).To(Equal(false))
			})
		})

		Describe("Title.Delete()", func() {
			It("Deletes a record", func() {
				c := models.Count(db, "titles")
				t := models.Title{InternetArchiveID: util.UUID4()}
				t.CreateOrSelect(db)
				Expect(models.Count(db, "titles")).To(Equal(c + 1))
				t.Delete(db)
				Expect(models.Count(db, "titles")).To(Equal(c))
			})
		})
	})
})
