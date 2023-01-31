package output

import (
	"strconv"

	"github.com/gnames/gnfmt"
)

func (o OutputOccurrenceShort) Name() string {
	return "occurrences"
}

func (o OutputOccurrenceShort) header() []string {
	return []string{"NameId", "DetectedNameVerbatim", "PageId"}
}

func (o OutputOccurrenceShort) csvOutput(sep rune) string {
	pageID := strconv.Itoa(o.PageID)
	s := []string{o.NameID, o.DetectedVerbatim, pageID}

	return gnfmt.ToCSV(s, sep)
}

type occurShort struct {
	NameID           string `json:"nameId"`
	DetectedVerbatim string `json:"detectedVerbatim"`
	PageID           int    `json:"pageId"`
}

func (o OutputOccurrenceShort) jsonOutput(pretty bool) string {
	oShort := occurShort{
		NameID:           o.NameID,
		DetectedVerbatim: o.DetectedVerbatim,
		PageID:           o.PageID,
	}
	enc := gnfmt.GNjson{Pretty: pretty}
	res, _ := enc.Encode(oShort)
	return string(res)
}

func (o OutputOccurrenceShort) PageNameIDs() (string, string) {
	return strconv.Itoa(o.PageID), o.NameID
}
