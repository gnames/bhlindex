package output

import "github.com/gnames/gnfmt"

// OutputName provides fields for data-dump of unique name-strings. The data
// also contains reconciliation and resolution data according to the best
// matches do a variety of data-sources.
type OutputName struct {
	// NameID is an UUID v5 of the name. It is derived from Name.
	NameID string `json:"nameId"`

	// DetectedName is a normalized version of a detected name-string.
	DetectedName string `json:"detectedName"`

	// Cardinality is the number of words in the simplest canonical form of a
	// detected name. For example `Aus` has cardinality 1, `Aus bus` has
	// cardinality 2, `Aus bus var. cus` has cardinality 3.
	Cardinality int `json:"cardinality"`

	// OccurrencesNumber is the total number of occurrences of a particular name
	// in the BHL corpus.
	OccurrencesNumber int `json:"occurrencesNumber"`

	// OddsLog10 is a logarithm with base 10 of odds that a detected string is
	// actually a scientific name according to a Naive Bayes algorithm.
	OddsLog10 float64 `json:"oddsLog10"`

	// SortOrder is used for sorting multiple matches from the same name.
	// The best results always has SortOrder = 0.
	SortOrder int `json:"matchSortOrder"`

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

	// EditDistance shows how much difference exists between name-string and a
	// match according to Levenshtein algorithm.
	EditDistance int `json:"editDistance"`

	// StemEditDistance shows how much difference exists between name-string and
	// a match according to Levenshtein algorithm.
	StemEditDistance int `json:"stemEditDistance"`

	// MatchedCanonical provides canonical form of the matched name-string.
	MatchedCanonical string `json:"matchedCanonical"`

	// MatchedFullName provides the complete complete name-string.
	MatchedFullName string `json:"matchedFullName"`

	// MatchedCardinality is the cardinality of matched name (see Cardinality
	// field for explanation).
	MatchedCardinality int `json:"matchedCardinality"`

	// CurrentCanonical is a canonical form of the currently accepted name of
	// the match.
	CurrentCanonical string `json:"currentCanonical"`

	// CurrentFullName is the full currently accepted name of the match
	// provided by the DataSource.
	CurrentFullName string `json:"currentFullName"`

	// CurrentCardinality is a cardinality of the currently accepted
	// name of the match.
	CurrentCardinality int `json:"currentCardinality"`

	// Classification contains a classification to the name provided by the
	// DataSource.
	Classification string `json:"classification"`

	// ClassificationRanks provide data about ranks used in the classificaiton.
	ClassificationRanks string `json:"classificationRanks"`

	// ClassificationIDs provides data about IDs a DataSource assigns to
	// taxons in the classification.
	ClassificationIDs string `json:"classificationIds"`

	// RecordID is the ID assigned by the DataSource to the name.
	RecordID string `json:"recordId"`

	// DataSourceID is the ID of the data-source according to GNverifier.
	// The mapping of IDs to data-sources can be found at
	// https://verifier.globalnames.org/data_sources
	// site.
	DataSourceID int `json:"dataSourceId"`

	// DataSource provides a title of the data-source that matched the
	// name-string.
	DataSource string `json:"dataSource"`

	// DataSourcesNumber is the number of dataSources that matched the name.
	DataSourcesNumber int `json:"dataSourcesNumber"`

	// Curation provides information about a level of curation according to
	// GNverifier. The following categories are supported:
	//
	// NotCurated -- None of data-sources that matched a name-string are marked as curated.
	// Curated -- Some data-sources with a match are marked as curated.
	// AutoCurated -- Some data-sources have automatic quality control, but not much human curation.
	Curation string `json:"Curation"`

	// VerifError contains error that happened during verification. If this field
	// is empty then verification was completed successfully for the name-string.
	VerifError string `json:"verificationError"`
}

type OutputNameShort struct {
	OutputName
}

// OutputOccurrence provides fields for data-dump of detected names.
type OutputOccurrence struct {
	// NameID is an UUID v5 of the name. It is derived from DetectedName.
	NameID string `json:"nameId"`

	// PageID is a number extracted from the page file name. The page
	// filename consists of this number and its Item's barcode.
	PageID int `json:"pageId"`

	// ItemID is an Archive ID of an Item where a name appeared.
	ItemID int `json:"itemId"`

	// DetectedName is a normalized version of a detected name-string.
	DetectedName string `json:"detectedName"`

	// DetectedVerbatim is a detected name-string without normalization.
	// On rare occasions verbatim name will be truncated if it has too
	// much "junk" and exceeds the length of 225 characters.
	DetectedVerbatim string `json:"detectedVerbatim"`

	// OddsLog10 is a logarithm with base 10 of odds that a detected string is
	// actually a scientific name according to a Naive Bayes algorithm.
	OddsLog10 float64 `json:"oddsLog10"`

	// OffsetStart provides the number of UTF-8 characters from the page start to
	// the start of the name-string.
	OffsetStart int `json:"start"`

	// OffsetEnd provides the number of UTF-8 characters from the page start to
	// the end of the name-string.
	OffsetEnd int `json:"end"`

	// EndsNextPage is true when a name starts on one page and continues on the
	// next page.
	EndsNextPage bool `json:"endsNextPage"`

	// Annotation is a normalized annotation of `sp. nov.`, `subsp. nov.` etc.
	// that was located after the DetectedName.
	Annotation string `json:"annotNomen"`
}

type OutputOccurrenceShort struct {
	OutputOccurrence
}

type OutputOddsVerification struct {
	OddsLog10    int     `json:"oddsLog10"`
	NamesNum     int     `json:"namesNum"`
	VerifPercent float64 `json:"verifPercent"`
}

// CSVHeader takes any object that implements Output interface
// and generates TSV or CSV header.
func CSVHeader[O Output](o O, f gnfmt.Format) string {
	sep := rune(',')
	if f == gnfmt.TSV {
		sep = rune('\t')
	}
	return gnfmt.ToCSV(o.header(), sep)
}

func Format[O Output](o O, f gnfmt.Format) string {
	switch f {
	case gnfmt.TSV:
		return o.csvOutput('\t')
	case gnfmt.CompactJSON:
		return o.jsonOutput(false)
	default:
		return o.csvOutput(',')
	}
}
