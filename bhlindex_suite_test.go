package main_test

import (
	"database/sql"
	"fmt"

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
	env := util.EnvVars()
	params := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		env.DbUser, env.DbPass, env.DbHost, env.Db)
	var err error
	db, err = sql.Open("postgres", params)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := db.Close()
	Expect(err).NotTo(HaveOccurred())
})
