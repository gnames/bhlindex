package bhlindex

import (
	"github.com/gnames/bhlindex/ent/finder"
	"github.com/gnames/bhlindex/ent/loader"
	"github.com/gnames/gnlib/ent/gnvers"
)

type BHLindex interface {
	FindNames(loader.Loader, finder.Finder) error
	VerifyNames() error
	GetVersion() gnvers.Version
}
