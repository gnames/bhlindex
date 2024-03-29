package verifio

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gnames/bhlindex/internal/ent/name"
	"github.com/gnames/gnfmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func (vrf verifio) maxOdds() (float64, error) {
	var res float64
	q := "SELECT MAX(odds_log10) FROM verified_names LIMIT 1"
	err := vrf.db.QueryRow(q).Scan(&res)
	return res, err
}

func (vrf verifio) oddsVerif(minOdds, maxOdds int) error {
	q := `SELECT
  match_type, count(*)
FROM verified_names
  WHERE odds_log10 > $1 AND odds_log10 >= $2
    AND name !~ '\A[A-Za-z]{0,4}\.'
GROUP BY match_type
`

	rows, err := vrf.db.Query(q, minOdds, maxOdds)
	if err != nil {
		return fmt.Errorf("-> Query %w", err)
	}
	defer rows.Close()

	var verif, notVerif int
	for rows.Next() {
		var match string
		var counts int
		err = rows.Scan(&match, &counts)
		if err != nil {
			return err
		}
		if match == "NoMatch" {
			notVerif += counts
		} else {
			verif += counts
		}
	}
	pcent := float64(verif) / float64(verif+notVerif)
	q = `INSERT INTO odds_verifications
  (odds_log10, names_num, verif_percent)
  VALUES ($1, $2, $3)
  `
	_, err = vrf.db.Query(q, maxOdds, verif+notVerif, pcent)
	if err != nil {
		return fmt.Errorf("-> Query %w", err)
	}
	return nil
}

func (vrf verifio) showOddsVerif() error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	header := table.Row{"Odds Log10", "Names Number", "% Verified"}
	t.AppendHeader(header)

	q := `SELECT odds_log10, names_num, verif_percent
FROM odds_verifications order by odds_log10`

	rows, err := vrf.db.Query(q)
	if err != nil {
		return fmt.Errorf("-> Query %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ov name.OddsVerification
		err = rows.Scan(&ov.OddsLog10, &ov.NamesNum, &ov.VerifPercent)
		if err != nil {
			return err
		}

		row := table.Row{
			ov.OddsLog10,
			ov.NamesNum,
			ov.VerifPercent,
		}
		t.AppendRow(row)
	}
	t.Render()
	return nil
}

func (vrf verifio) ExtractUniqueNames() error {
	log.Info().Msg("Extracting unique name-strings. It will take a while.")
	q := `INSERT INTO unique_names (name, odds_log10, occurrences)
          SELECT name, AVG(odds_log10), count(*)
            FROM detected_names GROUP BY name
            ORDER BY name`

	stmt, err := vrf.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("-> db.Prepare %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("-> Exec %w", err)
	}

	err = stmt.Close()
	if err != nil {
		return fmt.Errorf("-> Close %w", err)
	}
	return nil
}

func (vrf verifio) checkForDetectedNames() error {
	noNames, err := vrf.noDetectedNames()
	if err != nil {
		return fmt.Errorf("-> noDetectedNames %w", err)
	}
	if noNames {
		err = errors.New("detected_names table is empty")
		log.Warn().Err(err).Msg("Run 'bhlindex find' before 'bhlindex verify'")
		return err
	}
	return nil
}

func (vrf verifio) noDetectedNames() (bool, error) {
	var page_id string
	q := "SELECT page_id FROM detected_names LIMIT 1"
	err := vrf.db.QueryRow(q).Scan(&page_id)
	return page_id == "", err
}

func (vrf verifio) numberOfNames() (int, error) {
	q := "SELECT count(*) FROM unique_names"
	var namesNum int
	err := vrf.db.QueryRow(q).Scan(&namesNum)
	return namesNum, err
}

func (vrf verifio) truncateVerifTables() error {
	tables := []string{"unique_names", "verified_names", "odds_verifications"}
	for _, v := range tables {
		q := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY", v)
		_, err := vrf.db.Exec(q)
		if err != nil {
			return err
		}
	}
	q := "ALTER TABLE verified_names DROP CONSTRAINT IF EXISTS verified_names_pkey"
	_, err := vrf.db.Exec(q)
	return err
}

