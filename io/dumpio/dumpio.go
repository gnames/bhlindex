package dumpio

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gnames/bhlindex/config"
	"github.com/gnames/bhlindex/ent/output"
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

// Dump reads data for detected and verified names and sends it to output channel.
func (d *dumpio) Dump(ctx context.Context, ch chan<- []output.Output) error {
	err := d.checkForVerifiedNames()
	if err != nil {
		err = fmt.Errorf("dumpio.Dump: %w", err)
		log.Warn().Msg("Run `bhlindex verify` before `bhlindex dump`.")
		return err
	}

	namesTotal, occTotal, itemsTotal, err := d.stats()
	if err != nil {
		return fmt.Errorf("dumpio.Dump: %w", err)
	}
	log.Info().Msgf("Dumping %s occurrences of %s names found in %s items",
		humanize.Comma(int64(occTotal)),
		humanize.Comma(int64(namesTotal)),
		humanize.Comma(int64(itemsTotal)),
	)
	id := 1
	limit := 100
	var outputs []output.Output
	for id <= itemsTotal {
		select {
		case <-ctx.Done():
			return fmt.Errorf("dumpio.Dump: %w", ctx.Err())
		default:
			outputs, err = d.outputs(id, limit)
			if err != nil {
				return fmt.Errorf("dumpio.Dump: %w", err)
			}
			ch <- outputs
			id += limit
			itemsNum := id - 1
			fmt.Fprintf(os.Stderr, "\r%s", strings.Repeat(" ", 35))
			fmt.Fprintf(os.Stderr, "\rDumped names from %s items", humanize.Comma(int64(itemsNum)))
			if itemsNum%10_000 == 0 {
				makeLog(itemsTotal, itemsNum)
			}
		}
	}
	fmt.Fprint(os.Stderr, "\r")
	log.Info().Msgf("Dumped %s occurrences of %s names",
		humanize.Comma(int64(occTotal)),
		humanize.Comma(int64(namesTotal)),
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
