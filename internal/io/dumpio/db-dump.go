package dumpio

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gnames/bhlindex/internal/ent/output"
	"github.com/gnames/gnuuid"
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
	var nameID int
	q := "select id from verified_names limit 1"
	err := d.db.QueryRow(q).Scan(&nameID)
	return nameID == 0, err
}

func (d *dumpio) stats(ds []int) (int, int, int, error) {
	var allNames, names, items int
	dataSources := getDataSources(ds)
	nameQ := fmt.Sprintf("SELECT count(*) as count FROM verified_names WHERE 1=1 %s",
		dataSources)
	err := d.db.QueryRow(nameQ).Scan(&names)
	if err == nil {
		err = d.db.QueryRow("SELECT max(id) FROM verified_names").Scan(&allNames)
	}
	if err == nil {
		err = d.db.QueryRow("SELECT max(id) FROM items").Scan(&items)
	}
	if err != nil {
		err = fmt.Errorf("stats: %w", err)
	}
	return allNames, names, items, err
}

func (d *dumpio) outputPages(id, limit int) ([]output.OutputPage, error) {
	var rows *sql.Rows
	var err error

	q := `
SELECT
  i.internet_archive_id, p.id
  FROM items i
    JOIN pages p
      ON i.id = p.item_id
  WHERE i.id >= $1 and i.id < $2
ORDER by i.id, p.item_id
`
	rows, err = d.db.Query(q, id, id+limit)
	if err != nil {
		return nil, fmt.Errorf("outputPages: %w", err)
	}
	defer rows.Close()

	var count int
	res := make([]output.OutputPage, 0, limit)
	for rows.Next() {
		o := output.OutputPage{}
		var pageBarcode string
		err := rows.Scan(
			&o.ItemBarcode, &pageBarcode,
		)
		if err != nil {
			return nil, fmt.Errorf("outputPages: %w", err)
		}

		o.PageBarcodeNum, err = pageNum(pageBarcode)
		if err != nil {
			return nil, fmt.Errorf("outputPages: %w", err)
		}

		res = append(res, o)
		count++
	}

	return res, nil
}

func (d *dumpio) outputNames(id, limit int, ds []int) ([]output.OutputName, error) {
	var rows *sql.Rows
	var err error

	dataSources := getDataSources(ds)

	q := fmt.Sprintf(`
SELECT
  name, cardinality, occurrences, odds_log10, match_type, edit_distance,
  stem_edit_distance, matched_canonical, matched_name, matched_cardinality,
  current_canonical, current_name, current_cardinality, classification,
  classification_ranks, classification_ids, record_id, data_source_id,
  data_source_title, data_sources_number, curation, error, sort_order
  FROM verified_names
  WHERE id >= $1 and id < $2
  %s
  ORDER by id
`, dataSources)
	rows, err = d.db.Query(q, id, id+limit)
	if err != nil {
		return nil, fmt.Errorf("outputNames: %w", err)
	}
	defer rows.Close()

	res := make([]output.OutputName, 0, limit)
	for rows.Next() {
		o := output.OutputName{}
		err := rows.Scan(
			&o.DetectedName, &o.Cardinality, &o.OccurrencesNumber, &o.OddsLog10,
			&o.MatchType, &o.EditDistance, &o.StemEditDistance, &o.MatchedCanonical,
			&o.MatchedFullName, &o.MatchedCardinality, &o.CurrentCanonical,
			&o.CurrentFullName, &o.CurrentCardinality, &o.Classification,
			&o.ClassificationRanks, &o.ClassificationIDs, &o.RecordID,
			&o.DataSourceID, &o.DataSource, &o.DataSourcesNumber, &o.Curation,
			&o.VerifError, &o.SortOrder,
		)
		if err != nil {
			return nil, fmt.Errorf("outputNames: %w", err)
		}

		o.NameID = gnuuid.New(o.DetectedName).String()

		res = append(res, o)
	}

	return res, nil
}

func (d *dumpio) outputOccurs(id, limit int, ds []int) ([]output.OutputOccurrence, error) {
	var rows *sql.Rows
	var err error

	dataSources := getDataSources(ds)

	q := fmt.Sprintf(`
SELECT
  dn.page_id, i.internet_archive_id, vn.name,
  dn.name_verbatim, vn.odds_log10, dn.offset_start,
  dn.offset_end, dn.ends_next_page, dn.annot_nomen_type

  FROM items i
    JOIN detected_names dn
      ON i.id = dn.item_id
    JOIN verified_names vn
      ON dn.name = vn.name
  WHERE i.id >= $1 and i.id < $2
  %s
ORDER by i.id
`, dataSources)
	rows, err = d.db.Query(q, id, id+limit)
	if err != nil {
		return nil, fmt.Errorf("outputOccurs: %w", err)
	}
	defer rows.Close()

	var count int
	res := make([]output.OutputOccurrence, 0, limit)
	for rows.Next() {
		o := output.OutputOccurrence{}
		var pageBarcode string
		err := rows.Scan(
			&pageBarcode, &o.ItemBarcode, &o.DetectedName,
			&o.DetectedVerbatim, &o.OddsLog10, &o.OffsetStart,
			&o.OffsetEnd, &o.EndsNextPage, &o.Annotation,
		)
		if err != nil {
			return nil, fmt.Errorf("outputOccurs: %w", err)
		}

		o.PageBarcodeNum, err = pageNum(pageBarcode)
		if err != nil {
			return nil, fmt.Errorf("outputOccurs: %w", err)
		}
		o.NameID = gnuuid.New(o.DetectedName).String()

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

func getDataSources(ds []int) string {
	var dataSources string
	if len(ds) > 0 {
		dsStr := make([]string, len(ds))
		for i := range ds {
			dsStr[i] = strconv.Itoa(ds[i])
		}
		dataSources = fmt.Sprintf("AND data_source_id  IN (%s)",
			strings.Join(dsStr, ", "))
	}
	return dataSources
}
