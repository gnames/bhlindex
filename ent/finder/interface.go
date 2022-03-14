package finder

import (
	"context"
	"sync"

	"github.com/gnames/bhlindex/ent/item"
	"github.com/gnames/bhlindex/ent/name"
)

// Finder detects scientific names in items.
type Finder interface {
	// FindNames detects scientific names in an item. When a name is detected, it
	// is associated with the parent Item, a Page where it is located, as well as
	// its position on the page which is represented by offset in unicode
	// characters from the start of the page. Optionally, it provides words
	// that are located immediately before and after a name. If there is an
	// nomenclatural annotation after the name, such annotation is normalized
	// and registered.
	FindNames(
		<-chan *item.Item,
		chan<- []name.DetectedName,
		*sync.WaitGroup,
	)

	// SaveNames is an aggregator method that should always run in one instance.
	// It collects detected names data and saves them into storage.
	SaveNames(context.Context, <-chan []name.DetectedName) error
}
