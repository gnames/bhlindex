package bhlindex_test

import (
	"fmt"
	"testing"

	bhlindex "github.com/gnames/bhlindex/internal"
	"github.com/gnames/bhlindex/internal/config"
	"github.com/gnames/bhlindex/internal/io/dbio"
	"github.com/gnames/bhlindex/internal/io/finderio"
	"github.com/gnames/bhlindex/internal/io/loaderio"
	"github.com/gnames/bhlindex/internal/io/verifio"
	"github.com/stretchr/testify/assert"
)

func getBHLindex() (config.Config, bhlindex.BHLindex) {
	opts := []config.Option{
		config.OptBHLdir("./testdata/bhl"),
		config.OptPgDatabase("bhlindex_test"),
	}

	cfg := config.New(opts...)
	bi := bhlindex.New(cfg)
	return cfg, bi
}

func TestVerifyNames(t *testing.T) {
	cfg, bi := getBHLindex()
	db := dbio.NewWithInit(cfg)
	ldr := loaderio.New(cfg, db)
	fdr := finderio.New(cfg, db)
	vdr := verifio.New(cfg, db)

	err := bi.FindNames(ldr, fdr)
	assert.Nil(t, err)
	err = bi.VerifyNames(vdr)
	assert.Nil(t, err)
}

func BenchmarkInit(b *testing.B) {
	path := "./testdata/bhl/"
	opt := config.OptBHLdir(path)
	cfg := config.New(opt)
	db := dbio.NewWithInit(cfg)
	ldr := loaderio.New(cfg, db)
	fdr := finderio.New(cfg, db)
	vdr := verifio.New(cfg, db)
	bhli := bhlindex.New(cfg)

	b.Run("FindNames", func(b *testing.B) {
		var err error
		for i := 0; i < b.N; i++ {
			err = bhli.FindNames(ldr, fdr)
			assert.Nil(b, err)
		}
		_ = fmt.Sprintf("%v", err)
	})

	b.Run("VerifyNames", func(b *testing.B) {
		var err error
		for i := 0; i < b.N; i++ {
			err = bhli.FindNames(ldr, fdr)
			assert.Nil(b, err)
			err = bhli.VerifyNames(vdr)
			assert.Nil(b, err)
		}
		_ = fmt.Sprintf("%v", err)
	})
}
