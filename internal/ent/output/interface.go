package output

import (
	"context"
)

// Dumper interface contains methods for saving scientific names detected in
// Biodiversity Heritage Library as a flat CSV file on the file-system.
type Dumper interface {
	// DumpNames traverses database and outputs verified names in JSON, TSV or CSV format.
	DumpNames(context.Context, chan<- []Output, []int) error

	// DumpOccurrences traverses database and outputs names occurrences in JSON, TSV or CSV format.
	DumpOccurrences(context.Context, chan<- []Output, []int, bool) error

	// DumpOddsVerification gets result of mapping between Odds values and the percentage of
	// successful verifications.
	DumpOddsVerification() ([]Output, error)
}

// Output interface provides generic functions outputs of verified names, and
// names occurrences.
type Output interface {
	// Name returns the name of output file.
	Name() string

	// header returns list of fields for CSF or TSV format.
	header() []string

	// csvOutput provides a method to generate a CSV ot TSV row.
	csvOutput(rune) string

	// jsonOutput provides a method to generate a JSON output.
	jsonOutput(bool) string

	// PageNameIDs returns PageID and Name, if available, empty int and string
	// otherwize
	PageNameIDs() (string, string)
}
