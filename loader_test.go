package bhlindex_test

import (
	"github.com/GlobalNamesArchitecture/bhlindex/loader"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Loader", func() {

	Describe("Path", func() {
		It("returns ok", func() {
			c := make(chan string)
			count := 0
			go loader.Path(c)
			for {
				if _, mode := <-c; mode {
					count += 1
				} else {
					break
				}
			}
			Expect(count).To(Equal(6263))
		})
	})

})
