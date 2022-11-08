package name

import (
	"math"
	"time"

	"github.com/gnames/bhlindex/internal/ent/page"
	gnfout "github.com/gnames/gnfinder/ent/output"
)

// DetectedName holds information about a name-string returned by a
// name-finder.
type DetectedName struct {
	// ID is an identifier assigned to the occurrence by the database.
	ID uint `json:"id"`

	// PageID is the filename of the page without file extension.
	PageID string `json:"pageId" gorm:"type:varchar(255);index"`

	// ItemID is an identifier generated by database for the item.
	ItemID int `json:"itemId" gorm:"not null;index"`

	// Name is a normalized version of the detected name.
	Name string `json:"name" sql:"type:CHARACTER VARYING(255) COLLATE \"C\" NOT NULL" gorm:"index"`

	// NameVerbatim is a verbatim version of the detected name found on the page.
	NameVerbatim string `json:"nameVerbatim" sql:"type:CHARACTER VARYING(255) COLLATE \"C\" NOT NULL"`

	// AnnotNomen is a nomenclatural annotation found right after the name.
	AnnotNomen string `json:"annotNomen" gorm:"type:varchar(255)"`

	// AnnotNomenType is the normalized nomenclatural annotation.
	AnnotNomenType string `json:"annotNomenType" gorm:"type:varchar(255)"`

	// OffsetStart is the position the beginning of the detectedName.
	OffsetStart int `json:"offsetStart" gorm:"not null;default:0"`

	// OffsetEnd is the position the end of the detectedName.
	OffsetEnd int `json:"offsetEnd" gorm:"not null;default:0"`

	// EndsNextPage indicates if name ended on a different page than it started.
	EndsNextPage bool `json:"endsNextPage" gorm:"type:bool;default:false"`

	// OddsLog10 is a base 10 logarithm of odds generated by Bayes algorithm.
	OddsLog10 float64 `json:"oddsLog10" gorm:"type:float;not null;default:0"`

	// UpdatedAt is a timestamp when this record was generated.
	UpdatedAt time.Time `json:"updatedAt" sql:"type:timestamp without time zone"`
}

