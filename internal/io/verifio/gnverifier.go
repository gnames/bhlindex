package verifio

import (
	"context"
	"log"
	"sync"

	"github.com/gnames/bhlindex/internal/ent/name"
	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnverifier/pkg/ent/verifier"
)

type verifiedBatch struct {
	names []name.VerifiedName
	// in case if we get all results instead of only the best, we need to
	// know how many unique names were processed.
	namesNum int
}

func (vrf verifio) sendToVerify(
	ctx context.Context,
	vrfREST verifier.Verifier,
	chIn <-chan []name.UniqueName,
	chOut chan<- verifiedBatch,
	wgExt *sync.WaitGroup,
) {
	defer wgExt.Done()

	var wg sync.WaitGroup
	jobs := 4
	wg.Add(jobs)

	for i := 0; i < jobs; i++ {
		go vrf.verifyWorker(ctx, vrfREST, chIn, chOut, &wg)
	}

	wg.Wait()
}

func (vrf verifio) verifyWorker(
	ctx context.Context,
	vrfREST verifier.Verifier,
	chIn <-chan []name.UniqueName,
	chOut chan<- verifiedBatch,
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
		input := vlib.Input{NameStrings: names, WithAllMatches: vrf.cfg.VerifAllResults}
		output := vrfREST.Verify(ctx, input)
		if len(output.Names) < 1 {
			log.Fatalf("Did not get results from verifier")
		}
		outNames, namesNum := prepareNames(output.Names, uns)

		select {
		case <-ctx.Done():
			return
		case chOut <- verifiedBatch{names: outNames, namesNum: namesNum}:
		}
	}
}

func prepareNames(vns []vlib.Name, uns []name.UniqueName) ([]name.VerifiedName, int) {
	res := make([]name.VerifiedName, 0, len(vns))
	var count int
	for i := range vns {
		count++
		n := name.VerifiedName{
			ID:                uns[i].ID,
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
			res = append(res, n)
		} else if rs := vns[i].Results; len(rs) > 0 {
			var sortNum int
			for _, v := range rs {
				n.RecordID = v.RecordID
				n.EditDistance = v.EditDistance
				n.StemEditDistance = v.StemEditDistance
				n.MatchedName = v.MatchedName
				n.MatchedCanonical = v.MatchedCanonicalFull
				n.MatchedCardinality = v.MatchedCardinality
				n.CurrentName = v.CurrentName
				n.CurrentCanonical = v.CurrentCanonicalFull
				n.CurrentCardinality = v.CurrentCardinality
				n.Classification = v.ClassificationPath
				n.ClassificationRanks = v.ClassificationRanks
				n.ClassificationIDs = v.ClassificationIDs
				n.DataSourceID = v.DataSourceID
				n.DataSourceTitle = v.DataSourceTitleShort
				n.SortOrder = sortNum
				nCopy := n
				res = append(res, nCopy)
				sortNum++
			}
		} else {
			res = append(res, n)
		}
	}
	return res, count
}
