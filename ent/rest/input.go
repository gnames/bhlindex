package rest

import (
	"github.com/gnames/bhlindex/ent/item"
	"github.com/gnames/bhlindex/ent/name"
	"github.com/gnames/bhlindex/ent/page"
)

type Input struct {
	Offset      int   `json:"offset" query:"offset"`
	Limit       int   `json:"limit" query:"limit"`
	DataSources []int `json:"dataSources" query:"data_sources"`
}

type OutputItems struct {
	Items []item.Item `json:"items"`
}

type OutputPages struct {
	Pages []page.Page `json:"pages"`
}

type OutputOccurrences struct {
	Occurrences []name.DetectedName `json:"occurrences"`
}

type OutputNames struct {
	Names []name.VerifiedName `json:"names"`
}
