package bhlindex

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/gnames/bhlindex/internal/config"
	"github.com/gnames/bhlindex/internal/ent/finder"
	"github.com/gnames/bhlindex/internal/ent/item"
	"github.com/gnames/bhlindex/internal/ent/loader"
	"github.com/gnames/bhlindex/internal/ent/name"
	"github.com/gnames/bhlindex/internal/ent/output"
	"github.com/gnames/bhlindex/internal/ent/verif"
	"github.com/gnames/gnfmt"
	"github.com/gnames/gnlib/ent/gnvers"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type bhlindex struct {
	config.Config
	loader.Loader
	finder.Finder
}

// New sets up BHLindex interface using bhlindex instance.
func New(
	cfg config.Config,
) BHLindex {
	return &bhlindex{
		Config: cfg,
	}
}

// FindNames traverses BHL corpus directory structure, assembling texts,
// detecting names, saving data to storage.
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
		return fmt.Errorf("FindNames: %w", err)
	}

	wgFind.Wait()
	close(namesCh)

	err = gSave.Wait()
	if err != nil {
		return fmt.Errorf("FindNames: %w", err)
	}

	log.Info().Msg("Finding names finished successfully")
	return nil
}

// Verify names runs verification on unique detected names and saves the
// results to a local storage.
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

func (bi *bhlindex) DumpPages(dmp output.Dumper) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan []output.OutputPage)
	gDump, ctxDump := errgroup.WithContext(ctx)
	gOut, ctxOut := errgroup.WithContext(ctx)

	gDump.Go(func() error {
		return dmp.DumpPages(ctxDump, ch)
	})

	gOut.Go(func() error {
		return processOutput(bi, ctxOut, ch)
	})

	err := gDump.Wait()
	if err != nil {
		return fmt.Errorf("dumpPages: %w", err)
	}
	close(ch)

	err = gOut.Wait()
	if err != nil {
		return fmt.Errorf("dumpPages: %w", err)
	}
	return nil
}

func (bi *bhlindex) DumpNames(dmp output.Dumper) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan []output.OutputName)
	gDump, ctxDump := errgroup.WithContext(ctx)
	gOut, ctxOut := errgroup.WithContext(ctx)

	gDump.Go(func() error {
		return dmp.DumpNames(ctxDump, ch, bi.OutputDataSourceIDs)
	})

	gOut.Go(func() error {
		return processOutput(bi, ctxOut, ch)
	})

	err := gDump.Wait()
	if err != nil {
		return fmt.Errorf("bhlindex: %w", err)
	}
	close(ch)

	err = gOut.Wait()
	if err != nil {
		return fmt.Errorf("bhlindex: %w", err)
	}
	return nil
}

// DumpOccurrences creates output with detected and verified names in CSV,
// TSV, or JSON formats.
func (bi *bhlindex) DumpOccurrences(dmp output.Dumper) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan []output.OutputOccurrence)
	gDump, ctxDump := errgroup.WithContext(ctx)
	gOut, ctxOut := errgroup.WithContext(ctx)

	gDump.Go(func() error {
		return dmp.DumpOccurrences(ctxDump, ch, bi.OutputDataSourceIDs)
	})

	gOut.Go(func() error {
		return processOutput(bi, ctxOut, ch)
	})

	err := gDump.Wait()
	if err != nil {
		return fmt.Errorf("bhlindex: %w", err)
	}
	close(ch)

	err = gOut.Wait()
	if err != nil {
		return fmt.Errorf("bhlindex: %w", err)
	}
	return nil
}

// GetVersion outputs the version of BHLindex.
func (bi *bhlindex) GetVersion() gnvers.Version {
	return gnvers.Version{Version: Version, Build: Build}
}

// GetConfig returns an instance of configuration fields.
func (bi *bhlindex) GetConfig() config.Config {
	return bi.Config
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

func (bi *bhlindex) extension() string {
	switch bi.OutputFormat {
	case gnfmt.CSV:
		return ".csv"
	case gnfmt.TSV:
		return ".tsv"
	default:
		return ".json"
	}
}

func processOutput[O output.Output](
	bi *bhlindex,
	ctx context.Context,
	ch <-chan []O,
) error {
	var o O
	path := filepath.Join(bi.OutputDir, o.Name()+bi.extension())
	w, err := os.Create(path)
	if err != nil {
		return err
	}
	if bi.OutputFormat != gnfmt.CompactJSON {
		_, err = w.WriteString(output.CSVHeader(o, bi.OutputFormat) + "\n")
		if err != nil {
			return err
		}
	}
	var count int
	for rows := range ch {
		count++
		if count%500_000 == 0 {
			log.Info().Msgf("Processed %d %s", count, o.Name())
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("processOutput: %w", ctx.Err())
		default:
			for i := range rows {
				_, err = w.WriteString(output.Format(rows[i], bi.OutputFormat) + "\n")
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
