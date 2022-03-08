package finderio

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/gnames/bhlindex/config"
	"github.com/gnames/bhlindex/ent/finder"
	"github.com/gnames/bhlindex/ent/item"
	"github.com/gnames/bhlindex/ent/name"
	"github.com/gnames/gnfinder"
	gnfcfg "github.com/gnames/gnfinder/config"
	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/ent/output"
	"github.com/gnames/gnfinder/io/dict"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
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
			log.Warn().Err(err).Msg("Cannot save detected names")
			return err
		}
	}
	return nil
}

func (fdr finderio) savePageNameStrings(names []name.DetectedName) error {
	now := time.Now()
	columns := []string{"page_id", "item_id", "name_string", "annot_nomen",
		"annot_nomen_type", "offset_start", "offset_end",
		"ends_next_page", "odds_log10", "cardinality", "updated_at"}
	transaction, err := fdr.db.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	stmt, err := transaction.Prepare(pq.CopyIn("detected_names", columns...))
	if err != nil {
		return err
	}

	for _, v := range names {
		_, err = stmt.Exec(v.PageID, v.ItemID, v.NameString,
			v.AnnotNomen, v.AnnotNomenType, v.OffsetStart, v.OffsetEnd,
			v.EndsNextPage, v.OddsLog10, v.Cardinality, now)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return transaction.Commit()
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
