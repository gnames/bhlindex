package finder

import (
	"database/sql"
	"log"
	"sync"

	"github.com/GlobalNamesArchitecture/bhlindex/loader"
	"github.com/GlobalNamesArchitecture/bhlindex/models"
)

func ProcessTitles(db *sql.DB) {
	titleIDs := make(chan int)
	counter := make(chan int)
	var wg sync.WaitGroup

	go counterLog(counter)

	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go titleWorker(db, &wg, titleIDs, counter)
	}

	loader.ImportTitles(db, titleIDs)
	wg.Wait()
	close(counter)
}

func counterLog(counter <-chan int) {
	count := 0
	for i := range counter {
		count += i
		if count%10000 == 0 {
			log.Printf("Processed %d titles", count)
		}
	}
}

func titleWorker(db *sql.DB, wg *sync.WaitGroup, titleIDs <-chan int,
	counter chan<- int) {
	defer wg.Done()
	for titleID := range titleIDs {
		title := models.TitleFind(db, titleID)
		loader.ImportPages(db, &title)
		counter <- 1
	}
}
