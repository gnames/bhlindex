package verifio

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gnames/bhlindex/ent/name"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func (vrf verifio) ExtractUniqueNames() error {
	log.Info().Msg("Extracting unique name-strings. It will take a while.")
	q := `INSERT INTO unique_names (name, odds_log10, occurrences)
          SELECT name, AVG(odds_log10), count(*)
            FROM detected_names GROUP BY name
            ORDER BY name`

	stmt, err := vrf.db.Prepare(q)
	if err != nil {
		return fmt.Errorf("ExtractUniqueNames: %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("ExtractUniqueNames: %w", err)
	}

	err = stmt.Close()
	if err != nil {
		return fmt.Errorf("ExtractUniqueNames: %w", err)
	}
	return nil
}

func (vrf verifio) checkForDetectedNames() error {
	noNames, err := vrf.noDetectedNames()
	if err != nil {
		return fmt.Errorf("checkForDetectedNames: %w", err)
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
	q := "select page_id from detected_names limit 1"
	err := vrf.db.QueryRow(q).Scan(&page_id)
	return page_id == "", err
}

func (vrf verifio) numberOfNames() (int, error) {
	q := "select count(*) from unique_names"
	var namesNum int
	err := vrf.db.QueryRow(q).Scan(&namesNum)
	return namesNum, err
}

func (vrf verifio) truncateVerifTables() error {
	tables := []string{"unique_names", "verified_names"}
	for _, v := range tables {
		q := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY", v)
		_, err := vrf.db.Exec(q)
		if err != nil {
			return fmt.Errorf("truncateVerifTables: %w", err)
		}
	}
	q := "ALTER TABLE verified_names DROP CONSTRAINT IF EXISTS verified_names_pkey"
	_, err := vrf.db.Exec(q)
	if err != nil {
		return fmt.Errorf("truncateVerifTables: %w", err)
	}
	return nil
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
			return fmt.Errorf("loadNames: %w", err)
		}
		var n string
		var odds float64
		var id, occur int
		uns := make([]name.UniqueName, 0, limit)
		for rows.Next() {
			err := rows.Scan(&id, &n, &odds, &occur)
			if err != nil {
				rows.Close()
				return fmt.Errorf("loadNames: %w", err)
			}
			uns = append(uns, name.UniqueName{
				ID: id, Name: n, OddsLog10: odds, Occurrences: occur,
			})
		}
		rows.Close()

		offset += limit
		select {
		case <-ctx.Done():
			return fmt.Errorf("loadNames: %w", ctx.Err())
		case chIn <- uns:
		}
	}
	return nil
}

func makeLog(start time.Time, namesNum, count int) {
	names := humanize.Comma(int64(count))
	total := humanize.Comma(int64(namesNum))
	perHour := namesPerHour(start, count)
	percent := 100 * float64(count) / float64(namesNum)
	log.Info().
		Str("names", names+"/"+total).
		Str("names/hr", perHour).
		Msgf("Verified %0.1f%%", percent)
}

func namesPerHour(start time.Time, count int) string {
	dur := float64(time.Since(start)) / float64(time.Hour)
	rate := float64(count) / dur
	return humanize.Comma(int64(rate))
}

func (vrf verifio) saveVerif(
	ctx context.Context,
	chVer <-chan []name.VerifiedName,
	namesNum int,
	start time.Time,
) error {

	var count int
	for vns := range chVer {
		err := vrf.saveNamesToDB(vns)
		if err != nil {
			return fmt.Errorf("saveVerif: %w", err)
		}
		count = incrLog(start, namesNum, count, len(vns))
	}
	fmt.Fprint(os.Stderr, "\r")
	makeLog(start, namesNum, count)
	return nil
}

func incrLog(start time.Time, total, count, incr int) int {
	count += incr
	if count%100_000 == 0 {
		fmt.Fprint(os.Stderr, "\r")
		makeLog(start, total, count)
	} else if count%100 == 0 {
		fmt.Fprintf(os.Stderr, "\r%s", strings.Repeat(" ", 60))
		fmt.Fprintf(
			os.Stderr,
			"\rVerified %s names, %s items/hour",
			humanize.Comma(int64(count)),
			namesPerHour(start, count),
		)
	}
	return count
}

func (vrf verifio) saveNamesToDB(names []name.VerifiedName) error {
	now := time.Now()
	columns := []string{
		"name_id", "name", "record_id", "match_type", "edit_distance",
		"stem_edit_distance", "matched_name", "matched_canonical",
		"matched_cardinality", "current_name", "current_canonical",
		"current_cardinality", "classification", "classification_ranks",
		"classification_ids", "data_source_id", "data_source_title",
		"data_sources_number", "curation", "odds_log10", "occurrences", "retries",
		"error", "updated_at",
	}
	transaction, err := vrf.db.Begin()
	if err != nil {
		return fmt.Errorf("saveNamesToDB: %w", err)
	}
	defer transaction.Rollback()

	stmt, err := transaction.Prepare(pq.CopyIn("verified_names", columns...))
	if err != nil {
		return fmt.Errorf("saveNamesToDB: %w", err)
	}

	for _, v := range names {
		_, err = stmt.Exec(
			v.NameID, v.Name, v.RecordID, v.MatchType, v.EditDistance,
			v.StemEditDistance, v.MatchedName, v.MatchedCanonical,
			v.MatchedCardinality, v.CurrentName, v.CurrentCanonical,
			v.CurrentCardinality, v.Classification, v.ClassificationRanks,
			v.ClassificationIDs, v.DataSourceID, v.DataSourceTitle,
			v.DataSourcesNumber, v.Curation, v.OddsLog10, v.Occurrences, v.Retries,
			v.Error, now,
		)
		if err != nil {
			return fmt.Errorf("saveNamesToDB: %w", err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("saveNamesToDB: %w", err)
	}

	err = stmt.Close()
	if err != nil {
		return fmt.Errorf("saveNamesToDB: %w", err)
	}

	return transaction.Commit()
}

func (vrf verifio) setPrimaryKey() error {
	q := `ALTER TABLE verified_names ADD PRIMARY KEY (name_id)`
	_, err := vrf.db.Exec(q)
	return err
}
