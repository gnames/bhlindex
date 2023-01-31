package output

import (
	"strconv"

	"github.com/gnames/gnfmt"
)

func (on OutputNameShort) Name() string {
	return "names"
}

func (on OutputNameShort) header() []string {
	return []string{
		"NameID", "DetectedName", "OccurrencesNumber", "MatchedCanonical",
		"MatchedFullName", "RecordID", "DataSourceID", "MatchType",
		"MatchSortOrder", "OddsLog10", "Curation", "Error",
	}
}

func (on OutputNameShort) csvOutput(sep rune) string {
	var odds string
	if on.OddsLog10 > 0 {
		odds = strconv.FormatFloat(on.OddsLog10, 'f', 5, 64)
	}

	occNum := strconv.Itoa(on.OccurrencesNumber)
	dsID := strconv.Itoa(on.DataSourceID)
	sortOrder := strconv.Itoa(on.SortOrder)

	s := []string{
		on.NameID, on.DetectedName, occNum, on.MatchedCanonical,
		on.MatchedFullName, on.RecordID, dsID, on.MatchType, sortOrder,
		odds, on.Curation, on.VerifError,
	}

	return gnfmt.ToCSV(s, sep)
}

type nameShort struct {
	NameID            string  `json:"nameID"`
	DetectedName      string  `json:"detectedName"`
	OccurrencesNumber int     `json:"occurrencesNumber"`
	MatchedCanonical  string  `json:"matchedCanonical"`
	MatchedFullName   string  `json:"matchedFullName"`
	RecordID          string  `json:"recordID"`
	DataSourceID      int     `json:"dataSourceId"`
	MatchType         string  `json:"matchType"`
	SortOrder         int     `json:"sortOrder"`
	OddsLog10         float64 `json:"oddsLog10"`
	Curation          string  `json:"curation"`
	VerifError        string  `json:"verificationError"`
}

func (on OutputNameShort) jsonOutput(pretty bool) string {
	out := nameShort{
		NameID:            on.NameID,
		DetectedName:      on.DetectedName,
		OccurrencesNumber: on.OccurrencesNumber,
		MatchedCanonical:  on.MatchedCanonical,
		MatchedFullName:   on.MatchedFullName,
		RecordID:          on.RecordID,
		DataSourceID:      on.DataSourceID,
		MatchType:         on.MatchType,
		SortOrder:         on.SortOrder,
		OddsLog10:         on.OddsLog10,
		Curation:          on.Curation,
		VerifError:        on.VerifError,
	}
	enc := gnfmt.GNjson{Pretty: pretty}
	res, _ := enc.Encode(out)
	return string(res)
}

func (on OutputNameShort) PageNameIDs() (string, string) {
	return "", ""
}
