package output

import (
	"strconv"

	"github.com/gnames/gnfmt"
)

func (o *Output) Format(f gnfmt.Format) string {
	switch f {
	case gnfmt.TSV:
		return o.csvOutput('\t')
	case gnfmt.CompactJSON:
		return o.jsonOutput(false)
	default:
		return o.csvOutput(',')
	}
}

// CSVHeader returns the header string for CSV output format.
func CSVHeader(f gnfmt.Format) string {
	sep := rune(',')
	if f == gnfmt.TSV {
		sep = rune('\t')
	}

	res := []string{
		"Id", "NameId", "ItemBarcode", "PageBarcodeNum", "DetectedName",
		"Occurrences", "OddsLog10", "Start", "End", "EndsNextPage", "Cardinality",
		"MatchType", "EditDistance", "MatchedCanonical", "MatchedFullName",
		"MatchedCardinality", "DataSourceID", "DataSource", "InCuratedSources",
		"VerifError",
	}
	return gnfmt.ToCSV(res, sep)
}

func (o *Output) csvOutput(sep rune) string {
	var odds string
	if o.OddsLog10 > 0 {
		odds = strconv.FormatFloat(o.OddsLog10, 'f', 5, 64)
	}

	start := strconv.Itoa(o.OffsetStart)
	end := strconv.Itoa(o.OffsetEnd)
	barcodeNum := strconv.Itoa(o.PageBarcodeNum)
	s := []string{
		strconv.Itoa(o.ID), strconv.Itoa(o.NameID), o.ItemBarcode, barcodeNum,
		o.DetectedName, strconv.Itoa(o.Occurrences), odds, start, end,
		strconv.FormatBool(o.EndsNextPage), strconv.Itoa(o.Cardinality),
		o.MatchType, strconv.Itoa(o.EditDistance), o.MatchedCanonical,
		o.MatchedFullName, strconv.Itoa(o.MatchedCardinality),
		strconv.Itoa(o.DataSourceID), o.DataSource, o.Curation, o.VerifError,
	}

	return gnfmt.ToCSV(s, sep)
}

func (o *Output) jsonOutput(pretty bool) string {
	enc := gnfmt.GNjson{Pretty: pretty}
	res, _ := enc.Encode(o)
	return string(res)
}
