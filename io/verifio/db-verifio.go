package verifio

import (
	"context"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gnames/bhlindex/ent/name"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

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
	chNames chan<- []name.UniqueName,
	start time.Time) error {
	var count int
	batchSize := 5_000
	threshold := batchSize * 1

	q := `SELECT name, odds_log10, occurrences
FROM unique_names
OFFSET $1
LIMIT $2
`
	for count < namesNum {
		if count%threshold == 0 {
			makeLog(start, count, namesNum)
		}
		rows, err := vrf.db.Query(q, count, batchSize)
		if err != nil {
			return err
		}

		var n string
		var odds float64
		var occur int
		names := make([]name.UniqueName, batchSize)
		var i int
		for rows.Next() {
			err := rows.Scan(&n, &odds, &occur)
			if err != nil {
				rows.Close()
				return err
			}
			names[i] = name.UniqueName{Name: n, OddsLog10: odds, Occurrences: occur}
			i++
		}
		rows.Close()

		count = count + batchSize
		if count > namesNum {
			names = names[0:i]
			makeLog(start, namesNum, namesNum)
			log.Info().Msg("Finishing verification jobs")
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case chNames <- names:
		}
	}
	return nil
}

func makeLog(start time.Time, count, namesNum int) {
	if count == 0 {
		return
	}

	names := humanize.Comma(int64(count))
	perHour := namesPerHour(start, count)
	percent := 100 * float64(count) / float64(namesNum)
	log.Info().
		Str("names", names).
		Str("names/hr", perHour).
		Msgf("Verified %0.1f%% of names", percent)
}

func namesPerHour(start time.Time, count int) string {
	dur := float64(time.Since(start)) / float64(time.Hour)
	rate := float64(count) / dur
	return humanize.Comma(int64(rate))
}

func (vrf verifio) saveVerif(
	ctx context.Context,
	chVer <-chan []name.VerifiedName,
) error {
	for vns := range chVer {
		err := vrf.saveNamesToDB(vns)
		if err != nil {
			return err
		}
	}
	return nil
}

func (vrf verifio) saveNamesToDB(names []name.VerifiedName) error {
	now := time.Now()
	columns := []string{
		"name", "record_id", "match_type", "edit_distance", "stem_edit_distance",
		"matched_name", "matched_canonical", "current_name", "current_canonical",
		"classification", "data_source_id", "data_source_title",
		"data_sources_number", "curation", "odds_log10", "occurrences", "retries",
		"error", "updated_at",
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
			v.Name, v.RecordID, v.MatchType, v.EditDistance, v.StemEditDistance,
			v.MatchedName, v.MatchedCanonical, v.CurrentName, v.CurrentCanonical,
			v.Classification, v.DataSourceID, v.DataSourceTitle, v.DataSourcesNumber,
			v.Curation, v.OddsLog10, v.Occurrences, v.Retries, v.Error, now,
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
