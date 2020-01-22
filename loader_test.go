package bhlindex_test

import (
	"github.com/gnames/bhlindex/loader"
	"github.com/gnames/bhlindex/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = BeforeEach(func() {
	truncateAll()
})

var _ = Describe("Loader", func() {

	Describe("FindItems", func() {
		It("gets all items", func() {
			c := make(chan string)
			count := 0
			go loader.FindItems(c)
			for range c {
				count += 1
			}
			// There are 20 items total.
			Expect(count).To(Equal(20))
		})
	})

	Describe("ImportItems", func() {
		It("saves items to database", func() {
			itemsChan := make(chan int)
			// save chan from blocking
			go func(titesChan <-chan int) {
				for range itemsChan {
				}
			}(itemsChan)
			loader.ImportItems(db, itemsChan)
			Expect(models.Count(db, "items")).To(Equal(20))
		})
	})
})
