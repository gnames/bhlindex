package loader

import (
	"context"

	"github.com/gnames/bhlindex/internal/ent/item"
)

type Loader interface {
	// LoadItems walks BHL corpus directory, creates `items` in the database and
	// pushes created `*item.Item` to channel for furhter use.
	LoadItems(context.Context, chan<- *item.Item) error
}
