package finderio

import (
	"fmt"
	"math"
	"time"

	"github.com/gnames/bhlindex/ent/name"
	"github.com/lib/pq"
)

func (fdr finderio) savePageNameStrings(names []name.DetectedName) error {
	now := time.Now()
	columns := []string{"page_id", "item_id", "name", "annot_nomen",
		"annot_nomen_type", "offset_start", "offset_end",
		"ends_next_page", "odds_log10", "cardinality", "updated_at"}
	transaction, err := fdr.db.Begin()
	if err != nil {
		return fmt.Errorf("savePageNameStrings: %w", err)
	}
	defer transaction.Rollback()

	stmt, err := transaction.Prepare(pq.CopyIn("detected_names", columns...))
	if err != nil {
		return fmt.Errorf("savePageNameStrings: %w", err)
	}

	for _, v := range names {
		if math.IsInf(v.OddsLog10, -1) {
			v.OddsLog10 = 0
		}
		_, err = stmt.Exec(v.PageID, v.ItemID, v.Name,
			v.AnnotNomen, v.AnnotNomenType, v.OffsetStart, v.OffsetEnd,
			v.EndsNextPage, v.OddsLog10, v.Cardinality, now)
		if err != nil {
			return fmt.Errorf("savePageNameStrings: %w", err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("savePageNameStrings: %w", err)
	}

	err = stmt.Close()
	if err != nil {
		return fmt.Errorf("savePageNameStrings: %w", err)
	}

	return transaction.Commit()
}
