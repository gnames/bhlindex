package finder

import (
	"database/sql"

	"github.com/gnames/bhlindex"
)

type verificationStatus int

const (
	Init verificationStatus = iota
	AllNamesHarvested
	AllNamesVerified
	AllErrorsProcessed
)

type metaData struct {
	Status        verificationStatus
	StaleErrTries int
}

func readStatus(db *sql.DB) metaData {
	q := `SELECT status, stale_errors_tries FROM metadata`
	var status int
	var errNum int

	err := db.QueryRow(q).Scan(&status, &errNum)
	if err == sql.ErrNoRows {
		return initStatus(db)
	}
	bhlindex.Check(err)
	return metaData{Status: verificationStatus(status), StaleErrTries: errNum}
}

func initStatus(db *sql.DB) metaData {
	q := `INSERT INTO metadata (status, stale_errors_tries) VALUES (0, 0)`
	stmt, err := db.Prepare(q)
	bhlindex.Check(err)
	_, err = stmt.Exec()
	bhlindex.Check(err)

	err = stmt.Close()
	bhlindex.Check(err)

	return metaData{}
}

func updateStatus(db *sql.DB, s *metaData) {
	readStatus(db) //to init status if does not exist
	status := int(s.Status)
	errNum := s.StaleErrTries
	q := `UPDATE metadata set status=$1, stale_errors_tries=$2`
	stmt, err := db.Prepare(q)
	bhlindex.Check(err)
	_, err = stmt.Exec(status, errNum)
	bhlindex.Check(err)
	err = stmt.Close()
	bhlindex.Check(err)
}
