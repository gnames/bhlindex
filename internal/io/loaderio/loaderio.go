package loaderio

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
		err = fmt.Errorf("-> checkRoot %w", err)
		return l, err
	}

	log.Info().Msg("Pre-processing to detect PageID duplicates")

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
				fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", 80))
				fmt.Fprintf(os.Stderr, "Pre-processing page %s", pages)
			}

			return nil
		})

	pages := humanize.Comma(int64(l.pagesTotal))
	fmt.Fprint(os.Stderr, "\r")
	log.Info().Msgf("Pre-processed %s pages", pages)

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
	var wg sync.WaitGroup
	itemJobs := 3
	wg.Add(itemJobs)

	log.Info().Int("gnfinderJobs", l.Jobs).Msg("Finding names in BHL items")

	gImp.Go(func() error {
		err = l.importItems(context.Background(), itemCh)
		if err != nil {
			err = fmt.Errorf("-> importItems %w", err)
		}
		return err
	})

	for i := 0; i < itemJobs; i++ {
		go l.processItemWorker(ctx, itemCh, dbItemCh, &wg)
	}

	err = gImp.Wait()
	close(itemCh)
	if err != nil {
		return err
	}

	wg.Wait()
	fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", 80))
	log.Info().Msgf("Processed %s pages", humanize.Comma(int64(l.pagesTotal)))
	return nil
}
