package item

import (
	"time"

	"github.com/gnames/bhlindex/ent/page"
)

type Item struct {
	ID                int
	Path              string       `gorm:"not null"`
	InternetArchiveID string       `gorm:"unique_index;not null"`
	Status            int          `gorm:"index;not null;default:0"`
	UpdatedAt         time.Time    `sql:"type:timestamp without time zone"`
	Pages             []*page.Page `gorm:"-"`
	Text              []byte       `gorm:"-"`
}
