package bhlindex

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	var err error
	itemCh := make(chan *item.Item, 10)
	namesCh := make(chan []name.DetectedName)
	var wgFind sync.WaitGroup
	gLoad, ctx := errgroup.WithContext(context.Background())
	gSave := new(errgroup.Group)

	// check for problems and find the number of pages
	ldr, err = ldr.DetectPageDups()
	if err != nil {
		err = fmt.Errorf("-> DetectPageDups %w", err)
		return err
	}

	gLoad.Go(func() error {
		err = ldr.LoadItems(ctx, itemCh)
		if err != nil {
			err = fmt.Errorf("-> LoadItems %w", err)
		}
		return err
	})

	wgFind.Add(1)
	go fdr.FindNames(itemCh, namesCh, &wgFind)

	gSave.Go(func() error {
		err = fdr.SaveNames(ctx, namesCh)
		if err != nil {
			err = fmt.Errorf("-> SaveNames %w", err)
		}
		return err
	})

	err = gLoad.Wait()
	close(itemCh)
	if err != nil {
		return err
	}

	wgFind.Wait()
	close(namesCh)

	err = gSave.Wait()
	if err != nil {
		return err
	}

	log.Info().Msg("Finding names finished successfully")
	return nil
}

// Verify names runs verification on unique detected names and saves the
// results to a local storage.
func (bi *bhlindex) VerifyNames(vrf verif.VerifierBHL) (err error) {
	err = vrf.Reset()
	if err != nil {
		err = fmt.Errorf("-> Reset %w", err)
		return err
	}

	err = vrf.ExtractUniqueNames()
	if err != nil {
		err = fmt.Errorf("-> ExtractUniqueNames %w", err)
		return err
	}

	err = vrf.Verify()
	if err != nil {
		err = fmt.Errorf("-> Verify %w", err)
	}

	return err
}

func (bi *bhlindex) CalcOddsVerif(vrf verif.VerifierBHL) (err error) {
	return vrf.CalcOddsVerif()
}

func (bi *bhlindex) DumpNames(dmp output.Dumper) error {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan []output.Output)
	gDump, ctxDump := errgroup.WithContext(ctx)
	gOut, ctxOut := errgroup.WithContext(ctx)

	gDump.Go(func() error {
		err = dmp.DumpNames(ctxDump, ch, bi.OutputDataSourceIDs)
		if err != nil {
			err = fmt.Errorf("-> DumpNames %w", err)
		}
		return err
	})

	gOut.Go(func() error {
		err = processOutput(bi, ctxOut, ch)
		if err != nil {
			err = fmt.Errorf("-> processOutput %w", err)
		}
		return err
	})

	err = gDump.Wait()
	if err != nil {
		return err
	}
	close(ch)

	err = gOut.Wait()
	if err != nil {
		return err
	}
	return nil
}

// DumpOccurrences creates output with detected and verified names in CSV,
// TSV, or JSON formats.
func (bi *bhlindex) DumpOccurrences(dmp output.Dumper) error {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan []output.Output)
	gDump, ctxDump := errgroup.WithContext(ctx)
	gOut, ctxOut := errgroup.WithContext(ctx)

	gDump.Go(func() error {
		err = dmp.DumpOccurrences(ctxDump, ch, bi.OutputDataSourceIDs)
		if err != nil {
			err = fmt.Errorf("-> DumpOccurrences %w", err)
		}
		return err

	})

	gOut.Go(func() error {
		err = processOutput(bi, ctxOut, ch)
		if err != nil {
			err = fmt.Errorf("-> processOutput %w", err)
		}
		return err
	})

	err = gDump.Wait()
	if err != nil {
		return err
	}
	close(ch)

	err = gOut.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (bi *bhlindex) DumpOddsVerification(dmp output.Dumper) error {
	var w *os.File
	outs, err := dmp.DumpOddsVerification()
	if err != nil {
		err = fmt.Errorf("-> DumpOddsVerification: %w", err)
		return err
	}
	for i := range outs {
		o := outs[i]
		if i == 0 {
			path := filepath.Join(bi.OutputDir, o.Name()+bi.extension())
			w, err = os.Create(path)
			if err != nil {
				return err
			}
			if bi.OutputFormat != gnfmt.CompactJSON {
				_, err = w.WriteString(output.CSVHeader(o, bi.OutputFormat) + "\n")
				if err != nil {
					return err
				}
			}
		}
		_, err = w.WriteString(output.Format(outs[i], bi.OutputFormat) + "\n")
		if err != nil {
			return err
		}
	}
	log.Info().Msgf("Dumped %d odds/verification records", len(outs))
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
	var err error
	var w *os.File
	var count int
	for rows := range ch {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for i := range rows {
				o = rows[i]
				if count == 0 {
					path := filepath.Join(bi.OutputDir, o.Name()+bi.extension())
					w, err = os.Create(path)
					if err != nil {
						return err
					}
					if bi.OutputFormat != gnfmt.CompactJSON {
						_, err = w.WriteString(output.CSVHeader(o, bi.OutputFormat) + "\n")
						if err != nil {
							return err
						}
					}
				}
				count++
				if count%25_000_000 == 0 {
					fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", 80))
					log.Info().Msgf("Processed %d %s", count, o.Name())
				}
				_, err = w.WriteString(output.Format(rows[i], bi.OutputFormat) + "\n")
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
