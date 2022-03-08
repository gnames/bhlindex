package verifio

import (
	"context"
	"sync"

	vlib "github.com/gnames/gnlib/ent/verifier"
)

func (vrf verifio) sendToVerify(
	ctx context.Context,
	chNames <-chan []string,
	chVer chan<- vlib.Output,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for range chNames {
		chVer <- vlib.Output{}
	}
}
