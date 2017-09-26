package bhlindex_test

import (
	"github.com/GlobalNamesArchitecture/bhlindex/finder"
	"github.com/GlobalNamesArchitecture/bhlindex/models"
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
			finder.ProcessTitles(db)
			Expect(models.Count(db, "pages")).To(Equal(6234))
		})
	})
})
