package verif

// VerifierBHL interface provides reconciliation and resolution of scientific
// names. Reconciliation matches name-string to all found lexical variants of
// the string. Resolution uses information in taxonomic databases such as
// Catalogue of Life to determing currently accepted name according to the
// database.
type VerifierBHL interface {
	// Reset removes saved data from the previous verification.
	Reset() error

	// Verify method organizes names detected in BHL, verifies them against
	// many scientific name databases, and returns reconciliation/resolution
	// results.
	Verify() error

	// ExtractUniqueNames runs after name detection is finished. It goes through
	// the detected names and saves the unique list off all names, their odds,
	// and the number of their occurrences.
	ExtractUniqueNames() error
}
