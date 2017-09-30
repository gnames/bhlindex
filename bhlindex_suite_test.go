package bhlindex_test

import (
	"database/sql"

	"github.com/GlobalNamesArchitecture/bhlindex"
	"github.com/GlobalNamesArchitecture/gnfinder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	_ "github.com/lib/pq"
)

func TestBhlindex(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bhlindex Suite")
}

var db *sql.DB
var dict *gnfinder.Dictionary

var _ = BeforeSuite(func() {
	var err error
	db, err = bhlindex.DbInit()
	Expect(err).NotTo(HaveOccurred())
	d := gnfinder.LoadDictionary()
	dict = &d
})

var _ = AfterSuite(func() {
	err := db.Close()
	Expect(err).NotTo(HaveOccurred())
})
