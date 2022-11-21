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
func (d *dumpio) DumpNames(ctx context.Context, ch chan<- []output.OutputName, ds []int) error {
	err := d.checkForVerifiedNames()
	if err != nil {
		err = fmt.Errorf("dumpio.DumpNames: %w", err)
		log.Warn().Msg("Run `bhlindex verify` before `bhlindex dump`.")
		return err
	}

	namesTotal, namesNum, itemsTotal, err := d.stats(ds)
	if err != nil {
		return fmt.Errorf("dumpio.DumpNames: %w", err)
	}
	log.Info().Msgf("Dumping %s names found in %s items",
		humanize.Comma(int64(namesNum)),
		humanize.Comma(int64(itemsTotal)),
	)
	if len(ds) > 0 {
		var suffix string
		if len(ds) > 1 {
			suffix = "s"
		}
		log.Info().Msgf("Names are filtered by %d data-source%s", len(ds), suffix)
	}

	id := 1
	limit := 100_000
	var outputs []output.OutputName
	var count int
	for id <= namesTotal {
		select {
		case <-ctx.Done():
			return fmt.Errorf("dumpio.DumpNames: %w", ctx.Err())
		default:
			outputs, err = d.outputNames(id, limit, ds)
			count += len(outputs)
			if err != nil {
				return fmt.Errorf("dumpio.DumpNames: %w", err)
			}
			ch <- outputs
			id += limit
			fmt.Fprintf(os.Stderr, "\r%s", strings.Repeat(" ", 35))
			fmt.Fprintf(os.Stderr, "\rDumped %s names out of %s",
				humanize.Comma(int64(count)),
				humanize.Comma(int64(namesNum)))
		}
	}
	fmt.Fprintf(os.Stderr, "\r%s", strings.Repeat(" ", 35))
	fmt.Fprint(os.Stderr, "\r")
	log.Info().Msgf("Dumped %s names%s",
		humanize.Comma(int64(namesNum)),
		strings.Repeat(" ", 15),
	)
	return nil
}

// DumpOccurrences reads data for detected and verified names and sends it to output channel.
func (d *dumpio) DumpOccurrences(ctx context.Context, ch chan<- []output.OutputOccurrence, ds []int) error {
	err := d.checkForVerifiedNames()
	if err != nil {
		err = fmt.Errorf("dumpio.DumpOccurrences: %w", err)
		log.Warn().Msg("Run `bhlindex verify` before `bhlindex dump`.")
		return err
	}

	_, namesNum, itemsTotal, err := d.stats(ds)
	if err != nil {
		return fmt.Errorf("dumpio.DumpOccurrences: %w", err)
	}
	log.Info().Msgf("Dumping occurrences of %s names found in %s items",
		humanize.Comma(int64(namesNum)),
		humanize.Comma(int64(itemsTotal)),
	)
	if len(ds) > 0 {
		var suffix string
		if len(ds) > 1 {
			suffix = "s"
		}
		log.Info().Msgf("Occurrences are filtered by %d data-source%s", len(ds), suffix)
	}
	id := 1
	limit := 100
	var count int
	var outputs []output.OutputOccurrence
	for id <= itemsTotal {
		select {
		case <-ctx.Done():
			return fmt.Errorf("dumpio.DumpOccurrences: %w", ctx.Err())
		default:
			outputs, err = d.outputOccurs(id, limit, ds)
			if err != nil {
				return fmt.Errorf("dumpio.DumpOccurrences: %w", err)
			}
			count += len(outputs)
			ch <- outputs
			id += limit
			itemsNum := id - 1
			fmt.Fprintf(os.Stderr, "\r%s", strings.Repeat(" ", 35))
			fmt.Fprintf(os.Stderr, "\rDumped occurrences from %s items", humanize.Comma(int64(itemsNum)))
			if itemsNum%10_000 == 0 {
				makeLog(itemsTotal, itemsNum)
			}
		}
	}
	fmt.Fprint(os.Stderr, "\r")
	log.Info().Msgf("Dumped %s occurrences from %s items",
		humanize.Comma(int64(count)),
		humanize.Comma(int64(itemsTotal)),
	)
	return nil
}

func makeLog(itemsTotal, itemsNum int) {
	fmt.Fprint(os.Stderr, "\r")
	items := humanize.Comma(int64(itemsNum))
	total := humanize.Comma(int64(itemsTotal))
	percent := 100 * float64(itemsNum) / float64(itemsTotal)
	log.Info().
		Str("items", items+"/"+total).
		Msgf("Dumped %0.1f%%", percent)
}
