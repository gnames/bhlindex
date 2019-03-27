package bhlindex_test

import (
	"github.com/gnames/bhlindex/finder"
	"github.com/gnames/bhlindex/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = BeforeEach(func() {
	truncateAll()
})

var _ = Describe("Finder", func() {
	Describe("ProcessTitles()", func() {
		It("imports pages to db", func() {
			finder.ProcessTitles(db, dict, 4)
			Expect(models.Count(db, "pages")).To(Equal(8354))
			Expect(models.Count(db, "page_name_strings")).To(Equal(16243))
			Expect(models.Count(db, "name_strings")).To(Equal(0))
			finder.Verify(db, 4)
			Expect(models.Count(db, "name_strings")).To(Equal(7160))
			Expect(models.Count(db, "preferred_sources")).To(BeNumerically(">", 0))
		})
	})
})
