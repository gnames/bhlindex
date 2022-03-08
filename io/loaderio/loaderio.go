package loaderio

import (
	"context"
	"database/sql"

	"github.com/gnames/bhlindex/config"
	"github.com/gnames/bhlindex/ent/item"
	"github.com/gnames/bhlindex/ent/loader"
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

// LoadItems saves items into databsase and removes
// duplicates. This process preceeds actual work on name resolution. After
// a item is impored to the database its id goes into itemIDs channel
func (l loaderio) LoadItems(ctx context.Context, dbItemCh chan<- *item.Item) error {
	itemCh := make(chan *item.Item)

	// in case of an error ctx will send a signal to kill all workers
	gImp := new(errgroup.Group)
	gProc, ctx := errgroup.WithContext(ctx)
	log.Info().Int("gnfinderJobs", l.Jobs).Msg("Finding names in BHL items")

	gImp.Go(func() error {
		return l.importItems(ctx, itemCh)
	})

	for i := 0; i < 3; i++ {
		gProc.Go(func() error {
			return l.processItemWorker(ctx, itemCh, dbItemCh)
		})
	}

	err := gImp.Wait()
	close(itemCh)
	if err != nil {
		return err
	}

	err = gProc.Wait()
	return err
}
