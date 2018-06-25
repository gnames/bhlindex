package finder

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gnames/bhlindex"
	"github.com/gnames/bhlindex/models"
	"github.com/gnames/gnfinder/resolver"
	gfutil "github.com/gnames/gnfinder/util"
	"github.com/lib/pq"
)

var (
	empty     struct{}
	uniqNames = make(map[string]struct{})
)

func Verify(db *sql.DB, foundNames <-chan []models.DetectedName,
	wg *sync.WaitGroup) {
	m := gfutil.NewModel()
	m.Workers = 15
	resChan := make(chan []string)
	defer wg.Done()

	go prepareVerification(resChan, foundNames)

	count := 0
	for names := range resChan {
		count += len(names)
		fmt.Println("Matching", count, "names", time.Now())
		verified := resolver.Verify(names, m)
		saveVerifiedNameStrings(db, verified)
	}
}

func saveVerifiedNameStrings(db *sql.DB, verified resolver.VerifyOutput) {
	var errStr sql.NullString
	now := time.Now()
	columns := []string{"name", "match_type", "edit_distance",
		"matched_name", "current_name", "classification", "datasource_id",
		"datasources_number", "retries", "error", "updated_at"}
	transaction, err := db.Begin()
	bhlindex.Check(err)

	stmt, err := transaction.Prepare(pq.CopyIn("name_strings", columns...))
	bhlindex.Check(err)

	for name, v := range verified {
		if v.Error == nil {
			errStr = sql.NullString{}
		} else {
			errStr.Scan(v.Error.Error())
		}
		_, err = stmt.Exec(name, v.MatchType, v.EditDistance, v.MatchedName,
			v.CurrentName, v.ClassificationPath, v.DataSourceID, v.DatabasesNum,
			v.Retries, errStr, now)
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

func prepareVerification(resChan chan<- []string, foundNames <-chan []models.DetectedName) {
	var names []string
	for fNames := range foundNames {
		for _, n := range fNames {
			name := n.NameString
			if _, ok := uniqNames[name]; ok {
				continue
			}

			uniqNames[name] = empty
			names = append(names, name)
		}
		if len(names) > 50000 {
			chunk := make([]string, len(names))
			copy(chunk, names)
			resChan <- chunk
			names = names[:0]
		}
	}
	resChan <- names
	close(resChan)
	fmt.Println("Unique names", len(uniqNames))
}
