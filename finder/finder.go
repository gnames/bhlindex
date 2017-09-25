package finder

import (
	"database/sql"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/GlobalNamesArchitecture/bhlindex"
	"github.com/GlobalNamesArchitecture/bhlindex/loader"
	"github.com/GlobalNamesArchitecture/bhlindex/models"
	"github.com/lib/pq"
)

func FindNames(db *sql.DB, nf bhlindex.NameFinder) {
	c := make(chan string)
	var wg sync.WaitGroup

	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go findWorker(&wg, c, db, nf)
	}

	loader.Path(c)
	wg.Wait()
}

func findWorker(wg *sync.WaitGroup, c <-chan string, db *sql.DB,
	nf bhlindex.NameFinder) {
	defer wg.Done()
	batchSize := 1000
	batch := make([]models.Title, batchSize)
	i := 0
	for p := range c {
		if i < batchSize {
			batch[i] = titleFromPath(p)
		} else {
			i = 0
			batch[i] = titleFromPath(p)
			titleBatchSave(db, batch)
		}
		i++
	}
	i--
	titleBatchSave(db, batch[0:i])
}

func titleFromPath(p string) models.Title {
	return models.Title{Path: p, InternetArchiveID: filepath.Base(p)}
}

func titleBatchSave(db *sql.DB, batch []models.Title) {
	now := time.Now()
	columns := []string{"path", "internet_archive_id", "updated_at"}
	txn, err := db.Begin()
	bhlindex.Check(err)

	stmt, err := txn.Prepare(pq.CopyIn("titles", columns...))
	bhlindex.Check(err)

	for _, t := range batch {
		_, err = stmt.Exec(t.Path, t.InternetArchiveID, now)
		bhlindex.Check(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Println(`
Bulk import of titles data failed, probably you need to empty all data
and start with empty database.
`)
		log.Fatal(err)
	}

	err = stmt.Close()
	bhlindex.Check(err)

	err = txn.Commit()
	bhlindex.Check(err)
}
