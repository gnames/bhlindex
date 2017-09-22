package main_test

import (
	"github.com/GlobalNamesArchitecture/bhlindex/finder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Finder", func() {

	Describe("FindNames", func() {
		It("Finds names", func() {
			res := finder.FindNames()
			Expect(res).To(Equal("ok"))
		})
	})

})