// VerifiedName holds information about verification results for a name-string.
type VerifiedName struct {
	// ID is the identifier automatically generated by the database.
	ID int `json:"id"`

	// Name is a normalized version of a detected name.
	Name string `json:"name" sql:"type:CHARACTER VARYING(255) COLLATE \"C\" NOT NULL" gorm:"index"`

	// Cardinality is the number of words in the simplest canonical form of a
	// detected name. For example `Aus` has cardinality 1, `Aus bus` has
	// cardinality 2, `Aus bus var. cus` has cardinality 3.
	Cardinality int `json:"cardinality" gorm:"not null;default:0"`

	// RecordID is the ID assigned by the DataSource to the name.
	RecordID string `json:"recordID" sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`

	// MatchType describes a resulting kind of a name-string match.
	// The following match types are possible:
	//
	// NoMatch - GNverifier did not find a match for the name-string.
	// Exact - Canonical form of a name matched exactly
	// PartialExact - Canonical form matched exactly after removal of some words.
	// Fuzzy - Canonical form matched, but with some differences.
	// PartialFuzzy - Canonical form matched with differences after removal of some words.
	// Virus - Name-string matched as a virus name.
	MatchType string `json:"matchType" gorm:"type:varchar(100)"`

	// EditDistance shows how much difference exists between name-string and a
	// match according to Levenshtein algorithm.
	EditDistance int `json:"editDistance" gorm:"not null;default:0"`

	// StemEditDistance shows how much difference exists between name-string and
	// a match according to Levenshtein algorithm.
	StemEditDistance int `json:"stemEditDistance" gorm:"not null;default:0"`

	// MatchedName provides the complete complete name-string.
	MatchedName string `json:"matchedName" sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`

	// MatchedCanonical provides canonical form of the matched name-string.
	MatchedCanonical string `json:"matchedCanonical" sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`

	// MatchedCardinality is the cardinality of matched name (see Cardinality
	// field for explanation).
	MatchedCardinality int `json:"matchedCardinality"`

	// CurrentName is the full currently accepted name of the match
	// provided by the DataSource.
	CurrentName string `json:"currentName" sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`

	// CurrentCanonical is a canonical form of the currently accepted name of
	// the match.
	CurrentCanonical string `json:"currentCanonical" sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`

	// CurrentCardinality is a cardinality of the currently accepted
	// name of the match.
	CurrentCardinality int `json:"currentCardinality"`

	// Classification contains a classification to the name provided by the
	// DataSource.
	Classification string `json:"classification" sql:"type:CHARACTER VARYING COLLATE \"C\""`

	// ClassificationRanks contains a ranks of each classification entry.
	ClassificationRanks string `json:"classificationRanks" sql:"type:CHARACTER VARYING COLLATE \"C\""`

	// ClassificationIDs are identifiers assigned to each classification entry.
	ClassificationIDs string `json:"classificationIds" sql:"type:CHARACTER VARYING COLLATE \"C\""`

	// DataSourceID is the ID of the data-source according to GNverifier.
	// The mapping of IDs to data-sources can be found at
	// https://verifier.globalnames.org/data_sources
	// site.
	DataSourceID int `json:"dataSourceID" gorm:"index"`

	// DataSourceTitle provides a title of the data-source that matched the
	// name-string.
	DataSourceTitle string `json:"dataSourceTitle" sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`

	// DataSourcesNumber is the number of dataSources that matched the name.
	DataSourcesNumber int `json:"dataSourcesNumber"`

	// Curation provides information about a level of curation according to
	// GNverifier. The following categories are supported:
	//
	// NotCurated -- None of data-sources that matched a name-string are marked as curated.
	// Curated -- Some data-sources with a match are marked as curated.
	// AutoCurated -- Some data-sources have automatic quality control, but not much human curation.
	Curation string `json:"curation"`

	// OddsLog10 is a logarithm with base 10 of odds that a detected string is
	// actually a scientific name according to a Naive Bayes algorithm.
	OddsLog10 float64 `json:"oddsLog10" gorm:"type:float;not null;default:0"`

	// SortOrder allows to compare the quality of name verification from many
	// data-sources from the best to the worst (0-N). The best result always has
	// the sort order 0. In case if only the best verification results are
	// returned, SortOrder is always equal 0.
	SortOrder int `json:"SortOrder" gorm:"not null, default 0"`

	// Occurrences is the number of times this name appeared in BHL texts.
	Occurrences int `json:"occurrences" gorm:"not null;default:0"`

	// Number of tries for verification process. If the number is larger than
	// 1, there were difficulties to get verification from the remote site.
	Retries int `json:"-" gorm:"not null;default:0"`

	// Error contains error that happened during verification. If this field
	// is empty then verification was completed successfully for the name-string.
	Error string `json:"error"`

	// UpdatedAt is the timestemp of the record.
	UpdatedAt time.Time `json:"updatedAt" sql:"type:timestamp without time zone"`
}

// UniqueName summarizes informaion after finding unique normalized names
// found in BHL texts.
type UniqueName struct {
	// ID is an identifier automatically generated by the database.
	ID int `json:"id" gorm:"primary_key"`

	// Name is a normalized version of detected name.
	Name string `json:"name" sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`

	// OddsLog10 is an average of odds results for a detected name.
	OddsLog10 float64 `json:"oddsLog10" gorm:"type:float;not null;default:0"`

	// Occurrences is the number of times a name-string was detected in BHL.
	Occurrences int `json:"occurrences" gorm:"not null;default:0"`
}

// New is a factory method for generating a record of a detected name.
func New(itemID int, p *page.Page, n gnfout.Name) DetectedName {
	var endsNextPage bool
	var end int
	start := n.OffsetStart - p.Offset
	if n.OffsetEnd < p.OffsetNext {
		end = n.OffsetEnd - p.Offset
	} else {
		end = n.OffsetEnd - p.OffsetNext
		endsNextPage = true
	}
	dn := DetectedName{
		PageID:         p.ID,
		ItemID:         itemID,
		Name:           n.Name,
		NameVerbatim:   n.Verbatim,
		OffsetStart:    start,
		OffsetEnd:      end,
		AnnotNomen:     n.AnnotNomen,
		AnnotNomenType: n.AnnotNomenType,
		EndsNextPage:   endsNextPage,
		OddsLog10:      math.Log10(n.Odds),
		UpdatedAt:      time.Now(),
	}
	return dn
}
