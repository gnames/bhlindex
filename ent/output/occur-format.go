package output

import (
	"strconv"

	"github.com/gnames/gnfmt"
)

func (o *OutputOccurrence) Format(f gnfmt.Format) string {
	switch f {
	case gnfmt.TSV:
		return o.csvOutput('\t')
	case gnfmt.CompactJSON:
		return o.jsonOutput(false)
	default:
		return o.csvOutput(',')
	}
}

// CSVHeaderOccur returns the header string for CSV output format.
func CSVHeaderOccur(f gnfmt.Format) string {
	sep := rune(',')
	if f == gnfmt.TSV {
		sep = rune('\t')
	}

	res := []string{
		"ItemBarcode", "PageBarcodeNum", "NameId",
		"DetectedName", "DetectedNameVerbatim", "OddsLog10",
		"Start", "End", "EndsNextPage", "Annotation",
	}
	return gnfmt.ToCSV(res, sep)
}

func (o *OutputOccurrence) csvOutput(sep rune) string {
	var odds string
	if o.OddsLog10 > 0 {
		odds = strconv.FormatFloat(o.OddsLog10, 'f', 5, 64)
	}

	start := strconv.Itoa(o.OffsetStart)
	end := strconv.Itoa(o.OffsetEnd)
	barcodeNum := strconv.Itoa(o.PageBarcodeNum)
	s := []string{
		o.ItemBarcode, barcodeNum, o.NameID,
		o.DetectedName, o.DetectedVerbatim, odds,
		start, end, strconv.FormatBool(o.EndsNextPage),
		o.Annotation,
	}

	return gnfmt.ToCSV(s, sep)
}

func (o *OutputOccurrence) jsonOutput(pretty bool) string {
	enc := gnfmt.GNjson{Pretty: pretty}
	res, _ := enc.Encode(o)
	return string(res)
}
