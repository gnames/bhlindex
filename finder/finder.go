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

func ProcessItems(db *sql.DB, d *dict.Dictionary, w int) {
	itemIDs := make(chan int)
	findQueue := make(chan *models.Item)
	counter := make(chan int)
	results := make(chan []models.DetectedName)
	workersNum := w
	var wgLoad sync.WaitGroup
	var wgFind sync.WaitGroup
	var wgSave sync.WaitGroup

	go counterLog(counter)

	for i := 1; i <= 10; i++ {
		wgLoad.Add(1)
		go itemWorker(db, &wgLoad, itemIDs, findQueue, counter)
	}

	for i := 1; i <= workersNum; i++ {
		wgFind.Add(1)
		go finderWorker(&wgFind, findQueue, results, d)
	}

	wgSave.Add(1)
	go saveFoundNames(db, &wgSave, results)

	loader.ImportItems(db, itemIDs)

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
	columns := []string{"page_id", "item_id", "words_before",
		"name_string", "words_after", "annot_nomen", "name_offset_start",
		"name_offset_end", "ends_next_page", "odds", "kind", "updated_at"}
	transaction, err := db.Begin()
	bhlindex.Check(err)
	defer transaction.Rollback()

	stmt, err := transaction.Prepare(pq.CopyIn("page_name_strings", columns...))
	bhlindex.Check(err)

	for _, v := range names {
		wordsBefore := v.WordsBefore
		if len([]rune(wordsBefore)) > 250 {
			wordsBefore = string([]rune(wordsBefore)[:250])
		}
		wordsAfter := v.WordsAfter
		if len([]rune(wordsAfter)) > 250 {
			wordsAfter = string([]rune(wordsAfter)[:250])
		}

		kind := "N/A"
		switch v.Cardinality {
		case 1:
			kind = "Uninomial"
		case 2:
			kind = "Binomial"
		case 3:
			kind = "Trinomial"
		}

		_, err = stmt.Exec(v.PageID, v.ItemID, wordsBefore, v.NameString,
			wordsAfter, v.AnnotNomen, v.OffsetStart, v.OffsetEnd,
			v.EndsNextPage, v.Odds, kind, now)
		bhlindex.Check(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Println(`
Bulk import of items data failed, probably you need to empty all data
and start with empty database.`)
		log.Fatal(err)
	}

	err = stmt.Close()
	bhlindex.Check(err)

	err = transaction.Commit()
	bhlindex.Check(err)
}

func finderWorker(wg *sync.WaitGroup, findQueue <-chan *models.Item,
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
			log.Printf("Processed %d items", count)
		}
	}
}

func itemWorker(db *sql.DB, wg *sync.WaitGroup, itemIDs <-chan int,
	resQueue chan<- *models.Item, counter chan<- int) {
	defer wg.Done()
	for itemID := range itemIDs {
		item := models.ItemFind(db, itemID)
		loader.ImportPages(db, &item)
		resQueue <- &item
		counter <- 1
	}
}
