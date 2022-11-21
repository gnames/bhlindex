package item

import (
	"time"

	"github.com/gnames/bhlindex/internal/ent/page"
)

// Item in BHL is physically separate entity, such as a book, journal's volume etc.
type Item struct {
	// ID is an identifier from BHL database.
	ID int `json:"id"`
	// Path is a path to the Item's content.
	Path string `json:"path" gorm:"not null"`
	// UpdatedAt is the time when the item was modified last time in BHLindex.
	UpdatedAt time.Time `json:"updatedAt" sql:"type:timestamp without time zone"`
	// Pages is the slice of all pages located inside of the Item.
	Pages []*page.Page `json:"pages" gorm:"-"`
	// Text is the text of the item where texts of all pages are merged together.
	Text []byte `json:"-" gorm:"-"`
}
