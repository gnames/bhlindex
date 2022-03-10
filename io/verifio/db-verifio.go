package verifio

import (
	"context"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gnames/bhlindex/ent/name"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func (vrf verifio) loadNames(
	ctx context.Context,
	namesNum int,
	chNames chan<- []string,
	start time.Time) error {
	var count int
	batchSize := 5_000
	threshold := batchSize * 1

	q := `SELECT name
FROM name_statuses
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

		var name string
		names := make([]string, batchSize)
		var i int
		for rows.Next() {
			err := rows.Scan(&name)
			if err != nil {
				rows.Close()
				return err
			}
			names[i] = name
			i++
		}
		rows.Close()

		count = count + batchSize
		if count > namesNum {
			names = names[0:i]
			makeLog(start, namesNum, namesNum)
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
	chVer <-chan []vlib.Name,
) error {
	for namesVerif := range chVer {
		names := prepareNames(namesVerif)
		err := vrf.saveNamesToDB(names)
		if err != nil {
			return err
		}
	}
	return nil
}

func prepareNames(namesVerif []vlib.Name) []name.VerifiedName {
	res := make([]name.VerifiedName, len(namesVerif))
	for i, v := range namesVerif {
		n := name.VerifiedName{
			Name:              v.Name,
			MatchType:         v.MatchType.String(),
			Curation:          v.Curation.String(),
			DataSourcesNumber: v.DataSourcesNum,
			Error:             v.Error,
		}
		if br := v.BestResult; br != nil {
			n.RecordID = br.RecordID
			n.EditDistance = br.EditDistance
			n.StemEditDistance = br.StemEditDistance
			n.MatchedName = br.MatchedName
			n.MatchedCanonical = br.MatchedCanonicalFull
			n.CurrentName = br.CurrentName
			n.Classification = br.ClassificationPath
			n.DataSourceID = br.DataSourceID
			n.DataSourceTitle = br.DataSourceTitleShort
		}
		res[i] = n
	}
	return res
}

func (vrf verifio) saveNamesToDB(names []name.VerifiedName) error {
	now := time.Now()
	columns := []string{
		"name", "record_id", "match_type", "edit_distance", "stem_edit_distance",
		"matched_name", "matched_canonical", "current_name", "classification",
		"data_source_id", "data_source_title", "data_sources_number", "curation",
		"retries", "error", "updated_at",
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
			v.MatchedName, v.MatchedCanonical, v.CurrentName, v.Classification,
			v.DataSourceID, v.DataSourceTitle, v.DataSourcesNumber, v.Curation,
			v.Retries, v.Error, now,
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
