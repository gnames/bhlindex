package loaderio

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/gnames/bhlindex/ent/item"
	"github.com/gnames/bhlindex/ent/page"
	"github.com/lib/pq"
)

func pageFromPath(path string) (*page.Page, error) {
	id := pageID(filepath.Base(path))
	return &page.Page{ID: id}, nil
}

func updatePages(itm *item.Item) error {
	var itemText []byte
	var offset int
	for i := range itm.Pages {
		itm.Pages[i].ItemID = itm.ID
		itm.Pages[i].Offset = offset
		path := filepath.Join(itm.Path, itm.Pages[i].ID+".txt")
		text, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		itemText = append(itemText, text...)
		pageUTF := []rune(string(text))
		offset += len(pageUTF)
		itm.Pages[i].OffsetNext = offset
	}
	itm.Text = itemText
	return nil
}

func (l loaderio) insertPages(itm *item.Item) error {
	var stmt *sql.Stmt
	batch := itm.Pages
	columns := []string{"id", "item_id", "offset"}

	transaction, err := l.db.Begin()
	if err != nil {
		return err
	}

	stmt, err = transaction.Prepare(pq.CopyIn("pages", columns...))
	if err != nil {
		return err
	}

	for _, p := range batch {
		_, err = stmt.Exec(p.ID, p.ItemID, p.Offset)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return transaction.Commit()
}

func isPageFile(f string) bool {
	res, _ := filepath.Match("*_[0-9][0-9][0-9][0-9].txt", f)
	return res
}

func pageID(f string) string {
	extLen := len(filepath.Ext(f))
	idLen := len(f) - extLen
	return f[0:idLen]
}
