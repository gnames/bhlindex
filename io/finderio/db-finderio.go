package finderio

import (
	"context"
	"time"

	"github.com/gnames/bhlindex/ent/name"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func (fdr finderio) SaveNames(
	ctx context.Context,
	namesCh <-chan []name.DetectedName,
) error {
	for v := range namesCh {
		_ = v
		err := fdr.savePageNameStrings(v)
		if err != nil {
			log.Warn().Err(err).Msg("Cannot save detected names")
			return err
		}
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
