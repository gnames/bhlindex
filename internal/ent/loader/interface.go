package loader

import (
	"context"

	"github.com/gnames/bhlindex/internal/ent/item"
)

type Loader interface {
	// DetectPageDups checks directory structure for errors and finds the number of
	// pages in the corpus.
	DetectPageDups() (Loader, error)

	// LoadItems walks BHL corpus directory, creates `items` in the database and
	// pushes created `*item.Item` to channel for furhter use.
	LoadItems(context.Context, chan<- *item.Item) error
}
