package restio

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gnames/bhlindex"
	"github.com/gnames/bhlindex/ent/item"
	"github.com/gnames/bhlindex/ent/name"
	"github.com/gnames/bhlindex/ent/page"
	"github.com/gnames/bhlindex/ent/rest"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

var (
	maxLimit = 50_000
	apiPath  = "/api/v0/"
)

type restio struct {
	bi bhlindex.BHLindex
	db *sql.DB
}

func New(bi bhlindex.BHLindex, db *sql.DB) rest.REST {
	res := restio{bi: bi, db: db}
	return res
}

// Run starts HTTP/1 service on a given port for scientific names verification.
func (r restio) Run(port int) {
	log.Info().Int("port", port).Msg("Starting HTTP API server")
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	r.setLogger(e)

	e.GET(apiPath+"ping", r.Ping())
	e.GET(apiPath+"version", r.Version())
	e.GET(apiPath+"items", r.Items())
	e.GET(apiPath+"pages", r.Pages())
	e.GET(apiPath+"occurrences", r.Occurrences())
	e.GET(apiPath+"names", r.Names())
	e.GET(apiPath+"names/last_id", r.NamesLastID())

	addr := fmt.Sprintf(":%d", port)
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}

func (r restio) Ping() func(echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	}
}

func (r restio) Version() func(echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, r.bi.GetVersion())
	}
}

func (r restio) Items() func(echo.Context) error {
	return func(c echo.Context) error {
		var err error
		var items []item.Item
		var input rest.Input

		ctx, cancel := getContext(c)
		defer cancel()

		if err = c.Bind(&input); err != nil {
			log.Warn().Err(err)
			return err
		}
		if input.Limit > maxLimit {
			input.Limit = maxLimit
		}
		select {
		case <-ctx.Done():
			log.Warn().Err(ctx.Err()).Msg("Forced cancellation")
			return ctx.Err()
		default:
			if items, err = r.items(ctx, input); err != nil {
				log.Warn().Err(err)
				return err
			}
		}

		return c.JSON(http.StatusOK, items)
	}
}

func (r restio) Pages() func(echo.Context) error {
	return func(c echo.Context) error {
		var err error
		var pages []page.Page
		var input rest.Input

		ctx, cancel := getContext(c)
		defer cancel()

		if err = c.Bind(&input); err != nil {
			log.Warn().Err(err)
			return err
		}
		if input.Limit > maxLimit {
			input.Limit = maxLimit
		}
		select {
		case <-ctx.Done():
			log.Warn().Err(ctx.Err()).Msg("Forced cancellation")
			return ctx.Err()
		default:
			if pages, err = r.pages(ctx, input); err != nil {
				log.Warn().Err(err)
				return err
			}
		}

		return c.JSON(http.StatusOK, pages)
	}
}

func (r restio) Occurrences() func(echo.Context) error {
	return func(c echo.Context) error {
		var err error
		var occurrences []name.DetectedName
		var input rest.Input

		ctx, cancel := getContext(c)
		defer cancel()

		if err = c.Bind(&input); err != nil {
			log.Warn().Err(err)
			return err
		}
		if input.Limit > maxLimit {
			input.Limit = maxLimit
		}
		select {
		case <-ctx.Done():
			log.Warn().Err(ctx.Err()).Msg("Forced cancellation")
			return ctx.Err()
		default:
			if occurrences, err = r.occurrences(ctx, input); err != nil {
				log.Warn().Err(err)
				return err
			}
		}
		return c.JSON(http.StatusOK, occurrences)
	}
}

func (r restio) Names() func(echo.Context) error {
	return func(c echo.Context) error {
		var err error
		var names []name.VerifiedName
		var input rest.Input

		ctx, cancel := getContext(c)
		defer cancel()

		if err = c.Bind(&input); err != nil {
			log.Warn().Err(err)
			return err
		}
		if input.Limit > maxLimit {
			input.Limit = maxLimit
		}
		select {
		case <-ctx.Done():
			log.Warn().Err(ctx.Err()).Msg("Forced cancellation")
			return ctx.Err()
		default:
			if names, err = r.names(ctx, input); err != nil {
				log.Warn().Err(err)
				return err
			}
		}
		return c.JSON(http.StatusOK, names)
	}
}

func (r restio) NamesLastID() func(echo.Context) error {
	return func(c echo.Context) error {
		var lastID int
		var err error
		if lastID, err = r.namesLastID(); err != nil {
			return err
		}
		return c.String(http.StatusOK, strconv.Itoa(lastID))
	}
}

func getContext(c echo.Context) (ctx context.Context, cancel func()) {
	ctx = c.Request().Context()
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	return ctx, cancel
}

func (r restio) setLogger(e *echo.Echo) {
	// log.Logger = log.Output(os.Stdout)
	if r.bi.GetConfig().WithWebLogs {
		e.Use(middleware.Logger())
	}
}