package finder

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gnames/bhlindex"
	"github.com/gnames/bhlindex/models"
	"github.com/gnames/gnfinder/verifier"
	"github.com/lib/pq"
)

const BATCH_SIZE = 50000

func Verify(db *sql.DB, workersNum int) {
	namesVerified := make(chan bool)
	counter := make(chan int)

	models.Truncate(db, "name_strings")
	models.Truncate(db, "name_statuses")
	populateUniqueNames(db)
	go verifyLog(counter)
	go verifyNames(db, counter, namesVerified, workersNum)

	// wait till all names are imported
	<-namesVerified
	close(counter)
}

func verifyNames(db *sql.DB, counter chan<- int,
	namesVerified chan<- bool, workersNum int) {
	for {
		time.Sleep(200 * time.Microsecond)
		verifyNamesQuery(db, counter, workersNum)
		status := readStatus(db)
		if status.Status == AllNamesVerified {
			break
		}
	}
	namesVerified <- true
}

func verifyNamesQuery(db *sql.DB, counter chan<- int, workersNum int) int {
	env := bhlindex.EnvVars()
	verif := verifier.NewVerifier(verifier.OptWorkers(workersNum),
		verifier.OptSources(env.PrefSources))
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
	verified := verif.Run(names)
	if namesSize > 0 {
		timeSpent := float64(time.Now().UnixNano()-time1) / 1000000000
		speed := int(float64(namesSize) / timeSpent)
		log.Printf("\033[40;32;1mVerification speed %d names/sec\033[0m\n", speed)
		saveVerifiedNameStrings(db, verified)
		savePreferredSources(db, verified)
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

func saveVerifiedNameStrings(db *sql.DB, verified verifier.Output) {
	var errStr sql.NullString
	now := time.Now()
	columns := []string{
		"name", "taxon_id", "match_type", "edit_distance",
		"stem_edit_distance", "matched_name", "matched_canonical",
		"current_name", "classification", "datasource_id", "datasource_title",
		"datasources_number", "curation", "retries",
		"error", "updated_at"}
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

		br := v.BestResult
		_, err = stmt.Exec(name, br.TaxonID, br.MatchType, br.EditDistance,
			br.StemEditDistance, br.MatchedName, br.MatchedCanonical,
			br.CurrentName, br.ClassificationPath, br.DataSourceID, br.DataSourceTitle,
			v.DataSourcesNum, v.DataSourceQuality, v.Retries,
			errStr, now)
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

func savePreferredSources(db *sql.DB, verified verifier.Output) {
	columns := []string{"name", "taxon_id", "match_type", "edit_distance",
		"stem_edit_distance", "matched_name", "matched_canonical", "current_name",
		"classification", "datasource_id", "datasource_title"}

	transaction, err := db.Begin()
	bhlindex.Check(err)
	stmt, err := transaction.Prepare(pq.CopyIn("preferred_sources",
		columns...))
	bhlindex.Check(err)

	for name, v := range verified {

		if v.BestResult.MatchType == "NoMatch" {
			continue
		}
		for _, vv := range v.PreferredResults {
			_, err = stmt.Exec(name, vv.TaxonID, vv.MatchType, vv.EditDistance,
				vv.StemEditDistance, vv.MatchedName, vv.MatchedCanonical, vv.CurrentName,
				vv.ClassificationPath, vv.DataSourceID, vv.DataSourceTitle)
			bhlindex.Check(err)
		}
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

func populateUniqueNames(db *sql.DB) {

	log.Println("\033[40;32;1mExtract unique name-strings to verify.\033[0m")
	q := `INSERT INTO name_statuses
          SELECT name_string, AVG(odds), count(*), false
            FROM page_name_strings GROUP BY name_string
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
