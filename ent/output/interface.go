package output

import "context"

// Dumper interface contains methods for saving scientific names detected in
// Biodiversity Heritage Library as a flat CSV file on the file-system.
type Dumper interface {
	// DumpNames traverses database and outputs verified names in JSON, TSV or CSV format.
	DumpNames(context.Context, chan<- []OutputName, []int) error

	// DumpOccurrences traverses database and outputs names occurrences in JSON, TSV or CSV format.
	DumpOccurrences(context.Context, chan<- []OutputOccurrence, []int) error
}
