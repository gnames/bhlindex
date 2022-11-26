package loaderio

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gnames/bhlindex/internal/ent/item"
	"github.com/lib/pq"
)

// insertItem  add data from a item to bhlindex database and returns newly
// UpdatedAt fields in the Items object.
// If an item is a duplicate its ID is 0.
func (l loaderio) insertItem(item *item.Item) error {
	var id int
	updatedAt := time.Now()
	q := `
INSERT INTO items
  (id, path, updated_at)
	VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id`
	err := l.db.QueryRow(q, item.ID, item.Path, updatedAt).Scan(&id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			id = 0
		} else {
			return err
		}
	}
	item.UpdatedAt = updatedAt
	return nil
}

func (l loaderio) insertPages(itm *item.Item) error {
	var stmt *sql.Stmt
	pgs := itm.Pages
	columns := []string{"id", "item_id", "file_num", "file_name", "offset"}

	transaction, err := l.db.Begin()
	if err != nil {
		return err
	}

	stmt, err = transaction.Prepare(pq.CopyIn("pages", columns...))
	if err != nil {
		return err
	}

	for _, p := range pgs {
		_, err = stmt.Exec(p.ID, p.ItemID, p.FileNum, p.FileName, p.Offset)
		if err != nil {
			err = fmt.Errorf("-> Exec PageID %d %w", p.ID, err)
			return err
		}
	}

	// Flush COPY FROM to db.
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
