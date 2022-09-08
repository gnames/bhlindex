package output

import (
	"strconv"

	"github.com/gnames/gnfmt"
)

func (op *OutputPage) Format(f gnfmt.Format) string {
	switch f {
	case gnfmt.TSV:
		return op.csvOutput('\t')
	case gnfmt.CompactJSON:
		return op.jsonOutput(false)
	default:
		return op.csvOutput(',')
	}
}

// CSVHeaderPage returns the header string for CSV output format.
func CSVHeaderPage(f gnfmt.Format) string {
	sep := rune(',')
	if f == gnfmt.TSV {
		sep = rune('\t')
	}

	res := []string{
		"ItemBarcode", "PageBarcodeNum",
	}
	return gnfmt.ToCSV(res, sep)
}

func (op *OutputPage) csvOutput(sep rune) string {
	s := []string{
		op.ItemBarcode, strconv.Itoa(op.PageBarcodeNum),
	}

	return gnfmt.ToCSV(s, sep)
}

func (op *OutputPage) jsonOutput(pretty bool) string {
	enc := gnfmt.GNjson{Pretty: pretty}
	res, _ := enc.Encode(op)
	return string(res)
}
