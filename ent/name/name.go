package name

import (
	"math"
	"time"

	"github.com/gnames/bhlindex/ent/page"
	gnfout "github.com/gnames/gnfinder/ent/output"
)

// DetectedName holds information about a name-string returned by a
// name-finder.
type DetectedName struct {
	PageID         string    `gorm:"type:varchar(255);index"`
	ItemID         int       `gorm:"not null;index"`
	NameString     string    `sql:"type:CHARACTER VARYING(255) COLLATE \"C\" NOT NULL" gorm:"index"`
	AnnotNomen     string    `gorm:"type:varchar(255)"`
	AnnotNomenType string    `gorm:"type:varchar(255)"`
	OffsetStart    int       `gorm:"not null;default:0"`
	OffsetEnd      int       `gorm:"not null;default:0"`
	EndsNextPage   bool      `gorm:"type:bool;default:false"`
	OddsLog10      float64   `gorm:"type:float;not null;default:0"`
	Cardinality    int       `gorm:"not null;default:0"`
	UpdatedAt      time.Time `sql:"type:timestamp without time zone"`
}

type NameString struct {
	ID                int    `gorm:"primary_key"`
	Name              string `sql:"type:CHARACTER VARYING(255) COLLATE \"C\" NOT NULL"`
	TaxonID           string `sql:"type:CHARACTER VARYING(255) COLLATE \"C\"" gorm:"unique_index"`
	MatchType         string `gorm:"type:varchar(100)"`
	EditDistance      int    `gorm:"not null;default:0"`
	StemEditDistance  int    `gorm:"not null;default:0"`
	MatchedName       string `sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`
	MatchedCanonical  string `sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`
	CurrentName       string `sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`
	Classification    string `sql:"type:CHARACTER VARYING COLLATE \"C\""`
	DataSourceID      int
	DatasourceTitle   string `sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`
	DataSourcesNumber int
	Curation          string
	Retries           int `gorm:"not null;default:0"`
	Error             string
	UpdatedAt         time.Time `sql:"type:timestamp without time zone"`
}

type NameStatus struct {
	Name        string  `sql:"type:CHARACTER VARYING(255) COLLATE \"C\"" gorm:"primary_key;auto_increment:false"`
	OddsLog10   float64 `gorm:"type:float;not null;default:0"`
	Occurrences int     `gorm:"not null;default:0"`
	Processed   bool    `gorm:"not null;default:false;index"`
}

func New(itemID int, p *page.Page, n gnfout.Name) DetectedName {
	var endsNextPage bool
	var end int
	start := n.OffsetStart - p.Offset
	if n.OffsetEnd < p.OffsetNext {
		end = n.OffsetEnd - p.Offset
	} else {
		end = n.OffsetEnd - p.OffsetNext
		endsNextPage = true
	}
	dn := DetectedName{
		PageID:         p.ID,
		ItemID:         itemID,
		NameString:     n.Name,
		OffsetStart:    start,
		OffsetEnd:      end,
		AnnotNomen:     n.AnnotNomen,
		AnnotNomenType: n.AnnotNomenType,
		EndsNextPage:   endsNextPage,
		OddsLog10:      math.Log10(n.Odds),
		Cardinality:    n.Cardinality,
		UpdatedAt:      time.Now(),
	}
	return dn
}
