package finder

import (
	"database/sql"
	"sync"

	"github.com/GlobalNamesArchitecture/bhlindex/loader"
	"github.com/GlobalNamesArchitecture/bhlindex/models"
)

func ProcessTitles(db *sql.DB) {
	titleIDs := make(chan int)
	var wg sync.WaitGroup

	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go titleWorker(&wg, titleIDs, db)
	}

	loader.ImportTitles(db, titleIDs)
	wg.Wait()
}

func titleWorker(wg *sync.WaitGroup, c <-chan int, db *sql.DB) {
	defer wg.Done()
	for titleID := range c {
		title := models.TitleFind(db, titleID)
		loader.ImportPages(db, &title)
	}
}
