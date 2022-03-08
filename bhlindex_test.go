package bhlindex_test

import (
	"fmt"
	"testing"

	"github.com/gnames/bhlindex"
	"github.com/gnames/bhlindex/config"
	"github.com/gnames/bhlindex/io/dbio"
	"github.com/gnames/bhlindex/io/finderio"
	"github.com/gnames/bhlindex/io/loaderio"
	"github.com/gnames/bhlindex/io/verifio"
	"github.com/stretchr/testify/assert"
)

func TestFindNames(t *testing.T) {
	path := "./testdata/bhl/"
	opt := config.OptBHLdir(path)
	cfg := config.New(opt)
	db := dbio.NewWithInit(cfg)
	ldr := loaderio.New(cfg, db)
	fdr := finderio.New(cfg, db)

	bhli := bhlindex.New(cfg)
	err := bhli.FindNames(ldr, fdr)
	assert.Nil(t, err)
}

func TestVerifyNames(t *testing.T) {
	path := "./testdata/bhl/"
	opt := config.OptBHLdir(path)
	cfg := config.New(opt)
	db := dbio.NewWithInit(cfg)
	ldr := loaderio.New(cfg, db)
	fdr := finderio.New(cfg, db)
	vdr := verifio.New(cfg, db)

	bhli := bhlindex.New(cfg)
	err := bhli.FindNames(ldr, fdr)
	assert.Nil(t, err)
	err = bhli.VerifyNames(vdr)
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
