package item

import (
	"time"

	"github.com/gnames/bhlindex/internal/ent/page"
)

type Item struct {
	ID                int          `json:"id"`
	Path              string       `json:"path" gorm:"not null"`
	InternetArchiveID string       `json:"internetArchiveId" gorm:"unique_index;not null"`
	UpdatedAt         time.Time    `json:"updatedAt" sql:"type:timestamp without time zone"`
	Pages             []*page.Page `json:"pages" gorm:"-"`
	Text              []byte       `json:"-" gorm:"-"`
}
