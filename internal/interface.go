package bhlindex

import (
	"github.com/gnames/bhlindex/internal/config"
	"github.com/gnames/bhlindex/internal/ent/finder"
	"github.com/gnames/bhlindex/internal/ent/loader"
	"github.com/gnames/bhlindex/internal/ent/output"
	"github.com/gnames/bhlindex/internal/ent/verif"
	"github.com/gnames/gnlib/ent/gnvers"
)

// BHLindex defines main functionality of BHLindex project.
type BHLindex interface {
	// FindNames traverses BHL corpus directory structure, assembling texts,
	// detecting names, saving the resulting data to storage.
	FindNames(loader.Loader, finder.Finder) error

	// Verify names runs verification on unique detected names and saves the
	// results to storage.
	VerifyNames(verif.VerifierBHL) error

	// CalcOddsVerif calculates relationship between odds and the percentage
	// of successful verifications.
	CalcOddsVerif(verif.VerifierBHL) error

	// DumpOccurrences exports data about detected names and their position
	// on BHL pages in CSV, TSV, or JSON formats.
	DumpOccurrences(output.Dumper) error

	// DumpNames exports results of names verification in CSV,
	// TSV, or JSON formats.
	DumpNames(output.Dumper) error

	// DumpOddsVerification exports mapping of Odds values to verification
	// percentage
	DumpOddsVerification(output.Dumper) error

	// GetVersion outputs the version of BHLindex.
	GetVersion() gnvers.Version

	// GetConfig returns an instance of configuration fields.
	GetConfig() config.Config
}
