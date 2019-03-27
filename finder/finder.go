package finder

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/gnames/bhlindex"
	"github.com/gnames/bhlindex/loader"
	"github.com/gnames/bhlindex/models"
	"github.com/gnames/gnfinder/dict"
	"github.com/lib/pq"
)

func ProcessTitles(db *sql.DB, d *dict.Dictionary, w int) {
	titleIDs := make(chan int)
	findQueue := make(chan *models.Title)
	counter := make(chan int)
	results := make(chan []models.DetectedName)
	workersNum := w
	var wgLoad sync.WaitGroup
	var wgFind sync.WaitGroup
	var wgSave sync.WaitGroup

	go counterLog(counter)

	for i := 1; i <= 10; i++ {
		wgLoad.Add(1)
		go titleWorker(db, &wgLoad, titleIDs, findQueue, counter)
	}

	for i := 1; i <= workersNum; i++ {
		wgFind.Add(1)
		go finderWorker(&wgFind, findQueue, results, d)
	}

	wgSave.Add(1)
	go saveFoundNames(db, &wgSave, results)

	loader.ImportTitles(db, titleIDs)

	wgLoad.Wait()
	close(counter)
	close(findQueue)
	wgFind.Wait()
	close(results)
	wgSave.Wait()
}

func saveFoundNames(db *sql.DB, wgSave *sync.WaitGroup,
	results <-chan []models.DetectedName) {
	defer wgSave.Done()
	for v := range results {
		savePageNameStrings(db, v)
	}
}

func savePageNameStrings(db *sql.DB, names []models.DetectedName) {
	now := time.Now()
	columns := []string{"page_id", "title_id", "name_string", "name_offset_start",
		"name_offset_end", "ends_next_page", "odds", "kind", "updated_at"}
	transaction, err := db.Begin()
	bhlindex.Check(err)
	defer transaction.Rollback()

	stmt, err := transaction.Prepare(pq.CopyIn("page_name_strings", columns...))
	bhlindex.Check(err)

	for _, v := range names {
		_, err = stmt.Exec(v.PageID, v.TitleID, v.NameString, v.OffsetStart, v.OffsetEnd,
			v.EndsNextPage, v.Odds, v.Kind, now)
		bhlindex.Check(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Println(`
Bulk import of titles data failed, probably you need to empty all data
and start with empty database.`)
		log.Fatal(err)
	}

	err = stmt.Close()
	bhlindex.Check(err)

	err = transaction.Commit()
	bhlindex.Check(err)
}

func finderWorker(wg *sync.WaitGroup, findQueue <-chan *models.Title,
	results chan<- []models.DetectedName, d *dict.Dictionary) {
	defer wg.Done()
	for t := range findQueue {
		results <- t.FindNames(d)
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
