package verifio

import (
	"context"
	"log"
	"sync"

	"github.com/gnames/bhlindex/ent/name"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnverifier/ent/verifier"
)

func (vrf verifio) sendToVerify(
	ctx context.Context,
	vrfREST verifier.Verifier,
	chIn <-chan []name.UniqueName,
	chOut chan<- []name.VerifiedName,
	wgExt *sync.WaitGroup,
) {
	defer wgExt.Done()

	var wg sync.WaitGroup
	jobs := 4
	wg.Add(jobs)

	for i := 0; i < jobs; i++ {
		go verifyWorker(ctx, vrfREST, chIn, chOut, &wg)
	}

	wg.Wait()
}

func verifyWorker(
	ctx context.Context,
	vrfREST verifier.Verifier,
	chIn <-chan []name.UniqueName,
	chOut chan<- []name.VerifiedName,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for uns := range chIn {
		if len(uns) == 0 {
			continue
		}
		names := make([]string, len(uns))
		for i := range uns {
			names[i] = uns[i].Name
		}
		input := vlib.Input{NameStrings: names}
		output := vrfREST.Verify(ctx, input)
		if len(output.Names) < 1 {
			log.Fatalf("Did not get results from verifier")
		}
		outNames := prepareNames(output.Names, uns)

		select {
		case <-ctx.Done():
			return
		case chOut <- outNames:
		}
	}
}

func prepareNames(vns []vlib.Name, uns []name.UniqueName) []name.VerifiedName {
	res := make([]name.VerifiedName, len(vns))
	for i := range vns {
		n := name.VerifiedName{
			NameID:            uns[i].ID,
			Name:              vns[i].Name,
			Cardinality:       vns[i].Cardinality,
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
			n.MatchedCardinality = br.MatchedCardinality
			n.CurrentName = br.CurrentName
			n.CurrentCanonical = br.CurrentCanonicalFull
			n.CurrentCardinality = br.CurrentCardinality
			n.Classification = br.ClassificationPath
			n.ClassificationRanks = br.ClassificationRanks
			n.ClassificationIDs = br.ClassificationIDs
			n.DataSourceID = br.DataSourceID
			n.DataSourceTitle = br.DataSourceTitleShort
		}
		res[i] = n
	}
	return res
}
