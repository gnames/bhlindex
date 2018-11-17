package finder

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gnames/bhlindex"
	"github.com/gnames/bhlindex/models"
	gfutil "github.com/gnames/gnfinder/util"
	"github.com/gnames/gnfinder/verifier"
	"github.com/lib/pq"
)

const BATCH_SIZE = 50000

func Verify(db *sql.DB, workers int) {
	namesVerified := make(chan bool)
	counter := make(chan int)

	models.Truncate(db, "name_strings")
	models.Truncate(db, "name_statuses")
	populateUniqueNames(db)
	go verifyLog(counter)
	go verifyNames(db, counter, namesVerified, workers)

	// wait till all names are imported
	<-namesVerified
	close(counter)
}

func verifyNames(db *sql.DB, counter chan<- int,
	namesVerified chan<- bool, workers int) {
	for {
		time.Sleep(200 * time.Microsecond)
		verifyNamesQuery(db, counter, workers)
		status := readStatus(db)
		if status.Status == AllNamesVerified {
			break
		}
	}
	namesVerified <- true
}

func verifyNamesQuery(db *sql.DB, counter chan<- int, workers int) int {
	m := gfutil.NewModel()
	m.Workers = workers
	q := `
    WITH temp AS (
      SELECT name FROM name_statuses
        WHERE processed=false LIMIT $1)
    UPDATE name_statuses ns SET processed=true
    FROM temp
      WHERE ns.name = temp.name
    RETURNING temp.name`

	rows, err := db.Query(q, BATCH_SIZE)
	bhlindex.Check(err)
	defer rows.Close()

	var name string
	names := make([]string, 0, BATCH_SIZE)

	for rows.Next() {
		rows.Scan(&name)
		names = append(names, name)
	}
	err = rows.Err()
	bhlindex.Check(err)

	namesSize := len(names)

	status := readStatus(db)
	if status.Status == AllNamesVerified && namesSize == 0 {
		updateStatus(db, &metaData{Status: AllErrorsProcessed})
	} else if status.Status == AllNamesHarvested && namesSize == 0 {
		updateStatus(db, &metaData{Status: AllNamesVerified})
	}

	counter <- namesSize
	time1 := time.Now().UnixNano()
	verified := verifier.Verify(names, m)
	if namesSize > 0 {
		timeSpent := float64(time.Now().UnixNano()-time1) / 1000000000
		speed := int(float64(namesSize) / timeSpent)
		log.Printf("\033[40;32;1mVerification speed %d names/sec\033[0m\n", speed)
		saveVerifiedNameStrings(db, verified)
	}
	return namesSize
}

func verifyLog(counter <-chan int) {
	total := 0
	count_same := 0
	for i := range counter {
		total += i
		if i > 0 {
			count_same = 0
			log.Printf("\033[40;32;1mVerifying %d name-strings\033[0m\n", total)
		} else {
			count_same++
			if count_same%100 == 0 {
				fmt.Printf("\033[40;31;1mVerification que is empty %d times\033[0m\n",
					count_same)
			}
		}
	}
}

func saveVerifiedNameStrings(db *sql.DB, verified verifier.VerifyOutput) {
	var errStr sql.NullString
	now := time.Now()
	columns := []string{"name", "match_type", "edit_distance",
		"matched_name", "current_name", "classification", "datasource_id",
		"datasources_number", "curation", "retries", "error", "updated_at"}
	transaction, err := db.Begin()
	bhlindex.Check(err)

	stmt, err := transaction.Prepare(pq.CopyIn("name_strings", columns...))
	bhlindex.Check(err)

	for name, v := range verified {
		if v.Error == "" {
			errStr = sql.NullString{}
		} else {
			errStr.Scan(v.Error)
		}
		_, err = stmt.Exec(name, v.MatchType, v.EditDistance, v.MatchedName,
			v.CurrentName, v.ClassificationPath, v.DataSourceID, v.DataSourcesNum,
			v.DataSourceQuality, v.Retries, errStr, now)
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

func populateUniqueNames(db *sql.DB) {

	log.Println("\033[40;32;1mExtract unique name-strings to verify.\033[0m")
	q := `INSERT INTO name_statuses
          SELECT DISTINCT name_string, false
            FROM page_name_strings
            ORDER BY name_string`

	stmt, err := db.Prepare(q)
	bhlindex.Check(err)

	_, err = stmt.Exec()
	bhlindex.Check(err)

	err = stmt.Close()
	bhlindex.Check(err)
	updateStatus(db, &metaData{Status: AllNamesHarvested})
	log.Println("\033[40;32;1mVerification started.\033[0m")
}
