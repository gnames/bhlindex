package bhlindex_test

import (
	"github.com/gnames/bhlindex/finder"
	"github.com/gnames/bhlindex/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = BeforeEach(func() {
	models.Truncate(db, "titles")
	models.Truncate(db, "pages")
})

var _ = Describe("Finder", func() {
	Describe("ProcessTitles()", func() {
		It("imports pages to db", func() {
			finder.ProcessTitles(db, dict)
			Expect(models.Count(db, "pages")).To(Equal(8354))
		})
	})
})
