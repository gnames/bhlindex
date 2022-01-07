package bhlindex_test

import (
	"fmt"
	"testing"

	"github.com/gnames/bhlindex"
	"github.com/gnames/bhlindex/config"
	"github.com/gnames/bhlindex/ent/finder"
	"github.com/gnames/bhlindex/io/dbio"
	"github.com/gnames/bhlindex/io/loaderio"
	"github.com/stretchr/testify/assert"
)

func TestFindNames(t *testing.T) {
	path := "./testdata/bhl/"
	opt := config.OptBHLdir(path)
	cfg := config.New(opt)
	db := dbio.NewWithInit(cfg)
	ldr := loaderio.New(cfg, db)
	fdr := finder.New(cfg, db)

	bhli := bhlindex.New(cfg)
	err := bhli.FindNames(ldr, fdr)
	assert.Nil(t, err)

}

func BenchmarkInit(b *testing.B) {
	path := "./testdata/bhl/"
	opt := config.OptBHLdir(path)
	cfg := config.New(opt)
	db := dbio.NewWithInit(cfg)
	ldr := loaderio.New(cfg, db)
	fdr := finder.New(cfg, db)
	bhli := bhlindex.New(cfg)

	b.Run("Load Items", func(b *testing.B) {
		var err error
		for i := 0; i < b.N; i++ {
			err = bhli.FindNames(ldr, fdr)
			assert.Nil(b, err)
		}
		_ = fmt.Sprintf("%v", err)
	})

}
