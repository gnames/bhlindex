package output

import (
	"strconv"

	"github.com/gnames/gnfmt"
)

func (on *OutputName) Format(f gnfmt.Format) string {
	switch f {
	case gnfmt.TSV:
		return on.csvOutput('\t')
	case gnfmt.CompactJSON:
		return on.jsonOutput(false)
	default:
		return on.csvOutput(',')
	}
}

// CSVHeaderName returns the header string for CSV output format.
func CSVHeaderName(f gnfmt.Format) string {
	sep := rune(',')
	if f == gnfmt.TSV {
		sep = rune('\t')
	}

	res := []string{
		"NameID", "DetectedName", "Cardinality", "OccurrencesNumber", "OddsLog10",
		"MatchType", "EditDistance", "StemEditDistance", "MatchedCanonical",
		"MatchedFullName", "MatchedCardinality", "CurrentCanonical",
		"CurrentFullName", "CurrentCardinality", "Classification", "RecordID",
		"DataSourceID", "DataSource", "DataSourcesNumber", "Curation", "Error",
	}
	return gnfmt.ToCSV(res, sep)
}

func (on *OutputName) csvOutput(sep rune) string {
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

	s := []string{
		on.NameID, on.DetectedName, card, occNum, odds, on.MatchType,
		eDist, stEDist, on.MatchedCanonical, on.MatchedFullName, matchCard,
		on.CurrentCanonical, on.CurrentFullName, currCard, on.Classification,
		on.RecordID, dsID, on.DataSource, dsNum, on.Curation, on.VerifError,
	}

	return gnfmt.ToCSV(s, sep)
}

func (on *OutputName) jsonOutput(pretty bool) string {
	enc := gnfmt.GNjson{Pretty: pretty}
	res, _ := enc.Encode(on)
	return string(res)
}
