package models

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/GlobalNamesArchitecture/bhlindex"
)

// Title respresents BHL title data. Title in BHL can be a book, a journal etc.
// Title is the name of a directory that has pages files. All pages files end
// follow "*_0001.txt" pattern.
type Title struct {
	ID                int
	Path              string
	InternetArchiveID string
	GnrdURL           string
	Status            int
	Language          string
	EnglishDetected   bool
	UpdatedAt         time.Time
	Content           Content
}

type Content struct {
	Pages []Page
	Text  []byte
}

func (c *Content) Concatenate(ps []Page, path string) {
	c.Pages = ps
	var text []byte
	offset := 0
	for i, p := range c.Pages {
		c.Pages[i].Offset = offset
		f := fmt.Sprintf("%s/%s.txt", path, p.ID)
		pageText, err := ioutil.ReadFile(f)
		bhlindex.Check(err)
		text = append(text, pageText...)
		pageUTF := []rune(string(pageText))
		offset += len(pageUTF)
	}
	c.Text = text
}

// Insert add data from a title to bhlindex database and returns newly
// a newly generated ID. If a title is duplicate instead of ID it returns 0.
func (t *Title) Insert(db *sql.DB) int {
	var id int
	q := `
INSERT INTO titles 
  (path, internet_archive_id, gnrd_url, status, language, english_detected,
	 updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING RETURNING id`
	err := db.QueryRow(q, t.Path, t.InternetArchiveID, t.GnrdURL, t.Status,
		t.Language, t.EnglishDetected, time.Now()).Scan(&id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			id = 0
		} else {
			bhlindex.Check(err)
		}
	}
	return id
}

func TitleFind(db *sql.DB, id int) Title {
	var status int
	var path, internetArchiveID string
	var gnrdURL, language sql.NullString
	var englishDetected bool
	var updatedAt time.Time

	err := db.QueryRow("SELECT * FROM titles WHERE id = $1", id).Scan(&id, &path,
		&internetArchiveID, &gnrdURL, &status, &language, &englishDetected,
		&updatedAt)
	bhlindex.Check(err)
	title := Title{id, path, internetArchiveID, gnrdURL.String,
		status, language.String, englishDetected, updatedAt, Content{}}
	return title
}

// Disable these temporarily, we need something like this later
// func (t *Title) CreateOrSelect(db *sql.DB) {
// 	var id, status int
// 	var path, internetArchiveID, gnrdURL, language sql.NullString
// 	var englishDetected bool
// 	var updatedAt time.Time
// 	q := `
// WITH new_row AS (
// 	INSERT INTO titles (path, internet_archive_id, gnrd_url, status, language,
// 		english_detected, updated_at)
// 		SELECT $1, CAST($2 AS VARCHAR), $3, $4, $5, $6, $7
// 			WHERE NOT EXISTS (SELECT * FROM titles WHERE internet_archive_id = $2)
// 				RETURNING *
// 	)
// 	SELECT * FROM new_row
// 	 	UNION
// 	SELECT * FROM titles WHERE internet_archive_id = $2
// `
// 	err := db.QueryRow(q, t.Path, t.InternetArchiveID, t.GnrdURL, t.Status,
// 		t.Language, t.EnglishDetected, time.Now()).Scan(&id, &path,
// 		&internetArchiveID, &gnrdURL, &status, &language,
// 		&englishDetected, &updatedAt)
// 	bhlindex.Check(err)
// 	t.ID = id
// 	t.Path = path.String
// 	t.InternetArchiveID = internetArchiveID.String
// 	t.GnrdURL = gnrdURL.String
// 	t.Status = status
// 	t.Language = language.String
// 	t.EnglishDetected = englishDetected
// 	t.UpdatedAt = updatedAt
// }
//
// func (t *Title) Delete(db *sql.DB) {
// 	var err error
// 	if t.ID == 0 {
// 		_, err = db.Query(`DELETE FROM titles WHERE ID = $1`, t.ID)
// 	} else {
// 		_, err = db.Query(`DELETE FROM titles WHERE internet_archive_id = $1`,
// 			t.InternetArchiveID)
// 	}
// 	bhlindex.Check(err)
// }
