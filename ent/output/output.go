package output

type Output struct {
	ID                 int     `json:"id"`
	NameID             int     `json:"nameId"`
	PageBarcode        string  `json:"pageBarcode"`
	ItemBarcode        string  `json:"itemBarcode"`
	DetectedName       string  `json:"detectedName"`
	Occurrences        int     `json:"occurrences"`
	OddsLog10          float64 `json:"oddsLog10"`
	OffsetStart        int     `json:"start"`
	OffsetEnd          int     `json:"end"`
	EndsNextPage       bool    `json:"endsNextPage"`
	Cardinality        int     `json:"cardinality"`
	MatchType          string  `json:"matchType"`
	EditDistance       int     `json:"editDistance"`
	MatchedCanonical   string  `json:"matchedCanonical"`
	MatchedFullName    string  `json:"matchedFullName"`
	MatchedCardinality int     `json:"matchedCardinality"`
	DataSourceID       int     `json:"dataSourceID"`
	DataSource         string  `json:"dataSource"`
	Curation           string  `json:"Curation"`
	VerifError         string  `json:"verificationError"`
}
