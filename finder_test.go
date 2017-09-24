package main_test

import (
	"github.com/GlobalNamesArchitecture/bhlindex/finder"
	"github.com/GlobalNamesArchitecture/bhlindex/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Finder", func() {

	Describe("FindNames", func() {
		It("Finds names", func() {
			finder.FindNames(db)
			Expect(models.Count(db, "titles")).To(BeNumerically(">", 5000))
		})
	})

})
