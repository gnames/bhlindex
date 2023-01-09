package output

import (
	"strconv"

	"github.com/gnames/gnfmt"
)

func (o OutputOccurrenceShort) Name() string {
	return "occurrences"
}

func (o OutputOccurrenceShort) header() []string {
	return []string{"NameId", "PageId"}
}

func (o OutputOccurrenceShort) csvOutput(sep rune) string {
	pageID := strconv.Itoa(o.PageID)
	s := []string{o.NameID, pageID}

	return gnfmt.ToCSV(s, sep)
}

type occurShort struct {
	NameID string `json:"nameId"`
	PageID int    `json:"pageId"`
}

func (o OutputOccurrenceShort) jsonOutput(pretty bool) string {
	oShort := occurShort{NameID: o.NameID, PageID: o.PageID}
	enc := gnfmt.GNjson{Pretty: pretty}
	res, _ := enc.Encode(oShort)
	return string(res)
}
