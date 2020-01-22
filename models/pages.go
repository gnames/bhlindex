package models

import "path/filepath"

// Page is a representation of a page file. ID is the filename, ItemID
// is id of the parent item, Offset, number of runes the page is away from
// the start of the text.
type Page struct {
	ID         string
	ItemID     int
	Offset     int
	OffsetNext int
}

func IsPageFile(f string) bool {
	res, _ := filepath.Match("*_[0-9][0-9][0-9][0-9].txt", f)
	return res
}

func PageID(f string) string {
	extLen := len(filepath.Ext(f))
	idLen := len(f) - extLen
	return f[0:idLen]
}
