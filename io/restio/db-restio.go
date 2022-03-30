package restio

import (
	"context"
	"fmt"

	"github.com/gnames/bhlindex/ent/item"
	"github.com/gnames/bhlindex/ent/name"
	"github.com/gnames/bhlindex/ent/page"
	"github.com/gnames/bhlindex/ent/rest"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func (r restio) items(
	ctx context.Context,
	inp rest.Input,
) ([]item.Item, error) {
	args := []any{inp.OffsetID, inp.OffsetID + inp.Limit}
	q := `SELECT
  id, path, internet_archive_id, updated_at
  FROM items
  WHERE id >= $1
    AND id < $2`

	return r.itemsQuery(ctx, q, args, inp.Limit)
}

func (r restio) itemsQuery(
	ctx context.Context,
	q string,
	args []any,
	limit int,
) ([]item.Item, error) {
	res := make([]item.Item, limit)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("itemsQuery: %w", err)
	}
	defer rows.Close()

	var i int
	for rows.Next() {
		var itm item.Item
		if err = rows.Scan(
			&itm.ID, &itm.Path, &itm.InternetArchiveID, &itm.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("itemsQuery: %w", err)
		}
		res[i] = itm

		i++
	}
	if i < limit-1 {
		res = res[0:i]
	}
	return res, nil
}

func (r restio) pages(
	ctx context.Context,
	inp rest.Input,
) ([]page.Page, error) {
	args := []any{inp.OffsetID, inp.OffsetID + inp.Limit}
	q := `SELECT
  id, item_id
  FROM pages
  WHERE item_id >= $1
    AND item_id < $2`

	return r.pagesQuery(ctx, q, args, inp.Limit)
}

func (r restio) pagesQuery(
	ctx context.Context,
	q string,
	args []any,
	limit int,
) ([]page.Page, error) {
	var res []page.Page

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("pagesQuery: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pg page.Page
		if err = rows.Scan(
			&pg.ID, &pg.ItemID,
		); err != nil {
			return nil, fmt.Errorf("pagesQuery: %w", err)
		}
		res = append(res, pg)
	}
	return res, nil
}

func (r restio) occurrences(
	ctx context.Context,
	inp rest.Input,
) ([]name.DetectedName, error) {
	args := []any{inp.OffsetID, inp.OffsetID + inp.Limit}
	q := `SELECT
  id, page_id, item_id, name, annot_nomen,
  annot_nomen_type, offset_start, offset_end,
  ends_next_page, odds_log10, cardinality,
  updated_at
  FROM detected_names
  where id >= $1
    AND id < $2`

	return r.occurrencesQuery(ctx, q, args, inp.Limit)
}

func (r restio) occurrencesQuery(
	ctx context.Context,
	q string,
	args []any,
	limit int,
) ([]name.DetectedName, error) {
	res := make([]name.DetectedName, limit)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("occurrencesQuery: %w", err)
	}
	defer rows.Close()

	var i int
	for rows.Next() {
		var dn name.DetectedName
		if err = rows.Scan(
			&dn.ID, &dn.PageID, &dn.ItemID, &dn.Name, &dn.AnnotNomen,
			&dn.AnnotNomenType, &dn.OffsetStart, &dn.OffsetEnd,
			&dn.EndsNextPage, &dn.OddsLog10, &dn.Cardinality,
			&dn.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("occurrencesQuery: %w", err)
		}
		res[i] = dn
		i++
	}
	if i < limit-1 {
		res = res[0:i]
	}
	return res, nil
}

func (r restio) names(
	ctx context.Context,
	inp rest.Input,
) ([]name.VerifiedName, error) {
	args := []any{inp.OffsetID, inp.OffsetID + inp.Limit}
	q := `SELECT
  name_id, name, record_id, match_type, edit_distance,
  stem_edit_distance, matched_name, matched_canonical,
  current_name, current_canonical, classification,
  data_source_id, data_source_title, data_sources_number,
  curation, odds_log10, occurrences, error, updated_at
  FROM verified_names
  WHERE name_id >= $1
    AND name_id < $2`

	if len(inp.DataSources) > 0 {
		args = append(args, pq.Array(inp.DataSources))
		q += "\n  AND data_source_id = any($3::int[])"
	}

	select {
	case <-ctx.Done():
		log.Warn().Err(ctx.Err()).Msg("Forced cancellation")
		return nil, ctx.Err()
	default:
		return r.namesQuery(ctx, q, args, inp.Limit)
	}
}

func (r restio) namesQuery(
	ctx context.Context,
	q string,
	args []any,
	limit int,
) ([]name.VerifiedName, error) {
	res := make([]name.VerifiedName, limit)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("namesQuery: %w", err)
	}
	defer rows.Close()

	var i int
	for rows.Next() {
		var vn name.VerifiedName
		if err = rows.Scan(
			&vn.NameID, &vn.Name, &vn.RecordID, &vn.MatchType, &vn.EditDistance,
			&vn.StemEditDistance, &vn.MatchedName, &vn.MatchedCanonical,
			&vn.CurrentName, &vn.CurrentCanonical, &vn.Classification,
			&vn.DataSourceID, &vn.DataSourceTitle, &vn.DataSourcesNumber,
			&vn.Curation, &vn.OddsLog10, &vn.Occurrences, &vn.Error, &vn.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("namesQuery: %w", err)
		}
		res[i] = vn
		i++
	}
	if i < limit-1 {
		res = res[0:i]
	}
	return res, nil
}

func (r restio) namesLastID() (int, error) {
	var lastID int
	err := r.db.QueryRow("SELECT max(name_id) FROM verified_names").
		Scan(&lastID)
	return lastID, err
}
