package verifio

import (
	"context"
	"sync"

	"github.com/gnames/bhlindex/ent/name"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnverifier"
	gnvconfig "github.com/gnames/gnverifier/config"
	"github.com/gnames/gnverifier/io/verifrest"
)

func (vrf verifio) sendToVerify(
	ctx context.Context,
	chNames <-chan []name.UniqueName,
	chVer chan<- []name.VerifiedName,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	opts := []gnvconfig.Option{
		gnvconfig.OptVerifierURL(vrf.cfg.VerifierURL),
	}
	gnvcfg := gnvconfig.New(opts...)
	gnvRest := verifrest.New(vrf.cfg.VerifierURL)
	gnv := gnverifier.New(gnvcfg, gnvRest)
	for uns := range chNames {
		// TODO: parallelize this
		names := make([]string, len(uns))
		for i := range uns {
			names[i] = uns[i].Name
		}
		out := gnv.VerifyBatch(names)
		vns := prepareNames(out, uns)
		chVer <- vns
	}
}

func prepareNames(vns []vlib.Name, uns []name.UniqueName) []name.VerifiedName {
	res := make([]name.VerifiedName, len(vns))
	for i := range vns {
		n := name.VerifiedName{
			Name:              vns[i].Name,
			MatchType:         vns[i].MatchType.String(),
			Curation:          vns[i].Curation.String(),
			DataSourcesNumber: vns[i].DataSourcesNum,
			OddsLog10:         uns[i].OddsLog10,
			Occurrences:       uns[i].Occurrences,
			Error:             vns[i].Error,
		}
		if br := vns[i].BestResult; br != nil {
			n.RecordID = br.RecordID
			n.EditDistance = br.EditDistance
			n.StemEditDistance = br.StemEditDistance
			n.MatchedName = br.MatchedName
			n.MatchedCanonical = br.MatchedCanonicalFull
			n.CurrentName = br.CurrentName
			n.CurrentCanonical = br.CurrentCanonicalFull
			n.Classification = br.ClassificationPath
			n.DataSourceID = br.DataSourceID
			n.DataSourceTitle = br.DataSourceTitleShort
		}
		res[i] = n
	}
	return res
}
