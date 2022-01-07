package bhlindex

import "github.com/gnames/gnlib/ent/gnvers"

type bhlindex struct{}

func New() BHLindex {
	return &bhlindex{}
}

func (bi *bhlindex) FindNames() error {
	return nil
}

func (bi *bhlindex) VerifyNames() error {
	return nil
}

func (bi *bhlindex) GetVersion() gnvers.Version {

}
