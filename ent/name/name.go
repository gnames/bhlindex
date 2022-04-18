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
	ID             uint      `json:"id"`
	PageID         string    `gorm:"type:varchar(255);index"`
	ItemID         int       `gorm:"not null;index"`
	Name           string    `sql:"type:CHARACTER VARYING(255) COLLATE \"C\" NOT NULL" gorm:"index"`
	NameVerbatim   string    `json:"nameVerbatim" sql:"type:CHARACTER VARYING(255) COLLATE \"C\" NOT NULL"`
	AnnotNomen     string    `gorm:"type:varchar(255)"`
	AnnotNomenType string    `gorm:"type:varchar(255)"`
	OffsetStart    int       `gorm:"not null;default:0"`
	OffsetEnd      int       `gorm:"not null;default:0"`
	EndsNextPage   bool      `gorm:"type:bool;default:false"`
	OddsLog10      float64   `gorm:"type:float;not null;default:0"`
	Cardinality    int       `gorm:"not null;default:0"`
	UpdatedAt      time.Time `sql:"type:timestamp without time zone"`
}

type VerifiedName struct {
	NameID              int       `json:"id"`
	Name                string    `json:"name" sql:"type:CHARACTER VARYING(255) COLLATE \"C\" NOT NULL" gorm:"index"`
	RecordID            string    `json:"recordID" sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`
	MatchType           string    `json:"matchType" gorm:"type:varchar(100)"`
	EditDistance        int       `json:"editDistance" gorm:"not null;default:0"`
	StemEditDistance    int       `json:"stemEditDistance" gorm:"not null;default:0"`
	MatchedName         string    `json:"matchedName" sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`
	MatchedCanonical    string    `json:"matchedCanonical" sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`
	MatchedCardinality  int       `json:"matchedCardinality"`
	CurrentName         string    `json:"currentName" sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`
	CurrentCanonical    string    `json:"currentCanonical" sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`
	CurrentCardinality  int       `json:"currentCardinality"`
	Classification      string    `json:"classification" sql:"type:CHARACTER VARYING COLLATE \"C\""`
	ClassificationRanks string    `json:"classificationRanks" sql:"type:CHARACTER VARYING COLLATE \"C\""`
	ClassificationIDs   string    `json:"classificationIds" sql:"type:CHARACTER VARYING COLLATE \"C\""`
	DataSourceID        int       `json:"dataSourceID" gorm:"index"`
	DataSourceTitle     string    `json:"dataSourceTitle" sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`
	DataSourcesNumber   int       `json:"dataSourcesNumber"`
	Curation            string    `json:"curation"`
	OddsLog10           float64   `json:"oddsLog10" gorm:"type:float;not null;default:0"`
	Occurrences         int       `json:"occurrences" gorm:"not null;default:0"`
	Retries             int       `json:"-" gorm:"not null;default:0"`
	Error               string    `json:"error"`
	UpdatedAt           time.Time `json:"updatedAt" sql:"type:timestamp without time zone"`
}

type UniqueName struct {
	ID          int     `json:"id" gorm:"primary_key"`
	Name        string  `sql:"type:CHARACTER VARYING(255) COLLATE \"C\""`
	OddsLog10   float64 `gorm:"type:float;not null;default:0"`
	Occurrences int     `gorm:"not null;default:0"`
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
		Name:           n.Name,
		NameVerbatim:   n.Verbatim,
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
