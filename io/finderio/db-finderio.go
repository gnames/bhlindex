package finderio

import (
	"time"

	"github.com/gnames/bhlindex/ent/name"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func (fdr finderio) ExtractUniqueNames() error {
	log.Info().Msg("Extracting unique name-strings. It will take a while.")
	q := `INSERT INTO unique_names
          SELECT name_string, AVG(odds_log10), count(*)
            FROM detected_names GROUP BY name_string
            ORDER BY name_string`

	stmt, err := fdr.db.Prepare(q)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

func (fdr finderio) savePageNameStrings(names []name.DetectedName) error {
	now := time.Now()
	columns := []string{"page_id", "item_id", "name_string", "annot_nomen",
		"annot_nomen_type", "offset_start", "offset_end",
		"ends_next_page", "odds_log10", "cardinality", "updated_at"}
	transaction, err := fdr.db.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	stmt, err := transaction.Prepare(pq.CopyIn("detected_names", columns...))
	if err != nil {
		return err
	}

	for _, v := range names {
		_, err = stmt.Exec(v.PageID, v.ItemID, v.NameString,
			v.AnnotNomen, v.AnnotNomenType, v.OffsetStart, v.OffsetEnd,
			v.EndsNextPage, v.OddsLog10, v.Cardinality, now)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return transaction.Commit()
}
