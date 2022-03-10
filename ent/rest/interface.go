package rest

import echo "github.com/labstack/echo/v4"

// REST interface describes functionality of RESTful API for BHL's
// scientific names index.
type REST interface {
	// Run creates a service to BHLindex running on.
	Run(port int)

	// Ping checks connection to the RESTful server.
	Ping() func(echo.Context) error

	// Version returns bhlindex's version.
	Version() func(echo.Context) error

	// Items returns BHL's items metadata. An item can be a book, a journal,
	// a bulletin etc.
	Items() func(echo.Context) error

	// Pages returns metadata for pages from BHL's items.
	Pages() func(echo.Context) error

	// Names returns a batch of verified names from the index.
	Names() func(echo.Context) error

	// Occurrences returns a batch of names occurrences and their metadata.
	Occurrences() func(echo.Context) error
}
