package restio_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gnames/bhlindex"
	"github.com/gnames/bhlindex/config"
	"github.com/gnames/bhlindex/ent/rest"
	"github.com/gnames/bhlindex/io/dbio"
	"github.com/gnames/bhlindex/io/finderio"
	"github.com/gnames/bhlindex/io/loaderio"
	"github.com/gnames/bhlindex/io/restio"
	"github.com/gnames/bhlindex/io/verifio"
	"github.com/gnames/gnfmt"
	echo "github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var r rest.REST = getREST()

func TestPing(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := r.Ping()(c)
	assert.Nil(t, err)
	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, rec.Body.String(), "pong")
}

func TestVersion(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := r.Version()(c)
	assert.Nil(t, err)
	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Contains(t, rec.Body.String(), "version")
}

func TestItems(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/items?offset=0&limit=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := r.Items()(c)
	assert.Nil(t, err)
	assert.Equal(t, rec.Code, http.StatusOK)
	var items rest.OutputItems
	err = gnfmt.GNjson{}.Decode(rec.Body.Bytes(), &items)
	assert.Nil(t, err)
	assert.Equal(t, len(items.Items), 10)
}

func TestPages(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/pages?offset=0&limit=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := r.Pages()(c)
	assert.Nil(t, err)
	assert.Equal(t, rec.Code, http.StatusOK)
	var pages rest.OutputPages
	err = gnfmt.GNjson{}.Decode(rec.Body.Bytes(), &pages)
	assert.Nil(t, err)
	assert.Equal(t, len(pages.Pages), 10)
}

func TestOccurrences(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/occurrences?offset=0&limit=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := r.Occurrences()(c)
	assert.Nil(t, err)
	assert.Equal(t, rec.Code, http.StatusOK)
	var ocrs rest.OutputOccurrences
	err = gnfmt.GNjson{}.Decode(rec.Body.Bytes(), &ocrs)
	assert.Nil(t, err)
	assert.Equal(t, len(ocrs.Occurrences), 10)
}

func TestNames(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/names?offset=0&limit=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := r.Names()(c)
	assert.Nil(t, err)
	assert.Equal(t, rec.Code, http.StatusOK)
	var names rest.OutputNames
	err = gnfmt.GNjson{}.Decode(rec.Body.Bytes(), &names)
	assert.Nil(t, err)
	assert.Equal(t, len(names.Names), 10)
}

func dataExists(r rest.REST) bool {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/names?offset=0&limit=1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := r.Names()(c)
	if err != nil {
		log.Warn().Err(err).Msg("Cannot query names")
		return false
	}
	out := rec.Body.String()
	return strings.Contains(out, `"id":`)
}

func getREST() rest.REST {
	opts := []config.Option{
		config.OptBHLdir("../../testdata/bhl"),
		config.OptPgDatabase("bhlindex_test"),
	}

	cfg := config.New(opts...)
	bi := bhlindex.New(cfg)
	db := dbio.New(cfg)
	r := restio.New(bi, db)
	if !dataExists(r) {
		log.Info().Msg("No data found, creating...")
		_ = dbio.NewWithInit(cfg)
		ldr := loaderio.New(cfg, db)
		fdr := finderio.New(cfg, db)
		err := bi.FindNames(ldr, fdr)
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot find names")
		}
		vrf := verifio.New(cfg, db)
		err = bi.VerifyNames(vrf)
		if err != nil {
			log.Fatal().Err(err).Msg("Name verification failed")
		}
	}
	return r
}
