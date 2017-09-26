package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/GlobalNamesArchitecture/bhlindex"
	"github.com/lib/pq"
)

// NameFinder interface determines behavior of scientific name finders
type NameFinder interface {
	FindNames(title *Title) ([]DetectedName, error)
}

// DetectedName helds information about a name-string returned by a
// name-finder.
type DetectedName struct {
	PageID       string
	NameString   string
	NameID       int
	OffsetStart  int
	OffsetEnd    int
	EndsNextPage bool
	UpdatedAt    time.Time
}

// Count returns number of rows in a table
func Count(db *sql.DB, table string) int {
	var count int
	q := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, pq.QuoteIdentifier(table))
	err := db.QueryRow(q).Scan(&count)
	bhlindex.Check(err)
	return count
}

// Truncate removes all the rows from a table
func Truncate(db *sql.DB, table string) {
	q := fmt.Sprintf("TRUNCATE TABLE %s", pq.QuoteIdentifier(table))
	_, err := db.Query(q)
	bhlindex.Check(err)
}
