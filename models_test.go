package bhlindex_test

import (
	"github.com/GlobalNamesArchitecture/bhlindex/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = BeforeEach(func() {
	models.Truncate(db, "titles")
})

var _ = Describe("Models", func() {
	Describe("Title", func() {
		It("creates defaults", func() {
			t := models.Title{InternetArchiveID: "test"}
			Expect(t.InternetArchiveID).To(Equal("test"))
			Expect(t.Status).To(Equal(0))
		})
	})

	Describe("title.Insert()", func() {
		It("inserts title returns id", func() {
			t := models.Title{InternetArchiveID: "test"}
			Expect(t.Insert(db)).To(BeNumerically(">", 0))
		})

		It("ignores duplicates", func() {
			t := models.Title{InternetArchiveID: "test"}
			id := t.Insert(db)
			Expect(id).To(BeNumerically(">", 0))
			t2 := models.Title{InternetArchiveID: "test"}
			id2 := t2.Insert(db)
			Expect(id2).To(Equal(0))
		})
	})
})
