package verifio

import (
	"context"
	"database/sql"
	"fmt"
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

// Reset cleans up all stored data for verifications.
func (vrf verifio) Reset() error {
	log.Info().Msg("Cleaning up previous verification results if they exist")
	return vrf.truncateVerifTables()
}

// Verify verifies all detected names and stores the data localy.
func (vrf verifio) Verify() error {
	err := vrf.checkForDetectedNames()
	if err != nil {
		return fmt.Errorf("Verify: %w", err)
	}

	namesNum, err := vrf.numberOfNames()
	if err != nil {
		return fmt.Errorf("Verify: %w", err)
	}
	log.Info().Msgf("Verifying %s names", humanize.Comma(int64(namesNum)))

	vrfREST := verifrest.New(vrf.cfg.VerifierURL)
	ctx, cancel := context.WithCancel(context.Background())
	// this cancel will stop all child contexts
	defer cancel()

	start := time.Now()
	chIn := make(chan []name.UniqueName)
	chOut := make(chan []name.VerifiedName)
	gLoad, ctxLoad := errgroup.WithContext(ctx)
	gSave, ctxSave := errgroup.WithContext(ctx)
	var wg sync.WaitGroup
	wg.Add(1)

	gLoad.Go(func() error {
		err = vrf.loadNames(ctxLoad, namesNum, chIn)
		return err
	})

	go vrf.sendToVerify(ctx, vrfREST, chIn, chOut, &wg)

	gSave.Go(func() error {
		return vrf.saveVerif(ctxSave, chOut, namesNum, start)
	})

	err = gLoad.Wait()
	if err != nil {
		return fmt.Errorf("Verify: %w", err)
	}
	close(chIn)

	wg.Wait()
	close(chOut)

	gSave.Wait()
	if err != nil {
		return fmt.Errorf("Verify: %w", err)
	}

	return vrf.setPrimaryKey()
}
