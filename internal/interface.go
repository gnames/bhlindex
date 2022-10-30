package bhlindex

import (
	"github.com/gnames/bhlindex/internal/config"
	"github.com/gnames/bhlindex/internal/ent/finder"
	"github.com/gnames/bhlindex/internal/ent/loader"
	"github.com/gnames/bhlindex/internal/ent/output"
	"github.com/gnames/bhlindex/internal/ent/verif"
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

	// DumpPages creates output with pages info in CSV, TSV, or JSON formats.
	DumpPages(output.Dumper) error

	// GetVersion outputs the version of BHLindex.
	GetVersion() gnvers.Version

	// GetConfig returns an instance of configuration fields.
	GetConfig() config.Config
}
