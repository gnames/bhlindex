package verifio

import (
	"context"
	"sync"

	vlib "github.com/gnames/gnlib/ent/verifier"
	"github.com/gnames/gnverifier"
	gnvconfig "github.com/gnames/gnverifier/config"
	"github.com/gnames/gnverifier/io/verifrest"
)

func (vrf verifio) sendToVerify(
	ctx context.Context,
	chNames <-chan []string,
	chVer chan<- []vlib.Name,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	opts := []gnvconfig.Option{
		gnvconfig.OptVerifierURL(vrf.cfg.VerifierURL),
	}
	gnvcfg := gnvconfig.New(opts...)
	gnvRest := verifrest.New(vrf.cfg.VerifierURL)
	gnv := gnverifier.New(gnvcfg, gnvRest)
	for names := range chNames {
		// TODO: parallelize this
		out := gnv.VerifyBatch(names)
		chVer <- out
	}
}
