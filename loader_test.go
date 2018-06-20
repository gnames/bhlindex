package bhlindex_test

import (
	"github.com/gnames/bhlindex/loader"
	"github.com/gnames/bhlindex/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = BeforeEach(func() {
	models.Truncate(db, "titles")
	models.Truncate(db, "pages")
})

var _ = Describe("Loader", func() {

	Describe("FindTitles", func() {
		It("gets all titles", func() {
			c := make(chan string)
			count := 0
			go loader.FindTitles(c)
			for _ = range c {
				count += 1
			}
			// There are 20 titles total.
			Expect(count).To(Equal(20))
		})
	})

	Describe("ImportTitles", func() {
		It("saves titles to database", func() {
			titlesChan := make(chan int)
			go func() {
				for _ = range titlesChan {
				}
			}()
			loader.ImportTitles(db, titlesChan)
			Expect(models.Count(db, "titles")).To(Equal(20))
		})
	})
})
