package verifio

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gnames/bhlindex/ent/name"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func (vrf verifio) ExtractUniqueNames() error {
	log.Info().Msg("Extracting unique name-strings. It will take a while.")
	q := `INSERT INTO unique_names (name, odds_log10, occurrences, processed)
          SELECT name, AVG(odds_log10), count(*), false
            FROM detected_names GROUP BY name
            ORDER BY name`

	stmt, err := vrf.db.Prepare(q)
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

func (vrf verifio) checkForDetectedNames() error {
	noNames, err := vrf.noDetectedNames()
	if err != nil {
		return err
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
			return err
		}
	}
	return nil
}

func (vrf verifio) loadNames(
	ctx context.Context,
	namesNum int,
	chIn chan<- []name.UniqueName,
) error {
	var count int
	batchSize := 5_000

	q := `
    WITH temp AS (
      SELECT id, name, odds_log10, occurrences FROM unique_names
        WHERE processed=false LIMIT $1)
    UPDATE unique_names ns SET processed=true
    FROM temp
      WHERE ns.name = temp.name
    RETURNING temp.id, temp.name, temp.odds_log10, temp.occurrences`

	for count < namesNum {
		rows, err := vrf.db.Query(q, batchSize)
		if err != nil {
			return err
		}
		var n string
		var odds float64
		var id, occur int
		uns := make([]name.UniqueName, batchSize)
		var i int
		for rows.Next() {
			err := rows.Scan(&id, &n, &odds, &occur)
			if err != nil {
				rows.Close()
				return err
			}
			uns[i] = name.UniqueName{
				ID: id, Name: n, OddsLog10: odds, Occurrences: occur,
			}
			i++
		}
		rows.Close()

		count = count + batchSize
		if count > namesNum {
			uns = uns[0:i]
		}
		select {
		case <-ctx.Done():
			log.Warn().Err(ctx.Err()).Msg("Exiting loading names")
			return ctx.Err()
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
			return err
		}
		count = incrLog(start, namesNum, count, len(vns))
	}
	fmt.Print("\r")
	makeLog(start, namesNum, count)
	return nil
}

func incrLog(start time.Time, total, count, incr int) int {
	count += incr
	if count%100_000 == 0 {
		fmt.Print("\r")
		makeLog(start, total, count)
	} else if count%100 == 0 {
		fmt.Printf("\r%s", strings.Repeat(" ", 60))
		fmt.Printf(
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
		"stem_edit_distance", "matched_name", "matched_canonical", "current_name",
		"current_canonical", "classification", "data_source_id",
		"data_source_title", "data_sources_number", "curation", "odds_log10",
		"occurrences", "retries", "error", "updated_at",
	}
	transaction, err := vrf.db.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	stmt, err := transaction.Prepare(pq.CopyIn("verified_names", columns...))
	if err != nil {
		return err
	}

	for _, v := range names {
		_, err = stmt.Exec(
			v.NameID, v.Name, v.RecordID, v.MatchType, v.EditDistance,
			v.StemEditDistance, v.MatchedName, v.MatchedCanonical, v.CurrentName,
			v.CurrentCanonical, v.Classification, v.DataSourceID, v.DataSourceTitle,
			v.DataSourcesNumber, v.Curation, v.OddsLog10, v.Occurrences, v.Retries,
			v.Error, now,
		)
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

func (vrf verifio) setPrimaryKey() error {
	q := `ALTER TABLE verified_names ADD PRIMARY KEY (name_id)`
	_, err := vrf.db.Exec(q)
	return err
}
