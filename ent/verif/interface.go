package verif

// VerifierBHL interface provides reconciliation and resolution of scientific
// names. Reconciliation matches name-string to all found lexical variants of
// the string. Resolution uses information in taxonomic databases such as
// Catalogue of Life to determing currently accepted name according to the
// database.
type VerifierBHL interface {
	// Verify method organizes names detected in BHL, verifies them against
	// many scientific name databases, and returns reconciliation/resolution
	// results.
	Verify() error
}