package loaderio

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/gnames/bhlindex/internal/config"
	"github.com/gnames/bhlindex/internal/ent/item"
	"github.com/gnames/bhlindex/internal/ent/loader"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type loaderio struct {
	config.Config
	db *sql.DB
}

func New(cfg config.Config, db *sql.DB) loader.Loader {
	res := loaderio{
		Config: cfg,
		db:     db,
	}
	return res
}

// LoadItems saves items into databsase duplicates. This process preceeds
// actual work on name resolution. After a item is impored to the database its
// id goes into itemIDs channel
func (l loaderio) LoadItems(ctx context.Context, dbItemCh chan<- *item.Item) error {
	var err error
	itemCh := make(chan *item.Item)

	// in case of an error ctx will send a signal to kill all workers
	gImp := errgroup.Group{}
	gProc, ctx := errgroup.WithContext(ctx)
	log.Info().Int("gnfinderJobs", l.Jobs).Msg("Finding names in BHL items")

	gImp.Go(func() error {
		err = l.importItems(ctx, itemCh)
		if err != nil {
			log.Warn().Err(err).Msg("importItems")
		}
		return err
	})

	for i := 0; i < 3; i++ {
		gProc.Go(func() error {
			err = l.processItemWorker(ctx, itemCh, dbItemCh)
			if err != nil {
				log.Warn().Err(err).Msg("processItemWorker")
			}
			return err
		})
	}

	err = gImp.Wait()
	close(itemCh)
	if err != nil {
		log.Warn().Err(err).Msg("gImp.Wait:")
		return fmt.Errorf("LoadItems: %w", err)
	}

	err = gProc.Wait()
	if err != nil {
		log.Warn().Err(err).Msg("gProc.Wait:")
		return fmt.Errorf("LoadItems: %w", err)
	}
	fmt.Fprint(os.Stderr, "\r")
	log.Info().Msg("All items are loaded into name-finder.")
	return nil
}
