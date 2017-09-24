package finder

import (
	"database/sql"
	"path/filepath"
	"sync"
	"time"

	"github.com/GlobalNamesArchitecture/bhlindex/loader"
	"github.com/GlobalNamesArchitecture/bhlindex/models"
	"github.com/GlobalNamesArchitecture/bhlindex/util"
	"github.com/lib/pq"
)

func FindNames(db *sql.DB) string {
	c := make(chan string)
	var wg sync.WaitGroup

	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go findWorker(&wg, c, db)
	}

	loader.Path(c)
	wg.Wait()
	return "ok"
}

func findWorker(wg *sync.WaitGroup, c <-chan string, db *sql.DB) {
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
	util.Check(err)

	stmt, err := txn.Prepare(pq.CopyIn("titles", columns...))
	util.Check(err)

	for _, t := range batch {
		_, err = stmt.Exec(t.Path, t.InternetArchiveID, now)
		util.Check(err)
	}

	_, err = stmt.Exec()
	util.Check(err)

	err = stmt.Close()
	util.Check(err)

	err = txn.Commit()
	util.Check(err)
}
