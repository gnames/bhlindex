package bhlindex_test

import (
	"database/sql"

	"github.com/gnames/bhlindex"
	"github.com/gnames/bhlindex/models"
	dictionary "github.com/gnames/gnfinder/dict"
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
var dict *dictionary.Dictionary

var _ = BeforeSuite(func() {
	var err error
	db, err = bhlindex.DbInit()
	Expect(err).NotTo(HaveOccurred())
	d := dictionary.LoadDictionary()
	dict = &d
})

var _ = AfterSuite(func() {
	err := db.Close()
	Expect(err).NotTo(HaveOccurred())
})

func truncateAll() {
	models.Truncate(db, "titles")
	models.Truncate(db, "pages")
	models.Truncate(db, "name_strings")
	models.Truncate(db, "page_name_strings")
	models.Truncate(db, "name_statuses")
	models.Truncate(db, "metadata")
}
