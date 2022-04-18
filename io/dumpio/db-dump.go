package dumpio

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/gnames/bhlindex/ent/output"
	"github.com/rs/zerolog/log"
)

func (d *dumpio) checkForVerifiedNames() error {
	noNames, err := d.noVerifiedNames()
	if err != nil {
		return fmt.Errorf("checkForVerifiedNames: %w", err)
	}
	if noNames {
		err = errors.New("verified_names table is empty")
		log.Warn().Err(err).Msg("Run 'bhlindex verify' before 'bhlindex dump'")
		return err
	}
	return nil
}

func (d *dumpio) noVerifiedNames() (bool, error) {
	var name_id int
	q := "select name_id from verified_names limit 1"
	err := d.db.QueryRow(q).Scan(&name_id)
	return name_id == 0, err
}

func (d *dumpio) stats() (int, int, int, error) {
	var names, occurrs, items int
	err := d.db.QueryRow("SELECT max(name_id) from verified_names").Scan(&names)
	if err == nil {
		err = d.db.QueryRow("SELECT max(id) from detected_names").Scan(&occurrs)
	}
	if err == nil {
		err = d.db.QueryRow("SELECT max(id) from items").Scan(&items)
	}
	if err != nil {
		err = fmt.Errorf("stats: %w", err)
	}
	return names, occurrs, items, err
}

func (d *dumpio) outputs(id, limit int) ([]output.Output, error) {
	var rows *sql.Rows
	var err error
	q := `
SELECT
  dn.id, vn.name_id, dn.page_id, i.internet_archive_id, vn.name,
  vn.occurrences, vn.odds_log10, dn.offset_start, dn.offset_end,
  dn.ends_next_page, dn.cardinality, vn.match_type, vn.edit_distance,
  vn.matched_canonical, vn.matched_name, vn.matched_cardinality,
  vn.data_source_id, vn.data_source_title, vn.curation, vn.error
  FROM items i
    JOIN detected_names dn
      ON i.id = dn.item_id
    JOIN verified_names vn
      ON dn.name = vn.name
  WHERE i.id >= $1 and i.id < $2
ORDER by i.id
`
	rows, err = d.db.Query(q, id, id+limit)
	if err != nil {
		return nil, fmt.Errorf("outputs: %w", err)
	}
	defer rows.Close()

	var count int
	res := make([]output.Output, 0, limit)
	for rows.Next() {
		o := output.Output{}
		var pageBarcode string
		err := rows.Scan(
			&o.ID, &o.NameID, &pageBarcode, &o.ItemBarcode, &o.DetectedName,
			&o.Occurrences, &o.OddsLog10, &o.OffsetStart, &o.OffsetEnd,
			&o.EndsNextPage, &o.Cardinality, &o.MatchType, &o.EditDistance,
			&o.MatchedCanonical, &o.MatchedFullName, &o.MatchedCardinality,
			&o.DataSourceID, &o.DataSource, &o.Curation,
			&o.VerifError,
		)
		if err != nil {
			return nil, fmt.Errorf("outputs: %w", err)
		}

		o.PageBarcodeNum, err = pageNum(pageBarcode)
		if err != nil {
			return nil, fmt.Errorf("outputs: %w", err)
		}

		res = append(res, o)
		count++
	}

	return res, nil
}

func pageNum(barCode string) (int, error) {
	l := len(barCode)
	num := barCode[l-4 : l]
	return strconv.Atoi(num)
}
