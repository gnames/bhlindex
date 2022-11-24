package loaderio

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gnames/bhlindex/internal/config"
	"github.com/gnames/bhlindex/internal/ent/item"
	"github.com/gnames/bhlindex/internal/ent/loader"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type loaderio struct {
	config.Config
	pagesTotal   int
	pagesCount   int
	ignoredItems map[int]struct{}
	db           *sql.DB
}

func New(cfg config.Config, db *sql.DB) loader.Loader {
	res := loaderio{
		Config:       cfg,
		ignoredItems: make(map[int]struct{}),
		db:           db,
	}
	return res
}

func (l loaderio) DetectPageDups() (loader.Loader, error) {
	var err error
	pageIDs := make(map[int]struct{})
	rootDir := l.BHLdir

	err = checkRoot(rootDir)
	if err != nil {
		err = fmt.Errorf("checkRoot: %w", err)
		return l, err
	}

	log.Info().Msg("Preprocessing to detect PageID duplicates")

	err = filepath.WalkDir(rootDir,
		func(path string, d fs.DirEntry, e error) error {
			err = e
			if d.IsDir() || !isPageFile(d.Name()) {
				return err
			}
			_, itemID, pageID, _ := parseFileName(path)
			if _, ok := pageIDs[pageID]; ok {
				l.ignoredItems[itemID] = struct{}{}
				return err
			}
			pageIDs[pageID] = struct{}{}
			l.pagesTotal++

			if l.pagesTotal%100_000 == 0 {
				pages := humanize.Comma(int64(l.pagesTotal))
				fmt.Fprintf(os.Stderr, "\r%s", strings.Repeat(" ", 80))
				fmt.Fprintf(os.Stderr, "\rPreprocessing page %s", pages)
			}

			return nil
		})

	pages := humanize.Comma(int64(l.pagesTotal))
	fmt.Fprint(os.Stderr, "\r")
	log.Info().Msgf("Preprocessed %s pages", pages)

	if len(l.ignoredItems) == 0 {
		log.Info().Msg("No PageID duplicates were found.")
	} else {
		for k := range l.ignoredItems {
			log.Warn().Msgf("Item %d will be ignored due to PageID duplicates.", k)
		}
	}
	return l, err
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
			err = fmt.Errorf("importItems: %w", err)
		}
		return err
	})

	for i := 0; i < 3; i++ {
		gProc.Go(func() error {
			err = l.processItemWorker(ctx, itemCh, dbItemCh)
			if err != nil {
				// This error is masked by gImp.Wait, so we
				// need to show it in a warning.
				err = fmt.Errorf("processItemWorker: %w", err)
				if !strings.Contains(err.Error(), "context canceled") {
					fmt.Fprint(os.Stderr, "\r")
					log.Warn().Err(err).Msg("")
				}
			}
			return err
		})
	}

	err = gImp.Wait()
	close(itemCh)
	if err != nil {
		return err
	}

	err = gProc.Wait()
	if err != nil {
		return err
	}
	fmt.Fprint(os.Stderr, "\r")
	log.Info().Msg("All items are loaded into name-finder.")
	return nil
}
