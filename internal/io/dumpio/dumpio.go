package dumpio

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gnames/bhlindex/internal/config"
	"github.com/gnames/bhlindex/internal/ent/output"
	"github.com/rs/zerolog/log"
)

type dumpio struct {
	cfg config.Config
	db  *sql.DB
}

// New creates a new instance of Dumper.
func New(cfg config.Config, db *sql.DB) output.Dumper {
	return &dumpio{cfg: cfg, db: db}
}

// DumpNames outputs data about verified names.
func (d *dumpio) DumpNames(ctx context.Context, ch chan<- []output.Output, ds []int) error {
	err := d.checkForVerifiedNames()
	if err != nil {
		err = fmt.Errorf("-> checkForVerifiedNames %w", err)
		log.Warn().Msg("Run `bhlindex verify` before `bhlindex dump`.")
		return err
	}

	namesTotal, namesNum, _, err := d.stats(ds)
	if err != nil {
		return fmt.Errorf("-> stats %w", err)
	}
	log.Info().Msgf("Dumping verified names")
	if len(ds) > 0 {
		var suffix string
		if len(ds) > 1 {
			suffix = "s"
		}
		log.Info().Msgf("Names are filtered by %d data-source%s", len(ds), suffix)
	}

	id := 1
	limit := 100_000
	var outputs []output.Output
	var count int
	for id <= namesTotal {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			outputs, err = d.outputNames(id, limit, ds)
			count += len(outputs)
			if err != nil {
				return fmt.Errorf("-> outputNames %w", err)
			}
			ch <- outputs
			id += limit
			percent := float64(count) * 100 / float64(namesNum)
			fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", 80))
			fmt.Fprintf(os.Stderr, "Dumped %s names (%0.1f%%)",
				humanize.Comma(int64(count)),
				percent,
			)
		}
	}
	fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", 80))
	log.Info().Msgf("Dumped %s verified names",
		humanize.Comma(int64(namesNum)),
	)
	return nil
}

// DumpOccurrences reads data for detected and verified names and sends it to output channel.
func (d *dumpio) DumpOccurrences(ctx context.Context, ch chan<- []output.Output, ds []int) error {
	err := d.checkForVerifiedNames()
	if err != nil {
		err = fmt.Errorf("-> checkForVerifiedNames %w", err)
		log.Warn().Msg("Run `bhlindex verify` before `bhlindex dump`.")
		return err
	}

	_, _, occursTotal, err := d.stats(ds)
	if err != nil {
		return fmt.Errorf("-> stats %w", err)
	}
	log.Info().Msg("Dumping name occurrences")
	if len(ds) > 0 {
		var suffix string
		if len(ds) > 1 {
			suffix = "s"
		}
		log.Info().Msgf("Occurrences are filtered by %d data-source%s", len(ds), suffix)
	}
	id := 1
	limit := 100_000
	var count int
	var outputs []output.Output
	for id <= occursTotal {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			outputs, err = d.outputOccurs(id, limit, ds)
			if err != nil {
				return fmt.Errorf("-> outputOccurs %w", err)
			}
			count += len(outputs)
			ch <- outputs
			id += limit
			occursNum := id - 1
			percent := float64(occursNum) * 100 / float64(occursTotal)
			fmt.Fprintf(os.Stderr, "\r%s", strings.Repeat(" ", 80))
			fmt.Fprintf(os.Stderr, "\rDumped %s occurrences (%0.1f%%)",
				humanize.Comma(int64(occursNum)),
				percent,
			)
		}
	}
	fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", 80))
	log.Info().Msgf("Dumped %s occurrences",
		humanize.Comma(int64(count)),
	)
	return nil
}

func (d dumpio) DumpOddsVerification() ([]output.Output, error) {
	log.Info().Msg("Dumping odds vs verification percentage")
	return d.getOccurVerif()
}
