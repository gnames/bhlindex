package bhlindex

import (
	"github.com/gnames/bhlindex/config"
	"github.com/gnames/bhlindex/ent/finder"
	"github.com/gnames/bhlindex/ent/loader"
	"github.com/gnames/bhlindex/ent/output"
	"github.com/gnames/bhlindex/ent/verif"
	"github.com/gnames/gnlib/ent/gnvers"
)

// BHLindex us the main usecase interface that defines functionality of BHLindex
type BHLindex interface {
	// FindNames traverses BHL corpus directory structure, assembling texts,
	// detecting names, saving data to storage.
	FindNames(loader.Loader, finder.Finder) error

	// Verify names runs verification on unique detected names and saves the
	// results to a local storage.
	VerifyNames(verif.VerifierBHL) error

	// DumpOccurrences creates output with detected names in CSV,
	// TSV, or JSON formats.
	DumpOccurrences(output.Dumper) error

	// DumpNames creates output with verified names in CSV,
	// TSV, or JSON formats.
	DumpNames(output.Dumper) error

	// GetVersion outputs the version of BHLindex.
	GetVersion() gnvers.Version

	// GetConfig returns an instance of configuration fields.
	GetConfig() config.Config
}
