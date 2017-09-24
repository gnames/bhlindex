package main_test

import (
	"database/sql"

	"github.com/GlobalNamesArchitecture/bhlindex/util"
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

var _ = BeforeSuite(func() {
	var err error
	db, err = util.DbInit()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := db.Close()
	Expect(err).NotTo(HaveOccurred())
})
