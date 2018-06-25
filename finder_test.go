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
	models.Truncate(db, "name_strings")
	models.Truncate(db, "page_name_strings")
})

var _ = Describe("Finder", func() {
	Describe("ProcessTitles()", func() {
		It("imports pages to db", func() {
			finder.ProcessTitles(db, dict)
			Expect(models.Count(db, "pages")).To(Equal(8354))
			Expect(models.Count(db, "page_name_strings")).To(Equal(16899))
		})
	})
})
