package output

import "context"

// Dumper interface contains methods for saving scientific names detected in
// Biodiversity Heritage Library as a flat CSV file on the file-system.
type Dumper interface {
	// Dump traverses database and outputs its data in JSON, TSV or CSV format.
	Dump(context.Context, chan<- []Output) error
}
