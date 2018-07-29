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
		FIt("imports pages to db", func() {
			finder.ProcessTitles(db, dict)
			Expect(models.Count(db, "pages")).To(Equal(8354))
			Expect(models.Count(db, "page_name_strings")).To(Equal(16899))
			Expect(models.Count(db, "name_strings")).To(Equal(7215))
		})
	})
})
