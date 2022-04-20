package output

// Output provides fields for data-dump of detected and verified names.
type Output struct {
	// ID of a name occurrence. It is automatically generated by database
	ID int `json:"id"`

	// NameID is an ID of a verified name. It is automatically generated by database.
	NameID int `json:"nameId"`

	// ItemBarcode is an Archive ID of an Item where a name appeared.
	ItemBarcode string `json:"itemBarcode"`

	// PageBarcodeNum is a number extracted from the page file name. The page filename
	// consists of this number and its Item's barcode.
	PageBarcodeNum int `json:"pageBarcodeNum"`

	// DetectedName is a normalized version of a detected name-string.
	DetectedName string `json:"detectedName"`

	// DetectedVerbatim is a detected name-string without normalization.
	// On rare occasions verbatim name will be truncated if it has too
	// much "junk" and exceeds the length of 225 characters.
	DetectedVerbatim string `json:"detectedVerbatim"`

	// OccurrencesTotal is the total number of occurrences of a particular name
	// in the BHL corpus.
	OccurrencesTotal int `json:"occurrencesTotal"`

	// OddsLog10 is a logarithm with base 10 of odds that a detected string is actually a
	// scientific name according to a Naive Bayes algorithm.
	OddsLog10 float64 `json:"oddsLog10"`

	// OffsetStart provides the number of UTF-8 characters from the page start to the
	// start of the name-string.
	OffsetStart int `json:"start"`

	// OffsetEnd provides the number of UTF-8 characters from the page start to the
	// end of the name-string.
	OffsetEnd int `json:"end"`

	// EndsNextPage is true when a name starts on one page and continues on the
	// next page.
	EndsNextPage bool `json:"endsNextPage"`

	// Cardinality is the number of words in the simplest canonical form of a name.
	// For example `Aus` has cardinality 1, `Aus bus` has cardinality 2,
	// `Aus bus var. cus` has cardinality 3.
	Cardinality int `json:"cardinality"`

	// MatchType describes a resulting kind of a name-string match.
	// The following match types are possible:
	//
	// NoMatch - GNverifier did not find a match for the name-string.
	// Exact - Canonical form of a name matched exactly
	// PartialExact - Canonical form matched exactly after removal of some words.
	// Fuzzy - Canonical form matched, but with some differences.
	// PartialFuzzy - Canonical form matched with differences after removal of some words.
	// Virus - Name-string matched as a virus name.
	MatchType string `json:"matchType"`

	// EditDistance shows how much difference exists between name-string and a match
	// according to Levenshtein algorithm.
	EditDistance int `json:"editDistance"`

	// MatchedCanonical provides canonical form of the matched name-string.
	MatchedCanonical string `json:"matchedCanonical"`

	// MatchedFullName provides the complete complete name-string.
	MatchedFullName string `json:"matchedFullName"`

	// MatchedCardinality is the cardinality of matched name (see Cardinality field
	// for explanation).
	MatchedCardinality int `json:"matchedCardinality"`

	// DataSourceID is the ID of the data-source according to GNverifier.
	// The mapping of IDs to data-sources can be found at
	// https://verifier.globalnames.org/data_sources
	DataSourceID int `json:"dataSourceID"`

	// DataSource provides a title of the data-source that matched the
	// name-string.
	DataSource string `json:"dataSource"`

	// Curation provides information about a level of curation according to
	// GNverifier. The following categories are supported:
	//
	// NotCurated -- None of data-sources that matched a name-string are marked as curated.
	// Curated -- Some data-sources with a match are marked as curated.
	// AutoCurated -- Some data-sources have automatic quality control, but not much human curation.
	Curation string `json:"Curation"`

	// VerifError contains error that happened during verification. If this field is empty
	// then verification was completed successfully for the name-string.
	VerifError string `json:"verificationError"`
}
