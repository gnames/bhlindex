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
		"PageId", "ItemId", "NameId",
		"DetectedName", "DetectedNameVerbatim", "OddsLog10",
		"Start", "End", "EndsNextPage", "Annotation",
	}
}

func (o OutputOccurrence) csvOutput(sep rune) string {
	var odds string
	if o.OddsLog10 > 0 {
		odds = strconv.FormatFloat(o.OddsLog10, 'f', 5, 64)
	}
	pageID := strconv.Itoa(o.PageID)
	itemID := strconv.Itoa(o.ItemID)
	start := strconv.Itoa(o.OffsetStart)
	end := strconv.Itoa(o.OffsetEnd)
	s := []string{
		pageID, itemID, o.NameID,
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
