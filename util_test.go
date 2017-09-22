package main_test

import (
	"errors"

	"github.com/GlobalNamesArchitecture/bhlindex/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Util", func() {
	Describe("Check()", func() {
		It("ignores `nil` errors", func() {
			err := error(nil)
			a := "one"
			util.Check(err)
			Expect(a).To(Equal("one"))
		})

		It("panics if err is not `nil`", func() {
			defer func() {
				if r := recover(); r != nil {
					e := r.(error)
					Expect(e).To(Equal(errors.New("My error")))
				}
			}()
			err := errors.New("My error")
			util.Check(err)
		})
	})
})
