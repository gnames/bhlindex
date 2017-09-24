package bhlindex_test

import (
	"errors"

	"github.com/GlobalNamesArchitecture/bhlindex"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("bhlindex", func() {
	Describe("Check()", func() {
		It("ignores `nil` errors", func() {
			err := error(nil)
			a := "one"
			bhlindex.Check(err)
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
			bhlindex.Check(err)
		})
	})

	Describe("EnvVars", func() {
		env := bhlindex.EnvVars()
		It("Returns envaronment variables", func() {
			Expect(env.BHLDir).To(Equal("./testdata/"))
			Expect(env.Db).To(Equal("bhlindex"))
			Expect(env.DbUser).To(Equal("postgres"))
			Expect(env.DbPass).To(Equal(""))
			Expect(env.DbHost).To(Equal("pg"))
		})
	})
})
