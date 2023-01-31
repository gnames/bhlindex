package output

import (
	"strconv"

	"github.com/gnames/gnfmt"
)

func (on OutputName) Name() string {
	return "names"
}

func (on OutputName) header() []string {
	return []string{
		"NameID", "DetectedName", "Cardinality", "OccurrencesNumber", "OddsLog10",
		"MatchType", "MatchSortOrder", "EditDistance", "StemEditDistance",
		"MatchedCanonical", "MatchedFullName", "MatchedCardinality",
		"CurrentCanonical", "CurrentFullName", "CurrentCardinality",
		"Classification", "ClassificationRanks", "ClassificationIDs", "RecordID",
		"DataSourceID", "DataSource", "DataSourcesNumber", "Curation", "Error",
	}
}

func (on OutputName) csvOutput(sep rune) string {
	card := strconv.Itoa(on.Cardinality)
	occNum := strconv.Itoa(on.OccurrencesNumber)

	var odds string
	if on.OddsLog10 > 0 {
		odds = strconv.FormatFloat(on.OddsLog10, 'f', 5, 64)
	}

	eDist := strconv.Itoa(on.EditDistance)
	stEDist := strconv.Itoa(on.StemEditDistance)
	matchCard := strconv.Itoa(on.MatchedCardinality)
	currCard := strconv.Itoa(on.CurrentCardinality)
	dsID := strconv.Itoa(on.DataSourceID)
	dsNum := strconv.Itoa(on.DataSourcesNumber)
	sortOrder := strconv.Itoa(on.SortOrder)

	s := []string{
		on.NameID, on.DetectedName, card, occNum, odds, on.MatchType, sortOrder,
		eDist, stEDist, on.MatchedCanonical, on.MatchedFullName, matchCard,
		on.CurrentCanonical, on.CurrentFullName, currCard, on.Classification,
		on.ClassificationRanks, on.ClassificationIDs, on.RecordID, dsID,
		on.DataSource, dsNum, on.Curation, on.VerifError}

	return gnfmt.ToCSV(s, sep)
}

func (on OutputName) jsonOutput(pretty bool) string {
	enc := gnfmt.GNjson{Pretty: pretty}
	res, _ := enc.Encode(on)
	return string(res)
}

func (on OutputName) PageNameIDs() (string, string) {
	return "", ""
}
