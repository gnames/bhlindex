package finderio

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/gnames/bhlindex/config"
	"github.com/gnames/bhlindex/ent/finder"
	"github.com/gnames/bhlindex/ent/item"
	"github.com/gnames/bhlindex/ent/name"
	"github.com/gnames/gnfinder"
	gnfcfg "github.com/gnames/gnfinder/config"
	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/ent/output"
	"github.com/gnames/gnfinder/io/dict"
)

type finderio struct {
	config.Config
	db *sql.DB
}

func New(cfg config.Config, db *sql.DB) finder.Finder {
	return finderio{Config: cfg, db: db}
}

func (fdr finderio) FindNames(
	itemCh <-chan *item.Item,
	namesCh chan<- []name.DetectedName,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	var wgWrkr sync.WaitGroup
	wgWrkr.Add(fdr.Jobs)
	for i := 0; i < fdr.Jobs; i++ {
		go fdr.finderWorker(itemCh, namesCh, &wgWrkr)
	}
	wgWrkr.Wait()
}

func (fdr finderio) SaveNames(
	ctx context.Context,
	namesCh <-chan []name.DetectedName,
) error {
	for v := range namesCh {
		_ = v
		err := fdr.savePageNameStrings(v)
		if err != nil {
			return fmt.Errorf("SaveNames: %w", err)
		}
	}
	return nil
}

func (fdr finderio) finderWorker(
	chItem <-chan *item.Item,
	chNames chan<- []name.DetectedName,
	wgWrkr *sync.WaitGroup,
) {
	defer wgWrkr.Done()

	opts := []gnfcfg.Option{
		gnfcfg.OptWithPlainInput(true),
	}
	cfg := gnfcfg.New(opts...)
	d := dict.LoadDictionary()
	weights := nlp.BayesWeights()
	gnf := gnfinder.New(cfg, d, weights)
	for itm := range chItem {
		o := gnf.Find("", string(itm.Text))
		chNames <- namesToDetectedNames(itm, o.Names)
	}
}

func namesToDetectedNames(itm *item.Item, names []output.Name) []name.DetectedName {
	ns := make([]name.DetectedName, len(names))
	j := 0
	if j >= len(names) {
		return ns
	}
	n := names[j]
	for _, page := range itm.Pages {
		for {
			if n.OffsetStart <= page.OffsetNext {
				ns[j] = name.New(itm.ID, page, n)
				j++
				if j >= len(names) {
					return ns
				}
				n = names[j]
			} else {
				break
			}
		}
	}
	return ns
}
