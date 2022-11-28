package item

import (
	"github.com/gnames/bhlindex/internal/ent/page"
)

// Item in BHL is physically separate entity, such as a book, journal's volume etc.
type Item struct {
	// ID is an identifier from BHL database.
	ID int
	// Path is a path to the Item's content.
	Path string
	// Pages is the slice of all pages located inside of the Item.
	Pages []*page.Page
	// Text is the text of the item where texts of all pages are merged together.
	Text []byte
}
