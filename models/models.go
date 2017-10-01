package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gnames/bhlindex"
	"github.com/gnames/gnfinder"
	"github.com/lib/pq"
)

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

func NewDetectedName(p Page, n gnfinder.Name) DetectedName {
	var endsNextPage bool
	var end int
	start := n.OffsetStart - p.Offset
	if n.OffsetEnd < p.OffsetNext {
		end = n.OffsetEnd - p.Offset
	} else {
		end = n.OffsetEnd - p.OffsetNext
		endsNextPage = true
	}
	dn := DetectedName{
		PageID:       p.ID,
		NameString:   n.ScientificName,
		OffsetStart:  start,
		OffsetEnd:    end,
		EndsNextPage: endsNextPage,
		UpdatedAt:    time.Now(),
	}
	return dn
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
