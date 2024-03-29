package verif

// VerifierBHL interface describes methods for reconciliation and resolution of
// scientific names. Reconciliation matches a name-string to all found lexical
// variants of the string. Resolution uses information in taxonomic databases
// such as Catalogue of Life to determine currently accepted name according to
// the database.
type VerifierBHL interface {
	// Reset removes saved data from the previous verification.
	Reset() error

	// Verify method organizes names detected in BHL, verifies them against
	// many scientific name databases, and returns reconciliation/resolution
	// results.
	Verify() error

	// ExtractUniqueNames runs after name detection is finished. It goes through
	// the detected names and saves the unique list of all names, their odds,
	// and the number of their occurrences.
	ExtractUniqueNames() error

	// CalcOddsVerif calculates the relationship between Odds and Verifications.
	CalcOddsVerif() error
}
