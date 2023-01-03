package verifio

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gnames/bhlindex/internal/config"
	"github.com/gnames/bhlindex/internal/ent/name"
	"github.com/gnames/bhlindex/internal/ent/verif"
	"github.com/gnames/gnverifier/io/verifrest"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type verifio struct {
	cfg config.Config
	db  *sql.DB
}

// New returns an instance of VerifierBHL.
func New(cfg config.Config, db *sql.DB) verif.VerifierBHL {
	res := verifio{
		cfg: cfg,
		db:  db,
	}
	return res
}

func (vrf verifio) CalcOddsVerif() error {
	log.Info().Msg("Calculating relationship between odds and verification")
	maxOdds, err := vrf.maxOdds()
	if err != nil {
		return err
	}
	for i := 0; i < int(math.Floor(maxOdds)); i++ {
		err = vrf.oddsVerif(i, i+1)
		if err != nil {
			return err
		}
	}
	log.Info().Msg("Odds vs verification percentage (without abbreviated names)")

	return vrf.showOddsVerif()
}

// Reset cleans up all stored data for verifications.
func (vrf verifio) Reset() error {
	log.Info().Msg("Cleaning up previous verification results if they exist")
	err := vrf.truncateVerifTables()
	if err != nil {
		err = fmt.Errorf("-> truncateVerifTables %w", err)
	}
	return err
}

// Verify verifies all detected names and stores the data localy.
func (vrf verifio) Verify() error {
	err := vrf.checkForDetectedNames()
	if err != nil {
		return fmt.Errorf("-> checkForDetectedNames %w", err)
	}
	namesNum, err := vrf.numberOfNames()
	if err != nil {
		return fmt.Errorf("-> numberOfNames %w", err)
	}
	log.Info().Msgf("Verifying %s names", humanize.Comma(int64(namesNum)))

	vrfREST := verifrest.New(vrf.cfg.VerifierURL)
	ctx, cancel := context.WithCancel(context.Background())
	// this cancel will stop all child contexts
	defer cancel()

	start := time.Now()
	chIn := make(chan []name.UniqueName)
	chOut := make(chan verifiedBatch)
	gLoad, ctxLoad := errgroup.WithContext(ctx)
	gSave, ctxSave := errgroup.WithContext(ctx)
	var wg sync.WaitGroup
	wg.Add(1)

	gLoad.Go(func() error {
		err = vrf.loadNames(ctxLoad, namesNum, chIn)
		if err != nil {
			err = fmt.Errorf("-> loadNames %w", err)
		}
		return err
	})

	go vrf.sendToVerify(ctx, vrfREST, chIn, chOut, &wg)

	gSave.Go(func() error {
		err = vrf.saveVerif(ctxSave, chOut, namesNum, start)
		if err != nil {
			err = fmt.Errorf("-> saveVerif %w", err)
		}
		return err
	})

	err = gLoad.Wait()
	if err != nil {
		return err
	}
	close(chIn)

	wg.Wait()
	close(chOut)

	gSave.Wait()
	if err != nil {
		return err
	}

	return nil
}
