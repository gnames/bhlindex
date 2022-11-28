package loaderio

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gnames/bhlindex/internal/ent/item"
	"github.com/gnames/bhlindex/internal/ent/page"
	"github.com/gnames/gnfmt"
	"github.com/gnames/gnsys"
	"github.com/rs/zerolog/log"
)

// importItems starts from the root of BHL directory and traverses its children
// collecting directories that correspond to BHL items.
func (l loaderio) importItems(
	ctx context.Context,
	itemCh chan<- *item.Item,
) error {
	rootDir := l.BHLdir
	var err error
	err = checkRoot(rootDir)
	if err != nil {
		err = fmt.Errorf("-> checkRoot %w", err)
		return err
	}

	currentDir := ""
	var itm *item.Item
	var pages []*page.Page

	start := time.Now()
	var count int

	// Walk traverses files in lexical order. It means we do not need to
	// sort pages after they are collected.
	err = filepath.WalkDir(rootDir,
		func(path string, d fs.DirEntry, e error) error {
			err = e
			var pg *page.Page
			if d.IsDir() || !isPageFile(d.Name()) {
				return err
			}

			if dir := filepath.Dir(path); dir != currentDir {
				if itm != nil {
					itm.Pages = pages
					l.pagesCount += len(pages)
					count = l.countIncr(start, count)
					select {
					case <-ctx.Done():
						return ctx.Err()
					case itemCh <- itm:
					}
				}

				itm = itemFromPath(path)
				pg = pageFromPath(path)
				pages = []*page.Page{pg}
				currentDir = dir
			} else {
				pg = pageFromPath(path)
				pages = append(pages, pg)
			}
			return err
		})

	itm.Pages = pages
	select {
	case <-ctx.Done():
		return ctx.Err()
	case itemCh <- itm:
	}
	return err
}

// check if root BHL directory exists and is not empty.
func checkRoot(rootDir string) error {
	exists, empty, err := gnsys.DirExists(rootDir)
	if err != nil {
		return fmt.Errorf("DirExists: %w", err)
	}
	if !exists {
		return fmt.Errorf("directory '%s' does not exist", rootDir)
	}
	if empty {
		return fmt.Errorf("directory '%s' is empty", rootDir)
	}
	return nil
}

func (l loaderio) processItemWorker(
	ctx context.Context,
	itemCh <-chan *item.Item,
	dbItemCh chan<- *item.Item,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	var err error

	for itm := range itemCh {
		// ignore items with duplicated PageIDs
		if _, ok := l.ignoredItems[itm.ID]; ok {
			fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", 80))
			log.Warn().Msgf("Skipping Item %d because of duplicates", itm.ID)
			continue
		}

		err = updatePages(itm)
		if err != nil {
			err = fmt.Errorf("-> updatePages Item %d %w", itm.ID, err)
			log.Warn().Err(err).Msg("")
		}

		// if any go-routine returns an error, ctx will cancel all
		// other go-routines in the errgroup.
		select {
		case dbItemCh <- itm:
		case <-ctx.Done():
			log.Warn().Err(ctx.Err()).Msg("")
		}

	}
}

func (l loaderio) countIncr(start time.Time, count int) int {
	var pages, pagesPerHour, eta, percent string
	count++
	if count%10_000 == 0 {
		pages, pagesPerHour, eta, percent = l.pagesPerHour(start)
		fmt.Fprint(os.Stderr, "\r")
		log.Info().
			Msgf("Processed %s pages (%s), %s pages/hr: ETA %s",
				pages, percent, pagesPerHour, eta,
			)
	} else if count%100 == 0 {
		pages, pagesPerHour, eta, percent = l.pagesPerHour(start)
		fmt.Fprintf(os.Stderr, "\r%s", strings.Repeat(" ", 80))
		fmt.Fprintf(
			os.Stderr,
			"\rProcessed %s pages (%s), %s pages/hr: ETA %s",
			pages, percent, pagesPerHour, eta,
		)
	}
	return count
}

func (l loaderio) pagesPerHour(start time.Time) (string, string, string, string) {
	pages := humanize.Comma(int64(l.pagesCount))

	dur := float64(time.Since(start)) / float64(time.Hour)
	rate := float64(l.pagesCount) / dur
	pagesPerHour := humanize.Comma(int64(rate))
	eta := 3600 * float64(l.pagesTotal-l.pagesCount) / rate
	etaStr := gnfmt.TimeString(eta)
	percent := 100 * float64(l.pagesCount) / float64(l.pagesTotal)
	percentStr := fmt.Sprintf("%0.1f%%", percent)
	return pages, pagesPerHour, etaStr, percentStr
}

func itemFromPath(path string) *item.Item {
	path, _ = filepath.Split(path)
	_, itm := filepath.Split(path[0 : len(path)-1])
	id, err := strconv.Atoi(itm)
	if err != nil {
		log.Warn().Err(err).Msgf("cannot convert '%s' to int", itm)
	}
	return &item.Item{ID: id, Path: path}
}
