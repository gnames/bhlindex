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
		return fmt.Errorf("-> noVerifiedNames %w", err)
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
	var verifNames, verifNamesDataSources, occurs int
	dataSources := getDataSources(ds)
	nameQ := fmt.Sprintf("SELECT count(*) as count FROM verified_names WHERE 1=1 %s",
		dataSources)
	err := d.db.QueryRow(nameQ).Scan(&verifNamesDataSources)
	if err == nil {
		err = d.db.QueryRow("SELECT count(*) FROM verified_names").Scan(&verifNames)
	}
	if err == nil {
		err = d.db.QueryRow("SELECT max(id) FROM detected_names").Scan(&occurs)
	}
	if err != nil {
		err = fmt.Errorf("-> stats %w", err)
	}
	return verifNames, verifNamesDataSources, occurs, err
}

func (d *dumpio) outputNames(id, limit int, ds []int) ([]output.Output, error) {
	var rows *sql.Rows
	var err error

	q := `
SELECT
  name, cardinality, occurrences, odds_log10, match_type, edit_distance,
  stem_edit_distance, matched_canonical, matched_name, matched_cardinality,
  current_canonical, current_name, current_cardinality, classification,
  classification_ranks, classification_ids, record_id, data_source_id,
  data_source_title, data_sources_number, curation, error, sort_order
  FROM verified_names
  WHERE id >= $1 and id < $2
  ORDER by id `
	rows, err = d.db.Query(q, id, id+limit)
	if err != nil {
		return nil, fmt.Errorf("-> Query %w", err)
	}
	defer rows.Close()

	res := make([]output.Output, 0, limit)
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
			return nil, fmt.Errorf("-> Scan %w", err)
		}

		o.NameID = gnuuid.New(o.DetectedName).String()
		if d.cfg.OutputShort {
			res = append(res, output.OutputNameShort{OutputName: o})
		} else {
			res = append(res, o)
		}
	}

	return res, nil
}

func (d *dumpio) outputOccurs(
	id, limit int,
	ds []int, normVerb bool,
) ([]output.Output, error) {
	var rows *sql.Rows
	var err error

	q := `
  SELECT
    page_id, item_id, name,
    name_verbatim, odds_log10, offset_start,
    offset_end, ends_next_page, annot_nomen_type
  FROM detected_names dn
  WHERE id >= $1 and id < $2
  ORDER BY id`
	rows, err = d.db.Query(q, id, id+limit)
	if err != nil {
		return nil, fmt.Errorf("-> Query %w", err)
	}
	defer rows.Close()

	var count int
	res := make([]output.Output, 0, limit)
	for rows.Next() {
		o := output.OutputOccurrence{}
		err := rows.Scan(
			&o.PageID, &o.ItemID, &o.DetectedName,
			&o.DetectedVerbatim, &o.OddsLog10, &o.OffsetStart,
			&o.OffsetEnd, &o.EndsNextPage, &o.Annotation,
		)
		if err != nil {
			return nil, fmt.Errorf("-> Scan %w", err)
		}
		o.NameID = gnuuid.New(o.DetectedName).String()

		if normVerb {
			o.NormalizeVerbatim()
		}

		if d.cfg.OutputShort {
			res = append(res, output.OutputOccurrenceShort{OutputOccurrence: o})
		} else {
			res = append(res, o)
		}
		count++
	}

	return res, nil
}

func (d dumpio) getOccurVerif() ([]output.Output, error) {
	var rows *sql.Rows
	var err error

	q := `
  SELECT
    odds_log10, names_num, verif_percent
  FROM odds_verifications
  ORDER BY odds_log10`
	rows, err = d.db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("-> Query %w", err)
	}
	defer rows.Close()

	var count int
	var res []output.Output
	for rows.Next() {
		o := output.OutputOddsVerification{}
		err := rows.Scan(
			&o.OddsLog10, &o.NamesNum, &o.VerifPercent,
		)
		if err != nil {
			return nil, fmt.Errorf("-> Scan %w", err)
		}
		res = append(res, o)
		count++
	}

	return res, nil
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
