package finder

import (
	"sync"

	"github.com/GlobalNamesArchitecture/bhlindex/loader"
)

func FindNames() string {
	c := make(chan string)
	var wg sync.WaitGroup

	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go findWorker(&wg, c)
	}

	loader.Path(c)
	wg.Wait()
	return "ok"
}

func findWorker(wg *sync.WaitGroup, c <-chan string) {
	defer wg.Done()
	for p := range c {
		p = p + ""
	}
}