func (vrf verifio) loadNames(
	ctx context.Context,
	namesNum int,
	chIn chan<- []name.UniqueName,
) error {
	offset := 1
	limit := 5_000

	q := `
SELECT id, name, odds_log10, occurrences
  FROM unique_names
  WHERE id >= $1 and id < $2
`
	for offset < namesNum {
		rows, err := vrf.db.Query(q, offset, offset+limit)
		if err != nil {
			return err
		}
		var n string
		var odds float64
		var id, occur int
		uns := make([]name.UniqueName, 0, limit)
		for rows.Next() {
			err := rows.Scan(&id, &n, &odds, &occur)
			if err != nil {
				rows.Close()
				return fmt.Errorf("-> Scan %w", err)
			}
			uns = append(uns, name.UniqueName{
				ID: id, Name: n, OddsLog10: odds, Occurrences: occur,
			})
		}
		rows.Close()

		offset += limit
		select {
		case <-ctx.Done():
			return ctx.Err()
		case chIn <- uns:
		}
	}
	return nil
}

func logStr(start time.Time, namesNum, count int) string {
	rate := namesPerHour(start, count)
	namesStr := humanize.Comma(int64(count))
	perHourStr := humanize.Comma(int64(rate))
	percent := 100 * float64(count) / float64(namesNum)
	eta := 3600 * float64(namesNum-count) / rate
	etaStr := gnfmt.TimeString(eta)
	return fmt.Sprintf(
		"%s verified names (%0.1f%%), %s names/hr: ETA %s",
		namesStr, percent, perHourStr, etaStr,
	)
}

func namesPerHour(start time.Time, count int) float64 {
	dur := float64(time.Since(start)) / float64(time.Hour)
	return float64(count) / dur
}

func (vrf verifio) saveVerif(
	ctx context.Context,
	chVer <-chan verifiedBatch,
	namesNum int,
	start time.Time,
) error {
	var count int
	for vns := range chVer {
		err := vrf.saveNamesToDB(vns.names)
		if err != nil {
			return fmt.Errorf("-> saveNamesToDB %w", err)
		}
		count = incrLog(start, namesNum, count, vns.namesNum)
	}
	fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", 80))
	log.Info().Msgf("Verified %s names", humanize.Comma(int64(namesNum)))

	return nil
}

func incrLog(start time.Time, total, count, incr int) int {
	count += incr
	if count%1_000_000 == 0 {
		fmt.Fprint(os.Stderr, "\r")
		log.Info().Msg(logStr(start, total, count))
	} else if count%100 == 0 {
		fmt.Fprintf(os.Stderr, "\r%s", strings.Repeat(" ", 80))
		fmt.Fprint(os.Stderr, "\r"+logStr(start, total, count))
	}
	return count
}

func (vrf verifio) saveNamesToDB(names []name.VerifiedName) error {
	now := time.Now()
	columns := []string{
		"id", "name", "cardinality", "record_id", "match_type",
		"edit_distance", "stem_edit_distance", "matched_name", "matched_canonical",
		"matched_cardinality", "current_name", "current_canonical",
		"current_cardinality", "classification", "classification_ranks",
		"classification_ids", "data_source_id", "data_source_title",
		"data_sources_number", "curation", "odds_log10", "sort_order",
		"occurrences", "retries", "error", "updated_at",
	}
	transaction, err := vrf.db.Begin()
	if err != nil {
		return fmt.Errorf("-> Begin %w", err)
	}
	defer transaction.Rollback()

	stmt, err := transaction.Prepare(pq.CopyIn("verified_names", columns...))
	if err != nil {
		return fmt.Errorf("-> Prepare %w", err)
	}

	for _, v := range names {
		_, err = stmt.Exec(
			v.ID, v.Name, v.Cardinality, v.RecordID, v.MatchType, v.EditDistance,
			v.StemEditDistance, v.MatchedName, v.MatchedCanonical,
			v.MatchedCardinality, v.CurrentName, v.CurrentCanonical,
			v.CurrentCardinality, v.Classification, v.ClassificationRanks,
			v.ClassificationIDs, v.DataSourceID, v.DataSourceTitle,
			v.DataSourcesNumber, v.Curation, v.OddsLog10, v.SortOrder,
			v.Occurrences, v.Retries, v.Error, now,
		)
		if err != nil {
			return fmt.Errorf("-> Exec %w", err)
		}
	}

	// Flush COPY FROM to db.
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("-> Exec flush %w", err)
	}

	err = stmt.Close()
	if err != nil {
		return fmt.Errorf("-> Close %w", err)
	}

	return transaction.Commit()
}
