package output

import (
	"strconv"

	"github.com/gnames/gnfmt"
)

func (o OutputOddsVerification) Name() string {
	return "odds-verification"
}

func (o OutputOddsVerification) header() []string {
	return []string{"OddsLog10", "NamesNum", "VerifPercent"}
}

func (o OutputOddsVerification) csvOutput(sep rune) string {
	odds := strconv.Itoa(o.OddsLog10)
	namesNum := strconv.Itoa(o.NamesNum)
	verifPercent := strconv.FormatFloat(o.VerifPercent, 'f', 5, 64)
	s := []string{odds, namesNum, verifPercent}
	return gnfmt.ToCSV(s, sep)
}

func (o OutputOddsVerification) jsonOutput(pretty bool) string {
	enc := gnfmt.GNjson{Pretty: pretty}
	res, _ := enc.Encode(o)
	return string(res)
}
