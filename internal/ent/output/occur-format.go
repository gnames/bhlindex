package output

import (
	"strconv"

	"github.com/gnames/gnfmt"
)

func (o OutputOccurrence) Name() string {
	return "occurrences"
}

func (o OutputOccurrence) header() []string {
	return []string{
		"ItemBarcode", "PageBarcodeNum", "NameId",
		"DetectedName", "DetectedNameVerbatim", "OddsLog10",
		"Start", "End", "EndsNextPage", "Annotation",
	}
}

func (o OutputOccurrence) csvOutput(sep rune) string {
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

func (o OutputOccurrence) jsonOutput(pretty bool) string {
	enc := gnfmt.GNjson{Pretty: pretty}
	res, _ := enc.Encode(o)
	return string(res)
}
