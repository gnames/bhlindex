package models

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gnames/bhlindex"
	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/output"
)

// Item respresents BHL item data. Item in BHL can be a book, a journal etc.
// Item's internet_archive_id is the name of a directory that has pages files.
// All pages files names follow "{item_name}_0001.txt" pattern.
type Item struct {
	ID                int
	Path              string
	InternetArchiveID string
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
		c.Pages[i].OffsetNext = offset
	}
	c.Text = text
}

// Insert add data from a item to bhlindex database and returns newly
// a newly generated ID. If a item is duplicate instead of ID it returns 0.
func (t *Item) Insert(db *sql.DB) int {
	var id int
	q := `
INSERT INTO items
  (path, internet_archive_id, status, language, english_detected,
	 updated_at)
	VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING RETURNING id`
	err := db.QueryRow(q, t.Path, t.InternetArchiveID, t.Status,
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

func (t *Item) FindNames(d *dict.Dictionary) []DetectedName {
	text := t.Content.Text
	gnf := gnfinder.NewGNfinder(gnfinder.OptDict(d), gnfinder.OptBayes(true))
	output := gnf.FindNames(text)
	detectedNames := namesToDetectedNames(t, output.Names)
	return detectedNames
}

func ItemFind(db *sql.DB, id int) Item {
	var status int
	var path, internetArchiveID string
	var language sql.NullString
	var englishDetected bool
	var updatedAt time.Time

	err := db.QueryRow("SELECT * FROM items WHERE id = $1", id).Scan(&id, &path,
		&internetArchiveID, &status, &language, &englishDetected,
		&updatedAt)
	bhlindex.Check(err)
	item := Item{id, path, internetArchiveID,
		status, language.String, englishDetected, updatedAt, Content{}}
	return item
}

func namesToDetectedNames(t *Item, names []output.Name) []DetectedName {
	ns := make([]DetectedName, len(names))
	j := 0
	if j >= len(names) {
		return ns
	}
	name := names[j]
	for _, page := range t.Content.Pages {
		for {
			if name.OffsetStart <= page.OffsetNext {
				ns[j] = NewDetectedName(t.ID, page, name)
				j++
				if j >= len(names) {
					return ns
				}
				name = names[j]
			} else {
				break
			}
		}
	}
	return ns
}
