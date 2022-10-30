package loaderio

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gnames/bhlindex/internal/ent/item"
	"github.com/gnames/bhlindex/internal/ent/page"
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
	err := checkRoot(rootDir)
	if err != nil {
		return fmt.Errorf("importItems: %w", err)
	}

	currentDir := ""
	var itm *item.Item
	var pages []*page.Page

	start := time.Now()
	var count int

	// Walk traverses files in lexical order. It means we do not need to
	// sort pages after they are collected.
	err = filepath.WalkDir(rootDir,
		func(path string, d fs.DirEntry, _ error) error {
			var err error
			var pg *page.Page
			if d.IsDir() || !isPageFile(d.Name()) {
				return nil
			}

			if dir := filepath.Dir(path); dir != currentDir {
				if itm != nil {
					itm.Pages = pages
					select {
					case <-ctx.Done():
						return ctx.Err()
					case itemCh <- itm:
					}
				}

				itm = itemFromPath(path)
				pg, err = pageFromPath(path)
				if err != nil {
					return fmt.Errorf("WalkDir: %w", err)
				}
				pages = []*page.Page{pg}
				currentDir = dir
				count = countIncr(start, count)
			} else {
				pg, err = pageFromPath(path)
				if err != nil {
					return fmt.Errorf("WalkDir: %w", err)
				}
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
		return fmt.Errorf("checkRoot: %w", err)
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
) error {
	var err error
	for itm := range itemCh {

		err = l.insertItem(itm)
		if err != nil {
			return fmt.Errorf("processItemWorker: %w", err)
		}
		if itm.ID == 0 {
			continue
		}

		err = updatePages(itm)
		if err != nil {
			return fmt.Errorf("processItemWorker: %w", err)
		}

		err = l.insertPages(itm)
		if err != nil {
			return fmt.Errorf("processItemWorker: %w", err)
		}

		// if any go-routine returns an error, ctx will cancel all
		// other go-routines in the errgroup.
		select {
		case <-ctx.Done():
			log.Info().Err(ctx.Err()).Msg("Item processing cancelled")
			return ctx.Err()
		case dbItemCh <- itm:
		}

	}
	return nil
}

func countIncr(start time.Time, count int) int {
	count++
	if count%10_000 == 0 {
		fmt.Fprint(os.Stderr, "\r")
		log.Info().
			Str("items", humanize.Comma(int64(count))).
			Str("items/hour", itemsPerHour(count, start)).
			Msg("Finding names in BHL items")
	} else if count%100 == 0 {
		fmt.Fprintf(os.Stderr, "\r%s", strings.Repeat(" ", 40))
		fmt.Fprintf(
			os.Stderr,
			"\rProcessed %s items, %s items/hour",
			humanize.Comma(int64(count)),
			itemsPerHour(count, start),
		)
	}
	return count
}

func itemsPerHour(itemsNum int, start time.Time) string {
	dur := float64(time.Since(start)) / float64(time.Hour)
	rate := float64(itemsNum) / dur
	return humanize.Comma(int64(rate))
}

func itemFromPath(path string) *item.Item {
	path, _ = filepath.Split(path)
	_, itm := filepath.Split(path[0 : len(path)-1])
	return &item.Item{Path: path, InternetArchiveID: itm}
}
