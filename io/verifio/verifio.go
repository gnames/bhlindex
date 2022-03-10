package verifio

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/gnames/bhlindex/config"
	"github.com/gnames/bhlindex/ent/name"
	"github.com/gnames/bhlindex/ent/verif"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type verificationStatus int

const (
	Init verificationStatus = iota
	AllNamesHarvested
	AllNamesVerified
	AllErrorsProcessed
)

var verJobs = 5

type metaData struct {
	Status        verificationStatus
	StaleErrTries int
}

type verifio struct {
	cfg config.Config
	db  *sql.DB
	metaData
}

func New(cfg config.Config, db *sql.DB) verif.VerifierBHL {
	res := verifio{
		cfg: cfg,
		db:  db,
	}
	return res
}

func (vrf verifio) Verify() error {
	noNames, err := vrf.noDetectedNames()
	if err != nil {
		return err
	}
	if noNames {
		err = errors.New("detected_names table is empty")
		log.Warn().Err(err).Msg("Run 'bhlindex find' before 'bhlindex verify'")
		return err
	}

	err = vrf.dedupeNames()
	if err != nil {
		return err
	}

	return vrf.verifyNames()
}

func (vrf verifio) Reset() error {
	log.Info().Msg("Cleaning up previous verification results")
	return vrf.truncateVerifTables()
}

func (vrf verifio) verifyNames() error {
	namesNum, err := vrf.numberOfNames()
	if err != nil {
		return err
	}
	log.Info().Msgf("Verifying %d names", namesNum)

	start := time.Now()
	chNames := make(chan []name.UniqueName)
	chVer := make(chan []name.VerifiedName)
	gLoad, ctx := errgroup.WithContext(context.Background())
	gSave := new(errgroup.Group)
	var wg sync.WaitGroup
	wg.Add(verJobs)

	gLoad.Go(func() error {
		return vrf.loadNames(ctx, namesNum, chNames, start)
	})

	for i := 0; i < verJobs; i++ {
		go vrf.sendToVerify(ctx, chNames, chVer, &wg)
	}

	gSave.Go(func() error {
		return vrf.saveVerif(ctx, chVer)
	})

	err = gLoad.Wait()
	if err != nil {
		return err
	}
	close(chNames)

	wg.Wait()
	close(chVer)

	gSave.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (vrf verifio) noDetectedNames() (bool, error) {
	var page_id string
	q := "select page_id from detected_names limit 1"
	err := vrf.db.QueryRow(q).Scan(&page_id)
	return page_id == "", err
}

func (vrf verifio) numberOfNames() (int, error) {
	q := "select count(*) from unique_names"
	var namesNum int
	err := vrf.db.QueryRow(q).Scan(&namesNum)
	return namesNum, err
}

func (vfr verifio) dedupeNames() error {
	log.Info().Msg("Extracting unique name-strings. It will take a while.")
	q := `INSERT INTO unique_names
          SELECT name_string, AVG(odds_log10), count(*)
            FROM detected_names GROUP BY name_string
            ORDER BY name_string`

	stmt, err := vfr.db.Prepare(q)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}
	vfr.Status = AllNamesHarvested
	return nil
}
