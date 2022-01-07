package bhlindex

import "github.com/gnames/gnlib/ent/gnvers"

type BHLindex interface {
	FindNames() error
	VerifyNames() error
	GetVersion() gnvers.Version
}
