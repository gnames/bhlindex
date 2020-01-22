package loader

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/gnames/bhlindex"
	"github.com/gnames/bhlindex/models"
	"github.com/lib/pq"
)

var env = bhlindex.EnvVars()

// FindItems starts from the root bhl directory and traverses its children
// collecting directories that correspond to BHL items.
func FindItems(c chan<- string) {
	root := env.BHLDir
	currentDir := ""

	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if models.IsPageFile(filepath.Base(path)) {
				if dir := filepath.Dir(path); dir != currentDir {
					currentDir = dir
					c <- dir
				}
			}
			return nil
		})
	bhlindex.Check(err)
	close(c)
}

// ImportItems saves items into databsase and removes
// duplicates. This process preceeds actual work on name resolution. After
// a item is impored to the database its id goes into itemIDs channel
func ImportItems(db *sql.DB, itemIDs chan<- int) {
	c := make(chan string)
	var wg sync.WaitGroup

	wg.Add(1)
	go itemsWorker(&wg, c, itemIDs, db)
	FindItems(c)
	wg.Wait()
}

func itemsWorker(wg *sync.WaitGroup, c <-chan string, itemIDs chan<- int,
	db *sql.DB) {
	defer wg.Done()
	for path := range c {
		item := itemFromPath(path)
		id := item.Insert(db)
		if id > 0 {
			itemIDs <- id
		}
	}
	close(itemIDs)
}

func itemFromPath(p string) models.Item {
	dirs := strings.Split(p, "/")
	item := dirs[len(dirs)-1]
	return models.Item{Path: p, InternetArchiveID: item}
}

// ImportPages finds all pages files and saves them to the database
func ImportPages(db *sql.DB, t *models.Item) {
	d, err := os.Open(t.Path)
	bhlindex.Check(err)
	defer func() {
		err := d.Close()
		bhlindex.Check(err)
	}()

	names, err := d.Readdirnames(-1)
	bhlindex.Check(err)
	var pageIDs []string
	for _, name := range names {
		if models.IsPageFile(name) {
			id := models.PageID(name)
			pageIDs = append(pageIDs, id)
		}
	}
	pages := generatePages(pageIDs, t.ID)
	t.Content.Concatenate(pages, t.Path)
	savePages(db, t)
}

func generatePages(pageIDs []string, itemID int) []models.Page {
	pages := make([]models.Page, len(pageIDs))
	sort.Strings(pageIDs)
	for i, v := range pageIDs {
		pages[i] = models.Page{ID: v, ItemID: itemID}
	}
	return pages
}

func savePages(db *sql.DB, t *models.Item) {
	batch := t.Content.Pages
	columns := []string{"page_id", "item_id", "page_offset"}
	transaction, err := db.Begin()
	bhlindex.Check(err)

	stmt, err := transaction.Prepare(pq.CopyIn("pages", columns...))
	bhlindex.Check(err)

	for _, p := range batch {
		_, err = stmt.Exec(p.ID, p.ItemID, p.Offset)
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
