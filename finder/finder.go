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
	findQueue := make(chan *models.Title)
	counter := make(chan int)
	results := make(chan *models.Title)
	var wgLoad sync.WaitGroup
	var wgFind sync.WaitGroup
	var wgFinish sync.WaitGroup

	go counterLog(counter)

	for i := 1; i <= 10; i++ {
		wgLoad.Add(1)
		go titleWorker(db, &wgLoad, titleIDs, findQueue, counter)
	}

	for i := 1; i <= 50; i++ {
		wgFind.Add(1)
		go finderWorker(&wgFind, findQueue, results)
	}

	wgFinish.Add(1)
	go saveFoundNames(db, &wgFinish, results)

	loader.ImportTitles(db, titleIDs)

	wgLoad.Wait()
	close(counter)
	close(findQueue)
	wgFind.Wait()
	close(results)
	wgFinish.Wait()
}

func saveFoundNames(db *sql.DB, wg *sync.WaitGroup,
	results <-chan *models.Title) {
	defer wg.Done()
	for _ = range results {
	}
}

func finderWorker(wg *sync.WaitGroup, findQueue <-chan *models.Title,
	results chan<- *models.Title) {
	defer wg.Done()
	for t := range findQueue {
		results <- t
	}
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
	resQueue chan<- *models.Title, counter chan<- int) {
	defer wg.Done()
	for titleID := range titleIDs {
		title := models.TitleFind(db, titleID)
		loader.ImportPages(db, &title)
		resQueue <- &title
		counter <- 1
	}
}
