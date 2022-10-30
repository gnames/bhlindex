package output

import (
	"strconv"

	"github.com/gnames/gnfmt"
)

func (op OutputPage) Name() string {
	return "pages"
}

func (op OutputPage) header() []string {
	return []string{
		"ItemBarcode", "PageBarcodeNum",
	}
}

func (op OutputPage) csvOutput(sep rune) string {
	s := []string{
		op.ItemBarcode, strconv.Itoa(op.PageBarcodeNum),
	}

	return gnfmt.ToCSV(s, sep)
}

func (op OutputPage) jsonOutput(pretty bool) string {
	enc := gnfmt.GNjson{Pretty: pretty}
	res, _ := enc.Encode(op)
	return string(res)
}
