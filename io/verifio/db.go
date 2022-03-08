package verifio

import (
	"context"
	"time"

	"github.com/dustin/go-humanize"
	vlib "github.com/gnames/gnlib/ent/verifier"
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
	chVer <-chan vlib.Output,
) error {
	for range chVer {
	}
	return nil
}
