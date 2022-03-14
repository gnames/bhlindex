package bhlindex

import (
	"context"
	"sync"

	"github.com/gnames/bhlindex/config"
	"github.com/gnames/bhlindex/ent/finder"
	"github.com/gnames/bhlindex/ent/item"
	"github.com/gnames/bhlindex/ent/loader"
	"github.com/gnames/bhlindex/ent/name"
	"github.com/gnames/bhlindex/ent/verif"
	"github.com/gnames/gnlib/ent/gnvers"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type bhlindex struct {
	config.Config
	loader.Loader
	finder.Finder
}

func New(
	cfg config.Config,
) BHLindex {
	return &bhlindex{
		Config: cfg,
	}
}

func (bi *bhlindex) FindNames(
	ldr loader.Loader,
	fdr finder.Finder,
) error {
	itemCh := make(chan *item.Item, 10)
	namesCh := make(chan []name.DetectedName)
	var wgFind sync.WaitGroup
	gLoad, ctx := errgroup.WithContext(context.Background())
	gSave := new(errgroup.Group)

	gLoad.Go(func() error {
		return ldr.LoadItems(ctx, itemCh)
	})

	wgFind.Add(1)
	go fdr.FindNames(itemCh, namesCh, &wgFind)

	gSave.Go(func() error {
		return fdr.SaveNames(ctx, namesCh)
	})

	err := gLoad.Wait()
	close(itemCh)
	if err != nil {
		return err
	}

	wgFind.Wait()
	close(namesCh)

	return gSave.Wait()
}

func (bi *bhlindex) VerifyNames(vrf verif.VerifierBHL) (err error) {
	err = vrf.Reset()
	if err == nil {
		err = vrf.ExtractUniqueNames()
	}

	if err == nil {
		err = vrf.Verify()
	}
	return err
}

func (bi *bhlindex) GetVersion() gnvers.Version {
	return gnvers.Version{Version: Version, Build: Build}
}

func counterLog(counter <-chan int) {
	count := 0
	for i := range counter {
		count += i
		if count%10000 == 0 {
			log.Info().Msgf("Processed %d items", count)
		}
	}
}

func (bi *bhlindex) GetConfig() config.Config {
	return bi.Config
}
